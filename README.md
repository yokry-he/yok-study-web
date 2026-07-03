# Vue 前端学习文档站

这是一个基于 VitePress 的 Vue 前端学习文档站，第一阶段聚焦 Vue 3、前端基础、工程化、Vue Admin 实战和速查手册。

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
├─ javascript/
├─ vue/
├─ engineering/
├─ projects/
├─ cheatsheets/
└─ contribute/
```

## 内容方向

第一阶段文档以“详细、一看就懂、能解决实际项目问题”为目标。核心文档按概念、示例、真实问题、解决方案和最佳实践组织；项目问题沉淀在 `docs/projects/real-world-issues.md`。

## 启动

```bash
npm install
npm run docs:dev
```

## 构建

```bash
npm run docs:build
```

## 主题说明

主题基于 VitePress 默认主题扩展，保留默认文档能力，并通过 `.vitepress/theme/styles.css` 定义清新、友好、专业的视觉变量。首页、学习路线、技术卡片和实践提示块使用自定义 Vue 组件实现。

## 样式约定

- 优先覆盖 VitePress 官方 CSS 变量。
- 业务样式使用明确 class，例如 `.custom-home__title`、`.learning-path__card`。
- 避免使用宽泛后代选择器污染默认主题或后续组件库样式。
- 固定尺寸视觉元素需要设置稳定宽高和不可压缩行为。
