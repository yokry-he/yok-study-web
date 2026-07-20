# 后端接口与服务问题

## 适合谁看

这篇适合开始接触 Node.js、Java、Go、Python 后端项目，或者需要和后端联调的前端同学。这里收录语言无关的接口契约、鉴权、幂等、事务、日志和并发问题；如果根因位于 Node.js 模块加载、事件循环、线程池、Stream 或进程生命周期，请转到 [Node.js 真实项目问题库](/projects/issues-node)。

## 使用方式

后端问题通常不要只看接口返回值，还要一起看：

- 请求参数。
- 请求头和登录态。
- 日志链路。
- 数据库写入结果。
- 重试和并发场景。
- 前端是否重复提交。

## 问题 1：前端明明传了参数，后端却说参数为空

### 问题现象

- 前端 Network 里能看到参数。
- 后端接口返回“参数不能为空”。
- 有些接口正常，有些接口不正常。

### 影响范围

所有接口联调场景，尤其是登录、查询、创建、批量操作和文件上传。

### 常见根因

前后端对参数位置理解不一致：

- 前端放在 query，后端从 body 读。
- 前端提交 JSON，后端按表单格式解析。
- 前端字段名是 `userId`，后端接收 `user_id`。
- 前端数组传法和后端解析规则不一致。
- 请求头 `Content-Type` 不正确。

### 解决方案

先固定接口契约，不要靠猜。

```text
GET /api/users?page=1&pageSize=20&keyword=tom
POST /api/users
Content-Type: application/json

{
  "username": "tom",
  "roleIds": [1, 2]
}
```

前端请求封装要明确区分 query 和 body。

```ts
export function getUserList(params: UserListQuery) {
  return request.get('/api/users', { params })
}

export function createUser(data: CreateUserPayload) {
  return request.post('/api/users', data)
}
```

后端校验错误要返回具体字段。

```json
{
  "code": "VALIDATION_ERROR",
  "message": "参数校验失败",
  "fields": [
    { "field": "username", "message": "用户名不能为空" }
  ]
}
```

### 预防方式

- 每个接口都有明确的 method、path、query、body、response。
- 字段命名统一，不要同一个项目混用多种风格。
- 联调时保存一份成功请求样例。
- 参数校验错误返回字段级提示，前端才能正确展示。

## 问题 2：用户连续点击提交，创建了两条重复数据

### 问题现象

- 用户点击“保存”后页面卡住，又点了一次。
- 后端创建了两条相同数据。
- 前端禁用按钮后仍偶尔复现，因为用户可能刷新、重试或网络层重复发送。

### 影响范围

所有有副作用的接口：

- 创建订单。
- 支付回调。
- 提交审批。
- 创建用户。
- 发放优惠券。
- 扣减库存。

### 常见根因

系统只做了“前端防重复点击”，没有做后端幂等控制。前端可以减少误触，但不能保证请求只到达一次。

### 解决方案

前端层面：提交中禁用按钮。

```ts
const submitting = ref(false)

async function submit() {
  if (submitting.value) return

  submitting.value = true
  try {
    await api.createOrder(form.value)
  } finally {
    submitting.value = false
  }
}
```

后端层面：关键操作必须有幂等键。

```text
POST /api/orders
Idempotency-Key: 20260702-user-1001-submit-001
```

服务端处理逻辑：

```text
1. 检查 Idempotency-Key 是否已处理。
2. 如果已处理，直接返回上次结果。
3. 如果未处理，执行业务写入。
4. 把 key 和结果保存起来。
```

数据库层面：对业务唯一性加约束。

```sql
CREATE UNIQUE INDEX uk_order_request
ON orders(user_id, request_no);
```

### 预防方式

- 有副作用的接口不要只依赖前端禁用按钮。
- 关键业务必须有业务唯一键或幂等键。
- 支付、库存、优惠券、审批必须考虑重复请求。
- 错误重试要区分“重试查询结果”和“重新执行业务动作”。

## 问题 3：接口偶尔很慢，但本地和测试环境都正常

### 问题现象

- 大部分请求 100ms 内返回。
- 偶尔某些请求超过 3 秒，甚至超时。
- 重启服务后短暂恢复，过一会又出现。

### 影响范围

列表查询、导出、报表、复杂详情页、第三方接口聚合页。

### 常见根因

常见原因包括：

- 数据量增长后 SQL 没有命中索引。
- 接口内串行调用多个外部服务。
- 连接池耗尽。
- 日志、文件、对象存储或第三方 API 慢。
- 某些异常分支没有及时释放资源。

### 解决方案

先把一次请求拆成可观察的阶段。

```text
request_start
auth_checked
db_query_done
third_party_done
response_sent
```

日志里必须带请求 ID。

```json
{
  "requestId": "req_123",
  "path": "/api/reports",
  "userId": 1001,
  "dbCostMs": 820,
  "thirdPartyCostMs": 2100,
  "totalCostMs": 3060
}
```

如果 SQL 慢，先看执行计划和索引，不要直接加缓存。

```sql
EXPLAIN SELECT *
FROM orders
WHERE user_id = 1001
ORDER BY created_at DESC
LIMIT 20;
```

如果外部服务慢，要设置超时和降级。

```ts
const result = await requestThirdParty({
  timeout: 2000,
  fallback: null
})
```

### 预防方式

- 所有后端接口记录耗时、状态码、请求 ID。
- 连接数据库、Redis、第三方 API 都要有超时。
- 慢查询要进日志或监控。
- 报表和导出不要和核心在线请求抢资源。

## 问题 4：错误提示只有“服务器错误”，前端无法处理

### 问题现象

