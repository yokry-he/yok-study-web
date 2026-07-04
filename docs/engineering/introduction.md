# 前端工程化学习导览

## 这个模块解决什么

前端工程化解决的是“项目如何长期稳定开发、协作、构建、测试和交付”的问题。

如果只会写页面，项目小的时候还能靠个人经验推进；但项目一旦进入多人协作、多个环境、频繁发布、组件复用、权限控制和线上排错阶段，就必须依赖工程化能力。

工程化不是堆工具，而是把下面这些事情变成稳定流程：

- 项目怎么创建和组织目录。
- 代码怎么保持一致风格。
- 环境变量怎么管理。
- 依赖怎么升级和锁定。
- 测试怎么覆盖核心行为。
- 多包项目怎么拆分和复用。
- 构建产物怎么部署和回滚。
- 出问题时怎么快速定位。

## 适合谁看

适合已经能写 Vue、React 或普通前端页面，但准备进入真实项目开发的学习者。

你可能正在遇到这些问题：

- 本地能跑，测试环境或生产环境不正常。
- 依赖升级后项目突然启动失败。
- 团队每个人格式化结果不一样。
- 构建偶尔失败，但不知道从哪里查。
- 项目越来越大，组件、接口、工具函数到处散落。
- 想做组件库、工具库或多应用项目，但不知道 Monorepo 怎么组织。

## 学习顺序

推荐按这个顺序学习：

```text
Vite 工程基础
↓
图解前端工程化核心概念
↓
前端工程化从零到项目落地
↓
代码规范
↓
环境配置
↓
依赖管理
↓
测试策略
↓
Monorepo
↓
组件库工程从零到项目
↓
构建与部署
↓
工程化常见问题
```

不要一开始就追复杂工具链。先让项目具备最基本的稳定性，再逐步加入测试、发布、包管理和多包治理。

## 工程化能力地图

| 能力 | 解决的问题 | 对应文档 |
| --- | --- | --- |
| 工程化全局图 | 从开发、质量门禁、构建、部署到回滚的整体链路 | [图解前端工程化核心概念](/engineering/visual-guide) |
| 从零到项目 | 把脚本、目录、环境、规范、测试、CI、构建、发布和回滚串成可交付项目 | [前端工程化从零到项目落地](/engineering/project-from-zero) |
| 项目启动和构建 | 项目如何跑起来、如何输出生产产物 | [Vite 工程基础](/engineering/vite) |
| 代码规范 | 多人协作时风格和低级错误如何控制 | [代码规范](/engineering/eslint-prettier) |
| 环境配置 | 本地、测试、预发、生产如何隔离 | [环境配置](/engineering/env-config) |
| 依赖管理 | npm、pnpm、lockfile、升级和安全如何处理 | [依赖管理](/engineering/package-management) |
| 测试策略 | 单元测试、组件测试、端到端测试怎么落地 | [测试策略](/engineering/testing) |
| Monorepo | 多应用、多包、组件库、工具库如何组织 | [Monorepo 项目组织](/engineering/monorepo) |
| 组件库工程 | 如何组织组件包、主题 token、文档站、测试、构建和版本发布 | [组件库工程从零到项目](/engineering/component-library-project) |
| 构建部署 | 静态资源、路由、缓存、回滚如何处理 | [构建与部署](/engineering/build-deploy) |
| 包体积分析 | 如何定位大依赖、拆包、懒加载和 chunk 警告 | [包体积分析](/engineering/bundle-analysis) |
| 模块联邦 | 多团队、多应用如何做微前端和独立发布 | [模块联邦与微前端](/engineering/module-federation) |
| 工程性能 | 安装、启动、HMR、类型检查、测试和 CI 如何提速 | [工程性能优化](/engineering/performance-optimization) |
| 项目问题库 | 安装、环境、CI、构建、部署缓存、依赖升级、Monorepo 和回滚怎么排查 | [前端工程化真实项目问题库](/projects/issues-engineering) |
| 排错 | 启动、构建、样式、依赖、部署问题怎么查 | [工程化常见问题](/engineering/troubleshooting) |

