# PostgreSQL 入门与项目实践

## 适合谁看

适合想学习现代关系型数据库能力的人。PostgreSQL 在类型系统、约束、索引、JSONB、全文搜索、事务和扩展能力上很强，适合复杂业务系统、SaaS、数据平台和后端服务。

如果你已经会 MySQL，学习 PostgreSQL 时重点关注：

- 更丰富的数据类型。
- 更强的约束表达。
- JSONB 和 GIN 索引。
- EXPLAIN ANALYZE。
- 事务隔离和 MVCC。
- RLS、扩展和函数能力。

## 基础表设计示例

```sql
create table public.organization (
  -- 使用 identity 是 SQL 标准方式，适合单库业务自增主键。
  id bigint generated always as identity primary key,

  -- 组织名称用于后台展示和搜索；同一租户体系内要求唯一时，应加唯一约束。
  name text not null,

  -- 状态限制在固定集合内，避免应用层写入未知状态。
  status text not null default 'active'
    check (status in ('active', 'disabled')),

  -- 扩展配置用 jsonb，适合少量非核心、变化较快的组织设置。
  settings jsonb not null default '{}'::jsonb,

  -- 创建时间使用 timestamptz，保留时区语义，便于跨地区系统统一处理。
  created_at timestamptz not null default now(),

  -- 更新时间由应用或触发器维护，用于审计、同步和缓存失效。
  updated_at timestamptz not null default now()
);

comment on table public.organization is '组织表：保存 SaaS 租户或企业组织的基础信息、启停状态和少量扩展配置';
comment on column public.organization.id is '组织主键，数据库内部唯一标识，使用 identity 便于单库自增和索引局部性';
comment on column public.organization.name is '组织名称，用于后台展示、搜索和审计日志，不允许为空';
comment on column public.organization.status is '组织状态：active=正常使用，disabled=停用；通过 check 约束阻止未知状态';
comment on column public.organization.settings is '组织扩展配置，保存低频变化的非核心设置；核心查询字段不应长期藏在 JSONB 中';
```

## 约束

PostgreSQL 的约束能力很重要：

| 约束 | 作用 |
| --- | --- |
| primary key | 唯一定位记录 |
| not null | 保证关键字段存在 |
| unique | 保证业务唯一 |
| foreign key | 保证关联关系 |
| check | 保证字段值域 |

示例：

```sql
alter table public.organization
  add constraint organization_name_unique unique (name);
```

迁移里要注意：PostgreSQL 不支持 `ADD CONSTRAINT IF NOT EXISTS`。如果要写幂等迁移，需要先查 `pg_constraint`，再决定是否添加。

## 主键策略

常见选择：

| 场景 | 推荐 |
| --- | --- |
| 单库业务表 | `bigint generated always as identity` |
| 外部暴露 ID | 增加 public_id 或使用有序 UUID/ULID |
| 分布式写入 | 有序 UUID、ULID、雪花 ID |

随机 UUID v4 在大表主键上可能造成索引碎片。除非业务确实需要，不要把随机 UUID 当作默认主键。

## 外键和索引

PostgreSQL 不会自动给外键列创建索引。外键列如果用于 JOIN、查询或级联删除，应该显式创建索引。

```sql
create table public.member (
  id bigint generated always as identity primary key,

  -- 关联组织表，表示成员属于哪个组织。
  organization_id bigint not null references public.organization(id),

  user_id bigint not null,
  role text not null default 'member',
  created_at timestamptz not null default now()
);

comment on column public.member.organization_id is '成员所属组织 ID；用于组织成员列表查询和组织删除前的关联检查';

create index member_organization_id_idx
on public.member (organization_id);
```

## JSONB 使用边界

JSONB 适合：

- 扩展配置。
- 不稳定字段。
- 第三方响应快照。
- 低频查询的附加信息。

不适合：

- 核心筛选条件。
- 强约束字段。
- 频繁 JOIN 字段。
- 需要稳定统计和报表的字段。

如果某个 JSONB 字段开始频繁出现在 WHERE 条件中，要考虑拆成正式列或建立合适索引。

## EXPLAIN ANALYZE

PostgreSQL 查询慢时，不要靠猜。

```sql
explain (analyze, buffers)
select *
from public.member
where organization_id = 10
order by created_at desc
limit 20;
```

重点看：

- 是否 Seq Scan 扫大表。
- 是否使用预期索引。
- 实际行数和预估行数差距是否过大。
- Buffers read 是否很高。
- Sort 是否使用磁盘。

`EXPLAIN ANALYZE` 会实际执行查询。对更新、删除等写操作使用前要谨慎，可以放在事务中并回滚，或先用只读查询分析。

## 实际项目问题

### 问题：组织成员列表越来越慢

**原因**

`member.organization_id` 是外键，但没有索引。

**解决方案**

```sql
create index concurrently member_organization_id_created_at_idx
on public.member (organization_id, created_at desc);
```

这个索引服务两个场景：按组织筛选成员、按创建时间倒序分页。

### 问题：状态字段写入了奇怪值

**原因**

只在后端代码里限制枚举，数据库没有兜底。

**解决方案**

```sql
alter table public.organization
  add constraint organization_status_check
  check (status in ('active', 'disabled'));
```

迁移前要先清理历史脏数据，否则添加约束会失败。

## 最佳实践

- 新表写清楚 table 和 column comment。
- 核心业务规则尽量用数据库约束兜底。
- 外键列按查询场景显式加索引。
- JSONB 不要替代正常建模。
- 慢查询先用 EXPLAIN ANALYZE 找证据。
- 大表加索引优先考虑并发创建和低峰执行。

## 参考资料

- [PostgreSQL: Constraints](https://www.postgresql.org/docs/current/ddl-constraints.html)
- [PostgreSQL: CREATE INDEX](https://www.postgresql.org/docs/current/sql-createindex.html)
- [PostgreSQL: Using EXPLAIN](https://www.postgresql.org/docs/current/using-explain.html)
- [PostgreSQL: Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)

## 下一步学习

继续学习 [Redis 缓存与数据结构](/database/redis)。
