# 程序员技术学习文档站

这是一个基于 VitePress 的中文技术文档站。内容以“先看图建立模型，再学核心知识，然后完成项目、练习和真实问题复盘”为主线，当前已覆盖前端基础、JavaScript、TypeScript、CSS、Vue、React、Nuxt/Next、Node.js、Java、Go、数据库、浏览器、工程化、DevOps 和 AI 工程。

## 项目结构

```text
docs/
├─ .vitepress/
│  ├─ config.ts
│  └─ theme/
│     ├─ index.ts
│     ├─ styles.css
│     └─ components/
├─ roadmap/
├─ frontend/
├─ css/
├─ javascript/
├─ typescript/
├─ vue/
├─ react/
├─ meta-frameworks/
├─ node/
├─ java/
├─ go/
├─ database/
├─ browser/
├─ engineering/
├─ devops/
├─ ai-engineering/
├─ projects/
├─ cheatsheets/
├─ technologies/
└─ contribute/
examples/
├─ java-admin-api/          # Java 25 + Spring Boot + PostgreSQL 用户角色 API
└─ go-task-api/             # Go 1.26.5 + PostgreSQL 18.4 用户任务 API
scripts/
├─ check-docs.mjs
├─ check-visual-assets.mjs
└─ audit-doc-visuals.mjs
```

## 内容方向

文档以“详细、一看就懂、能解决实际项目问题”为目标。成熟模块通常包含：

- 学习导览和推荐顺序。
- 图解核心概念。
- 基础与进阶章节。
- 从零到项目落地。
- 专项练习和验收清单。
- 真实项目问题库与常见问题。

全站模块成熟度记录在 `docs/contribute/module-status.md`，问题库总入口是 `docs/projects/real-world-issues.md`。

## 启动

```bash
npm install
npm run docs:dev
```

开发服务器默认地址：`http://127.0.0.1:6173`。

## 文档检查

```bash
npm run docs:check
```

该命令检查内部路由、成熟模块结构、核心页面必备章节、配置入口和视觉资产登记；本地图片缺失、空 alt、未登记或超出体积限制都会失败。内容深度、侧边栏顺序和图示运行时渲染仍需要人工验收。

## 构建

```bash
npm run docs:build
```

生产构建后可以本地预览：

```bash
npm run docs:preview
```

## 配套示例验证

Go 普通单元测试不依赖 Docker：

```bash
cd examples/go-task-api
go test ./...
go test -race ./...
go vet ./...
```

Java 示例当前测试套件包含 Testcontainers 集成测试；下面两组数据库测试都会启动真实 PostgreSQL，需要本机 Docker daemon 可用：

```bash
cd examples/java-admin-api
mvn -B -ntp test

cd ../go-task-api
go test -tags=integration ./... -count=1 -v
```

Go 示例还可执行短时 Fuzz 与完整容器 smoke：

```bash
cd examples/go-task-api
go test ./internal/platform/httpx -run '^$' -fuzz '^FuzzDecodeJSON$' -fuzztime 10s
POSTGRES_PORT=55432 docker compose -p go-task-api up -d --build
docker compose -p go-task-api ps
docker compose -p go-task-api down -v --remove-orphans
```

完整前置条件、迁移命令、接口请求和清理风险分别记录在两个示例目录的 README 中。

## 主题说明

主题基于 VitePress 默认主题扩展，保留默认文档能力，并通过 `.vitepress/theme/styles.css` 定义清新、友好、专业的视觉变量。首页、学习路线、技术卡片和实践提示块使用自定义 Vue 组件实现。

## 样式约定

- 优先覆盖 VitePress 官方 CSS 变量。
- 业务样式使用明确 class，例如 `.custom-home__title`、`.learning-path__card`。
- 避免使用宽泛后代选择器污染默认主题或后续组件库样式。
- 固定尺寸视觉元素需要设置稳定宽高和不可压缩行为。

## 新增内容约定

扩展模块时优先补齐现有模块，不只增加零散页面。明显扩展后同步更新：

- 对应模块导览。
- `docs/.vitepress/config.ts` 导航和侧边栏。
- `docs/technologies/index.md` 技术库入口。
- `docs/technologies/expansion-plan.md` 扩展路线。
- `docs/contribute/module-status.md` 模块状态。

新增 Mermaid 图示后，除了运行构建，还要在浏览器中确认每张图生成 SVG、没有 `.mermaid-diagram__error`，并检查 390px 移动端没有页面级横向滚动。
