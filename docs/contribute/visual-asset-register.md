# 视觉资产登记

本页是 `docs/public/images` 中所有教学图片的唯一登记表。准确流程优先使用 Mermaid，真实界面优先使用可复现截图，只有抽象心智模型才使用生成图片。

## 登记规则

每个本地资产都必须包含一行机器可解析注释：

```md
<!-- asset: /images/vue/admin-list-filtered-result.webp | type: live-screenshot | license: project-owned | status: verified -->
```

随后记录：

- 使用页面和具体章节。
- 教学目的，而不是“让页面更好看”。
- 能独立描述有效信息的中文 alt。
- 解释图与相邻正文关系的中文图注。
- 截图复现步骤或生成图片完整 prompt。
- 来源、许可、生成工具、实际日期和人工核对结论。

普通图片不得超过 500 KB。只有确实无法进一步压缩且细节必须保留时，才在登记状态中加入 `large-approved`，上限为 1.5 MB。

## 媒介选择

### `diagram-sufficient`

状态、数据流、依赖和时序由 Mermaid 表达更准确时，不增加位图。完成标准是图可渲染、正文逐步解释、移动端不溢出。

### `needs-live-screenshot`

读者必须观察真实 UI 或浏览器结果时使用。场景必须能通过仓库内演示页或真实项目重现，并固定视口、设备像素比和数据状态。

### `needs-annotated-screenshot`

必须指出界面区域、前后变化或因果关系时使用。标注不能遮住原始信息，且正文逐一解释每个标注。

### `needs-generated-visual`

只用于所有权、并发协作、请求旅程等抽象类比。生成图不能承载必须精确阅读的代码、状态名或中文标签；准确规则仍由代码和 Mermaid 提供。

### `needs-official-source`

工具界面、浏览器面板或标准事实无法自行可靠复现时，使用官方来源。必须登记来源 URL、许可、版本和核对日期。

### `needs-mermaid-refactor`

已有图过密、重复或移动端不可读时，先拆分 Mermaid，不用截图掩盖结构问题。

## 当前资产

### Go API 请求往返路径

<!-- asset: /images/go/go-api-request-journey.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/go/visual-guide` 的“Handler、Service、Repository 的依赖方向”。
- **教学目的：** 把请求进入、分层处理、数据库写入和统一响应返回放在同一张图中，帮助初学者识别每层边界。
- **中文 alt：** Go API 请求从客户端依次经过中间件、严格 JSON 解码、Handler、Service、Repository 和 PostgreSQL，再携带状态码、统一 JSON 与同一个 Request ID 返回客户端。
- **图注：** 绿色箭头表示请求进入，黄色路径表示响应返回；准确职责以相邻 Mermaid、正文和真实源码为准。
- **复现来源：** `docs/public/visual-demos/go-api-request-journey.html`，在 `1600 × 900`、设备像素比 1 下通过本地 VitePress 预览截图，再用 `cwebp -q 86` 转换。
- **工具与日期：** Codex in-app Browser 截图、cwebp，2026-07-21。
- **人工核对：** 已确认七个阶段顺序、往返箭头、中文标签和 Request ID 含义正确；1600 × 900 WebP 为 51,580 bytes，无裁切和错误文字。

### Go 并发所有权工作坊类比

<!-- asset: /images/go/go-concurrency-workshop.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/go/visual-guide` 的“有界 worker pool”。
- **教学目的：** 用工作流类比解释 channel 关闭权、有界队列背压、worker 退出条件和 Context 取消传播。
- **中文 alt：** 一个协调者创建并唯一关闭容量为五的任务通道，三个 Worker 从通道取任务，并在通道关闭或 Context 取消时分别退出并报告完成。
- **图注：** 这是帮助理解所有权和退出条件的类比；准确的 channel、Context 与 WaitGroup 规则以相邻 Mermaid 和代码为准。
- **复现来源：** `docs/public/visual-demos/go-concurrency-workshop.html`，在 `1600 × 900`、设备像素比 1 下通过本地 VitePress 预览截图，再用 `cwebp -q 86` 转换。
- **工具与日期：** Codex in-app Browser 截图、cwebp，2026-07-21。
- **人工核对：** 已确认只有协调者拥有 close、缓冲容量为五、三个 Worker 都有退出路径；1600 × 900 WebP 为 56,508 bytes，无裁切和错误文字。

## 前端基础教学截图

### HTML 语义化文章结构

<!-- asset: /images/frontend/semantic-article.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/html-semantics` 的“页面区域怎样选择”。
- **教学目的：** 对照真实排版识别 `header`、`nav`、`main`、`article`、`aside` 和 `footer` 的职责。
- **中文 alt：** 语义化技术文章页面标出页头、导航、主体、文章、补充信息和页脚区域。
- **图注：** 语义标签表达内容职责，不是为了得到默认样式。
- **复现来源：** `/visual-demos/frontend/index.html?scene=semantic-article`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 六个页面区域和阅读顺序正确；46,418 bytes，无裁切或横向溢出。

### 表单有效状态

<!-- asset: /images/frontend/form-valid.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/forms-media-accessibility` 的“校验失败时要解释原因”。
- **教学目的：** 展示标签、帮助文本和成功反馈在正常提交前的协作方式。
- **中文 alt：** 注册表单的姓名和邮箱字段填写有效，页面展示绿色校验通过提示。
- **图注：** 有效状态要保持字段说明可见，并用文字和颜色共同反馈。
- **复现来源：** `/visual-demos/frontend/index.html?scene=form-valid`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 标签、输入值、帮助文本与成功提示对应正确；36,216 bytes。

### 表单无效状态

<!-- asset: /images/frontend/form-invalid.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/forms-media-accessibility` 的“校验失败时要解释原因”。
- **教学目的：** 说明错误摘要、字段错误和 `aria-describedby` 应如何共同定位问题。
- **中文 alt：** 注册表单的邮箱格式错误，字段下方显示原因并在顶部汇总一个待修正问题。
- **图注：** 错误反馈要说明具体原因，并与对应字段建立可访问关联。
- **复现来源：** `/visual-demos/frontend/index.html?scene=form-invalid`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 错误摘要与邮箱字段一致，未用颜色作为唯一线索；40,276 bytes。

### 响应式桌面布局

<!-- asset: /images/frontend/responsive-desktop.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/html-css` 的“响应式布局”。
- **教学目的：** 展示宽屏下导航、主内容和辅助栏的信息层级。
- **中文 alt：** 宽屏技术文档页面采用顶部导航、双栏主体和三列学习卡片布局。
- **图注：** 宽屏增加并行信息量，但主内容仍保持明确的阅读顺序。
- **复现来源：** `/visual-demos/frontend/index.html?scene=responsive-desktop`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 主次区域、列宽和间距稳定；30,088 bytes，无横向溢出。

