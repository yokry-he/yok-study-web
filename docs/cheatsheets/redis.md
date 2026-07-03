# Redis 速查

## 连接

| 命令 | 用途 |
| --- | --- |
| `redis-cli` | 连接本机 Redis |
| `redis-cli -h host -p 6379` | 指定地址和端口 |
| `AUTH password` | 密码认证 |
| `PING` | 检查连接 |
| `SELECT 1` | 切换数据库 |

生产环境不要在聊天记录或截图里暴露 Redis 密码。

## Key 通用命令

| 命令 | 用途 |
| --- | --- |
| `EXISTS key` | key 是否存在 |
| `DEL key` | 删除 key |
| `TYPE key` | 查看类型 |
| `TTL key` | 查看剩余过期秒数 |
| `EXPIRE key 60` | 设置 60 秒过期 |
| `PERSIST key` | 移除过期时间 |
| `SCAN 0 MATCH user:* COUNT 100` | 分批扫描 key |

线上不要用 `KEYS *` 扫全库，优先用 `SCAN`。

## String

| 命令 | 用途 |
| --- | --- |
| `SET token abc EX 3600` | 设置值并过期 |
| `GET token` | 获取值 |
| `MGET a b c` | 批量获取 |
| `INCR counter` | 自增 |
| `DECR counter` | 自减 |
| `SETNX lock 1` | 不存在时设置 |

缓存登录态或验证码时，一定要设置过期时间。

## Hash

| 命令 | 用途 |
| --- | --- |
| `HSET user:1 name Tom` | 设置字段 |
| `HGET user:1 name` | 获取字段 |
| `HGETALL user:1` | 获取全部字段 |
| `HMGET user:1 name age` | 获取多个字段 |
| `HDEL user:1 name` | 删除字段 |
| `HINCRBY user:1 score 1` | 字段自增 |

Hash 适合保存对象的多个字段，但过期时间作用在整个 key 上。

## List

| 命令 | 用途 |
| --- | --- |
| `LPUSH queue item` | 左侧加入 |
| `RPUSH queue item` | 右侧加入 |
| `LPOP queue` | 左侧弹出 |
| `RPOP queue` | 右侧弹出 |
| `LRANGE queue 0 -1` | 查看范围 |
| `LLEN queue` | 列表长度 |

List 可用于简单队列，但复杂任务队列通常需要更完整的可靠性设计。

## Set

| 命令 | 用途 |
| --- | --- |
| `SADD roles admin` | 添加成员 |
| `SREM roles admin` | 删除成员 |
| `SMEMBERS roles` | 查看成员 |
| `SISMEMBER roles admin` | 判断成员 |
| `SCARD roles` | 成员数量 |

Set 适合标签、去重、权限集合等场景。

## ZSet

| 命令 | 用途 |
| --- | --- |
| `ZADD rank 100 user1` | 添加分数 |
| `ZRANGE rank 0 9 WITHSCORES` | 正序排名 |
| `ZREVRANGE rank 0 9 WITHSCORES` | 倒序排名 |
| `ZSCORE rank user1` | 查看分数 |
| `ZREM rank user1` | 删除成员 |

ZSet 适合排行榜、延迟任务、按时间排序的数据。

## 缓存常见模式

读取：

```text
读缓存
↓
命中直接返回
↓
未命中读数据库
↓
写缓存并设置过期
```

更新：

```text
更新数据库
↓
删除缓存
```

优先删除缓存，而不是先更新缓存。这样更容易避免缓存和数据库不一致。

## 常见问题

| 问题 | 处理 |
| --- | --- |
| 缓存穿透 | 参数校验、缓存空值、布隆过滤器 |
| 缓存击穿 | 热点 key 互斥重建、逻辑过期 |
| 缓存雪崩 | 过期时间加随机值、限流、降级 |
| 内存上涨 | 设置 maxmemory、淘汰策略、排查大 key |
| key 太多 | 统一命名、设置过期、避免无界增长 |

## 参考资料

- [Redis commands](https://redis.io/docs/latest/commands/)
- [Redis SET](https://redis.io/docs/latest/commands/set/)
- [Redis EXPIRE](https://redis.io/docs/latest/commands/expire/)

## 延伸学习

- [Redis 缓存与数据结构](/database/redis)
- [数据库与缓存问题](/projects/issues-database)
