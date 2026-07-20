# Nuxt 项目实践

## 适合谁看

适合已经学过 Vue 3，并希望用 Vue 做官网、内容站、文档站、门户站、全栈前端项目或 SSR 项目的学习者。

Nuxt 不是“更复杂的 Vue”，而是一套围绕 Vue 的应用框架。它把路由、布局、数据获取、服务端渲染、接口和部署约定组合在一起，让你少写很多基础工程代码。

## Nuxt 适合什么项目

| 项目类型 | 适合原因 |
| --- | --- |
| 官网和营销页 | SEO、静态生成、首屏速度重要 |
| 博客和内容站 | 文件路由、数据预取、元信息管理方便 |
| 文档站 | 内容结构清晰，适合静态部署 |
| 门户和用户中心 | 需要布局、登录态和服务端数据 |
| Vue 全栈项目 | 可以在同一项目里组织页面和 server routes |

如果只是一个内部后台管理系统，且不需要 SEO，普通 Vite + Vue Admin 可能更简单。Nuxt 的优势主要在页面渲染、内容、部署形态和全栈能力。

## 基本目录结构

常见结构：

```text
app.vue
nuxt.config.ts
pages/
├─ index.vue
├─ about.vue
└─ products/
   └─ [id].vue
layouts/
├─ default.vue
components/
composables/
server/
└─ api/
   └─ users.get.ts
public/
```

关键目录：

| 目录 | 作用 |
| --- | --- |
| `pages/` | 文件路由 |
| `layouts/` | 页面布局 |
| `components/` | 可复用组件 |
| `composables/` | 组合式逻辑 |
| `server/api/` | 服务端 API |
| `public/` | 直接公开的静态资源 |

## 文件路由

`pages/index.vue` 对应首页：

```text
/
```

`pages/about.vue` 对应：

```text
/about
```

动态路由：

```text
pages/products/[id].vue
```

对应：

```text
/products/1001
```

页面里读取参数：

```vue
<script setup lang="ts">
const route = useRoute()
const productId = route.params.id
</script>
```

## 布局

默认布局：

```vue
<!-- layouts/default.vue -->
<template>
  <header>站点头部</header>
  <main>
    <slot />
  </main>
  <footer>站点底部</footer>
</template>
```

页面使用指定布局：

```vue
<script setup lang="ts">
definePageMeta({
  layout: 'default'
})
</script>
```

布局适合放：

- 顶部导航。
- 页脚。
- 侧边栏。
- 页面容器。
- 登录后统一外壳。

不要把具体页面业务逻辑写进布局。

## 数据获取

Nuxt 常用数据获取方式包括 `useFetch`、`useAsyncData` 和 `$fetch`。它们的核心差异在于是否绑定响应式状态、是否适合 SSR、是否需要手动管理 key。

基础示例：

```vue
<script setup lang="ts">
const { data, pending, error, refresh } = await useFetch('/api/products')
</script>

<template>
  <p v-if="pending">加载中...</p>
  <p v-else-if="error">加载失败</p>
  <ProductList v-else :items="data?.items || []" />
</template>
```

项目里要注意：

- 服务端渲染时，请求可能发生在服务端。
- 浏览器端导航时，请求也可能重新触发。
- 登录态、cookie、header 需要明确传递。
- 不要在组件各处散落重复请求逻辑。

## server api

示例：

```ts
// server/api/health.get.ts
export default defineEventHandler(() => {
  return {
    ok: true
  }
})
```

访问：

```text
GET /api/health
```

server api 适合：

- 聚合后端接口。
- 处理服务端密钥。
- 做轻量 BFF。
- 读取服务端运行时配置。

不适合：

- 把复杂后端系统全部塞进前端项目。
- 在没有鉴权和日志的情况下直接操作关键数据。
- 把数据库密钥暴露到客户端。

## 环境变量和运行时配置

Nuxt 项目要区分公开配置和私有配置。

```ts
export default defineNuxtConfig({
  runtimeConfig: {
    apiSecret: '',
    public: {
      apiBase: '/api'
    }
  }
})
```

页面中读取公开配置：

```ts
const config = useRuntimeConfig()
const apiBase = config.public.apiBase
```

服务端可以读取私有配置，但客户端不能。

## 实际项目建议

- 官网和内容站优先考虑静态生成和 SEO。
- 登录后系统要谨慎处理服务端渲染和用户态。
- API 请求封装要区分服务端和客户端。
- server api 要有日志、错误处理和权限校验。
- 部署前确认目标平台支持的 Nuxt 能力。

## 常见坑

| 问题 | 处理 |
| --- | --- |
| 页面刷新数据不一致 | 检查数据是在服务端取还是客户端取 |
| hydration mismatch | 避免服务端和客户端渲染不同随机值 |
| cookie 没带上 | 检查服务端请求 header 转发 |
| 静态部署后接口不可用 | 静态站点没有 Node server，需要独立 API |
| 环境变量泄漏 | 私有配置不要放到 public |

## 下一步学习

- [图解 Nuxt / Next 元框架核心概念](/meta-frameworks/visual-guide)
- [路由、布局与数据获取](/meta-frameworks/routing-data)
- [课程内容平台从零到项目](/meta-frameworks/project-from-zero)
- [部署、缓存与运行时](/meta-frameworks/deployment)
- [常见问题](/meta-frameworks/troubleshooting)
- [Vue 学习导览](/vue/introduction)
