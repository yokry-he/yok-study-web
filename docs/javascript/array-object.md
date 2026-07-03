# 数组与对象处理

## 适合谁看

适合经常处理列表、表格、树结构、下拉选项和接口数据转换的学习者。

真实前端项目里，大量业务逻辑其实是数组和对象处理。把这部分写清楚，Vue 页面会简单很多。

## 常用数组方法

| 方法 | 用途 |
| --- | --- |
| `map` | 转换每一项 |
| `filter` | 筛选 |
| `find` | 找一个 |
| `some` | 是否存在 |
| `every` | 是否全部满足 |
| `reduce` | 汇总 |

## map：转换结构

```ts
const options = users.map((user) => ({
  label: user.username,
  value: user.id
}))
```

适合把接口数据转成选择器选项。

## filter：筛选数据

```ts
const enabledUsers = users.filter((user) => user.enabled)
```

组合筛选：

```ts
function filterUsers(users: User[], query: UserQuery) {
  return users.filter((user) => {
    const matchKeyword = !query.keyword || user.username.includes(query.keyword)
    const matchStatus = query.enabled == null || user.enabled === query.enabled
    return matchKeyword && matchStatus
  })
}
```

## find 和 some

```ts
const currentUser = users.find((user) => user.id === currentUserId)
const hasAdmin = users.some((user) => user.roles.includes('admin'))
```

权限判断中很常见：

```ts
function can(permissions: string[], code: string) {
  return permissions.includes(code)
}
```

## reduce：汇总

```ts
const total = orders.reduce((sum, order) => {
  return sum + order.amount
}, 0)
```

分组：

```ts
const usersByDepartment = users.reduce<Record<string, User[]>>((map, user) => {
  const key = user.departmentName || '未分组'
  map[key] ||= []
  map[key].push(user)
  return map
}, {})
```

## 对象复制

浅拷贝：

```ts
const copied = { ...user }
```

只复制第一层。嵌套对象仍然共享引用。

表单建议明确复制字段：

```ts
function createUserForm(user: User): UserForm {
  return {
    id: user.id,
    username: user.username,
    mobile: user.mobile,
    enabled: user.enabled,
    roleIds: user.roles.map((role) => role.id)
  }
}
```

## 实际项目常见问题

### 1. 直接修改 props 传入对象

**解决方案**

复制成表单对象，再编辑。

### 2. map 中忘记 return

```ts
const options = users.map((user) => {
  label: user.username
  value: user.id
})
```

上面写法会出错。对象字面量需要加括号：

```ts
const options = users.map((user) => ({
  label: user.username,
  value: user.id
}))
```

### 3. filter 条件越来越复杂

**解决方案**

拆成命名判断函数：

```ts
function matchKeyword(user: User, keyword: string) {
  return !keyword || user.username.includes(keyword)
}
```

### 4. 树结构处理混乱

菜单、部门、权限树都属于树结构。建议先定义清楚类型，再写递归函数。

```ts
interface TreeNode {
  id: number
  label: string
  children?: TreeNode[]
}
```

## 最佳实践

- 数据转换放到函数里，不堆在模板中。
- 表单对象明确复制字段。
- 复杂筛选拆成多个命名条件。
- 树结构先定义类型再递归。
- 不要为了省事直接修改外部传入对象。

## 下一步学习

继续学习 [DOM 事件](/javascript/dom-events)，理解页面交互、事件委托和组件卸载清理。
