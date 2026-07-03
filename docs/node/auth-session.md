# 鉴权与会话

## 适合谁看

适合已经能写 Node.js HTTP API，但开始处理登录、权限、token、Cookie、401/403 时不够稳定的人：

- 登录成功后不知道前后端应该保存什么。
- 401 和 403 经常混用。
- token 放哪里、过期怎么处理没有统一方案。
- 管理后台按钮权限和接口权限不一致。
- 不知道 Session、JWT、OAuth、OIDC 的边界。

鉴权不是“登录接口返回一个 token”这么简单。真实项目要同时考虑身份认证、权限判断、会话过期、退出登录、刷新策略和审计日志。

## 先区分认证和授权

| 概念 | 问题 | 示例 |
| --- | --- | --- |
| Authentication | 你是谁 | 用户名密码登录、短信登录、SSO |
| Authorization | 你能做什么 | 能否查看用户、删除订单、导出数据 |

项目里常见错误是：只要用户登录了，就允许访问所有接口。正确做法是每个敏感接口都要做权限判断。

## 常见登录态方案

| 方案 | 特点 | 适合 |
| --- | --- | --- |
| Cookie Session | 服务端保存会话，浏览器自动带 Cookie | 传统 Web、后台系统 |
| JWT Access Token | token 自包含，服务端验证签名 | 前后端分离 API、移动端 |
| Access + Refresh Token | 短 token + 刷新 token | 登录态较长的应用 |
| OIDC / SSO | 第三方身份提供商 | 企业统一登录 |

初学项目不要一开始就追求复杂 OAuth。后台管理系统可以先把 Cookie Session 或 Access Token 做清楚。

## 推荐接口分层

```text
router
↓
auth middleware
↓
permission middleware
↓
controller
↓
service
↓
database
```

职责：

| 层 | 负责 |
| --- | --- |
| auth middleware | 解析登录态，识别当前用户 |
| permission middleware | 判断当前用户是否有权限 |
| controller | 处理请求参数和响应 |
| service | 业务规则 |
| database | 数据查询和持久化 |

不要把权限判断散落在控制器各处，更不要只放在前端菜单里。

## 401 和 403

| 状态码 | 含义 | 前端常见处理 |
| --- | --- | --- |
| 401 | 未登录或登录过期 | 清理登录态，跳登录 |
| 403 | 已登录但没有权限 | 展示无权限提示 |

不要把所有权限问题都返回 401，否则前端会错误地跳登录。

## 登录接口示例结构

示意代码：

```ts
app.post('/api/login', async (req, res, next) => {
  try {
    const { username, password } = req.body
    const user = await userService.verifyPassword(username, password)

    if (!user) {
      res.status(401).json({ message: '用户名或密码错误' })
      return
    }

    const token = await authService.createAccessToken(user)

    res.json({
      accessToken: token,
      user: {
        id: user.id,
        username: user.username,
        roles: user.roles
      }
    })
  } catch (error) {
    next(error)
  }
})
```

注意：

- 密码校验放 service。
- 响应里不要返回密码 hash。
- 登录失败提示不要暴露“用户名存在但密码错”这类信息。
- 生产环境要限制登录频率。

## 鉴权中间件

```ts
async function authRequired(req, res, next) {
  const header = req.headers.authorization
  const token = header?.startsWith('Bearer ') ? header.slice(7) : ''

  if (!token) {
    res.status(401).json({ message: '未登录' })
    return
  }

  try {
    const user = await authService.verifyAccessToken(token)
    req.currentUser = user
    next()
  } catch {
    res.status(401).json({ message: '登录已过期' })
  }
}
```

真实项目中要给 `req.currentUser` 补类型声明，避免后续代码到处 `any`。

## 权限中间件

```ts
function permissionRequired(action) {
  return (req, res, next) => {
    const user = req.currentUser

    if (!user.permissions.includes(action)) {
      res.status(403).json({ message: '没有操作权限' })
      return
    }

    next()
  }
}
```

使用：

```ts
app.delete(
  '/api/users/:id',
  authRequired,
  permissionRequired('user:delete'),
  deleteUserController
)
```

按钮权限只能改善前端体验，接口权限才是安全边界。

## Cookie Session 注意点

如果使用 Cookie：

```http
Set-Cookie: sid=xxx; HttpOnly; Secure; SameSite=Lax; Path=/
```

关键属性：

| 属性 | 作用 |
| --- | --- |
| `HttpOnly` | 前端 JS 不能读取，降低 XSS 后 token 被直接窃取的风险 |
| `Secure` | 只在 HTTPS 发送 |
| `SameSite` | 降低 CSRF 风险 |
| `Path` / `Domain` | 限制 Cookie 作用范围 |

跨域 Cookie 还要同时处理 CORS 和 `credentials`，细节可回到浏览器模块学习。

## Token 过期和刷新

常见设计：

```text
access token：短期有效
refresh token：长期一点，仅用于刷新
```

注意：

- refresh token 要能失效。
- 退出登录要清理服务端状态或黑名单。
- 多端登录要考虑设备维度。
- 刷新接口要防并发重复刷新。
- 发现异常登录要能强制下线。

如果只是学习项目，可以先实现短 access token + 重新登录，不要一开始做复杂刷新链路。

## 密码存储

不要明文存储密码。

服务端只保存密码哈希，并使用适合密码存储的算法。业务文档里至少要写清：

- 密码不明文存储。
- 登录只比较哈希结果。
- 密码重置要有有效期。
- 管理员也不能查看用户原密码。

## 实际项目常见问题

### 1. 前端隐藏按钮，但接口仍然能调用

**原因**

只做了前端权限，没有后端权限。

**解决方案**

后端接口必须检查当前用户权限。前端权限只负责体验。

### 2. 登录过期后接口一直重复请求

**原因**

多个接口同时返回 401，前端重复刷新 token 或重复跳登录。

**后端建议**

明确返回 401 和统一错误码。刷新接口和普通接口分开。

### 3. 403 被前端当成 401

**原因**

后端状态码设计不清。

**解决方案**

未登录返回 401，无权限返回 403。

### 4. token 永不过期

**风险**

token 泄漏后长期可用。

**解决方案**

设置过期时间，并提供退出、刷新、强制失效策略。

## 最佳实践

- 区分认证和授权。
- 每个敏感接口都做服务端权限判断。
- 401 和 403 明确区分。
- 密码不明文保存。
- Cookie 设置 `HttpOnly`、`Secure`、`SameSite`。
- token 要有过期和失效策略。
- 登录、退出、权限变更要记录安全日志。

## 学习检查

学完本节后，你应该能回答：

- 认证和授权有什么区别。
- 401 和 403 应该怎么用。
- 为什么前端按钮权限不是安全边界。
- Cookie Session 和 JWT 各自适合什么场景。
- token 过期、刷新和退出登录要考虑什么。

## 参考资料

- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [OWASP Authorization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html)
- [OWASP Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)

## 下一步学习

继续学习 [数据库集成](/node/database-integration)，把用户、角色、权限和业务数据持久化到数据库中。
