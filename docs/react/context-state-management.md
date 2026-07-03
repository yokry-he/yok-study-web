# Context 与状态管理

## 适合谁看

适合不清楚状态应该放组件、Context 还是全局状态库的学习者。

React 官方文档强调，状态共享最常见方式是“状态提升”：把状态移动到最近共同父组件，再通过 props 传递。Context 适合跨层传递，但不是所有状态的默认容器。

## 状态放哪里

| 状态 | 推荐位置 |
| --- | --- |
| 输入框内容 | 当前组件 |
| 弹窗开关 | 当前页面或局部组件 |
| 表格筛选 | 页面组件 |
| 当前用户 | 全局状态或 Context |
| 主题和语言 | Context 或全局状态 |
| 权限码 | 全局状态 |

## 状态提升

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

如果两个子组件需要同步状态，把状态放到它们最近的共同父组件。

## Context

```tsx
const AuthContext = createContext<AuthContextValue | null>(null)

function AuthProvider({ children }: { children: React.ReactNode }) {
  const [profile, setProfile] = useState<UserProfile | null>(null)

  return (
    <AuthContext.Provider value={{ profile, setProfile }}>
      {children}
    </AuthContext.Provider>
  )
}
```

读取：

```tsx
function useAuth() {
  const value = useContext(AuthContext)
  if (!value) {
    throw new Error('useAuth must be used inside AuthProvider')
  }
  return value
}
```

## Context 适合什么

适合：

- 主题。
- 语言。
- 当前用户。
- 权限上下文。
- 少量全局配置。

不适合：

- 高频变化的大对象。
- 表单输入。
- 表格筛选。
- 每行列表状态。

## 全局状态库

当 Context 变得复杂，可以考虑 Zustand、Redux Toolkit、Jotai 等状态库。选择前先确认问题：

- 是否跨很多页面共享。
- 是否需要持久化。
- 是否有复杂派生状态。
- 是否需要 DevTools 和调试能力。

## 实际项目常见问题

### 1. Context 更新导致大面积重渲染

**原因**

Provider value 每次渲染都创建新对象，或 Context 放了过多高频状态。

**解决方案**

- 拆分 Context。
- 高频状态留在局部。
- 必要时使用状态库。

### 2. props 层层传递很痛苦

先判断是不是组件结构问题。不要一遇到 props drilling 就立刻全局状态。

### 3. 全局状态残留

退出登录或切换组织时，要清理用户、权限、菜单和缓存数据。

## 最佳实践

- 状态默认放最近使用它的地方。
- 多个子组件共享时状态提升。
- 跨层稳定上下文用 Context。
- 高频业务状态谨慎放 Context。
- 全局状态要有清理策略。

## 下一步

继续学习 [路由与项目结构](/react/router-structure)。
