# JavaScript 基础

## 适合谁看

适合已经能写基础页面，准备理解 Vue 响应式、组件逻辑、接口请求和业务数据处理的学习者。

Vue 的模板、事件、响应式、请求封装、权限判断，本质上都离不开 JavaScript。学 Vue 之前，不需要成为 JavaScript 专家，但必须掌握真实项目里最常用的语法和思维方式。

## 你会学到什么

- 变量、对象、数组怎么用于业务数据。
- 函数和模块如何组织代码。
- Promise 和 async/await 如何处理接口请求。
- 常见数据转换怎么写。
- 实际项目里的空值、异步、重复请求问题怎么处理。

## 变量和常量

优先使用 `const`，只有需要重新赋值时才使用 `let`。

```ts
const appName = 'Vue Admin'
let page = 1

page += 1
```

不要使用 `var`。它的作用域规则容易制造意外问题。

## 对象：描述一个业务实体

用户对象：

```ts
const user = {
  id: 1001,
  username: 'alice',
  enabled: true,
  roles: ['admin']
}
```

读取：

```ts
console.log(user.username)
```

复制并修改：

```ts
const updatedUser = {
  ...user,
  enabled: false
}
```

真实项目中，复制对象很常见，例如编辑表单不要直接修改表格行对象。

## 数组：处理列表数据

后台项目经常处理列表：

```ts
const users = [
  { id: 1, username: 'alice', enabled: true },
  { id: 2, username: 'bob', enabled: false }
]
```

筛选启用用户：

```ts
const enabledUsers = users.filter((user) => user.enabled)
```

提取用户名：

```ts
const usernames = users.map((user) => user.username)
```

查找某个用户：

```ts
const target = users.find((user) => user.id === 1)
```

判断是否有管理员：

```ts
const hasAdmin = users.some((user) => user.roles?.includes('admin'))
```

## 空值处理

接口数据不一定总是完整的。直接访问可能报错：

```ts
console.log(user.profile.name)
```

使用可选链：

```ts
console.log(user?.profile?.name)
```

提供默认值：

```ts
const nickname = user?.profile?.nickname ?? '未命名用户'
```

`??` 只在左侧是 `null` 或 `undefined` 时使用默认值，比 `||` 更适合处理数字 0 和空字符串。

## 函数：把动作命名

不要把复杂逻辑全写在事件里，应该提取成有名字的函数。

```ts
function normalizeKeyword(value: string) {
  return value.trim().toLowerCase()
}

function canDeleteUser(user: User) {
  return user.enabled === false && !user.roles.includes('admin')
}
```

函数名应该表达业务意图，而不是实现细节。

## 模块化

把不同职责放到不同文件：

```text
utils/format.ts
utils/permission.ts
api/user.ts
services/auth.ts
```

导出：

```ts
export function formatDate(value: string) {
  return value.slice(0, 10)
}
```

导入：

```ts
import { formatDate } from '@/utils/format'
```

## Promise 和 async/await

接口请求是异步的。推荐用 `async/await`：

```ts
async function fetchUsers() {
  loading.value = true

  try {
    const result = await userApi.getList()
    users.value = result.items
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}
```

`finally` 很适合恢复 loading 状态，不管请求成功还是失败都会执行。

## 实际项目常见问题

### 1. Cannot read properties of undefined

**原因**

数据还没加载完成，就访问了深层字段。

**解决方案**

```vue
<template>
  <div v-if="user">
    {{ user.profile?.nickname ?? user.username }}
  </div>
</template>
```

### 2. forEach 里 await 不按预期等待

**问题代码**

```ts
users.forEach(async (user) => {
  await updateUser(user)
})
```

`forEach` 不会等待内部 async。

**解决方案**

顺序执行：

```ts
for (const user of users) {
  await updateUser(user)
}
```

并发执行：

```ts
await Promise.all(users.map((user) => updateUser(user)))
```

### 3. 接口返回旧数据覆盖新数据

**原因**

多个请求并发，旧请求更晚返回。

**解决方案**

```ts
let requestId = 0

async function fetchList() {
  const currentId = ++requestId
  const result = await getList()

  if (currentId !== requestId) return

  list.value = result.items
}
```

### 4. 对象复制后修改仍然互相影响

**原因**

浅拷贝只复制第一层，嵌套对象仍然共享引用。

```ts
const copied = { ...user }
```

如果有嵌套表单，需要逐层复制关键字段，或者使用项目中确认过的深拷贝方案。

## 最佳实践

- 业务数据处理优先使用 `map`、`filter`、`find`、`some`。
- 异步请求必须考虑 loading、错误和 finally。
- 接口数据访问前考虑空值。
- 不要在页面里写大量无名复杂表达式。
- 复杂数据转换放到明确函数里，并写清楚输入输出。

## 下一步学习

继续学习 [数据类型与判断](/javascript/types)。
