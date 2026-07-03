# Vue 项目集成

## 适合谁看

适合准备在 Vue 3 项目中实际使用 TypeScript 的学习者。

## ref 类型

数组：

```ts
const users = ref<User[]>([])
```

可能为空：

```ts
const currentUser = ref<User | null>(null)
```

DOM 引用：

```ts
const inputRef = ref<HTMLInputElement | null>(null)
```

## props 类型

```vue
<script setup lang="ts">
interface Props {
  users: User[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})
</script>
```

## emits 类型

```ts
const emit = defineEmits<{
  edit: [user: User]
  remove: [id: number]
  submit: [payload: UserForm]
}>()
```

调用时如果参数错误，编辑器会提示。

## defineModel 类型

Vue 3.4+：

```ts
const visible = defineModel<boolean>('visible', { default: false })
const keyword = defineModel<string>({ default: '' })
```

适合自定义表单组件和弹窗显示状态。

## Pinia 类型

```ts
interface UserProfile {
  id: number
  username: string
  permissions: string[]
}

export const useUserStore = defineStore('user', () => {
  const profile = ref<UserProfile | null>(null)
  const permissions = computed(() => profile.value?.permissions ?? [])

  return {
    profile,
    permissions
  }
})
```

解构：

```ts
const userStore = useUserStore()
const { profile, permissions } = storeToRefs(userStore)
```

## API 请求类型

```ts
interface UserQuery {
  page: number
  pageSize: number
  keyword?: string
}

function getUserList(params: UserQuery) {
  return request.get<PageResult<User>>('/users', { params })
}
```

页面使用：

```ts
const users = ref<User[]>([])

async function fetchUsers() {
  const result = await getUserList(query.value)
  users.value = result.items
}
```

## 表单类型

```ts
interface UserForm {
  id?: number
  username: string
  mobile: string
  enabled: boolean
}

function createDefaultUserForm(): UserForm {
  return {
    username: '',
    mobile: '',
    enabled: true
  }
}
```

## 实际项目常见问题

### 1. props 默认值写不对

使用 `withDefaults`：

```ts
withDefaults(defineProps<Props>(), {
  loading: false
})
```

### 2. emit 参数传错

用类型约束：

```ts
const emit = defineEmits<{
  success: []
  submit: [payload: UserForm]
}>()
```

### 3. ref 初始 null 后访问报错

先判断：

```ts
if (!currentUser.value) return

console.log(currentUser.value.username)
```

### 4. 请求返回 any

**问题**

页面使用接口结果时没有提示，字段写错也不报错。

**解决方案**

给 request 方法传泛型，或者在 API 函数声明返回类型。

## 最佳实践

- 每个业务实体有 interface。
- 表单和实体分开。
- props、emits、ref 都加类型。
- API 函数明确返回类型。
- Pinia 解构使用 `storeToRefs`。
- 不用 any 跳过问题。

## 下一步

继续学习 [常见问题](/typescript/troubleshooting)。
