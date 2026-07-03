# Nginx 速查

## 常用命令

```bash
nginx -t
nginx -s reload
nginx -s stop
```

| 命令 | 用途 |
| --- | --- |
| `nginx -t` | 检查配置是否正确 |
| `nginx -s reload` | 平滑重载配置 |
| `nginx -s stop` | 停止 Nginx |
| `tail -f access.log` | 查看访问日志 |
| `tail -f error.log` | 查看错误日志 |

修改配置后先 `nginx -t`，再 reload。不要在配置未检查时直接重载生产服务。

## 静态站点配置

```nginx
server {
  listen 80;
  server_name example.com;

  root /var/www/app;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }
}
```

`try_files` 用于支持前端 history 路由。没有它，刷新 `/users`、`/orders/1` 这类前端路由会 404。

## 反向代理

```nginx
location /api/ {
  proxy_pass http://backend-service/;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}
```

注意 `proxy_pass` 末尾斜杠会影响路径转发。

```nginx
location /api/ {
  proxy_pass http://backend/;
}
```

请求 `/api/users` 会转发到：

```text
http://backend/users
```

如果写成：

```nginx
location /api/ {
  proxy_pass http://backend;
}
```

请求路径通常会保留 `/api/users`。

## 子路径部署

前端部署到：

```text
https://example.com/admin/
```

Nginx：

```nginx
location /admin/ {
  alias /var/www/admin/;
  try_files $uri $uri/ /admin/index.html;
}
```

前端也要同步：

```ts
export default defineConfig({
  base: '/admin/'
})
```

Router：

```ts
createWebHistory('/admin/')
```

## 缓存策略

`index.html` 不建议强缓存：

```nginx
location = /index.html {
  add_header Cache-Control "no-cache, no-store, must-revalidate";
}
```

带 hash 的静态资源可以长期缓存：

```nginx
location /assets/ {
  add_header Cache-Control "public, max-age=31536000, immutable";
}
```

发布顺序建议：

```text
先上传 assets
↓
再上传 index.html
↓
清理 CDN 的 index.html
↓
验证核心路由
```

## gzip

```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml image/svg+xml;
gzip_min_length 1024;
```

不要压缩已经压缩过的资源，例如大多数图片、zip、字体等。

## 常见问题

| 问题 | 处理 |
| --- | --- |
| 刷新二级路由 404 | 配置 `try_files` fallback |
| 接口 404 | 检查 `location` 和 `proxy_pass` 路径 |
| 用户看到旧页面 | 检查 `index.html` 缓存 |
| 子路径资源 404 | 同步 Vite `base` 和 Router base |
| reload 失败 | 先看 `nginx -t` 和 error log |

## 项目建议

- 每个环境保留独立 Nginx 配置说明。
- 发布前检查 `nginx -t`。
- 静态资源和入口文件使用不同缓存策略。
- 生产代理要转发必要请求头。
- 关键配置改动要写入部署文档和回滚说明。

## 下一步学习

- [Nginx 静态部署与代理](/devops/nginx)
- [构建与部署](/engineering/build-deploy)
- [部署、缓存与 DevOps 问题](/projects/issues-deployment)
