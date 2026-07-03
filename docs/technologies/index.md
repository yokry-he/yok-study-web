# 技术库总览

## 这个页面解决什么

技术库是整个文档站的总入口。它负责回答三个问题：

- 现在已经有哪些技术模块。
- 每个模块适合谁学。
- 后续应该按什么顺序继续扩展。

第一阶段不追求“所有技术都铺开”，而是先把前端 Vue 方向做深，再用同一套结构扩展到 React、Node、数据库、DevOps、AI 工程等方向。

## 当前已完成模块

<TechGrid :items="[
  { title: '前端基础', description: 'HTML/CSS、TypeScript、浏览器与网络，帮助理解页面、类型和运行环境。', link: '/frontend/html-css', level: '基础' },
  { title: 'CSS', description: '独立样式模块，覆盖图解入门、盒模型、Flex、Grid、响应式、动画、可访问性、设计 token、样式架构和排错。', link: '/css/introduction', level: '样式' },
  { title: '浏览器与网络', description: '覆盖图解入门、HTTP、跨域、缓存、存储、安全、PWA、Web API、实时通信、原生组件、Wasm、WebGPU、自动化调试和性能排错。', link: '/browser/introduction', level: '基础' },
  { title: 'JavaScript', description: '独立语言模块，覆盖图解入门、数据类型、函数、原型链、数组对象、DOM 事件、正则、异步、事件循环、错误处理、内存管理、模块化和项目落地。', link: '/javascript/introduction', level: '语言' },
  { title: 'TypeScript', description: '独立类型系统模块，覆盖图解入门、基础类型、接口、泛型、类型收窄、工具类型、tsconfig 工程配置、Vue 集成和项目落地。', link: '/typescript/introduction', level: '类型' },
  { title: 'Vue', description: '完整 Vue 3 模块，覆盖图解入门、模板、响应式、组件、路由、状态、表单、项目落地、性能、测试和排错。', link: '/vue/introduction', level: '框架' },
  { title: 'React', description: 'React 稳定模块，覆盖图解入门、JSX、组件、Hooks、Effect、表单、请求、状态、路由、性能、测试、管理台项目和常见问题。', link: '/react/introduction', level: '框架' },
  { title: 'Nuxt / Next 元框架', description: '稳定框架生态模块，覆盖 SSR、文件路由、数据获取、部署缓存、服务端鉴权、SEO、国际化、内容站案例和问题排查。', link: '/meta-frameworks/introduction', level: '框架生态' },
  { title: 'Node.js', description: '后端 JavaScript 稳定模块，覆盖图解入门、运行时、包管理、HTTP API、鉴权会话、数据库集成、错误日志、测试、安全、部署和权限 API 项目。', link: '/node/introduction', level: '后端' },
  { title: 'Java', description: '企业后端语言模块，覆盖图解入门、JDK、面向对象、集合泛型、异常日志、并发虚拟线程、JVM、Spring Boot、事务、测试和排错。', link: '/java/introduction', level: '后端' },
  { title: 'Go', description: '云原生后端语言模块，覆盖图解入门、Go Modules、类型函数、接口组合、错误处理、并发、Context、HTTP、数据库、测试、部署和性能诊断。', link: '/go/introduction', level: '后端' },
  { title: '数据库', description: '稳定数据模块，覆盖图解入门、MySQL、PostgreSQL、Redis、权限系统数据层落地、建模、索引、事务、迁移、ORM、备份恢复、安全审计、脱敏和排错。', link: '/database/introduction', level: '数据' },
  { title: 'AI 工程', description: '稳定 AI 工程模块，覆盖图解入门、LLM API、提示词、结构化输出、函数调用、多模态、RAG、MCP、Agent、文档问答项目、评测、上线治理和问题排查。', link: '/ai-engineering/introduction', level: 'AI' },
  { title: '前端工程化', description: '稳定工程模块，覆盖图解入门、Vite、规范、环境、依赖、测试、Monorepo、组件库工程、构建部署、包体积、模块联邦、工程性能和排错。', link: '/engineering/introduction', level: '工程' },
  { title: 'DevOps 与部署', description: '稳定运维模块，覆盖图解入门、Linux、Nginx、Docker、CI/CD、项目上线全流程、发布回滚、可观测性、Kubernetes、云服务、对象存储和线上排查。', link: '/devops/introduction', level: '运维' },
  { title: '实战项目', description: '稳定实战模块，覆盖 Vue Admin、组件库、项目阶段任务、权限、权限运营、组织架构、审批流、文件中心、数据看板、多租户、消息通知、导入导出、支付订单、会员订阅、搜索中心、任务调度、消息队列、开放平台、工作流、低代码流程、审计中心、运营活动、财务对账、渠道结算、渠道费用稽核、渠道费用 ROI 复盘、渠道费用预算优化、渠道费用异常预警、渠道费用策略灰度、渠道策略效果复盘、渠道策略对照实验、渠道策略版本治理、渠道策略审批矩阵、渠道策略发布审计、渠道策略回滚治理、渠道策略异常仲裁、渠道策略仲裁复盘、渠道策略裁决标准库、渠道策略标准效果监控、渠道策略标准灰度发布、渠道策略标准版本回滚、渠道策略标准回滚演练、渠道价格稽核、渠道窜货监控、渠道信用评级、渠道返利风控、渠道政策模拟、渠道利润模拟、渠道价格弹性、主数据、客户主数据、低代码表单、报表、智能 BI、客服工单、客服质检、客服知识运营、系统集成、国际化、数据治理、数据质量、数据资产运营、数据安全运营、规则引擎、灰度发布、灾备、风控、合同、合同付款、合同变更、合同续签、客户合同风险预警、客户合同收入预测、知识库、配置中心、合规审计、客户成功、客户生命周期价值、客户流失预警、客户续费挽回、客户续约定价、客户分群运营、客户触达自动化、客户权益运营、客户投诉闭环、工单自动化、计费、数据交换、企业门户、资产、预算、资金计划、费用报销、员工借款、税务、发票协同、采购、采购寻源、供应商准入、供应商合同协同、供应商协同门户、供应商门户权限审计、供应商索赔、供应商绩效、供应链计划、项目管理、研发需求池、报价、价格审批、上线事故、运维值班、库存、渠道库存、备件库存、备件补货、售后备件周转、备件旧件返修、仓储、售后、售后远程诊断、售后专家协同、售后知识自动推荐、售后知识质量治理、售后知识智能检索优化、售后知识问答助手、售后知识自动质检、售后知识专家审核、售后知识发布灰度、售后知识回滚治理、售后知识影响追踪、售后知识客户通知治理、售后知识外部服务商通知协同、售后知识服务商培训闭环、售后知识培训效果复盘、售后知识培训认证治理、售后知识认证派单联动、售后知识认证质量稽核、售后知识认证风险画像、退换货质检、退款风控、售后结算、现场服务收费、售后备件成本核算、售后成本毛利、售后服务成本优化、售后 SLA 赔付、售后服务商评级、售后维修质量复盘、售后投诉根因分析、报修派单、服务网点、数据权限审计、门店零售、CRM、客户账期、客户授信风控、客户回款风险预测、客户坏账处置策略、客户应收催收自动化、销售回款预测调度、销售现金流预警、销售回款策略模拟、销售风险动作编排、销售风险处置复盘、销售风险预案演练、销售风险指标治理、销售风险指标血缘审计、销售风险指标异常根因、销售风险指标自动修复、销售风险指标治理成熟度、销售风险指标治理运营看板、销售风险指标治理成本收益评估、销售回款计划、销售预测复盘、销售目标拆解、销售佣金、销售返利、会员营销、生产制造、生产排程、产能负荷预测、生产计划达成、生产停线损失复盘、质量异常、生产异常 CAPA、生产过程审核、生产巡检移动端、生产现场安全隐患、生产安全培训闭环、生产安全考试认证、生产安全风险画像、生产安全应急演练、生产安全事故复盘、生产安全风险整改复查、生产安全整改看板、生产安全整改 SLA、生产安全整改成本复盘、生产安全整改预算预测、生产安全整改资源排期、生产安全整改产线影响评估、生产安全整改多方案决策、生产安全整改决策复盘、生产安全整改决策知识库、生产安全整改决策智能推荐、生产良率、生产瓶颈、生产换型损失、生产设备异常、生产能耗分析、生产成本、制造成本差异、IoT 和教育培训。', link: '/projects/real-world-issues', level: '实战' },
  { title: '速查手册', description: '稳定速查模块，覆盖 Vue、JavaScript、TypeScript、CSS、正则、Node、Java、Go、HTTP、Git、Linux、Docker、Nginx、SQL、Redis、常用命令和调试工具。', link: '/cheatsheets/', level: '速查' }
]" />

