# 内置组件

## 适合谁看

适合已经能写普通组件，准备处理动画、缓存、弹窗挂载、异步组件加载等更复杂页面体验的学习者。

Vue 提供了一些内置组件，例如 `Transition`、`TransitionGroup`、`KeepAlive`、`Teleport`、`Suspense`。它们不是每个页面都要用，但在合适场景能明显提升体验和结构清晰度。

## 你会学到什么

- 每个内置组件解决什么问题。
- 后台项目里哪些场景适合使用。
- 使用时容易踩哪些坑。
- 什么时候不该用。

## Transition

`Transition` 用于元素或组件进入、离开 DOM 时的动画。

```vue
<Transition name="fade">
  <div v-if="visible" class="message-panel">
    保存成功
  </div>
</Transition>
```

CSS：

```css
.fade-enter-active,
.fade-leave-active {
  transition: opacity 160ms ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
```

适合：

- 提示条。
- 弹层内容。
- 小范围状态切换。

不适合：

- 大量列表项同时复杂动画。
- 影响用户效率的长动画。

## TransitionGroup

`TransitionGroup` 用于列表项插入、删除和移动动画。

```vue
<TransitionGroup name="list" tag="ul">
  <li v-for="item in items" :key="item.id">
    {{ item.name }}
  </li>
</TransitionGroup>
```

关键点：必须有稳定唯一 `key`。

适合：

- 待办列表。
- 通知列表。
- 小数据量排序或删除。

后台大表格通常不建议做复杂列表动画，优先保证性能和可读性。

## KeepAlive

`KeepAlive` 用于缓存组件实例。Vue 官方文档说明，它可以在多个动态组件之间切换时缓存组件实例。

```vue
<KeepAlive :include="['UserList']">
  <component :is="currentView" />
</KeepAlive>
```

在路由中常见：

```vue
<RouterView v-slot="{ Component }">
  <KeepAlive>
    <component :is="Component" />
  </KeepAlive>
</RouterView>
```

适合：

- 列表页返回后保留筛选条件。
- Tab 页面切换保留状态。

注意：

- 被缓存组件不会频繁卸载。
- 需要使用 `onActivated` 和 `onDeactivated` 处理恢复和暂停。
- 不要无脑缓存所有页面。

## Teleport

`Teleport` 可以把组件模板的一部分渲染到当前 DOM 层级之外。官方文档常见例子是弹窗、浮层等视觉上需要脱离父容器的内容。

```vue
<Teleport to="body">
  <div v-if="visible" class="modal">
    弹窗内容
  </div>
</Teleport>
```

适合：

- Modal。
- Drawer。
- 全局提示。
- 浮层菜单。

如果项目使用组件库，弹窗、抽屉通常已经内部处理了 Teleport，不需要你自己再写。

## Suspense

`Suspense` 用于协调异步依赖，可以在等待嵌套异步组件时显示 loading。Vue 官方文档标注它仍是实验性特性，因此业务项目里要谨慎使用。

```vue
<Suspense>
  <template #default>
    <AsyncDashboard />
  </template>

  <template #fallback>
    <div>加载中...</div>
  </template>
</Suspense>
```

建议：

- 学习了解即可。
- 生产业务中优先使用明确的 loading 状态。
- 除非团队确认 API 风险，否则不要把关键流程强依赖 Suspense。

## 实际项目常见问题

### 1. KeepAlive 后页面数据不刷新

**原因**

缓存页面再次显示时不会重新执行 `onMounted`。

**解决方案**

```ts
onActivated(() => {
  fetchList()
})
```

### 2. 弹窗被父容器 overflow 裁剪

**原因**

弹窗渲染在有 `overflow: hidden` 的父容器内。

**解决方案**

使用组件库弹窗，或使用 `Teleport to="body"`。

### 3. 动画导致布局抖动

**原因**

动画改变了宽高、margin 等会触发布局计算的属性。

**解决方案**

优先使用 `opacity` 和 `transform`。

### 4. 缓存页面占用过多内存

**原因**

过多页面被 KeepAlive 缓存，表格、图表和大对象没有释放。

**解决方案**

限制缓存范围：

```vue
<KeepAlive :include="['UserList', 'RoleList']" :max="10">
  <component :is="Component" />
</KeepAlive>
```

## 最佳实践

- 动画要短、轻、可关闭，尊重 `prefers-reduced-motion`。
- KeepAlive 只缓存明确需要恢复状态的页面。
- 弹窗和浮层优先使用组件库能力。
- Suspense 仍需谨慎，不作为第一阶段核心依赖。
- 内置组件是增强体验的工具，不是所有页面的默认配置。

## 下一步学习

继续学习 [性能优化](/vue/performance)。
