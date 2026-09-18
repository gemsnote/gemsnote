# 快速开始

本文只介绍全新 PostgreSQL 安装。已有 Leanote 数据、MongoDB 部署和生产环境细节请参阅 [安装和部署指南](DEPLOYMENT.md)。

## 方式一：使用 Release 包

### 1. 准备 PostgreSQL

安装 PostgreSQL 12 或更高版本，然后创建用户和空数据库：

```sql
CREATE USER gemsnote WITH PASSWORD '请替换为强密码';
CREATE DATABASE gemsnote OWNER gemsnote;
```

下载与系统架构匹配的 Release 包并解压。Linux/macOS 包为
`gemsnote-<os>-<arch>-v<version>.tar.gz`，Windows 包为 `.zip`。进入解压后的
`gemsnote` 目录。Gemsnote 首次连接真正的空数据库时会自动执行包内的
`database/schema.sql` 和 `database/seed.sql`，不需要手工运行 `psql`。

### 2. 配置并启动

编辑 `conf/app.conf`：

```ini
db.type=postgresql
db.host=127.0.0.1
db.port=5432
db.dbname=gemsnote
db.username=gemsnote
db.password=请替换为强密码
```

按实际访问地址修改 `site.url`。默认 `app.secret` 技术上可以启动，但正式或生产
使用必须在首次启动前将它改为随机长字符串。创建持久文件目录后启动：

```bash
mkdir -p files public/upload
./run.sh
```

macOS 首次运行若被系统拦截，需要在“系统设置 → 隐私与安全性”中允许该程序。

Windows 在 PowerShell 或命令提示符中执行：

```bat
mkdir files
mkdir public\upload
run.bat
```

`run.bat` 首次启动会创建运行时目录联接；若创建失败，请以管理员身份运行一次，
或在 Windows 设置中启用开发人员模式。

## 方式二：使用 Docker Compose

需要 Docker Engine 和 Compose 插件。两个 Compose 文件已默认持久化图片、附件和旧上传
资源。创建持久化目录，再启动服务：

```bash
mkdir -p files public/upload
docker compose -f docker-compose.postgres.yml up -d --build
docker compose -f docker-compose.postgres.yml ps
```

当前 Compose 配置使用 PostgreSQL 18，数据库数据保存在 `./data`，实际 PGDATA
位于 `./data/18/docker`。仅当该 PostgreSQL 数据目录尚未初始化时，官方镜像才会
自动依次执行 `database/schema.sql` 和 `database/seed.sql`。默认数据库名、用户和
密码均为 `gemsnote`（兼容内置 seed 数据；首次登录后请立即修改）。

应用读取 `conf/app.docker-postgres.conf`，默认只在宿主机
`127.0.0.1:9000` 提供服务。正式使用前请至少修改：

- Compose 中的 `POSTGRES_PASSWORD`；
- `conf/app.docker-postgres.conf` 中对应的 `db.password`；
- `conf/app.docker-postgres.conf` 中的 `app.secret` 和 `site.url`。

`files/` 保存正文图片和附件，`public/upload/` 用于兼容头像等旧上传资源；上述
挂载可确保重建应用容器后文件仍然存在。

查看服务日志：

```bash
docker compose -f docker-compose.postgres.yml logs -f gemsnote
```

## 使用 Web 界面

默认访问地址为 <http://127.0.0.1:9000>，初始账号为：

```text
用户名：admin
密码：gemsnote
```

首次登录后立即在“账号”中修改密码。随后可以创建笔记本和笔记，或在“管理”中
配置系统。若修改了 `site.url` 或端口，请使用对应地址访问。

服务版本接口为：

```bash
curl http://127.0.0.1:9000/api/system/version
```

停止服务：

```bash
# Release：在前台按 Ctrl+C

# Docker（保留数据库和附件）
docker compose -f docker-compose.postgres.yml down
```

不要使用 `down -v`，也不要删除 `./data` 或 `./files`，除非确认不再需要其中的
数据。
