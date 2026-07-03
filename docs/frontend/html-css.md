# HTML 与 CSS

## 适合谁看

适合准备学习 Vue，但页面结构、布局、响应式和样式组织还不稳定的学习者。

Vue 组件最终仍然会渲染成 HTML 和 CSS。如果基础不牢，后面做组件、后台管理页面、移动端适配时，很容易出现布局错位、样式污染、按钮变形、表格横向溢出等问题。

## 你会学到什么

- HTML 语义化为什么重要。
- Flexbox 和 Grid 分别适合什么布局。
- 响应式页面如何从一开始就考虑。
- Vue 项目里 CSS 如何避免污染组件库。
- 常见布局问题如何排查。

## HTML：先写清楚结构

HTML 不只是把内容放到页面上，还要表达内容含义。

```html
<main>
  <section>
    <h1>用户管理</h1>
    <p>管理系统用户、角色和启用状态。</p>
  </section>
</main>
```

常见语义标签：

| 标签 | 适合场景 |
| --- | --- |
| `header` | 页面或区域头部 |
| `nav` | 导航区域 |
| `main` | 页面主体 |
| `section` | 一个独立内容区块 |
| `article` | 可独立阅读的内容 |
| `button` | 执行动作 |
| `a` | 跳转链接 |

重要规则：**跳转用 `a`，操作用 `button`。**

例如：

```html
<a href="/users">进入用户管理</a>
<button type="button">删除用户</button>
```

## CSS：先掌握盒模型

每个元素都可以理解为一个盒子：

```text
content 内容
padding 内边距
border 边框
margin 外边距
```

建议全局使用：

```css
*,
*::before,
*::after {
  box-sizing: border-box;
}
```

这样元素宽度会包含 padding 和 border，更容易控制布局。

## Flexbox：一维布局

Flexbox 适合横向或纵向排列一组元素。

工具栏示例：

```html
<div class="user-toolbar">
  <div class="user-toolbar__filters">筛选条件</div>
  <div class="user-toolbar__actions">操作按钮</div>
</div>
```

```css
.user-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.user-toolbar__actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
```

`flex-shrink: 0` 很重要。头像、图标按钮、操作按钮、状态点这类固定尺寸元素，不应该被挤压。

## Grid：二维布局

Grid 适合卡片网格、表单布局、仪表盘。

```css
.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

@media (max-width: 900px) {
  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 520px) {
  .metric-grid {
    grid-template-columns: 1fr;
  }
}
```

`minmax(0, 1fr)` 可以减少内容撑破网格导致横向滚动的问题。

## 响应式布局

响应式不是把桌面页面缩小，而是重新安排信息优先级。

后台页面常见处理：

| 桌面端 | 移动端 |
| --- | --- |
| 固定侧边栏 | 顶部菜单按钮 + 抽屉 |
| 表格多列展示 | 关键字段卡片或横向滚动 |
| 工具栏横排 | 筛选折叠或换行 |
| 大弹窗 | 全屏抽屉或底部弹层 |

移动端不要把整块桌面侧边栏直接堆到首屏上方，否则用户第一屏看不到核心内容。

## Vue 项目中的 CSS 组织

推荐：

```text
src/styles/
├─ reset.css
├─ variables.css
├─ layout.css
└─ utilities.css
```

组件内部样式：

```vue
<style scoped>
.user-card {
  display: grid;
  gap: 12px;
}

.user-card__title {
  font-weight: 700;
}
</style>
```

## 避免样式污染

不要写：

```css
.page div {
  margin-bottom: 12px;
}

.toolbar button {
  width: 100%;
}

.panel * {
  line-height: 1.5;
}
```

这些选择器会影响组件库内部 DOM，导致按钮、表格、弹窗、开关等控件异常。

推荐写明确业务 class：

```css
.user-search-form__field {
  min-width: 180px;
}

.toolbar-action {
  flex-shrink: 0;
}

.metric-card__value {
  font-size: 28px;
  font-weight: 800;
}
```

如果确实需要调整组件库，优先使用组件库的主题 token、组件 props、CSS 变量或官方 API。

## 实际项目常见问题

### 1. 页面出现横向滚动条

**常见原因**

- 固定宽度元素超出屏幕。
- Grid 子项内容撑破容器。
- 表格或代码块没有处理溢出。
- 使用了 `width: 100vw`，叠加滚动条宽度后溢出。

**解决方案**

```css
.page-container {
  max-width: 100%;
  overflow-x: hidden;
}

.card-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.code-panel {
  overflow-x: auto;
}
```

### 2. 头像被压成椭圆

**原因**

Flex 布局中头像被压缩。

**解决方案**

```css
.user-avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}
```

### 3. 操作按钮在表格里被挤坏

**解决方案**

```css
.table-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.table-actions__button {
  flex-shrink: 0;
}
```

### 4. 组件库样式突然全乱了

**排查顺序**

1. 搜索宽泛选择器。
2. 检查是否覆盖了组件库内部 class。
3. 检查全局样式是否影响 `button`、`input`、`table`。
4. 检查 HMR 或旧进程缓存。

## 最佳实践

- 先写清楚 HTML 结构，再写样式。
- 页面布局优先用 Grid 和 Flexbox。
- 固定尺寸视觉元素设置稳定宽高和不可压缩行为。
- 业务样式使用明确 class。
- 修改全局样式后必须检查关键页面和移动端。

## 下一步学习

继续学习 [JavaScript 基础](/javascript/fundamentals)。
