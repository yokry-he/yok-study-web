# 项目上线全流程实践

## 适合谁看

适合已经能本地运行 Vue、Node、Java 或 Go 项目，但不知道如何把项目稳定发布到线上环境的人。

这篇用一个常见组合做案例：Vue Admin 前端、后端 API、MySQL、Redis、Nginx、Docker Compose、CI/CD。它不是让你一次学会所有云平台，而是先掌握通用上线链路。

## 这篇解决什么

很多上线问题不是代码不会写，而是链路没想清楚：

- 前端构建产物放在哪里。
- `/api` 到底代理到哪个后端。
- 环境变量是构建时注入还是运行时读取。
- Docker 容器端口和宿主机端口怎么对应。
- 发布失败后怎么回滚。
- 线上出错时先看哪一层日志。

完整上线链路可以先理解成：

```mermaid
flowchart LR
  A[开发提交代码] --> B[CI 拉取代码]
  B --> C[安装依赖]
  C --> D[检查和测试]
  D --> E[构建前端和后端]
  E --> F[生成镜像或静态产物]
  F --> G[发布到服务器]
  G --> H[Nginx 对外提供访问]
  H --> I[发布后验证]
  I --> J{是否异常}
  J -- 否 --> K[观察日志和指标]
  J -- 是 --> L[回滚到上一版本]
```

## 案例架构

目标是上线一个后台管理系统：

- 前端：Vue Admin，构建后生成静态资源。
- 后端：API 服务，提供登录、用户、角色、菜单等接口。
- 数据库：MySQL，保存业务数据。
- 缓存：Redis，保存登录态、权限缓存或验证码。
- 网关：Nginx，负责静态资源、history fallback、接口反向代理和缓存头。

```mermaid
flowchart TD
  Browser[浏览器] --> CDN[CDN 可选]
  CDN --> Nginx[Nginx]
  Nginx --> Static[前端静态资源]
  Nginx --> API[后端 API 容器]
  API --> MySQL[MySQL]
  API --> Redis[Redis]
```

学习阶段可以先不用 CDN，但要提前理解 CDN 会缓存静态资源，后续上线白屏、旧页面、资源 404 常常和它有关。

## 第一步：本地生产构建

不要只验证 `npm run dev`。上线前至少要跑生产构建：

```bash
npm run build
```

如果项目有检查命令，也要一起跑：

```bash
npm run lint
npm run test
```

前端构建后要确认：

- `dist/index.html` 存在。
- `dist/assets` 有带 hash 的 JS/CSS 文件。
- 构建日志里的 API 地址符合目标环境。
- 如果部署在子路径，Vite `base` 配置正确。

后端构建后要确认：

- 程序能读取生产环境变量。
- 数据库连接地址不是本地地址。
- 日志会输出到 stdout 或指定日志目录。
- 健康检查接口可访问。

## 第二步：区分构建时配置和运行时配置

这是前端上线最容易混乱的地方。

Vite 的 `import.meta.env` 是构建时配置。也就是说：

```ts
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
```

构建结束后，API 地址已经写进 JS 文件里。你把静态文件复制到另一台服务器，再改服务器环境变量，不会改变这个地址。

如果同一份前端包要部署到多个环境，建议用运行时配置：

```html
<script src="/runtime-config.js"></script>
```

```js
window.__APP_CONFIG__ = {
  apiBaseUrl: '/api'
}
```

两种方案对比：

| 方案 | 优点 | 风险 |
| --- | --- | --- |
| 构建时配置 | 简单，适合每个环境单独构建 | 构建包不能随便跨环境复用 |
| 运行时配置 | 一份包可部署多个环境 | 需要额外管理 `runtime-config.js` 缓存 |

初学阶段可以先用构建时配置，但要在发布文档里写清楚“哪个环境对应哪个构建命令”。

## 第三步：Nginx 部署前端

最小可用配置：

```nginx
server {
  listen 80;
  server_name admin.example.com;

  root /var/www/admin/current;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }

  location /api/ {
    proxy_pass http://127.0.0.1:3000/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}
```

这段配置做了三件事：

- `root` 指向当前前端版本目录。
- `try_files` 解决 Vue Router history 刷新 404。
- `/api/` 代理到后端服务。

注意 `proxy_pass` 末尾的 `/` 会影响路径是否保留。实际项目里要和后端接口路径统一测试。

