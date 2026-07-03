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

| 属性 | 用途 |
| --- | --- |
| `align-items` | 交叉轴对齐 |
| `justify-content` | 主轴分布 |
| `gap` | 子项间距 |
| `flex-wrap` | 是否换行 |
| `flex-shrink` | 是否压缩 |
| `flex: 1` | 占据剩余空间 |

## Grid：二维布局

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