## 推荐学习路径

如果你还不确定从哪里开始，先看 [阅读顺序与使用方法](/roadmap/reading-guide) 和 [学习路线总览](/roadmap/introduction)。它们会说明如何按路线、技术模块、练习包、问题库和速查手册组合使用。

如果你是前端方向，建议按这个顺序：

```text
HTML/CSS
↓
JavaScript
↓
TypeScript
↓
Vue 3
↓
Vue Router + Pinia
↓
请求、权限、表单
↓
工程化与部署
↓
Vue Admin 实战
↓
真实项目问题库
```

如果你已经有基础，可以直接进入对应模块：

| 目标 | 入口 |
| --- | --- |
| 学会使用本站 | [阅读顺序与使用方法](/roadmap/reading-guide) |
| 选择学习路线 | [学习路线总览](/roadmap/introduction) |
| 先看懂技术之间如何连接 | [图解学习地图](/roadmap/visual-learning-map) |
| 做阶段练习 | [学习路径练习包](/roadmap/practice-labs) |
| 做 Vue Admin 练习 | [Vue Admin 专项练习](/roadmap/vue-admin-practice) |
| 系统完成 Vue Admin 项目 | [Vue Admin 学习地图与交付清单](/roadmap/vue-admin-learning-map) |
| 细化列表搜索分页表格 | [Vue Admin 列表、搜索、分页与表格闭环实战](/vue/admin-list-search-table) |
| 细化表单新增编辑校验 | [Vue Admin 表单弹窗、新增编辑与校验闭环实战](/vue/admin-form-modal-crud) |
| 细化详情状态操作记录 | [Vue Admin 详情页、状态流转与操作记录闭环实战](/vue/admin-detail-status-audit) |
| 细化文件上传导入导出 | [Vue Admin 文件上传、下载、导入导出与异步任务闭环实战](/vue/admin-file-import-export) |
| 细化工作台数据看板 | [Vue Admin 工作台、统计卡片、图表看板与数据刷新闭环实战](/vue/admin-dashboard-analytics) |
| 细化审批流状态机 | [Vue Admin 审批流、状态机、待办与审计闭环实战](/vue/admin-approval-workflow) |
| 细化用户管理模块 | [Vue Admin 用户模块实现手册](/vue/admin-user-module) |
| 细化角色权限模块 | [Vue Admin 角色权限模块实现手册](/vue/admin-permission-module) |
| 细化菜单动态路由模块 | [Vue Admin 菜单与动态路由实现手册](/vue/admin-menu-route-module) |
| 细化组织数据权限模块 | [Vue Admin 组织架构与数据权限实现手册](/vue/admin-organization-data-permission) |
| 细化请求错误处理模块 | [Vue Admin 请求封装与错误处理闭环手册](/vue/admin-request-error-handling) |
| 不知道项目问题从哪查 | [项目排障方法论](/projects/debugging-playbook) |
| 补 JavaScript | [JavaScript 基础](/javascript/fundamentals) |
| 做 JS 任务看板 | [JavaScript 任务看板从零到项目](/javascript/task-board-project) |
| 补 CSS | [CSS 学习导览](/css/introduction) |
| 补浏览器与网络 | [浏览器学习导览](/browser/introduction) |
| 补 TypeScript | [TypeScript 学习导览](/typescript/introduction) |
| 做 TS 类型边界项目 | [TypeScript 类型边界从零到项目](/typescript/type-boundary-project) |
| 查 TS 类型问题 | [TypeScript 类型边界问题](/projects/issues-typescript) |
| 系统学 Vue | [Vue 学习导览](/vue/introduction) |
| 看懂 Vue Admin 架构 | [图解 Vue Admin 项目架构](/vue/admin-architecture-visual-guide) |
| 从 mock 切真实接口 | [Vue Admin Mock 到真实接口联调实战](/vue/admin-mock-to-api) |
| 做好列表搜索分页表格 | [Vue Admin 列表、搜索、分页与表格闭环实战](/vue/admin-list-search-table) |
| 做好表单新增编辑校验 | [Vue Admin 表单弹窗、新增编辑与校验闭环实战](/vue/admin-form-modal-crud) |
| 做好详情状态操作记录 | [Vue Admin 详情页、状态流转与操作记录闭环实战](/vue/admin-detail-status-audit) |
| 做好文件上传导入导出 | [Vue Admin 文件上传、下载、导入导出与异步任务闭环实战](/vue/admin-file-import-export) |
| 做好工作台统计图表 | [Vue Admin 工作台、统计卡片、图表看板与数据刷新闭环实战](/vue/admin-dashboard-analytics) |
| 做好审批流待办状态机 | [Vue Admin 审批流、状态机、待办与审计闭环实战](/vue/admin-approval-workflow) |
| 做 Vue Admin 权限路由闭环 | [Vue Admin 权限路由闭环实战](/vue/admin-permission-route-flow) |
| 查 Vue 项目问题 | [Vue 真实项目问题](/projects/issues-vue) |
| 查 Vue Admin 请求权限问题 | [Vue Admin 请求、权限与数据问题排查专题](/projects/issues-vue-admin-request) |
| 入门 React | [React 学习导览](/react/introduction) |
| 做 React 管理台 | [React 管理台从零到项目](/react/project-admin) |
| 学 Nuxt / Next | [Nuxt / Next 元框架学习导览](/meta-frameworks/introduction) |
| 入门 Node.js | [Node.js 学习导览](/node/introduction) |
| 做权限 API | [Node 权限 API 从零到项目](/node/permission-api-project) |
| 学 Java 后端 | [Java 学习导览](/java/introduction) |
| 学 Go 后端 | [Go 学习导览](/go/introduction) |
| 学数据库 | [数据库学习导览](/database/introduction) |
| 入门 AI 工程 | [AI 工程学习导览](/ai-engineering/introduction) |
| 做 AI 文档问答 | [AI 文档问答从零到项目](/ai-engineering/doc-qa-project) |
| 学部署上线 | [DevOps 学习导览](/devops/introduction) |
| 做后台项目 | [Vue Admin 实战](/projects/vue-admin) |
| 做权限运营 | [权限运营项目案例](/projects/permission-operation-case) |
| 做组织模块 | [组织架构项目案例](/projects/organization-case) |
| 做审批模块 | [Vue Admin 审批流、状态机、待办与审计闭环实战](/vue/admin-approval-workflow)、[审批流项目案例](/projects/approval-workflow-case) |
| 做文件中心 | [文件中心项目案例](/projects/file-center-case) |
| 做数据看板 | [数据看板项目案例](/projects/analytics-dashboard-case) |
| 做多租户权限 | [多租户权限项目案例](/projects/multi-tenant-permission-case) |
| 做消息通知 | [消息通知项目案例](/projects/notification-center-case) |
| 做导入导出 | [数据导入导出项目案例](/projects/import-export-case) |
| 做支付订单 | [支付订单项目案例](/projects/payment-order-case) |
| 做会员订阅 | [会员订阅项目案例](/projects/subscription-billing-case) |
| 做搜索中心 | [搜索中心项目案例](/projects/search-center-case) |
| 做任务调度 | [任务调度项目案例](/projects/task-scheduler-case) |
| 做消息队列 | [消息队列项目案例](/projects/message-queue-case) |
| 做开放平台 | [第三方开放平台项目案例](/projects/open-platform-case) |
| 做工作流配置器 | [工作流配置器项目案例](/projects/workflow-builder-case) |
| 做低代码流程平台 | [低代码流程平台项目案例](/projects/low-code-workflow-case) |
| 做审计中心 | [审计中心项目案例](/projects/audit-center-case) |
| 做运营活动 | [运营活动项目案例](/projects/marketing-campaign-case) |
| 做财务对账 | [复杂财务对账项目案例](/projects/finance-reconciliation-case) |
| 做渠道结算 | [渠道结算项目案例](/projects/channel-settlement-case) |
| 做渠道费用稽核 | [渠道费用稽核项目案例](/projects/channel-expense-audit-case) |
| 做渠道费用 ROI 复盘 | [渠道费用 ROI 复盘项目案例](/projects/channel-expense-roi-review-case) |
| 做渠道费用预算优化 | [渠道费用预算优化项目案例](/projects/channel-expense-budget-optimization-case) |
| 做渠道费用异常预警 | [渠道费用异常预警项目案例](/projects/channel-expense-anomaly-warning-case) |
| 做渠道费用策略灰度 | [渠道费用策略灰度项目案例](/projects/channel-expense-strategy-gray-release-case) |
| 做渠道策略效果复盘 | [渠道策略效果复盘项目案例](/projects/channel-strategy-effect-review-case) |
| 做渠道策略对照实验 | [渠道策略对照实验项目案例](/projects/channel-strategy-ab-experiment-case) |
| 做渠道策略版本治理 | [渠道策略版本治理项目案例](/projects/channel-strategy-version-governance-case) |
| 做渠道策略审批矩阵 | [渠道策略审批矩阵项目案例](/projects/channel-strategy-approval-matrix-case) |
| 做渠道策略发布审计 | [渠道策略发布审计项目案例](/projects/channel-strategy-release-audit-case) |
| 做渠道策略回滚治理 | [渠道策略回滚治理项目案例](/projects/channel-strategy-rollback-governance-case) |
| 做渠道策略异常仲裁 | [渠道策略异常仲裁项目案例](/projects/channel-strategy-exception-arbitration-case) |
| 做渠道策略仲裁复盘 | [渠道策略仲裁复盘项目案例](/projects/channel-strategy-arbitration-review-case) |
| 做渠道策略裁决标准库 | [渠道策略裁决标准库项目案例](/projects/channel-strategy-decision-standard-library-case) |
| 做渠道策略标准效果监控 | [渠道策略标准效果监控项目案例](/projects/channel-strategy-standard-effect-monitoring-case) |
| 做渠道策略标准灰度发布 | [渠道策略标准灰度发布项目案例](/projects/channel-strategy-standard-gray-release-case) |
| 做渠道策略标准版本回滚 | [渠道策略标准版本回滚项目案例](/projects/channel-strategy-standard-version-rollback-case) |
| 做渠道策略标准回滚演练 | [渠道策略标准回滚演练项目案例](/projects/channel-strategy-standard-rollback-drill-case) |
| 做渠道策略标准灾备切换 | [渠道策略标准灾备切换项目案例](/projects/channel-strategy-standard-disaster-recovery-switch-case) |
| 做渠道价格稽核 | [渠道价格稽核项目案例](/projects/channel-price-audit-case) |
| 做渠道窜货监控 | [渠道窜货监控项目案例](/projects/channel-diversion-monitor-case) |
| 做渠道信用评级 | [渠道信用评级项目案例](/projects/channel-credit-rating-case) |
| 做渠道返利风控 | [渠道返利风控项目案例](/projects/channel-rebate-risk-control-case) |
| 做渠道政策模拟 | [渠道政策模拟项目案例](/projects/channel-policy-simulation-case) |
| 做渠道利润模拟 | [渠道利润模拟项目案例](/projects/channel-profit-simulation-case) |
| 做渠道价格弹性分析 | [渠道价格弹性分析项目案例](/projects/channel-price-elasticity-analysis-case) |
| 做主数据管理 | [主数据管理项目案例](/projects/master-data-case) |
| 做客户主数据 | [客户主数据项目案例](/projects/customer-master-data-case) |
| 做低代码表单 | [低代码表单项目案例](/projects/low-code-form-case) |
| 做报表配置器 | [报表配置器项目案例](/projects/report-builder-case) |
| 做智能报表与 BI | [智能报表与 BI 分析项目案例](/projects/smart-bi-dashboard-case) |
| 做客服工单 | [客服工单项目案例](/projects/support-ticket-case) |
| 做客服质检 | [客服质检项目案例](/projects/customer-service-quality-case) |
| 做系统集成 | [集团级系统集成项目案例](/projects/enterprise-integration-case) |
| 做国际化后台 | [国际化后台项目案例](/projects/i18n-admin-case) |
| 做数据治理平台 | [数据治理平台项目案例](/projects/data-governance-case) |
| 做数据质量专项 | [数据质量专项项目案例](/projects/data-quality-special-case) |
| 做数据资产运营 | [数据资产运营项目案例](/projects/data-asset-operation-case) |
| 做数据安全运营 | [数据安全运营项目案例](/projects/data-security-operation-case) |
| 做规则引擎 | [规则引擎项目案例](/projects/rule-engine-case) |
| 做灰度发布 | [灰度发布后台项目案例](/projects/gray-release-admin-case) |
| 做跨区域灾备 | [跨区域灾备管理项目案例](/projects/disaster-recovery-case) |
| 做风控中心 | [风控中心项目案例](/projects/risk-control-center-case) |
| 做合同管理 | [合同管理项目案例](/projects/contract-management-case) |
| 做合同履约 | [合同履约项目案例](/projects/contract-fulfillment-case) |
| 做合同付款 | [合同付款项目案例](/projects/contract-payment-case) |
| 做合同变更 | [合同变更项目案例](/projects/contract-change-case) |
| 做合同续签 | [合同续签项目案例](/projects/contract-renewal-case) |
| 做客户合同风险预警 | [客户合同风险预警项目案例](/projects/customer-contract-risk-warning-case) |
| 做客户合同收入预测 | [客户合同收入预测项目案例](/projects/customer-contract-revenue-forecast-case) |
| 做知识库平台 | [知识库平台项目案例](/projects/knowledge-base-case) |
| 做客服知识运营 | [客服知识运营项目案例](/projects/customer-knowledge-operation-case) |
| 做统一配置中心 | [统一配置中心项目案例](/projects/config-center-case) |
| 做合规审计 | [行业合规审计项目案例](/projects/compliance-audit-case) |
| 做客户成功平台 | [客户成功平台项目案例](/projects/customer-success-case) |
| 做客户生命周期价值分析 | [客户生命周期价值分析项目案例](/projects/customer-lifetime-value-analysis-case) |
| 做客户流失预警 | [客户流失预警项目案例](/projects/customer-churn-warning-case) |
| 做客户续费挽回 | [客户续费挽回项目案例](/projects/customer-renewal-recovery-case) |
| 做客户续约定价策略 | [客户续约定价策略项目案例](/projects/customer-renewal-pricing-strategy-case) |
| 做客户分群运营 | [客户分群运营项目案例](/projects/customer-segmentation-operation-case) |
| 做客户触达自动化 | [客户触达自动化项目案例](/projects/customer-touch-automation-case) |
| 做客户权益运营 | [客户权益运营项目案例](/projects/customer-benefit-operation-case) |
| 做客户投诉闭环 | [客户投诉闭环项目案例](/projects/customer-complaint-closed-loop-case) |
| 做工单自动化 | [工单自动化项目案例](/projects/ticket-automation-case) |
| 做计费中台 | [计费中台项目案例](/projects/billing-platform-case) |
| 做数据交换平台 | [数据交换平台项目案例](/projects/data-exchange-platform-case) |
| 做企业门户 | [企业门户项目案例](/projects/enterprise-portal-case) |
| 做资产管理 | [资产管理项目案例](/projects/asset-management-case) |
| 做预算管理 | [预算管理项目案例](/projects/budget-management-case) |
| 做资金计划 | [资金计划项目案例](/projects/cash-flow-planning-case) |
| 做费用报销 | [费用报销项目案例](/projects/expense-reimbursement-case) |
| 做员工借款 | [员工借款项目案例](/projects/employee-loan-case) |
| 做税务管理 | [税务管理项目案例](/projects/tax-management-case) |
| 做发票协同 | [发票协同项目案例](/projects/invoice-collaboration-case) |
| 做采购管理 | [采购管理项目案例](/projects/procurement-management-case) |
| 做采购寻源 | [采购寻源项目案例](/projects/procurement-sourcing-case) |
| 做供应商准入 | [供应商准入项目案例](/projects/supplier-onboarding-case) |
| 做供应商合同协同 | [供应商合同协同项目案例](/projects/supplier-contract-collaboration-case) |
| 做供应商协同门户 | [供应商协同门户项目案例](/projects/supplier-collaboration-portal-case) |
| 做供应商门户权限审计 | [供应商门户权限审计项目案例](/projects/supplier-portal-permission-audit-case) |
| 做供应商索赔 | [供应商索赔项目案例](/projects/supplier-claim-case) |
| 做供应商绩效 | [供应商绩效项目案例](/projects/supplier-performance-case) |
| 做供应链计划 | [供应链计划项目案例](/projects/supply-chain-planning-case) |
| 做项目管理 | [项目管理项目案例](/projects/project-management-case) |
| 做研发需求池 | [研发需求池项目案例](/projects/rd-requirement-pool-case) |
| 做报价中心 | [报价中心项目案例](/projects/quotation-center-case) |
| 做价格审批中心 | [价格审批中心项目案例](/projects/price-approval-center-case) |
| 查上线事故 | [上线事故案例库](/projects/production-incident-cases) |
| 做运维值班 | [运维值班项目案例](/projects/operations-oncall-case) |
| 做库存管理 | [库存管理项目案例](/projects/inventory-management-case) |
| 做渠道库存协同 | [渠道库存协同项目案例](/projects/channel-inventory-collaboration-case) |
| 做备件库存 | [备件库存项目案例](/projects/spare-parts-inventory-case) |
| 做备件补货 | [备件补货项目案例](/projects/spare-parts-replenishment-case) |
| 做售后备件周转分析 | [售后备件周转分析项目案例](/projects/after-sales-spare-parts-turnover-case) |
| 做备件旧件返修 | [备件旧件返修项目案例](/projects/spare-parts-return-repair-case) |
| 做仓储物流 | [仓储物流项目案例](/projects/warehouse-logistics-case) |
| 做售后服务 | [售后服务项目案例](/projects/after-sales-service-case) |
| 做售后远程诊断 | [售后远程诊断项目案例](/projects/after-sales-remote-diagnosis-case) |
| 做售后专家协同 | [售后专家协同项目案例](/projects/after-sales-expert-collaboration-case) |
| 做售后知识自动推荐 | [售后知识自动推荐项目案例](/projects/after-sales-knowledge-recommendation-case) |
| 做售后知识质量治理 | [售后知识质量治理项目案例](/projects/after-sales-knowledge-quality-governance-case) |
| 做售后知识智能检索优化 | [售后知识智能检索优化项目案例](/projects/after-sales-knowledge-search-optimization-case) |
| 做售后知识问答助手 | [售后知识问答助手项目案例](/projects/after-sales-knowledge-qa-assistant-case) |
| 做售后知识自动质检 | [售后知识自动质检项目案例](/projects/after-sales-knowledge-auto-quality-inspection-case) |
| 做售后知识专家审核 | [售后知识专家审核项目案例](/projects/after-sales-knowledge-expert-review-case) |
| 做售后知识发布灰度 | [售后知识发布灰度项目案例](/projects/after-sales-knowledge-release-gray-case) |
| 做售后知识回滚治理 | [售后知识回滚治理项目案例](/projects/after-sales-knowledge-rollback-governance-case) |
| 做售后知识影响追踪 | [售后知识影响追踪项目案例](/projects/after-sales-knowledge-impact-trace-case) |
| 做售后知识客户通知治理 | [售后知识客户通知治理项目案例](/projects/after-sales-knowledge-customer-notification-governance-case) |
| 做售后知识外部服务商通知协同 | [售后知识外部服务商通知协同项目案例](/projects/after-sales-knowledge-provider-notification-collaboration-case) |
| 做售后知识服务商培训闭环 | [售后知识服务商培训闭环项目案例](/projects/after-sales-knowledge-provider-training-closed-loop-case) |
| 做售后知识培训效果复盘 | [售后知识培训效果复盘项目案例](/projects/after-sales-knowledge-training-effect-review-case) |
| 做售后知识培训认证治理 | [售后知识培训认证治理项目案例](/projects/after-sales-knowledge-training-certification-governance-case) |
| 做售后知识认证派单联动 | [售后知识认证派单联动项目案例](/projects/after-sales-knowledge-certification-dispatch-linkage-case) |
| 做售后知识认证质量稽核 | [售后知识认证质量稽核项目案例](/projects/after-sales-knowledge-certification-quality-audit-case) |
| 做售后知识认证风险画像 | [售后知识认证风险画像项目案例](/projects/after-sales-knowledge-certification-risk-profile-case) |
| 做售后知识认证服务商整改 | [售后知识认证服务商整改项目案例](/projects/after-sales-knowledge-certification-provider-rectification-case) |
| 做客户退换货质检 | [客户退换货质检项目案例](/projects/customer-return-quality-inspection-case) |
| 做客户退款风控 | [客户退款风控项目案例](/projects/customer-refund-risk-control-case) |
| 做售后结算 | [售后结算项目案例](/projects/after-sales-settlement-case) |
| 做现场服务收费 | [现场服务收费项目案例](/projects/field-service-charging-case) |
| 做售后备件成本核算 | [售后备件成本核算项目案例](/projects/after-sales-spare-part-cost-case) |
| 做售后成本毛利分析 | [售后成本毛利分析项目案例](/projects/after-sales-cost-margin-case) |
| 做售后服务成本优化 | [售后服务成本优化项目案例](/projects/after-sales-service-cost-optimization-case) |
| 做售后 SLA 赔付分析 | [售后 SLA 赔付分析项目案例](/projects/after-sales-sla-compensation-case) |
| 做售后服务商评级 | [售后服务商评级项目案例](/projects/after-sales-provider-rating-case) |
| 做售后维修质量复盘 | [售后维修质量复盘项目案例](/projects/after-sales-repair-quality-review-case) |
| 做售后投诉根因分析 | [售后投诉根因分析项目案例](/projects/after-sales-complaint-root-cause-case) |
| 做报修派单 | [报修派单项目案例](/projects/repair-dispatch-case) |
| 做服务网点 | [服务网点项目案例](/projects/service-outlet-case) |
| 做数据权限审计 | [数据权限审计项目案例](/projects/data-permission-audit-case) |
| 做门店零售 | [门店零售管理项目案例](/projects/retail-store-management-case) |
| 做 CRM 销售管理 | [CRM 销售管理项目案例](/projects/crm-sales-management-case) |
| 做客户账期 | [客户账期项目案例](/projects/customer-credit-term-case) |
| 做客户授信风控 | [客户授信风控项目案例](/projects/customer-credit-risk-control-case) |
| 做客户回款风险预测 | [客户回款风险预测项目案例](/projects/customer-payment-risk-prediction-case) |
| 做客户坏账处置策略 | [客户坏账处置策略项目案例](/projects/customer-bad-debt-disposal-case) |
| 做客户应收催收自动化 | [客户应收催收自动化项目案例](/projects/customer-receivable-collection-automation-case) |
| 做销售回款预测调度 | [销售回款预测调度项目案例](/projects/sales-payment-prediction-scheduling-case) |
| 做销售现金流预警 | [销售现金流预警项目案例](/projects/sales-cash-flow-warning-case) |
| 做销售回款策略模拟 | [销售回款策略模拟项目案例](/projects/sales-collection-strategy-simulation-case) |
| 做销售风险动作编排 | [销售风险动作编排项目案例](/projects/sales-risk-action-orchestration-case) |
| 做销售风险处置复盘 | [销售风险处置复盘项目案例](/projects/sales-risk-disposal-review-case) |
| 做销售风险预案演练 | [销售风险预案演练项目案例](/projects/sales-risk-contingency-drill-case) |
| 做销售风险指标治理 | [销售风险指标治理项目案例](/projects/sales-risk-metric-governance-case) |
| 做销售风险指标血缘审计 | [销售风险指标血缘审计项目案例](/projects/sales-risk-metric-lineage-audit-case) |
| 做销售风险指标异常根因 | [销售风险指标异常根因项目案例](/projects/sales-risk-metric-anomaly-root-cause-case) |
| 做销售风险指标自动修复 | [销售风险指标自动修复项目案例](/projects/sales-risk-metric-auto-repair-case) |
| 做销售风险指标治理成熟度 | [销售风险指标治理成熟度项目案例](/projects/sales-risk-metric-governance-maturity-case) |
| 做销售风险指标治理运营看板 | [销售风险指标治理运营看板项目案例](/projects/sales-risk-metric-governance-operations-dashboard-case) |
| 做销售风险指标治理成本收益评估 | [销售风险指标治理成本收益评估项目案例](/projects/sales-risk-metric-governance-cost-benefit-evaluation-case) |
| 做销售风险指标治理预算审批 | [销售风险指标治理预算审批项目案例](/projects/sales-risk-metric-governance-budget-approval-case) |
| 做销售回款计划 | [销售回款计划项目案例](/projects/sales-collection-plan-case) |
| 做销售预测复盘 | [销售预测复盘项目案例](/projects/sales-forecast-review-case) |
| 做销售目标拆解 | [销售目标拆解项目案例](/projects/sales-target-breakdown-case) |
| 做销售佣金核算 | [销售佣金核算项目案例](/projects/sales-commission-settlement-case) |
| 做销售返利政策 | [销售返利政策项目案例](/projects/sales-rebate-policy-case) |
| 做会员营销 | [会员营销项目案例](/projects/member-marketing-case) |
| 做生产制造 | [生产制造项目案例](/projects/manufacturing-execution-case) |
| 做生产排程 | [生产排程项目案例](/projects/production-scheduling-case) |
| 做产能负荷预测 | [产能负荷预测项目案例](/projects/capacity-load-forecast-case) |
| 做生产计划达成分析 | [生产计划达成分析项目案例](/projects/production-plan-attainment-case) |
| 做生产停线损失复盘 | [生产停线损失复盘项目案例](/projects/production-line-stop-loss-review-case) |
| 做质量追溯 | [质量追溯项目案例](/projects/quality-traceability-case) |
| 做生产质量异常 | [生产质量异常项目案例](/projects/production-quality-exception-case) |
| 做生产异常 CAPA | [生产异常 CAPA 项目案例](/projects/production-exception-capa-case) |
| 做生产过程审核 | [生产过程审核项目案例](/projects/production-process-audit-case) |
| 做生产巡检移动端 | [生产巡检移动端项目案例](/projects/production-mobile-inspection-case) |
| 做生产现场安全隐患 | [生产现场安全隐患项目案例](/projects/production-safety-hazard-case) |
| 做生产安全培训闭环 | [生产安全培训闭环项目案例](/projects/production-safety-training-closed-loop-case) |
| 做生产安全考试认证 | [生产安全考试认证项目案例](/projects/production-safety-exam-certification-case) |
| 做生产安全风险画像 | [生产安全风险画像项目案例](/projects/production-safety-risk-profile-case) |
| 做生产安全应急演练 | [生产安全应急演练项目案例](/projects/production-safety-emergency-drill-case) |
| 做生产安全事故复盘 | [生产安全事故复盘项目案例](/projects/production-safety-incident-review-case) |
| 做生产安全风险整改复查 | [生产安全风险整改复查项目案例](/projects/production-safety-risk-rectification-review-case) |
| 做生产安全整改看板 | [生产安全整改看板项目案例](/projects/production-safety-rectification-dashboard-case) |
| 做生产安全整改 SLA | [生产安全整改 SLA 项目案例](/projects/production-safety-rectification-sla-case) |
| 做生产安全整改成本复盘 | [生产安全整改成本复盘项目案例](/projects/production-safety-rectification-cost-review-case) |
| 做生产安全整改预算预测 | [生产安全整改预算预测项目案例](/projects/production-safety-rectification-budget-forecast-case) |
| 做生产安全整改资源排期 | [生产安全整改资源排期项目案例](/projects/production-safety-rectification-resource-scheduling-case) |
| 做生产安全整改产线影响评估 | [生产安全整改产线影响评估项目案例](/projects/production-safety-rectification-line-impact-assessment-case) |
| 做生产安全整改多方案决策 | [生产安全整改多方案决策项目案例](/projects/production-safety-rectification-multi-scenario-decision-case) |
| 做生产安全整改决策复盘 | [生产安全整改决策复盘项目案例](/projects/production-safety-rectification-decision-review-case) |
| 做生产安全整改决策知识库 | [生产安全整改决策知识库项目案例](/projects/production-safety-rectification-decision-knowledge-base-case) |
| 做生产安全整改决策智能推荐 | [生产安全整改决策智能推荐项目案例](/projects/production-safety-rectification-decision-intelligent-recommendation-case) |
| 做生产安全整改决策推荐评测 | [生产安全整改决策推荐评测项目案例](/projects/production-safety-rectification-recommendation-evaluation-case) |
| 做生产良率分析 | [生产良率分析项目案例](/projects/production-yield-analysis-case) |
| 做生产瓶颈分析 | [生产瓶颈分析项目案例](/projects/production-bottleneck-analysis-case) |
| 做生产换型损失分析 | [生产换型损失分析项目案例](/projects/production-changeover-loss-analysis-case) |
| 做生产设备异常 | [生产设备异常项目案例](/projects/production-equipment-exception-case) |
| 做生产能耗分析 | [生产能耗分析项目案例](/projects/production-energy-analysis-case) |
| 做生产成本核算 | [生产成本核算项目案例](/projects/production-cost-accounting-case) |
| 做制造成本差异分析 | [制造成本差异分析项目案例](/projects/manufacturing-cost-variance-case) |
| 做 IoT 设备管理 | [IoT 设备管理项目案例](/projects/iot-device-management-case) |
| 做设备维保 | [设备维保项目案例](/projects/equipment-maintenance-case) |
| 做教育培训平台 | [教育培训平台项目案例](/projects/education-training-platform-case) |
| 查项目问题 | [真实项目问题库](/projects/real-world-issues) |
| 先定位问题层级 | [项目排障方法论](/projects/debugging-playbook) |
| 做交付检查 | [项目交付检查清单](/projects/delivery-checklist) |
| 快速查写法 | [速查手册总览](/cheatsheets/) |
| 查前端状态问题 | [前端页面与状态问题](/projects/issues-frontend) |
| 查后端接口问题 | [后端接口与服务问题](/projects/issues-backend) |
| 查数据库缓存问题 | [数据库与缓存问题](/projects/issues-database) |
| 查部署上线问题 | [部署、缓存与 DevOps 问题](/projects/issues-deployment) |
| 查 AI 工程问题 | [AI 工程问题](/projects/issues-ai) |
| 学前端工程化 | [前端工程化学习导览](/engineering/introduction) |
| 做组件库工程 | [组件库工程从零到项目](/engineering/component-library-project) |
| 准备上线 | [构建与部署](/engineering/build-deploy) |

