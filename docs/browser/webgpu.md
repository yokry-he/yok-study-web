# WebGPU

## 适合谁看

适合已经了解浏览器基础能力，并开始接触这些场景的人：

- 页面需要复杂 2D/3D 图形渲染。
- 要在浏览器里做大量并行计算。
- 项目涉及数据可视化、地图、游戏、图像处理或 AI 推理。
- 听过 WebGL、Canvas、Three.js，但不知道 WebGPU 的定位。

WebGPU 是浏览器访问 GPU 能力的新一代 API。它比 WebGL 更贴近现代 GPU 模型，也支持图形渲染之外的通用计算。

## 它解决什么问题

CPU 适合处理通用逻辑，GPU 更适合处理大量相似计算，例如像素、矩阵、顶点、纹理和并行数据处理。

WebGPU 的目标是让 Web 应用可以更直接、更高效地使用 GPU：

| 能力 | 说明 |
| --- | --- |
| 图形渲染 | 绘制复杂 2D/3D 场景 |
| 通用计算 | 用 GPU 做并行计算 |
| 现代 GPU 模型 | 更接近 Vulkan、Metal、Direct3D 12 等现代图形 API |
| 更明确的资源管理 | 通过 device、buffer、pipeline、command 等概念组织执行 |

你可以把 WebGPU 理解成浏览器里的底层图形和并行计算能力。普通业务页面不需要直接写 WebGPU。

## 和 Canvas、WebGL、Three.js 的关系

| 技术 | 更适合 |
| --- | --- |
| Canvas 2D | 简单图形、图表、图片绘制 |
| WebGL | 成熟 3D 渲染、已有生态项目 |
| Three.js | 高层 3D 场景开发 |
| WebGPU | 新一代高性能渲染、GPU 计算、底层能力探索 |

多数业务项目不应该直接从 WebGPU 开始。更常见的路径是：

```text
简单图形：Canvas / SVG
↓
复杂 3D：Three.js
↓
底层性能或新能力：WebGPU
```

如果团队不是做图形、游戏、可视化或计算平台，了解 WebGPU 的定位就够了。

## 基本概念

WebGPU 的概念比普通 Web API 更底层。

| 概念 | 含义 |
| --- | --- |
| `navigator.gpu` | WebGPU 的入口 |
| Adapter | 当前设备可用的 GPU 适配器 |
| Device | 与 GPU 交互的逻辑设备 |
| Buffer | 存放数据的 GPU 缓冲区 |
| Texture | 存放图像或纹理数据 |
| Shader | 在 GPU 上执行的程序 |
| Pipeline | 描述 GPU 如何处理数据 |
| Command Encoder | 记录 GPU 命令 |
| Queue | 提交命令给 GPU 执行 |

这套模型的核心思想是：提前准备资源和执行管线，再提交命令给 GPU。

## 最小初始化流程

概念示例：

```ts
if (!navigator.gpu) {
  throw new Error('当前浏览器不支持 WebGPU')
}

const adapter = await navigator.gpu.requestAdapter()

if (!adapter) {
  throw new Error('当前设备没有可用 GPU adapter')
}

const device = await adapter.requestDevice()
```

如果要渲染到页面，通常还需要 canvas context：

```ts
const canvas = document.querySelector<HTMLCanvasElement>('#canvas')!
const context = canvas.getContext('webgpu')!
const format = navigator.gpu.getPreferredCanvasFormat()

context.configure({
  device,
  format,
  alphaMode: 'premultiplied'
})
```

真实渲染还需要 shader、pipeline、buffer、command encoder，这也是 WebGPU 学习曲线比较陡的原因。

## WGSL 是什么

WebGPU 使用 WGSL 编写 shader。

简单理解：

```text
TypeScript 运行在 JavaScript 引擎里
WGSL shader 运行在 GPU 上
```

示意：

```wgsl
@vertex
fn vertexMain(@builtin(vertex_index) index: u32) -> @builtin(position) vec4f {
  var positions = array<vec2f, 3>(
    vec2f(0.0, 0.5),
    vec2f(-0.5, -0.5),
    vec2f(0.5, -0.5)
  );

  let position = positions[index];
  return vec4f(position, 0.0, 1.0);
}
```

初学时不用急着掌握全部 WGSL 语法。先理解 shader 是在 GPU 上并行执行的小程序。

## WebGPU 的项目落地方式

实际项目里常见三种方式：

| 方式 | 适合场景 |
| --- | --- |
| 直接写 WebGPU | 学习底层、做引擎、做性能敏感能力 |
| 使用图形库封装 | 产品级 3D、数据可视化、地图 |
| 使用现成 SDK | AI 推理、图像处理、专业工具 |

对于前端工程师，重点不是立刻手写完整渲染管线，而是知道：

- WebGPU 能解决什么问题。
- 它的运行和部署条件是什么。
- 什么时候应该用高层库。
- 如何排查兼容性和性能问题。

