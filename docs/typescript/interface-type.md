# 对象、接口与 type

## 适合谁看

适合准备给接口数据、表单、权限、路由和组件 props 建模的学习者。

## interface：描述对象结构

```ts
interface User {
  id: number
  username: string
  mobile: string
  enabled: boolean
  roles: string[]
}
```

使用：

```ts
const user: User = {
  id: 1001,
  username: 'alice',
  mobile: '13800000000',
  enabled: true,
  roles: ['admin']
}
```

## 表单类型和实体类型分开

后端用户对象：

```ts
interface User {
  id: number
  username: string
  mobile: string
  enabled: boolean
  createdAt: string
}
```

表单对象：

```ts
interface UserForm {
  id?: number
  username: string
  mobile: string
  enabled: boolean
}
```

不要把完整实体直接当表单，因为有些字段不应该由用户编辑。

## type：联合、映射和组合

联合类型：

```ts
type UserStatus = 'enabled' | 'disabled' | 'locked'
```

组合类型：

```ts
type UserWithPermissions = User & {
  permissions: string[]
}
```

## interface 和 type 怎么选

| 场景 | 推荐 |
| --- | --- |
| 描述对象结构 | `interface` |
| 描述联合类型 | `type` |
| 组合工具类型 | `type` |
| 公共业务实体 | `interface` |

不用过度纠结，团队统一即可。

## 接口响应建模

```ts
interface ApiResult<T> {
  code: number
  message: string
  data: T
}

interface PageResult<T> {
  items: T[]
  total: number
}
```

用户列表：

```ts
type UserListResult = PageResult<User>
```

## 权限码建模

```ts
export const PermissionCode = {
  UserCreate: 'system:user:create',
  UserUpdate: 'system:user:update',
  UserDelete: 'system:user:delete'
} as const

export type PermissionCode =
  typeof PermissionCode[keyof typeof PermissionCode]
```

这样权限码既集中，又能获得类型提示。

## 实际项目常见问题

### 1. 一个类型文件越来越大

**解决方案**

按业务拆分：

```text
types/user.ts
types/role.ts
types/permission.ts
```

不要把所有类型都放到 `types/index.ts`。

### 2. 接口字段变化导致页面 undefined

**解决方案**

在 service 层做适配：

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

### 3. 为了省事把所有字段设为可选

**问题**

类型失去约束，哪里都要判断空值。

**建议**

只有确实可能不存在的字段才写 `?`。

## 最佳实践

- 实体、表单、接口响应分开建模。
- 权限码、状态值集中定义。
- 类型按业务模块拆分。
- 不要滥用可选字段。
- 后端原始结构和前端页面结构不一致时，在 service 层转换。

## 下一步

继续学习 [泛型](/typescript/generics)。
