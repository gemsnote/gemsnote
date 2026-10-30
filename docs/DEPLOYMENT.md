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

Release 默认 `http.addr=0.0.0.0`、`http.port=9000`，会监听所有 IPv4 网卡。
`http://127.0.0.1:9000` 只是本机访问地址，不表示仅本机可连接。
仅本机使用或同机反向代理部署时将 `http.addr` 设为 `127.0.0.1`；远程访问需配置
防火墙和 HTTPS。此处与 Compose 将宿主机端口绑定到 `127.0.0.1` 的默认行为不同。

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
会自动执行 Release 包中的 `database/schema.sql` 和 `database/seed.sql`。如果已有
schema，但 `public` 下所有业务基础表仍为空，也会补做初始化（忽略版本记录表）。
只要其中任意业务表已有数据，就不会自动导入 seed。普通新安装无需手工导入 SQL；
迁移库不要导入 seed，也不要在数据迁移完成前启动应用，以免生成默认账号和笔记。

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

下面是 Bash 示例。先检查是否存在同名容器；若存在，核实其用途并为本次恢复选择
另一个名称（同时替换后续命令），不要强制删除已有容器。`recovery-db` 必须是本次
恢复专用的空目录，不要与上次恢复数据混用。

```bash
mkdir -p recovery-db logical-dump
cp -a /path/to/leanote-db-backup/. recovery-db/
docker container ls -a --filter name=leanote-recovery-mongo
```

确认名称未占用、备份副本准备完成后启动临时容器：

```bash
docker run -d --name leanote-recovery-mongo \
  -p 127.0.0.1:27018:27017 \
  -v "$PWD/recovery-db:/data/db" mongo:4.2 --bind_ip_all
# 最多等待约 60 秒；超时会终止本段命令，不执行导出或清理。
(
  ready=false
  for attempt in {1..30}; do
    if docker exec leanote-recovery-mongo mongo --quiet --eval 'quit(db.adminCommand({ping:1}).ok ? 0 : 1)' >/dev/null 2>&1; then
      ready=true
      break
    fi
    sleep 2
  done
  if [ "$ready" != true ]; then
    docker logs leanote-recovery-mongo
    echo 'MongoDB 尚未就绪，请检查版本兼容性和日志，暂勿继续恢复。' >&2
    exit 1
  fi
)
```

仅当就绪检查成功后，执行导出；任何一步失败都应停止并检查，不要继续删除容器：

