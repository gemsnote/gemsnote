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

迁移完成后，将附件和图片文件复制到 Gemsnote 的 `files/`、`public/upload/`，再按下面的 Release 或 Docker 方式配置 PostgreSQL 启动。完整参数和校验方式见[数据库迁移指南](development/MIGRATION_GUIDE.md)。

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

Windows 执行 `run.bat`。服务默认监听 `http://127.0.0.1:9000`。

## 三、使用 GHCR Docker 镜像启动

需要 Docker Engine 和 Compose 插件。先将与版本匹配的 GHCR 镜像拉取并标记为 Compose 文件使用的本地名称：

```bash
mkdir -p files public/upload
docker pull ghcr.io/gemsnote/gemsnote:1.0.0
docker tag ghcr.io/gemsnote/gemsnote:1.0.0 gemsnote:1.0.0
docker compose -f docker-compose.postgres.yml up -d
docker compose -f docker-compose.postgres.yml logs -f gemsnote
```

首次初始化时 PostgreSQL 自动执行 schema 和 seed，默认账号为 `admin`，默认密码为 `gemsnote`；首次登录后立即修改密码。正式使用前修改 Compose 中的数据库密码，以及 `conf/app.docker-postgres.conf` 的 `app.secret` 和 `site.url`。MongoDB 部署使用 `docker-compose.mongodb.yml`，先将上面的 GHCR 镜像标记为 `gemsnote:1.0.0`，再按其中的 MongoDB 配置连接或恢复 Leanote 备份。

Compose 会挂载 `conf/app.docker-postgres.conf`、`conf/app.docker-mongodb.conf`、`files/` 和 `public/upload/`。启动前确认配置文件已经存在且是普通文件；不要使用 `docker compose ... up --build`，否则会重新构建本地镜像而不是使用 GHCR 镜像。

## Web 使用

访问 <http://127.0.0.1:9000>，默认账号为 `admin`，默认密码为 `gemsnote`。首次登录后立即修改密码。服务版本可通过以下接口检查：

```bash
curl http://127.0.0.1:9000/api/system/version
```

附件和图片位于 `files/`、`public/upload/`，升级或重建容器时不要删除这些目录。更多生产部署、测试和自动化构建说明见[部署与交付文档](DEPLOYMENT.md)。
