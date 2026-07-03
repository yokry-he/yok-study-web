# 组件设计

## 适合谁看

适合已经能写单个 Vue 页面，但页面越写越长、组件之间传值混乱、不知道什么时候拆组件的学习者。

Vue 项目维护成本高，通常不是因为某个 API 很难，而是因为组件边界不清楚。组件设计的核心不是“拆得越多越好”，而是让每个组件职责清楚、输入输出稳定。

## 你会学到什么

- 组件到底解决什么问题。
- props、emits、slot 分别适合什么场景。
- 页面组件、业务组件、基础组件怎么区分。
- 表单、弹窗、表格类组件怎么拆。
- 实际项目里组件失控时怎么治理。

## 组件是什么

组件可以理解为“带逻辑的页面积木”。它把一块 UI 和这块 UI 需要的状态、事件组织在一起。

例如用户管理页面可以拆成：

```text
UserPage
├─ UserSearchForm
├─ UserTable
└─ UserFormDrawer
```

每个组件只关心自己的事情：

| 组件 | 职责 |
| --- | --- |
| `UserPage` | 组织整体流程：查询、打开弹窗、刷新列表 |
| `UserSearchForm` | 收集筛选条件 |
| `UserTable` | 展示列表和操作按钮 |
| `UserFormDrawer` | 新增和编辑用户 |

## 组件分类

### 页面组件

放在 `views/` 目录下，负责一个路由页面。

页面组件可以知道业务流程，例如：

- 当前页面需要请求哪个接口。
- 搜索后如何刷新列表。
- 点击编辑时打开哪个弹窗。
- 删除成功后如何更新页面。

### 业务组件

放在页面目录或 `components/` 下，承载某个业务片段。

例如：

- `UserFormDrawer`
- `RolePermissionTree`
- `DepartmentSelector`

业务组件可以理解业务字段，但不应该知道整个系统的所有流程。

### 基础组件

基础组件是跨业务复用的 UI 能力，例如：

- `AppPage`
- `StatusTag`
- `ConfirmButton`
- `EmptyState`

如果项目已经选定组件库，按钮、输入框、表格、弹窗、抽屉等基础控件应优先使用组件库，不要重复手写。

## Props：父组件传给子组件的数据

Props 适合传入“子组件展示或处理所需的数据”。

```vue
<script setup lang="ts">
interface User {
  id: number
  username: string
  mobile: string
  enabled: boolean
}

defineProps<{
  users: User[]
  loading: boolean
}>()
</script>

<template>
  <div v-if="loading">加载中...</div>
  <ul v-else>
    <li v-for="user in users" :key="user.id">
      {{ user.username }} - {{ user.mobile }}
    </li>
  </ul>
</template>
```

Props 要保持只读。子组件不要直接修改父组件传进来的对象。

**错误示例**

```ts
const props = defineProps<{ user: User }>()
props.user.username = 'new name'
```

**推荐做法**

如果需要编辑，复制一份到当前组件的表单状态：

```ts
const props = defineProps<{ user: User }>()

const form = ref({
  username: props.user.username,
  mobile: props.user.mobile
})
```

## Emits：子组件通知父组件发生了什么

Emits 适合表达“用户在子组件里做了一个动作”。

```vue
<script setup lang="ts">
interface User {
  id: number
  username: string
}

defineProps<{
  users: User[]
}>()

const emit = defineEmits<{
  edit: [user: User]
  remove: [id: number]
}>()
</script>

<template>
  <ul>
    <li v-for="user in users" :key="user.id">
      {{ user.username }}
      <button type="button" @click="emit('edit', user)">编辑</button>
      <button type="button" @click="emit('remove', user.id)">删除</button>
    </li>
  </ul>
</template>
```

父组件接收事件：

```vue
<UserTable
  :users="users"
  @edit="openEditDrawer"
  @remove="confirmRemove"
/>
```

命名建议：

| 动作 | 推荐事件名 |
| --- | --- |
| 打开编辑 | `edit` |
| 删除 | `remove` |
| 提交表单 | `submit` |
| 关闭弹窗 | `close` |
| 更新值 | `update:modelValue` |

## Slot：父组件传入一段 UI

Slot 适合“结构由子组件控制，但部分内容由父组件决定”的场景。

例如页面容器：

```vue
<template>
  <section class="app-page">
    <header class="app-page__header">
      <h1>{{ title }}</h1>
      <slot name="actions" />
    </header>

    <main class="app-page__body">
      <slot />
    </main>
  </section>
</template>
```

