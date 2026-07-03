# AI 工程路线

## 适合谁看

适合已经有编程和项目基础，想把大语言模型接入真实业务系统的人。AI 工程路线不适合只停留在“写几个 prompt”，它强调 API、RAG、Agent、评测、成本、安全和上线治理。

<LearningPath :steps="[
  { title: '编程和后端基础', description: '掌握 JavaScript/TypeScript、HTTP、Node.js 和基本服务分层。', link: '/roadmap/node-backend', badge: '基础' },
  { title: 'LLM API 调用', description: '封装模型调用、处理结构化输出、错误、重试、日志和成本。', link: '/ai-engineering/llm-api', badge: 'API' },
  { title: '提示词工程', description: '把任务、输入、输出、边界、示例和失败策略写成可维护 prompt。', link: '/ai-engineering/prompt-engineering', badge: 'Prompt' },
  { title: 'RAG 知识库', description: '完成文档切分、向量检索、引用来源、权限过滤和拒答策略。', link: '/ai-engineering/rag', badge: '知识库' },
  { title: 'Agent 工作流', description: '理解工具调用、状态、审批、最大步骤、trace 和风险边界。', link: '/ai-engineering/agents', badge: 'Agent' },
  { title: '评测与上线', description: '建立评测集、回归测试、成本监控、延迟优化和安全治理。', link: '/ai-engineering/evaluation', badge: '质量' }
]" />

## 推荐项目

推荐从低风险项目开始：

- 文档问答助手。
- 工单分类助手。
- 后台运营文案助手。
- 代码规范问答助手。
- 项目问题库检索助手。

不建议初学阶段直接做：

- 自动审批。
- 自动退款。
- 自动删除数据。
- 自动执行生产命令。

这些场景需要更强的权限、审计、审批和回滚设计。

## 阶段验收

| 阶段 | 能力结果 |
| --- | --- |
| API 调用 | 能把模型调用封装到后端服务，并处理错误和日志 |
| Prompt | 能版本化 prompt，并解释每段约束的作用 |
| RAG | 能返回基于来源的答案，并处理资料不足 |
| Agent | 能定义工具边界、审批和最大步骤 |
| 评测 | 能用样例集判断 prompt 或模型改动是否退化 |
| 上线 | 能控制成本、延迟、安全和用户反馈闭环 |

## 学习顺序建议

1. 先做一个普通模型调用接口。
2. 再做结构化输出和校验。
3. 再接入文档知识库。
4. 再加入工具调用。
5. 再评估是否需要 Agent。
6. 最后补评测、成本和安全治理。

## 常见误区

### Demo 效果好就直接上线

Demo 通常只覆盖理想输入。上线前必须有评测集、失败兜底、权限过滤和日志。

### 把权限交给模型判断

权限必须由业务系统判断。模型可以生成建议，但不能决定用户能不能访问数据或执行操作。

### 只改 prompt 不看检索

RAG 问题经常发生在检索层。先看检索片段，再看生成答案。

## 下一步学习

继续进入 [AI 工程学习导览](/ai-engineering/introduction)，按 API、Prompt、RAG、Agent、评测和上线顺序推进。
