# 数据库集成

## 适合谁看

适合已经能写 Node.js API，但开始连接 MySQL、PostgreSQL、Redis 或 ORM 时不够清楚的人：

- 不知道连接池是什么。
- 每个请求都新建数据库连接。
- SQL 参数直接字符串拼接。
- 事务写了但偶尔不生效。
- 数据库错误和业务错误混在一起。

Node.js 做后端时，数据库集成不是“能查出数据”就结束。项目里要考虑连接池、参数化查询、事务、错误处理、迁移、环境配置和数据边界。

## 推荐分层

```text
route
↓
controller
↓
service
↓
repository / dao
↓
database client / orm
```

职责：

| 层 | 负责 |
| --- | --- |
| controller | 处理 HTTP 输入输出 |
| service | 业务规则、权限、事务编排 |
| repository | 数据查询和持久化 |
| database client | 连接池、SQL、驱动能力 |

不要在路由里直接写一堆 SQL。那样很快会让鉴权、校验、事务和响应混在一起。

## 连接池

数据库连接是有限资源。连接池负责复用连接。

PostgreSQL 示例：

```ts
import pg from 'pg'

const { Pool } = pg

export const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 10
})
```

查询：

```ts
const result = await pool.query('select * from users where id = $1', [id])
```

连接池一般在应用启动时创建一次，不要在每个请求里创建。

## 参数化查询

不要拼接 SQL：

```ts
await pool.query(`select * from users where username = '${username}'`)
```

使用参数：

```ts
await pool.query('select * from users where username = $1', [username])
```

好处：

- 降低 SQL 注入风险。
- 让数据库驱动正确处理转义。
- 查询结构更清晰。

## repository 示例

```ts
type UserRecord = {
  id: number
  username: string
  enabled: boolean
}

export async function findUserById(id: number): Promise<UserRecord | null> {
  const result = await pool.query<UserRecord>(
    'select id, username, enabled from users where id = $1',
    [id]
  )

  return result.rows[0] ?? null
}
```

repository 只处理数据，不负责 HTTP 状态码，也不直接返回前端响应。

## service 示例

```ts
export async function getUserProfile(id: number) {
  const user = await userRepository.findUserById(id)

  if (!user) {
    throw new NotFoundError('用户不存在')
  }

  return {
    id: user.id,
    username: user.username,
    enabled: user.enabled
  }
}
```

业务错误在 service 层表达，再由统一错误处理中间件转换为 HTTP 响应。

## 事务

涉及多条写操作时需要事务。

重要规则：同一个事务必须使用同一个数据库 client。

```ts
const client = await pool.connect()

try {
  await client.query('begin')
  await client.query('insert into orders(user_id) values($1)', [userId])
  await client.query('update users set order_count = order_count + 1 where id = $1', [userId])
  await client.query('commit')
} catch (error) {
  await client.query('rollback')
  throw error
} finally {
  client.release()
}
```

不要在事务中混用 `pool.query`，因为它可能拿到不同连接。

## ORM 怎么看

ORM 可以提高效率，但不能代替数据库理解。

常见收益：

- 类型提示。
- CRUD 快速开发。
- migration 管理。
- 关系查询封装。

常见风险：

- 不理解生成 SQL。
- N+1 查询。
- 复杂查询难优化。
- migration 随意生成但没人审。

建议：

- 初学时至少能读懂 SQL。
- 重要查询要看执行计划。
- migration 要写清变更原因。
- 不要让 ORM 模型和接口响应完全绑定。

## 环境变量

不要把数据库密码写进代码。

```env
DATABASE_URL=postgres://user:password@localhost:5432/app
```

读取：

```ts
const databaseUrl = process.env.DATABASE_URL

if (!databaseUrl) {
  throw new Error('缺少 DATABASE_URL')
}
```

启动时就检查关键配置，不要等到第一个请求才报错。

## 实际项目常见问题

### 1. 数据库连接数很快打满

**原因**

每个请求都创建连接，或者连接没有释放。

**解决方案**

- 使用全局连接池。
- 事务后 `finally release()`。
- 设置连接池上限。
- 排查慢查询和连接泄漏。

### 2. 事务没有生效

**原因**

事务中使用了 `pool.query`，不同语句跑在不同连接上。

**解决方案**

事务内全部使用同一个 `client`。

### 3. SQL 注入风险

**原因**

把用户输入拼进 SQL 字符串。

**解决方案**

使用参数化查询或 ORM 安全 API。

### 4. 接口字段和数据库字段耦合

**问题**

数据库字段一改，前端响应也跟着乱。

**解决方案**

repository 返回 record，service 转换成业务对象或响应 DTO。

### 5. 数据库错误直接返回给前端

不要把 SQL、表名、栈信息直接返回用户。记录日志，给用户返回明确但不泄漏内部结构的错误。

## 最佳实践

- 应用启动时创建连接池。
- 所有用户输入都参数化。
- 多写操作使用事务。
- 事务内使用同一个 client。
- controller、service、repository 分层。
- migration 要写清原因和回滚风险。
- 数据库错误进入日志，不直接暴露给用户。

## 学习检查

学完本节后，你应该能回答：

- 为什么不能每个请求创建一个数据库连接。
- 参数化查询解决什么问题。
- 为什么事务必须使用同一个 client。
- ORM 能解决什么，不能替你解决什么。
- 数据库错误应该在哪一层转换成 HTTP 响应。

## 参考资料

- [node-postgres: Pool API](https://node-postgres.com/apis/pool)
- [node-postgres: Transactions](https://node-postgres.com/features/transactions)
- [Node.js: Environment Variables](https://nodejs.org/api/environment_variables.html)
- [Node.js Learn: Read environment variables](https://nodejs.org/learn/command-line/how-to-read-environment-variables-from-nodejs)

## 下一步学习

继续学习 [错误处理与日志](/node/error-logging)，把数据库错误、业务错误和系统错误纳入统一排查流程。
