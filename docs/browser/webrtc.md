# WebRTC 实时音视频

## 适合谁看

适合想在浏览器里实现音视频通话、屏幕共享、实时会议、点对点数据通道，或者需要理解 WebRTC 为什么比普通 WebSocket 更复杂的学习者。

WebRTC 是实时通信能力集合，不只是一个 API。它涉及媒体采集、点对点连接、网络穿透、编解码、信令、权限、安全上下文和浏览器兼容。

## WebRTC 能做什么

常见场景：

- 视频会议。
- 语音通话。
- 屏幕共享。
- 在线客服。
- 远程协作。
- 点对点文件传输。
- 低延迟数据通道。

核心能力：

| 能力 | API |
| --- | --- |
| 摄像头和麦克风 | `getUserMedia` |
| 屏幕共享 | `getDisplayMedia` |
| 点对点连接 | `RTCPeerConnection` |
| 数据通道 | `RTCDataChannel` |
| 媒体播放 | `HTMLMediaElement` |

## WebRTC 和 WebSocket 的区别

| 对比 | WebSocket | WebRTC |
| --- | --- | --- |
| 通信对象 | 浏览器和服务器 | 浏览器和浏览器，或浏览器和媒体服务 |
| 适合内容 | 消息、状态、进度 | 音视频、实时媒体、低延迟数据 |
| 是否需要信令 | 可以不需要独立信令 | 必须需要信令 |
| 网络复杂度 | 相对低 | 高，需要 NAT 穿透 |
| 常见服务 | Socket 服务 | STUN、TURN、信令、媒体服务器 |

WebRTC 的媒体可以点对点传输，但双方建立连接前需要先交换连接信息，这个过程叫信令。信令通常用 WebSocket 或 HTTP 实现。

## 基础流程

简化流程：

```text
获取本地媒体
↓
创建 RTCPeerConnection
↓
添加本地音视频轨道
↓
通过信令交换 offer / answer
↓
交换 ICE candidate
↓
连接建立
↓
播放远端媒体
```

这就是 WebRTC 比普通请求复杂的原因：它不只是“调用接口”，而是两个浏览器之间协商一条实时媒体通道。

## 获取摄像头和麦克风

```ts
const localStream = await navigator.mediaDevices.getUserMedia({
  video: true,
  audio: true
})

const video = document.querySelector('video')!
video.srcObject = localStream
```

注意：

- 通常需要 HTTPS。
- 用户可能拒绝权限。
- 设备可能不存在。
- 浏览器可能限制后台或 iframe 权限。

错误处理：

```ts
try {
  const stream = await navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true
  })
} catch (error) {
  showError('无法访问摄像头或麦克风')
}
```

## 创建连接

```ts
const peer = new RTCPeerConnection({
  iceServers: [
    { urls: 'stun:stun.l.google.com:19302' }
  ]
})

localStream.getTracks().forEach((track) => {
  peer.addTrack(track, localStream)
})

peer.ontrack = (event) => {
  remoteVideo.srcObject = event.streams[0]
}

peer.onicecandidate = (event) => {
  if (event.candidate) {
    sendSignal({
      type: 'candidate',
      candidate: event.candidate
    })
  }
}
```

STUN 用来帮助发现公网地址。TURN 用来在点对点失败时中继流量。生产环境不能只依赖公共 STUN。

## 信令

信令不是 WebRTC 标准里的具体协议。你可以用 WebSocket、HTTP 或其他方式交换：

- offer。
- answer。
- ICE candidate。
- 房间信息。
- 用户加入离开。
- 静音、关闭摄像头等状态。

示意：

```ts
sendSignal({
  type: 'offer',
  roomId: 'room-001',
  sdp: offer
})
```

项目建议：

- 信令消息要有稳定类型。
- 房间和用户身份要鉴权。
- 断线重连要重新同步状态。
- 不要把媒体流本身通过信令服务器转发。

## 数据通道

WebRTC 不只传音视频，也可以传点对点数据。

```ts
const channel = peer.createDataChannel('chat')

channel.onopen = () => {
  channel.send('hello')
}

channel.onmessage = (event) => {
  console.log(event.data)
}
```

适合：

- 低延迟状态同步。
- 小文件分片传输。
- 白板协作数据。

如果只是普通业务消息，WebSocket 更简单。

## 实际项目常见问题

### 1. 本地能通，线上不能通

常见原因：

- 线上不是 HTTPS。
- 摄像头或麦克风权限被拒绝。
- NAT 穿透失败。
- 没有 TURN 服务。
- 信令服务器连接失败。

处理：

- 确认 HTTPS。
- 检查权限状态。
- 增加 TURN。
- 打印 ICE connection state。
- 检查信令消息是否完整交换。

### 2. 一对一可以，多人会议很卡

点对点 Mesh 模式下，每个用户都要给其他用户发流。人数增加后带宽和 CPU 快速上升。

多人会议通常需要媒体服务器，例如 SFU。

项目取舍：

| 规模 | 建议 |
| --- | --- |
| 一对一 | P2P 可以考虑 |
| 小房间 | 可先评估 P2P 或 SFU |
| 多人会议 | 优先 SFU |
| 录制、转码、审计 | 需要服务端媒体能力 |

### 3. 用户关闭页面后摄像头还亮

常见原因是没有停止媒体轨道。

处理：

```ts
localStream.getTracks().forEach((track) => {
  track.stop()
})
```

页面卸载、离开房间、切换账号时都要释放媒体资源。

## 安全和隐私

WebRTC 涉及摄像头、麦克风、屏幕和用户网络信息，必须谨慎：

- 明确告知用户正在采集什么。
- 用户离开房间后停止媒体。
- 不默认开启摄像头和麦克风。
- 屏幕共享前给清晰提示。
- 信令和房间加入必须鉴权。
- 日志不要记录敏感 SDP 或用户隐私数据。

## 项目建议

- 先做一对一最小闭环，再扩多人。
- 生产环境准备 TURN，不只依赖 STUN。
- 信令协议要可追踪、可重试、可排查。
- 监控连接状态、失败原因和设备权限。
- 离开页面时释放媒体轨道和连接。

## 下一步学习

- [WebSocket 实时通信](/browser/websocket)
- [浏览器安全基础](/browser/security)
- [常用 Web API](/browser/web-apis)
- [部署、缓存与 DevOps 问题](/projects/issues-deployment)
