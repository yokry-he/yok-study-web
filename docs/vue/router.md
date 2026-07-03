# 路由与页面

## 适合谁看

适合已经能写 Vue 组件，但对“页面如何组织”“URL 和组件如何对应”“登录权限怎么拦截”还不清楚的学习者。

Vue Router 是 Vue 官方路由方案。它让单页应用可以通过 URL 切换页面，而不是每次都从服务器重新加载整个页面。

## 你会学到什么

- 路由和页面组件的关系。
- 静态路由、动态路由、嵌套路由怎么写。
- `meta` 应该放什么信息。
- 登录拦截和权限路由的基本流程。
- 实际项目中刷新 404、动态路由重复、参数变化不刷新怎么解决。

## 路由是什么

路由可以理解为一张表：

| URL | 页面组件 |
| --- | --- |
| `/login` | 登录页 |
| `/dashboard` | 工作台 |
| `/users` | 用户管理 |
| `/users/1001` | 用户详情 |

配置示例：

```ts
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue')
    },
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

export default router
```

## 静态路由和动态路由

### 静态路由

所有用户都能访问的路由，例如登录页、404 页、首页：

```ts
const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/404',
    component: () => import('@/views/error/404.vue')
  }
]
```

### 动态路由

根据用户权限动态添加的路由，例如用户管理、角色管理、系统设置：

```ts
const asyncRoutes = [
  {
    path: '/system/users',
    name: 'SystemUsers',
    component: () => import('@/views/system/users/index.vue'),
    meta: {
      title: '用户管理',
      permission: 'system:user:list'
    }
  }
]
```

Vue Router 官方动态路由文档建议，如果在导航守卫中添加路由，需要通过返回目标地址触发重定向，而不是手动调用 `router.replace()`。

## 嵌套路由

后台系统通常有一个布局壳：

```text
Layout
├─ Dashboard
├─ System
│  ├─ Users
│  └─ Roles
└─ Settings
```

配置：

```ts
const routes = [
  {
    path: '/',
    component: () => import('@/layouts/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台', requiresAuth: true }
      },
      {
        path: 'system/users',
        name: 'SystemUsers',
        component: () => import('@/views/system/users/index.vue'),
        meta: { title: '用户管理', requiresAuth: true }
      }
    ]
  }
]
```

注意：子路由 path 不以 `/` 开头时，会拼接父路由路径。

## Route Meta 放什么

`meta` 适合放路由相关的元信息：

```ts
meta: {
  title: '用户管理',
  requiresAuth: true,
  permission: 'system:user:list',
  keepAlive: true,
  icon: 'users'
}
```

不要把复杂业务状态放进 `meta`。例如当前筛选条件、表单数据、接口结果都不应该放这里。

## 登录拦截流程

常见流程：

```text
进入页面
↓
是否需要登录？
↓
没有 token -> 去登录页
↓
有 token 但没有用户信息 -> 拉取用户信息和权限
↓
生成动态路由和菜单
↓
进入目标页面
```

示例：

```ts
router.beforeEach(async (to) => {
  const userStore = useUserStore()

  if (!to.meta.requiresAuth) {
    return true
  }

  if (!userStore.token) {
    return {
      path: '/login',
      query: { redirect: to.fullPath }
    }
  }

  if (!userStore.profile) {
    await userStore.fetchProfile()
  }

  return true
})
```

登录成功后回到原页面：

```ts
const redirect = route.query.redirect?.toString() || '/dashboard'
router.replace(redirect)
```

## 路由参数

详情页通常使用动态参数：

```ts
{
  path: '/users/:id',
  name: 'UserDetail',
  component: () => import('@/views/users/detail.vue')
}
```

页面中读取：

```ts
const route = useRoute()

watch(
  () => route.params.id,
  async (id) => {
    await fetchUserDetail(String(id))
  },
  { immediate: true }
)
```

为什么不用 `onMounted`？因为从 `/users/1` 切到 `/users/2` 时，可能复用同一个组件实例，`onMounted` 不会再次执行。

## 实际项目常见问题

### 1. 刷新页面后 404

**症状**

访问 `/system/users` 正常，但刷新后服务器返回 404。

**原因**

前端路由是浏览器里的路由。刷新时，浏览器会向服务器请求 `/system/users`。如果服务器没有配置回退到 `index.html`，就会 404。

**解决方案**

Nginx 示例：

```nginx
location / {
  try_files $uri $uri/ /index.html;
}
```

如果项目部署在子路径，例如 `/admin/`，还需要配置 Vite 的 `base` 和 Router 的 history base。

### 2. 动态路由重复注册

**症状**

登录后菜单重复，控制台提示重复路由名，或者页面跳转异常。

**原因**

每次进入守卫都重新 `addRoute`，没有判断是否已经注册。

**解决方案**

```ts
if (!router.hasRoute(route.name)) {
  router.addRoute(route)
}
```

退出登录时要清理用户相关状态。复杂系统里可以记录动态路由名称，退出时逐个移除。

### 3. 路由参数变化但页面数据没刷新

**症状**

从用户 A 详情跳到用户 B 详情，URL 变了，但页面仍显示 A。

**原因**

组件被复用，`onMounted` 没有重新执行。

**解决方案**

监听参数：

```ts
watch(
  () => route.params.id,
  (id) => fetchDetail(String(id)),
  { immediate: true }
)
```

### 4. 登录后跳回登录页

**常见原因**

- token 保存了，但请求用户信息失败。
- 请求拦截器没有带上 token。
- 权限接口返回慢，守卫没有等待完成。
- 登录页本身也被错误标记为需要登录。

**排查顺序**

1. 看 Network 里用户信息接口是否成功。
2. 看请求头是否带了 Authorization。
3. 看 Pinia 中 token 和 profile 是否正确。
4. 看路由 meta 是否配置错误。

## 最佳实践

- 页面导航必须使用真实路由，不用首页锚点冒充页面。
- 路由名 `name` 要稳定唯一。
- `meta` 只放路由相关信息。
- 动态路由注册要防重复。
- 参数变化的数据请求使用 `watch`。
- 部署前必须验证刷新非首页路径不会 404。

## 下一步学习

继续学习 [Pinia 状态管理](/vue/pinia)，把登录态、用户信息、菜单和权限状态组织起来。如果你正在做后台项目，继续看 [Vue Admin 菜单与动态路由实现手册](/vue/admin-menu-route-module)。
