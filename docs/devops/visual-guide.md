# 图解 DevOps 核心概念

## 这个页面解决什么

DevOps 容易被误解成“会点 Linux、Docker、Nginx 命令”。真实项目里更重要的是理解请求如何到达服务、镜像如何构建、发布如何灰度、出了问题如何回滚。

## 适合谁看

适合已经能完成项目开发，但对服务器、Nginx、Docker、CI/CD、灰度发布、监控和线上排错缺少全链路理解的人。

## 一张图理解线上请求链路

```mermaid
flowchart LR
  A["用户浏览器"] --> B["DNS"]
  B --> C["CDN"]
  C --> D["负载均衡"]
  D --> E["Nginx / Ingress"]
  E --> F["前端静态资源"]
  E --> G["后端 API 服务"]
  G --> H[("数据库")]
  G --> I[("Redis")]
  G --> J["日志 / 指标 / 链路追踪"]
```

排查线上问题时，要先判断请求卡在哪一层：

- 域名解析。
- CDN 缓存。
- Nginx 代理。
- 前端静态资源。
- 后端接口。
- 数据库或缓存。

## 一张图理解 Nginx 反向代理

```mermaid
sequenceDiagram
  participant B as Browser
  participant N as Nginx
  participant FE as Static Files
  participant API as Backend API

  B->>N: GET /
  N->>FE: 读取 index.html
  FE-->>N: HTML/CSS/JS
  N-->>B: 静态资源

  B->>N: GET /api/users
  N->>API: proxy_pass
  API-->>N: JSON
  N-->>B: JSON
```

常见问题：

- `base` 路径不对导致资源 404。
- API 代理路径多一段或少一段。
- 刷新页面 404，没有回退到 `index.html`。
- 缓存策略导致旧前端资源没更新。

## 一张图理解 Docker 镜像分层

```mermaid
flowchart TD
  A["FROM 基础镜像"] --> B["安装依赖"]
  B --> C["复制 package / go.mod / pom.xml"]
  C --> D["下载依赖"]
  D --> E["复制源码"]
  E --> F["构建产物"]
  F --> G["运行镜像"]
```

镜像优化的关键：

- 依赖层和源码层分开，提升缓存命中。
- 多阶段构建，只把产物放进运行镜像。
- 不把 `.env`、密钥、无关文件打进镜像。
- 镜像标签要能追踪 commit。

## 一张图理解 CI/CD

```mermaid
flowchart LR
  A["提交代码"] --> B["安装依赖"]
  B --> C["Lint / Docs Check"]
  C --> D["单元测试"]
  D --> E["构建"]
  E --> F["构建镜像"]
  F --> G["推送制品库"]
  G --> H["部署测试环境"]
  H --> I["验收"]
  I --> J["灰度生产"]
  J --> K["监控"]
```

CI/CD 的目标不是“自动点发布”，而是让每一步都可重复、可追踪、可回滚。

## 一张图理解蓝绿发布和灰度发布

```mermaid
flowchart TD
  A["用户流量"] --> B{"流量分配"}
  B -- "90%" --> C["旧版本 v1"]
  B -- "10%" --> D["新版本 v2"]
  D --> E{"指标正常吗"}
  E -- "是" --> F["逐步增加 v2 流量"]
  E -- "否" --> G["回滚到 v1"]
```

灰度时至少观察：

- 错误率。
- 接口延迟。
- 业务转化。
- 日志异常。
- 数据库慢查询。
- 资源使用率。

## 一张图理解可观测性

```mermaid
flowchart TD
  A["一次请求"] --> B["日志 Logs<br/>发生了什么"]
  A --> C["指标 Metrics<br/>数量和趋势"]
  A --> D["链路 Traces<br/>经过哪些服务"]
  B --> E["定位单次错误"]
  C --> F["发现异常趋势"]
  D --> G["定位慢在哪一段"]
```

日志、指标、链路追踪解决的问题不同：

- 日志适合看单次请求细节。
- 指标适合看趋势和告警。
- 链路追踪适合看跨服务耗时。

## 一张图理解上线故障排查

```mermaid
flowchart TD
  A["上线后异常"] --> B{"是否影响用户"}
  B -- "影响严重" --> C["先回滚或降级"]
  B -- "影响可控" --> D["保留现场继续定位"]
  C --> E["收集版本、日志、指标"]
  D --> E
  E --> F{"异常层次"}
  F --> G["前端资源 / CDN"]
  F --> H["Nginx / 网关"]
  F --> I["后端服务"]
  F --> J["数据库 / Redis"]
  F --> K["第三方服务"]
```

上线排查原则：

1. 先止血。
2. 再保留证据。
3. 再定位根因。
4. 最后补测试、监控和复盘。

## 下一步学习

继续学习 [Linux 与 Shell 基础](/devops/linux-shell)，或进入 [Nginx 静态部署与代理](/devops/nginx)。
