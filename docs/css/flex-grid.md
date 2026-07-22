# Flex 与 Grid

## 适合谁看

适合准备系统掌握现代 CSS 布局的学习者。

简单理解：

- Flex 适合一行或一列。
- Grid 适合行和列同时控制。

## Flex：一维布局

工具栏：

```css
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.toolbar-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
```

适合：

- 导航栏。
- 按钮组。
- 表单行。
- 卡片内部左右布局。

## Flex 常见属性

先把主轴和交叉轴认清，再记属性。`flex-direction` 决定主轴方向，`justify-content` 沿主轴分配空间，`align-items` 沿交叉轴对齐项目。

<DocFigure
  src="/images/css/flex-main-cross-axis.webp"
  alt="Flex 容器的主轴、交叉轴以及 justify-content 和 align-items 的作用方向"
  caption="先判断轴，再选择对齐属性，可以避免靠反复试值完成布局。"
  :width="1440"
  :height="900"
/>

当子项内容过长时，Flex 项目的默认最小宽度可能阻止它继续收缩。图中对比了溢出状态与设置 `min-width: 0` 后的截断状态。

<DocFigure
  src="/images/css/flex-overflow.webp"
  alt="Flex 子项长文本溢出与设置 min-width 0 后正确截断的对比"
  caption="项目中常见的“省略号不生效”，根因往往不是 text-overflow，而是 Flex 子项仍保留内容最小宽度。"
  :width="1440"
  :height="900"
/>

| 属性 | 用途 |
| --- | --- |
| `align-items` | 交叉轴对齐 |
| `justify-content` | 主轴分布 |
| `gap` | 子项间距 |
| `flex-wrap` | 是否换行 |
| `flex-shrink` | 是否压缩 |
| `flex: 1` | 占据剩余空间 |

## Grid：二维布局

Grid 适合同时控制行和列。具名区域让页面结构直接体现在 CSS 中，比依赖网格序号更容易维护。

<DocFigure
  src="/images/css/grid-template-areas.webp"
  alt="使用 header sidebar main 和 aside 命名区域组成的 CSS Grid 页面布局"
  caption="grid-template-areas 把页面骨架写成可阅读的二维地图，适合后台与内容型页面。"
  :width="1440"
  :height="900"
/>

卡片网格通常不需要手写多个媒体查询。`repeat(auto-fit, minmax(...))` 可以让列数随容器宽度自然变化。

<DocFigure
  src="/images/css/grid-minmax.webp"
  alt="CSS Grid 使用 auto-fit 与 minmax 在不同容器宽度下自动调整卡片列数"
  caption="minmax 为单列设置可接受的宽度范围，auto-fit 再决定当前能放下多少列。"
  :width="1440"
  :height="900"
/>

卡片网格：

```css
.card-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
```

响应式：

```css
@media (max-width: 900px) {
  .card-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 560px) {
  .card-grid {
    grid-template-columns: 1fr;
  }
}
```

`minmax(0, 1fr)` 可以减少内容撑破网格。

## 表单布局

```css
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.form-field--full {
  grid-column: 1 / -1;
}
```

移动端：

```css
@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
```

## 实际项目常见问题

### 1. 按钮被压扁

**解决方案**

```css
.toolbar-action {
  flex-shrink: 0;
}
```

### 2. 左侧文字太长导致右侧按钮消失

**解决方案**

```css
.toolbar-title {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toolbar-actions {
  flex-shrink: 0;
}
```

### 3. Grid 卡片撑出页面

**解决方案**

使用 `minmax(0, 1fr)`，并处理长文本换行。

## 最佳实践

- 工具栏、按钮组、行内对齐用 Flex。
- 卡片、表单、仪表盘用 Grid。
- 固定操作区设置 `flex-shrink: 0`。
- Grid 列建议使用 `minmax(0, 1fr)`。
- 响应式布局从结构阶段就考虑。

## 下一步

继续学习 [响应式设计](/css/responsive)。
