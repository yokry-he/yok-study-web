# JavaScript 学习导览

## 适合谁看

适合已经能写一点页面交互，但对数据类型、函数、作用域、数组对象、DOM 事件、异步、错误处理和模块化还没有系统理解的人。

JavaScript 是前端、Node.js、工程脚本和很多自动化工具的共同基础。Vue 和 React 都能帮你组织 UI，但真正决定代码质量的，仍然是你对 JavaScript 运行机制和数据处理方式的理解。

## 你会学到什么

- 基础语法和常见运行规则。
- 数据类型、空值、真假值和类型判断。
- 函数、作用域、闭包、this 和原型链的边界。
- 数组、对象和树结构的项目处理方式。
- DOM 事件、事件委托、监听清理和常见交互问题。
- 正则表达式在校验、搜索、替换和日志解析里的用法。
- Promise、async/await、事件循环和并发请求。
- 错误处理、请求失败、全局兜底和日志定位。
- 内存管理、组件卸载清理、缓存边界和泄漏排查。
- ESM 模块化和项目目录职责。

## 学习顺序

```text
图解 JavaScript 核心概念
↓
JavaScript 基础
↓
数据类型与判断
↓
函数、作用域与闭包
↓
原型与原型链
↓
数组与对象处理
↓
DOM 事件
↓
正则表达式
↓
异步编程
↓
事件循环
↓
错误处理
↓
内存管理
↓
模块化与工程实践
↓
项目落地实践
↓
任务看板从零到项目
```

## 章节地图

| 章节 | 解决的问题 |
| --- | --- |
| [图解 JavaScript 核心概念](/javascript/visual-guide) | 用图理解执行上下文、闭包、原型链、事件循环、Promise 和模块加载 |
| [JavaScript 基础](/javascript/fundamentals) | 变量、条件、循环、函数和基础语法 |
| [数据类型与判断](/javascript/types) | 类型判断、空值处理、接口数据规范化 |
| [函数、作用域与闭包](/javascript/functions-scope) | 函数设计、闭包、作用域和 this |
| [原型与原型链](/javascript/prototype-chain) | 对象方法查找、class 底层模型、原型污染和属性判断 |
| [数组与对象处理](/javascript/array-object) | 列表转换、表单复制、树结构处理 |
| [DOM 事件](/javascript/dom-events) | 事件监听、冒泡捕获、事件委托、组件卸载清理 |
| [正则表达式](/javascript/regular-expressions) | 表单校验、搜索高亮、文本替换、日志解析和动态转义 |
| [异步编程](/javascript/async) | Promise、async/await、并发和错误处理 |
| [事件循环](/javascript/event-loop) | 同步代码、任务、微任务、渲染时机和页面卡顿 |
| [错误处理](/javascript/error-handling) | try/catch、请求错误、HTTP 状态、全局兜底和日志 |
| [内存管理](/javascript/memory-management) | 垃圾回收、引用保留、事件定时器清理、内存泄漏排查 |
| [模块化与工程实践](/javascript/modules) | ESM、目录职责、导入路径和副作用 |
| [项目落地实践](/javascript/project-practice) | 把接口转换、列表筛选、表单复制、异步请求、事件清理和权限判断放进真实页面 |
| [任务看板从零到项目](/javascript/task-board-project) | 用原生 JavaScript 串联状态、DOM、事件委托、localStorage、异步加载、错误处理和模块拆分 |

## 实际项目建议

学习 JavaScript 不要只刷语法题。建议结合真实项目练习：

- 把接口返回值转换成页面需要的数据结构。
- 把表单默认值封装成函数。
- 把复杂筛选拆成多个命名判断函数。
- 把列表点击、弹窗外部点击和快捷键监听封装成可清理的逻辑。
- 把复杂正则命名并写测试，不在页面里散落魔法字符串。
- 用 async/await 处理 loading、error、finally。
- 用事件循环理解 loading 不显示、DOM 更新延迟和页面卡顿。
- 把请求错误、业务错误和日志上报分层处理。
- 组件卸载时清理事件、定时器、Socket、Observer、图表实例。
- 用模块边界区分 API、service、utils、store。
- 做一个任务看板，把状态、渲染、事件、存储和错误处理串成完整闭环。

## 常见误区

### 只会写页面，不会处理数据

后台管理系统里大量代码都在处理列表、表单、权限、树结构和接口响应。如果数组和对象处理不扎实，模板会越来越乱。

### 把异步错误吞掉

不处理错误会导致用户只看到页面没反应。所有请求都应该考虑失败、加载中和最终清理。

### 只背事件循环题，不会排查卡顿

事件循环不是面试题专用概念。真实项目中，loading 不显示、输入卡顿、大列表阻塞、DOM 更新时机都和主线程、任务和微任务有关。

### 复制正则，不知道边界

表单校验、搜索高亮和日志解析都可能用到正则。复杂正则必须命名、注释和测试，动态正则必须转义用户输入。

### 以为自动 GC 就不会泄漏

JavaScript 会自动垃圾回收，但事件监听、定时器、全局缓存、闭包和第三方实例仍然可能长期持有引用。

### utils 变成杂物间

只有通用、无业务状态、无副作用的函数才适合放到 utils。业务流程应该放到 service 或 composable。

## 最佳实践

- 写函数前先明确输入和输出。
- 外部接口数据进入页面前先规范化。
- 异步请求必须处理 loading、error 和 finally。
- 复杂同步计算要关注主线程阻塞，必要时分页、分片或使用 Worker。
- 错误处理要分层，用户提示和技术日志不要混在一起。
- DOM 事件、定时器、Observer、Socket 必须有生命周期清理点。
- 复杂数据转换放在命名函数里，不堆在模板中。
- 模块目录按职责划分，不按“方便”堆文件。

## 下一步学习

第一次进入 JavaScript 模块，建议先看 [图解 JavaScript 核心概念](/javascript/visual-guide)，再学习 [JavaScript 基础](/javascript/fundamentals)。学完 DOM、异步和模块化后，继续做 [任务看板从零到项目](/javascript/task-board-project)，再进入 [Vue 从零到项目落地](/vue/project-from-zero)。
