# 原型与原型链

## 适合谁看

适合已经会写对象、数组、函数和 class，但遇到这些问题时还不够清楚的人：

- `obj.xxx` 明明对象上没有，为什么还能访问到。
- `Array.prototype.map`、`Object.prototype.toString` 到底是什么。
- `class` 和原型是什么关系。
- 为什么不建议随便修改内置对象的 prototype。
- 面试里听过原型链，但项目里不知道怎么理解它。

原型链不是为了让你每天手写继承代码。它更重要的价值是帮助你理解 JavaScript 对象、方法查找、class 语法、内置对象和调试时看到的 `[[Prototype]]`。

## 一个对象为什么能调用方法

示例：

```ts
const users = ['Tom', 'Jerry']

users.map((name) => name.toUpperCase())
```

`users` 自己并没有直接保存一个 `map` 函数。浏览器会沿着它的原型链去找：

```text
users
↓
Array.prototype
↓
Object.prototype
↓
null
```

找到 `Array.prototype.map` 后，就能调用。

## prototype 和 `[[Prototype]]`

初学者容易把两个概念混在一起：

| 名称 | 含义 |
| --- | --- |
| `prototype` | 函数对象上的属性，通常给实例共享方法 |
| `[[Prototype]]` | 对象内部指向原型对象的链接 |
| `__proto__` | 访问 `[[Prototype]]` 的历史写法，不建议业务代码依赖 |

示例：

```ts
function User(name: string) {
  this.name = name
}

User.prototype.sayHi = function () {
  return `Hi, ${this.name}`
}

const user = new User('Tom')

user.sayHi()
```

查找过程：

```text
user.sayHi
↓ user 自己没有
User.prototype.sayHi
↓ 找到并执行
```

## class 不是另一套对象模型

`class` 只是更容易理解的语法，本质仍然基于原型。

```ts
class User {
  constructor(name: string) {
    this.name = name
  }

  sayHi() {
    return `Hi, ${this.name}`
  }
}
```

可以理解为：

```text
实例属性：放在 user 自己身上
实例方法：放在 User.prototype 上
静态方法：放在 User 构造函数本身
```

所以：

```ts
const user = new User('Tom')

user.sayHi === User.prototype.sayHi // true
```

项目里大多数时候直接使用 `class` 或对象组合就够了，不需要手动操作原型。

## 方法查找规则

当访问 `obj.name` 时，JavaScript 会按顺序查找：

1. 对象自身有没有 `name`。
2. 自身没有，就去原型对象找。
3. 还没有，就继续沿原型的原型找。
4. 找到 `null` 还没有，就返回 `undefined`。

示例：

```ts
const baseUser = {
  role: 'user'
}

const admin = Object.create(baseUser)

admin.name = 'Tom'

console.log(admin.name) // Tom
console.log(admin.role) // user
```

`role` 来自原型，不是 `admin` 自己的属性。

## 判断自有属性

项目里处理接口数据时，最好区分“对象自己有这个字段”和“从原型上找到这个字段”。

推荐：

```ts
Object.hasOwn(user, 'id')
```

兼容旧环境时：

```ts
Object.prototype.hasOwnProperty.call(user, 'id')
```

不要直接写：

```ts
user.hasOwnProperty('id')
```

因为接口数据可能覆盖或没有这个方法。

## 原型污染是什么

原型污染是指外部输入修改了对象原型，影响到其他对象。

危险示例：

```ts
function merge(target: any, source: any) {
  for (const key in source) {
    target[key] = source[key]
  }
}

merge({}, JSON.parse('{"__proto__":{"admin":true}}'))

console.log(({} as any).admin)
```

真实项目里，如果把用户输入直接深合并到对象上，就可能造成安全风险。

防护建议：

- 合并对象时过滤 `__proto__`、`constructor`、`prototype`。
- 使用成熟工具库并保持版本更新。
- 接口入参做白名单校验。
- 权限判断不要依赖前端对象上的任意字段。

## 不要随便改内置原型

不推荐：

```ts
Array.prototype.first = function () {
  return this[0]
}
```

原因：

- 可能和未来标准或第三方库冲突。
- 会影响所有数组。
- 枚举、序列化、测试时可能出现意外行为。
- 团队成员很难知道方法来自哪里。

推荐写普通工具函数：

```ts
function first<T>(items: T[]) {
  return items[0]
}
```

## 实际项目常见问题

### 1. 控制台看到 `[[Prototype]]`，以为接口多返回了字段

`[[Prototype]]` 是浏览器展示对象原型链，不是接口响应字段。判断接口字段时看对象自身属性。

### 2. `for...in` 遍历出意外字段

`for...in` 会遍历可枚举的继承属性。处理普通对象时更推荐：

```ts
for (const key of Object.keys(data)) {
  console.log(key, data[key])
}
```

### 3. class 方法作为回调时 this 丢失

```ts
class UserService {
  name = 'user'

  logName() {
    console.log(this.name)
  }
}

const service = new UserService()
setTimeout(service.logName)
```

方法被单独传出去后，调用者变了，`this` 也会变。

解决方式：

```ts
setTimeout(() => service.logName())
```

或者在 class 中使用箭头函数字段，但要理解它会成为实例属性，不是原型方法。

### 4. 深拷贝后对象方法丢失

JSON 深拷贝只保留可序列化数据，不会保留原型方法。

```ts
const copied = JSON.parse(JSON.stringify(user))
```

项目里建议接口数据保持为普通数据对象，业务方法单独放在 service、utils 或 class 中。

## 最佳实践

- 理解原型链用于读懂对象模型，不要为了炫技手写复杂继承。
- 项目里优先使用组合、普通函数和明确的数据结构。
- 不修改内置对象原型。
- 处理外部输入时防止原型污染。
- 判断对象字段时优先用 `Object.hasOwn`。
- class 方法作为回调传递时关注 `this`。

## 学习检查

学完本节后，你应该能回答：

- 对象访问属性时会按什么顺序查找。
- `prototype` 和对象内部原型链接有什么区别。
- `class` 和原型链是什么关系。
- 为什么不建议修改 `Array.prototype`。
- 原型污染为什么是安全风险。

## 参考资料

- [MDN: Inheritance and the prototype chain](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Inheritance_and_the_prototype_chain)
- [MDN: Object prototypes](https://developer.mozilla.org/en-US/docs/Learn_web_development/Extensions/Advanced_JavaScript_objects/Object_prototypes)

## 下一步学习

继续学习 [数组与对象处理](/javascript/array-object)，把对象模型落实到真实数据处理场景。
