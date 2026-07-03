# 测试策略

## 适合谁看

适合已经能写 Node.js API，但每次改接口都靠手动 Postman 或前端页面验证的人：

- 改了登录逻辑，不知道有没有影响其他接口。
- 数据库 service 没有测试。
- 错误处理中间件改动后只靠手动点。
- CI 里没有测试门禁。
- 不知道单元测试、集成测试、端到端测试怎么分。

Node.js 后端测试的目标不是追求覆盖率数字，而是把核心业务、权限、安全边界和错误场景固定下来，防止回归。

## 测试分层

| 类型 | 测什么 | 示例 |
| --- | --- | --- |
| 单元测试 | 纯函数、service 小逻辑 | 权限判断、参数转换 |
| 集成测试 | 多层协作 | API + service + database |
| 端到端测试 | 真实流程 | 登录 -> 创建用户 -> 查询列表 |
| 契约测试 | 接口输入输出 | 响应字段、错误码 |

初学项目可以先覆盖：

- service 业务规则。
- auth / permission 中间件。
- 关键 API 成功和失败路径。
- 数据库事务或 repository。

## Node 内置测试运行器

Node.js 提供 `node:test`。

```ts
import test from 'node:test'
import assert from 'node:assert/strict'

function sum(a: number, b: number) {
  return a + b
}

test('sum adds two numbers', () => {
  assert.equal(sum(1, 2), 3)
})
```

运行：

```bash
node --test
```

项目也可以使用 Vitest、Jest 等工具。关键不是工具名，而是测试边界是否清楚。

## service 单元测试

权限判断适合单元测试：

```ts
function canDeleteUser(currentUser: User, targetUser: User) {
  if (targetUser.roles.includes('admin')) return false
  return currentUser.permissions.includes('user:delete')
}
```

测试：

```ts
test('cannot delete admin user', () => {
  const currentUser = {
    permissions: ['user:delete']
  } as User

  const targetUser = {
    roles: ['admin']
  } as User

  assert.equal(canDeleteUser(currentUser, targetUser), false)
})
```

这种测试不需要启动 HTTP 服务，也不需要数据库，速度快、定位准。

## API 集成测试

接口测试要覆盖：

- 正常成功。
- 参数错误。
- 未登录。
- 无权限。
- 资源不存在。
- 服务端错误兜底。

示例断言结构：

```ts
test('GET /api/users requires auth', async () => {
  const response = await request(app).get('/api/users')

  assert.equal(response.status, 401)
  assert.equal(response.body.message, '未登录')
})
```

如果项目没有引入 HTTP 测试库，也可以先对 controller/service 做集成测试。

## 数据库测试

数据库测试要解决两个问题：

1. 数据从哪里来。
2. 测试后如何清理。

常见策略：

| 策略 | 适合 |
| --- | --- |
| 测试数据库 | 最接近真实情况 |
| 每个测试事务回滚 | 保持数据干净 |
| 测试前 seed | 固定基础数据 |
| repository mock | 单元测试 service |

不要让测试依赖线上或共享开发库。测试数据必须可控。

## mock 的边界

不要所有东西都 mock。否则测试只证明 mock 写得对。

适合 mock：

- 邮件、短信、支付。
- 第三方 API。
- 时间。
- 随机数。
- 文件上传。

不建议 mock：

- 核心权限判断。
- 核心业务规则。
- 你正想验证的数据库查询。

## 错误路径测试

很多项目只测成功路径，这是不够的。

应该补：

```text
参数缺失 -> 400
未登录 -> 401
无权限 -> 403
不存在 -> 404
重复提交 -> 409
内部异常 -> 500
```

错误响应结构要稳定：

```json
{
  "code": "USER_NOT_FOUND",
  "message": "用户不存在"
}
```

前端才能可靠处理。

## CI 中运行测试

常见脚本：

```json
{
  "scripts": {
    "test": "node --test",
    "typecheck": "tsc --noEmit",
    "build": "npm run typecheck && npm test"
  }
}
```

实际项目中 build 是否执行 test 要看耗时，但 CI 至少应该执行：

- lint 或格式检查。
- typecheck。
- test。
- build。

## 实际项目常见问题

### 1. 测试很慢，没人愿意跑

**解决方案**

- 单元测试和集成测试分开命令。
- 核心 service 测试保持无数据库。
- 数据库测试只覆盖关键路径。

### 2. 测试依赖执行顺序

**原因**

测试之间共享状态。

**解决方案**

每个测试准备自己的数据，或测试后回滚/清理。

### 3. 本地通过，CI 失败

**常见原因**

- 环境变量缺失。
- 测试数据库未准备。
- 时间、时区、随机数不稳定。
- 测试依赖真实外部服务。

### 4. 只测成功路径

权限、登录态、错误码、事务回滚这些失败路径更容易在线上出事故。它们应该进入测试。

## 最佳实践

- service 业务规则优先写单元测试。
- auth、permission、错误处理中间件要测试失败路径。
- 数据库测试使用独立测试库。
- 第三方服务用 mock 或测试沙箱。
- CI 至少跑 typecheck、test、build。
- 修复线上 bug 后补回归测试。

## 学习检查

学完本节后，你应该能回答：

- 单元测试、集成测试、端到端测试分别测什么。
- 为什么不要所有东西都 mock。
- API 测试为什么必须覆盖 401/403/404/500。
- 数据库测试如何保持数据干净。
- CI 里为什么要有测试门禁。

## 参考资料

- [Node.js Test Runner](https://nodejs.org/api/test.html)

## 下一步学习

继续学习 [Node.js 安全基础](/node/security)，把测试覆盖和安全边界结合起来。
