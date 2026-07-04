# 内容扩展路线

## 为什么需要扩展路线

这个文档站的最终目标是覆盖程序员常用技术，但不能一开始就铺成低质量百科。正确方式是：

1. 先做深一个方向。
2. 建立稳定模板。
3. 用模板复制到新技术。
4. 持续补真实项目问题。

当前已经完成第一步：Vue 前端方向已经形成较完整样板。

## 阶段一：前端 Vue 方向

当前状态：已基本完成。

已覆盖：

- 前端基础。
- JavaScript 图解核心概念、语言基础、原型链、DOM 事件、正则表达式、事件循环、错误处理、内存管理、模块化和任务看板从零到项目。
- CSS 图解核心概念、盒模型、Flex/Grid、响应式、动画、可访问性、设计 token、样式架构和 CSS 从零到项目落地。
- TypeScript 图解核心概念、基础、接口、泛型、类型收窄、工具类型边界、tsconfig 工程配置、Vue 集成和类型边界从零到项目。
- Vue 3 图解核心概念和核心章节。
- Vue Router。
- Pinia。
- 表单、请求、权限。
- Vue 从零到项目落地、JavaScript 项目落地实践、JavaScript 任务看板从零到项目、TypeScript 项目类型边界实践和 TypeScript 类型边界从零到项目。
- 工程化图解核心概念、Vite、规范、环境、依赖、测试、Monorepo、组件库工程从零到项目、构建部署、包体积、模块联邦和工程性能。
- 部署上线。
- 浏览器图解核心概念、安全、Service Worker/PWA、常用 Web API。
- WebSocket、WebRTC、Web Components。
- WebAssembly、WebGPU、浏览器自动化调试。
- Vue Admin 实战。
- 阅读顺序与使用方法。
- 学习路径练习包。
- 项目交付检查清单。
- 真实项目问题库。
- 真实项目问题库已经拆分为前端、后端、数据库、部署和 AI 工程分类，并补充 React 重复请求、组件库升级错位、权限半更新、401/403 混乱、AI 引用不可信和 Agent 工具循环等项目问题。
- 速查手册已经覆盖 Vue、Router、Pinia、Vite、JavaScript、TypeScript、CSS、正则、Node、Java、Go、HTTP、Git、Linux、常用命令、调试工具、Docker、Nginx、SQL 和 Redis。

后续微调：

- CSS 已补从零到项目落地、容器查询和样式系统练习，后续补打印样式、更多组件样式案例和大型项目主题治理。
- TypeScript 后续补声明文件、发布类型、库类型设计和更多复杂业务案例。
- JavaScript 后续补 DOM API 深入、日期时间、数字精度和更多项目题。
- 浏览器模块继续补浏览器测试策略、WebView 兼容和更多性能案例。
- 前端工程化已经补充包体积分析、工程性能和模块联邦。
- 增加更多 Vue Admin 案例页面。
- 真实项目问题库已经补充故障复盘、联调问题、权限案例、项目阶段任务、React/Node/AI/组件库工程高频问题。
- 学习路线已经补充阅读指南、图解学习地图和阶段练习包，帮助读者从“看懂技术关系”进入“做练习和复盘”。
- 速查手册已经补齐 Linux、Redis、正则、常用命令和调试工具。

## 阶段二：React 方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
react/
├─ introduction
├─ visual-guide
├─ quick-start
├─ component-jsx
├─ hooks-state
├─ effects
├─ router-structure
├─ project-admin
└─ troubleshooting
```

后续继续扩展：

- Next.js 与服务端渲染已经放入元框架模块。
- 后续补 React 生态里的 Remix、TanStack Router 或更完整状态库对比。

已覆盖：

- 图解组件树、props、state、Effect、服务端数据和排错路径。
- JSX、组件拆分、Hooks、Effect、表单、请求、Context、路由、性能和测试。
- React 管理台从零到项目：登录、路由、请求、表格、表单、权限、测试和部署说明。

## 阶段二补充：Nuxt / Next 元框架方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
meta-frameworks/
├─ introduction
├─ nuxt
├─ next
├─ routing-data
├─ deployment
├─ server-auth
├─ seo-metadata
├─ i18n
├─ content-site-case
└─ troubleshooting
```

已覆盖：

- 元框架适用场景和学习路线。
- Nuxt 项目实践。
- Next.js 项目实践。
- 文件路由、布局和数据获取。
- SSR、静态生成、Node server、Serverless 和 Edge 部署取舍。
- 服务端鉴权、登录态、接口保护和用户态缓存边界。
- SEO、metadata、Open Graph、sitemap 和结构化数据。
- 国际化、多语言路由、多语言 SEO 和翻译治理。
- 内容站案例：技术博客、官网、路由、内容模型和发布缓存。
- hydration、缓存、环境变量、部署后 404 等常见问题。