## 第四步：静态资源缓存策略

前端资源通常分两类：

```mermaid
flowchart TD
  A[前端构建产物] --> B[index.html]
  A --> C[assets 带 hash 文件]
  B --> D[不强缓存]
  C --> E[长期缓存]
```

推荐配置：

```nginx
location = /index.html {
  add_header Cache-Control "no-cache, no-store, must-revalidate";
}

location /assets/ {
  add_header Cache-Control "public, max-age=31536000, immutable";
}
```

原因很简单：

- `index.html` 是入口，必须尽快拿到新版本。
- `assets/app.abc123.js` 带 hash，内容变了文件名也会变，可以长期缓存。

如果把 `index.html` 也强缓存，用户可能一直引用旧 JS 文件，出现旧页面或白屏。

## 第五步：Docker Compose 组织服务

学习阶段推荐用 Docker Compose 先把后端、MySQL、Redis 跑起来。

```yaml
services:
  api:
    image: admin-api:1.0.0
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: mysql://admin:password@mysql:3306/admin
      REDIS_URL: redis://redis:6379
      NODE_ENV: production
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8
    environment:
      MYSQL_DATABASE: admin
      MYSQL_USER: admin
      MYSQL_PASSWORD: password
      MYSQL_ROOT_PASSWORD: root-password
    volumes:
      - mysql-data:/var/lib/mysql

  redis:
    image: redis:7
    command: ["redis-server", "--appendonly", "yes"]
    volumes:
      - redis-data:/data

volumes:
  mysql-data:
  redis-data:
```

这里有几个关键点：

- API 连接 MySQL 时使用服务名 `mysql`，不是 `localhost`。
- 数据库和 Redis 使用 volume 保存数据。
- 容器内端口和宿主机端口要分清。
- 密码不要提交到公开仓库，生产环境用平台密钥或 `.env` 管理。

## 第六步：CI/CD 发布流水线

第一版流水线不需要复杂，但至少要包含检查、构建、发布和验证。

```mermaid
flowchart TD
  A[push main] --> B[安装依赖]
  B --> C[代码检查]
  C --> D[单元测试]
  D --> E[构建前端]
  E --> F[构建后端镜像]
  F --> G[上传产物]
  G --> H[切换 current 版本]
  H --> I[健康检查]
```

前端静态发布建议使用版本目录：

```text
/var/www/admin/releases/2026-07-02-153000
/var/www/admin/releases/2026-07-02-160500
/var/www/admin/current -> /var/www/admin/releases/2026-07-02-160500
```

这样回滚时只需要把 `current` 指回上一个版本。

## 第七步：发布后验证

发布不是“文件传完就结束”。必须验证用户真实访问路径。

| 验证项 | 怎么验证 |
| --- | --- |
| 首页 | 打开域名，确认没有白屏 |
| 二级路由 | 直接访问 `/users`、`/roles` 并刷新 |
| 登录 | 使用测试账号登录 |
| API 代理 | Network 中确认请求到 `/api` 且状态正确 |
| 静态资源 | JS/CSS 没有 404 |
| 权限菜单 | 登录后菜单和按钮权限正确 |
| 健康检查 | 后端 `/health` 返回正常 |
| 日志 | Nginx 和后端没有新增错误 |

发布后至少观察 10 到 30 分钟。小项目可以人工观察，大项目要接入日志、指标和告警。

## 第八步：回滚方案

回滚不是失败时才临时想。上线前就要写好。

前端静态资源回滚：

```bash
ln -sfn /var/www/admin/releases/2026-07-02-153000 /var/www/admin/current
nginx -s reload
```

后端容器回滚：

```bash
docker compose pull api
docker compose up -d api
```

如果用镜像标签，要保留上一个稳定版本，例如：

```text
admin-api:2026-07-02-153000
admin-api:2026-07-02-160500
```

数据库迁移最难回滚。上线前必须区分：

- 只新增表或字段，通常容易回滚。
- 删除字段、改字段含义、批量修改数据，风险很高。
- 已经被新代码写入的数据，不能简单反向执行。

数据库变更要优先设计向前兼容方案：

1. 先加新字段，旧字段继续保留。
2. 双写或回填数据。
3. 新代码切换读取新字段。
4. 稳定后再删除旧字段。

## 第九步：线上排查顺序

当用户说“打不开”时，不要直接改代码。先按链路排查：

