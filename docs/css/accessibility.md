# CSS 可访问性

## 适合谁看

适合已经能完成页面样式，但希望页面对更多用户可用、也能通过基础可访问性检查的人：

- 自定义按钮后键盘焦点看不见。
- 颜色看起来好看，但文字对比度不足。
- 表单错误只靠红色提示。
- hover 状态有，focus 状态没有。
- 动画、滚动或闪烁让部分用户不舒服。

CSS 可访问性不是额外装饰，而是保证界面能被键盘、屏幕阅读器、低视力用户和动作敏感用户正常使用的基础工程质量。

## CSS 负责哪些可访问性

CSS 主要影响：

| 方向 | CSS 影响 |
| --- | --- |
| 可见性 | 文字是否清晰、状态是否可辨认 |
| 焦点 | 键盘用户是否知道当前位置 |
| 对比度 | 文本、边框、图标、状态是否足够明显 |
| 响应式 | 放大字体和窄屏下是否可用 |
| 动效 | 是否尊重减少动态效果 |
| 隐藏内容 | 是否错误隐藏了辅助技术需要的信息 |

语义结构主要由 HTML 负责，但 CSS 不能破坏语义带来的可用性。

## 不要移除焦点样式

危险写法：

```css
button:focus {
  outline: none;
}
```

如果移除浏览器默认焦点，就必须提供更清楚的替代样式。

推荐：

```css
.primary-button:focus-visible {
  outline: 3px solid #14b89a;
  outline-offset: 3px;
}
```

`focus-visible` 更适合区分键盘焦点和鼠标点击，减少鼠标用户看到不必要焦点框。

## 焦点样式要稳定

焦点样式应该：

观察下图中的键盘焦点：轮廓位于控件外侧，不会挤压文字或改变布局；颜色同时兼顾亮色背景与操作按钮。

<DocFigure
  src="/images/frontend/accessible-focus.webp"
  alt="表单中使用清晰 focus-visible 外轮廓展示键盘焦点"
  caption="稳定的焦点环既要明显，也不能通过改变边框宽度导致页面抖动。"
  :width="1440"
  :height="900"
/>

- 清晰可见。
- 不依赖颜色微小变化。
- 不改变布局尺寸。
- 和背景有足够对比。
- 不被 `overflow: hidden` 裁掉。

示例：

```css
.form-input:focus-visible {
  border-color: #0f9f86;
  box-shadow: 0 0 0 3px rgb(20 184 154 / 24%);
  outline: none;
}
```

如果使用 `box-shadow` 做焦点环，要确认它不会被父容器裁剪。

## 颜色不能是唯一信息

不推荐只靠颜色表达状态：

```css
.field-message {
  color: #dc2626;
}
```

更稳的做法：

```html
<p class="field-message field-message--error">
  <span class="field-message__icon" aria-hidden="true">!</span>
  手机号格式不正确
</p>
```

```css
.field-message--error {
  color: #b42318;
}

.field-message__icon {
  display: inline-grid;
  width: 1rem;
  height: 1rem;
  place-items: center;
  border: 1px solid currentColor;
  border-radius: 999px;
  font-size: 0.75rem;
}
```

颜色、图标、文案一起表达状态，用户更容易理解。

## 文本对比度

文本颜色要和背景保持足够对比。尤其注意：

- 浅灰文字。
- 禁用态文字。
- 占位符。
- 彩色标签。
- 低透明度文字。
- 深色模式。

不要为了“高级感”把正文颜色压得太浅。

常见建议：

```css
:root {
  --text-primary: #17211f;
  --text-secondary: #4f635f;
  --text-muted: #6b7f7a;
  --surface: #ffffff;
}
```

正文优先使用主文本色，辅助信息才使用次级文本色。禁用态也要能读，不要接近背景。

## 隐藏内容要小心

`display: none` 会让内容从视觉和辅助技术中都消失。

如果内容只是视觉隐藏，但要保留给屏幕阅读器，可以使用 visually hidden 模式：

```css
.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
  border: 0;
}
```

适合：

