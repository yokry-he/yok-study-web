# Go Task API 接口契约

本文档描述 `go-task-api` 的完整 HTTP 契约。示例默认服务地址为 `http://127.0.0.1:8080`，命令可直接在 Bash、Zsh 或兼容终端中执行。

> 该服务没有认证和授权，只适合本地学习。不要直接暴露到公网或放入生产环境。

## 1. 快速准备

```bash
export BASE_URL=http://127.0.0.1:8080
```

所有接口均使用 UTF-8。`POST`、`PUT`、`PATCH` 请求必须带：

```http
Content-Type: application/json
```

请求体上限为 1 MiB。JSON 使用严格解码：未知字段、空请求体、字段类型错误、两个连续 JSON 值都会被拒绝。

## 2. 统一响应格式

### 2.1 成功 envelope

除 `DELETE 204` 外，成功响应都使用同一结构：

```json
{
  "success": true,
  "data": {},
  "requestId": "req_8a6287f8c2332a4830da410f03173f10"
}
```

### 2.2 失败 envelope

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "任务参数校验失败",
    "fields": {
      "title": "标题长度必须为 1 到 128 个字符"
    }
  },
  "requestId": "req_8a6287f8c2332a4830da410f03173f10"
}
```

`fields` 只在存在字段级错误时返回。数据库错误、panic、连接串和内部堆栈不会出现在响应中。

### 2.3 Request ID

- 客户端可以发送 `X-Request-ID`，允许字符为字母、数字、`.`、`_`、`:`、`-`，长度 1 到 128。
- 缺失或格式不合法时，服务端生成 `req_` 开头的 ID。
- 响应 Header 和 JSON 中的 `requestId` 必须一致。
- 排障时应保存该值，用它关联访问日志。

```bash
curl -i "$BASE_URL/health/live" \
  -H 'X-Request-ID: local-check-001'
