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

### 导入笔记

点击目标笔记本右侧的三点菜单，选择“导入笔记/笔记本”：

- 选择 `.md` 或 `.markdown` 文件，将它导入所选的个人笔记本，文件名作为笔记标题。
- 选择 `.rtf` 或 `.docx` 文件，转换为富文本笔记；导入前显示转换注意事项。
- 选择 `.zip` 文件，在该笔记本下创建一个新子笔记本。可修改名称和父级；ZIP 内的目录按原层级创建子笔记本，MD、RTF、DOCX 文件导入对应目录。

确认预览中的数量后点击“开始导入”。进度窗口支持暂停；出错后可在同一窗口点击
“继续导入”，保留已创建内容并避免重复创建。关闭窗口或刷新页面后不再保留本次进度；
此时重新选择文件会发起一次新的导入，请先检查已创建的内容。

Markdown 使用 UTF-8，正文原样保留。RTF/DOCX 在浏览器本地转换为富文本 HTML，支持常见文字格式、段落、列表、表格和安全链接。复杂排版可能简化，内嵌图片/对象以占位文字保留；不执行宏、不抓取外部资源。
不单独导入 ZIP 中的图片或附件，也不改写相对链接。文档及转换后正文每篇最大 4 MiB，ZIP 最大 32 MiB，总正文最大 64 MiB；每个 ZIP
最多 5,000 个条目、2,000 篇笔记和 1,000 个新笔记本（含顶层），目录深度最多 20 层。
加密、多卷和 ZIP64 压缩包不支持；无效文件或不安全路径会在创建内容前被拒绝。

邮件功能需要额外启用 Celery + Redis；笔记功能可独立运行。修改邮箱、邮箱验证和找回密码的配置见[异步邮件部署](MAIL.md)。

不提供自行注册功能。新增账号只能由管理员在 Web「头像 → 管理 → 用户管理」中创建；新账号默认为普通用户。


## 导出笔记与笔记本

在 Web 中打开笔记后点击“导出”：Markdown 笔记下载为 Hugo Markdown 文件（`.md`），富文本笔记下载为 RTF（`.rtf`）。Markdown 文件头包含
`title`（标题）、`date`（创建时间）、`lastmod`（修改时间）、`tags` 和 `categories`
（所属笔记本名称），格式见 [Hugo front matter](https://gohugo.io/content-management/front-matter/)。
时间沿用笔记记录并转换为 UTC；缺失时间标为 `null`，不填入导出时间。Markdown 正文原样保留，
RTF 保留常见文字格式、段落、列表、表格和链接，复杂样式可能简化，图片保留说明及原链接而不内嵌。

笔记本右侧三点菜单 → **导出笔记本（ZIP）**：导出该笔记本及全部子笔记本，目录结构与笔记本层级一致，
包括空目录；Markdown 保存为 `.md`、富文本保存为 `.rtf`，不包括回收站笔记。重名文件自动加序号，读取失败会提示并中止，不生成缺失内容的压缩包。
单次限制为正文合计 128 MiB、10000 篇笔记、3000 个笔记本及 100 层目录。
图片和附件不嵌入压缩包，正文中的链接保留；此功能用于内容导出，不代替数据库与文件备份。
