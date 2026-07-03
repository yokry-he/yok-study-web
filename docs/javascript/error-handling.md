# 错误处理

## 适合谁看

适合已经会写接口请求和业务函数，但经常遇到这些问题的人：

- 请求失败后页面没有提示。
- loading 一直转圈。
- `try...catch` 包了代码，但异步错误还是没捕获到。
- 后端返回 400、401、500，不知道前端该怎么分层处理。
- 线上用户报错，但控制台或日志里没有足够信息。

错误处理不是在每段代码外面套一层 `try...catch`。真正的目标是让失败可预期、可提示、可恢复、可定位。

## 错误分几类

项目里建议先按来源分类：

| 类型 | 示例 | 处理重点 |
| --- | --- | --- |
| 编程错误 | 空对象取属性、变量未定义 | 开发期修复，线上记录 |
| 业务错误 | 表单校验失败、余额不足 | 给用户明确提示 |
| 网络错误 | 断网、超时、DNS 失败 | 提示重试或降级 |
| HTTP 错误 | 401、403、404、500 | 按状态码分层处理 |
| 数据错误 | 接口字段缺失、类型不符合预期 | 规范化和兜底 |

不要把所有错误都显示成“系统异常”。用户需要知道能不能重试、要不要重新登录、是否需要修改输入。

## try...catch

基础写法：

```ts
try {
  riskyTask()
} catch (error) {
  console.error(error)
} finally {
  cleanup()
}
```

`finally` 无论成功还是失败都会执行，适合恢复 loading、释放锁、关闭临时状态。

```ts
async function submit() {
  submitting.value = true

  try {
    await saveForm()
    showSuccess('保存成功')
  } catch (error) {
    showError(getErrorMessage(error))
  } finally {
    submitting.value = false
  }
}
```

## try...catch 捕获不到什么

同步 `try...catch` 捕获不到另一个任务里的错误。

```ts
try {
  setTimeout(() => {
    throw new Error('boom')
  }, 0)
} catch (error) {
  console.log('不会执行')
}
```

因为定时器回调不是在当前这段同步代码里执行。

正确做法是在异步边界内部处理：

```ts
setTimeout(() => {
  try {
    riskyTask()
  } catch (error) {
    reportError(error)
  }
}, 0)
```

Promise 错误要用 `await + try...catch` 或 `.catch()`：

```ts
try {
  await fetchUser()
} catch (error) {
  showError(getErrorMessage(error))
}
```

## fetch 的特殊点

`fetch` 只有网络层失败时才会 reject。HTTP 400 或 500 默认仍然是成功返回 response。

```ts
const response = await fetch('/api/users')

if (!response.ok) {
  throw new Error(`HTTP ${response.status}`)
}

const data = await response.json()
```

如果不检查 `response.ok`，接口 500 也可能继续进入后续逻辑。

## 请求错误的分层

推荐分层：

```text
request 基础层
↓
api 模块
↓
service / composable
↓
页面组件
```

每层职责不同：

| 层 | 职责 |
| --- | --- |
| request | 处理 HTTP 状态、解析响应、统一错误对象 |
| api | 声明接口函数，不写 UI 逻辑 |
| service/composable | 组织业务流程、loading、重试 |
| 页面组件 | 展示状态和用户提示 |

不要在 request 层直接写大量页面跳转和弹窗逻辑，否则会导致全局副作用难以控制。

## 设计统一错误对象

建议把不同来源的错误整理成统一结构：

```ts
type AppError = {
  code: string
  message: string
  status?: number
  cause?: unknown
  recoverable: boolean
}
```

示例：

```ts
function normalizeError(error: unknown): AppError {
  if (error instanceof Error) {
    return {
      code: 'UNKNOWN_ERROR',
      message: error.message,
      cause: error,
      recoverable: false
    }
  }

  return {
    code: 'UNKNOWN_ERROR',
    message: '未知错误',
    cause: error,
    recoverable: false
  }
}
```

页面只关心怎么展示，日志系统保留原始错误。

## 401、403、500 怎么处理