```

### 2.4 时间、排序与分页

- 时间字段使用 RFC 3339，例如 `2026-08-01T10:00:00Z`。
- 所有列表按 `createdAt DESC, id DESC` 稳定排序。
- `page` 默认为 `1`，`pageSize` 默认为 `20`，最大 `100`。
- `total` 是过滤后总数，不是当前页条数。
- 超出最后一页时 `items` 为空，但 `total` 保持真实总数。

### 2.5 乐观锁

用户状态修改、任务更新、任务状态修改和任务删除都要求客户端回传最后读取到的 `version`：

1. 先读取资源，记录 `version`。
2. 写请求传入 `expectedVersion`。
3. 写入成功后，服务端把 `version` 加一。
4. 若其他请求先完成写入，旧版本请求返回 `409`。
5. 收到冲突后重新读取资源，由用户或业务规则决定是否重试；不要盲目覆盖。

## 3. 路由总览

| # | 方法 | 路径 | 成功状态 | 作用 |
|---|---|---|---:|---|
| 1 | `GET` | `/health/live` | 200 | 进程存活探针 |
| 2 | `GET` | `/health/ready` | 200 | 数据库与关闭状态就绪探针 |
| 3 | `GET` | `/api/users` | 200 | 用户列表 |
| 4 | `POST` | `/api/users` | 201 | 创建用户 |
| 5 | `GET` | `/api/users/{id}` | 200 | 用户详情 |
| 6 | `PATCH` | `/api/users/{id}/status` | 200 | 修改用户状态 |
| 7 | `GET` | `/api/tasks` | 200 | 任务列表 |
| 8 | `POST` | `/api/tasks` | 201 | 创建任务 |
| 9 | `GET` | `/api/tasks/{id}` | 200 | 任务详情 |
| 10 | `PUT` | `/api/tasks/{id}` | 200 | 完整替换任务可编辑字段 |
| 11 | `PATCH` | `/api/tasks/{id}/status` | 200 | 推进任务状态 |
| 12 | `DELETE` | `/api/tasks/{id}` | 204 | 删除任务 |

对已存在路径使用不受支持的方法时返回 `405`，并通过 `Allow` Header 给出可用方法。未知路径返回 `404 NOT_FOUND`。

## 4. 健康检查

### 4.1 进程存活

```http
GET /health/live
```

该接口不依赖数据库。只要 HTTP 进程仍能处理请求，就返回：

```json
{
  "success": true,
  "data": { "status": "alive" },
  "requestId": "local-check-001"
}
```

```bash
curl --fail-with-body "$BASE_URL/health/live"
```

### 4.2 服务就绪

```http
GET /health/ready
```

服务端会实际执行数据库 `PingContext`。数据库可用且应用没有进入关闭阶段时返回 `200`：

```json
{
  "success": true,
  "data": { "status": "ready" },
  "requestId": "req_..."
}
```

数据库不可用或应用正在关闭时返回 `503`：

```json
{
  "success": false,
  "error": {
    "code": "NOT_READY",
    "message": "服务尚未就绪"
  },
  "requestId": "req_..."
}
```

```bash
curl --fail-with-body "$BASE_URL/health/ready"
```

## 5. 用户接口

### 5.1 用户对象

```json
{
  "id": 1,
  "name": "张三",
  "email": "user@example.com",
  "status": "ACTIVE",
  "version": 0,
  "createdAt": "2026-07-21T08:00:00Z",
  "updatedAt": "2026-07-21T08:00:00Z"
}
```

字段约束：

| 字段 | 约束 |
|---|---|
| `name` | 去除首尾空白后 2 到 64 个 Unicode 字符；不能包含控制字符 |
| `email` | 去除首尾空白并转小写；1 到 254 字节；必须是单个规范邮箱地址 |
| `status` | `ACTIVE` 或 `DISABLED` |
| `version` | 从 0 开始；每次成功修改状态加一 |

邮箱唯一性不区分大小写，因此 `USER@example.com` 与 `user@example.com` 冲突。

### 5.2 查询用户列表

```http
GET /api/users?page=1&pageSize=20&status=ACTIVE
```

查询参数：

| 参数 | 必填 | 默认值 | 说明 |
|---|---|---|---|
| `page` | 否 | 1 | 正整数页码 |
| `pageSize` | 否 | 20 | 1 到 100 |
| `status` | 否 | 无 | `ACTIVE` 或 `DISABLED` |

成功响应中的 `data`：

```json
{
  "items": [
    {
      "id": 1,
      "name": "张三",
      "email": "user@example.com",
      "status": "ACTIVE",
      "version": 0,
      "createdAt": "2026-07-21T08:00:00Z",
      "updatedAt": "2026-07-21T08:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "total": 1
}
```

错误：非整数查询参数返回 `400 INVALID_ARGUMENT`；越界分页或未知状态返回 `422 VALIDATION_FAILED`。

```bash
curl --fail-with-body \
  "$BASE_URL/api/users?page=1&pageSize=20&status=ACTIVE"
```

### 5.3 创建用户

```http
POST /api/users
Content-Type: application/json
```

请求：

```json
{
  "name": "张三",
  "email": "user@example.com"
}
```

成功返回 `201`、完整用户对象，并设置：

```http
Location: /api/users/1
```

常见错误：

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 400 | `INVALID_JSON` | JSON 语法错误、空 body 或多个 JSON 值 |
| 400 | `UNKNOWN_FIELD` | 请求包含未声明字段 |
| 400 | `BODY_TOO_LARGE` | 请求体超过 1 MiB |
| 415 | `UNSUPPORTED_MEDIA_TYPE` | Content-Type 不是 JSON |
| 422 | `VALIDATION_FAILED` | 姓名或邮箱不符合约束 |
| 409 | `EMAIL_CONFLICT` | 邮箱已经存在 |

```bash
curl --fail-with-body -i \
  -X POST "$BASE_URL/api/users" \
  -H 'Content-Type: application/json' \
  -d '{"name":"张三","email":"user@example.com"}'
```

### 5.4 查询用户详情

```http
GET /api/users/{id}
```

`id` 必须是大于 0 且不超过 `int64` 上限的十进制整数。成功返回 `200` 和完整用户对象。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 400 | `INVALID_ARGUMENT` | ID 不是正整数或发生溢出 |
| 404 | `USER_NOT_FOUND` | 用户不存在 |

```bash
curl --fail-with-body "$BASE_URL/api/users/1"
```

### 5.5 修改用户状态

```http
PATCH /api/users/{id}/status
Content-Type: application/json
```

请求：

```json
{
  "status": "DISABLED",
  "expectedVersion": 0
}
```

成功返回 `200` 和新用户对象；本例新 `version` 为 `1`。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 400 | `INVALID_ARGUMENT` | ID 或 JSON 字段类型不正确 |
| 404 | `USER_NOT_FOUND` | 用户不存在 |
| 409 | `USER_VERSION_CONFLICT` | `expectedVersion` 已过期 |
| 422 | `VALIDATION_FAILED` | 状态非法、版本为负数或缺少版本 |

```bash
curl --fail-with-body \
  -X PATCH "$BASE_URL/api/users/1/status" \
  -H 'Content-Type: application/json' \
  -d '{"status":"DISABLED","expectedVersion":0}'
```

禁用用户后不能再为其创建任务；已有任务不会被自动删除。

## 6. 任务接口

### 6.1 任务对象

```json
{
  "id": 8,
  "ownerId": 1,
  "title": "完成 Go 文档",
  "description": "补齐示例、测试与排障说明",
  "status": "TODO",
  "dueAt": "2026-08-01T10:00:00Z",
  "version": 0,
  "createdAt": "2026-07-21T08:30:00Z",
  "updatedAt": "2026-07-21T08:30:00Z"
}
```

字段约束：

| 字段 | 约束 |
|---|---|
| `ownerId` | 必须引用一个存在且状态为 `ACTIVE` 的用户 |
| `title` | 去除首尾空白后 1 到 128 个 Unicode 字符；不能包含控制字符 |
| `description` | 可为字符串或 `null`；字符串会去除首尾空白 |
| `dueAt` | 可为 RFC 3339 时间或 `null`；不能早于服务端当前时间 |
| `status` | `TODO`、`DOING`、`DONE`、`CANCELLED` |
| `version` | 从 0 开始；更新字段或状态后加一 |

状态转换规则：

| 当前状态 | 可以转换到 |
|---|---|
| `TODO` | `DOING`、`CANCELLED` |
| `DOING` | `TODO`、`DONE`、`CANCELLED` |
| `DONE` | 无 |
| `CANCELLED` | 无 |

### 6.2 查询任务列表

```http
GET /api/tasks?page=1&pageSize=20&ownerId=1&status=TODO
```

| 参数 | 必填 | 默认值 | 说明 |
|---|---|---|---|
| `page` | 否 | 1 | 正整数页码 |
| `pageSize` | 否 | 20 | 1 到 100 |
| `ownerId` | 否 | 无 | 正整数用户 ID |
| `status` | 否 | 无 | 四种任务状态之一 |

成功响应中的 `data`：

```json
{
  "items": [
    {
      "id": 8,
      "ownerId": 1,
      "title": "完成 Go 文档",
      "description": null,
      "status": "TODO",
      "dueAt": null,
      "version": 0,
      "createdAt": "2026-07-21T08:30:00Z",
      "updatedAt": "2026-07-21T08:30:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "total": 1
}
```

```bash
curl --fail-with-body \
  "$BASE_URL/api/tasks?ownerId=1&status=TODO&page=1&pageSize=20"
```

### 6.3 创建任务

```http
POST /api/tasks
Content-Type: application/json
```

请求：

```json
{
  "ownerId": 1,
  "title": "完成 Go 文档",
  "description": "补齐示例与排障说明",
  "dueAt": "2026-08-01T10:00:00Z"
}
```

`description` 和 `dueAt` 可以省略或传 `null`。成功返回 `201`、状态为 `TODO` 的完整任务对象，并设置：

```http
Location: /api/tasks/8
```

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 404 | `TASK_OWNER_NOT_FOUND` | 负责人不存在 |
| 422 | `TASK_OWNER_DISABLED` | 负责人已禁用 |
| 422 | `VALIDATION_FAILED` | 标题、截止时间或 ownerId 不合法 |

```bash
curl --fail-with-body -i \
  -X POST "$BASE_URL/api/tasks" \
  -H 'Content-Type: application/json' \
  -d '{"ownerId":1,"title":"完成 Go 文档","description":"补齐示例与排障说明"}'
```

### 6.4 查询任务详情

```http
GET /api/tasks/{id}
```

成功返回 `200` 和完整任务对象。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 400 | `INVALID_ARGUMENT` | ID 非法或溢出 |
| 404 | `TASK_NOT_FOUND` | 任务不存在 |

```bash
curl --fail-with-body "$BASE_URL/api/tasks/8"
```

### 6.5 更新任务字段

```http
PUT /api/tasks/{id}
Content-Type: application/json
```

请求：

```json
{
  "title": "完成 Go 全量文档",
  "description": "补齐 API、部署和排障说明",
  "dueAt": null,
  "expectedVersion": 0
}
```

这是完整替换可编辑字段的 `PUT`：`title` 必填；省略或传 `null` 的 `description`、`dueAt` 会被保存为 `null`。成功返回 `200` 和版本加一后的任务。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 404 | `TASK_NOT_FOUND` | 任务不存在 |
| 409 | `TASK_VERSION_CONFLICT` | 版本已过期 |
| 422 | `VALIDATION_FAILED` | 标题、截止时间或版本非法；或缺少版本 |

```bash
curl --fail-with-body \
  -X PUT "$BASE_URL/api/tasks/8" \
  -H 'Content-Type: application/json' \
  -d '{"title":"完成 Go 全量文档","description":null,"dueAt":null,"expectedVersion":0}'
```

### 6.6 修改任务状态

```http
PATCH /api/tasks/{id}/status
Content-Type: application/json
```

请求：

```json
{
  "status": "DOING",
  "expectedVersion": 1
}
```

服务端先检查版本，再检查状态转换。这保证旧版本请求即使目标状态也非法，仍稳定返回版本冲突。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 404 | `TASK_NOT_FOUND` | 任务不存在 |
| 409 | `TASK_VERSION_CONFLICT` | 版本已过期 |
| 409 | `TASK_INVALID_TRANSITION` | 当前状态不能转换到目标状态 |
| 422 | `VALIDATION_FAILED` | 状态、版本非法或缺少版本 |

```bash
curl --fail-with-body \
  -X PATCH "$BASE_URL/api/tasks/8/status" \
  -H 'Content-Type: application/json' \
  -d '{"status":"DOING","expectedVersion":1}'
```

### 6.7 删除任务

```http
DELETE /api/tasks/{id}?expectedVersion=2
```

`expectedVersion` 是必填查询参数。成功返回 `204 No Content`，响应体为空。

| 状态 | 错误码 | 原因 |
|---:|---|---|
| 400 | `INVALID_ARGUMENT` | ID 或版本不是整数 |
| 404 | `TASK_NOT_FOUND` | 任务不存在 |
| 409 | `TASK_VERSION_CONFLICT` | 版本已过期 |
| 422 | `VALIDATION_FAILED` | 缺少版本或版本为负数 |

```bash
curl --fail-with-body -i \
  -X DELETE "$BASE_URL/api/tasks/8?expectedVersion=2"
```

## 7. 错误码索引

| HTTP | 错误码 | 客户端处理建议 |
|---:|---|---|
| 400 | `INVALID_JSON` | 修正 JSON 语法，确保只有一个值 |
| 400 | `UNKNOWN_FIELD` | 删除拼写错误或未定义字段 |
| 400 | `BODY_TOO_LARGE` | 缩小请求体；大文件不要放入该 API |
| 400 | `INVALID_ARGUMENT` | 修正 ID、查询参数或字段类型 |
| 404 | `NOT_FOUND` | 检查路由路径 |
| 404 | `USER_NOT_FOUND` | 刷新用户列表 |
| 404 | `TASK_NOT_FOUND` | 刷新任务列表 |
| 404 | `TASK_OWNER_NOT_FOUND` | 选择存在的负责人 |
| 409 | `EMAIL_CONFLICT` | 更换邮箱或读取已有用户 |
| 409 | `USER_VERSION_CONFLICT` | 重新读取用户后再决定是否修改 |
| 409 | `TASK_VERSION_CONFLICT` | 重新读取任务后再决定是否修改 |
| 409 | `TASK_INVALID_TRANSITION` | 根据状态表选择合法目标状态 |
| 415 | `UNSUPPORTED_MEDIA_TYPE` | 设置 `Content-Type: application/json` |
| 422 | `VALIDATION_FAILED` | 按 `fields` 修正输入 |
| 422 | `TASK_OWNER_DISABLED` | 启用负责人或选择其他负责人 |
| 405 | `METHOD_NOT_ALLOWED` | 查看 `Allow` Header |
| 503 | `NOT_READY` | 等待数据库和迁移就绪 |
| 504 | `DEADLINE_EXCEEDED` | 缩小工作量，检查数据库慢查询后重试 |
| 500 | `INTERNAL_ERROR` | 保存 request ID 并检查服务日志 |

## 8. 一次完整的 curl 流程

下面示例依赖 `jq`：

```bash
USER_JSON=$(curl --fail-with-body -sS \
  -X POST "$BASE_URL/api/users" \
  -H 'Content-Type: application/json' \
  -d '{"name":"张三","email":"user@example.com"}')
USER_ID=$(printf '%s' "$USER_JSON" | jq -r '.data.id')

TASK_JSON=$(curl --fail-with-body -sS \
  -X POST "$BASE_URL/api/tasks" \
  -H 'Content-Type: application/json' \
  -d "{\"ownerId\":$USER_ID,\"title\":\"完成 Go 文档\"}")
TASK_ID=$(printf '%s' "$TASK_JSON" | jq -r '.data.id')
VERSION=$(printf '%s' "$TASK_JSON" | jq -r '.data.version')

TASK_JSON=$(curl --fail-with-body -sS \
  -X PUT "$BASE_URL/api/tasks/$TASK_ID" \
  -H 'Content-Type: application/json' \
  -d "{\"title\":\"完成 Go 全量文档\",\"description\":null,\"dueAt\":null,\"expectedVersion\":$VERSION}")
VERSION=$(printf '%s' "$TASK_JSON" | jq -r '.data.version')

TASK_JSON=$(curl --fail-with-body -sS \
  -X PATCH "$BASE_URL/api/tasks/$TASK_ID/status" \
  -H 'Content-Type: application/json' \
  -d "{\"status\":\"DOING\",\"expectedVersion\":$VERSION}")
VERSION=$(printf '%s' "$TASK_JSON" | jq -r '.data.version')

curl --fail-with-body -sS \
  "$BASE_URL/api/tasks?ownerId=$USER_ID&status=DOING&page=1&pageSize=20" | jq

curl --fail-with-body -i \
  -X DELETE "$BASE_URL/api/tasks/$TASK_ID?expectedVersion=$VERSION"
```

## 9. 兼容性约定

- 字段名、状态字符串、错误码和 HTTP 状态属于公开契约，修改它们需要同步更新集成测试与本文档。
- 新增可选响应字段是向后兼容修改；删除字段、改变类型或改变错误优先级不是。
- 客户端应依据 `error.code` 分支，不应解析中文 `message`。
- 客户端应忽略自己不使用的响应字段，但请求端不能发送未知字段。
- 当前没有批量接口、认证、幂等键和自动重试协议。
