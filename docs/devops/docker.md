# Docker 容器化

## 适合谁看

适合想把前端、Node API 或整套开发环境打包运行的人。

Docker 解决的核心问题是：让应用和运行环境一起交付，减少“我本地可以，服务器不行”的差异。

## 基本概念

| 概念 | 含义 |
| --- | --- |
| Dockerfile | 构建镜像的说明书 |
| Image | 构建出来的镜像模板 |
| Container | 镜像运行后的实例 |
| Volume | 持久化数据或挂载目录 |
| Network | 容器之间通信的网络 |
| Compose | 用一个文件管理多个服务 |

关系：

```text
Dockerfile -> docker build -> image -> docker run -> container
```

## 前端项目 Dockerfile

Vite 前端常见多阶段构建：

```dockerfile
FROM node:22-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

第一阶段负责构建，第二阶段只保留 Nginx 和构建产物。这样生产镜像不会包含完整 Node 构建环境和源码。

## Node API Dockerfile

```dockerfile
FROM node:22-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci

FROM node:22-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production

COPY --from=deps /app/node_modules ./node_modules
COPY . .

EXPOSE 3000
CMD ["node", "server.js"]
```

实际项目里还要考虑：

- 是否需要 TypeScript 编译。
- 是否要只安装 production dependencies。
- 是否需要健康检查。
- 是否要使用非 root 用户运行。

## .dockerignore

不要把无关文件复制进镜像：

```text
node_modules
dist
.git
.env
*.log
```

`.env` 不应该打进镜像。环境变量应在运行时注入。

## Docker Compose

Compose 状态需要同时阅读进程状态、健康检查、一次性迁移退出码和重启次数。`running` 只说明主进程仍在，不代表服务已经就绪。

<DocFigure
  src="/images/devops/docker-container-state.webp"
  alt="Docker Compose 状态报告展示 PostgreSQL 和 API healthy、迁移 exited 0、Worker 反复 restarting"
  caption="重启策略可能掩盖持续崩溃；看到 restarting 时先保存第一条错误和退出码。"
  :width="1440"
  :height="900"
/>

迁移容器正常完成后应 `exited 0`，而长驻 API 应保持 `running (healthy)`；两者不能使用同一健康判断。

Compose 适合本地或中小型部署管理多服务：

```yaml
services:
  web:
    build: ./frontend
    ports:
      - "8080:80"

  api:
    build: ./backend
    environment:
      NODE_ENV: production
      DATABASE_URL: postgres://app:secret@db:5432/app
    depends_on:
      - db

  db:
    image: postgres:18
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: app
    volumes:
      - db-data:/var/lib/postgresql/data

volumes:
  db-data:
```

服务之间可以用服务名通信，例如 API 连接数据库时使用 `db:5432`。

## 构建和运行

```bash
docker build -t vue-admin:latest .
docker run -p 8080:80 vue-admin:latest
```

Compose：

```bash
docker compose up -d --build
docker compose logs -f
docker compose ps
docker compose down
```

## 环境变量

构建时变量和运行时变量要分清楚。

前端 Vite 项目的 `VITE_*` 通常在构建时就会被写进产物：

```bash
VITE_API_BASE_URL=/api npm run build
```

Node 服务可以在运行时读取：

```bash
docker run -e DATABASE_URL=postgres://... api:latest
```

不要以为修改容器环境变量后，已构建好的前端静态文件也会自动变化。前端静态资源需要重新构建，除非你设计了运行时配置文件。

## 实际项目问题

### 问题：镜像很大

**原因**

- 单阶段构建把源码、依赖缓存、构建工具都带进了镜像。
- 没有 `.dockerignore`。
- 复制了 `node_modules`、日志或临时文件。

**解决方案**

- 使用多阶段构建。
- 添加 `.dockerignore`。
- 只把运行时需要的文件复制到最终镜像。

### 问题：容器里服务启动了，但外部访问不到

**排查**

1. 容器内服务是否监听 `0.0.0.0`。
2. 是否做了端口映射。
3. `docker ps` 中端口是否正确。
4. 防火墙或安全组是否放行。

Node 服务不要只监听 `localhost`：

```ts
server.listen(3000, '0.0.0.0')
```

### 问题：容器删除后数据库数据没了

**原因**

数据库数据没有挂载 volume。

**解决方案**

```yaml
volumes:
  - db-data:/var/lib/postgresql/data
```

## 最佳实践

- 生产镜像使用多阶段构建。
- 使用 `.dockerignore` 控制上下文。
- 密钥通过运行时环境变量或密钥管理注入，不写进镜像。
- 数据库必须挂载 volume。
- 前端构建变量和后端运行变量要分开理解。
- 先学 Dockerfile 和 Compose，再考虑 Kubernetes。

## 参考资料

- [Docker Docs: Multi-stage builds](https://docs.docker.com/build/building/multi-stage/)
- [Docker Docs: Building best practices](https://docs.docker.com/build/building/best-practices/)
- [Docker Docs: Compose file reference](https://docs.docker.com/reference/compose-file/)
- [Docker Docs: Dockerfile reference](https://docs.docker.com/reference/dockerfile/)

## 下一步学习

继续学习 [CI/CD 自动化发布](/devops/ci-cd)。
