# 请求与接口封装

## 适合谁看

适合已经会写页面，但接口请求散落在组件里、错误处理重复、登录失效逻辑混乱的学习者。

请求封装的目标不是把 axios 包一层就结束，而是建立清晰的数据访问边界：组件负责交互，API 模块负责请求，service 负责业务流程，store 负责跨页面状态。

## 你会学到什么

- 请求层应该解决哪些问题。
- API、service、store、view 如何分工。
- 如何处理 token、错误提示、登录失效、重复提交。
- 实际项目中接口结构变化、并发请求、401 循环跳转怎么处理。

## 请求层负责什么

常见职责：

- 设置 `baseURL`。
- 自动携带 token。
- 统一处理响应结构。
- 统一处理错误。
- 处理登录失效。
- 处理超时。
- 规范接口函数命名和类型。

不建议请求层直接做的事：

- 直接控制页面弹窗开关。
- 直接修改页面表格数据。
- 直接依赖某个具体页面组件。

## 推荐分层

```text
src/
├─ api/
│  ├─ request.ts       请求实例和拦截器
│  └─ user.ts          用户相关接口
├─ services/
│  └─ auth.ts          登录、登出、初始化用户上下文
├─ stores/
│  └─ user.ts          token、用户信息、权限
└─ views/
   └─ login/index.vue  页面展示和交互
```

分层说明：

| 层 | 负责 | 不负责 |
| --- | --- | --- |
| `api` | 发送请求，描述接口入参和返回值 | 页面交互 |
| `service` | 编排业务流程 | 具体 UI 展示 |
| `store` | 保存跨页面状态 | 所有接口请求 |
| `view` | 表单、按钮、页面状态 | 通用请求规则 |

## 定义响应类型

很多后端接口会包一层结构：

```ts
interface ApiResult<T> {
  code: number
  message: string
  data: T
}

interface PageResult<T> {
  items: T[]
  total: number
}
```

用户列表类型：

```ts
interface User {
  id: number
  username: string
  mobile: string
  enabled: boolean
}

interface UserQuery {
  page: number
  pageSize: number
  keyword?: string
}
```

## API 模块示例

```ts
// api/user.ts
import { request } from './request'

export function getUserList(params: UserQuery) {
  return request.get<PageResult<User>>('/users', { params })
}

export function createUser(payload: CreateUserPayload) {
  return request.post<User>('/users', payload)
}

export function updateUser(id: number, payload: UpdateUserPayload) {
  return request.put<User>(`/users/${id}`, payload)
}

export function removeUser(id: number) {
  return request.delete<void>(`/users/${id}`)
}
```

接口函数命名建议：

| 操作 | 命名 |
| --- | --- |
| 查询列表 | `getUserList` |
| 查询详情 | `getUserDetail` |
| 新增 | `createUser` |
| 更新 | `updateUser` |
| 删除 | `removeUser` |

## 请求拦截器做什么

请求拦截器适合加 token、trace id、租户 id 等统一请求头。

```ts
request.interceptors.request.use((config) => {
  const userStore = useUserStore()

  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }

  return config
})
```

注意：如果在组件外使用 Pinia，需要确保 Pinia 已经安装到应用。Pinia 官方文档也提醒，组件外使用 store 时要注意 pinia 实例的注入时机。

## 响应拦截器做什么

响应拦截器适合统一拆包和处理错误：

```ts
request.interceptors.response.use(
  (response) => {
    const result = response.data as ApiResult<unknown>

    if (result.code !== 0) {
      throw new Error(result.message || '请求失败')
    }

    return result.data
  },
  (error) => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      router.replace('/login')
    }

    return Promise.reject(error)
  }
)
```

真实项目里，是否在拦截器里弹错误提示要谨慎。如果所有错误都自动弹，会导致表单校验、静默刷新、轮询接口也弹出干扰信息。

## 页面如何使用

```ts
const users = ref<User[]>([])
const loading = ref(false)

async function fetchUsers() {
  loading.value = true

  try {
    const result = await getUserList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: keyword.value
    })

    users.value = result.items
    pagination.total = result.total
  } finally {
    loading.value = false
  }
}
```

页面只关心：

- 什么时候请求。
- 加载中怎么展示。
- 请求成功后更新哪个页面状态。
- 请求失败后是否需要页面级处理。

## 实际项目常见问题

### 1. 接口请求写在每个组件里，重复很多

**症状**

每个页面都重复写 baseURL、token、错误处理、loading。

**解决方案**

建立统一 `request.ts`，页面只调用 API 函数。

### 2. 401 后无限跳登录

**症状**

登录失效后页面不断跳转，或者多个请求同时失败导致重复弹窗。

**原因**

多个并发请求都收到了 401，每个请求都执行一次登出和跳转。

**解决方案**

增加一次性处理标记：

```ts
let isHandlingUnauthorized = false

function handleUnauthorized() {
  if (isHandlingUnauthorized) return

  isHandlingUnauthorized = true
  const userStore = useUserStore()
  userStore.logout()
  router.replace('/login').finally(() => {
    isHandlingUnauthorized = false
  })
}
```

### 3. 重复点击提交按钮创建了多条数据

**症状**

用户连续点“保存”，后端创建了重复记录。

**解决方案**

前端禁用按钮，后端也要做幂等或唯一约束。

```ts
const submitting = ref(false)

async function submit() {
  if (submitting.value) return

  submitting.value = true
  try {
    await createUser(form.value)
  } finally {
    submitting.value = false
  }
}
```

### 4. 后端返回字段变化导致页面大量报错

**症状**

后端把 `user_name` 改成 `username`，多个页面同时坏掉。

**解决方案**

在 API 或 service 层做数据适配，不要让页面直接依赖后端原始结构。

```ts
interface RawUser {
  user_name: string
  phone_no: string
}

function normalizeUser(raw: RawUser): User {
  return {
    username: raw.user_name,
    mobile: raw.phone_no
  }
}
```

### 5. 搜索接口返回顺序错乱

**症状**

用户连续输入关键字，旧请求比新请求晚返回，页面显示了旧结果。

**解决方案**

记录请求序号，只接受最后一次请求结果：

```ts
let requestId = 0

async function fetchList() {
  const currentId = ++requestId
  const result = await getUserList(query.value)

  if (currentId !== requestId) return

  users.value = result.items
}
```

## 最佳实践

- 组件里不要直接拼接大量 URL。
- API 函数必须有清晰入参和返回类型。
- 请求层处理通用规则，页面层处理具体交互。
- 登录失效处理要防重复。
- 提交按钮要防重复点击。
- 后端数据结构变化时，优先在 service 层适配。

## 下一步学习

继续学习 [权限与菜单](/vue/permission)，把请求拿到的用户权限转换成路由、菜单和按钮控制。
