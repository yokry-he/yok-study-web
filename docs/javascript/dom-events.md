# DOM 事件

## 适合谁看

适合已经会写按钮点击、输入框监听和 Vue/React 事件绑定，但对原生 DOM 事件机制还不够清楚的人：

- 不知道事件冒泡和捕获是什么。
- 不知道什么时候用事件委托。
- 组件卸载后事件监听没有清理，导致重复触发。
- `event.target`、`event.currentTarget` 经常分不清。
- 想理解框架事件绑定背后的浏览器机制。

Vue 和 React 帮你封装了常见事件写法，但底层仍然是浏览器事件系统。理解 DOM 事件后，处理弹窗、下拉菜单、拖拽、快捷键、点击外部关闭、列表事件委托会更稳。

## 事件监听的基本写法

原生 DOM 事件使用 `addEventListener`：

```ts
const button = document.querySelector<HTMLButtonElement>('#save')!

function handleClick() {
  console.log('save')
}

button.addEventListener('click', handleClick)
```

移除监听：

```ts
button.removeEventListener('click', handleClick)
```

注意：移除时必须传入同一个函数引用。

不推荐这样写：

```ts
button.addEventListener('click', () => {
  console.log('save')
})

button.removeEventListener('click', () => {
  console.log('save')
})
```

这两个箭头函数不是同一个引用，无法移除之前的监听。

## 事件对象

事件回调会收到一个事件对象：

```ts
button.addEventListener('click', (event) => {
  console.log(event.type)
  console.log(event.target)
  console.log(event.currentTarget)
})
```

常用属性：

| 属性 | 含义 |
| --- | --- |
| `type` | 事件类型，例如 `click`、`input` |
| `target` | 真正触发事件的元素 |
| `currentTarget` | 当前正在执行监听器的元素 |
| `preventDefault()` | 阻止默认行为 |
| `stopPropagation()` | 阻止继续冒泡 |

`target` 和 `currentTarget` 是事件委托里最容易混淆的两个概念。

## 冒泡和捕获

一次点击不是只到达目标按钮。事件会沿祖先链捕获到目标，再沿同一条路径冒泡，这也是列表事件委托能够工作的基础。

<DocFigure
  src="/images/javascript/dom-event-path.webp"
  alt="DOM 点击事件从 window 和 document 捕获到 button，再经过列表项向 document 冒泡"
  caption="target 表示最初触发节点，currentTarget 表示当前监听器绑定节点；两者不能混用。"
  :width="1440"
  :height="900"
/>

如果中途调用 `stopPropagation()`，后续路径上的监听器将失去执行机会；排查“点击没反应”时要同时检查传播阶段和阻止位置。

页面结构：

```html
<div id="panel">
  <button id="save">保存</button>
</div>
```

点击按钮时，事件大致经历：

```text
捕获阶段：window -> document -> panel -> button
目标阶段：button
冒泡阶段：button -> panel -> document -> window
```

默认监听通常发生在冒泡阶段：

```ts
panel.addEventListener('click', () => {
  console.log('panel')
})

button.addEventListener('click', () => {
  console.log('button')
})
```

点击按钮时，先执行按钮，再冒泡到 panel。

如果要在捕获阶段监听：

```ts
document.addEventListener(
  'click',
  () => {
    console.log('capture')
  },
  { capture: true }
)
```

## 事件委托

事件委托是把监听器挂到父元素上，通过事件冒泡处理子元素事件。

适合列表：

```html
<ul id="user-list">
  <li data-id="1">Tom</li>
  <li data-id="2">Jerry</li>
</ul>
```

```ts
const list = document.querySelector<HTMLUListElement>('#user-list')!

list.addEventListener('click', (event) => {
  const target = event.target as HTMLElement
  const item = target.closest<HTMLLIElement>('li[data-id]')

  if (!item || !list.contains(item)) return

  console.log(item.dataset.id)
})
```

好处：

- 不需要给每一项都绑定监听。
- 动态新增列表项也能响应。
- 减少大量监听器造成的维护成本。

不要滥用事件委托。复杂交互、独立组件、需要清晰生命周期的逻辑，仍然建议放在组件内部处理。

## 阻止默认行为和阻止冒泡

阻止默认行为：

```ts
form.addEventListener('submit', (event) => {
  event.preventDefault()
  submitForm()
})
```

阻止冒泡：

```ts
dropdown.addEventListener('click', (event) => {
  event.stopPropagation()
})
```

