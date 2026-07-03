# 测试策略

## 适合谁看

适合已经能写前端页面，但希望项目改动后更有把握、不再完全依赖人工点页面的学习者。

测试不是为了追求覆盖率数字，而是为了保护核心业务行为。真实项目里，测试应该优先覆盖那些“经常改、影响大、靠人工容易漏”的部分。

## 你会学到什么

- 前端测试分哪些层级。
- 单元测试、组件测试、端到端测试分别测什么。
- Vue 项目中哪些代码最值得优先测试。
- 如何给表单、权限、请求和工具函数写测试。
- 测试失败时如何定位问题。

## 测试金字塔

前端常见测试分层：

| 类型 | 关注点 | 示例工具 | 成本 |
| --- | --- | --- | --- |
| 单元测试 | 函数、组合式逻辑、数据转换 | Vitest | 低 |
| 组件测试 | 组件渲染、事件、props、插槽 | Vue Test Utils | 中 |
| 端到端测试 | 真实浏览器流程 | Playwright、Cypress | 高 |
| 静态检查 | 类型、lint、格式 | TypeScript、ESLint | 低 |

推荐顺序：

```text
静态检查
↓
单元测试
↓
组件测试
↓
少量关键 E2E
```

不要一开始就把所有页面都写成端到端测试。E2E 很有价值，但维护成本也更高。

## 哪些代码最值得测试

优先测试：

- 权限判断。
- 表单校验。
- 数据转换。
- 请求错误处理。
- 分页、排序、筛选参数组装。
- 金额、时间、状态转换。
- 复杂 composable。
- 关键业务流程。

可以少测或不测：

- 纯展示静态页面。
- 只包了一层样式的组件。
- 频繁变化的临时页面。
- 已经由组件库保证的基础交互。

## Vitest 基础示例

假设有状态转换函数：

```ts
export function getOrderStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    canceled: '已取消'
  }

  return map[status] || '未知状态'
}
```

测试：

```ts
import { describe, expect, it } from 'vitest'
import { getOrderStatusText } from './order'

describe('getOrderStatusText', () => {
  it('returns known status text', () => {
    expect(getOrderStatusText('paid')).toBe('已支付')
  })

  it('returns fallback for unknown status', () => {
    expect(getOrderStatusText('archived')).toBe('未知状态')
  })
})
```

这类测试成本低，但能避免状态文案、兜底逻辑和数据转换被改坏。

## 测试 composable

例如分页逻辑：

```ts
import { computed, reactive } from 'vue'

export function usePagination() {
  const pagination = reactive({
    page: 1,
    pageSize: 20,
    total: 0
  })

  const offset = computed(() => {
    return (pagination.page - 1) * pagination.pageSize
  })

  function resetPage() {
    pagination.page = 1
  }

  return {
    pagination,
    offset,
    resetPage
  }
}
```

测试：

```ts
import { describe, expect, it } from 'vitest'
import { usePagination } from './usePagination'

describe('usePagination', () => {
  it('calculates offset', () => {
    const { pagination, offset } = usePagination()

    pagination.page = 3
    pagination.pageSize = 20

    expect(offset.value).toBe(40)
  })

  it('resets page', () => {
    const { pagination, resetPage } = usePagination()

    pagination.page = 5
    resetPage()

    expect(pagination.page).toBe(1)
  })
})
```

可复用逻辑越复杂，越应该测试。这样页面改动时不用反复人工验证所有分支。

## 组件测试测什么

组件测试不要重复测试浏览器和组件库已经保证的能力，而要测试你的业务约定。

例如用户表单组件：

```text
UserForm
├─ 新增模式不显示 id
├─ 编辑模式回填 username 和 mobile
├─ 手机号为空时阻止提交
├─ 点击保存时 emit submit
└─ 提交中按钮禁用
```

测试重点：

- props 变化后展示是否正确。
- 用户输入是否更新表单。
- 点击按钮是否触发事件。
- 表单校验是否阻止错误提交。
- loading、disabled、empty 等状态是否正确。

## 请求测试怎么做

请求相关测试不要直接打真实后端。应该 mock 网络层或 API 函数。

例如 service 层：

```ts
export async function loadUserOptions(keyword: string) {
  const res = await userApi.getOptions({ keyword })

  return res.items.map((item) => ({
    label: item.username,
    value: item.id
  }))
}
```

测试重点不是 HTTP 本身，而是：

- 参数是否正确。
- 返回数据是否被转换成组件需要的格式。
- 空数据是否能处理。
- 错误是否向上抛出或转成用户提示。

## E2E 测试适合覆盖什么

端到端测试适合覆盖少量核心流程：

- 登录。
- 进入核心页面。
- 搜索列表。
- 新增一条数据。
- 编辑并保存。
- 权限用户看不到某按钮。
- 刷新二级路由不 404。

不建议把每个按钮都写成 E2E。页面变化频繁时，E2E 维护成本会迅速上升。

## CI 中怎么安排测试

推荐质量门禁：

```text
install
↓
lint
↓
type-check
↓
unit test
↓
build
↓
关键 E2E
```

脚本示例：

```json
{
  "scripts": {
    "lint": "eslint .",
    "type-check": "vue-tsc --noEmit",
    "test": "vitest run",
    "build": "vite build",
    "test:e2e": "playwright test"
  }
}
```

如果项目刚开始补测试，先不要强求所有测试类型一步到位。可以先把 `lint`、`type-check` 和少量核心单元测试纳入 CI。

## 实际项目常见问题

### 1. 测试太难写

### 常见原因

- 组件里混了请求、状态、数据转换和 UI。
- 没有把业务逻辑拆到 composable 或 service。
- 代码直接依赖全局对象，不方便 mock。

### 解决方案

先重构边界：

```text
组件负责展示和交互
composable 负责状态逻辑
service 负责业务流程
api 负责请求
utils 负责纯函数
```

边界清楚后，测试会自然变简单。

### 2. 测试一改页面就挂

### 常见原因

测试依赖了太多 DOM 细节，例如第几个 `div`、第几个按钮。

### 解决方案

测试用户能感知的行为：

- 文案是否出现。
- 按钮是否可点击。
- 表单错误是否展示。
- 事件是否触发。

不要测试组件库内部 DOM 结构。

### 3. 覆盖率很高但线上还是出问题

### 常见原因

覆盖率只说明代码被执行过，不说明关键业务被验证过。

### 解决方案

用业务风险排序测试优先级：

1. 登录和权限。
2. 数据提交和修改。
3. 金额、状态、时间。
4. 错误处理。
5. 发布后最容易出问题的页面。

## 最佳实践

- 先测稳定业务逻辑，再测频繁变化页面。
- 组件测试关注业务行为，不依赖组件库内部 DOM。
- 请求测试 mock API，不打真实后端。
- E2E 只覆盖关键主流程。
- CI 至少包含 lint、type-check、test 和 build。
- 测试失败要能快速定位，不要只追求数量。

## 下一步学习

继续学习 [依赖管理](/engineering/package-management)、[Monorepo 项目组织](/engineering/monorepo) 和 [构建与部署](/engineering/build-deploy)。
