# Node.js 常见问题

## 1. 端口被占用

### 症状

启动时报 `EADDRINUSE`。

### 解决方案

换端口，或停止占用端口的进程。

```bash
lsof -i :3000
```

## 2. 环境变量是 undefined

### 排查

- `.env` 是否存在。
- 是否加载 dotenv。
- 变量名是否拼错。
- 部署平台是否注入变量。

## 3. async 错误导致进程异常

### 解决方案

路由中捕获错误，并交给统一错误处理。

Express：

```ts
app.get('/users', async (req, res, next) => {
  try {
    const users = await getUsers()
    res.json(users)
  } catch (error) {
    next(error)
  }
})
```

## 4. CORS 跨域

### 开发环境

可以配置 CORS 或前端代理。

### 生产环境

优先使用 Nginx 同域反向代理，减少浏览器跨域问题。

## 5. 请求体读取不到

Express 需要：

```ts
app.use(express.json())
```

Fastify 默认会处理 JSON，但仍需要确认请求头 `Content-Type: application/json`。

## 6. 日志看不出问题

### 解决方案

结构化日志至少包含：

- requestId。
- method。
- path。
- status。
- duration。
- error stack。

## 7. 线上返回错误堆栈

### 风险

泄露内部路径和实现。

### 解决方案

生产只返回通用错误，详细错误写日志。

## 8. 服务越来越慢

### 排查

- 是否有同步 CPU 密集任务。
- 数据库查询是否慢。
- 日志是否阻塞。
- 是否存在内存泄漏。
- 事件循环是否被阻塞。

## 快速排查表

| 问题 | 优先检查 |
| --- | --- |
| 启动失败 | 端口、环境变量、依赖 |
| 接口 500 | 日志和错误堆栈 |
| 请求体为空 | JSON 中间件和 Content-Type |
| 跨域 | CORS 或 Nginx 代理 |
| 线上慢 | 数据库、CPU、事件循环 |

## 最佳实践

- 启动时校验配置。
- 每个接口错误都进入统一处理。
- 日志结构化。
- 不在请求里执行长时间同步任务。
- 生产环境隐藏错误细节。
