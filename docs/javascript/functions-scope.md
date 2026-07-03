# 函数、作用域与闭包

## 适合谁看

适合能写函数，但对变量作用范围、闭包、回调、事件处理和 composable 里的状态保存还不够清楚的学习者。

Vue 项目中，函数不仅用于事件处理，也用于数据转换、权限判断、请求封装和组合式逻辑。理解作用域和闭包，能帮助你写出更稳定的 composable。

## 函数的职责

函数应该表达一个明确动作：

```ts
function formatUsername(user: User) {
  return user.nickname || user.username || '未命名用户'
}

function canDeleteUser(user: User) {
  return !user.enabled && !user.roles.includes('admin')
}
```

如果函数名需要用“并且”才能描述，通常说明它做了太多事。

## 参数和返回值

好的函数输入输出清楚：

```ts
function toUserOption(user: User) {
  return {
    label: user.username,
    value: user.id
  }
}
```

不要让函数偷偷依赖太多外部变量。

## 作用域

```ts
function submit() {
  const message = '保存成功'
  console.log(message)
}

console.log(message) // 访问不到
```

变量只在它声明的作用域内可见。使用 `const` 和 `let` 可以减少意外污染。

## 闭包

闭包可以理解为：函数记住了它创建时的外部变量。

```ts
function createCounter() {
  let count = 0

  return function increment() {
    count += 1
    return count
  }
}

const counter = createCounter()
counter() // 1
counter() // 2
```

Vue composable 里经常使用闭包保存状态：

```ts
export function useRequestLock() {
  const locked = ref(false)

  async function run(task: () => Promise<void>) {
    if (locked.value) return

    locked.value = true
    try {
      await task()
    } finally {
      locked.value = false
    }
  }

  return { locked, run }
}
```

## 实际项目常见问题

### 1. 循环里异步拿到错误的值

使用 `let` 可以避免很多旧代码中 `var` 造成的问题。

```ts
for (let index = 0; index < users.length; index++) {
  setTimeout(() => {
    console.log(users[index])
  })
}
```

### 2. 函数依赖外部状态太多

**症状**

函数很难复用，也很难测试。

**解决方案**

把依赖作为参数传入：

```ts
function filterUsers(users: User[], keyword: string) {
  return users.filter((user) => user.username.includes(keyword))
}
```

### 3. 回调函数里 this 丢失

Vue 3 `<script setup>` 中优先使用普通函数和箭头函数，不依赖 `this`。

```ts
const submit = async () => {
  await save()
}
```

## 最佳实践

- 函数只做一件清楚的事。
- 输入输出明确，少依赖外部状态。
- 复杂页面逻辑拆成多个命名函数。
- Composable 用闭包保存局部状态，但不要隐藏过多副作用。
- Vue 3 中少用 `this`，优先使用组合式 API。

## 下一步学习

继续学习 [原型与原型链](/javascript/prototype-chain)，理解对象方法查找、class 底层模型和原型污染风险。
