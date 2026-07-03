# 迁移、种子与版本治理

## 适合谁看

适合开始参与真实项目表结构变更的人。数据库迁移是高风险操作，因为它会影响已有数据和线上服务。

一个成熟的迁移不只是 SQL 文件，还应该包含：

- 变更原因。
- 影响范围。
- 字段和约束含义。
- 回滚策略。
- 数据修复方案。
- 上线顺序。
- 验证方式。

## migration 是什么

migration 是数据库结构变更脚本，例如：

- 新增表。
- 新增字段。
- 修改字段类型。
- 添加索引。
- 添加约束。
- 清理历史数据。
- 初始化基础数据。

迁移应该进入版本管理，和代码一起评审。

## 一个好的迁移文件

```sql
-- 迁移目标：
-- 为订单表增加支付流水号，用于对接第三方支付平台的账务核对。
--
-- 业务背景：
-- 原系统只保存订单状态，无法在支付平台账单中快速定位对应流水。
--
-- 数据约束：
-- payment_no 由支付平台返回，同一支付平台内唯一；历史未支付订单允许为空。
--
-- 上线注意：
-- 先新增 nullable 字段并发布后端写入逻辑，再回填历史数据；确认无空值后再考虑增加唯一约束。

alter table orders
  add column payment_no text null;

comment on column orders.payment_no is '第三方支付流水号，用于账务核对；历史未支付订单允许为空，支付成功后必须写入';

create unique index concurrently orders_payment_no_unique_idx
on orders (payment_no)
where payment_no is not null;
```

这个例子说明了：

- 为什么改。
- 字段含义是什么。
- 是否允许为空。
- 为什么唯一索引是部分索引。
- 上线顺序是什么。

## 不要一次做太多

高风险做法：

```text
新增字段
↓
立即改成 not null
↓
立即改后端
↓
立即删除旧字段
```

更稳的方式：

1. 新增兼容字段。
2. 后端同时写新旧字段。
3. 回填历史数据。
4. 切读新字段。
5. 验证稳定。
6. 删除旧字段。

这叫渐进式迁移。

## 添加约束

PostgreSQL 不支持：

```sql
alter table users
add constraint if not exists users_email_unique unique (email);
```

可以先检查：

```sql
do $$
begin
  if not exists (
    select 1
    from pg_constraint
    where conname = 'users_email_unique'
      and conrelid = 'users'::regclass
  ) then
    alter table users
      add constraint users_email_unique unique (email);
  end if;
end $$;
```

添加约束前要先检查历史数据是否满足约束。

## 添加索引

大表添加索引要谨慎。

PostgreSQL 可考虑：

```sql
create index concurrently orders_created_at_idx
on orders (created_at desc);
```

注意：`create index concurrently` 不能放在普通事务块里执行。

MySQL 大表加索引也要评估锁表、执行时间和业务低峰窗口。

## seed 是什么

seed 是种子数据，用来初始化系统必要数据。

适合 seed：

- 默认管理员角色。
- 基础权限码。
- 字典项。
- 系统配置。

不适合 seed：

- 真实用户数据。
- 生产业务数据。
- 随机测试垃圾数据。

示例：

```sql
insert into permission (code, name, type)
values
  ('system:user:view', '查看用户', 'button'),
  ('system:user:create', '新增用户', 'button')
on conflict (code) do update
set name = excluded.name,
    type = excluded.type;
```

业务含义：权限 code 是稳定标识，seed 可以更新展示名称，但不应该随意改 code。

## 回滚策略

不是所有迁移都能安全回滚。

| 变更 | 回滚难度 |
| --- | --- |
| 新增 nullable 字段 | 低 |
| 新增索引 | 低 |
| 新增表 | 中 |
| 删除字段 | 高 |
| 修改字段类型 | 高 |
| 清理历史数据 | 高 |

删除字段前要确认：

- 代码不再读写。
- 报表不再依赖。
- 数据备份可恢复。
- 回滚期间旧代码不会访问。

## 实际项目问题

### 问题：上线后字段 not null 迁移失败

**原因**

历史数据里有 null。

**解决方案**

拆成多步：

1. 新增 nullable 字段。
2. 发布应用写入新字段。
3. 回填历史数据。
4. 检查是否还有 null。
5. 再改 not null。

### 问题：迁移脚本没人敢改

**原因**

没有业务注释，不知道字段和约束为什么存在。

**解决方案**

每次迁移必须写清楚：

- 业务背景。
- 字段含义。
- 约束原因。
- 回滚风险。
- 验证 SQL。

## 最佳实践

- 数据库迁移必须进入代码评审。
- 结构变更要有详细注释和文档。
- 大表变更先评估锁和执行时间。
- 删除字段和改类型要分阶段。
- seed 使用稳定业务 key，保证可重复执行。
- 每次上线前准备验证 SQL 和回滚方案。

## 下一步学习

继续学习 [ORM 实战](/database/orm-practice)，把迁移、模型、查询和事务放到应用代码中管理。
