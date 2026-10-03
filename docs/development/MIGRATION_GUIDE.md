# Gemsnote 数据库迁移指南

Gemsnote 支持 MongoDB 和 PostgreSQL。两种后端统一使用 MongoDB ObjectId：24 位十六进制字符串。PostgreSQL 使用 `CHAR(24)` 保存 ID，因此迁移不会生成 UUID，也不需要额外的 ID 映射表。

## 准备

1. 迁移前备份源数据库和目标数据库。
2. 停止 Leanote/Gemsnote 等所有源库写入；目标 Gemsnote 在迁移完成前也不要启动，以免初始化 seed。
3. 优先使用独立的空目标数据库，不导入 seed；向非空库迁移前必须阅读下方“非空目标库与失败处理”。
4. PostgreSQL schema 位于 `database/schema.sql`。MongoDB→PostgreSQL 时工具默认会执行它。

## MongoDB → PostgreSQL

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

如需先预检，目标 PostgreSQL 必须已具有 schema；dry-run 不会自动建表。
在仓库或 Release 根目录执行以下命令，仅导入 schema，不导入 seed：

```bash
psql 'host=127.0.0.1 port=5432 user=gemsnote password=请替换为强密码 dbname=gemsnote sslmode=disable' \
  -v ON_ERROR_STOP=1 -f database/schema.sql
```

然后只读取并检查连接和表是否存在：

```bash
go run ./tools/migration \
  -direction mongo_to_pg \
  -dry-run -validate=false \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

如果 PostgreSQL schema 已由其他工具管理，可传入 `-schema ''` 禁止自动执行。

MongoDB→PostgreSQL 的 dry-run 只遍历源记录并检查目标表存在，不执行字段转换、
约束检查或实际写入，不是完整的迁移模拟。`-validate` 默认会比较源/目标数量，
尚未写入的空目标库必然与非空源库不同，因此预检需显式关闭。正式迁移去掉
`-dry-run -validate=false`，恢复默认记录数校验；示例连接中的密码需替换为实际密码。
Release 包可将命令入口 `go run ./tools/migration` 换为 `./bin/gemsnote-migrate`。

## PostgreSQL → MongoDB

先确保目标 MongoDB 已备份，然后执行：

```bash
go run ./tools/migration \
  -direction pg_to_mongo \
  -mongo-url 'mongodb://127.0.0.1:27017/gemsnote' \
  -postgres-url 'host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable'
```

工具按照 `_id` 执行 upsert，并将 PostgreSQL 的 `CHAR(24)` ID 转回 `bson.ObjectId`。

## 非空目标库与失败处理

迁移是逐条写入，不是覆盖式镜像恢复，也不是全程事务：

- 同 ID 的目标记录会被更新；目标独有记录不会被删除。因此非空目标可能保留额外数据，
  导致默认记录数校验失败。仅确认“允许更新同 ID”不足以保证迁移后的库与源库一致。
- 中途发生连接、转换或写入错误时，此前成功写入的记录会保留，不会自动回滚。
- 最终数量校验在复制之后执行；即使校验失败并返回非零退出状态，数据也已经写入。
  不要用 `-validate=false` 掩盖正式迁移的差异，也不要把数量相等当作内容完全一致的证明。

失败后保持目标应用停止，保存错误日志并核对源、目标连接信息及具体失败集合。
优先修复原因后改用新的空目标库重做迁移；需要撤销非空目标的变更时，从迁移前备份
恢复到经过核实的目标环境，不要删除唯一备份或清理未确认的数据目录。

只有确认源数据保持冻结、目标没有其他写入，并核查部分写入及目标独有记录后，才可
在原目标上重试。工具支持按 ID 再次写入，但这不是带检查点的断点续传，也不会清除
前次迁移留下的额外数据。重新执行默认记录数校验及业务验证，均通过后再切换服务。

## 校验

默认启用 `-validate=true`，迁移结束后比较所有 collection/table 的记录数。任何表不一致都会返回非零退出状态。

记录数校验不能替代业务验证。切换数据库前还应检查：

- 登录与用户配置；
- 笔记本、笔记、正文和历史版本；
- 标签、搜索和同步 USN；
- 分享权限；
- 图片与附件元数据。

工具仍保留旧博客、评论、点赞等集合的数据转换，但 Gemsnote 不再提供这些产品功能；
需要保留历史记录时检查数据库中的迁移结果，不应将旧博客页面作为上线验收入口。

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

切换后先限制外部访问（例如在反向代理处阻断公众流量，仅允许运维访问），完成登录、读取和写入冒烟测试，再恢复外部流量；这里不是指应用存在某个“维护模式”配置。迁移工具不会复制附件文件本身，只迁移数据库中的附件元数据；`files/` 和需要的 `public/upload/` 必须单独保持或迁移。
