# 销售风险指标治理成本收益评估项目案例

## 适合谁看

- 想理解销售风险指标治理如何证明投入产出的前端开发者。
- 正在做销售风控、指标平台、数据治理、经营分析、CRM 或治理运营看板的团队。
- 希望避免“治理做了很多任务，但管理层看不到节省了多少成本、降低了多少风险”的项目负责人。

## 业务目标

销售风险指标治理运营看板能展示指标健康度和任务闭环，但治理团队还需要回答一个更现实的问题：治理投入是否值得。成本收益评估要把治理人力、系统建设、规则维护、异常处理和自动修复投入，与风险降低、人工节省、决策效率提升、审计成本下降和业务损失减少放在同一套口径下评估。

成本收益评估要解决：

- 指标治理投入如何量化。
- 治理收益如何从风险、效率、审计和业务损失维度计算。
- 哪些指标应该优先治理，哪些指标暂时维持低成本治理。
- 治理动作是否真正降低异常和人工处理成本。
- 如何向管理层解释治理价值和下一阶段预算。

## 成本收益评估链路

```mermaid
flowchart LR
  A[治理动作] --> B[投入归集]
  A --> C[收益归集]
  B --> D[成本模型]
  C --> E[收益模型]
  D --> F[ROI 计算]
  E --> F
  F --> G[优先级排序]
  G --> H[预算建议]
  H --> I[治理复盘]
```

治理 ROI 不是单纯财务公式。很多收益是避免损失、减少返工和降低审计风险，需要用业务可接受的估算口径表达。

## 核心概念

| 概念 | 说明 |
| --- | --- |
| 治理成本 | 指标梳理、口径修订、质量规则、血缘建设、自动修复和任务运营产生的投入。 |
| 治理收益 | 减少异常、减少人工核查、降低坏账风险、提升决策效率和降低审计成本带来的价值。 |
| 避免损失 | 因指标更可信而提前发现风险、避免错误决策或减少业务损失。 |
| 人工节省 | 异常定位、数据核查、复算和解释工作减少的工时。 |
| ROI 分层 | 按指标、业务线、治理动作和时间周期分别评估。 |
| 预算建议 | 根据收益、风险和成熟度差距给出下一阶段治理投入建议。 |

## 数据模型

```mermaid
erDiagram
  METRIC_ASSET ||--o{ GOVERNANCE_COST_RECORD : costs
  METRIC_ASSET ||--o{ GOVERNANCE_BENEFIT_RECORD : benefits
  GOVERNANCE_COST_RECORD ||--o{ COST_ALLOCATION_ITEM : allocated_by
  GOVERNANCE_BENEFIT_RECORD ||--o{ BENEFIT_EVIDENCE : proves
  METRIC_ASSET ||--o{ ROI_EVALUATION : evaluated_by
  ROI_EVALUATION ||--o{ ROI_DIMENSION_SCORE : contains
  ROI_EVALUATION ||--o{ GOVERNANCE_BUDGET_SUGGESTION : suggests

  METRIC_ASSET {
    string id
    string metric_code
    string business_domain
    string risk_level
  }
  GOVERNANCE_COST_RECORD {
    string id
    string metric_id
    string cost_type
    decimal cost_amount
  }
  GOVERNANCE_BENEFIT_RECORD {
    string id
    string metric_id
    string benefit_type
    decimal benefit_amount
  }
  ROI_EVALUATION {
    string id
    string metric_id
    string period
    decimal roi_value
  }
```

成本和收益要分开建模。成本通常来自任务、人力和系统，收益通常来自异常、风控、审计和业务结果。

## 推荐表结构

| 表 | 作用 | 关键字段 |
| --- | --- | --- |
| `governance_cost_record` | 保存治理成本 | `metric_id`、`cost_type`、`cost_amount`、`period` |
| `cost_allocation_item` | 保存成本分摊 | `cost_id`、`allocation_target`、`allocation_ratio`、`reason` |
| `governance_benefit_record` | 保存治理收益 | `metric_id`、`benefit_type`、`benefit_amount`、`confidence_level` |
| `benefit_evidence` | 保存收益证据 | `benefit_id`、`evidence_type`、`evidence_value`、`source_system` |
| `roi_evaluation` | 保存 ROI 评估 | `metric_id`、`period`、`roi_value`、`payback_period` |
| `roi_dimension_score` | 保存维度评分 | `evaluation_id`、`dimension`、`score`、`summary` |
| `governance_budget_suggestion` | 保存预算建议 | `evaluation_id`、`suggested_budget`、`priority`、`reason` |

## ROI 计算流程

```mermaid
sequenceDiagram
  participant Ops as 治理运营
  participant Task as 任务中心
  participant Quality as 质量中心
  participant Risk as 风控系统
  participant Finance as 财务口径
  participant ROI as ROI 服务

  Ops->>ROI: 发起周期评估
  ROI->>Task: 归集治理任务工时和投入
  ROI->>Quality: 统计异常下降和修复效率
  ROI->>Risk: 统计风险发现和避免损失
  ROI->>Finance: 获取成本和收益折算口径
  ROI->>ROI: 计算 ROI 和预算建议
  ROI-->>Ops: 输出评估报告
```

