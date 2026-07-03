# 快速开始

## 适合谁看

适合刚开始学习 Vue 3，或者之前只零散写过 Vue 页面、想重新建立完整开发流程的同学。

这一节不追求一次讲完所有 API，而是让你先知道一个 Vue 项目从创建、启动、写组件到拆分目录的基本过程。你学完后应该能独立创建一个小项目，并写出一个列表页面。

## 你会学到什么

- Vue 项目通常由哪些文件组成。
- 单文件组件 `.vue` 的三个区域分别做什么。
- `ref`、事件、模板绑定的最小用法。
- 一个页面什么时候该拆成组件。
- 新手常见启动问题怎么排查。

## 第一步：创建项目

推荐使用官方脚手架创建 Vue 项目：

```bash
npm create vue@latest
```

创建过程中会出现一些选项。第一阶段建议这样选：

| 选项 | 建议 | 原因 |
| --- | --- | --- |
| TypeScript | 是 | 真实项目里更容易维护接口、表单、状态 |
| JSX | 否 | 初学阶段先掌握模板语法 |
| Vue Router | 是 | 多页面项目一定会用到路由 |
| Pinia | 是 | 后续管理登录态、用户信息、菜单 |
| Vitest | 暂时否 | 等业务结构稳定后再补 |
| ESLint | 是 | 尽早养成规范习惯 |
| Prettier | 是 | 自动格式化，减少无意义争论 |

创建后进入项目并启动：

```bash
cd your-vue-project
npm install
npm run dev
```

看到本地地址后，在浏览器打开即可。

## 第二步：理解项目目录

一个常见 Vue 3 项目可以这样组织：

```text
src/
├─ api/             接口请求函数，例如 userApi.getList()
├─ assets/          图片、字体等静态资源
├─ components/      可复用组件，例如 UserForm、PageHeader
├─ composables/     可复用逻辑，例如 usePagination
├─ router/          路由表和路由守卫
├─ stores/          Pinia 全局状态
├─ styles/          全局样式、变量、重置样式
├─ types/           TypeScript 类型
├─ utils/           工具函数，例如 formatDate
└─ views/           页面，例如 users/index.vue
```

先记住一个简单原则：

> 页面放 `views`，可复用 UI 放 `components`，可复用逻辑放 `composables`，跨页面状态放 `stores`，接口函数放 `api`。

这样做的好处是后续项目变大时，不会把请求、状态、页面展示和工具函数全部堆在一个文件里。

## 第三步：认识单文件组件

Vue 的 `.vue` 文件通常由三部分组成：

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
    点击次数：{{ count }}
  </button>
</template>

<style scoped>
button {
  padding: 8px 12px;
}
</style>
```

| 区域 | 作用 | 初学理解 |
| --- | --- | --- |
| `<script setup>` | 写数据、函数、导入组件 | 页面逻辑 |
| `<template>` | 写页面结构 | 页面长什么样 |
| `<style scoped>` | 写当前组件样式 | 当前组件怎么显示 |

`ref(0)` 创建了一个响应式数据。脚本里改它时要写 `count.value`，模板里可以直接写 `count`。

## 第四步：写一个真实一点的列表

先定义类型和初始数据：

```vue
<script setup lang="ts">
import { computed, ref } from 'vue'

interface TodoItem {
  id: number
  title: string
  done: boolean
}

const keyword = ref('')
const todos = ref<TodoItem[]>([
  { id: 1, title: '学习 Vue 组件', done: true },
  { id: 2, title: '整理 Pinia 状态', done: false },
  { id: 3, title: '完成后台列表页', done: false }
])

const filteredTodos = computed(() => {
  return todos.value.filter((item) => item.title.includes(keyword.value))
})

function toggleTodo(id: number) {
  const target = todos.value.find((item) => item.id === id)
  if (target) {
    target.done = !target.done
  }
}
</script>

<template>
  <section>
    <input v-model="keyword" placeholder="搜索任务" />

    <ul>
      <li v-for="todo in filteredTodos" :key="todo.id">
        <label>
          <input
            type="checkbox"
            :checked="todo.done"
            @change="toggleTodo(todo.id)"
          />
          <span>{{ todo.title }}</span>
        </label>
      </li>
    </ul>
  </section>
</template>
```

这个例子包含了真实项目最常见的几个动作：

- 用接口或本地数据生成列表。
- 用 `v-model` 绑定搜索条件。
- 用 `computed` 得到过滤后的列表。
- 用 `v-for` 渲染数据。
- 用事件函数修改状态。

## 第五步：什么时候拆组件

初学者很容易过早拆组件，或者完全不拆。可以按下面的标准判断：

| 情况 | 是否拆 | 原因 |
| --- | --- | --- |
| 一段 UI 在多个页面重复出现 | 拆 | 复用价值高 |
| 一个页面文件超过 250 行 | 考虑拆 | 阅读成本开始变高 |
| 表单、弹窗、表格操作区逻辑独立 | 拆 | 便于维护和测试 |
| 只有两三行简单 HTML | 先不拆 | 过度拆分会增加跳转成本 |

例如用户管理页面可以拆成：

```text
views/users/index.vue          页面入口，组织筛选、列表、弹窗
views/users/UserSearch.vue     搜索表单
views/users/UserTable.vue      用户表格
views/users/UserFormDrawer.vue 新增/编辑抽屉表单
```

## 实际项目常见问题

### 1. 启动后页面空白

**常见原因**

- 控制台有运行时报错。
- 路由配置的页面路径写错。
- 组件导入路径大小写不一致。
- `main.ts` 没有正确挂载应用。

**排查顺序**

1. 打开浏览器控制台，看第一条红色错误。
2. 检查终端是否有编译错误。
3. 检查 `router/index.ts` 中的组件路径。
4. 检查 `main.ts` 是否有 `createApp(App).mount('#app')`。

### 2. 修改代码后页面没有变化

**常见原因**

- 开发服务器没有启动在当前项目目录。
- 浏览器打开的是旧端口。
- 文件保存失败。
- HMR 出现临时异常。

**解决方案**

- 确认终端里的本地地址和浏览器地址一致。
- 强制刷新页面。
- 停掉开发服务器，重新执行 `npm run dev`。

### 3. 路径别名 `@` 报错

**常见原因**

Vite 和 TypeScript 没有同时配置路径别名。

**解决思路**

Vite 负责运行时解析，TypeScript 负责编辑器类型提示。两边都要配置，否则可能出现“能运行但编辑器报错”或“编辑器不报错但构建失败”。

## 最佳实践

- 新项目优先使用 `<script setup lang="ts">`。
- 页面先能跑通，再抽组件和 composable。
- 不要把接口请求、权限判断、表单校验、弹窗状态全部写进一个大函数。
- `v-for` 必须使用稳定唯一的 `key`，不要用数组下标。
- 每完成一个页面，就补一段“这个页面的数据从哪里来、状态怎么流转”的说明。

## 下一步学习

继续学习 [响应式基础](/vue/reactivity)。响应式是 Vue 的核心，理解它之后，组件、Pinia 和表单逻辑都会更容易。
