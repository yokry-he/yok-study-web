# 真实项目问题库

## 这个页面解决什么

这个页面是项目问题库总入口。它先记录 Vue 前端和后台项目中最常见、最容易复现的问题，再把问题按领域拆到独立页面，方便后续持续补充。

每个问题都按下面格式整理：

```text
问题现象
影响范围
根因分析
解决方案
预防方式
```

它适合在项目开发、上线、排错时快速查阅。

如果你还不确定问题属于前端、接口、状态、权限、缓存还是部署，先看 [项目排障方法论](/projects/debugging-playbook)。如果问题已经确定发生在前端项目里，但不知道该查 Vue、请求、CSS、TypeScript、浏览器还是工程化，进入 [前端项目排障图谱](/projects/frontend-debugging-map)。它会帮你从现象收集证据，再进入下面的具体分类。

## 分类入口

如果你已经知道问题属于哪个方向，优先进入对应分类：

| 分类 | 适合排查的问题 | 入口 |
| --- | --- | --- |
| 项目排障方法 | 不知道问题在哪一层、缺少复现证据、修复后不知道如何回归和沉淀 | [项目排障方法论](/projects/debugging-playbook) |
| 前端项目排障图谱 | 白屏、数据旧值、权限错位、样式污染、构建缓存等前端综合问题如何分流到对应专题 | [前端项目排障图谱](/projects/frontend-debugging-map) |
| Vue 项目专项 | 动态菜单刷新丢失、动态路由 404、Pinia 解构失去响应式、弹窗污染列表、KeepAlive 缓存、权限按钮错位 | [Vue 真实项目问题](/projects/issues-vue) |
| Vue Admin 请求权限 | 401/403、旧请求覆盖、重复提交、数据范围裁剪、导出不一致、trace id 丢失 | [Vue Admin 请求、权限与数据问题排查专题](/projects/issues-vue-admin-request) |
| Vue Admin 消息通知 | 未读数不准、重复通知、切换账号污染、实时连接重连、通知跳转 404/403、全部已读回退 | [Vue Admin 消息通知、未读数与实时提醒问题排查专题](/projects/issues-vue-admin-notification) |
| JavaScript 专项 | 数字精度、日期时区、异步循环、旧请求覆盖、事件监听泄漏、深拷贝、数组变异、本地缓存和动态正则 | [JavaScript 真实项目问题库](/projects/issues-javascript) |
| CSS 专项 | 横向溢出、组件库污染、固定元素变形、表格操作列压缩、长文本撑破、弹层遮挡、移动端导航和动画卡顿 | [CSS 真实项目问题库](/projects/issues-css) |
| HTML 与无障碍 | 鼠标可点但键盘不可用、表单标签与错误、重复提交、图片跳动和过大、焦点丢失、资源路径与深层 URL | [HTML 与无障碍真实项目问题库](/projects/issues-html-accessibility) |
| React 专项 | Strict Mode 重复请求、Effect 循环、旧闭包、请求竞态、key 错行、派生状态、表单污染、Context 性能、路由刷新、401/403、权限绕过和 chunk 404 | [React 真实项目问题库](/projects/issues-react) |
| 前端页面与状态 | 搜索旧结果、弹窗污染列表、动态路由丢失、React 重复请求、组件库升级错位、列表卡顿 | [前端页面与状态问题](/projects/issues-frontend) |
| Nuxt / Next 元框架 | hydration、重复取数、旧缓存、用户数据串号、SSR 会话、接口越权、重定向、静态部署和多实例一致性 | [Nuxt / Next 真实项目问题库](/projects/issues-meta-frameworks) |
| TypeScript 类型边界 | DTO 泄漏到页面、表单和提交参数混用、权限码拼错、RouteMeta 未声明、Store 类型污染、typecheck 缺失 | [TypeScript 类型边界问题](/projects/issues-typescript) |
| 前端工程化专项 | 依赖安装、环境变量、CI、构建、部署缓存、组件库升级、Monorepo、包体积和回滚失败 | [前端工程化真实项目问题库](/projects/issues-engineering) |
| Node.js 专项 | ESM/CJS、事件循环阻塞、线程池饱和、Stream 背压、异步上下文、进程退出、open handles 和多实例状态 | [Node.js 真实项目问题库](/projects/issues-node) |
| Java 专项 | 字节码、类路径、Bean、事务代理、JPA、连接池、线程、GC、Metaspace 和优雅停机 | [Java 真实项目问题库](/projects/issues-java) |
| Go 专项 | typed nil、slice 持有、map 竞态、goroutine 泄漏、channel 关闭、Context、连接池、事务、乐观锁和优雅停机 | [Go 真实项目问题库](/projects/issues-go) |
| 后端接口与服务 | 参数接收失败、重复提交、接口偶发慢、错误码不清晰、权限半更新、401/403 混乱 | [后端接口与服务问题](/projects/issues-backend) |
| 前后端联调 | 参数位置、认证信息、分页结构、枚举、文件上传和导出问题 | [前后端联调排查](/projects/integration-debugging) |
| 数据库与缓存 | 慢查询、事务失效、权限缓存、迁移发布顺序、复合索引、N+1 查询、Redis 内存上涨 | [数据库与缓存问题](/projects/issues-database) |
| 部署、缓存与 DevOps | 二级路由 404、旧页面缓存、环境变量错误、Nginx 代理路径、资源 404、CI 旧包、回滚兼容 | [部署、缓存与 DevOps 问题](/projects/issues-deployment) |
| AI 工程 | 输出不稳定、RAG 答非所问、成本升高、资料越权、引用不可信、Agent 工具循环 | [AI 工程问题](/projects/issues-ai) |
| 故障复盘 | 权限越权、缓存事故、慢查询、AI 越权等问题如何沉淀改进项 | [故障复盘模板](/projects/incident-review) |

如果你还不确定问题在哪一层，建议按这个顺序排查：

```text
浏览器 Console 和 Network
↓
前端状态和路由
↓
接口请求与响应
↓
后端日志和错误码
↓
数据库、缓存和第三方服务
↓
部署层、CDN、Nginx、容器和环境变量
```

如果你不是在排查已经发生的问题，而是在上线前做检查，先看 [项目交付检查清单](/projects/delivery-checklist)。如果你还在拆项目阶段，先看 [项目阶段任务](/projects/project-stage-tasks)。

## 复杂项目案例入口

如果你不是在查某一个问题，而是在准备做一个完整后台模块，优先看下面这些案例。它们会把业务目标、数据模型、页面拆分、接口链路、权限边界和常见问题串起来：

