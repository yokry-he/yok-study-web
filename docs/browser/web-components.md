# Web Components

## 适合谁看

适合已经会 Vue 或 React 组件，但希望理解浏览器原生组件模型、跨框架组件复用、Shadow DOM 样式隔离、自定义元素生命周期的学习者。

Web Components 不是 Vue 或 React 的替代品。它更适合做跨技术栈复用的组件、设计系统底层能力、嵌入式组件、小部件和平台级 UI 能力。

## Web Components 包含什么

Web Components 通常由几类能力组成：

| 能力 | 作用 |
| --- | --- |
| Custom Elements | 定义自定义 HTML 标签 |
| Shadow DOM | 封装内部 DOM 和样式 |
| HTML Templates | 定义可复用模板 |
| Slots | 让外部内容插入组件内部 |

示例标签：

```html
<user-card name="Tom"></user-card>
```

浏览器可以像普通 HTML 元素一样识别这个自定义元素。

## Custom Elements

基础示例：

```ts
class UserCard extends HTMLElement {
  connectedCallback() {
    const name = this.getAttribute('name') || '未命名'

    this.innerHTML = `
      <article>
        <strong>${name}</strong>
      </article>
    `
  }
}

customElements.define('user-card', UserCard)
```

使用：

```html
<user-card name="Tom"></user-card>
```

自定义元素名称必须包含短横线，例如 `user-card`，不能叫 `usercard`。

## 生命周期

常见生命周期：

| 方法 | 触发时机 |
| --- | --- |
| `connectedCallback` | 元素插入 DOM |
| `disconnectedCallback` | 元素从 DOM 移除 |
| `attributeChangedCallback` | 监听的属性变化 |
| `adoptedCallback` | 元素被移动到新文档 |

监听属性：

```ts
class UserBadge extends HTMLElement {
  static get observedAttributes() {
    return ['status']
  }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    if (name === 'status') {
      this.render()
    }
  }

  connectedCallback() {
    this.render()
  }

  render() {
    this.textContent = this.getAttribute('status') || 'unknown'
  }
}

customElements.define('user-badge', UserBadge)
```

## Shadow DOM

Shadow DOM 用于封装内部结构和样式，避免被页面外部样式随意影响。

```ts
class AppButton extends HTMLElement {
  connectedCallback() {
    const shadow = this.attachShadow({ mode: 'open' })

    shadow.innerHTML = `
      <style>
        button {
          border: 0;
          border-radius: 6px;
          padding: 8px 12px;
          background: #14b89a;
          color: white;
        }
      </style>
      <button><slot></slot></button>
    `
  }
}

customElements.define('app-button', AppButton)
```

使用：

```html
<app-button>保存</app-button>
```

Shadow DOM 能降低样式污染风险，但也会带来主题、可访问性和测试的额外设计成本。

## Template 和 Slot

template 不会直接渲染：

```html
<template id="card-template">
  <style>
    .card {
      border: 1px solid #dfe7e3;
      padding: 12px;
    }
  </style>
  <article class="card">
    <slot name="title"></slot>
    <slot></slot>
  </article>
</template>
```

组件里克隆：

```ts
const template = document.querySelector<HTMLTemplateElement>('#card-template')!

class InfoCard extends HTMLElement {
  connectedCallback() {
    const shadow = this.attachShadow({ mode: 'open' })
    shadow.append(template.content.cloneNode(true))
  }
}

customElements.define('info-card', InfoCard)
```

使用：

```html
<info-card>
  <h3 slot="title">标题</h3>
  <p>内容</p>
</info-card>
```

## 什么时候适合用

适合：

- 跨 Vue、React、普通 HTML 复用组件。
- 给外部系统嵌入小组件。
- 做不依赖框架的基础 UI 元素。
- 封装富交互但边界清晰的控件。
- 在多个技术栈间共享设计系统底层能力。

不适合：

- 普通 Vue 或 React 单项目里的所有业务组件。
- 复杂状态管理页面。
- 和框架深度绑定的路由页面。
- 团队还不熟悉原生 DOM 和可访问性时大规模使用。

## 和 Vue / React 的关系

Web Components 可以被 Vue 或 React 使用，但不是完全无缝。

需要注意：

- props 和 attributes 的差异。
- 事件命名和监听方式。
- 表单控件集成。
- 主题传递。
- SSR 和 hydration。
- TypeScript 类型声明。

如果组件只在 Vue 项目里使用，优先写 Vue 组件。如果组件需要跨多个框架复用，再考虑 Web Components。

## 实际项目常见问题

### 1. 外部样式改不了组件内部样式

Shadow DOM 的目标就是隔离。要想开放定制能力，应设计：

- CSS custom properties。
- `::part`。
- slots。
- 明确的属性。

示例：

```css
app-button {
  --app-button-bg: #1677ff;
}
```

### 2. 组件在 React 中事件监听不符合预期

CustomEvent 和 React 合成事件不是完全一回事。复杂事件要测试框架集成方式。

建议：

- 使用标准 DOM 事件。
- 文档写清事件名和 detail 结构。
- 为 React/Vue 提供轻量 wrapper。

### 3. 自定义元素重复注册报错

同一个名称不能重复注册。

处理：

```ts
if (!customElements.get('app-button')) {
  customElements.define('app-button', AppButton)
}
```

## 可访问性建议

Web Components 不会自动获得可访问性。

你仍然需要：

- 使用语义化元素。
- 设置合适 ARIA。
- 支持键盘操作。
- 管理焦点。
- 暴露 label 和 description。
- 测试屏幕阅读器行为。

不要因为 Shadow DOM 封装了结构，就忽略用户如何操作和理解组件。

## 项目建议

- 小范围试点，不要直接重写全部组件库。
- 明确哪些属性、事件、slot、CSS 变量是公共 API。
- 提供 Vue/React 使用示例。
- 设计主题和可访问性策略。
- 对外发布前写完整文档和版本变更说明。

## 下一步学习

- [常用 Web API](/browser/web-apis)
- [浏览器安全基础](/browser/security)
- [CSS 项目样式架构](/css/architecture)
- [组件库实战](/projects/component-library)
