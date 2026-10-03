# Gemsnote API2

## 笔记本移动

`POST /api2/notebook/moveNotebook` 接受 JSON：
`{"notebookId":"源笔记本ID","parentNotebookId":"目标父笔记本ID"}`。
`parentNotebookId` 为空表示移到根目录。Web 使用已登录会话和
`X-Requested-With: XMLHttpRequest`，返回 `{"Ok":true}` 或失败信息。
只修改父级，保留子树及笔记归属；拒绝其他账号、已删除、自身或子孙目标。
本接口为 API2 增量扩展，旧 `/api` 不变。

API2 是 Web 和新客户端使用的接口空间，前缀为 `/api2`，`1.0.0` 为首个定稿版本。
旧 `/api` 作为 Leanote 客户端兼容层继续保留，已移除的旧 Web/博客功能不在此范围。

## 笔记修改时间

API2 客户端笔记新增/更新接口中，省略可选的 `UpdatedTime` 表示使用服务端时间。
同步版本判断使用 USN；本地缓存实现见[客户端开发指南](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/DEVELOPMENT.md)。
服务端为同次请求的笔记元数据和正文使用同一个时间，更新响应不再用另一个
“响应生成时间”覆盖实际写入时间。显式传入历史时间的既有客户端合同及旧 `/api`
保持不变；无需新增必填字段或改变路径。`CreatedTime` 仍保留原始创建时间。

## 认证

Web 登录使用 `POST /api2/auth/session`，请求体为 `application/json`（`email`、`pwd`、可选 `captcha`），成功后设置同源 session cookie。Token 客户端使用 `POST /api2/auth/login`。该接口在一次响应中返回认证 token、用户资料和服务端协议版本，登录流程不需要再请求 `/api2/user/info` 或 `/api2/system/version`：

```json
{
  "Ok": true,
  "Token": "...",
  "User": {"UserId": "...", "Username": "admin", "Email": "...", "Verified": true, "Logo": "..."},
  "Server": {"Name": "gemsnote", "Version": "1.0.0", "MinVersion": ""}
}
```

两种登录均推荐 JSON，字段为 `email`、`pwd`（服务端也接受表单）。Web 后续请求
携带同源 cookie，受保护的写操作还需 `X-Requested-With: XMLHttpRequest`；
客户端同步请求携带 `token` 查询参数。不要将 token 认证和浏览器会话认证混用。
`/api/doLogin` 仅供旧 Leanote 兼容客户端使用。

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
| GET | `/api2/groups` | 当前用户分组 |
| POST | `/api2/admin/data` | 管理数据查询（JSON） |
| POST | `/api2/attachments` | 附件列表（JSON） |
| POST | `/api2/attachments/delete` | 删除附件（JSON） |

文件上传 `/api2/avatar`、`/api2/attachments/upload` 使用 `multipart/form-data`，
文件字段为 `file`，不能把文件编码为普通 JSON 请求。文件读取使用
`/api2/file/getImage`、`/api2/file/getAttach` 和 `/api2/attachments/download`。

API2 并非所有接口均接受 JSON。例如 `/api2/user/updatePwd`、
`/api2/notebook/updateNotebookTitle` 和部分共享/管理动作仍使用
`application/x-www-form-urlencoded`。客户端笔记写入 `/api2/note/addNote`、
`/api2/note/updateNote` 及 `/api2/client/*` 动作也复用表单绑定。
完整路由以 `conf/routes` 为准，请求编码可对照 `frontend/src/api.ts` 和对应控制器；
不要仅根据路径前缀推断请求格式或将浏览器动作当成 token 同步接口。

## 客户端正文同步合同

以下读取接口使用 token 认证：

| GET 路径 | 正文语义 |
| --- | --- |
| `/api2/note/getSyncNotes` | 增量元数据，不加载正文 |
| `/api2/note/getNote` | 单篇元数据，不加载正文 |
| `/api2/note/getNoteContent` | 通过 `noteId` 读取真实正文，空字符串是有效正文 |
| `/api2/note/getSyncNotesWithContent` | 元数据与真实正文分页快照，`maxEntry` 最大为 50 |

前两项响应中的 `Content:""` 是共享响应模型产生的占位值，**不表示正文为空或已下载**。
客户端不能根据 `Content` 字段存在就覆盖本地正文或跳过正文下载。处理远端版本时，
应在正文成功写入本地后才确认该版本已完整同步；修复历史缓存时不能仅凭 USN 相同跳过正文。
正文快照返回的空正文则是真实值，不应当作缺失。

正文快照根据 `Accept-Encoding` 协商 gzip 压缩，不依赖全局压缩开关，解压后 JSON
结构不变。该优化只应用于 API2 快照接口，旧 `/api` 响应不变。

## 错误

当前没有统一的 HTTP 状态码与错误包装合同。客户端必须同时检查 HTTP 状态和
端点的 JSON 业务结果：部分认证/参数错误仍以 HTTP 200 返回 `{"Ok":false,"Msg":"..."}`，
例如 token 校验的 `NOTLOGIN`；部分浏览器防护会返回 403/405，历史动作还可能返回
布尔值。不能将 HTTP 200 一律视为成功。

`Msg` 既可能是固定错误码，也可能是本地化文本（如登录密码错误）。仅对已明确约定的
错误码进行逻辑分支，其余作为提示展示，不能假设所有 `Msg` 都是稳定、未翻译的机器码。

## 版本策略

`/api` 是 legacy 接口，当前版本不会移除。API2 `1.0.0` 已定稿，后续合同变更须升级版本。
新客户端通过 `/api2/system/version` 获取 `{"server":"gemsnote","version":"1.0.0","min_version":""}`，
或使用登录响应的 `Server` 字段（注意两者字段大小写不同）。当未来设置 `min_version`
后，低于该版本的客户端需提示升级；空值表示保留旧 Leanote 客户端兼容性。
