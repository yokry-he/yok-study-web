# Vue 速查

## 响应式 API

| API | 用途 | 示例 |
| --- | --- | --- |
| `ref` | 定义基础状态或数组 | `const count = ref(0)` |
| `reactive` | 定义结构稳定的对象 | `const form = reactive({ name: '' })` |
| `computed` | 定义派生状态 | `const total = computed(() => price.value * count.value)` |
| `watch` | 状态变化后执行副作用 | `watch(keyword, fetchList)` |
| `nextTick` | 等待 DOM 更新完成 | `await nextTick()` |

## 组件 API

```ts
defineProps<{
  title: string
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: FormData]
  cancel: []
}>()
```

## 常用指令

| 指令 | 用途 | 注意 |
| --- | --- | --- |
| `v-if` | 条件渲染 | 切换时会创建/销毁 DOM |
| `v-show` | 显示隐藏 | DOM 一直存在 |
| `v-for` | 列表渲染 | 必须写稳定 `key` |
| `v-model` | 双向绑定 | 常用于表单 |
| `v-bind` / `:` | 绑定属性 | `:disabled="loading"` |
| `v-on` / `@` | 绑定事件 | `@click="submit"` |

## 常见模板

### 列表

```vue
<ul>
  <li v-for="user in users" :key="user.id">
    {{ user.username }}
  </li>
</ul>
```

### 加载状态

```vue
<div v-if="loading">加载中...</div>
<UserTable v-else :users="users" />
```

### 空状态

```vue
<EmptyState v-if="!loading && users.length === 0" />
```

## 常见坑

| 问题 | 正确处理 |
| --- | --- |
| `reactive` 解构后不更新 | 使用 `toRefs` |
| Pinia 解构后不更新 | 使用 `storeToRefs` |
| `computed` 里修改状态 | 改到事件函数或 action |
| `v-for` 使用 index 做 key | 使用业务唯一 id |
| 路由参数变化页面不刷新 | `watch(() => route.params.id)` |

## 选择建议

| 场景 | 推荐 |
| --- | --- |
| 简单状态 | `ref` |
| 表单对象 | `reactive` 或 `ref<Form>()` |
| 展示值由状态计算 | `computed` |
| 状态变化后请求接口 | `watch` |
| 多页面共享状态 | Pinia |