## 模块标准结构

后续每个技术模块都应该尽量使用统一结构：

```text
overview          技术概览
quick-start       快速开始
core-concepts     核心概念
practice          实战场景
best-practices    最佳实践
troubleshooting   常见问题
cheatsheet        速查手册
resources         延伸资料
```

这样做的好处是：

- 用户知道每个模块从哪里开始。
- 后续扩展 React、Node、数据库时不需要重新设计结构。
- 贡献者可以按模板补内容。
- 搜索结果更容易理解。

## 后续扩展顺序

优先级建议：

1. `TypeScript` 独立模块：类型系统、泛型、工程实践、Vue 集成。
2. `CSS` 独立模块：布局、响应式、设计系统、组件库样式边界。
3. `React` 模块：组件、Hooks、状态、路由、Next.js。
4. `Node.js` 模块：运行时、包管理、API、鉴权、数据库、测试、安全和部署。
5. `Java` 模块：JDK、JVM、Spring Boot、并发、事务、测试和排错。
6. `Go` 模块：Go Modules、并发、Context、HTTP、数据库、测试、部署和性能诊断。
7. `数据库` 模块：MySQL、PostgreSQL、Redis、建模、性能、ORM、备份和安全。
8. `DevOps` 模块：Linux、Docker、Nginx、CI/CD、观测、Kubernetes 和云部署。
9. `AI 工程` 模块：LLM API、结构化输出、函数调用、多模态、RAG、MCP、Agent、评测和上线。

