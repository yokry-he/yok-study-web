# 常见问题

## 使用方式

这一页按“症状、常见原因、解决方案、预防建议”的格式整理。遇到问题时，先找到最像的症状，再按排查顺序处理。

如果问题已经进入项目链路，例如动态菜单刷新丢失、Pinia 解构后页面不更新、编辑弹窗污染列表、KeepAlive 缓存导致数据不刷新、权限按钮错位，优先看 [Vue 真实项目问题库](/projects/issues-vue)。

## 1. 页面刷新后 404

### 症状

- 从菜单点击 `/system/users` 可以进入。
- 在 `/system/users` 刷新浏览器后变成 404。
- 直接复制链接给别人，别人打开也是 404。

### 常见原因

Vue Router 使用 history 模式时，页面路由由前端接管。刷新浏览器时，服务器会收到 `/system/users` 请求。如果服务器没有配置回退到 `index.html`，就会返回 404。

### 解决方案

Nginx：

```nginx
location / {
  try_files $uri $uri/ /index.html;
}
```

如果部署在 `/admin/` 子路径：

```ts
// vite.config.ts
export default defineConfig({
  base: '/admin/'
})
```

```ts
// router/index.ts
createRouter({
  history: createWebHistory('/admin/'),
  routes
})
```

### 预防建议

每次上线前都测试：

- 首页刷新。
- 二级路由刷新。
- 深层详情页刷新。
- 复制链接到新标签页打开。

## 2. 数据已经变了，但页面不更新

### 症状

- 控制台打印数据已变化。
- 页面显示还是旧值。

### 常见原因

- 解构 `reactive` 后丢失响应式。
- 子组件复制了 props，但没有监听 props 变化。
- 修改了非响应式普通变量。
- 使用了错误的数组或对象引用。

### 解决方案

`reactive` 解构使用 `toRefs`：

```ts
const state = reactive({ name: 'alice' })
const { name } = toRefs(state)
```

Pinia 解构使用 `storeToRefs`：

```ts
const userStore = useUserStore()
const { profile } = storeToRefs(userStore)
```

编辑表单复制 props 时，如果 props 会变化，需要监听：

```ts
watch(
  () => props.user,
  (user) => {
    form.value = user ? createFormByUser(user) : defaultForm()
  },
  { immediate: true }
)
```

## 3. watch 不触发

### 症状

状态变化了，但 `watch` 回调没有执行。

### 常见原因

监听了普通值，而不是响应式来源。

错误示例：

```ts
const state = reactive({ count: 0 })

watch(state.count, () => {
  // 不会按预期工作
})
```

### 解决方案

监听 getter：

```ts
watch(
  () => state.count,
  () => {
    console.log('count changed')
  }
)
```

如果监听多个值：

```ts
watch(
  [() => query.keyword, () => query.status],
  () => {
    fetchList()
  }
)
```

## 4. 路由参数变化但页面不刷新

### 症状

从 `/users/1` 跳到 `/users/2`，URL 变了，但页面内容仍然是用户 1。

### 原因

Vue Router 可能复用同一个组件实例，`onMounted` 不会再次执行。

### 解决方案

监听路由参数：

```ts
const route = useRoute()

watch(
  () => route.params.id,
  (id) => {
    fetchUserDetail(String(id))
  },
  { immediate: true }
)
```

## 5. 登录后又跳回登录页

### 症状

登录接口成功，但页面马上回到登录页。

### 排查顺序

1. token 是否保存成功。
2. 请求用户信息接口是否成功。
3. 请求拦截器是否带了 token。
4. 路由守卫是否等待用户信息请求完成。
5. 登录页是否错误设置了 `requiresAuth`。

### 常见解决方案

```ts
if (userStore.token && !userStore.profile) {
  await userStore.fetchProfile()
}
```

如果用户信息请求失败，要清理 token 并给出明确提示。

## 6. 动态菜单刷新后丢失

### 症状

登录后菜单正常，刷新页面后只剩首页或空白菜单。

