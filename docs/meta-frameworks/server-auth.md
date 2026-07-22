# 服务端鉴权与登录态

## 适合谁看

适合已经能写 Nuxt 或 Next 页面，但开始处理登录、权限、服务端数据获取和 API Route 安全时容易混乱的人：

- 页面上隐藏了按钮，但接口仍然能被调用。
- SSR 页面不知道怎么读取登录态。
- Token 放 localStorage，服务端渲染时拿不到。
- API Route 没有鉴权就直接访问数据库。
- 登录后缓存导致不同用户看到错误数据。

元框架的鉴权难点在于：代码可能运行在浏览器，也可能运行在服务端。你必须先判断当前逻辑在哪个运行环境执行，再决定如何读取会话、跳转和缓存。

## 先区分三个概念

| 概念 | 说明 |
| --- | --- |
| Authentication 认证 | 用户是谁，例如是否登录 |
| Authorization 授权 | 用户能做什么，例如是否能删除订单 |
| Session 登录态 | 服务端和浏览器之间保持登录状态的机制 |

不要把“能访问页面”和“能调用接口”混为一谈。页面保护和接口保护都要做。

## 推荐登录态位置

SSR 和服务端数据获取场景，优先使用 HttpOnly Cookie 承载会话标识。

原因：

- 服务端可以从请求 Cookie 读取登录态。
- 浏览器 JavaScript 不能直接读取 HttpOnly Cookie。
- 避免把 Token 暴露给前端脚本。
- 页面刷新、SSR、API Route 都能使用同一套会话。

不推荐把核心访问令牌只放在 localStorage。localStorage 只在浏览器可用，SSR 阶段拿不到，并且更容易受到 XSS 影响。

## 页面保护和接口保护

服务端应该在读取和输出敏感数据之前完成会话与权限判断。下图展示同一个受保护路由如何在无权限时直接重定向，在有权限时才读取页面数据。

<DocFigure
  src="/images/meta-frameworks/server-auth-redirect.webp"
  alt="服务端收到管理页面请求后读取会话、检查权限，无权限重定向，有权限才渲染受保护数据"
  caption="客户端跳转不能充当安全边界；页面和业务接口都必须独立鉴权。"
  :width="1440"
  :height="900"
/>

无权限请求不应先执行昂贵查询或序列化敏感结果，再靠页面隐藏；授权结论要尽量靠近数据访问边界。

完整鉴权至少有两层：

```text
页面访问保护
↓
API / Server Route 保护
```

页面保护负责用户体验：

- 未登录跳转登录页。
- 无权限展示 403 页面。
- 登录后回到原页面。

接口保护负责安全：

- 未登录返回 401。
- 无权限返回 403。
- 不能只相信前端页面是否展示按钮。

## Nuxt 中的鉴权位置

Nuxt 常见位置：

| 位置 | 用途 |
| --- | --- |
| route middleware | 页面跳转前检查登录态 |
| server middleware | Nitro 请求进入服务端路由前处理 |
| server/api | 服务端接口，必须自己校验权限 |
| composable | 封装 `useSession`、`useUser` 等状态 |

页面中可以声明 middleware：

```vue
<script setup lang="ts">
definePageMeta({
  middleware: ['auth']
})
</script>
```

示例 route middleware：

```ts
export default defineNuxtRouteMiddleware(async (to) => {
  const { data: session } = await useFetch('/api/session')

  if (!session.value?.user) {
    return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})
```

服务端接口仍然要检查：

```ts
export default defineEventHandler(async (event) => {
  const session = await requireUserSession(event)

  if (!session.user.permissions.includes('order:read')) {
    throw createError({ statusCode: 403, statusMessage: '无权限' })
  }

  return orderService.findOrders(session.user.id)
})
```

## Next 中的鉴权位置

Next App Router 常见位置：

