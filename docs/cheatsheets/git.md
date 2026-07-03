# Git 速查

## 常用命令

```bash
git status
git add .
git commit -m "feat: add user page"
git pull
git push
```

| 命令 | 用途 |
| --- | --- |
| `git status` | 查看当前状态 |
| `git diff` | 查看未暂存改动 |
| `git diff --staged` | 查看已暂存改动 |
| `git add <file>` | 暂存文件 |
| `git commit -m` | 提交 |
| `git log --oneline` | 查看提交历史 |

## 分支

创建并切换分支：

```bash
git switch -c feat/user-management
```

切换分支：

```bash
git switch main
```

查看分支：

```bash
git branch
```

删除本地分支：

```bash
git branch -d feat/user-management
```

## 拉取和推送

拉取远端：

```bash
git pull
```

推送当前分支：

```bash
git push
```

第一次推送新分支：

```bash
git push -u origin feat/user-management
```

查看远端：

```bash
git remote -v
```

## 撤销改动

取消暂存：

```bash
git restore --staged file.md
```

恢复某个未暂存文件：

```bash
git restore file.md
```

恢复命令会丢弃本地改动。执行前先 `git diff` 确认。

修改上一条提交信息：

```bash
git commit --amend
```

## 合并和变基

合并：

```bash
git merge main
```

变基：

```bash
git rebase main
```

简单理解：

| 操作 | 特点 |
| --- | --- |
| `merge` | 保留分支合并记录 |
| `rebase` | 历史更线性，但会改写提交 |

团队协作中，不要随意 rebase 已经推送且别人正在使用的共享分支。

## 解决冲突

冲突标记：

```text
<<<<<<< HEAD
当前分支内容
=======
合入分支内容
>>>>>>> main
```

处理步骤：

```text
1. 打开冲突文件。
2. 保留正确内容，删除冲突标记。
3. 运行测试或构建。
4. git add 冲突文件。
5. 继续 merge 或 rebase。
```

继续 rebase：

```bash
git rebase --continue
```

放弃 rebase：

```bash
git rebase --abort
```

## stash

临时保存改动：

```bash
git stash push -m "work in progress"
```

查看 stash：

```bash
git stash list
```

恢复：

```bash
git stash pop
```

适合临时切分支、拉代码前保存半成品。不要把 stash 当长期备份。

## 提交信息建议

```text
feat: add user management page
fix: handle expired token once
docs: update deployment checklist
refactor: split request service
test: add permission guard tests
chore: update dependencies
```

| 类型 | 用途 |
| --- | --- |
| `feat` | 新功能 |
| `fix` | 修复问题 |
| `docs` | 文档 |
| `refactor` | 重构 |
| `test` | 测试 |
| `chore` | 工程杂项 |

## 常见问题

| 问题 | 处理 |
| --- | --- |
| 不知道改了什么 | `git status` 和 `git diff` |
| commit 前漏文件 | `git status` 检查 |
| 合并冲突 | 手动处理后继续 merge/rebase |
| 分支太久没同步 | 先拉 main，再解决冲突 |
| 不小心提交敏感信息 | 立即通知团队并轮换密钥 |

## 项目建议

- 每个任务使用独立分支。
- 提交前运行检查和构建。
- PR 描述写清变更、验证和风险。
- 不把 `.env`、密钥、构建产物随意提交。
- 合并前确认文档和代码同步更新。

## 下一步学习

- [项目交付检查清单](/projects/delivery-checklist)
- [文档治理](/contribute/governance)
- [质量检查](/contribute/quality-check)
