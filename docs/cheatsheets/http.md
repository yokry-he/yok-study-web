# HTTP 速查

## 常见方法

| 方法 | 用途 | 示例 |
| --- | --- | --- |
| `GET` | 查询资源 | 查询用户列表 |
| `POST` | 创建资源或提交动作 | 创建订单 |
| `PUT` | 整体更新资源 | 更新用户完整信息 |
| `PATCH` | 局部更新资源 | 修改用户状态 |
| `DELETE` | 删除资源 | 删除角色 |

接口设计要让方法和语义匹配，不要所有操作都用 `POST`。

## 常见状态码

| 状态码 | 含义 | 常见处理 |
| ---: | --- | --- |
| `200` | 成功 | 正常读取响应 |
| `201` | 创建成功 | 创建资源后返回 |
| `204` | 成功但无内容 | 删除成功常用 |
| `400` | 参数错误 | 展示字段错误 |
| `401` | 未登录 | 跳登录或刷新 token |
| `403` | 无权限 | 展示无权限 |
| `404` | 不存在 | 展示空状态或错误页 |
| `409` | 业务冲突 | 提示重复、状态变化 |
| `500` | 服务异常 | 记录 requestId 并提示稍后再试 |

不要把业务冲突都返回 500。前端需要稳定状态码和错误码来处理。

## 请求头

```http
Authorization: Bearer <token>
Content-Type: application/json
Accept: application/json
X-Request-Id: req_123
```

| 请求头 | 用途 |
| --- | --- |
| `Authorization` | 携带 token |
| `Content-Type` | 声明请求体格式 |
| `Accept` | 声明期望响应格式 |
| `Cookie` | 携带 cookie |
| `X-Request-Id` | 链路追踪 |

JSON 请求：

```ts
fetch('/api/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    username: 'Tom'
  })
})
```

## Query 和 Body

Query 适合查询条件：

```text
GET /api/users?page=1&pageSize=20&keyword=tom
```

Body 适合提交数据：

```http
POST /api/users
Content-Type: application/json

{
  "username": "Tom",
  "mobile": "13800000000"
}
```

常见联调问题是：前端把参数放 query，后端从 body 读，或者字段名不一致。

## REST 风格示例

```text
GET    /api/users             查询用户列表
GET    /api/users/1           查询用户详情
POST   /api/users             创建用户
PUT    /api/users/1           更新用户
PATCH  /api/users/1/status    修改用户状态
DELETE /api/users/1           删除用户
```

批量操作：

```http
POST /api/users/batch-delete
Content-Type: application/json

{
  "ids": [1, 2, 3]
}
```

## CORS

浏览器跨域由同源策略限制。

常见响应头：

```http
Access-Control-Allow-Origin: https://example.com
Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE
Access-Control-Allow-Headers: Content-Type,Authorization
Access-Control-Allow-Credentials: true
```

如果携带 cookie：

- `Access-Control-Allow-Credentials` 必须为 `true`。
- `Access-Control-Allow-Origin` 不能是 `*`。
- 前端请求要设置 `credentials`。

```ts
fetch('/api/profile', {
  credentials: 'include'
})
```

## 缓存

常见响应头：

```http
Cache-Control: no-cache
Cache-Control: public, max-age=31536000, immutable
ETag: "abc123"
```

| 场景 | 推荐 |
| --- | --- |
| `index.html` | 不强缓存 |
| hash 静态资源 | 长缓存 |
| 用户接口 | 通常不缓存或短缓存 |
| 公共配置 | 根据更新频率设置 |

前端上线旧页面问题通常和入口文件缓存有关。

## 常见问题

| 问题 | 处理 |
| --- | --- |
| 401 重复弹窗 | 全局只处理一次未登录 |
| 403 被当成 500 | 前后端区分权限错误 |
| 参数为空 | 检查 query、body、字段名和 Content-Type |
| 跨域失败 | 看预检请求和 CORS 响应头 |
| 用户看到旧页面 | 检查 Cache-Control 和 CDN |

## 项目建议

- 接口契约写清 method、path、query、body、response。
- 错误响应包含稳定 code 和 requestId。
- 前端不要用 message 文案判断业务逻辑。
- 所有全局请求错误处理都考虑并发。
- 联调问题先看 Network，再看后端日志。

## 下一步学习

- [HTTP 与请求流程](/browser/http-request)
- [跨域与登录态](/browser/cors-auth)
- [后端接口与服务问题](/projects/issues-backend)
- [前端页面与状态问题](/projects/issues-frontend)
