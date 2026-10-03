# 快速开始

Gemsnote 支持 PostgreSQL 和 MongoDB 两种数据库。新部署推荐 PostgreSQL；需要继续使用旧 Leanote 数据时可以直接连接 Leanote MongoDB，或先迁移到 PostgreSQL。正式部署前应修改 `app.secret`、数据库密码和站点地址。

## 一、从 Leanote 迁移

迁移前停止旧 Leanote 写入，并备份 MongoDB 业务库和 `files/`。只导出 Leanote 业务库即可，不需要迁移 `admin`、`config`、`local` 数据库。

### Release 包中的迁移工具

Release 包包含编译好的 `bin/gemsnote-migrate`（Windows 为 `.exe`）。先准备空的 PostgreSQL 数据库，然后执行：

```bash
bin/gemsnote-migrate \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/leanote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable'
```

迁移及校验成功后，将附件和图片文件复制到 Gemsnote 的 `files/`、`public/upload/`。
使用 Release 部署时，将 `conf/app.conf` 指向本次迁移的目标 PostgreSQL，不要另建业务库。
完整参数、校验及失败处理见[数据库迁移指南](development/MIGRATION_GUIDE.md)。

下面的默认 Docker 快速启动流程用于全新安装：它会启动自己的 PostgreSQL，应用连接
`postgres` 容器，不会自动使用上面 `127.0.0.1:5432` 的已迁移数据库。需要 Docker 部署
迁移数据时，应按[部署指南中的 Docker 迁移流程](DEPLOYMENT.md#docker-环境迁移-mongodb-到-postgresql)
准备容器目标库并向其中迁移，不要直接执行默认 `docker compose up -d`。
若已迁移到外部 PostgreSQL 并打算继续使用它，按同节的外部库说明调整配置。

## 二、使用 Release 包启动

下载与系统和 CPU 架构匹配的服务端 Release 包并解压。创建持久化目录：

```bash
mkdir -p files public/upload
```

使用 PostgreSQL 时先创建空数据库：

```sql
CREATE USER gemsnote WITH PASSWORD '请替换为强密码';
CREATE DATABASE gemsnote OWNER gemsnote;
```

编辑 `conf/app.conf`，选择数据库并填写配置。PostgreSQL 示例：

```ini
db.type=postgresql
db.host=127.0.0.1
db.port=5432
db.dbname=gemsnote
db.username=gemsnote
db.password=请替换为强密码
app.secret=请替换为随机长字符串
site.url=http://127.0.0.1:9000
```

MongoDB 示例：

```ini
db.type=mongodb
db.host=127.0.0.1
db.port=27017
db.dbname=gemsnote
db.username=
db.password=
app.secret=请替换为随机长字符串
```

连接旧 Leanote MongoDB 时将 `db.dbname` 改为 `leanote`，或直接配置完整的 `db.url`。新 MongoDB 实例可使用 Release 包中的初始化快照；PostgreSQL 空数据库会在首次启动时自动执行 schema 和 seed。

使用 Release 包内的 MongoDB 初始化快照时，先执行 `mongorestore --drop --db gemsnote <快照目录>`，再启动服务。恢复旧 Leanote 备份时只恢复 `leanote` 业务库；如果备份是 MongoDB 的物理 `db` 目录，必须先启动兼容版本的 MongoDB，再导出逻辑备份。

启动服务：

```bash
./run.sh
```

Windows 执行 `run.bat`。本机可访问 <http://127.0.0.1:9000>，但 Release 默认配置为
`http.addr=0.0.0.0`、`http.port=9000`，即监听所有 IPv4 网卡，并非仅限本机。
仅本机试用或使用同机反向代理时，将 `conf/app.conf` 的 `http.addr` 改为 `127.0.0.1`。
需要远程访问时配合防火墙及 HTTPS 反向代理限制暴露范围；`site.url` 不控制监听地址。

## 三、使用 GHCR Docker 镜像启动

需要 Docker Engine 和 Compose 插件，不需要 Go 或 Node.js。从对应版本 GitHub Release 下载 **Source code** 压缩包并解压，进入包含 `docker-compose.yml` 的目录（这里只使用部署文件，不编译源码）。至少保留 `docker-compose.yml`、`conf/app.docker-postgres.conf`、`database/schema.sql` 和 `database/seed.sql`。

首次启动前，在该目录创建 `.env`：

```dotenv
GEMSNOTE_VERSION=1.0.0
GEMSNOTE_POSTGRES_PASSWORD=请替换为强密码
```

编辑 `conf/app.docker-postgres.conf`：`db.password` 必须与上面的密码一致，保留 `db.host=postgres`，并修改 `app.secret` 为随机长字符串、`site.url` 为实际访问地址。`.env` 只供 Compose 插值，不会自动替换应用配置中的密码；不要提交包含密码的配置文件。

```bash
mkdir -p files public/upload
docker compose pull
docker compose up -d
docker compose logs -f gemsnote
```

首次初始化时 PostgreSQL 自动执行 schema 和 seed。默认 Compose 直接使用 GHCR 镜像，没有 `build` 配置，无需手工重命名镜像。MongoDB 和迁移部署见[详细部署指南](DEPLOYMENT.md)。

Compose 会挂载 PostgreSQL 配置、初始化 SQL、`files/` 和 `public/upload/`，数据库持久化在 `data/`。启动前确认配置和 SQL 已存在且是普通文件。端口仅绑定宿主机 `127.0.0.1:9000`；远程访问应配置 HTTPS 反向代理。

## Web 使用

访问 <http://127.0.0.1:9000>，全新内置数据的默认账号为 `admin`，默认密码为 `gemsnote`。迁移或直接连接旧库时保留原账号密码。首次登录后立即修改默认密码。服务版本可通过以下接口检查：

```bash
curl http://127.0.0.1:9000/api2/system/version
```

附件和图片位于 `files/`、`public/upload/`，升级或重建容器时不要删除这些目录。更多生产部署、测试和自动化构建说明见[部署与交付文档](DEPLOYMENT.md)。
