# 图解 TypeScript 核心概念

## 这个页面解决什么

TypeScript 的难点不是“多写类型”，而是理解类型如何在编辑器、编译阶段和项目边界上帮助你发现问题。

## 适合谁看

适合已经会 JavaScript，准备在 Vue、React、Node 或后端工具项目中系统使用 TypeScript 的学习者。

## 一张图理解 TypeScript 在项目中的位置

```mermaid
flowchart LR
  A[".ts / .vue 源码"] --> B["TypeScript 类型检查"]
  B --> C["构建工具<br/>Vite / tsup / vue-tsc"]
  C --> D["JavaScript 输出"]
  D --> E["浏览器 / Node.js 运行"]
```

TypeScript 类型只在开发和构建阶段生效。运行时执行的仍然是 JavaScript。

所以：

- 类型能提前发现很多错误。
- 类型不能替代运行时校验。
- 接口返回数据、用户输入、localStorage 内容仍要校验。

## 一张图理解类型推导

```mermaid
flowchart TD
  A["const name = 'Ada'"] --> B["推导为字符串字面量或 string"]
  C["let age = 18"] --> D["推导为 number"]
  E["function add(a: number, b: number)"] --> F["返回值可推导为 number"]
  G["复杂对象"] --> H["建议显式声明接口"]
```

类型不是写得越多越好。简单变量交给推导，跨模块边界和复杂对象显式声明。

## 一张图理解 interface、type、class

```mermaid
flowchart TD
  A["interface"] --> A1["描述对象形状<br/>适合扩展和公共契约"]
  B["type"] --> B1["类型别名<br/>适合联合、交叉、工具类型"]
  C["class"] --> C1["运行时也存在<br/>有构造器和方法"]
  A1 --> D["API DTO / 组件 Props"]
  B1 --> E["状态枚举 / 组合类型"]
  C1 --> F["需要实例行为的模型"]
```

项目里最常见：

- API 数据结构用 `interface`。
- 联合状态用 `type`。
- 前端业务很少必须用 `class`。

## 一张图理解联合类型和类型收窄

```mermaid
flowchart TD
  A["value: string | number | null"] --> B{"if value === null"}
  B -- "是" --> C["value: null"]
  B -- "否" --> D["value: string | number"]
  D --> E{"typeof value === 'string'"}
  E -- "是" --> F["value: string"]
  E -- "否" --> G["value: number"]
```

类型收窄的目标是让代码每个分支里的类型更准确，减少强制断言。

## 一张图理解泛型

```mermaid
flowchart LR
  A["输入类型 T"] --> B["函数 / 组件 / 工具类型"]
  B --> C["输出仍保留 T 的信息"]
  D["User"] --> B
  B --> E["User[] / ApiResult<User>"]
  F["Order"] --> B
  B --> G["Order[] / ApiResult<Order>"]
```

泛型不是为了变复杂，而是为了保留输入和输出之间的类型关系。

典型场景：

```ts
interface ApiResult<T> {
  code: string
  data: T
  message: string
}
```

## 一张图理解 Vue 与 TypeScript

```mermaid
flowchart TD
  A["defineProps<Props>()"] --> B["模板 props 类型"]
  C["defineEmits<Emits>()"] --> D["事件名和参数类型"]
  E["ref<User | null>()"] --> F["状态类型"]
  G["computed(() => ...)"] --> H["派生类型"]
  I["API 返回 DTO"] --> E
```

Vue 项目里类型链路最好从 API DTO 开始：

```text
后端响应 DTO
↓
请求函数返回值
↓
页面 state
↓
组件 props
↓
模板展示
```

任何一层用了 `any`，后面的类型保护都会变弱。

## 一张图理解类型问题排查

```mermaid
flowchart TD
  A["TypeScript 报错"] --> B{"错误来源"}
  B --> C["对象字段不存在"]
  B --> D["可能为 undefined/null"]
  B --> E["联合类型未收窄"]
  B --> F["泛型推不出来"]
  B --> G["第三方库类型缺失"]

  C --> C1["检查接口定义和后端字段"]
  D --> D1["补默认值、判空、可选链"]
  E --> E1["使用 typeof、in、判别字段"]
  F --> F1["显式传泛型或调整函数签名"]
  G --> G1["安装类型包或补声明文件"]
```

## 下一步学习

继续学习 [基础类型](/typescript/basic-types)，或进入 [类型收窄与类型守卫](/typescript/narrowing-guards)。