## 能力检测和降级

WebGPU 仍然需要认真做能力检测。

```ts
export async function createGpuDevice() {
  if (!navigator.gpu) {
    return null
  }

  const adapter = await navigator.gpu.requestAdapter()

  if (!adapter) {
    return null
  }

  return adapter.requestDevice()
}
```

业务代码不要直接假设所有用户都有 WebGPU。

降级策略：

- 不支持 WebGPU 时使用 WebGL。
- 不支持复杂渲染时切回 Canvas 或静态图。
- AI 推理场景切回 WebAssembly 或服务端推理。
- 给用户明确提示，而不是白屏。

## 安全上下文和兼容性

WebGPU 属于比较敏感的硬件能力，通常需要安全上下文。上线前要检查：

- 是否使用 HTTPS。
- 浏览器版本是否支持。
- 目标用户设备是否支持。
- 是否运行在受限 WebView 或企业浏览器中。
- 是否需要特殊 feature policy 或浏览器开关。

内部后台、数据大屏和专业工具上线前，应先统计真实用户浏览器环境。

## 性能排查思路

WebGPU 性能问题不能只看 JavaScript 时间。

排查方向：

- CPU 是否在频繁创建资源。
- 是否每帧重复创建 pipeline、buffer、shader。
- 数据是否频繁在 CPU 和 GPU 间复制。
- canvas 尺寸是否远超实际展示尺寸。
- 是否没有控制帧率。
- GPU 任务是否太重导致掉帧。

常见原则：

- 能复用的 GPU 资源尽量复用。
- 初始化阶段创建 pipeline。
- 每帧只更新必要数据。
- 纹理和 buffer 要控制大小。
- 对低端设备提供简化效果。

## 实际项目常见问题

### 1. `navigator.gpu` 是 undefined

**原因**

当前浏览器、系统、设备、上下文或 WebView 不支持 WebGPU。

**解决方案**

- 使用 HTTPS。
- 检查目标浏览器版本。
- 增加能力检测。
- 提供 WebGL、Canvas 或服务端降级。

### 2. 本地可用，线上不可用

**原因**

线上可能不是安全上下文，或被 iframe、WebView、企业策略限制。

**排查**

- 看页面是否 HTTPS。
- 看浏览器控制台安全错误。
- 检查是否嵌入在 iframe。
- 检查目标用户浏览器和设备。

### 3. 页面加载后白屏

**原因**

初始化失败后没有 fallback，或者 canvas 尺寸、context 配置、shader 编译失败。

**解决方案**

- 初始化过程全部加错误捕获。
- 显示明确错误状态。
- 日志记录 adapter/device/context 创建结果。
- 对 shader 编译错误做开发期暴露。

### 4. GPU 页面让电脑风扇狂转

**原因**

渲染循环没有停止，或者每帧执行过重。

**解决方案**

- 页面不可见时暂停。
- 组件卸载时取消 `requestAnimationFrame`。
- 降低分辨率或效果质量。
- 控制帧率。

### 5. 和 Vue/React 集成后资源泄漏

**原因**

组件卸载时没有停止循环，也没有释放或断开相关资源引用。

**解决方案**

在组件生命周期里明确管理：

- 初始化。
- resize。
- render loop。
- pause/resume。
- dispose。

Vue 示例结构：

```ts
onMounted(async () => {
  await initWebGpu()
  startRenderLoop()
})

onBeforeUnmount(() => {
  stopRenderLoop()
  disposeResources()
})
```

## 最佳实践

- 普通业务页面不要直接引入 WebGPU。
- 产品项目优先考虑成熟图形库或 SDK。
- 初始化必须做能力检测和 fallback。
- WebGPU 渲染循环要和页面生命周期绑定。
- 不要每帧创建昂贵资源。
- 移动端和低端设备要有降级质量。
- 关键图形能力上线前做真实设备验证。

## 学习检查

学完本节后，你应该能回答：

- WebGPU 和 WebGL、Canvas 的区别是什么。
- 为什么 WebGPU 需要 adapter、device、pipeline。
- 为什么不应该在所有项目里直接写 WebGPU。
- WebGPU 页面白屏时应该查什么。
- Vue/React 集成 GPU 渲染时为什么要处理卸载。

## 参考资料

- [MDN: WebGPU API](https://developer.mozilla.org/en-US/docs/Web/API/WebGPU_API)
- [MDN: GPU](https://developer.mozilla.org/en-US/docs/Web/API/GPU)
- [W3C: WebGPU](https://www.w3.org/TR/webgpu/)
- [Chrome for Developers: Overview of WebGPU](https://developer.chrome.com/docs/web-platform/webgpu/overview)
- [WebGPU.org](https://webgpu.org/)

## 下一步学习

继续学习 [浏览器自动化调试](/browser/browser-automation-debugging)，把页面问题从人工点击排查升级为可重复验证。