| 状态 | 含义 | 常见处理 |
| --- | --- | --- |
| 400 | 请求参数错误 | 展示字段或业务提示 |
| 401 | 未登录或登录过期 | 清理登录态，引导重新登录 |
| 403 | 已登录但无权限 | 展示无权限，不要重复跳登录 |
| 404 | 资源不存在 | 显示空状态或不存在页面 |
| 409 | 数据冲突 | 提示刷新或重新提交 |
| 500 | 服务端异常 | 提示稍后重试，记录日志 |

认证错误要避免死循环。登录接口、刷新 token 接口、退出接口通常需要白名单。

## 不要吞掉错误

危险写法：

```ts
try {
  await save()
} catch (error) {
  // 什么都不做
}
```

这样用户不知道发生了什么，开发也没有证据定位。

至少要做一件事：

- 显示用户提示。
- 记录日志。
- 返回明确失败状态。
- 重新抛出给上层处理。

## 重新抛出错误

底层函数不知道怎么展示错误时，可以补充上下文后重新抛出。

```ts
async function loadUserProfile(userId: string) {
  try {
    return await userApi.getProfile(userId)
  } catch (error) {
    throw new Error(`加载用户资料失败：${userId}`, {
      cause: error
    })
  }
}
```

上层页面决定怎么提示，日志里还能看到原始原因。

## 全局错误兜底

前端项目一般需要几类兜底：

- `window.onerror`：未捕获同步错误。
- `window.onunhandledrejection`：未处理 Promise rejection。
- 框架错误边界：Vue `app.config.errorHandler`、React Error Boundary。
- 请求拦截器：HTTP 和业务错误规范化。
- 日志上报：带上路由、用户、版本、环境和关键上下文。

兜底不能替代局部错误处理。用户正在提交表单时，页面仍然应该给出明确反馈。

## 实际项目常见问题

### 1. loading 一直不结束

**原因**

失败分支没有恢复状态。

**解决方案**

把恢复逻辑放到 `finally`：

```ts
loading.value = true

try {
  await loadData()
} finally {
  loading.value = false
}
```

### 2. 接口 500 进入了成功逻辑

**原因**

使用 `fetch` 时没有检查 `response.ok`。

**解决方案**

在 request 层统一检查 HTTP 状态。

### 3. 多个接口同时 401，弹出多个提示

**解决方案**

给登录过期处理加锁：

```ts
let redirectingLogin = false

function handleUnauthorized() {
  if (redirectingLogin) return

  redirectingLogin = true
  clearAuth()
  router.replace('/login')
}
```

### 4. 错误提示太技术化

不要直接把后端栈信息、SQL 错误、英文异常展示给用户。用户提示要可理解，技术细节进入日志。

### 5. catch 里又抛错

错误处理代码也可能出错。日志上报、提示组件、路由跳转都要考虑失败。

保持 catch 简单：

```ts
catch (error) {
  const message = getErrorMessage(error)
  showError(message)
  reportError(error)
}
```

## 最佳实践

- 每个异步流程都设计 success、error、finally。
- request 层统一 HTTP 错误和业务错误。
- 页面层负责用户提示和可恢复操作。
- 不吞错误，至少记录或重新抛出。
- 用户提示讲人话，技术细节进日志。
- 401 和刷新 token 要防并发和死循环。
- 上线后用全局错误兜底捕获漏网问题。

## 学习检查

学完本节后，你应该能回答：

- `try...catch` 能捕获哪些错误，捕获不到哪些错误。
- 为什么 fetch 遇到 500 不一定进入 catch。
- request、api、service、页面分别应该处理什么错误。
- 为什么 `finally` 对 loading 和提交锁很重要。
- 线上错误日志应该保留哪些上下文。

## 参考资料

- [MDN: Control flow and error handling](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Control_flow_and_error_handling)
- [MDN: try...catch](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/try...catch)
- [MDN: Error](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error)
- [MDN: Using promises](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Using_promises)

## 下一步学习

继续学习 [内存管理](/javascript/memory-management)，理解事件监听、定时器、缓存和组件卸载中的泄漏风险。
