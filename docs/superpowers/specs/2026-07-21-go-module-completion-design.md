# Go 模块全量完善设计

## 1. 背景

当前站点已经有 15 篇 Go 文档，覆盖语法、接口、并发、HTTP、数据库、测试、部署、性能和 gRPC，也已经有一篇较长的 HTTP API 项目说明。

现阶段的主要问题不是“没有内容”，而是内容之间还没有形成可验证的学习闭环：

- 图解页只有 9 张图，slice、interface、channel、事务、错误链和服务关闭等高频难点仍主要依赖文字理解。
- HTTP API 项目页包含较多独立代码片段，但没有连接到一个可以构建、测试和启动的完整示例工程。
- 项目说明中的技术选择仍保留“标准库或轻量路由库”“多个迁移工具任选”等分叉，初学者难以判断应该照哪一条路线实践。
- Go 的真实项目问题散落在通用后端问题库里，缺少针对 goroutine、context、连接池、nil interface 和部署问题的系统排查手册。
- 缺少从语言基础逐步过渡到可交付 API 的专项练习和明确验收标准。
- 现有项目文档中存在需要一并修正的内容错误，例如 ER 图中的重复字段、超时中间件示例中的重复函数声明。

本次完善采用“标准库优先”路线。目标不是证明框架没有价值，而是先让读者看懂 Go 标准能力和真实请求生命周期，再解释什么时候需要 Chi、Gin、GORM 或 sqlc。

## 2. 目标

本次交付完成后，读者应当能够：

1. 从零安装 Go、创建模块、理解包边界并掌握常用语法。
2. 理解 slice、map、interface、error、goroutine、channel 和 context 的运行模型。
3. 使用标准库 `net/http` 实现有统一错误契约、中间件、超时和优雅关闭的 API。
4. 使用 `database/sql` 和 PostgreSQL 完成连接池配置、事务、分页和乐观锁。
5. 运行一个与文档代码完全一致的任务管理 API 示例。
6. 使用 `testing`、`httptest`、race detector、Fuzzing 和 Testcontainers 验证核心行为。
7. 根据现象、证据、根因、修复和回归验证五个步骤处理真实项目问题。
8. 完成 12 个递进练习，并根据命令输出和接口行为判断自己是否真正掌握。

## 3. 非目标

本批次不实现以下内容：

- 用户登录、JWT、OAuth2、RBAC 或多租户权限系统。
- 消息队列、分布式事务、服务发现、服务网格和 Kubernetes 编排。
- 完整实现第二套 Gin、Chi、GORM 或 sqlc 示例工程。
- 把 gRPC 示例扩展为多服务生产平台。
- 提供云厂商专属基础设施代码。

这些主题可以在后续独立模块中展开。本批次只在选型说明中解释它们的适用边界，不让额外技术栈干扰第一条完整学习路径。

## 4. 技术基线

### 4.1 版本

- 文档基线：Go 1.26.5，这是 2026-07-21 编写本规格时的稳定版本；后续只在核对官方发布记录后更新版本说明。
- 示例工程的 `go.mod`：声明 `go 1.26.0`，保证使用 Go 1.26 语言和标准库能力。
- 数据库：PostgreSQL 18。
- 容器构建：构建阶段使用 `golang:1.26.5-bookworm`，运行阶段使用 `gcr.io/distroless/static-debian12:nonroot`，以非 root 用户运行静态二进制。

版本相关事实只引用 Go、PostgreSQL 和依赖项目的官方资料。页面不把某个补丁版本写成永久不变的结论，并明确团队项目应以 CI 和生产镜像版本为准。

### 4.2 应用依赖

主示例使用以下组合：