收益折算口径要提前和财务、业务负责人确认，否则 ROI 报告很容易被质疑。

## 评估状态设计

```mermaid
stateDiagram-v2
  [*] --> Draft
  Draft --> CollectingCost: 归集成本
  CollectingCost --> CollectingBenefit: 归集收益
  CollectingBenefit --> Calculating: 计算 ROI
  Calculating --> Reviewing: 业务复核
  Reviewing --> Published: 发布报告
  Reviewing --> Adjusting: 调整口径
  Adjusting --> Calculating: 重新计算
  Published --> BudgetPlanning: 形成预算建议
  BudgetPlanning --> Closed: 关闭
  Closed --> [*]
```

ROI 评估要允许复核和调整口径，但每次调整都要保留原因。

## 成本收益维度拆解

```mermaid
flowchart TD
  A[治理成本收益] --> B[成本维度]
  A --> C[收益维度]
  B --> D[人力成本]
  B --> E[系统成本]
  B --> F[运营成本]
  C --> G[风险降低]
  C --> H[效率提升]
  C --> I[审计节省]
  C --> J[业务损失减少]
```

成本收益不要只看总额。不同指标的收益来源不同，有的偏风险，有的偏效率，有的偏审计。

## 优先级矩阵

```mermaid
flowchart LR
  A[指标候选池] --> B{风险高低}
  B -->|高风险| C{治理收益}
  B -->|低风险| D{治理成本}
  C -->|高收益| E[优先投入]
  C -->|低收益| F[保底治理]
  D -->|低成本| G[标准治理]
  D -->|高成本| H[暂缓或自动化]
```

高风险不等于一定大投入。若治理收益低或成本极高，应先做保底治理和监控。

## 前端页面拆分

| 页面 | 核心内容 | 设计重点 |
| --- | --- | --- |
| ROI 总览 | 总投入、总收益、ROI、回收周期、预算建议 | 管理层优先看趋势和结论。 |
| 指标 ROI 列表 | 指标、风险等级、成本、收益、ROI、优先级 | 支持按业务线和负责人筛选。 |
| 评估详情 | 成本明细、收益证据、计算口径、复核记录 | 让 ROI 结果可解释。 |
| 预算建议 | 下一阶段投入、预期收益、风险说明、审批状态 | 把评估结果转为资源申请。 |
| 口径配置 | 成本类型、收益类型、折算规则、置信等级 | 防止不同周期口径不一致。 |

## 接口拆分建议

| 接口 | 作用 |
| --- | --- |
| `GET /api/sales-risk-metric-governance-roi-dashboard` | 查询 ROI 总览。 |
| `GET /api/sales-risk-metric-governance-roi-evaluations` | 查询 ROI 评估列表。 |
| `POST /api/sales-risk-metric-governance-roi-evaluations` | 创建评估。 |
| `GET /api/sales-risk-metric-governance-roi-evaluations/:id` | 查询评估详情。 |
| `POST /api/sales-risk-metric-governance-roi-evaluations/:id/calculate` | 计算 ROI。 |
| `POST /api/sales-risk-metric-governance-roi-evaluations/:id/review` | 提交业务复核。 |
| `POST /api/sales-risk-metric-governance-roi-evaluations/:id/budget-suggestions` | 创建预算建议。 |

## 实际项目常见问题

### 1. 只统计成本不统计收益

治理看起来永远是成本中心。解决方式是同步记录异常减少、人工节省、风险降低和审计节省。

### 2. 收益口径太主观

业务不认可节省金额。解决方式是为每类收益配置证据和置信等级。

### 3. ROI 只看短期

指标治理前期投入大，短期 ROI 低。解决方式是同时展示回收周期和长期收益。

### 4. 成本无法分摊到指标

平台级建设服务多个指标。解决方式是按使用量、风险等级或任务量分摊成本。

### 5. 评估报告不能转成预算

报告看完就结束。解决方式是评估详情直接生成预算建议和审批材料。

## 权限与审计

| 权限 | 说明 |
| --- | --- |
| 查看 ROI | 可以查看成本收益评估和趋势。 |
| 维护口径 | 可以配置成本和收益折算规则。 |
| 发起评估 | 可以创建周期性 ROI 评估。 |
| 复核收益 | 可以确认收益证据和金额。 |
| 提交预算建议 | 可以将评估结果转成预算申请。 |

成本归集、收益证据、口径调整、ROI 计算、复核意见和预算建议都要保留审计。

## 验收清单

- 能归集指标治理成本。
- 能记录风险降低、效率提升、审计节省和业务损失减少收益。
- 能按指标、业务线和周期计算 ROI。
- 能解释成本和收益来源。
- 能生成治理优先级排序。
- 能根据 ROI 输出预算建议。
- 能保留评估口径和复核记录。

## 下一步学习

- [销售风险指标治理运营看板项目案例](/projects/sales-risk-metric-governance-operations-dashboard-case)
- [销售风险指标治理成熟度项目案例](/projects/sales-risk-metric-governance-maturity-case)
- [数据资产运营项目案例](/projects/data-asset-operation-case)
