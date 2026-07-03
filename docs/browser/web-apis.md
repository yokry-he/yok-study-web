# 常用 Web API

## 适合谁看

适合已经能写业务页面，但希望理解浏览器提供的常用 API，而不是所有能力都依赖框架和第三方库的学习者。

Web API 很多，不需要一次全部学完。真实项目里优先掌握那些高频、稳定、容易带来性能或体验提升的能力。

## 学习方式

先按场景学习：

| 场景 | API |
| --- | --- |
| 元素进入视口 | IntersectionObserver |
| 大计算不阻塞页面 | Web Workers |
| 多标签页通信 | BroadcastChannel、storage event |
| 复制文本 | Clipboard API |
| 文件选择和预览 | File、FileReader、URL.createObjectURL |
| 页面可见性 | Page Visibility API |
| 网络状态 | Navigator.onLine、online/offline event |
| 性能观测 | Performance API |

不要为了“会用 API”而强行使用。先判断它是否解决了真实问题。

## IntersectionObserver

适合：

- 图片懒加载。
- 无限滚动。
- 曝光埋点。
- 目录高亮。
- 动画进入视口后触发。

示例：

```ts
const observer = new IntersectionObserver((entries) => {
  for (const entry of entries) {
    if (entry.isIntersecting) {
      console.log('进入视口', entry.target)
      observer.unobserve(entry.target)
    }
  }
})

observer.observe(document.querySelector('.target')!)
```

项目建议：

- 组件卸载时取消观察。
- 不要在回调里做重计算。
- 列表很多时复用 observer。

## Web Workers

适合把耗时计算放到独立线程，避免阻塞主线程。

适合：

- 大量数据计算。
- 文件解析。
- 图片处理。
- 复杂加密或压缩。

主线程：

```ts
const worker = new Worker(new URL('./worker.ts', import.meta.url), {
  type: 'module'
})

worker.postMessage({ items })

worker.onmessage = (event) => {
  console.log(event.data)
}
```

worker：

```ts
self.onmessage = (event) => {
  const result = heavyCalculate(event.data.items)
  self.postMessage(result)
}
```

注意：

- Worker 不能直接操作 DOM。
- 数据传递有序列化成本。
- 组件卸载时终止 worker。

## BroadcastChannel

适合多标签页通信。

示例：

```ts
const channel = new BroadcastChannel('auth')

channel.postMessage({
  type: 'logout'
})

channel.onmessage = (event) => {
  if (event.data.type === 'logout') {
    location.href = '/login'
  }
}
```

适合：

- 多标签页退出登录同步。
- 多标签页主题同步。
- 多窗口任务状态同步。

如果只需要监听 localStorage 变化，也可以使用 `storage` 事件。

## Clipboard API

复制文本：

```ts
async function copyText(text: string) {
  await navigator.clipboard.writeText(text)
}
```

注意：

- 通常需要安全上下文。
- 可能需要用户手势触发。
- 失败时要给出提示或兜底。

示例：

```ts
async function copyInviteLink(link: string) {
  try {
    await navigator.clipboard.writeText(link)
    showSuccess('复制成功')
  } catch {
    showError('复制失败，请手动复制')
  }
}
```

## File API

文件选择：

```ts
function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  console.log(file.name, file.size, file.type)
}
```

图片预览：

```ts
const url = URL.createObjectURL(file)

// 使用后释放
URL.revokeObjectURL(url)
```

注意：

- 限制文件大小。
- 校验文件类型。
- 上传 `FormData` 时不要手动设置错误的 multipart boundary。
- 预览 URL 用完要释放。

## Page Visibility API

判断页面是否可见：

```ts
document.addEventListener('visibilitychange', () => {
  if (document.visibilityState === 'visible') {
    refreshData()
  }
})
```

适合：

- 页面重新可见时刷新数据。
- 页面隐藏时暂停轮询。
- 降低后台标签页资源消耗。

## 网络状态

```ts
window.addEventListener('online', () => {
  showSuccess('网络已恢复')
})

window.addEventListener('offline', () => {
  showWarning('网络已断开')
})
```

注意：`navigator.onLine` 只能提供粗略判断，不代表业务接口一定可用。

项目里更可靠的方式是请求自己的健康检查或关键接口。

## Performance API

简单记录耗时：

```ts
performance.mark('list-start')

await fetchList()

performance.mark('list-end')
performance.measure('fetch-list', 'list-start', 'list-end')
```

读取：

```ts
const measures = performance.getEntriesByName('fetch-list')
```

适合：

- 记录关键流程耗时。
- 分析首屏和交互性能。
- 配合 Performance 面板定位长任务。

## 常见问题

### 1. API 在本地可用，线上不可用

常见原因：

- API 需要安全上下文。
- 浏览器权限被用户拒绝。
- 部分浏览器不支持。
- 被 iframe 权限策略限制。

处理：

- 检查 HTTPS。
- 检查浏览器兼容性。
- 检查权限状态。
- 做降级和错误提示。

### 2. 页面越来越卡

常见原因：

- observer 没取消。
- worker 没终止。
- object URL 没释放。
- 事件监听没移除。
- 定时器或轮询没停止。

处理：

- 组件卸载时清理资源。
- 长任务移到 worker。
- 使用 Performance 面板定位。

## 项目建议

- 使用 Web API 前先确认浏览器支持和权限要求。
- 对失败场景提供明确提示。
- 组件卸载时清理 observer、worker、timer、event listener。
- 处理用户隐私和敏感权限时保持克制。
- 需要长期兼容时查 MDN 兼容性表。

## 下一步学习

- [浏览器学习导览](/browser/introduction)
- [渲染与性能](/browser/rendering-performance)
- [浏览器存储](/browser/storage)
- [前端页面与状态问题](/projects/issues-frontend)