| 案例 | 适合学习的内容 | 入口 |
| --- | --- | --- |
| 权限运营 | 权限申请、风险审批、临时授权、复核回收和敏感操作审计 | [权限运营项目案例](/projects/permission-operation-case) |
| 组织架构 | 部门树、岗位、员工、直属上级、数据范围权限、离职交接 | [组织架构项目案例](/projects/organization-case) |
| 审批流 | 流程定义、节点实例、审批动作、并发审批、审批详情页 | [审批流项目案例](/projects/approval-workflow-case) |
| 文件中心 | 上传链路、对象存储、临时签名 URL、文件权限、无主文件清理 | [文件中心项目案例](/projects/file-center-case) |
| 数据看板 | 指标口径、图表请求、缓存策略、权限过滤、数据不一致排查 | [数据看板项目案例](/projects/analytics-dashboard-case) |
| 多租户权限 | 租户上下文、成员角色、数据隔离、缓存隔离、跨租户审计 | [多租户权限项目案例](/projects/multi-tenant-permission-case) |
| 消息通知 | 站内信、消息模板、接收人计算、未读数量、发送重试、实时提醒排障 | [消息通知项目案例](/projects/notification-center-case)、[Vue Admin 消息通知排障](/projects/issues-vue-admin-notification) |
| 数据导入导出 | 模板下载、字段校验、错误行、异步任务、权限范围导出 | [数据导入导出项目案例](/projects/import-export-case) |
| 支付订单 | 订单状态机、支付回调、幂等、退款、对账、金额精度 | [支付订单项目案例](/projects/payment-order-case) |
| 会员订阅 | 套餐权益、试用期、续费、到期降级、权益校验 | [会员订阅项目案例](/projects/subscription-billing-case) |
| 搜索中心 | 全局搜索、索引同步、权限过滤、排序高亮、搜索日志 | [搜索中心项目案例](/projects/search-center-case) |
| 任务调度 | 定时任务、分布式锁、执行记录、失败重试、任务告警 | [任务调度项目案例](/projects/task-scheduler-case) |
| 消息队列 | 事件发布、消费幂等、失败重试、死信队列、积压监控 | [消息队列项目案例](/projects/message-queue-case) |
| 第三方开放平台 | 应用密钥、签名、scope、限流、Webhook、API 版本治理 | [第三方开放平台项目案例](/projects/open-platform-case) |
| 工作流配置器 | 可视化流程、草稿发布、版本运行、条件节点、执行日志 | [工作流配置器项目案例](/projects/workflow-builder-case) |
| 低代码流程平台 | 业务事件、系统节点、人工任务、条件分支、补偿和运行轨迹 | [低代码流程平台项目案例](/projects/low-code-workflow-case) |
| 审计中心 | 操作日志、变更明细、敏感脱敏、风险规则、审计导出 | [审计中心项目案例](/projects/audit-center-case) |
| 运营活动 | 活动规则、资格校验、库存扣减、奖励发放、活动复盘 | [运营活动项目案例](/projects/marketing-campaign-case) |
| 复杂财务对账 | 渠道账单、本地快照、差异处理、财务导出、资金审计 | [复杂财务对账项目案例](/projects/finance-reconciliation-case) |
| 渠道结算 | 渠道档案、分润规则、账单生成、对账确认和付款归档 | [渠道结算项目案例](/projects/channel-settlement-case) |
| 渠道费用稽核 | 费用政策、预算、发票证据、稽核规则、异常扣减、结算和效果复盘 | [渠道费用稽核项目案例](/projects/channel-expense-audit-case) |
| 渠道费用 ROI 复盘 | 费用投入、销售基线、增量归因、毛利 ROI、政策反馈和预算优化 | [渠道费用 ROI 复盘项目案例](/projects/channel-expense-roi-review-case) |
| 渠道费用预算优化 | 预算池、分配规则、ROI 权重、预算模拟、执行监控和复盘调优 | [渠道费用预算优化项目案例](/projects/channel-expense-budget-optimization-case) |
| 渠道费用异常预警 | 费用申请、异常规则、预警事件、阻断补证、复核和渠道画像 | [渠道费用异常预警项目案例](/projects/channel-expense-anomaly-warning-case) |
| 渠道费用策略灰度 | 策略版本、灰度范围、命中条件、风险阈值、回滚和发布复盘 | [渠道费用策略灰度项目案例](/projects/channel-expense-strategy-gray-release-case) |
| 渠道策略效果复盘 | 策略目标、费用投入、销售产出、基线对比、异常剔除和复盘动作 | [渠道策略效果复盘项目案例](/projects/channel-strategy-effect-review-case) |
| 渠道策略对照实验 | 实验假设、样本分组、对照组、护栏指标、效果归因和推广决策 | [渠道策略对照实验项目案例](/projects/channel-strategy-ab-experiment-case) |
| 渠道策略版本治理 | 策略版本、变更差异、影响分析、审批发布、执行追踪和回滚治理 | [渠道策略版本治理项目案例](/projects/channel-strategy-version-governance-case) |
| 渠道策略审批矩阵 | 策略变更、影响分析、审批条件、矩阵匹配、审批任务和审计追踪 | [渠道策略审批矩阵项目案例](/projects/channel-strategy-approval-matrix-case) |
| 渠道策略发布审计 | 策略版本、审批记录、发布范围、审计清单、命中监控和回滚归档 | [渠道策略发布审计项目案例](/projects/channel-strategy-release-audit-case) |
| 渠道策略回滚治理 | 策略异常、回滚版本、影响评估、规则刷新、补偿任务和回滚复盘 | [渠道策略回滚治理项目案例](/projects/channel-strategy-rollback-governance-case) |
| 渠道策略异常仲裁 | 异常事件、争议申请、证据归集、规则复算、仲裁会审和裁决执行 | [渠道策略异常仲裁项目案例](/projects/channel-strategy-exception-arbitration-case) |
| 渠道策略仲裁复盘 | 已结案仲裁、根因归类、证据完整度、裁决一致性、执行闭环和优化任务 | [渠道策略仲裁复盘项目案例](/projects/channel-strategy-arbitration-review-case) |
| 渠道策略裁决标准库 | 裁决标准、适用条件、证据门槛、规则版本、案件引用和偏离审计 | [渠道策略裁决标准库项目案例](/projects/channel-strategy-decision-standard-library-case) |
| 渠道策略标准效果监控 | 标准引用率、偏离率、处理效率、申诉改善、健康度评分和修订建议 | [渠道策略标准效果监控项目案例](/projects/channel-strategy-standard-effect-monitoring-case) |
| 渠道策略标准灰度发布 | 标准版本、灰度范围、适用对象、护栏指标、回滚和效果监控 | [渠道策略标准灰度发布项目案例](/projects/channel-strategy-standard-gray-release-case) |
| 渠道策略标准版本回滚 | 异常版本、影响案件、回滚范围、缓存刷新、补偿处理和回滚复盘 | [渠道策略标准版本回滚项目案例](/projects/channel-strategy-standard-version-rollback-case) |
| 渠道策略标准回滚演练 | 演练计划、沙箱数据、异常场景、验证点、缺口整改和发布门禁 | [渠道策略标准回滚演练项目案例](/projects/channel-strategy-standard-rollback-drill-case) |
| 渠道策略标准灾备切换 | 灾备快照、切换策略、只读裁决、引用暂存、回切校验和故障复盘 | [渠道策略标准灾备切换项目案例](/projects/channel-strategy-standard-disaster-recovery-switch-case) |
| 渠道价格稽核 | 价格政策、授权价、合同价、成交价、低价异常、补证和处罚 | [渠道价格稽核项目案例](/projects/channel-price-audit-case) |
| 渠道窜货监控 | 授权区域、商品码、流向证据、跨区销售、核查单、申诉和处罚 | [渠道窜货监控项目案例](/projects/channel-diversion-monitor-case) |
| 渠道信用评级 | 回款逾期、风险事件、信用评分、额度账期、订单控制和整改复核 | [渠道信用评级项目案例](/projects/channel-credit-rating-case) |
| 渠道返利风控 | 返利政策、基数快照、退货冲量、窜货低价、冻结申诉和结算扣减 | [渠道返利风控项目案例](/projects/channel-rebate-risk-control-case) |
| 渠道政策模拟 | 政策草稿、模拟样本、成本收益、风险扫描、方案对比和偏差复盘 | [渠道政策模拟项目案例](/projects/channel-policy-simulation-case) |
| 渠道利润模拟 | 渠道产品价格、折扣返利、成本构成、低毛利风险、方案对比和偏差复盘 | [渠道利润模拟项目案例](/projects/channel-profit-simulation-case) |
| 渠道价格弹性分析 | 价格销量样本、弹性估算、调价模拟、毛利影响、试点监控和偏差复盘 | [渠道价格弹性分析项目案例](/projects/channel-price-elasticity-analysis-case) |
| 主数据管理 | 统一编码、质量规则、重复合并、版本快照、系统同步 | [主数据管理项目案例](/projects/master-data-case) |
| 客户主数据 | 统一客户身份、来源映射、字段可信来源、去重合并和同步审计 | [客户主数据项目案例](/projects/customer-master-data-case) |
| 低代码表单 | 表单设计器、字段权限、条件显隐、版本兼容、数据收集 | [低代码表单项目案例](/projects/low-code-form-case) |
| 报表配置器 | 数据集、指标维度、查询生成器、权限过滤、导出任务 | [报表配置器项目案例](/projects/report-builder-case) |
| 智能报表与 BI | 指标语义层、自然语言问数、异常洞察、AI 解释和证据引用 | [智能报表与 BI 分析项目案例](/projects/smart-bi-dashboard-case) |
| 客服工单 | 工单流转、分派规则、SLA、内部备注、满意度统计 | [客服工单项目案例](/projects/support-ticket-case) |
| 客服质检 | 会话归档、抽检策略、规则评分、复核申诉、整改闭环 | [客服质检项目案例](/projects/customer-service-quality-case) |
| 集团级系统集成 | 外部系统、字段映射、协议转换、补偿处理、集成监控 | [集团级系统集成项目案例](/projects/enterprise-integration-case) |
| 国际化后台 | 多语言词条、时区币种、通知模板、导出语言、缺失翻译 | [国际化后台项目案例](/projects/i18n-admin-case) |
| 数据治理平台 | 数据目录、血缘、质量规则、敏感分级、权限申请 | [数据治理平台项目案例](/projects/data-governance-case) |
| 数据质量专项 | 质量规则、检查任务、问题样本、整改复查和发布门禁 | [数据质量专项项目案例](/projects/data-quality-special-case) |
| 数据资产运营 | 资产目录、资产申请、使用统计、健康度、热度和下线影响 | [数据资产运营项目案例](/projects/data-asset-operation-case) |
| 数据安全运营 | 敏感识别、脱敏策略、访问审计、异常检测和整改闭环 | [数据安全运营项目案例](/projects/data-security-operation-case) |
| 规则引擎 | 规则集、条件动作、试算、灰度、命中日志、回滚 | [规则引擎项目案例](/projects/rule-engine-case) |
| 灰度发布后台 | 功能开关、分流规则、发布审批、命中日志、指标回滚 | [灰度发布后台项目案例](/projects/gray-release-admin-case) |
| 跨区域灾备 | RTO/RPO、备份复制、演练、故障切换、恢复报告 | [跨区域灾备管理项目案例](/projects/disaster-recovery-case) |
| 风控中心 | 风险策略、名单库、实时决策、命中日志、人工审核 | [风控中心项目案例](/projects/risk-control-center-case) |
| 合同管理 | 合同生命周期、模板、审批、电子签署、归档和到期提醒 | [合同管理项目案例](/projects/contract-management-case) |
| 合同履约 | 履约节点、交付验收、收付款条件、延期变更和风险跟踪 | [合同履约项目案例](/projects/contract-fulfillment-case) |
| 合同付款 | 付款条款、验收发票、付款申请、资金排期、执行和风险 | [合同付款项目案例](/projects/contract-payment-case) |
| 合同变更 | 金额、期限、付款条款、履约范围、补充协议和版本影响 | [合同变更项目案例](/projects/contract-change-case) |
| 合同续签 | 到期提醒、续签机会、客户评估、续签审批、签署和续签率分析 | [合同续签项目案例](/projects/contract-renewal-case) |
| 客户合同风险预警 | 合同履约、回款逾期、条款风险、客户风险、处置任务和风险复盘 | [客户合同风险预警项目案例](/projects/customer-contract-risk-warning-case) |
| 客户合同收入预测 | 合同收入项、确认规则、履约验收、预测版本、风险概率和偏差复盘 | [客户合同收入预测项目案例](/projects/customer-contract-revenue-forecast-case) |
| 知识库平台 | 知识发布、分类标签、权限范围、搜索、反馈和版本 | [知识库平台项目案例](/projects/knowledge-base-case) |
| 客服知识运营 | 知识缺口、坐席辅助、机器人问答、效果看板和优化任务 | [客服知识运营项目案例](/projects/customer-knowledge-operation-case) |
| 统一配置中心 | 配置作用域、版本发布、灰度、回滚、客户端缓存 | [统一配置中心项目案例](/projects/config-center-case) |
| 行业合规审计 | 合规控制项、证据采集、整改闭环、审计报告 | [行业合规审计项目案例](/projects/compliance-audit-case) |
| 客户成功平台 | 客户健康度、续费机会、流失预警、跟进任务、客户触达 | [客户成功平台项目案例](/projects/customer-success-case) |
| 客户生命周期价值分析 | 收入毛利、服务成本、留存概率、客户分层、运营策略和效果回写 | [客户生命周期价值分析项目案例](/projects/customer-lifetime-value-analysis-case) |
| 客户流失预警 | 使用下降、合同到期、投诉、回款异常、风险评分、跟进任务和挽回复盘 | [客户流失预警项目案例](/projects/customer-churn-warning-case) |
| 客户续费挽回 | 到期客户、健康度、风险信号、挽回任务、报价续签和续费复盘 | [客户续费挽回项目案例](/projects/customer-renewal-recovery-case) |
| 客户续约定价策略 | 续费机会、客户价值、定价建议、毛利校验、价格审批和定价复盘 | [客户续约定价策略项目案例](/projects/customer-renewal-pricing-strategy-case) |
| 客户分群运营 | 标签、分群规则、人群包、触达计划、频控、效果回写和复盘 | [客户分群运营项目案例](/projects/customer-segmentation-operation-case) |
| 客户触达自动化 | 行为事件、触达规则、频控互斥、渠道模板、发送记录和效果回写 | [客户触达自动化项目案例](/projects/customer-touch-automation-case) |
| 客户权益运营 | 权益定义、权益包、资格规则、发放实例、核销流水、成本和效果复盘 | [客户权益运营项目案例](/projects/customer-benefit-operation-case) |
| 客户投诉闭环 | 投诉受理、分级分派、处理方案、客户回访、赔付和整改复盘 | [客户投诉闭环项目案例](/projects/customer-complaint-closed-loop-case) |
| 工单自动化 | 自动分派、SLA 升级、回复建议、重复工单、命中日志 | [工单自动化项目案例](/projects/ticket-automation-case) |
| 计费中台 | 价格版本、用量计费、账单生成、调整、支付和欠费处理 | [计费中台项目案例](/projects/billing-platform-case) |
| 数据交换平台 | 数据资源、交换申请、字段映射、脱敏、回执和审计 | [数据交换平台项目案例](/projects/data-exchange-platform-case) |
| 企业门户 | 统一入口、应用导航、待办聚合、公告、搜索和访问审计 | [企业门户项目案例](/projects/enterprise-portal-case) |
| 资产管理 | 资产入库、领用、调拨、维修、盘点、报废和审计 | [资产管理项目案例](/projects/asset-management-case) |
| 预算管理 | 年度预算、预算占用、预算执行、调整审批和超预算预警 | [预算管理项目案例](/projects/budget-management-case) |
| 资金计划 | 收付款预测、账户余额、付款排期、资金缺口和现金流复盘 | [资金计划项目案例](/projects/cash-flow-planning-case) |
| 费用报销 | 费用标准、发票验真、预算占用、审批、付款和财务入账 | [费用报销项目案例](/projects/expense-reimbursement-case) |
| 员工借款 | 备用金、差旅借款、付款、报销冲销、还款和逾期提醒 | [员工借款项目案例](/projects/employee-loan-case) |
| 税务管理 | 税率配置、发票、红冲、申报期、税务风险和财税对账 | [税务管理项目案例](/projects/tax-management-case) |
| 发票协同 | 开票申请、收票登记、验真查重、业务匹配、红冲和税务归档 | [发票协同项目案例](/projects/invoice-collaboration-case) |
| 采购管理 | 采购申请、供应商、比价、订单、验收、发票和付款 | [采购管理项目案例](/projects/procurement-management-case) |
| 采购寻源 | 供应商池、询价招标、报价轮次、评标、定标和审计 | [采购寻源项目案例](/projects/procurement-sourcing-case) |
| 供应商准入 | 注册申请、资质审核、风险检查、准入审批、准入范围和复审 | [供应商准入项目案例](/projects/supplier-onboarding-case) |
| 供应商合同协同 | 合同草稿、供应商确认、条款协商、电子签署、履约和付款控制 | [供应商合同协同项目案例](/projects/supplier-contract-collaboration-case) |
| 供应商协同门户 | 外部供应商注册、报价、接单、发货、对账、开票、索赔和门户权限隔离 | [供应商协同门户项目案例](/projects/supplier-collaboration-portal-case) |
| 供应商门户权限审计 | 门户账号、角色范围、字段可见性、访问日志、敏感导出和整改闭环 | [供应商门户权限审计项目案例](/projects/supplier-portal-permission-audit-case) |
| 供应商索赔 | 质量、交付、合同和财务差异索赔，供应商确认、申诉、扣款和绩效影响 | [供应商索赔项目案例](/projects/supplier-claim-case) |
| 供应商绩效 | 交付、质量、价格、服务、申诉、整改和等级策略 | [供应商绩效项目案例](/projects/supplier-performance-case) |
| 供应链计划 | 需求预测、补货建议、采购计划、缺货预警和计划复盘 | [供应链计划项目案例](/projects/supply-chain-planning-case) |
| 项目管理 | 项目立项、成员、里程碑、任务、风险、工时和复盘 | [项目管理项目案例](/projects/project-management-case) |
| 研发需求池 | 需求收集、去重、评审、优先级、版本排期和上线复盘 | [研发需求池项目案例](/projects/rd-requirement-pool-case) |
| 报价中心 | 价格版本、折扣审批、报价版本、有效期和转合同 | [报价中心项目案例](/projects/quotation-center-case) |
| 价格审批中心 | 标准价、底价、折扣权限、特殊价格审批、价格授权和成交价复盘 | [价格审批中心项目案例](/projects/price-approval-center-case) |
| 上线事故案例 | 白屏、401/403、慢查询、权限缓存、队列积压、AI 越权 | [上线事故案例库](/projects/production-incident-cases) |
| 运维值班 | 值班排班、告警路由、升级通知、交接班、SLO 和复盘 | [运维值班项目案例](/projects/operations-oncall-case) |
| 库存管理 | 可用库存、锁定库存、入库、出库、调拨、盘点和流水 | [库存管理项目案例](/projects/inventory-management-case) |
| 渠道库存协同 | 渠道库存快照、可用锁定库存、缺货积压预警、补货调拨和库存对账 | [渠道库存协同项目案例](/projects/channel-inventory-collaboration-case) |
| 备件库存 | 备件适配、网点库存、工程师领用、工单消耗和旧件回收 | [备件库存项目案例](/projects/spare-parts-inventory-case) |
| 备件补货 | 安全库存、消耗预测、采购补货、网点调拨和补货复盘 | [备件补货项目案例](/projects/spare-parts-replenishment-case) |
| 售后备件周转分析 | 库存金额、周转天数、缺货率、呆滞风险、调拨补货和清理建议 | [售后备件周转分析项目案例](/projects/after-sales-spare-parts-turnover-case) |
| 备件旧件返修 | 旧件回收、交接追踪、检测返修、供应商维修、翻新入库和报废 | [备件旧件返修项目案例](/projects/spare-parts-return-repair-case) |
| 仓储物流 | 入库上架、拣货、复核、打包、发货、物流轨迹和异常件 | [仓储物流项目案例](/projects/warehouse-logistics-case) |
| 售后服务 | 退货、换货、维修、退款、补发、质检和售后原因分析 | [售后服务项目案例](/projects/after-sales-service-case) |
| 售后远程诊断 | 诊断会话、故障信号、客户自检、远程命令、派单建议和误判复盘 | [售后远程诊断项目案例](/projects/after-sales-remote-diagnosis-case) |
| 售后专家协同 | 协同单、证据包、专家组、会诊结论、执行反馈和知识沉淀 | [售后专家协同项目案例](/projects/after-sales-expert-collaboration-case) |
| 售后知识自动推荐 | 工单上下文、知识召回、排序证据、使用反馈、知识缺口和质量治理 | [售后知识自动推荐项目案例](/projects/after-sales-knowledge-recommendation-case) |
| 售后知识质量治理 | 质量指标、负反馈、版本复审、冲突处理、推荐权重和治理任务 | [售后知识质量治理项目案例](/projects/after-sales-knowledge-quality-governance-case) |
| 售后知识智能检索优化 | 查询理解、同义词、语义召回、范围过滤、搜索反馈和零结果治理 | [售后知识智能检索优化项目案例](/projects/after-sales-knowledge-search-optimization-case) |
| 售后知识问答助手 | 工单上下文、知识召回、答案生成、引用校验、低置信度降级和反馈闭环 | [售后知识问答助手项目案例](/projects/after-sales-knowledge-qa-assistant-case) |
| 售后知识自动质检 | 质检规则、语义冲突、负反馈、质量评分、治理任务和重新索引 | [售后知识自动质检项目案例](/projects/after-sales-knowledge-auto-quality-inspection-case) |
| 售后知识专家审核 | 知识版本、专家池、审核任务、会审、紧急发布、索引更新和追溯 | [售后知识专家审核项目案例](/projects/after-sales-knowledge-expert-review-case) |
| 售后知识发布灰度 | 知识版本、灰度范围、搜索问答引用、反馈监控、护栏和回滚 | [售后知识发布灰度项目案例](/projects/after-sales-knowledge-release-gray-case) |
| 售后知识回滚治理 | 异常反馈、版本回滚、索引缓存刷新、影响追踪、通知和修订复审 | [售后知识回滚治理项目案例](/projects/after-sales-knowledge-rollback-governance-case) |
| 售后知识影响追踪 | 知识引用日志、搜索问答、工单使用、风险分级、复核通知和处理闭环 | [售后知识影响追踪项目案例](/projects/after-sales-knowledge-impact-trace-case) |
| 售后知识客户通知治理 | 通知必要性、客户分组、文案审批、消息回执、客户反馈和跟进闭环 | [售后知识客户通知治理项目案例](/projects/after-sales-knowledge-customer-notification-governance-case) |
| 售后知识外部服务商通知协同 | 服务商范围、通知方案、工程师确认、执行反馈、逾期升级和协同归档 | [售后知识外部服务商通知协同项目案例](/projects/after-sales-knowledge-provider-notification-collaboration-case) |
| 售后知识服务商培训闭环 | 知识变更、课程包、学习任务、考试认证、现场验证和复训闭环 | [售后知识服务商培训闭环项目案例](/projects/after-sales-knowledge-provider-training-closed-loop-case) |
| 售后知识培训效果复盘 | 培训批次、效果窗口、现场工单、质检投诉、知识点归因和改进动作 | [售后知识培训效果复盘项目案例](/projects/after-sales-knowledge-training-effect-review-case) |
| 售后知识培训认证治理 | 认证规则、工程师资质、派单资格、质量联动、暂停恢复和续期复训 | [售后知识培训认证治理项目案例](/projects/after-sales-knowledge-training-certification-governance-case) |
| 售后知识认证派单联动 | 工程师认证、派单资格、工单匹配、异常放行、暂停恢复和审计 | [售后知识认证派单联动项目案例](/projects/after-sales-knowledge-certification-dispatch-linkage-case) |
| 售后知识认证质量稽核 | 服务质量信号、稽核规则、认证暂停、整改复训、恢复复核和质量看板 | [售后知识认证质量稽核项目案例](/projects/after-sales-knowledge-certification-quality-audit-case) |
| 售后知识认证风险画像 | 认证覆盖、质量表现、派单风险、复训结果、服务商缺口和预警任务 | [售后知识认证风险画像项目案例](/projects/after-sales-knowledge-certification-risk-profile-case) |
| 售后知识认证服务商整改 | 风险画像、整改计划、服务商执行、证据审核、整改验收和评级派单联动 | [售后知识认证服务商整改项目案例](/projects/after-sales-knowledge-certification-provider-rectification-case) |
| 客户退换货质检 | 退换申请、收货登记、外观功能质检、责任判定、退款换货和库存去向 | [客户退换货质检项目案例](/projects/customer-return-quality-inspection-case) |
| 客户退款风控 | 可退金额锁定、退款风险规则、人工审核、渠道退款、权益回滚和审计 | [客户退款风控项目案例](/projects/customer-refund-risk-control-case) |
| 售后结算 | 退款、维修收费、备件成本、服务商结算、对账和调整 | [售后结算项目案例](/projects/after-sales-settlement-case) |
| 现场服务收费 | 上门费、工时费、备件费、客户确认、收款、退款和服务结算 | [现场服务收费项目案例](/projects/field-service-charging-case) |
| 售后备件成本核算 | 工单备件消耗、成本取价、旧件抵减、供应商索赔和售后毛利分析 | [售后备件成本核算项目案例](/projects/after-sales-spare-part-cost-case) |
| 售后成本毛利分析 | 工单收入、备件人工服务商成本、赔付索赔、毛利异常和策略复盘 | [售后成本毛利分析项目案例](/projects/after-sales-cost-margin-case) |
| 售后服务成本优化 | 工单成本、备件差旅、服务商结算、重复上门、优化任务和体验复盘 | [售后服务成本优化项目案例](/projects/after-sales-service-cost-optimization-case) |
| 售后 SLA 赔付分析 | 服务等级、响应到场修复计时、违约原因、赔付规则和责任分摊 | [售后 SLA 赔付分析项目案例](/projects/after-sales-sla-compensation-case) |
| 售后服务商评级 | 响应、到场、一次修复、满意度、成本、投诉、派单策略和整改淘汰 | [售后服务商评级项目案例](/projects/after-sales-provider-rating-case) |
| 售后维修质量复盘 | 一次修复、返修、重复故障、根因归因、整改验证和产品反馈 | [售后维修质量复盘项目案例](/projects/after-sales-repair-quality-review-case) |
| 售后投诉根因分析 | 投诉分级、关联订单产品服务、根因责任、整改验证、赔付和复发监控 | [售后投诉根因分析项目案例](/projects/after-sales-complaint-root-cause-case) |
| 报修派单 | 设备报修、派单规则、工程师调度、备件、SLA 和回访 | [报修派单项目案例](/projects/repair-dispatch-case) |
| 服务网点 | 服务区域、网点能力、工程师、备件、SLA 和服务质量 | [服务网点项目案例](/projects/service-outlet-case) |
| 数据权限审计 | 数据访问、敏感字段、导出审计、权限复核和越权检测 | [数据权限审计项目案例](/projects/data-permission-audit-case) |
| 门店零售管理 | 门店、商品、库存、收银、会员、促销、日结和经营看板 | [门店零售管理项目案例](/projects/retail-store-management-case) |
| CRM 销售管理 | 线索、客户、商机、跟进、销售漏斗、报价和回款预测 | [CRM 销售管理项目案例](/projects/crm-sales-management-case) |
| 客户账期 | 客户授信、额度占用、应收账款、回款核销、逾期和催收 | [客户账期项目案例](/projects/customer-credit-term-case) |
| 客户授信风控 | 授信额度、账期、额度占用、逾期冻结、临时额度、复审和审计 | [客户授信风控项目案例](/projects/customer-credit-risk-control-case) |
| 客户回款风险预测 | 应收账款、风险评分、风险原因、处置任务、承诺回款和偏差复盘 | [客户回款风险预测项目案例](/projects/customer-payment-risk-prediction-case) |
| 客户坏账处置策略 | 长期逾期、催收证据、坏账评估、计提核销、追偿和信用回写 | [客户坏账处置策略项目案例](/projects/customer-bad-debt-disposal-case) |
| 客户应收催收自动化 | 应收账龄、客户分层、催收策略、触达记录、承诺回款和升级复盘 | [客户应收催收自动化项目案例](/projects/customer-receivable-collection-automation-case) |
| 销售回款预测调度 | 应收计划、到账概率、回款日历、调度任务、销售承诺和偏差复盘 | [销售回款预测调度项目案例](/projects/sales-payment-prediction-scheduling-case) |
| 销售现金流预警 | 回款预测、支出计划、安全线、缺口日期、预警分级、处置任务和复盘 | [销售现金流预警项目案例](/projects/sales-cash-flow-warning-case) |
| 销售回款策略模拟 | 样本选择、策略参数、试算对比、风险评估、执行成本和采纳决策 | [销售回款策略模拟项目案例](/projects/sales-collection-strategy-simulation-case) |
| 销售风险动作编排 | 风险信号、动作模板、任务派发、执行反馈、风险升级和效果回写 | [销售风险动作编排项目案例](/projects/sales-risk-action-orchestration-case) |
| 销售风险处置复盘 | 处置结果、复盘批次、效果指标、归因标签、策略建议和版本迭代 | [销售风险处置复盘项目案例](/projects/sales-risk-disposal-review-case) |
| 销售风险预案演练 | 风险场景、预案模板、模拟事件、演练任务、缺口整改和预案升级 | [销售风险预案演练项目案例](/projects/sales-risk-contingency-drill-case) |
| 销售风险指标治理 | 指标定义、口径版本、数据血缘、质量校验、使用审计和指标下线 | [销售风险指标治理项目案例](/projects/sales-risk-metric-governance-case) |
| 销售风险指标血缘审计 | 上游来源、加工链路、下游引用、运行质量、影响评估和审计报告 | [销售风险指标血缘审计项目案例](/projects/sales-risk-metric-lineage-audit-case) |
| 销售风险指标异常根因 | 质量告警、影响范围、血缘定位、根因分类、修复验证和预防规则 | [销售风险指标异常根因项目案例](/projects/sales-risk-metric-anomaly-root-cause-case) |
| 销售风险指标自动修复 | 修复策略、护栏校验、重跑回补、指标重算、质量复验和下游刷新 | [销售风险指标自动修复项目案例](/projects/sales-risk-metric-auto-repair-case) |
| 销售风险指标治理成熟度 | 治理等级、能力评分、差距分析、改进路线、复评和成熟度看板 | [销售风险指标治理成熟度项目案例](/projects/sales-risk-metric-governance-maturity-case) |
| 销售风险指标治理运营看板 | 健康快照、成熟度分布、质量异常、任务 SLA、风险下钻和运营复盘 | [销售风险指标治理运营看板项目案例](/projects/sales-risk-metric-governance-operations-dashboard-case) |
| 销售风险指标治理成本收益评估 | 治理成本、收益证据、ROI 计算、优先级排序、预算建议和复核记录 | [销售风险指标治理成本收益评估项目案例](/projects/sales-risk-metric-governance-cost-benefit-evaluation-case) |
| 销售风险指标治理预算审批 | ROI 评估、预算明细、收益证据、审批矩阵、预算占用和执行复盘 | [销售风险指标治理预算审批项目案例](/projects/sales-risk-metric-governance-budget-approval-case) |
| 销售回款计划 | 应收、回款计划、收款认领、核销、逾期催收和现金流预测 | [销售回款计划项目案例](/projects/sales-collection-plan-case) |
| 销售预测复盘 | 商机预测、预测版本、主管评审、实际回写、准确率和偏差原因复盘 | [销售预测复盘项目案例](/projects/sales-forecast-review-case) |
| 销售目标拆解 | 公司目标、区域团队个人配额、过程指标、目标确认、达成跟踪和复盘 | [销售目标拆解项目案例](/projects/sales-target-breakdown-case) |
| 销售佣金核算 | 计佣口径、归属规则、阶梯提成、退款冲减、销售确认和财务发放 | [销售佣金核算项目案例](/projects/sales-commission-settlement-case) |
| 销售返利政策 | 返利政策、适用对象、达成口径、阶梯规则、扣减、确认和发放 | [销售返利政策项目案例](/projects/sales-rebate-policy-case) |
| 会员营销 | 会员等级、积分、优惠券、人群包、触达和效果分析 | [会员营销项目案例](/projects/member-marketing-case) |
| 生产制造 | BOM、生产计划、工单、工序、报工、质检和成品入库 | [生产制造项目案例](/projects/manufacturing-execution-case) |
| 生产排程 | 工单排程、设备产能、物料齐套、插单重排和影响分析 | [生产排程项目案例](/projects/production-scheduling-case) |
| 产能负荷预测 | 需求预测、订单、工艺路线、产能日历、负荷率、缺口和调度建议 | [产能负荷预测项目案例](/projects/capacity-load-forecast-case) |
| 生产计划达成分析 | 计划版本、工单执行、报工入库、达成率、偏差原因和改善任务 | [生产计划达成分析项目案例](/projects/production-plan-attainment-case) |
| 生产停线损失复盘 | 停线事件、影响工单、损失计算、根因责任、改善任务和复发验证 | [生产停线损失复盘项目案例](/projects/production-line-stop-loss-review-case) |
| 质量追溯 | 原料批次、生产过程、质检结果、发货流向和召回分析 | [质量追溯项目案例](/projects/quality-traceability-case) |
| 生产质量异常 | 异常上报、批次隔离、根因分析、处置方案、CAPA 和复发监控 | [生产质量异常项目案例](/projects/production-quality-exception-case) |
| 生产异常 CAPA | 异常分级、影响范围、临时处置、根因分析、纠正预防和有效性验证 | [生产异常 CAPA 项目案例](/projects/production-exception-capa-case) |
| 生产过程审核 | 审核计划、检查表、现场问题、整改任务、复查验证和 CAPA 升级 | [生产过程审核项目案例](/projects/production-process-audit-case) |
| 生产巡检移动端 | 巡检计划、扫码定位、移动表单、离线同步、异常处理和巡检看板 | [生产巡检移动端项目案例](/projects/production-mobile-inspection-case) |
| 生产现场安全隐患 | 隐患上报、风险分级、临时控制、整改复查、重复隐患和安全培训 | [生产现场安全隐患项目案例](/projects/production-safety-hazard-case) |
| 生产安全培训闭环 | 培训需求、岗位资质、考试实操、上岗限制、补训和效果复盘 | [生产安全培训闭环项目案例](/projects/production-safety-training-closed-loop-case) |
| 生产安全考试认证 | 认证要求、题库考试、实操验证、证书有效期、上岗授权和复审 | [生产安全考试认证项目案例](/projects/production-safety-exam-certification-case) |
| 生产安全风险画像 | 隐患、巡检、培训、设备、违规记录、风险评分、预警和整改闭环 | [生产安全风险画像项目案例](/projects/production-safety-risk-profile-case) |
| 生产安全应急演练 | 应急预案、演练计划、现场签到、过程记录、评估问题和整改闭环 | [生产安全应急演练项目案例](/projects/production-safety-emergency-drill-case) |
| 生产安全事故复盘 | 事故上报、应急控制、调查取证、根因分析、整改验证和经验沉淀 | [生产安全事故复盘项目案例](/projects/production-safety-incident-review-case) |
| 生产安全风险整改复查 | 风险问题、整改任务、证据检查、现场复查、退回升级和风险回写 | [生产安全风险整改复查项目案例](/projects/production-safety-risk-rectification-review-case) |
| 生产安全整改看板 | 整改任务、指标快照、高风险未闭环、重复隐患、督办升级和趋势复盘 | [生产安全整改看板项目案例](/projects/production-safety-rectification-dashboard-case) |
| 生产安全整改 SLA | 整改时限、SLA 规则、预警提醒、逾期升级、暂停计时和结果归档 | [生产安全整改 SLA 项目案例](/projects/production-safety-rectification-sla-case) |
| 生产安全整改成本复盘 | 整改成本、停线损失、风险效果、重复隐患成本、预算占用和投入建议 | [生产安全整改成本复盘项目案例](/projects/production-safety-rectification-cost-review-case) |
| 生产安全整改预算预测 | 成本基线、风险任务池、预算版本、优先级排序、执行跟踪和偏差分析 | [生产安全整改预算预测项目案例](/projects/production-safety-rectification-budget-forecast-case) |
| 生产安全整改资源排期 | 整改任务池、资源需求、产线窗口、冲突检查、排期发布和执行跟踪 | [生产安全整改资源排期项目案例](/projects/production-safety-rectification-resource-scheduling-case) |
| 生产安全整改产线影响评估 | 整改任务、产线依赖、订单影响、停线方案、风险成本权衡和排期约束 | [生产安全整改产线影响评估项目案例](/projects/production-safety-rectification-line-impact-assessment-case) |
| 生产安全整改多方案决策 | 候选方案、风险测算、产能成本交付评分、多角色评审和执行复盘 | [生产安全整改多方案决策项目案例](/projects/production-safety-rectification-multi-scenario-decision-case) |
| 生产安全整改决策复盘 | 决策方案、预测实际差异、风险成本交付偏差、根因和模型改进 | [生产安全整改决策复盘项目案例](/projects/production-safety-rectification-decision-review-case) |
| 生产安全整改决策知识库 | 历史复盘、相似案例、适用场景、决策引用、使用反馈和知识版本 | [生产安全整改决策知识库项目案例](/projects/production-safety-rectification-decision-knowledge-base-case) |
| 生产安全整改决策智能推荐 | 整改上下文、相似案例、资源约束、风险规则、候选方案和效果回写 | [生产安全整改决策智能推荐项目案例](/projects/production-safety-rectification-decision-intelligent-recommendation-case) |
| 生产安全整改决策推荐评测 | 评测集、专家标注、离线评测、上线门禁、在线反馈和失败样本分析 | [生产安全整改决策推荐评测项目案例](/projects/production-safety-rectification-recommendation-evaluation-case) |
| 生产良率分析 | 良率、一次合格率、返工报废、缺陷原因、下钻分析和改善闭环 | [生产良率分析项目案例](/projects/production-yield-analysis-case) |
| 生产瓶颈分析 | 节拍、负荷、等待、在制品、瓶颈原因、改善任务和产能复盘 | [生产瓶颈分析项目案例](/projects/production-bottleneck-analysis-case) |
| 生产换型损失分析 | 换型计划、清线换模、调机、首件确认、爬坡损失和改善任务 | [生产换型损失分析项目案例](/projects/production-changeover-loss-analysis-case) |
| 生产设备异常 | 设备停机、维修派工、备件领用、恢复确认、MTTR/MTBF 和复盘 | [生产设备异常项目案例](/projects/production-equipment-exception-case) |
| 生产能耗分析 | 能源计量、单耗分析、峰谷成本、异常告警、成本分摊和节能验证 | [生产能耗分析项目案例](/projects/production-energy-analysis-case) |
| 生产成本核算 | 标准成本、实际成本、材料人工能耗、制造费用分摊、成本差异和结账 | [生产成本核算项目案例](/projects/production-cost-accounting-case) |
| 制造成本差异分析 | 标准实际对比、价格用量效率差异、超阈值任务、整改和调标建议 | [制造成本差异分析项目案例](/projects/manufacturing-cost-variance-case) |
| IoT 设备管理 | 设备接入、遥测、命令、告警、固件升级和运维记录 | [IoT 设备管理项目案例](/projects/iot-device-management-case) |
| 设备维保 | 设备台账、点检保养、维修工单、备件领用和停机分析 | [设备维保项目案例](/projects/equipment-maintenance-case) |
| 教育培训平台 | 课程、班级、报名、学习进度、作业、考试和证书 | [教育培训平台项目案例](/projects/education-training-platform-case) |

