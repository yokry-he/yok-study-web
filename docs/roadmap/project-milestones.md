# 项目里程碑

## 适合谁看

适合想用项目驱动学习，但不知道每个项目应该做到什么程度的人。

里程碑的作用是把学习路线变成可验收的作品。每个里程碑都应该有明确产出，而不是只写“学完 Vue”“学完 Node”。

## 里程碑 1：静态页面作品

目标：证明你能把页面结构、样式和基础交互做稳定。

建议项目：

- 个人主页。
- 产品介绍页。
- 文档首页。
- 简单移动端列表页。

必须包含：

- 语义化结构。
- 响应式布局。
- 基础交互。
- 可访问的按钮和链接。
- 清晰目录结构。

验收清单：

- 移动端没有横向溢出。
- 文本不重叠。
- 图片有稳定尺寸。
- 控制台无错误。
- README 写清启动方式。

推荐入口：

- [前端基础学习导览](/frontend/introduction)
- [前端基础从零到项目](/frontend/project-from-zero)
- [HTML 与无障碍真实项目问题库](/projects/issues-html-accessibility)
- [CSS 学习导览](/css/introduction)
- [浏览器学习导览](/browser/introduction)

## 里程碑 2：Vue Admin 前端

目标：完成一个可扩展后台前端。

必须包含：

- 登录页。
- 首页布局。
- 用户列表。
- 用户编辑表单。
- 角色权限页面。
- 动态菜单。
- 按钮权限。
- 请求封装。
- 错误提示。
- 构建部署说明。

验收清单：

- 路由、状态、请求、组件分层清楚。
- 未登录自动跳转。
- 无权限不展示入口，并能处理接口 403。
- 表单校验清楚。
- 列表支持筛选、分页和加载状态。
- 项目文档说明目录和数据流。

推荐入口：

- [Vue 学习导览](/vue/introduction)
- [Vue Admin 学习地图与交付清单](/roadmap/vue-admin-learning-map)
- [权限与菜单](/vue/permission)
- [Vue Admin 实战](/projects/vue-admin)

## 里程碑 3：Node 权限 API

目标：让前端项目接入真实后端，而不是长期依赖 mock。

必须包含：

- `/health`。
- 登录接口。
- 当前用户接口。
- 用户列表。
- 用户新增和编辑。
- 角色列表。
- 权限分配。
- 操作日志。

验收清单：

- 参数校验明确。
- 401 和 403 区分。
- controller、service、repository 分层。
- 数据库迁移有注释。
- 关键操作写审计日志。
- API 错误响应统一。

推荐入口：

- [Node.js 学习导览](/node/introduction)
- [鉴权与会话](/node/auth-session)
- [数据库集成](/node/database-integration)

## 里程碑 4：完整全栈后台

目标：把前端、后端、数据库、部署串起来。

必须包含：

- Vue Admin 前端。
- Node API 后端。
- MySQL 或 PostgreSQL。
- Redis 可选。
- Nginx 反向代理。
- Docker Compose。
- CI 构建。
- 发布检查清单。

验收清单：

- 一条命令能启动本地依赖。
- 前端接口地址按环境配置。
- 数据库结构可迁移。
- 生产构建可预览。
- 二级路由刷新不 404。
- 有回滚说明。
- 有常见问题记录。

推荐入口：

- [全栈工程师路线](/roadmap/fullstack)
- [DevOps 学习导览](/devops/introduction)
- [项目交付检查清单](/projects/delivery-checklist)

## 里程碑 5：内容站或官网

目标：掌握 SSR、SSG、SEO、内容模型和多语言。

建议项目：

- 技术博客。
- 产品官网。
- 文档站。
- 多语言内容站。

必须包含：

- 首页。
- 文章列表。
- 文章详情。
- 分类或标签。
- SEO metadata。
- sitemap。
- 分享图。
- 部署缓存策略。

验收清单：

- 每个详情页有独立 title 和 description。
- sitemap 不包含草稿和后台页面。
- `index.html` 缓存策略正确。
- 多语言页面有稳定路径。
- 内容更新后缓存能失效。

推荐入口：

- [Nuxt / Next 元框架学习导览](/meta-frameworks/introduction)
- [SEO、Metadata 与结构化数据](/meta-frameworks/seo-metadata)
- [内容站案例](/meta-frameworks/content-site-case)

## 里程碑 6：AI 文档问答助手

目标：把 AI 能力接入真实业务，并建立评测和安全边界。

必须包含：

- 文档导入。
- 文档切分。
- 检索。
- 模型回答。
- 来源引用。
- 权限过滤。
- 评测样例。
- 调用日志。
- 成本记录。

验收清单：

- 资料不足时能拒答。
- 回答能引用来源。
- 不越权回答用户不可见文档。
- prompt 有版本。
- 评测样例能重复运行。
- 成本和延迟可观察。

推荐入口：

- [AI 工程路线](/roadmap/ai-engineering)
- [RAG 检索增强生成](/ai-engineering/rag)
- [评测与质量保障](/ai-engineering/evaluation)

## 作品集记录模板

每个项目建议记录：

```md
## 项目名称

### 解决什么问题

### 技术栈

### 核心功能

### 项目结构

### 关键难点

### 如何启动

### 如何部署

### 遇到的问题和解决方案
```

作品集不是堆截图。能解释架构、边界、问题和取舍，才说明这个项目真的属于你。

## 下一步学习

如果里程碑做不出来，先回到 [学习路径练习包](/roadmap/practice-labs) 补对应练习。如果项目已经能跑，继续进入 [能力自测](/roadmap/self-assessment)，判断每个里程碑是否真正达标。
