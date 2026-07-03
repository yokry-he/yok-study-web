# Vue Admin 实战

## 项目目标

Vue Admin 是最适合练习企业级 Vue 项目的实战类型。它能把 Vue 3、Vue Router、Pinia、请求封装、权限控制、表格、表单、弹窗、工程化和部署全部串起来。

这一节不是让你马上写完整系统，而是先建立“后台管理系统应该怎么拆”的思路。

## 你会学到什么

- 后台管理系统的基础模块。
- 登录、权限、菜单、路由的关系。
- 表格、筛选、分页、弹窗表单如何组织。
- 前端项目如何分层。
- 真实项目中高频问题如何提前规避。

如果你想按阶段完成整个后台项目，先看 [Vue Admin 学习地图与交付清单](/roadmap/vue-admin-learning-map)。如果你还不清楚后台项目怎么分层，先看 [图解 Vue Admin 项目架构](/vue/admin-architecture-visual-guide)。如果你已经完成基础页面，下一步优先看 [Vue Admin 权限路由闭环实战](/vue/admin-permission-route-flow)。它会把登录态、菜单、动态路由、按钮权限、接口 403、数据权限和刷新恢复串成一条完整链路。

## 功能模块

| 模块 | 学习重点 | 产出 |
| --- | --- | --- |
| 登录 | 表单校验、token、用户上下文 | 登录页 |
| 布局 | 顶部栏、侧边栏、移动端菜单 | 管理台壳子 |
| 权限 | 动态路由、菜单、按钮权限 | 权限模型 |
| 用户管理 | 筛选、表格、分页、弹窗表单 | CRUD 页面 |
| 角色管理 | 权限树、批量勾选、保存变更 | 授权页面 |
| 请求层 | 拦截器、错误处理、登录失效 | API 基建 |
| 部署 | base、fallback、代理、缓存 | 可上线产物 |

## 推荐目录结构

```text
src/
├─ api/
│  ├─ auth.ts
│  ├─ user.ts
│  └─ role.ts
├─ components/
│  ├─ AppPage.vue
│  ├─ PermissionButton.vue
│  └─ StatusTag.vue
├─ composables/
│  ├─ useLoading.ts
│  ├─ usePagination.ts
│  └─ usePermission.ts
├─ layouts/
│  └─ AdminLayout.vue
├─ router/
│  ├─ index.ts
│  ├─ constantRoutes.ts
│  └─ asyncRoutes.ts
├─ services/
│  └─ auth.ts
├─ stores/
│  ├─ user.ts
│  ├─ permission.ts
│  └─ app.ts
└─ views/
   ├─ login/
   ├─ dashboard/
   └─ system/
      ├─ users/
      └─ roles/
```

## 登录流程

```text
用户输入账号密码
↓
调用登录接口
↓
保存 token
↓
请求当前用户信息
↓
生成权限路由和菜单
↓
跳转到 redirect 或工作台
```

代码结构建议：

```ts
// services/auth.ts
export async function login(payload: LoginPayload) {
  const userStore = useUserStore()
  const permissionStore = usePermissionStore()

  const result = await authApi.login(payload)
  userStore.setToken(result.token)

  await userStore.fetchProfile()
  await permissionStore.generateRoutes(userStore.permissions)
}
```

页面只做表单和调用：

```ts
async function submit() {
  await login(form.value)
  router.replace(route.query.redirect?.toString() || '/dashboard')
}
```

## 用户管理页面

页面职责：

- 保存筛选条件。
- 请求用户列表。
- 打开新增/编辑弹窗。
- 删除用户后刷新列表。

```text
views/system/users/
├─ index.vue
├─ UserSearchForm.vue
├─ UserTable.vue
└─ UserFormDrawer.vue
```

`index.vue` 负责流程：

```vue
<template>
  <UserSearchForm v-model="query" @search="search" />

  <UserTable
    :users="users"
    :loading="loading"
    @create="openCreate"
    @edit="openEdit"
    @remove="removeUser"
  />

  <UserFormDrawer
    v-model:visible="drawerVisible"
    :user="editingUser"
    @success="fetchList"
  />
</template>
```

这样后续换表格组件、换弹窗样式、调整搜索条件时，不会互相污染。

## 角色授权页面

角色授权通常比用户管理复杂，因为它涉及权限树。

常见数据结构：

```ts
interface PermissionNode {
  id: number
  label: string
  code: string
  type: 'menu' | 'button'
  children?: PermissionNode[]
}
```

保存时不要直接提交整棵树，通常只提交选中的权限 id 或 code：

```ts
interface AssignPermissionPayload {
  roleId: number
  permissionCodes: string[]
}
```

## 移动端布局

桌面端可以使用固定侧边栏：

```text
左侧菜单 + 顶部栏 + 内容区
```

移动端不要把整块桌面侧边栏堆到首屏上方。推荐：

- 顶部显示菜单按钮。
- 点击后打开抽屉菜单。
- 首屏优先显示当前页面核心内容。
- 表格在移动端可以切换为卡片列表或横向滚动容器。

## 实际项目常见问题

### 1. 用户管理页面一个文件写到 800 行

**原因**

搜索表单、表格、弹窗、请求、权限按钮、数据转换全部写在同一个文件。

**解决方案**

拆分：

- 搜索表单组件。
- 表格组件。
- 表单抽屉组件。
- 分页 composable。
- 请求 API 模块。

页面保留流程，不保留所有细节。

### 2. 权限按钮到处写 `v-if`

**问题**

权限逻辑散落，权限码改名时很难维护。

**解决方案**

封装：

```vue
<PermissionButton code="system:user:create" @click="openCreate">
  新增用户
</PermissionButton>
```

或者：

```ts
const { can } = usePermission()
```

### 3. 编辑弹窗污染列表数据

**原因**

表单直接绑定了表格行对象。

**解决方案**

打开编辑时复制数据，保存成功后重新请求列表。

### 4. 删除后当前页为空

**场景**

当前在第 5 页，删除最后一条后，列表为空，但其实第 4 页还有数据。

**解决方案**

删除成功后，如果当前页没有数据且页码大于 1，页码减 1 再请求。

```ts
async function handleRemove(id: number) {
  await userApi.removeUser(id)
  await fetchList()

  if (users.value.length === 0 && pagination.page > 1) {
    pagination.page -= 1
    await fetchList()
  }
}
```

### 5. 页面级 loading 和按钮 loading 混乱

**建议**

区分：

- 页面列表 loading：影响表格区域。
- 提交按钮 submitting：影响保存按钮。
- 删除按钮 deletingId：只影响当前行。

不要用一个 `loading` 控制所有交互。

## 验收标准

完成 Vue Admin 第一阶段时，应能做到：

- 登录后能进入工作台。
- 刷新页面后用户信息和菜单能恢复。
- 无权限页面不能进入。
- 无权限按钮不展示或不可用。
- 用户列表支持筛选、分页、新增、编辑、删除。
- 构建后能通过 HTTP 服务预览。
- 部署到非首页路由刷新不 404。

## 下一步

遇到真实问题后，把处理过程沉淀到 [真实项目问题库](/projects/real-world-issues)。
