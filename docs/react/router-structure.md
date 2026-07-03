# 路由与项目结构

## 适合谁看

适合准备把 React 从单页面组件推进到真实多页面项目的学习者。

## 推荐目录

```text
src/
├─ api/
├─ components/
├─ hooks/
├─ layouts/
├─ pages/
├─ router/
├─ services/
├─ stores/
├─ styles/
└─ types/
```

## React Router

React Router 官方文档提供了路由能力，可以组织 URL、页面和布局。

基础结构：

```tsx
import { createBrowserRouter } from 'react-router'

export const router = createBrowserRouter([
  {
    path: '/',
    Component: AppLayout,
    children: [
      { path: 'dashboard', Component: DashboardPage },
      { path: 'users', Component: UsersPage }
    ]
  }
])
```

## 受保护路由

```tsx
function RequireAuth({ children }: { children: React.ReactNode }) {
  const token = useAuthStore((state) => state.token)

  if (!token) {
    return <Navigate to="/login" replace />
  }

  return children
}
```

## Next.js App Router

如果使用 Next.js，官方 App Router 是基于文件系统的路由，并支持 Server Components 等 React 新能力。学习时要区分：

- 普通 React SPA。
- React Router 项目。
- Next.js App Router 项目。

不要把三者的路由规则混在一起。

## 实际项目常见问题

### 1. 页面和组件目录混乱

**解决方案**

页面放 `pages` 或 `routes`，可复用组件放 `components`，业务流程放 `services`。

### 2. 权限判断散落

封装受保护路由或权限组件，不要每个页面重复判断。

### 3. 刷新 404

React SPA 使用 history 路由时，同样需要服务器 fallback 到 `index.html`。

## 下一步

继续学习 [常见问题](/react/troubleshooting)。