| 能力 | 选择 | 采用原因 |
| --- | --- | --- |
| HTTP | 标准库 `net/http`、`http.ServeMux` | 直接学习路由、请求、响应、中间件和服务生命周期 |
| JSON | 标准库 `encoding/json` | 明确请求体限制、未知字段、空值和错误处理 |
| 日志 | 标准库 `log/slog` | 使用结构化字段记录 request id、耗时和错误 |
| 配置 | 标准库 `os`、`time`、`net/url` | 示例配置规模较小，不引入配置框架 |
| 数据访问 | `database/sql` + `github.com/jackc/pgx/v5/stdlib` | 保留标准连接池与事务接口，使用维护成熟的 PostgreSQL 驱动 |
| 数据迁移 | `github.com/golang-migrate/migrate/v4` | 迁移可追踪、可在命令行与容器流程中复用，不自行实现迁移引擎 |
| 单元和 HTTP 测试 | `testing`、`httptest` | 先掌握 Go 原生测试组织方式 |
| 数据库集成测试 | `testcontainers-go` 的 PostgreSQL 模块 | 对真实 PostgreSQL 约束、事务和并发行为做回归验证 |

所有第三方模块在实施时锁定到经官方发布页确认的稳定版本，并提交 `go.sum`。生产运行依赖仅保留 PostgreSQL 驱动和迁移库，Testcontainers 只参与测试。

### 4.3 明确不采用的主路线

| 方案 | 本批次不作为主路线的原因 | 后续适用场景 |
| --- | --- | --- |
| Gin + GORM | 同时隐藏 HTTP 和 SQL 两层关键机制，初学者容易只记框架 API | 需要快速 CRUD、团队已有统一 Gin/GORM 基建 |
| Chi + sqlc | 工程质量高，但一次引入路由抽象、SQL 代码生成和生成物管理，第一项目认知负担偏高 | SQL 较复杂、重视编译期查询类型和薄运行时抽象 |
| 纯手写迁移执行器 | 会把教学重点转移到锁、版本表、失败恢复和并发执行等迁移工具内部问题 | 不采用；使用成熟迁移工具 |

文档会在读者完成标准库项目后给出选型对照，让读者知道何时升级工具，而不是把框架描述为错误选择。

## 5. 内容架构

现有 15 篇 Go 页面保留路由，避免已有链接失效。完善时按以下层次组织：

### 5.1 第一层：导览与心智模型

- `docs/go/introduction.md`：目标、前置知识、版本、学习顺序、每阶段产出和自测入口。
- `docs/go/visual-guide.md`：用 26 张 Mermaid 图建立语言、并发、HTTP、数据库、测试和部署模型。

### 5.2 第二层：语言和工程基础

- `setup-modules.md`：安装、环境变量、模块、工作区、依赖校验、私有模块和 CI。
- `syntax-types.md`：值、指针、数组、slice、map、struct、方法、泛型和常见边界。
- `interfaces-composition.md`：隐式实现、最小接口、组合、nil interface 和依赖边界。
- `errors-logging-config.md`：错误包装、`errors.Is/As`、日志字段、配置解析和启动失败策略。

### 5.3 第三层：并发和服务开发

- `concurrency.md`：goroutine 生命周期、channel 所有权、select、同步原语、race 和泄漏。
- `context-http.md`：context 传播、ServeMux 路由、中间件、请求限制、客户端和服务端超时。
- `database-transaction.md`：连接池、查询、事务、隔离、分页、乐观锁、幂等和慢 SQL。

### 5.4 第四层：测试、性能和交付

- `testing.md`：表格测试、httptest、集成测试、race、Fuzzing、Benchmark 和测试替身。
- `performance.md`：指标优先、pprof、trace、逃逸分析和可复现基准。
- `project-deployment.md`：目录、构建参数、容器、信号、优雅关闭、健康检查和回滚。
- `troubleshooting.md`：按环境、编译、测试、运行、数据库和容器快速定位问题。
- `grpc-service-communication.md`：保留独立进阶路线，并补足与 HTTP/JSON 选择边界和失败处理。

### 5.5 第五层：完整闭环

