# 开发文档

本目录集中维护 Gemsnote 的开发、架构、接口、迁移、UI 和构建发布文档。面向需要修改代码、扩展数据库后端、开发客户端或维护发布流程的贡献者。

## 数据库与迁移

- [数据库抽象指南](DATABASE_ABSTRACTION_README.md)：数据库后端选择、统一 ID、查询兼容层、版本迁移和集成测试。
- [数据库迁移指南](MIGRATION_GUIDE.md)：MongoDB 与 PostgreSQL 之间的数据迁移、备份、校验及附件迁移。

## 接口与客户端

- [API2 新接口](API2.md)：新 Web 和后续客户端使用的 `/api2` 接口规范。
- [Legacy API 兼容说明](API.md)：旧 `/api` 接口、Leanote 客户端兼容范围及限制。

## Web 与 UI

- [统一 UI 设计](UI.md)：Web 和 desktop 共享的界面结构与交互设计。
- [Vue Web UI 重写与部署计划](WEB_UI_MIGRATION_PLAN.md)：Web 前后端分离实现记录和维护边界。

## 构建与发布

- [自动化构建与发布](RELEASE.md)：本地构建、GitHub Actions、多平台客户端、服务端及镜像发布流程。

部署和日常使用文档仍位于上一级 `docs` 目录：

- [快速开始](../QUICK_START.md)
- [部署、测试与交付](../DEPLOYMENT.md)
