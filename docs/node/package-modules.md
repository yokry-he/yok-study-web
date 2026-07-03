# 包管理与模块化

## 适合谁看

适合不熟悉 `package.json`、npm scripts、依赖版本、ESM/CommonJS 的学习者。

## package.json

`package.json` 是 Node 项目的说明书。

```json
{
  "name": "node-api",
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "tsx watch src/main.ts",
    "build": "tsc",
    "start": "node dist/main.js"
  },
  "dependencies": {
    "fastify": "^5.0.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0"
  }
}
```

## dependencies 和 devDependencies

| 字段 | 用途 |
| --- | --- |
| dependencies | 运行时需要 |
| devDependencies | 开发和构建需要 |

例如 API 服务运行需要 Fastify，所以放 dependencies。TypeScript 编译工具通常放 devDependencies。

## npm scripts

常见脚本：

```json
{
  "scripts": {
    "dev": "tsx watch src/main.ts",
    "build": "tsc",
    "start": "node dist/main.js",
    "lint": "eslint ."
  }
}
```

项目 README 必须写清楚这些脚本用途。

## ESM 和 CommonJS

ESM：

```ts
import { readFile } from 'node:fs/promises'
export function loadConfig() {}
```

CommonJS：

```js
const fs = require('node:fs')
module.exports = {}
```

新项目建议统一使用 ESM，不要混用。

## 实际项目常见问题

### 1. Cannot use import statement outside a module

**原因**

项目没有配置 ESM，或文件扩展名、运行命令不匹配。

**解决方案**

在 `package.json` 中设置：

```json
{
  "type": "module"
}
```

并统一使用 ESM。

### 2. 本地能跑，服务器缺包

**排查**

- 是否提交了 package-lock。
- 是否在服务器执行了 `npm ci`。
- 依赖是否放错 dependencies/devDependencies。

### 3. 版本不一致导致问题

**建议**

使用 lockfile，部署时使用 `npm ci`，并在 README 写清楚 Node 版本要求。

## 最佳实践

- 新项目统一模块规范。
- 提交 lockfile。
- 部署使用 `npm ci`。
- scripts 命名清楚。
- README 写明 Node 版本和启动命令。

## 下一步

继续学习 [HTTP API 开发](/node/http-api)。