## 内容质量要求

每篇文档都应该尽量包含：

- 适合谁看。
- 你会学到什么。
- 核心概念。
- 基础示例。
- 实际项目问题。
- 解决方案。
- 最佳实践。
- 下一步学习。

不要只写 API 清单。API 清单适合速查手册，学习文档必须解释场景、原因和取舍。

## 当前阶段结论

现在站点已经具备长期扩展基础：

- VitePress 架构稳定。
- 首页和导航可扩展。
- Vue 模块已经较完整。
- 学习路线已经补齐阅读顺序、阶段任务、学习路径练习包、项目里程碑、能力自测和维护规则。
- JavaScript 已成为独立模块，覆盖语言基础、运行机制、错误处理和工程实践。
- 前端工程化已经扩展为稳定学习路径。
- Nuxt / Next 元框架模块已扩展为稳定学习入口。
- Java 和 Go 已作为独立后端语言模块接入，分别覆盖企业后端和云原生后端的核心学习链路。
- AI 工程模块已经覆盖 LLM API、结构化输出、函数调用、多模态、RAG、MCP、Agent、产品协作、评测和上线治理。
- 实战项目模块已经补齐项目阶段任务、联调排查、权限案例、权限运营、组织架构、审批流、文件中心、数据看板、多租户权限、消息通知、数据导入导出、支付订单、会员订阅、搜索中心、任务调度、消息队列、开放平台、工作流配置器、低代码流程平台、审计中心、运营活动、财务对账、渠道结算、渠道费用稽核、渠道费用 ROI 复盘、渠道费用预算优化、渠道费用异常预警、渠道费用策略灰度、渠道策略效果复盘、渠道策略对照实验、渠道策略版本治理、渠道策略审批矩阵、渠道策略发布审计、渠道策略回滚治理、渠道策略异常仲裁、渠道策略仲裁复盘、渠道策略裁决标准库、渠道策略标准效果监控、渠道策略标准灰度发布、渠道策略标准版本回滚、渠道策略标准回滚演练、渠道价格稽核、渠道窜货监控、渠道信用评级、主数据管理、客户主数据、低代码表单、报表配置器、智能报表与 BI、客服工单、客服质检、客服知识运营、集团系统集成、国际化后台、数据治理平台、数据质量专项、数据资产运营、数据安全运营、规则引擎、灰度发布后台、跨区域灾备、风控中心、合同管理、合同履约、合同付款、合同变更、知识库平台、统一配置中心、行业合规审计、客户成功平台、客户生命周期价值分析、客户流失预警、客户分群运营、客户触达自动化、工单自动化、计费中台、数据交换平台、企业门户、资产管理、预算管理、资金计划、费用报销、员工借款、税务管理、采购管理、采购寻源、供应商绩效、供应链计划、项目管理、研发需求池、报价中心、上线事故案例库、运维值班、库存管理、备件库存、备件补货、售后备件周转分析、仓储物流、售后服务、售后知识自动推荐、售后知识质量治理、售后知识智能检索优化、售后知识问答助手、售后知识自动质检、售后知识专家审核、售后知识发布灰度、售后知识回滚治理、售后知识影响追踪、售后知识客户通知治理、售后知识外部服务商通知协同、售后知识服务商培训闭环、售后知识培训效果复盘、售后知识培训认证治理、售后知识认证派单联动、售后知识认证质量稽核、售后知识认证风险画像、售后结算、售后 SLA 赔付分析、售后服务商评级、售后维修质量复盘、报修派单、服务网点、数据权限审计、门店零售管理、CRM 销售管理、客户账期、客户回款风险预测、客户坏账处置策略、客户应收催收自动化、销售回款预测调度、销售现金流预警、销售回款策略模拟、销售风险动作编排、销售风险处置复盘、销售风险预案演练、销售风险指标治理、销售风险指标血缘审计、销售风险指标异常根因、销售风险指标自动修复、销售风险指标治理成熟度、销售风险指标治理运营看板、销售风险指标治理成本收益评估、会员营销、生产制造、生产排程、产能负荷预测、生产计划达成分析、生产现场安全隐患、生产安全培训闭环、生产安全考试认证、生产安全风险画像、生产安全应急演练、生产安全事故复盘、生产安全风险整改复查、生产安全整改看板、生产安全整改 SLA、生产安全整改成本复盘、生产安全整改预算预测、生产安全整改资源排期、生产安全整改产线影响评估、生产安全整改多方案决策、生产安全整改决策复盘、生产安全整改决策知识库、生产安全整改决策智能推荐、质量追溯、生产瓶颈分析、生产换型损失分析、IoT 设备管理、设备维保、教育培训平台、交付检查清单、故障复盘和分类问题库。
- 速查手册已经覆盖前端开发、工程部署、联调协作、Linux、Redis、正则、调试工具和 SQL 高频回查场景。
- 后续可以按统一模板和路线治理规则扩展更多技术。
