# Nuxt / Next 元框架学习导览

## 这个模块解决什么

Nuxt 和 Next 解决的是“只靠浏览器端 SPA 不够用时，前端项目如何获得服务端渲染、静态生成、文件路由、接口能力、SEO 和更完整部署能力”的问题。

如果 Vue 和 React 是前端 UI 框架，那么 Nuxt 和 Next 更像应用框架。它们会帮你处理：

- 文件系统路由。
- 页面布局。
- 服务端渲染。
- 静态站点生成。
- 数据获取。
- 接口路由。
- 页面元信息和 SEO。
- 部署到 Node、静态托管、Serverless 或 Edge 环境。

这个模块不会把 Nuxt 和 Next 混成一门技术，而是先建立共同心智模型，再分别学习 Nuxt 和 Next 的项目实践。

## 适合谁看

适合已经具备以下基础的学习者：

- 会写 Vue 或 React 组件。
- 理解路由、状态、请求和构建。
- 做过普通 Vite SPA 项目。
- 开始关心 SEO、首屏速度、内容站、官网、B 端门户、全栈前端或服务端部署。

如果你还没系统学过 Vue 或 React，先看：

- [Vue 学习导览](/vue/introduction)
- [React 学习导览](/react/introduction)
- [前端工程化学习导览](/engineering/introduction)

## 什么时候需要 Nuxt 或 Next

适合使用元框架的场景：

| 场景 | 为什么适合 |
| --- | --- |
| 官网、博客、文档、营销页 | 需要 SEO、快首屏和静态生成 |
| 内容型产品 | 页面多、数据可预取、分享链接重要 |
| 中大型 Web 应用 | 需要布局、路由、数据获取和部署约定 |
| 全栈前端项目 | 希望前端仓库里包含接口、页面和服务端逻辑 |
| 多端部署 | 需要 Node、静态、Serverless、Edge 等不同形态 |

不一定需要元框架的场景：

- 纯后台管理系统，SEO 不重要。
- 项目主要在登录后使用。
- 团队只熟悉客户端 SPA，且交付周期很短。
- 后端已经提供完整页面和 API，前端只做局部交互。

后台系统也可以用 Nuxt 或 Next，但不要为了“看起来高级”引入 SSR。先判断业务是否真的需要服务端渲染、静态生成或全栈路由。

## 学习路线

推荐顺序：

```text
Vue 或 React 基础
↓
Vite 与工程化
↓
浏览器与 HTTP
↓
Nuxt / Next 共同模型
↓
Nuxt 或 Next 选择一个深入
↓
路由、布局、数据获取
↓
部署和缓存
↓
真实项目问题排查
```

如果你是 Vue 方向，优先学 Nuxt。

如果你是 React 方向，优先学 Next。

如果你已经会 Vue 和 React，可以先学本模块的共同模型，再根据项目类型选择。

## 共同心智模型

Nuxt 和 Next 虽然生态不同，但核心问题非常相似。

| 能力 | Nuxt | Next |
| --- | --- | --- |
| 基础框架 | Vue | React |
| 路由方式 | 文件路由 | 文件路由 |
| 布局 | layouts | layouts |
| 数据获取 | `useFetch`、`useAsyncData` 等 | Server Components、fetch、Server Actions 等 |
| 接口能力 | server routes | route handlers、server functions |
| 渲染方式 | SSR、SSG、客户端渲染 | SSR、SSG、动态渲染、流式渲染等 |
| 部署 | Node、静态、Serverless、Edge | Node、Docker、静态导出、平台适配 |

学习时不要只背 API。要先理解这几个问题：

- 页面是在服务器渲染，还是浏览器渲染。
- 数据是在服务器取，还是浏览器取。
- 哪些内容可以缓存，哪些内容必须实时。
- 代码会运行在 Node、浏览器，还是 Edge。
- 页面是否需要被搜索引擎和社交分享正确识别。

## 目录结构的变化

普通 SPA 项目通常是：

```text
src/
├─ router/
├─ views/
├─ components/
├─ stores/
└─ api/
```

元框架项目更强调约定：

```text
pages 或 app/      页面和路由
layouts/           布局
components/        组件
server/            服务端接口或逻辑
composables/       可复用逻辑
middleware/        中间件或守卫
public/            静态资源
```

目录本身就是框架能力的一部分。不要把所有文件都塞进 `components`，也不要把服务端逻辑写进客户端组件。

## 实际项目中的取舍

### 官网和内容站

优先考虑：

- 静态生成。
- 页面元信息。
- 图片优化。
- 缓存策略。
- 构建和预览流程。

### 业务系统

优先考虑：

- 登录态和鉴权。
- 服务端和客户端状态边界。
- 权限路由。
- 接口错误处理。
- 部署形态。

### 全栈前端项目

优先考虑：

- 服务端接口边界。
- 数据库访问是否放在框架内。
- 日志和错误追踪。
- 环境变量隔离。
- 安全和权限校验。

## 本模块内容

| 文档 | 解决的问题 |
| --- | --- |
| [Nuxt 项目实践](/meta-frameworks/nuxt) | Vue 方向如何进入 Nuxt |
| [Next.js 项目实践](/meta-frameworks/next) | React 方向如何进入 Next |
| [路由、布局与数据获取](/meta-frameworks/routing-data) | 元框架最核心的页面和数据模型 |
| [部署、缓存与运行时](/meta-frameworks/deployment) | 如何选择部署方式和缓存策略 |
| [服务端鉴权与登录态](/meta-frameworks/server-auth) | SSR、API Route、Cookie 会话、401/403 和缓存边界 |
| [SEO、Metadata 与结构化数据](/meta-frameworks/seo-metadata) | title、description、Open Graph、sitemap、robots 和结构化数据 |
| [国际化与多语言站点](/meta-frameworks/i18n) | locale 路由、翻译字典、多语言 SEO 和缓存 |
| [内容站案例：技术博客与官网](/meta-frameworks/content-site-case) | 把路由、内容、SEO、国际化、缓存和部署组合成完整项目 |
| [常见问题](/meta-frameworks/troubleshooting) | SSR、hydration、环境变量、缓存等问题排查 |

## 下一步学习

如果你是 Vue 方向，继续看 [Nuxt 项目实践](/meta-frameworks/nuxt)。

如果你是 React 方向，继续看 [Next.js 项目实践](/meta-frameworks/next)。

如果你想先理解共性，继续看 [路由、布局与数据获取](/meta-frameworks/routing-data)。
