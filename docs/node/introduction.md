# Node.js 学习导览

## 适合谁看

适合已经掌握 JavaScript，准备进入后端 API、脚本工具、服务端应用或全栈开发的学习者。

Node.js 官方 Learn 文档说明，Node.js 是一个开源、跨平台的 JavaScript 运行时，可以用于服务器、Web 应用、命令行工具等项目。学习 Node.js 的重点不是“又一种语法”，而是理解 JavaScript 在浏览器之外如何运行。

## 你会学到什么

- Node.js 运行时和事件循环。
- npm、package.json 和模块化。
- HTTP API 的基本开发方式。
- 登录鉴权、会话保持、接口权限和 401/403 的边界。
- 数据库连接池、仓储层、事务和环境变量管理。
- Express/Fastify 中间件、路由、错误处理的基本概念。
- 测试策略、安全边界、日志、项目结构和部署边界。

## 学习顺序

<LearningPath :steps="[
  { title: '图解 Node.js 核心概念', description: '先用图理解运行时、事件循环、HTTP 请求链路、鉴权、数据库和排错路径。', link: '/node/visual-guide', badge: '图解' },
  { title: '运行时与事件循环', description: '理解 Node.js 如何运行 JavaScript，以及非阻塞 I/O 的基本模型。', link: '/node/runtime-event-loop', badge: '基础' },
  { title: '包管理与模块化', description: '掌握 package.json、npm scripts、ESM/CommonJS 和依赖管理。', link: '/node/package-modules', badge: '工程' },
  { title: 'HTTP API 开发', description: '学习路由、请求、响应、中间件和参数校验。', link: '/node/http-api', badge: '后端' },
  { title: '鉴权与会话', description: '区分认证、授权、会话、Token、刷新机制和按钮权限。', link: '/node/auth-session', badge: '权限' },
  { title: '数据库集成', description: '学习连接池、仓储层、参数化查询、事务和环境变量配置。', link: '/node/database-integration', badge: '数据' },
  { title: '错误处理与日志', description: '建立统一错误响应、日志记录和排错流程。', link: '/node/error-logging', badge: '质量' },
  { title: '测试策略', description: '用单元测试、接口测试和数据库测试守住服务行为。', link: '/node/testing', badge: '测试' },
  { title: 'Node.js 安全基础', description: '处理输入校验、注入、权限、上传、依赖和敏感配置风险。', link: '/node/security', badge: '安全' },
  { title: '项目结构与部署', description: '组织 API 项目目录、环境变量、启动命令和部署检查。', link: '/node/project-deployment', badge: '交付' },
  { title: 'Node 权限 API 从零到项目', description: '用用户、角色、菜单和按钮权限案例串联分层、鉴权、事务、日志和部署。', link: '/node/permission-api-project', badge: '实战' },
  { title: '常见问题', description: '排查端口占用、异步错误、环境变量、跨域和进程崩溃。', link: '/node/troubleshooting', badge: '排错' }
]" />

## Node.js 能做什么

| 场景 | 示例 |
| --- | --- |
| API 服务 | 登录、用户管理、订单接口 |
| BFF | 给前端聚合多个后端接口 |
| 命令行工具 | 代码生成、文件处理、自动化脚本 |
| 构建工具 | Vite、Webpack、ESLint 都运行在 Node 上 |
| 实时服务 | WebSocket、消息推送 |

## 学习建议

前端开发者学 Node.js，优先从 API 服务开始：

```text
HTTP 基础
↓
路由
↓
请求参数
↓
响应结构
↓
鉴权与权限
↓
数据库与事务
↓
错误处理与日志
↓
测试与安全
↓
部署
```

不要一开始就同时学微服务、消息队列、Kubernetes。先写出一个稳定 API 服务。

## 推荐项目

如果你已经学完 HTTP、鉴权和数据库，继续做 [Node 权限 API 从零到项目](/node/permission-api-project)。这个项目会把用户、角色、菜单、按钮权限、事务、审计日志和错误处理串起来。

## 下一步

从 [图解 Node.js 核心概念](/node/visual-guide) 开始，再进入 [运行时与事件循环](/node/runtime-event-loop)。