### 响应式移动布局

<!-- asset: /images/frontend/responsive-mobile.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/html-css` 的“响应式布局”。
- **教学目的：** 对照桌面图理解窄屏下导航收纳、内容重排和操作区保留。
- **中文 alt：** 同一技术文档页面在手机宽度下改为单栏内容和紧凑顶部导航。
- **图注：** 移动端不是等比缩小桌面页，而是按任务优先级重新排列内容。
- **复现来源：** `/visual-demos/frontend/index.html?scene=responsive-mobile`，390 × 844、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 单栏顺序、按钮触达区和文本换行正确；18,608 bytes，无横向溢出。

### 键盘焦点样式

<!-- asset: /images/frontend/accessible-focus.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/frontend/forms-media-accessibility` 的“键盘操作”和 `/css/accessibility` 的“焦点样式要稳定”。
- **教学目的：** 展示不改变布局的 `focus-visible` 外轮廓。
- **中文 alt：** 表单控件获得键盘焦点时显示清晰的绿色外轮廓和辅助说明。
- **图注：** 焦点环既要明显，也不能通过改变边框宽度造成页面抖动。
- **复现来源：** `/visual-demos/frontend/index.html?scene=accessible-focus`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 焦点目标、轮廓和提示文本清楚，颜色对比可辨；39,110 bytes。

## CSS 教学截图

### 盒模型四层结构

<!-- asset: /images/css/box-model.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/box-model-layout` 的“盒模型”。
- **教学目的：** 把 content、padding、border、margin 从内向外的关系可视化。
- **中文 alt：** CSS 盒模型四层结构依次标出内容区、内边距、边框和外边距。
- **图注：** 盒模型是逐层包裹的空间结构，不是四个孤立属性。
- **复现来源：** `/visual-demos/css/index.html?scene=box-model`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四层顺序与尺寸计算规则一致；30,720 bytes。

### 垂直外边距折叠

<!-- asset: /images/css/margin-collapse.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/box-model-layout` 的“盒模型”。
- **教学目的：** 解释相邻块级元素垂直 margin 不一定相加。
- **中文 alt：** 两个上下排列的块级元素外边距折叠，最终间距取较大的 32 像素。
- **图注：** 普通文档流中的特定垂直外边距会折叠，Flex、Grid 和 BFC 场景不同。
- **复现来源：** `/visual-demos/css/index.html?scene=margin-collapse`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 24px、32px 与最终 32px 的关系标注正确；38,570 bytes。

### Flex 主轴与交叉轴

<!-- asset: /images/css/flex-main-cross-axis.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/flex-grid` 的“Flex 常见属性”。
- **教学目的：** 帮助读者按轴向选择 `justify-content` 与 `align-items`。
- **中文 alt：** Flex 容器标出主轴、交叉轴以及两个对齐属性的作用方向。
- **图注：** `flex-direction` 决定主轴，其他对齐属性要基于轴向理解。
- **复现来源：** `/visual-demos/css/index.html?scene=flex-main-cross-axis`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 两条轴和属性方向正确；33,430 bytes。

### Flex 长文本溢出

<!-- asset: /images/css/flex-overflow.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/flex-grid` 的“Flex 常见属性”。
- **教学目的：** 对比默认最小内容宽度与 `min-width: 0` 修复结果。
- **中文 alt：** Flex 子项长文本溢出与设置 min-width 0 后正确省略的并排对比。
- **图注：** 省略号不生效时，要先检查 Flex 子项是否允许收缩。
- **复现来源：** `/visual-demos/css/index.html?scene=flex-overflow`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 问题态和修复态差异清晰；48,186 bytes。

### Grid 具名区域

<!-- asset: /images/css/grid-template-areas.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/flex-grid` 的“Grid：二维布局”。
- **教学目的：** 展示 `grid-template-areas` 如何表达页面二维骨架。
- **中文 alt：** CSS Grid 使用 header、sidebar、main 和 aside 具名区域组成页面布局。
- **图注：** 具名区域把页面骨架写成可阅读的二维地图。
- **复现来源：** `/visual-demos/css/index.html?scene=grid-template-areas`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 区域名称与位置一一对应；25,424 bytes。

### Grid 自动列数

<!-- asset: /images/css/grid-minmax.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/flex-grid` 的“Grid：二维布局”。
- **教学目的：** 对比不同容器宽度下 `auto-fit` 与 `minmax` 的自动排布。
- **中文 alt：** CSS Grid 卡片在宽、中、窄容器中自动从多列调整为少列。
- **图注：** `minmax` 定义单列范围，`auto-fit` 决定当前能容纳的列数。
- **复现来源：** `/visual-demos/css/index.html?scene=grid-minmax`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 三种容器宽度的列数变化正确；38,022 bytes。

### 响应式内容容器

<!-- asset: /images/css/responsive-container.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/responsive` 的“常见断点”。
- **教学目的：** 说明断点应由内容可读性触发，而不是设备名称触发。
- **中文 alt：** 同一内容容器在宽屏、中屏和窄屏下从三列变为两列再变为一列。
- **图注：** 断点保护内容空间和操作效率，具体数值应根据真实布局确定。
- **复现来源：** `/visual-demos/css/index.html?scene=responsive-container`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 三种重排状态无裁切或溢出；33,798 bytes。

### 减少动态效果

<!-- asset: /images/css/reduced-motion.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/css/animation-transition` 的“尊重减少动态效果”。
- **教学目的：** 对比正常动画和 `prefers-reduced-motion` 模式的反馈差异。
- **中文 alt：** 普通动画偏好与减少动态效果模式并排展示位移和状态反馈差异。
- **图注：** 减少动态不是删除状态反馈，而是取消非必要位移和持续动画。
- **复现来源：** `/visual-demos/css/index.html?scene=reduced-motion`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 两种偏好的行为说明准确；42,254 bytes。

## Vue Admin 项目截图

### 工作台总览

<!-- asset: /images/vue/admin-dashboard.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-dashboard-analytics` 的“这一页最终要做到什么”。
- **教学目的：** 展示筛选区、指标卡、趋势、排行和待办区块组成的工作台层级。
- **中文 alt：** Vue Admin 工作台包含时间筛选、指标卡片、趋势图、排行榜和待办任务。
- **图注：** 工作台由多个独立数据区块组成，每块应独立加载、失败和刷新。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=dashboard`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 指标、图表、排行和待办层级清晰；40,200 bytes。

### 用户列表默认态

