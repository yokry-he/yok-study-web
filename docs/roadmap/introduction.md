# 学习路线总览

## 适合谁看

适合不知道该从哪条路线开始，或者已经掌握一部分技术、需要判断下一阶段补什么的人。

这个文档站已经覆盖前端、框架、Node.js、Java、Go、数据库、DevOps、AI 工程和项目实战。路线总览的作用是把这些模块组织成可执行的成长路径，而不是让你在技术库里迷路。

## 路线选择

| 目标 | 推荐路线 | 适合阶段 |
| --- | --- | --- |
| 想成为 Vue 前端工程师 | [Vue 前端工程师路线](/roadmap/vue-frontend) | 前端入门到项目交付 |
| 想写后端 API 和服务 | [Node 后端工程师路线](/roadmap/node-backend) | 前端转后端、全栈入门 |
| 想走企业后端 | [Java 学习导览](/java/introduction) | Spring Boot、事务、JVM |
| 想走云原生后端 | [Go 学习导览](/go/introduction) | 高并发、服务端工具、云原生 |
| 想独立交付完整 Web 项目 | [全栈工程师路线](/roadmap/fullstack) | 已有前端基础 |
| 想做 AI 应用和智能工具 | [AI 工程路线](/roadmap/ai-engineering) | 已有编程和项目基础 |

## 推荐总路径

```text
前端基础
↓
JavaScript / TypeScript
↓
Vue 或 React
↓
浏览器与网络
↓
Node.js / Java / Go
↓
数据库
↓
DevOps 与部署
↓
AI 工程
↓
项目实战和问题库
```

不要把这条路径理解成必须一次学完。更实际的方式是围绕一个项目迭代：

1. 先做一个前端后台页面。
2. 再补请求、权限、表单和状态。
3. 接一个 Node.js、Java 或 Go API。
4. 接数据库和 Redis。
5. 部署上线并沉淀问题。
6. 最后把 AI 能力接入某个真实业务场景。

## 学习路线地图

<TechGrid :items="[
  { title: '阅读顺序与使用方法', description: '说明本站如何按路线、图解、项目、练习、问题库和速查手册组合使用。', link: '/roadmap/reading-guide', level: '导读' },
  { title: '学习工作流与笔记模板', description: '把读导览、看图解、做实验、接项目、查问题和写复盘整理成一套可执行工作流。', link: '/roadmap/study-workflow', level: '方法' },
  { title: '图解学习地图', description: '把浏览器、JavaScript、TypeScript、Vue、工程化、后端、数据库、DevOps 和 AI 的图解串成一条学习主线。', link: '/roadmap/visual-learning-map', level: '图解' },
  { title: 'Vue 前端工程师路线', description: '从 HTML/CSS、JavaScript、TypeScript 到 Vue 3、权限、工程化和 Vue Admin 实战。', link: '/roadmap/vue-frontend', level: '前端' },
  { title: 'Node 后端工程师路线', description: '从 JavaScript、HTTP、Node.js 到数据库、接口设计、日志、部署和排错。', link: '/roadmap/node-backend', level: '后端' },
  { title: 'Java 后端学习入口', description: '从 JDK、面向对象、JVM、Spring Boot 到事务、测试、部署和线上排错。', link: '/java/introduction', level: '后端' },
  { title: 'Go 后端学习入口', description: '从 Go Modules、接口组合、并发、Context、HTTP 到数据库、部署和性能诊断。', link: '/go/introduction', level: '后端' },
  { title: '全栈工程师路线', description: '串联前端、后端语言、数据库、DevOps 和项目实战，目标是独立交付完整应用。', link: '/roadmap/fullstack', level: '全栈' },
  { title: 'AI 工程路线', description: '从 LLM API、提示词、RAG、Agent 到评测、成本、安全和上线治理。', link: '/roadmap/ai-engineering', level: 'AI' },
  { title: '阶段任务清单', description: '把学习路线拆成可执行任务、阶段产出、验收标准和常见卡点。', link: '/roadmap/phase-tasks', level: '任务' },
  { title: '学习路径练习包', description: '用静态页面、JS 列表、TS 类型、Vue Admin、数据库、后端 API、上线和复盘练习验证能力。', link: '/roadmap/practice-labs', level: '练习' },
  { title: '前端基础专项练习', description: '用 10 个练习掌握语义 HTML、表单、图片、键盘操作、渐进增强、故障注入和生产验收。', link: '/roadmap/frontend-foundation-practice', level: '基础' },
  { title: '前端综合实战练习', description: '用一个 Vue Admin 工作台串联 CSS、浏览器、工程化、Vue、问题库和交付验收。', link: '/roadmap/frontend-capstone-lab', level: '综合' },
  { title: 'Nuxt / Next 专项练习', description: '用课程内容平台训练路由、SSR、数据边界、Cookie 会话、缓存、SEO、排错和部署。', link: '/roadmap/meta-framework-practice', level: '练习' },
  { title: 'Vue Admin 学习地图', description: '把 Vue Admin 文档组织成从入门、项目、权限、问题库到交付验收的完整路线。', link: '/roadmap/vue-admin-learning-map', level: '地图' },
  { title: 'Vue Admin 专项练习', description: '用 14 天计划完成 Vue Admin 用户管理模块，覆盖路由、Pinia、请求、表单、权限、测试和复盘。', link: '/roadmap/vue-admin-practice', level: '练习' },
  { title: '项目里程碑', description: '用静态页面、Vue Admin、Node API、全栈后台、内容站和 AI 助手验证能力。', link: '/roadmap/project-milestones', level: '项目' },
  { title: '能力自测', description: '按前端、Vue、Node、数据库、工程化和 AI 工程评估是否能进入下一阶段。', link: '/roadmap/self-assessment', level: '自测' },
  { title: '路线维护规则', description: '定义路线内容如何跟随模块成熟度持续更新，避免过时和重复。', link: '/roadmap/roadmap-governance', level: '治理' }
]" />

