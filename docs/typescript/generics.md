# 泛型

## 适合谁看

适合已经理解基础类型和 interface，但看到 `<T>`、`ApiResult<T>`、`PageResult<T>` 时不太确定含义的学习者。

泛型可以理解为：**先留一个类型位置，使用时再填进去**。

## 最小例子

```ts
function identity<T>(value: T): T {
  return value
}

const name = identity<string>('alice')
const age = identity<number>(18)
```

实际项目里通常不需要写这种函数，但它能帮助理解泛型。

## 接口响应泛型

后端统一响应：

```ts
interface ApiResult<T> {
  code: number
  message: string
  data: T
}
```

当接口返回用户：

```ts
type UserResult = ApiResult<User>
```

当接口返回角色：

```ts
type RoleResult = ApiResult<Role>
```

`ApiResult<T>` 不关心具体 data 是什么，使用时填进去。

## 分页泛型

```ts
interface PageResult<T> {
  items: T[]
  total: number
}
```

用户分页：

```ts
type UserPageResult = PageResult<User>
```

角色分页：

```ts
type RolePageResult = PageResult<Role>
```

请求函数：

```ts
function getUserList(params: UserQuery) {
  return request.get<PageResult<User>>('/users', { params })
}
```

## 泛型约束

有时泛型需要满足某些条件：

```ts
function getById<T extends { id: number }>(items: T[], id: number) {
  return items.find((item) => item.id === id)
}
```

这里要求 `T` 必须有 `id: number`。

## Composable 中的泛型

编辑抽屉：

```ts
export function useEditDrawer<T>() {
  const visible = ref(false)
  const editingRecord = ref<T | null>(null)

  function openEdit(record: T) {
    editingRecord.value = record
    visible.value = true
  }

  return {
    visible,
    editingRecord,
    openEdit
  }
}
```

用户页面：

```ts
const { visible, editingRecord, openEdit } = useEditDrawer<User>()
```

角色页面：

```ts
const drawer = useEditDrawer<Role>()
```

## 实际项目常见问题

### 1. 泛型写得太复杂

**症状**

类型看起来很强大，但团队没人敢改。

**建议**

泛型优先用于：

- API 响应。
- 分页结果。
- 通用 composable。
- 通用组件 props。

不要为了炫技写复杂类型体操。

### 2. 泛型没有约束导致访问字段报错

```ts
function getId<T>(item: T) {
  return item.id
}
```

TypeScript 不知道 T 一定有 id。

解决：

```ts
function getId<T extends { id: number }>(item: T) {
  return item.id
}
```

### 3. 泛型推断失败

可以显式传入：

```ts
const users = ref<User[]>([])
const drawer = useEditDrawer<User>()
```

## 最佳实践

- 先理解 `ApiResult<T>` 和 `PageResult<T>`。
- 通用函数访问字段时加泛型约束。
- 泛型服务复用，不服务炫技。
- 类型过复杂时优先拆小，而不是继续嵌套。

## 下一步

继续学习 [类型收窄与类型守卫](/typescript/narrowing-guards)，理解如何安全处理外部数据和联合类型。