- `http-api-project-from-zero.md`：改造成真实示例工程的逐步实施手册，所有关键代码从 `examples/go-task-api` 导入。
- `docs/projects/issues-go.md`：16 个 Go 真实项目问题。
- `docs/roadmap/go-practice.md`：12 个递进练习。

## 6. 图解体系

`docs/go/visual-guide.md` 最终包含 26 张独立 Mermaid 图，每张图都必须有“先看什么”“图中发生了什么”“项目中如何使用”和“常见误区”说明。

图解主题如下：

1. 源码、测试、构建和二进制运行链路。
2. `go.mod`、`go.sum`、模块缓存和代理。
3. workspace 与多个模块的本地协作。
4. package、`internal` 和导入边界。
5. 值传递、指针和方法接收者。
6. slice header、底层数组、扩容和共享修改。
7. map 的并发限制和同步选择。
8. interface 的动态类型、动态值和 nil 陷阱。
9. error 创建、包装、判断和 API 映射。
10. `defer`、`panic`、`recover` 的边界。
11. goroutine 从创建到退出或泄漏。
12. channel 所有权、发送、接收和关闭。
13. select、超时、取消和 worker 退出。
14. happens-before、数据竞争和 race detector。
15. context 从 HTTP 传播到 SQL 和外部请求。
16. `net/http` 请求生命周期。
17. 中间件执行顺序与响应提交时机。
18. Handler、Service、Repository 的职责和依赖方向。
19. JSON 解码、校验和统一错误响应。
20. `database/sql` 连接池借用、等待和归还。
21. 事务提交、回滚和失败路径。
22. 乐观锁避免丢失更新。
23. 单元、HTTP、集成和端到端测试分层。
24. 指标、pprof、trace 和证据链。
25. 多阶段构建和最小运行镜像。
26. readiness、信号处理和优雅关闭。

复杂图在窄屏下可以横向滚动，但不得造成页面级横向滚动。每张图必须实际渲染为 SVG，禁止只通过构建成功推断图示可用。

### 6.1 图像媒介选择

本模块不把所有视觉内容都做成 Mermaid，也不为满足数量添加装饰图片。按以下顺序选择媒介：

1. 请求链路、状态转换、依赖关系和并发时序使用 Mermaid，保证技术内容可修改、可搜索并适配深浅色主题。
2. API 的实际输出、容器运行状态和可观测性结果使用示例工程运行后生成的真实截图，保证读者看到的结果可以复现。
3. 只有当空间关系、运行场景或整体心智模型无法通过结构图清楚表达时，才使用 `imagegen` 生成 `scientific-educational` 教学图片；生成图不能承载必须精确读取的代码、命令或中文标签。
4. 只有官方图片或许可明确、允许再使用的外部图片才能本地保存；必须在视觉资产登记中记录来源 URL、作者或机构、许可证和访问日期。

所有项目图片保存在 `docs/public/images/go/`，使用稳定的英文文件名、中文替代文本和紧邻图片的解释。生成图片同时登记最终 prompt、生成方式和生成日期；真实截图登记复现命令、页面地址和视口。

## 7. 可运行示例工程

### 7.1 目录与边界

新增 `examples/go-task-api`，使用按业务能力分包、平台能力集中管理的结构：

```text
examples/go-task-api/
├─ cmd/
│  ├─ api/
│  ├─ migrate/
│  └─ healthcheck/
├─ internal/
│  ├─ app/
│  ├─ platform/
│  │  ├─ database/
│  │  └─ httpx/
│  ├─ user/
│  └─ task/
├─ migrations/
├─ go.mod
├─ go.sum
├─ Dockerfile
├─ compose.yaml
├─ .env.example
├─ README.md
├─ API_CONTRACT.md
└─ TROUBLESHOOTING.md
```

依赖方向固定为：

```text
cmd/api -> internal/app -> user/task -> repository contract
                            |             |
                            +------> platform/database
internal/platform/httpx 为 HTTP 适配层提供通用协议能力，不依赖业务包
```

