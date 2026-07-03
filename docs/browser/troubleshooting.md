# 浏览器与网络常见问题

## 使用方式

遇到线上问题时，不要直接按错误文案搜索复制答案。先按下面顺序定位：

1. 是请求问题、缓存问题、存储问题还是渲染问题。
2. 能否在本地复现。
3. DevTools 哪个面板能看到证据。
4. 问题发生在浏览器、网关、后端还是前端业务代码。
5. 修复后如何防止再次发生。

## 1. 本地正常，线上跨域

**现象**

本地开发接口正常，部署后浏览器控制台出现 CORS 错误。

**原因**

本地用了 Vite proxy，浏览器看到的是同源请求；线上前端域名和接口域名不同，进入真实跨源场景。

**解决方案**

优先做同域反向代理：

```text
https://admin.example.com/api -> http://api-service
```

如果必须跨域，后端需要正确返回 CORS 响应头，并处理 `OPTIONS` 预检请求。

## 2. 登录接口返回 Set-Cookie，但后续接口没带 Cookie

**排查**

1. Application 面板是否保存了 Cookie。
2. Cookie 是否被浏览器拒绝。
3. 后续接口 Request Headers 是否有 Cookie。
4. `SameSite`、`Secure`、`Domain`、`Path` 是否匹配。
5. 前端是否开启 `credentials: 'include'` 或 `withCredentials`。

**解决方案**

跨源且需要 Cookie 时，通常需要：

```ts
fetch('https://api.example.com/user', {
  credentials: 'include'
})
```

后端：

```http
Access-Control-Allow-Origin: https://admin.example.com
Access-Control-Allow-Credentials: true
Set-Cookie: session_id=abc; HttpOnly; Secure; SameSite=None; Path=/
```

## 3. 刷新页面 404

**现象**

点击菜单进入页面正常，刷新当前页面后服务器返回 404。

**原因**

Vue Router 或 React Router 使用 history 模式，浏览器刷新时会请求真实路径，服务器没有回退到 `index.html`。

**解决方案**

Nginx：

```nginx
location / {
  try_files $uri $uri/ /index.html;
}
```

注意接口路径不要被错误回退到前端 HTML。

## 4. 发布后部分用户白屏

**常见原因**

- 用户缓存了旧 `index.html`。
- 旧 HTML 引用的 hash js 已被删除。
- CDN 部分节点还没更新。
- Service Worker 返回了旧资源。

**排查**

1. 控制台是否有 js 404。
2. Network 看 `index.html` 是否来自缓存。
3. 查看 HTML 里引用的 js 文件名。
4. 检查服务器或 CDN 是否仍保留旧资源。
5. Application 面板检查 Service Worker。

**解决方案**

- `index.html` 不长缓存。
- hash 静态资源长缓存。
- 发布时保留旧版本静态资源。
- Service Worker 做版本更新提示。

## 5. 接口 401 死循环跳登录

**原因**

请求拦截器遇到 401 后跳登录，但登录页或刷新 token 接口也触发 401，造成循环。

**解决方案**

给认证相关接口设置白名单：

```ts
const authFreeUrls = ['/api/login', '/api/refresh-token']

function shouldRedirectLogin(url: string, status: number) {
  return status === 401 && !authFreeUrls.some(item => url.includes(item))
}
```

同时要避免多个接口同时 401 时重复弹窗或重复跳转。

## 6. 页面看起来更新了，但接口还是旧数据

**原因**

可能是接口缓存、前端状态缓存、请求库缓存或后端缓存。

**排查**

- Network 看接口是否真的重新请求。
- 看响应头是否有 `Cache-Control`、`ETag`。
- 看 Pinia/Redux/TanStack Query 是否复用缓存。
- 后端是否使用 Redis 或本地缓存。

**解决方案**

把缓存策略写清楚。不要为了省事给所有接口加时间戳参数。

## 7. 上传文件后端收不到

**常见原因**

前端手动设置了错误的 `Content-Type`：

```ts
await fetch('/api/upload', {
  method: 'POST',
  headers: {
    'Content-Type': 'multipart/form-data'
  },
  body: formData
})
```

这样会缺少 boundary。

**解决方案**

上传 `FormData` 时让浏览器自动设置：

```ts
await fetch('/api/upload', {
  method: 'POST',
  body: formData
})
```

## 8. 多标签页退出登录不同步

**原因**

一个标签页清了 token，另一个标签页内存里的用户状态还在。

**解决方案**

使用 `storage` 事件同步：

```ts
window.addEventListener('storage', event => {
  if (event.key === 'token' && !event.newValue) {
    location.href = '/login'
  }
})
```

如果使用 Cookie Session，可以在关键接口返回 401 时同步清理前端状态。

## 9. 滚动页面卡顿

**排查**

- Performance 面板录制滚动。
- 看是否有长任务。
- 看是否频繁 Layout。
- 看滚动事件里是否做了复杂计算。
- 看列表 DOM 是否过多。

**解决方案**

- 大列表虚拟滚动。
- 滚动监听节流。
- 避免滚动时同步读写布局。
- 图片懒加载。

## 10. 移动端出现横向滚动条

**常见原因**

- 固定宽度超过屏幕。
- 表格、代码块、长单词未处理。
- 绝对定位元素超出视口。
- 图片没有 `max-width: 100%`。

**排查**

在控制台运行：

```js
[...document.querySelectorAll('body *')]
  .filter(el => el.getBoundingClientRect().right > document.documentElement.clientWidth)
  .map(el => ({
    tag: el.tagName,
    className: el.className,
    width: el.getBoundingClientRect().width
  }))
```

**解决方案**

给业务容器设置明确的响应式约束，不要用宽泛选择器强行覆盖组件库内部 DOM。

## 排查清单

| 问题 | 优先看哪里 |
| --- | --- |
| 接口失败 | Network |
| 登录态异常 | Network + Application |
| 旧版本 | Network + CDN 配置 |
| 刷新 404 | 服务器路由 fallback |
| 页面卡顿 | Performance |
| 本地正常线上异常 | 环境变量、代理、域名、缓存 |

## 下一步学习

继续进入 [Vue 请求与接口封装](/vue/request)，把浏览器请求规则落到项目代码里。
