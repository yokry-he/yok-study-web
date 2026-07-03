# 包体积分析

## 适合谁看

适合项目已经能构建上线，但开始出现首屏慢、构建提示 chunk 过大、依赖越来越多的人：

- `vite build` 提示 chunk 超过 500 kB。
- 页面首屏加载慢，但不知道大在哪里。
- 组件库、图表库、编辑器、日期库全都进了首页。
- 改了懒加载，体积还是没降。
- 不知道该不该用 `manualChunks`。

包体积优化的第一步不是立刻拆包，而是先看清楚：谁大、为什么大、用户是否首屏必须下载。

## 先看三个指标

| 指标 | 说明 |
| --- | --- |
| 原始体积 | 构建产物未压缩大小，便于定位来源 |
| gzip / brotli 体积 | 用户真实传输更接近这个值 |
| 首屏必须加载体积 | 影响首屏速度的关键体积 |

不要只看单个 chunk 文件大小。一个大 chunk 如果是后台页面懒加载，不一定影响首页；很多小 chunk 如果造成网络瀑布，也可能很慢。

## 常见体积来源

| 来源 | 示例 |
| --- | --- |
| 组件库 | 全量引入 Element Plus、Ant Design Vue |
| 图表库 | echarts、antv、three |
| 编辑器 | monaco-editor、codemirror、富文本编辑器 |
| 日期工具 | moment、dayjs locale 全量 |
| 工具库 | lodash 全量引入 |
| 图标库 | 全量图标集合 |
| 多语言包 | 所有语言同时打进主包 |

要先确认真实来源，再决定按需引入、懒加载、拆 chunk 或替换方案。

## 分析工具

常见方式：

```bash
npm run build
```

配合可视化分析：

```bash
npx vite-bundle-visualizer
```

或在项目里接入 Rollup 可视化插件。

分析时看：

- 哪些依赖最大。
- 它们属于首页还是某个业务页。
- 是否重复打包。
- 是否包含不需要的语言包或主题包。
- 是否可被懒加载。

## 路由级懒加载

Vue Router 示例：

```ts
const routes = [
  {
    path: '/reports',
    component: () => import('@/pages/reports/index.vue')
  }
]
```

React 示例：

```tsx
const ReportsPage = lazy(() => import('./pages/reports'))
```

后台系统常见重页面：

- 报表。
- 可视化大屏。
- 富文本编辑。
- 低代码设计器。
- 权限配置矩阵。
- 文件预览。

这些页面不应该进入首页主包。

## 组件按需加载

组件库优化优先看官方按需方案。

原则：

- 使用组件库推荐的自动导入或插件。
- 样式按需加载。
- 图标按需导入。
- 不依赖内部 DOM 结构写样式。

不要为了省几 KB 破坏组件库升级能力。

## 动态导入重依赖

某些功能只在用户点击后使用，可以动态加载。

```ts
async function openEditor() {
  const { createEditor } = await import('@/features/editor/createEditor')
  createEditor()
}
```

适合：

- 富文本编辑器。
- Excel 导入导出。
- PDF 预览。
- 图表大屏。
- 地图。
- 代码编辑器。

不适合把每个小组件都动态导入。拆得太碎会增加请求和维护成本。

## manualChunks 怎么用

Vite 生产构建基于 Rollup，可以通过 `manualChunks` 控制拆包。

示例：

```ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vue: ['vue', 'vue-router', 'pinia'],
          charts: ['echarts']
        }
      }
    }
  }
})
```

使用原则：

- 先分析，再拆包。
- 稳定基础库可以单独拆。
- 重型业务依赖可以单独拆。
- 不要把所有 node_modules 粗暴拆成很多块。
- 拆完要看请求数量、缓存命中和首屏耗时。

## chunk 警告怎么判断

Vite 的 chunk size warning 是提醒，不是构建失败。

处理流程：

1. 看警告对应 chunk。
2. 用分析工具确认内容。
3. 判断是否首屏必须加载。
4. 如果不是首屏，确认是否已懒加载。
5. 如果是首屏，考虑按需引入、替换依赖或拆分。
6. 如果确认合理，再调整 `chunkSizeWarningLimit`。

不要第一反应就把警告阈值调大。阈值调大只是不提醒，体积没有变小。

## 实际项目问题

### 1. 后台首页加载了 ECharts

**原因**

报表路由和首页一起同步引入。

**解决方案**

- 报表页面路由懒加载。
- 图表组件内部动态导入 ECharts。
- 首页只加载指标摘要。

### 2. 图标库让主包很大

**原因**

全量导入图标。

**解决方案**

- 使用按需图标导入。
- 建立业务图标集合。
- 未使用图标不进入构建。

### 3. 拆包后首屏更慢

**原因**

拆成太多小 chunk，造成网络瀑布。

**解决方案**

- 合并稳定基础依赖。
- 保留路由级大块。
- 使用浏览器网络面板观察请求顺序。

### 4. 本地快，线上慢

**原因**

本地没有真实网络、CDN、gzip/brotli、缓存策略差异。

**解决方案**

- 用生产构建预览。
- 检查服务器压缩。
- 检查静态资源缓存。
- 使用 Lighthouse 或浏览器性能面板。

## 最佳实践

- 先分析包体积，再优化。
- 首屏不需要的重依赖必须懒加载。
- 组件库、图标库、语言包优先按需。
- `manualChunks` 要服务缓存和首屏，不要盲目拆。
- 构建产物通过 HTTP 服务验证。
- 包体积变化进入 PR 或发布检查。

## 参考资料

- [Vite Building for Production](https://vite.dev/guide/build)
- [Vite Build Options](https://vite.dev/config/build-options)
- [Rollup output.manualChunks](https://rollupjs.org/configuration-options/)

## 下一步学习

继续学习 [模块联邦与微前端](/engineering/module-federation)，理解多应用协作时如何拆分和共享运行时代码。
