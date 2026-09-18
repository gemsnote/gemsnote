# 安装和部署指南

Gemsnote 支持 PostgreSQL 和 MongoDB。新安装推荐 PostgreSQL；MongoDB 主要用于直接接入现有 Leanote 数据。数据库中的附件记录不包含文件实体，任何迁移都必须同时处理项目根目录下的 `files/`。

## 部署前准备

1. 备份数据库和 `files/`，迁移期间停止旧服务写入。
2. 修改 `app.secret`，并设置真实的 `site.url`、数据库密码和 HTTPS。
3. 确保运行账号对 `files/` 有读写权限。

Gemsnote 会在 PostgreSQL 的 `gemsnote_schema_migrations` 表或 MongoDB 的同名集合记录数据库版本。旧 Leanote 数据库没有版本记录时会登记为 `1.0.0`，不会因此重写业务数据。

## 一、使用 Release 包部署

Release 包包含服务端、Web 前端、数据库文件、文档和编译好的迁移工具 `bin/gemsnote-migrate`（Windows 为 `.exe`）。包名为：

```text
gemsnote-linux-amd64-v<version>.tar.gz
gemsnote-linux-arm64-v<version>.tar.gz
gemsnote-darwin-amd64-v<version>.tar.gz
gemsnote-darwin-arm64-v<version>.tar.gz
gemsnote-windows-amd64-v<version>.zip
```

### Linux 和 macOS

```bash
tar -xzf gemsnote-<os>-<arch>-v<version>.tar.gz
cd gemsnote
mkdir -p files public/upload
```

macOS 首次运行若被系统拦截，在“系统设置 → 隐私与安全性”中允许该程序。

### Windows

解压 ZIP 后进入 `gemsnote` 目录：

```bat
mkdir files
mkdir public\upload
```

`run.bat` 首次运行使用 `mklink /J` 创建运行时目录联接。若失败，请以管理员身份运行一次或启用 Windows 开发人员模式。

完成下方 PostgreSQL 或 MongoDB 配置后再启动服务。Linux/macOS 执行
`./run.sh`，Windows 执行 `run.bat`。

### Release + PostgreSQL

创建用户和空数据库：

```sql
CREATE USER gemsnote WITH PASSWORD '请替换为强密码';
CREATE DATABASE gemsnote OWNER gemsnote;
```

编辑 `conf/app.conf`：

```ini
db.type=postgresql
db.host=127.0.0.1
db.port=5432
db.dbname=gemsnote
db.username=gemsnote
db.password=请替换为强密码
```

也可使用 `db.url` 配置 PostgreSQL URL 或参数式 DSN。按实际访问地址修改
`site.url`。仓库默认 `app.secret` 技术上可以启动，但正式或生产使用必须在首次
启动前将它改为随机长字符串。

配置完成后直接运行 `run.sh` 或 `run.bat`。Gemsnote 首次连接真正的空数据库时
会自动执行 Release 包中的 `database/schema.sql` 和 `database/seed.sql`；已有表的
数据库不会自动导入初始数据，因此不需要也不应手工执行这两个 SQL 文件。

### Release + MongoDB

全新 MongoDB 可从 Release 内置数据初始化：

```bash
mongorestore --drop --db gemsnote mongodb_backup/gemsnote_install_data
```

配置 `conf/app.conf`：

```ini
db.type=mongodb
db.host=127.0.0.1
db.port=27017
db.dbname=gemsnote
db.username=
db.password=
```

启用认证或副本集时推荐使用完整 URL：

```ini
db.url=mongodb://用户名:密码@主机:27017/gemsnote?authSource=admin
```

### Release 直接连接旧 Leanote MongoDB

无需改库名。停止旧 Leanote 写入并备份后，把 Gemsnote 指向原业务库：

```ini
db.type=mongodb
db.url=mongodb://用户名:密码@数据库主机:27017/leanote?authSource=admin
```

未启用认证时可配置 `db.host`、`db.port` 和 `db.dbname=leanote`。Gemsnote 会继续使用旧数据中的 24 位 ObjectId。首次启动会写入版本记录，因此切换前必须备份。这是由 Gemsnote 接管旧数据库，不是让两个服务长期共享同一业务库；Gemsnote 启动后应保持旧 Leanote 服务停止，避免并发写入和 USN 冲突。

