# 模块化与工程实践

## 适合谁看

适合已经会写函数，但不知道项目里如何拆文件、如何导入导出、如何组织工具函数和业务模块的学习者。

模块化能让项目从“一个文件能跑”变成“多个功能能长期维护”。Vue 项目里的 API、service、store、router、utils 都依赖模块化。

## 导出和导入

命名导出：

```ts
export function formatDate(value: string) {
  return value.slice(0, 10)
}

export function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`
}
```

导入：

```ts
import { formatDate, formatMoney } from '@/utils/format'
```

默认导出：

```ts
export default router
```

导入：

```ts
import router from '@/router'
```

团队里建议统一规则，不要混乱使用。

## 文件职责

```text
api/user.ts          用户接口
services/auth.ts     登录流程
stores/user.ts       用户状态
utils/format.ts      通用格式化函数
types/user.ts        用户类型
```

一个文件只负责一类事情。

## 避免循环依赖

循环依赖例子：

```text
store/user.ts -> router/index.ts -> store/user.ts
```

可能导致初始化时拿到空对象或报错。

解决思路：

- 把共享逻辑抽到第三个模块。
- 避免在模块顶层立即执行复杂逻辑。
- 在函数内部延迟读取依赖。

## 入口文件

可以用 `index.ts` 统一导出：

```ts
export * from './format'
export * from './permission'
```

但不要滥用。过大的 barrel file 会让依赖关系变模糊。

## 实际项目常见问题

### 1. utils 变成垃圾桶

**症状**

所有东西都放 `utils`，包括接口、权限、业务流程。

**解决方案**

按职责拆分：

- 通用无状态函数放 `utils`。
- 请求放 `api`。
- 业务流程放 `services`。
- 跨页面状态放 `stores`。

### 2. 导入路径太深

**解决方案**

配置 `@` 别名：

```ts
import { formatDate } from '@/utils/format'
```

Vite 和 TypeScript 都要配置。

### 3. 模块初始化时访问不到 Pinia

**原因**

在 app 安装 Pinia 之前，模块顶层就调用了 store。

**解决方案**

在函数内部调用，确保运行时 Pinia 已安装。

## 最佳实践

- 文件名表达职责。
- 不要把业务流程放进 utils。
- 减少模块顶层副作用。
- 警惕循环依赖。
- 导入路径使用项目统一别名。

## 下一步学习

继续进入 [TypeScript 学习导览](/typescript/introduction)，把 JavaScript 代码逐步升级为更稳的类型化工程代码。