- 图标按钮的文本标签。
- 跳过导航链接。
- 表格或表单的辅助说明。

不适合把一堆可见内容强行塞给屏幕阅读器。

## 响应式和放大

可访问性不仅是颜色和键盘。用户可能：

- 放大浏览器到 200%。
- 使用较大系统字体。
- 用窄屏设备访问。
- 使用横屏或分屏。

CSS 要避免：

- 固定高度导致文字溢出。
- 按钮文字被压缩。
- 表格操作列挤成一团。
- 长单词撑破容器。

常用处理：

```css
.action-button {
  min-height: 40px;
  padding: 0.5rem 0.875rem;
  white-space: normal;
}

.content-panel {
  overflow-wrap: anywhere;
}
```

固定尺寸视觉元素要防止被压缩：

```css
.status-dot {
  width: 0.5rem;
  height: 0.5rem;
  flex: 0 0 0.5rem;
  border-radius: 999px;
}
```

## 动效可访问性

明显运动要提供减少动态效果：

```css
@media (prefers-reduced-motion: reduce) {
  .page-enter,
  .skeleton-line,
  .parallax-layer {
    animation: none;
    transition: none;
  }
}
```

避免：

- 大面积视差。
- 快速闪烁。
- 自动循环且无法停止的动画。
- 页面进入时过度缩放和旋转。

## 组件库项目里的 CSS 可访问性

使用组件库时，不要用宽泛 CSS 破坏组件库状态样式。

危险：

```css
.page button {
  outline: none;
  border: 0;
}
```

这会影响所有按钮，包括组件库按钮、弹窗按钮和表单按钮。

更好的方式：

- 使用组件库主题 token。
- 使用组件库提供的 props。
- 只写明确业务 class。
- 保留 focus、disabled、error、active 等状态样式。

## 实际项目常见问题

### 1. 键盘 tab 到按钮没有任何提示

**原因**

全局 CSS 移除了 outline，或者自定义按钮没有 focus 样式。

**解决方案**

恢复默认 outline，或提供 `:focus-visible` 样式。

### 2. 禁用按钮文字看不清

**原因**

透明度太低，背景和文字对比不足。

**解决方案**

不要只用 `opacity: 0.3`。为禁用态设计明确的文本色、边框色和背景色。

### 3. 表单错误只显示红框

**问题**

色弱用户或屏幕阅读器用户可能无法理解。

**解决方案**

错误状态同时提供文案、图标、边框和必要的 ARIA 关联。

### 4. 移动端字号过小

**解决方案**

正文和表单控件保持可读尺寸。移动端不要为了塞更多内容把字体压得过小。

### 5. 主题切换后对比度失效

**原因**

只设计了浅色主题 token，没有验证深色主题组合。

**解决方案**

主题 token 要成组验证，尤其是文本、背景、边框、焦点和状态色。

## 最佳实践

- 不移除焦点样式，除非提供更好的替代。
- 文本、图标、边框和焦点状态都要考虑对比度。
- 状态信息不要只靠颜色表达。
- 响应式要考虑文字放大和长内容。
- 动效要尊重 `prefers-reduced-motion`。
- 组件库项目不要用宽泛选择器污染内部控件。
- 每次改全局样式后检查键盘焦点、禁用态、错误态和移动端。

## 学习检查

学完本节后，你应该能回答：

- 为什么不能直接 `outline: none`。
- `:focus-visible` 适合解决什么问题。
- 为什么状态信息不能只靠颜色。
- 视觉隐藏和 `display: none` 有什么区别。
- 主题切换后为什么要重新检查对比度。

## 参考资料

- [W3C WAI: Focus Appearance](https://www.w3.org/WAI/WCAG22/Understanding/focus-appearance.html)
- [W3C WAI: Non-text Contrast](https://www.w3.org/WAI/WCAG21/Understanding/non-text-contrast.html)
- [MDN: prefers-reduced-motion](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/At-rules/%40media/prefers-reduced-motion)

## 下一步学习

继续学习 [设计 Token 与主题](/css/design-tokens)，把颜色、间距、圆角、阴影和状态样式沉淀成可维护的系统。
