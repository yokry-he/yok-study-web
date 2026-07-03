# 内存管理

## 适合谁看

适合已经能写 Vue、React 或普通前端页面，但开始遇到这些问题的人：

- 页面用久了越来越卡。
- 弹窗、图表、大列表反复打开后内存上涨。
- 组件已经卸载，事件、定时器或请求回调还在执行。
- 图片、文件、Blob、Canvas 占用很大，不知道什么时候释放。
- 听说 JavaScript 有垃圾回收，所以以为完全不用管内存。

JavaScript 会自动进行垃圾回收，但这不等于前端开发者不用关心内存。只要代码仍然持有对象引用，垃圾回收就不能释放它。

## 内存生命周期

可以简单理解为三步：

```text
分配内存
↓
使用内存
↓
不再需要时释放
```

JavaScript 自动处理释放动作，但前提是对象不再可达。

示例：

```ts
function createUsers() {
  const users = new Array(10000).fill(null).map((_, index) => ({
    id: index,
    name: `user-${index}`
  }))

  return users
}
```

调用函数会创建数组和对象。只要外部还引用 `users`，这些对象就不能被释放。

## 什么是可达

如果一个对象还能从全局对象、当前调用栈、闭包、DOM、定时器、事件监听、缓存等路径访问到，它就是可达的。

常见引用来源：

| 来源 | 示例 |
| --- | --- |
| 全局变量 | `window.cache = data` |
| 闭包 | 回调函数保存了外部变量 |
| DOM 引用 | JS 保存了已删除节点 |
| 事件监听 | listener 引用了组件状态 |
| 定时器 | interval 持续持有回调 |
| Map/数组缓存 | 缓存只增不删 |
| 未完成请求 | 回调里引用页面状态 |

内存泄漏通常不是“垃圾回收坏了”，而是代码还留着引用。

## 常见泄漏场景

### 1. 事件监听未清理

```ts
onMounted(() => {
  window.addEventListener('resize', handleResize)
})
```

如果组件卸载不移除监听，`handleResize` 可能仍然引用组件里的状态。

修复：

```ts
onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})
```

或者使用 `AbortController` 批量清理。

### 2. 定时器未清理

```ts
const timer = setInterval(() => {
  refresh()
}, 5000)
```

修复：

```ts
onBeforeUnmount(() => {
  clearInterval(timer)
})
```

定时器不仅占内存，还会持续执行逻辑，造成重复请求和性能问题。

### 3. 大对象缓存只增不删

```ts
const detailCache = new Map<string, Detail>()

function saveDetail(id: string, detail: Detail) {
  detailCache.set(id, detail)
}
```

如果没有容量上限和淘汰策略，页面用久了缓存会越来越大。

改进：

- 限制缓存数量。
- 切换用户或路由时清理。
- 设置过期时间。
- 使用服务端分页，避免缓存全量数据。

### 4. 保存了已删除 DOM

```ts
const nodes: Element[] = []

function collectNode(node: Element) {
  nodes.push(node)
}
```

即使 DOM 从页面上移除了，只要数组还引用它，相关对象仍然不能释放。

### 5. Blob URL 没有释放

```ts
const url = URL.createObjectURL(file)
preview.src = url
```

不用后应释放：

```ts
URL.revokeObjectURL(url)
```

上传、预览、导出文件、图片裁剪工具里很常见。

## Vue / React 里的清理边界

Vue：

```ts
onMounted(() => {
  const timer = window.setInterval(loadData, 5000)

  onBeforeUnmount(() => {
    clearInterval(timer)
  })
})
```

React：

```tsx
useEffect(() => {
  const timer = window.setInterval(loadData, 5000)

  return () => {
    clearInterval(timer)
  }
}, [])
```

凡是创建了外部资源，都要问一句：组件卸载时谁负责清理？

外部资源包括：

- DOM 事件。
- 定时器。
- WebSocket。
- Observer。
- Worker。
- 图表实例。
- 地图实例。
- Blob URL。
- 第三方 SDK 实例。

