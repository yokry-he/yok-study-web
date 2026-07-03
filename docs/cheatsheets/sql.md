# SQL 速查

## 查询基础

```sql
SELECT id, username, created_at
FROM users
WHERE enabled = true
ORDER BY created_at DESC
LIMIT 20;
```

| 子句 | 用途 |
| --- | --- |
| `SELECT` | 选择字段 |
| `FROM` | 选择表 |
| `WHERE` | 过滤条件 |
| `ORDER BY` | 排序 |
| `LIMIT` | 限制条数 |
| `OFFSET` | 跳过条数 |

项目里不要默认 `SELECT *`，列表接口只查需要展示的字段。

## 条件查询

```sql
SELECT id, order_no, status
FROM orders
WHERE status = 'paid'
  AND created_at >= '2026-07-01'
  AND user_id IN (1, 2, 3);
```

模糊搜索：

```sql
SELECT id, username
FROM users
WHERE username LIKE 'tom%';
```

普通索引更适合前缀匹配。`%tom%` 通常难以有效使用普通索引。

## JOIN

```sql
SELECT u.id, u.username, r.name AS role_name
FROM users u
LEFT JOIN user_roles ur ON ur.user_id = u.id
LEFT JOIN roles r ON r.id = ur.role_id
WHERE u.enabled = true;
```

| JOIN | 含义 |
| --- | --- |
| `INNER JOIN` | 两边都有匹配才返回 |
| `LEFT JOIN` | 保留左表数据 |
| `RIGHT JOIN` | 保留右表数据 |

后台列表常用 `LEFT JOIN`，避免关联数据缺失时主数据也消失。

## 聚合统计

```sql
SELECT status, COUNT(*) AS count
FROM orders
GROUP BY status;
```

带条件聚合：

```sql
SELECT user_id, SUM(amount) AS total_amount
FROM orders
WHERE status = 'paid'
GROUP BY user_id
HAVING SUM(amount) > 1000;
```

`WHERE` 在分组前过滤，`HAVING` 在分组后过滤。

## 分页

普通分页：

```sql
SELECT id, order_no, created_at
FROM orders
ORDER BY created_at DESC
LIMIT 20 OFFSET 40;
```

深分页建议改成游标分页：

```sql
SELECT id, order_no, created_at
FROM orders
WHERE created_at < '2026-07-01 12:00:00'
ORDER BY created_at DESC
LIMIT 20;
```

数据量大时，`OFFSET` 越大越慢。

## 索引

单列索引：

```sql
CREATE INDEX idx_users_mobile ON users(mobile);
```

组合索引：

```sql
CREATE INDEX idx_orders_status_created
ON orders(status, created_at DESC);
```

唯一索引：

```sql
CREATE UNIQUE INDEX uk_users_mobile ON users(mobile);
```

常见原则：

- 高频过滤字段适合建索引。
- 高频排序字段要结合查询条件设计。
- 唯一业务约束用唯一索引保证。
- 索引不是越多越好，写入也要维护索引。

## 事务

```sql
BEGIN;

UPDATE accounts
SET balance = balance - 100
WHERE id = 1;

UPDATE accounts
SET balance = balance + 100
WHERE id = 2;

COMMIT;
```

失败时：

```sql
ROLLBACK;
```

事务适合保证多步数据库写入的一致性。不要在事务里执行慢网络请求。

## EXPLAIN

```sql
EXPLAIN
SELECT id, order_no
FROM orders
WHERE status = 'paid'
ORDER BY created_at DESC
LIMIT 20;
```

重点看：

| 信息 | 关注点 |
| --- | --- |
| 是否使用索引 | 没用索引可能慢 |
| 扫描行数 | 扫描太多说明过滤差 |
| 排序方式 | 是否额外文件排序 |
| 查询条件 | 是否和索引顺序匹配 |

## 常见坑

| 问题 | 处理 |
| --- | --- |
| 列表慢 | 看 EXPLAIN 和索引 |
| 深分页慢 | 改游标分页 |
| 重复数据 | 加唯一索引和幂等 |
| 事务无效 | 确认所有写入在同一事务 |
| 缓存旧数据 | 写数据库后删除相关缓存 |

## 项目建议

- 每个表和字段写清业务含义。
- 重要约束交给数据库保证，不只靠前端校验。
- 列表接口设计时同步设计索引。
- 迁移脚本写清变更原因和回滚风险。
- 慢查询先看执行计划，不要直接加缓存。

## 下一步学习

- [数据库学习导览](/database/introduction)
- [索引与查询优化](/database/indexes)
- [事务、锁与并发](/database/transactions)
- [数据库与缓存问题](/projects/issues-database)
