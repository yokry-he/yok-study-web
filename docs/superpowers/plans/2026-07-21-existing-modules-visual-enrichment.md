# 既有模块视觉讲解完善 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为现有技术文档建立可维护、可追踪、可验证的图片讲解体系，并按 P0、P1、P2 顺序给真正需要视觉说明的既有页面补充真实截图、教学图片或更清晰的结构图。

**Architecture:** 新增 `DocFigure` 组件、视觉资产登记和检查脚本，所有本地图片从 `docs/public/images` 提供。准确流程继续使用 Mermaid；视觉结果使用可复现截图；抽象心智模型才使用 built-in imagegen；外部图片只接受官方或明确许可来源。

**Tech Stack:** VitePress 1.6.4、Vue 3.5、Mermaid 11.16、built-in imagegen、Codex in-app Browser、HTML/CSS 复现页面、Node.js 检查脚本。

---

## 执行约束

- Task 1-3 在 Go 模块 Task 13 之前执行，为 Go 教学图片提供 `DocFigure` 和资产检查；Task 4-12 在 Go 模块最终验收后执行，避免同时修改主题、Go 内容和大量旧页面。
- 不自动提交、推送或部署；用户明确要求后再执行 Git 外部操作。
- 不要求每篇文档有位图；`diagram-sufficient` 是合法且推荐的审计结论。
- 不使用库存照片和纯装饰图片，不用生成图片承载必须准确阅读的代码或中文标签。
- `imagegen` 使用 built-in 模式；每个最终项目资产复制到仓库并登记最终 prompt。
- 网络图片必须先确认许可和来源；找不到明确许可时改为自行复现或 Mermaid。
- 所有截图先清除令牌、邮箱、连接串、本机绝对路径和真实用户数据。

## 第一阶段文件结构

```text
docs/.vitepress/theme/components/DocFigure.vue
docs/.vitepress/theme/index.ts
docs/.vitepress/theme/styles.css
docs/contribute/visual-audit.md
docs/contribute/visual-asset-register.md
docs/public/images/<module>/*.{png,jpg,jpeg,webp}
docs/public/visual-demos/frontend/index.html
docs/public/visual-demos/css/index.html
docs/public/visual-demos/vue-admin/index.html
docs/public/visual-demos/shared/demo.css
scripts/check-visual-assets.mjs
scripts/audit-doc-visuals.mjs
package.json
```

## Task 1: 用失败夹具定义视觉资产检查规则

**Files:**
- Create: `scripts/check-visual-assets.mjs`
- Create temporarily during test: `tmp/visual-check-fixtures/`
- Modify: `package.json`

- [ ] **Step 1: 写可从命令行指定根目录的检查器入口**

```js
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const rootDir = path.resolve(process.argv[2] ?? process.cwd())
const docsDir = path.join(rootDir, 'docs')
const publicDir = path.join(docsDir, 'public')
const registerPath = path.join(docsDir, 'contribute', 'visual-asset-register.md')

export function checkVisualAssets({ docsDir, publicDir, registerPath }) {
  const errors = []
  return errors
}

if (fileURLToPath(import.meta.url) === path.resolve(process.argv[1])) {
  const errors = checkVisualAssets({ docsDir, publicDir, registerPath })
  if (errors.length > 0) {
    console.error(errors.join('\n'))
    process.exit(1)
  }
  console.log('视觉资产检查通过')
}
```

- [ ] **Step 2: 创建四个失败夹具并确认检查器会失败**

夹具分别包含：不存在的 `/images/vue/missing.png`、空 `alt`、未登记图片、登记但磁盘不存在的图片。每个夹具使用真实 `docs/page.md`、`docs/public/images/...` 和 `visual-asset-register.md` 结构。

Run: `node scripts/check-visual-assets.mjs tmp/visual-check-fixtures/missing-file`

Expected: 非零退出并包含 `图片文件不存在：/images/vue/missing.png`。

- [ ] **Step 3: 实现 DocFigure 和裸 Markdown 图片提取**

```js
const componentRe = /<DocFigure\s+([\s\S]*?)\/>/g
const propRe = /(?:^|\s)(src|alt|caption|source-url)="([^"]*)"/g
const markdownImageRe = /!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g

function stripFencedCode(markdown) {
  return markdown.replace(/^(`{3,}|~{3,})[^\n]*\n[\s\S]*?^\1[ \t]*$/gm, '')
}

