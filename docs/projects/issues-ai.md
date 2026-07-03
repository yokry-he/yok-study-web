# AI 工程问题

## 适合谁看

这篇适合正在把 LLM、RAG、Agent 或 AI 辅助功能接入真实项目的同学。你不需要先成为算法工程师，但需要理解 AI 功能和普通接口不同：输出不完全稳定、成本受输入输出影响、质量需要评测，安全边界也更复杂。

## 使用方式

AI 工程问题排查时，至少记录：

- 用户输入。
- 最终 prompt。
- 模型名称和参数。
- 检索命中的文档。
- 模型输出。
- 延迟、token、成本。
- 用户是否采纳结果。

没有这些记录，很难判断问题来自模型、提示词、检索、上下文、数据质量还是产品流程。

## 问题 1：同一个问题，模型每次回答都不一样

### 问题现象

- 用户重复问同一个问题，答案有明显差异。
- 有时格式正确，有时格式错乱。
- 测试同学很难写稳定用例。

### 影响范围

问答助手、文案生成、代码生成、摘要、分类、信息抽取。

### 常见根因

LLM 默认具有生成随机性。如果 prompt 不够明确，模型会根据不同采样路径给出不同表达。即使温度较低，也不能把模型当成传统确定性函数。

### 解决方案

降低随机性。

```json
{
  "temperature": 0.2,
  "top_p": 0.9
}
```

把输出格式写清楚。

```text
请只输出 JSON，不要输出解释文字。
字段：
- title: 字符串，不超过 20 个字
- summary: 字符串，不超过 80 个字
- tags: 字符串数组，最多 3 个
```

对结构化任务增加校验和重试。

```ts
const result = await callModel(prompt)
const parsed = safeParseJson(result)

if (!parsed.ok) {
  return retryWithFormatRepair(result)
}
```

把“可接受答案范围”写成评测用例，不要只比较字符串是否完全相等。

### 预防方式

- 创意任务允许变化，结构化任务必须校验。
- 不把模型输出直接写入关键业务状态。
- 对核心 AI 功能建立样例集和评测脚本。
- 产品上明确哪些内容是“建议”，哪些内容需要人工确认。

## 问题 2：RAG 问答经常答非所问

### 问题现象

- 用户问公司制度，模型回答通用知识。
- 明明知识库里有答案，却没有引用到。
- 回答看起来流畅，但和内部文档不一致。

### 影响范围

知识库问答、客服助手、文档助手、企业内部助手。

### 常见根因

问题不一定在模型，而可能在检索链路：

- 文档切片太大或太小。
- 向量检索召回了不相关内容。
- 没有保留标题、章节、时间等元数据。
- prompt 没要求模型只基于检索内容回答。
- 检索结果没有进入最终上下文。

### 解决方案

检索结果要可观测。

```json
{
  "query": "年假怎么申请",
  "hits": [
    {
      "title": "员工休假制度",
      "score": 0.82,
      "chunk": "年假申请需要提前三个工作日..."
    }
  ]
}
```

prompt 明确回答边界。

```text
你只能根据下面的资料回答。
如果资料中没有答案，请回答“当前资料中没有找到明确答案”，不要编造。
```

切片要保留上下文信息。

```text
文档标题：员工休假制度
章节：年假申请
更新时间：2026-06-01
正文：年假申请需要提前三个工作日提交...
```

对无答案场景做单独处理。

```ts
if (retrievalHits.length === 0 || retrievalHits[0].score < 0.55) {
  return '当前资料中没有找到明确答案。'
}
```

### 预防方式

- 记录每次回答使用了哪些检索片段。
- 文档入库时保留标题、章节、时间、权限等元数据。
- 建立“应该答出”和“应该拒答”的测试问题。
- 定期抽样检查低评分命中和用户差评问题。

## 问题 3：AI 功能上线后成本突然变高

### 问题现象

- 用户量变化不大，但 token 成本快速上涨。
- 响应变慢。
- 某些页面重复触发 AI 请求。

### 影响范围

对话助手、智能总结、批量生成、Agent 工作流、RAG 问答。

### 常见根因

- 每次请求都带完整历史上下文。
- 检索结果过多，prompt 过长。
- 前端重复触发请求。
- Agent 循环没有限制步数。
- 没有缓存相同输入的结果。

### 解决方案

记录每次调用 token。

```json
{
  "feature": "doc_qa",
  "model": "example-model",
  "inputTokens": 4200,
  "outputTokens": 800,
  "latencyMs": 3200
}
```

限制上下文长度。

```ts
const recentMessages = messages.slice(-8)
```

限制 RAG 召回数量。

```ts
const hits = await retrieve(query, {
  topK: 5,
  minScore: 0.55
})
```

限制 Agent 最大步数。

```ts
const maxSteps = 6
for (let step = 0; step < maxSteps; step++) {
  await runAgentStep()
}
```

对相同输入和稳定知识库结果做缓存。

### 预防方式

- AI 调用必须记录 feature、model、tokens、latency、userId。
- 每个 AI 功能上线前估算单次成本和日成本上限。
- 前端按钮提交中禁用，避免重复调用。
- Agent 必须有最大步数、最大耗时和失败出口。

