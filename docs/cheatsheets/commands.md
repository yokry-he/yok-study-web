# 常用命令速查

## 项目启动

| 命令 | 用途 |
| --- | --- |
| `npm install` | 安装依赖 |
| `npm run dev` | 启动开发服务 |
| `npm run build` | 生产构建 |
| `npm run preview` | 预览构建产物 |
| `npm run test` | 运行测试 |
| `npm run lint` | 代码检查 |
| `npm run typecheck` | 类型检查 |

优先使用项目 `package.json` 里已有 scripts，不要凭经验猜命令。

## VitePress

| 命令 | 用途 |
| --- | --- |
| `npm run docs:dev -- --port 6173` | 指定端口启动文档站 |
| `npm run docs:build` | 构建文档站 |
| `npm run docs:preview` | 预览构建产物 |
| `npm run docs:check` | 运行文档检查 |

本项目约定开发端口为 `6173`。

## Git

| 命令 | 用途 |
| --- | --- |
| `git status` | 查看改动 |
| `git diff` | 查看未暂存差异 |
| `git diff --cached` | 查看已暂存差异 |
| `git add file` | 暂存文件 |
| `git commit -m \"message\"` | 提交 |
| `git log --oneline -5` | 查看最近提交 |
| `git branch` | 查看分支 |

执行会丢改动的命令前，先看 `git status` 和 `git diff`。

## 查文件和文本

| 命令 | 用途 |
| --- | --- |
| `rg \"keyword\"` | 搜索文本 |
| `rg --files` | 列出文件 |
| `rg -n \"keyword\" docs` | 带行号搜索 |
| `sed -n '1,80p' file` | 查看指定行 |
| `nl -ba file` | 带行号查看 |
| `wc -l file` | 行数统计 |

项目中优先用 `rg`，比 `grep` 更适合代码库搜索。

## 端口和进程

| 命令 | 用途 |
| --- | --- |
| `lsof -i :6173` | 查端口占用 |
| `ps aux | grep node` | 查 Node 进程 |
| `kill <pid>` | 停止进程 |
| `node -v` | 查看 Node 版本 |
| `npm -v` | 查看 npm 版本 |

## HTTP 检查

如果有 `curl`：

```bash
curl -I http://127.0.0.1:6173
```

如果没有 `curl`：

```bash
node -e \"fetch('http://127.0.0.1:6173').then(r=>console.log(r.status))\"
```

批量检查：

```bash
node - <<'NODE'
const paths = ['/','/technologies/']
for (const path of paths) {
  const response = await fetch(`http://127.0.0.1:6173${path}`)
  console.log(response.status, path)
}
NODE
```

## 构建前检查顺序

```text
安装依赖
↓
类型检查
↓
lint
↓
测试
↓
构建
↓
HTTP 预览
```

## 危险命令提醒

谨慎执行：

| 命令 | 风险 |
| --- | --- |
| `rm -rf` | 删除文件不可恢复 |
| `git reset --hard` | 丢弃本地改动 |
| `git clean -fd` | 删除未跟踪文件 |
| `kill -9` | 强制结束进程，可能丢状态 |
| `docker system prune -a` | 清理镜像和缓存 |

## 延伸学习

- [Git 速查](/cheatsheets/git)
- [Linux 速查](/cheatsheets/linux)
- [Vite 速查](/cheatsheets/vite)