<!-- asset: /images/vue/admin-list-default.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-list-search-table` 的“最终目标”。
- **教学目的：** 展示筛选、总数、表格、行操作和分页的默认协作状态。
- **中文 alt：** Vue Admin 用户列表默认状态包含筛选表单、表格、状态标签和分页。
- **图注：** 默认态要同时说明当前条件、结果范围和可执行操作。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=list-default`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 查询区、表格和分页状态一致；44,266 bytes。

### 用户列表筛选态

<!-- asset: /images/vue/admin-list-filtered.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-list-search-table` 的“最终目标”。
- **教学目的：** 说明筛选条件改变后，结果、总数、分页和导出范围要同步。
- **中文 alt：** Vue Admin 用户列表应用部门和状态筛选后展示匹配结果。
- **图注：** 所有结果区域必须消费同一份标准化查询参数。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=list-filtered`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 已选条件与结果数量相符；35,166 bytes。

### 用户列表空态

<!-- asset: /images/vue/admin-list-empty.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-list-search-table` 的“最终目标”。
- **教学目的：** 区分无结果与接口错误，并提供清除筛选恢复入口。
- **中文 alt：** Vue Admin 用户列表没有匹配结果，显示空状态和清除筛选按钮。
- **图注：** 空态要解释当前没有匹配数据，而不是只留一块空白。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=list-empty`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 空状态原因和恢复动作清晰；31,008 bytes。

### 新增用户表单

<!-- asset: /images/vue/admin-form-create.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-form-modal-crud` 的“最终目标”。
- **教学目的：** 展示新增模式的字段分区、初始值和固定操作区。
- **中文 alt：** Vue Admin 新增用户抽屉表单包含基础信息、角色、状态和保存操作。
- **图注：** 新增态应由独立 FormState 初始化，不能复用上一次编辑数据。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=form-create`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 字段分组、必填标识和操作区完整；28,794 bytes。

### 表单校验错误

<!-- asset: /images/vue/admin-form-validation.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-form-modal-crud` 的“最终目标”。
- **教学目的：** 对比必填、格式与服务端冲突错误的字段级展示。
- **中文 alt：** Vue Admin 用户表单在对应字段附近展示必填、格式和账号冲突错误。
- **图注：** 全局错误提示负责总结，字段级错误负责帮助用户直接修正。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=form-validation`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 三类错误均有文字说明且位置正确；31,964 bytes。

### 详情与审计时间线

<!-- asset: /images/vue/admin-detail-audit.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-detail-status-audit` 的“最终目标”。
- **教学目的：** 展示状态、主信息、允许操作与历史证据的阅读顺序。
- **中文 alt：** Vue Admin 订单详情展示当前状态、客户信息、金额、操作按钮和审计时间线。
- **图注：** 用户应先判断现状和可执行动作，再沿时间线追溯原因。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=detail-audit`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 状态、详情字段与时间线逻辑一致；38,450 bytes。

### 待审批任务

<!-- asset: /images/vue/admin-approval-pending.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-approval-workflow` 的“最终要做到什么”。
- **教学目的：** 把业务摘要、流程位置、审批意见和节点动作放在同一页面理解。
- **中文 alt：** Vue Admin 采购审批详情展示业务摘要、审批节点、处理意见和同意驳回操作。
- **图注：** 审批动作由任务归属、节点状态、权限和流程版本共同决定。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=approval-pending`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 当前节点、操作区和历史节点关系正确；43,490 bytes。

### 消息中心未读态

<!-- asset: /images/vue/admin-notification-unread.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-notification-center` 的“最终要做到什么”。
- **教学目的：** 展示未读总数、消息分类、列表状态和批量操作的一致性。
- **中文 alt：** Vue Admin 消息中心展示未读计数、消息分类、未读列表和批量已读操作。
- **图注：** 未读数来自服务端记录与增量同步，不能用当前页长度代替。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=notification-unread`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 未读标签、总数和操作入口清楚；39,330 bytes。

### 权限拒绝页

<!-- asset: /images/vue/admin-permission-denied.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-permission-route-flow` 的“最终要做成什么”。
- **教学目的：** 展示确认 403 后的安全反馈与恢复入口。
- **中文 alt：** Vue Admin 无权限页面展示 403 状态、原因说明和返回工作台按钮。
- **图注：** 无权限时不能先渲染业务数据再隐藏，接口也必须再次校验。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=permission-denied`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 页面未暴露业务内容，原因和返回动作明确；31,260 bytes。

### 权限恢复加载态

<!-- asset: /images/vue/admin-permission-loading.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-permission-route-flow` 的“最终要做成什么”。
- **教学目的：** 解释刷新深层路由时恢复用户、菜单和动态路由的启动阶段。
- **中文 alt：** Vue Admin 刷新页面时依次恢复用户信息、菜单权限和动态路由。
- **图注：** 权限恢复期间应显示独立 loading，防止路由过早判定和页面闪烁。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=permission-loading`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 三个恢复阶段顺序准确；23,786 bytes。

### 文件上传进度

<!-- asset: /images/vue/admin-file-upload-progress.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/vue/admin-file-import-export` 的“这一页最终要做到什么”。
- **教学目的：** 展示上传任务需要暴露的文件、阶段、百分比、速度、剩余时间和取消能力。
- **中文 alt：** Vue Admin 文件上传任务展示进度百分比、速度、剩余时间和取消按钮。
- **图注：** 上传进度属于具体文件任务，重试、取消和页面切换都要保留明确状态。
- **复现来源：** `/visual-demos/vue-admin/index.html?scene=file-upload-progress`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 文件信息、进度和操作对应正确；33,430 bytes。

## 浏览器证据实验台截图

### Network 请求头证据

<!-- asset: /images/browser/network-request-headers.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/http-request` 的“Network 面板排查流程”和 `/browser/project-from-zero` 的“DevTools 证据链”。
- **教学目的：** 把 URL、Method、Status、请求头、响应头和 Request ID 放进同一条证据链。
- **中文 alt：** 浏览器网络请求证据面板展示 URL、GET 方法、200 状态、请求头和响应头。
- **图注：** 先核对请求事实，再解释页面为什么成功或失败。
- **复现来源：** `/visual-demos/browser/index.html?scene=network-request-headers`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 请求与响应使用同一 `req_demo_7f29`，无真实令牌与用户数据；50,858 bytes。

### Network Timing 证据

