# Linux 速查

## 常用目录

| 路径 | 用途 |
| --- | --- |
| `/etc` | 系统和服务配置 |
| `/var/log` | 日志 |
| `/var/www` | 常见 Web 部署目录 |
| `/usr/local/bin` | 本机安装命令 |
| `/tmp` | 临时文件 |
| `~` | 当前用户目录 |

## 文件和目录

| 命令 | 用途 |
| --- | --- |
| `pwd` | 查看当前目录 |
| `ls -lah` | 查看文件和权限 |
| `cd /path` | 切换目录 |
| `mkdir -p logs/app` | 创建多级目录 |
| `cp source target` | 复制文件 |
| `cp -r source target` | 复制目录 |
| `mv old new` | 移动或重命名 |
| `rm file` | 删除文件 |
| `rm -rf dir` | 删除目录，危险操作 |
| `touch app.log` | 创建空文件或更新时间 |

执行 `rm -rf` 前先 `pwd` 和 `ls`，确认当前目录。

## 查看文件

| 命令 | 用途 |
| --- | --- |
| `cat file` | 输出完整文件 |
| `less file` | 分页查看 |
| `head -n 50 file` | 查看前 50 行 |
| `tail -n 100 file` | 查看后 100 行 |
| `tail -f app.log` | 实时跟踪日志 |
| `grep \"error\" app.log` | 搜索文本 |
| `grep -R \"keyword\" .` | 递归搜索 |

项目排错优先用 `tail -f` 看实时日志，用 `grep` 定位错误关键字。

## 权限

| 命令 | 用途 |
| --- | --- |
| `whoami` | 当前用户 |
| `id` | 用户和用户组 |
| `chmod +x deploy.sh` | 赋予执行权限 |
| `chmod 644 file` | 文件常见权限 |
| `chmod 755 dir` | 目录常见权限 |
| `chown user:group file` | 修改所有者 |
| `sudo command` | 使用管理员权限执行 |

权限常见问题：

- Nginx 403：文件权限、目录权限或用户不匹配。
- 脚本不能执行：缺少 `+x`。
- 服务读不到配置：运行用户没有读权限。

## 进程和端口

| 命令 | 用途 |
| --- | --- |
| `ps aux | grep node` | 查 Node 进程 |
| `top` | 查看系统资源 |
| `kill <pid>` | 结束进程 |
| `kill -9 <pid>` | 强制结束，谨慎使用 |
| `lsof -i :6173` | 查看端口占用 |
| `ss -lntp` | 查看监听端口 |

端口占用排查：

```bash
lsof -i :6173
kill <pid>
```

## 磁盘和内存

| 命令 | 用途 |
| --- | --- |
| `df -h` | 查看磁盘空间 |
| `du -sh *` | 查看当前目录各项大小 |
| `free -h` | 查看内存 |
| `top` | 查看 CPU 和内存进程 |

磁盘满常见影响：

- 日志写不进去。
- 数据库异常。
- 构建失败。
- 上传失败。

## 网络

| 命令 | 用途 |
| --- | --- |
| `ping example.com` | 检查连通性 |
| `curl -I https://example.com` | 查看响应头 |
| `curl http://127.0.0.1:3000/health` | 检查健康接口 |
| `dig example.com` | 查 DNS |
| `ip addr` | 查看本机 IP |

本环境如果没有 `curl`，可以用 Node `fetch` 替代：

```bash
node -e \"fetch('http://127.0.0.1:6173').then(r=>console.log(r.status))\"
```

## 服务管理

| 命令 | 用途 |
| --- | --- |
| `systemctl status nginx` | 查看服务状态 |
| `systemctl restart nginx` | 重启服务 |
| `systemctl reload nginx` | 重新加载配置 |
| `journalctl -u nginx -n 100` | 查看服务日志 |

修改 Nginx 配置后先测试：

```bash
nginx -t
systemctl reload nginx
```

## 项目排错顺序

```text
进程是否存在
↓
端口是否监听
↓
健康接口是否正常
↓
日志是否报错
↓
磁盘和内存是否异常
↓
Nginx 或网关是否代理正确
```

## 参考资料

- [GNU Coreutils](https://www.gnu.org/software/coreutils/)

## 延伸学习

- [Linux 与 Shell 基础](/devops/linux-shell)
- [DevOps 常见问题](/devops/troubleshooting)
