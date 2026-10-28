# 开发文档

本目录集中维护 Gemsnote 的开发、架构、接口、迁移、UI 和构建发布文档。面向需要修改代码、扩展数据库后端、开发客户端或维护发布流程的贡献者。

## 数据库与迁移

- [数据库抽象指南](DATABASE_ABSTRACTION_README.md)：数据库后端选择、统一 ID、查询兼容层、版本迁移和集成测试。
- [数据库迁移指南](MIGRATION_GUIDE.md)：MongoDB 与 PostgreSQL 之间的数据迁移、备份、校验及附件迁移。

## 接口与客户端

- [API2 新接口](API2.md)：Web 和新客户端使用的 `/api2` 接口规范，分别列出 1.1.0 新增合同与 1.0.0 原有接口。
- [Legacy API 兼容说明](LEGACY_API.md)：旧 `/api` 接口、Leanote 客户端兼容范围及限制。

Web 认证回归使用 `frontend/tests/web-auth.mjs`，连接带内置种子的独立临时
PostgreSQL/API2 实例（仅限 `127.0.0.1:19090`），验证旧 Cookie、过期/撤销 token、
Bearer 登录和页面刷新。显式确认目标是一次性环境后运行：

```sh
GEMSNOTE_DISPOSABLE_API2=1 GEMSNOTE_API2_URL=http://127.0.0.1:19090 \
  npm run test:web-auth --prefix frontend
```

不要将该用例指向实际使用的服务；测试会创建和撤销该临时实例的登录会话。

## Web 与 UI

- [统一 UI 设计](UI.md)：Web 界面结构与交互设计。
- [桌面客户端文档](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/DEVELOPMENT.md)：客户端开发、缓存与同步；[桌面 UI](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/UI.md) 与 [Release 构建](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/RELEASE.md) 在客户端仓库维护。

## 构建与发布

- [自动化构建与发布](RELEASE.md)：本地构建、GitHub Actions、多平台客户端、服务端及镜像发布流程。

部署和日常使用文档仍位于上一级 `docs` 目录：

- [快速开始](../QUICK_START.md)
- [部署、测试与交付](../DEPLOYMENT.md)
