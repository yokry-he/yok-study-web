# 全栈工程师路线

## 适合谁看

适合已经能写前端页面，想独立完成一个完整 Web 项目的人。全栈不是“什么都浅学一点”，而是能把前端、后端、数据库、部署和问题排查串起来。

全栈路线的目标是：能从一个业务需求出发，完成页面、接口、数据表、权限、部署和回滚。

<LearningPath :steps="[
  { title: '前端页面与框架', description: '使用 Vue 或 React 完成页面、组件、路由、状态和表单。', link: '/roadmap/vue-frontend', badge: '前端' },
  { title: '浏览器与网络', description: '理解跨域、登录态、缓存、存储和性能排查。', link: '/browser/introduction', badge: '运行环境' },
  { title: '后端 API', description: '选择 Node.js、Java 或 Go 设计路由、业务服务、错误处理和日志。', link: '/roadmap/node-backend', badge: '后端' },
  { title: '数据库建模', description: '设计表结构、约束、索引、事务和迁移。', link: '/database/modeling', badge: '数据' },
  { title: 'DevOps 与部署', description: '完成 Nginx、Docker、CI/CD、版本发布和回滚。', link: '/devops/introduction', badge: '交付' },
  { title: '真实问题库', description: '把联调、上线、缓存、权限、慢查询等问题沉淀成可复用经验。', link: '/projects/real-world-issues', badge: '排错' }
]" />

## 推荐项目

推荐用“后台管理系统”作为全栈练习项目：

```text
Vue / React 前端
↓
Node.js / Java / Go API
↓
PostgreSQL 或 MySQL
↓
Redis 缓存
↓
Nginx 反向代理
↓
Docker Compose
↓
CI/CD 发布
```

这个项目可以覆盖：

- 登录。
- 用户管理。
- 角色权限。
- 动态菜单。
- 表格筛选。
- 表单提交。
- 操作日志。
- 部署上线。

## 阶段验收

| 阶段 | 能力结果 |
| --- | --- |
| 前端 | 能拆组件、管理状态、处理请求和权限 |
| 后端 | 能写 API、校验参数、处理业务错误 |
| 数据库 | 能说明表设计、索引原因和迁移风险 |
| 部署 | 能构建、部署、代理、回滚和验证 |
| 排错 | 能定位前端、后端、数据库、Nginx 或缓存层问题 |

## 学习顺序建议

不要同时开太多技术。建议：

1. 先完成前端页面和静态 mock。
2. 再接 Node.js、Java 或 Go API；如果选择 Java，优先用 [Spring Boot 从零到项目落地](/java/spring-boot-project-from-zero) 做出用户角色 API。
3. 再接数据库。
4. 再补权限和日志。
5. 再容器化和部署。
6. 最后补测试和问题库。

## 常见误区

### 每层都只学框架

前端不只是 Vue，后端不只是 Express，数据库不只是 CRUD，部署不只是复制文件。全栈能力来自跨层理解。

### 没有边界

全栈项目也要分层：UI、状态、请求、业务服务、数据访问、数据库、部署配置不能混在一起。

### 没有上线意识

项目不是本地跑通就结束。要能说明环境变量、构建产物、反向代理、缓存策略和回滚方式。

## 下一步学习

如果你还没完成前端项目，先走 [Vue 前端工程师路线](/roadmap/vue-frontend)。如果你已经能做前端页面，继续进入 [Node 后端工程师路线](/roadmap/node-backend)。
