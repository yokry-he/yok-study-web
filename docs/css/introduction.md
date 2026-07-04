# CSS 学习导览

## 适合谁看

适合已经能写简单样式，但在真实项目里经常遇到布局错位、响应式困难、组件库样式被污染、移动端横向滚动等问题的学习者。

CSS 学习不能只背属性。真实项目更需要理解布局模型、响应式策略、样式边界和排错方法。

## 你会学到什么

- 盒模型和布局基本规则。
- Flex 与 Grid 的使用场景。
- 响应式设计如何从一开始规划。
- 动画和过渡如何服务状态变化，而不是干扰用户。
- 焦点、对比度、减少动态效果等 CSS 可访问性基础。
- 设计 token 和主题变量如何支撑长期维护。
- Vue 项目中如何组织样式。
- 如何避免污染组件库内部 DOM。
- 常见样式问题如何定位和修复。

## 学习顺序

<LearningPath :steps="[
  { title: '图解 CSS 核心概念', description: '先用图理解选择器、层叠、盒模型、布局、响应式、主题和排错路径。', link: '/css/visual-guide', badge: '图解' },
  { title: '盒模型与布局基础', description: '理解 content、padding、border、margin、display 和尺寸计算。', link: '/css/box-model-layout', badge: '基础' },
  { title: 'Flex 与 Grid', description: '掌握一维布局和二维布局，能做工具栏、卡片网格和表单布局。', link: '/css/flex-grid', badge: '核心' },
  { title: '响应式设计', description: '处理桌面、平板、移动端的信息优先级和布局变化。', link: '/css/responsive', badge: '移动端' },
  { title: '动画与过渡', description: '用 transition、animation 和 reduced motion 做克制、稳定、可访问的状态变化。', link: '/css/animation-transition', badge: '动效' },
  { title: 'CSS 可访问性', description: '处理焦点可见、颜色对比、视觉隐藏、响应式放大和状态表达。', link: '/css/accessibility', badge: '可访问' },
  { title: '设计 Token 与主题', description: '把颜色、间距、圆角、阴影和组件库主题沉淀成可维护变量。', link: '/css/design-tokens', badge: '系统' },
  { title: '项目样式架构', description: '建立全局样式、业务 class、设计变量和组件库边界。', link: '/css/architecture', badge: '工程' },
  { title: 'CSS 真实项目问题库', description: '排查横向溢出、组件库污染、固定元素变形、表格压缩、层级遮挡、移动端导航、动画卡顿和主题变量问题。', link: '/projects/issues-css', badge: '问题库' },
  { title: '常见问题', description: '排查横向溢出、头像变形、表格压缩、组件库样式污染等问题。', link: '/css/troubleshooting', badge: '排错' }
]" />

## CSS 在项目中的定位

CSS 负责界面呈现，但它不能没有边界。真实项目中，样式问题往往不是某个属性不会写，而是：

- 全局样式影响了组件库。
- 业务 class 不清晰。
- 固定尺寸元素没有稳定尺寸。
- 响应式只在最后临时补。
- 页面布局没有考虑内容增长。
- 动画只追求好看，没有考虑性能和减少动态效果。
- 深色模式、品牌换肤和组件库主题没有统一 token。

## 学习建议

先掌握布局，再追求视觉细节。对于后台、SaaS、文档站这类产品，样式目标应该是：

- 清晰。
- 稳定。
- 响应式可用。
- 焦点和状态清楚。
- 动效克制。
- 不污染组件库。
- 方便长期维护。

## 下一步

从 [图解 CSS 核心概念](/css/visual-guide) 开始，再进入 [盒模型与布局基础](/css/box-model-layout)。如果你已经在项目里遇到样式异常，直接看 [CSS 真实项目问题库](/projects/issues-css)。
