# Docker 速查

## 常用命令

```bash
docker ps
docker ps -a
docker images
docker logs app
docker stop app
docker rm app
```

| 命令 | 用途 |
| --- | --- |
| `docker ps` | 查看运行中的容器 |
| `docker ps -a` | 查看所有容器 |
| `docker images` | 查看镜像 |
| `docker logs <container>` | 查看日志 |
| `docker exec -it <container> sh` | 进入容器 |
| `docker stop <container>` | 停止容器 |
| `docker rm <container>` | 删除容器 |

## 运行容器

```bash
docker run --name web -p 8080:80 nginx
```

含义：

```text
宿主机 8080 -> 容器内 80
```

后台运行：

```bash
docker run -d --name web -p 8080:80 nginx
```

挂载目录：

```bash
docker run -d \
  --name web \
  -p 8080:80 \
  -v $(pwd)/dist:/usr/share/nginx/html \
  nginx
```

## Dockerfile 基础

前端构建示例：

```dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
```

关键点：

- 先复制依赖文件，再安装依赖，利用缓存。
- 构建阶段和运行阶段分离。
- 前端静态资源最终由 Nginx 或静态服务托管。

## docker compose

```yaml
services:
  web:
    build: .
    ports:
      - '8080:80'
    depends_on:
      - api

  api:
    image: node:20-alpine
    working_dir: /app
    command: npm run start
    volumes:
      - ./server:/app
```

启动：

```bash
docker compose up -d
```

查看日志：

```bash
docker compose logs -f web
```

停止：

```bash
docker compose down
```

## 常见排查

容器运行但访问不到：

```bash
docker ps
docker logs app
docker exec -it app sh
```

容器内检查服务：

```bash
curl http://127.0.0.1:3000/health
```

宿主机检查映射：

```bash
curl http://127.0.0.1:8080/health
```

清理未使用资源：

```bash
docker system prune
```

## 常见坑

| 问题 | 处理 |
| --- | --- |
| 容器内服务只监听 127.0.0.1 | 改成监听 `0.0.0.0` |
| 端口访问不到 | 检查 `宿主机端口:容器端口` |
| 镜像很大 | 使用多阶段构建 |
| 构建每次都很慢 | 优化 COPY 顺序 |
| 容器数据丢失 | 使用 volume 持久化 |

## 项目建议

- 每个服务提供 `/health`。
- README 写清容器端口和宿主机端口。
- 不把密钥写进 Dockerfile。
- 生产镜像只包含运行必需文件。
- 前端 history 路由仍需要 Nginx fallback。

## 下一步学习

- [Docker 容器化](/devops/docker)
- [Nginx 静态部署与代理](/devops/nginx)
- [部署、缓存与 DevOps 问题](/projects/issues-deployment)
