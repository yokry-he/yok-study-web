# 浏览器存储

## 适合谁看

适合想搞清楚 Cookie、LocalStorage、SessionStorage、IndexedDB、Cache API 怎么选的人。

前端存储不是“哪里方便放哪里”。不同存储的生命周期、容量、同步异步、安全边界、是否会随请求发送都不同。错误选择会导致登录失效、数据泄漏、页面卡顿或缓存难以更新。

## 常见存储方式

| 存储 | 生命周期 | 容量 | 是否随请求发送 | 常见用途 |
| --- | --- | --- | --- | --- |
| Cookie | 可设置过期时间 | 小 | 是，匹配规则满足时自动发送 | 会话、服务端登录态 |
| LocalStorage | 手动清除前保留 | 较小 | 否 | 主题、语言、非敏感偏好 |
| SessionStorage | 标签页关闭后清除 | 较小 | 否 | 临时表单、单标签页状态 |
| IndexedDB | 长期保存 | 较大 | 否 | 离线数据、大结构化数据 |
| Cache API | 按缓存策略保存 | 较大 | 否 | PWA、离线资源、请求响应缓存 |
| Memory | 页面生命周期 | 取决于内存 | 否 | Pinia/Redux 状态、临时数据 |

## Cookie

Cookie 最适合和服务端会话绑定。浏览器会根据 Domain、Path、SameSite、Secure 等规则决定是否发送。

服务端设置 Cookie 示例：

```http
Set-Cookie: session_id=abc123; HttpOnly; Secure; SameSite=Lax; Path=/
```

关键点：

- `HttpOnly` Cookie 不能被 JavaScript 读取，更适合保存会话标识。
- `Secure` 表示只在 HTTPS 下发送。
- `SameSite` 影响跨站请求是否携带。
- Cookie 会自动随请求发送，过多 Cookie 会增加请求体积。

## LocalStorage

LocalStorage API 简单：

```ts
localStorage.setItem('theme', 'light')
const theme = localStorage.getItem('theme')
localStorage.removeItem('theme')
```

适合：

- 主题。
- 语言。
- 表格密度。
- 用户非敏感偏好。

不适合：

- 密码。
- 长期有效 token。
- 大量列表数据。
- 高频写入数据。

LocalStorage 是同步 API。大量读写会阻塞主线程。

## SessionStorage

SessionStorage 和 LocalStorage 类似，但生命周期通常限制在当前标签页。

适合：

- 多步骤表单临时草稿。
- 从列表进入详情后保存返回位置。
- 当前标签页内的临时筛选状态。

不适合跨标签页共享状态。

## IndexedDB

IndexedDB 是浏览器里的结构化数据库，API 是异步的，适合保存较大、较复杂的数据。

适合：

- 离线应用数据。
- 大列表缓存。
- 编辑器草稿。
- 客户端搜索索引。
- 文件元信息。

不适合为了保存一个主题配置就引入。简单偏好用 LocalStorage 更直接。

## Cache API

Cache API 常和 Service Worker 一起使用，用来缓存请求和响应。

适合：

- PWA 离线资源。
- 静态资源缓存。
- 特定接口响应缓存。

如果项目没有离线需求，不要一开始就引入 Service Worker。它会接管请求，调试和更新策略都更复杂。

## 登录态怎么存

没有绝对统一答案，按业务风险选择。

### Cookie Session

适合：

- 后台管理系统。
- 同域部署。
- 服务端希望统一管理会话。

建议：

- 使用 `HttpOnly`。
- 使用 HTTPS 和 `Secure`。
- 设置合适的 `SameSite`。
- 后端做 CSRF 防护。

### Token

适合：

- 多端共用 API。
- 移动端、小程序、第三方客户端。
- 前端明确控制请求头。

建议：

- token 有合理过期时间。
- refresh token 策略清晰。
- 尽量降低 XSS 风险。
- 敏感权限以服务端校验为准。

## 实际项目问题

### 问题：刷新页面后 Pinia 状态丢失

**原因**

Pinia 默认状态在内存里，刷新页面会重新加载 JavaScript，内存状态自然丢失。

**解决方案**

需要持久化的状态写入合适存储：

```ts
import { watch } from 'vue'
import { useUserStore } from './user'

const userStore = useUserStore()

watch(
  () => userStore.preferences,
  preferences => {
    localStorage.setItem('preferences', JSON.stringify(preferences))
  },
  { deep: true }
)
```

登录态是否持久化要按安全策略决定，不要默认把所有 store 都存进 LocalStorage。

### 问题：多个标签页登录状态不同步

**原因**

内存状态不能跨标签页。LocalStorage 更新可以触发其他标签页的 `storage` 事件，但当前标签页不会触发自己的事件。

**解决方案**

```ts
window.addEventListener('storage', event => {
  if (event.key === 'token' && !event.newValue) {
    // 其他标签页退出登录，当前标签页同步退出
    location.href = '/login'
  }
})
```

### 问题：LocalStorage 里存了对象，读取后类型不对

**原因**

LocalStorage 只能保存字符串。

**解决方案**

封装读写：

```ts
export function setStorage<T>(key: string, value: T) {
  localStorage.setItem(key, JSON.stringify(value))
}

export function getStorage<T>(key: string, fallback: T): T {
  const raw = localStorage.getItem(key)

  if (!raw) return fallback

  try {
    return JSON.parse(raw) as T
  } catch {
    return fallback
  }
}
```

## 最佳实践

- 敏感会话优先考虑 `HttpOnly` Cookie。
- LocalStorage 只放非敏感、小体积、低频读写的数据。
- 大体积结构化数据用 IndexedDB。
- 离线资源缓存再考虑 Cache API 和 Service Worker。
- 所有存储读取都要考虑解析失败、空值和版本迁移。
- 存储不是权限来源，关键权限必须由服务端校验。

## 参考资料

- [MDN: Web Storage API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)
- [MDN: Storage Access API](https://developer.mozilla.org/en-US/docs/Web/API/Storage_Access_API)
- [MDN: Set-Cookie header](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie)

## 下一步学习

继续学习 [渲染与性能](/browser/rendering-performance)。