关键约束：

- `main` 只读取配置、组装依赖、启动服务和处理关闭信号。
- Handler 只处理 HTTP 协议、JSON、路径参数、查询参数和状态码。
- Service 只处理业务规则、事务边界和跨仓储编排。
- Repository 只处理 SQL、结果扫描和数据库错误转换。
- interface 定义在使用方附近，只为真实替换边界服务，不为每个 struct 机械创建接口。
- `context.Context` 作为每层第一个参数传递，不保存在 struct 中，不用于传业务参数。
- `cmd/migrate` 固定支持 `up`、`down 1` 和 `version`；迁移文件使用成对的 `.up.sql`、`.down.sql`，不在 API 进程启动时自动改表。
- `cmd/healthcheck` 是只使用标准库的最小探针程序，供无 shell 的 distroless 运行镜像执行容器健康检查。

### 7.2 领域模型

示例包含 `users` 和 `tasks` 两张表：

- 用户字段：数据库生成的 `bigint` 主键、姓名、邮箱、状态、版本号、创建时间、更新时间。
- 任务字段：数据库生成的 `bigint` 主键、负责人、标题、描述、状态、截止时间、版本号、创建时间、更新时间。
- 邮箱使用 `lower(email)` 唯一索引，确保大小写变化不能绕过唯一约束。
- 任务状态由数据库约束限制为 `TODO`、`DOING`、`DONE`、`CANCELLED`。
- 用户状态限制为 `ACTIVE`、`DISABLED`。
- 更新用户和任务必须携带 `expectedVersion`，使用 `where id = $1 and version = $2` 实现乐观锁。
- 外键和删除行为在迁移中明确声明，示例不使用隐式级联删除。

所有表、字段、约束和索引都添加中文 SQL 注释，并在迁移文档中说明业务含义、空值语义、索引场景、升级和回滚方式。

### 7.3 HTTP API

示例至少提供以下路由：

| 方法 | 路由 | 行为 |
| --- | --- | --- |
| `GET` | `/health/live` | 进程存活检查，不依赖数据库 |
| `GET` | `/health/ready` | 数据库就绪检查 |
| `GET` | `/api/users` | 分页和状态筛选 |
| `POST` | `/api/users` | 创建用户 |
| `GET` | `/api/users/{id}` | 获取用户详情 |
| `PATCH` | `/api/users/{id}/status` | 请求体携带 `status` 和 `expectedVersion` |
| `GET` | `/api/tasks` | 分页、负责人和状态筛选 |
| `POST` | `/api/tasks` | 创建任务 |
| `GET` | `/api/tasks/{id}` | 获取任务详情 |
| `PUT` | `/api/tasks/{id}` | 请求体携带完整可编辑字段和 `expectedVersion` |
| `PATCH` | `/api/tasks/{id}/status` | 请求体携带 `status` 和 `expectedVersion` |
| `DELETE` | `/api/tasks/{id}?expectedVersion=3` | 查询参数携带版本号后删除任务 |

示例不包含认证，因此本地 Compose 端口只绑定 `127.0.0.1`，README 明确说明它不是可直接暴露到公网的生产 API。

### 7.4 响应和错误契约

成功响应统一包含 `success`、`data` 和 `requestId`；列表响应的 `data` 中包含 `items`、`page`、`pageSize` 和 `total`。服务同时在 `X-Request-ID` 响应头中返回相同值；调用方传入合法的 `X-Request-ID` 时复用，否则使用 `crypto/rand` 生成不可预测标识。

错误响应固定为：

```json
{
  "success": false,
  "error": {
    "code": "TASK_VERSION_CONFLICT",
    "message": "任务已被其他请求修改，请刷新后重试",
    "fields": {}
  },
  "requestId": "req_01..."
}
```

错误映射至少覆盖：

