# Gemsnote（珠玑笔记）

Gemsnote（珠玑笔记）是一套支持私有部署的开源笔记与知识管理系统，提供笔记本、标签、Markdown、富文本编辑、分享与协作等功能。Web 端采用 Vue 前后端分离架构；本项目不再提供 Leanote 的公开博客、评论和博客主题功能。

当前版本：`1.0.0`。

## 项目来源

本项目基于开源项目 [Leanote](https://github.com/leanote/leanote) 修改开发。感谢 Leanote 原作者及所有贡献者奠定的基础。

Gemsnote 在保留 Leanote 核心功能和客户端兼容性的基础上，主要进行了以下扩展：

- 增加 PostgreSQL 数据库支持，同时保留 MongoDB 支持；
- 提供 MongoDB 与 PostgreSQL 双向数据迁移工具；
- 两种数据库统一使用 MongoDB ObjectId 兼容的 24 位十六进制 ID；
- 抽象数据库访问层，可通过配置选择 MongoDB 或 PostgreSQL；
- 项目名称、默认数据库、配置、界面和相关资源更名为 Gemsnote／珠玑笔记；
- 保持原有 HTTP API 路径、请求参数、响应字段及 MongoDB BSON 字段稳定，以兼容旧版 Leanote 客户端和已有数据。

本项目继续遵循仓库中的开源许可证。使用、修改和分发时，请同时遵守原项目及本项目的许可证要求。

## 主要功能

- 笔记本、笔记和标签管理；
- Markdown 编辑器与富文本编辑器；
- Vim 和 Emacs 编辑模式；
- 笔记分享与多人协作；
- Vue Web、桌面客户端兼容访问；
- PDF 导出；
- 批量笔记操作；
- MongoDB 和 PostgreSQL 双数据库后端；
- MongoDB ↔ PostgreSQL 双向迁移。

## 数据库支持

数据库后端通过 `conf/app.conf` 中的 `db.type` 选择。

### PostgreSQL

PostgreSQL 是新部署的默认推荐后端：

```ini
db.type=postgresql
db.url=host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable
```

### MongoDB

继续使用 MongoDB：

```ini
db.type=mongodb
db.url=mongodb://127.0.0.1:27017/gemsnote
```

Gemsnote 的 MongoDB 集合名和 BSON 字段与原版 Leanote 保持兼容。升级已有 Leanote 实例时，可以直接连接原来的 `leanote` 数据库，不需要先修改数据库名：

```ini
db.type=mongodb
db.url=mongodb://127.0.0.1:27017/leanote
```

连接已有数据库前请先完成备份。Gemsnote 不会自动发现 `leanote` 数据库，必须在配置中明确指定数据库名。

## 版本与自动升级

Gemsnote 从 `1.0.0` 起使用独立版本号。应用版本集中定义在 `app/version/version.go`，管理后台和数据库迁移流程均读取该版本。

服务启动连接数据库后，会自动检查并执行尚未应用的数据库迁移：

- MongoDB 在 `gemsnote_schema_migrations` 集合记录已应用版本；
- PostgreSQL 在 `gemsnote_schema_migrations` 表记录已应用版本；
- 原 Leanote MongoDB 或早期 Gemsnote 数据库没有版本记录时，会作为升级基线登记为 `1.0.0`，不会改写现有业务数据；
- 如果数据库版本高于当前应用版本，服务会拒绝启动，防止旧程序破坏新结构。

升级前仍应备份数据库和附件文件。新增数据库结构变更时，需要在 `app/db/migrations.go` 中增加幂等迁移，不能只修改 `database/schema.sql`。

## 从 Leanote 迁移

从旧 Leanote MongoDB 迁移至 Gemsnote PostgreSQL：

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/leanote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

迁移工具会保留原 MongoDB ObjectId，并默认执行 PostgreSQL Schema 和逐表记录数校验。迁移不会复制磁盘上的图片、附件实体文件，这些文件目录需要单独迁移或继续使用原挂载路径。

迁移不是服务启动时的自动操作。正式迁移前应停止旧服务写入，同时备份源数据库、目标数据库及附件文件。

## 快速启动

### Docker + PostgreSQL

```bash
docker-compose -f docker-compose.postgres.yml up -d
```

空 PostgreSQL 数据卷会依次执行 `database/schema.sql` 和 `database/seed.sql`，不需要 MongoDB 即可完成首次初始化。初始管理员为 `admin`／`gemsnote`，这是兼容内置 seed 数据的初始密码；首次登录后必须立即修改密码。初始化 SQL 只在 PostgreSQL 数据目录为空时执行。

### Docker + MongoDB

```bash
docker-compose -f docker-compose.mongodb.yml up -d
```

### 从源码运行

安装 Go、Node.js 22 和对应数据库后，先构建 Vue 前端，再启动 Revel：

```bash
npm ci --prefix frontend
npm run build --prefix frontend
revel run github.com/gemsnote/gemsnote
```

默认访问地址为：<http://localhost:9000>

首次部署前必须修改 `conf/app.conf` 中的 `app.secret`，并根据实际环境设置数据库账号、密码、监听地址和站点 URL。

## 文档

- [快速开始](docs/QUICK_START.md)
- [部署指南](docs/DEPLOYMENT.md)
- [MongoDB ↔ PostgreSQL 数据迁移指南](docs/MIGRATION_GUIDE.md)
- [数据库抽象概览](docs/DATABASE_ABSTRACTION_README.md)
- [数据库抽象完整指南](docs/DATABASE_ABSTRACTION_GUIDE.md)
- [自动化构建与发布](docs/RELEASE.md)
- [统一 UI 设计](docs/UI.md)
- [Vue Web UI 重写与部署计划](docs/WEB_UI_MIGRATION_PLAN.md)
- [API2 新接口](docs/API2.md)
- [Legacy API 兼容说明](docs/API.md)

## API 与客户端兼容性

Gemsnote 保留了 Leanote 原有 HTTP API 路径、参数和响应字段。旧版 Leanote Electron 客户端可以在登录界面选择自建服务并填写 Gemsnote 服务端地址进行连接。

旧客户端默认地址仍可能指向 `leanote.com`，因此连接私有部署时必须明确填写服务器地址。

## 测试

运行全部 Go 测试：

```bash
go test ./...
```

数据库集成测试和双向迁移测试所需的环境变量及命令，请参阅[数据库抽象完整指南](docs/DATABASE_ABSTRACTION_GUIDE.md)。

## 自动化构建与发布

每次 push 和 Pull Request 都会自动安装、测试并构建 Vue 前端，执行全部 Go 测试，生成并编译真正的 Revel 服务端入口，同时验证 Docker 镜像可以构建。

推送与 `app/version/version.go` 一致的正式 SemVer 标签（例如 `v1.0.0`）时，GitHub Actions 会自动创建 GitHub Release，并附带 Linux amd64／arm64、Windows amd64、macOS amd64／arm64 发布包及 SHA-256 校验文件，同时向 GHCR 发布 Linux amd64／arm64 多架构镜像。详细操作见[自动化构建与发布](docs/RELEASE.md)。

## 相关项目

- [Leanote 原项目](https://github.com/leanote/leanote)
- [Leanote Desktop 原项目](https://github.com/leanote/desktop-app)

## 参与贡献

欢迎通过 Issue 和 Pull Request 报告问题、提出建议或贡献代码。提交涉及数据库结构、迁移逻辑或 API 的修改时，请同时补充相应测试，并确保 MongoDB、PostgreSQL 和旧客户端兼容性不受影响。
