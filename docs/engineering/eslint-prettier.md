# 代码规范

## 适合谁看

适合准备多人协作开发 Vue 项目，或者希望项目长期维护不失控的学习者。

代码规范不是为了统一个人审美，而是为了减少低级错误、降低协作成本、让项目在几个月后仍然容易改。

## 你会学到什么

- ESLint 和 Prettier 分别负责什么。
- Vue 项目应该规范哪些内容。
- 命名、目录、组件和样式如何约定。
- 提交前应该检查什么。
- 实际项目中规范冲突、自动格式化失效怎么处理。

## ESLint 和 Prettier 的区别

| 工具 | 负责 |
| --- | --- |
| ESLint | 发现代码问题，例如未使用变量、错误写法、Vue 规则 |
| Prettier | 统一格式，例如缩进、换行、引号 |

简单理解：

> ESLint 管“代码有没有问题”，Prettier 管“代码长什么样”。

## 推荐检查命令

```bash
npm run lint
npm run build
```

如果项目有类型检查：

```bash
npm run type-check
```

发布前至少保证构建通过。

## 命名规范

### 文件命名

| 类型 | 推荐 |
| --- | --- |
| 页面目录 | `users/index.vue` |
| Vue 组件 | `UserFormDrawer.vue` |
| composable | `usePagination.ts` |
| store | `user.ts` |
| API 模块 | `user.ts` |
| 类型文件 | `user.types.ts` |

### 组件命名

组件名要表达业务含义：

```text
UserTable
UserFormDrawer
RolePermissionTree
PermissionButton
```

不要使用：

```text
Table1
MyDialog
CommonComponent
```

## 目录规范

目录职责要稳定：

```text
api/          请求函数
services/     业务流程
stores/       全局状态
views/        页面
components/   可复用组件
composables/  可复用逻辑
utils/        无状态工具函数
```

不要把业务请求写到 `utils`，也不要把页面流程塞进 `components`。

## CSS 规范

业务样式必须命中明确 class：

```css
.user-search-form__actions {
  display: flex;
  gap: 8px;
}

.permission-switch-row {
  display: flex;
  align-items: center;
}
```

禁止宽泛选择器污染组件库：

```css
.page button {}
.content div {}
.panel * {}
```

如果必须调整组件库样式，优先使用主题变量、组件 props、CSS 变量或官方 API。

## 提交前检查清单

- 新增页面是否有真实路由。
- 新增组件是否职责清楚。
- 是否把接口请求放在 API 层。
- 是否把跨页面状态放在 Pinia。
- 是否避免了宽泛 CSS 选择器。
- 是否更新了相关文档。
- `npm run build` 是否通过。

## 实际项目常见问题

### 1. ESLint 和 Prettier 互相打架

**症状**

保存后格式变化，运行 lint 又改回去。

**解决方案**

使用官方推荐的集成配置，避免 ESLint 和 Prettier 同时管理格式规则。团队统一编辑器格式化设置。

### 2. 有人提交了构建失败代码

**解决方案**

增加提交前或 CI 检查：

```text
lint
type-check
build
```

本地可以用 Git hooks，团队项目更应该依赖 CI。

### 3. 项目越写越乱

**常见原因**

- 没有目录职责。
- 组件边界不清。
- 接口、状态、页面逻辑混在一起。
- 文档没有同步更新。

**解决方案**

在 README 写清楚项目结构和开发约定，并在评审时持续执行。

## 最佳实践

- 规范要少而稳定，不要一次加太多难以执行的规则。
- 自动化能解决的，不靠人工记忆。
- 构建和类型检查是最低质量门槛。
- 修改功能时同步更新对应文档。
- 样式规范必须保护组件库和响应式布局。

## 下一步学习

继续学习 [测试策略](/engineering/testing)、[依赖管理](/engineering/package-management) 和 [工程化常见问题](/engineering/troubleshooting)。
