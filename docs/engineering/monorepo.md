# Monorepo 项目组织

## 适合谁看

适合已经做过单个前端项目，并准备组织多个应用、组件库、工具库或业务包的学习者。

Monorepo 不是大公司专属。只要一个仓库里同时存在多个相互依赖的包，就需要考虑如何组织代码、安装依赖、运行脚本、发布版本和控制边界。

## 你会学到什么

- Monorepo 解决什么问题。
- 什么时候该用，什么时候不该用。
- apps 和 packages 如何划分。
- pnpm workspace 怎么配置。
- 包之间如何依赖。
- 多包项目常见问题怎么处理。

## Monorepo 是什么

Monorepo 是“一个仓库管理多个项目或包”的方式。

示例：

```text
repo/
├─ apps/
│  ├─ admin/          后台管理系统
│  └─ docs/           文档站
├─ packages/
│  ├─ ui/             组件库
│  ├─ utils/          通用工具
│  └─ config/         共享工程配置
├─ package.json
├─ pnpm-workspace.yaml
└─ tsconfig.base.json
```

它和普通多仓库的区别是：多个应用和包共享同一个代码仓库、依赖安装和提交历史。

## 什么时候适合用 Monorepo

适合：

- 一个后台系统和一个文档站共享组件库。
- 多个业务应用共享请求库、工具库和类型。
- 项目里有组件库、主题包、图标包。
- 希望一次提交同时修改应用和共享包。
- 希望统一 lint、测试、构建和发布流程。

不适合：

- 只有一个小应用。
- 团队还没有稳定目录和工程规范。
- 多个项目完全无关。
- 权限隔离要求非常强。
- 团队还没有能力维护构建和依赖边界。

不要为了“显得高级”使用 Monorepo。它解决的是复用和协作问题，也会带来边界治理成本。

## 推荐目录结构

基础结构：

```text
apps/
├─ admin/
└─ docs/

packages/
├─ ui/
├─ utils/
├─ request/
└─ config/
```

目录职责：

| 目录 | 作用 |
| --- | --- |
| `apps/` | 可独立运行和部署的应用 |
| `packages/` | 被应用复用的库 |
| `packages/ui` | 组件库 |
| `packages/utils` | 无业务状态的工具函数 |
| `packages/request` | 请求封装和错误处理 |
| `packages/config` | ESLint、TS、Vite 等共享配置 |

业务专属代码不要放进通用包。通用包一旦混入业务逻辑，会让所有应用都被迫理解同一套业务假设。

## pnpm workspace 配置

创建 `pnpm-workspace.yaml`：

```yaml
packages:
  - 'apps/*'
  - 'packages/*'
```

根 `package.json`：

```json
{
  "name": "frontend-workspace",
  "private": true,
  "packageManager": "pnpm@9.15.0",
  "scripts": {
    "dev:admin": "pnpm --filter @app/admin dev",
    "build": "pnpm -r build",
    "test": "pnpm -r test",
    "lint": "pnpm -r lint"
  }
}
```

应用包：

```json
{
  "name": "@app/admin",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "dependencies": {
    "@repo/ui": "workspace:*",
    "@repo/utils": "workspace:*"
  }
}
```

组件库包：

```json
{
  "name": "@repo/ui",
  "version": "0.1.0",
  "type": "module",
  "main": "./dist/index.js",
  "types": "./dist/index.d.ts",
  "scripts": {
    "build": "vite build",
    "test": "vitest run"
  }
}
```

## 包之间如何依赖

推荐依赖方向：

```text
apps/admin
  ↓
packages/ui
  ↓
packages/utils
```

不要让底层包反过来依赖业务应用。

错误方向：

```text
packages/ui 依赖 apps/admin 的路由、store、接口
```

这样会导致组件库无法独立构建，也无法复用到其他应用。

## TypeScript 配置

根目录放基础配置：

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "baseUrl": ".",
    "paths": {
      "@repo/ui": ["packages/ui/src/index.ts"],
      "@repo/utils": ["packages/utils/src/index.ts"]
    }
  }
}
```

每个包继承：

```json
{
  "extends": "../../tsconfig.base.json",
  "include": ["src"]
}
```

路径别名要和构建工具保持一致，否则会出现编辑器正常、构建失败的问题。

## 脚本怎么组织

常用命令：

```bash
pnpm --filter @app/admin dev
pnpm --filter @repo/ui build
pnpm -r test
pnpm -r lint
```

含义：

| 命令 | 作用 |
| --- | --- |
| `--filter @app/admin` | 只运行某个包的命令 |
| `-r` | 递归运行所有包 |
| `workspace:*` | 使用当前 workspace 中的包 |

根目录脚本应该包装常用操作，让新人不用记复杂命令。

## 实际项目常见问题

### 1. 改了组件库，应用没有更新

### 常见原因

- 应用依赖的是已发布版本，不是 workspace 包。
- 组件库没有正确导出源码或构建产物。
- Vite 对 workspace 依赖预构建缓存没有刷新。

### 解决方案

应用中使用：

```json
{
  "dependencies": {
    "@repo/ui": "workspace:*"
  }
}
```

必要时重启 dev server，清理 Vite 缓存。

### 2. 包之间循环依赖

### 问题现象

- 构建顺序混乱。
- 类型解析异常。
- 某些导入在运行时是 `undefined`。

### 解决方案

画出依赖方向。如果两个包互相依赖，通常说明边界错了。

可以抽出更底层的包：

```text
packages/shared-types
packages/utils
```

让两个包都依赖底层包，而不是互相依赖。

### 3. 所有东西都被放进 shared

### 问题现象

`shared` 里什么都有：组件、业务常量、接口、权限、工具函数、样式。

### 根因

没有定义共享包的职责，导致它变成新的垃圾桶目录。

### 解决方案

共享包按职责拆：

```text
packages/ui
packages/utils
packages/request
packages/constants
packages/config
```

每个包都要有 README，说明它能放什么，不能放什么。

## 最佳实践

- 先有稳定复用需求，再引入 Monorepo。
- `apps` 放应用，`packages` 放可复用包。
- 包之间依赖方向必须清楚。
- 通用包不要依赖业务应用。
- 根目录脚本包装常用命令。
- 每个包都有独立 README 和构建脚本。
- CI 中按包运行 lint、test、build。

## 下一步学习

继续学习 [依赖管理](/engineering/package-management)、[测试策略](/engineering/testing)、[构建与部署](/engineering/build-deploy) 和 [工程化常见问题](/engineering/troubleshooting)。
