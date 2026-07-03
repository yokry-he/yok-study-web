# 构建与部署

## 适合谁看

适合已经能本地开发 Vue 项目，但准备上线、部署到 Nginx、对象存储、静态托管平台或企业内网服务器的学习者。

构建部署不是最后一步才考虑的事情。路由模式、接口前缀、静态资源路径、缓存策略都会影响线上结果。

如果你想系统学习服务器、Nginx、Docker、CI/CD、发布和回滚，请继续阅读 [DevOps 学习导览](/devops/introduction)。本页保留为前端工程化里的部署入门。

## 构建命令

```bash
npm run build
```

本地预览构建产物：

```bash
npm run preview
```

不要直接双击打开 `dist/index.html`。Vite 官方排错文档也说明，构建产物通过 `file://` 打开时可能出现 CORS 或资源加载问题，应使用 HTTP 服务预览。

## 部署前检查清单

| 检查项 | 为什么重要 |
| --- | --- |
| `npm run build` 通过 | 确认生产构建可用 |
| 非首页路由刷新正常 | 避免 history 模式 404 |
| API 代理正确 | 避免本地正常、线上 404 |
| 静态资源路径正确 | 避免 js/css 加载失败 |
| 环境变量正确 | 避免请求错环境 |
| 缓存策略明确 | 避免用户拿到旧版本 |

## Nginx 基础配置

部署在根路径：

```nginx
server {
  listen 80;
  server_name example.com;

  root /var/www/vue-admin;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }

  location /api/ {
    proxy_pass http://backend-service/;
  }
}
```

`try_files` 的作用是：如果服务器找不到真实文件，就回退到 `index.html`，让 Vue Router 接管路由。

## 部署到子路径

如果项目部署在：

```text
https://example.com/admin/
```

Vite 配置：

```ts
export default defineConfig({
  base: '/admin/'
})
```

Router 配置：

```ts
createRouter({
  history: createWebHistory('/admin/'),
  routes
})
```

Nginx：

```nginx
location /admin/ {
  alias /var/www/vue-admin/;
  try_files $uri $uri/ /admin/index.html;
}
```

## 缓存策略

构建后的 js/css 文件通常带 hash，可以长期缓存。`index.html` 不建议强缓存，否则用户可能一直加载旧入口。

示例：

```nginx
location /assets/ {
  expires 1y;
  add_header Cache-Control "public, immutable";
}

location = /index.html {
  add_header Cache-Control "no-cache";
}
```

## 实际项目常见问题

### 1. 上线后白屏

**排查顺序**

1. 打开浏览器 Console，看 js 是否加载失败。
2. 打开 Network，看资源路径是否 404。
3. 检查 Vite `base` 是否和部署路径一致。
4. 检查服务器是否返回了旧的 `index.html`。

### 2. 本地 preview 正常，Nginx 上刷新 404

**原因**

Nginx 没有配置 history fallback。

**解决方案**

```nginx
try_files $uri $uri/ /index.html;
```

### 3. 用户反馈“我这里还是旧版本”

**原因**

浏览器或 CDN 缓存了旧的 `index.html`。

**解决方案**

- `index.html` 设置 `no-cache`。
- 静态资源使用 hash 文件名。
- 发布后刷新 CDN。
- 必要时做版本检测提示用户刷新。

### 4. 接口跨域

**原因**

生产环境没有使用和本地一样的代理策略。

**解决方案**

优先使用同域反向代理：

```text
前端请求 /api/users
Nginx 转发到 http://backend-service/users
```

如果必须跨域，由后端配置 CORS，前端不要试图“绕过”浏览器安全策略。

### 5. 构建内存不足

**症状**

大型项目构建时 Node 进程内存不足。

**解决方案**

临时提高内存：

```bash
NODE_OPTIONS=--max-old-space-size=4096 npm run build
```

长期应分析包体积、路由懒加载和依赖拆分。

## 回滚策略

上线前应明确：

- 当前版本号。
- 构建产物保存位置。
- 上一个稳定版本保存位置。
- 如何切换 Nginx 指向。
- 回滚后是否需要清理 CDN。

简单目录结构：

```text
/var/www/releases/
├─ 2026-07-01-1200/
├─ 2026-07-01-1800/
└─ current -> 2026-07-01-1800
```

Nginx 指向 `current`，回滚时切换软链接。

## 最佳实践

- 构建产物必须通过 HTTP 服务验证。
- 部署文档写清楚路径、代理、环境变量和回滚方式。
- `index.html` 不强缓存，hash 静态资源可长期缓存。
- history 模式必须配置 fallback。
- 每次上线后验证首页、二级路由、详情页、接口请求和刷新行为。

## 下一步学习

继续学习 [包体积分析](/engineering/bundle-analysis)、[项目交付检查清单](/projects/delivery-checklist)、[部署、缓存与 DevOps 问题](/projects/issues-deployment) 和 [工程化常见问题](/engineering/troubleshooting)。