## 问题 1：本地正常，线上刷新 404

### 问题现象

- 本地 `npm run dev` 正常。
- 线上从菜单进入页面正常。
- 刷新 `/users`、`/roles`、`/dashboard` 后 404。

### 影响范围

所有使用 Vue Router history 模式的页面。

### 根因分析

前端路由只存在于浏览器中。刷新页面时，浏览器请求真实服务器路径。服务器没有对应文件，也没有回退到 `index.html`。

### 解决方案

Nginx 增加：

```nginx
location / {
  try_files $uri $uri/ /index.html;
}
```

### 预防方式

上线检查必须包含：

- 直接访问二级路由。
- 刷新二级路由。
- 复制详情页链接到新标签页。

## 问题 2：搜索列表偶尔显示旧结果

### 问题现象

快速输入搜索关键字，最后显示的结果不是最后一次搜索的结果。

### 根因分析

多个请求并发，旧请求后返回，覆盖了新请求结果。

### 解决方案

使用请求序号：

```ts
let requestId = 0

async function fetchList() {
  const currentId = ++requestId
  const result = await getList(query.value)

  if (currentId !== requestId) return

  list.value = result.items
}
```

### 预防方式

所有“会被快速连续触发”的请求都要考虑：

- 防抖。
- 取消请求。
- 请求序号。

