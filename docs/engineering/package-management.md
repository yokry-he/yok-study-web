# 依赖管理

## 适合谁看

适合已经能使用 `npm install`，但经常被依赖安装、版本冲突、lockfile、升级风险和安全漏洞困扰的学习者。

依赖管理不是“缺什么装什么”。真实项目里，依赖会影响启动速度、构建结果、包体积、安全漏洞、CI 稳定性和团队协作。

## 你会学到什么

- `dependencies` 和 `devDependencies` 的区别。
- 版本号中的 `^`、`~`、精确版本是什么意思。
- lockfile 为什么必须提交。
- npm、pnpm、yarn 如何选择。
- 依赖升级怎么降低风险。
- 安全漏洞和废弃依赖怎么处理。

## dependencies 和 devDependencies

| 类型 | 含义 | 示例 |
| --- | --- | --- |
| `dependencies` | 运行时需要的依赖 | `vue`、`pinia`、`axios` |
| `devDependencies` | 开发、构建、测试阶段使用 | `vite`、`eslint`、`vitest`、`typescript` |

判断方式很简单：如果生产代码运行时需要它，就放 `dependencies`；如果只是开发或构建工具，就放 `devDependencies`。

例如 Vue 应用：

```json
{
  "dependencies": {
    "vue": "^3.5.0",
    "pinia": "^2.2.0"
  },
  "devDependencies": {
    "vite": "^6.0.0",
    "typescript": "^5.6.0",
    "vitest": "^2.0.0"
  }
}
```

## 版本号怎么读

语义化版本通常是：

```text
major.minor.patch
主版本.次版本.修订版本
```

例如：

```text
3.5.12
```

| 写法 | 含义 |
| --- | --- |
| `3.5.12` | 锁定精确版本 |
| `^3.5.12` | 允许升级到 `3.x.x`，不升级到 `4.0.0` |
| `~3.5.12` | 允许升级到 `3.5.x`，不升级到 `3.6.0` |
| `latest` | 每次可能拿到最新版本，不适合业务项目 |

业务项目建议谨慎使用宽松版本，尤其是构建工具、组件库、样式库和框架核心依赖。

## lockfile 为什么重要

lockfile 记录的是实际安装的完整依赖树。

常见文件：

```text
package-lock.json
pnpm-lock.yaml
yarn.lock
```

没有 lockfile 时，两个开发者即使 `package.json` 一样，也可能安装到不同子依赖版本，导致：

- 本地能跑，CI 失败。
- A 同事构建正常，B 同事构建失败。
- 某个子依赖更新后引入 bug。
- 线上构建和本地构建不一致。

团队项目应该提交 lockfile，并且统一包管理器。

## npm、pnpm、yarn 怎么选

| 工具 | 特点 | 适合场景 |
| --- | --- | --- |
| npm | Node 默认自带，学习成本低 | 小项目、教学项目 |
| pnpm | 安装快、磁盘复用好、依赖隔离更严格 | 中大型项目、Monorepo |
| yarn | 生态成熟，有不同主版本差异 | 已经使用 yarn 的团队 |

如果是新项目，推荐优先考虑 pnpm，尤其是未来可能做 Monorepo、组件库或多应用仓库时。

项目里可以用 `packageManager` 固定工具：

```json
{
  "packageManager": "pnpm@9.15.0"
}
```

并在 README 写清楚：

```bash
corepack enable
pnpm install
pnpm dev
```

## 安装依赖的基本规则

安装运行时依赖：

```bash
pnpm add axios
```

安装开发依赖：

```bash
pnpm add -D vitest
```

删除依赖：

```bash
pnpm remove axios
```

查看依赖为什么被安装：

```bash
pnpm why lodash
```

检查过期依赖：

```bash
pnpm outdated
```

不要手动改 `node_modules`。它是安装结果，不是源码。

## 依赖升级流程

真实项目不要把所有依赖一次性全升级。推荐流程：

```text
确认升级原因
↓
阅读 changelog
↓
升级小范围依赖
↓
运行测试和构建
↓
检查关键页面
↓
记录变更和回滚方式
```

升级前先看当前版本：

```bash
pnpm list vue vite typescript
```

升级单个依赖：

```bash
pnpm up vite
```

升级到指定版本：

```bash
pnpm add vite@6.4.3 -D
```

升级后至少运行：

```bash
pnpm lint
pnpm test
pnpm build
```

如果项目没有这些脚本，至少运行能代表真实交付的最小检查命令。

## 实际项目常见问题

### 1. 删除 node_modules 后仍然安装失败

### 问题现象

- `npm install` 或 `pnpm install` 报依赖解析错误。
- 删除 `node_modules` 后问题仍然存在。
- 不同电脑表现不一致。

### 常见原因

- lockfile 和 `package.json` 不一致。
- Node 版本不一致。
- 包管理器不一致。
- 私有源或镜像源配置不一致。

### 解决方案

先确认版本：

```bash
node -v
pnpm -v
```

再确认项目要求：

```json
{
  "engines": {
    "node": ">=20"
  },
  "packageManager": "pnpm@9.15.0"
}
```

如果 lockfile 明显异常，由团队统一重新生成，不要每个人各自删。

### 2. CI 安装依赖和本地不一致

### 常见原因

CI 没有使用 lockfile 的严格安装方式。

推荐：

```bash
pnpm install --frozen-lockfile
```

它能保证 CI 不会悄悄改 lockfile。

### 3. 依赖安全漏洞很多，不知道是否都要修

### 处理方式

先分类：

| 类型 | 处理优先级 |
| --- | --- |
| 生产依赖高危漏洞 | 最高 |
| 构建工具漏洞但不进入生产运行 | 中等 |
| 示例、测试、间接依赖低危 | 根据影响评估 |

不要看到 audit 报告就盲目 `force` 升级。强制升级可能引入更大的破坏。

### 4. 组件库升级后样式变了

### 常见原因

- 组件库默认样式调整。
- 主题 token 变化。
- 旧的局部样式依赖了组件库内部 DOM。

### 解决方案

- 阅读迁移指南。
- 检查主题配置。
- 检查是否有 `.xxx div`、`.xxx button` 这类污染选择器。
- 浏览器验证表格、表单、弹窗、开关、按钮等关键控件。

## 最佳实践

- 团队统一包管理器，不混用 npm、pnpm、yarn。
- 提交 lockfile。
- CI 使用 frozen lockfile。
- 不使用 `latest` 作为业务项目依赖版本。
- 依赖升级分批进行，并记录升级原因。
- 组件库、构建工具、框架核心依赖升级后必须浏览器验证。
- 定期清理不用的依赖，减少包体积和安全风险。

## 下一步学习

继续学习 [测试策略](/engineering/testing)、[Monorepo 项目组织](/engineering/monorepo) 和 [工程化常见问题](/engineering/troubleshooting)。
