# Pinia 速查

## Setup Store

```ts
import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const roles = ref<string[]>([])

  const isLoggedIn = computed(() => Boolean(token.value))

  function setToken(value: string) {
    token.value = value
  }

  function logout() {
    token.value = ''
    roles.value = []
  }

  return {
    token,
    roles,
    isLoggedIn,
    setToken,
    logout
  }
})
```

## 组件中使用

```ts
const userStore = useUserStore()

userStore.setToken('token-value')
```

解构状态：

```ts
const { token, roles } = storeToRefs(userStore)
```

action 可以直接解构：

```ts
const { logout } = userStore
```

## 什么适合放 Pinia

| 状态 | 是否适合 |
| --- | --- |
| token | 适合 |
| 当前用户信息 | 适合 |
| 权限码 | 适合 |
| 菜单 | 适合 |
| 主题偏好 | 可以 |
| 表单输入 | 通常不适合 |
| 弹窗显示 | 通常不适合 |
| 表格筛选 | 通常不适合 |

## 持久化 token

```ts
const token = ref(localStorage.getItem('token') || '')

function setToken(value: string) {
  token.value = value
  localStorage.setItem('token', value)
}

function clearToken() {
  token.value = ''
  localStorage.removeItem('token')
}
```

## 常见问题

| 问题 | 解决方案 |
| --- | --- |
| 刷新后状态丢失 | 必要状态持久化，刷新后重新请求 |
| 解构后不更新 | 使用 `storeToRefs` |
| store 互相依赖报错 | 避免初始化阶段循环读取 |
| 页面旧状态残留 | 局部状态不要放 Pinia |
| 退出后菜单还在 | logout 时清理 user、permission、router |