后续继续扩展：

- 更深入的性能优化、边缘渲染和图片优化。
- 内容管理系统接入案例。
- 多租户门户和会员内容案例。

## 阶段三：Node.js 后端方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
node/
├─ introduction
├─ visual-guide
├─ runtime-event-loop
├─ package-modules
├─ http-api
├─ auth-session
├─ database-integration
├─ error-logging
├─ testing
├─ security
├─ project-deployment
├─ permission-api-project
└─ troubleshooting
```

已覆盖：

- 图解运行时、事件循环、HTTP 请求链路、鉴权、数据库、错误日志和排错路径。
- 鉴权与会话。
- 数据库集成与事务。
- 测试策略。
- 安全基础。
- Node 权限 API 从零到项目：用户、角色、菜单、按钮权限、事务、审计日志和部署说明。

后续继续扩展：

- 缓存、队列与任务调度。
- 文件上传与对象存储。
- WebSocket 服务。
- 性能优化和压测。
- Express/Fastify 深入。

## 阶段三补充：Java 后端方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
java/
├─ introduction
├─ visual-guide
├─ setup-tooling
├─ syntax-oop
├─ collections-generics
├─ exceptions-logging
├─ streams-lambda
├─ concurrency-virtual-threads
├─ jvm-memory-gc
├─ spring-boot-api
├─ persistence-transaction
├─ testing-deployment
└─ troubleshooting
```

已覆盖：

- 图解 JDK、JVM、对象引用、调用栈、Spring 请求链路、事务、线程和排错路径。
- JDK、JVM、JRE、Maven、Gradle 和 IDE 工具链。
- Java 语法、面向对象、接口、record 和 sealed class。
- 集合、泛型、Optional、equals/hashCode 和常见集合选择。
- 异常、统一错误响应、结构化日志和脱敏。
- Stream、Lambda、分组、映射、批量查询和并行边界。
- 平台线程、线程池、CompletableFuture 和虚拟线程。
- JVM 内存、GC、线程 dump、heap dump 和类冲突排查。
- Spring Boot API、分层架构、配置、参数校验和请求链路。
- 数据库访问、事务、ORM、连接池和 N+1 查询。
- 测试、打包、部署、健康检查和上线验证。
- Java 速查和常见问题。

后续继续扩展：

- Spring Security 和权限体系。
- 消息队列、异步任务和事件驱动。
- 微服务治理、配置中心和注册发现。
- 更多企业后台、支付、审批流和报表案例。

## 阶段三补充：Go 后端方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
go/
├─ introduction
├─ visual-guide
├─ setup-modules
├─ syntax-types
├─ interfaces-composition
├─ errors-logging-config
├─ concurrency
├─ context-http
├─ database-transaction
├─ testing
├─ project-deployment
├─ performance
└─ troubleshooting
```

已覆盖：

- 图解模块构建、包边界、Handler-Service-Repository、goroutine 生命周期、channel、context、连接池和性能排查。
- Go 安装、go.mod、go.sum、workspace、replace 和 GOPRIVATE。
- 变量、零值、函数、结构体、方法、指针接收者和泛型。
- 小接口、隐式实现、组合和按业务域组织包。
- error 包装、日志、配置来源和启动校验。
- goroutine、channel、select、WaitGroup、数据竞争和并发控制。
- context、HTTP handler、中间件、请求超时和优雅关闭。
- database/sql、事务、连接池、rows 关闭和仓储层。
- 单元测试、表格测试、benchmark、fuzzing 和 race 检查。
- 项目结构、构建、Docker、部署检查和版本信息。
- pprof、goroutine profile、heap profile 和性能诊断。
- Go 速查和常见问题。

后续继续扩展：

- gRPC、Protobuf 和服务间通信。
- 云原生组件、Kubernetes Operator 和控制器模式。
- CLI 工具、后台任务和任务队列。
- 更多高并发接口、网关和基础设施案例。

## 阶段四：数据库方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
database/
├─ introduction
├─ visual-guide
├─ mysql
├─ postgresql
├─ redis
├─ project-practice
├─ modeling
├─ indexes
├─ transactions
├─ migration
├─ orm-practice
├─ backup-recovery
├─ security-audit
└─ troubleshooting
```

数据库模块必须特别重视注释和文档，所有表、字段、索引、约束、迁移原因都要写清楚。

