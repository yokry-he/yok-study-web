# DevOps 学习导览

## 适合谁看

适合已经能完成前端或 Node.js 项目开发，但对上线、服务器、Nginx、Docker、CI/CD、回滚和线上排错还不熟的人。

DevOps 不是“运维专属技能”。对前端和全栈开发者来说，至少要能回答这些问题：

- 项目构建后应该放到哪里。
- Nginx 为什么能解决 history 刷新 404 和接口代理。
- Docker 镜像和容器分别是什么。
- CI/CD 为什么能减少手工发布错误。
- 线上出问题后如何快速回滚。
- 服务器磁盘、端口、日志、进程怎么排查。

## 学习目标

学完第一版 DevOps 模块，你应该能独立完成一个中小型 Web 项目的基础上线流程：

```text
图解 DevOps 核心概念
↓
本地构建
↓
服务器准备
↓
Nginx 静态部署
↓
接口反向代理
↓
Docker 容器化
↓
CI/CD 自动构建
↓
项目上线全流程实践
↓
发布验证
↓
问题回滚
```

## DevOps 学习路线

| 阶段 | 重点 | 你要能做到 |
| --- | --- | --- |
| Linux 与 Shell | 文件、进程、端口、日志、权限 | 登录服务器后不慌，能查问题 |
| Nginx | 静态资源、history fallback、反向代理、缓存 | 能部署前端并代理接口 |
| Docker | 镜像、容器、Dockerfile、Compose | 能把应用稳定打包运行 |
| CI/CD | workflow、构建、测试、发布 | 能减少手工发布步骤 |
| 发布治理 | 环境、版本、回滚、健康检查 | 能控制上线风险 |
| 排错 | 日志、网络、缓存、端口、权限 | 能定位常见线上故障 |

## 你需要先理解的几个词

### 服务器

服务器可以是云主机、容器平台、企业内网机器或静态托管平台。学习阶段不需要一开始就追求 Kubernetes，先理解一台 Linux 服务器如何运行项目更重要。

### 进程

运行中的程序就是进程。Node 服务、Nginx、数据库、Docker 守护进程都以进程形式运行。线上排错常常要看进程是否存在、端口是否监听、日志是否报错。

### 端口

端口是服务对外提供访问的位置。例如：

```text
Nginx: 80 / 443
Node API: 3000
Vite dev server: 5173 / 6173
```

访问失败时，要先确认端口是否监听，再看防火墙、反向代理和服务日志。

### 镜像和容器

镜像像一个打包好的应用模板，容器是镜像运行起来后的实例。

```text
Dockerfile -> image -> container
```

不要把容器理解成虚拟机。容器更轻量，但也更依赖正确的环境变量、网络、卷和启动命令。

## 模块章节

| 章节 | 解决的问题 |
| --- | --- |
| [图解 DevOps 核心概念](/devops/visual-guide) | 用图理解线上请求链路、Nginx、Docker 分层、CI/CD、灰度发布和可观测性 |
| [Linux 与 Shell 基础](/devops/linux-shell) | 服务器基础命令、日志、进程、端口、权限 |
| [Nginx 静态部署与代理](/devops/nginx) | 前端部署、history fallback、反向代理、缓存 |
| [Docker 容器化](/devops/docker) | Dockerfile、多阶段构建、Compose、镜像体积 |
| [CI/CD 自动化发布](/devops/ci-cd) | GitHub Actions、构建测试、部署流水线 |
| [项目上线全流程实践](/devops/project-deployment-practice) | 从构建、Nginx、Docker Compose、缓存策略、发布验证到回滚的完整上线案例 |
| [发布、回滚与环境治理](/devops/deployment-strategy) | 版本、环境变量、回滚、上线检查 |
| [可观测性](/devops/observability) | 日志、指标、链路追踪、告警和发布观察窗口 |
| [Kubernetes 入门](/devops/kubernetes-basics) | Pod、Deployment、Service、Ingress 和健康检查 |
| [云服务与对象存储部署](/devops/cloud-deployment) | 静态托管、对象存储、CDN、云数据库和成本治理 |
| [常见问题](/devops/troubleshooting) | 真实项目里高频部署问题和解决方案 |

## 初学者容易踩的坑

### 只会本地运行，不会生产运行

本地 `npm run dev` 是开发模式，生产环境通常运行构建产物或 Node 服务。开发服务器自带的能力不能当作生产部署方案。

### 只会复制配置，不理解路径

部署前端时，最容易错的是路径：

- Vite `base`。
- Router base。
- Nginx `root` 和 `alias`。
- `try_files` fallback。
- `/api` 代理路径是否保留前缀。

### 手工发布没有回滚

只会把文件覆盖到服务器是不够的。真正的发布方案必须知道：

- 当前版本是哪一个。
- 上一个稳定版本在哪里。
- 如何切回去。
- 回滚后是否需要刷新缓存。

## 建议实践项目

可以用一个 Vue Admin 或 Node API 项目练习：

1. 本地构建。
2. Nginx 部署静态资源。
3. 配置 `/api` 反向代理。
4. 用 Dockerfile 打包前端或 Node 服务。
5. 用 Docker Compose 启动前端、后端和数据库。
6. 用 GitHub Actions 跑构建。
7. 写发布检查清单和回滚步骤。

## 下一步学习

继续学习 [Linux 与 Shell 基础](/devops/linux-shell)。
