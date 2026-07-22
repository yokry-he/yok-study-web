# 渲染与性能

## 适合谁看

适合想理解页面为什么会白屏、卡顿、滚动不流畅、输入延迟的人。

框架性能优化最终都落到浏览器渲染上。Vue 的 `computed`、React 的 memo、虚拟列表、懒加载、代码分割，本质上都是为了减少不必要的计算、网络传输、DOM 更新和渲染成本。

## 浏览器如何把代码变成页面

浏览器渲染一个页面大致经历：

先用全景图建立阶段位置：HTML、CSS 和 JavaScript 形成结构与规则，浏览器计算可见节点的几何位置，生成绘制指令，最后由合成阶段把图层变成屏幕像素。

<DocFigure
  src="/images/browser/browser-rendering-pipeline.webp"
  alt="HTML CSS 和 JavaScript 依次经过 DOM 与 CSSOM、渲染树、布局、绘制和图层合成成为像素"
  caption="这是阶段地图；实际页面可能跳过某些工作或重复其中一段，准确成本要结合 Performance 录制判断。"
  :width="1440"
  :height="900"
/>

看图时不要形成“任何样式修改都会完整重跑所有阶段”的误解：改变尺寸通常影响 Layout，改变颜色通常影响 Paint，`transform` 等属性可能主要由 Composite 处理，但最终仍以浏览器实际生成的图层和录制结果为准。

```text
HTML
↓
DOM
↓
CSS
↓
CSSOM
↓
Render Tree
↓
Layout
↓
Paint
↓
Composite
```

含义：

| 阶段 | 做什么 | 常见性能问题 |
| --- | --- | --- |
| DOM | 解析 HTML 结构 | DOM 太大、节点过多 |
| CSSOM | 解析样式规则 | CSS 过大、选择器复杂 |
| Render Tree | 合并可见节点和样式 | 隐藏/显示频繁变化 |
| Layout | 计算元素几何位置 | 频繁读写布局 |
| Paint | 绘制文本、颜色、边框、阴影 | 大面积重绘 |
| Composite | 图层合成 | 图层过多、动画不合理 |

## 首屏性能

用户打开页面时，影响首屏的常见因素：

- HTML 响应慢。
- JavaScript 包过大。
- CSS 阻塞渲染。
- 字体、图片、接口请求过慢。
- 首屏依赖太多同步任务。
- 路由懒加载没有拆好。

优化方向：

| 问题 | 方案 |
| --- | --- |
| 包太大 | 路由级代码分割、按需引入 |
| 图片太大 | 压缩、响应式图片、懒加载 |
| 接口阻塞 | 骨架屏、并行请求、缓存 |
| 首屏 JS 执行久 | 拆分任务、减少初始化逻辑 |
| 第三方脚本慢 | 延迟加载、异步加载、评估必要性 |

## 交互性能

用户点击、输入、滚动时卡顿，通常是主线程太忙。

常见原因：

- 一次渲染太多 DOM。
- 输入时同步过滤大数组。
- 滚动事件里频繁计算布局。
- 动画修改 `width`、`height`、`top`、`left` 等布局属性。
- 组件树重渲染范围太大。

优化方向：

- 大列表使用虚拟滚动。
- 搜索输入使用防抖。
- 大计算放到 Web Worker 或拆分任务。
- 动画优先使用 `transform` 和 `opacity`。
- Vue 中合理使用 `computed`、`v-memo`、组件拆分。
- React 中减少不必要 state 提升和 Context 大范围更新。

## Layout Thrashing

Layout Thrashing 指的是代码反复交替读取布局和修改布局，导致浏览器不断强制重新计算布局。

不推荐：

```ts
items.forEach(item => {
  const height = item.offsetHeight
  item.style.height = `${height + 10}px`
})
```

更好的方式是先读后写：

```ts
const heights = items.map(item => item.offsetHeight)

items.forEach((item, index) => {
  item.style.height = `${heights[index] + 10}px`
})
```

在复杂页面里，这种差异会非常明显。