不要随手 `stopPropagation`。它会影响外层组件、埋点、全局快捷键、点击外部关闭等逻辑。

## once、passive、signal

`addEventListener` 支持 options。

只执行一次：

```ts
button.addEventListener('click', handleClick, {
  once: true
})
```

滚动事件中声明不调用 `preventDefault`：

```ts
window.addEventListener('scroll', handleScroll, {
  passive: true
})
```

用 `AbortController` 批量取消监听：

```ts
const controller = new AbortController()

window.addEventListener('resize', handleResize, {
  signal: controller.signal
})

document.addEventListener('keydown', handleKeydown, {
  signal: controller.signal
})

controller.abort()
```

这在组件卸载清理多个原生事件时很有用。

## Vue / React 项目里的事件清理

Vue：

```ts
const controller = new AbortController()

onMounted(() => {
  window.addEventListener('resize', handleResize, {
    signal: controller.signal
  })
})

onBeforeUnmount(() => {
  controller.abort()
})
```

React：

```tsx
useEffect(() => {
  const controller = new AbortController()

  window.addEventListener('resize', handleResize, {
    signal: controller.signal
  })

  return () => {
    controller.abort()
  }
}, [])
```

组件卸载时不清理监听，最常见后果是重复触发、状态更新异常和内存泄漏。

## 点击外部关闭

下拉菜单常见需求：

```ts
function useClickOutside(container: HTMLElement, close: () => void) {
  const controller = new AbortController()

  document.addEventListener(
    'click',
    (event) => {
      const target = event.target as Node

      if (!container.contains(target)) {
        close()
      }
    },
    { signal: controller.signal }
  )

  return () => controller.abort()
}
```

注意：

- 需要判断点击是否发生在容器内部。
- 弹窗、传送门、iframe 会让边界更复杂。
- 不要用全局 `stopPropagation` 解决所有问题。

## 实际项目常见问题

### 1. 事件监听重复触发

**原因**

组件每次打开都绑定事件，但关闭或卸载时没有清理。

**解决方案**

- 保留函数引用。
- 在生命周期卸载阶段 remove。
- 或使用 `AbortController` 批量取消。

### 2. 点击子元素时拿不到预期数据

**原因**

`event.target` 可能是内部的 `span`、`svg`、`path`，不是你以为的按钮。

**解决方案**

使用 `closest` 找到业务节点：

```ts
const button = (event.target as HTMLElement).closest<HTMLButtonElement>('button[data-action]')
```

### 3. 滚动事件导致卡顿

**原因**

滚动频繁触发，回调里做了复杂计算或同步布局读取。

**解决方案**

- 节流。
- 使用 `passive: true`。
- 优先考虑 IntersectionObserver。
- 避免滚动中频繁读写布局。

### 4. stopPropagation 破坏外层逻辑

**现象**

外层点击关闭、埋点或快捷键突然失效。

**解决方案**

先分析事件边界，优先用条件判断，不要随手阻止冒泡。

### 5. 移动端 click 延迟或触摸问题

现代浏览器已经大幅改善 click 延迟，但复杂拖拽、滑动、缩放仍然要关注 pointer/touch 事件差异。项目里优先使用成熟手势库或框架组件，不要临时拼一堆事件。

## 最佳实践

- 组件外部绑定的事件必须有清理策略。
- 列表大量子项适合事件委托。
- 区分 `target` 和 `currentTarget`。
- 不随手阻止冒泡。
- 滚动、输入、resize 这类高频事件要节流或换更合适的 API。
- 复杂交互优先使用成熟组件或封装 composable/hook。

## 学习检查

学完本节后，你应该能回答：

- 事件冒泡和捕获的顺序是什么。
- `target` 和 `currentTarget` 有什么区别。
- 事件委托适合什么场景。
- 为什么组件卸载时要清理原生事件。
- `once`、`passive`、`signal` 分别解决什么问题。

## 参考资料

- [MDN: EventTarget.addEventListener](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener)
- [MDN: EventTarget.removeEventListener](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener)
- [MDN: Event bubbling](https://developer.mozilla.org/en-US/docs/Learn_web_development/Core/Scripting/Event_bubbling)
- [MDN: AbortSignal abort event](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal/abort_event)

## 下一步学习

继续学习 [正则表达式](/javascript/regular-expressions)，处理表单校验、搜索、替换和日志解析场景。
