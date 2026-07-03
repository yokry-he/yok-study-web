# 浏览器安全基础

## 适合谁看

适合已经能写前端页面和接口请求，但对 XSS、CSRF、CSP、Cookie 安全属性、HTTPS、安全上下文还不够熟悉的学习者。

前端安全不是“后端的事”。浏览器负责执行用户拿到的 HTML、CSS 和 JavaScript，也负责限制跨源访问、Cookie 发送、剪贴板、定位、摄像头、Service Worker 等敏感能力。理解浏览器安全边界，才能避免把 token、用户数据和页面控制权暴露出去。

## 你会学到什么

- 浏览器为什么要有同源策略和安全上下文。
- XSS、CSRF、点击劫持分别是什么。
- Cookie 的 `HttpOnly`、`Secure`、`SameSite` 怎么影响登录态。
- CSP 能解决什么，不能解决什么。
- 前端项目里哪些写法容易引入安全风险。

## 同源策略

同源由三个部分共同决定：

```text
协议 + 域名 + 端口
```

例如：

```text
https://example.com:443
```

只要协议、域名或端口有一个不同，就是跨源。浏览器会限制跨源脚本读取响应内容，但不等于请求完全不能发出。

项目里常见影响：

- 跨域接口需要 CORS。
- Cookie 的 Domain、Path、SameSite 会影响登录态。
- iframe、图片、脚本、字体等资源加载有不同规则。
- 某些 Web API 只允许在安全上下文中使用。

## 安全上下文

很多敏感 API 只在安全上下文中可用，通常意味着页面需要通过 HTTPS 访问。

常见受影响能力：

- Service Worker。
- Clipboard API。
- Geolocation。
- Web Crypto。
- Notification。
- Camera 和 microphone。

本地开发时 `localhost` 通常被浏览器特殊对待，但线上必须使用 HTTPS。

## XSS

XSS 是攻击者让恶意脚本在用户浏览器里执行。

常见来源：

- 把用户输入直接当 HTML 渲染。
- 使用 `innerHTML` 插入未清洗内容。
- Markdown、富文本、评论内容没有净化。
- 第三方脚本被污染。
- URL 参数被拼进页面。

危险示例：

```ts
contentEl.innerHTML = userInput
```

更安全的方式：

```ts
contentEl.textContent = userInput
```

如果业务必须渲染 HTML，需要使用可靠的 HTML sanitizer，并限制允许的标签和属性。

### Vue / React 中的注意点

Vue 的模板插值和 React 的 JSX 文本默认会转义内容，但下面这些能力需要格外谨慎：

```vue
<div v-html="html"></div>
```

```tsx
<div dangerouslySetInnerHTML={{ __html: html }} />
```

只有在内容来源可信且经过净化时才使用。

## CSRF

CSRF 是攻击者诱导用户浏览器带着已有登录态去发起请求。

它常见于 Cookie Session 场景，因为浏览器会自动带 Cookie。

降低风险的方式：

- Cookie 设置合适的 `SameSite`。
- 关键写操作校验 CSRF token。
- 写操作使用非简单请求并校验自定义请求头。
- 后端校验 Origin 或 Referer。
- 不用 GET 做有副作用的操作。

错误示例：

```text
GET /api/delete-user?id=1
```

删除、支付、审批、修改权限这类操作必须使用明确的写操作方法，并在服务端做权限和 CSRF 防护。

## Cookie 安全属性

常见 Set-Cookie：

```http
Set-Cookie: session_id=abc; HttpOnly; Secure; SameSite=Lax; Path=/
```

| 属性 | 作用 |
| --- | --- |
| `HttpOnly` | JavaScript 不能读取，降低 XSS 窃取风险 |
| `Secure` | 只在 HTTPS 下发送 |
| `SameSite` | 限制跨站请求携带 Cookie |
| `Domain` | 限定可发送的域 |
| `Path` | 限定可发送的路径 |
| `Max-Age` / `Expires` | 控制有效期 |

跨站点嵌入或跨站接口如果确实需要 Cookie，通常会涉及：

```http
SameSite=None; Secure
```

同时 CORS 不能使用 `Access-Control-Allow-Origin: *` 搭配凭证。

## CSP

CSP 是 Content Security Policy，用来限制页面可以加载哪些脚本、样式、图片、字体、iframe 等资源。

示例：

```http
Content-Security-Policy:
  default-src 'self';
  script-src 'self';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  object-src 'none';
  base-uri 'self';
  frame-ancestors 'none';
```

CSP 能降低 XSS 的破坏范围，但不能替代输入校验和输出转义。

落地建议：

- 先用 Report-Only 模式收集违规。
- 再逐步收紧策略。
- 第三方脚本要列清来源。
- 禁止不必要的 `unsafe-inline` 和 `unsafe-eval`。
- 对 iframe 嵌入使用 `frame-ancestors`。

## 点击劫持

点击劫持是把你的页面放进透明 iframe，诱导用户点击。

防护方式：

```http
Content-Security-Policy: frame-ancestors 'none'
```

或旧方案：

```http
X-Frame-Options: DENY
```

后台系统、支付页面、权限配置页面通常不应该允许被第三方 iframe 嵌入。

## 前端项目安全检查清单

| 检查项 | 要点 |
| --- | --- |
| token 存储 | 不把长期敏感 token 暴露在易被 XSS 读取的位置 |
| HTML 渲染 | 不直接渲染未净化的用户输入 |
| Cookie | 登录 Cookie 设置 `HttpOnly`、`Secure`、`SameSite` |
| 写操作 | 不用 GET 做删除、支付、审批等操作 |
| CORS | 凭证请求不能使用通配 Origin |
| CSP | 至少为生产站点设计基础策略 |
| 第三方脚本 | 控制来源，避免无限制注入 |
| 错误日志 | 不上传密码、token、身份证号等敏感数据 |

## 常见问题

### 1. 使用 localStorage 存 token 是否安全

localStorage 容易被 XSS 读取。它不是绝对不能用，但必须理解风险。

如果项目安全要求较高，优先考虑服务端 Session + `HttpOnly` Cookie，或者短 token + refresh 策略，并配合 XSS 防护。

### 2. 有了框架转义是不是就不会 XSS

不是。框架默认转义普通文本，但 `v-html`、`dangerouslySetInnerHTML`、第三方富文本、Markdown 渲染、外部脚本仍然可能引入 XSS。

### 3. CSP 会不会影响正常功能

会。CSP 本质是白名单策略，配置过严可能阻止第三方脚本、图片、样式或埋点。建议先用 Report-Only 观察，再逐步启用。

## 下一步学习

- [跨域与登录态](/browser/cors-auth)
- [HTTP 与请求流程](/browser/http-request)
- [HTTP 速查](/cheatsheets/http)
- [后端接口与服务问题](/projects/issues-backend)
