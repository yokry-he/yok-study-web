# 索引与查询优化

## 适合谁看

适合遇到接口慢、列表慢、搜索慢，或者不知道什么时候该加索引的人。

索引的核心作用是减少扫描范围。没有索引时，数据库可能要从表头扫到表尾；有合适索引时，可以更快定位目标数据。

## 索引不是越多越好

索引的收益：

- 加速 WHERE 查询。
- 加速 JOIN。
- 加速 ORDER BY。
- 支持唯一约束。

索引的成本：

- 占用磁盘。
- 写入、更新、删除时要维护索引。
- 索引设计不匹配查询时可能用不上。

所以索引必须对应查询场景。

## 哪些列应该加索引

优先考虑：

| 场景 | 示例 |
| --- | --- |
| 高频筛选 | `where user_id = ?` |
| JOIN 条件 | `join orders on orders.user_id = users.id` |
| 排序分页 | `order by created_at desc limit 20` |
| 唯一业务规则 | `unique email` |
| 外键列 | `order.user_id` |

不要盲目索引：

- 值很少的字段，例如单独的 `enabled`。
- 很少用于查询的字段。
- 大文本字段。
- 频繁更新但查询少的字段。

## 单列索引

```sql
create index orders_customer_id_idx
on orders (customer_id);
```

适合：

```sql
select *
from orders
where customer_id = 1001;
```

## 复合索引

如果查询同时按多个字段过滤，复合索引通常比多个单列索引更合适。

```sql
create index orders_status_created_at_idx
on orders (status, created_at desc);
```

适合：

```sql
select *
from orders
where status = 'pending'
order by created_at desc
limit 20;
```

列顺序很重要。一般把等值条件放前面，范围或排序字段放后面。

## 外键索引

很多数据库不会自动给外键列创建索引。外键列经常用于 JOIN 和级联操作，应该显式评估。

```sql
create index orders_customer_id_idx
on orders (customer_id);
```

如果没有这个索引，查询某个用户订单和删除用户前检查关联都可能变慢。

## 覆盖索引和回表

如果查询需要的字段都在索引中，数据库可能不需要再读取完整行。

示例：

```sql
create index orders_customer_status_created_idx
on orders (customer_id, status, created_at desc);
```

适合：

```sql
select customer_id, status, created_at
from orders
where customer_id = 1001
  and status = 'paid'
order by created_at desc
limit 20;
```

但不要为了覆盖所有查询创建很宽的索引。宽索引写入成本更高。

## N+1 查询

N+1 是后端项目常见性能问题。

错误模式：

```ts
const users = await db.query('select id, name from users')

for (const user of users) {
  user.orders = await db.query('select * from orders where user_id = $1', [user.id])
}
```

如果有 100 个用户，就会产生 101 次数据库请求。

改成批量查询：

```sql
select *
from orders
where user_id = any($1::bigint[]);
```

或者 JOIN：

```sql
select u.id, u.name, o.id as order_id, o.total
from users u
left join orders o on o.user_id = u.id
where u.status = 'active';
```

## EXPLAIN 怎么看

PostgreSQL：

```sql
explain (analyze, buffers)
select *
from orders
where customer_id = 1001
  and status = 'paid';
```

MySQL：

```sql
explain
select *
from orders
where customer_id = 1001
  and status = 'paid';
```

重点看：

- 是否全表扫描。
- 是否使用了预期索引。
- 扫描行数是否过大。
- 排序是否使用临时表或磁盘。
- 实际耗时是否集中在某个节点。

## 实际项目问题

### 问题：后台列表第一页很快，越往后越慢

**原因**

偏移分页：

```sql
select *
from orders
order by created_at desc
limit 20 offset 100000;
```

数据库仍然要跳过大量记录。

**解决方案**

对于无限滚动或时间线，用游标分页：

```sql
select *
from orders
where created_at < $1
order by created_at desc
limit 20;
```

配套索引：

```sql
create index orders_created_at_idx
on orders (created_at desc);
```

### 问题：加了索引但查询还是慢

**可能原因**

- 查询条件顺序和复合索引不匹配。
- 数据选择性太低。
- 函数包裹了索引列。
- 统计信息过旧。
- 返回数据太多。

**排查**

用 EXPLAIN 看真实执行计划，不要只看索引是否存在。

## 最佳实践

- 索引围绕查询场景设计。
- 高频 WHERE、JOIN、ORDER BY 优先评估索引。
- 外键列不要忘记索引。
- 复合索引注意列顺序。
- 慢查询先 EXPLAIN，再改 SQL 或加索引。
- 删除不再使用的重复索引，避免写入成本膨胀。

## 参考资料

- [PostgreSQL: Indexes](https://www.postgresql.org/docs/current/indexes.html)
- [PostgreSQL: Multicolumn Indexes](https://www.postgresql.org/docs/current/indexes-multicolumn.html)
- [PostgreSQL: Using EXPLAIN](https://www.postgresql.org/docs/current/using-explain.html)

## 下一步学习

继续学习 [事务、锁与并发](/database/transactions)。
