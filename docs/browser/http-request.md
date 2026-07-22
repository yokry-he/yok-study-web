# HTTP 与请求流程

## 适合谁看

适合想把“接口请求到底发生了什么”搞清楚的人。学完这一节，你应该能看懂 Network 面板，并能判断一个接口问题到底是 URL、请求方法、参数、请求头、响应头、状态码还是业务响应的问题。

## 一次请求经历什么

以前端调用接口为例：

```ts
const response = await fetch('/api/users?page=1')
const data = await response.json()
```

浏览器大致会做这些事：

1. 根据当前页面地址和请求地址计算最终 URL。
2. 判断是否跨源。
3. 根据缓存策略判断是否可以复用缓存。
4. 必要时发起 DNS、TCP、TLS 连接。
5. 发送请求行、请求头和请求体。
6. 接收状态码、响应头和响应体。
7. 根据 CORS、MIME、缓存等规则决定前端代码能否读取响应。
8. JavaScript 继续处理响应数据。

框架里的 `axios`、`fetch` 封装只是请求入口，不会绕过浏览器规则。

## 请求方法

| 方法 | 常见用途 | 是否应有副作用 |
| --- | --- | --- |
| `GET` | 查询列表、详情、配置 | 不应该修改数据 |
| `POST` | 创建、提交表单、复杂查询 | 通常会修改数据 |
| `PUT` | 全量更新资源 | 会修改数据 |
| `PATCH` | 局部更新资源 | 会修改数据 |
| `DELETE` | 删除资源 | 会修改数据 |
| `OPTIONS` | 预检请求 | 浏览器或服务器协商用 |

实际项目中不要把所有接口都写成 `POST`。方法语义清楚后，缓存、网关、日志、权限和接口文档都更容易维护。

## 状态码怎么判断

| 状态码 | 含义 | 前端常见处理 |
| --- | --- | --- |
| `200` | 请求成功 | 继续判断业务 code |
| `201` | 创建成功 | 表单提交后提示成功 |
| `204` | 成功但无响应体 | 删除、开关更新常见 |
| `301/302` | 重定向 | 登录跳转、网关跳转 |
| `304` | 缓存仍有效 | 浏览器复用本地缓存 |
| `400` | 参数错误 | 显示校验提示或检查请求参数 |
| `401` | 未登录或登录过期 | 清理登录态并跳转登录 |
| `403` | 已登录但无权限 | 显示无权限页面或按钮禁用 |
| `404` | 资源不存在 | 检查接口路径或前端路由 fallback |
| `409` | 数据冲突 | 版本冲突、重复提交 |
| `422` | 业务校验失败 | 表单字段提示 |
| `429` | 请求过于频繁 | 限流提示、退避重试 |
| `500` | 服务端异常 | 记录日志并提示稍后重试 |
| `502/503/504` | 网关或服务不可用 | 检查网关、服务实例、超时 |

很多团队还有业务 code，例如：

```json
{
  "code": 0,
  "message": "ok",
  "data": []
}
```

状态码表示 HTTP 层结果，业务 code 表示业务层结果。不要混在一起理解。

## 请求头和响应头

### 常见请求头

| Header | 作用 |
| --- | --- |
| `Authorization` | token、Bearer token、签名 |
| `Content-Type` | 请求体格式，例如 JSON 或表单 |
| `Accept` | 期望响应格式 |
| `Cookie` | 浏览器自动携带的 Cookie |
| `Origin` | 跨源请求来源 |
| `Referer` | 当前页面来源 |

### 常见响应头

| Header | 作用 |
| --- | --- |
| `Content-Type` | 响应体格式 |
| `Set-Cookie` | 服务端设置 Cookie |
| `Cache-Control` | 缓存策略 |
| `ETag` | 缓存校验标识 |
| `Last-Modified` | 资源最后修改时间 |
| `Access-Control-Allow-Origin` | CORS 允许的来源 |
| `Access-Control-Allow-Credentials` | 是否允许携带凭据的跨源请求 |

## Content-Type 对请求的影响

### JSON

```ts
await fetch('/api/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'Alice'
  })
})
```

适合普通业务对象。后端需要按 JSON 解析请求体。

### FormData

```ts
const formData = new FormData()
formData.append('file', file)
formData.append('name', 'avatar')

await fetch('/api/upload', {
  method: 'POST',
  body: formData
})
```

上传文件时不要手动设置 `Content-Type`，浏览器会自动补上 `multipart/form-data` 的 boundary。手动写反而容易导致后端解析失败。

