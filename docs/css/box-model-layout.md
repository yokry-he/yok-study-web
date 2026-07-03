# 盒模型与布局基础

## 适合谁看

适合经常对元素宽高、间距、溢出和对齐感到困惑的学习者。

## 盒模型

每个元素都可以理解为一个盒子：

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
