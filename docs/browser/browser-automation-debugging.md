# 浏览器自动化调试

## 适合谁看

适合已经会用 DevTools 排查问题，但希望把浏览器验证变得更稳定、更可重复的人：

- 每次改完页面都靠手动点，容易漏测。
- 想自动检查页面是否白屏、路由是否能打开、关键按钮是否存在。
- 想在 CI 里验证登录、表单、上传、下载、权限、响应式。
- 想理解 Playwright、Puppeteer、Chrome DevTools Protocol 的区别。

浏览器自动化调试不是只写端到端测试。它还可以用来做本地开发验证、故障复现、截图对比、性能采样和线上巡检。

## 为什么需要自动化

前端问题有一个特点：代码看起来对，不代表浏览器里真的对。

常见人工检查容易漏掉：

- 路由刷新 404。
- 页面局部白屏。
- 移动端横向滚动。
- 弹窗被遮挡。
- 按钮被压缩。
- 登录态丢失。
- 控制台有错误但页面看似正常。
- 构建产物路径线上不一致。

自动化调试的目标是把这些检查变成脚本，每次都按同样步骤执行。

## 工具分层

| 工具 | 定位 |
| --- | --- |
| Chrome DevTools | 手动调试、性能分析、网络排查 |
| Chrome DevTools Protocol | 底层浏览器调试协议 |
| Puppeteer | 基于 Chromium 的自动化库 |
| Playwright | 跨 Chromium、Firefox、WebKit 的自动化和测试框架 |
| Lighthouse | 性能、可访问性、最佳实践审计 |

如果你是普通前端项目，优先学 Playwright。它覆盖测试、脚本、截图、移动端模拟和多浏览器验证，落地成本更低。

## 最小 Playwright 脚本

安装后可以写一个最小检查：

```ts
import { chromium } from 'playwright'

const browser = await chromium.launch()
const page = await browser.newPage()

await page.goto('http://127.0.0.1:6173/browser/introduction')

const title = await page.locator('h1').textContent()

if (!title?.includes('浏览器学习导览')) {
  throw new Error('页面标题不符合预期')
}

await browser.close()
```

这个脚本验证的是：浏览器能打开页面，并且页面内容不是空的。

## 检查白屏

白屏不一定会导致 HTTP 失败。页面可能返回 200，但 JavaScript 运行时异常。

可以同时检查：

- HTTP 状态码。
- `h1` 或主内容是否存在。
- 控制台是否有 error。
- 页面是否有大面积空白。
- 核心容器是否有高度。

示例：

```ts
const errors: string[] = []

page.on('console', message => {
  if (message.type() === 'error') {
    errors.push(message.text())
  }
})

const response = await page.goto('http://127.0.0.1:6173/vue/introduction')

if (!response?.ok()) {
  throw new Error(`页面请求失败：${response?.status()}`)
}

await page.locator('h1').waitFor()

if (errors.length > 0) {
  throw new Error(errors.join('\n'))
}
```

这比只看 `curl 200` 更接近真实用户体验。

## 检查移动端横向溢出

很多页面桌面端正常，移动端出现横向滚动条。

```ts
await page.setViewportSize({
  width: 390,
  height: 844
})

await page.goto('http://127.0.0.1:6173/browser/introduction')

const hasOverflow = await page.evaluate(() => {
  return document.documentElement.scrollWidth > window.innerWidth
})

if (hasOverflow) {
  throw new Error('页面存在横向溢出')
}
```

文档站、后台系统和表格页面都应该把这个检查纳入常规验证。

## 检查关键交互

比如搜索框、导航、按钮、弹窗：

```ts
await page.goto('http://127.0.0.1:6173/')

await page.getByRole('button', { name: /搜索/ }).click()
await page.getByPlaceholder(/搜索/).fill('Vue')

await page.keyboard.press('Enter')

await page.getByText('Vue').first().waitFor()
```

优先使用面向用户的定位方式：

- `getByRole`
- `getByLabel`
- `getByPlaceholder`
- `getByText`

少用依赖 DOM 结构的复杂 CSS 选择器。DOM 结构变了，不代表用户行为变了。

## 检查网络请求

自动化脚本可以配合 Network 验证接口行为。

```ts
const userResponse = page.waitForResponse(response => {
  return response.url().includes('/api/user') && response.status() === 200
})

await page.getByRole('button', { name: '刷新用户' }).click()

await userResponse
```

也可以拦截接口：

```ts
await page.route('**/api/user', route => {
  route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify({ name: 'Tom' })
  })
})
```

这适合验证边界状态，例如空列表、接口 500、权限不足、网络超时。

## 和 DevTools 的关系

手动 DevTools 适合定位一次问题：

