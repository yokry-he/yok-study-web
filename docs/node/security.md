# Node.js 安全基础

## 适合谁看

适合已经能写 Node.js API，但还没有系统考虑安全边界的人：

- 参数直接进入 SQL 或查询条件。
- 错误栈直接返回给前端。
- 文件上传只看扩展名。
- 依赖包长期不更新。
- 管理后台接口只靠前端菜单隐藏。

安全不是上线前最后加一层中间件。它应该贯穿输入校验、鉴权授权、数据库访问、错误处理、依赖管理、日志审计和部署配置。

## 安全边界在哪里

Node.js API 常见边界：

```text
外部请求
↓
输入校验
↓
认证
↓
授权
↓
业务规则
↓
数据库 / 第三方服务
↓
响应输出
```

每一层都可能出问题。不要把安全完全寄托在某一个库或框架上。

## 输入校验

所有外部输入都不可信：

- body。
- query。
- params。
- headers。
- cookie。
- 文件。
- 第三方回调。

建议在 controller 或 validation 层先做校验：

```ts
function parsePageQuery(query) {
  const page = Number(query.page || 1)
  const pageSize = Number(query.pageSize || 20)

  if (!Number.isInteger(page) || page < 1) {
    throw new BadRequestError('page 不合法')
  }

  if (!Number.isInteger(pageSize) || pageSize < 1 || pageSize > 100) {
    throw new BadRequestError('pageSize 不合法')
  }

  return { page, pageSize }
}
```

不要把原始 query 直接传进 service 或 SQL。

## SQL 注入防护

危险：

```ts
await pool.query(`select * from users where name = '${name}'`)
```

安全：

```ts
await pool.query('select * from users where name = $1', [name])
```

ORM 也要使用安全 API，不要拼 raw SQL。如果必须 raw SQL，参数仍然要绑定。

## 鉴权和授权

安全项目里必须区分：

- 未登录：401。
- 没权限：403。
- 不存在或不可见资源：按业务决定 404 或 403。

接口层必须做授权，不能只靠前端隐藏按钮。

```ts
app.delete(
  '/api/users/:id',
  authRequired,
  permissionRequired('user:delete'),
  deleteUserController
)
```

权限判断要靠服务端当前用户和服务端数据，不要相信前端传来的角色或权限。

## 错误处理不要泄漏内部信息

不推荐：

```json
{
  "message": "select * from users where ...",
  "stack": "Error at ..."
}
```

推荐：

```json
{
  "code": "INTERNAL_ERROR",
  "message": "服务暂时不可用"
}
```

日志里可以记录：

- request id。
- user id。
- route。
- error stack。
- 环境。
- 版本号。

响应给用户的信息要克制，不能泄漏表名、SQL、文件路径、密钥、内部 IP。

## 文件上传安全

文件上传常见风险：

- 超大文件拖垮服务。
- 伪造扩展名。
- 上传脚本文件。
- 文件名路径穿越。
- 公开访问敏感文件。

建议：

- 限制大小。
- 校验 MIME 和真实内容。
- 文件名重新生成。
- 存储路径和静态服务隔离。
- 图片处理放到隔离流程。
- 上传结果记录审计日志。

不要把用户上传文件直接放进可执行目录。

## 依赖安全

Node.js 项目高度依赖 npm 生态。

建议：

- 锁定 lockfile。
- 定期更新依赖。
- 删除不用的包。
- 使用 `npm audit` 或供应链扫描。
- 不随意安装低维护、低下载、权限过大的包。
- CI 中检查高危依赖。

依赖包不是越多越好。一个小功能引入大包，要评估维护状态和安全风险。

## 环境变量和密钥

不要提交：

- 数据库密码。
- JWT secret。
- 第三方 API key。
- 云服务凭证。

启动时检查必需配置：

```ts
function requireEnv(name: string) {
  const value = process.env[name]

  if (!value) {
    throw new Error(`缺少环境变量：${name}`)
  }

  return value
}
```

生产密钥要由部署平台或密钥管理系统注入，不要写进代码仓库。

## HTTP 安全头

常见安全头：

- `Content-Security-Policy`
- `X-Content-Type-Options`
- `Referrer-Policy`
- `Strict-Transport-Security`

如果使用 Express，可以考虑成熟中间件统一设置。但安全头要结合业务验证，不能盲目复制。

## 限流和防爆破

登录、验证码、短信、导出、搜索等接口要考虑限流。

```text
同一 IP 登录失败过多 -> 限制
同一账号登录失败过多 -> 限制
短信发送频率过高 -> 限制
导出任务过多 -> 排队或拒绝
```

没有限流的登录接口容易被撞库，没有限制的导出接口容易拖垮数据库。

## 日志与审计

安全相关操作要记录：

- 登录成功和失败。
- 退出登录。
- 密码修改。
- 权限变更。
- 删除、导出、审批。
- 管理员操作。

日志要避免记录明文密码、完整 token、敏感个人信息。

## 实际项目常见问题

### 1. 接口只靠前端权限控制

**风险**

用户可以直接调用接口。

**解决方案**

服务端每个敏感接口做权限判断。

### 2. 错误响应暴露 SQL

**解决方案**

统一错误处理中间件区分日志和用户响应。

### 3. 文件上传导致服务器磁盘打满

**解决方案**

限制大小、类型、数量，存储到对象存储或隔离目录，并做清理策略。

### 4. .env 被提交

**解决方案**

`.env` 加入 `.gitignore`，提供 `.env.example`，密钥通过部署平台注入。

### 5. npm 包漏洞长期不处理

**解决方案**

定期依赖审计，优先处理可利用的高危漏洞。

## 最佳实践

- 所有外部输入先校验。
- SQL 使用参数化查询。
- 鉴权和授权都在服务端执行。
- 错误响应不泄漏内部信息。
- 文件上传限制大小、类型和存储路径。
- 密钥不进入代码仓库。
- 登录、导出、短信等高风险接口限流。
- 依赖安全进入常规维护流程。

## 学习检查

学完本节后，你应该能回答：

- Node API 的安全边界有哪些。
- 为什么参数化查询能降低 SQL 注入风险。
- 为什么权限不能只放前端。
- 错误日志和错误响应应该有什么区别。
- 文件上传和依赖包为什么是高风险点。

## 参考资料

- [OWASP Nodejs Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Nodejs_Security_Cheat_Sheet.html)
- [OWASP REST Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)

## 下一步学习

继续学习 [项目结构与部署](/node/project-deployment)，把安全配置、环境变量、启动命令和上线检查串起来。