使用：

```vue
<AppPage title="用户管理">
  <template #actions>
    <button type="button">新增用户</button>
  </template>

  <UserTable :users="users" />
</AppPage>
```

Slot 不适合用来绕过组件边界。如果父组件需要控制子组件内部太多细节，说明组件 API 可能设计错了。

## v-model：适合双向绑定表单类组件

自定义组件支持 `v-model`：

```vue
<script setup lang="ts">
const model = defineModel<string>()
</script>

<template>
  <input v-model="model" />
</template>
```

使用：

```vue
<KeywordInput v-model="keyword" />
```

适合使用 `v-model` 的组件：

- 输入框。
- 选择器。
- 开关。
- 日期选择。
- 弹窗显示状态。

不适合滥用 `v-model` 的情况：

- 一个复杂业务组件暴露十几个双向绑定字段。
- 子组件内部直接修改父组件业务对象。

## 表格页面拆分示例

用户管理页面可以这样写：

```vue
<script setup lang="ts">
import UserFormDrawer from './UserFormDrawer.vue'
import UserSearchForm from './UserSearchForm.vue'
import UserTable from './UserTable.vue'

const query = ref({
  keyword: '',
  enabled: undefined as boolean | undefined
})

const users = ref<User[]>([])
const drawerVisible = ref(false)
const editingUser = ref<User | null>(null)

function openCreateDrawer() {
  editingUser.value = null
  drawerVisible.value = true
}

function openEditDrawer(user: User) {
  editingUser.value = user
  drawerVisible.value = true
}
</script>

<template>
  <UserSearchForm v-model="query" @search="fetchUsers" />

  <UserTable
    :users="users"
    @create="openCreateDrawer"
    @edit="openEditDrawer"
  />

  <UserFormDrawer
    v-model:visible="drawerVisible"
    :user="editingUser"
    @success="fetchUsers"
  />
</template>
```

这个页面的好处是：父组件负责流程，子组件负责具体 UI。

## 实际项目常见问题

### 1. 组件 props 越传越多

**症状**

一个组件接收二三十个 props，每次使用都要写很长一串。

**原因**

组件职责太大，或者把多个场景强行塞进一个组件。

**解决方案**

- 拆成多个更小的组件。
- 把相关配置合并成对象，例如 `tableOptions`。
- 区分“基础组件”和“业务组件”，不要让基础组件理解太多业务。

### 2. 子组件偷偷改父组件数据

**症状**

关闭弹窗后，列表里的数据已经被改了，但用户还没点保存。

**原因**

编辑表单直接使用了父组件传入的对象引用。

**解决方案**

打开弹窗时深拷贝一份表单数据：

```ts
function openEdit(user: User) {
  form.value = {
    id: user.id,
    username: user.username,
    mobile: user.mobile,
    enabled: user.enabled
  }
}
```

保存成功后再通知父组件刷新列表。

### 3. 组件里既请求接口又控制路由又改全局状态

**症状**

组件很难复用，也很难测试。改一个按钮可能影响登录态或路由。

**原因**

业务流程没有分层。

**解决方案**

- 接口请求放到 `api` 或 `service`。
- 跨页面状态放到 Pinia。
- 页面组件组织业务流程。
- 展示组件只接收 props 和发出 emits。

### 4. 移动端布局被挤坏

**症状**

头像变成椭圆，操作按钮被压缩成窄条，表格操作列换行错乱。

**原因**

固定尺寸元素没有设置稳定宽高和不可压缩行为。

**解决方案**

```css
.user-avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}

.toolbar-action {
  flex-shrink: 0;
}
```

同时避免使用宽泛选择器覆盖组件库内部 DOM。

## 组件设计检查清单

- 这个组件的职责能用一句话说清楚吗？
- 它的 props 是必要的吗？
- 它发出的 emits 是业务动作，而不是内部细节吗？
- 它是否直接修改了父组件对象？
- 它是否混入了不该属于自己的请求、路由、全局状态？
- 它在移动端是否仍然可读可操作？

## 最佳实践

- 页面组件组织业务流程，展示组件保持输入输出清晰。
- props 只读，修改通过 emits 通知父组件。
- 表单编辑时复制数据，不直接改列表对象。
- 使用组件库时不要依赖内部 DOM 写样式。
- 组件超过 250 行或职责超过 2 个时，优先考虑拆分。

## 下一步学习

继续学习 [组合式 API](/vue/composition-api)，把可复用逻辑从组件中抽出来。
