# WebAssembly

## 适合谁看

适合已经能写前端业务页面，但开始遇到这些场景的人：

- 页面里有大量计算，JavaScript 执行时间明显偏长。
- 需要在浏览器里复用 Rust、C、C++ 等语言写好的能力。
- 需要处理图片、音视频、压缩、加密、CAD、地图、游戏、AI 推理等重计算任务。
- 听过 Wasm，但不知道它和 JavaScript、Worker、npm 包到底是什么关系。

WebAssembly 通常简称 Wasm。它不是用来替代 JavaScript 写普通页面逻辑的，而是把适合高性能计算的代码编译成浏览器可以执行的二进制模块，再由 JavaScript 调用。

## 它解决什么问题

普通前端业务里，JavaScript 已经足够好用。但当任务变成大量循环、矩阵计算、像素处理、音频处理或算法库移植时，JavaScript 可能不是最合适的载体。

WebAssembly 的价值主要有三点：

| 价值 | 说明 |
| --- | --- |
| 性能更稳定 | 适合计算密集型逻辑，执行模型更接近底层 |
| 复用已有库 | 可以把 Rust、C、C++ 等生态里的成熟库带到浏览器 |
| 沙箱执行 | Wasm 在浏览器安全沙箱中运行，不能随意访问系统资源 |

不要把 Wasm 理解成“更快的 JavaScript”。它更像浏览器里的一个高性能计算模块，适合被 JavaScript 调度。

## 基本工作流

典型流程：

```text
Rust / C / C++
↓ 编译
.wasm 文件
↓ 加载
JavaScript WebAssembly API
↓ 调用
页面业务逻辑
```

项目里通常不是手写 Wasm，而是使用语言工具链生成。例如 Rust 常见工具链会生成 `.wasm` 文件和一层 JavaScript 胶水代码，方便前端直接 import。

## 最小调用模型

浏览器提供了 `WebAssembly` 全局对象，可以加载和实例化 Wasm 模块。

概念示例：

```ts
const response = await fetch('/modules/calc.wasm')
const bytes = await response.arrayBuffer()
const result = await WebAssembly.instantiate(bytes)

const add = result.instance.exports.add as CallableFunction

console.log(add(1, 2))
```

真实项目里通常不会直接这样写，因为：

- 工具链会生成类型和初始化代码。
- Vite、Webpack 等构建工具可能需要插件处理 `.wasm`。
- 复杂模块还涉及内存、字符串、对象传递。

但这个示例能说明核心关系：JavaScript 负责加载模块，Wasm 暴露函数，JavaScript 调用这些函数。

## JavaScript 和 Wasm 怎么分工

建议这样分工：

| 层 | 适合负责 |
| --- | --- |
| JavaScript / TypeScript | UI、状态、请求、路由、权限、用户交互 |
| Web Worker | 把重计算移出主线程 |
| WebAssembly | 计算密集型算法、底层库、跨语言复用 |

很多项目会把 Wasm 放到 Worker 里执行：

```text
Vue / React 页面
↓ postMessage
Web Worker
↓ 调用
WebAssembly 模块
↓ 返回结果
页面更新状态
```

这样能避免 Wasm 计算阻塞主线程，减少页面卡顿。

## 什么时候值得用

适合：

- 图片裁剪、滤镜、压缩、格式转换。
- 音视频编解码、波形分析、降噪。
- 加密、哈希、压缩算法。
- 复杂图形、游戏、物理模拟。
- 浏览器端 CAD、GIS、科学计算。
- 已经存在成熟 Rust/C/C++ 库，希望复用到 Web。
- 需要离线执行的高性能逻辑。

不适合：

- 普通表单、列表、权限、管理后台页面。
- 简单数据处理。
- 只是为了“看起来高级”。
- 团队没有工具链维护能力。
- 需要频繁和 DOM 交互的逻辑。

Wasm 和 DOM 不是直接配合得很好。DOM 操作仍然交给 JavaScript 更合理。

## 与 Vite 项目的关系

在现代前端项目里，Wasm 通常通过构建工具接入。

常见方式：

```ts
import init, { calculate } from './pkg/calc'

await init()

const result = calculate(100)
```

具体写法取决于 Wasm 的生成工具和构建配置。

在 Vite 项目里要重点确认：

- `.wasm` 文件是否能被正确打包。
- 生产环境资源路径是否正确。
- 是否需要异步初始化。
- 是否要把 Wasm 放到 Worker 中。
- 服务端是否返回正确的 `Content-Type`。

## 资源加载和 MIME

线上常见问题是本地可用，部署后加载失败。

排查顺序：

1. Network 面板看 `.wasm` 是否返回 200。
2. 看响应体是不是被网关错误返回成 HTML。
3. 看资源路径是否受 base path 影响。
4. 看响应头是否有合适的 MIME 类型。
5. 看是否被 CSP 拦截。

