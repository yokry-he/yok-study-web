# 项目样式架构

## 适合谁看

适合准备在 Vue 项目中长期维护样式，或者已经遇到全局样式污染、组件库样式异常的学习者。

## 推荐目录

```text
src/styles/
├─ reset.css
├─ variables.css
├─ layout.css
├─ utilities.css
└─ theme.css
```

组件样式写在组件内：

```vue
<style scoped>
.user-card {
  display: grid;
  gap: 12px;
}
</style>
```

## CSS 变量

```css
:root {
  --app-color-primary: #1d9a78;
  --app-color-text: #17231f;
  --app-radius-card: 8px;
  --app-space-md: 16px;
}
```

使用：

```css
.metric-card {
  border-radius: var(--app-radius-card);
  color: var(--app-color-text);
  padding: var(--app-space-md);
}
```

## 业务 class 命名

推荐：

```css
.user-search-form {}
.user-search-form__actions {}
.permission-switch-row {}
.metric-card__value {}
```

不推荐：

```css
.page div {}
.content span {}
.panel * {}
```

## 组件库样式边界

不要假设组件库内部 DOM 稳定。

不推荐：

```css
.user-page .n-data-table-td div {}
```

优先：

- 组件库主题 token。
- 组件 props。
- CSS 变量。
- 官方暴露 class 或 API。
- 外层业务容器样式。

## 全局样式应该少

适合全局：

- reset。
- CSS 变量。
- body 基础字体。
- 通用布局容器。
- 少量工具类。

不适合全局：

- 某个页面的按钮样式。
- 某个表格的行高。
- 某个弹窗的内容布局。

## 实际项目常见问题

### 1. 一个全局 button 样式毁掉组件库

**问题**

```css
button {
  width: 100%;
  border: none;
}
```

**解决方案**

只作用于业务 class：

```css
.login-submit-button {
  width: 100%;
}
```

### 2. scoped 后仍然影响子组件

`scoped` 能限制当前组件样式，但如果使用深度选择器，仍然可能影响子组件。深度选择器要谨慎，并写明原因。

### 3. 主题色到处写死

**解决方案**

抽成变量：

```css
--app-color-primary: #1d9a78;
```

## 最佳实践

- 全局样式越少越好。
- 业务样式使用明确 class。
- 组件库样式通过官方机制调整。
- 主题色、间距、圆角抽变量。
- 修改全局样式后检查关键页面。

## 下一步

继续学习 [常见问题](/css/troubleshooting)。
