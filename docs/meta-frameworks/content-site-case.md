# 内容站案例：技术博客与官网

## 适合谁看

适合已经学过 Nuxt / Next 的路由、数据获取、部署和 SEO，但还没有把它们组合成完整内容站的人：

- 会写页面，但不知道内容模型怎么设计。
- 文章详情、分类、标签、搜索、SEO 分散在各处。
- 不知道内容站应该 SSR、SSG 还是 ISR。
- 发布内容后缓存和 sitemap 没同步。
- 多语言内容没有治理方式。

这个案例不追求覆盖所有业务，而是建立一个可复用的内容站骨架：官网首页、文章列表、文章详情、分类标签、SEO、国际化、部署和排错。

## 项目目标

做一个技术内容站：

- 首页展示品牌、主推内容和最新文章。
- 文章列表支持分类、标签和分页。
- 文章详情支持 SEO、目录、相关推荐。
- 支持中英文。
- 支持 sitemap。
- 支持静态生成或增量更新。
- 有上线检查和回滚方案。

## 推荐目录模型

Nuxt 内容站示例：

```text
content/
├─ zh/
│  └─ articles/
│     └─ vue-reactivity.md
└─ en/
   └─ articles/
      └─ vue-reactivity.md
```

Next 内容站示例：

```text
src/
├─ app/
│  └─ [lang]/
│     ├─ page.tsx
│     └─ articles/
│        ├─ page.tsx
│        └─ [slug]/
│           └─ page.tsx
├─ content/
└─ lib/
   └─ content.ts
```

关键不是目录名字，而是内容读取、路由生成、SEO 和缓存要有统一入口。

## 内容模型

一篇文章至少包含：

```yaml
title: Vue 3 响应式入门
description: 用简单示例理解 ref、reactive 和 computed。
slug: vue-reactivity
locale: zh
cover: /images/articles/vue-reactivity.png
category: Vue
tags:
  - Vue
  - 响应式
publishedAt: 2026-07-01
updatedAt: 2026-07-02
draft: false
```

建议区分：

| 字段 | 用途 |
| --- | --- |
| title | 页面标题和 H1 |
| description | SEO 描述和列表摘要 |
| slug | URL |
| locale | 语言 |
| cover | 分享图和列表图 |
| category | 分类 |
| tags | 标签 |
| draft | 是否草稿 |
| updatedAt | 内容更新和 sitemap |

SEO 字段不要完全依赖正文自动截取，重要页面应该人工维护。

## 路由设计

推荐：

```text
/zh
/zh/articles
/zh/articles/vue-reactivity
/zh/categories/vue
/zh/tags/vue

/en
/en/articles
/en/articles/vue-reactivity
```

不要把语言、分类、标签都放到 query 里。公开内容更适合稳定路径，方便分享、缓存和 SEO。

## 页面渲染策略

| 页面 | 推荐策略 |
| --- | --- |
| 首页 | SSG 或 ISR |
| 文章列表 | SSG/ISR，内容量大时分页生成 |
| 文章详情 | SSG/ISR |
| 搜索页 | 客户端搜索或服务端搜索 |
| 用户中心 | SSR 或客户端渲染 |

公开内容优先静态生成或增量更新。登录态内容不要静态生成成公共页面。

## 数据读取层

把内容读取集中到 `content service`。

```ts
export async function getArticle(locale: string, slug: string) {
  const article = await contentRepository.findArticle(locale, slug)

  if (!article || article.draft) {
    return null
  }

  return article
}
```

页面不应该到处直接读文件或请求 CMS。集中封装后，后续从 Markdown 切到 CMS、数据库或远程 API 更容易。

## SEO 生成

详情页 metadata 来自文章模型：

```ts
function createArticleSeo(article: Article) {
  return {
    title: article.title,
    description: article.description,
    openGraph: {
      title: article.title,
      description: article.description,
      images: [article.cover]
    }
  }
}
```

每个详情页都要能生成：

- title。
- description。
- canonical。
- og:image。
- 结构化数据。
- 多语言 hreflang。

## sitemap 更新

sitemap 应包含：

- 首页。
- 文章列表页。
- 文章详情页。
- 分类页。
- 标签页。

不要包含：

- 草稿。
- 预览页。
- 登录页。
- 后台页面。
- 搜索组合页。

内容发布后，要确保 sitemap 更新或重新生成。

## 缓存策略

内容站常见缓存：

| 层 | 策略 |
| --- | --- |
| HTML | ISR 或短缓存 |
| JS/CSS | 文件 hash 长缓存 |
| 图片 | 长缓存，必要时换文件名 |
| CMS API | 短缓存或按内容版本缓存 |
| sitemap | 内容发布后更新 |

内容更新后页面没变，通常是缓存没失效。

解决方向：

- 发布内容后触发重建或 revalidate。
- 图片更新换 URL 或刷新 CDN。
- 详情页缓存 key 包含 locale 和 slug。

## 多语言内容治理

内容模型要记录语言关系：

```yaml
translationKey: vue-reactivity
locale: zh
slug: vue-reactivity
```

英文版本：

```yaml
translationKey: vue-reactivity
locale: en
slug: vue-reactivity
```

这样可以实现：

- 语言切换保持当前文章。
- hreflang 自动生成。
- 找出缺失翻译。
- 按语言独立发布。

## 实际项目问题

### 1. 文章更新后页面还是旧内容

**原因**

静态生成或 CDN 缓存没有失效。

**解决方案**

- 内容发布后触发构建或 revalidate。
- 刷新对应路径 CDN。
- 列表页和详情页都要更新。

### 2. 草稿被搜索引擎收录

**原因**

构建时没有过滤 draft，或者预览环境可公开访问。

**解决方案**

- 构建 sitemap 时排除 draft。
- 草稿页加鉴权。
- 预发环境 noindex。

### 3. 分类页越来越慢

**原因**

每次请求都扫描全部文章并动态过滤。

**解决方案**

- 构建时生成分类索引。
- 数据读取层缓存文章列表。
- 内容量大时使用搜索服务或数据库索引。

### 4. 中英文文章无法互相切换

**原因**

只按 slug 匹配，没有 translationKey。

**解决方案**

- 内容模型增加 translationKey。
- 语言切换根据 translationKey 找目标语言文章。
- 缺失翻译时展示提示或返回语言首页。

## 最佳实践

- 内容读取、SEO 生成、sitemap 生成都要有统一工具层。
- 公开内容优先 SSG/ISR，不要把登录态内容静态化。
- 内容模型包含 SEO、语言、发布时间和草稿状态。
- sitemap 排除草稿、预览和后台页面。
- 多语言内容使用 translationKey 关联。
- 发布内容时同步处理缓存刷新和 sitemap。
- 案例项目要包含上线检查和回滚方案。

## 下一步学习

继续学习 [Nuxt / Next 常见问题](/meta-frameworks/troubleshooting)，把 SSR、缓存、hydration、环境变量和部署问题放到实际排查链路里。
