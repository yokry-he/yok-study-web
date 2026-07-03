# CSS 常见问题

## 1. 页面横向溢出

### 症状

移动端或桌面端底部出现横向滚动条。

### 常见原因

- 使用 `width: 100vw`。
- 固定宽度元素超出容器。
- 表格、代码块、图片太宽。
- Grid 子项内容撑破。

### 解决方案

```css
.card-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.code-block {
  overflow-x: auto;
}

.image {
  max-width: 100%;
}
```

## 2. 头像变椭圆

### 原因

Flex 容器中被压缩。

### 解决方案

```css
.avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}
```

## 3. 表格操作列被压缩

### 解决方案

```css
.table-actions {
  display: inline-flex;
  gap: 8px;
  white-space: nowrap;
}

.table-actions__button {
  flex-shrink: 0;
}
```

## 4. 组件库控件突然变形

### 排查顺序

1. 搜索宽泛选择器。
2. 检查全局 button、input、table 样式。
3. 检查是否覆盖组件库内部 class。
4. 检查 HMR 或旧进程缓存。

搜索命令：

```bash
rg "(\\.\\w+\\s+(div|span|button|\\*)|div > div|\\.\\w+ \\*)" src
```

## 5. 文本撑破卡片

### 解决方案

```css
.card-title {
  min-width: 0;
  overflow-wrap: anywhere;
}
```

如果是单行：

```css
.card-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
```

## 6. 移动端首屏内容被导航挤下去

### 原因

桌面导航在移动端直接堆叠显示。

### 解决方案

移动端使用：

- 抽屉。
- 顶部菜单按钮。
- 底部导航。
- 更多菜单。

首屏优先展示当前页面核心内容。

## 7. hover 动画导致卡顿

### 原因

动画修改了宽高、阴影过重或触发布局。

### 解决方案

优先使用：

```css
transform: translateY(-2px);
opacity: 0.9;
```

避免频繁动画 `width`、`height`、`top`、`left`。

## 快速排查表

| 问题 | 优先检查 |
| --- | --- |
| 横向滚动 | 固定宽度、100vw、表格、代码块 |
| 控件变形 | 全局选择器、组件库内部覆盖 |
| 头像变形 | width、height、flex-shrink |
| 文本溢出 | min-width、overflow-wrap |
| 移动端混乱 | 是否直接复用桌面布局 |

## 最佳实践

- 先定位溢出元素，再修改。
- 不用高优先级覆盖掩盖根因。
- 固定尺寸元素必须不可压缩。
- 修改全局样式后做桌面和移动端验证。
