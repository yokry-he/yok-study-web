# 正则表达式

## 适合谁看

适合已经会写字符串处理，但遇到校验、搜索、替换时经常复制正则却不敢改的人：

- 表单里要校验手机号、邮箱、编码或文件名。
- 搜索框要支持关键词匹配和高亮。
- 需要从日志、URL、文本里提取关键信息。
- `.*`、`?`、`+`、`[]`、`() `经常分不清。
- 正则写出来能跑，但边界条件很多。

正则表达式是处理字符串模式的工具。它很强，但也容易变成不可维护的“魔法字符串”。项目里应该把正则用于合适场景，并给复杂规则命名和注释。

## 正则能做什么

常见用途：

| 场景 | 示例 |
| --- | --- |
| 判断格式 | 邮箱、手机号、订单号 |
| 查找文本 | 找出所有标签、关键词 |
| 替换内容 | 脱敏手机号、替换空格 |
| 提取字段 | 从日志中提取时间、状态码 |
| 拆分字符串 | 按多个分隔符切分 |

不适合：

- 解析复杂 HTML。
- 解析完整编程语言。
- 处理层级很深的嵌套结构。
- 替代后端严肃校验。

前端正则可以提升体验，但安全和业务最终约束必须由后端保证。

## 创建正则

字面量写法：

```ts
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
```

构造函数写法：

```ts
const keyword = 'vue'
const pattern = new RegExp(keyword, 'i')
```

固定规则优先用字面量。动态规则才使用 `RegExp` 构造函数。

动态关键词要先转义特殊字符：

```ts
function escapeRegExp(input: string) {
  return input.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

const pattern = new RegExp(escapeRegExp(keyword), 'gi')
```

否则用户输入 `.`、`*`、`[` 这类字符会改变匹配含义。

## 常用方法

| 方法 | 用途 |
| --- | --- |
| `pattern.test(text)` | 判断是否匹配 |
| `text.match(pattern)` | 获取匹配结果 |
| `text.matchAll(pattern)` | 获取所有匹配及分组 |
| `text.replace(pattern, value)` | 替换内容 |
| `text.split(pattern)` | 按模式拆分 |

判断：

```ts
const isOrderNo = /^ORD-\d{6}$/.test('ORD-123456')
```

替换：

```ts
function maskPhone(phone: string) {
  return phone.replace(/^(\d{3})\d{4}(\d{4})$/, '$1****$2')
}
```

提取：

```ts
const text = 'status=500 duration=230ms'
const match = text.match(/status=(\d+)/)

console.log(match?.[1])
```

## 基础语法

| 写法 | 含义 |
| --- | --- |
| `.` | 任意单个字符，通常不匹配换行 |
| `\d` | 数字 |
| `\w` | 字母、数字、下划线 |
| `\s` | 空白字符 |
| `^` | 开头 |
| `$` | 结尾 |
| `[]` | 字符集合 |
| `[^]` | 排除集合 |
| `()` | 分组 |
| `|` | 或 |
| `?` | 0 次或 1 次 |
| `*` | 0 次或多次 |
| `+` | 1 次或多次 |
| `{n,m}` | 次数范围 |

示例：

```ts
/^\d{4}-\d{2}-\d{2}$/.test('2026-07-02')
```

表示：

```text
开头
4 位数字
-
2 位数字
-
2 位数字
结尾
```

## 贪婪和非贪婪

默认是贪婪匹配，会尽量多匹配。

```ts
const text = '<span>Vue</span><span>React</span>'

text.match(/<span>.*<\/span>/)?.[0]
```

会匹配到：

```text
<span>Vue</span><span>React</span>
```

非贪婪：

```ts
text.match(/<span>.*?<\/span>/)?.[0]
```

会尽量少匹配：

```text
<span>Vue</span>
```

但不要用正则解析复杂 HTML。这里只是说明贪婪行为。

## 分组和命名分组

普通分组：

```ts
const match = '2026-07-02'.match(/^(\d{4})-(\d{2})-(\d{2})$/)

console.log(match?.[1])
```

命名分组更清楚：

```ts
const match = '2026-07-02'.match(
  /^(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})$/
)

console.log(match?.groups?.year)
```

命名分组适合日志解析、URL 解析、复杂文本提取。

## 常用 flags

| flag | 含义 |
| --- | --- |
| `g` | 全局匹配 |
| `i` | 忽略大小写 |
| `m` | 多行模式 |
| `s` | 让 `.` 匹配换行 |
| `u` | Unicode 模式 |
| `y` | 从指定位置粘连匹配 |

高频组合：

```ts
const pattern = /vue/gi
```

表示全局、不区分大小写地查找 `vue`。

## 表单校验建议

不要追求一个正则解决所有场景。

手机号：

```ts
const phonePattern = /^1[3-9]\d{9}$/
```

邮箱可以用相对宽松规则：

```ts
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
```

原因：真实邮箱规则非常复杂，前端只需要挡明显错误，最终以后端校验和邮件验证为准。

## 搜索高亮

```ts
function highlight(text: string, keyword: string) {
  if (!keyword.trim()) return text

  const pattern = new RegExp(`(${escapeRegExp(keyword)})`, 'gi')

  return text.replace(pattern, '<mark>$1</mark>')
}
```

如果结果要插入 DOM，要注意 XSS。更安全的做法是拆成文本片段，由框架渲染，而不是直接拼 HTML。

## 实际项目常见问题

### 1. 用户输入导致正则报错

**原因**

把用户输入直接拼到 `RegExp` 里，特殊字符破坏了正则语法。

**解决方案**

使用 `escapeRegExp` 转义。

### 2. 正则校验过严导致合法数据被拒绝

**例子**

邮箱、姓名、地址、国际手机号都很复杂。

**解决方案**

前端做基本格式提示，复杂规则交给后端或专门库。

### 3. `g` 模式下 test 结果忽真忽假

带 `g` 的正则会维护 `lastIndex`：

```ts
const pattern = /vue/g

pattern.test('vue') // true
pattern.test('vue') // false
```

判断单个字符串时不要加 `g`：

```ts
/vue/.test('vue')
```

### 4. 正则太长没人敢维护

**解决方案**

- 拆成多个命名规则。
- 给复杂规则写注释和测试用例。
- 用更明确的解析函数替代正则。

### 5. 正则造成性能问题

某些复杂回溯会导致页面卡顿。用户输入参与正则时，要避免危险模式，必要时限制输入长度。

## 最佳实践

- 固定规则用正则字面量，动态规则用 `RegExp` 并转义输入。
- 前端校验负责体验，后端校验负责最终可信。
- 复杂正则必须命名、注释、加测试。
- 不用正则解析复杂嵌套结构。
- 单次判断不要使用带 `g` 的正则。
- 搜索高亮不要直接拼未转义 HTML。

## 学习检查

学完本节后，你应该能回答：

- `test`、`match`、`replace` 分别适合什么。
- 为什么动态正则要转义用户输入。
- 贪婪和非贪婪有什么区别。
- 为什么邮箱校验不应该追求极端严格。
- 带 `g` 的正则为什么可能影响连续 `test` 结果。

## 参考资料

- [MDN: Regular expressions guide](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions)
- [MDN: RegExp](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp)
- [MDN: Regular expression syntax cheat sheet](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions/Cheatsheet)

## 下一步学习

继续学习 [异步编程](/javascript/async)，把字符串处理和接口请求、表单提交结合起来。