## 阶段验收

| 阶段 | 能力结果 |
| --- | --- |
| 前端基础 | 能写可维护页面，理解浏览器、网络、缓存和登录态 |
| 框架开发 | 能用 Vue 或 React 组织组件、路由、状态和表单 |
| 后端与数据 | 能使用 Node.js、Java 或 Go 写 API，设计表结构，处理事务和慢查询 |
| 工程化部署 | 能构建、部署、回滚、排查线上问题 |
| AI 工程 | 能把模型能力接入业务，并评估效果和风险 |

## 实际项目建议

最推荐的练习项目是后台管理系统，因为它能覆盖：

- 登录态。
- 动态菜单。
- 按钮权限。
- 表格筛选。
- 表单校验。
- 请求封装。
- 数据库建模。
- 部署上线。
- 真实问题排查。

如果想加入 AI 能力，可以做“后台运营助手”或“文档问答助手”，风险比自动审批、自动删除数据低。

## 下一步学习

如果你是第一次使用本站，先看 [阅读顺序与使用方法](/roadmap/reading-guide)、[学习工作流与笔记模板](/roadmap/study-workflow) 和 [图解学习地图](/roadmap/visual-learning-map)。如果 HTML、表单、图片和键盘操作还不稳定，先完成 [前端基础专项练习](/roadmap/frontend-foundation-practice)。如果你以 Vue 为主线，继续进入 [Vue 前端工程师路线](/roadmap/vue-frontend)。如果你目标是完成后台项目，进入 [Vue Admin 学习地图与交付清单](/roadmap/vue-admin-learning-map)。如果你已经有前端基础，可以直接进入 [全栈工程师路线](/roadmap/fullstack)。如果你想按任务推进，进入 [阶段任务清单](/roadmap/phase-tasks) 和 [学习路径练习包](/roadmap/practice-labs)；完成基础练习后，用 [前端综合实战练习](/roadmap/frontend-capstone-lab) 把 CSS、浏览器、工程化、Vue 和问题库串成一个可交付项目。如果你的项目需要 SSR、SEO 或全栈前端能力，进入 [Nuxt / Next 专项练习](/roadmap/meta-framework-practice)。
