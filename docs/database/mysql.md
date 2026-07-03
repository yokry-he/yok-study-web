# MySQL 入门与项目实践

## 适合谁看

适合要做后台管理、业务系统、电商、CRM、权限系统的人。MySQL 在国内业务系统里非常常见，很多 Java、Go、Node 后端都会使用它。

学习 MySQL 的重点不是背语法，而是理解 InnoDB、主键、索引、事务、字符集和慢查询。

## MySQL 和 InnoDB

现代 MySQL 项目默认使用 InnoDB。InnoDB 提供：

- 事务。
- 行级锁。
- 崩溃恢复。
- 外键约束。
- 聚簇索引。

如果你没有特别理由，新表应使用 InnoDB。

## 基础表设计示例

下面是后台系统中的用户表示例。示例刻意写了详细注释，实际项目中数据库迁移也应该保留类似说明。

```sql
create table sys_user (
  -- 主键使用 bigint，便于长期增长；业务上不暴露自增规律时，可额外提供 public_id。
  id bigint unsigned not null auto_increment comment '用户主键，系统内部唯一标识，不直接作为公开身份凭证',

  -- 用户名用于登录，必须唯一；长度限制来自产品规则和登录表单限制。
  username varchar(64) not null comment '登录用户名，业务要求全局唯一，创建后谨慎修改',

  -- 手机号可能用于登录或找回密码；如果业务允许未绑定手机号，则可为空。
  mobile varchar(20) null comment '用户手机号，可用于登录、通知或找回密码，允许未绑定',

  -- 状态使用小整数，避免直接散落字符串；具体含义应在代码枚举和文档中同步维护。
  status tinyint unsigned not null default 1 comment '用户状态：1=启用，2=禁用，3=冻结',

  -- 创建时间由数据库生成，保证不同服务写入时规则一致。
  created_at datetime not null default current_timestamp comment '记录创建时间，用于审计和排序',

  -- 更新时间自动维护，便于排查数据最近一次变化。
  updated_at datetime not null default current_timestamp on update current_timestamp comment '记录更新时间，用于审计和缓存失效判断',

  primary key (id),

  -- 用户名登录查询高频，唯一约束同时承担业务约束和查询加速。
  unique key uk_sys_user_username (username),

  -- 手机号查询可能用于登录或客服检索；如果手机号允许重复或为空，要按业务重新设计。
  key idx_sys_user_mobile (mobile),

  -- 后台列表常按状态筛选并按创建时间排序，复合索引服务该查询场景。
  key idx_sys_user_status_created (status, created_at)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci comment='系统用户表，保存后台用户登录身份、联系信息和启停状态';
```

## 主键选择

常见主键：

| 类型 | 优点 | 风险 |
| --- | --- | --- |
| 自增 bigint | 简单、性能稳定、索引局部性好 | 暴露增长趋势，不适合公开 ID |
| UUID | 分布式生成方便 | 随机 UUID 可能导致索引碎片和空间变大 |
| 雪花 ID | 分布式、趋势递增 | 依赖生成器、时钟和实现规范 |

中小型业务系统优先自增 bigint。需要公开给外部时，可以增加一个业务编号或 public_id，而不是直接暴露内部主键。

## 字符集

推荐使用 `utf8mb4`，它能支持 emoji 和更完整的 Unicode 字符。

如果字符集混乱，可能出现：

- 用户昵称保存失败。
- 排序规则不一致。
- 联表比较时索引失效。
- 迁移时乱码。

## 索引基础

索引用来加速查询，但会增加写入成本。

适合加索引的字段：

- 高频 WHERE 条件。
- JOIN 字段。
- ORDER BY 字段。
- 唯一业务字段。

不适合盲目加索引：

- 低选择性的布尔字段。
- 很少查询的字段。
- 频繁更新但查询少的字段。

## EXPLAIN

MySQL 可以用 EXPLAIN 查看执行计划：

```sql
explain
select id, username, mobile
from sys_user
where status = 1
order by created_at desc
limit 20;
```

重点看：

| 字段 | 关注点 |
| --- | --- |
| `type` | 是否出现 `ALL` 全表扫描 |
| `key` | 是否使用预期索引 |
| `rows` | 预计扫描行数是否过大 |
| `Extra` | 是否出现 filesort、temporary |

不要看到查询慢就立刻加索引。先看执行计划，再判断索引是否匹配查询。

## 事务示例

转账、下单、库存扣减这类操作必须使用事务。

```sql
start transaction;

-- 锁定订单，确保同一订单不会被两个支付回调同时更新。
select id, status
from orders
where id = 1001
for update;

-- 只允许待支付订单变为已支付，避免重复回调造成重复入账。
update orders
set status = 'paid',
    paid_at = now()
where id = 1001
  and status = 'pending';

commit;
```

事务内不要做 HTTP 调用、文件上传、远程 RPC 等慢操作。先完成外部调用，再开启短事务更新数据库。

## 实际项目问题

### 问题：后台列表数据多了以后变慢

**原因**

列表常见查询没有合适索引：

```sql
select *
from sys_user
where status = 1
order by created_at desc
limit 20;
```

**解决方案**

增加匹配筛选和排序的复合索引：

```sql
create index idx_sys_user_status_created
on sys_user (status, created_at);
```

并记录索引用途：服务后台用户列表按状态筛选、按创建时间倒序分页查询。

### 问题：删除用户时报外键错误

**原因**

用户已经被订单、角色或审计日志引用。

**解决方案**

不要直接物理删除核心业务实体。常见做法是软删除：

```sql
alter table sys_user
  add column deleted_at datetime null comment '软删除时间，为空表示未删除；保留用户历史关联和审计记录';
```

## 最佳实践

- 新表默认 InnoDB 和 utf8mb4。
- 主键优先使用 bigint 自增或有序分布式 ID。
- 每个唯一业务规则都应有唯一约束。
- 列表查询上线前用 EXPLAIN 看执行计划。
- 核心表不要轻易物理删除。
- 每次表结构变更写清楚业务原因和回滚风险。

## 参考资料

- [MySQL: Introduction to InnoDB](https://dev.mysql.com/doc/refman/9.1/en/innodb-introduction.html)
- [MySQL: Clustered and Secondary Indexes](https://dev.mysql.com/doc/refman/9.7/en/innodb-index-types.html)
- [MySQL: InnoDB Locking](https://dev.mysql.com/doc/en/innodb-locking.html)

## 下一步学习

继续学习 [PostgreSQL 入门与项目实践](/database/postgresql)。
