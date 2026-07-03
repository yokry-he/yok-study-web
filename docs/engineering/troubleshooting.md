# 工程化常见问题

## 适合谁看

适合在前端项目启动、安装依赖、构建、测试、部署时遇到问题，但不知道从哪一层开始排查的学习者。

工程化问题通常不是业务代码单点错误，而是 Node 版本、包管理器、lockfile、环境变量、Vite 配置、缓存、路径和部署平台共同影响的结果。

## 排查总原则

遇到工程化问题，先确认四件事：

```text
当前命令是什么
当前 Node 和包管理器版本是什么
错误发生在安装、启动、构建、测试还是部署
最近改过哪些配置或依赖
```

不要一上来就删除所有文件。先保留错误信息，再定位问题发生在哪个阶段。

## 问题 1：依赖安装失败

### 问题现象

- `npm install` 或 `pnpm install` 报错。
- 换一台电脑正常，当前电脑不正常。
- CI 安装失败，本地安装成功。

### 常见原因

- Node 版本不符合项目要求。
- npm、pnpm、yarn 混用。
- lockfile 和包管理器不匹配。
- 依赖源或私有仓库权限异常。
- 某个依赖版本被废弃或无法下载。

### 解决方案

先确认版本：

```bash
node -v
pnpm -v
npm -v
```

检查项目声明：

```json
{
  "engines": {
    "node": ">=20"
  },
  "packageManager": "pnpm@9.15.0"
}
```

CI 中使用严格安装：

```bash
pnpm install --frozen-lockfile
```

如果项目使用 pnpm，就不要再提交 `package-lock.json` 或 `yarn.lock`。

### 预防方式

- README 写清 Node 版本和包管理器。
- 提交 lockfile。
- 团队不要混用多个包管理器。
- CI 使用 frozen lockfile。

## 问题 2：本地启动成功，但页面空白

### 问题现象

- `npm run dev` 没有明显报错。
- 浏览器打开是白屏。
- Console 有运行时错误或资源加载错误。

### 常见原因

- 入口文件报错。
- 路径别名配置不一致。
- 环境变量缺失。
- API 请求阻塞了启动流程。
- 老的 dev server 或缓存影响了页面。

### 解决方案

按顺序排查：

```text
1. 看终端是否有编译错误。
2. 看浏览器 Console 第一条错误。
3. 看 Network 是否有 JS、CSS、接口 404。
4. 重启 dev server。
5. 检查 .env 是否重启后生效。
```

如果是路径别名问题，同时检查：

```text
vite.config.ts
tsconfig.json
实际文件大小写
```

macOS 对大小写不敏感时，本地可能正常，Linux CI 会失败。

### 预防方式

- 构建前跑类型检查。
- 不在应用启动时阻塞等待非必要接口。
- 环境变量集中读取并提供明确缺省值。
- 文件名大小写保持一致。

## 问题 3：构建失败但开发环境正常

### 问题现象

- `npm run dev` 正常。
- `npm run build` 报错。
- 报错可能来自 TypeScript、Rollup、Vite 插件或依赖。

### 常见原因

- 开发环境没有执行完整类型检查。
- 动态导入路径无法被构建工具静态分析。
- 某些依赖只支持 Node 环境，不能进入浏览器包。
- 文件大小写在 CI 环境不一致。
- 环境变量只在本地存在。

### 解决方案

本地复现生产构建：

```bash
npm run build
```

如果项目有类型检查：

```bash
npm run type-check
```

动态导入要让构建工具能识别：

```ts
const modules = import.meta.glob('../views/**/*.vue')
```

不要在浏览器代码里直接使用 Node API：

```ts
// 不适合浏览器业务代码
import fs from 'node:fs'
```

### 预防方式

- PR 或合并前必须跑 build。
- 类型检查和构建都进入 CI。
- 文件名大小写在导入和实际文件中保持一致。
- 浏览器代码不要混入 Node-only API。

## 问题 4：修改环境变量后不生效

### 问题现象

- 改了 `.env`。
- 页面仍然读取旧值。
- `import.meta.env.xxx` 是 `undefined`。

### 常见原因

- Vite 需要重启 dev server 才读取新的环境变量。
- 变量没有以 `VITE_` 开头。
- 写错了 mode 对应的文件。
- 构建使用了错误 mode。

### 解决方案

前端可暴露变量必须以 `VITE_` 开头：

```ini
VITE_API_BASE_URL=/api
```

重启开发服务器。

构建时指定 mode：

```bash
vite build --mode production
vite build --mode test
```

### 预防方式

- 环境变量集中到 `src/config/app.ts`。
- README 写清每个环境的构建命令。
- 构建日志输出当前 mode 和 API 前缀。

## 问题 5：测试在本地通过，CI 失败

### 问题现象

- 本地 `npm run test` 通过。
- CI 测试失败。
- 错误和时间、时区、随机数、网络或文件路径有关。

### 常见原因

- 测试依赖真实时间。
- 测试依赖真实网络。
- 测试顺序不独立。
- 本地和 CI Node 版本不同。
- 快照包含不稳定内容。

### 解决方案

不要让单元测试依赖真实网络：

```ts
vi.mock('../api/user', () => ({
  getUser: vi.fn()
}))
```

时间相关逻辑固定时间：

```ts
vi.setSystemTime(new Date('2026-07-01T00:00:00Z'))
```

每个测试后清理 mock：

```ts
afterEach(() => {
  vi.restoreAllMocks()
})
```

### 预防方式

- 测试之间互相独立。
- mock 网络、时间和随机数。
- CI 和本地 Node 版本一致。
- 不把不稳定输出写进快照。

## 问题 6：上线后页面白屏或静态资源 404

### 问题现象

- 构建成功。
- 线上打开白屏。
- Network 中 JS 或 CSS 资源 404。
- 子路径部署更容易出现。

### 常见原因

- Vite `base` 配错。
- Router history base 配错。
- Nginx 没有配置 fallback。
- `index.html` 引用了不存在的旧资源。
- CDN 还缓存旧入口。

### 解决方案

检查部署路径。如果项目部署在 `/admin/`：

```ts
export default defineConfig({
  base: '/admin/'
})
```

Router 也要同步：

```ts
createWebHistory('/admin/')
```

Nginx 配置 fallback：

```nginx
location /admin/ {
  try_files $uri $uri/ /admin/index.html;
}
```

### 预防方式

- 上线前用真实路径预览。
- 发布后检查首页、二级路由和详情页刷新。
- `index.html` 不强缓存。
- 保留旧静态资源一段时间，降低灰度期间 404 风险。

## 下一步学习

继续学习 [Vite 工程基础](/engineering/vite)、[环境配置](/engineering/env-config)、[依赖管理](/engineering/package-management)、[测试策略](/engineering/testing) 和 [构建与部署](/engineering/build-deploy)。
