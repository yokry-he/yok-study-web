# React 常见问题

## 1. Effect 无限循环

### 症状

页面不断请求接口或不断重新渲染。

### 常见原因

- Effect 里更新了依赖状态。
- 依赖数组里放了每次渲染都会新建的对象或函数。

### 解决方案

- 能在事件里做的逻辑不要放 Effect。
- 派生数据直接计算。
- 需要稳定函数时再考虑 `useCallback`。

## 2. 状态更新后读取还是旧值

### 原因

state 更新会在下一次渲染体现。

### 解决方案

依赖上一次状态时使用函数式更新：

```tsx
setCount((count) => count + 1)
```

## 3. 列表输入框错乱

### 原因

使用 index 作为 key。

### 解决方案

使用稳定业务 id。

```tsx
items.map((item) => <ItemRow key={item.id} item={item} />)
```

## 4. 输入框无法编辑

### 原因

受控输入写了 `value`，但没有 `onChange` 更新。

### 解决方案

```tsx
<input value={keyword} onChange={(event) => setKeyword(event.target.value)} />
```

## 5. Hook 调用报错

### 原因

Hook 写在条件、循环、普通函数或事件函数里。

### 解决方案

Hook 只在组件或自定义 Hook 顶层调用。

## 6. 组件拆分后 props 层层传递

### 解决思路

- 先确认是否真的需要全局状态。
- 小范围可以组件组合。
- 跨多层共享可考虑 Context 或状态库。
- 不要一开始就把所有状态放全局。

## 7. 搜索结果显示旧数据

### 原因

多个请求并发，旧请求后返回，覆盖了新请求。

### 解决方案

使用请求序号、AbortController 或请求管理库。

```tsx
const requestIdRef = useRef(0)

async function fetchList() {
  const requestId = ++requestIdRef.current
  const result = await getList()

  if (requestId !== requestIdRef.current) return

  setList(result.items)
}
```

## 8. Context 更新导致很多组件重渲染

### 原因

Context 中放了高频变化状态，或 Provider value 每次都是新对象。

### 解决方案

- 拆分 Context。
- 高频状态下沉。
- 必要时使用状态库。

## 9. 表单关闭后还有旧数据

### 原因

关闭或打开新增时没有重置表单。

### 解决方案

使用默认值函数：

```tsx
setForm(createDefaultForm())
```

## 快速排查表

| 问题 | 优先检查 |
| --- | --- |
| 无限请求 | useEffect 依赖 |
| 输入不动 | value + onChange |
| 列表错乱 | key 是否稳定 |
| Hook 报错 | 是否顶层调用 |
| 旧状态 | 是否需要函数式更新 |
| 旧请求覆盖 | 是否需要请求序号或取消 |
| Context 重渲染 | 是否状态放太大 |

## 最佳实践

- 遵守 Hooks 规则。
- 不滥用 Effect。
- key 使用稳定 id。
- 状态放在最近共同父组件。
- 页面、组件、hooks、services 分层清楚。
