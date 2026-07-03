# 前端页面与状态问题

## 适合谁看

这篇适合已经能写 Vue 或 React 页面，但在真实项目里经常遇到下面问题的同学：

- 页面偶尔显示旧数据。
- 弹窗、表单、列表互相影响。
- 路由、菜单、权限状态刷新后丢失。
- 组件库样式或响应式布局在某些页面突然变形。
- 页面性能不好，但不知道该从哪里排查。

前端问题不要只看“代码有没有报错”。很多真实问题来自状态边界、请求并发、组件复用、布局约束和浏览器运行环境。

## 使用方式

每个问题都按固定结构阅读：

```text
问题现象
影响范围
常见根因
解决方案
预防方式
```

遇到线上问题时，先对照“问题现象”；修复后再看“预防方式”，把规则沉淀到组件、工具函数或项目规范中。

## 问题 1：搜索框输入很快时，列表显示的不是最后一次结果

### 问题现象

- 用户输入 `a`、`ab`、`abc`，列表最后显示的是 `ab` 的结果。
- Network 里能看到多个请求都成功了。
- 后端没有报错，本地慢网速更容易复现。

### 影响范围

所有“会被快速连续触发”的请求：

- 搜索列表。
- 远程下拉选择器。
- 自动补全。
- 切换筛选条件后自动刷新表格。
- 输入地址后自动计算运费、地区、坐标。

### 常见根因

请求是异步的，不保证先发出的请求先返回。如果旧请求比新请求更晚返回，就会覆盖新请求的数据。

错误写法通常是：

```ts
async function search() {
  loading.value = true
  const res = await api.search(keyword.value)
  list.value = res.items
  loading.value = false
}
```

这段代码只关心“请求是否成功”，没有关心“这个请求是不是当前最新请求”。

### 解决方案

方案一：用请求序号保护最终写入。

```ts
let requestSeq = 0

async function search() {
  const currentSeq = ++requestSeq
  loading.value = true

  try {
    const res = await api.search(keyword.value)

    if (currentSeq !== requestSeq) return

    list.value = res.items
  } finally {
    if (currentSeq === requestSeq) {
      loading.value = false
    }
  }
}
```

方案二：如果请求库支持取消请求，切换条件时取消上一次请求。

```ts
let controller: AbortController | null = null

async function search() {
  controller?.abort()
  controller = new AbortController()

  const res = await fetch(`/api/search?q=${keyword.value}`, {
    signal: controller.signal
  })

  list.value = await res.json()
}
```

### 预防方式

- 对搜索、筛选、联想类请求统一封装 `latestOnly` 工具。
- 输入类搜索默认加防抖，例如 300ms。
- 组件卸载时取消未完成请求。
- 所有请求状态写入都要确认“当前请求仍然有效”。

## 问题 2：编辑弹窗修改字段后，背后的表格行也变了

### 问题现象

- 点击编辑，弹出表单。
- 用户还没有点保存，只是在输入框里改了字段。
- 表格里的同一行数据已经跟着变化。
- 关闭弹窗后，页面显示被改过，但刷新后又恢复。

### 影响范围

常见于后台管理系统：

- 用户编辑。
- 角色编辑。
- 商品编辑。
- 订单备注编辑。
- 配置项编辑。

### 常见根因

弹窗表单直接引用了列表行对象。

```ts
function openEdit(row: User) {
  form.value = row
}
```

`row` 是对象，`form.value` 指向同一个引用。表单改字段，本质上就是改列表里的对象。

### 解决方案

打开弹窗时复制一份干净表单。

```ts
function openEdit(row: User) {
  form.value = {
    id: row.id,
    username: row.username,
    mobile: row.mobile,
    enabled: row.enabled,
    roleIds: [...row.roleIds]
  }

  visible.value = true
}
```

复杂对象可以写一个转换函数，不要在页面里到处手写复制逻辑。

```ts
function toUserForm(row: User): UserForm {
  return {
    id: row.id,
    username: row.username,
    mobile: row.mobile,
    enabled: row.enabled,
    roleIds: [...row.roleIds]
  }
}
```

保存成功后再刷新列表，或者只替换列表中的当前行。

```ts
async function submit() {
  await api.updateUser(form.value)
  visible.value = false
  await fetchList()
}
```

### 预防方式

- 表单状态和列表状态分开管理。
- 弹窗表单不要直接修改 props。
- 列表行对象不要直接作为可编辑状态。
- 给每个复杂表单建立 `toForm` 和 `toPayload` 转换函数。

## 问题 3：刷新页面后菜单、按钮权限或动态路由丢失

### 问题现象

