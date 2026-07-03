# 运行时与事件循环

## 适合谁看

适合知道 JavaScript，但不清楚 Node.js 和浏览器有什么区别、为什么 Node 能处理很多请求的学习者。

## Node.js 是什么

Node.js 让 JavaScript 可以在浏览器之外运行。它使用 V8 引擎执行 JavaScript，并提供文件系统、网络、进程等服务器能力。

浏览器提供：

- DOM。
- BOM。
- 页面事件。
- fetch。

Node.js 提供：

- 文件系统。
- HTTP 服务。
- 进程管理。
- 路径处理。
- stream。

## 非阻塞 I/O

Node.js 官方文档说明，事件循环让 Node.js 可以在默认单线程 JavaScript 执行模型下处理非阻塞 I/O。

简单理解：

```text
JavaScript 线程处理业务逻辑
↓
耗时 I/O 交给系统
↓
I/O 完成后回调或 Promise 继续执行
```

这让 Node 很适合处理大量 I/O 型任务，例如 API 请求、数据库查询、文件读写。

## async/await

```ts
async function getUser() {
  const user = await userRepository.findById(1)
  return user
}
```

`await` 不等于阻塞整个 Node 进程。它暂停当前 async 函数，等待 Promise 结果。

## CPU 密集型任务

Node 不适合在主线程长时间执行 CPU 密集计算，例如大规模图片处理、复杂加密、超大数据计算。

这些任务会阻塞事件循环，导致其他请求延迟。

解决方式：

- 使用 Worker Threads。
- 交给专门服务。
- 拆分任务。
- 放到队列异步处理。

## 实际项目常见问题

### 1. 一个接口很慢导致其他接口也慢

**原因**

接口里有同步 CPU 密集任务或同步文件操作。

**解决方案**

避免在请求处理里执行长时间同步任务。使用异步 API 或任务队列。

### 2. await 会不会阻塞服务

`await` 不会阻塞整个进程，但如果 await 的 Promise 内部执行的是同步重计算，仍然会阻塞。

### 3. setTimeout 不准

如果事件循环被阻塞，定时器也会延迟。

## 最佳实践

- API 服务里避免长时间同步计算。
- 文件、数据库、网络使用异步 API。
- 慢任务放后台队列。
- 监控接口耗时和事件循环延迟。

## 下一步

继续学习 [包管理与模块化](/node/package-modules)。