## 问题 3：编辑表单还没保存，列表数据已经变了

### 问题现象

打开编辑弹窗，修改表单字段，还没点保存，背后的表格行已经变化。

### 根因分析

弹窗表单直接绑定了表格行对象。对象是引用类型，修改表单等于修改列表数据。

### 解决方案

打开弹窗时复制一份：

```ts
function openEdit(row: User) {
  form.value = {
    id: row.id,
    username: row.username,
    mobile: row.mobile,
    enabled: row.enabled
  }
}
```

保存成功后重新请求列表。

### 预防方式

表单编辑永远使用独立表单对象，不直接修改 props 或列表行对象。

## 问题 4：刷新后菜单丢失

### 问题现象

登录后菜单正常。刷新页面后菜单为空，或者进入页面后变成 404。

### 根因分析

菜单和动态路由保存在内存里。刷新后 Pinia 重置，动态路由没有重新注册。

### 解决方案

路由守卫中恢复：

```ts
if (userStore.token && !permissionStore.ready) {
  await initUserContext()
  return to.fullPath
}
```

### 预防方式

把“恢复用户上下文”作为应用启动流程的一部分，而不是只在登录成功后执行。

## 问题 5：401 后重复弹出多个登录失效提示

### 问题现象

登录过期后，页面同时弹出多个错误提示，甚至反复跳登录页。

