# 工程性能优化

## 适合谁看

适合项目变大后，开发体验和交付速度明显变差的人：

- `npm install` 很慢。
- `npm run dev` 启动慢。
- HMR 更新慢。
- TypeScript 检查慢。
- CI 构建越来越久。
- 每次改一行代码都要等很久。

工程性能优化不是页面运行性能，而是开发、构建、测试、发布这些工程流程的速度和稳定性。

## 先测量再优化

先记录基线：

```bash
time npm install
time npm run dev
time npm run typecheck
time npm run test
time npm run build
```

还要记录：

- Node 版本。
- 包管理器版本。
- 操作系统。
- CI 机器规格。
- 缓存是否命中。
- 分支依赖变化。

没有基线就无法判断优化是否有效。

## 安装速度

影响安装速度的因素：

- 依赖数量。
- lockfile 是否稳定。
- 是否频繁删除缓存。
- 包管理器选择。
- 私有源速度。
- postinstall 脚本。

优化方向：

- 使用稳定 lockfile。
- CI 缓存包管理器目录。
- 删除无用依赖。
- 避免不必要的 postinstall。
- 大依赖按需引入。

## Vite dev server 启动慢

Vite 首次启动会做依赖预构建。依赖越复杂，启动越可能变慢。

排查：

- 是否引入了很大的 CommonJS 依赖。
- 是否 workspace 包没有正确构建。
- 是否频繁清理 Vite 缓存。
- 是否插件太多。
- 是否自动导入扫描范围过大。

优化：

```ts
export default defineConfig({
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia']
  }
})
```

不要盲目配置。先观察 Vite 输出和实际慢在哪里。

## HMR 慢

常见原因：

- 单个文件过大。
- 组件依赖链太长。
- 全局状态改动影响太多页面。
- 自动导入或插件扫描成本高。
- 样式文件全局影响过大。

解决方向：

- 拆分大组件。
- 降低全局状态粒度。
- 页面级模块边界清晰。
- 业务样式局部化。
- 插件扫描目录收窄。

## TypeScript 慢

常见原因：

- 类型过度复杂。
- 大量深层泛型。
- `include` 范围过大。
- 生成文件也进入类型检查。
- Monorepo 没有项目引用。

优化：

- `tsconfig` 排除构建产物和生成目录。
- 复杂类型拆分。
- Monorepo 使用 project references。
- CI 类型检查和本地增量检查分层。
- 业务代码不要追求炫技类型。

类型系统服务业务，不是展示技巧。

## 测试慢

测试慢常见原因：

- 单元测试启动了真实浏览器。
- 大量测试依赖真实网络。
- 每个用例重复初始化大对象。
- 不区分单元测试和 E2E。
- 快照过多且频繁变化。

优化：

- 单元测试优先纯函数和组件边界。
- API、数据库、E2E 分层执行。
- CI 按变更范围执行。
- 失败时保留日志和截图。
- 慢测试单独标记。

## CI 慢

CI 优化思路：

```text
install
↓
lint
↓
typecheck
↓
test
↓
build
↓
deploy
```

可优化：

- 缓存依赖。
- 缓存构建产物。
- 并行 lint、typecheck、test。
- 只在主分支部署。
- PR 阶段不跑全量 E2E。
- 对文档修改跳过无关任务。

但不要为了速度删掉质量门禁。要分层，而不是取消。

## Monorepo 性能

Monorepo 需要关注任务编排：

- 哪些包受影响。
- 哪些任务可以缓存。
- 哪些任务可以并行。
- 构建顺序是否正确。

典型策略：

```text
只构建受影响 package
↓
复用远程缓存
↓
共享 lint / tsconfig
↓
应用和包独立发布
```

如果 Monorepo 只是把所有项目放一起，但没有任务缓存和边界，反而会变慢。

## 实际项目问题

### 1. dev server 第一次启动很慢

**原因**

依赖预构建成本高，或者依赖图过大。

**解决方案**

- 查看 Vite 性能建议。
- 删除不用依赖。
- 明确 optimizeDeps。
- 避免每次启动都清缓存。

### 2. 改一个基础组件导致整站刷新

**原因**

基础组件被大量页面依赖，HMR 影响范围大。

**解决方案**

- 基础组件保持稳定。
- 拆分大组件。
- 避免基础组件引用业务模块。
- 将业务逻辑下沉到页面或 feature。

### 3. CI 每次都跑半小时

**原因**

所有任务串行、无缓存、无变更范围判断。

**解决方案**

- 依赖缓存。
- lint、typecheck、test 并行。
- E2E 只在关键分支或定时任务跑全量。
- Monorepo 使用 affected 构建。

### 4. 类型检查越来越慢

**原因**

tsconfig include 范围过大，或复杂类型蔓延。

**解决方案**

- 排除 dist、coverage、generated。
- 拆分 tsconfig。
- 减少超复杂泛型。
- 对自动生成类型单独管理。

## 最佳实践

- 所有优化先有基线数据。
- 依赖数量、插件数量和扫描范围都要控制。
- 工程配置要分层，不把业务逻辑塞进构建配置。
- CI 通过缓存、并行和变更范围提速。
- Monorepo 必须有清晰包边界和任务缓存。
- 类型系统要服务可维护性，不追求无意义复杂度。
- 每次工程性能优化都记录原因、前后数据和影响范围。

## 参考资料

- [Vite Performance](https://vite.dev/guide/performance)
- [Vite Dependency Pre-Bundling](https://vite.dev/guide/dep-pre-bundling)
- [Vite Build Options](https://vite.dev/config/build-options)

## 下一步学习

继续学习 [工程化常见问题](/engineering/troubleshooting)，把安装、启动、构建、测试和部署问题串起来排查。
