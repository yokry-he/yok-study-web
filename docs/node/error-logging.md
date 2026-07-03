# 错误处理与日志

## 适合谁看

适合准备让 Node API 服务稳定运行、方便排错的学习者。

## 错误分类

| 类型 | 示例 | 处理 |
| --- | --- | --- |
| 参数错误 | 缺少 username | 返回 400 |
| 未登录 | token 缺失 | 返回 401 |
| 无权限 | 权限不足 | 返回 403 |
| 不存在 | 用户不存在 | 返回 404 |
| 系统错误 | 数据库异常 | 返回 500 并记录日志 |

## Express 错误处理中间件

Express 官方文档说明，错误处理中间件需要四个参数。

```ts
app.use((error, req, res, next) => {
  console.error(error)

  res.status(500).json({
    code: 500,
    message: '服务器内部错误'
  })
})
```

即使不用 `next`，也要保留四个参数，否则 Express 不会把它识别为错误处理中间件。

## Fastify 错误处理

```ts
app.setErrorHandler((error, request, reply) => {
  request.log.error(error)

  reply.status(500).send({
    code: 500,
    message: '服务器内部错误'
  })
})
```

## 日志应该记录什么

- 请求路径。
- 请求方法。
- 状态码。
- 耗时。
- 错误堆栈。
- 用户 id。
- trace id。

不要记录：

- 密码。
- token。
- 身份证号。
- 支付信息。
- 其他敏感数据。

## 实际项目常见问题

### 1. async 错误没有被捕获

Express 中需要 `try/catch` 后 `next(error)`，或使用支持异步错误处理的框架/封装。

### 2. 生产环境返回了错误堆栈

**风险**

泄露路径、依赖和内部实现。

**解决方案**

生产环境只返回通用错误消息，详细堆栈写日志。

### 3. 日志太多但没用

**原因**

没有结构化字段，无法检索。

**解决方案**

使用结构化日志，至少包含 requestId、userId、method、path、status、duration。

## 最佳实践

- 统一错误响应。
- 业务错误和系统错误区分。
- 生产不返回堆栈。
- 日志不记录敏感信息。
- 每个请求有 trace id 或 request id。

## 下一步

继续学习 [测试策略](/node/testing)。
