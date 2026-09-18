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

容器初始化脚本只对空数据目录执行完整 Schema。已有数据库由服务启动时的版本迁移机制增量升级；升级前仍必须备份。

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

完整参数、备份和校验流程见 [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md)。

## 测试

普通测试不依赖本机数据库：

```bash
go test ./...
```

PostgreSQL CRUD 集成测试：

```bash
GEMSNOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable' \
  go test ./app/db -run TestPostgresCRUDContract -v
```

双向迁移往返测试：

```bash
GEMSNOTE_INTEGRATION_MONGO_URL='mongodb://127.0.0.1:27017/gemsnote' \
GEMSNOTE_INTEGRATION_POSTGRES_URL='host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable' \
  go test ./tools/migration -run TestRoundTripMigration -v
```
