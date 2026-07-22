# 动画与过渡

## 适合谁看

适合已经会写基础样式，但在真实项目里遇到这些问题的人：

- hover 动效突兀，页面显得生硬。
- 弹窗、下拉、抽屉出现和消失不自然。
- 动画导致页面卡顿或移动端掉帧。
- 不知道 `transition` 和 `animation` 怎么选。
- 不知道如何照顾开启“减少动态效果”的用户。

动画不是为了炫技，而是为了帮助用户理解状态变化。好的动效应该短、轻、明确，不抢内容本身的注意力。

## transition 和 animation 怎么选

| 能力 | 适合场景 |
| --- | --- |
| `transition` | 两个状态之间的平滑变化，例如 hover、focus、展开收起 |
| `animation` | 更完整的时间线，例如 loading、骨架屏、循环提示、入场动画 |

简单判断：

```text
状态 A -> 状态 B：优先 transition
多个关键帧或循环：使用 animation + @keyframes
```

不要为了简单 hover 写复杂 `@keyframes`，也不要用一堆 transition 拼复杂时间线。

## transition 基础

```css
.action-button {
  background: #14b89a;
  transform: translateY(0);
  transition:
    background-color 160ms ease,
    transform 160ms ease;
}

.action-button:hover {
  background: #0f9f86;
  transform: translateY(-1px);
}
```

建议只过渡明确属性，不要写：

```css
transition: all 200ms ease;
```

`all` 会让不该动的属性也进入过渡，增加排查成本。

## animation 基础

```css
@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(8px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.toast {
  animation: fade-in-up 180ms ease-out both;
}
```

`both` 表示动画开始前和结束后都应用关键帧状态，常用于入场动画。

## 优先动这些属性

更适合动画：

- `opacity`
- `transform`
- `filter` 少量使用

谨慎动画：

- `width`
- `height`
- `top`
- `left`
- `margin`
- `padding`

原因是尺寸和位置变化更容易触发布局计算。多数 UI 动效可以通过 `transform` 模拟位移、缩放。

## 常用交互动效

### 按钮反馈

```css
.toolbar-action {
  transform: translateY(0);
  transition:
    transform 120ms ease,
    box-shadow 120ms ease;
}

.toolbar-action:hover {
  transform: translateY(-1px);
  box-shadow: 0 8px 20px rgb(15 23 42 / 10%);
}

.toolbar-action:active {
  transform: translateY(0);
}
```

按钮反馈应该很轻，不要让按钮跳动太大。

### 弹窗入场

```css
@keyframes dialog-enter {
  from {
    opacity: 0;
    transform: translateY(10px) scale(0.98);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.dialog-panel {
  animation: dialog-enter 180ms ease-out both;
}
```

弹窗动效要短。用户目标是处理内容，不是看动画。

### 骨架屏

```css
@keyframes skeleton-shimmer {
  from {
    background-position: 100% 0;
  }

  to {
    background-position: -100% 0;
  }
}

.skeleton-line {
  background: linear-gradient(90deg, #edf2f7 25%, #f8fafc 37%, #edf2f7 63%);
  background-size: 400% 100%;
  animation: skeleton-shimmer 1.2s ease-in-out infinite;
}
```

骨架屏不要过亮、过快，否则会造成视觉干扰。

## 尊重减少动态效果

用户可能在系统里开启减少动态效果。CSS 可以通过 `prefers-reduced-motion` 检测。

下图对比了正常动效和减少动态效果两种系统偏好。减少动态并不等于删除所有反馈，而是移除大幅位移、持续循环和容易引起不适的过渡。

<DocFigure
  src="/images/css/reduced-motion.webp"
  alt="普通动画偏好与 prefers-reduced-motion 减少动态效果模式的对比"
  caption="保留状态变化，缩短或取消非必要位移动画，让不同感知需求的用户都能顺利操作。"
  :width="1440"
  :height="900"
/>

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    scroll-behavior: auto !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

项目里更推荐对关键动画做精细降级：

```css
@media (prefers-reduced-motion: reduce) {
  .dialog-panel {
    animation: none;
  }

  .skeleton-line {
    animation: none;
    background: #edf2f7;
  }
}
```

不要强迫所有用户接受大幅运动、视差滚动或长时间循环动画。

## 和 Vue / React 的关系

框架负责状态切换，CSS 负责视觉变化。

Vue 中可以配合 `<Transition>`：

```css
.fade-enter-active,
.fade-leave-active {
  transition: opacity 160ms ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
```

React 中常见做法是通过 class 切换状态：

```tsx
<div className={open ? 'panel panel--open' : 'panel'} />
```

```css
.panel {
  opacity: 0;
  transform: translateY(8px);
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.panel--open {
  opacity: 1;
  transform: translateY(0);
}
```

不要把所有动效都写成 JavaScript 定时器。优先使用 CSS 和框架提供的过渡机制。

## 实际项目常见问题

### 1. 动画导致页面卡顿

**原因**

动画了 `height`、`top`、`left` 等容易触发布局的属性，或者页面节点太多。

**解决方案**

- 优先使用 `transform` 和 `opacity`。
- 降低同时动画的节点数量。
- 避免在动画过程中频繁读写布局。
- 用 Performance 面板录制验证。

### 2. 关闭弹窗时动画没播放

**原因**

状态一变，组件直接卸载，CSS 没机会播放离场动画。

**解决方案**

- 使用 Vue `<Transition>`、React 动画库或保留一段离场状态。
- 不要直接 `v-if` / 条件渲染后马上销毁，除非框架过渡已经处理。

### 3. hover 动效移动端无意义

移动端没有稳定 hover。按钮 hover 可以保留，但关键交互状态不能只依赖 hover。

需要同时考虑：

- `:focus-visible`
- `:active`
- 明确选中状态

### 4. 动画太多导致页面廉价

动效应该服务信息层级。后台、文档站、管理台更适合克制、短促、低幅度动效。

### 5. 忘记减少动态效果

涉及大面积位移、缩放、循环动效时，都应该提供 `prefers-reduced-motion` 降级。

## 最佳实践

- hover/focus 状态优先使用 `transition`。
- loading、骨架屏、复杂入场才使用 `animation`。
- 不写 `transition: all`。
- 优先动画 `opacity` 和 `transform`。
- 动画时间通常控制在 120ms 到 240ms。
- 循环动画要谨慎，避免干扰阅读。
- 为明显运动提供 `prefers-reduced-motion` 降级。

## 学习检查

学完本节后，你应该能回答：

- `transition` 和 `animation` 怎么选。
- 为什么不推荐 `transition: all`。
- 为什么优先动画 `transform` 和 `opacity`。
- 弹窗离场动画为什么可能不播放。
- `prefers-reduced-motion` 解决什么问题。

## 参考资料

- [MDN: CSS transitions](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Transitions)
- [MDN: transition](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/transition)
- [MDN: CSS animations](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Animations)
- [MDN: @keyframes](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/%40keyframes)
- [MDN: prefers-reduced-motion](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/%40media/prefers-reduced-motion)

## 下一步学习

继续学习 [CSS 可访问性](/css/accessibility)，确保视觉样式不仅好看，也能被键盘用户和低视力用户正常使用。
