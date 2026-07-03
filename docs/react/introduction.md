# React 学习导览

## 适合谁看

适合已经掌握 JavaScript 和基础前端，准备学习 React 组件开发、Hooks、路由和真实项目组织方式的学习者。

React 官方文档现在强调“用组件描述 UI、用状态驱动界面、用 Hooks 使用 React 能力”。学习 React 的重点不是背 API，而是理解组件渲染、状态更新和副作用边界。

## 你会学到什么

- JSX 和组件如何表达 UI。
- `useState` 如何管理局部状态。
- `useEffect` 什么时候该用、什么时候不该用。
- React Router 如何组织页面。
- React 项目中常见的重复渲染、Effect 循环、状态放错位置如何处理。

## 学习顺序

<LearningPath :steps="[
  { title: '图解 React 核心概念', description: '先用图理解组件树、props、state、Effect、服务端数据和排错路径。', link: '/react/visual-guide', badge: '图解' },
  { title: '快速开始', description: '创建 React 项目，理解组件、状态和事件。', link: '/react/quick-start', badge: '入门' },
  { title: '组件与 JSX', description: '学习 JSX、props、列表渲染、条件渲染和组件拆分。', link: '/react/component-jsx', badge: '核心' },
  { title: 'Hooks 与状态', description: '掌握 useState、自定义 Hook 和状态提升。', link: '/react/hooks-state', badge: '核心' },
  { title: 'Effect 与副作用', description: '理解同步外部系统、请求、订阅和清理逻辑。', link: '/react/effects', badge: '难点' },
  { title: '表单与请求', description: '处理受控表单、提交、防重复、接口请求和数据流。', link: '/react/forms', badge: '项目' },
  { title: 'Context 与状态管理', description: '理解 Context 适用边界、局部状态和全局状态的取舍。', link: '/react/context-state-management', badge: '状态' },
  { title: '路由与项目结构', description: '使用 React Router 组织页面、布局和受保护路由。', link: '/react/router-structure', badge: '项目' },
  { title: '性能与测试', description: '学习 memo、useMemo、React DevTools、组件测试和 E2E 策略。', link: '/react/performance', badge: '质量' },
  { title: 'React 管理台从零到项目', description: '用用户管理案例串联登录、路由、请求、表格、表单、权限、测试和部署说明。', link: '/react/project-admin', badge: '实战' },
  { title: '常见问题', description: '排查无限循环、状态不更新、key 错乱和重复请求。', link: '/react/troubleshooting', badge: '排错' }
]" />

## React 和 Vue 学习差异

| 主题 | Vue | React |
| --- | --- | --- |
| 模板 | Vue 模板语法 | JSX |
| 状态 | ref/reactive | useState/useReducer |
| 逻辑复用 | composable | custom hook |
| 副作用 | watch/onMounted | useEffect |
| 组件通信 | props/emits | props/callback |

## 下一步

从 [图解 React 核心概念](/react/visual-guide) 开始，再进入 [快速开始](/react/quick-start)。

## 章节地图

| 阶段 | 章节 |
| --- | --- |
| 图解 | 图解 React 核心概念 |
| 入门 | 快速开始、组件与 JSX |
| 核心 | Hooks 与状态、Effect 与副作用 |
| 项目 | 表单处理、请求与数据流、Context 与状态管理、路由与项目结构、React 管理台从零到项目 |
| 质量 | 性能优化、测试策略、最佳实践、常见问题 |
