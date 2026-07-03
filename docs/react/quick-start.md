# React 快速开始

## 适合谁看

适合第一次写 React，希望先跑通一个可交互页面的学习者。

## 创建项目

常见方式是使用 Vite：

```bash
npm create vite@latest my-react-app -- --template react-ts
cd my-react-app
npm install
npm run dev
```

## 第一个组件

```tsx
import { useState } from 'react'

export function Counter() {
  const [count, setCount] = useState(0)

  return (
    <button type="button" onClick={() => setCount(count + 1)}>
      点击次数：{count}
    </button>
  )
}
```

React 组件本质上是返回 UI 的函数。

## 事件处理

```tsx
function UserSearch() {
  const [keyword, setKeyword] = useState('')

  function search() {
    console.log(keyword)
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault()
        search()
      }}
    >
      <input value={keyword} onChange={(event) => setKeyword(event.target.value)} />
      <button type="submit">搜索</button>
    </form>
  )
}
```

## 实际项目常见问题

### 1. 输入框无法输入

**原因**

写了 `value`，但没有正确更新 state。

**解决方案**

```tsx
<input value={keyword} onChange={(event) => setKeyword(event.target.value)} />
```

### 2. 点击后状态没按预期连续增加

如果依赖上一次状态，使用函数式更新：

```tsx
setCount((count) => count + 1)
```

## 下一步

继续学习 [组件与 JSX](/react/component-jsx)。
