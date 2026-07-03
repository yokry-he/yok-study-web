# 组件与 JSX

## 适合谁看

适合已经知道 React 是组件化框架，但还不熟悉 JSX、Props、条件渲染、列表渲染和组件拆分的人。

## JSX 是什么

JSX 是 JavaScript 里的 UI 描述语法：

```tsx
const title = '用户管理'

export function PageTitle() {
  return <h1>{title}</h1>
}
```

JSX 中使用 `{}` 写 JavaScript 表达式。

## Props

```tsx
interface UserCardProps {
  username: string
  enabled: boolean
}

export function UserCard({ username, enabled }: UserCardProps) {
  return (
    <article>
      <strong>{username}</strong>
      <span>{enabled ? '启用' : '停用'}</span>
    </article>
  )
}
```

使用：

```tsx
<UserCard username="alice" enabled={true} />
```

## 条件渲染

```tsx
{loading ? <Loading /> : <UserTable users={users} />}
```

权限按钮：

```tsx
{can('system:user:create') && (
  <button type="button">新增用户</button>
)}
```

## 列表渲染

```tsx
{users.map((user) => (
  <UserCard key={user.id} username={user.username} enabled={user.enabled} />
))}
```

`key` 必须稳定唯一，不要使用 index。

## 组件拆分

用户管理页面：

```text
UsersPage
├─ UserSearchForm
├─ UserTable
└─ UserFormDialog
```

父组件组织流程，子组件负责展示和局部交互。

## 实际项目常见问题

### 1. map 渲染列表忘记 key

**影响**

React 无法稳定识别列表项，可能导致输入框状态错乱。

**解决方案**

使用业务 id：

```tsx
users.map((user) => <UserRow key={user.id} user={user} />)
```

### 2. JSX 里写太多复杂逻辑

**解决方案**

把复杂计算提前到变量或函数：

```tsx
const enabledUsers = users.filter((user) => user.enabled)
```

## 下一步

继续学习 [Hooks 与状态](/react/hooks-state)。
