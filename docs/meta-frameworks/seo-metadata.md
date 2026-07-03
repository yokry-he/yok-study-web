# SEO、Metadata 与结构化数据

## 适合谁看

适合准备做官网、博客、文档、产品详情页、营销页或内容站的人：

- 页面能打开，但搜索结果标题和描述不对。
- 分享到微信、Slack、X 时没有正确卡片。
- 多语言页面没有 hreflang。
- 动态详情页所有 metadata 都一样。
- 不知道 SSR、SSG 和 SEO 的关系。

SEO 不是只写几个 meta 标签。它涉及页面是否能被抓取、内容是否稳定、标题描述是否准确、结构化数据是否可信、站点地图是否更新和多语言关系是否正确。

## 元框架为什么适合 SEO

普通 SPA 的内容主要在浏览器运行后生成。搜索引擎虽然能执行部分 JavaScript，但内容站更推荐让首屏 HTML 就包含核心内容和 metadata。

Nuxt / Next 常见优势：

- 服务端渲染页面内容。
- 构建时生成静态页面。
- 每个路由生成独立 metadata。
- 支持动态详情页 metadata。
- 更容易生成 sitemap、canonical、Open Graph。

如果页面需要被搜索、分享、收录和预览，就应该认真处理 SEO。

## 每个页面至少要有什么

| 内容 | 作用 |
| --- | --- |
| title | 搜索结果标题、浏览器标签 |
| description | 搜索结果摘要和分享描述 |
| canonical | 告诉搜索引擎标准 URL |
| og:title | 社交分享标题 |
| og:description | 社交分享描述 |
| og:image | 分享卡片图片 |
| robots | 控制是否索引和跟随链接 |

不要所有页面共用同一个 title 和 description。

## Nuxt SEO

Nuxt 推荐使用 `useSeoMeta` 管理 SEO meta。

```vue
<script setup lang="ts">
useSeoMeta({
  title: 'Vue 学习路线',
  description: '从 HTML、JavaScript、TypeScript 到 Vue 3 项目实战的系统学习路线。',
  ogTitle: 'Vue 学习路线',
  ogDescription: '适合前端开发者的 Vue 3 系统学习路线。',
  ogImage: '/images/vue-roadmap-cover.png'
})
</script>
```

动态详情页：

```vue
<script setup lang="ts">
const route = useRoute()
const { data: article } = await useAsyncData(
  `article-${route.params.slug}`,
  () => $fetch(`/api/articles/${route.params.slug}`)
)

useSeoMeta({
  title: () => article.value?.title || '文章',
  description: () => article.value?.summary || '技术文章',
  ogTitle: () => article.value?.title,
  ogDescription: () => article.value?.summary,
  ogImage: () => article.value?.cover
})
</script>
```

公开内容的 SEO 数据应该来自内容模型或 CMS，而不是硬编码在组件里。

## Next SEO

Next App Router 支持静态 metadata 和 `generateMetadata`。

静态 metadata：

```tsx
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: '前端学习路线',
  description: '系统学习前端基础、Vue、工程化和项目实战。'
}
```

动态 metadata：

```tsx
import type { Metadata } from 'next'

type Props = {
  params: Promise<{ slug: string }>
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params
  const article = await getArticle(slug)

  return {
    title: article.title,
    description: article.summary,
    openGraph: {
      title: article.title,
      description: article.summary,
      images: [article.cover]
    }
  }
}
```

metadata 生成要考虑缓存和数据源。如果文章更新后 metadata 没更新，用户分享和搜索结果都会滞后。

## 结构化数据

结构化数据帮助搜索引擎理解页面内容。

常见类型：

- Article。
- Product。
- BreadcrumbList。
- FAQPage。
- Organization。
- WebSite。

示例：

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "Vue 3 响应式入门",
  "description": "用简单示例理解 ref、reactive 和 computed。",
  "author": {
    "@type": "Person",
    "name": "Vue 学习站"
  }
}
</script>
```

结构化数据必须和页面真实内容一致，不要堆无关关键词。

## sitemap 和 robots

内容站需要：

| 文件 | 作用 |
| --- | --- |
| `sitemap.xml` | 告诉搜索引擎有哪些可抓取页面 |
| `robots.txt` | 告诉爬虫哪些路径可抓取或不可抓取 |

常见策略：

- 公开文章进入 sitemap。
- 登录态页面不进 sitemap。
- 搜索结果页、筛选组合页谨慎索引。
- 预发环境禁止索引。
- 删除内容后更新 sitemap 和返回状态。

## 多语言 SEO

多语言站点要处理：

- 每个语言版本的 URL。
- `hreflang`。
- canonical。
- 默认语言。
- 翻译后的 title 和 description。
- 多语言 sitemap。

不要多个语言页面使用同一个 canonical，否则搜索引擎可能只保留一个语言版本。

## 实际项目问题

### 1. 所有详情页搜索标题都一样

**原因**

metadata 写在 layout，详情页没有根据数据生成。

**解决方案**

- Nuxt 用页面数据驱动 `useSeoMeta`。
- Next 用 `generateMetadata`。
- 内容模型里保存 title、description、cover。

### 2. 分享卡片没有图片

**原因**

缺少 Open Graph 图片，或图片不是可公开访问的绝对地址。

**解决方案**

- 配置 `og:image`。
- 使用稳定公开 URL。
- 图片尺寸符合平台要求。
- 发布后用平台调试工具刷新缓存。

### 3. 预发环境被搜索引擎收录

**原因**

预发环境没有 robots 限制，也没有鉴权。

**解决方案**

- 预发环境加访问控制。
- 设置 `robots: noindex, nofollow`。
- robots.txt 禁止抓取。

### 4. 页面收录了，但内容过期

**原因**

sitemap、缓存或 metadata 没有随内容更新。

**解决方案**

- 内容发布后刷新缓存。
- 动态页面设置合理 revalidate。
- 更新 sitemap。
- 删除页面返回 404 或 410。

## 最佳实践

- 公开页面单独设计 title、description 和分享图。
- SEO 数据进入内容模型，不要散落在组件里。
- 详情页 metadata 必须根据详情数据生成。
- 登录态和后台页面不要进入 sitemap。
- 预发环境禁止索引。
- 多语言页面处理 hreflang 和 canonical。
- 结构化数据必须与页面真实内容一致。

## 参考资料

- [Nuxt SEO and Meta](https://nuxt.com/docs/4.x/getting-started/seo-meta)
- [Nuxt useSeoMeta](https://nuxt.com/docs/4.x/api/composables/use-seo-meta)
- [Nuxt useServerSeoMeta](https://nuxt.com/docs/4.x/api/composables/use-server-seo-meta)
- [Next.js generateMetadata](https://nextjs.org/docs/app/api-reference/functions/generate-metadata)

## 下一步学习

继续学习 [国际化与多语言站点](/meta-frameworks/i18n)，把 SEO、路由和内容翻译结合起来。
