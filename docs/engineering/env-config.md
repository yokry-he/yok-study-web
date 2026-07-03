# 环境配置

## 适合谁看

适合准备把项目从本地开发推进到测试、预发、生产环境的学习者。

很多线上事故不是业务代码写错，而是环境变量、接口地址、部署路径、缓存策略不一致。环境配置文档的目标是让“项目在哪个环境怎么运行”变得清楚。

## 常见环境

| 环境 | 用途 | 常见特点 |
| --- | --- | --- |
| development | 本地开发 | 使用 Vite dev server 和本地代理 |
| test | 测试环境 | 接测试后端、测试账号、测试数据 |
| staging | 预发环境 | 尽量接近生产 |
| production | 生产环境 | 面向真实用户 |

## Vite 环境文件

```text
.env
.env.development
.env.test
.env.production
```

示例：

```ini
VITE_APP_TITLE=Vue Admin
VITE_API_BASE_URL=/api
VITE_UPLOAD_URL=/api/files/upload
```

读取：

```ts
export const appConfig = {
  title: import.meta.env.VITE_APP_TITLE,
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  uploadUrl: import.meta.env.VITE_UPLOAD_URL
}
```

## 什么可以放前端环境变量

可以放：

- API 前缀。
- 应用标题。
- 静态资源前缀。
- 是否开启 mock。
- 公开的埋点项目 id。

不要放：

- 数据库密码。
- API 私钥。
- 云服务密钥。
- JWT 签名密钥。
- 后端内部地址和敏感配置。

原因很简单：前端构建产物会发到浏览器，用户可以看到。

## 环境配置集中管理

不要在页面里到处写 `import.meta.env`。推荐集中到一个配置文件：

```ts
// src/config/app.ts
export const appConfig = {
  title: import.meta.env.VITE_APP_TITLE || 'Vue Admin',
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || '/api',
  uploadUrl: import.meta.env.VITE_UPLOAD_URL || '/api/files/upload'
}
```

其他地方使用：

```ts
request.defaults.baseURL = appConfig.apiBaseUrl
```

好处：

- 缺省值集中。
- 类型和命名集中。
- 排查问题时不用全项目搜索。

## 实际项目常见问题

### 1. 测试环境请求到了生产接口

**原因**

- 构建命令使用了错误 mode。
- `.env.production` 被错误复用。
- 部署平台注入了旧变量。

**解决方案**

明确构建命令：

```bash
vite build --mode test
vite build --mode production
```

并在页面角落或构建信息中展示当前环境，仅测试环境可见：

```ts
const mode = import.meta.env.MODE
```

### 2. 修改 `.env` 后页面没变化

**原因**

Vite 启动时读取环境变量。运行中修改 `.env` 不会自动全部生效。

**解决方案**

重启开发服务器。

### 3. 环境变量在代码里是 `undefined`

**排查**

- 是否以 `VITE_` 开头。
- 是否重启了 dev server。
- 是否写在了正确的 `.env.[mode]` 文件。
- 构建命令的 mode 是否正确。

### 4. 不同环境接口返回结构不一致

**症状**

本地正常，测试环境字段缺失，生产环境又是另一种结构。

**解决方案**

- 在 API 层定义类型。
- 在 service 层做数据适配。
- 和后端确认接口版本。
- 对关键字段做兜底和错误提示。

## 环境文档模板

项目 README 应至少记录：

```text
环境名称：
访问地址：
构建命令：
API 前缀：
后端服务：
是否启用 mock：
账号来源：
部署负责人：
回滚方式：
```

## 最佳实践

- 环境变量集中读取，不散落在页面中。
- 前端变量不保存密钥。
- 每个环境有明确构建命令。
- 测试和生产接口地址必须可追踪。
- 修改环境变量、代理、部署路径时同步更新文档。

## 下一步学习

继续学习 [Vite 工程基础](/engineering/vite)、[构建与部署](/engineering/build-deploy) 和 [工程化常见问题](/engineering/troubleshooting)。
