# TypeScript 学习导览

## 适合谁看

适合已经会写 JavaScript，准备在 Vue、React、Node 或大型前端项目中使用 TypeScript 的学习者。

TypeScript 的价值不是“写更多代码”，而是提前描述数据形状，让编辑器和构建工具在运行前发现问题。真实项目里，接口数据、表单、权限、路由、状态和组件 API 都非常依赖类型。

## 你会学到什么

- 基础类型、联合类型和类型推断。
- interface 和 type 的使用场景。
- 泛型如何表达“可复用但仍有类型约束”的逻辑。
- 类型收窄和类型守卫如何处理接口、路由、LocalStorage 等不可信数据。
- 工具类型如何服务表单、接口 payload 和页面 ViewModel，同时避免类型体操失控。
- tsconfig 如何影响严格检查、路径别名、构建质量和项目边界。
- TypeScript 如何集成到 Vue 组件、props、emits、ref 和 API 请求中。
- 项目里 any 滥用、类型过度复杂、接口字段变化如何处理。

## 学习顺序

<LearningPath :steps="[
  { title: '图解 TypeScript 核心概念', description: '先理解类型检查位置、类型推导、interface/type、联合类型收窄、泛型和 Vue 类型流。', link: '/typescript/visual-guide', badge: '图解' },
  { title: '基础类型', description: '理解 string、number、boolean、数组、联合类型、可选字段和类型推断。', link: '/typescript/basic-types', badge: '入门' },
  { title: '对象、接口与 type', description: '给业务对象、表单、接口响应和状态建模。', link: '/typescript/interface-type', badge: '核心' },
  { title: '泛型', description: '理解 ApiResult<T>、PageResult<T>、通用函数和组件类型。', link: '/typescript/generics', badge: '进阶' },
  { title: '类型收窄与类型守卫', description: '用 typeof、in、判别联合和自定义守卫安全处理 unknown 与外部数据。', link: '/typescript/narrowing-guards', badge: '安全' },
  { title: '工具类型与类型边界', description: '使用 Partial、Pick、Omit、Record 等工具表达表单、payload 和 ViewModel。', link: '/typescript/utility-types-boundary', badge: '边界' },
  { title: 'tsconfig 与工程配置', description: '理解 strict、paths、include、typecheck、project references 和构建检查。', link: '/typescript/tsconfig-engineering', badge: '工程' },
  { title: 'Vue 项目集成', description: '掌握 ref、props、emits、Pinia、API 请求和表单类型。', link: '/typescript/vue-integration', badge: '实战' },
  { title: '项目落地实践', description: '用 DTO、ViewModel、FormModel、Payload 和权限码类型串联真实项目类型边界。', link: '/typescript/project-practice', badge: '项目' },
  { title: '类型边界从零到项目', description: '用用户权限管理项目串联 DTO、ViewModel、FormModel、Payload、权限码、路由 meta、Pinia 和 typecheck。', link: '/typescript/type-boundary-project', badge: '项目' },
  { title: '常见问题', description: '处理 any、unknown、类型不一致、空值、第三方库类型等问题。', link: '/typescript/troubleshooting', badge: '排错' }
]" />

## TypeScript 在项目中解决什么

| 场景 | 没有类型的问题 | 有类型的收益 |
| --- | --- | --- |
| 接口请求 | 字段写错运行时才发现 | 编辑器提前提示 |
| 表单 | 提交字段和接口不一致 | 表单和 payload 边界清楚 |
| 组件 props | 父子组件传参混乱 | 调用组件时自动提示 |
| 权限码 | 字符串到处散落 | 可集中约束和复用 |
| 状态管理 | store 字段不清楚 | 状态结构明确 |
| 工程配置 | 本地和 CI 检查不一致 | typecheck、paths、strict 有统一规则 |

## 学习重点

第一阶段不要追求把 TypeScript 所有高级类型学完。优先掌握：

- 类型推断。
- interface。
- 联合类型。
- 泛型。
- `unknown`。
- 类型收窄。
- 常用工具类型。
- tsconfig 基础配置。
- Vue 组件中的类型写法。

高级条件类型和递归类型可以在项目需要时再补，不要一开始就追求类型体操。

## 最佳学习方式

不要单独刷语法。建议把 TypeScript 放到真实 Vue 项目里学：

1. 给接口响应加类型。
2. 给表单加类型。
3. 给 props 和 emits 加类型。
4. 给 Pinia store 加类型。
5. 处理构建时报出的类型错误。
6. 给 typecheck 接入构建或 CI。

## 下一步

第一次进入 TypeScript 模块，建议先看 [图解 TypeScript 核心概念](/typescript/visual-guide)，再从 [基础类型](/typescript/basic-types) 开始。学完 Vue 集成后，继续做 [类型边界从零到项目](/typescript/type-boundary-project)，把接口、表单、权限、路由和 Store 的类型关系串起来。