function parseProps(raw) {
  return Object.fromEntries([...raw.matchAll(propRe)].map(match => [match[1], match[2]]))
}
```

检查器遍历所有 Markdown，先调用 `stripFencedCode`，避免把教学代码片段误判成真实图片引用；发现 Markdown 裸图片直接报错；`DocFigure` 缺少 `src`、`alt` 或 `caption` 报错；本地 `src` 必须以 `/images/` 开头并解析到 `docs/public`。

- [ ] **Step 4: 实现大小、文件名和登记一致性检查**

允许扩展名：`.png`、`.jpg`、`.jpeg`、`.webp`；文件名匹配 `/^[a-z0-9]+(?:-[a-z0-9]+)*\.(png|jpe?g|webp)$/`；普通图片上限 500 KB，登记行包含 `large-approved` 时上限 1.5 MB；扫描图片目录找出未登记孤立文件。

登记表每个资产使用机器可解析行：

```markdown
<!-- asset: /images/go/go-api-request-journey.png | type: generated | license: project-generated | status: verified -->
```

- [ ] **Step 5: 跑通四个失败夹具和一个成功夹具**

Run:

```bash
for fixture in missing-file empty-alt unregistered missing-registered; do
  if node scripts/check-visual-assets.mjs "tmp/visual-check-fixtures/$fixture"; then
    echo "错误夹具意外通过：$fixture"
    exit 1
  fi
done
node scripts/check-visual-assets.mjs tmp/visual-check-fixtures/valid
```

Expected: 四个错误夹具按预期失败，`valid` 夹具输出 `视觉资产检查通过`。完成后删除 `tmp/visual-check-fixtures`。

- [ ] **Step 6: 接入现有检查命令**

```json
{
  "scripts": {
    "docs:check": "node scripts/check-docs.mjs && node scripts/check-visual-assets.mjs"
  }
}
```

Run: `npm run docs:check`

Expected: 当前无图片时也通过。

## Task 2: 实现可访问的 DocFigure 组件

**Files:**
- Create: `docs/.vitepress/theme/components/DocFigure.vue`
- Modify: `docs/.vitepress/theme/index.ts`
- Modify: `docs/.vitepress/theme/styles.css`

- [ ] **Step 1: 创建语义化组件**

```vue
<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  src: string
  alt: string
  caption: string
  sourceUrl?: string
  sourceLabel?: string
  width?: number
  height?: number
  zoomable?: boolean
}>(), {
  sourceLabel: '图片来源',
  zoomable: true
})

const open = ref(false)
const failed = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const dialog = ref<HTMLElement | null>(null)
const closeButton = ref<HTMLButtonElement | null>(null)
let previousBodyOverflow = ''

function close() {
  open.value = false
}

function onDialogKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    event.preventDefault()
    close()
    return
  }
  if (event.key !== 'Tab' || !dialog.value) return
  const focusable = [...dialog.value.querySelectorAll<HTMLElement>('button, [href], [tabindex]:not([tabindex="-1"])')]
  if (focusable.length === 0) return
  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

watch(open, async value => {
  if (typeof document === 'undefined') return
  if (value) {
    previousBodyOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    await nextTick()
    closeButton.value?.focus()
  } else {
    document.body.style.overflow = previousBodyOverflow
    trigger.value?.focus()
  }
})

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') document.body.style.overflow = previousBodyOverflow
})
</script>

<template>
  <figure class="doc-figure">
    <button
      v-if="zoomable && !failed"
      ref="trigger"
      class="doc-figure__trigger"
      type="button"
      :aria-label="`放大查看：${alt}`"
      @click="open = true"
    >
      <img class="doc-figure__image" :src="src" :alt="alt" :width="width" :height="height" loading="lazy" decoding="async" @error="failed = true">
    </button>
    <img v-else-if="!failed" class="doc-figure__image" :src="src" :alt="alt" :width="width" :height="height" loading="lazy" decoding="async">
    <p v-else class="doc-figure__error" role="status">图片加载失败，请根据图注继续阅读。</p>
    <figcaption class="doc-figure__caption">
      {{ caption }}
      <a v-if="sourceUrl" class="doc-figure__source" :href="sourceUrl" target="_blank" rel="noreferrer">{{ sourceLabel }}</a>
    </figcaption>
    <Teleport to="body">
      <div v-if="open" ref="dialog" class="doc-figure-lightbox" role="dialog" aria-modal="true" :aria-label="alt" @click.self="close" @keydown="onDialogKeydown">
        <button ref="closeButton" class="doc-figure-lightbox__close" type="button" aria-label="关闭图片预览" @click="close">×</button>
        <img class="doc-figure-lightbox__image" :src="src" :alt="alt">
      </div>
    </Teleport>
  </figure>
