# 正则速查

## 基础写法

| 写法 | 含义 |
| --- | --- |
| `/abc/` | 匹配 abc |
| `/abc/i` | 忽略大小写 |
| `/abc/g` | 全局匹配 |
| `/abc/m` | 多行模式 |
| `test()` | 判断是否匹配 |
| `match()` | 返回匹配结果 |
| `replace()` | 替换 |
| `split()` | 分割 |

示例：

```ts
/^\\d+$/.test('123')
```

## 字符类

| 写法 | 含义 |
| --- | --- |
| `.` | 任意字符，通常不含换行 |
| `\\d` | 数字 |
| `\\D` | 非数字 |
| `\\w` | 字母、数字、下划线 |
| `\\W` | 非 `\\w` |
| `\\s` | 空白字符 |
| `\\S` | 非空白字符 |
| `[abc]` | a、b、c 任意一个 |
| `[^abc]` | 非 a、b、c |
| `[a-z]` | 小写字母范围 |

## 数量词

| 写法 | 含义 |
| --- | --- |
| `*` | 0 次或多次 |
| `+` | 1 次或多次 |
| `?` | 0 次或 1 次 |
| `{3}` | 正好 3 次 |
| `{2,5}` | 2 到 5 次 |
| `{2,}` | 至少 2 次 |
| `+?` | 非贪婪匹配 |

示例：

```ts
'<span>hello</span>'.match(/<.*?>/)
```

## 边界

| 写法 | 含义 |
| --- | --- |
| `^` | 字符串开头 |
| `$` | 字符串结尾 |
| `\\b` | 单词边界 |
| `(?=x)` | 后面是 x |
| `(?!x)` | 后面不是 x |

校验整段字符串时，通常要同时使用 `^` 和 `$`。

```ts
/^1\\d{10}$/.test('13800000000')
```

## 分组和捕获

| 写法 | 含义 |
| --- | --- |
| `(abc)` | 捕获分组 |
| `(?:abc)` | 非捕获分组 |
| `(a|b)` | a 或 b |
| `(?<name>\\w+)` | 命名捕获 |
| `$1` | replace 中引用第 1 组 |

示例：

```ts
'2026-07-02'.replace(/(\\d{4})-(\\d{2})-(\\d{2})/, '$1/$2/$3')
```

## 常用场景

手机号基础校验：

```ts
/^1\\d{10}$/
```

邮箱基础校验：

```ts
/^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$/
```

去除首尾空白：

```ts
text.replace(/^\\s+|\\s+$/g, '')
```

提取 URL 参数片段：

```ts
/[?&]id=([^&]+)/
```

## 项目注意事项

- 正则适合格式校验，不适合解析复杂 HTML。
- 用户名、手机号、邮箱的最终合法性应以业务和后端校验为准。
- 正则过复杂时要拆分并写注释。
- 表单校验要给用户可理解的错误提示。
- 注意灾难性回溯，避免在长文本上使用危险正则。

## 参考资料

- [MDN Regular expressions guide](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions)
- [MDN RegExp](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp)
- [MDN Regular expression syntax cheat sheet](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_expressions/Cheatsheet)

## 延伸学习

- [正则表达式](/javascript/regular-expressions)
- [JavaScript 速查](/cheatsheets/javascript)
