# 盒模型与布局基础

## 适合谁看

适合经常对元素宽高、间距、溢出和对齐感到困惑的学习者。

## 盒模型

每个元素都可以理解为一个盒子：

先观察下图，再回到代码理解四层空间。内容区保存真正的数据，`padding` 扩大内容与边框的距离，`border` 画出边界，`margin` 负责元素与外界的距离。

<DocFigure
  src="/images/css/box-model.webp"
  alt="CSS 盒模型四层结构图，依次标出内容区、内边距、边框和外边距"
  caption="盒模型不是四个互不相关的属性，而是从内容向外逐层包裹的空间结构。"
  :width="1440"
  :height="900"
/>

两个普通块级元素上下排列时，垂直 `margin` 可能发生折叠。图中可以看到，间距取较大的 `32px`，而不是把 `24px` 与 `32px` 相加得到 `56px`。

<DocFigure
  src="/images/css/margin-collapse.webp"
  alt="两个上下排列的块级元素发生外边距折叠，最终垂直间距为 32 像素"
  caption="外边距折叠只发生在特定的普通文档流场景；Flex、Grid 和建立 BFC 后的结果会不同。"
  :width="1440"
  :height="900"
/>

```text
content 内容
padding 内边距
border 边框
margin 外边距
```

推荐全局设置：

```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```

这样 `width` 会包含 padding 和 border，更容易控制布局。

## display

常见值：

| 值 | 用途 |
| --- | --- |
| `block` | 块级布局，占一行 |
| `inline` | 行内内容 |
| `inline-block` | 行内但可设置宽高 |
| `flex` | 一维布局 |
| `grid` | 二维布局 |
| `none` | 不渲染 |

## width 和 max-width

不要所有地方都写死宽度。

```css
.page-container {
  width: 100%;
  max-width: 1180px;
  margin: 0 auto;
  padding: 0 24px;
}
```

`max-width` 保证大屏不太散，`width: 100%` 保证小屏能收缩。

## overflow

内容超出容器时会产生溢出。

代码块适合横向滚动：

```css
.code-panel {
  overflow-x: auto;
}
```

页面整体不应该随便 `overflow-x: hidden` 掩盖问题。先找出哪个元素撑破了页面。

## 实际项目常见问题

### 1. 设置 width: 100vw 后出现横向滚动

**原因**

`100vw` 包含滚动条宽度，可能比页面可用宽度更宽。

**解决方案**

大多数容器用：

```css
width: 100%;
```

### 2. padding 撑大了元素

**原因**

没有设置 `box-sizing: border-box`。

**解决方案**

全局设置 border-box。

### 3. 内容撑破卡片

**解决方案**

```css
.card-title {
  min-width: 0;
  overflow-wrap: anywhere;
}
```

## 最佳实践

- 全局使用 `box-sizing: border-box`。
- 容器优先用 `width: 100%` 和 `max-width`。
- 不用 `overflow-x: hidden` 掩盖布局错误。
- 固定尺寸元素设置明确宽高和不可压缩行为。

## 下一步

继续学习 [Flex 与 Grid](/css/flex-grid)。
