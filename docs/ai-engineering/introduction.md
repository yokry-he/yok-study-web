# AI 工程学习导览

## 适合谁看

适合已经会做 Web、Node.js 或后台系统，但想把大语言模型接入真实项目的人。

AI 工程不是“会调用一次模型接口”就结束。真正上线时，你会遇到：

- 用户输入不可控。
- 模型输出不稳定。
- 私有知识库答不上来。
- 工具调用失败。
- 成本和延迟超预算。
- 线上效果无法评估。
- 提示词改动后旧能力退化。
- 数据安全和权限边界不清楚。

本模块的目标是把 AI 能力当成工程系统来设计，而不是只写几个 prompt。

## 你会学到什么

- 如何用 LLM API 完成文本生成、结构化输出和多轮对话。
- 如何写可维护、可调试、可评审的提示词。
- 如何用 RAG 把模型连接到私有知识库。
- 如何区分普通 API 调用、工具调用和 Agent 工作流。
- 如何设计评测集、回归测试和人工验收。
- 如何控制成本、延迟、权限、安全和上线风险。

## 学习路线

```text
LLM API 基础调用
↓
图解 AI 工程核心概念
↓
提示词工程
↓
结构化输出和工具调用
↓
RAG 检索增强生成
↓
Agent 工作流
↓
评测与质量保障
↓
AI 文档问答从零到项目
↓
上线、成本、安全和监控
↓
真实问题排查
```

## AI 工程常见能力分层

| 层级 | 解决什么问题 | 典型技术 |
| --- | --- | --- |
| 模型调用 | 让模型生成文本、JSON、摘要、分类 | Responses API、Chat Completions |
| 提示词 | 约束任务、输入、输出、边界 | system/developer/user prompt、few-shot |
| 结构化输出 | 让结果可被程序消费 | JSON schema、类型校验 |
| 工具调用 | 让模型触发业务能力 | function calling、tools |
| RAG | 接入私有知识 | embeddings、vector store、file search |
| Agent | 多步规划、调用工具、状态管理 | Agents SDK、工作流编排 |
| 评测 | 判断效果是否真的变好 | eval set、回归测试、人工标注 |
| 生产治理 | 控制成本、延迟、安全、监控 | rate limit、日志、权限、缓存 |

## 模块章节

| 章节 | 解决的问题 |
| --- | --- |
| [图解 AI 工程核心概念](/ai-engineering/visual-guide) | 用图理解模型调用、Prompt、结构化输出、工具、RAG、Agent、评测和上线治理 |
| [LLM API 调用](/ai-engineering/llm-api) | 如何从项目代码里稳定调用模型 |
| [提示词工程](/ai-engineering/prompt-engineering) | 如何让模型更稳定地完成任务 |
| [结构化输出与函数调用](/ai-engineering/structured-outputs-tools) | 如何让模型输出可验证 JSON，并安全调用业务工具 |
| [多模态 AI 应用](/ai-engineering/multimodal) | 如何处理图片、语音、文件和多模态成本权限 |
| [RAG 检索增强生成](/ai-engineering/rag) | 如何让模型回答私有知识库问题 |
| [MCP 与企业工具集成](/ai-engineering/mcp-integration) | 如何用 MCP 连接企业内部工具、资源和提示模板 |
| [Agent 工作流](/ai-engineering/agents) | 什么时候需要 Agent，如何控制边界 |
| [AI 产品设计与人机协作](/ai-engineering/product-workflow) | 如何设计人工确认、反馈闭环和灰度上线 |
| [评测与质量保障](/ai-engineering/evaluation) | 如何证明提示词和模型改动没有退化 |
| [AI 文档问答从零到项目](/ai-engineering/doc-qa-project) | 用企业内部文档问答案例串联导入、切分、检索、权限、回答、引用、评测和上线治理 |
| [上线、成本与安全](/ai-engineering/deployment) | 如何把 AI 功能稳定上线 |
| [常见问题](/ai-engineering/troubleshooting) | 真实项目里的高频故障和处理方式 |

## 先做什么项目

建议从低风险场景开始：

- 文档问答。
- 工单分类。
- 文本摘要。
- SQL 或代码解释辅助。
- 表单内容润色。
- 后台运营助手。

不要一开始就做自动下单、自动审批、自动删除数据这类高风险闭环。AI 功能上线初期应保留人工确认。

## 核心原则

- 模型输出不等于事实，关键结论要有来源或校验。
- 模型建议不等于权限，业务权限仍由系统判断。
- Prompt 是代码的一部分，需要版本管理和评审。
- AI 功能要可观测，至少记录输入摘要、输出、耗时、成本、错误。
- 任何自动执行动作都要有边界、审批和回滚策略。

## 参考资料

- [OpenAI: Prompt engineering](https://developers.openai.com/api/docs/guides/prompt-engineering)
- [OpenAI: TypeScript SDK and Responses API](https://developers.openai.com/api/reference/typescript/)
- [OpenAI: Agents SDK](https://developers.openai.com/api/docs/guides/agents)

## 下一步学习

继续学习 [图解 AI 工程核心概念](/ai-engineering/visual-guide)，再进入 [LLM API 调用](/ai-engineering/llm-api)。
