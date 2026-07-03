# 图解 Go 核心概念

## 这个页面解决什么

Go 语法看起来简单，但真正写后端项目时，难点通常在模块、包边界、接口、goroutine、channel、context、连接池和性能诊断。

这一页用图先建立整体理解。你可以把它当作 Go 模块的“地图页”。

## 一张图理解 Go 项目从代码到运行

```mermaid
flowchart LR
  A[".go 源码"] --> B["go test<br/>执行测试"]
  A --> C["go build<br/>构建"]
  C --> D["单个二进制文件"]
  D --> E["容器 / 服务器运行"]
  F["go.mod"] --> C
  G["go.sum"] --> C
```

Go 项目通常构建成一个二进制文件，部署时不需要像 Node.js 一样把 `node_modules` 一起带上，也不需要像 Java 一样依赖 JVM。

但这不表示部署可以随意：

- 仍要管理配置。
- 仍要有数据库迁移。
- 仍要有健康检查。
- 仍要处理证书、时区、日志和优雅关闭。

## 一张图理解 go.mod、go.sum、模块缓存

```mermaid
flowchart TD
  A["当前模块"] --> B["go.mod<br/>声明模块路径和依赖"]
  A --> C["go.sum<br/>依赖校验信息"]
  B --> D["直接依赖"]
  D --> E["间接依赖"]
  D --> F["模块缓存"]
  E --> F
  G["go mod tidy"] --> B
  G --> C
```

关键理解：

- `go.mod` 决定项目依赖什么。
- `go.sum` 用来校验依赖内容是否一致。
- `go mod tidy` 会根据源码重新计算依赖。
- CI 失败时，先检查 Go 版本、私有依赖权限、`go.sum` 和 `replace`。

## 一张图理解包边界和 internal

```mermaid
flowchart TD
  A["cmd/server<br/>程序入口"] --> B["internal/user<br/>用户业务"]
  A --> C["internal/order<br/>订单业务"]
  B --> D["internal/platform/db<br/>数据库基础设施"]
  C --> D
  E["外部项目"] -. "不能导入" .-> B
  F["pkg/logger<br/>确实要复用的公共包"] --> A
  E --> F
```

`internal` 是 Go 的真实边界机制。放在 `internal` 下的包不能被外部模块导入。

这能帮助你避免两类问题：

- 应用内部实现被其他项目依赖，后续不敢改。
- `pkg` 变成什么都放的公共垃圾桶。

## 一张图理解 Handler、Service、Repository

```mermaid
flowchart TD
  A["HTTP Router"] --> B["Handler<br/>解析请求、返回响应"]
  B --> C["Service<br/>业务规则、事务编排"]
  C --> D["Repository Interface<br/>数据访问契约"]
  E["MySQL Repository<br/>具体实现"] -. "实现接口" .-> D
  E --> F[("Database")]
  C --> G["External Client<br/>外部服务"]
```

代码放置建议：

| 逻辑 | 位置 |
| --- | --- |
| 解析 URL、JSON、Header | Handler |
| 参数基础校验 | Handler 或 Request DTO |
| 业务规则 | Service |
| 事务编排 | Service |
| SQL | Repository |
| 外部 HTTP/gRPC 调用 | Client |

Go 不需要为了“架构感”写很多层，但基本边界要清楚。

## 一张图理解 goroutine 生命周期

```mermaid
stateDiagram-v2
  [*] --> Created: go func()
  Created --> Running: scheduler 调度
  Running --> Waiting: 等待 channel / I/O / lock / timer
  Waiting --> Running: 条件满足
  Running --> Done: 函数返回
  Done --> [*]
  Waiting --> Leaked: 没有取消信号或没人接收
```

每个 goroutine 都必须有退出条件。常见泄漏原因：

- channel 发送后没人接收。
- 外部 HTTP 请求没有超时。
- 后台循环没有监听 context。
- 定时任务没有停止机制。

## 一张图理解 channel 发送和接收

```mermaid
sequenceDiagram
  participant G1 as goroutine A
  participant CH as channel
  participant G2 as goroutine B

  G1->>CH: send value
  Note over G1,CH: 无缓冲 channel 会等待接收方
  G2->>CH: receive value
  CH-->>G2: value
  G1-->>G1: send 完成
```

无缓冲 channel 的核心是同步：发送方和接收方要同时准备好。

如果你只是要保护一个共享计数器，`sync.Mutex` 可能比 channel 更直观。channel 更适合表达任务传递、结果收集、取消信号和 worker pool。

## 一张图理解 context 取消传播

```mermaid
flowchart TD
  A["HTTP Request Context"] --> B["Handler"]
  B --> C["Service"]
  C --> D["Repository QueryContext"]
  C --> E["External HTTP Request"]
  C --> F["goroutine worker"]

  G["客户端断开 / 超时"] --> A
  A -. "Done 关闭" .-> B
  A -. "Done 关闭" .-> C
  A -. "取消 SQL" .-> D
  A -. "取消 HTTP" .-> E
  A -. "通知退出" .-> F
```

context 的作用不是“传所有参数”，而是传请求生命周期。

推荐规则：

- Handler 从 `r.Context()` 获取 ctx。
- Service、Repository、Client 都接收 ctx。
- 数据库调用使用 `QueryContext`、`ExecContext`。
- 外部 HTTP 请求绑定 ctx。
- 后台 goroutine 监听 `ctx.Done()`。

## 一张图理解 database/sql 连接池

```mermaid
flowchart TD
  A["Handler 请求"] --> B["Repository"]
  B --> C["*sql.DB<br/>连接池句柄"]
  C --> D{"是否有空闲连接"}
  D -- "有" --> E["复用连接"]
  D -- "没有但未超上限" --> F["创建新连接"]
  D -- "达到上限" --> G["等待连接"]
  E --> H[("Database")]
  F --> H
  G --> H
  H --> I["rows.Close / tx.Commit / tx.Rollback"]
  I --> C
```

`*sql.DB` 不是单个连接，而是连接池。

常见问题：

- 忘记 `rows.Close()`，连接不能归还。
- 事务太长，连接被占用。
- 慢 SQL 导致连接池排队。
- 并发无限制，瞬间打满数据库。

## 一张图理解 Go 性能排查

```mermaid
flowchart TD
  A["性能问题"] --> B{"先判断现象"}
  B --> C["CPU 高"]
  B --> D["内存高"]
  B --> E["goroutine 数增长"]
  B --> F["接口慢"]

  C --> C1["CPU profile"]
  C1 --> C2["热点函数、循环、JSON、压缩"]

  D --> D1["heap profile"]
  D1 --> D2["大对象、slice、缓存、字符串拼接"]

  E --> E1["goroutine profile"]
  E1 --> E2["channel 阻塞、无超时 I/O、后台循环"]

  F --> F1["链路耗时拆分"]
  F1 --> F2["SQL、Redis、外部接口、连接池、锁"]
```

不要一开始就优化代码。先用指标和 profile 确认瓶颈，再修改。

## 建议阅读顺序

如果你第一次学 Go，建议：

1. 先读本页，建立模块、包、并发、context 和连接池模型。
2. 再读 [环境、模块与工作区](/go/setup-modules)。
3. 再读 [语法、类型与函数](/go/syntax-types)。
4. 学 HTTP 服务前，回来看 “Handler、Service、Repository” 和 “context 取消传播”。
5. 排查性能时，回来看 “Go 性能排查”。

## 下一步学习

继续学习 [环境、模块与工作区](/go/setup-modules)。