</template>
```

- [ ] **Step 2: 注册全局组件**

```ts
import DocFigure from './components/DocFigure.vue'

app.component('DocFigure', DocFigure)
```

- [ ] **Step 3: 添加明确业务 class 样式**

样式要求：8px 以内圆角、稳定边框、浅色与深色主题变量、触发按钮完整焦点环、图片 `max-width: 100%`、lightbox 固定全屏且图片使用 `max-width: min(94vw, 1440px)` 和 `max-height: 90vh`。打开时把焦点移到关闭按钮并锁定背景滚动，关闭后恢复触发按钮焦点；Tab 焦点保持在对话框内。移动端关闭按钮保持 44x44 可点击区域，不使用宽泛后代选择器。

- [ ] **Step 4: 用一张 1x1 测试图片验证组件后删除夹具**

临时创建 `/images/test/doc-figure-test.png` 和测试页面，验证正常加载、失败状态、Enter/Space 打开、Escape 关闭、点击遮罩关闭和焦点可见；完成后删除测试图片和页面，资产检查无孤立文件。

- [ ] **Step 5: 运行构建**

Run: `npm run docs:check && npm run docs:build`

Expected: 通过，无 SSR `window is not defined` 错误。

## Task 3: 建立全站视觉审计和资产登记

**Files:**
- Create: `scripts/audit-doc-visuals.mjs`
- Create: `docs/contribute/visual-audit.md`
- Create: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 实现审计脚本**

脚本遍历 Markdown，输出 route、模块、Mermaid 数、DocFigure 数、代码块数和初始优先级：

```js
const p0Modules = new Set(['frontend', 'css', 'vue', 'browser'])
const p1Modules = new Set([
  'javascript', 'typescript', 'react', 'meta-frameworks', 'node', 'java',
  'go', 'database', 'engineering', 'devops', 'ai-engineering'
])

function priorityFor(moduleName, relativePath) {
  if (p0Modules.has(moduleName)) return 'P0'
  if (p1Modules.has(moduleName)) return 'P1'
  if (moduleName === 'projects') return 'P2'
  return 'P2'
}
```

脚本排除 `docs/superpowers/` 和生成目标 `docs/contribute/visual-audit.md`；其他正式内容初始状态统一为 `review-required`，每次按 route 排序，保证 diff 稳定。

- [ ] **Step 2: 生成并人工复核基线数字**

Run: `node scripts/audit-doc-visuals.mjs > docs/contribute/visual-audit.md`

Expected: 包含全部 Markdown 路由；基线说明记录 2026-07-21 时 524 篇正式内容文档、327 个包含 Mermaid 的页面和 0 个旧位图引用。规格与计划页面不计入正式内容基线，但仍接受资产检查。

- [ ] **Step 3: 建立登记格式**

```markdown
## `/images/vue/admin-list-filtered-result.png`

<!-- asset: /images/vue/admin-list-filtered-result.png | type: live-screenshot | license: project-owned | status: verified -->

