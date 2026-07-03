# WebSocket 实时通信

## 适合谁看

适合已经理解 HTTP 请求，但项目里开始遇到实时消息、在线状态、通知、聊天、协同编辑、进度推送等需求的学习者。

WebSocket 解决的是“客户端和服务器需要持续双向通信”的问题。它不是用来替代所有 HTTP 接口的，也不是只要实时就一定要用。学习它的重点是理解连接生命周期、消息协议、鉴权、重连、心跳和降级策略。

## WebSocket 解决什么问题

普通 HTTP 请求是一次请求一次响应：

```text
浏览器发请求
↓
服务器返回响应
↓
连接结束或复用
```

如果页面想知道服务器有没有新消息，常见做法是轮询：

```text
每 5 秒请求一次 /api/messages
```

轮询简单，但有缺点：

- 延迟取决于轮询间隔。
- 没有消息时也会发请求。
- 高频轮询会增加服务器压力。
- 用户越多，浪费越明显。

WebSocket 建立连接后可以双向发送消息：

```text
浏览器 ⇄ WebSocket 连接 ⇄ 服务器
```

适合：

- 聊天。
- 在线协作。
- 实时通知。
- 订单状态推送。
- 大任务进度推送。
- 股票、监控、日志等实时流。

不适合：

- 普通 CRUD。
- 低频数据查询。
- 一次性提交表单。
- 不需要实时性的后台列表。

## 基础用法

```ts
const socket = new WebSocket('wss://example.com/ws')

socket.addEventListener('open', () => {
  console.log('连接已建立')
  socket.send(JSON.stringify({ type: 'ping' }))
})

socket.addEventListener('message', (event) => {
  const message = JSON.parse(event.data)
  console.log('收到消息', message)
})

socket.addEventListener('close', () => {
  console.log('连接已关闭')
})

socket.addEventListener('error', (event) => {
  console.error('连接错误', event)
})
```

线上环境通常使用 `wss://`，不要在 HTTPS 页面里连接不安全的 `ws://`。

## 消息协议

不要直接发送散乱字符串。项目里应该定义稳定消息结构。

```ts
type SocketMessage =
  | { type: 'ping'; payload?: never }
  | { type: 'notification'; payload: { id: string; title: string } }
  | { type: 'task-progress'; payload: { taskId: string; percent: number } }
  | { type: 'error'; payload: { code: string; message: string } }
```

发送：

```ts
socket.send(JSON.stringify({
  type: 'task-progress',
  payload: {
    taskId: 'export-001',
    percent: 60
  }
}))
```

处理：

```ts
function handleMessage(message: SocketMessage) {
  switch (message.type) {
    case 'notification':
      showNotification(message.payload)
      break
    case 'task-progress':
      updateTaskProgress(message.payload)
      break
    case 'error':
      showError(message.payload.message)
      break
  }
}
```

## 鉴权

WebSocket 连接也需要鉴权。常见方式：

| 方式 | 说明 | 注意 |
| --- | --- | --- |
| Cookie | 连接同源或正确跨域时自动带上 | 受 SameSite、Domain、Secure 影响 |
| query token | URL 中携带短期 token | URL 可能进入日志，避免长期敏感 token |
| 首条消息鉴权 | 连接后先发送 auth 消息 | 未鉴权前不能订阅业务消息 |
| 子协议 | 使用 `Sec-WebSocket-Protocol` | 需要前后端协议一致 |

不要把长期 token 直接暴露在 URL 中。

更稳妥的方式是使用短期连接 token：

```text
1. 普通 HTTPS 接口申请短期 socket token。
2. 使用短期 token 建立 WebSocket。
3. 服务端验证通过后允许订阅业务频道。
```

## 心跳和重连

真实项目必须考虑连接断开。

常见断开原因：

- 网络切换。
- 浏览器标签页休眠。
- 网关空闲连接超时。
- 服务重启。
- 用户登录态过期。

心跳示例：

```ts
let heartbeatTimer: number | undefined

function startHeartbeat(socket: WebSocket) {
  heartbeatTimer = window.setInterval(() => {
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'ping' }))
    }
  }, 30000)
}

function stopHeartbeat() {
  if (heartbeatTimer) {
    window.clearInterval(heartbeatTimer)
  }
}
```

重连要有退避策略，不要断开后疯狂重连。

```ts
let retryCount = 0

function getRetryDelay() {
  return Math.min(1000 * 2 ** retryCount, 30000)
}
```

## 和 HTTP 的分工

推荐分工：

| 任务 | 推荐方式 |
| --- | --- |
| 查询列表 | HTTP |
| 创建、编辑、删除 | HTTP |
| 实时通知 | WebSocket |
| 导出进度 | WebSocket 或 SSE |
| 聊天消息 | WebSocket |
| 在线协作 | WebSocket |

不要为了“统一”把所有业务都走 WebSocket。HTTP 更适合可审计、可缓存、可重试的普通请求。

## 实际项目常见问题

### 1. 本地能连，线上连接失败

常见原因：

- HTTPS 页面连接了 `ws://`。
- Nginx 没配置 WebSocket upgrade。
- 网关超时太短。
- 跨域和鉴权配置不一致。

Nginx 示例：

```nginx
location /ws/ {
  proxy_pass http://socket-service/;
  proxy_http_version 1.1;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  proxy_set_header Host $host;
}
```

### 2. 页面开久了收不到消息

常见原因：

- 连接被网关断开。
- 浏览器标签页进入后台。
- 心跳没有实现。
- 服务端没有感知断线。

处理：

- 加心跳。
- 加断线重连。
- 页面重新可见时检查连接状态。
- 服务端清理失效连接。

### 3. 重连后收到重复消息

常见原因：

- 重新订阅时没有去重。
- 服务端重复推送历史消息。
- 客户端没有使用消息 id。

处理：

- 消息带唯一 id。
- 客户端按 id 去重。
- 订阅协议定义从哪个位置继续。

## 项目建议

- 定义稳定消息协议，不发散乱字符串。
- 鉴权、订阅、业务消息分阶段处理。
- 必须有心跳、重连和关闭清理。
- 业务消息带 id，便于去重和追踪。
- 普通 CRUD 继续使用 HTTP。
- 上线前验证 Nginx 或网关是否支持 WebSocket upgrade。

## 下一步学习

- [HTTP 与请求流程](/browser/http-request)
- [跨域与登录态](/browser/cors-auth)
- [Nginx 速查](/cheatsheets/nginx)
- [后端接口与服务问题](/projects/issues-backend)
