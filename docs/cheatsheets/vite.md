# Vite 速查

## 常用命令

```bash
npm run dev
npm run build
npm run preview
```

## Vue 插件

```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()]
})
```

## 路径别名

```ts
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
```

TypeScript 也要配置：

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

## 环境变量

```ini
VITE_API_BASE_URL=/api
VITE_APP_TITLE=Vue Admin
```

读取：

```ts
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
const mode = import.meta.env.MODE
```

## 本地代理

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

## 部署子路径

```ts
export default defineConfig({
  base: '/admin/'
})
```

Router 也要同步：

```ts
createWebHistory('/admin/')
```

## 常见问题

| 问题 | 解决方案 |
| --- | --- |
| 修改 `.env` 不生效 | 重启 dev server |
| `@` 路径运行报错 | 检查 Vite alias |
| `@` 路径编辑器报错 | 检查 tsconfig paths |
| 本地接口正常线上 404 | 生产环境配置 Nginx 或网关 |
| 构建后资源 404 | 检查 `base` |
| 不能直接打开 dist | 用 `npm run preview` 或 HTTP 服务 |