## 闭包导致的引用保留

闭包本身不是问题，但闭包会保留它用到的外部变量。

```ts
function createHandler(bigData: BigData) {
  return function handleClick() {
    console.log(bigData.name)
  }
}
```

只要 `handleClick` 还在，`bigData` 就可能不能释放。

项目里要注意：

- 不要在长期存在的回调里捕获巨大对象。
- 不再需要时移除回调。
- 大数据只传必要字段。

## 异步请求和卸载

组件卸载后，请求回来还更新状态，会导致警告、无效更新或保留引用。

思路：

```ts
const controller = new AbortController()

async function loadData() {
  const response = await fetch('/api/list', {
    signal: controller.signal
  })

  return response.json()
}

onBeforeUnmount(() => {
  controller.abort()
})
```

如果请求库不支持取消，也可以用请求序号或 mounted 标记避免旧结果写入。

## 如何排查内存问题

浏览器 DevTools 常用方法：

1. 打开 Memory 面板。
2. 记录初始快照。
3. 重复执行可疑操作，例如打开关闭弹窗 10 次。
4. 手动触发 GC 或等待一段时间。
5. 再记录快照。
6. 对比对象数量和 retained size。

也可以用 Performance 面板观察长时间交互中的内存曲线。

排查方向：

- 哪类对象数量持续增加。
- 是否有 Detached DOM tree。
- 是否有大量 listener、timer、closure。
- 图表、地图、编辑器实例是否 dispose。
- 缓存 Map 是否无限增长。

## 实际项目常见问题

### 1. 弹窗反复打开后越来越卡

**原因**

每次打开都创建图表或事件监听，关闭时没有销毁。

**解决方案**

- 弹窗关闭时销毁图表实例。
- 组件卸载时 remove listener。
- 检查是否重复 setInterval。

### 2. 大列表切换筛选后内存上涨

**原因**

保存了多份大数组，旧数据仍被引用。

**解决方案**

- 只保留当前展示所需数据。
- 使用分页或虚拟滚动。
- 缓存设置上限。

### 3. 图片预览页面占用很大

**原因**

图片文件、Canvas、Blob URL 没释放。

**解决方案**

- 预览关闭时 revoke object URL。
- Canvas 不用时断开引用。
- 大图压缩放到 Worker。

### 4. WebSocket 页面离开后还在收消息

**原因**

路由离开或组件卸载时没有关闭连接。

**解决方案**

在生命周期清理：

```ts
onBeforeUnmount(() => {
  socket.close()
})
```

### 5. 全局 store 越来越大

**原因**

把页面临时数据、搜索结果、详情缓存都放进全局状态。

**解决方案**

- 全局 store 只放跨页面共享状态。
- 页面临时状态留在页面或 composable。
- 路由离开时重置不需要保留的数据。

## 最佳实践

- 创建外部资源时，同时写清理逻辑。
- 组件卸载时清理事件、定时器、Observer、Worker、Socket、图表实例。
- 缓存必须有边界：数量、时间或业务清理点。
- 大数据避免放进长期全局状态。
- Blob URL、Canvas、文件预览要明确释放。
- 用 DevTools Memory 对比快照，不靠感觉判断泄漏。

## 学习检查

学完本节后，你应该能回答：

- 为什么 JavaScript 自动 GC 不等于不用管内存。
- 什么样的对象不能被回收。
- 事件监听、定时器、闭包为什么容易造成泄漏。
- 组件卸载时应该清理哪些资源。
- 如何用 DevTools 排查内存持续上涨。

## 参考资料

- [MDN: Memory management](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Memory_management)
- [MDN: Performance.memory](https://developer.mozilla.org/en-US/docs/Web/API/Performance/memory)
- [MDN: EventTarget.removeEventListener](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener)

## 下一步学习

继续学习 [模块化与工程实践](/javascript/modules)，把事件、缓存、请求和清理逻辑放到合适的模块边界里。