## 一个真实项目的工程化最低配置

一个可交付的前端项目，至少应该具备：

```text
package.json              统一脚本入口
vite.config.ts            构建、代理、路径别名
tsconfig.json             类型检查和路径映射
.env.*                    环境配置
eslint.config.*           代码质量规则
prettier.config.*         格式化规则
src/config/               应用配置集中读取
src/api/                  接口请求
src/services/             业务流程
src/stores/               全局状态
src/router/               路由和守卫
src/styles/               全局样式边界
tests/ 或 __tests__/       测试用例
README.md                 启动、构建、部署说明
```

脚本建议保持简单明确：

```json
{
  "scripts": {
    "dev": "vite",
    "type-check": "vue-tsc --noEmit",
    "lint": "eslint .",
    "test": "vitest",
    "build": "vite build",
    "preview": "vite preview"
  }
}
```

初学者不要把脚本命名得太花。团队里每个人都应该能一眼看懂这些命令的作用。

## 工程化和业务代码的边界

工程化配置服务于业务代码，但不要和业务逻辑混在一起。

推荐边界：

| 内容 | 应该放哪里 | 不应该放哪里 |
| --- | --- | --- |
| API 前缀 | `src/config/app.ts` | 每个页面里手写 |
| 请求拦截 | `src/api/request.ts` | 业务组件里重复写 |
| 权限码 | `src/constants/permissions.ts` | 模板里写字符串 |
| 环境变量 | `.env.*` 和配置模块 | 组件里散落读取 |
| 构建配置 | `vite.config.ts` | 业务代码里判断打包环境 |
| 测试工具 | `tests/setup.ts` | 每个测试重复初始化 |

如果一个配置会被多个地方使用，就应该集中管理。

## 常见误区

### 误区 1：工程化就是装很多插件

插件越多，项目越难升级。只有当问题真实存在，并且团队愿意维护配置时，才引入新工具。

### 误区 2：本地能跑就可以上线

本地开发服务器和生产部署完全不同。Vite proxy、热更新、未压缩资源、开发环境变量都不能代表生产环境。

### 误区 3：没有测试也能靠人工点一遍

人工验收适合发现体验问题，不适合稳定覆盖核心逻辑。权限判断、表单校验、数据转换、工具函数和请求错误处理应该逐步测试化。

### 误区 4：依赖能升级就直接升级

依赖升级可能带来构建变化、类型变化、浏览器兼容变化和样式变化。真实项目需要锁版本、看 changelog、分批升级和回滚方案。

## 阶段目标

### 初级阶段

你应该能做到：

- 使用 Vite 创建和运行项目。
- 看懂 `vite.config.ts`。
- 会配置路径别名和环境变量。
- 能解释本地代理和生产代理的区别。
- 能执行 `build` 并预览产物。

### 中级阶段

你应该能做到：

- 制定目录规范。
- 使用 ESLint、Prettier 和 TypeScript 检查项目。
- 把环境配置集中管理。
- 编写基础单元测试和组件测试。
- 排查常见启动、构建和部署问题。

### 项目阶段

你应该能做到：

- 为团队项目设计脚本和质量门禁。
- 管理依赖升级和 lockfile。
- 组织组件库、工具库和多应用 Monorepo。
- 为发布、回滚和故障复盘提供文档。
- 让新人按 README 就能启动和参与开发。

## 下一步学习

从 [图解前端工程化核心概念](/engineering/visual-guide) 开始，再进入 [前端工程化从零到项目落地](/engineering/project-from-zero)，把 [Vite 工程基础](/engineering/vite)、[代码规范](/engineering/eslint-prettier)、[环境配置](/engineering/env-config)、[依赖管理](/engineering/package-management) 和 [测试策略](/engineering/testing) 串成项目流程。
