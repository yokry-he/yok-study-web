# CI/CD 自动化发布

## 适合谁看

适合想把“拉代码、装依赖、构建、测试、部署”从手工操作变成自动流程的人。

CI/CD 的价值不是炫技，而是减少人为错误：

- 避免忘记跑测试。
- 避免本地环境和发布环境不一致。
- 避免手工复制文件出错。
- 让每次发布都有记录。
- 失败时能快速定位到哪一步。

## CI 和 CD

| 缩写 | 含义 | 常见任务 |
| --- | --- | --- |
| CI | Continuous Integration，持续集成 | 安装依赖、Lint、测试、构建 |
| CD | Continuous Delivery/Deployment，持续交付或部署 | 打包、上传、部署、通知、回滚 |

学习阶段先把 CI 做稳，再接 CD。

## GitHub Actions 基础结构

Workflow 文件放在：

```text
.github/workflows/ci.yml
```

示例：

```yaml
name: CI

on:
  pull_request:
  push:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build
```

工作流由触发条件、job 和 step 组成。

## 前端项目推荐 CI

```yaml
name: Frontend CI

on:
  pull_request:
  push:
    branches:
      - main

jobs:
  quality:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm

      - run: npm ci
      - run: npm run lint
      - run: npm run test
      - run: npm run build
```

如果项目暂时没有测试，也应该至少保留构建检查。文档站项目可以用：

```yaml
- run: npm run docs:build
```

## 缓存依赖

`actions/setup-node` 支持 npm、yarn、pnpm 缓存。缓存能减少依赖安装时间，但不要缓存 `node_modules` 后跳过锁文件安装。

推荐：

```yaml
with:
  node-version: 22
  cache: npm
```

然后仍然执行：

```bash
npm ci
```

`npm ci` 会严格按 lock 文件安装，更适合 CI。

## 部署静态站

常见方式：

| 方式 | 说明 |
| --- | --- |
| 上传到服务器 | rsync、scp、对象存储 |
| 部署到静态托管 | Vercel、Netlify、Cloudflare Pages |
| 构建 Docker 镜像 | 推送镜像仓库，服务器拉取 |

服务器部署示例思路：

```text
CI 构建 dist
↓
上传到 /var/www/releases/<version>
↓
切换 current 软链接
↓
reload Nginx
↓
访问健康检查 URL
```

不要直接覆盖线上目录。覆盖过程中用户可能访问到半更新状态。

## 环境变量和密钥

CI 中的密钥应放在平台 Secrets，不要写进仓库。

常见密钥：

- 服务器 SSH 私钥。
- Docker Registry token。
- 云平台访问密钥。
- 部署 webhook。

前端公开变量要谨慎。以 `VITE_` 开头的变量会进入前端产物，不能放私密信息。

## 实际项目问题

### 问题：本地 build 通过，CI build 失败

**常见原因**

- 本地没有清理依赖，CI 是干净环境。
- Node 版本不同。
- lock 文件没提交。
- 大小写路径在 macOS 不敏感，但 Linux 敏感。
- 环境变量缺失。

**解决方案**

- 固定 Node 版本。
- 使用 `npm ci`。
- 提交 lock 文件。
- 在本地用干净安装复现。
- 给必要环境变量提供示例和校验。

### 问题：CI 很慢

**排查**

- 依赖安装耗时。
- 测试耗时。
- 构建耗时。
- Docker build 上下文太大。

**解决方案**

- 开启包管理器缓存。
- 拆分 job。
- 优化 Dockerfile 层缓存。
- 使用 `.dockerignore`。

### 问题：发布成功但线上没变化

**常见原因**

- 上传到了错误目录。
- Nginx 指向的不是新目录。
- CDN 或浏览器缓存旧 HTML。
- 容器没重启或拉的不是新镜像。

**解决方案**

发布后必须做自动验证：

```bash
curl -I https://example.com/
curl -I https://example.com/assets/index.xxxxx.js
```

还可以写版本文件：

```text
/version.json
```

发布后请求版本号确认是否生效。

## 最佳实践

- CI 至少包括安装依赖和生产构建。
- 固定 Node 版本，使用 lock 文件。
- 密钥放 Secrets，不写进仓库。
- 部署不要直接覆盖线上目录。
- 每次发布生成版本号和构建记录。
- CD 后必须做健康检查和关键页面验证。

## 参考资料

- [GitHub Actions documentation](https://docs.github.com/actions)
- [GitHub Actions workflow syntax](https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions)
- [GitHub Actions dependency caching](https://docs.github.com/en/actions/reference/workflows-and-actions/dependency-caching)
- [GitHub Actions: Building and testing Node.js](https://docs.github.com/actions/guides/building-and-testing-nodejs)

## 下一步学习

继续学习 [发布、回滚与环境治理](/devops/deployment-strategy)。