### 根因分析

多个接口同时返回 401，每个响应拦截器都执行一次退出登录。

### 解决方案

加全局处理锁：

```ts
let isHandlingUnauthorized = false

function handleUnauthorized() {
  if (isHandlingUnauthorized) return

  isHandlingUnauthorized = true
  userStore.logout()
  router.replace('/login').finally(() => {
    isHandlingUnauthorized = false
  })
}
```

### 预防方式

所有全局错误处理都要考虑并发场景。

## 问题 6：组件库样式突然变形

### 问题现象

- 表格行高异常。
- Switch 变形。
- 按钮尺寸变小。
- 弹窗内容被压缩。

### 根因分析

常见原因是业务 CSS 使用了宽泛选择器，例如：

```css
.page button {}
.content div {}
.panel * {}
```

这些选择器污染了组件库内部 DOM。

### 解决方案

改成明确业务 class：

```css
.user-search-form__actions {}
.permission-switch-row {}
.metric-card__value {}
```

如需影响组件库样式，优先使用组件库主题 token、props、CSS 变量或官方 API。

### 预防方式

提交前搜索：

```bash
rg "(\\.\\w+\\s+(div|span|button|\\*)|div > div|\\.\\w+ \\*)" src
```

## 问题 7：删除最后一条数据后页面空白

