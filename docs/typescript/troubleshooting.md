# TypeScript 常见问题

## 1. 到处写 any

### 症状

类型报错消失了，但运行时错误变多，编辑器也不提示字段。

### 原因

`any` 会关闭类型检查。

### 解决方案

不确定外部数据时用 `unknown`：

```ts
function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    return error.message
  }

  return '未知错误'
}
```

## 2. 空数组推断错误

### 问题代码

```ts
const users = ref([])
```

### 解决方案

```ts
const users = ref<User[]>([])
```

## 3. 类型和后端返回不一致

### 症状

类型看起来正确，但页面显示 `undefined`。

### 原因

前端类型是假设出来的，后端实际返回不同字段。

### 解决方案

定义 Raw 类型并转换：

```ts
interface RawUser {
  user_name: string
  phone_no: string
}

function normalizeUser(raw: RawUser): User {
  return {
    id: 0,
    username: raw.user_name,
    mobile: raw.phone_no,
    enabled: true,
    roles: []
  }
}
```

## 4. 可选字段太多导致到处判断

### 原因

为了省事把所有字段都写成 `?`。

### 解决方案

只给真正可选的字段加 `?`。如果接口返回后一定有值，就不要写可选。

## 5. 第三方库没有类型

### 解决方案

优先查是否有官方类型或 `@types/xxx`。

如果没有，可以先写最小声明：

```ts
declare module 'legacy-lib' {
  export function init(options: Record<string, unknown>): void
}
```

不要为了一个库把整个项目降成 any。

## 6. 类型太复杂，没人敢改

### 症状

类型嵌套很多层，错误信息很长。

### 解决方案

- 拆成多个命名类型。
- 减少过度泛型。
- 优先让业务代码清晰。
- 复杂类型加注释说明意图。

## 7. Vue 模板里类型不准确

### 常见原因

- props 类型没写。
- ref 初始值没有泛型。
- 组件导入类型丢失。

### 解决方案

补齐 props、emits、ref 和 API 返回类型。

## 快速排查表

| 问题 | 优先检查 |
| --- | --- |
| 字段无提示 | 是否用了 any |
| 空数组报错 | 是否写了 `ref<User[]>([])` |
| 可能为空报错 | 是否需要 `User | null` |
| props 默认值报错 | 是否使用 `withDefaults` |
| emit 参数报错 | `defineEmits` 类型是否正确 |
| 后端字段不一致 | 是否需要 Raw 类型和 normalize |

如果问题已经进入真实项目层面，例如 DTO 泄漏到页面、表单和 Payload 混用、权限码散落、RouteMeta 没有类型、Store 被局部状态污染，继续看 [TypeScript 类型边界问题库](/projects/issues-typescript)。

## 最佳实践

- any 是最后手段，不是默认方案。
- 接口数据要以真实响应为准。
- 类型复杂时先拆命名类型。
- 业务类型要靠近业务模块。
- 类型错误不要简单绕过，优先理解它在提醒什么。

## 下一步学习

继续学习 [TypeScript 类型边界从零到项目](/typescript/type-boundary-project) 和 [TypeScript 类型边界问题库](/projects/issues-typescript)，把类型设计和真实排错连起来。