## Network 面板排查流程

先观察一条已经成功返回的请求，确认五类事实：请求是否发出、目标 URL 是否正确、方法是否正确、身份凭证是否带上、响应类型和 Request ID 是否能与后端日志对应。

<DocFigure
  src="/images/browser/network-request-headers.webp"
  alt="浏览器网络请求证据面板展示 URL、GET 方法、200 状态、请求头和响应头"
  caption="Headers 面板把一次请求的输入和输出放在一起；先核对事实，再解释页面为什么成功或失败。"
  :width="1440"
  :height="900"
/>

不依赖图片的读取路径：打开 DevTools → Network → 选择目标请求 → General 查看 URL、Method、Status → Request Headers 查看身份与内容类型 → Response Headers 查看缓存、跨域、响应类型和链路 ID。

当问题是“请求为什么慢”时，不要只看总耗时。下面的示例把 218 ms 拆成 DNS、连接与 TLS、等待首字节和下载四段，其中 142 ms 都花在 Waiting。

<DocFigure
  src="/images/browser/network-request-timing.webp"
  alt="浏览器 Network Timing 将 218 毫秒请求拆成 DNS、连接、等待首字节和下载阶段"
  caption="Waiting 占比最高时，应先用 Request ID 对齐服务端日志，而不是先优化很小的响应体。"
  :width="1440"
  :height="900"
/>

不依赖图片的读取路径：选择请求 → Timing → 比较 Queueing、DNS、Initial connection、SSL、Waiting 和 Content Download；记录占比最大的阶段与具体毫秒数。

### 1. 找到请求

按接口路径、方法或 `Fetch/XHR` 过滤。先确认请求有没有发出去。

如果没有发出去，问题可能在：

- 代码分支没有执行。
- 按钮被禁用。
- 请求被防抖、节流、缓存或取消。
- 前端运行时报错中断。

### 2. 看 URL

重点看：

- 域名是否正确。
- 端口是否正确。
- 路径是否多了或少了 `/api`。
- 查询参数是否为空。
- 环境变量是否用了错误环境。

### 3. 看 Request Headers

重点看：

- `Authorization` 是否存在。
- Cookie 是否带出去了。
- `Content-Type` 是否和请求体匹配。
- 跨域时 `Origin` 是什么。

### 4. 看 Response Headers

重点看：

- 是否有 `Set-Cookie`。
- 是否有 CORS 相关响应头。
- 是否有缓存头。
- 响应类型是否正确。

### 5. 看 Response Body

很多接口虽然 HTTP 状态码是 200，但业务 code 表示失败。此时应该按业务错误处理，而不是按请求成功处理。

## 实际项目问题

### 问题：接口偶尔被取消

**现象**

Network 显示 `canceled`，页面数据没有更新。

**常见原因**

- 路由切换时组件卸载，请求被取消。
- 用户连续输入搜索，旧请求被 AbortController 取消。
- 表格筛选条件变化太快，多个请求互相覆盖。

**解决方案**

给请求建立明确策略：

```ts
let currentController: AbortController | null = null

export async function searchUsers(keyword: string) {
  currentController?.abort()

  const controller = new AbortController()
  currentController = controller

  const response = await fetch(`/api/users?keyword=${encodeURIComponent(keyword)}`, {
    signal: controller.signal
  })

  return response.json()
}
```

如果是搜索框，这是合理行为；如果是保存表单，就不应该随便取消。

### 问题：接口 200 但页面提示失败

**原因**

HTTP 层成功，不代表业务成功。例如：

```json
{
  "code": 10001,
  "message": "库存不足",
  "data": null
}
```

**解决方案**

请求封装里明确区分网络错误、HTTP 错误和业务错误：

```ts
async function request<T>(url: string): Promise<T> {
  const response = await fetch(url)

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  const result = await response.json()

  if (result.code !== 0) {
    throw new Error(result.message || '业务处理失败')
  }

  return result.data
}
```

## 最佳实践

- 请求封装只做通用处理，不要把具体业务逻辑塞进去。
- 接口错误要分层：网络错误、HTTP 错误、业务错误、渲染错误。
- 文件上传不要手写 `multipart/form-data` 的 `Content-Type`。
- 重要写操作要处理重复提交、超时和服务端幂等。
- 排查接口问题时先看 Network，再看代码。

## 下一步学习

继续学习 [跨域与登录态](/browser/cors-auth)。
