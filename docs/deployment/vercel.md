# Vercel 部署记录

本文记录 `yok-study-web` 文档站部署到 Vercel 的配置、命令、结果和后续维护方式。

## 部署目标

- 项目类型：VitePress 静态文档站
- 部署平台：Vercel
- Vercel 团队：`lieyankuis-projects`
- Vercel 项目：`yok-study-web`
- 构建命令：`npm run docs:build`
- 静态输出目录：`docs/.vitepress/dist`

## 部署配置

项目根目录新增 `vercel.json`，显式声明 VitePress 的构建命令和输出目录：

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "buildCommand": "npm run docs:build",
  "outputDirectory": "docs/.vitepress/dist",
  "installCommand": "npm install"
}
```

项目根目录新增 `.vercelignore`，避免 CLI 部署时上传本地依赖、构建产物、缓存、环境变量和临时目录。

## 本次部署过程

时间：2026-07-15

### 1. 检查项目状态

执行命令：

```bash
pwd
git remote get-url origin
cat .vercel/project.json 2>/dev/null || cat .vercel/repo.json 2>/dev/null
vercel whoami 2>/dev/null
vercel teams list --format json 2>/dev/null
```

检查结果：

- 当前项目目录：`/Users/yokry/Documents/Codex/2026-07-01/yok-study-web`
- GitHub 远程地址：`git@github.com:yokry-he/yok-study-web.git`
- 初始状态没有 `.vercel/project.json` 或 `.vercel/repo.json`
- 初始状态系统未安装全局 `vercel` CLI

### 2. 本地验证

执行命令：

```bash
git diff --check
npm run docs:check
npm run docs:build
```

验证结果：

- `git diff --check` 通过
- `npm run docs:check` 通过
  - 初次部署前：Markdown 文档 505 篇，内部路由 505 个，配置路由 525 个
  - 新增部署记录后：Markdown 文档 506 篇，内部路由 506 个，配置路由 525 个
- `npm run docs:build` 通过
  - VitePress 版本：`1.6.4`
  - 初次构建完成时间约 14.68s
  - 新增部署配置后再次构建完成时间约 12.96s

构建时有两类非阻塞警告：

- npm 提示本机配置里的 `electron_mirror`、`electron_builder_binaries_mirror` 未来版本可能不再支持。
- VitePress 提示部分代码块语言 `env` 未加载，会按 `txt` 高亮；另有部分 chunk 超过 500 KB。

### 3. 登录 Vercel

全局安装 `vercel` CLI 超过两分钟无输出后终止，改用：

```bash
npx vercel@latest login
```

第一次登录失败：

```text
Error: An unexpected error occurred in login: TypeError: fetch failed
```

定位原因：当前 shell 中存在本地代理变量：

```text
http_proxy=http://127.0.0.1:33210
https_proxy=http://127.0.0.1:33210
all_proxy=socks5://127.0.0.1:33211
```

处理方式：仅对 Vercel CLI 命令临时移除代理环境变量：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest login
```

登录结果：

- 登录成功
- 当前账号：`lieyankui`
- 当前团队：`lieyankuis-projects`

### 4. 链接 Vercel 项目

执行命令：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest link --repo --scope lieyankuis-projects
```

交互选择：

- 确认将当前 Git 仓库链接到 Vercel 项目
- 创建项目：`lieyankuis-projects/yok-study-web`

链接结果：

- Vercel 项目创建成功
- 本地生成 `.vercel/repo.json`
- Vercel CLI 自动将 `.vercel` 写入 `.gitignore`

### 5. 执行第一次部署

执行命令：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest deploy --yes --no-wait --scope lieyankuis-projects
```

第一次部署结果：

- Deployment ID：`dpl_H6f91E295TjvQrUb3BvVXr1GWEbj`
- Target：`production`
- 状态：`Ready`
- 部署地址：`https://yok-study-jn8txugkc-lieyankuis-projects.vercel.app`
- 生产别名：`https://yok-study-web.vercel.app`
- 团队别名：`https://yok-study-web-lieyankuis-projects.vercel.app`
- 上传体积：约 115.8 MB

说明：虽然命令未显式使用 `--prod`，Vercel CLI 在当前已链接项目和默认分支上下文中将该部署标记为了 `production`。后续如果必须只创建 Preview，应在非生产分支或 Vercel 项目设置中调整目标分支策略后再部署。

### 6. 增加部署配置后再次部署

第一次部署后补充了 `vercel.json` 和 `.vercelignore`。其中 `.vercelignore` 用来避免上传本地依赖、已有构建产物、缓存和临时目录。

重新执行：

```bash
git diff --check
npm run docs:check
npm run docs:build
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest deploy --yes --no-wait --scope lieyankuis-projects
```

第二次部署结果：

- Deployment ID：`dpl_Huq7xCRq83aNruicKLubY5bRSrdE`
- Target：`preview`
- 状态：`Ready`
- Preview 地址：`https://yok-study-coyyb0fys-lieyankuis-projects.vercel.app`
- 预览别名：`https://yok-study-web-yokry-lieyankuis-projects.vercel.app`
- 上传体积：约 5 KB

这次部署是带有 `vercel.json` 和 `.vercelignore` 的最终验证版本。

### 7. 检查部署状态

执行命令：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest inspect yok-study-jn8txugkc-lieyankuis-projects.vercel.app --scope lieyankuis-projects
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest inspect yok-study-coyyb0fys-lieyankuis-projects.vercel.app --scope lieyankuis-projects
```

检查结果：

- Production 部署：从 `Building` 变为 `Ready`
- Preview 部署：从 `Building` 变为 `Ready`

## 后续维护

如果希望以后通过 GitHub 自动部署，需要在 Vercel 控制台确认仓库连接：

1. 项目：`yok-study-web`
2. 仓库：`yokry-he/yok-study-web`
3. Build Command：`npm run docs:build`
4. Output Directory：`docs/.vitepress/dist`
5. Install Command：`npm install`

手动预览部署：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest deploy --yes --scope lieyankuis-projects
```

手动生产部署：

```bash
env -u http_proxy -u https_proxy -u all_proxy -u HTTP_PROXY -u HTTPS_PROXY -u ALL_PROXY npx vercel@latest deploy --prod --yes --scope lieyankuis-projects
```
