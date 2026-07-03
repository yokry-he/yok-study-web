# 设计 Token 与主题

## 适合谁看

适合已经能写 CSS，但项目开始出现这些问题的人：

- 同一个主色在不同文件里写了很多遍。
- 深色模式、品牌换肤、主题切换越来越难维护。
- 组件库 token、CSS 变量、设计稿变量不知道怎么对应。
- 设计师改一个颜色，开发要全局搜索替换。
- 项目里颜色、间距、圆角、阴影没有统一命名。

设计 token 是把设计决策用稳定名称记录下来。它不是单纯的 CSS 变量，而是设计系统和代码之间的协议。

## 什么是设计 token

可以理解为：

```text
设计决策
↓
命名 token
↓
CSS 变量 / 组件库主题 / 多端样式
```

示例：

```css
:root {
  --color-brand-primary: #14b89a;
  --color-text-primary: #17211f;
  --radius-control: 8px;
  --space-3: 12px;
}
```

使用：

```css
.primary-button {
  border-radius: var(--radius-control);
  background: var(--color-brand-primary);
  color: white;
  padding: var(--space-3);
}
```

好处是后续换主题时改 token，而不是到处改业务样式。

## token 分层

不要一开始就把所有变量都平铺。

推荐分层：

| 层级 | 示例 | 用途 |
| --- | --- | --- |
| 基础 token | `--palette-mint-500` | 原始色板、基础间距 |
| 语义 token | `--color-success-bg` | 表达业务语义 |
| 组件 token | `--button-primary-bg` | 组件级定制 |

示例：

```css
:root {
  --palette-mint-500: #14b89a;
  --palette-red-600: #dc2626;

  --color-brand-primary: var(--palette-mint-500);
  --color-danger-text: var(--palette-red-600);

  --button-primary-bg: var(--color-brand-primary);
}
```

项目早期可以只做基础 token 和语义 token。组件 token 在组件数量变多后再细化。

## 常见 token 类型

| 类型 | 示例 |
| --- | --- |
| 颜色 | 品牌色、文本色、背景色、边框色、状态色 |
| 字体 | 字号、行高、字重 |
| 间距 | padding、gap、布局间隔 |
| 圆角 | 按钮、卡片、弹窗 |
| 阴影 | 浮层、卡片、下拉 |
| 尺寸 | 控件高度、图标尺寸、侧栏宽度 |
| 动效 | 时长、缓动曲线 |
| z-index | 顶栏、抽屉、弹窗、提示 |

不要把所有 CSS 属性都 token 化。token 应该表达会复用、会统一、会被主题影响的设计决策。

## CSS 变量实现主题

浅色主题：

```css
:root {
  --color-bg-page: #f7fbf9;
  --color-bg-surface: #ffffff;
  --color-text-primary: #17211f;
  --color-border-subtle: #dce8e3;
  --color-brand-primary: #14b89a;
}
```

深色主题：

```css
.dark {
  --color-bg-page: #101816;
  --color-bg-surface: #16211e;
  --color-text-primary: #eef8f4;
  --color-border-subtle: #2f4640;
  --color-brand-primary: #62d9c2;
}
```

业务样式只引用 token：

```css
.metric-card {
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border-subtle);
  color: var(--color-text-primary);
}
```

这样同一个组件可以跟随主题变化。

## 命名规则

命名要稳定、可读、可扩展。

推荐：

```text
--color-text-primary
--color-text-secondary
--color-bg-page
--color-bg-surface
--color-border-subtle
--space-1
--space-2
--radius-control
--shadow-popover
```

不推荐：

```text
--green
--big
--left-card-color
--new-color-1
--designer-blue
```

token 名称应该表达用途，而不是某一次设计稿里的位置。

## 和组件库主题的关系

如果项目使用 Naive UI、Element Plus、Ant Design Vue、Arco Design 等组件库，优先使用组件库主题能力。

常见策略：

```text
设计 token
↓
组件库主题 token
↓
业务 CSS 变量
```

不要用宽泛选择器强行改组件库内部 DOM：

```css
.page .n-button span {
  color: red;
}
```

这类写法依赖内部结构，不稳定。

更好的做法：

- 通过组件库主题配置改全局风格。
- 通过组件 props 设置状态。
- 业务区域使用自己的明确 class。
- 必须局部覆盖时限制在业务容器，并说明原因。

## 设计 token 文件示例

可以用 JSON 保存跨工具 token：

```json
{
  "color": {
    "brand": {
      "primary": {
        "$value": "#14b89a",
        "$type": "color"
      }
    }
  },
  "radius": {
    "control": {
      "$value": "8px",
      "$type": "dimension"
    }
  }
}
```

这类结构可以再转换成 CSS 变量、设计工具变量、移动端资源等。

项目早期不一定要引入复杂构建链路，但要先统一命名和使用方式。

## 实际项目常见问题

### 1. token 太多，没人知道用哪个

**原因**

没有分层，也没有使用规则。

**解决方案**

- 先控制 token 数量。
- 优先语义 token。
- 写清楚每个 token 的用途。
- 废弃 token 要有迁移说明。

### 2. 业务样式绕过 token

**现象**

页面里大量硬编码颜色：

```css
.user-card {
  color: #333;
  border-color: #eee;
}
```

**解决方案**

要求新增业务样式优先使用 token。确实需要新颜色时，先判断它是否应该成为 token。

### 3. 深色模式只反转背景

**问题**

深色主题不是把背景变黑就结束，还要检查文本、边框、阴影、焦点、状态色、图表色。

**解决方案**

按页面状态验证：

- 默认态。
- hover。
- active。
- focus。
- disabled。
- error。
- selected。

### 4. token 命名和业务耦合太深

不推荐：

```css
--dashboard-card-title-color: #17211f;
```

如果多个地方都需要主文本色，应使用：

```css
--color-text-primary: #17211f;
```

业务组件可以再按需要映射组件 token。

### 5. 组件库主题和业务 token 不一致

**解决方案**

把品牌色、文本色、边框色等核心 token 同步到组件库主题配置，避免组件库按钮和业务卡片像两个系统。

## 最佳实践

- token 表达设计决策，不是给所有属性起变量名。
- 先建立颜色、间距、圆角、阴影、字号这些高频 token。
- 命名优先表达语义和用途。
- 业务 CSS 不硬编码可复用颜色和尺寸。
- 深色模式要成组验证状态色和对比度。
- 和组件库主题能力对齐，不依赖内部 DOM 覆盖。
- token 变更要更新文档和迁移说明。

## 学习检查

学完本节后，你应该能回答：

- 设计 token 和普通 CSS 变量有什么关系。
- 基础 token、语义 token、组件 token 怎么区分。
- 为什么 token 命名要表达用途而不是颜色名。
- 深色主题为什么不能只改背景色。
- 组件库项目里应该如何让 token 和主题配置一致。

## 参考资料

- [MDN: Using CSS custom properties](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Cascading_variables/Using_custom_properties)
- [MDN: CSS custom properties](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Properties/--%2A)
- [MDN: var()](https://developer.mozilla.org/en-US/docs/Web/CSS/Reference/Values/var)
- [Design Tokens Format Module](https://www.designtokens.org/tr/drafts/format/)
- [W3C Design Tokens Community Group](https://www.w3.org/community/design-tokens/)

## 下一步学习

继续学习 [项目样式架构](/css/architecture)，把 token、业务 class、全局样式和组件库边界组织到可维护的结构里。