### 从 Leanote 的 MongoDB 原始 db 目录恢复

包含 `WiredTiger`、`collection-*.wt` 等文件的是 MongoDB 物理数据目录，不能直接交给 `mongorestore`。必须操作副本，并先用与旧实例兼容的 MongoDB 临时启动，再导出逻辑备份：

```bash
mkdir -p recovery-db logical-dump
cp -a /path/to/leanote-db-backup/. recovery-db/
docker rm -f leanote-recovery-mongo 2>/dev/null || true
docker run -d --name leanote-recovery-mongo \
  -p 127.0.0.1:27018:27017 \
  -v "$PWD/recovery-db:/data/db" mongo:4.2 --bind_ip_all
docker logs leanote-recovery-mongo
docker exec leanote-recovery-mongo mongodump --db leanote --out /tmp/leanote-dump
docker cp leanote-recovery-mongo:/tmp/leanote-dump/leanote ./logical-dump/
docker stop leanote-recovery-mongo
docker rm leanote-recovery-mongo
```

`mongo:4.2` 是当前 Compose 使用的版本，不保证兼容所有旧备份。若日志提示存储格式或 `featureCompatibilityVersion` 不兼容，应停止容器并改用接近原实例的版本；不要对唯一备份执行 `--repair`。MMAPv1 数据（如 `leanote.ns`、`leanote.0`）通常需旧版 MongoDB 启动后导出。

只导出业务库 `leanote`，不要迁移 `admin`、`config`、`local`。恢复并配置相同库名：

```bash
mongorestore --drop --db leanote logical-dump/leanote
```

### Release：MongoDB 迁移到 PostgreSQL

创建空的目标 PostgreSQL 数据库，但不要执行 `schema.sql` 或 `seed.sql`；迁移工具默认自行应用 `database/schema.sql`。工具迁移固定支持的 Gemsnote/Leanote 业务集合，保留 ObjectId，并默认比较源、目标记录数：

```bash
./bin/gemsnote-migrate \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/leanote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable'
```

Windows：

```bat
bin\gemsnote-migrate.exe -direction mongo_to_pg -mongo-url "mongodb://127.0.0.1:27017/leanote" -postgres-url "host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable"
```

可先用相同参数加 `-dry-run` 验证连接和读取而不写入；由于 dry-run 不会自动应用
schema，目标 PostgreSQL 必须预先具有 `database/schema.sql` 中的表结构。迁移成功后
把 `conf/app.conf` 切换为 PostgreSQL，再启动服务。不要补导 `seed.sql`。

## 二、使用 Docker Compose 部署

仓库当前提供：

- `docker-compose.postgres.yml`：PostgreSQL 18 + Gemsnote，配置为 `conf/app.docker-postgres.conf`；
- `docker-compose.mongodb.yml`：MongoDB 4.2 + Gemsnote，配置为 `conf/app.docker-mongodb.conf`。

Web 默认只发布到宿主机 `127.0.0.1:9000`。
镜像构建不会带入本地 `conf/app.conf`、Compose 专用配置、`files/` 或 `public/upload/` 中的数据；镜像内的 `conf/app.conf` 来自 `conf/app.conf-default`。Compose 在运行时挂载对应配置及文件目录，因此部署前仍需编辑宿主机上的配置并备份持久化数据。

### Docker + PostgreSQL

首次启动前创建应用文件持久化目录：

```bash
mkdir -p files public/upload
```

```bash
docker compose -f docker-compose.postgres.yml up -d --build
docker compose -f docker-compose.postgres.yml ps
docker compose -f docker-compose.postgres.yml logs -f gemsnote
```

