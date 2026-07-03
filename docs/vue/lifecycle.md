# 生命周期

## 适合谁看

适合已经会写组件，但不清楚“什么时候请求数据”“什么时候能访问 DOM”“什么时候清理定时器和监听器”的学习者。

Vue 官方文档说明，每个组件实例在创建时会经历初始化、挂载、更新和卸载等步骤，并在这些阶段运行生命周期钩子。生命周期就是 Vue 给你的“时机入口”。

## 你会学到什么

- 常用生命周期钩子分别什么时候执行。
- 接口请求、DOM 操作、定时器、事件监听应该放哪里。
- `nextTick` 什么时候用。
- KeepAlive 下生命周期有什么不同。
- 实际项目里重复请求、内存泄漏、DOM 获取不到怎么排查。

## 常用生命周期

组合式 API 中常用：

| 钩子 | 触发时机 | 常见用途 |
| --- | --- | --- |
| `onMounted` | 组件挂载到 DOM 后 | 首次请求、访问 DOM |
| `onUpdated` | 响应式更新导致 DOM 更新后 | 少量 DOM 依赖逻辑 |
| `onUnmounted` | 组件卸载后 | 清理定时器、监听器 |
| `onActivated` | KeepAlive 缓存组件重新激活 | 刷新缓存页数据 |
| `onDeactivated` | KeepAlive 缓存组件离开 | 暂停轮询、暂停监听 |

## 首次请求数据

```ts
onMounted(async () => {
  await fetchList()
})
```

如果请求依赖路由参数，并且参数变化时也要重新请求，更推荐 `watch`：

```ts
const route = useRoute()

watch(
  () => route.params.id,
  (id) => {
    fetchDetail(String(id))
  },
  { immediate: true }
)
```

## 访问 DOM

```vue
<script setup lang="ts">
import { onMounted, ref } from 'vue'

const inputRef = ref<HTMLInputElement | null>(null)

onMounted(() => {
  inputRef.value?.focus()
})
</script>

<template>
  <input ref="inputRef" />
</template>
```

如果数据变化后要等待 DOM 更新：

```ts
items.value = await fetchItems()
await nextTick()

listRef.value?.scrollTo({ top: 0 })
```

## 清理副作用

定时器必须清理：

```ts
let timer: number | undefined

onMounted(() => {
  timer = window.setInterval(() => {
    fetchUnreadCount()
  }, 10000)
})

onUnmounted(() => {
  window.clearInterval(timer)
})
```

事件监听也要清理：

```ts
function handleResize() {
  width.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
```

## KeepAlive 场景

被 `<KeepAlive>` 缓存的组件离开时不会卸载，因此 `onUnmounted` 不会马上执行。

```ts
onActivated(() => {
  fetchList()
})

onDeactivated(() => {
  stopPolling()
})
```

适合缓存的页面：

- 列表页返回后保留筛选和滚动位置。
- Tab 页面频繁切换。

不适合缓存的页面：

- 高实时性详情页。
- 权限或数据变化频繁的页面。

## 实际项目常见问题

### 1. 页面打开请求了两次

**常见原因**

- `onMounted(fetchList)` 调了一次。
- `watch(query, fetchList, { immediate: true })` 又调了一次。

**解决方案**

保留一个入口。依赖路由或筛选条件的请求，优先用 `watch` 统一处理。

### 2. onMounted 中获取不到元素

**原因**

元素被 `v-if` 控制，当前还没有渲染；或者数据更新后 DOM 还没刷新。

**解决方案**

```ts
visible.value = true
await nextTick()
inputRef.value?.focus()
```

### 3. 页面离开后还在请求接口

**原因**

定时器、轮询、事件监听没有清理。

**解决方案**

在 `onUnmounted` 或 `onDeactivated` 中清理。

### 4. KeepAlive 页面数据不刷新

**原因**

缓存组件再次显示时不会重新执行 `onMounted`。

**解决方案**

使用 `onActivated`。

## 最佳实践

- 首次请求可以放 `onMounted`，依赖参数变化的请求用 `watch`。
- DOM 更新后再操作元素时使用 `nextTick`。
- 所有定时器、事件监听、订阅都要清理。
- KeepAlive 页面用 `onActivated` 和 `onDeactivated` 管理恢复与暂停。
- 不要把大量业务逻辑都塞进生命周期钩子，复杂流程应拆成函数。

## 下一步学习

继续学习 [路由与页面](/vue/router)。