- 使用页面：`/vue/admin-list-search-table`
- 教学目的：对比筛选前后表格、分页和空状态的联动。
- Alt：用户列表筛选后仅保留启用状态，分页总数同步减少。
- 图注：筛选条件、列表请求参数和分页总数必须来自同一查询状态。
- 复现：`/visual-demos/vue-admin/?scene=list-filtered`，视口 `1440x900`，设备像素比 2。
- 来源：项目自建演示页面。
```

- [ ] **Step 4: 在贡献文档说明媒介选择规则**

明确 `diagram-sufficient`、`needs-live-screenshot`、`needs-annotated-screenshot`、`needs-generated-visual`、`needs-official-source`、`needs-mermaid-refactor` 六种结果及完成定义。

- [ ] **Step 5: 运行检查**

Run: `npm run docs:check && git diff --check`

Expected: 通过。

## Task 4: 创建可复现的前端与 Vue Admin 视觉演示页

**Files:**
- Create: `docs/public/visual-demos/shared/demo.css`
- Create: `docs/public/visual-demos/frontend/index.html`
- Create: `docs/public/visual-demos/css/index.html`
- Create: `docs/public/visual-demos/vue-admin/index.html`

- [ ] **Step 1: 建立演示页共同规范**

`demo.css` 使用项目现有清新配色，但包含白色、炭灰、绿色、青色和琥珀色，不使用单一色系。固定 `box-sizing`、字体回退、焦点环、状态色、表格和表单尺寸；所有场景由 URL `?scene=` 选择，页面只展示目标画面，不显示使用说明。

- [ ] **Step 2: 创建前端基础 6 个场景**

场景名固定为：

```text
semantic-article
form-valid
form-invalid
responsive-desktop
responsive-mobile
accessible-focus
```

每个场景使用真实 HTML 元素；无障碍焦点场景必须通过 `:focus-visible` 呈现；表单错误使用文字和 `aria-describedby`，不只依赖红色。

- [ ] **Step 3: 创建 CSS 8 个场景**

```text
box-model
margin-collapse
flex-main-cross-axis
flex-overflow
grid-template-areas
grid-minmax
responsive-container
reduced-motion
```

场景在页面中显示实际渲染结果，并用可维护 HTML 标签标出 padding、border、gap、轨道和断点；不把 CSS 代码绘制进图片。

- [ ] **Step 4: 创建 Vue Admin 10 个场景**

```text
dashboard
list-default
list-filtered
list-empty
form-create
form-validation
detail-audit
approval-pending
notification-unread
permission-denied
```

所有中文文本由 HTML 真实渲染；表格、抽屉、弹窗、状态标签、按钮和图表使用稳定尺寸；移动场景使用紧凑顶部菜单，不堆叠完整桌面侧栏。

- [ ] **Step 5: 浏览器逐场景验证**

在 1440x900 和 390x844 验证所有场景无横向溢出、文字不重叠、焦点可见、图表和表格稳定；记录场景 URL 供截图复现。

## Task 5: 给前端基础和 CSS 页面补真实截图

**Files:**
- Create: `docs/public/images/frontend/semantic-article.webp`
- Create: `docs/public/images/frontend/form-valid.webp`
- Create: `docs/public/images/frontend/form-invalid.webp`
- Create: `docs/public/images/frontend/responsive-desktop.webp`
- Create: `docs/public/images/frontend/responsive-mobile.webp`
- Create: `docs/public/images/frontend/accessible-focus.webp`
- Create: `docs/public/images/css/box-model.webp`
- Create: `docs/public/images/css/margin-collapse.webp`
- Create: `docs/public/images/css/flex-main-cross-axis.webp`
- Create: `docs/public/images/css/flex-overflow.webp`
- Create: `docs/public/images/css/grid-template-areas.webp`
- Create: `docs/public/images/css/grid-minmax.webp`
- Create: `docs/public/images/css/responsive-container.webp`
- Create: `docs/public/images/css/reduced-motion.webp`
- Modify: `docs/frontend/html-semantics.md`
- Modify: `docs/frontend/forms-media-accessibility.md`
- Modify: `docs/frontend/html-css.md`
- Modify: `docs/frontend/project-from-zero.md`
- Modify: `docs/css/box-model-layout.md`
- Modify: `docs/css/flex-grid.md`
- Modify: `docs/css/responsive.md`
- Modify: `docs/css/accessibility.md`
- Modify: `docs/css/animation-transition.md`
- Modify: `docs/css/project-from-zero.md`
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 截取前端基础场景**

用 in-app Browser 打开本地 6173 端口的演示页，以 2x 像素密度输出六张临时 PNG，截图只包含演示区域，不包含浏览器地址栏。使用 `cwebp -quiet -q 88 input.png -o docs/public/images/frontend/<scene>.webp` 转换，视觉核对后删除临时 PNG。

- [ ] **Step 2: 截取 CSS 场景**

输出八张 WebP；盒模型、Flex 和 Grid 同时截桌面目标场景；响应式章节用两张同数据、不同视口的图片并排解释，不用缩放桌面图伪装移动端。

- [ ] **Step 3: 使用 DocFigure 插入对应章节**

```md
<DocFigure
  src="/images/css/flex-main-cross-axis.webp"
  alt="三个项目沿主轴水平排列，交叉轴控制垂直对齐"
  caption="先由 flex-direction 决定主轴，再分别使用 justify-content 和 align-items 控制两个方向。"
  :width="1440"
  :height="900"
