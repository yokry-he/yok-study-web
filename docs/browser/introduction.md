# 浏览器学习导览

## 适合谁看

适合已经能写 Vue、React 或普通前端页面，但遇到这些问题时还不够有把握的人：

- 接口在 Postman 能通，浏览器里不通。
- 本地开发正常，部署后跨域、刷新 404、缓存旧文件。
- 登录态一会儿有效一会儿失效。
- 页面首屏慢、滚动卡、接口重复请求。
- 不知道 Cookie、LocalStorage、IndexedDB、Cache API 应该怎么选。

前端工程师最终写出来的代码都运行在浏览器里。框架负责组织页面，浏览器负责请求资源、执行 JavaScript、管理安全边界、缓存文件、渲染像素和存储数据。很多线上问题不是 Vue 或 React 的问题，而是浏览器规则没有理解清楚。

## 学习路线

```text
URL 与 HTTP 请求
↓
图解浏览器核心概念
↓
浏览器与网络从零到项目落地
↓
跨域与登录态
↓
HTTP 缓存与发布缓存
↓
Cookie、Web Storage、IndexedDB
↓
浏览器安全、CSP、XSS、CSRF
↓
Service Worker 与 PWA
↓
常用 Web API
↓
WebSocket、WebRTC、Web Components
↓
WebAssembly、WebGPU
↓
DOM、CSSOM、渲染流水线
↓
Network、Application、Performance、自动化调试
↓
真实项目问题定位
```

## 你需要先掌握的概念

### URL

```text
https://example.com:443/admin/users?page=1#detail
```

| 部分 | 含义 | 项目中常见影响 |
| --- | --- | --- |
| `https` | 协议 | 是否安全上下文、Cookie 是否能带 `Secure` |
| `example.com` | 域名 | 是否同源、Cookie Domain 是否匹配 |
| `443` | 端口 | 同源判断的一部分 |
| `/admin/users` | 路径 | 前端路由、服务端 fallback、接口路径 |
| `page=1` | 查询参数 | 列表筛选、分页、分享链接 |
| `#detail` | hash | hash 路由、页面内定位 |

同源要求协议、域名、端口都相同。只要其中一个不同，浏览器就会按跨源场景处理。

### 浏览器不是简单的 HTTP 客户端

Postman、curl 和浏览器都能发请求，但浏览器多了安全策略：

- 会执行同源策略。
- 会根据 CORS 响应头决定是否允许前端读取响应。
- 会自动管理 Cookie，但是否发送 Cookie 受 SameSite、Secure、Domain、Path、credentials 等条件共同影响。
- 会使用 HTTP 缓存、内存缓存、Service Worker 缓存。
- 会阻止部分不安全能力，例如 HTTPS 页面请求 HTTP 资源。

因此“接口本身可用”不等于“浏览器前端可以正常调用”。

## 模块章节

