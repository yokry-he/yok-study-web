# Kubernetes 入门

## 适合谁看

适合已经理解 Docker，但看到 Kubernetes、Pod、Service、Ingress、Deployment 时容易混乱的人：

- 不知道 Pod 和容器是什么关系。
- 不知道 Service 和 Ingress 分别做什么。
- 不知道为什么要写 readinessProbe。
- 以为 Kubernetes 只是“更复杂的 Docker Compose”。
- 不知道前端、后端和数据库是否都应该放进集群。

Kubernetes 是容器编排系统，负责部署、扩缩容、服务发现、滚动更新和故障恢复。初学阶段不需要一上来掌握全部对象，先理解最常见的运行链路。

## 最小心智模型

```text
Container
↓
Pod
↓
Deployment
↓
Service
↓
Ingress
↓
User
```

含义：

| 对象 | 作用 |
| --- | --- |
| Container | 真正运行应用进程 |
| Pod | Kubernetes 最小调度单元，里面可以有一个或多个容器 |
| Deployment | 管理 Pod 副本、滚动更新和回滚 |
| Service | 给一组 Pod 提供稳定访问入口 |
| Ingress | 把外部 HTTP/HTTPS 流量路由到 Service |

用户通常不会直接访问 Pod，而是通过 Ingress 和 Service。

## Pod

Pod 是 Kubernetes 调度的最小单位。

一个常见后端服务 Pod：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: api
spec:
  containers:
    - name: api
      image: example/api:1.0.0
      ports:
        - containerPort: 3000
```

真实项目很少直接创建裸 Pod，通常用 Deployment 管理。

## Deployment

Deployment 用来声明期望状态。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - name: api
          image: example/api:1.0.0
          ports:
            - containerPort: 3000
```

意思是：我希望一直有 3 个 `app=api` 的 Pod 在运行。

如果某个 Pod 崩溃，Deployment 会尝试拉起新的 Pod。

## Service

Pod 会重建，IP 不稳定。Service 给一组 Pod 提供稳定入口。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api
spec:
  selector:
    app: api
  ports:
    - port: 80
      targetPort: 3000
```

其他服务可以访问 `http://api`，不用关心后面有几个 Pod。

## Ingress

Ingress 负责把外部 HTTP/HTTPS 路由到集群内部服务。

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api
                port:
                  number: 80
```

Ingress 需要集群里有 Ingress Controller。只写 Ingress 配置，不安装控制器，外部流量不会自动进来。

## ConfigMap 和 Secret

应用配置不要写死在镜像里。

| 对象 | 用途 |
| --- | --- |
| ConfigMap | 普通配置，例如环境名、公开 API 地址 |
| Secret | 敏感配置，例如数据库密码、Token |

注意：Kubernetes Secret 不是“天然绝对安全”。还要结合访问权限、加密存储和密钥管理。

## 健康检查

Kubernetes 常用三类探针：

| 探针 | 作用 |
| --- | --- |
| startupProbe | 应用启动慢时，判断是否已经启动完成 |
| readinessProbe | 判断是否可以接收流量 |
| livenessProbe | 判断是否需要重启容器 |

示例：

```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 3000
  initialDelaySeconds: 5
  periodSeconds: 10

livenessProbe:
  httpGet:
    path: /health
    port: 3000
  initialDelaySeconds: 30
  periodSeconds: 20
```

不要把 readiness 和 liveness 混为一谈。

- readiness 失败：先不要给它流量。
- liveness 失败：认为进程坏了，需要重启。

如果 liveness 检查依赖数据库，数据库短暂抖动可能导致应用被反复重启。通常 liveness 应该检查进程是否活着，readiness 才检查关键依赖是否可用。

## 前端项目怎么放 Kubernetes

常见方案：

1. 前端构建成静态文件，放 Nginx 镜像里。
2. 用 Deployment 运行 Nginx。
3. 用 Service 暴露内部端口。
4. 用 Ingress 绑定域名和 HTTPS。

但如果只是静态官网或文档站，云静态托管、对象存储 + CDN、Vercel、Netlify、Cloudflare Pages 可能更简单。不是所有项目都必须上 Kubernetes。

## 数据库要不要放 Kubernetes

学习环境可以放，但生产环境要谨慎。

数据库需要：

- 持久化存储。
- 备份恢复。
- 高可用。
- 监控告警。
- 权限管理。
- 升级策略。

很多团队会把无状态应用放 Kubernetes，把数据库交给云数据库或专业数据库平台。

## 实际项目问题

### 1. Pod 一直重启

**常见原因**

- 镜像启动命令错误。
- 环境变量缺失。
- 端口配置错误。
- livenessProbe 过于激进。
- 应用启动时间比探针等待时间长。

**解决方案**

- 查看 Pod 事件。
- 查看容器日志。
- 临时放宽探针时间。
- 区分启动失败和健康检查失败。

### 2. Service 访问不到 Pod

**常见原因**

- Service selector 和 Pod label 不匹配。
- targetPort 写错。
- Pod 没有 ready。
- 应用只监听 127.0.0.1。

**解决方案**

- 检查 label。
- 检查 endpoints。
- 确认容器监听 `0.0.0.0`。
- 检查 readinessProbe。

### 3. Ingress 配了但外部访问 404

**常见原因**

- Ingress Controller 没安装。
- host 不匹配。
- pathType 或 path 配置错误。
- Service 名称或端口错误。

**解决方案**

- 先确认 Ingress Controller 工作正常。
- 再确认 Ingress 规则。
- 最后检查 Service 和 Pod。

### 4. 发布后新版本不接流量

**常见原因**

readinessProbe 失败。

**解决方案**

- 查看 readiness 失败原因。
- 确认 `/health` 不依赖过慢初始化。
- 把启动期检查放到 startupProbe。
- 发布后观察 ready 副本数。

## 最佳实践

- 先学 Docker 和 Compose，再学 Kubernetes。
- 无状态服务更适合先上 Kubernetes。
- 用 Deployment 管理副本和滚动更新。
- Service 提供稳定内部入口。
- Ingress 处理外部 HTTP/HTTPS 路由。
- readiness 和 liveness 要分清职责。
- Secret 还需要配合权限和密钥治理。
- 生产数据库是否进集群要谨慎评估。

## 参考资料

- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Kubernetes Service](https://kubernetes.io/docs/concepts/services-networking/service/)
- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Configure Liveness, Readiness and Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

## 下一步学习

继续学习 [云服务与对象存储部署](/devops/cloud-deployment)，理解项目从服务器部署走向云上托管时的选择。