<!-- asset: /images/browser/network-request-timing.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/http-request` 的“Network 面板排查流程”。
- **教学目的：** 把请求总耗时拆成 DNS、连接、等待首字节和下载阶段。
- **中文 alt：** 浏览器 Network Timing 将 218 毫秒请求拆成 DNS、连接、等待首字节和下载阶段。
- **图注：** Waiting 占比最大时先对齐服务端日志，而不是先压缩很小的响应体。
- **复现来源：** `/visual-demos/browser/index.html?scene=network-request-timing`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四阶段相加为 218 ms，时间条顺序一致；46,696 bytes。

### 缓存来源证据

<!-- asset: /images/browser/cache-memory-disk.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/cache` 的“缓存分几层”。
- **教学目的：** 区分 memory cache、disk cache、304 协商缓存和正常网络响应。
- **中文 alt：** 浏览器资源列表对比 memory cache、disk cache、304 协商缓存和正常网络响应。
- **图注：** 资源来源、缓存头和文件名 hash 必须一起判断。
- **复现来源：** `/visual-demos/browser/index.html?scene=cache-memory-disk`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四类资源状态与示例缓存策略对应；45,966 bytes。

### LocalStorage 与 SessionStorage

<!-- asset: /images/browser/storage-local-session.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/storage` 的“常见存储方式”。
- **教学目的：** 对比同源多标签页中的共享范围和生命周期。
- **中文 alt：** 浏览器 Application 面板对比 LocalStorage 跨同源标签共享与 SessionStorage 按标签隔离。
- **图注：** LocalStorage 的值跨同源标签共享，SessionStorage 通常只属于当前标签页。
- **复现来源：** `/visual-demos/browser/index.html?scene=storage-local-session`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 示例 key 均为非敏感演示数据，共享与隔离状态准确；48,566 bytes。

### CORS 预检证据

<!-- asset: /images/browser/cors-preflight.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/cors-auth` 的“简单请求和预检请求”。
- **教学目的：** 展示 OPTIONS、允许响应和真实 POST 的先后与条件关系。
- **中文 alt：** 跨域请求先发送 OPTIONS 预检，服务端返回允许规则后再发送带 Authorization 的 POST。
- **图注：** OPTIONS 失败时，真正的业务请求不会发送。
- **复现来源：** `/visual-demos/browser/index.html?scene=cors-preflight`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** Origin、Method 与 Header 的询问和允许项匹配；56,674 bytes。

### Performance 长任务证据

<!-- asset: /images/browser/performance-long-task.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/rendering-performance` 的“DevTools 性能排查”。
- **教学目的：** 展示超过 50 ms 的主线程任务如何推迟输入处理。
- **中文 alt：** 浏览器 Performance 主线程轨道显示 121 毫秒长任务并推迟输入响应。
- **图注：** 红色长任务说明主线程被占用，具体根因仍要通过调用栈判断。
- **复现来源：** `/visual-demos/browser/index.html?scene=performance-long-task`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** Main 与 Interactions 时间关系一致，未把颜色直接解释为根因；45,638 bytes。

### Playwright 定位失败证据

<!-- asset: /images/browser/automation-locator-failure.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/browser-automation-debugging` 的“选择器一改就全挂”。
- **教学目的：** 对比 Locator 的测试意图与同一时刻 DOM 快照。
- **中文 alt：** Playwright 定位失败日志寻找提交按钮，而同一时刻 DOM 快照只存在保存按钮。
- **图注：** 应修正文案契约或定位意图，不能退化为依赖按钮序号的选择器。
- **复现来源：** `/visual-demos/browser/index.html?scene=automation-locator-failure`，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 错误日志、DOM 可访问名称和解决结论一致；63,010 bytes。

## 后端项目真实运行证据

### Java Admin API Readiness

<!-- asset: /images/java/java-admin-api-ready.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/java/spring-boot-project-from-zero` 的“运行和手工联调”。
- **教学目的：** 展示 PostgreSQL、Flyway、非 root 应用与 readiness 的完整启动证据链。
- **中文 alt：** Java Admin API 容器依次完成 PostgreSQL 健康检查、Flyway V2 迁移、非 root 启动并返回 readiness UP。
- **图注：** Readiness 证明数据库、迁移和 Spring Context 已可服务，不只是 Java 进程存活。
- **复现来源：** `maven:3.9.11-eclipse-temurin-25` 运行 9 个测试；`java-admin-api:evidence` 连接 `postgres:18-alpine`；请求 `/actuator/health/readiness`；再由 `/visual-demos/backend-evidence/index.html?scene=java-admin-api-ready` 重排脱敏证据，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Docker 29.2.1、Java 25.0.1、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** Flyway V1/V2、PostgreSQL 18.4、容器 uid 10001 和 `{"status":"UP"}` 均来自本次真实运行；56,836 bytes。

### Java Admin API 版本冲突

<!-- asset: /images/java/java-admin-version-conflict.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/java/spring-boot-project-from-zero` 的“乐观锁实验”。
- **教学目的：** 说明两个客户端持有同一旧版本时，第二次更新为何必须返回 409。
- **中文 alt：** Java Admin API 两个客户端同时持有 version 1，第一个更新到 version 2，第二个收到 409 STALE_VERSION。
- **图注：** 收到 409 后重新读取并决定如何处理，不能自动重放旧意图。
- **复现来源：** 创建 `evidence@example.test` 用户，第一次 PATCH 使用 `expectedVersion=1` 返回 version 2，再次使用 1 返回 `409 STALE_VERSION`；由 `/visual-demos/backend-evidence/index.html?scene=java-admin-version-conflict` 重排脱敏证据，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** curl、jq、Docker、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 200/409、版本 1→2 和错误码均与实际响应一致；48,892 bytes。

### Go Task API Readiness

<!-- asset: /images/go/go-task-api-ready.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/go/http-api-project-from-zero` 的“执行真实 smoke test”。
- **教学目的：** 展示 PostgreSQL healthy、migration exit 0、API healthy 与数据库 readiness 的依赖顺序。
- **中文 alt：** Go Task API 依次完成 PostgreSQL 健康检查、数据库迁移、非 root 启动并返回 ready。
- **图注：** Ready 会执行数据库 PingContext，数据库不可用或关闭中应返回 503。
- **复现来源：** `POSTGRES_PORT=55434 docker compose -p go-task-evidence up -d --no-build`，请求 `/health/ready` 并检查容器用户；由 `/visual-demos/backend-evidence/index.html?scene=go-task-api-ready` 重排脱敏证据，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** Go 1.26.5 镜像、PostgreSQL 18.4、Docker、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** migration 版本 1、`nonroot:nonroot`、healthy 与 `data.status=ready` 均来自本次真实运行；62,968 bytes。

### Go Task API 版本冲突

<!-- asset: /images/go/go-task-version-conflict.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/go/http-api-project-from-zero` 的“创建用户和任务”。
- **教学目的：** 说明条件更新如何阻止旧页面覆盖已经更新的任务。
- **中文 alt：** Go Task API 两个客户端同时持有任务 version 0，第一个更新到 version 1，第二个收到 409 TASK_VERSION_CONFLICT。
- **图注：** `UPDATE WHERE version=0` 最多成功一次，第二个旧版本请求得到稳定错误码。
- **复现来源：** 创建 task 1 得到 version 0，第一次 PUT 返回 version 1，第二次使用 `expectedVersion=0` 返回 `409 TASK_VERSION_CONFLICT`；由 `/visual-demos/backend-evidence/index.html?scene=go-task-version-conflict` 重排证据，1440 × 900、DPR 1、`cwebp -q 88`。
- **工具与日期：** curl、jq、Docker、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 200/409、版本 0→1、Request ID 和错误码均与实际响应一致；49,256 bytes。

