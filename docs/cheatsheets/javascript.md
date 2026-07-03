# JavaScript 速查

## 变量和基础判断

```ts
const name = 'Tom'
let count = 0

const hasName = Boolean(name)
const isEmpty = value == null || value === ''
```

| 写法 | 用途 |
| --- | --- |
| `const` | 默认使用，引用不重新赋值 |
| `let` | 需要重新赋值 |
| `typeof value` | 判断基础类型 |
| `Array.isArray(value)` | 判断数组 |
| `value == null` | 同时判断 `null` 和 `undefined` |

## 常用数组方法

| 方法 | 用途 | 示例 |
| --- | --- | --- |
| `map` | 转换每一项 | `list.map(item => item.id)` |
| `filter` | 过滤 | `list.filter(item => item.enabled)` |
| `find` | 找第一项 | `list.find(item => item.id === id)` |
| `some` | 是否存在 | `list.some(item => item.checked)` |
| `every` | 是否全部满足 | `list.every(item => item.enabled)` |
| `reduce` | 汇总 | `list.reduce((sum, item) => sum + item.amount, 0)` |

常见列表转换：

```ts
const options = users.map((user) => ({
  label: user.username,
  value: user.id
}))
```

去掉无效项：

```ts
const validIds = ids.filter((id) => id != null)
```

## 常用对象写法

浅拷贝：

```ts
const nextForm = { ...form }
```

合并默认值：

```ts
const query = {
  page: 1,
  pageSize: 20,
  ...params
}
```

动态字段：

```ts
const field = 'username'
const payload = {
  [field]: 'Tom'
}
```

安全读取：

```ts
const city = user.profile?.address?.city ?? '未知'
```

## 异步写法

基础写法：

```ts
async function fetchUser(id: number) {
  const user = await api.getUser(id)
  return user
}
```

并发请求：

```ts
const [profile, permissions] = await Promise.all([
  api.getProfile(),
  api.getPermissions()
])
```

失败处理：

```ts
try {
  await api.submit(form)
} catch (error) {
  showError(error)
} finally {
  loading.value = false
}
```

## 模块导入导出

命名导出：

```ts
export function formatDate() {}
export const statusMap = {}
```

命名导入：

```ts
import { formatDate, statusMap } from './format'
```

默认导出：

```ts
export default function request() {}
```

默认导入：

```ts
import request from './request'
```

项目中推荐业务工具使用命名导出，重构和搜索更清晰。

## 常见坑

| 问题 | 正确处理 |
| --- | --- |
| `0` 被当成空值 | 用 `value == null` 判断缺失 |
| 数组直接改原对象 | 表单编辑时复制对象 |
| `forEach` 里 `await` 不等待 | 用 `for...of` 或 `Promise.all` |
| `JSON.parse` 报错 | 用 `try/catch` 包住 |
| 浮点数精度问题 | 金额用分或专门 decimal 方案 |

## 项目建议

- 页面里少写复杂数据转换，抽到 service 或 utils。
- 请求并发结果要防止旧请求覆盖新请求。
- 表单编辑不要直接绑定列表行对象。
- 复杂条件判断抽成具名函数。
- 不要用中文文案做业务逻辑判断，使用稳定 code。

## 下一步学习

- [JavaScript 学习导览](/javascript/introduction)
- [数组与对象处理](/javascript/array-object)
- [异步编程](/javascript/async)
- [前端页面与状态问题](/projects/issues-frontend)
