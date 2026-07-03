# Node.js 速查

## 常用命令

```bash
node -v
npm -v
npm init -y
npm install
npm run dev
npm run build
```

| 命令 | 用途 |
| --- | --- |
| `node -v` | 查看 Node 版本 |
| `npm install` | 安装依赖 |
| `npm run <script>` | 运行 `package.json` 脚本 |
| `npx <command>` | 临时执行包命令 |
| `corepack enable` | 启用 pnpm/yarn 版本管理 |

## package.json 常见字段

```json
{
  "name": "my-app",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "start": "node server.js"
  },
  "dependencies": {},
  "devDependencies": {},
  "engines": {
    "node": ">=20"
  },
  "packageManager": "pnpm@9.15.0"
}
```

| 字段 | 用途 |
| --- | --- |
| `scripts` | 项目命令入口 |
| `dependencies` | 运行时依赖 |
| `devDependencies` | 开发和构建依赖 |
| `type` | 模块类型，`module` 表示 ESM |
| `engines` | 声明 Node 版本要求 |
| `packageManager` | 固定包管理器 |

## CommonJS 和 ESM

CommonJS：

```js
const path = require('node:path')
module.exports = {}
```

ESM：

```js
import path from 'node:path'
export default {}
```

如果 `package.json` 中有：

```json
{
  "type": "module"
}
```

那么 `.js` 默认按 ESM 解析。需要 CommonJS 时可以使用 `.cjs`。

## 常用内置模块

| 模块 | 用途 |
| --- | --- |
| `node:path` | 处理文件路径 |
| `node:fs` | 读写文件 |
| `node:url` | 处理 URL 和 file URL |
| `node:http` | 创建 HTTP 服务 |
| `node:crypto` | 哈希、随机值、加密相关 |
| `node:process` | 进程参数、环境变量 |

路径示例：

```js
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
```

## 环境变量

读取环境变量：

```js
const port = process.env.PORT || 3000
```

启动时传入：

```bash
PORT=3000 node server.js
```

项目建议：

- 环境变量集中读取。
- 必填变量启动时校验。
- 不把密钥写进代码仓库。
- 前端变量和后端变量分开管理。

## 简单 HTTP 服务

```js
import http from 'node:http'

const server = http.createServer((req, res) => {
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' })
    res.end(JSON.stringify({ ok: true }))
    return
  }

  res.writeHead(404)
  res.end('Not Found')
})

server.listen(3000, '0.0.0.0')
```

容器或服务器环境中，服务通常要监听 `0.0.0.0`，不要只监听 `127.0.0.1`。

## 常见问题

| 问题 | 处理 |
| --- | --- |
| Node 版本不一致 | 用 `engines`、`.nvmrc` 或 Volta 固定 |
| ESM 中没有 `__dirname` | 用 `fileURLToPath(import.meta.url)` |
| 端口被占用 | 换端口或结束旧进程 |
| 容器内服务访问不到 | 监听 `0.0.0.0` 并检查端口映射 |
| 本地能跑 CI 失败 | 检查 lockfile、Node 版本和环境变量 |

## 项目建议

- 所有常用命令写进 `scripts`。
- README 写清 Node 版本和包管理器。
- 后端服务必须提供 `/health`。
- 日志中带 requestId，方便联调排查。
- 不在业务代码里散落读取环境变量。

## 下一步学习

- [Node.js 学习导览](/node/introduction)
- [运行时与事件循环](/node/runtime-event-loop)
- [包管理与模块化](/node/package-modules)
- [依赖管理](/engineering/package-management)