## 抽象技术心智模型

### 浏览器渲染阶段地图

<!-- asset: /images/browser/browser-rendering-pipeline.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/browser/rendering-performance` 的“浏览器如何把代码变成页面”。
- **教学目的：** 建立 DOM/CSSOM、Render Tree、Layout、Paint 与 Composite 的空间顺序。
- **中文 alt：** HTML CSS 和 JavaScript 依次经过 DOM 与 CSSOM、渲染树、布局、绘制和图层合成成为像素。
- **图注：** 这是阶段地图，是否跳过或重复某阶段仍以 Performance 实际录制为准。
- **复现来源：** `/visual-demos/concepts/index.html?scene=browser-rendering-pipeline`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 阶段顺序、DOM/CSSOM 合流和重排/重绘边界说明准确；49,542 bytes。

### JVM 内存区域地图

<!-- asset: /images/java/jvm-memory-regions.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/java/jvm-memory-gc` 的“JVM 运行结构”。
- **教学目的：** 区分线程私有栈、共享堆、Metaspace 与 GC Roots 可达性追踪。
- **中文 alt：** JVM 多个线程拥有独立调用栈，共享对象堆和 Metaspace，GC 从 Roots 追踪对象可达性。
- **图注：** 对象数量是教学示意，实际内存和回收行为要通过 JVM 诊断工具确认。
- **复现来源：** `/visual-demos/concepts/index.html?scene=jvm-memory-regions`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 私有/共享边界、堆对象和 Metaspace 职责未混淆；51,332 bytes。

### Node.js 事件循环与 I/O

<!-- asset: /images/node/event-loop-io-workshop.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/node/runtime-event-loop` 的“非阻塞 I/O”。
- **教学目的：** 区分主线程协调、外部 I/O 等待、完成回调队列和 CPU 阻塞。
- **中文 alt：** Node.js 事件循环协调文件数据库和网络 I O，完成回调返回队列，而 CPU 长任务阻塞主线程。
- **图注：** 非阻塞只描述 I/O 等待；JavaScript 回调和 CPU 计算仍占用主线程。
- **复现来源：** `/visual-demos/concepts/index.html?scene=node-event-loop-io-workshop`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** I/O 不占用主线程空转、回调不抢占、CPU 阻塞三条边界准确；55,486 bytes。

### RAG 入库与检索旅程

<!-- asset: /images/ai-engineering/rag-retrieval-journey.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/ai-engineering/rag` 的“基本流程”。
- **教学目的：** 分开说明知识入库与用户问答两条链路，并把引用与权限放入验收范围。
- **中文 alt：** RAG 入库链将来源文档切分并写入语义索引，问答链检索少量证据后生成带引用答案。
- **图注：** 整体图不能替代召回、重排、权限、回答与引用的独立评测。
- **复现来源：** `/visual-demos/concepts/index.html?scene=rag-retrieval-journey`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 入库与查询分离、top-k 证据和带来源回答关系准确；56,378 bytes。

## P1 核心模块教学图

### JavaScript 事件循环调试时间线

<!-- asset: /images/javascript/event-loop-devtools.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/javascript/event-loop`；**教学目的：** 把同步任务、微任务、渲染机会和定时器放到同一条时间线上。
- **中文 alt：** JavaScript 同步脚本执行后清空 Promise 微任务，浏览器获得渲染机会，随后执行 setTimeout。
- **图注：** 当前同步任务结束后会先清空微任务，再进入后续渲染和宏任务阶段。
- **复现来源：** `/visual-demos/p1/index.html?scene=event-loop-devtools`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 队列执行顺序、渲染机会和 DevTools 时间标记准确；44,206 bytes。

### DOM 事件传播路径

<!-- asset: /images/javascript/dom-event-path.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/javascript/dom-events`；**教学目的：** 解释捕获、目标和冒泡三个阶段及事件委托的位置。
- **中文 alt：** DOM 点击事件从 window 和 document 捕获到 button，再经过列表项向 document 冒泡。
- **图注：** `event.target` 是触发节点，`currentTarget` 是当前正在执行监听器的节点。
- **复现来源：** `/visual-demos/p1/index.html?scene=dom-event-path`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 捕获、目标、冒泡方向以及委托监听位置准确；39,526 bytes。

### JavaScript 任务看板状态

<!-- asset: /images/javascript/task-board-states.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/javascript/task-board-project`；**教学目的：** 展示业务页面必须覆盖的加载、成功、空和失败状态。
- **中文 alt：** JavaScript 任务看板展示加载中、正常列表、空列表和离线失败四种状态。
- **图注：** 状态分支应由统一状态模型驱动，而不是在多个 DOM 操作中零散判断。
- **复现来源：** `/visual-demos/p1/index.html?scene=task-board-states`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四种状态互斥、文案和操作入口完整；47,048 bytes。

### TypeScript 控制流收窄

<!-- asset: /images/typescript/narrowing-control-flow.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/typescript/narrowing-guards`；**教学目的：** 展示外部 `unknown` 数据如何经过守卫变成可安全使用的联合类型。
- **中文 alt：** TypeScript 从 unknown 输入经过结构校验、判别联合和穷尽检查逐步收窄类型。
- **图注：** 类型断言不会验证运行时数据，边界数据必须先经过实际校验。
- **复现来源：** `/visual-demos/p1/index.html?scene=narrowing-control-flow`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** `unknown`、守卫、判别字段和 `never` 穷尽检查关系准确；42,252 bytes。

### TypeScript 数据边界与错误状态