- JSON 语法错误、未知字段、请求体过大和多余 JSON 值：`400`。
- 字段校验失败：`422`，并返回 `fields`。
- 资源不存在：`404`。
- 邮箱冲突和版本冲突：`409`。
- 不支持的方法：`405`，包含 `Allow` 响应头。
- 不支持的媒体类型：`415`。
- 服务设置的请求截止时间到期且响应尚未提交：`504`；客户端主动断开时停止下游工作并记录取消原因，不再尝试写响应。
- 未分类内部错误：`500`，客户端不接收内部堆栈、SQL 或连接串。

中间件顺序固定为：request id、访问日志、recover、请求截止时间、内容类型与请求体限制、路由处理。截止时间中间件只向 request context 增加 deadline，不额外启动无法回收的 handler goroutine；Handler、Service 和 Repository 在收到 `context.DeadlineExceeded` 后交给统一错误映射返回 `504`。如果响应已经提交，只记录错误，不能再改写状态码。JSON 解码固定使用 `http.MaxBytesReader`、`DisallowUnknownFields`，并确认请求体中只有一个 JSON 值。

### 7.5 配置、超时和关闭

配置从环境变量读取，并在启动时一次性完成解析和校验：

- HTTP 地址、读取 Header 超时、读取超时、写入超时、空闲超时、单请求截止时间和整体关闭超时。
- 数据库 DSN、最大连接数、最大空闲连接数、连接最大生命周期和空闲生命周期。
- 日志级别和运行环境。

无效配置直接让进程启动失败，并指出字段名称，不回显数据库密码。

HTTP Server 必须显式设置超时；接收到 `SIGINT` 或 `SIGTERM` 后停止接收新请求，使用有截止时间的 context 执行 `Shutdown`，等待在途请求，最后关闭数据库。readiness 在关闭阶段先变为不可用，避免负载均衡继续发送请求。

Compose 中 API 和 PostgreSQL 的宿主机端口都只绑定到 `127.0.0.1`。迁移作为独立一次性服务在 API 启动前执行，API 不拥有生产表结构升级权限；健康检查通过 `cmd/healthcheck` 编译出的二进制访问容器内 `/health/live`。

## 8. 测试策略

测试按风险分层：

### 8.1 单元测试

- Service 使用表格测试覆盖正常流程、参数边界、状态机和仓储错误。
- 配置解析覆盖默认值、非法持续时间、非法数字和缺失 DSN。
- 错误映射覆盖已知业务错误和未分类错误。

### 8.2 HTTP 测试

- 使用 `httptest` 验证状态码、响应体、Header、未知字段、错误 JSON、请求体过大、非法或溢出的整数 ID、405 和 415。
- 验证 request id 在响应和日志上下文中保持一致。
- 验证 recover 不泄漏 panic 内容。

### 8.3 PostgreSQL 集成测试

- 集成测试使用 `integration` build tag 与普通测试分开，Testcontainers 启动 PostgreSQL 18，并执行真实迁移。
- 验证 `lower(email)` 唯一索引，而不是只测试 Service 层预检查。
- 验证事务成功提交、错误回滚和 context 取消。
- 验证两个并发更新只有一个版本号匹配，另一个得到版本冲突。
- 验证外键、状态约束和分页排序。

### 8.4 并发和健壮性

- `go test -race ./...` 必须通过。
- 至少提供一个针对 JSON 或筛选参数的 Fuzz 测试，保证异常输入不会 panic。
- 需要基准的代码先固定输入和预期结果，Benchmark 只用于教学，不宣称没有对照数据的性能提升。

## 9. 真实项目问题库

新增 `docs/projects/issues-go.md`，固定收录 16 个问题：

