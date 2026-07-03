# Vue Router 速查

## 基础配置

```ts
import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index.vue'),
      meta: {
        title: '工作台',
        requiresAuth: true
      }
    }
  ]
})
```

## 跳转

```ts
const router = useRouter()

router.push('/dashboard')
router.replace('/login')
router.push({ name: 'UserDetail', params: { id: 1 } })
```

## 读取路由

```ts
const route = useRoute()

const id = route.params.id
const redirect = route.query.redirect
```

## 路由守卫

```ts
router.beforeEach((to) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.token) {
    return {
      path: '/login',
      query: { redirect: to.fullPath }
    }
  }

  return true
})
```

## 动态路由

```ts
if (!router.hasRoute(route.name)) {
  router.addRoute(route)
}
```

动态注册后，如果当前地址需要重新匹配：

```ts
return to.fullPath
```

## 常见问题

| 问题 | 解决方案 |
| --- | --- |
| 刷新 404 | 服务器配置 fallback 到 `index.html` |
| 参数变化页面不刷新 | watch 路由参数 |
| 动态路由重复 | `router.hasRoute()` 判断 |
| 登录后跳不回原页面 | 使用 `redirect=to.fullPath` |
| 部署子路径失效 | Vite `base` 和 Router base 保持一致 |

## 参数变化监听

```ts
watch(
  () => route.params.id,
  (id) => {
    fetchDetail(String(id))
  },
  { immediate: true }
)
```
