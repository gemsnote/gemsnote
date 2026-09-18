# 数据库抽象与迁移

当前实现说明统一维护在：

- [数据库抽象说明](DATABASE_ABSTRACTION_README.md)
- [双向迁移指南](MIGRATION_GUIDE.md)

重要约定：MongoDB 和 PostgreSQL 均使用 24 位 MongoDB ObjectId，不使用 UUID，也不需要 ID 映射表。
