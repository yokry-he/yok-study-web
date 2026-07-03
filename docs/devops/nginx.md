# Nginx 静态部署与代理

## 适合谁看

适合要把 Vue、React、VitePress 或普通静态站部署到服务器的人。

Nginx 在 Web 项目里常做三件事：

- 直接返回静态资源。
- 把前端路由回退到 `index.html`。
- 把 `/api` 请求反向代理到后端服务。

## 最小前端部署配置

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

核心点：

| 配置 | 作用 |
| --- | --- |
| `root` | 静态资源根目录 |
| `index` | 默认入口文件 |
| `try_files` | 找不到真实文件时回退到前端入口 |

`try_files $uri $uri/ /index.html;` 是 history 路由刷新不 404 的关键。

## root 和 alias 的区别

`root` 会把请求路径拼到根目录后面：

```nginx
location /admin/ {
  root /var/www;
}
```

请求 `/admin/index.html` 时，会找：

```text
/var/www/admin/index.html
```

`alias` 会用指定目录替换匹配路径：

```nginx
location /admin/ {
  alias /var/www/vue-admin/;
}
```

请求 `/admin/index.html` 时，会找：

```text
/var/www/vue-admin/index.html
```

部署子路径时，`alias` 末尾的 `/` 很重要。

## 反向代理 API

常见配置：

```nginx
location /api/ {
  proxy_pass http://backend-service/;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
}
```

注意 `proxy_pass` 结尾是否带 `/` 会影响路径转发。

示例：

```nginx
location /api/ {
  proxy_pass http://backend/;
}
```

`/api/users` 转发成：

```text
http://backend/users
```

如果是：

```nginx
location /api/ {
  proxy_pass http://backend;
}
```

`/api/users` 通常会保留 `/api/users`。

路径是否保留要和后端接口前缀一致。

## 缓存配置

前端构建产物通常分为入口 HTML 和 hash 静态资源：

```text
index.html
assets/index.8f3a1c.js
assets/index.71d2c.css
```

推荐：

```nginx
location = /index.html {
  add_header Cache-Control "no-cache";
}

location /assets/ {
  add_header Cache-Control "public, max-age=31536000, immutable";
}
```

不要让 `index.html` 长期强缓存，否则用户可能一直拿到旧入口。

## gzip 压缩

常见配置：

```nginx
gzip on;
gzip_types text/plain text/css application/javascript application/json image/svg+xml;
gzip_min_length 1024;
```

压缩能减少传输体积，但图片、视频这类已经压缩的资源收益较小。

## 配置检查和重载

修改 Nginx 配置后先检查：

```bash
nginx -t
```

通过后重载：

```bash
systemctl reload nginx
```

不要配置没检查就直接重启。配置错误可能导致服务不可用。

## 实际项目问题

### 问题：刷新二级路由 404

**原因**

缺少 history fallback。

**解决方案**

```nginx
location / {
  try_files $uri $uri/ /index.html;
}
```

### 问题：接口变成前端 HTML

**现象**

Network 里 `/api/users` 返回的 `Content-Type` 是 `text/html`，内容是前端页面。

**原因**

`/api` 没有单独代理，被 `location /` fallback 到了 `index.html`。

**解决方案**

把 `/api/` location 放清楚，并代理到后端：

```nginx
location /api/ {
  proxy_pass http://backend-service/;
}
```

### 问题：部署到 `/admin/` 后资源 404

**原因**

Vite `base`、路由 base、Nginx 子路径配置不一致。

**解决方案**

三处统一：

```ts
// vite.config.ts
export default defineConfig({
  base: '/admin/'
})
```

```ts
createWebHistory('/admin/')
```

```nginx
location /admin/ {
  alias /var/www/vue-admin/;
  try_files $uri $uri/ /admin/index.html;
}
```

## 最佳实践

- 前端路由和 API 代理分开配置。
- `index.html` 不长缓存，hash 静态资源长缓存。
- 每次改配置先 `nginx -t`。
- 部署子路径时统一 Vite base、Router base 和 Nginx 路径。
- 遇到线上问题先看 Nginx access log 和 error log。

## 参考资料

- [NGINX: Serve Static Content](https://docs.nginx.com/nginx/admin-guide/web-server/serving-static-content/)
- [nginx.org: ngx_http_core_module](https://nginx.org/en/docs/http/ngx_http_core_module.html)

## 下一步学习

继续学习 [Docker 容器化](/devops/docker)。