| 章节 | 解决的问题 |
| --- | --- |
| [图解浏览器核心概念](/browser/visual-guide) | 用图理解 URL 到渲染、请求跨域、Cookie、缓存、存储、渲染和排错路径 |
| [浏览器与网络从零到项目落地](/browser/project-from-zero) | 用请求诊断工作台串起 Fetch、CORS、Cookie、缓存、存储、Service Worker 和 DevTools 证据链 |
| [HTTP 与请求流程](/browser/http-request) | URL、请求方法、状态码、请求头、响应头、Network 面板 |
| [跨域与登录态](/browser/cors-auth) | CORS、预检请求、Cookie、token、401/403 |
| [缓存策略](/browser/cache) | 强缓存、协商缓存、CDN 缓存、前端发布缓存 |
| [浏览器存储](/browser/storage) | Cookie、LocalStorage、SessionStorage、IndexedDB、Cache API |
| [浏览器安全基础](/browser/security) | 同源策略、XSS、CSRF、CSP、Cookie 安全属性 |
| [Service Worker 与 PWA](/browser/service-worker-pwa) | 离线缓存、Service Worker 生命周期、PWA 更新策略 |
| [常用 Web API](/browser/web-apis) | IntersectionObserver、Web Workers、BroadcastChannel、Clipboard、File API |
| [WebSocket 实时通信](/browser/websocket) | 双向连接、消息协议、鉴权、心跳、重连 |
| [WebRTC 实时音视频](/browser/webrtc) | 摄像头、麦克风、点对点连接、信令、TURN/STUN |
| [Web Components](/browser/web-components) | Custom Elements、Shadow DOM、template、slot、跨框架组件 |
| [WebAssembly](/browser/webassembly) | 浏览器高性能计算、跨语言复用、Worker 搭配和部署排错 |
| [WebGPU](/browser/webgpu) | GPU 渲染、并行计算、能力检测、降级和生命周期管理 |
| [浏览器自动化调试](/browser/browser-automation-debugging) | Playwright、CDP、白屏检查、路由巡检、移动端溢出检查 |
| [渲染与性能](/browser/rendering-performance) | DOM、CSSOM、布局、绘制、卡顿和性能排查 |
| [常见问题](/browser/troubleshooting) | 真实项目里高频浏览器问题和处理方案 |

## 实际项目最常用的判断

### 接口问题先看 Network

不要先猜 store、路由或组件。先打开 DevTools 的 Network：

1. 请求有没有发出去。
2. URL 是否正确。
3. Method 是否正确。
4. Status Code 是多少。
5. Request Headers 是否带了 token 或 Cookie。
6. Response Headers 是否有 CORS、缓存、Set-Cookie。
7. Response Body 是否是业务错误。

很多问题到第 5 步就能定位。

### 登录态问题同时看前端和响应头

登录态一般有两类方案：

| 方案 | 特点 | 常见风险 |
| --- | --- | --- |
| Cookie Session | 服务端通过 Cookie 识别用户 | SameSite、Domain、Secure、跨域 credentials 配错 |
| Token | 前端把 token 放到请求头 | token 泄漏、刷新丢失、过期刷新策略缺失 |

不要只看前端是否保存了 token。还要看请求是否真的带出去了，后端是否接受，响应是否返回新的登录状态。

### 缓存问题先分清缓存对象

用户看到旧页面时，不要只让用户清缓存。先判断是哪一层缓存：

| 缓存对象 | 常见位置 | 处理思路 |
| --- | --- | --- |
| `index.html` | 浏览器、CDN、网关 | 不建议强缓存，发布后要能及时更新 |
| `app.[hash].js` | 浏览器、CDN | 文件名带 hash 后可长期缓存 |
| 接口响应 | 浏览器、代理、服务端 | 根据业务设置 `Cache-Control` |
| Service Worker | 浏览器 Application 面板 | 检查更新策略和缓存版本 |

## 学习建议

浏览器模块不要死记概念。建议每学完一节，都在 DevTools 里验证一次：

- Network 看请求和缓存。
- Application 看 Cookie、Storage、Cache。
- Performance 录制一次页面交互。
- Lighthouse 或浏览器性能面板看首屏和阻塞资源。

## 参考资料

- [MDN: Cross-Origin Resource Sharing](/browser/cors-auth#参考资料)
- [MDN: Cache-Control](/browser/cache#参考资料)
- [MDN: Web Storage API](/browser/storage#参考资料)
- [MDN: Critical rendering path](/browser/rendering-performance#参考资料)
- [MDN: WebAssembly](/browser/webassembly#参考资料)
- [MDN: WebGPU API](/browser/webgpu#参考资料)
- [Playwright](/browser/browser-automation-debugging#参考资料)

## 下一步学习

继续学习 [图解浏览器核心概念](/browser/visual-guide)，再进入 [浏览器与网络从零到项目落地](/browser/project-from-zero)，用项目把 [HTTP 与请求流程](/browser/http-request)、[跨域与登录态](/browser/cors-auth)、[缓存策略](/browser/cache) 和 [浏览器存储](/browser/storage) 串起来。
