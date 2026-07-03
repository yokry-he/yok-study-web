# 响应式基础

## 适合谁看

适合已经能写简单 Vue 页面，但对 `ref`、`reactive`、`computed`、`watch` 什么时候使用还不清楚的学习者。

响应式是 Vue 最核心的概念。你可以把它理解为：**数据变了，页面自动跟着变**。我们写 Vue 时，大多数工作都围绕“定义状态、修改状态、根据状态显示页面”展开。

## 你会学到什么

- `ref` 和 `reactive` 的区别。
- 为什么脚本里要写 `.value`。
- `computed` 和普通函数有什么区别。
- `watch` 适合解决什么问题。
- 实际项目里响应式最常见的坑怎么排查。

## 什么是响应式

先看一个最小例子：

```vue
<script setup lang="ts">
import { ref } from 'vue'

const count = ref(0)

function increase() {
  count.value++
}
</script>

<template>
  <button type="button" @click="increase">
    {{ count }}
  </button>
</template>
```

当点击按钮时，`count.value` 发生变化，模板中的 `{{ count }}` 自动更新。你不需要手动查找 DOM，也不需要写 `document.querySelector()`。

## ref：最常用的状态容器

`ref` 适合保存一个独立的值：

```ts
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const selectedIds = ref<number[]>([])
```

脚本中读取和修改要写 `.value`：

```ts
loading.value = true
page.value += 1
selectedIds.value.push(1001)
```

模板中不需要 `.value`：

```vue
<template>
  <p>当前页：{{ page }}</p>
</template>
```

新手可以先记住：**不知道用什么时，优先用 `ref`。**

## reactive：适合结构稳定的对象

`reactive` 适合保存一个不会整体替换的对象，例如表单：

```ts
import { reactive } from 'vue'

const form = reactive({
  username: '',
  mobile: '',
  enabled: true
})

form.username = 'alice'
```

不要这样做：

```ts
let form = reactive({
  username: '',
  mobile: ''
})

form = reactive({
  username: 'alice',
  mobile: '13800000000'
})
```

上面这种整体替换会让维护变复杂。真实项目里，如果你经常需要整体替换对象，更推荐：

```ts
const form = ref({
  username: '',
  mobile: ''
})

form.value = {
  username: 'alice',
  mobile: '13800000000'
}
```

## computed：用来写“由状态推导出的值”

`computed` 适合表达“这个值不是用户直接输入的，而是根据其他状态算出来的”。

```ts
const price = ref(99)
const quantity = ref(2)

const total = computed(() => price.value * quantity.value)
```

实际项目例子：

```ts
interface User {
  id: number
  name: string
  enabled: boolean
}

const users = ref<User[]>([])
const keyword = ref('')

const enabledUsers = computed(() => {
  return users.value.filter((user) => user.enabled)
})

const filteredUsers = computed(() => {
  return enabledUsers.value.filter((user) => user.name.includes(keyword.value))
})
```

`computed` 的特点是：只有依赖的数据变化时才重新计算。Vue 官方文档也强调，计算属性会基于响应式依赖缓存结果。

## watch：用来处理副作用

副作用是指“状态变化后，要额外做一件事”，例如：

- 搜索关键字变化后重新请求列表。
- 路由参数变化后重新加载详情。
- 弹窗打开后初始化表单。
- 选中组织变化后清空部门选择。

```ts
watch(keyword, async () => {
  page.value = 1
  await fetchUserList()
})
```

监听对象里的某个字段时，推荐写 getter：

```ts
const query = reactive({
  page: 1,
  keyword: ''
})

watch(
  () => query.keyword,
  async () => {
    query.page = 1
    await fetchList()
  }
)
```

如果页面加载时也要执行一次，可以加 `immediate`：

```ts
watch(
  () => route.params.id,
  async (id) => {
    await fetchDetail(String(id))
  },
  { immediate: true }
)
```

## watchEffect：先少用

`watchEffect` 会自动收集依赖，写起来短，但对初学者来说不够明确。

```ts
watchEffect(() => {
  document.title = keyword.value || '用户列表'
})
```

第一阶段建议优先使用 `watch`，因为它能清楚表达“我正在监听什么”。

## 实际项目常见问题

### 1. 解构后页面不更新

**问题代码**

```ts
const state = reactive({
  username: 'alice',
  age: 18
})

const { username } = state
```

`username` 被解构成普通变量后，和原来的响应式对象断开了。

**解决方案**

```ts
import { toRefs } from 'vue'

const state = reactive({
  username: 'alice',
  age: 18
})

const { username } = toRefs(state)
```

如果是 Pinia store，使用 `storeToRefs`：

```ts
const userStore = useUserStore()
const { profile, roles } = storeToRefs(userStore)
```

### 2. computed 里修改了状态

**问题代码**

```ts
const total = computed(() => {
  count.value++
  return count.value * 2
})
```

`computed` 应该是纯计算，不应该在里面修改其他状态。否则可能导致难以理解的重复更新。

**解决方案**

把修改状态放进事件函数或 action：

```ts
const total = computed(() => count.value * 2)

function increase() {
  count.value++
}
```

### 3. watch 触发太频繁

**场景**

搜索框每输入一个字就请求接口，接口压力大，页面也卡。

**解决方案**

加防抖：

```ts
let timer: number | undefined

watch(keyword, () => {
  window.clearTimeout(timer)
  timer = window.setTimeout(() => {
    fetchList()
  }, 300)
})
```

后续项目中可以封装成 `useDebounceFn` 或使用成熟工具库。

### 4. 表单重置不彻底

**常见原因**

编辑弹窗和新增弹窗共用同一个表单对象，关闭弹窗时没有恢复默认值。

**推荐写法**

```ts
interface UserForm {
  id?: number
  username: string
  mobile: string
  enabled: boolean
}

const defaultForm = (): UserForm => ({
  username: '',
  mobile: '',
  enabled: true
})

const form = ref<UserForm>(defaultForm())

function resetForm() {
  form.value = defaultForm()
}
```

用函数返回默认值，可以避免多个地方共用同一个对象引用。

## 选择建议

| 需求 | 推荐 |
| --- | --- |
| 一个字符串、数字、布尔值 | `ref` |
| 数组列表 | `ref<T[]>([])` |
| 表单对象，字段固定 | `reactive` 或 `ref<Form>()` |
| 需要整体替换对象 | `ref` |
| 根据状态计算展示值 | `computed` |
| 状态变化后请求接口或同步外部系统 | `watch` |

## 最佳实践

- 初学阶段优先用 `ref`，减少心智负担。
- `computed` 只做计算，不做请求、不改状态。
- `watch` 只在确实需要副作用时使用。
- 表单默认值用函数创建，避免对象引用污染。
- Pinia store 解构时使用 `storeToRefs`。

## 下一步学习

继续学习 [组件设计](/vue/component)。响应式负责“数据怎么变”，组件负责“页面怎么组织”。
