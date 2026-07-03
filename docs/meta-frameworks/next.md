# Next.js 项目实践

## 适合谁看

适合已经学过 React，并希望用 React 构建官网、内容站、门户、全栈前端应用、SSR 应用或需要 SEO 的产品页面的学习者。

Next.js 不只是 React 的脚手架。它提供了文件路由、布局、Server Components、数据获取、缓存、接口路由、图片和字体优化、部署约定等能力。

## Next 适合什么项目

| 项目类型 | 适合原因 |
| --- | --- |
| 官网和营销页 | SEO、元信息、首屏速度重要 |
| 内容站和博客 | 静态生成和动态渲染都可支持 |
| SaaS 前台 | 路由、布局、数据获取、登录态复杂 |
| 全栈前端应用 | 可组织页面、接口和服务端逻辑 |
| 需要渐进加载的页面 | 支持流式渲染和更细粒度加载状态 |

如果项目是纯内部后台管理系统，普通 React + Vite 可能更直接。Next 的价值在于页面渲染模型、服务器能力和部署体系。

## App Router 心智模型

现代 Next 项目通常优先理解 App Router。

常见结构：

```text
app/
├─ layout.tsx
├─ page.tsx
├─ loading.tsx
├─ error.tsx
├─ users/
│  ├─ page.tsx
│  └─ [id]/
│     └─ page.tsx
components/
lib/
public/
next.config.ts
```

核心文件：

| 文件 | 作用 |
| --- | --- |
| `layout.tsx` | 布局 |
| `page.tsx` | 页面 |
| `loading.tsx` | 加载状态 |
| `error.tsx` | 错误边界 |
| `not-found.tsx` | 404 页面 |
| `route.ts` | Route Handler |

## Server Component 和 Client Component

默认情况下，App Router 中的组件倾向于服务端组件。需要浏览器交互时，使用 `"use client"`。

服务端组件适合：

- 读取数据库或服务端 API。
- 渲染静态内容。
- 减少客户端 JavaScript。
- 做 SEO 页面。

客户端组件适合：

- 点击、输入、拖拽等交互。
- 使用浏览器 API。
- 使用 React state 和 effect。
- 使用依赖浏览器环境的组件库。

客户端组件：

```tsx
'use client'

import { useState } from 'react'

export function Counter() {
  const [count, setCount] = useState(0)

  return <button onClick={() => setCount(count + 1)}>{count}</button>
}
```

不要为了省事把整个页面都标成客户端组件。这样会丢掉很多服务端渲染和性能收益。

## 文件路由

首页：

```text
app/page.tsx -> /
```

用户列表：

```text
app/users/page.tsx -> /users
```

用户详情：

```text
app/users/[id]/page.tsx -> /users/1001
```

动态参数：

```tsx
export default async function UserPage({
  params
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  return <div>User {id}</div>
}
```

不同版本和配置下参数类型写法可能有差异，项目中应以当前 Next 官方模板和类型提示为准。

## 数据获取

服务端组件中可以直接 `await` 数据。

```tsx
async function getUsers() {
  const res = await fetch('https://api.example.com/users')

  if (!res.ok) {
    throw new Error('Failed to fetch users')
  }

  return res.json()
}

export default async function UsersPage() {
  const users = await getUsers()

  return <UserList users={users} />
}
```

项目里要明确：

- 数据是否可以缓存。
- 页面是否需要实时数据。
- 错误如何展示。
- 登录态如何传递。
- 哪些请求只能在服务端执行。

## Route Handler

```ts
// app/api/health/route.ts
export async function GET() {
  return Response.json({ ok: true })
}
```

适合：

- 轻量 API。
- BFF 聚合。
- Webhook。
- 读取服务端环境变量。
- 和 Server Actions 配合处理业务动作。

不适合：

- 无边界地替代完整后端。
- 没有鉴权就直接写关键数据。
- 把大量业务逻辑散落在多个页面目录里。

## 缓存和重新验证

Next 的缓存能力很强，但也容易误用。

学习时先问三个问题：

```text
这个页面是否允许缓存？
缓存多久可以接受？
用户操作后是否需要刷新数据？
```

常见场景：

| 场景 | 建议 |
| --- | --- |
| 官网内容 | 可以静态化或较长缓存 |
| 商品详情 | 可缓存，但要考虑更新 |
| 用户个人信息 | 通常不应公共缓存 |
| 后台列表 | 更偏动态请求 |

不要在不了解缓存规则时盲目上线关键业务页面。

## 实际项目建议

- 默认先区分 Server Component 和 Client Component。
- 交互组件尽量小，不要把整页变成客户端组件。
- 数据获取封装到 `lib/` 或 service 层。
- Route Handler 要有鉴权、错误处理和日志。
- 部署前确认目标平台支持你的渲染和缓存能力。

## 常见坑

| 问题 | 处理 |
| --- | --- |
| 浏览器 API 报错 | 只在客户端组件或 effect 中使用 |
| 页面数据没更新 | 检查缓存和 revalidate 策略 |
| hydration mismatch | 避免服务端和客户端初始输出不同 |
| 组件库不能在服务端渲染 | 放到客户端组件边界内 |
| API 密钥泄漏 | 私密变量只能在服务端使用 |

## 下一步学习

- [路由、布局与数据获取](/meta-frameworks/routing-data)
- [部署、缓存与运行时](/meta-frameworks/deployment)
- [常见问题](/meta-frameworks/troubleshooting)
- [React 学习导览](/react/introduction)