```bash
(
set -e
docker exec leanote-recovery-mongo mongodump --db leanote --out /tmp/leanote-dump
docker cp leanote-recovery-mongo:/tmp/leanote-dump/leanote ./logical-dump/
docker stop leanote-recovery-mongo
docker rm leanote-recovery-mongo
)
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

可先用相同参数加 `-dry-run -validate=false` 验证连接和读取而不写入；由于 dry-run 不会自动应用
schema，目标 PostgreSQL 必须预先仅导入 `database/schema.sql`，不要导入 seed。
dry-run 不模拟转换和写入约束；默认记录数校验也不适合尚未复制数据的目标库。
正式迁移时移除这两个参数，恢复默认记录数校验。迁移成功后
把 `conf/app.conf` 切换为 PostgreSQL，再启动服务。不要补导 `seed.sql`。

## 二、使用 Docker Compose 部署

仓库当前提供：

- `docker-compose.yml`：直接使用 GHCR 发布镜像 + PostgreSQL 18，无本地构建；
- `docker-compose.postgres.yml`：PostgreSQL 18 + Gemsnote，配置为 `conf/app.docker-postgres.conf`；
- `docker-compose.mongodb.yml`：MongoDB 4.2 + Gemsnote，配置为 `conf/app.docker-mongodb.conf`。

Web 默认只发布到宿主机 `127.0.0.1:9000`。
镜像构建不会带入本地 `conf/app.conf`、Compose 专用配置、`files/` 或 `public/upload/` 中的数据；镜像内的 `conf/app.conf` 来自 `conf/app.conf-default`。Compose 在运行时挂载对应配置及文件目录，因此部署前仍需编辑宿主机上的配置并备份持久化数据。

### 使用 GHCR 发布镜像

正式版本镜像发布在 `ghcr.io/gemsnote/gemsnote:<版本>`。推荐使用默认 `docker-compose.yml`。
从对应 Release 下载 Source code 压缩包并解压，进入其根目录；不需要编译源码，
但必须提供宿主机上的 Compose、配置及 SQL 文件，不能只下载镜像。

首次启动前创建 `.env`（不要提交真实密码）：

```dotenv
GEMSNOTE_VERSION=1.0.0
GEMSNOTE_POSTGRES_PASSWORD=请替换为强密码
```

编辑 `conf/app.docker-postgres.conf`，设置同样的 `db.password`、随机 `app.secret`
和实际 `site.url`，保留 `db.host=postgres`。`.env` 不会自动修改应用配置。
所需文件为 `docker-compose.yml`、`conf/app.docker-postgres.conf`、
`database/schema.sql`、`database/seed.sql`：

```bash
mkdir -p files public/upload
docker compose config --quiet
docker compose pull
docker compose up -d
docker compose ps
docker compose logs -f gemsnote
```

数据库通过健康检查后才启动应用。持久化路径与下面的源码 Compose 相同；不要同时
运行两套 Compose 访问同一个数据库数据目录。默认文件没有 `build` 配置。
升级时备份数据，更新 `.env` 中的版本并检查新版本配置，再执行 `pull` 和 `up -d`。
已有库改密码需要先执行 `ALTER ROLE`，再同步修改 `.env` 和应用配置，不能仅改环境变量。

如果需要用现有 `docker-compose.postgres.yml` 或 `docker-compose.mongodb.yml`
运行 GHCR 镜像，可先拉取并标记为这两个文件中的应用镜像名：

```bash
docker pull ghcr.io/gemsnote/gemsnote:1.0.0
docker tag ghcr.io/gemsnote/gemsnote:1.0.0 gemsnote:1.0.0
```

然后将下面示例中的 `--build` 替换为 `--no-build`，使用已拉取的镜像；这两个旧文件保留源码构建用途。GHCR 包若为私有，需要先执行 `docker login ghcr.io`；公开包可直接拉取。镜像支持 Linux amd64 和 arm64，Docker 会根据主机架构选择对应镜像。

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

当前 MongoDB 配置默认使用 `db.host=mongodb` 和 `db.dbname=gemsnote`，全新安装按此库名初始化：

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
使用 Compose 的 `./data/mongo/db`。保留 `db.host=mongodb`，并明确将默认库名改为原业务库名（例如 `db.dbname=leanote`）后启动：

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

先按 Release 章节的流程，把物理备份转成仅含 `leanote` 的逻辑备份，再恢复。
启动应用前将 `conf/app.docker-mongodb.conf` 的 `db.dbname` 从 `gemsnote` 改为
`leanote`（或实际恢复库名），否则会连到另一个库：

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
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=请替换为已配置的数据库密码 dbname=gemsnote sslmode=disable'
docker compose -f docker-compose.postgres-migration.yml up -d --build gemsnote
```

迁移命令中的密码必须与目标 PostgreSQL 实际密码一致，并与迁移 Compose 中的
`POSTGRES_PASSWORD`、`conf/app.docker-postgres.conf` 的 `db.password` 对应。
等待目标 PostgreSQL 可连接后再执行迁移；迁移失败不要启动应用。迁移后可重新注释 PostgreSQL 端口发布。

迁移完成后继续使用该迁移 Compose 文件及其 `data-pg-migration` 挂载。
不要直接切回默认 `docker-compose.yml`：它挂载的是 `./data`，并不会自动使用迁移库。
如需改用默认文件，必须先停止原部署，核对并调整数据库挂载、凭据及应用配置，
确保仍指向同一份已迁移的数据，且没有两个 PostgreSQL 进程同时访问同一目录。