### 原因

菜单和动态路由只存在内存中，刷新后没有重新生成。

### 解决方案

在路由守卫中恢复用户上下文：

```ts
if (userStore.token && !permissionStore.ready) {
  await initUserContext()
  return to.fullPath
}
```

返回 `to.fullPath` 是为了让刚注册的动态路由重新匹配当前地址。

## 7. 接口请求重复发送

### 症状

打开页面时列表接口请求了两次。

### 常见原因

- `onMounted(fetchList)` 调用一次。
- `watch(query, fetchList, { immediate: true })` 又调用一次。
- 父组件和子组件都请求了同一份数据。

### 解决方案

确定唯一的数据入口。列表页通常由页面组件统一请求，子组件通过 props 接收数据。

## 8. 搜索结果偶尔显示旧数据

### 症状

快速输入关键字，页面最终显示的不是最后一次输入的结果。

### 原因

多个请求并发，旧请求比新请求更晚返回。

### 解决方案

记录请求序号：

```ts
let requestId = 0

async function fetchList() {
  const currentId = ++requestId
  const result = await getList(query.value)

  if (currentId !== requestId) return

  list.value = result.items
}
```

## 9. 表单关闭后再次打开仍然是旧数据

### 症状

编辑用户 A 后关闭弹窗，再新增用户时表单里还残留用户 A 的数据。

### 原因

关闭弹窗时没有重置表单，或者默认值对象被复用。

### 解决方案

```ts
function createDefaultForm() {
  return {
    username: '',
    mobile: '',
    enabled: true
  }
}

const form = ref(createDefaultForm())

function resetForm() {
  form.value = createDefaultForm()
}
```

## 10. 组件样式异常

### 症状

- 组件库按钮尺寸突然变小。
- Switch、Table、Drawer 样式变形。
- 表格操作列被压缩。
- 移动端出现横向滚动。

### 排查顺序

1. 搜索是否有宽泛选择器，例如 `.page button`、`.content div`、`.xxx *`。
2. 检查是否覆盖了组件库内部 class。
3. 检查 flex 子项是否缺少 `flex-shrink: 0`。
4. 检查表格、工具栏、按钮组是否设置了稳定宽度。
5. 强制刷新，排除 HMR 或旧进程缓存。

### 解决方案

业务样式命中明确 class：

```css
.permission-switch-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-action {
  flex-shrink: 0;
}
```

如果确实需要改组件库样式，优先使用主题 token、组件 props、CSS 变量或官方暴露 API。

## 11. 构建成功但打开 dist 空白

### 常见原因

- 直接用 `file://` 打开构建产物。
- 部署子路径和 Vite `base` 不一致。
- 静态资源被服务器错误缓存。
- 环境变量缺失。

Vite 官方排错文档也说明，构建产物不能直接通过 `file` 协议正常运行，应通过 HTTP 服务访问。

### 解决方案

本地预览：

```bash
npm run build
npm run preview
```

部署到子路径时配置：

```ts
export default defineConfig({
  base: '/admin/'
})
```

## 12. 删除按钮隐藏了但接口仍能删除

### 原因

前端权限只控制展示，不是安全边界。

### 解决方案

后端必须校验接口权限。前端隐藏按钮只是减少无效操作。

## 快速排查口诀

| 问题类型 | 先看哪里 |
| --- | --- |
| 页面空白 | 浏览器 Console 第一条红错 |
| 接口失败 | Network 请求、状态码、响应体 |
| 权限异常 | token、profile、permissions、路由 meta |
| 页面不更新 | ref/reactive、toRefs、storeToRefs |
| 刷新 404 | 服务器 fallback 和 Vite base |
| 样式异常 | 宽泛选择器和组件库覆盖 |

## 参考工具

- Vue DevTools：查看组件树、props、state、Pinia。
- 浏览器 Network：检查请求、响应、请求头。
- 浏览器 Performance：分析卡顿和重复渲染。
- 终端构建日志：定位类型错误、导入路径错误、构建配置问题。
