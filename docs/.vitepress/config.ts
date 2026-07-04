import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'zh-CN',
  title: '程序员技术学习站',
  description: '面向程序员的系统化技术学习路线、工程实践、实战项目和问题库。',
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ['link', { rel: 'shortcut icon', href: '/favicon.ico' }],
    ['link', { rel: 'icon', href: '/logo.svg', type: 'image/svg+xml' }],
    ['meta', { name: 'theme-color', content: '#7edfc6' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0' }]
  ],
  markdown: {
    config(md) {
      const defaultFence = md.renderer.rules.fence?.bind(md.renderer.rules)

      md.renderer.rules.fence = (tokens, idx, options, env, self) => {
        const token = tokens[idx]
        const language = token.info.trim().split(/\s+/)[0]

        if (language === 'mermaid') {
          const encoded = Buffer.from(token.content, 'utf8').toString('base64')

          return `<MermaidDiagram encoded="${encoded}" />`
        }

        return defaultFence ? defaultFence(tokens, idx, options, env, self) : self.renderToken(tokens, idx, options)
      }
    }
  },
  themeConfig: {
    logo: '/logo.svg',
    siteTitle: '技术学习站',
    search: {
      provider: 'local',
      options: {
        async _render(src, env, md) {
          const searchableContent = src
            .replace(/```[\s\S]*?```/g, '')
            .split('\n')
            .filter((line) => {
              const trimmed = line.trim()

              if (trimmed.startsWith('|')) return false
              if (/^<\/?[A-Z][\w-]*/.test(trimmed)) return false

              return true
            })
            .join('\n')
          const html = md.renderAsync
            ? await md.renderAsync(searchableContent, env)
            : md.render(searchableContent, env)

          return env.frontmatter?.search === false ? '' : html
        },
        locales: {
          root: {
            translations: {
              button: {
                buttonText: '搜索文档',
                buttonAriaLabel: '搜索文档'
              },
              modal: {
                displayDetails: '显示详情',
                resetButtonTitle: '清除搜索',
                backButtonTitle: '关闭搜索',
                noResultsText: '没有找到相关内容',
                footer: {
                  selectText: '选择',
                  navigateText: '切换',
                  closeText: '关闭'
                }
              }
            }
          }
        }
      }
    },
    nav: [
      { text: '学习路线', link: '/roadmap/introduction' },
      { text: '技术库', link: '/technologies/' },
      {
        text: '前端基础',
        items: [
          { text: 'HTML 与 CSS', link: '/frontend/html-css' },
          { text: 'CSS', link: '/css/introduction' },
          { text: 'JavaScript', link: '/javascript/fundamentals' },
          { text: 'TypeScript', link: '/typescript/introduction' },
          { text: '浏览器与网络', link: '/browser/introduction' }
        ]
      },
      {
        text: '框架',
        items: [
          { text: 'Vue', link: '/vue/introduction' },
          { text: 'React', link: '/react/introduction' },
          { text: 'Nuxt / Next 元框架', link: '/meta-frameworks/introduction' }
        ]
      },
      {
        text: '后端与数据',
        items: [
          { text: 'Node.js', link: '/node/introduction' },
          { text: 'Java', link: '/java/introduction' },
          { text: 'Go', link: '/go/introduction' },
          { text: '数据库', link: '/database/introduction' },
          { text: 'AI 工程', link: '/ai-engineering/introduction' }
        ]
      },
      {
        text: '工程化',
        items: [
          { text: '前端工程化', link: '/engineering/introduction' },
          { text: 'DevOps 与部署', link: '/devops/introduction' }
        ]
      },
      { text: '实战项目', link: '/projects/vue-admin' }
    ],
    sidebar: {
      '/technologies/': [
        {
          text: '技术库',
          items: [
            { text: '技术库总览', link: '/technologies/' },
            { text: '前端技术图谱', link: '/technologies/frontend-map' },
            { text: '内容扩展路线', link: '/technologies/expansion-plan' }
          ]
        }
      ],
      '/roadmap/': [
        {
          text: '学习路线',
          items: [
            { text: '学习路线总览', link: '/roadmap/introduction' },
            { text: '阅读顺序与使用方法', link: '/roadmap/reading-guide' },
            { text: '图解学习地图', link: '/roadmap/visual-learning-map' },
            { text: 'Vue 前端工程师路线', link: '/roadmap/vue-frontend' },
            { text: 'Node 后端工程师路线', link: '/roadmap/node-backend' },
            { text: '全栈工程师路线', link: '/roadmap/fullstack' },
            { text: 'AI 工程路线', link: '/roadmap/ai-engineering' },
            { text: '阶段任务清单', link: '/roadmap/phase-tasks' },
            { text: '学习路径练习包', link: '/roadmap/practice-labs' },
            { text: 'Vue Admin 学习地图', link: '/roadmap/vue-admin-learning-map' },
            { text: 'Vue Admin 专项练习', link: '/roadmap/vue-admin-practice' },
            { text: '项目里程碑', link: '/roadmap/project-milestones' },
            { text: '能力自测', link: '/roadmap/self-assessment' },
            { text: '路线维护规则', link: '/roadmap/roadmap-governance' }
          ]
        }
      ],
      '/frontend/': [
        {
          text: '前端基础',
          items: [
            { text: 'HTML 与 CSS', link: '/frontend/html-css' },
            { text: 'TypeScript 基础', link: '/frontend/typescript' },
            { text: '浏览器与网络', link: '/browser/introduction' }
          ]
        }
      ],
      '/browser/': [
        {
          text: '浏览器与网络',
          items: [
            { text: '浏览器学习导览', link: '/browser/introduction' },
            { text: '图解浏览器核心概念', link: '/browser/visual-guide' },
            { text: 'HTTP 与请求流程', link: '/browser/http-request' },
            { text: '跨域与登录态', link: '/browser/cors-auth' },
            { text: '缓存策略', link: '/browser/cache' },
            { text: '浏览器存储', link: '/browser/storage' },
            { text: '浏览器安全基础', link: '/browser/security' },
            { text: 'Service Worker 与 PWA', link: '/browser/service-worker-pwa' },
            { text: '常用 Web API', link: '/browser/web-apis' },
            { text: 'WebSocket 实时通信', link: '/browser/websocket' },
            { text: 'WebRTC 实时音视频', link: '/browser/webrtc' },
            { text: 'Web Components', link: '/browser/web-components' },
            { text: 'WebAssembly', link: '/browser/webassembly' },
            { text: 'WebGPU', link: '/browser/webgpu' },
            { text: '浏览器自动化调试', link: '/browser/browser-automation-debugging' },
            { text: '渲染与性能', link: '/browser/rendering-performance' },
            { text: '常见问题', link: '/browser/troubleshooting' }
          ]
        }
      ],
      '/javascript/': [
        {
          text: 'JavaScript',
          items: [
            { text: 'JavaScript 学习导览', link: '/javascript/introduction' },
            { text: '图解 JavaScript 核心概念', link: '/javascript/visual-guide' },
            { text: 'JavaScript 基础', link: '/javascript/fundamentals' },
            { text: '数据类型与判断', link: '/javascript/types' },
            { text: '函数、作用域与闭包', link: '/javascript/functions-scope' },
            { text: '原型与原型链', link: '/javascript/prototype-chain' },
            { text: '数组与对象处理', link: '/javascript/array-object' },
            { text: 'DOM 事件', link: '/javascript/dom-events' },
            { text: '正则表达式', link: '/javascript/regular-expressions' },
            { text: '异步编程', link: '/javascript/async' },
            { text: '事件循环', link: '/javascript/event-loop' },
            { text: '错误处理', link: '/javascript/error-handling' },
            { text: '内存管理', link: '/javascript/memory-management' },
            { text: '模块化与工程实践', link: '/javascript/modules' },
            { text: '项目落地实践', link: '/javascript/project-practice' },
            { text: '任务看板从零到项目', link: '/javascript/task-board-project' },
            { text: 'JavaScript 真实项目问题库', link: '/projects/issues-javascript' }
          ]
        }
      ],
      '/typescript/': [
        {
          text: 'TypeScript',
          items: [
            { text: 'TypeScript 学习导览', link: '/typescript/introduction' },
            { text: '图解 TypeScript 核心概念', link: '/typescript/visual-guide' },
            { text: '基础类型', link: '/typescript/basic-types' },
            { text: '对象、接口与 type', link: '/typescript/interface-type' },
            { text: '泛型', link: '/typescript/generics' },
            { text: '类型收窄与类型守卫', link: '/typescript/narrowing-guards' },
            { text: '工具类型与类型边界', link: '/typescript/utility-types-boundary' },
            { text: 'tsconfig 与工程配置', link: '/typescript/tsconfig-engineering' },
            { text: 'Vue 项目集成', link: '/typescript/vue-integration' },
            { text: '项目落地实践', link: '/typescript/project-practice' },
            { text: '类型边界从零到项目', link: '/typescript/type-boundary-project' },
            { text: '常见问题', link: '/typescript/troubleshooting' }
          ]
        }
      ],
      '/css/': [
        {
          text: 'CSS',
          items: [
            { text: 'CSS 学习导览', link: '/css/introduction' },
            { text: '图解 CSS 核心概念', link: '/css/visual-guide' },
            { text: '盒模型与布局基础', link: '/css/box-model-layout' },
            { text: 'Flex 与 Grid', link: '/css/flex-grid' },
            { text: '响应式设计', link: '/css/responsive' },
            { text: '动画与过渡', link: '/css/animation-transition' },
            { text: 'CSS 可访问性', link: '/css/accessibility' },
            { text: '设计 Token 与主题', link: '/css/design-tokens' },
            { text: '项目样式架构', link: '/css/architecture' },
            { text: 'CSS 真实项目问题库', link: '/projects/issues-css' },
            { text: '常见问题', link: '/css/troubleshooting' }
          ]
        }
      ],
      '/react/': [
        {
          text: 'React',
          items: [
            { text: 'React 学习导览', link: '/react/introduction' },
            { text: '图解 React 核心概念', link: '/react/visual-guide' },
            { text: '快速开始', link: '/react/quick-start' },
            { text: '组件与 JSX', link: '/react/component-jsx' },
            { text: 'Hooks 与状态', link: '/react/hooks-state' },
            { text: 'Effect 与副作用', link: '/react/effects' },
            { text: '表单处理', link: '/react/forms' },
            { text: '请求与数据流', link: '/react/request-data-flow' },
            { text: 'Context 与状态管理', link: '/react/context-state-management' },
            { text: '路由与项目结构', link: '/react/router-structure' },
            { text: '性能优化', link: '/react/performance' },
            { text: '测试策略', link: '/react/testing' },
            { text: '最佳实践', link: '/react/best-practices' },
            { text: 'React 管理台从零到项目', link: '/react/project-admin' },
            { text: '常见问题', link: '/react/troubleshooting' }
          ]
        }
      ],
      '/meta-frameworks/': [
        {
          text: 'Nuxt / Next 元框架',
          items: [
            { text: '元框架学习导览', link: '/meta-frameworks/introduction' },
            { text: 'Nuxt 项目实践', link: '/meta-frameworks/nuxt' },
            { text: 'Next.js 项目实践', link: '/meta-frameworks/next' },
            { text: '路由、布局与数据获取', link: '/meta-frameworks/routing-data' },
            { text: '部署、缓存与运行时', link: '/meta-frameworks/deployment' },
            { text: '服务端鉴权与登录态', link: '/meta-frameworks/server-auth' },
            { text: 'SEO、Metadata 与结构化数据', link: '/meta-frameworks/seo-metadata' },
            { text: '国际化与多语言站点', link: '/meta-frameworks/i18n' },
            { text: '内容站案例', link: '/meta-frameworks/content-site-case' },
            { text: '常见问题', link: '/meta-frameworks/troubleshooting' }
          ]
        }
      ],
      '/node/': [
        {
          text: 'Node.js',
          items: [
            { text: 'Node.js 学习导览', link: '/node/introduction' },
            { text: '图解 Node.js 核心概念', link: '/node/visual-guide' },
            { text: '运行时与事件循环', link: '/node/runtime-event-loop' },
            { text: '包管理与模块化', link: '/node/package-modules' },
            { text: 'HTTP API 开发', link: '/node/http-api' },
            { text: '鉴权与会话', link: '/node/auth-session' },
            { text: '数据库集成', link: '/node/database-integration' },
            { text: '错误处理与日志', link: '/node/error-logging' },
            { text: '测试策略', link: '/node/testing' },
            { text: 'Node.js 安全基础', link: '/node/security' },
            { text: '项目结构与部署', link: '/node/project-deployment' },
            { text: 'Node 权限 API 从零到项目', link: '/node/permission-api-project' },
            { text: '常见问题', link: '/node/troubleshooting' }
          ]
        }
      ],
      '/database/': [
        {
          text: '数据库',
          items: [
            { text: '数据库学习导览', link: '/database/introduction' },
            { text: '图解数据库核心概念', link: '/database/visual-guide' },
            { text: 'MySQL 入门与项目实践', link: '/database/mysql' },
            { text: 'PostgreSQL 入门与项目实践', link: '/database/postgresql' },
            { text: 'Redis 缓存与数据结构', link: '/database/redis' },
            { text: '数据库项目落地实践', link: '/database/project-practice' },
            { text: '数据建模与表设计', link: '/database/modeling' },
            { text: '索引与查询优化', link: '/database/indexes' },
            { text: '事务、锁与并发', link: '/database/transactions' },
            { text: '迁移、种子与版本治理', link: '/database/migration' },
            { text: 'ORM 实战', link: '/database/orm-practice' },
            { text: '备份与恢复', link: '/database/backup-recovery' },
            { text: '数据安全、审计与脱敏', link: '/database/security-audit' },
            { text: '常见问题', link: '/database/troubleshooting' }
          ]
        }
      ],
      '/java/': [
        {
          text: 'Java',
          items: [
            { text: 'Java 学习导览', link: '/java/introduction' },
            { text: '图解 Java 核心概念', link: '/java/visual-guide' },
            { text: '环境、JDK 与构建工具', link: '/java/setup-tooling' },
            { text: '语法与面向对象', link: '/java/syntax-oop' },
            { text: '集合、泛型与常用类库', link: '/java/collections-generics' },
            { text: '异常、日志与编码规范', link: '/java/exceptions-logging' },
            { text: 'Stream、Lambda 与数据处理', link: '/java/streams-lambda' },
            { text: '并发、线程池与虚拟线程', link: '/java/concurrency-virtual-threads' },
            { text: 'JVM 内存、GC 与诊断', link: '/java/jvm-memory-gc' },
            { text: 'Spring Boot API 开发', link: '/java/spring-boot-api' },
            { text: '数据库、事务与 ORM', link: '/java/persistence-transaction' },
            { text: '测试、打包与部署', link: '/java/testing-deployment' },
            { text: '常见问题', link: '/java/troubleshooting' }
          ]
        }
      ],
      '/go/': [
        {
          text: 'Go',
          items: [
            { text: 'Go 学习导览', link: '/go/introduction' },
            { text: '图解 Go 核心概念', link: '/go/visual-guide' },
            { text: '环境、模块与工作区', link: '/go/setup-modules' },
            { text: '语法、类型与函数', link: '/go/syntax-types' },
            { text: '接口、组合与项目建模', link: '/go/interfaces-composition' },
            { text: '错误处理、日志与配置', link: '/go/errors-logging-config' },
            { text: '并发：goroutine、channel、select', link: '/go/concurrency' },
            { text: 'Context、HTTP 服务与中间件', link: '/go/context-http' },
            { text: '数据库、事务与仓储层', link: '/go/database-transaction' },
            { text: '测试、Benchmark 与 Fuzzing', link: '/go/testing' },
            { text: '项目结构、构建与部署', link: '/go/project-deployment' },
            { text: '性能分析与线上诊断', link: '/go/performance' },
            { text: '常见问题', link: '/go/troubleshooting' }
          ]
        }
      ],
      '/ai-engineering/': [
        {
          text: 'AI 工程',
          items: [
            { text: 'AI 工程学习导览', link: '/ai-engineering/introduction' },
            { text: '图解 AI 工程核心概念', link: '/ai-engineering/visual-guide' },
            { text: 'LLM API 调用', link: '/ai-engineering/llm-api' },
            { text: '提示词工程', link: '/ai-engineering/prompt-engineering' },
            { text: '结构化输出与函数调用', link: '/ai-engineering/structured-outputs-tools' },
            { text: '多模态 AI 应用', link: '/ai-engineering/multimodal' },
            { text: 'RAG 检索增强生成', link: '/ai-engineering/rag' },
            { text: 'MCP 与企业工具集成', link: '/ai-engineering/mcp-integration' },
            { text: 'Agent 工作流', link: '/ai-engineering/agents' },
            { text: 'AI 产品设计与人机协作', link: '/ai-engineering/product-workflow' },
            { text: '评测与质量保障', link: '/ai-engineering/evaluation' },
            { text: 'AI 文档问答从零到项目', link: '/ai-engineering/doc-qa-project' },
            { text: '上线、成本与安全', link: '/ai-engineering/deployment' },
            { text: '常见问题', link: '/ai-engineering/troubleshooting' }
          ]
        }
      ],
      '/vue/': [
        {
          text: 'Vue 核心',
          items: [
            { text: 'Vue 学习导览', link: '/vue/introduction' },
            { text: '图解 Vue 核心概念', link: '/vue/visual-guide' },
            { text: '快速开始', link: '/vue/quick-start' },
            { text: '模板语法', link: '/vue/template-syntax' },
            { text: '响应式基础', link: '/vue/reactivity' },
            { text: '组件设计', link: '/vue/component' },
            { text: '组合式 API', link: '/vue/composition-api' },
            { text: '生命周期', link: '/vue/lifecycle' },
            { text: '路由与页面', link: '/vue/router' },
            { text: 'Pinia 状态管理', link: '/vue/pinia' },
            { text: '表单处理', link: '/vue/forms' },
            { text: '请求与接口封装', link: '/vue/request' },
            { text: '权限与菜单', link: '/vue/permission' },
            { text: '内置组件', link: '/vue/built-ins' },
            { text: '性能优化', link: '/vue/performance' },
            { text: '测试策略', link: '/vue/testing' },
            { text: '最佳实践', link: '/vue/best-practices' },
            { text: '从零到项目落地', link: '/vue/project-from-zero' },
            { text: 'Vue Admin 阅读索引', link: '/vue/admin-reading-guide' },
            { text: '图解 Vue Admin 架构', link: '/vue/admin-architecture-visual-guide' },
            { text: 'Vue Admin Mock 到真实接口', link: '/vue/admin-mock-to-api' },
            { text: 'Vue Admin 列表搜索表格', link: '/vue/admin-list-search-table' },
            { text: 'Vue Admin 表单新增编辑', link: '/vue/admin-form-modal-crud' },
            { text: 'Vue Admin 详情状态记录', link: '/vue/admin-detail-status-audit' },
            { text: 'Vue Admin 文件导入导出', link: '/vue/admin-file-import-export' },
            { text: 'Vue Admin 工作台看板', link: '/vue/admin-dashboard-analytics' },
            { text: 'Vue Admin 审批流闭环', link: '/vue/admin-approval-workflow' },
            { text: 'Vue Admin 消息通知闭环', link: '/vue/admin-notification-center' },
            { text: 'Vue Admin 权限路由闭环', link: '/vue/admin-permission-route-flow' },
            { text: 'Vue Admin 用户模块实现', link: '/vue/admin-user-module' },
            { text: 'Vue Admin 角色权限实现', link: '/vue/admin-permission-module' },
            { text: 'Vue Admin 菜单与动态路由', link: '/vue/admin-menu-route-module' },
            { text: 'Vue Admin 组织与数据权限', link: '/vue/admin-organization-data-permission' },
            { text: 'Vue Admin 请求与错误处理', link: '/vue/admin-request-error-handling' },
            { text: 'Vue 真实项目问题库', link: '/projects/issues-vue' },
            { text: '常见问题', link: '/vue/troubleshooting' }
          ]
        }
      ],
      '/engineering/': [
        {
          text: '前端工程化',
          items: [
            { text: '工程化学习导览', link: '/engineering/introduction' },
            { text: '图解前端工程化核心概念', link: '/engineering/visual-guide' },
            { text: 'Vite 工程基础', link: '/engineering/vite' },
            { text: '代码规范', link: '/engineering/eslint-prettier' },
            { text: '环境配置', link: '/engineering/env-config' },
            { text: '依赖管理', link: '/engineering/package-management' },
            { text: '测试策略', link: '/engineering/testing' },
            { text: 'Monorepo 项目组织', link: '/engineering/monorepo' },
            { text: '组件库工程从零到项目', link: '/engineering/component-library-project' },
            { text: '构建与部署', link: '/engineering/build-deploy' },
            { text: '包体积分析', link: '/engineering/bundle-analysis' },
            { text: '模块联邦与微前端', link: '/engineering/module-federation' },
            { text: '工程性能优化', link: '/engineering/performance-optimization' },
            { text: '常见问题', link: '/engineering/troubleshooting' }
          ]
        }
      ],
      '/devops/': [
        {
          text: 'DevOps 与部署',
          items: [
            { text: 'DevOps 学习导览', link: '/devops/introduction' },
            { text: '图解 DevOps 核心概念', link: '/devops/visual-guide' },
            { text: 'Linux 与 Shell 基础', link: '/devops/linux-shell' },
            { text: 'Nginx 静态部署与代理', link: '/devops/nginx' },
            { text: 'Docker 容器化', link: '/devops/docker' },
            { text: 'CI/CD 自动化发布', link: '/devops/ci-cd' },
            { text: '项目上线全流程实践', link: '/devops/project-deployment-practice' },
            { text: '发布、回滚与环境治理', link: '/devops/deployment-strategy' },
            { text: '可观测性', link: '/devops/observability' },
            { text: 'Kubernetes 入门', link: '/devops/kubernetes-basics' },
            { text: '云服务与对象存储部署', link: '/devops/cloud-deployment' },
            { text: '常见问题', link: '/devops/troubleshooting' }
          ]
        }
      ],
      '/projects/': [
        {
          text: '项目实战',
          items: [
            { text: 'Vue Admin 实战', link: '/projects/vue-admin' },
            { text: '组件库实战', link: '/projects/component-library' },
            { text: '项目阶段任务', link: '/projects/project-stage-tasks' },
            { text: '权限系统案例', link: '/projects/permission-case-study' },
            { text: '权限运营案例', link: '/projects/permission-operation-case' },
            { text: '组织架构案例', link: '/projects/organization-case' },
            { text: '审批流案例', link: '/projects/approval-workflow-case' },
            { text: '文件中心案例', link: '/projects/file-center-case' },
            { text: '数据看板案例', link: '/projects/analytics-dashboard-case' },
            { text: '多租户权限案例', link: '/projects/multi-tenant-permission-case' },
            { text: '消息通知案例', link: '/projects/notification-center-case' },
            { text: '导入导出案例', link: '/projects/import-export-case' },
            { text: '支付订单案例', link: '/projects/payment-order-case' },
            { text: '会员订阅案例', link: '/projects/subscription-billing-case' },
            { text: '搜索中心案例', link: '/projects/search-center-case' },
            { text: '任务调度案例', link: '/projects/task-scheduler-case' },
            { text: '消息队列案例', link: '/projects/message-queue-case' },
            { text: '开放平台案例', link: '/projects/open-platform-case' },
            { text: '工作流配置器案例', link: '/projects/workflow-builder-case' },
            { text: '低代码流程平台案例', link: '/projects/low-code-workflow-case' },
            { text: '审计中心案例', link: '/projects/audit-center-case' },
            { text: '运营活动案例', link: '/projects/marketing-campaign-case' },
            { text: '财务对账案例', link: '/projects/finance-reconciliation-case' },
            { text: '渠道结算案例', link: '/projects/channel-settlement-case' },
            { text: '渠道费用稽核案例', link: '/projects/channel-expense-audit-case' },
            { text: '渠道费用 ROI 复盘案例', link: '/projects/channel-expense-roi-review-case' },
            { text: '渠道费用预算优化案例', link: '/projects/channel-expense-budget-optimization-case' },
            { text: '渠道费用异常预警案例', link: '/projects/channel-expense-anomaly-warning-case' },
            { text: '渠道费用策略灰度案例', link: '/projects/channel-expense-strategy-gray-release-case' },
            { text: '渠道策略效果复盘案例', link: '/projects/channel-strategy-effect-review-case' },
            { text: '渠道策略对照实验案例', link: '/projects/channel-strategy-ab-experiment-case' },
            { text: '渠道策略版本治理案例', link: '/projects/channel-strategy-version-governance-case' },
            { text: '渠道策略审批矩阵案例', link: '/projects/channel-strategy-approval-matrix-case' },
            { text: '渠道策略发布审计案例', link: '/projects/channel-strategy-release-audit-case' },
            { text: '渠道策略回滚治理案例', link: '/projects/channel-strategy-rollback-governance-case' },
            { text: '渠道策略异常仲裁案例', link: '/projects/channel-strategy-exception-arbitration-case' },
            { text: '渠道策略仲裁复盘案例', link: '/projects/channel-strategy-arbitration-review-case' },
            { text: '渠道策略裁决标准库案例', link: '/projects/channel-strategy-decision-standard-library-case' },
            { text: '渠道策略标准效果监控案例', link: '/projects/channel-strategy-standard-effect-monitoring-case' },
            { text: '渠道策略标准灰度发布案例', link: '/projects/channel-strategy-standard-gray-release-case' },
            { text: '渠道策略标准版本回滚案例', link: '/projects/channel-strategy-standard-version-rollback-case' },
            { text: '渠道策略标准回滚演练案例', link: '/projects/channel-strategy-standard-rollback-drill-case' },
            { text: '渠道策略标准灾备切换案例', link: '/projects/channel-strategy-standard-disaster-recovery-switch-case' },
            { text: '渠道价格稽核案例', link: '/projects/channel-price-audit-case' },
            { text: '渠道窜货监控案例', link: '/projects/channel-diversion-monitor-case' },
            { text: '渠道信用评级案例', link: '/projects/channel-credit-rating-case' },
            { text: '渠道返利风控案例', link: '/projects/channel-rebate-risk-control-case' },
            { text: '渠道政策模拟案例', link: '/projects/channel-policy-simulation-case' },
            { text: '渠道利润模拟案例', link: '/projects/channel-profit-simulation-case' },
            { text: '渠道价格弹性分析案例', link: '/projects/channel-price-elasticity-analysis-case' },
            { text: '主数据管理案例', link: '/projects/master-data-case' },
            { text: '客户主数据案例', link: '/projects/customer-master-data-case' },
            { text: '低代码表单案例', link: '/projects/low-code-form-case' },
            { text: '报表配置器案例', link: '/projects/report-builder-case' },
            { text: '智能报表与 BI 案例', link: '/projects/smart-bi-dashboard-case' },
            { text: '客服工单案例', link: '/projects/support-ticket-case' },
            { text: '客服质检案例', link: '/projects/customer-service-quality-case' },
            { text: '集团系统集成案例', link: '/projects/enterprise-integration-case' },
            { text: '国际化后台案例', link: '/projects/i18n-admin-case' },
            { text: '数据治理平台案例', link: '/projects/data-governance-case' },
            { text: '数据质量专项案例', link: '/projects/data-quality-special-case' },
            { text: '数据资产运营案例', link: '/projects/data-asset-operation-case' },
            { text: '数据安全运营案例', link: '/projects/data-security-operation-case' },
            { text: '规则引擎案例', link: '/projects/rule-engine-case' },
            { text: '灰度发布后台案例', link: '/projects/gray-release-admin-case' },
            { text: '跨区域灾备案例', link: '/projects/disaster-recovery-case' },
            { text: '风控中心案例', link: '/projects/risk-control-center-case' },
            { text: '合同管理案例', link: '/projects/contract-management-case' },
            { text: '合同履约案例', link: '/projects/contract-fulfillment-case' },
            { text: '合同付款案例', link: '/projects/contract-payment-case' },
            { text: '合同变更案例', link: '/projects/contract-change-case' },
            { text: '合同续签案例', link: '/projects/contract-renewal-case' },
            { text: '客户合同风险预警案例', link: '/projects/customer-contract-risk-warning-case' },
            { text: '客户合同收入预测案例', link: '/projects/customer-contract-revenue-forecast-case' },
            { text: '知识库平台案例', link: '/projects/knowledge-base-case' },
            { text: '客服知识运营案例', link: '/projects/customer-knowledge-operation-case' },
            { text: '统一配置中心案例', link: '/projects/config-center-case' },
            { text: '行业合规审计案例', link: '/projects/compliance-audit-case' },
            { text: '客户成功平台案例', link: '/projects/customer-success-case' },
            { text: '客户生命周期价值分析案例', link: '/projects/customer-lifetime-value-analysis-case' },
            { text: '客户流失预警案例', link: '/projects/customer-churn-warning-case' },
            { text: '客户续费挽回案例', link: '/projects/customer-renewal-recovery-case' },
            { text: '客户续约定价策略案例', link: '/projects/customer-renewal-pricing-strategy-case' },
            { text: '客户分群运营案例', link: '/projects/customer-segmentation-operation-case' },
            { text: '客户触达自动化案例', link: '/projects/customer-touch-automation-case' },
            { text: '客户权益运营案例', link: '/projects/customer-benefit-operation-case' },
            { text: '客户投诉闭环案例', link: '/projects/customer-complaint-closed-loop-case' },
            { text: '工单自动化案例', link: '/projects/ticket-automation-case' },
            { text: '计费中台案例', link: '/projects/billing-platform-case' },
            { text: '数据交换平台案例', link: '/projects/data-exchange-platform-case' },
            { text: '企业门户案例', link: '/projects/enterprise-portal-case' },
            { text: '资产管理案例', link: '/projects/asset-management-case' },
            { text: '预算管理案例', link: '/projects/budget-management-case' },
            { text: '资金计划案例', link: '/projects/cash-flow-planning-case' },
            { text: '费用报销案例', link: '/projects/expense-reimbursement-case' },
            { text: '员工借款案例', link: '/projects/employee-loan-case' },
            { text: '税务管理案例', link: '/projects/tax-management-case' },
            { text: '发票协同案例', link: '/projects/invoice-collaboration-case' },
            { text: '采购管理案例', link: '/projects/procurement-management-case' },
            { text: '采购寻源案例', link: '/projects/procurement-sourcing-case' },
            { text: '供应商准入案例', link: '/projects/supplier-onboarding-case' },
            { text: '供应商合同协同案例', link: '/projects/supplier-contract-collaboration-case' },
            { text: '供应商协同门户案例', link: '/projects/supplier-collaboration-portal-case' },
            { text: '供应商门户权限审计案例', link: '/projects/supplier-portal-permission-audit-case' },
            { text: '供应商索赔案例', link: '/projects/supplier-claim-case' },
            { text: '供应商绩效案例', link: '/projects/supplier-performance-case' },
            { text: '供应链计划案例', link: '/projects/supply-chain-planning-case' },
            { text: '项目管理案例', link: '/projects/project-management-case' },
            { text: '研发需求池案例', link: '/projects/rd-requirement-pool-case' },
            { text: '报价中心案例', link: '/projects/quotation-center-case' },
            { text: '价格审批中心案例', link: '/projects/price-approval-center-case' },
            { text: '库存管理案例', link: '/projects/inventory-management-case' },
            { text: '渠道库存协同案例', link: '/projects/channel-inventory-collaboration-case' },
            { text: '备件库存案例', link: '/projects/spare-parts-inventory-case' },
            { text: '备件补货案例', link: '/projects/spare-parts-replenishment-case' },
            { text: '售后备件周转分析案例', link: '/projects/after-sales-spare-parts-turnover-case' },
            { text: '备件旧件返修案例', link: '/projects/spare-parts-return-repair-case' },
            { text: '仓储物流案例', link: '/projects/warehouse-logistics-case' },
            { text: '售后服务案例', link: '/projects/after-sales-service-case' },
            { text: '售后远程诊断案例', link: '/projects/after-sales-remote-diagnosis-case' },
            { text: '售后专家协同案例', link: '/projects/after-sales-expert-collaboration-case' },
            { text: '售后知识自动推荐案例', link: '/projects/after-sales-knowledge-recommendation-case' },
            { text: '售后知识质量治理案例', link: '/projects/after-sales-knowledge-quality-governance-case' },
            { text: '售后知识智能检索优化案例', link: '/projects/after-sales-knowledge-search-optimization-case' },
            { text: '售后知识问答助手案例', link: '/projects/after-sales-knowledge-qa-assistant-case' },
            { text: '售后知识自动质检案例', link: '/projects/after-sales-knowledge-auto-quality-inspection-case' },
            { text: '售后知识专家审核案例', link: '/projects/after-sales-knowledge-expert-review-case' },
            { text: '售后知识发布灰度案例', link: '/projects/after-sales-knowledge-release-gray-case' },
            { text: '售后知识回滚治理案例', link: '/projects/after-sales-knowledge-rollback-governance-case' },
            { text: '售后知识影响追踪案例', link: '/projects/after-sales-knowledge-impact-trace-case' },
            { text: '售后知识客户通知治理案例', link: '/projects/after-sales-knowledge-customer-notification-governance-case' },
            { text: '售后知识外部服务商通知协同案例', link: '/projects/after-sales-knowledge-provider-notification-collaboration-case' },
            { text: '售后知识服务商培训闭环案例', link: '/projects/after-sales-knowledge-provider-training-closed-loop-case' },
            { text: '售后知识培训效果复盘案例', link: '/projects/after-sales-knowledge-training-effect-review-case' },
            { text: '售后知识培训认证治理案例', link: '/projects/after-sales-knowledge-training-certification-governance-case' },
            { text: '售后知识认证派单联动案例', link: '/projects/after-sales-knowledge-certification-dispatch-linkage-case' },
            { text: '售后知识认证质量稽核案例', link: '/projects/after-sales-knowledge-certification-quality-audit-case' },
            { text: '售后知识认证风险画像案例', link: '/projects/after-sales-knowledge-certification-risk-profile-case' },
            { text: '售后知识认证服务商整改案例', link: '/projects/after-sales-knowledge-certification-provider-rectification-case' },
            { text: '客户退换货质检案例', link: '/projects/customer-return-quality-inspection-case' },
            { text: '客户退款风控案例', link: '/projects/customer-refund-risk-control-case' },
            { text: '售后结算案例', link: '/projects/after-sales-settlement-case' },
            { text: '现场服务收费案例', link: '/projects/field-service-charging-case' },
            { text: '售后备件成本核算案例', link: '/projects/after-sales-spare-part-cost-case' },
            { text: '售后成本毛利分析案例', link: '/projects/after-sales-cost-margin-case' },
            { text: '售后服务成本优化案例', link: '/projects/after-sales-service-cost-optimization-case' },
            { text: '售后 SLA 赔付分析案例', link: '/projects/after-sales-sla-compensation-case' },
            { text: '售后服务商评级案例', link: '/projects/after-sales-provider-rating-case' },
            { text: '售后维修质量复盘案例', link: '/projects/after-sales-repair-quality-review-case' },
            { text: '售后投诉根因分析案例', link: '/projects/after-sales-complaint-root-cause-case' },
            { text: '报修派单案例', link: '/projects/repair-dispatch-case' },
            { text: '服务网点案例', link: '/projects/service-outlet-case' },
            { text: '数据权限审计案例', link: '/projects/data-permission-audit-case' },
            { text: '门店零售管理案例', link: '/projects/retail-store-management-case' },
            { text: 'CRM 销售管理案例', link: '/projects/crm-sales-management-case' },
            { text: '客户账期案例', link: '/projects/customer-credit-term-case' },
            { text: '客户授信风控案例', link: '/projects/customer-credit-risk-control-case' },
            { text: '客户回款风险预测案例', link: '/projects/customer-payment-risk-prediction-case' },
            { text: '客户坏账处置策略案例', link: '/projects/customer-bad-debt-disposal-case' },
            { text: '客户应收催收自动化案例', link: '/projects/customer-receivable-collection-automation-case' },
            { text: '销售回款预测调度案例', link: '/projects/sales-payment-prediction-scheduling-case' },
            { text: '销售现金流预警案例', link: '/projects/sales-cash-flow-warning-case' },
            { text: '销售回款策略模拟案例', link: '/projects/sales-collection-strategy-simulation-case' },
            { text: '销售风险动作编排案例', link: '/projects/sales-risk-action-orchestration-case' },
            { text: '销售风险处置复盘案例', link: '/projects/sales-risk-disposal-review-case' },
            { text: '销售风险预案演练案例', link: '/projects/sales-risk-contingency-drill-case' },
            { text: '销售风险指标治理案例', link: '/projects/sales-risk-metric-governance-case' },
            { text: '销售风险指标血缘审计案例', link: '/projects/sales-risk-metric-lineage-audit-case' },
            { text: '销售风险指标异常根因案例', link: '/projects/sales-risk-metric-anomaly-root-cause-case' },
            { text: '销售风险指标自动修复案例', link: '/projects/sales-risk-metric-auto-repair-case' },
            { text: '销售风险指标治理成熟度案例', link: '/projects/sales-risk-metric-governance-maturity-case' },
            { text: '销售风险指标治理运营看板案例', link: '/projects/sales-risk-metric-governance-operations-dashboard-case' },
            { text: '销售风险指标治理成本收益评估案例', link: '/projects/sales-risk-metric-governance-cost-benefit-evaluation-case' },
            { text: '销售风险指标治理预算审批案例', link: '/projects/sales-risk-metric-governance-budget-approval-case' },
            { text: '销售回款计划案例', link: '/projects/sales-collection-plan-case' },
            { text: '销售预测复盘案例', link: '/projects/sales-forecast-review-case' },
            { text: '销售目标拆解案例', link: '/projects/sales-target-breakdown-case' },
            { text: '销售佣金核算案例', link: '/projects/sales-commission-settlement-case' },
            { text: '销售返利政策案例', link: '/projects/sales-rebate-policy-case' },
            { text: '会员营销案例', link: '/projects/member-marketing-case' },
            { text: '生产制造案例', link: '/projects/manufacturing-execution-case' },
            { text: '生产排程案例', link: '/projects/production-scheduling-case' },
            { text: '产能负荷预测案例', link: '/projects/capacity-load-forecast-case' },
            { text: '生产计划达成分析案例', link: '/projects/production-plan-attainment-case' },
            { text: '生产停线损失复盘案例', link: '/projects/production-line-stop-loss-review-case' },
            { text: '质量追溯案例', link: '/projects/quality-traceability-case' },
            { text: '生产质量异常案例', link: '/projects/production-quality-exception-case' },
            { text: '生产异常 CAPA 案例', link: '/projects/production-exception-capa-case' },
            { text: '生产过程审核案例', link: '/projects/production-process-audit-case' },
            { text: '生产巡检移动端案例', link: '/projects/production-mobile-inspection-case' },
            { text: '生产现场安全隐患案例', link: '/projects/production-safety-hazard-case' },
            { text: '生产安全培训闭环案例', link: '/projects/production-safety-training-closed-loop-case' },
            { text: '生产安全考试认证案例', link: '/projects/production-safety-exam-certification-case' },
            { text: '生产安全风险画像案例', link: '/projects/production-safety-risk-profile-case' },
            { text: '生产安全应急演练案例', link: '/projects/production-safety-emergency-drill-case' },
            { text: '生产安全事故复盘案例', link: '/projects/production-safety-incident-review-case' },
            { text: '生产安全风险整改复查案例', link: '/projects/production-safety-risk-rectification-review-case' },
            { text: '生产安全整改看板案例', link: '/projects/production-safety-rectification-dashboard-case' },
            { text: '生产安全整改 SLA 案例', link: '/projects/production-safety-rectification-sla-case' },
            { text: '生产安全整改成本复盘案例', link: '/projects/production-safety-rectification-cost-review-case' },
            { text: '生产安全整改预算预测案例', link: '/projects/production-safety-rectification-budget-forecast-case' },
            { text: '生产安全整改资源排期案例', link: '/projects/production-safety-rectification-resource-scheduling-case' },
            { text: '生产安全整改产线影响评估案例', link: '/projects/production-safety-rectification-line-impact-assessment-case' },
            { text: '生产安全整改多方案决策案例', link: '/projects/production-safety-rectification-multi-scenario-decision-case' },
            { text: '生产安全整改决策复盘案例', link: '/projects/production-safety-rectification-decision-review-case' },
            { text: '生产安全整改决策知识库案例', link: '/projects/production-safety-rectification-decision-knowledge-base-case' },
            { text: '生产安全整改决策智能推荐案例', link: '/projects/production-safety-rectification-decision-intelligent-recommendation-case' },
            { text: '生产安全整改决策推荐评测案例', link: '/projects/production-safety-rectification-recommendation-evaluation-case' },
            { text: '生产良率分析案例', link: '/projects/production-yield-analysis-case' },
            { text: '生产瓶颈分析案例', link: '/projects/production-bottleneck-analysis-case' },
            { text: '生产换型损失分析案例', link: '/projects/production-changeover-loss-analysis-case' },
            { text: '生产设备异常案例', link: '/projects/production-equipment-exception-case' },
            { text: '生产能耗分析案例', link: '/projects/production-energy-analysis-case' },
            { text: '生产成本核算案例', link: '/projects/production-cost-accounting-case' },
            { text: '制造成本差异分析案例', link: '/projects/manufacturing-cost-variance-case' },
            { text: 'IoT 设备管理案例', link: '/projects/iot-device-management-case' },
            { text: '设备维保案例', link: '/projects/equipment-maintenance-case' },
            { text: '教育培训平台案例', link: '/projects/education-training-platform-case' },
            { text: '运维值班案例', link: '/projects/operations-oncall-case' },
            { text: '项目交付检查清单', link: '/projects/delivery-checklist' }
          ]
        },
        {
          text: '真实项目问题库',
          items: [
            { text: '问题库总览', link: '/projects/real-world-issues' },
            { text: '项目排障方法论', link: '/projects/debugging-playbook' },
            { text: 'Vue 项目专项', link: '/projects/issues-vue' },
            { text: 'Vue Admin 请求权限排障', link: '/projects/issues-vue-admin-request' },
            { text: 'Vue Admin 消息通知排障', link: '/projects/issues-vue-admin-notification' },
            { text: 'JavaScript 专项', link: '/projects/issues-javascript' },
            { text: 'CSS 专项', link: '/projects/issues-css' },
            { text: '前端页面与状态', link: '/projects/issues-frontend' },
            { text: 'TypeScript 类型边界', link: '/projects/issues-typescript' },
            { text: '后端接口与服务', link: '/projects/issues-backend' },
            { text: '前后端联调排查', link: '/projects/integration-debugging' },
            { text: '数据库与缓存', link: '/projects/issues-database' },
            { text: '部署、缓存与 DevOps', link: '/projects/issues-deployment' },
            { text: 'AI 工程问题', link: '/projects/issues-ai' },
            { text: '上线事故案例库', link: '/projects/production-incident-cases' },
            { text: '故障复盘模板', link: '/projects/incident-review' }
          ]
        }
      ],
      '/cheatsheets/': [
        {
          text: '速查手册总览',
          items: [
            { text: '速查手册总览', link: '/cheatsheets/' }
          ]
        },
        {
          text: '前端开发',
          items: [
            { text: 'Vue 速查', link: '/cheatsheets/vue' },
            { text: 'Vue Router 速查', link: '/cheatsheets/vue-router' },
            { text: 'Pinia 速查', link: '/cheatsheets/pinia' },
            { text: 'JavaScript 速查', link: '/cheatsheets/javascript' },
            { text: '正则速查', link: '/cheatsheets/regex' },
            { text: 'TypeScript 速查', link: '/cheatsheets/typescript' },
            { text: 'CSS 速查', link: '/cheatsheets/css' }
          ]
        },
        {
          text: '工程与部署',
          items: [
            { text: 'Vite 速查', link: '/cheatsheets/vite' },
            { text: 'Git 速查', link: '/cheatsheets/git' },
            { text: 'Linux 速查', link: '/cheatsheets/linux' },
            { text: '常用命令速查', link: '/cheatsheets/commands' },
            { text: '调试工具速查', link: '/cheatsheets/debugging-tools' },
            { text: 'Docker 速查', link: '/cheatsheets/docker' },
            { text: 'Nginx 速查', link: '/cheatsheets/nginx' }
          ]
        },
        {
          text: '后端与数据',
          items: [
            { text: 'Node.js 速查', link: '/cheatsheets/node' },
            { text: 'Java 速查', link: '/cheatsheets/java' },
            { text: 'Go 速查', link: '/cheatsheets/go' },
            { text: 'HTTP 速查', link: '/cheatsheets/http' },
            { text: 'SQL 速查', link: '/cheatsheets/sql' },
            { text: 'Redis 速查', link: '/cheatsheets/redis' }
          ]
        }
      ],
      '/contribute/': [
        {
          text: '贡献指南',
          items: [
            { text: '贡献指南', link: '/contribute/' },
            { text: '内容写作规范', link: '/contribute/content-style' },
            { text: '模块模板', link: '/contribute/module-template' },
            { text: '文档治理', link: '/contribute/governance' },
            { text: '模块状态', link: '/contribute/module-status' },
            { text: '质量检查', link: '/contribute/quality-check' }
          ]
        }
      ]
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/vuejs/vitepress' }
    ],
    footer: {
      message: '系统化整理程序员技术学习路线、工程实践和真实项目问题。',
      copyright: 'Copyright © 2026 Programmer Learning Docs'
    },
    outline: {
      label: '本页目录',
      level: [2, 3]
    },
    docFooter: {
      prev: '上一节',
      next: '下一节'
    },
    lastUpdated: {
      text: '最后更新'
    },
    darkModeSwitchLabel: '切换深色模式',
    sidebarMenuLabel: '菜单',
    returnToTopLabel: '回到顶部'
  }
})