<!-- asset: /images/typescript/type-boundary-error-state.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/typescript/type-boundary-project`；**教学目的：** 区分接口 DTO、领域映射、页面模型和表单错误。
- **中文 alt：** UserDTO 经过运行时校验和 mapper 转成 UserViewModel，表单错误使用独立状态保存。
- **图注：** 服务端结构不应直接成为页面状态，映射层负责隔离变化和补齐语义。
- **复现来源：** `/visual-demos/p1/index.html?scene=type-boundary-error-state`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** DTO、运行时校验、ViewModel 和错误状态职责边界准确；44,172 bytes。

### TypeScript 模块解析追踪

<!-- asset: /images/typescript/tsconfig-trace.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/typescript/tsconfig-engineering`；**教学目的：** 说明 `extends`、`paths`、包 `exports` 与真实解析结果的关系。
- **中文 alt：** TypeScript 模块解析报告展示 tsconfig 继承、paths 命中、包 exports 与缺失模块候选。
- **图注：** `paths` 主要影响类型解析，运行时打包器仍需配置对应别名。
- **复现来源：** `/visual-demos/p1/index.html?scene=tsconfig-trace`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 配置继承、别名匹配、包入口和失败候选的顺序准确；55,134 bytes。

### React Effect 请求竞态

<!-- asset: /images/react/effect-request-race.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/react/effects`；**教学目的：** 展示依赖变化后清理旧副作用如何避免过期响应覆盖新结果。
- **中文 alt：** React 查询从 v 变为 vue 后清理旧 Effect，取消第一个请求并只接收第二个响应。
- **图注：** 清理函数必须对应本次 Effect 创建的请求、订阅或计时器。
- **复现来源：** `/visual-demos/p1/index.html?scene=effect-request-race`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** Effect 创建、cleanup、AbortController 与最终响应顺序准确；38,454 bytes。

### React Profiler 报告

<!-- asset: /images/react/performance-profiler.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/react/performance`；**教学目的：** 演示如何从提交耗时定位慢组件和重新渲染原因。
- **中文 alt：** React Profiler 报告显示一次 18.4 毫秒提交中 OrderTable 最慢，并列出 props state 和 context 变化。
- **图注：** 先用 Profiler 找到真实瓶颈，再决定是否拆分组件、缓存计算或稳定引用。
- **复现来源：** `/visual-demos/p1/index.html?scene=performance-profiler`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** flame chart、提交耗时和重渲染原因提示可对应阅读；48,104 bytes。

### React Admin 页面状态

<!-- asset: /images/react/admin-states.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/react/project-admin`；**教学目的：** 展示后台列表的成功、空、失败和无权限完整状态。
- **中文 alt：** React Admin 页面展示列表就绪、无结果、接口失败和 403 无权限四种状态。
- **图注：** 错误与无权限不是同一种状态，恢复动作也应不同。
- **复现来源：** `/visual-demos/p1/index.html?scene=admin-states`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 状态、解释文案和操作入口互不混淆；50,108 bytes。

### 全栈路由数据边界

<!-- asset: /images/meta-frameworks/route-data-boundaries.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/meta-frameworks/routing-data`；**教学目的：** 解释一次路由请求中匹配、鉴权、加载、渲染与激活的边界。
- **中文 alt：** 全栈框架路由依次完成参数匹配、服务端鉴权、数据加载、页面渲染和客户端交互。
- **图注：** 权限判断应尽量发生在读取敏感数据之前。
- **复现来源：** `/visual-demos/p1/index.html?scene=route-data-boundaries`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 服务端和客户端职责、失败边界与顺序准确；44,642 bytes。

### 服务端鉴权与重定向

<!-- asset: /images/meta-frameworks/server-auth-redirect.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/meta-frameworks/server-auth`；**教学目的：** 展示受保护路由在服务端的认证、授权和重定向分支。
- **中文 alt：** 服务端收到管理页面请求后读取会话、检查权限，无权限重定向，有权限才渲染受保护数据。
- **图注：** 隐藏按钮不能替代服务端授权，数据加载层仍需校验访问范围。
- **复现来源：** `/visual-demos/p1/index.html?scene=server-auth-redirect`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 未登录、无权限和允许访问三条路径清晰；43,064 bytes。

### 课程平台数据边界

<!-- asset: /images/meta-frameworks/course-platform.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/meta-frameworks/project-from-zero`；**教学目的：** 用完整页面说明服务端数据和客户端交互的拆分方式。
- **中文 alt：** 课程平台将课程概要、用户学习进度、客户端视频播放和练习提交拆成独立数据边界。
- **图注：** 首屏内容优先由服务端提供，强交互区域再在客户端激活。
- **复现来源：** `/visual-demos/p1/index.html?scene=course-platform`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 课程、进度、播放和提交边界与页面状态一致；55,016 bytes。

### Node 权限接口响应

<!-- asset: /images/node/permission-api-response.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/node/permission-api-project`；**教学目的：** 展示操作权限与数据范围必须同时进入响应和查询条件。
- **中文 alt：** Node 权限接口报告列出 users read write、roles grant 和 audit read 的允许结果及部门数据范围。
- **图注：** 能否执行操作和能看哪些数据是两套判断，不能只保留布尔权限。
- **复现来源：** `/visual-demos/p1/index.html?scene=permission-api-response`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 权限项、拒绝项和部门范围的语义完整；47,598 bytes。

### Node 缓存与队列看板

<!-- asset: /images/node/cache-queue-dashboard.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/node/cache-queue-project`；**教学目的：** 把缓存命中、回源、队列积压和失败任务放在同一运行视图中。
- **中文 alt：** Node 缓存队列运行看板展示缓存命中率、回源延迟、队列积压、失败任务和各队列状态。
- **图注：** 命中率正常不代表系统健康，还需同时观察回源延迟与积压趋势。
- **复现来源：** `/visual-demos/p1/index.html?scene=cache-queue-dashboard`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 指标、告警和队列状态的关联准确；47,482 bytes。

### Java Testcontainers 测试运行

<!-- asset: /images/java/testcontainers-run.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/java/testing-deployment`；**教学目的：** 证明迁移和集成测试在真实 PostgreSQL 容器中执行。
- **中文 alt：** Java Testcontainers 使用 PostgreSQL 18.4 执行 Flyway V1 V2 并运行 9 个测试全部通过。
- **图注：** 集成测试应验证真实数据库约束与迁移，而不只使用内存替身。
- **复现来源：** Maven 3.9.11 + Temurin 25 容器运行 `mvn -B -ntp test`；由 `/visual-demos/p1/index.html?scene=testcontainers-run` 重排脱敏证据，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** Maven、Java 25.0.1、Testcontainers、PostgreSQL 18.4、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** Flyway V1/V2、容器版本和 9 项测试结果均来自真实运行；43,870 bytes。

### PostgreSQL 索引执行计划

<!-- asset: /images/database/index-explain.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/database/indexes`；**教学目的：** 教读者从扫描方式、估算行数、实际行数、耗时和缓冲区判断索引效果。
- **中文 alt：** PostgreSQL EXPLAIN ANALYZE 报告展示索引扫描、实际 20 行、1.84 毫秒和缓冲命中。
- **图注：** 是否使用索引不是唯一目标，应结合返回规模、随机 I/O 和总耗时判断。
- **复现来源：** `/visual-demos/p1/index.html?scene=index-explain`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 节点层级、rows、time、buffers 和过滤条件可对应阅读；51,882 bytes。

