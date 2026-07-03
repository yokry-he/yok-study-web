# Vue 学习导览

## 适合谁看

适合已经具备 HTML、CSS、JavaScript 基础，准备系统学习 Vue 3 的前端开发者。

## 你会学到什么

- Vue 3 的组件模型。
- 响应式数据与模板渲染。
- 组合式 API 的组织方式。
- 路由、状态管理、请求封装和权限控制的项目实践。

<TechGrid :items="[
  { title: '图解 Vue 核心概念', description: '先用图理解应用结构、响应式更新、组件通信、路由权限、Pinia 和请求封装。', link: '/vue/visual-guide', level: '图解' },
  { title: '快速开始', description: '创建 Vue 应用，理解单文件组件和基础开发流程。', link: '/vue/quick-start', level: '入门' },
  { title: '模板语法', description: '掌握插值、属性绑定、事件、条件、列表和 v-model。', link: '/vue/template-syntax', level: '基础' },
  { title: '响应式基础', description: '理解 ref、reactive、computed、watch 的职责。', link: '/vue/reactivity', level: '核心' },
  { title: '组件设计', description: '学习 props、emits、slot 和组件边界设计。', link: '/vue/component', level: '核心' },
  { title: '组合式 API', description: '用 composable 抽离状态、逻辑和副作用。', link: '/vue/composition-api', level: '进阶' },
  { title: '生命周期', description: '理解请求、DOM、定时器、KeepAlive 的执行时机。', link: '/vue/lifecycle', level: '核心' },
  { title: '路由与页面', description: '组织真实页面、嵌套路由、路由守卫和权限入口。', link: '/vue/router', level: '项目' },
  { title: 'Pinia 状态管理', description: '管理登录态、用户信息、菜单和跨页面状态。', link: '/vue/pinia', level: '项目' },
  { title: '表单处理', description: '覆盖默认值、编辑回显、校验、转换和防重复提交。', link: '/vue/forms', level: '项目' },
  { title: 'Vue Admin 学习地图', description: '按阶段串起从零到项目、用户模块、权限路由、问题库和交付验收。', link: '/roadmap/vue-admin-learning-map', level: '路线' },
  { title: '图解 Vue Admin 架构', description: '用图理解后台项目分层、目录、数据流、权限流、请求流、表单流和交付链路。', link: '/vue/admin-architecture-visual-guide', level: '图解' },
  { title: '从零到项目落地', description: '用后台管理系统串联目录、路由、请求、列表、表单、权限和验收清单。', link: '/vue/project-from-zero', level: '实战' },
  { title: 'Vue Admin Mock 到真实接口', description: '把本地 mock、环境变量、Vite proxy、DTO 转换、分页参数、401/403 和 traceId 联调串成闭环。', link: '/vue/admin-mock-to-api', level: '实战' },
  { title: 'Vue Admin 列表搜索表格', description: '把搜索条件、QueryState、分页、表格列、批量选择、行操作、导出和空态错误态串成列表页闭环。', link: '/vue/admin-list-search-table', level: '实战' },
  { title: 'Vue Admin 表单新增编辑', description: '把新增、编辑、复制、FormState、Payload、422 回填、关闭确认和防重复提交串成表单闭环。', link: '/vue/admin-form-modal-crud', level: '实战' },
  { title: 'Vue Admin 详情状态记录', description: '把详情页、状态流转、操作按钮、时间线、审计日志、隐藏路由和列表同步串成闭环。', link: '/vue/admin-detail-status-audit', level: '实战' },
  { title: 'Vue Admin 文件导入导出', description: '把文件上传、下载、模板导入、异步导出任务、进度轮询、权限审计和错误处理串成闭环。', link: '/vue/admin-file-import-export', level: '实战' },
  { title: 'Vue Admin 工作台看板', description: '把指标口径、统计卡片、趋势图、排行榜、待办、自动刷新、权限范围和图表性能串成闭环。', link: '/vue/admin-dashboard-analytics', level: '实战' },
  { title: 'Vue Admin 审批流闭环', description: '把流程模板、实例、任务、待办、同意驳回、转办撤回、状态机、权限和审计时间线串成闭环。', link: '/vue/admin-approval-workflow', level: '实战' },
  { title: 'Vue Admin 消息通知闭环', description: '把业务事件、站内信、未读数量、实时提醒、已读未读、通知偏好、重连补偿和权限范围串成闭环。', link: '/vue/admin-notification-center', level: '实战' },
  { title: 'Vue Admin 权限路由闭环', description: '把登录态、用户信息、菜单、动态路由、按钮权限、接口权限和刷新恢复串成完整链路。', link: '/vue/admin-permission-route-flow', level: '实战' },
  { title: 'Vue Admin 用户模块', description: '用用户管理模块学习列表、搜索、弹窗、DTO、权限按钮和刷新恢复。', link: '/vue/admin-user-module', level: '实战' },
  { title: 'Vue Admin 角色权限', description: '把角色、权限树、按钮权限、API 权限和数据范围串成一条链路。', link: '/vue/admin-permission-module', level: '实战' },
  { title: 'Vue Admin 菜单动态路由', description: '把后端菜单变成侧边栏、动态路由、面包屑、标签页和刷新恢复流程。', link: '/vue/admin-menu-route-module', level: '实战' },
  { title: 'Vue Admin 组织数据权限', description: '把部门树、员工归属、角色数据范围和业务列表查询连成闭环。', link: '/vue/admin-organization-data-permission', level: '实战' },
  { title: 'Vue Admin 请求错误处理', description: '统一处理 401、403、并发请求、重复提交、导出任务和页面错误状态。', link: '/vue/admin-request-error-handling', level: '实战' },
  { title: 'Vue 真实项目问题库', description: '排查动态路由、Pinia、表单污染、重复请求、KeepAlive、权限按钮和组件库样式问题。', link: '/projects/issues-vue', level: '排错' },
  { title: '性能与测试', description: '建立性能排查、测试分层和项目质量检查习惯。', link: '/vue/performance', level: '进阶' }
]" />

## 推荐学习顺序

先学模板、响应式和组件，再学路由、状态、表单、请求和权限。不要一开始就把权限、动态路由、请求拦截器全部塞进项目，应该逐步叠加。

## 文档结构

每篇 Vue 文档都按“核心概念、基础示例、实战场景、常见问题、最佳实践”的结构组织，方便后续持续扩写。

## 完整章节

| 阶段 | 章节 |
| --- | --- |
| 图解 | 图解 Vue 核心概念 |
| 入门 | 快速开始、模板语法、响应式基础 |
| 组件 | 组件设计、组合式 API、生命周期 |
| 项目 | 路由与页面、Pinia、表单处理、请求封装、权限与菜单、从零到项目落地 |
| Vue Admin 实战 | 图解项目架构、Mock 到真实接口联调、列表搜索表格闭环、表单新增编辑闭环、详情状态记录闭环、文件导入导出闭环、工作台数据看板闭环、审批流状态机闭环、消息通知闭环、权限路由闭环、用户模块、角色权限模块、菜单与动态路由模块、组织架构与数据权限模块、请求与错误处理模块 |
| 进阶 | 内置组件、性能优化、测试策略、最佳实践、常见问题、Vue 真实项目问题库 |
