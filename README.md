# 程序员技术学习文档站

这是一个基于 VitePress 的中文技术文档站。内容以“先看图建立模型，再学核心知识，然后完成项目、练习和真实问题复盘”为主线，当前已覆盖前端基础、JavaScript、TypeScript、CSS、Vue、React、Nuxt/Next、Node.js、Java、Go、数据库、浏览器、工程化、DevOps 和 AI 工程。

## 项目结构

```text
docs/
├─ .vitepress/
│  ├─ config.ts
│  └─ theme/
│     ├─ index.ts
│     ├─ styles.css
│     └─ components/
├─ roadmap/
├─ frontend/
├─ css/
├─ javascript/
├─ typescript/
├─ vue/
├─ react/
├─ meta-frameworks/
├─ node/
├─ java/
├─ go/
├─ database/
├─ browser/
├─ engineering/
├─ devops/
├─ ai-engineering/
├─ projects/
├─ cheatsheets/
├─ technologies/
└─ contribute/
```

## 内容方向

文档以“详细、一看就懂、能解决实际项目问题”为目标。成熟模块通常包含：

- 学习导览和推荐顺序。
- 图解核心概念。
- 基础与进阶章节。
- 从零到项目落地。
- 专项练习和验收清单。
- 真实项目问题库与常见问题。

全站模块成熟度记录在 `docs/contribute/module-status.md`，问题库总入口是 `docs/projects/real-world-issues.md`。

## 启动

```bash
npm install
npm run docs:dev
```

开发服务器默认地址：`http://127.0.0.1:6173`。

## 文档检查

```bash
npm run docs:check
```

该命令检查内部路由是否存在、成熟模块是否有导览和快速排错页、核心页面是否有必备章节，以及配置路由和技术库入口是否指向真实页面。内容深度、侧边栏顺序和图示运行时渲染仍需要人工验收。

## 构建

```bash
npm run docs:build
```

生产构建后可以本地预览：

```bash
npm run docs:preview
```

## 主题说明

主题基于 VitePress 默认主题扩展，保留默认文档能力，并通过 `.vitepress/theme/styles.css` 定义清新、友好、专业的视觉变量。首页、学习路线、技术卡片和实践提示块使用自定义 Vue 组件实现。

## 样式约定

- 优先覆盖 VitePress 官方 CSS 变量。
- 业务样式使用明确 class，例如 `.custom-home__title`、`.learning-path__card`。
- 避免使用宽泛后代选择器污染默认主题或后续组件库样式。
- 固定尺寸视觉元素需要设置稳定宽高和不可压缩行为。

## 新增内容约定

扩展模块时优先补齐现有模块，不只增加零散页面。明显扩展后同步更新：

- 对应模块导览。
- `docs/.vitepress/config.ts` 导航和侧边栏。
- `docs/technologies/index.md` 技术库入口。
- `docs/technologies/expansion-plan.md` 扩展路线。
- `docs/contribute/module-status.md` 模块状态。

新增 Mermaid 图示后，除了运行构建，还要在浏览器中确认每张图生成 SVG、没有 `.mermaid-diagram__error`，并检查 390px 移动端没有页面级横向滚动。
