# 基础类型

## 适合谁看

适合第一次系统学习 TypeScript，或者之前只会简单写 `string`、`number` 的学习者。

## 类型推断

很多时候 TypeScript 能自己推断类型：

```ts
const username = 'alice'
const page = 1
const enabled = true
```

不需要写成：

```ts
const username: string = 'alice'
```

只有当类型不够明确、可能为空、数组初始为空时，才需要手动标注。

## 基础类型

```ts
const username: string = 'alice'
const age: number = 18
const enabled: boolean = true
const roles: string[] = ['admin', 'editor']
```

数组另一种写法：

```ts
const ids: Array<number> = [1, 2, 3]
```

项目里常用 `number[]`，更短。

## 联合类型

联合类型表示“可能是这些值之一”：

```ts
type UserStatus = 'enabled' | 'disabled' | 'locked'

const status: UserStatus = 'enabled'
```

好处是避免写错：

```ts
const status: UserStatus = 'enable'
// 类型错误
```

## 可选字段

```ts
interface UserProfile {
  id: number
  nickname?: string
  avatarUrl?: string
}
```

`nickname?: string` 表示这个字段可能不存在。

读取时要处理：

```ts
const displayName = profile.nickname ?? '未命名用户'
```

## null 和 undefined

加载中的数据通常写成：

```ts
const currentUser = ref<User | null>(null)
```

这样可以明确表达：当前用户可能还没加载。

模板或逻辑里需要判断：

```ts
if (!currentUser.value) {
  return
}

console.log(currentUser.value.username)
```

## 字面量类型

权限动作可以这样定义：

```ts
type Action = 'create' | 'read' | 'update' | 'delete'
```

状态文案映射：

```ts
const statusLabel: Record<UserStatus, string> = {
  enabled: '启用',
  disabled: '停用',
  locked: '锁定'
}
```

如果遗漏某个状态，TypeScript 会提示。

## unknown 和 any

`any` 会关闭类型检查：

```ts
let value: any = {}
value.foo.bar.baz()
```

`unknown` 更安全：

```ts
function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    return error.message
  }

  return '未知错误'
}
```

不确定外部数据时，优先用 `unknown`，再通过判断缩小类型。

## 实际项目常见问题

### 1. 空数组推断成 never[]

```ts
const users = ref([])
```

后续 push 用户可能报错。推荐：

```ts
const users = ref<User[]>([])
```

### 2. 字符串状态到处写

**问题**

```ts
if (user.status === 'enable') {}
```

拼错了也不容易发现。

**解决方案**

定义联合类型：

```ts
type UserStatus = 'enabled' | 'disabled'
```

### 3. 接口数据可能为空

明确写出来：

```ts
const profile = ref<UserProfile | null>(null)
```

不要假设接口一定立即返回。

## 最佳实践

- 能推断就不手写。
- 空数组和空对象要补类型。
- 固定状态用联合类型。
- 不确定外部数据用 `unknown`。
- 可能为空就明确写 `null` 或可选字段。

## 下一步

继续学习 [对象、接口与 type](/typescript/interface-type)。