## 图片性能

图片是很多页面最大的资源。

建议：

- 使用合适尺寸，不要上传 4000px 图片再用 CSS 缩成 200px。
- 使用现代格式，例如 WebP 或 AVIF。
- 首屏关键图片优先加载。
- 非首屏图片使用懒加载。
- 设置宽高或 `aspect-ratio`，减少布局偏移。

示例：

```html
<img
  src="/banner.webp"
  width="1200"
  height="480"
  loading="lazy"
  alt="产品首页横幅"
/>
```

## 字体性能

字体加载也会影响首屏和布局稳定性。

建议：

- 控制字体文件数量。
- 中文字体谨慎自托管，文件可能很大。
- 使用 `font-display` 控制字体交换策略。
- 重要页面避免大量不同字重。

## DevTools 性能排查

### Network

看资源加载：

- 哪个资源最大。
- 哪个请求最慢。
- 是否命中缓存。
- 是否被阻塞。

### Performance

录制交互：

- Main 线程是否有长任务。
- 是否频繁 Layout。
- 点击后多久有响应。
- 哪些脚本执行最久。

下图中 `filterAndRender` 连续占用主线程 121 ms，输入事件虽然已经发生，却只能等长任务结束后处理。红色长任务是“主线程被占用”的证据，还需要展开调用栈才能判断成本来自业务计算、Vue 更新还是 Layout。

<DocFigure
  src="/images/browser/performance-long-task.webp"
  alt="浏览器 Performance 主线程轨道显示 121 毫秒长任务并推迟输入响应"
  caption="先定位超过 50 ms 的长任务，再向下展开调用栈；颜色提示耗时位置，不直接等于根因。"
  :width="1440"
  :height="900"
/>

不依赖图片的读取路径：Performance → 开始录制 → 执行一次卡顿操作 → 停止 → 在 Main 轨道寻找带红色角标或超过 50 ms 的 Task → 展开 Bottom-up 与 Call tree，记录最耗时函数及其调用来源。

### Lighthouse

适合做基础体检，但不要只追分数。分数只是信号，真正要解决的是用户路径里的具体问题。

## 实际项目问题

### 问题：表格 5000 行直接渲染导致页面卡死

**原因**

DOM 节点过多，渲染、布局和事件处理成本都变高。

**解决方案**

- 服务端分页。
- 虚拟滚动。
- 后端搜索和筛选。
- 导出功能走异步任务，不要把所有数据塞进表格。

### 问题：输入框搜索时每输入一个字都卡

**原因**

输入事件里同步触发大数组过滤、接口请求或复杂渲染。

**解决方案**

```ts
let timer: ReturnType<typeof setTimeout> | undefined

function handleInput(keyword: string) {
  clearTimeout(timer)

  timer = setTimeout(() => {
    search(keyword)
  }, 300)
}
```

Vue 项目也可以使用 VueUse 的 `useDebounceFn`。

### 问题：动画不流畅

**原因**

动画频繁修改布局属性。

不推荐：

```css
.panel {
  transition: left 0.2s;
}
```

推荐：

```css
.panel {
  transition: transform 0.2s;
}
```

`transform` 和 `opacity` 更容易交给合成线程处理。

## 最佳实践

- 首屏先减少必须加载和执行的内容。
- 大列表优先分页或虚拟滚动。
- 动画优先使用 `transform` 和 `opacity`。
- 不要在循环里交替读写布局。
- 性能优化必须用 DevTools 验证，不要只凭感觉。
- 组件性能问题要结合框架 DevTools 和浏览器 Performance 一起看。

## 参考资料

- [MDN: Critical rendering path](https://developer.mozilla.org/en-US/docs/Web/Performance/Guides/Critical_rendering_path)
- [web.dev: Understand the critical path](https://web.dev/learn/performance/understanding-the-critical-path)
- [web.dev: Rendering performance](https://web.dev/articles/rendering-performance)

## 下一步学习

继续学习 [常见问题](/browser/troubleshooting)。
