# Gemsnote 数据库迁移指南

Gemsnote 支持 MongoDB 和 PostgreSQL。两种后端统一使用 MongoDB ObjectId：24 位十六进制字符串。PostgreSQL 使用 `CHAR(24)` 保存 ID，因此迁移不会生成 UUID，也不需要额外的 ID 映射表。

## 准备

1. 迁移前备份源数据库和目标数据库。
2. 停止 Gemsnote 写入，避免迁移期间源数据变化。
3. 确认目标数据库为空，或确认同 ID 记录可以被更新。
4. PostgreSQL schema 位于 `database/schema.sql`。MongoDB→PostgreSQL 时工具默认会执行它。

## MongoDB → PostgreSQL

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

先只读取并检查连接和表结构：

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -dry-run \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

如果 PostgreSQL schema 已由其他工具管理，可传入 `-schema ''` 禁止自动执行。

## PostgreSQL → MongoDB

先确保目标 MongoDB 已备份，然后执行：

```bash
go run ./tools/migration \
  -direction pg_to_mongo \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

工具按照 `_id` 执行 upsert，并将 PostgreSQL 的 `CHAR(24)` ID 转回 `bson.ObjectId`。

## 校验

默认启用 `-validate=true`，迁移结束后比较所有 collection/table 的记录数。任何表不一致都会返回非零退出状态。

记录数校验不能替代业务验证。切换数据库前还应检查：

- 登录与用户配置；
- 笔记本、笔记、正文和历史版本；
- 标签、搜索和同步 USN；
- 分享权限；
- 博客、评论和点赞；
- 图片与附件元数据。

## 切换应用后端

MongoDB：

```ini
db.type=mongodb
db.url=mongodb://127.0.0.1:27017/gemsnote
```

PostgreSQL：

```ini
db.type=postgresql
db.url=host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable
```

切换后先以维护模式启动，完成登录、读取和写入冒烟测试，再恢复外部流量。迁移工具不会复制附件文件本身，只迁移数据库中的附件元数据；文件目录需要单独保持或迁移。
