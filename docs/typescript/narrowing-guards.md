# 类型收窄与类型守卫

## 适合谁看

适合已经会写联合类型，但遇到这些问题时还不够稳定的人：

- 接口返回 `unknown`，不知道怎么安全使用。
- `User | null` 到处报可能为空。
- 联合类型里不同字段不能直接访问。
- 表单、路由参数、LocalStorage 数据类型不可信。
- 想少写 `as`，但不知道怎么让 TypeScript 理解判断逻辑。

类型收窄的目标是：先承认数据可能有多种形态，再通过运行时判断把类型缩小到安全范围。

## 为什么需要收窄

TypeScript 只在编译期检查类型，但真实数据来自运行时：

- 接口响应。
- 表单输入。
- URL query。
- LocalStorage。
- 第三方 SDK。
- JSON.parse。

这些数据不一定可信。

```ts
const raw: unknown = JSON.parse(localStorage.getItem('user') || 'null')

console.log(raw.name)
```

这里不能直接访问 `name`，因为 `unknown` 需要先判断。

## 常见收窄方式

| 方式 | 适合 |
| --- | --- |
| `typeof` | 判断 string、number、boolean |
| `instanceof` | 判断 Error、Date、类实例 |
| `in` | 判断对象是否有某个字段 |
| 字面量判断 | 判断状态、类型枚举 |
| 自定义类型守卫 | 复用复杂判断 |
| 判空 | 处理 `null`、`undefined` |

## typeof

```ts
function formatValue(value: string | number) {
  if (typeof value === 'string') {
    return value.trim()
  }

  return value.toFixed(2)
}
```

进入 `if` 后，TypeScript 知道 `value` 是 string。

## 判空

```ts
type User = {
  id: number
  name: string
}

function getUsername(user: User | null) {
  if (!user) {
    return '未登录'
  }

  return user.name
}
```

不要为了省判断乱用非空断言：

```ts
user!.name
```

`!` 只是告诉 TypeScript “相信我”，不会减少运行时风险。

## in 判断对象字段

```ts
type PasswordLogin = {
  type: 'password'
  username: string
  password: string
}

type SmsLogin = {
  type: 'sms'
  mobile: string
  code: string
}

type LoginPayload = PasswordLogin | SmsLogin

function submitLogin(payload: LoginPayload) {
  if ('password' in payload) {
    return loginByPassword(payload.username, payload.password)
  }

  return loginBySms(payload.mobile, payload.code)
}
```

实际项目中，更推荐使用明确的判别字段。

## 判别联合

给每种形态加一个稳定字段：

```ts
type ApiState<T> =
  | { status: 'idle' }
  | { status: 'loading' }
  | { status: 'success'; data: T }
  | { status: 'error'; message: string }
```

使用：

```ts
function renderState(state: ApiState<User[]>) {
  switch (state.status) {
    case 'idle':
      return '等待加载'
    case 'loading':
      return '加载中'
    case 'success':
      return state.data.length
    case 'error':
      return state.message
  }
}
```

判别联合适合：

- 请求状态。
- 表单模式。
- 权限动作。
- 消息类型。
- 弹窗状态。

## 自定义类型守卫

自定义类型守卫让判断逻辑可复用。

```ts
type User = {
  id: number
  name: string
}

function isUser(value: unknown): value is User {
  if (typeof value !== 'object' || value === null) return false

  const record = value as Record<string, unknown>

  return typeof record.id === 'number' && typeof record.name === 'string'
}
```

使用：

```ts
const raw: unknown = await fetchUser()

if (!isUser(raw)) {
  throw new Error('用户数据格式不正确')
}

console.log(raw.name)
```

进入 `if` 后，TypeScript 知道 `raw` 是 User。

## 断言函数

有些场景希望不满足条件时直接抛错：

```ts
function assertUser(value: unknown): asserts value is User {
  if (!isUser(value)) {
    throw new Error('用户数据格式不正确')
  }
}
```

使用：

```ts
const raw: unknown = await fetchUser()

assertUser(raw)

console.log(raw.name)
```

断言函数适合在 service 层或数据入口做硬校验。

## 穷尽检查

当联合类型新增一种状态时，希望编译器提醒你更新处理逻辑。

```ts
function assertNever(value: never): never {
  throw new Error(`未处理的状态：${JSON.stringify(value)}`)
}

function getStateText(state: ApiState<User[]>) {
  switch (state.status) {
    case 'idle':
      return '等待'
    case 'loading':
      return '加载中'
    case 'success':
      return '成功'
    case 'error':
      return state.message
    default:
      return assertNever(state)
  }
}
```

如果之后新增 `status: 'empty'`，没有处理时会出现类型错误。

## 实际项目常见问题

### 1. 接口数据直接写成业务类型

**问题**

```ts
const user = await request.get<User>('/user')
```

这只是告诉 TypeScript “它是 User”，并没有验证后端真实返回。

**建议**

关键数据入口使用 Raw 类型、normalize 或类型守卫：

```ts
type RawUser = {
  user_id: number
  user_name: string
}

function normalizeUser(raw: RawUser): User {
  return {
    id: raw.user_id,
    name: raw.user_name
  }
}
```

### 2. 为了消除报错到处 as

`as` 可以用于边界转换，但不能变成绕过类型系统的默认方案。

优先顺序：

1. 运行时判断。
2. 类型守卫。
3. normalize。
4. 最后才是局部 `as`。

### 3. 可选字段导致模板到处判断

如果数据进入页面前可以补默认值，就不要让页面到处处理 `undefined`。

```ts
function normalizeUser(raw: RawUser): User {
  return {
    id: raw.id,
    name: raw.name || '未命名用户',
    roles: raw.roles ?? []
  }
}
```

### 4. 权限码只是 string

不推荐：

```ts
function can(action: string) {}
```

推荐：

```ts
type PermissionAction = 'user:create' | 'user:update' | 'user:delete'

function can(action: PermissionAction) {}
```

这样写错权限码时编辑器会提示。

## 最佳实践

- 外部输入先用 `unknown` 或 Raw 类型承接。
- 通过类型守卫、normalize 或断言函数进入业务类型。
- 联合类型优先设计判别字段。
- 不要滥用 `as` 和非空断言。
- 关键状态用穷尽检查防止漏处理。
- 类型收窄逻辑放在数据入口或 service，不要散落在模板里。

## 学习检查

学完本节后，你应该能回答：

- `unknown` 和 `any` 的区别是什么。
- `typeof`、`in`、判别联合分别适合什么。
- 自定义类型守卫为什么返回 `value is User`。
- 为什么接口响应类型不等于运行时校验。
- `never` 怎么帮助检查漏处理状态。

## 参考资料

- [TypeScript Handbook: Narrowing](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)

## 下一步学习

继续学习 [工具类型与类型边界](/typescript/utility-types-boundary)，把常见类型转换控制在可维护范围内。
