# HTTP API 开发

## 适合谁看

适合准备用 Node.js 写后端接口的学习者。

## API 服务基本结构

```text
请求
↓
路由
↓
参数校验
↓
业务服务
↓
数据访问
↓
响应
```

## Fastify 示例

Fastify 官方文档支持用 JSON Schema 校验请求。

```ts
import Fastify from 'fastify'

const app = Fastify({ logger: true })

app.get('/health', async () => {
  return { ok: true }
})

app.post('/users', {
  schema: {
    body: {
      type: 'object',
      required: ['username'],
      properties: {
        username: { type: 'string' },
        mobile: { type: 'string' }
      }
    }
  }
}, async (request) => {
  return {
    id: 1,
    ...request.body
  }
})

await app.listen({ port: 3000 })
```

## Express 示例

```ts
import express from 'express'

const app = express()

app.use(express.json())

app.get('/health', (req, res) => {
  res.json({ ok: true })
})

app.post('/users', async (req, res, next) => {
  try {
    const user = await createUser(req.body)
    res.status(201).json(user)
  } catch (error) {
    next(error)
  }
})
```

## 推荐分层

```text
src/
├─ routes/
├─ services/
├─ repositories/
├─ schemas/
├─ config/
└─ main.ts
```

| 层 | 职责 |
| --- | --- |
| routes | HTTP 路由、参数读取、响应 |
| services | 业务规则 |
| repositories | 数据访问 |
| schemas | 参数校验 |
| config | 环境配置 |

## 响应结构

建议统一：

```ts
interface ApiResult<T> {
  code: number
  message: string
  data: T
}
```

或者直接使用 HTTP 状态码表达结果。团队需要统一，不要混用多套风格。

## 实际项目常见问题

### 1. 参数没有校验

**后果**

业务层要处理大量脏数据。

**解决方案**

在路由入口做 schema 校验。

### 2. 路由里写了太多业务

**解决方案**

路由只处理 HTTP 细节，业务放 service。

### 3. 前端跨域

开发环境可配置 CORS，生产更推荐同域反向代理。

## 最佳实践

- 路由、业务、数据访问分层。
- 请求参数必须校验。
- 错误响应结构统一。
- 健康检查接口保留 `/health`。
- API 文档和前端类型保持同步。

## 下一步

继续学习 [鉴权与会话](/node/auth-session)。
