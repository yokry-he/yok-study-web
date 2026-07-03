# React 性能优化

## 适合谁看

适合遇到 React 页面输入卡顿、列表渲染慢、组件频繁重渲染、包体积过大的学习者。

性能优化的第一原则是：先确认瓶颈，再优化。React DevTools 可以帮助检查组件渲染和性能问题。

## 常见性能问题

| 类型 | 表现 | 常见原因 |
| --- | --- | --- |
| 首屏慢 | 页面白屏时间长 | 包体积大、接口慢 |
| 输入卡 | 输入框延迟 | 渲染范围太大 |
| 列表卡 | 滚动不顺 | DOM 数量过多 |
| 重复渲染 | 子组件频繁更新 | props 不稳定 |

## 避免不必要的派生状态

React 官方文档强调，很多数据转换不需要 Effect。可以直接在渲染时计算：

```tsx
const enabledUsers = users.filter((user) => user.enabled)
```

如果计算很重，再考虑 `useMemo`：

```tsx
const enabledUsers = useMemo(() => {
  return users.filter((user) => user.enabled)
}, [users])
```

## React.memo

```tsx
const UserRow = memo(function UserRow({ user }: { user: User }) {
  return <div>{user.username}</div>
})
```

只有在子组件渲染成本明显、props 相对稳定时使用。不要无脑给所有组件加 `memo`。

## useMemo 和 useCallback

`useMemo` 缓存计算结果：

```tsx
const options = useMemo(() => {
  return users.map((user) => ({
    label: user.username,
    value: user.id
  }))
}, [users])
```

`useCallback` 缓存函数引用：

```tsx
const handleEdit = useCallback((user: User) => {
  setEditingUser(user)
}, [])
```

这些工具用于解决具体性能问题，不是默认写法。

## 列表优化

- 使用稳定 key。
- 服务端分页。
- 大列表使用虚拟滚动。
- 列表项组件拆分。
- 避免每行创建复杂对象。

## 包体积优化

- 路由级懒加载。
- 图表、富文本、编辑器异步加载。
- 图标按需导入。
- 避免引入整个工具库。

```tsx
const CodeEditor = lazy(() => import('./CodeEditor'))
```

## 实际项目常见问题

### 1. 输入框每输入一个字都卡

**原因**

输入状态放在过高层级，导致大页面重渲染。

**解决方案**

- 状态下沉到表单组件。
- 拆分重组件。
- 对重列表分页或虚拟化。

### 2. memo 没有效果

**原因**

传给子组件的对象或函数每次都是新的。

**解决方案**

先确认是否真的需要 memo。如果需要，再稳定 props。

### 3. 首屏加载大

**解决方案**

路由懒加载，重组件异步加载。

## 最佳实践

- 先用 DevTools 定位问题。
- 不要默认使用 memo、useMemo、useCallback。
- 状态放在尽量低的位置。
- 列表优先分页。
- 重依赖异步加载。

## 下一步

继续学习 [测试策略](/react/testing)。