### 问题现象

列表第 3 页只有 1 条数据，删除后页面显示空，但其实第 2 页还有数据。

### 根因分析

删除后仍请求当前页，而当前页已经没有数据。

### 解决方案

```ts
async function removeRecord(id: number) {
  await api.remove(id)
  await fetchList()

  if (list.value.length === 0 && pagination.page > 1) {
    pagination.page -= 1
    await fetchList()
  }
}
```

### 预防方式

分页删除、批量删除后都要考虑当前页是否仍有数据。

## 问题 8：构建成功但页面白屏

### 问题现象

`npm run build` 成功，但部署后页面白屏。

### 根因分析

常见原因：

- 静态资源路径错误。
- 部署子路径和 `base` 不一致。
- `index.html` 被强缓存。
- 环境变量缺失。

### 解决方案

排查：

1. 看 Network 中 js/css 是否 404。
2. 看 Console 第一条错误。
3. 检查 `base`。
4. 检查部署环境变量。
5. 清理 CDN 或浏览器缓存。

### 预防方式

上线前用真实部署路径预览，不只在本地根路径预览。

## 问题 9：表格操作列在小屏幕被压扁

### 问题现象

操作按钮换行、图标变形、头像变椭圆。

### 根因分析

固定尺寸元素缺少稳定宽高和 `flex-shrink: 0`。

### 解决方案

```css
.table-action {
  display: inline-flex;
  gap: 8px;
  flex-shrink: 0;
}

.user-avatar {
  width: 32px;
  height: 32px;
  flex: 0 0 32px;
  border-radius: 50%;
}
```

### 预防方式

每次修改布局后，至少检查桌面和较窄视口。

## 问题 10：权限码改名后多个页面坏掉

### 问题现象

后端把权限码改名，多个按钮、菜单和路由同时失效。

### 根因分析

权限码直接散落在页面模板中。

### 解决方案

集中维护权限码：

```ts
export const PermissionCode = {
  UserCreate: 'system:user:create',
  UserUpdate: 'system:user:update',
  UserDelete: 'system:user:delete'
} as const
```

使用：

```vue
<PermissionButton :code="PermissionCode.UserCreate">
  新增用户
</PermissionButton>
```

### 预防方式

权限码上线后保持稳定。必须改名时，前后端一起做兼容迁移。

## 持续沉淀规则

每次遇到线上或联调问题，都应该补充：

- 问题现象截图或描述。
- 最小复现步骤。
- 根因。
- 最终修复方案。
- 如何避免再次发生。
