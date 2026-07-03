# 前端技术图谱

## 图谱目标

前端技术很多，容易学成碎片。这个图谱把前端能力分成 8 层，帮助你判断自己缺哪一层。

```text
页面基础
↓
编程语言
↓
框架开发
↓
状态与数据
↓
工程化
↓
质量保障
↓
部署上线
↓
架构与协作
```

## 1. 页面基础

| 技术 | 学习目标 |
| --- | --- |
| HTML | 语义结构、表单、可访问性基础 |
| CSS | 盒模型、Flex、Grid、响应式、样式隔离 |
| 浏览器 | HTTP、跨域、缓存、存储、渲染、安全限制 |

当前文档入口：

- [HTML 与 CSS](/frontend/html-css)
- [CSS 学习导览](/css/introduction)
- [浏览器学习导览](/browser/introduction)

## 2. 编程语言

| 技术 | 学习目标 |
| --- | --- |
| JavaScript | 数据、函数、异步、模块化 |
| TypeScript | 类型建模、泛型、工程集成 |

当前文档入口：

- [JavaScript 基础](/javascript/fundamentals)
- [TypeScript 学习导览](/typescript/introduction)

## 3. 框架开发

| 技术 | 学习目标 |
| --- | --- |
| Vue | 响应式、组件、路由、状态、表单 |
| React | 组件、Hooks、状态、路由、服务端渲染 |
| 框架生态 | Nuxt、Next、组件库、工具链 |

当前已完成：

- [Vue 学习导览](/vue/introduction)

后续扩展：

- React 模块。
- Nuxt 模块。
- Next.js 模块。

## 4. 状态与数据

| 能力 | 典型技术 |
| --- | --- |
| 局部状态 | `ref`、`reactive`、组件 state |
| 全局状态 | Pinia、Redux、Zustand |
| 请求管理 | axios、fetch、TanStack Query |
| 权限与菜单 | 路由守卫、权限码、动态路由 |

当前文档入口：

- [Pinia 状态管理](/vue/pinia)
- [请求与接口封装](/vue/request)
- [权限与菜单](/vue/permission)

## 5. 工程化

| 能力 | 学习目标 |
| --- | --- |
| 构建工具 | Vite、Webpack 基础 |
| 代码规范 | ESLint、Prettier、TypeScript |
| 环境配置 | 多环境变量、代理、构建模式 |
| 组件库 | 组件封装、设计系统、文档示例 |

当前文档入口：

- [Vite 工程基础](/engineering/vite)
- [代码规范](/engineering/eslint-prettier)
- [环境配置](/engineering/env-config)
- [组件库实战](/projects/component-library)

## 6. 质量保障

| 能力 | 学习目标 |
| --- | --- |
| 单元测试 | 工具函数、权限判断、数据转换 |
| 组件测试 | 组件输入输出和交互 |
| E2E 测试 | 登录、核心业务流程 |
| 性能排查 | Network、Performance、Vue DevTools |

当前文档入口：

- [测试策略](/vue/testing)
- [性能优化](/vue/performance)

## 7. 部署上线

| 能力 | 学习目标 |
| --- | --- |
| 静态部署 | dist、HTTP 服务、缓存 |
| Nginx | history fallback、反向代理、缓存策略 |
| CI/CD | 自动构建、发布、回滚 |
| 监控 | 错误日志、性能指标、用户反馈 |

当前文档入口：

- [构建与部署](/engineering/build-deploy)
- [真实项目问题库](/projects/real-world-issues)

## 8. 架构与协作

| 能力 | 学习目标 |
| --- | --- |
| 分层 | UI、状态、请求、业务流程分离 |
| 文档 | README、模块说明、变更记录 |
| 设计系统 | 组件库、主题、样式边界 |
| 团队协作 | 规范、评审、测试、发布流程 |

当前文档入口：

- [最佳实践](/vue/best-practices)
- [贡献指南](/contribute/)

## 学习建议

不要同时展开所有技术。建议按项目驱动学习：

1. 先完成一个 Vue Admin 项目。
2. 遇到问题时沉淀到问题库。
3. 项目稳定后再抽组件库。
4. 再扩展 React、Node、数据库等横向能力。
