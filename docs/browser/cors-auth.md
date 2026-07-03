# 跨域与登录态

## 适合谁看

适合经常遇到“本地接口正常，浏览器报 CORS”“Cookie 没带上”“登录后刷新又退出”“401 和 403 分不清”的学习者。

跨域和登录态是前端项目上线后最容易出问题的地方。原因是它们同时涉及浏览器安全策略、前端请求配置、后端响应头、网关代理、Cookie 属性和部署域名。

## 同源策略

同源要求三个部分完全一致：

| 部分 | 示例 |
| --- | --- |
| 协议 | `https` |
| 域名 | `admin.example.com` |
| 端口 | `443` |

下面这些都不是同源：

```text
http://localhost:5173   和 http://localhost:8080
https://a.example.com   和 https://b.example.com
https://example.com     和 http://example.com
```

同源策略限制的是浏览器环境。服务端请求服务端不受这个限制。

## CORS 是什么

CORS 是浏览器和服务器之间的一套跨源访问协商机制。浏览器会根据请求类型和响应头判断前端代码能不能读取响应。

一个最小的 CORS 响应可能是：

```http
Access-Control-Allow-Origin: https://admin.example.com
```

如果需要跨源携带 Cookie，还需要：

```http
Access-Control-Allow-Credentials: true
```

前端也必须明确开启凭据：

```ts
await fetch('https://api.example.com/user', {
  credentials: 'include'
})
```

axios 对应配置：

```ts
axios.create({
  baseURL: 'https://api.example.com',
  withCredentials: true
})
```

只配一边不够。浏览器、前端请求和后端响应必须同时满足条件。

## 简单请求和预检请求

浏览器遇到某些跨源请求，会先发一个 `OPTIONS` 预检请求，询问服务器是否允许真正的请求。

常见触发原因：

- 使用 `PUT`、`PATCH`、`DELETE` 等方法。
- 自定义请求头，例如 `Authorization`。
- `Content-Type` 不是简单类型。

预检请求失败时，真正的业务请求不会发出去。此时后端日志里可能看不到你的 `POST /api/users`，只能看到 `OPTIONS /api/users`。

## Cookie 登录态

Cookie 登录态通常是：

1. 用户登录。
2. 服务端返回 `Set-Cookie`。
3. 浏览器保存 Cookie。
4. 后续请求自动携带 Cookie。
5. 服务端根据 Cookie 找到会话。

常见 Cookie 属性：

| 属性 | 作用 |
| --- | --- |
| `HttpOnly` | 禁止 JavaScript 读取，降低 XSS 窃取风险 |
| `Secure` | 只在 HTTPS 下发送 |
| `SameSite` | 控制跨站请求是否携带 |
| `Domain` | 控制哪些域名可用 |
| `Path` | 控制哪些路径可用 |
| `Max-Age` / `Expires` | 控制过期时间 |

### SameSite 的实际影响

| 值 | 含义 | 常见场景 |
| --- | --- | --- |
| `Strict` | 只在同站请求发送 | 安全性强，但登录跳转场景可能受影响 |
| `Lax` | 顶层导航等场景可发送 | 普通站点默认较常见 |
| `None` | 跨站也可发送 | 前后端不同站点、嵌入式系统 |

如果使用 `SameSite=None`，通常还需要 `Secure`，也就是 HTTPS。

## Token 登录态

Token 方案通常是：

1. 用户登录，后端返回 token。
2. 前端保存 token。
3. 请求时放到 `Authorization`。
4. 后端校验 token。

```ts
await fetch('/api/user', {
  headers: {
    Authorization: `Bearer ${token}`
  }
})
```

Token 方案的优点是前端控制明确，移动端、小程序、多端复用方便。缺点是 token 如果放在 LocalStorage，发生 XSS 时更容易被读取。

## 401 和 403

| 状态码 | 含义 | 前端处理 |
| --- | --- | --- |
| `401` | 没登录、登录失效、token 无效 | 清登录态，跳登录页 |
| `403` | 已登录，但没有权限 | 显示无权限或隐藏入口 |

不要把 403 也跳到登录页。否则用户会以为账号异常，实际只是没有某个操作权限。

## 开发环境代理

本地开发常用 Vite 代理：

```ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

代理让浏览器看到的是同源 `/api`，因此本地不跨域。但生产环境如果没有同域网关或 Nginx 反向代理，仍然会跨域。

## 生产环境推荐方案

后台管理系统优先推荐同域反向代理：

```text
https://admin.example.com
  /              -> 前端静态资源
  /api           -> 后端 API 网关
```

好处：

- 浏览器视角下同源，跨域问题少。
- Cookie 策略更简单。
- 运维可统一做鉴权、日志、限流、缓存。

如果业务必须前后端不同域，再严谨配置 CORS 和 Cookie。

## 实际项目问题

### 问题：后端已经允许 CORS，但浏览器仍然报错

**排查**

1. 响应头里的 `Access-Control-Allow-Origin` 是否等于当前页面 origin。
2. 是否把 `Access-Control-Allow-Origin` 写成 `*`，同时又开启了 credentials。
3. 预检请求 `OPTIONS` 是否返回了允许的方法和请求头。
4. 网关是否拦截了 `OPTIONS`。
5. 前端是否设置了 `withCredentials` 或 `credentials: 'include'`。

**解决方案**

带 Cookie 的跨源请求不能只写 `*`。应该返回明确来源：

```http
Access-Control-Allow-Origin: https://admin.example.com
Access-Control-Allow-Credentials: true
Vary: Origin
```

### 问题：Set-Cookie 有返回，但浏览器没保存

**排查**

- 是否 HTTPS。
- 是否设置了 `Secure`。
- `SameSite=None` 是否同时有 `Secure`。
- `Domain` 是否匹配当前站点。
- 是否跨源但前端没有开启 credentials。
- 是否被浏览器隐私策略或第三方 Cookie 策略限制。

### 问题：登录后接口还是 401

**排查顺序**

1. 登录接口 Response Headers 是否有 `Set-Cookie` 或 token。
2. Application 面板里是否真的保存成功。
3. 业务接口 Request Headers 是否带出 Cookie 或 Authorization。
4. 后端是否读取到了同一个会话。
5. token 是否过期或签名环境不一致。

## 最佳实践

- 后台管理项目优先同域反向代理。
- 401 和 403 必须分开处理。
- Cookie 登录态不要只看前端代码，要看 `Set-Cookie` 和后续请求头。
- Token 不要长期裸放敏感权限信息。
- CORS 配置要覆盖预检请求，不要只处理业务接口。
- 本地代理配置必须在部署文档中说明生产替代方案。

## 参考资料

- [MDN: Cross-Origin Resource Sharing](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CORS)
- [MDN: Set-Cookie header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie)

## 下一步学习

继续学习 [缓存策略](/browser/cache)。
