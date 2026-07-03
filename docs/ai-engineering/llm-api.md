# LLM API 调用

## 适合谁看

适合想把大语言模型接入 Node.js、后台系统、文档站、客服系统或内部工具的人。

这一节不追求覆盖所有参数，而是先建立工程调用模型：请求从哪里来、模型怎么选、输出如何校验、错误怎么处理、成本怎么记录。

## 基础调用形态

在项目里，模型调用通常不是直接散落在页面或路由里，而是封装成 service：

```ts
interface GenerateSummaryInput {
  title: string
  content: string
}

interface GenerateSummaryResult {
  summary: string
  risks: string[]
}

export async function generateSummary(input: GenerateSummaryInput): Promise<GenerateSummaryResult> {
  // 真实项目中这里调用 LLM API，并对返回结果做结构校验。
  return {
    summary: `${input.title} 的摘要`,
    risks: []
  }
}
```

这样做的好处：

- 页面不关心模型供应商细节。
- 后续可以替换模型或 API。
- 方便加日志、重试、限流和评测。
- 方便单元测试。

## 推荐分层

```text
route/controller
↓
ai service
↓
prompt builder
↓
provider client
↓
model API
```

| 层 | 职责 |
| --- | --- |
| route | 接收 HTTP 请求、鉴权、参数校验 |
| ai service | 组织业务流程、选择 prompt 和模型 |
| prompt builder | 拼接任务说明、上下文和输出格式 |
| provider client | 调用模型 API、处理错误和重试 |
| model API | 生成结果 |

不要把 prompt、业务判断、API 调用和数据库写入全部塞在一个路由函数里。

## 输入设计

用户输入必须先校验：

```ts
interface AskDocumentInput {
  question: string
  documentIds: string[]
}

function validateAskDocumentInput(input: AskDocumentInput) {
  if (!input.question.trim()) {
    throw new Error('问题不能为空')
  }

  if (input.documentIds.length === 0) {
    throw new Error('至少选择一份文档')
  }
}
```

不要把完整用户请求无控制地塞给模型。应限制长度、过滤无关字段，并保留必要上下文。

## 输出设计

如果输出要被程序继续处理，尽量要求结构化结果，并做运行时校验。

```ts
interface TicketClassification {
  category: 'bug' | 'feature' | 'question'
  priority: 'low' | 'medium' | 'high'
  reason: string
}
```

模型返回 JSON 不代表一定可信。仍然要解析、校验、兜底。

## 错误处理

常见错误：

- 网络超时。
- 认证失败。
- 速率限制。
- 上下文过长。
- 输出不是预期 JSON。
- 安全策略拒绝。

错误应该分层处理：

```text
用户可修复：输入太长、问题为空
系统可重试：网络波动、临时限流
系统不可重试：API key 错误、配置缺失
业务需降级：模型不可用、输出无法解析
```

## 日志和追踪

AI 调用至少记录：

- 请求 ID。
- 用户或租户 ID。
- 场景名称。
- 模型名称。
- 输入长度。
- 输出长度。
- 耗时。
- 是否命中缓存。
- 是否失败。

不要直接记录敏感原文。可以记录摘要、hash 或脱敏后的片段。

## 实际项目问题

### 问题：前端直接调用模型 API

**风险**

- API key 暴露。
- 用户可以绕过权限。
- 无法统一限流和审计。
- 成本不可控。

**解决方案**

前端调用自己的后端接口，由后端鉴权、校验、限流、记录日志后再调用模型。

### 问题：模型输出 JSON 偶尔解析失败

**原因**

模型输出天然有不确定性，prompt 约束不够或缺少结构化输出校验。

**解决方案**

- 使用结构化输出能力或明确 JSON schema。
- 对结果做运行时校验。
- 失败时进行一次受控重试。
- 多次失败后返回可理解错误，不要让页面白屏。

## 最佳实践

- 模型调用放在后端。
- Prompt 和模型参数集中管理。
- 输出进入业务流程前必须校验。
- 所有 AI 调用都要有超时、错误和日志。
- API key 和模型配置放在服务端环境变量中。
- 高成本场景要记录 token、耗时和调用来源。

## 参考资料

- [OpenAI TypeScript SDK](https://developers.openai.com/api/reference/typescript/)
- [OpenAI Responses API](https://developers.openai.com/api/reference/typescript/)

## 下一步学习

继续学习 [提示词工程](/ai-engineering/prompt-engineering)。
