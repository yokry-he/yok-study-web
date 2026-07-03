# 表单处理

## 适合谁看

适合准备开发登录页、搜索表单、新增编辑弹窗、复杂业务表单的学习者。

表单是后台管理系统中最常见、也最容易出错的部分。它涉及默认值、校验、编辑回显、提交转换、重复提交、错误提示和重置。

## 你会学到什么

- 表单状态怎么定义。
- 新增和编辑如何共用表单。
- 表单默认值为什么要用函数。
- 提交前如何校验和转换。
- 实际项目中表单残留、污染列表、重复提交怎么解决。

## 表单类型

```ts
interface UserForm {
  id?: number
  username: string
  mobile: string
  enabled: boolean
  roleIds: number[]
}
```

不要直接把后端完整用户对象当表单类型。表单只包含用户可以输入或修改的字段。

## 默认值

推荐用函数创建默认值：

```ts
function createDefaultForm(): UserForm {
  return {
    username: '',
    mobile: '',
    enabled: true,
    roleIds: []
  }
}

const form = ref<UserForm>(createDefaultForm())
```

好处是每次重置都会得到一个全新对象，避免引用污染。

## 新增和编辑

```ts
const visible = ref(false)
const form = ref<UserForm>(createDefaultForm())

function openCreate() {
  form.value = createDefaultForm()
  visible.value = true
}

function openEdit(user: User) {
  form.value = {
    id: user.id,
    username: user.username,
    mobile: user.mobile,
    enabled: user.enabled,
    roleIds: user.roles.map((role) => role.id)
  }
  visible.value = true
}
```

编辑时不要直接绑定列表行对象。

## v-model 绑定

```vue
<input v-model="form.username" placeholder="请输入用户名" />
<input v-model="form.mobile" placeholder="请输入手机号" />
<input v-model="form.enabled" type="checkbox" />
```

组件库项目中，优先使用组件库表单组件，例如 Naive UI、Element Plus、Ant Design Vue 等已经提供的输入框、选择器、校验能力。

## 基础校验

提交前做必要校验：

```ts
function validateForm() {
  if (!form.value.username.trim()) {
    return '请输入用户名'
  }

  if (!/^1\d{10}$/.test(form.value.mobile)) {
    return '请输入正确手机号'
  }

  return ''
}
```

提交：

```ts
async function submit() {
  const message = validateForm()
  if (message) {
    showError(message)
    return
  }

  await saveUser(form.value)
}
```

实际项目应优先使用组件库表单校验系统，避免手写重复校验。

## 提交转换

页面表单结构不一定等于接口结构：

```ts
function toCreatePayload(form: UserForm): CreateUserPayload {
  return {
    username: form.username.trim(),
    mobile: form.mobile,
    enabled: form.enabled,
    roleIds: form.roleIds
  }
}
```

把转换逻辑抽出来，便于测试和维护。

## 防重复提交

```ts
const submitting = ref(false)

async function submit() {
  if (submitting.value) return

  submitting.value = true
  try {
    await saveUser(toCreatePayload(form.value))
    visible.value = false
    emit('success')
  } finally {
    submitting.value = false
  }
}
```

前端防重复只是用户体验，后端仍然应该有幂等、唯一约束或业务校验。

## 实际项目常见问题

### 1. 新增表单残留上次编辑的数据

**原因**

新增时没有重置表单。

**解决方案**

`openCreate()` 中重新赋值 `createDefaultForm()`。

### 2. 编辑表单还没保存，列表已经变化

**原因**

表单直接引用了列表行对象。

**解决方案**

打开编辑时复制字段，保存成功后刷新列表。

### 3. 关闭弹窗后校验错误还在

**原因**

组件库表单校验状态没有清理。

**解决方案**

关闭时重置表单数据和校验状态。不同组件库 API 不同，通常会提供 `resetFields` 或 `clearValidate`。

### 4. 数字输入变成字符串

**原因**

HTML input 默认输入字符串。

**解决方案**

提交前转换：

```ts
const payload = {
  ...form.value,
  age: Number(form.value.age)
}
```

或使用组件库数字输入组件。

## 最佳实践

- 表单类型和接口类型分开。
- 默认表单用函数创建。
- 新增和编辑入口都明确初始化表单。
- 提交前校验、转换、防重复。
- 编辑时复制数据，不直接修改 props 或列表行。
- 组件库表单校验状态要在关闭时清理。

## 下一步学习

继续学习 [请求与接口封装](/vue/request)。