/>
```

每张图前说明观察任务，图后解释现象和对应 CSS；不连续堆放图片。

- [ ] **Step 4: 更新审计和登记**

对应页面从 `review-required` 改为 `verified`；不需要位图的 frontend/css 页面标记 `diagram-sufficient` 并写一句理由。

- [ ] **Step 5: 检查资源和页面**

Run: `npm run docs:check && npm run docs:build`

Expected: 所有图片低于大小上限、无孤立资产、页面构建成功。

## Task 6: 给 Vue 与 Vue Admin 既有功能补界面讲解

**Files:**
- Create: `docs/public/images/vue/admin-dashboard.webp`
- Create: `docs/public/images/vue/admin-list-default.webp`
- Create: `docs/public/images/vue/admin-list-filtered.webp`
- Create: `docs/public/images/vue/admin-list-empty.webp`
- Create: `docs/public/images/vue/admin-form-create.webp`
- Create: `docs/public/images/vue/admin-form-validation.webp`
- Create: `docs/public/images/vue/admin-detail-audit.webp`
- Create: `docs/public/images/vue/admin-approval-pending.webp`
- Create: `docs/public/images/vue/admin-notification-unread.webp`
- Create: `docs/public/images/vue/admin-permission-denied.webp`
- Create: `docs/public/images/vue/admin-file-upload-progress.webp`
- Create: `docs/public/images/vue/admin-permission-loading.webp`
- Modify: `docs/vue/admin-dashboard-analytics.md`
- Modify: `docs/vue/admin-list-search-table.md`
- Modify: `docs/vue/admin-form-modal-crud.md`
- Modify: `docs/vue/admin-detail-status-audit.md`
- Modify: `docs/vue/admin-approval-workflow.md`
- Modify: `docs/vue/admin-notification-center.md`
- Modify: `docs/vue/admin-permission-route-flow.md`
- Modify: `docs/vue/admin-file-import-export.md`
- Modify: `docs/vue/project-from-zero.md`
- Modify: `docs/vue/forms.md`
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 截取 10 个 Vue Admin 场景**

场景截图必须显示与正文一致的数据状态：筛选条件与总数一致、表单错误贴近字段、详情审计按时间排序、审批按钮受状态控制、未读数与列表一致、无权限状态不显示业务数据。

- [ ] **Step 2: 为每个闭环页面建立“观察 -> 原理 -> 代码”结构**

每张截图前给 2 到 4 个观察点，截图后连接到已有 Mermaid、请求参数、响应字段或 Vue 组件代码。界面图片不能代替状态机或权限判断说明。

- [ ] **Step 3: 为上传和权限页面补专用状态**

在演示页增加 `file-upload-progress` 和 `permission-loading` 两个状态，再截取并插入对应文档；加载状态不得闪现无权限页，上传进度与异步任务状态分开。

- [ ] **Step 4: 为 Vue 核心页面选择 Mermaid 或真实截图**

`reactivity.md`、`lifecycle.md`、`composition-api.md` 保持 Mermaid；`forms.md` 和 `project-from-zero.md` 使用真实界面；其余页面根据是否需要观察渲染结果标记审计结论。

- [ ] **Step 5: 验证桌面和移动版**

至少检查 `/vue/admin-list-search-table`、`/vue/admin-form-modal-crud`、`/vue/admin-dashboard-analytics`、`/vue/admin-permission-route-flow` 的 1440 和 390 视口，图片、图注和 lightbox 无重叠。

## Task 7: 给浏览器与调试章节补证据型图片

**Files:**
- Create: `docs/public/images/browser/network-request-headers.webp`
- Create: `docs/public/images/browser/network-request-timing.webp`
- Create: `docs/public/images/browser/cache-memory-disk.webp`
- Create: `docs/public/images/browser/storage-local-session.webp`
- Create: `docs/public/images/browser/cors-preflight.webp`
- Create: `docs/public/images/browser/performance-long-task.webp`
- Create: `docs/public/images/browser/automation-locator-failure.webp`
- Modify: `docs/browser/http-request.md`
- Modify: `docs/browser/cache.md`
- Modify: `docs/browser/storage.md`
- Modify: `docs/browser/rendering-performance.md`
- Modify: `docs/browser/browser-automation-debugging.md`
- Modify: `docs/browser/cors-auth.md`
- Modify: `docs/browser/project-from-zero.md`
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 先判断是否能自行复现**

请求瀑布、Cache-Control、Storage、CORS 和性能问题优先用本地演示页与浏览器工具自行复现。截图前统一清空历史记录并使用固定测试数据，保证图中状态与正文一致。

- [ ] **Step 2: 只在无法自行复现时搜索官方图片**

搜索限定 Chrome Developers、MDN、Firefox Source Docs 等官方域名；打开原始页面核对版本、上下文和许可。Google Developers 图片只有在页面许可允许时使用，并在 DocFigure 的 `sourceUrl` 和登记中记录；不引用搜索结果缩略图。

- [ ] **Step 3: 建立 7 张证据图**

目标主题固定为：请求 Headers、请求 Timing、memory/disk cache 区别、local/session storage、CORS 预检、Performance 长任务、自动化定位器失败证据。每张图用编号图注解释先看哪个区域，不修改截图中的事实数据。

- [ ] **Step 4: 更新正文排查步骤**

每张工具截图后写出相同证据的文本读取路径和命令替代方案，保证屏幕阅读器用户不依赖图片完成排查。

- [ ] **Step 5: 验证来源和图片请求**

Run: `npm run docs:check && npm run docs:build`

Expected: 所有外部来源有许可证字段；本地图片请求无 404。

## Task 8: 给已有可运行项目补最终状态和运行证据

**Files:**
- Create: `docs/public/images/java/java-admin-api-ready.webp`
- Create: `docs/public/images/java/java-admin-version-conflict.webp`
- Create: `docs/public/images/go/go-task-api-ready.webp`
- Create: `docs/public/images/go/go-task-version-conflict.webp`
- Modify: `docs/java/spring-boot-project-from-zero.md`
- Modify: `docs/go/http-api-project-from-zero.md`
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 先运行每个真实示例**

只有仓库中存在且测试通过的项目才能生成“可运行项目”截图。本批只处理 `examples/java-admin-api` 和完成后的 `examples/go-task-api`，分别截 readiness 和旧版本更新冲突。JavaScript、React 和 Node.js 项目页在没有对应仓库示例前不能标记为真实运行截图。

- [ ] **Step 2: 截图不替代代码输出**

每个项目页使用 2 张高价值图片；curl、测试日志和错误码仍保留可复制文本。终端截图必须使用演示账号和脱敏路径。

- [ ] **Step 3: 为 Java 和 Go 后端登记复现命令**

Java 使用 `docker compose up --build` 和用户接口流程；Go 使用任务接口、版本冲突和优雅关闭流程。登记具体命令、镜像版本和截图时的提交状态。

- [ ] **Step 4: 浏览器检查 Java 和 Go 项目页**

图片必须展示正文承诺的功能，不使用通用后台模板冒充项目成品。

## Task 9: 使用 imagegen 补难以截图的抽象心智模型

**Files:**
- Create: `docs/public/images/browser/browser-rendering-pipeline.webp`
- Create: `docs/public/images/java/jvm-memory-regions.webp`
- Create: `docs/public/images/node/event-loop-io-workshop.webp`
- Create: `docs/public/images/ai-engineering/rag-retrieval-journey.webp`
- Modify: `docs/browser/rendering-performance.md`
- Modify: `docs/java/jvm-memory-gc.md`
- Modify: `docs/node/runtime-event-loop.md`
- Modify: `docs/ai-engineering/rag.md`
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/contribute/visual-asset-register.md`

