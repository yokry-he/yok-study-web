# Hooks 与状态

## 适合谁看

适合已经能写 React 组件，但对 useState、状态提升、自定义 Hook 和 Hook 调用规则还不稳定的人。

## Hooks 是什么

React 官方文档说明，Hooks 让你在组件中使用 React 的不同能力，也可以组合内置 Hooks 创建自己的 Hooks。

## useState

```tsx
const [keyword, setKeyword] = useState('')
```

状态更新后，组件会重新渲染。

## 状态提升

如果两个子组件都需要同一份状态，把状态放到它们共同的父组件。

```tsx
function UsersPage() {
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null)

  return (
    <>
      <UserTable onSelect={setSelectedUserId} />
      <UserDetail userId={selectedUserId} />
    </>
  )
}
```

## 自定义 Hook

```tsx
function useLoading() {
  const [loading, setLoading] = useState(false)

  async function run(task: () => Promise<void>) {
    setLoading(true)
    try {
      await task()
    } finally {
      setLoading(false)
    }
  }

  return { loading, run }
}
```

自定义 Hook 命名必须以 `use` 开头，并且内部可以使用其他 Hook。

## Hooks 规则

React 官方规则要求：Hook 只能在组件或自定义 Hook 顶层调用，不能放在循环、条件或嵌套函数里。

不推荐：

```tsx
if (visible) {
  const [count, setCount] = useState(0)
}
```

推荐：

```tsx
const [count, setCount] = useState(0)
```

## 实际项目常见问题

### 1. 状态更新后马上读取还是旧值

React 状态更新会触发下一次渲染，不是同步改当前变量。

如果依赖上一次状态：

```tsx
setCount((count) => count + 1)
```

### 2. 所有状态都放到全局

局部表单、弹窗、当前 tab 通常留在组件内。跨页面共享状态才考虑全局状态库。

## 下一步

继续学习 [Effect 与副作用](/react/effects)。
