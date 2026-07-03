# MCP 与企业工具集成

## 适合谁看

适合想让 AI 应用连接企业内部工具、数据源和操作系统的人：

- 想让模型查询 GitHub、Notion、数据库或工单。
- 想把内部 API 暴露给 AI 助手。
- 不知道 MCP server、client、tools、resources 是什么。
- 担心模型越权访问内部资料。
- 想避免每个 AI 应用都重复写一套工具接入。

MCP 的价值是把“AI 应用如何连接上下文和工具”标准化。

## 基本角色

| 角色 | 说明 |
| --- | --- |
| Host | 用户使用的 AI 应用，例如编辑器或聊天应用 |
| Client | Host 内部连接 MCP Server 的组件 |
| Server | 暴露工具、资源和提示模板 |
| Tool | 可被调用的动作，例如查 issue、查订单 |
| Resource | 可读取的上下文，例如文档、文件、记录 |
| Prompt | 可复用的提示模板或工作流 |

MCP 官方规范把 server 能提供的能力分为 resources、prompts、tools。

## 什么时候需要 MCP

适合：

- 多个 AI 应用需要复用同一套工具。
- 企业内部系统很多。
- 工具需要统一鉴权和审计。
- 希望工具定义和业务系统解耦。
- 需要把文件、知识库、工单、代码仓库接入 AI。

不一定需要：

- 只有一个简单工具。
- 只是一次性脚本。
- 没有稳定权限模型。
- 业务 API 还没成型。

## MCP 和函数调用的关系

函数调用是模型调用业务能力的一种方式。

MCP 更像工具和上下文的标准接入协议。

```text
模型
↓
工具调用
↓
MCP client
↓
MCP server
↓
企业系统
```

实际项目中可以把 MCP server 暴露的 tool 作为模型可调用工具。

## Tool 设计

工具要小而明确。

好例子：

- `search_docs`
- `get_ticket`
- `create_draft_reply`
- `list_user_permissions`

差例子：

- `do_everything`
- `query`
- `run_sql_anything`

工具描述要写清：

- 什么时候使用。
- 参数是什么。
- 返回什么。
- 需要什么权限。
- 是否有副作用。

## Resource 设计

Resource 适合提供上下文：

- 当前项目文件。
- 文档页面。
- 数据记录。
- 工单详情。
- API 文档。

Resource 不应该绕过权限。用户看不到的数据，AI 也不应该通过 Resource 看到。

## 权限和审计

企业 MCP 必须考虑：

- 用户身份如何传递。
- 工具是否按用户权限执行。
- 是否支持租户隔离。
- 工具调用是否记录日志。
- 敏感数据是否脱敏。
- 高风险工具是否需要审批。

不要让 MCP server 使用一个超级管理员账号替所有用户执行操作。

## 实际项目问题

### 1. MCP 工具太多，模型乱用

**原因**

一次暴露所有工具，工具名称和描述不清楚。

**解决方案**

- 按场景暴露工具。
- 工具名称明确。
- 工具描述写清适用范围。
- 对工具选择做评测。

### 2. 工具返回太多数据

**原因**

Resource 或 tool 没有分页、筛选和摘要。

**解决方案**

- 默认限制返回数量。
- 支持分页。
- 返回结构化摘要。
- 大内容走检索或分段读取。

### 3. AI 看到用户无权访问的数据

**原因**

MCP server 用系统账号查询，没带用户上下文。

**解决方案**

- 每次调用传递用户身份。
- 工具内部做权限过滤。
- 返回前做脱敏。
- 写审计日志。

### 4. 工具调用失败但用户看不懂

**原因**

工具错误直接暴露内部异常。

**解决方案**

- 工具返回稳定错误码。
- 模型或前端转成用户可理解提示。
- 日志记录内部错误和 request id。

## 最佳实践

- MCP server 按业务域设计。
- Tool 小而明确，避免万能工具。
- Resource 必须遵守用户权限。
- 工具调用要有日志和 request id。
- 高风险写操作需要确认或审批。
- 工具返回结构化、分页、可解释。
- MCP 接入前先明确安全边界。

## 参考资料

- [Model Context Protocol specification](https://modelcontextprotocol.io/specification/2025-06-18)
- [MCP Tools](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [OpenAI Using Tools](https://developers.openai.com/api/docs/guides/tools)

## 下一步学习

继续学习 [AI 产品设计与人机协作](/ai-engineering/product-workflow)，把模型能力、工具边界和人工确认组合成真实产品流程。
