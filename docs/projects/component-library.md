# 组件库实战

## 项目目标

组件库实战的目标不是“把所有 UI 都封装一遍”，而是学会沉淀稳定、重复、边界清晰的组件能力。

在 Vue Admin 项目里，最值得沉淀的通常不是按钮和输入框，因为这些基础控件应该优先使用组件库；更值得沉淀的是结合业务流程的组件，例如页面容器、权限按钮、状态标签、搜索工具栏。

## 你会学到什么

- 什么组件值得封装。
- 组件 API 如何设计。
- 如何写 props、emits、slot 和文档示例。
- 如何避免组件过度封装。
- 实际项目里组件库失控怎么治理。

## 第一批组件

| 组件 | 作用 | 是否依赖业务 |
| --- | --- | --- |
| `AppPage` | 页面容器，统一标题、操作区和内容区 | 低 |
| `SearchToolbar` | 搜索筛选区域 | 中 |
| `DataActionBar` | 表格批量操作区 | 中 |
| `PermissionButton` | 权限控制按钮 | 高 |
| `StatusTag` | 状态展示标签 | 中 |
| `EmptyState` | 空状态展示 | 低 |

## AppPage

页面容器负责统一页面结构：

```vue
<template>
  <section class="app-page">
    <header class="app-page__header">
      <div>
        <h1 class="app-page__title">{{ title }}</h1>
        <p v-if="description" class="app-page__description">
          {{ description }}
        </p>
      </div>

      <div class="app-page__actions">
        <slot name="actions" />
      </div>
    </header>

    <main class="app-page__body">
      <slot />
    </main>
  </section>
</template>
```

使用：

```vue
<AppPage title="用户管理" description="管理系统用户和启用状态">
  <template #actions>
    <PermissionButton code="system:user:create">
      新增用户
    </PermissionButton>
  </template>

  <UserTable :users="users" />
</AppPage>
```

## PermissionButton

权限按钮统一处理权限判断：

```vue
<script setup lang="ts">
const props = withDefaults(defineProps<{
  code: string
  disabled?: boolean
}>(), {
  disabled: false
})

const { can } = usePermission()
const allowed = computed(() => can(props.code))
</script>

<template>
  <button
    v-if="allowed"
    type="button"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>
```

如果项目使用组件库，应把内部 `button` 换成组件库按钮，例如 `NButton`、`ElButton` 或 `AButton`。

## StatusTag

状态标签负责把业务状态映射成稳定展示：

```ts
type UserStatus = 'enabled' | 'disabled' | 'locked'

const statusMap: Record<UserStatus, { label: string; tone: string }> = {
  enabled: { label: '启用', tone: 'success' },
  disabled: { label: '停用', tone: 'warning' },
  locked: { label: '锁定', tone: 'danger' }
}
```

使用：

```vue
<StatusTag :status="user.status" />
```

好处是状态文案和颜色集中维护，不会每个页面写一套。

## SearchToolbar

搜索区常见需求：

- 输入关键字。
- 选择状态。
- 点击查询。
- 点击重置。
- 移动端可换行。

组件 API 不要过度设计。第一阶段可以让它只负责布局：

```vue
<SearchToolbar>
  <input v-model="query.keyword" placeholder="搜索用户名" />
  <select v-model="query.enabled">
    <option :value="undefined">全部状态</option>
    <option :value="true">启用</option>
    <option :value="false">停用</option>
  </select>

  <template #actions>
    <button type="button" @click="search">查询</button>
    <button type="button" @click="reset">重置</button>
  </template>
</SearchToolbar>
```

## 组件文档应该写什么

每个组件至少写：

```text
组件用途
适合场景
Props
Emits
Slots
基础示例
常见错误
移动端表现
```

如果没有文档，组件库会变成“只有作者知道怎么用”的代码堆。

## 实际项目常见问题

### 1. 组件封装太早

**症状**

刚写一个页面就抽了很多通用组件，后续需求稍微变化就很难改。

**解决方案**

至少出现 2 到 3 次稳定重复后再抽公共组件。第一次先在页面内写清楚。

### 2. 组件 API 太复杂

**症状**

一个表格组件有几十个 props，使用成本比直接写表格还高。

**解决方案**

拆分职责。布局组件负责布局，业务组件负责业务，基础控件交给组件库。

### 3. 组件内部写死业务接口

**症状**

组件只能用于用户管理，不能用于角色管理。

**解决方案**

可复用组件通过 props 接收数据，通过 emits 通知事件，不直接请求具体业务接口。

### 4. 组件样式污染全局

**解决方案**

使用明确 class，避免宽泛选择器。组件库内部 DOM 不要依赖层级覆盖。

## 最佳实践

- 先沉淀业务重复，再抽组件。
- 组件 API 要少、稳定、可解释。
- 基础控件优先使用成熟组件库。
- 组件文档和示例必须同步维护。
- 组件要检查移动端和窄容器表现。

## 下一步

把 Vue Admin 中重复出现的页面容器、权限按钮、状态标签逐步沉淀出来，不要一次性大规模抽象。
