# 事件循环

## 适合谁看

适合已经会写 Promise、async/await、定时器和接口请求，但不清楚这些代码为什么按某个顺序执行的人：

- `setTimeout(..., 0)` 为什么不是立刻执行。
- `Promise.then` 为什么比定时器更早执行。
- 页面为什么会因为一段 JavaScript 计算而卡住。
- Vue 的 `nextTick`、React 的状态更新为什么不是马上反映到 DOM。
- 接口请求明明是异步，为什么页面仍然会卡。

事件循环是 JavaScript 运行机制的核心。理解它后，你能更准确地判断异步顺序、页面渲染时机和性能问题。

## 先记住一句话

JavaScript 主线程一次只能执行一段同步代码。

当同步代码执行完后，运行环境会从任务队列中取出下一段要执行的代码。这个不断取任务、执行任务、处理微任务、再进入下一轮的过程，就是事件循环。

## 同步代码先执行

```ts
console.log('A')

setTimeout(() => {
  console.log('B')
}, 0)

console.log('C')
```

输出：

```text
A
C
B
```

`setTimeout` 不是把回调立刻插入当前执行栈，而是等当前同步代码结束后，再进入后续任务。

## 任务和微任务

先沿时间顺序观察同步脚本、Promise 微任务、渲染机会和定时器任务分别在哪个检查点执行。

<DocFigure
  src="/images/javascript/event-loop-devtools.webp"
  alt="JavaScript 同步脚本执行后清空 Promise 微任务，浏览器获得渲染机会，随后执行 setTimeout"
  caption="微任务在当前任务结束后批量清空；递归创建微任务也可能推迟渲染。"
  :width="1440"
  :height="900"
/>

图是一次典型执行路径，不代表浏览器每轮都一定渲染；准确顺序仍要结合调用栈、任务来源和 Performance 录制判断。

常见分类：

| 类型 | 常见来源 |
| --- | --- |
| 任务 | `setTimeout`、用户点击、网络事件、脚本执行 |
| 微任务 | `Promise.then`、`queueMicrotask`、`MutationObserver` |

简化顺序：

```text
执行一个任务
↓
执行当前产生的所有微任务
↓
浏览器可能进行渲染
↓
进入下一个任务
```

示例：

```ts
console.log('A')

setTimeout(() => {
  console.log('B')
}, 0)

Promise.resolve().then(() => {
  console.log('C')
})

console.log('D')
```

输出：

```text
A
D
C
B
```

原因：

1. `A` 和 `D` 是同步代码。
2. Promise 回调进入微任务队列。
3. 定时器回调进入任务队列。
4. 当前任务结束后先清空微任务，所以 `C` 早于 `B`。

## async/await 的执行顺序

`await` 后面的代码可以理解为进入 Promise 的后续微任务。

```ts
async function run() {
  console.log('A')
  await Promise.resolve()
  console.log('B')
}

run()
console.log('C')
```

输出：

```text
A
C
B
```

`await` 之前仍然是同步执行，`await` 之后等待当前任务结束后继续。

## 为什么页面会卡

浏览器渲染也需要主线程参与。只要 JavaScript 长时间占用主线程，用户输入、动画和渲染都会被阻塞。

```ts
function heavyCalculate() {
  const start = Date.now()

  while (Date.now() - start < 3000) {
    // 模拟长任务
  }
}

heavyCalculate()
```

这 3 秒内页面无法正常响应。

解决方向：

- 拆分长任务。
- 使用 Web Worker。
- 减少同步循环和大数据一次性处理。
- 把大列表渲染改成虚拟滚动。
- 用 Performance 面板定位 Long Task。

## 渲染时机和 DOM 更新

很多前端框架不会每次状态变化都立刻更新 DOM，而是把多次变化合并，在合适时机统一更新。

Vue 示例：

```ts
count.value += 1

console.log(document.querySelector('#count')?.textContent)

await nextTick()

console.log(document.querySelector('#count')?.textContent)
```

状态改了，不代表 DOM 立即完成更新。需要等待框架调度完成。

## 定时器不是精准时钟

`setTimeout(fn, 100)` 表示至少等待 100ms 后可以进入任务队列，不保证 100ms 后立刻执行。

如果主线程正忙，定时器会延后。

```ts
setTimeout(() => {
  console.log('timer')
}, 100)

heavyCalculate()
```

如果 `heavyCalculate` 执行 3 秒，定时器也只能等它结束后再执行。

## 实际项目常见问题

### 1. Promise 回调比定时器先执行

这是微任务和任务顺序导致的，不是浏览器随机行为。

排查异步顺序时，先区分：

- 当前同步代码。
- 微任务。
- 下一轮任务。

### 2. 接口是异步的，页面还是卡

接口等待本身不占主线程，但接口回来后你可能做了大量同步处理：

```ts
const result = await fetchLargeList()

list.value = normalizeHugeTree(result)
```

如果 `normalizeHugeTree` 很重，页面仍然会卡。

解决方式：

- 后端分页。
- 前端分片处理。
- Web Worker 处理复杂转换。
- 减少一次性渲染节点数量。

### 3. loading 状态没有显示出来

```ts
loading.value = true
heavyCalculate()
loading.value = false
```

主线程一直被同步计算占着，浏览器没有机会渲染 loading。

解决方向：

- 把重计算放到下一轮任务或 Worker。
- 先让 DOM 更新，再开始重计算。
- 更根本地减少同步计算。

### 4. 多次状态更新只渲染一次

这是框架批处理的结果，通常是好事。不要依赖每次赋值后 DOM 立刻变化。

需要读取更新后的 DOM 时，使用框架提供的更新完成 API，例如 Vue 的 `nextTick`。

### 5. 无限微任务导致页面不渲染

如果不断递归创建微任务，浏览器可能长时间没有机会进入渲染阶段。

```ts
function loop() {
  Promise.resolve().then(loop)
}

loop()
```

不要用无限微任务做轮询。

## 最佳实践

- 把同步、任务、微任务分清楚。
- 不要用 `setTimeout(..., 0)` 当作可靠的立即执行工具。
- 大计算不要阻塞主线程。
- 需要读取更新后的 DOM 时等待框架更新完成。
- 线上卡顿优先用 Performance 面板找长任务。
- 数据量大时从分页、虚拟滚动、Worker 三个方向考虑。

## 学习检查

学完本节后，你应该能回答：

- 同步代码、任务、微任务的执行顺序是什么。
- 为什么 Promise 比 setTimeout 更早执行。
- 为什么异步请求回来后仍然可能卡页面。
- 为什么 loading 有时来不及显示。
- Vue 或 React 中为什么 DOM 更新不是立即可读。

## 参考资料

- [MDN: JavaScript execution model](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Execution_model)
- [MDN: In depth: Microtasks and the JavaScript runtime environment](https://developer.mozilla.org/en-US/docs/Web/API/HTML_DOM_API/Microtask_guide/In_depth)

## 下一步学习

继续学习 [错误处理](/javascript/error-handling)，把异步顺序和失败恢复结合起来。
