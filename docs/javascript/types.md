# 数据类型与判断

## 适合谁看

适合经常遇到 `undefined`、`null`、空字符串、数字字符串、数组对象判断混乱的学习者。

JavaScript 的数据类型看似简单，但真实项目里的接口数据、表单输入、权限状态、缓存数据都可能出现类型不一致。掌握类型判断，可以减少大量运行时错误。

## 基础类型

| 类型 | 示例 |
| --- | --- |
| string | `'alice'` |
| number | `18` |
| boolean | `true` |
| undefined | 未赋值 |
| null | 明确为空 |
| object | `{ id: 1 }` |
| array | `[1, 2, 3]` |
| function | `() => {}` |

## typeof

```ts
typeof 'hello' // 'string'
typeof 123 // 'number'
typeof true // 'boolean'
typeof undefined // 'undefined'
typeof function () {} // 'function'
```

注意：

```ts
typeof null // 'object'
typeof [] // 'object'
```

所以判断数组不要用 `typeof`。

## 判断数组

```ts
Array.isArray(value)
```

示例：

```ts
function normalizeRoles(value: unknown) {
  return Array.isArray(value) ? value : []
}
```

## null 和 undefined

| 值 | 含义 |
| --- | --- |
| `undefined` | 没有赋值或字段不存在 |
| `null` | 明确表示空 |

接口设计里建议统一约定。前端处理时要有兜底：

```ts
const nickname = user.nickname ?? '未命名用户'
```

## 真值和假值

假值包括：

```text
false
0
''
null
undefined
NaN
```

注意：

```ts
const pageSize = value || 10
```

如果 `value` 是 `0`，也会变成 `10`。如果只想处理 `null` 和 `undefined`，用：

```ts
const pageSize = value ?? 10
```

## 字符串数字转换

表单输入经常拿到字符串：

```ts
const age = Number(form.age)
```

判断是否有效：

```ts
const age = Number(form.age)

if (Number.isNaN(age)) {
  throw new Error('年龄必须是数字')
}
```

## 实际项目常见问题

### 1. 后端返回 `null`，页面直接报错

**解决方案**

```ts
const list = Array.isArray(result.items) ? result.items : []
const name = result.user?.name ?? '未知用户'
```

### 2. 状态值既有数字又有字符串

**症状**

`status === 1` 有时不生效，因为接口返回的是 `'1'`。

**解决方案**

在 service 层统一转换：

```ts
function normalizeStatus(value: string | number) {
  return Number(value)
}
```

### 3. 空字符串被误判成默认值

**原因**

使用了 `||`。

**解决方案**

根据业务选择：

- 空字符串也算无值：使用 `||`。
- 只处理 `null` 和 `undefined`：使用 `??`。

## 最佳实践

- 判断数组用 `Array.isArray`。
- 默认值优先考虑 `??`。
- 表单提交前转换数字和布尔值。
- 外部接口数据进入页面前做规范化。
- TypeScript 类型不能替代运行时校验。

## 下一步学习

继续学习 [函数、作用域与闭包](/javascript/functions-scope)。
