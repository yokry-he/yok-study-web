# 组合式 API

## 适合谁看

适合已经会写组件，但发现多个页面里有重复逻辑，例如分页、加载状态、弹窗状态、表单提交、权限判断的学习者。

组合式 API 的重点不是“新语法”，而是让你把相关逻辑放在一起，并把可复用逻辑抽成 composable。

## 你会学到什么

- `<script setup>` 的基本组织方式。
- composable 是什么，什么时候该抽。
- 如何封装加载状态、分页、弹窗和权限判断。
- 组合式 API 常见滥用问题怎么避免。

## 为什么需要组合式 API

假设一个用户列表页里有这些逻辑：

- 查询条件。
- 分页。
- 请求列表。
- 加载状态。
- 删除确认。
- 新增/编辑弹窗。

如果全部写在一个组件里，文件很快会变长。组合式 API 允许我们把“同一类逻辑”收在一起：

```text
useUserList()       用户列表和分页
useUserDrawer()     新增/编辑弹窗
usePermission()     权限判断
```

组件只负责组装这些能力。

## `<script setup>` 基础

```vue
<script setup lang="ts">
import { computed, ref } from 'vue'

const firstName = ref('Ada')
const lastName = ref('Lovelace')

const fullName = computed(() => `${firstName.value} ${lastName.value}`)

function updateName() {
  firstName.value = 'Grace'
  lastName.value = 'Hopper'
}
</script>

<template>
  <p>{{ fullName }}</p>
  <button type="button" @click="updateName">更新姓名</button>
</template>
```

在 `<script setup>` 中定义的变量和函数，可以直接在模板中使用。

## composable 是什么

Composable 是一个以 `use` 开头的函数，用来封装可复用状态和逻辑。

最简单的例子：

```ts
import { ref } from 'vue'

export function useLoading() {
  const loading = ref(false)

  async function run(task: () => Promise<void>) {
    loading.value = true
    try {
      await task()
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    run
  }
}
```

使用：

```ts
const { loading, run } = useLoading()

async function submit() {
  await run(async () => {
    await saveUser(form.value)
  })
}
```

这样每个提交按钮都不用重复写 `loading.value = true` 和 `finally`。

## 封装分页逻辑

真实后台项目里，列表分页非常常见：

```ts
import { reactive } from 'vue'

export function usePagination(defaultPageSize = 10) {
  const pagination = reactive({
    page: 1,
    pageSize: defaultPageSize,
    total: 0
  })

  function resetPage() {
    pagination.page = 1
  }

  function setTotal(total: number) {
    pagination.total = total
  }

  return {
    pagination,
    resetPage,
    setTotal
  }
}
```

页面使用：

```ts
const { pagination, resetPage, setTotal } = usePagination(20)

async function fetchList() {
  const result = await userApi.getList({
    page: pagination.page,
    pageSize: pagination.pageSize,
    keyword: keyword.value
  })

  users.value = result.items
  setTotal(result.total)
}

async function search() {
  resetPage()
  await fetchList()
}
```

## 封装弹窗状态

新增和编辑弹窗通常需要这些状态：

- 是否显示。
- 当前编辑对象。
- 打开新增。
- 打开编辑。
- 关闭并重置。

```ts
import { ref } from 'vue'

export function useEditDrawer<T>() {
  const visible = ref(false)
  const editingRecord = ref<T | null>(null)

  function openCreate() {
    editingRecord.value = null
    visible.value = true
  }

  function openEdit(record: T) {
    editingRecord.value = record
    visible.value = true
  }

  function close() {
    visible.value = false
    editingRecord.value = null
  }

  return {
    visible,
    editingRecord,
    openCreate,
    openEdit,
    close
  }
}
```

这个 composable 不知道“用户”或“角色”是什么，因此可以复用到多个页面。

## 什么逻辑适合抽成 composable

| 逻辑 | 是否适合 | 原因 |
| --- | --- | --- |
| 分页状态 | 适合 | 多个列表页重复出现 |
| 加载状态 | 适合 | 多个按钮和请求都需要 |
| 权限判断 | 适合 | 多个页面和按钮都需要 |
| 用户表单字段 | 不一定 | 如果只在一个页面使用，先留在组件内 |
| 页面专属业务流程 | 不适合 | 抽出去反而更难读 |

## 实际项目常见问题

### 1. composable 变成万能工具箱

**症状**

一个 `useUserPage()` 里包含请求、表格、弹窗、权限、路由跳转、消息提示，代码比原组件还长。

**原因**

只是把代码从组件搬到函数里，没有重新划分职责。

**解决方案**

按能力拆分：

```text
useUserList()
useUserFormDrawer()
useUserPermission()
```

页面再组合它们。

### 2. composable 里隐藏太多副作用

**症状**

调用一个函数后，路由跳了、消息弹了、store 也改了，调用者很难预期后果。

**解决方案**

Composable 的返回值和函数名要清楚表达行为。涉及路由跳转、全局状态修改、弹窗提示的逻辑，要谨慎封装，并在命名上说明。

### 3. 重复请求

**症状**

页面一打开请求两次列表。

**常见原因**

- `onMounted(fetchList)` 调了一次。
- `watch(query, fetchList, { immediate: true })` 又调了一次。

**解决方案**

选择一个入口。推荐用 `watch` 管理依赖变化，或者用 `onMounted` 做首次加载，不要两个同时做相同事情。

## 最佳实践

- Composable 命名用 `useXxx`。
- 返回响应式状态和明确动作函数。
- 不要为了抽而抽，只抽重复、稳定、有边界的逻辑。
- 请求、路由、状态、弹窗提示这类副作用要显式命名。
- 页面组件应该能看出主要业务流程，不要把所有流程都藏进 composable。

## 下一步学习

继续学习 [路由与页面](/vue/router)，把组件和业务页面连接到真实 URL。