1. interface 不等于 nil，但内部指针为 nil。
2. 小 slice 长期引用大底层数组导致内存不释放。
3. map 并发读写和数据竞争。
4. goroutine 数持续增长且无法退出。
5. channel 重复关闭、无人接收或死锁。
6. context 没有传播，客户端取消后下游仍工作。
7. HTTP Client 没有超时或忘记关闭响应体。
8. HTTP Server 超时配置缺失，慢客户端占用资源。
9. `database/sql` 连接池耗尽，请求大量等待。
10. `rows`、`stmt` 或事务未关闭导致连接泄漏。
11. 并发更新丢失、重试造成重复写入。
12. JSON 的零值、空值、缺失字段和 `omitempty` 语义混乱。
13. 私有模块、代理、校验和 CI 环境不一致。
14. CPU、内存或 goroutine 异常时不会使用 pprof 和 trace。
15. 优雅关闭顺序错误，发布时丢请求或健康检查误报。
16. 容器中的证书、时区、信号、权限和配置问题。

每个问题都使用统一模板：现象、最小复现、错误方向、证据采集、根因、修复、回归测试、预防清单。至少 10 个问题配独立 Mermaid 图，不复用装饰性流程图凑数量。

## 10. 专项练习

新增 `docs/roadmap/go-practice.md`，固定包含 12 个递进练习：

1. 安装、模块和依赖校验。
2. 类型、集合、指针和方法。
3. interface、组合和错误链。
4. goroutine、channel 和退出条件。
5. context 和可取消操作。
6. 标准库路由、中间件和统一响应。
7. SQL、连接池和 Repository。
8. 事务、并发更新和幂等。
9. 表格测试、httptest、race 和 Fuzzing。
10. slog、指标、pprof 和诊断证据。
11. Docker、健康检查和优雅关闭。
12. 完整任务 API 综合交付。

每个练习包含目标、起始条件、任务步骤、限制条件、验证命令、通过标准、常见失败和进阶挑战。至少 7 个练习配图，综合练习直接复用真实示例工程的 API 契约和验收脚本。

## 11. 文档与示例的单一事实来源

`http-api-project-from-zero.md` 不再复制可能漂移的大段代码，而是通过 VitePress 代码导入语法从 `examples/go-task-api` 引用经过测试的文件或代码区域。

约束如下：

- 文档中的启动命令必须来自示例 README，并在真实环境执行。
- API 请求和响应必须与 `API_CONTRACT.md` 一致。
- 数据库字段和约束必须与迁移文件一致。
- 故障处理必须链接到 `TROUBLESHOOTING.md` 和 Go 问题库。
- 示例行为改变时，同一批修改必须更新文档、测试和契约。

## 12. 导航与成熟度更新

实施完成后同步更新：

- `docs/.vitepress/config.ts`：Go 侧边栏加入问题库和专项练习。
- `docs/technologies/index.md`：Go 卡片描述和入口。
- `docs/technologies/expansion-plan.md`：Go 模块完成状态和后续边界。
- `docs/contribute/module-status.md`：记录页面数、图示、问题库、练习和可运行示例。
- `docs/projects/real-world-issues.md`、`issues-backend.md`：增加 Go 专项入口，避免内容重复。
- `docs/roadmap/introduction.md`、`practice-labs.md`、`reading-guide.md`：增加 Go 学习路径。
- `README.md`：项目结构中增加 `examples/`，补充 Go 示例的独立验证命令。
- `.gitignore`：忽略 Go 示例二进制、覆盖率文件和本地环境文件，不忽略源码、迁移和 `go.sum`。
- `docs/contribute/visual-asset-register.md`：登记 Go 模块新增截图、生成图片或外部图片的用途、替代文本、来源与复现方式。

## 13. 验收标准

### 13.1 示例工程

以下命令必须在 `examples/go-task-api` 中通过：

```bash
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
go test -race ./...
go test -run=^$ -fuzz=FuzzDecodeJSON -fuzztime=10s ./internal/platform/httpx
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT
go build -o "$build_dir/api" ./cmd/api
go build -o "$build_dir/migrate" ./cmd/migrate
go build -o "$build_dir/healthcheck" ./cmd/healthcheck
```

