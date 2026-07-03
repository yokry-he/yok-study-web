# 请求与数据流

## 适合谁看

适合在 React 项目中处理接口请求、列表加载、搜索、防重复请求和状态流转的学习者。

React 官方文档强调 Effect 用于同步外部系统。网络请求属于外部系统，但不是所有数据转换都应该放进 Effect。

## 推荐分层

```text
src/
├─ api/          请求函数
├─ services/     业务流程
├─ pages/        页面状态和交互
├─ components/   展示组件
└─ hooks/        可复用逻辑
```

页面负责什么时候请求，API 负责怎么请求。

## API 函数

```ts
export function getUserList(params: UserQuery) {
  return request.get<PageResult<User>>('/users', { params })
}
```

## 页面请求

```tsx
function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)

  async function fetchUsers() {
    setLoading(true)
    try {
      const result = await getUserList({ page: 1, pageSize: 20 })
      setUsers(result.items)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [])

  return <UserTable users={users} loading={loading} />
}
```

## 搜索请求

```tsx
const [keyword, setKeyword] = useState('')
const [page, setPage] = useState(1)

useEffect(() => {
  fetchUsers({ keyword, page })
}, [keyword, page])
```

如果输入频繁，应加防抖或使用请求管理库。

## 旧请求覆盖新请求

```tsx
const requestIdRef = useRef(0)

async function fetchUsers(query: UserQuery) {
  const requestId = ++requestIdRef.current
  const result = await getUserList(query)

  if (requestId !== requestIdRef.current) return

  setUsers(result.items)
}
```

## 自定义 Hook

```tsx
function useUserList() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)

  async function fetchUsers(query: UserQuery) {
    setLoading(true)
    try {
      const result = await getUserList(query)
      setUsers(result.items)
    } finally {
      setLoading(false)
    }
  }

  return { users, loading, fetchUsers }
}
```

Hook 不应该隐藏太多副作用。调用者应能看出什么时候请求。

## 实际项目常见问题

### 1. 页面打开请求两次

开发模式 Strict Mode 可能额外执行一次，帮助发现副作用问题。请求逻辑要具备幂等意识。

### 2. 搜索结果显示旧数据

使用请求序号、AbortController 或请求管理库。

### 3. 子组件也在请求同一份数据

统一数据入口。页面请求数据，子组件通过 props 接收。

## 最佳实践

- API 函数和页面状态分开。
- 请求必须有 loading 和错误处理。
- 快速输入请求要防抖或取消。
- 子组件不要重复请求父组件已有数据。
- 复杂服务端状态可以使用成熟请求管理库。

## 下一步

继续学习 [Context 与状态管理](/react/context-state-management)。
