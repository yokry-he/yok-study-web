# 工具类型与类型边界

## 适合谁看

适合已经理解 interface、type 和泛型，但看到 `Partial`、`Pick`、`Omit`、`Record` 时不知道怎么用，或者项目里类型越来越复杂的人。

工具类型能减少重复类型定义，但也容易被滥用成难以维护的类型体操。本节重点是项目里最常用的类型转换和边界。

## 先记住原则

工具类型服务业务表达，不服务炫技。

如果一个类型需要反复读很久才能看懂，通常应该拆成命名类型，或者回到更简单的数据结构。

## 常用工具类型

| 工具类型 | 作用 |
| --- | --- |
| `Partial<T>` | 把所有字段变成可选 |
| `Required<T>` | 把所有字段变成必填 |
| `Readonly<T>` | 把所有字段变成只读 |
| `Pick<T, K>` | 从类型里选择部分字段 |
| `Omit<T, K>` | 从类型里排除部分字段 |
| `Record<K, T>` | 创建键值映射类型 |
| `ReturnType<T>` | 获取函数返回值类型 |
| `Parameters<T>` | 获取函数参数类型 |
| `Awaited<T>` | 获取 Promise 解包后的类型 |

## Partial：更新表单

用户详情：

```ts
type User = {
  id: number
  username: string
  mobile: string
  enabled: boolean
}
```

编辑时只提交变化字段：

```ts
type UpdateUserPayload = Partial<Pick<User, 'username' | 'mobile' | 'enabled'>>
```

使用：

```ts
async function updateUser(id: number, payload: UpdateUserPayload) {
  return request.patch(`/users/${id}`, payload)
}
```

注意：`Partial<User>` 不一定合适，因为 `id` 通常不应该被更新。

## Pick 和 Omit：拆分边界

创建用户时没有 id：

```ts
type CreateUserPayload = Omit<User, 'id'>
```

列表只需要部分字段：

```ts
type UserListItem = Pick<User, 'id' | 'username' | 'enabled'>
```

建议：

- Payload 类型靠近接口模块。
- ViewModel 类型靠近页面模块。
- 不要所有场景都直接复用完整 `User`。

## Record：字典映射

权限文案：

```ts
type PermissionAction = 'create' | 'update' | 'delete'

const permissionLabels: Record<PermissionAction, string> = {
  create: '新增',
  update: '编辑',
  delete: '删除'
}
```

如果漏写 `delete`，TypeScript 会提示。

适合：

- 枚举到文案。
- 状态到颜色。
- 权限到配置。
- tab key 到组件。

## ReturnType 和 Parameters

从函数推导类型，减少重复声明。

```ts
function createUserQuery() {
  return {
    keyword: '',
    page: 1,
    pageSize: 20
  }
}

type UserQuery = ReturnType<typeof createUserQuery>
```

适合默认值函数和查询参数。

```ts
function openUserDialog(id: number, mode: 'view' | 'edit') {}

type OpenUserDialogArgs = Parameters<typeof openUserDialog>
```

这类工具要谨慎使用。对外公共类型最好显式命名，避免函数改动连带影响太大。

## Awaited：异步结果

```ts
async function fetchUsers() {
  return request.get<PageResult<User>>('/users')
}

type FetchUsersResult = Awaited<ReturnType<typeof fetchUsers>>
```

适合从 API 函数推导异步结果。

如果团队觉得这种组合太绕，也可以显式写：

```ts
type FetchUsersResult = PageResult<User>
```

可读性优先。

## 类型边界怎么划分

真实项目里常见几类类型：

| 类型 | 放在哪里 | 说明 |
| --- | --- | --- |
| API Raw 类型 | `api` 或 `model` | 后端原始返回 |
| Domain 类型 | `model` | 前端业务稳定结构 |
| Form 类型 | 页面或 feature | 表单状态 |
| Payload 类型 | `api` | 提交给接口的数据 |
| ViewModel 类型 | 页面或组件 | 页面展示所需结构 |

不要一个 `User` 类型走天下。

示例：

```ts
type RawUser = {
  user_id: number
  user_name: string
}

type User = {
  id: number
  name: string
}

type UserForm = {
  name: string
}

type CreateUserPayload = UserForm
```

这能让接口变化、页面变化和业务模型变化互相隔离。

## 类型体操边界

高级类型不是不能用，而是要有边界。

可以用：

- `Pick` / `Omit` 拆接口字段。
- `Record` 约束配置完整性。
- `ReturnType` 推导局部工厂函数结果。
- 简单条件类型封装通用库能力。

谨慎用：

- 多层递归类型。
- 过深条件类型。
- 很长的模板字面量类型。
- 为了少写几个字段制造复杂工具。

判断标准：

```text
这个类型能否被普通业务同事快速读懂？
错误信息是否还能定位？
改一个字段会不会牵动太多地方？
```

## 实际项目常见问题

### 1. Partial 滥用导致必填字段丢失

不推荐：

```ts
function submit(payload: Partial<User>) {}
```

这会让所有字段都可选，可能导致接口缺必填字段。

更好：

```ts
type UpdateUserPayload = Partial<Pick<User, 'username' | 'mobile'>>
```

### 2. Omit 套太多层看不懂

```ts
type A = Omit<Pick<User, 'id' | 'name' | 'roles'>, 'roles'>
```

这种类型不如直接写清楚：

```ts
type UserOption = {
  id: number
  name: string
}
```

### 3. Record key 写成 string

```ts
const labels: Record<string, string> = {}
```

这样无法约束必须有哪些 key。

更好：

```ts
type Status = 'enabled' | 'disabled'

const labels: Record<Status, string> = {
  enabled: '启用',
  disabled: '停用'
}
```

### 4. 从 API 函数推导类型导致耦合

`ReturnType<typeof api.xxx>` 很方便，但如果页面类型完全依赖 API 函数签名，接口层调整会影响页面。

公共业务模型建议显式命名。

## 最佳实践

- 工具类型优先用于减少重复，不用于炫技。
- `Partial` 只用于明确的更新场景。
- `Pick` / `Omit` 用于清晰的边界转换。
- `Record` 配合联合类型约束配置完整性。
- 类型太复杂时直接写命名类型。
- 区分 Raw、Domain、Form、Payload、ViewModel。
- 公共类型要稳定，局部推导要克制。

## 学习检查

学完本节后，你应该能回答：

- `Partial` 为什么不适合随便套完整业务类型。
- `Pick` 和 `Omit` 分别适合什么场景。
- `Record<'a' | 'b', T>` 比 `Record<string, T>` 多了什么约束。
- Raw、Domain、Payload、Form 为什么要区分。
- 什么样的高级类型应该停下来拆小。

## 参考资料

- [TypeScript Handbook: Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html)
- [TypeScript Handbook: Creating Types from Types](https://www.typescriptlang.org/docs/handbook/2/types-from-types.html)

## 下一步学习

继续学习 [tsconfig 与工程配置](/typescript/tsconfig-engineering)，把类型能力接入构建、编辑器和项目边界。
