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

Gemsnote 支持 PostgreSQL 和 MongoDB 两种数据库，数据库访问由抽象层统一处理。具体的数据库配置、全新初始化、从 Leanote 迁移和启动方式请参阅[快速开始](docs/QUICK_START.md)。

## 文档

- [快速开始](docs/QUICK_START.md)
- [部署、测试与交付](docs/DEPLOYMENT.md)
- [开发文档总览](docs/development/README.md)

## API 与客户端兼容性

Gemsnote 保留了 Leanote 原有 HTTP API 路径、参数和响应字段。旧版 Leanote Electron 客户端可以在登录界面选择自建服务并填写 Gemsnote 服务端地址进行连接。

旧客户端默认地址仍可能指向 `leanote.com`，因此连接私有部署时必须明确填写服务器地址。

## 相关项目

- [Leanote 原项目](https://github.com/leanote/leanote)
- [Leanote Desktop 原项目](https://github.com/leanote/desktop-app)

## 参与贡献

欢迎通过 Issue 和 Pull Request 报告问题、提出建议或贡献代码。提交涉及数据库结构、迁移逻辑或 API 的修改时，请同时补充相应测试，并确保 MongoDB、PostgreSQL 和旧客户端兼容性不受影响。
