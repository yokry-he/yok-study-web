# 调试工具速查

## 浏览器 DevTools

| 面板 | 用途 |
| --- | --- |
| Elements | 查看 DOM、CSS、布局 |
| Console | 查看日志和运行 JS |
| Network | 查看请求、响应、缓存、耗时 |
| Sources | 断点调试 JS |
| Application | 查看 storage、cookie、service worker |
| Performance | 分析运行性能 |
| Lighthouse | 做页面质量检查 |

## Network 排查

优先看：

| 项 | 说明 |
| --- | --- |
| Status | 200、400、401、403、500 |
| Method | GET、POST、PUT、DELETE |
| Request URL | 路径和 query 是否正确 |
| Request Headers | token、cookie、content-type |
| Payload | body、form-data 是否正确 |
| Response | 错误码和 message |
| Timing | DNS、连接、等待、下载耗时 |

接口联调时先保存 Network 证据，再讨论前后端问题。

## Console 排查

常见信息：

| 类型 | 处理 |
| --- | --- |
| `ReferenceError` | 变量不存在或作用域错误 |
| `TypeError` | 值类型不符合预期 |
| Promise rejection | 异步错误没有捕获 |
| CORS error | 跨域配置或请求头问题 |
| 资源 404 | 路径、base、部署目录错误 |

不要忽略红色错误。一个控制台错误可能导致后续组件不渲染。

## Sources 断点

常用能力：

- 点击行号设置断点。
- 使用 conditional breakpoint。
- 查看 call stack。
- 查看 scope 变量。
- 使用 step over / step into。
- 在异常处暂停。

适合排查：

- 点击后状态为什么没变。
- 表单提交参数为什么不对。
- 权限判断为什么返回 false。
- 异步请求回调顺序。

## Application 面板

常查：

| 区域 | 用途 |
| --- | --- |
| Local Storage | 本地状态 |
| Session Storage | 会话状态 |
| Cookies | 登录态 |
| IndexedDB | 本地数据库 |
| Cache Storage | PWA 或缓存 |
| Service Workers | SW 注册和更新 |

登录问题重点看 Cookie、Local Storage 和请求头是否一致。

## Vue DevTools

常用：

- 查看组件树。
- 查看 props。
- 查看 emits。
- 查看 Pinia store。
- 查看路由状态。

适合排查：

- props 是否传错。
- store 是否刷新丢失。
- 组件是否重复渲染。
- 路由参数是否正确。

## 性能排查

常用顺序：

```text
Network 看资源大小和请求瀑布
↓
Performance 录制交互
↓
检查长任务和重复渲染
↓
分析包体积
↓
做懒加载或拆包
```

不要只凭感觉优化。先定位是网络慢、JS 执行慢、渲染慢，还是接口慢。

## Node 和服务端排查

| 工具 | 用途 |
| --- | --- |
| `console.log` | 临时定位 |
| 结构化日志 | 线上排查 |
| request id | 串联前端和后端 |
| `node --inspect` | Node 调试 |
| 健康检查 | 判断服务是否可用 |

线上不要打印敏感信息，例如 token、cookie、密码、身份证。

## 推荐排查顺序

```text
复现问题
↓
记录 URL、用户、时间
↓
看 Console
↓
看 Network
↓
看前端状态
↓
看后端日志 request id
↓
看数据库、缓存、部署层
```

## 参考资料

- [Chrome DevTools](https://developer.chrome.com/docs/devtools)
- [Debug JavaScript in Chrome DevTools](https://developer.chrome.com/docs/devtools/javascript)
- [Chrome DevTools Performance reference](https://developer.chrome.com/docs/devtools/performance/reference)
- [MDN Browser Developer Tools](https://developer.mozilla.org/en-US/docs/Learn_web_development/Howto/Tools_and_setup/What_are_browser_developer_tools)

## 延伸学习

- [浏览器自动化调试](/browser/browser-automation-debugging)
- [前端页面与状态问题](/projects/issues-frontend)
- [前后端联调排查](/projects/integration-debugging)
