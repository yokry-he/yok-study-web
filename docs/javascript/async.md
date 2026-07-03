# 异步编程

## 适合谁看

适合经常写接口请求，但对 Promise、async/await、并发请求、错误处理和 loading 状态还不熟悉的学习者。

Vue 项目里，登录、列表、详情、提交、上传、权限初始化都依赖异步编程。异步写不好，页面就会出现重复请求、旧数据覆盖新数据、loading 不消失等问题。

## Promise 是什么

Promise 表示一个未来才会完成的结果：

```ts
const promise = fetch('/api/users')
```

它可能成功，也可能失败。

## async/await

推荐写法：

```ts
async function fetchUsers() {
  const result = await userApi.getList()
  users.value = result.items
}
```

加错误处理：

```ts
async function fetchUsers() {
  loading.value = true

  try {
    const result = await userApi.getList()
    users.value = result.items
  } catch (error) {
    showError(getErrorMessage(error))
  } finally {
    loading.value = false
  }
}
```

## 顺序执行和并发执行

顺序执行：

```ts
await fetchUser()
await fetchPermissions()
await fetchMenus()
```

后一个依赖前一个时使用。

并发执行：

```ts
const [profile, permissions] = await Promise.all([
  fetchProfile(),
  fetchPermissions()
])
```

互不依赖时使用，可以减少等待时间。

## 防重复提交

```ts
const submitting = ref(false)

async function submit() {
  if (submitting.value) return

  submitting.value = true
  try {
    await save()
  } finally {
    submitting.value = false
  }
}
```

## 处理旧请求覆盖新请求

```ts
let requestId = 0

async function fetchList() {
  const currentId = ++requestId
  const result = await getList(query.value)

  if (currentId !== requestId) return

  list.value = result.items
}
```

## 实际项目常见问题

### 1. loading 一直不消失

**原因**

请求失败时没有恢复 loading。

**解决方案**

使用 `finally`。

### 2. forEach 中 await 不等待

**解决方案**

顺序：

```ts
for (const item of list) {
  await saveItem(item)
}
```

并发：

```ts
await Promise.all(list.map((item) => saveItem(item)))
```

### 3. 页面离开后请求回来仍然改状态

**解决方案**

记录组件是否仍然有效，或使用请求取消机制。简单场景可以用请求序号避免旧结果覆盖。

### 4. 多个接口部分失败

`Promise.all` 有一个失败就整体失败。如果需要分别处理，可以使用 `Promise.allSettled`。

```ts
const results = await Promise.allSettled([
  fetchProfile(),
  fetchPermissions()
])
```

## 最佳实践

- 所有请求都考虑 loading、error、finally。
- 不依赖顺序的请求用并发。
- 快速变化的请求考虑防抖、取消或请求序号。
- 提交操作防重复。
- 异步错误不要吞掉，至少记录或提示。

## 下一步学习

继续学习 [事件循环](/javascript/event-loop)，理解 Promise、定时器、渲染和页面卡顿背后的执行顺序。
