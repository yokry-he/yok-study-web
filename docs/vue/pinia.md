# Pinia 状态管理

## 适合谁看

适合已经会写组件和路由，但不清楚哪些状态应该放到全局、哪些应该留在组件内的学习者。

Pinia 是 Vue 官方推荐的状态管理库。它可以理解为“跨页面共享的数据仓库”，例如登录令牌、当前用户、权限、菜单、主题偏好。

## 你会学到什么

- Store 是什么。
- 哪些状态适合放 Pinia。
- Option Store 和 Setup Store 怎么写。
- 登录态和用户信息如何组织。
- 实际项目中刷新丢状态、store 互相依赖、解构丢响应怎么处理。

## 什么是 Store

Pinia 官方文档说明，Store 用 `defineStore()` 定义，并且需要一个唯一名称。

简单示例：

```ts
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    username: ''
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token)
  },
  actions: {
    setToken(token: string) {
      this.token = token
    }
  }
})
```

组件中使用：

```ts
const userStore = useUserStore()

if (userStore.isLoggedIn) {
  console.log(userStore.username)
}
```

## 哪些状态应该放 Pinia

| 状态 | 是否放 Pinia | 原因 |
| --- | --- | --- |
| token | 是 | 多个页面和请求拦截器都需要 |
| 当前用户信息 | 是 | 导航、权限、个人中心都会用 |
| 菜单和权限码 | 是 | 路由和按钮权限需要 |
| 主题模式 | 可以 | 多页面共享偏好 |
| 表格搜索条件 | 通常否 | 多数只属于当前页面 |
| 弹窗显示状态 | 通常否 | 多数只属于当前组件 |
| 表单输入内容 | 通常否 | 放全局会增加清理成本 |

记住一句话：**跨页面共享、刷新后可能要恢复、多个模块都依赖的状态，才考虑放 Pinia。**

## Setup Store 写法

Setup Store 更接近组合式 API：

```ts
import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const roles = ref<string[]>([])
  const username = ref('')

  const isLoggedIn = computed(() => Boolean(token.value))

  function setToken(value: string) {
    token.value = value
  }

  function clearUser() {
    token.value = ''
    roles.value = []
    username.value = ''
  }

  return {
    token,
    roles,
    username,
    isLoggedIn,
    setToken,
    clearUser
  }
})
```

如果你已经熟悉组合式 API，Setup Store 会更自然。

## 登录态 Store 示例

真实项目里可以这样组织：

```ts
interface UserProfile {
  id: number
  username: string
  nickname: string
  roles: string[]
  permissions: string[]
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const profile = ref<UserProfile | null>(null)

  const isLoggedIn = computed(() => Boolean(token.value))
  const permissions = computed(() => profile.value?.permissions ?? [])

  function setToken(value: string) {
    token.value = value
    localStorage.setItem('token', value)
  }

  async function fetchProfile() {
    profile.value = await userApi.getProfile()
  }

  function hasPermission(code: string) {
    return permissions.value.includes(code)
  }

  function logout() {
    token.value = ''
    profile.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    profile,
    isLoggedIn,
    permissions,
    setToken,
    fetchProfile,
    hasPermission,
    logout
  }
})
```

## Store 和接口层的边界

Store 可以调用 service 或 API，但不要把所有接口都塞进 Store。

推荐：

```text
api/user.ts       只描述请求
services/auth.ts 组织登录业务流程
stores/user.ts   保存登录态和用户上下文
```

例如：

```ts
async function login(payload: LoginPayload) {
  const result = await authApi.login(payload)
  userStore.setToken(result.token)
  await userStore.fetchProfile()
}
```

这样做比在组件里到处写登录流程更稳定。

## 实际项目常见问题

### 1. 刷新页面后 Pinia 状态丢失

**原因**

Pinia 默认存在内存中，浏览器刷新会重新加载应用，内存状态自然会丢。

**解决方案**

只持久化必要状态，例如 token、主题、语言。不要把完整用户信息、表单、菜单树全部无脑存 localStorage。

```ts
const token = ref(localStorage.getItem('token') || '')

function setToken(value: string) {
  token.value = value
  localStorage.setItem('token', value)
}
```

刷新后通过 token 重新请求用户信息。

### 2. 解构 store 后不更新

**问题代码**

```ts
const userStore = useUserStore()
const { profile, permissions } = userStore
```

**解决方案**

```ts
import { storeToRefs } from 'pinia'

const userStore = useUserStore()
const { profile, permissions } = storeToRefs(userStore)
```

action 可以直接解构：

```ts
const { logout } = userStore
```

### 3. Store 之间循环依赖

**症状**

应用启动报错，或者某个 getter 无限触发。

**原因**

两个 store 在初始化阶段互相读取对方状态。

**解决方案**

- 避免在 setup 顶层互相读取。
- 把依赖读取放到 action 或 computed 中。
- 合并职责过近的 store。

### 4. 所有页面状态都放 Pinia

**症状**

页面关闭后旧筛选条件、旧表单内容仍然影响新页面。

**原因**

把局部状态放进全局后，清理成本变高。

**解决方案**

表格筛选、弹窗开关、表单输入优先留在页面组件或 composable 中。只有确实需要跨页面共享时再进入 Pinia。

## 最佳实践

- Store 名称要唯一且稳定，例如 `user`、`permission`、`app`。
- 只把跨页面共享状态放 Pinia。
- 持久化要克制，只保存必要信息。
- 解构状态使用 `storeToRefs`。
- 登录态变化时同步清理路由、菜单、权限和缓存。

## 下一步学习

继续学习 [请求与接口封装](/vue/request)，把接口请求、错误处理和登录态连接起来。如果你正在做后台项目，继续看 [Vue Admin 菜单与动态路由实现手册](/vue/admin-menu-route-module)，学习菜单、权限和动态路由状态如何放进 Pinia。
