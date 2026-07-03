# Linux 与 Shell 基础

## 适合谁看

适合第一次接触服务器，或者只会复制部署命令但不会排查问题的人。

学习 Linux 不需要一开始背很多命令。对 Web 开发者来说，优先掌握这些能力：

- 找文件。
- 看日志。
- 看进程。
- 看端口。
- 查磁盘。
- 改权限。
- 用 Shell 串起简单操作。

## 目录与文件

常见目录：

| 目录 | 常见用途 |
| --- | --- |
| `/var/www` | 静态站点目录 |
| `/etc/nginx` | Nginx 配置 |
| `/var/log` | 系统和服务日志 |
| `/home/<user>` | 普通用户目录 |
| `/opt` | 第三方应用安装目录 |
| `/tmp` | 临时文件 |

常用命令：

```bash
pwd
ls -lah
cd /var/www
mkdir -p releases
cp -r dist /var/www/app
rm -rf old-dist
```

危险命令要谨慎，尤其是：

```bash
rm -rf /
rm -rf *
chmod -R 777 .
```

不要在不确认当前目录时执行删除命令。

## 查看文件内容

```bash
cat nginx.conf
less app.log
tail -n 100 app.log
tail -f app.log
```

常见用法：

```bash
tail -f /var/log/nginx/error.log
```

`tail -f` 适合实时观察部署后是否报错。

## 搜索文本

推荐使用 `grep` 或 `rg`：

```bash
grep -R "proxy_pass" /etc/nginx
rg "VITE_API_BASE_URL" .
```

如果服务器没有 `rg`，用 `grep` 即可。

## 进程排查

查看进程：

```bash
ps aux | grep node
ps aux | grep nginx
```

查看资源占用：

```bash
top
```

如果使用 systemd：

```bash
systemctl status nginx
systemctl restart nginx
journalctl -u nginx -n 100
```

不要看到服务异常就立刻重启。先看日志，确认根因。重启可以恢复服务，但也可能掩盖问题。

## 端口排查

查看端口监听：

```bash
ss -lntp
```

常见判断：

| 现象 | 可能原因 |
| --- | --- |
| 80 没监听 | Nginx 没启动或配置加载失败 |
| 3000 没监听 | Node 服务没启动 |
| 端口被占用 | 旧进程未停止 |
| 本机能访问，外部不能访问 | 防火墙、安全组或代理配置问题 |

如果系统没有 `ss`，可以用：

```bash
lsof -i :80
```

## 权限问题

权限问题常见于：

- Nginx 读不到静态文件。
- Node 写不了日志。
- 上传目录不可写。
- Docker volume 挂载后用户不匹配。

查看权限：

```bash
ls -lah /var/www
```

修改属主：

```bash
chown -R www-data:www-data /var/www/app
```

修改权限：

```bash
chmod -R 755 /var/www/app
```

不要直接 `chmod -R 777`。这会把安全风险扩大到整个目录。

## 磁盘和内存

查看磁盘：

```bash
df -h
du -sh /var/www/*
```

查看内存：

```bash
free -h
```

构建失败、Docker 拉镜像失败、日志写入失败，都可能是磁盘满了。

## Shell 脚本基础

发布脚本应该显式失败，避免某一步失败后继续执行：

```bash
#!/usr/bin/env bash
set -euo pipefail

npm ci
npm run build
```

含义：

| 配置 | 作用 |
| --- | --- |
| `set -e` | 命令失败就退出 |
| `set -u` | 使用未定义变量时报错 |
| `set -o pipefail` | 管道中任意命令失败就视为失败 |

## 实际项目问题

### 问题：部署后 Nginx 访问 403

**常见原因**

- 目录没有读取权限。
- `root` 指向了错误目录。
- 缺少 `index.html`。
- SELinux 或容器权限限制。

**排查**

```bash
ls -lah /var/www/app
cat /var/log/nginx/error.log
```

### 问题：服务启动成功但外部访问不了

**排查**

1. 服务是否监听端口。
2. 是否监听 `127.0.0.1` 还是 `0.0.0.0`。
3. 防火墙是否放行。
4. 云服务器安全组是否放行。
5. Nginx 是否代理到正确端口。

## 最佳实践

- 每次改 Nginx 前先备份配置。
- 删除文件前先 `pwd` 和 `ls` 确认路径。
- 发布脚本开启 `set -euo pipefail`。
- 权限最小化，不要用 777 解决问题。
- 线上排错先保留证据，再重启服务。

## 下一步学习

继续学习 [Nginx 静态部署与代理](/devops/nginx)。
