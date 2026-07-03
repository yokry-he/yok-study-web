# 权限与菜单

## 适合谁看

适合准备开发后台管理系统，或者已经遇到“菜单怎么动态生成”“按钮权限怎么控制”“前端权限和后端权限怎么分工”的学习者。

后台权限是 Vue Admin 项目最容易混乱的部分。它不只是隐藏几个按钮，而是涉及用户、角色、权限点、路由、菜单和接口校验。

## 你会学到什么

- 权限系统的基本模型。
- 路由权限、菜单权限、按钮权限、数据权限的区别。
- 前端如何根据权限生成路由和菜单。
- 按钮权限如何封装。
- 实际项目中权限不生效、刷新丢菜单、前后端权限不一致怎么处理。

## 基本模型

常见权限模型：

```text
用户 -> 角色 -> 权限点 -> 页面 / 菜单 / 按钮 / 数据范围
```

例子：

| 用户 | 角色 | 权限 |
| --- | --- | --- |
| Alice | 系统管理员 | 用户新增、用户编辑、角色授权 |
| Bob | 运营人员 | 用户查看、订单查看 |
| Cindy | 审计人员 | 日志查看、报表查看 |

权限码建议稳定且可读：

```text
system:user:list
system:user:create
system:user:update
system:user:delete
system:role:assign
```

## 四类权限

| 类型 | 前端表现 | 最终校验 |
| --- | --- | --- |
| 路由权限 | 能否进入页面 | 后端仍需校验接口 |
| 菜单权限 | 是否显示菜单 | 前端控制显示 |
| 按钮权限 | 是否显示或禁用按钮 | 后端校验操作接口 |
| 数据权限 | 能看到哪些数据 | 必须后端控制 |

重要原则：

> 前端权限主要控制用户体验，后端权限负责安全边界。

用户即使在前端隐藏了删除按钮，也可能通过手写请求调用删除接口。所以后端必须校验。

## 路由和权限码

路由 meta 中放权限码：

```ts
const routes = [
  {
    path: '/system/users',
    name: 'SystemUsers',
    component: () => import('@/views/system/users/index.vue'),
    meta: {
      title: '用户管理',
      requiresAuth: true,
      permission: 'system:user:list'
    }
  }
]
```

过滤路由：

```ts
function hasRoutePermission(route: AppRoute, permissions: string[]) {
  const permission = route.meta?.permission
  return !permission || permissions.includes(permission)
}

function filterRoutes(routes: AppRoute[], permissions: string[]) {
  return routes
    .filter((route) => hasRoutePermission(route, permissions))
    .map((route) => ({
      ...route,
      children: route.children
        ? filterRoutes(route.children, permissions)
        : undefined
    }))
}
```

## 菜单生成

菜单可以由路由生成，也可以由后端返回。第一阶段推荐“后端返回权限码，前端根据本地路由表生成菜单”。

优点：

- 页面组件路径由前端掌控。
- 菜单标题、图标、排序可以统一维护。
- 后端只关心用户拥有哪些权限。

菜单项可以来自路由 meta：

```ts
meta: {
  title: '用户管理',
  icon: 'users',
  permission: 'system:user:list',
  order: 20
}
```

## 按钮权限

按钮权限不建议在模板里到处写复杂判断：

```vue
<button v-if="userStore.permissions.includes('system:user:create')">
  新增用户
</button>
```

可以封装函数：

```ts
export function usePermission() {
  const userStore = useUserStore()

  function can(code: string) {
    return userStore.permissions.includes(code)
  }

  return { can }
}
```

使用：

```vue
<script setup lang="ts">
const { can } = usePermission()
</script>

<template>
  <button v-if="can('system:user:create')" type="button">
    新增用户
  </button>
</template>
```

如果项目使用组件库，可以进一步封装 `PermissionButton`，统一处理展示、禁用、tooltip 和点击行为。

## 登录后初始化权限

流程：

```text
登录成功
↓
保存 token
↓
请求当前用户信息
↓
拿到 roles 和 permissions
↓
过滤动态路由
↓
注册路由
↓
生成菜单
↓
进入目标页面
```

示例：

```ts
async function initUserContext() {
  const userStore = useUserStore()
  const permissionStore = usePermissionStore()

  await userStore.fetchProfile()

  const routes = filterRoutes(asyncRoutes, userStore.permissions)
  permissionStore.setRoutes(routes)

  routes.forEach((route) => {
    if (!router.hasRoute(route.name)) {
      router.addRoute(route)
    }
  })
}
```

## 实际项目常见问题

### 1. 刷新后菜单没了

**原因**

菜单只存在内存里，刷新后 Pinia 重置，但没有重新拉用户信息和权限。

**解决方案**

路由守卫中判断：

- 有 token。
- 没有 profile 或菜单。
- 重新请求用户信息。
- 重新生成动态路由和菜单。

### 2. 有菜单但进入页面 404

**常见原因**

- 菜单数据和路由表不一致。
- 动态路由还没注册就跳转。
- 路由 name 重复。
- 后端返回了前端不存在的权限码。

**解决方案**

菜单生成要以可用路由为准。权限码只是过滤条件，不应该直接拼页面组件路径。

### 3. 按钮隐藏了，但用户仍然能调用接口

**原因**

前端隐藏按钮只影响界面，不等于安全。

**解决方案**

后端必须校验操作权限。前端权限只是减少无效操作入口。

### 4. 超级管理员权限难维护

**症状**

代码里到处写 `if role === 'admin'`。

**解决方案**

把超级管理员转换成权限判断能力：

```ts
function can(code: string) {
  return userStore.isSuperAdmin || userStore.permissions.includes(code)
}
```

不要在每个页面重复判断角色名。

### 5. 权限码命名混乱

**症状**

同一个能力有多个叫法：`user:add`、`system:user:create`、`user:create`。

**解决方案**

建立命名规则：

```text
模块:资源:动作
system:user:create
system:user:update
system:user:delete
system:role:assign
```

权限码一旦上线，尽量不要随意改名。需要变更时要做兼容和迁移说明。

## 最佳实践

- 前端控制展示，后端保证安全。
- 权限码使用稳定命名规则。
- 菜单从可用路由生成，避免菜单和页面脱节。
- 动态路由注册要防重复。
- 刷新后要能恢复用户上下文和菜单。
- 按钮权限封装成函数、指令或组件，不要散落在模板里。

## 下一步学习

继续查看 [Vue Admin 权限路由闭环实战](/vue/admin-permission-route-flow)、[Vue Admin 角色权限模块实现手册](/vue/admin-permission-module)、[Vue Admin 菜单与动态路由实现手册](/vue/admin-menu-route-module) 和 [常见问题](/vue/troubleshooting)，把基础权限概念接到真实后台项目里。
