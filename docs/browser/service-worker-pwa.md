# Service Worker 与 PWA

## 适合谁看

适合已经理解浏览器缓存和普通前端部署，但想进一步学习离线访问、资源预缓存、更新提示、PWA 安装和 Service Worker 调试的学习者。

Service Worker 很强，但也很容易制造“用户一直看到旧页面”的问题。学习它时要同时理解能力和风险。

## Service Worker 是什么

Service Worker 是运行在页面之外的浏览器工作线程。它可以拦截网络请求、管理缓存、处理离线访问、推送通知和后台同步等能力。

它和普通页面 JavaScript 的区别：

| 能力 | 页面脚本 | Service Worker |
| --- | --- | --- |
| 操作 DOM | 可以 | 不可以 |
| 拦截请求 | 不可以 | 可以 |
| 管理 Cache API | 可以 | 可以 |
| 页面关闭后运行 | 通常不可以 | 可被事件唤起 |
| 需要 HTTPS | 线上通常需要 | 需要安全上下文 |

Service Worker 不是万能缓存插件。它更像浏览器和网络之间的一层可编程代理。

## 生命周期

核心阶段：

```text
register
↓
install
↓
waiting
↓
activate
↓
fetch
```

含义：

| 阶段 | 说明 |
| --- | --- |
| register | 页面注册 Service Worker |
| install | 安装阶段，常用于预缓存资源 |
| waiting | 新版本等待旧页面关闭 |
| activate | 激活阶段，常用于清理旧缓存 |
| fetch | 拦截请求并决定走网络还是缓存 |

很多更新问题来自 `waiting` 阶段：新 Service Worker 已经下载，但旧页面还在使用旧版本。

## 注册示例

页面中注册：

```ts
if ('serviceWorker' in navigator) {
  window.addEventListener('load', async () => {
    try {
      await navigator.serviceWorker.register('/sw.js')
    } catch (error) {
      console.error('Service Worker 注册失败', error)
    }
  })
}
```

注意：

- 生产环境需要 HTTPS。
- sw 文件路径会影响控制范围。
- 注册失败要能在控制台看到错误。

## 缓存策略

常见策略：

| 策略 | 适合场景 | 风险 |
| --- | --- | --- |
| Cache First | hash 静态资源、字体、图片 | 内容更新不及时 |
| Network First | API、动态内容 | 离线体验依赖缓存兜底 |
| Stale While Revalidate | 内容可先旧后新 | 用户可能短暂看到旧数据 |
| Network Only | 敏感接口、实时接口 | 离线不可用 |
| Cache Only | 构建时确定的离线资源 | 版本更新要谨慎 |

不要给所有请求使用同一种缓存策略。

## 预缓存和运行时缓存

预缓存适合构建时确定的资源：

```text
app.hash.js
style.hash.css
logo.svg
offline.html
```

运行时缓存适合访问过程中产生的资源：

```text
图片
字体
部分 GET API
内容页
```

敏感接口、用户个人数据、权限数据不建议随意缓存。

## 更新策略

PWA 最大的真实项目风险是更新不可控。

常见现象：

- 发布新版本后用户仍看到旧页面。
- 旧 Service Worker 返回旧资源。
- 新旧资源混用导致白屏。

建议：

- 给缓存命名加版本号。
- 激活时清理旧缓存。
- 检测到新版本后提示用户刷新。
- 发布时保留旧静态资源一段时间。
- 不缓存 `index.html` 或谨慎缓存入口文件。

示意：

```ts
const CACHE_NAME = 'app-cache-v2'

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys
          .filter((key) => key !== CACHE_NAME)
          .map((key) => caches.delete(key))
      )
    })
  )
})
```

## PWA Manifest

PWA 通常还需要 manifest：

```json
{
  "name": "Docs App",
  "short_name": "Docs",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#7edfc6",
  "icons": [
    {
      "src": "/icons/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    }
  ]
}
```

它影响：

- 安装名称。
- 启动地址。
- 图标。
- 独立窗口显示方式。
- 主题色。

## DevTools 调试

打开 Chrome DevTools：

```text
Application
├─ Service Workers
├─ Cache Storage
├─ Manifest
└─ Storage
```

常用检查：

- Service Worker 是否注册成功。
- 当前控制页面的是哪个版本。
- 是否处于 waiting。
- Cache Storage 里有哪些缓存。
- 勾选 Update on reload 测试更新。
- 清理站点数据后是否恢复正常。

## 实际项目常见问题

### 1. 发布后用户一直看到旧版本

常见原因：

- Service Worker 缓存了旧入口。
- 新 Service Worker 处于 waiting。
- CDN 和 Service Worker 双重缓存。

处理：

- 检查 Application 面板。
- 清理旧 cache。
- 做新版本提示。
- 发布时保留旧资源并控制入口缓存。

### 2. 离线页面显示不完整

常见原因：

- 只缓存了 HTML，没缓存 CSS、JS、字体或图片。
- 动态 API 没有离线兜底。
- 路由 fallback 没处理。

处理：

- 明确离线页面范围。
- 对离线资源做预缓存。
- 动态数据给出离线状态提示。

### 3. 登录后数据被缓存给了其他用户

这是严重问题。用户私有数据不应该进入公共缓存。

处理：

- 不缓存带鉴权的私有 API。
- 缓存 key 必须考虑用户上下文。
- 对权限、个人中心、订单等页面谨慎使用运行时缓存。

## 项目建议

- 没有离线需求时，不要急着引入 Service Worker。
- 如果引入 PWA，必须设计更新策略。
- 不缓存用户私有数据和权限数据。
- 发布流程里加入 Service Worker 检查。
- 遇到“旧页面”问题时同时检查 CDN、浏览器缓存和 Service Worker。

## 下一步学习

- [缓存策略](/browser/cache)
- [浏览器存储](/browser/storage)
- [部署、缓存与 DevOps 问题](/projects/issues-deployment)
- [构建与部署](/engineering/build-deploy)