- Network 看请求。
- Application 看 Cookie、Storage、Cache。
- Performance 看长任务。
- Console 看运行时错误。

自动化适合把问题变成回归检查：

```text
手动发现问题
↓
写出复现步骤
↓
变成 Playwright 脚本
↓
以后每次构建都跑
```

不要只靠自动化，也不要只靠人工。正确做法是：人工定位根因，自动化防止复发。

## Chrome DevTools Protocol 是什么

Chrome DevTools Protocol，简称 CDP，是 Chrome 和 Chromium 暴露的一套底层调试协议。

它可以做：

- 页面导航。
- DOM 检查。
- 网络监听。
- 性能采集。
- 截图。
- 控制台日志读取。
- 调试器能力。

Puppeteer 和很多调试工具都建立在 CDP 能力之上。普通项目不需要直接写 CDP，但理解它有助于知道浏览器自动化的底层原理。

## 自动化调试的常见用例

### 路由可访问性巡检

文档站可以把重要路由逐个打开：

```ts
const routes = [
  '/',
  '/vue/introduction',
  '/browser/introduction',
  '/projects/real-world-issues'
]

for (const route of routes) {
  const response = await page.goto(`http://127.0.0.1:6173${route}`)

  if (!response?.ok()) {
    throw new Error(`${route} 打开失败`)
  }

  await page.locator('h1').waitFor()
}
```

### 构建后产物验证

不要只运行 build。build 成功只能说明静态产物生成成功，不代表路由、资源和交互都能运行。

建议至少验证：

- 首页。
- 核心模块入口。
- 动态路由或真实路由刷新。
- 搜索、导航、主题切换等全局交互。
- 移动端宽度。

### 线上问题复现

把用户操作步骤写成脚本：

```text
打开页面
↓
登录
↓
进入订单列表
↓
筛选状态
↓
点击详情
↓
检查金额和按钮状态
```

脚本能复现，说明问题可稳定观察；脚本不能复现，说明还缺环境、数据或权限条件。

## 实际项目常见问题

### 1. 自动化脚本本地通过，CI 失败

**常见原因**

- CI 机器慢，等待条件不稳定。
- 没有安装浏览器依赖。
- 依赖真实接口数据。
- 使用了固定时间等待。

**解决方案**

- 使用 `locator.waitFor()` 或断言等待。
- 使用接口 mock 或测试数据。
- 不依赖随机线上数据。
- CI 中保留截图、trace 和视频。

### 2. 脚本经常 flaky

**原因**

测试步骤依赖不稳定 DOM、动画、异步请求或固定延时。

**解决方案**

- 优先用 role、label、text 定位。
- 等待明确业务状态。
- 避免 `waitForTimeout`。
- 把复杂流程拆成稳定步骤。

### 3. 截图对比总是失败

**原因**

字体、系统、时间、动画、抗锯齿、数据内容都会影响截图。

**解决方案**

- 固定 viewport。
- 固定测试数据。
- 关闭或等待动画完成。
- 避免把动态时间和随机内容纳入截图区域。

### 4. 登录流程拖慢所有测试

**解决方案**

- 用 API 准备登录态。
- 复用 storage state。
- 把登录流程本身单独保留少量测试。
- 其他用例直接进入已登录状态。

### 5. 选择器一改就全挂

**原因**

脚本依赖了组件库内部 DOM 或复杂层级。

**解决方案**

- 优先使用可访问名称和语义角色。
- 必要时给业务关键节点加稳定 `data-testid`。
- 不要使用 `.n-button > span > span` 这类内部结构选择器。

## 最佳实践

- 每个关键页面至少检查 HTTP 状态、标题、控制台错误和移动端溢出。
- 自动化脚本要贴近用户行为，不要贴近组件内部 DOM。
- 能 mock 的接口尽量 mock，端到端全链路只保留核心路径。
- 失败时保留截图、trace、控制台日志和网络信息。
- 对文档站、官网、后台首页建立路由巡检。
- 把线上修复过的问题沉淀成回归脚本。

## 学习检查

学完本节后，你应该能回答：

- 为什么 build 成功不等于页面可用。
- Playwright、Puppeteer、CDP 分别适合什么。
- 如何用脚本检查白屏和横向溢出。
- 为什么要用面向用户的定位方式。
- 如何把线上问题转成回归验证。

## 参考资料

- [Playwright](https://playwright.dev/)
- [Playwright: Installation](https://playwright.dev/docs/intro)
- [Playwright: Browsers](https://playwright.dev/docs/browsers)
- [Chrome DevTools](https://developer.chrome.com/docs/devtools)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)

## 下一步学习

继续学习 [渲染与性能](/browser/rendering-performance)，把自动化检查和性能定位结合起来。