数据库挂载的宿主机父目录是 `./data`；当前 PostgreSQL 18 镜像的实际数据库目录是
`./data/18/docker`。`schema.sql` 和 `seed.sql` 只在空数据目录首次启动时执行。
首次启动前修改密码，需同步修改 Compose 的 `POSTGRES_PASSWORD` 与
`conf/app.docker-postgres.conf` 的 `db.password`。数据库已经初始化后，修改
`POSTGRES_PASSWORD` 不会改变数据库内的密码；应先在 PostgreSQL 中执行
`ALTER ROLE gemsnote WITH PASSWORD '新密码';`，再更新这两处配置。

### Docker + MongoDB

当前 MongoDB 配置默认使用 `db.host=mongodb` 和 `db.dbname=leanote`。全新安装时，先将配置中的库名改为 `gemsnote`，再初始化：

首次启动前创建应用文件持久化目录：

```bash
mkdir -p files public/upload
```

启动前确认配置挂载源确实是普通文件，避免 Docker 把缺失路径创建成目录：

```bash
test -f conf/app.docker-mongodb.conf
```

```bash
docker compose -f docker-compose.mongodb.yml up -d mongodb
docker cp mongodb_backup/gemsnote_install_data mongodb:/tmp/gemsnote-install
docker compose -f docker-compose.mongodb.yml exec mongodb \
  mongorestore --drop --db gemsnote /tmp/gemsnote-install
docker compose -f docker-compose.mongodb.yml up -d --build gemsnote
```

MongoDB 数据保存在 `./data/mongo/db`。

### Docker 直接连接旧 Leanote MongoDB

只有当旧库目录由当前 Compose 的 MongoDB 4.2 创建、并且曾正常停机时，才可直接
使用 Compose 的 `./data/mongo/db`。保持当前 `db.host=mongodb`、`db.dbname=leanote` 后启动：

```bash
docker compose -f docker-compose.mongodb.yml up -d --build
```

不要把任意 Leanote 物理 db 备份直接复制到该目录；不同 MongoDB 版本或存储引擎的
备份必须按下一节先转换为逻辑备份。

若 MongoDB 在其他机器，在 `conf/app.docker-mongodb.conf` 设置：

```ini
db.type=mongodb
db.url=mongodb://用户名:密码@可从容器访问的主机:27017/leanote?authSource=admin
```

容器内的 `127.0.0.1` 不是宿主机。Linux 可使用宿主机网关地址，或给 `gemsnote` 服务增加：

```yaml
    extra_hosts:
      - "host.docker.internal:host-gateway"
```

再在 URL 中使用 `host.docker.internal`。

### Docker 从原始 MongoDB db 目录恢复

先按 Release 章节的流程，把物理备份转成仅含 `leanote` 的逻辑备份，再恢复：

```bash
docker compose -f docker-compose.mongodb.yml up -d mongodb
docker cp logical-dump/leanote mongodb:/tmp/leanote
docker compose -f docker-compose.mongodb.yml exec mongodb \
  mongorestore --drop --db leanote /tmp/leanote
docker compose -f docker-compose.mongodb.yml up -d --build gemsnote
```

不要恢复 `admin`、`config`、`local`；认证用户应在新实例重新创建。

### Docker 环境迁移 MongoDB 到 PostgreSQL

迁移工具未打入应用镜像，请使用 Release 包中的二进制，或按源码章节运行。MongoDB Compose 已发布 `127.0.0.1:27017`；PostgreSQL Compose 默认未发布端口，迁移期间需临时取消其中以下配置的注释：

```yaml
    ports:
      - 127.0.0.1:5432:5432
```

目标 PostgreSQL 必须是空库，不能让 Compose 先导入 `seed.sql`。建议复制
`docker-compose.postgres.yml` 为 `docker-compose.postgres-migration.yml`，把 PostgreSQL
数据挂载改为独立目录：

```yaml
      - ./data-pg-migration:/var/lib/postgresql
```

同时在迁移专用文件中移除两个 `/docker-entrypoint-initdb.d/` 挂载，让 migration tool
自行应用 schema。创建全新的 `./data-pg-migration` 后再启动目标 PostgreSQL。绝对不要
清空共享父目录 `./data`，因为当前 MongoDB 源数据位于 `./data/mongo/db`；也不要删除任何
包含有效数据的目录。

