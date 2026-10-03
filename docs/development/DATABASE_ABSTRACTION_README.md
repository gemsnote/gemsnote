# Gemsnote 数据库抽象

数据库后端由 `db.type` 选择，可配置为 `mongodb` 或 `postgresql`。业务 service 不直接使用 `mgo.Collection`，统一通过 `app/db` 的 CRUD、条件查询、排序、分页和投影接口访问数据库。

## ID 约定

两个后端都使用 24 位 MongoDB ObjectId。应用模型继续使用 `bson.ObjectId`，PostgreSQL 适配层负责在 `bson.ObjectId` 与 `CHAR(24)` 之间转换。

这样可以保持现有 API、客户端同步协议和 MongoDB 数据兼容，也使双向迁移无需维护 UUID 映射。

## 查询兼容层

PostgreSQL 适配器支持当前业务使用的 BSON 条件和更新操作，包括：

- `$or`、`$and`；
- `$gt`、`$gte`、`$lt`、`$lte`、`$ne`；
- `$in`、`$nin`、`$all`、`$exists`、`$regex`；
- `$set`、`$inc`、`$push`、`$pull`、`$addToSet`；
- 排序、offset、limit 和字段投影。

未知操作符会生成恒假条件，不会退化为全表更新或删除。

## PostgreSQL schema

DDL 位于 `database/schema.sql`，安装初始数据位于 `database/seed.sql`，表和列与 `app/info` 中的持久化模型保持一致。`docker-compose.postgres.yml` 会在新的 PostgreSQL 数据目录中依次执行这两个文件，因此全新 PostgreSQL 安装不依赖 MongoDB。

容器初始化脚本只对空数据目录执行 schema 和 seed。应用启动也会检查 PostgreSQL：
没有业务表，或者 `public` 下所有业务基础表均为空（忽略 `gemsnote_schema_migrations`），
会执行 schema 和 seed，覆盖此前只建表未导种子的情况。任意业务表已有记录则跳过 seed。
迁移目标在数据导入完成前不要启动应用。已有业务数据库由版本迁移机制增量升级；升级前必须备份。

## 数据库版本迁移

当前应用及数据库版本为 `1.0.0`，集中定义在 `app/version/version.go`。启动时会依次执行 `app/db/migrations.go` 中尚未应用的迁移，并把结果写入：

- MongoDB：`gemsnote_schema_migrations` 集合；
- PostgreSQL：`gemsnote_schema_migrations` 表。

无版本记录的旧 Leanote MongoDB 会被识别为基线数据库，在不修改业务数据的情况下登记 `1.0.0`。数据库存在比应用更新的迁移记录时，应用会拒绝启动。

## 数据迁移

双向迁移统一使用：

```bash
go run ./tools/migration -direction mongo_to_pg ...
go run ./tools/migration -direction pg_to_mongo ...
```

完整参数、备份和校验流程见 [数据库迁移指南](MIGRATION_GUIDE.md)。

## 测试

未设置集成测试环境变量时，数据库集成测试会跳过，普通测试不依赖本机数据库：

```bash
go test ./...
```

### 集成测试隔离要求

以下测试会写入数据库，必须使用专用、可丢弃的测试实例/数据库，禁止填写生产库或
个人笔记库的连接信息。尤其是迁移往返测试会复制所有支持集合的数据，双向 upsert
可能覆盖同 ID 的记录；测试结束只清理自身生成的记录，不会撤销整个数据库的迁移。
不要使用部署章节中的业务库 `gemsnote`，也不要将测试环境变量长期保存在日常运行环境中。

准备专用 PostgreSQL 测试角色和空库（用具备建库权限的管理员执行）：

```sql
CREATE USER gemsnote_test WITH PASSWORD '请替换为测试专用密码';
CREATE DATABASE gemsnote_test OWNER gemsnote_test;
```

在仓库根目录，仅导入 schema，不导入 seed，也不要启动应用连接此库：

```bash
psql 'host=127.0.0.1 port=5432 user=gemsnote_test password=请替换为测试专用密码 dbname=gemsnote_test sslmode=disable' \
  -v ON_ERROR_STOP=1 -f database/schema.sql
```

往返测试本身不负责建表；还需准备独立的 MongoDB 测试实例，以下命令中的
`gemsnote_test` 应是空业务库。运行不同测试组时使用各自干净的测试库，避免 CRUD
测试初始化数据影响往返校验。端口、认证参数和密码按测试环境调整。

PostgreSQL CRUD 集成测试：

```bash
GEMSNOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=gemsnote_test password=请替换为测试专用密码 dbname=gemsnote_test sslmode=disable' \
  go test ./app/db -run TestPostgresCRUDContract -v
```

双向迁移往返测试：

```bash
GEMSNOTE_INTEGRATION_MONGO_URL='mongodb://127.0.0.1:27017/gemsnote_test' \
GEMSNOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=gemsnote_test password=请替换为测试专用密码 dbname=gemsnote_test sslmode=disable' \
  go test ./tools/migration -run TestRoundTripMigration -v
```