- 登录后页面正常。
- 刷新浏览器后菜单为空。
- 直接访问某个业务页面变成 404。
- 按钮权限全部隐藏，或者所有按钮都显示。

### 影响范围

常见于使用动态路由、动态菜单和按钮权限的后台系统。

### 常见根因

登录后只把用户信息、菜单和动态路由保存在内存里。刷新后 Pinia、Redux 或普通变量都会重新初始化。

很多项目只在登录成功后执行：

```ts
await userStore.fetchProfile()
await permissionStore.generateRoutes()
router.addRoute(...)
```

但没有在应用启动或路由守卫中恢复这些状态。

### 解决方案

把“恢复用户上下文”做成可重复调用的启动流程。

```ts
let isBootstrapping = false

router.beforeEach(async (to) => {
  if (!authStore.token) {
    return to.path === '/login' ? true : '/login'
  }

  if (!permissionStore.ready && !isBootstrapping) {
    isBootstrapping = true

    try {
      await authStore.fetchProfile()
      await permissionStore.fetchMenus()
      permissionStore.routes.forEach((route) => router.addRoute(route))
      return to.fullPath
    } finally {
      isBootstrapping = false
    }
  }

  return true
})
```

如果项目有按钮权限，权限码也应跟随用户上下文一起恢复。

```ts
const canCreateUser = computed(() => {
  return permissionStore.has('system:user:create')
})
```

### 预防方式

- 区分“登录动作”和“应用启动恢复动作”。
- 用户信息、菜单、按钮权限、动态路由要有统一初始化入口。
- 路由守卫中避免重复初始化，要有 `ready` 标记。
- 权限恢复失败时进入明确的错误页，而不是静默白屏。

## 问题 4：组件库样式在某个页面突然失控

### 问题现象

- 表格行高异常。
- Switch、按钮、输入框尺寸变形。
- 弹窗内容挤在一起。
- 同一个组件在其他页面正常，只在当前页面异常。

### 影响范围

所有使用组件库的项目都可能遇到，尤其是后台系统和低代码配置页。

### 常见根因

业务 CSS 使用了宽泛选择器，污染了组件库内部 DOM。

```css
.page button {
  height: 28px;
}

.panel div {
  display: flex;
}

.content * {
  box-sizing: border-box;
}
```

组件库内部也有 `button`、`div`、`span`。这些样式命中后，组件就会出现不可预期的变形。

### 解决方案

业务样式只命中明确业务 class。

```css
.user-page__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.user-page__toolbar-action {
  flex: 0 0 auto;
}

.permission-switch-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
```

固定尺寸元素要避免被压缩。

```css
.user-avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}
```

确实需要改组件库样式时，优先使用：

- 组件 props。
- 主题 token。
- CSS 变量。
- 官方暴露的 class 或 API。

### 预防方式

- 禁止业务代码写 `.xxx div`、`.xxx button`、`.xxx *` 这类选择器。
- 提交前搜索宽泛选择器。
- 样式异常时先查全局 CSS 和页面 CSS，不要直接叠更高优先级。
- 每次改全局样式后检查表格、弹窗、表单、开关、按钮。

## 问题 5：列表页面越来越慢，滚动和输入都有明显卡顿

### 问题现象

- 表格超过几百行后页面明显卡顿。
- 输入搜索关键字时每输入一个字符都会卡。
- 切换 tab 或筛选条件时浏览器短暂无响应。

### 影响范围

后台列表、数据看板、权限树、组织架构树、商品规格表格等数据密集页面。

### 常见根因

常见原因不是单点问题，而是多个小问题叠加：

- 一次性渲染太多 DOM。
- 每行都创建复杂计算属性或函数。
- 列表 key 不稳定。
- 筛选、排序、聚合每次输入都全量计算。
- 大对象被深度响应式代理。

### 解决方案

先确认瓶颈在“请求慢”还是“渲染慢”。如果接口已经返回，但页面很久才更新，多半是渲染或计算问题。

减少一次性渲染数量：

```ts
const pageSize = ref(20)
const page = ref(1)

const visibleRows = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return rows.value.slice(start, start + pageSize.value)
})
```

避免模板里做复杂计算：

```vue
<!-- 不推荐：每次渲染都执行复杂函数 -->
<span>{{ formatUserStatus(row.status, row.expiredAt, row.locked) }}</span>
```

改成提前转换展示模型：

```ts
const viewRows = computed(() => {
  return rows.value.map((row) => ({
    ...row,
    statusText: getStatusText(row),
    statusType: getStatusType(row)
  }))
})
```

超大列表使用虚拟滚动或后端分页，不要把几万条数据交给浏览器一次性渲染。

### 预防方式