### PostgreSQL 事务锁等待

<!-- asset: /images/database/transaction-lock-wait.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/database/transactions`；**教学目的：** 展示持锁事务如何阻塞另一个更新以及排查入口。
- **中文 alt：** PostgreSQL 锁等待报告显示 PID 8177 持有锁并处于长事务，PID 8421 等待 4.8 秒。
- **图注：** 处理阻塞前先确认业务事务和影响范围，不应看到等待就直接终止连接。
- **复现来源：** `/visual-demos/p1/index.html?scene=transaction-lock-wait`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** blocker、waiter、锁对象和长事务线索准确；48,028 bytes。

### 数据库权限项目模型

<!-- asset: /images/database/permission-project.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/database/project-practice`；**教学目的：** 说明用户、角色、权限与数据范围如何共同约束查询。
- **中文 alt：** 后台权限数据库从用户、用户角色、角色权限和数据范围生成带租户条件的业务查询。
- **图注：** RBAC 解决操作权限，租户和部门范围仍必须进入 SQL 条件。
- **复现来源：** `/visual-demos/p1/index.html?scene=permission-project`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 关联表、数据范围和最终查询条件的关系准确；44,548 bytes。

### 前端 Bundle 分析

<!-- asset: /images/engineering/bundle-analysis.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/engineering/bundle-analysis`；**教学目的：** 展示如何由体积报告发现大模块、重复依赖和拆包机会。
- **中文 alt：** 前端 Bundle 分析报告展示 612 KB 首屏 JavaScript、图表大模块、重复日期依赖和拆包建议。
- **图注：** 体积优化应回到用户加载路径，不能只追求压缩后的总数字更小。
- **复现来源：** `/visual-demos/p1/index.html?scene=bundle-analysis`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 入口体积、模块占比、重复依赖和行动建议一致；48,502 bytes。

### Module Federation 运行链路

<!-- asset: /images/engineering/module-federation-runtime.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/engineering/module-federation`；**教学目的：** 展示宿主发现远程入口、协商共享依赖和加载模块的运行过程。
- **中文 alt：** Module Federation 宿主路由获取 remoteEntry、协商共享依赖、加载远程模块并进入错误边界。
- **图注：** 远程加载失败必须落入可恢复错误边界，不能让整个宿主页白屏。
- **复现来源：** `/visual-demos/p1/index.html?scene=module-federation-runtime`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** remoteEntry、share scope、remote module 和失败分支顺序准确；44,046 bytes。

### 前端工程交付流水线

<!-- asset: /images/engineering/project-pipeline.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/engineering/project-from-zero`；**教学目的：** 把代码检查、测试、构建、预览和渐进发布串成完整交付门禁。
- **中文 alt：** 前端工程流水线依次执行静态检查、自动测试、构建产物、预览验收和渐进发布。
- **图注：** 任一门禁失败都应停止推进，并保留可定位的产物和日志。
- **复现来源：** `/visual-demos/p1/index.html?scene=project-pipeline`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 流水线顺序、失败阻断和回滚点明确；42,270 bytes。

### Docker Compose 容器状态

<!-- asset: /images/devops/docker-container-state.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/devops/docker`；**教学目的：** 教读者区分长期服务、一次性迁移任务和反复重启的异常容器。
- **中文 alt：** Docker Compose 状态报告展示 PostgreSQL 和 API healthy、迁移 exited 0、Worker 反复 restarting。
- **图注：** `exited 0` 对迁移容器可能是成功，容器状态必须结合职责解释。
- **复现来源：** `/visual-demos/p1/index.html?scene=docker-container-state`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** healthy、exited 和 restarting 三类状态解释准确；46,322 bytes。

### 可观测性运行看板

<!-- asset: /images/devops/observability-dashboard.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/devops/observability`；**教学目的：** 把指标、日志和 Trace 组合成一次延迟问题的定位路径。
- **中文 alt：** 可观测性看板展示请求量、错误率、P95、数据库连接池饱和和慢查询 Trace。
- **图注：** 单个指标只能提示异常，根因判断需要跨信号建立时间和调用关联。
- **复现来源：** `/visual-demos/p1/index.html?scene=observability-dashboard`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** P95、连接池、慢查询和 Trace 链路可形成闭环；48,904 bytes。

### 金丝雀发布阶段

<!-- asset: /images/devops/deployment-canary.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/devops/deployment-strategy`；**教学目的：** 展示每次扩大流量前的观测门槛和失败回滚位置。
- **中文 alt：** 金丝雀发布从 1% 逐步扩大到 100% 流量，每个阶段检查错误率延迟和业务指标。
- **图注：** 技术指标正常仍不够，关键业务指标也必须进入发布判定。
- **复现来源：** `/visual-demos/p1/index.html?scene=deployment-canary`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 1%、10%、50%、100% 阶段、门禁和回滚方向明确；46,172 bytes。

### AI 文档问答评测报告

<!-- asset: /images/ai-engineering/evaluation-report.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/ai-engineering/evaluation`；**教学目的：** 展示 RAG 评测必须同时覆盖检索、事实依据、引用和安全。
- **中文 alt：** 文档问答评测报告展示 Retrieval at 5、Groundedness、Citation 和安全泄漏指标及失败样本。
- **图注：** 总分不能替代失败样本分析，每个指标都要能回到具体问题和证据。
- **复现来源：** `/visual-demos/p1/index.html?scene=evaluation-report`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 指标、阈值、失败样本和改进方向互相对应；45,512 bytes。

### AI 文档问答引用结果

<!-- asset: /images/ai-engineering/doc-qa-citations.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/ai-engineering/doc-qa-project`；**教学目的：** 对比有证据回答与证据不足拒答的产品状态。
- **中文 alt：** 文档问答结果展示有据回答、两条具体来源和证据不足时的拒答状态。
- **图注：** 引用必须能定位到真实文档片段，低置信度时应明确拒答而不是补猜。
- **复现来源：** `/visual-demos/p1/index.html?scene=doc-qa-citations`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 答案、来源、置信提示和拒答状态信息完整；53,280 bytes。

## P2 综合项目案例图

### Vue Admin 运营工作台

