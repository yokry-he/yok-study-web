# 模板语法

## 适合谁看

适合刚开始写 Vue 组件，但对 `{{ }}`、`:属性`、`@事件`、`v-if`、`v-for`、`v-model` 还不够熟的学习者。

Vue 官方文档说明，Vue 使用基于 HTML 的模板语法，把渲染出来的 DOM 和组件里的 JavaScript 状态声明式绑定起来。简单说：你写状态和模板，Vue 负责根据状态更新页面。

## 你会学到什么

- 插值表达式怎么用。
- 属性绑定和事件绑定怎么写。
- 条件渲染、列表渲染和双向绑定怎么选。
- 模板里哪些逻辑不应该写。
- 实际项目中模板报错、列表错乱、条件判断异常怎么处理。

## 文本插值

最常见写法：

```vue
<script setup lang="ts">
const username = 'alice'
</script>

<template>
  <p>当前用户：{{ username }}</p>
</template>
```

`{{ }}` 中可以写简单表达式：

```vue
<p>{{ username || '未登录' }}</p>
<p>{{ count + 1 }}</p>
```

但不要写复杂业务逻辑：

```vue
<!-- 不推荐 -->
<p>{{ user.firstName.split(' ')[0] + user.lastName.toUpperCase() }}</p>
```

推荐放到 `computed`：

```ts
const displayName = computed(() => {
  return `${user.value.firstName} ${user.value.lastName}`.trim()
})
```

## 属性绑定

HTML 属性需要绑定变量时，用 `v-bind`，简写是 `:`。

```vue
<button type="button" :disabled="loading">
  保存
</button>
```

绑定 class：

```vue
<div
  class="user-row"
  :class="{ 'user-row--disabled': !user.enabled }"
>
  {{ user.username }}
</div>
```

绑定 style 时要谨慎。项目里更推荐用 class 表达状态，避免样式散落在模板中。

## 事件绑定

事件使用 `v-on`，简写是 `@`：

```vue
<button type="button" @click="submit">
  提交
</button>
```

传参数：

```vue
<button type="button" @click="removeUser(user.id)">
  删除
</button>
```

阻止默认提交：

```vue
<form @submit.prevent="submit">
  <input v-model="form.username" />
  <button type="submit">保存</button>
</form>
```

## 条件渲染

`v-if` 会创建和销毁 DOM：

```vue
<UserForm v-if="visible" />
```

`v-show` 只控制显示隐藏：

```vue
<UserForm v-show="visible" />
```

选择建议：

| 场景 | 推荐 |
| --- | --- |
| 很少切换 | `v-if` |
| 频繁切换 | `v-show` |
| 权限控制入口 | `v-if` |
| Tab 面板频繁切换 | `v-show` 或 KeepAlive |

## 列表渲染

```vue
<ul>
  <li v-for="user in users" :key="user.id">
    {{ user.username }}
  </li>
</ul>
```

`key` 必须稳定唯一。不要用数组下标：

```vue
<!-- 不推荐 -->
<li v-for="(user, index) in users" :key="index">
  {{ user.username }}
</li>
```

数组顺序变化、插入、删除时，使用 index 容易导致输入框、勾选状态、动画状态错乱。

## v-if 和 v-for 不要写在同一个元素

不推荐：

```vue
<li v-for="user in users" v-if="user.enabled" :key="user.id">
  {{ user.username }}
</li>
```

推荐先用 `computed` 过滤：

```ts
const enabledUsers = computed(() => {
  return users.value.filter((user) => user.enabled)
})
```

```vue
<li v-for="user in enabledUsers" :key="user.id">
  {{ user.username }}
</li>
```

## v-model

表单输入常用 `v-model`：

```vue
<input v-model="form.username" />
<input v-model="form.mobile" />
```

复选框：

```vue
<input v-model="form.enabled" type="checkbox" />
```

自定义组件：

```vue
<UserFormDrawer v-model:visible="drawerVisible" />
```

Vue 3.4+ 可以在子组件中用 `defineModel`：

```vue
<script setup lang="ts">
const visible = defineModel<boolean>('visible', { default: false })
</script>
```

## 实际项目常见问题

### 1. 模板里访问 undefined 报错

**症状**

控制台提示 `Cannot read properties of undefined`。

**原因**

接口数据还没回来，模板已经访问深层字段。

**解决方案**

```vue
<p>{{ user?.profile?.nickname ?? '未命名用户' }}</p>
```

或：

```vue
<UserProfile v-if="user" :user="user" />
```

### 2. 列表删除后输入框内容错乱

**原因**

`v-for` 使用了 index 作为 key。

**解决方案**

使用业务唯一 id：

```vue
<UserRow v-for="user in users" :key="user.id" :user="user" />
```

### 3. 按钮点了没反应

**排查顺序**

1. 事件函数是否写在 `<script setup>` 中。
2. 模板里函数名是否拼错。
3. 按钮是否被 `disabled`。
4. 是否被遮罩层挡住。
5. 控制台是否有前面的报错中断执行。

### 4. 权限判断散落在模板里

**问题**

模板里到处写 `permissions.includes(...)`，后续难维护。

**解决方案**

封装：

```ts
const { can } = usePermission()
```

```vue
<PermissionButton v-if="can('system:user:create')">
  新增用户
</PermissionButton>
```

## 最佳实践

- 模板只写简单表达式，复杂逻辑放到 `computed` 或函数。
- `v-for` 必须使用稳定唯一 key。
- `v-if` 和 `v-for` 不写在同一元素上。
- 表单使用 `v-model`，提交时再校验和转换。
- 权限、状态文案、复杂样式映射不要散落在模板里。

## 下一步学习

继续学习 [响应式基础](/vue/reactivity)。
