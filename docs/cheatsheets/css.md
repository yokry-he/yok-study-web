# CSS 速查

## 盒模型

```css
.card {
  box-sizing: border-box;
  width: 320px;
  padding: 16px;
  border: 1px solid #dfe7e3;
}
```

推荐全局设置：

```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```

## Flex 常用写法

水平排列：

```css
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
```

固定尺寸元素不压缩：

```css
.avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}
```

自动占满剩余空间：

```css
.search-input {
  flex: 1 1 auto;
  min-width: 0;
}
```

## Grid 常用写法

自适应卡片：

```css
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}
```

两栏布局：

```css
.layout {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 24px;
}
```

移动端改一栏：

```css
@media (max-width: 768px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
```

## 文本处理

单行省略：

```css
.title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
```

多行省略：

```css
.summary {
  display: -webkit-box;
  overflow: hidden;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
```

长单词换行：

```css
.content {
  overflow-wrap: anywhere;
}
```

## 响应式断点

常见断点：

| 宽度 | 典型设备 |
| ---: | --- |
| `390px` | 手机 |
| `768px` | 平板或窄屏 |
| `1024px` | 小桌面 |
| `1280px` | 常规桌面 |
| `1440px` | 宽桌面 |

移动端优先写法：

```css
.page {
  padding: 16px;
}

@media (min-width: 768px) {
  .page {
    padding: 24px;
  }
}
```

## 业务样式命名

推荐明确业务 class：

```css
.user-page__toolbar {}
.user-page__table-action {}
.permission-switch-row {}
.metric-card__value {}
```

避免宽泛选择器：

```css
.page div {}
.content button {}
.panel * {}
div > div {}
```

这些选择器容易污染组件库内部 DOM，导致按钮、开关、表格、弹窗变形。

## 常见布局坑

| 问题 | 处理 |
| --- | --- |
| flex 子元素撑破容器 | 给子元素 `min-width: 0` |
| 图标被压成椭圆 | 固定宽高和 `flex: 0 0 size` |
| 表格操作列换行 | 操作容器 `flex-shrink: 0` |
| 移动端横向滚动 | 检查固定宽度和长文本 |
| 组件库样式异常 | 先查宽泛选择器和全局样式 |

## 项目建议

- 页面布局用明确容器 class。
- 固定尺寸视觉元素都设置不可压缩。
- 不依赖组件库内部 DOM 层级写样式。
- 全局样式只放 reset、变量和基础排版。
- 修改全局样式后检查桌面和窄屏。

## 下一步学习

- [CSS 学习导览](/css/introduction)
- [Flex 与 Grid](/css/flex-grid)
- [响应式设计](/css/responsive)
- [项目样式架构](/css/architecture)