- [ ] **Step 1: 确认四个首批 `needs-generated-visual` 页面**

首批固定为浏览器渲染流水线、Java JVM 内存区域、Node 事件循环与 I/O、AI RAG 检索链。这四张图片只建立整体空间心智模型，准确顺序和边界仍由相邻 Mermaid 负责。Go goroutine/channel 教学图片已由 Go 模块计划创建，本计划只在审计中复核，不重复生成。

- [ ] **Step 2: 每个资产单独调用 built-in imagegen**

共同 prompt 约束：

```text
Use case: scientific-educational
Asset type: Chinese programming documentation learning visual
Style/medium: clean flat technical illustration, crisp geometric forms, professional educational quality
Composition/framing: wide 16:9 with clear visual hierarchy and generous padding
Constraints: no text, no letters, no code, no logos, no watermark; use visual grouping only; the document caption will carry exact terminology
Avoid: dark background, decorative gradients, fantasy elements, illegible micro-details, fake product UI
```

每张图在 `Primary request` 中只描述一个心智模型，不通过一次调用批量生成不同主题。

四次调用分别追加以下 Primary request：

```text
Browser: show HTML, CSS, and JavaScript entering a browser engine, becoming two structured trees, passing through layout, paint, layered composition, and finally visible pixels; normal left-to-right flow, no labels.

JVM: show several isolated thread stacks beside one shared object heap, a separate class metadata area, and a garbage collector tracing reachable objects from roots; make shared versus thread-private space visually unmistakable, no labels.

Node.js: show one event-loop coordinator dispatching file and network operations away from the main loop, with completed callbacks returning through a queue; show that long CPU work blocks the coordinator while external I/O does not, no labels.

RAG: show source documents split into chunks and stored in a semantic index; a user question retrieves a small relevant subset, which joins the question before an answer is produced; clearly separate ingestion from query time, no labels.
```

- [ ] **Step 3: 检查技术隐喻边界**

图注必须明确哪些元素是类比，准确执行顺序仍由 Mermaid、代码和正文定义。发现错误对象数量、错误方向或错误分组时重新生成，不能靠正文为错误图片辩解。