```mermaid
flowchart TD
  A[用户反馈打不开] --> B[浏览器 Console]
  B --> C[Network 请求]
  C --> D[CDN 缓存]
  D --> E[Nginx access/error log]
  E --> F[后端日志]
  F --> G[容器状态]
  G --> H[数据库和 Redis]
  H --> I[最近一次发布记录]
```

常用命令：

```bash
docker ps
docker logs admin-api --tail 100
docker compose ps
nginx -t
tail -n 100 /var/log/nginx/error.log
tail -n 100 /var/log/nginx/access.log
```

排查时先确认现象：

- 是所有用户都打不开，还是部分用户。
- 是首页打不开，还是某个路由打不开。
- 是页面白屏，还是接口报错。
- 是发布后立刻出现，还是运行一段时间后出现。

这些问题会决定你从缓存、Nginx、后端、数据库还是资源路径开始查。

## 实际项目常见问题

### 问题 1：容器正常运行，但 Nginx 代理 502

常见原因：

- API 容器实际没有监听 `0.0.0.0`。
- Nginx 代理地址写成了错误端口。
- API 容器在 Docker 网络里，但 Nginx 在宿主机上，地址没打通。
- 后端启动慢，健康检查还没通过。

处理方式：

1. 在服务器上访问 API 健康检查。
2. 看 Nginx error log。
3. 看 API 容器日志。
4. 确认 `proxy_pass` 地址和端口。

### 问题 2：新版本上线后只有部分用户白屏

优先检查缓存：

- `index.html` 是否被 CDN 或浏览器强缓存。
- 新 `index.html` 引用的 JS 文件是否已上传。
- 老版本 JS 文件是否被提前删除。
- CDN 是否只刷新了 assets，没刷新入口文件。

解决方式：

- `index.html` 设置 no-cache。
- assets 保留多版本。
- 发布顺序先上传 assets，再切换 index。
- 回滚时也要同步刷新入口文件。

### 问题 3：环境变量改了，但页面请求地址没变

如果是 Vite 前端，原因通常是变量在构建时已经写入静态资源。

解决方式：

- 重新用正确环境构建。
- 或改成运行时配置。
- 发布说明里明确“改 API 地址需要重新构建还是只改配置文件”。

### 问题 4：数据库迁移成功，后端启动后报字段不存在

常见原因：

- 后端连错数据库。
- CI/CD 跑迁移的环境和后端运行环境不一致。
- 多实例部署时，一部分实例还是旧代码。
- 迁移顺序和代码发布顺序不兼容。

处理方式：

1. 确认后端打印的数据库地址和环境名。
2. 在目标数据库执行 `SHOW COLUMNS` 或等价命令。
3. 确认所有后端实例版本一致。
4. 对破坏性变更改成兼容发布。

## 上线前总检查清单

| 类别 | 检查项 |
| --- | --- |
| 代码 | 构建、检查、测试都通过 |
| 配置 | API 地址、数据库地址、Redis 地址、环境名正确 |
| 前端 | `base`、路由刷新、静态资源缓存策略正确 |
| Nginx | `nginx -t` 通过，代理路径和 fallback 正确 |
| Docker | 端口映射、volume、网络、健康检查正确 |
| 数据库 | 迁移说明、备份、回滚策略明确 |
| 安全 | 密钥不进仓库，生产环境不暴露调试接口 |
| 发布 | 有版本号、发布时间、负责人、变更说明 |
| 回滚 | 上一版本可用，回滚命令已准备 |
| 观察 | 日志、错误率、接口耗时、资源 404 可观察 |

## 学完后你应该能做到

你不需要一开始就掌握 Kubernetes 或复杂云平台，但至少应该能做到：

- 解释浏览器到后端 API 的完整访问链路。
- 区分前端构建时配置和运行时配置。
- 用 Nginx 部署 SPA 并代理接口。
- 用 Docker Compose 组织后端、数据库和缓存。
- 写出基本 CI/CD 流程。
- 发布后验证核心路径。
- 出问题时按链路排查并回滚。

## 下一步学习

下一步可以继续学习 [Nginx 静态部署与代理](/devops/nginx)、[Docker 容器化](/devops/docker)、[CI/CD 自动化发布](/devops/ci-cd) 和 [部署、缓存与 DevOps 问题](/projects/issues-deployment)。
