# tsconfig 与工程配置

## 适合谁看

适合已经能写 TypeScript，但对项目里的 `tsconfig.json`、路径别名、构建检查和严格模式还不够清楚的人：

- 不知道 `strict` 开不开有什么影响。
- 编辑器不报错，但构建时报类型错误。
- 路径别名在 TypeScript 里能识别，运行时却找不到。
- Vue 项目里 `vue-tsc`、`tsc`、Vite 的关系分不清。
- Monorepo 或多包项目类型检查越来越慢。

TypeScript 不只是语法。真正进入项目后，配置决定了类型检查范围、严格程度、模块解析和构建质量门槛。

## tsconfig 是什么

`tsconfig.json` 表示一个 TypeScript 项目的根配置。它告诉 TypeScript：

- 检查哪些文件。
- 使用哪些编译选项。
- 如何解析模块。
- 是否生成代码。
- 是否启用严格检查。

最小示例：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "jsx": "preserve",
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"]
}
```

不同框架和构建工具会有自己的推荐配置，不能盲目复制旧项目。

## Vue + Vite 项目常见配置

Vue 项目通常会有：

```text
tsconfig.json
tsconfig.app.json
tsconfig.node.json
```

常见分工：

| 文件 | 作用 |
| --- | --- |
| `tsconfig.json` | 根配置和 references |
| `tsconfig.app.json` | 前端源码类型检查 |
| `tsconfig.node.json` | Vite、脚本、配置文件类型检查 |

这样做能把浏览器代码和 Node 配置代码分开，避免环境类型互相污染。

## strict 要不要开

建议新项目开启：

```json
{
  "compilerOptions": {
    "strict": true
  }
}
```

`strict` 会启用一组更严格的检查，能提前发现大量空值、隐式 any、函数参数和类型不匹配问题。

老项目可以分阶段开启：

1. 先减少明显 `any`。
2. 补接口和表单类型。
3. 打开 `strictNullChecks`。
4. 逐步打开完整 `strict`。
5. 用 CI 阻止新增类型债务。

不要一次性打开严格模式后，用大量 `as any` 把错误压下去。

## include 和 exclude

`include` 决定哪些文件进入类型检查。

```json
{
  "include": ["src/**/*.ts", "src/**/*.vue"]
}
```

如果文件不在 include 里，编辑器或构建可能不会按预期检查。

常见问题：

- 新增 `scripts/` 没有类型检查。
- `.vue` 没被包含。
- 测试文件使用了不同环境类型。
- 生成文件进入检查导致很慢。

建议按用途拆配置，而不是把所有文件都塞进一个 tsconfig。

## paths 和运行时解析

路径别名失败时先确定实际加载的配置和解析轨迹。图中 `@/api/user` 命中本地 paths，而 `@shared/missing` 的所有候选文件都不存在。

<DocFigure
  src="/images/typescript/tsconfig-trace.webp"
  alt="TypeScript 模块解析报告展示 tsconfig 继承、paths 命中、包 exports 与缺失模块候选"
  caption="使用 tsc --showConfig 确认最终配置，再用 --traceResolution 查看每个候选路径。"
  :width="1440"
  :height="900"
/>

`paths` 只影响 TypeScript 解析，不自动改写浏览器或 Node 的运行时路径；Vite、测试工具和部署环境也要配置同一别名规则。

TypeScript 的 `paths` 只负责类型检查和编辑器解析，不一定改变运行时或打包解析。

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  }
}
```

Vite 里还需要对应 alias：

```ts
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
```

如果只配 tsconfig，不配 Vite，编辑器可能不报错，但运行或构建会失败。

## noEmit 和构建检查

很多前端项目不使用 `tsc` 输出 JS，而由 Vite、esbuild、Babel 或框架工具处理构建。

TypeScript 只做类型检查：

```json
{
  "compilerOptions": {
    "noEmit": true
  }
}
```

常见脚本：

```json
{
  "scripts": {
    "typecheck": "vue-tsc --noEmit",
    "build": "vue-tsc --noEmit && vite build"
  }
}
```

构建链路里必须有真实类型检查。只跑 Vite build 可能无法发现完整类型错误。

## moduleResolution