| 位置 | 用途 |
| --- | --- |
| middleware.ts | 请求进入页面前做轻量判断 |
| layout / page server component | 服务端读取 session，决定渲染或 redirect |
| Route Handler | API 接口鉴权 |
| Server Action | 表单提交和服务端动作鉴权 |

页面服务端保护示例：

```tsx
import { redirect } from 'next/navigation'

export default async function DashboardPage() {
  const session = await getSession()

  if (!session) {
    redirect('/login')
  }

  return <Dashboard user={session.user} />
}
```

Route Handler 保护示例：

```ts
export async function GET() {
  const session = await getSession()

  if (!session) {
    return Response.json({ message: '未登录' }, { status: 401 })
  }

  if (!session.user.permissions.includes('order:read')) {
    return Response.json({ message: '无权限' }, { status: 403 })
  }

  return Response.json(await orderService.findOrders(session.user.id))
}
```

## SSR 数据获取的缓存风险

登录态页面最容易出的问题是缓存。

如果页面内容和用户有关，就不能当成全站公共缓存。

危险场景：

```text
用户 A 请求 /account
↓
服务端渲染 A 的数据
↓
页面被错误公共缓存
↓
用户 B 看到 A 的信息
```

处理原则：

- 用户态页面禁用公共缓存。
- 数据请求携带 Cookie 时谨慎缓存。
- 缓存 key 必须包含用户、租户、权限等边界。
- 静态生成页面只放公开内容。

## 401、403 和 404

| 状态 | 含义 |
| --- | --- |
| 401 | 未登录或登录态失效 |
| 403 | 已登录但无权限 |
| 404 | 资源不存在，或业务上不希望暴露存在性 |

后台系统一般可以明确返回 403。面向公网的资源详情页，有时会对不可见资源返回 404，避免暴露资源是否存在。

## 实际项目问题

### 1. SSR 页面登录后仍显示未登录

**原因**

登录态只存 localStorage，服务端渲染时拿不到。

**解决方案**

- 使用 Cookie 会话。
- 服务端从请求 Cookie 读取 session。
- 客户端状态只作为展示缓存，不作为唯一权威。

### 2. 前端页面保护了，接口仍然能调用

**原因**

只做了 route middleware，没有做 API Route 鉴权。

**解决方案**

- 每个服务端接口都校验 session。
- 权限判断放在服务端。
- 高风险操作写审计日志。

### 3. 登录用户看到其他用户数据

**原因**

用户态数据被公共缓存。

**解决方案**

- 禁止用户态页面公共缓存。
- 区分公开页面和私有页面。
- 缓存 key 包含用户或租户边界。

### 4. 无限重定向

**原因**

登录页也被 auth middleware 保护，或 redirect 参数没有过滤。

**解决方案**

- 登录页跳过 auth 检查。
- 已登录用户访问登录页时跳到首页。
- redirect 只允许站内路径，避免开放重定向。

## 最佳实践

- SSR 登录态优先使用 HttpOnly Cookie。
- 页面保护和接口保护都要做。
- 用户态页面不要被公共缓存。
- 服务端接口区分 401 和 403。
- 权限判断基于服务端 session 和数据库数据。
- redirect 参数必须校验。
- 登录、退出、权限变更和高风险操作写审计日志。

## 参考资料

- [Next.js Authentication Guide](https://nextjs.org/docs/app/guides/authentication)
- [Next.js Learn: Adding Authentication](https://nextjs.org/learn/dashboard-app/adding-authentication)
- [Nuxt Sessions and Authentication](https://nuxt.com/docs/4.x/guide/recipes/sessions-and-authentication)
- [Nuxt Middleware](https://nuxt.com/docs/4.x/directory-structure/app/middleware)
- [Nuxt Server Directory](https://nuxt.com/docs/4.x/directory-structure/server)

## 下一步学习

继续学习 [SEO、Metadata 与结构化数据](/meta-frameworks/seo-metadata)，理解公开页面如何被搜索引擎和社交平台正确识别。