```bash
docker compose -f docker-compose.mongodb.yml up -d mongodb
docker compose -f docker-compose.postgres-migration.yml up -d postgres
./bin/gemsnote-migrate -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/leanote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
docker compose -f docker-compose.postgres-migration.yml up -d --build gemsnote
```

迁移后可重新注释 PostgreSQL 端口发布。

### Docker 持久化和迁移 files

两个 Compose 文件已默认挂载附件和旧上传目录：

```yaml
      - ./files:/opt/gemsnote/files
      - ./public/upload:/opt/gemsnote/public/upload
```

停止旧服务写入后复制整个目录，保留相对路径：

```bash
mkdir -p files public/upload
cp -a /path/to/old-leanote/files/. ./files/
cp -a /path/to/old-leanote/public/upload/. ./public/upload/
```

主要图片和附件位于 `files/`；头像或旧主题资源可能位于 `public/upload/`。遗漏实体文件会导致数据库记录存在但文件无法访问。

## 三、从源码部署

### Linux、macOS 和 Windows

需要 Go 1.22+（当前 CI/Docker 使用 Go 1.24）和 Node.js 22，以及选定数据库的客户端。Linux/macOS 使用以下命令：

```bash
git clone <repository-url> gemsnote
cd gemsnote
npm ci --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend
go test ./...
go run github.com/revel/cmd/revel run -a . -m prod
```

Windows PowerShell：

```powershell
git clone <repository-url> gemsnote
Set-Location gemsnote
npm ci --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend
go test ./...
go run github.com/revel/cmd/revel run -a . -m prod
```

开发模式将最后的 `prod` 改为 `dev`。生产环境更适合构建 Release：

```bash
scripts/build-release.sh <version> <goos> <goarch> <output-dir>
scripts/build-release.sh 1.0.0 linux amd64 ./dist
```

Windows 原生环境若无 Bash，可在 Git Bash/WSL 中执行脚本，或使用项目 GitHub Release。

### 源码 + PostgreSQL

创建空数据库，把 `conf/app.conf` 的 `db.type` 设为 `postgresql` 并填写连接信息，
然后直接启动。Gemsnote 会在首次连接真正的空数据库时自动应用
`database/schema.sql` 和 `database/seed.sql`，三个平台都无需手工导入。默认
`app.secret` 可以启动，但正式或生产使用必须在首次启动前更换。

### 源码 + MongoDB 或直连 Leanote

全新 MongoDB：

```bash
mongorestore --drop --db gemsnote mongodb_backup/gemsnote_install_data
```

配置 `db.type=mongodb`、MongoDB 主机和 `db.dbname=gemsnote`。直连旧 Leanote 时改为 `db.dbname=leanote`，或设置完整 `db.url`。原始 db 目录仍须先转换为逻辑备份，不能直接使用 `mongorestore`。

### 源码迁移 Leanote MongoDB 到 PostgreSQL

创建空目标库，不执行 schema/seed，然后运行：

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/leanote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable'
```

Windows PowerShell 使用同一入口；参数应使用 PowerShell 的续行符：

```powershell
go run ./tools/migration `
  -direction mongo_to_pg `
  -mongo-url "mongodb://127.0.0.1:27017/leanote" `
  -postgres-url "host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable"
```

可选参数：`-schema <path>`（默认 `database/schema.sql`）、`-dry-run`、`-validate`（默认 `true`）。迁移成功后切换 `conf/app.conf` 到 PostgreSQL，并迁移 `files/` 与需要的 `public/upload/`。

## 验证、备份和安全

部署或迁移后至少验证：账号登录、笔记本和笔记、Markdown、图片、头像、附件、共享笔记与同步，以及：

```bash
curl http://127.0.0.1:9000/api/system/version
```

升级前停止写入，同时备份数据库、`files/`、`public/upload/` 和实际配置。PostgreSQL 使用 `pg_dump`，MongoDB 只需 `mongodump --db <业务库>`，不要迁移 `admin/config/local`。数据库端口不要暴露到公网，Web 服务通过反向代理提供 HTTPS，启用 HTTPS 后将 `cookie.secure=true`。

更专门的数据转换说明参阅 [迁移指南](MIGRATION_GUIDE.md)。
