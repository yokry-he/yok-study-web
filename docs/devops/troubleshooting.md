# DevOps 常见问题

## 使用方式

线上问题不要只看前端控制台。DevOps 排错要从链路看：

```text
浏览器
↓
CDN / DNS
↓
Nginx
↓
后端服务
↓
数据库 / 缓存
↓
服务器资源
```

先确认故障发生在哪一层，再改配置。

## 1. Nginx 配置改了但不生效

**排查**

```bash
nginx -t
systemctl reload nginx
nginx -T | grep server_name -n
```

**常见原因**

- 改了错误配置文件。
- 配置没有被 include。
- `nginx -t` 没通过。
- reload 失败但没注意。
- 多个 server block 匹配到了另一个。

**解决方案**

用 `nginx -T` 输出最终加载配置，确认真实生效内容。

## 2. 502 Bad Gateway

**含义**

Nginx 作为代理时，无法从上游服务拿到有效响应。

**排查**

1. 后端进程是否运行。
2. 后端端口是否监听。
3. `proxy_pass` 地址是否正确。
4. 容器网络是否能互通。
5. 后端是否启动慢或崩溃。

命令：

```bash
ss -lntp
curl http://127.0.0.1:3000/health
tail -n 100 /var/log/nginx/error.log
```

## 3. 504 Gateway Timeout

**含义**

代理等后端响应超时。

**常见原因**

- 后端接口慢。
- 数据库查询慢。
- 后端服务卡死。
- Nginx 超时时间太短。

**解决方案**

不要只把超时时间调大。先查后端接口耗时和数据库慢查询。确实是长任务时，改为异步任务或轮询结果。

## 4. Docker 容器一直重启

**排查**

```bash
docker ps -a
docker logs <container>
docker inspect <container>
```

**常见原因**

- 启动命令错误。
- 环境变量缺失。
- 端口被占用。
- 程序启动后立即退出。
- 健康检查失败。

**解决方案**

先看容器日志，不要直接反复 `docker compose up -d`。

## 5. Docker build 很慢

**常见原因**

- Dockerfile 层顺序不合理。
- 每次都复制整个项目后再安装依赖。
- 没有 `.dockerignore`。
- 镜像源慢。

**优化**

先复制 lock 文件安装依赖：

```dockerfile
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build
```

这样源码变动不会总是破坏依赖层缓存。

## 6. CI 里 npm install 失败

**常见原因**

- lock 文件和 package.json 不一致。
- Node 版本不同。
- 私有包 token 缺失。
- 网络或 registry 不稳定。

**解决方案**

- CI 使用 `npm ci`。
- 固定 Node 版本。
- 私有 registry token 放 Secrets。
- 依赖变更时同步提交 lock 文件。

## 7. 静态资源 404

**排查**

1. 构建产物里文件是否存在。
2. Nginx root/alias 是否指向正确目录。
3. Vite `base` 是否正确。
4. CDN 是否缓存了旧 HTML。
5. 旧版本 hash 文件是否被删除。

**解决方案**

静态资源路径问题要同时看构建配置和部署路径，不要只改其中一个。

## 8. 生产环境接口请求到测试环境

**原因**

前端构建时使用了错误环境变量，或者 CI/CD secrets 指向错环境。

**排查**

- 查看构建日志中的环境标识。
- 查看最终产物里的 API base。
- 访问 `/version.json`。
- 检查 CI job 和环境绑定。

**预防**

每个环境单独配置 secrets，发布前输出非敏感环境名。

## 9. 服务器磁盘满了

**症状**

- Docker 拉镜像失败。
- 日志写不进去。
- 数据库异常。
- 构建失败。

**排查**

```bash
df -h
du -sh /var/lib/docker/*
du -sh /var/log/*
```

**处理**

- 清理旧日志。
- 清理未使用镜像。
- 保留必要发布版本，删除过旧 releases。
- 给数据库和日志设置保留策略。

不要盲目清理数据库目录和 Docker volume。

## 10. 回滚后仍然是新版本

**常见原因**

- CDN 缓存没刷新。
- 浏览器缓存了 `index.html`。
- Nginx 没 reload。
- current 软链接没切成功。
- 容器仍然运行旧实例。

**排查**

```bash
readlink -f /var/www/my-app/current
curl -I https://example.com/
curl https://example.com/version.json
```

**解决方案**

回滚流程必须包含版本验证和缓存处理。

## 排查清单

| 问题 | 优先看 |
| --- | --- |
| 404 | 路径、root/alias、fallback、资源是否存在 |
| 502 | 后端进程、端口、proxy_pass、容器网络 |
| 504 | 后端耗时、数据库慢查询、超时配置 |
| 白屏 | Console、Network、静态资源、缓存 |
| CI 失败 | Node 版本、lock 文件、secrets、构建脚本 |
| Docker 异常 | container logs、环境变量、端口、volume |

## 最佳实践

- 每个故障先定位层级。
- 修改配置前备份，修改后验证。
- 自动化发布要有失败停止机制。
- 日志、版本号、健康检查是线上排错基础设施。
- 回滚流程要定期演练，不要等事故发生才写。