- 列表默认后端分页。
- 大表格默认限制列数和行数。
- 模板只做展示，不做重计算。
- 大对象如果不需要深度响应式，可考虑 `shallowRef`。
- 性能问题先用浏览器 Performance 面板定位，再改代码。

## 问题 6：React 页面重复请求，甚至出现请求死循环

### 问题现象

- 页面打开后接口请求不止一次。
- 切换筛选条件时请求次数明显超过预期。
- Network 里能看到相同 URL 被连续调用。
- 严重时页面一直 loading，浏览器和后端压力都升高。

### 影响范围

React 列表页、详情页、远程下拉、仪表盘、权限菜单加载、文档问答页面。

### 常见根因

最常见原因是 `useEffect` 依赖不稳定，或者把“事件触发逻辑”错误放进 Effect。

错误示例：

```tsx
function UsersPage() {
  const query = { page: 1, keyword }

  useEffect(() => {
    fetchUsers(query)
  }, [query])
}
```

`query` 每次渲染都会创建新对象，React 会认为依赖变化，于是再次执行 Effect。

还有一种常见错误：

```tsx
useEffect(() => {
  setQuery({ ...query, page: 1 })
}, [query])
```

Effect 修改了自己的依赖，容易形成循环。

### 解决方案

方案一：把依赖拆成稳定的原始值。

```tsx
useEffect(() => {
  fetchUsers({ page, pageSize, keyword })
}, [page, pageSize, keyword])
```

方案二：用户点击搜索时再触发请求，而不是所有输入变化都立即请求。

```tsx
function handleSearch() {
  setPage(1)
  fetchUsers({ page: 1, pageSize, keyword })
}
```

方案三：列表条件同步到 URL query，页面根据 URL 变化请求。

```tsx
const [searchParams, setSearchParams] = useSearchParams()

useEffect(() => {
  fetchUsers({
    page: Number(searchParams.get('page') ?? 1),
    keyword: searchParams.get('keyword') ?? ''
  })
}, [searchParams])
```

如果使用 React Query、SWR 等数据请求库，要让 query key 稳定且表达真实依赖。

### 预防方式

- Effect 只用于同步外部系统，不用于普通数据计算。
- 依赖数组里避免放每次渲染都新建的对象和函数。
- 请求触发时机要明确：进入页面、URL 变化、点击搜索、切换分页。
- 用 React DevTools 和 Network 一起定位重复渲染和重复请求。
- 对核心列表页做一次“打开页面只请求几次”的冒烟检查。

## 问题 7：组件库升级后，业务页面大面积错位

### 问题现象

- 升级组件库后，表格行高、按钮尺寸、弹窗宽度、表单间距都变化。
- 有些页面正常，有些页面错位。
- 回退组件库版本后恢复。
- 业务代码没有明显改动。

### 影响范围

自研组件库、二次封装组件库、Element Plus、Ant Design Vue、Arco Design、Naive UI、TDesign 等组件库项目。

### 常见根因

组件库升级可能改变：

- 默认 token。
- 组件尺寸。
- DOM 结构。
- class 命名。
- 默认插槽结构。
- 弹层挂载位置。
- CSS 变量名称。

业务项目如果依赖了组件库内部 DOM，升级就很容易出问题。

```css
/* 不推荐：依赖内部结构 */
.user-page .n-button__content span {
  font-size: 12px;
}
```

### 解决方案

先看变更来源，不要直接强行覆盖。

排查顺序：

1. 看组件库 changelog。
2. 对比 lockfile 和组件库版本。
3. 用 DevTools 查看异常元素的 Computed 样式。
4. 检查业务 CSS 是否命中组件库内部 class。
5. 检查主题 token 或 ConfigProvider 是否变化。
6. 在关键页面做桌面和移动端回归。

如果需要调整组件样式，优先顺序是：

```text
组件 props
↓
组件库主题 token
↓
官方 CSS 变量
↓
项目封装组件
↓
局部、明确、可说明原因的样式覆盖
```

### 预防方式

- 组件库升级必须单独提交，避免混在业务需求里。
- 对按钮、表格、弹窗、表单、下拉、开关建立示例页。
- 禁止业务页面依赖组件库内部 DOM 结构写样式。
- 固定尺寸元素设置稳定宽高和 `flex-shrink: 0`。
- 升级前后跑关键页面截图或人工回归。

## 下一步学习

- [Vue 性能优化](/vue/performance)
- [Vue 常见问题](/vue/troubleshooting)
- [React 管理台从零到项目](/react/project-admin)
- [组件库工程从零到项目](/engineering/component-library-project)
- [浏览器渲染与性能](/browser/rendering-performance)
- [CSS 项目样式架构](/css/architecture)
