# 事务、锁与并发

## 适合谁看

适合已经会写 SQL，但对“为什么要事务”“为什么会死锁”“为什么库存会扣成负数”还不清楚的人。

事务和锁解决的是并发下的数据正确性问题。项目越接近真实业务，越不能只靠单条 SQL 和应用层判断。

## ACID

| 特性 | 含义 |
| --- | --- |
| Atomicity 原子性 | 要么全部成功，要么全部失败 |
| Consistency 一致性 | 事务前后数据满足约束 |
| Isolation 隔离性 | 并发事务之间互相隔离 |
| Durability 持久性 | 提交后数据可靠保存 |

## 什么时候需要事务

需要事务的场景：

- 创建订单同时写订单明细。
- 支付成功后更新订单和流水。
- 扣库存同时记录库存日志。
- 修改角色权限同时写审计日志。
- 转账扣减一方余额并增加另一方余额。

不需要事务或不应放进事务的操作：

- 远程 HTTP 调用。
- 文件上传。
- 发送短信。
- 调用第三方支付接口。
- 大量慢查询。

这些操作可以在事务前后做，但不要长时间占用数据库锁。

## 短事务原则

错误示例：

```text
开启事务
↓
锁定订单
↓
调用支付平台
↓
更新订单
↓
提交事务
```

问题：支付平台耗时期间，订单行一直被锁住。

更好的做法：

```text
调用支付平台
↓
开启事务
↓
按条件更新订单
↓
写支付流水
↓
提交事务
```

事务只包住必须原子完成的数据库操作。

## 条件更新防并发

库存扣减示例：

```sql
update product_stock
set stock = stock - 1
where product_id = 1001
  and stock > 0;
```

应用层检查影响行数：

- 影响 1 行：扣减成功。
- 影响 0 行：库存不足或商品不存在。

不要先查库存再更新库存：

```text
select stock
if stock > 0
update stock = stock - 1
```

并发下两个请求都可能读到同一个库存值。

## 悲观锁

```sql
begin;

select id, status
from orders
where id = 1001
for update;

update orders
set status = 'paid'
where id = 1001
  and status = 'pending';

commit;
```

`for update` 会锁住选中的行。适合必须串行处理的关键记录。

## 乐观锁

给表增加版本号：

```sql
alter table orders
  add column version integer not null default 0;
```

更新时带版本：

```sql
update orders
set status = 'paid',
    version = version + 1
where id = 1001
  and version = 3;
```

影响 0 行表示数据已被别人修改，应用层重新读取后再决定。

## 隔离级别

常见隔离级别：

| 隔离级别 | 特点 |
| --- | --- |
| Read Committed | 多数数据库默认或常用，避免脏读 |
| Repeatable Read | 同一事务内多次读取结果更稳定 |
| Serializable | 最强隔离，成本最高，可能需要重试 |

不要为了“更安全”无脑使用最高隔离级别。隔离越强，并发成本越高。

## 死锁

死锁常见原因：两个事务以不同顺序锁定资源。

```text
事务 A：锁订单 1 -> 等订单 2
事务 B：锁订单 2 -> 等订单 1
```

预防方式：

- 所有代码按相同顺序访问资源。
- 事务尽量短。
- 批量更新按主键排序。
- 捕获死锁错误并安全重试。

## 实际项目问题

### 问题：支付回调重复导致订单状态异常

**原因**

支付平台可能重复回调，接口没有幂等保护。

**解决方案**

```sql
begin;

update orders
set status = 'paid',
    paid_at = now()
where id = $1
  and status = 'pending';

-- 应用层检查影响行数：只有第一次 pending -> paid 成功。

insert into payment_log(order_id, event_type, raw_payload, created_at)
values ($1, 'payment_callback', $2, now());

commit;
```

### 问题：批量导入时锁住后台操作

**原因**

单个事务导入大量数据，持锁时间太长。

**解决方案**

- 分批导入。
- 每批单独提交。
- 避开业务高峰。
- 对唯一冲突使用明确策略。

## 最佳实践

- 事务只包数据库原子操作。
- 不要在事务里等待外部服务。
- 使用条件更新处理库存、状态机和幂等。
- 锁定多条记录时保持固定顺序。
- 死锁要记录上下文并允许安全重试。
- 慢事务和长事务要进入监控。

## 参考资料

- [PostgreSQL: Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)
- [PostgreSQL: Transaction Management](https://www.postgresql.org/docs/current/tutorial-transactions.html)
- [MySQL: InnoDB Locking](https://dev.mysql.com/doc/en/innodb-locking.html)

## 下一步学习

继续学习 [迁移、种子与版本治理](/database/migration)。