<!-- asset: /images/projects/vue-admin-overview.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/vue-admin`；**教学目的：** 展示指标、列表、待办、权限范围和状态反馈组成的管理台闭环。
- **中文 alt：** Vue Admin 运营工作台同时展示指标卡、用户列表、待办任务、权限范围和异常状态。
- **图注：** 数据、权限、页面状态和操作反馈必须形成闭环，不能只完成一张成功列表。
- **复现来源：** `/visual-demos/projects/index.html?scene=vue-admin-overview`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四项指标、三条用户记录、四条待办和状态边界清晰；54,636 bytes。

### 审批实例状态

<!-- asset: /images/projects/approval-workflow-state.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/approval-workflow-case`；**教学目的：** 同时呈现流程节点、当前任务、版本和审计记录。
- **中文 alt：** 采购审批实例展示已提交、直属上级、财务复核、部门负责人和结束节点，并列出当前任务与审计记录。
- **图注：** 审批动作必须校验任务状态与版本，并发处理时只能有一个请求成功。
- **复现来源：** `/visual-demos/projects/index.html?scene=approval-workflow-state`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 当前节点、处理人、剩余时间和四条审计记录一致；57,530 bytes。

### 销售数据看板

<!-- asset: /images/projects/analytics-dashboard-overview.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/analytics-dashboard-case`；**教学目的：** 让指标、趋势、排行、筛选范围和口径说明可同时核对。
- **中文 alt：** 销售经营看板展示成交额、订单量、转化率、退款率、成交趋势、渠道排行、筛选范围和指标口径。
- **图注：** 筛选上下文、数据时效与指标口径决定看板数字能否用于业务决策。
- **复现来源：** `/visual-demos/projects/index.html?scene=analytics-dashboard-overview`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四项指标、八日趋势、渠道排行和口径版本完整；52,170 bytes。

### 工作流设计器画布

<!-- asset: /images/projects/workflow-builder-canvas.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/workflow-builder-case`；**教学目的：** 展示画布、节点配置、校验错误、草稿和发布版本的关系。
- **中文 alt：** 报销审批工作流设计器画布包含开始、审批、条件和财务复核节点，右侧显示配置校验和未连接结束节点错误。
- **图注：** 设计态可以修改；发布生成不可变版本，运行实例继续使用启动时绑定的版本。
- **复现来源：** `/visual-demos/projects/index.html?scene=workflow-builder-canvas`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 四个节点、条件、配置项与发布阻断均清晰；58,356 bytes。

### 文件中心上传任务

<!-- asset: /images/projects/file-center-task.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/file-center-case`；**教学目的：** 区分分片上传、扫描、元数据确认、业务绑定和访问控制。
- **中文 alt：** 企业文件中心展示合同分片上传进度、设备清单病毒扫描状态、图片业务绑定成功以及访问与生命周期规则。
- **图注：** 对象存储成功后仍要确认元数据、扫描结果与业务绑定，文件才算可用。
- **复现来源：** `/visual-demos/projects/index.html?scene=file-center-task`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 三种文件状态、72% 进度和四条生命周期规则可辨；44,082 bytes。

### 消息中心未读与通道

<!-- asset: /images/projects/notification-center-unread.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/notification-center-case`；**教学目的：** 区分用户读状态与站内信、WebSocket、邮件、短信发送状态。
- **中文 alt：** 消息中心展示审批提醒、报表完成和权限更新消息，以及站内信、WebSocket、邮件和短信通道状态。
- **图注：** 用户消息与发送任务共享业务事件，但不应混用同一个状态字段。
- **复现来源：** `/visual-demos/projects/index.html?scene=notification-center-unread`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 已读/未读与四种通道状态没有混淆；49,634 bytes。

### 多租户权限范围

<!-- asset: /images/projects/multi-tenant-permission-scope.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/multi-tenant-permission-case`；**教学目的：** 展示租户上下文如何进入 SQL、缓存 key、数据范围和审计日志。
- **中文 alt：** 多租户权限中心把用户会话、租户标识和权限码映射到 SQL 租户条件、缓存 key 与审计日志。
- **图注：** 操作权限与数据范围必须一起校验，超级管理员也不能默认取消隔离。
- **复现来源：** `/visual-demos/projects/index.html?scene=multi-tenant-permission-scope`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 输入上下文、三类隔离目标和四条权限检查结果准确；50,226 bytes。

### 财务对账差异处理

<!-- asset: /images/projects/finance-reconciliation-exception.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/finance-reconciliation-case`；**教学目的：** 区分金额、缺失、跨日退款与状态差异的处理方法。
- **中文 alt：** 财务对账中心展示账单批次、平账比例、差异金额以及金额差异、本地缺失、跨日退款和状态差异记录。
- **图注：** 人工调整不能覆盖原始流水，必须留下原值、调整值、原因和审批证据。
- **复现来源：** `/visual-demos/projects/index.html?scene=finance-reconciliation-exception`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 批次指标、四类差异、渠道/本地值与建议动作一致；53,464 bytes。

### 风控人工复核

<!-- asset: /images/projects/risk-control-case-review.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/risk-control-center-case`；**教学目的：** 展示决策分数、命中原因、处置动作和人工复核证据。
- **中文 alt：** 提现风险复核台展示风险分数、命中规则、正负分值、冻结处置、补充验证和人工复核时间线。
- **图注：** 风控和业务系统使用同一 decision_id，才能解释决策并追踪实际处置。
- **复现来源：** `/visual-demos/projects/index.html?scene=risk-control-case-review`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 82 分、四项特征贡献和四步处置时间线完整；50,102 bytes。

### 跨区域灾备切换

<!-- asset: /images/projects/disaster-recovery-switch.webp | type: project-rendered | license: project-owned | status: verified -->

- **使用页面：** `/projects/disaster-recovery-case`；**教学目的：** 同时呈现主备健康、复制延迟、流量方向和切换门禁。
- **中文 alt：** 灾备控制台展示上海主区域异常、北京备区域就绪、复制延迟 38 秒和七步切换检查清单。
- **图注：** 故障切换不是一次 DNS 修改，每一步都需要证据、负责人和回退方式。
- **复现来源：** `/visual-demos/projects/index.html?scene=disaster-recovery-switch`，1440 × 900、DPR 1、`cwebp -q 86`。
- **工具与日期：** 仓库内 HTML/CSS 教学画布、Codex in-app Browser、cwebp，2026-07-21。
- **人工核对：** 主备状态、RPO 38 秒、步骤 4/7 与六项清单明确；55,762 bytes。

> 2026-07-21 曾三次尝试使用 built-in image generation 生成无文字抽象插图，但服务请求均因网络错误未产生资产。为保证可复现性，最终改用仓库内 HTML/CSS 教学画布渲染，未使用外部来源或第三方版权素材。