现代 Vite 项目常见：

```json
{
  "compilerOptions": {
    "moduleResolution": "Bundler"
  }
}
```

它更贴近现代打包器解析方式。

Node 脚本、库项目、旧项目可能使用不同策略。遇到模块找不到时，不要只改 import，先确认：

- `module`
- `moduleResolution`
- `types`
- `paths`
- package 的 `exports`
- 当前文件属于哪个 tsconfig

## types 和环境类型

浏览器项目和 Node 脚本需要不同全局类型。

前端源码：

```json
{
  "compilerOptions": {
    "lib": ["ES2022", "DOM", "DOM.Iterable"]
  }
}
```

Node 配置：

```json
{
  "compilerOptions": {
    "types": ["node"]
  }
}
```

如果在浏览器代码里直接出现 Node 全局类型，可能说明配置边界混了。

## 项目引用

大型项目或 Monorepo 可以用 project references 拆分类型检查边界。

根配置：

```json
{
  "files": [],
  "references": [
    { "path": "./packages/ui" },
    { "path": "./packages/admin" }
  ]
}
```

适合：

- 多包仓库。
- 组件库 + 业务项目。
- 前端应用 + shared 工具包。
- 类型检查很慢，需要拆边界。

不要在小项目里过早引入复杂 references。

## 实际项目常见问题

### 1. 编辑器没报错，CI 报错

**原因**

编辑器加载的 tsconfig 和 CI 执行的检查脚本不同，或者本地没有跑完整 typecheck。

**解决方案**

- 明确 `npm run typecheck`。
- CI 和本地使用同一条命令。
- 检查文件是否进入 include。

### 2. 路径别名编辑器可用，浏览器运行失败

**原因**

只配置了 tsconfig paths，没有配置 Vite alias。

**解决方案**

tsconfig 和 Vite alias 保持一致。

### 3. 开 strict 后报错太多

**解决方案**

分阶段开启，不要一次性用 `any` 压错误。可以先从接口、表单、store、组件 props 这些收益最高的区域开始。

### 4. 第三方库类型污染全局

**原因**

`types` 或全局声明文件范围过大。

**解决方案**

- 缩小声明文件范围。
- 明确 tsconfig include。
- 不要把临时声明放到全局 `any`。

### 5. JSON import 报错

如果项目需要直接导入 JSON：

```json
{
  "compilerOptions": {
    "resolveJsonModule": true
  }
}
```

同时确认构建工具也支持对应导入方式。

## 推荐检查清单

新建或接手项目时检查：

- 是否有明确 typecheck 脚本。
- build 是否包含 typecheck。
- `strict` 当前状态和开启计划。
- `include` 是否覆盖 `.vue`、测试、脚本。
- `paths` 是否和构建工具 alias 对齐。
- 浏览器代码和 Node 配置是否分 tsconfig。
- CI 是否执行同一套检查。

## 最佳实践

- 新项目默认开启 `strict`。
- 类型检查命令要进入 CI。
- Vite alias 和 tsconfig paths 必须同步。
- 浏览器源码和 Node 配置拆分 tsconfig。
- 大项目再考虑 project references。
- 警告和类型错误要区分，真正阻断质量的是类型检查退出码。
- 配置变更要写清楚原因，避免后续没人敢改。

## 学习检查

学完本节后，你应该能回答：

- `tsconfig.json` 主要控制哪些事情。
- 为什么 Vite 项目只跑 build 不一定等于完整类型检查。
- `paths` 和 Vite alias 为什么要同时配置。
- `strict` 开启要注意什么。
- 什么时候需要 project references。

## 参考资料

- [TypeScript: What is a tsconfig.json](https://www.typescriptlang.org/docs/handbook/tsconfig-json.html)
- [TypeScript: TSConfig Reference](https://www.typescriptlang.org/tsconfig/)
- [TypeScript: Project References](https://www.typescriptlang.org/docs/handbook/project-references.html)
- [TypeScript TSConfig: references](https://www.typescriptlang.org/tsconfig/references.html)

## 下一步学习

继续学习 [Vue 项目集成](/typescript/vue-integration)，把类型配置落到组件、请求、表单和状态管理中。
