# 缓存策略

## 适合谁看

适合遇到这些问题的人：

- 发布后用户仍然看到旧页面。
- JavaScript 文件已经更新，但浏览器还加载旧版本。
- 接口数据不是最新的。
- CDN 刷新后仍然不生效。
- 不知道 `Cache-Control`、`ETag`、`Last-Modified` 怎么配。

缓存能显著提升加载速度，但配置错误会直接造成线上事故。学习缓存的目标不是记住所有 Header，而是能根据资源类型设计策略。

## 缓存分几层

前端项目常见缓存层：

| 层级 | 位置 | 常见影响 |
| --- | --- | --- |
| 浏览器内存缓存 | 当前页面会话 | 刷新、跳转时复用资源 |
| 浏览器磁盘缓存 | 本地磁盘 | 下次打开仍可复用 |
| Service Worker | 浏览器 Application | PWA、离线缓存、资源接管 |
| CDN 缓存 | 边缘节点 | 静态资源、HTML、图片 |
| 网关缓存 | Nginx、API Gateway | 接口或页面缓存 |
| 服务端缓存 | Redis、本地缓存 | 接口数据缓存 |

用户看到旧内容时，必须先判断是哪一层缓存。

在 Network 面板里先观察 `Size` 或 `Transferred` 列：`memory cache` 表示当前浏览器进程直接复用内存资源，`disk cache` 表示从本地磁盘读取，`304` 表示浏览器已经向服务器验证过资源没有变化。

<DocFigure
  src="/images/browser/cache-memory-disk.webp"
  alt="浏览器网络资源列表对比 memory cache、disk cache、304 协商缓存和正常网络响应"
  caption="状态码只是证据的一部分；资源来源、Cache-Control 和带 hash 的文件名要一起判断。"
  :width="1440"
  :height="900"
/>

不依赖图片的读取路径：Network → 勾选 Disable cache 前后各刷新一次 → 记录 Status、Size/Transferred、Cache-Control、ETag 和文件名 hash；若请求未出现，再检查 Service Worker 和 Cache Storage。

## 强缓存

强缓存表示浏览器在缓存有效期内可以直接复用本地资源，不必向服务器确认。

常见响应头：

```http
Cache-Control: max-age=31536000
```

含义是资源在指定秒数内保持新鲜。

适合长期缓存的资源：

- `app.8f3a1c.js`
- `style.71d2c.css`
- `logo.3a9f2.svg`

前提是文件名带 hash。内容变化后文件名变化，旧缓存不会影响新版本。

## 协商缓存

协商缓存表示缓存过期后，浏览器向服务器确认资源是否变化。如果没有变化，服务器返回 `304 Not Modified`，浏览器继续使用本地缓存。

常见响应头：

```http
ETag: "abc123"
Last-Modified: Wed, 01 Jul 2026 10:00:00 GMT
```

后续请求可能带：

```http
If-None-Match: "abc123"
If-Modified-Since: Wed, 01 Jul 2026 10:00:00 GMT
```

如果资源没变，服务器返回 304，响应体为空，节省传输成本。

## 前端发布推荐策略

Vite、Webpack 等构建工具通常会给静态资源加 hash：

```text
dist/
  index.html
  assets/
    index.8f3a1c.js
    index.71d2c.css
```

推荐策略：

| 资源 | 缓存策略 | 原因 |
| --- | --- | --- |
| `index.html` | 不强缓存或短缓存 | 它引用最新的 js/css 文件名 |
| `assets/*.js` | 长缓存 | 文件名带 hash，内容变文件名变 |
| `assets/*.css` | 长缓存 | 同上 |
| 图片、字体 | 长缓存 | 建议文件名带 hash |
| 接口响应 | 按业务配置 | 列表、权限、用户信息不能随便缓存 |

Nginx 示例：

```nginx
location = /index.html {
  add_header Cache-Control "no-cache";
}

location /assets/ {
  add_header Cache-Control "public, max-age=31536000, immutable";
}
```

`no-cache` 不是“不缓存”，而是使用前需要重新验证。真正不存储是 `no-store`。

## 接口缓存

接口缓存要看业务性质。

| 接口 | 是否适合缓存 | 说明 |
| --- | --- | --- |
| 当前用户信息 | 谨慎 | 权限变化要及时生效 |
| 权限菜单 | 谨慎 | 后台系统变更敏感 |
| 字典配置 | 适合短缓存 | 变化频率低 |
| 商品详情 | 可缓存 | 需考虑库存、价格实时性 |
| 报表数据 | 可缓存 | 可以按时间窗口缓存 |

如果接口返回用户私有数据，响应头应避免被共享缓存误缓存：

```http
Cache-Control: private, no-cache
```

敏感数据可使用：

```http
Cache-Control: no-store
```

## CDN 缓存

CDN 会在边缘节点缓存资源。它能提升访问速度，但也可能让发布回滚和刷新变复杂。

发布时要确认：

1. `index.html` 是否被 CDN 长缓存。
2. hash 静态资源是否可以长期缓存。
3. CDN 刷新路径是否覆盖入口 HTML。
4. 回滚时旧资源是否仍然存在。

不要发布新 `index.html` 后立刻删除旧 hash 资源。部分用户可能还打开着旧页面，或者 CDN 节点尚未更新，删除旧资源会导致白屏。

## 实际项目问题

### 问题：用户反馈还是旧版本

**排查顺序**

1. Network 勾选 Disable cache 后刷新，看是否变新。
2. 看 `index.html` 的 Response Headers。
3. 看 `index.html` 是从浏览器缓存、CDN 还是服务器返回。
4. 看 HTML 引用的 js/css 文件名是否是新 hash。
5. 检查 CDN 是否刷新了入口 HTML。

**解决方案**

- `index.html` 设置 `no-cache` 或短缓存。
- 静态资源使用内容 hash。
- 发布脚本刷新 CDN 的 HTML 和必要入口。
- 保留最近几个版本的静态资源，避免旧 HTML 引用文件 404。

### 问题：接口数据更新了，页面还是旧数据

**可能原因**

- 浏览器缓存了 GET 响应。
- 代理层缓存了接口。
- 前端请求库或状态库复用了旧数据。
- 后端服务缓存未失效。

**解决方案**

先确定缓存层。不要一上来给 URL 加随机数。随机数能绕开缓存，但会让缓存完全失效，也掩盖真实问题。

### 问题：刷新后白屏，控制台报旧 js 404

**原因**

用户手上有旧的 `index.html`，它引用了旧 hash 的 js。但服务器发布时删除了旧 js。

**解决方案**

- 发布时不要立即删除旧静态资源。
- 保留最近几个构建版本的 `assets`。
- `index.html` 不长缓存。
- 如果使用 CDN，刷新入口 HTML 后再观察流量。

## 最佳实践

- HTML 短缓存或协商缓存，hash 资源长缓存。
- 不要让 `index.html` 和 js/css 使用同一缓存策略。
- 发布系统要考虑旧版本资源保留。
- 接口缓存要按业务敏感度设计。
- 排查旧版本问题时先找缓存层，不要只让用户清缓存。

## 参考资料

- [MDN: Cache-Control](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cache-Control)
- [MDN: Last-Modified](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Last-Modified)
- [web.dev: Prevent unnecessary network requests with the HTTP Cache](https://web.dev/articles/http-cache)

## 下一步学习

继续学习 [浏览器存储](/browser/storage)。