- [ ] **Step 4: 保存和登记最终输出**

从 `$CODEX_HOME/generated_images` 复制最终选中 PNG 到 `tmp/imagegen`，使用 `cwebp -quiet -q 86` 转换为文件列表中的 WebP；人工核对后删除临时 PNG。记录最终 prompt、built-in 模式、日期、人工核对项和未采用变体。项目引用不能指向 `$CODEX_HOME`。

## Task 10: 完成 P1 稳定技术模块审计与补图

**Files:**
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/javascript/event-loop.md`
- Modify: `docs/javascript/dom-events.md`
- Modify: `docs/javascript/task-board-project.md`
- Modify: `docs/typescript/narrowing-guards.md`
- Modify: `docs/typescript/type-boundary-project.md`
- Modify: `docs/typescript/tsconfig-engineering.md`
- Modify: `docs/react/effects.md`
- Modify: `docs/react/performance.md`
- Modify: `docs/react/project-admin.md`
- Modify: `docs/meta-frameworks/routing-data.md`
- Modify: `docs/meta-frameworks/server-auth.md`
- Modify: `docs/meta-frameworks/project-from-zero.md`
- Modify: `docs/node/permission-api-project.md`
- Modify: `docs/node/cache-queue-project.md`
- Modify: `docs/java/spring-boot-project-from-zero.md`
- Modify: `docs/java/testing-deployment.md`
- Modify: `docs/database/indexes.md`
- Modify: `docs/database/transactions.md`
- Modify: `docs/database/project-practice.md`
- Modify: `docs/engineering/bundle-analysis.md`
- Modify: `docs/engineering/module-federation.md`
- Modify: `docs/engineering/project-from-zero.md`
- Modify: `docs/devops/docker.md`
- Modify: `docs/devops/observability.md`
- Modify: `docs/devops/deployment-strategy.md`
- Modify: `docs/ai-engineering/evaluation.md`
- Modify: `docs/ai-engineering/doc-qa-project.md`
- Create: `docs/public/images/javascript/event-loop-devtools.webp`
- Create: `docs/public/images/javascript/dom-event-path.webp`
- Create: `docs/public/images/javascript/task-board-states.webp`
- Create: `docs/public/images/typescript/narrowing-control-flow.webp`
- Create: `docs/public/images/typescript/type-boundary-error-state.webp`
- Create: `docs/public/images/typescript/tsconfig-trace.webp`
- Create: `docs/public/images/react/effect-request-race.webp`
- Create: `docs/public/images/react/performance-profiler.webp`
- Create: `docs/public/images/react/admin-states.webp`
- Create: `docs/public/images/meta-frameworks/route-data-boundaries.webp`
- Create: `docs/public/images/meta-frameworks/server-auth-redirect.webp`
- Create: `docs/public/images/meta-frameworks/course-platform.webp`
- Create: `docs/public/images/node/permission-api-response.webp`
- Create: `docs/public/images/node/cache-queue-dashboard.webp`
- Create: `docs/public/images/java/testcontainers-run.webp`
- Create: `docs/public/images/database/index-explain.webp`
- Create: `docs/public/images/database/transaction-lock-wait.webp`
- Create: `docs/public/images/database/permission-project.webp`
- Create: `docs/public/images/engineering/bundle-analysis.webp`
- Create: `docs/public/images/engineering/module-federation-runtime.webp`
- Create: `docs/public/images/engineering/project-pipeline.webp`
- Create: `docs/public/images/devops/docker-container-state.webp`
- Create: `docs/public/images/devops/observability-dashboard.webp`
- Create: `docs/public/images/devops/deployment-canary.webp`
- Create: `docs/public/images/ai-engineering/evaluation-report.webp`
- Create: `docs/public/images/ai-engineering/doc-qa-citations.webp`

- [ ] **Step 1: 每个模块先选三个最高认知成本页面**

固定评估顺序：运行结果必须观察、文字中出现空间/时间关系、问题库需要证据截图、项目页缺最终状态、现有 Mermaid 在 390px 不可读。每个模块先完成三个页面，再评估是否继续。

- [ ] **Step 2: 使用模块媒介规则**

- JavaScript/TypeScript：DevTools 值、DOM、事件、source map 使用真实截图，类型关系使用 Mermaid。
- React/Nuxt/Next：页面状态和 hydration 差异使用真实界面，渲染流程使用 Mermaid。
- Node/Java/Go：运行证据和诊断工具使用截图，并发、内存、事务使用 Mermaid或经过审计的教学图。
- Database/DevOps：执行计划、锁、监控和发布状态使用脱敏截图，拓扑与事务使用 Mermaid。
- AI Engineering：检索结果、评测和结构化输出使用可复现界面，模型内部原理只做明确标注的概念图。

- [ ] **Step 3: 每完成一个模块立即验证**

Run: `npm run docs:check && npm run docs:build`

Expected: 通过；该模块审计页面没有 `review-required` 遗留在已声明完成的范围内。

## Task 11: 按独特业务状态处理 P2 案例库

**Files:**
- Modify: `docs/contribute/visual-audit.md`
- Modify: `docs/projects/vue-admin.md`
- Modify: `docs/projects/approval-workflow-case.md`
- Modify: `docs/projects/analytics-dashboard-case.md`
- Modify: `docs/projects/workflow-builder-case.md`
- Modify: `docs/projects/file-center-case.md`
- Modify: `docs/projects/notification-center-case.md`
- Modify: `docs/projects/multi-tenant-permission-case.md`
- Modify: `docs/projects/finance-reconciliation-case.md`
- Modify: `docs/projects/risk-control-center-case.md`
- Modify: `docs/projects/disaster-recovery-case.md`
- Create: `docs/public/images/projects/vue-admin-overview.webp`
- Create: `docs/public/images/projects/approval-workflow-state.webp`
- Create: `docs/public/images/projects/analytics-dashboard-overview.webp`
- Create: `docs/public/images/projects/workflow-builder-canvas.webp`
- Create: `docs/public/images/projects/file-center-task.webp`
- Create: `docs/public/images/projects/notification-center-unread.webp`
- Create: `docs/public/images/projects/multi-tenant-permission-scope.webp`
- Create: `docs/public/images/projects/finance-reconciliation-exception.webp`
- Create: `docs/public/images/projects/risk-control-case-review.webp`
- Create: `docs/public/images/projects/disaster-recovery-switch.webp`

- [ ] **Step 1: 将案例按视觉模式分组**

分组固定为：列表筛选、表单审批、详情审计、数据看板、配置画布、文件任务、消息通知、权限数据、财务对账、监控处置。高度相似案例不复制同一张图，而是在共同模式图后补该领域独有状态。

- [ ] **Step 2: 完成首批十个有独特状态的案例**

首批固定为文件列表中的十个案例，分别覆盖基础后台、审批、看板、配置画布、文件任务、通知、多租户权限、财务对账、风控和灾备。后续批次仍按“状态机不少于 4 个状态、角色视图明显不同、存在双人复核或审计链、图表直接影响决策、存在异常处置闭环”排序，并为每个批次新增独立计划，不在本计划中隐式扩大写集。

- [ ] **Step 3: 禁止通用后台截图冒充业务证据**

每张案例图至少显示一个该案例独有字段、状态或操作；图注说明它如何连接到正文中的数据模型和状态机。

- [ ] **Step 4: 每批更新审计统计**

记录本批完成页面、`diagram-sufficient` 页面、图片数量、生成图片数量、外部来源数量和剩余高优先级页面。

## Task 12: 最终视觉与无障碍验收

**Files:**
- Verify only; only edit files when a specific defect is reproduced.

- [ ] **Step 1: 运行自动检查**

Run:

```bash
npm run docs:check
npm run docs:build
git diff --check
```

Expected: 全部通过；无丢失图片、空 alt、未登记资产、超大资源或裸 Markdown 图片。

- [ ] **Step 2: 检查图片请求和页面宽度**

在 1440x900 和 390x844 抽查每个已完成模块至少两页；断言图片 HTTP 状态 200、`naturalWidth > 0`、页面无横向溢出、图片只在自身容器缩放。

- [ ] **Step 3: 检查 lightbox 键盘操作**

使用 Tab 聚焦图片，Enter 打开，Escape 关闭；关闭按钮有可见焦点和中文可访问名称；打开时背景不能意外触发链接。

- [ ] **Step 4: 检查内容准确性和许可**

逐项核对登记：真实截图能按命令复现，生成图不冒充事实，外部图有原始页面与许可，图片中无敏感数据。

- [ ] **Step 5: 检查性能**

普通图片小于 500 KB，批准的大图小于 1.5 MB；首屏外使用 lazy loading；图片有 width/height，滚动时无明显布局跳动。

- [ ] **Step 6: 输出批次报告但不提交**

报告已审计页面数、各种审计结论数量、本地图片数、总字节、生成图片清单、外部来源清单、桌面/移动端结果和下一批 P1/P2 页面。只有用户明确要求时才提交、推送或部署。