- 后端统一返回 500。
- 前端只能弹“服务器错误”。
- 用户不知道是参数问题、权限问题、数据冲突还是系统异常。

### 影响范围

所有用户可感知操作，包括保存、删除、导入、审批、支付和登录。

### 常见根因

后端没有区分错误类型，或者把所有异常都吞成通用错误。

### 解决方案

错误响应要有稳定结构。

```json
{
  "code": "USER_MOBILE_DUPLICATED",
  "message": "手机号已被其他用户使用",
  "requestId": "req_123"
}
```

常见错误分类：

| 类型 | HTTP 状态 | 示例 |
| --- | ---: | --- |
| 参数错误 | 400 | 字段为空、格式不对 |
| 未登录 | 401 | token 失效 |
| 无权限 | 403 | 没有按钮或接口权限 |
| 资源不存在 | 404 | 用户不存在 |
| 业务冲突 | 409 | 重复创建、状态已变更 |
| 系统异常 | 500 | 数据库异常、代码异常 |

前端按错误码处理可恢复场景。

```ts
if (error.code === 'USER_MOBILE_DUPLICATED') {
  formError.mobile = '手机号已被使用'
  return
}
```

### 预防方式

- 后端错误码要稳定，不要频繁改文案当作逻辑判断。
- 前端不要用 message 文案判断业务分支。
- 所有 500 错误必须带 requestId，方便查日志。
- 业务冲突不要返回 500，应返回明确错误码。

## 问题 5：角色授权接口成功返回，但权限只更新了一半

### 问题现象

- 管理员给角色重新分配菜单和按钮权限。
- 接口返回成功。
- 刷新后菜单更新了，但按钮权限没更新；或者按钮权限更新了，菜单没更新。
- 数据库里关系表数据不一致。

### 影响范围

用户角色、角色菜单、按钮权限、组织权限、数据权限、审批流节点授权等多表写入场景。

### 常见根因

后端把一次业务操作拆成多次独立写入，没有事务保护。

```ts
await roleMenuRepository.deleteByRoleId(roleId)
await roleMenuRepository.createMany(roleId, menuIds)

await rolePermissionRepository.deleteByRoleId(roleId)
await rolePermissionRepository.createMany(roleId, permissionIds)
```

如果第三步失败，前两步已经提交，就会出现半更新。

另一类原因是：部分 repository 使用事务对象，部分仍然使用全局数据库连接，导致事务没有真正覆盖所有写入。

### 解决方案

把一次授权操作放进同一个事务上下文。

```ts
await db.transaction(async (tx) => {
  await roleRepository.ensureRoleExists(tx, roleId)

  await roleMenuRepository.replaceByRoleId(tx, roleId, menuIds)
  await rolePermissionRepository.replaceByRoleId(tx, roleId, permissionIds)

  await auditRepository.create(tx, {
    actorId: currentUser.id,
    action: 'role.assign-permissions',
    targetId: roleId,
    detail: { menuIds, permissionIds }
  })
})
```

事务提交后再处理缓存。

```ts
await permissionCache.invalidateUsersByRole(roleId)
```

不要在事务里做慢网络请求，也不要在事务未提交前通知前端或消息队列“授权成功”。

### 预防方式

- 业务上必须一起成功或一起失败的写入都要放进事务。
- repository 函数明确接收 `tx`，避免误用全局 `db`。
- 授权接口增加失败路径测试。
- 关键写操作写审计日志。
- 权限缓存失效放在事务提交后。

## 问题 6：401 和 403 混用，前端登录态和权限处理混乱

### 问题现象

- 用户没有登录时，有的接口返回 403，有的返回 401。
- 用户已登录但无权限时，前端却跳回登录页。
- token 过期后页面弹出多个“无权限”提示。
- 测试和前端无法判断应该重新登录还是展示无权限页。

### 影响范围

所有需要登录和授权的系统：后台管理、SaaS 控制台、企业内部系统、开放平台。

### 常见根因

认证和授权没有分清：

- 认证：你是谁。
- 授权：你能做什么。

常见错误是把 token 失效、无权限、账号禁用、角色缺失都统一返回一个状态码。

### 解决方案

明确状态码语义。

| 场景 | HTTP 状态 | 错误码 |
| --- | ---: | --- |
| 没有 token | 401 | `UNAUTHENTICATED` |
| token 过期 | 401 | `TOKEN_EXPIRED` |
| token 有效但账号禁用 | 403 | `ACCOUNT_DISABLED` |
| 已登录但缺权限 | 403 | `FORBIDDEN` |
| 资源不存在 | 404 | `RESOURCE_NOT_FOUND` |

后端中间件建议分两层：

```text
authenticate: 解析 token，得到 currentUser
authorize: 检查 currentUser 是否拥有权限码
```

前端响应策略：

```ts
if (error.status === 401) {
  authStore.logout()
  router.replace('/login')
}

if (error.status === 403) {
  showForbiddenMessage()
}
```

### 预防方式

- API 规范中写清 401 和 403 的语义。
- 权限接口必须有测试覆盖未登录、无权限、账号禁用。
- 前端不要用错误文案判断登录态。
- 登录失效处理要加全局锁，避免多个接口同时弹多次提示。

## 下一步学习

- [Node.js HTTP API 开发](/node/http-api)
- [错误处理与日志](/node/error-logging)
- [Node 权限 API 从零到项目](/node/permission-api-project)
- [Node.js 真实项目问题库](/projects/issues-node)
- [Node.js 专项练习](/roadmap/node-practice)
- [Node.js 常见问题](/node/troubleshooting)
- [数据库事务、锁与并发](/database/transactions)
