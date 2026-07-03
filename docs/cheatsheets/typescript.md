# TypeScript 速查

## 基础类型

```ts
const name: string = 'Tom'
const age: number = 18
const enabled: boolean = true
const tags: string[] = ['admin', 'editor']
```

对象类型：

```ts
type User = {
  id: number
  username: string
  mobile?: string
  enabled: boolean
}
```

可选字段：

```ts
type Form = {
  keyword?: string
}
```

联合类型：

```ts
type Status = 'pending' | 'success' | 'failed'
```

## interface 和 type

| 场景 | 推荐 |
| --- | --- |
| 业务对象结构 | `type` 或 `interface` 都可以 |
| 联合类型 | `type` |
| 函数类型 | `type` |
| 需要扩展第三方声明 | `interface` |

常用写法：

```ts
type UserListQuery = {
  page: number
  pageSize: number
  keyword?: string
}

type UserListResult = {
  items: User[]
  total: number
}
```

## 函数类型

```ts
function getUser(id: number): Promise<User> {
  return api.get(`/users/${id}`)
}
```

回调类型：

```ts
type SubmitHandler = (payload: UserForm) => Promise<void>
```

组件事件载荷：

```ts
const emit = defineEmits<{
  submit: [payload: UserForm]
  cancel: []
}>()
```

## 泛型

接口响应：

```ts
type ApiResult<T> = {
  code: string
  message: string
  data: T
}
```

分页结果：

```ts
type PageResult<T> = {
  items: T[]
  total: number
}
```

请求函数：

```ts
function getPage<T>(url: string, params: object): Promise<PageResult<T>> {
  return request.get(url, { params })
}
```

使用：

```ts
const result = await getPage<User>('/users', query)
```

## 常用工具类型

| 工具类型 | 用途 |
| --- | --- |
| `Partial<T>` | 所有字段变可选 |
| `Required<T>` | 所有字段变必填 |
| `Pick<T, K>` | 选部分字段 |
| `Omit<T, K>` | 排除部分字段 |
| `Record<K, V>` | key-value 对象 |
| `ReturnType<T>` | 获取函数返回类型 |

示例：

```ts
type UserForm = Pick<User, 'username' | 'mobile' | 'enabled'>
type UpdateUserPayload = Partial<UserForm> & { id: number }
type StatusTextMap = Record<Status, string>
```

## Vue 中的常用类型

ref：

```ts
const count = ref<number>(0)
const user = ref<User | null>(null)
```

表单：

```ts
const form = reactive<UserForm>({
  username: '',
  mobile: '',
  enabled: true
})
```

props：

```ts
const props = withDefaults(defineProps<{
  title: string
  loading?: boolean
}>(), {
  loading: false
})
```

模板 ref：

```ts
const inputRef = ref<HTMLInputElement | null>(null)
```

## 常见坑

| 问题 | 正确处理 |
| --- | --- |
| 到处写 `any` | 用 `unknown` 或补类型 |
| 接口返回类型不明确 | 为 API 写 response 类型 |
| 表单和接口 payload 混用 | 分开 `Form` 和 `Payload` |
| 可空值直接读取 | 先判断或用 `?.` |
| 类型和运行时校验混淆 | TypeScript 不会校验线上数据 |

## 项目建议

- API 层定义请求参数和响应类型。
- 页面表单类型不要直接复用数据库实体。
- 权限码、状态码用联合类型或 `as const`。
- 对后端不稳定字段做适配，不要让页面到处兜底。
- 不用 `any` 逃避问题，除非临时迁移并标注原因。

## 下一步学习

- [TypeScript 学习导览](/typescript/introduction)
- [基础类型](/typescript/basic-types)
- [泛型](/typescript/generics)
- [Vue 项目集成](/typescript/vue-integration)
