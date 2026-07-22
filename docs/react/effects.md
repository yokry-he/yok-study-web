# Effect 与副作用

## 适合谁看

适合已经会写组件状态，但不清楚什么时候该用 Effect、什么时候不需要 Effect，以及如何清理订阅、定时器和请求的人。

## useEffect 解决什么

React 官方文档说明，`useEffect` 用来让组件和外部系统同步，例如网络、浏览器 API、订阅、定时器、第三方库。

不要把所有数据计算都放进 Effect。很多时候你并不需要 Effect。

## 基础写法

```tsx
useEffect(() => {
  document.title = `搜索：${keyword}`
}, [keyword])
```

依赖数组中的值变化时，Effect 会重新执行。

## 请求数据

搜索条件快速变化时，旧请求可能比新请求更晚返回。下面的时间线展示 cleanup 取消旧请求，并只允许最新响应写入当前页面。

<DocFigure
  src="/images/react/effect-request-race.webp"
  alt="React 查询从 v 变为 vue 后清理旧 Effect，取消第一个请求并只接收第二个响应"
  caption="Effect cleanup 会在依赖变化前执行上一轮清理，不只在组件卸载时执行。"
  :width="1440"
  :height="900"
/>

无法真正取消的异步操作也要使用请求序号或失效标记忽略旧结果，避免“后返回的数据覆盖当前条件”。

```tsx
useEffect(() => {
  let ignore = false

  async function fetchUsers() {
    const result = await userApi.getList()
    if (!ignore) {
      setUsers(result.items)
    }
  }

  fetchUsers()

  return () => {
    ignore = true
  }
}, [])
```

清理函数用于避免组件卸载后继续设置状态。

## 定时器清理

```tsx
useEffect(() => {
  const timer = window.setInterval(() => {
    fetchUnreadCount()
  }, 10000)

  return () => {
    window.clearInterval(timer)
  }
}, [])
```

## 实际项目常见问题

### 1. Effect 无限循环

**原因**

Effect 中更新了依赖数组里的状态。

**解决方案**

重新审视依赖和数据流。能在事件里做的，不要放进 Effect。

### 2. 依赖数组缺值

**问题**

Effect 使用了某个变量，但依赖数组没写，可能读到旧值。

**解决方案**

遵守 eslint-plugin-react-hooks 的提示。不要随便关闭规则。

### 3. 请求重复

React 开发模式下 Strict Mode 可能帮助暴露副作用问题。请求逻辑要能承受重复触发，或者使用请求库管理缓存和去重。

## 最佳实践

- Effect 用来同步外部系统。
- 派生数据优先直接计算或 useMemo，不放 Effect。
- 订阅、定时器、事件监听必须清理。
- 依赖数组按规则写，不随便忽略 lint。

## 下一步

继续学习 [路由与项目结构](/react/router-structure)。
