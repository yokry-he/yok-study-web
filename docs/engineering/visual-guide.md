# 图解前端工程化核心概念

## 适合谁看

适合已经能写页面，但对 Vite、环境变量、依赖管理、代码规范、测试、构建、部署和 Monorepo 还没有系统认识的人。

工程化的目标不是装很多工具，而是让项目在多人协作、长期迭代和频繁发布时仍然稳定。

## 你会学到什么

- 一个前端项目从开发到上线会经过哪些环节。
- Vite、TypeScript、ESLint、测试、构建、部署分别解决什么问题。
- 环境变量、请求配置和发布配置如何区分。
- 依赖和 lockfile 为什么影响稳定性。
- 构建产物、缓存、回滚和线上排错如何串起来。

## 工程化全流程

```mermaid
flowchart LR
  A[创建项目] --> B[目录约定]
  B --> C[代码规范]
  C --> D[环境配置]
  D --> E[开发调试]
  E --> F[测试]
  F --> G[构建]
  G --> H[部署]
  H --> I[监控和回滚]
```

每一步都应该有明确脚本和文档。新人不应该靠问人才能启动项目。

## Vite 开发到构建

```mermaid
flowchart TD
  A[源码] --> B[Vite Dev Server]
  B --> C[浏览器按需加载模块]
  C --> D[HMR 更新]
  A --> E[Vite Build]
  E --> F[Rollup 打包]
  F --> G[dist 静态产物]
```

本地开发服务器和生产构建不是一回事：

- dev server 提供 HMR、开发代理、未压缩模块。
- build 输出静态资源、hash 文件、压缩产物。
- 本地 proxy 不能代表生产 Nginx 代理。

## 环境配置链路

```mermaid
flowchart TD
  A[.env.development] --> D[配置读取模块]
  B[.env.staging] --> D
  C[.env.production] --> D
  D --> E[request baseURL]
  D --> F[feature flags]
  D --> G[build base]
```

前端环境变量要区分：

| 类型 | 示例 | 生效时机 |
| --- | --- | --- |
| 构建时变量 | `VITE_API_BASE_URL` | 构建时写入产物 |
| 运行时配置 | `/runtime-config.js` | 浏览器加载时读取 |
| 服务端环境变量 | `DATABASE_URL` | 后端进程运行时读取 |

不要在业务组件里到处读取环境变量。集中到配置模块里。

## 代码质量门禁

```mermaid
flowchart TD
  A[提交代码] --> B[格式化]
  B --> C[ESLint]
  C --> D[TypeScript 检查]
  D --> E[单元测试]
  E --> F[构建]
  F --> G[允许合并或发布]
```

质量门禁要从低成本开始：

- 格式化统一风格。
- ESLint 发现低级错误。
- TypeScript 防止类型边界漂移。
- 测试覆盖核心函数和关键组件。
- 构建确认生产产物能生成。

## 依赖管理

```mermaid
flowchart TD
  A[package.json] --> B[版本范围]
  B --> C[lockfile]
  C --> D[确定安装结果]
  D --> E[CI 和本地一致]
  E --> F[可复现构建]
```

真实项目不要随意删除 lockfile。依赖升级要能回答：

- 升级了哪些包。
- 为什么升级。
- 是否有 breaking changes。
- 是否通过测试和构建。
- 如何回滚。

## Monorepo 关系

```mermaid
flowchart TD
  A[workspace root] --> B[apps/admin]
  A --> C[apps/docs]
  A --> D[packages/ui]
  A --> E[packages/utils]
  D --> B
  E --> B
  E --> C
```

Monorepo 适合多应用、多包、组件库和工具库共存。不要为了“看起来高级”过早引入。只有当复用和协作问题真实存在时，再让项目进入多包治理。

## 构建部署和回滚

```mermaid
flowchart TD
  A[npm run build] --> B[dist]
  B --> C[上传 assets]
  C --> D[切换 index.html 或 current]
  D --> E[验证核心路由]
  E --> F{是否异常}
  F -- 否 --> G[观察日志]
  F -- 是 --> H[回滚上一版本]
```

前端上线重点：

- `index.html` 不长期强缓存。
- assets 带 hash 后可长期缓存。
- 二级路由刷新不 404。
- 发布有版本号。
- 回滚能找到上一版本。

## 工程化排错路径

```mermaid
flowchart TD
  A[工程问题] --> B{启动失败}
  B -- 是 --> C[依赖、Node 版本、环境变量]
  B -- 否 --> D{构建失败}
  D -- 是 --> E[类型、导入路径、构建配置]
  D -- 否 --> F{部署异常}
  F -- 是 --> G[base、缓存、Nginx、代理]
  F -- 否 --> H[性能、测试、依赖升级]
```

## 实际项目常见问题

### 问题 1：本地能跑，CI 构建失败

检查 Node 版本、包管理器版本、lockfile、大小写路径、环境变量和是否依赖本地未提交文件。

### 问题 2：构建后接口地址错误

确认是构建时变量还是运行时配置。Vite 变量构建后已经写入静态资源。

### 问题 3：依赖升级后页面样式异常

先看组件库 changelog、全局 CSS、主题 token 和 lockfile diff，不要直接写更高优先级覆盖。

## 下一步学习

继续学习 [Vite 工程基础](/engineering/vite)、[代码规范](/engineering/eslint-prettier)、[环境配置](/engineering/env-config) 和 [构建与部署](/engineering/build-deploy)。