启用 Docker 后还必须完成：

- `go test -tags=integration ./...` 启动 PostgreSQL 18，真实集成测试全部通过且没有因缺少 Docker 被跳过。
- `docker compose up --build` 后 live 和 ready 均返回成功。
- 创建用户、创建任务、分页查询、合法更新和旧版本更新冲突均符合契约。
- 发送终止信号后服务在关闭时限内退出，数据库和容器无残留测试资源。

### 13.2 文档站

以下命令必须通过：

```bash
npm run docs:check
npm run docs:build
git diff --check
```

浏览器必须在 `1440x900` 和 `390x844` 两种视口验证：

- Go 图解页 26 张 Mermaid 全部生成非空 SVG。
- Go 项目页、问题库和练习页没有 `.mermaid-diagram__error`。
- 本地图片全部成功加载、替代文本非空，截图与其对应的示例状态一致。
- 页面无整体横向滚动，代码块和复杂图只在自己的容器内滚动。
- 导航、目录、上一页和下一页链接可访问。
- 控制台无新增错误。
- 标题、表格、代码、提示块和图示在移动端不重叠。

### 13.3 内容质量

- Go 项目页中的所有导入代码都来自可通过测试的示例文件。
- 每篇被扩展的页面都有明确目标、前置知识、步骤、常见错误、验证方式和下一步。
- 图示之后必须有解释，不能用图替代关键文字。
- 问题库严格为 16 个问题，练习严格为 12 个练习。
- 修正现有 ER 图重复字段和超时示例重复函数声明等已知错误。
- 不出现未解释的占位符、伪命令、无来源版本断言或不可执行代码片段。

## 14. 风险与控制

### 14.1 内容规模过大

风险：一次改动覆盖文档、示例、测试和导航，容易出现文档与代码不一致。

控制：先完成可运行示例及其测试，再从示例导入代码；随后扩展图解、问题库和练习；最后统一更新导航和状态页。

### 14.2 标准库路线被误解为排斥框架

风险：读者可能认为生产项目不应使用框架或代码生成工具。

控制：在导览和项目选型章节明确标准库是学习基线，并提供 Gin、Chi、GORM、sqlc 的决策表和迁移时机。

### 14.3 数据库集成测试依赖 Docker

风险：无 Docker 环境时无法执行真实 PostgreSQL 测试。

控制：单元和 HTTP 测试不依赖 Docker；集成测试使用清晰的环境探测和跳过说明；最终交付环境必须实际运行 Docker 测试，不能只报告跳过结果。

### 14.4 Mermaid 在移动端可读性下降

风险：复杂图过宽或节点文本过多。

控制：按单一概念拆图，控制节点文本长度，图容器内部滚动，并在两种视口逐图检查 SVG 和页面宽度。

## 15. 回滚策略

- 保留现有 Go 页面路由，不通过删除旧路由完成重构。
- 示例工程位于独立目录，出现问题时可以按目录回退，不影响 VitePress 主题和其他技术模块。
- 导航更新只在目标页面存在且通过文档检查后加入。
- 数据库迁移同时提供向上和向下脚本；本地示例数据不视为生产数据。
- 不修改 Java 示例和已完成的 Java 文档行为，只在跨模块入口中增加 Go 链接。

## 16. 参考资料边界

实现和文档优先参考以下官方资料：

- Go 发布历史与 Go 1.26 Release Notes。
- Go Modules、Workspaces、数据库、事务、Fuzzing、race detector 和 pprof 官方文档。
- `net/http`、`database/sql`、`context`、`log/slog` 标准库文档。
- PostgreSQL 18 官方文档。
- pgx、golang-migrate 和 Testcontainers for Go 官方仓库及发布说明。

第三方博客可以用于收集真实问题场景，但不能作为版本、API 或安全结论的唯一依据。