## 问题 4：模型把用户不该看的内部信息回答出来

### 问题现象

- 普通用户问到了管理员文档内容。
- AI 回答了其他部门的资料。
- RAG 检索命中了没有权限的文档。

### 影响范围

企业知识库、客服系统、内部文档助手、代码助手、数据分析助手。

### 常见根因

只在页面入口做了权限控制，没有在检索和生成链路里做权限过滤。AI 功能一旦能访问全部资料，就可能把不该看的内容组织成自然语言输出。

### 解决方案

检索前带上用户权限上下文。

```ts
const hits = await retrieve(query, {
  userId: currentUser.id,
  departments: currentUser.departments,
  roles: currentUser.roles
})
```

向量库或检索服务按元数据过滤。

```json
{
  "filter": {
    "department": { "$in": ["frontend", "platform"] },
    "visibility": "internal"
  }
}
```

最终回答也要检查引用来源。

```ts
const visibleHits = hits.filter((hit) => canRead(currentUser, hit.metadata))

if (visibleHits.length === 0) {
  return '当前没有可访问的资料。'
}
```

### 预防方式

- 权限必须进入检索层，不只进入页面层。
- 文档入库时写入 owner、department、visibility、expiresAt 等元数据。
- 日志中记录回答引用了哪些文档。
- 对越权问题建立专门测试集。

## 问题 5：答案带了来源，但来源和结论对不上

### 问题现象

- 回答里有“来源 1”“来源 2”。
- 点开来源后发现原文没有支持这个结论。
- 模型把多个片段混在一起，得出文档里没有的结论。

### 影响范围

文档问答、政策问答、客服知识库、代码库问答、企业内部助手。

### 常见根因

来源引用只是“格式上存在”，但没有被系统校验：

- prompt 要求引用，但模型随便引用。
- 后端没有检查引用 ID 是否来自本次检索结果。
- chunk 太大，引用定位不精确。
- 多个片段主题相近，模型把 A 文档结论归到 B 文档。

### 解决方案

回答结构中强制引用真实 `chunkId`。

```json
{
  "answer": "年假申请需要提前三个工作日提交。",
  "citations": [
    {
      "chunkId": "leave-policy:annual:003",
      "reason": "该片段说明年假申请提前时间"
    }
  ]
}
```

后端校验引用：

```ts
const hitIds = new Set(retrievalHits.map((hit) => hit.chunkId))

for (const citation of answer.citations) {
  if (!hitIds.has(citation.chunkId)) {
    throw new Error('模型引用了不存在的来源')
  }
}
```

如果要求更高，可以让模型先列出证据，再基于证据回答。

```text
请先列出支持结论的资料编号，再回答。
如果没有资料支持，不要回答该结论。
```

### 预防方式

- 引用必须绑定本次检索返回的 chunkId。
- chunk 保留标题、章节、段落位置。
- 评测集中加入“相似但结论不同”的文档。
- 对低置信度或来源不足的问题返回拒答。

## 问题 6：Agent 一直调用工具，成本和耗时失控

### 问题现象

- 一个用户问题触发十几次工具调用。
- 日志显示 Agent 在重复查询同一个接口。
- 页面等待很久才返回。
- 成本快速升高，但结果质量没有提升。

### 影响范围

多工具助手、数据分析 Agent、运维 Agent、文档检索 Agent、自动工单处理。

### 常见根因

Agent 没有明确停止条件和预算边界：

- 没有最大步数。
- 工具调用失败后反复重试。
- 工具结果没有被结构化总结。
- 模型不知道哪些信息已经获得。
- 高风险工具没有人工确认。

### 解决方案

给 Agent 设置硬边界。

```ts
const limits = {
  maxSteps: 6,
  maxToolCalls: 8,
  maxLatencyMs: 15000,
  maxEstimatedCost: 0.5
}
```

每一步记录状态。

```json
{
  "step": 3,
  "tool": "searchDocuments",
  "input": { "query": "权限缓存不生效" },
  "success": true,
  "resultSummary": "找到 3 条相关文档"
}
```

工具失败要分类：

| 失败类型 | 处理方式 |
| --- | --- |
| 参数错误 | 让模型修正一次 |
| 权限不足 | 立即停止并说明 |
| 外部服务超时 | 最多重试一次 |
| 高风险操作 | 等待人工确认 |

### 预防方式

- Agent 必须有最大步数、最大耗时、最大成本。
- 每个工具定义清楚输入、输出、失败类型。
- 高风险操作必须人工确认。
- 对重复工具调用做检测。
- 上线前用评测集覆盖工具失败、权限不足、资料不足场景。

## 下一步学习

- [AI 工程学习导览](/ai-engineering/introduction)
- [RAG 检索增强生成](/ai-engineering/rag)
- [AI 文档问答从零到项目](/ai-engineering/doc-qa-project)
- [评测与质量保障](/ai-engineering/evaluation)
- [上线、成本与安全](/ai-engineering/deployment)
