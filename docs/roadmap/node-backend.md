# Node 后端工程师路线

## 适合谁看

适合已经会 JavaScript 或 TypeScript，想开始写后端 API、权限服务、数据接口和部署服务的人。也适合前端开发者向全栈方向扩展。

Node 后端路线的重点不是只学 Express 或 Fastify，而是理解后端工程的完整链路：HTTP、路由、参数校验、业务服务、数据库、日志、错误处理、部署和排错。

<LearningPath :steps="[
  { title: 'JavaScript 与 TypeScript', description: '掌握异步、模块化、类型建模和运行时校验边界。', link: '/javascript/introduction', badge: '语言' },
  { title: '浏览器与 HTTP', description: '理解 URL、状态码、请求头、跨域、登录态和缓存。', link: '/browser/http-request', badge: '基础' },
  { title: 'Node.js 运行时', description: '理解事件循环、包管理、模块系统和服务运行方式。', link: '/node/introduction', badge: '核心' },
  { title: 'HTTP API 开发', description: '设计路由、参数校验、错误响应、健康检查和分层结构。', link: '/node/http-api', badge: 'API' },
  { title: '数据库', description: '掌握 MySQL、PostgreSQL、Redis、建模、索引和事务。', link: '/database/introduction', badge: '数据' },
  { title: '错误与日志', description: '建立统一错误处理、结构化日志和 requestId 追踪。', link: '/node/error-logging', badge: '质量' },
  { title: '部署与运维', description: '使用 Nginx、Docker、CI/CD、健康检查和回滚策略。', link: '/devops/introduction', badge: '交付' }
]" />

## 学习节奏

先写一个最小 API 服务，再逐步加真实能力：

1. `/health` 健康检查。
2. 用户列表查询。
3. 用户新增和参数校验。
4. 登录和权限校验。
5. 数据库连接和事务。
6. 错误日志和 requestId。
7. Docker 和部署。

## 阶段验收

| 阶段 | 能力结果 |
| --- | --- |
| Node 基础 | 能解释事件循环和异步 I/O，不在接口里写长时间同步任务 |
| API 开发 | 能设计路由、校验参数、返回统一错误结构 |
| 数据库 | 能设计表、写迁移、建立索引、处理事务 |
| 质量保障 | 能记录日志、定位慢接口、处理异常 |
| 部署上线 | 能用 Docker 或进程管理器运行服务并配置反向代理 |

## 实际项目建议

推荐做一个“用户权限 API”：

- 用户、角色、权限表。
- 登录接口。
- 用户列表和筛选。
- 角色授权。
- 操作日志。
- 健康检查。
- Docker Compose 启动数据库和服务。

这个项目能同时练到 Node、数据库、权限、日志和部署。

## 项目与专项闭环

不要把路线停在“看完章节”。按下面顺序完成可运行项目、故障注入和复盘：

1. 用 [Node 权限 API 从零到项目](/node/permission-api-project) 完成 Fastify、TypeScript、PostgreSQL、会话、RBAC、测试和容器交付。
2. 用 [Redis 缓存与 BullMQ 队列项目](/node/cache-queue-project) 处理缓存一致性、队列重试、幂等和失败任务。
3. 完成 [Node.js 专项练习](/roadmap/node-practice) 的 12 个运行时练习，并保存命令、日志、指标和回归测试。
4. 使用 [Node.js 真实项目问题库](/projects/issues-node) 复盘模块、事件循环、Stream、进程和多实例故障。

接口契约、401/403、事务和慢 SQL 等跨语言问题继续查 [后端接口与服务问题](/projects/issues-backend) 和 [数据库与缓存问题](/projects/issues-database)。

## 常见误区

### 把后端写成接口转发层

只把前端参数转给数据库，不做校验、权限、错误处理和日志，这样后续很难维护。

### 忽略数据库设计

后端工程能力很大一部分体现在数据建模、事务和查询优化上。不要只关注路由框架。

### 本地能跑就算完成

后端服务必须考虑生产运行：环境变量、端口、日志、健康检查、自动重启和回滚。

## 下一步学习

继续学习 [Node.js 学习导览](/node/introduction)，然后进入 [数据库学习导览](/database/introduction) 和 [Node.js 专项练习](/roadmap/node-practice)。

如果你的目标是企业后台和 Spring 生态，切到 [Java 学习导览](/java/introduction)。如果你的目标是云原生服务、基础设施工具或高并发网关，切到 [Go 学习导览](/go/introduction)。
