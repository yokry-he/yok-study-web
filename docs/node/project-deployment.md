# 项目结构与部署

## 适合谁看

适合已经能写 Node.js API，但还不知道如何组织目录、管理环境变量、启动生产进程和部署到服务器的人。

## 推荐结构

```text
src/
├─ config/
├─ routes/
├─ services/
├─ repositories/
├─ schemas/
├─ middlewares/
├─ utils/
└─ main.ts
```

## 环境变量

```ini
NODE_ENV=production
PORT=3000
DATABASE_URL=postgres://user:pass@host:5432/db
JWT_SECRET=change-me
```

不要把 `.env` 提交到仓库。提供 `.env.example`。

## 启动脚本

```json
{
  "scripts": {
    "dev": "tsx watch src/main.ts",
    "build": "tsc",
    "start": "node dist/main.js"
  }
}
```

## 部署检查

上线前检查：

- Node 版本。
- 环境变量。
- 数据库连接。
- 端口监听。
- 健康检查。
- 日志输出。
- 进程守护。
- 回滚方案。

## 进程管理

生产环境不能只手动执行 `node dist/main.js` 后关闭终端。

常见方式：

- systemd。
- PM2。
- Docker。
- 平台托管服务。

## Nginx 反向代理

```nginx
location /api/ {
  proxy_pass http://127.0.0.1:3000/;
}
```

## 实际项目常见问题

### 1. 本地能连数据库，线上连不上

**排查**

- 环境变量是否正确。
- 网络安全组是否放行。
- 数据库账号权限。
- 服务器 DNS 和端口。

### 2. 进程崩溃后服务不可用

**解决方案**

使用进程管理器或容器编排，配置自动重启和日志采集。

### 3. 环境变量缺失

启动时校验配置，不要等接口请求时才报错。

## 最佳实践

- 提供 `.env.example`。
- 启动时校验必要配置。
- 暴露 `/health`。
- 使用进程管理。
- 部署文档写清楚启动、停止、回滚。

## 下一步

继续学习 [常见问题](/node/troubleshooting)。