已覆盖：

- 图解数据库访问链路、表关系、索引、事务锁、慢查询、缓存和排错路径。
- MySQL 和 InnoDB 项目实践。
- PostgreSQL 类型、约束、JSONB 和 EXPLAIN。
- Redis 缓存、数据结构、过期和淘汰策略。
- 后台权限系统数据层落地：表设计、索引、迁移、种子、事务、缓存和上线检查。
- 数据建模、表设计、字段注释和约束说明。
- 索引、复合索引、外键索引、N+1 和执行计划。
- 事务、锁、并发、短事务和幂等更新。
- 迁移、种子数据、版本治理和回滚风险。
- ORM 分层、字段选择、关联查询、事务和慢查询排查。
- 备份恢复、RPO、RTO、恢复演练和误删数据处理。
- 数据安全、最小权限、敏感字段、审计日志和脱敏治理。
- 常见数据库线上问题排查。

后续继续扩展：

- 分库分表和读写分离。
- PostgreSQL RLS 和多租户权限。
- 数据仓库、ETL 和分析型查询。
- 更多云数据库运维案例。

## 阶段五：DevOps 方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
devops/
├─ introduction
├─ visual-guide
├─ linux-shell
├─ nginx
├─ docker
├─ ci-cd
├─ project-deployment-practice
├─ deployment-strategy
├─ observability
├─ kubernetes-basics
├─ cloud-deployment
└─ troubleshooting
```

已覆盖：

- 图解线上请求链路、Nginx 反向代理、Docker 镜像分层、CI/CD、灰度发布和可观测性。
- Linux 常用命令。
- Nginx 部署前端和反向代理。
- Docker 镜像和容器。
- CI/CD 发布。
- 项目上线全流程：构建、Nginx、Docker Compose、缓存策略、发布验证、回滚和线上排查。
- 发布、回滚和环境治理。
- 可观测性：日志、指标、链路追踪、告警和发布观察。
- Kubernetes 基础：Pod、Deployment、Service、Ingress 和健康检查。
- 云服务、对象存储、CDN、云数据库和成本治理。
- 常见线上故障排查。

后续继续扩展：

- 发布审计和安全基线。
- GitOps 和环境漂移治理。
- Kubernetes 配置管理和灰度发布。
- 更多云成本优化案例。

## 阶段六：AI 工程方向

当前状态：已扩展为稳定模块。

已建立结构：

```text
ai-engineering/
├─ introduction
├─ visual-guide
├─ llm-api
├─ prompt-engineering
├─ structured-outputs-tools
├─ multimodal
├─ rag
├─ mcp-integration
├─ agents
├─ product-workflow
├─ evaluation
├─ doc-qa-project
├─ deployment
└─ troubleshooting
```

已覆盖：

- 图解模型调用、Prompt、结构化输出、工具调用、RAG、Agent、评测和上线治理。
- LLM API 调用。
- 提示词工程。
- 结构化输出和函数调用。
- 多模态：图片、语音、文件输入输出。
- RAG 检索增强生成。
- MCP 和企业内部工具集成。
- Agent 工作流和工具边界。
- AI 产品设计和人机协作流程。
- 评测和可观测性。
- AI 文档问答从零到项目：导入、切分、检索、权限过滤、回答、引用、评测和上线治理。
- 成本、延迟和安全。
- 常见 AI 工程问题排查。

后续继续扩展：

- 更完整的 eval 自动化实践。
- 更复杂的多 Agent 协作案例。
- 企业知识权限、审计和数据治理案例。

## 每个模块的完成标准

一个模块达到可发布质量，需要满足：

- 至少 8 到 12 篇核心文档。
- 有快速开始。
- 有实际项目问题。
- 有最佳实践。
- 有速查手册。
- 有下一步学习指引。
- 能通过全站构建。
- 关键页面移动端无横向溢出。

## 当前推荐下一步

继续完善顺序建议：

1. 继续补从零到项目章节和复杂案例，例如 React 管理台、Node 权限 API、AI 文档问答、工程化组件库项目。
2. 继续为 stable 模块补深入专题，例如 Java Spring Security、Go gRPC、Node 缓存队列、数据库分库分表、DevOps GitOps 和 AI eval 自动化。
3. 为实战项目继续补更多完整项目章节，例如销售风险指标治理预算执行复盘、渠道策略标准灾备演练、生产安全整改推荐模型治理、售后知识认证服务商复评和更多行业项目案例。
4. 按读者反馈继续补充速查手册的浏览器命令、包管理器、云平台命令和团队协作命令。
