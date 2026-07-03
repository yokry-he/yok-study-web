# Redis 缓存与数据结构

## 适合谁看

适合已经会使用 MySQL 或 PostgreSQL，但需要缓存、会话、限流、排行榜、分布式锁或消息队列能力的人。

Redis 是内存数据库，速度很快，但不能简单理解为“更快的 MySQL”。Redis 更适合高频、短路径、可重建或特定数据结构场景。

## Redis 常见用途

| 用途 | 常见数据结构 | 示例 |
| --- | --- | --- |
| 缓存 | String、Hash | 用户信息、配置、字典 |
| 会话 | String、Hash | 登录会话、验证码 |
| 限流 | String、Sorted Set | 短信发送频率、接口限流 |
| 排行榜 | Sorted Set | 积分榜、热度榜 |
| 去重 | Set、Bitmap | 签到、浏览记录 |
| 队列 | List、Stream | 异步任务、事件流 |
| 分布式锁 | String | 避免重复任务执行 |

## 数据结构选择

### String

适合简单 key-value：

```text
user:profile:1001 -> {"id":1001,"name":"Alice"}
```

常用命令：

```bash
set user:profile:1001 '{"id":1001,"name":"Alice"}' ex 3600
get user:profile:1001
```

### Hash

适合对象字段：

```bash
hset user:1001 id 1001 name Alice status active
hgetall user:1001
```

如果对象字段经常局部更新，Hash 比整段 JSON String 更合适。

### Set

适合去重集合：

```bash
sadd article:liked_users:88 1001 1002
sismember article:liked_users:88 1001
```

### Sorted Set

适合排行榜：

```bash
zadd leaderboard 9800 user:1001
zrevrange leaderboard 0 9 withscores
```

### Stream

适合事件流和消费组。学习阶段可以先掌握 List 和普通队列，复杂异步任务再研究 Stream。

## 过期时间

缓存一定要考虑过期。

```bash
set product:detail:1001 '{"id":1001}' ex 600
```

没有过期时间的缓存会带来：

- 内存持续增长。
- 数据长期不一致。
- 清理困难。

## 淘汰策略

Redis 可以在达到内存上限后自动淘汰 key。常见策略包括：

| 策略 | 含义 |
| --- | --- |
| noeviction | 不淘汰，写入报错 |
| allkeys-lru | 从所有 key 中淘汰最近最少使用 |
| volatile-lru | 从设置了过期时间的 key 中淘汰最近最少使用 |
| allkeys-lfu | 从所有 key 中淘汰使用频率低的 |
| volatile-ttl | 优先淘汰快过期的 key |

如果 Redis 主要当缓存，通常要配置 maxmemory 和合适淘汰策略。如果 Redis 存会话或关键状态，淘汰策略要非常谨慎。

## 缓存一致性

典型读流程：

```text
读 Redis
↓ 未命中
读数据库
↓
写 Redis
↓
返回结果
```

典型写流程：

```text
更新数据库
↓
删除缓存
```

大多数业务里，写后删除缓存比写后更新缓存更稳，因为缓存可能由多个表组合而来，更新逻辑容易遗漏。

## 缓存穿透、击穿、雪崩

| 问题 | 含义 | 常见方案 |
| --- | --- | --- |
| 穿透 | 查询不存在的数据，每次都打到数据库 | 缓存空值、参数校验、布隆过滤器 |
| 击穿 | 热点 key 过期，大量请求同时打数据库 | 互斥锁、逻辑过期、预热 |
| 雪崩 | 大量 key 同时过期 | 过期时间加随机抖动、分批预热 |

## 实际项目问题

### 问题：缓存里还是旧数据

**原因**

更新数据库后没有删除缓存，或者删除了错误 key。

**解决方案**

统一缓存 key 生成函数，不要散落字符串：

```ts
export function userProfileKey(userId: number) {
  return `user:profile:${userId}`
}
```

写操作成功后删除对应缓存：

```ts
await updateUser(userId, payload)
await redis.del(userProfileKey(userId))
```

### 问题：Redis 内存突然满了

**排查**

- 是否有大 key。
- 是否有大量 key 没有 TTL。
- maxmemory 是否配置。
- 淘汰策略是否符合业务。

**处理**

先定位 key，不要直接 flushall。

## 最佳实践

- 缓存 key 命名统一，包含业务域和 ID。
- 缓存值要有版本意识，结构变化时避免旧数据解析失败。
- 大多数缓存设置 TTL。
- 更新数据库后优先删除缓存。
- 不要把 Redis 当成唯一持久化数据库，除非明确设计了持久化和恢复方案。
- 限流、锁和队列要考虑异常释放、超时和幂等。

## 参考资料

- [Redis: Key eviction](https://redis.io/docs/latest/develop/reference/eviction/)
- [Redis: Persistence](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/)

## 下一步学习

继续学习 [数据建模与表设计](/database/modeling)。
