# Vite 工程基础

## 适合谁看

适合准备把 Vue 从“能写页面”推进到“能做真实项目交付”的学习者。

Vite 不只是启动项目的工具。它还影响路径别名、环境变量、代理、构建、静态资源、插件和部署路径。很多 Vue 项目线上问题，最后都能追到工程配置上。

## 你会学到什么

- Vite 在 Vue 项目里负责什么。
- 项目目录如何分层。
- 路径别名怎么配置。
- 开发代理和生产网关有什么区别。
- 常见工程化问题怎么排查。

## Vite 负责什么

| 能力 | 说明 |
| --- | --- |
| 开发服务器 | 启动本地项目，支持热更新 |
| 模块解析 | 处理 `import`、路径别名、依赖预构建 |
| 插件体系 | 接入 Vue、自动导入、压缩、分析等能力 |
| 环境变量 | 区分开发、测试、生产配置 |
| 构建 | 输出生产环境静态资源 |
| 代理 | 本地开发时转发 API 请求 |

## 推荐目录结构

```text
src/
├─ api/             接口请求函数
├─ assets/          图片、字体、静态资源
├─ components/      跨页面复用组件
├─ composables/     组合式逻辑
├─ layouts/         页面布局
├─ router/          路由表和守卫
├─ services/        业务流程编排
├─ stores/          Pinia 状态
├─ styles/          全局样式和变量
├─ types/           全局类型
├─ utils/           工具函数
└─ views/           路由页面
```

目录职责要稳定。不要把接口请求写进 `utils`，也不要把业务流程塞进 `components`。

## 路径别名

Vite 配置：

```ts
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
```

TypeScript 配置也要同步：

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  }
}
```

如果只配 Vite，不配 TypeScript，可能运行正常但编辑器报错。反过来，只配 TypeScript，不配 Vite，编辑器不报错但运行失败。

## 环境变量

Vite 暴露到客户端的变量需要以 `VITE_` 开头：

```ini
VITE_API_BASE_URL=/api
VITE_APP_TITLE=Vue Admin
```

读取：

```ts
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
```

不要把密钥放进前端环境变量。前端变量最终会进入浏览器，用户可以看到。

## 本地代理

开发环境常用代理解决跨域：

```ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

请求写：

```ts
fetch('/api/users')
```

本地会代理到：

```text
http://localhost:8080/users
```

注意：Vite 代理只在本地开发服务器生效。生产环境要由 Nginx、网关或后端服务处理。

## 实际项目常见问题

### 1. 本地接口正常，线上接口 404

**原因**

本地依赖 Vite proxy，线上没有同样的反向代理规则。

**解决方案**

上线时同步配置 Nginx 或网关：

```nginx
location /api/ {
  proxy_pass http://backend-service/;
}
```

同时在部署文档中写清楚：

- 前端请求前缀是什么。
- 后端服务地址是什么。
- 是否需要去掉 `/api` 前缀。

### 2. `@/xxx` 导入报错

**原因**

Vite alias 和 TypeScript paths 没有同时配置，或者文件大小写不一致。

**解决方案**

- 检查 `vite.config.ts`。
- 检查 `tsconfig.json`。
- 检查实际文件名大小写。macOS 有时不敏感，Linux 构建环境可能敏感。

### 3. 修改 `.env` 后不生效

**原因**

环境变量在启动时读取，修改 `.env` 后需要重启开发服务器。

**解决方案**

停止 `npm run dev`，重新启动。

### 4. 构建后静态资源路径错误

**症状**

部署到 `/admin/` 后，页面请求 `/assets/xxx.js` 失败。

**原因**

`base` 没有配置成部署子路径。

**解决方案**

```ts
export default defineConfig({
  base: '/admin/'
})
```

同时 Router history 也要使用相同 base。

### 5. 热更新异常

**症状**

保存文件后页面不更新，或者样式状态很奇怪。

**排查**

1. 看终端是否有编译错误。
2. 浏览器强制刷新。
3. 重启开发服务器。
4. 检查是否有多个旧 dev server 占用端口。

## 工程化最佳实践

- 所有目录职责写入项目 README。
- 路径别名同时配置 Vite 和 TypeScript。
- 环境变量分开发、测试、生产。
- 本地代理和生产代理分别记录。
- 每次修改构建、代理、路径、部署配置，都同步更新文档。
- 构建前跑类型检查、lint 和 build。

## 下一步学习

继续学习 [环境配置](/engineering/env-config) 和 [构建与部署](/engineering/build-deploy)。