若已将数据迁移到宿主机或其他机器上的 PostgreSQL，可复制默认 Compose 为外部库
专用文件，移除其中的 `postgres` 服务以及 `gemsnote.depends_on`，保留 GHCR 应用镜像和
文件挂载。在挂载的 `conf/app.docker-postgres.conf` 中设置 `db.type=postgresql` 和指向
已迁移库的 `db.url`（主机、端口、库名、账号及 TLS 参数按实际环境填写），使用
`docker compose -f <外部库专用文件> up -d` 启动。容器内的 `127.0.0.1` 不是宿主机；
连接宿主机可按前文设置 `host.docker.internal:host-gateway`，同时确保 PostgreSQL
监听地址、防火墙及 `pg_hba.conf` 允许容器访问。仅监听宿主机回环地址的数据库不能
直接通过网关访问。外部库部署不运行 Compose 的 schema/seed 初始化，也不要再导入 seed。

迁移中途失败或最终校验失败都不会自动撤销已写入的数据；不要直接切换服务。
恢复备份或核查后重试的边界见[迁移失败处理](development/MIGRATION_GUIDE.md#非空目标库与失败处理)。

### Docker 持久化和迁移 files

三个 Compose 文件均已默认挂载附件和旧上传目录：

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

需要 Go 1.22+（当前 CI/Docker 使用 Go 1.24）和 Node.js 22，以及选定数据库的客户端。
推荐从源码构建完整 Release 后运行，避免遗漏 Revel 运行时资源。Linux/macOS 使用 Bash；
Windows 使用 Git Bash（需 Go、Node.js、zip 均在 PATH），构建 Windows 原生包：

```bash
git clone <repository-url> gemsnote
cd gemsnote
npm ci --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend
go test ./...
scripts/build-release.sh linux amd64 "$PWD/release"
```

按目标平台将 `linux amd64` 改为 `linux arm64`、`darwin amd64`、`darwin arm64`
或 `windows amd64`。解压产物后按 Release 章节配置数据库，再运行 `run.sh` / `run.bat`。
在 WSL 构建 Windows 包时同样显式指定 `windows amd64`，不要使用 WSL 的 Linux 产物。

脚本使用仓库内 `go run ./app/cmd build . <临时目录> prod` 生成入口，再编译 `./app/tmp`。
本地命令只实现 `build`，不要改为 `./app/cmd run`；也不要使用外部 Revel CLI 替代，
否则会绕过仓库对 `data/`、`docs/`、`tmp/`、`files/` 的扫描排除规则。
输出目录使用绝对路径，完整依赖与手工构建步骤见[发布文档](development/RELEASE.md)。

### 源码 + PostgreSQL

创建空数据库，把 `conf/app.conf` 的 `db.type` 设为 `postgresql` 并填写连接信息，
然后直接启动。Gemsnote 会在首次连接真正的空数据库时自动应用
`database/schema.sql` 和 `database/seed.sql`；已有 schema 但所有业务表为空时也会补初始化，
三个平台都无需手工导入。已有业务数据不会补导 seed。默认
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

可选参数：`-schema <path>`（默认 `database/schema.sql`）、`-dry-run`、`-validate`（默认 `true`）。预检需先准备 schema，并使用 `-dry-run -validate=false`；正式迁移恢复默认校验。迁移成功后切换 `conf/app.conf` 到 PostgreSQL，并迁移 `files/` 与需要的 `public/upload/`。

## 验证、备份和安全

部署或迁移后至少验证：账号登录、笔记本和笔记、Markdown、图片、头像、附件、共享笔记与同步，以及：

```bash
curl http://127.0.0.1:9000/api2/system/version
```

升级前停止写入，同时备份数据库、`files/`、`public/upload/` 和实际配置。PostgreSQL 使用 `pg_dump`，MongoDB 只需 `mongodump --db <业务库>`，不要迁移 `admin/config/local`。数据库端口不要暴露到公网，Web 服务通过反向代理提供 HTTPS，启用 HTTPS 后将 `cookie.secure=true`。

更专门的数据转换说明参阅 [数据库迁移指南](development/MIGRATION_GUIDE.md)。

### 附件或图片文件缺失

附件元数据存在不代表文件已恢复。若下载返回 `404 fileNotFound`，查看服务端日志中的
`Download file missing` 及其 `path`，检查该路径对应的文件是否位于应用根目录下，
以及 `files/`、`public/upload/` 的挂载和读取权限。路径是目录而非普通文件时也返回 404；
打开或读取文件属性失败时返回 `500 fileReadFailed`，日志会保留具体原因。

数据库迁移工具不会复制附件和图片。恢复原实例的文件目录时需保留相对目录结构，
不能只恢复数据库，也不要删除附件元数据来掩盖文件丢失。补齐文件后可重新下载，
无需重新导入数据库。

## 四、测试与自动化构建

源码修改后，在仓库根目录运行基础检查：

```bash
npm ci --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend
go test ./...
```

数据库 CRUD 和迁移往返测试需要按[数据库抽象指南](development/DATABASE_ABSTRACTION_README.md)设置对应的集成测试环境变量。服务端的 GitHub Actions 会在每次 push 和 Pull Request 中执行前端测试、Go 测试、构建校验和 Docker 构建。

推送与应用版本一致的 `MAJOR.MINOR.PATCH` Git 标签（如 `1.0.0`，不带 `v`）后，Actions 会构建 Release 包并发布同名 GHCR 多架构镜像，也可以手动指定已有 tag 发布。标签须指向包含新版工作流的提交；构建产物文件名仍使用 `-v<version>`。本地构建、tag 要求及发布产物详见[自动化构建与发布](development/RELEASE.md)。

## 后台任务服务

`gemsnote-worker` 是使用 Asynq 的独立 Go 后台进程，目前负责邮箱验证、更换邮箱和找回密码等邮件任务。服务端负责校验、渲染模板并入队，worker 通过 SMTP 投递。Redis 仅保存任务队列和发送状态，笔记及用户数据仍在业务数据库中。不启用 worker 不影响笔记功能，但新用户无法收到验证邮件，邮箱修改和找回密码也无法发送邮件。

### Docker 启动 worker

三个 Compose 文件都提供可选的 `mail` profile。使用 GHCR 镜像：

```bash
docker compose --profile mail up -d
```

这会增加 `redis` 和 `gemsnote-worker` 两个服务。`gemsnote` 和 `gemsnote-worker` 使用同一个 `ghcr.io/gemsnote/gemsnote:${GEMSNOTE_VERSION}` 镜像，分别运行在独立容器中。默认入口 `/opt/gemsnote/docker-entry-server.sh` 启动 Go 服务端；worker 通过 Compose 覆盖入口为 `/opt/gemsnote/docker-entry-gemsnote-worker.sh`，以非 root 用户启动 Asynq。两个容器都通过 `redis://redis:6379/0` 连接邮件队列。新版本发布前，本地测试使用源码 Compose：

```bash
docker compose -f docker-compose.postgres.yml --profile mail up -d --build
```

源码 MongoDB Compose 使用同样的 `--profile mail` 参数。Redis 不发布宿主机端口，任务队列使用独立的 `mail-redis-data` 卷持久化。排查发送失败：

```bash
docker compose --profile mail logs --tail=100 gemsnote-worker redis
```

源码部署时给这些命令加对应的 `-f` 参数。仅重建 Go 服务端镜像不会启动 Asynq，必须启用 `mail` profile。

`mail` profile 在三个 Compose 文件中默认关闭。服务端和 worker 是两个独立容器，需同时使用同一版本镜像；worker 不挂载业务配置、附件或图片目录。其入口脚本以 UID 10001 运行，默认并发数为 2，Compose 停止宽限期为 45 秒。需要调整并发时，可为 `gemsnote-worker` 服务设置 `command: ["-concurrency", "4"]`。

### Release 和源码启动 worker

Release 包附带原生可执行文件 `bin/gemsnote-worker`（Windows 为 `.exe`），无需 Python 环境。准备 Redis，在 Release 解压目录运行：

```bash
export GEMSNOTE_MAIL_REDIS_URL='redis://127.0.0.1:6379/0'
./bin/gemsnote-worker -concurrency 2
```

Windows PowerShell：

```powershell
$env:GEMSNOTE_MAIL_REDIS_URL = 'redis://127.0.0.1:6379/0'
.\bin\gemsnote-worker.exe -concurrency 2
```

源码部署可在项目根目录运行 `go build -o /tmp/gemsnote-worker ./gemsnote-worker`，设置同一环境变量后执行该文件。`-concurrency` 默认为 2；worker 只需 Redis 和 SMTP，不连接业务数据库，也不读取 `app.conf`。

将 worker 配置为系统服务长期运行。Go 服务端也必须配置同一个 `GEMSNOTE_MAIL_REDIS_URL`，或在 `conf/app.conf` 的全局配置段设置：

```ini
mail.redis_url=redis://127.0.0.1:6379/0
```

环境变量优先于配置文件。Redis 需专用于邮件队列，或使用独立数据库编号；可使用 `redis://:密码@主机:端口/数据库编号`。不要将 Redis 暴露到公网；其中短暂保存邮件正文、验证链接和 SMTP 凭据。敏感任务内容最多保存 10 分钟，发送完成后删除；发送状态保留 24 小时。Windows 可运行原生 Go worker，Redis 可部署在 Docker 或 Linux 主机。

服务端读取 `mail.redis_url` 后需要重启；worker 始终通过自己的 `GEMSNOTE_MAIL_REDIS_URL` 环境变量连接 Redis。两者必须指向同一主机、端口和数据库编号。Linux/macOS 可由系统服务管理进程，Windows 可由服务管理工具启动 worker；`run.sh` / `run.bat` 只启动 Web 服务端。

### SMTP 与站点地址

管理员在 Web 笔记页面左下角点击头像 → **管理** → **站点与邮件**，填写站点 URL、SMTP 主机、端口、用户名、密码和 SSL。

- 站点 URL 应是用户实际可以访问的 HTTPS 地址，用于生成验证链接。
- SMTP 用户名同时作为发件地址，必须是有效邮箱。
- SSL 开启表示隐式 TLS（通常为 465）；关闭表示使用 STARTTLS（通常为 587）。不允许明文传输 SMTP 密码。
- 自建 CA 可为 worker 配置 `SSL_CERT_FILE`；不要关闭证书验证。
- 账号页发送验证邮件可以检查完整发送链路。入队成功不等于送达收件箱；`sent` 表示 SMTP 服务器已接受邮件，仍可能进入垃圾邮件或随后退信。

验证邮件语言跟随触发操作时的界面语言。新账号及当前邮箱验证链接有效 48 小时，修改邮箱链接有效 2 小时；新用户验证邮箱后设置密码的凭据有效 2 小时。站点地址必须指向生成验证记录的同一实例，错误地址会让链接在另一个实例被报告为无效。用户操作见[快速开始](QUICK_START.md#账号与邮箱)。

### 发送状态与故障排查

页面显示 `queued`（排队）、`sending`（发送中）、`sent`（SMTP 已接受）和 `failed`（失败）。账号页可重新查询发送结果；管理员也能查看自己管理的账号的邮件任务。查看 worker 日志时可用任务 ID 对应页面结果。

| 错误或现象 | 检查内容 |
| --- | --- |
| `mailQueueNotConfigured` | 服务端的 Redis 环境变量或 `mail.redis_url` 是否已配置。 |
| `mailQueueUnavailable` | Redis 是否运行，两端连接地址、密码、数据库编号及容器网络是否一致。 |
| `mailNotConfigured` | 管理页面的 SMTP 主机、端口、用户名和密码是否完整。 |
| `mailConnectionFailed` | worker 到 SMTP 的网络、端口与 SSL 模式、证书和自建 CA 配置。 |
| `mailAuthenticationFailed` | SMTP 用户名、密码或邮件服务商授权凭据。 |
| `mailRecipientRejected` | 收件地址及服务商的收件人限制。 |
| `mailDeliveryFailed` | SMTP 投递失败或中断后的结果不确定；检查日志后由用户手动重发。 |
| 长时间排队或 `mailTaskExpired` | 是否启用了 `mail` profile、worker 是否运行并连接到同一 Redis、任务是否超过 10 分钟。 |

敏感任务内容最多保留 10 分钟，处理完成后删除；发送状态保留 24 小时。整个 SMTP 会话限时 30 秒。SMTP 失败不自动重发，Redis 基础设施错误最多重试 3 次；若进程中断后无法确认邮件是否已投递，将标记失败供用户手动重发，避免自动重复投递。每个账号发送验证邮件有 60 秒间隔限制。

管理员创建账号成功但邮件发送失败时，不要重复创建用户。修复配置后，在用户管理点击“邮箱未验证”确认重发，或让用户使用登录页的“找回密码”。邮件已被 SMTP 接受仍可能进入垃圾邮件或发生后续退信。

### 从旧任务服务升级

Asynq 与 Celery 的任务协议不兼容。先暂停创建账号、发送验证邮件等入队操作，保持旧 worker 运行直至旧任务处理完成，再同时升级服务端和 worker；不要混用新旧版本。无法处理的旧任务不会自动转入新队列，升级后通过页面重新发送。

切换时停止旧 Compose 项目，避免旧 `mail-redis` 与新 `redis` 同时使用同一个数据卷。以下以源码 PostgreSQL 部署为例：

```bash
docker compose -f docker-compose.postgres.yml --profile mail down --remove-orphans
docker compose -f docker-compose.postgres.yml --profile mail up -d --build
```

GHCR 部署使用 `docker compose --profile mail down --remove-orphans`，随后 `docker compose --profile mail pull` 和 `docker compose --profile mail up -d`。保持原来的项目名和部署目录，不要使用 `down -v`；服务名改为 `redis`，已有 `mail-redis-data` 卷名保留，避免丢失队列。此操作会短暂停止服务。

此前的 `mail-worker` 服务现已改名为 `gemsnote-worker`，可执行文件统一为 `gemsnote-worker`，入口脚本为 `docker-entry-gemsnote-worker.sh`。已使用 Asynq 的实例只需更新部署配置、重新构建或拉取镜像，再执行 `docker compose --profile mail up -d --remove-orphans`（源码部署加对应 `-f` 和 `--build`）；这会移除旧服务容器。邮件队列、Redis 数据卷和 `GEMSNOTE_MAIL_REDIS_URL` 配置不变，未处理的 Asynq 任务可继续消费。Release 部署需同步更新系统服务中的程序路径。

## 管理员与用户状态

Web 的“头像 → 管理 → 用户管理”列出每位用户的角色，可将其他用户设为“管理员”或“普通用户”。新建用户默认为普通用户。用户管理还可设置账号状态：`A`（可用）或 `D`（禁用）。管理员不能禁用自己；禁用后登录、已有会话及验证链接失效，重新启用需重新登录，笔记数据保留。

管理员不能将自己改为普通用户：界面禁用自己的角色下拉框，服务端也强制拒绝此请求。其他管理员可以撤销该用户的权限，下一次管理请求立即按最新角色检查。

升级到 1.2.0 时自动增加用户角色字段。首次初始化沿用已有 `adminUserId`，或按 `adminUsername` 定位原管理员。若旧管理员已改名且未保存 ID，只有在所有旧全局配置均指向同一个有效拥有者、该账号没有明确角色且当前没有管理员时，才按该拥有者 ID 恢复管理员；配置拥有者缺失或不一致时不会猜测。只有没有角色的旧账号才会被初始化为管理员。此后以用户记录中的角色为准，改用户名或邮箱不改变角色，重启也不会重新提升已被降为普通用户的账号。

若旧配置无法唯一识别原管理员（拥有者缺失、不一致或账号不存在），启动日志会报告管理员初始化失败。此时需在 `conf/app.conf` 中设置正确的 `adminUsername` 或 `adminUserId=<24 位用户 ID>`，用于初始化原管理员。配置不再是长期权限白名单，也不能覆盖已有明确角色。

用户状态字段随 1.3.0 数据库迁移自动增加，历史用户默认为 `A`（可用），角色字段与状态字段独立保存。用户邀请和邮箱操作见[快速开始](QUICK_START.md#管理员创建用户)，接口合同见 [API2](development/API2.md)。