常见服务器配置需要支持：

```http
Content-Type: application/wasm
```

如果 `.wasm` 请求返回的是 `index.html`，通常是前端路由 fallback 把资源路径错误回退了。

## 内存和数据传递

Wasm 和 JavaScript 之间传递数字比较简单，但字符串、数组、对象会复杂很多。

原因是 Wasm 有自己的线性内存模型。复杂数据通常需要：

- JavaScript 把数据写入 Wasm 内存。
- Wasm 返回指针或结果。
- JavaScript 再从内存读取结果。

工具链会帮你隐藏很多细节，但你仍然要知道跨边界调用不是免费的。

项目经验：

- 不要在 JS 和 Wasm 之间频繁传递小对象。
- 尽量批量传入数据，一次计算后批量返回。
- 高频调用要关注序列化和内存复制成本。
- 大数组适合用 `ArrayBuffer`、`TypedArray` 一类结构。

## 和 Web Worker 搭配

如果 Wasm 计算耗时明显，建议放到 Worker 里。

主线程：

```ts
const worker = new Worker(new URL('./calc.worker.ts', import.meta.url), {
  type: 'module'
})

worker.postMessage({
  type: 'calculate',
  payload: [1, 2, 3]
})

worker.onmessage = event => {
  console.log(event.data)
}
```

Worker 内部再初始化 Wasm：

```ts
import init, { calculate } from './pkg/calc'

let ready = false

self.onmessage = async event => {
  if (!ready) {
    await init()
    ready = true
  }

  if (event.data.type === 'calculate') {
    self.postMessage(calculate(event.data.payload))
  }
}
```

这样页面不会因为重计算而失去响应。

## 实际项目常见问题

### 1. Wasm 文件部署后 404

**原因**

构建产物路径、CDN 路径或前端 `base` 配置不一致。

**解决方案**

- 用 Network 面板确认真实请求路径。
- 检查 Vite `base`。
- 确认 `.wasm` 是否进入构建产物。
- 确认服务器静态资源目录包含该文件。

### 2. 请求返回 200，但加载报错

**原因**

请求可能返回的是 HTML 错误页，而不是 Wasm 二进制。

**解决方案**

打开 Network 的 Response 或 Preview，确认内容不是 `index.html`、登录页或错误页。

### 3. 页面依然卡顿

**原因**

Wasm 只是提升计算执行效率，不代表不会阻塞主线程。只要在主线程同步执行长任务，页面仍然会卡。

**解决方案**

把计算放入 Web Worker，并减少 JS 与 Wasm 的频繁跨边界调用。

### 4. 开发环境可以，生产环境不行

**排查**

- 生产资源路径是否正确。
- `.wasm` 是否被 CDN 正确缓存。
- 服务器是否设置了 `application/wasm`。
- 是否有 CSP 限制。
- 是否依赖了开发服务器特有的路径。

### 5. 引入后包体变大

**原因**

Wasm 文件本身可能很大，尤其是编解码、图像处理、AI 推理类库。

**解决方案**

- 按需加载。
- 路由级懒加载。
- Worker 内懒初始化。
- 给加载过程增加进度或 fallback。
- 检查是否引入了不需要的能力。

## 最佳实践

- 先确认瓶颈来自计算，不要一开始就引入 Wasm。
- 普通业务逻辑继续用 TypeScript 写。
- 重计算优先放到 Worker。
- 减少 JS 与 Wasm 的频繁小数据交换。
- 给 Wasm 初始化设计 loading、失败重试和降级方案。
- 生产环境重点验证资源路径、MIME、缓存和 CSP。
- 记录工具链版本和构建命令，避免后续维护没人知道怎么重新生成。

## 学习检查

学完本节后，你应该能回答：

- WebAssembly 适合解决什么问题。
- 为什么 Wasm 不能替代 Vue 或 React。
- 为什么重计算最好放到 Worker。
- 为什么 JS 和 Wasm 之间频繁传对象会有成本。
- 部署后 `.wasm` 加载失败应该看哪些证据。

## 参考资料

- [MDN: WebAssembly](https://developer.mozilla.org/en-US/docs/WebAssembly)
- [MDN: Using the WebAssembly JavaScript API](https://developer.mozilla.org/en-US/docs/WebAssembly/Guides/Using_the_JavaScript_API)
- [MDN: WebAssembly concepts](https://developer.mozilla.org/en-US/docs/WebAssembly/Guides/Concepts)
- [W3C: WebAssembly Web API](https://www.w3.org/TR/wasm-web-api-2/)
- [WebAssembly.org](https://webassembly.org/)

## 下一步学习

继续学习 [WebGPU](/browser/webgpu)，理解浏览器如何使用 GPU 进行图形渲染和通用计算。
