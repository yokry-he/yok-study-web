# React 表单处理

## 适合谁看

适合准备在 React 项目中开发登录、搜索、新增编辑弹窗和复杂业务表单的学习者。

React 表单的核心是“受控组件”：输入框的值来自 state，用户输入后通过事件更新 state。这样表单数据始终在 React 中可控。

## 你会学到什么

- 受控输入怎么写。
- 表单默认值、编辑回显和重置怎么处理。
- 提交前如何校验和转换。
- 如何防重复提交。
- 实际项目中输入不了、旧数据残留、数字变字符串怎么解决。

## 受控输入

```tsx
function UserSearchForm() {
  const [keyword, setKeyword] = useState('')

  return (
    <input
      value={keyword}
      onChange={(event) => setKeyword(event.target.value)}
      placeholder="搜索用户名"
    />
  )
}
```

写了 `value` 就必须写 `onChange`，否则输入框会变成无法编辑。

## 表单对象

```tsx
interface UserForm {
  id?: number
  username: string
  mobile: string
  enabled: boolean
}

function createDefaultForm(): UserForm {
  return {
    username: '',
    mobile: '',
    enabled: true
  }
}
```

状态：

```tsx
const [form, setForm] = useState<UserForm>(() => createDefaultForm())
```

更新字段：

```tsx
function updateField<K extends keyof UserForm>(key: K, value: UserForm[K]) {
  setForm((current) => ({
    ...current,
    [key]: value
  }))
}
```

## 新增和编辑

新增：

```tsx
function openCreate() {
  setForm(createDefaultForm())
  setVisible(true)
}
```

编辑：

```tsx
function openEdit(user: User) {
  setForm({
    id: user.id,
    username: user.username,
    mobile: user.mobile,
    enabled: user.enabled
  })
  setVisible(true)
}
```

不要直接把列表行对象作为表单状态修改。

## 提交校验

```tsx
function validateForm(form: UserForm) {
  if (!form.username.trim()) {
    return '请输入用户名'
  }

  if (!/^1\d{10}$/.test(form.mobile)) {
    return '请输入正确手机号'
  }

  return ''
}
```

提交：

```tsx
async function submit() {
  const message = validateForm(form)
  if (message) {
    showError(message)
    return
  }

  await saveUser(form)
}
```

## 防重复提交

```tsx
const [submitting, setSubmitting] = useState(false)

async function submit() {
  if (submitting) return

  setSubmitting(true)
  try {
    await saveUser(form)
    setVisible(false)
  } finally {
    setSubmitting(false)
  }
}
```

## 实际项目常见问题

### 1. 输入框无法输入

**原因**

写了 `value` 但没有正确更新 state。

**解决方案**

```tsx
<input value={form.username} onChange={(event) => updateField('username', event.target.value)} />
```

### 2. 新增表单残留编辑数据

**原因**

打开新增时没有重置默认值。

**解决方案**

`openCreate()` 中调用 `createDefaultForm()`。

### 3. 数字输入变成字符串

HTML input 的值默认是字符串。提交前转换：

```tsx
const payload = {
  ...form,
  age: Number(form.age)
}
```

或者使用成熟组件库的数字输入组件。

## 最佳实践

- 表单类型和接口类型分开。
- 默认值用函数创建。
- 新增和编辑入口都明确初始化表单。
- 提交前校验、转换、防重复。
- 复杂表单优先使用成熟表单库或组件库表单系统。

## 下一步

继续学习 [请求与数据流](/react/request-data-flow)。
