# Gemsnote API2

API2 是新客户端使用的版本化接口空间，前缀为 `/api2`。第一阶段先完成 Web 迁移；旧的 `/api` 以及旧控制器路由继续保留，供 Leanote/Desktop 兼容客户端使用，暂不删除。

## 认证

Web 登录使用 `POST /api2/auth/session`，请求体为 `application/json`（`email`、`pwd`、可选 `captcha`），成功后设置同源 session cookie。Desktop 与服务端之间使用 `POST /api2/auth/login`。该接口在一次响应中返回认证 token、用户资料和服务端协议版本，登录流程不需要再请求 `/api2/user/info` 或 `/api2/system/version`：

```json
{
  "Ok": true,
  "Token": "...",
  "User": {"UserId": "...", "Username": "admin", "Email": "...", "Verified": true, "Logo": "..."},
  "Server": {"Name": "gemsnote", "Version": "1.0.0", "MinVersion": ""}
}
```

API2 当前处于定稿前阶段，客户端应按上述结构实现，不提供对早期 Gemsnote API2 登录响应的兼容。旧版登录路径 `/api/doLogin` 仅供 Leanote 兼容客户端使用。

## Web 核心接口

以下写接口使用 `application/json`，成功响应沿用当前数据结构，以便 Web 平滑切换：

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api2/bootstrap` | 当前用户、笔记本、标签及统计信息 |
| POST | `/api2/notes` | 查询笔记；字段 `notebookId`、`key`、`tag`、`sort`、`trash`、`starred`、`page` |
| POST | `/api2/document` | 读取笔记及正文；字段 `noteId` |
| POST | `/api2/save` | 新建或保存笔记；字段 `noteId`、`notebookId`、`ownerId`、`title`、`content`、`tags`、`isNew`、`isMarkdown`、`usn` |
| POST | `/api2/star` | 设置星标；字段 `noteId`、`starred` |
| POST | `/api2/restore` | 从回收站恢复；字段 `noteId` |

登录和上传等过渡接口仍使用表单或 multipart；上传、分组、共享和管理接口在第一阶段通过 `/api2/<旧控制器>/<动作>` 暴露，Web 已统一从 API2 命名空间访问；这些接口的 JSON 化和更严格的资源语义将在后续小版本中逐步补齐。文件读取使用 `/api2/file/getImage`、`/api2/file/getAttach`。

## 错误

HTTP 状态码用于区分认证和请求错误；业务失败返回 `Ok:false`，并在 `Msg` 中提供稳定错误码。客户端不得依赖本地化后的错误文本。

## 版本策略

`/api` 是 legacy 接口，当前版本不会移除。新客户端应只使用 `/api2`，并通过 `/api2/system/version` 判断服务端版本；当未来版本设置 `min_version` 后，低于该版本的客户端必须提示升级。
