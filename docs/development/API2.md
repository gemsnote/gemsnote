# Gemsnote API2

## 文档基本信息

| 项目 | 说明 |
| --- | --- |
| 接口前缀 | `/api2` |
| 使用方 | Web、新桌面客户端和 Flutter 移动客户端 |
| 当前版本 | `1.2.0`；`1.0.0` 为首个定稿版本 |
| 版本来源 | [app/version/version.go](../../app/version/version.go)，同时用于服务端版本响应和数据库迁移 |
| 路由与编码 | [conf/routes](../../conf/routes)、[Web 请求实现](../../frontend/src/api.ts)及对应控制器 |
| Legacy 接口 | 旧 `/api` 继续保留，详见 [LEGACY_API.md](LEGACY_API.md) |

本文按 1.2.0、1.1.0、1.0.0 的顺序介绍新增与保留的接口。
1.2.0 保留 1.0.0、1.1.0 的既有调用方式，`min_version` 仍为空。
接口前缀不代表统一的编码格式：JSON、表单和 multipart 必须按各端点要求使用。
API2 只开放路由表明确声明的动作和 HTTP 方法，不通过旧通配路由补充别名。
认证、上传及文件访问的安全约束见 [API 安全审查与回归](API_SECURITY.md)。

## 1.2.0 增加的部分

### 用户角色

不提供用户自行注册功能。账号只能由管理员通过 Web 用户管理创建；`/register`、`/doRegister`、`/api2/auth/register` 均不是公开注册入口。

- `POST /api2/admin/role`：JSON `{"userId":"<用户 ID>","role":"admin"}`，`role` 仅支持 `admin`、`user`。仅管理员可调用。成功返回 `Ok:true`，普通用户调用返回 HTTP 403 和 `noAuth`。
- 禁止管理员将自己设为 `user`，返回 `Ok:false, Msg:cannotDemoteSelf`。修改别人的角色时重新检查操作者权限；PostgreSQL 在事务锁内完成，避免并发相互降级导致无管理员。
- `/api2/admin/data` 的用户列表新增 `Role` 字段，Web 提供角色设置入口。新用户默认 `user`。角色变化在后续请求立即生效，不要求退出登录。
- 数据库迁移为用户添加角色字段；启动时只初始化尚无角色的原管理员，不会覆盖明确的角色设置。旧 `/api` 的用户 JSON 不新增该字段。

### 邮箱验证任务与管理员身份

- `POST /api2/web/emailChange`：表单 `email`、`pwd` 验证当前密码并向新地址排队发送验证邮件；表单 `current=true` 则重发当前邮箱验证邮件。返回 `{"Ok":true,"TaskId":"<任务 ID>","State":"queued"}`。邮箱此时**尚未修改**。
- `POST /api2/user/emailStatus`：表单 `taskId`，必须使用本人 Bearer token；返回 `TaskId`、`State`（`queued`、`sending`、`sent`、`failed`）及失败时的 `Error`。`sent` 仅表示 SMTP 接受邮件。别人的任务和不存在的任务均返回 `Ok:false, Msg:mailTaskNotFound`。
- `POST /api2/web/verifyEmail`：现有公开验证接口，表单 `token`、`change=true`，完成新邮箱验证后才修改邮箱；成功消费令牌，不可重复使用。
- `GET /api2/bootstrap` 的 `IsAdmin` 按用户记录的 `Role=admin` 判断，不再依赖可修改的用户名；桌面端仍不提供管理界面。

发送错误包括 `mailNotConfigured`、`mailQueueNotConfigured`、`mailQueueUnavailable`、`mailRateLimited`、`mailSiteURLInvalid`、`mailTemplateInvalid`、`mailTokenFailed`；投递错误包括 `mailAuthenticationFailed`、`mailRecipientRejected`、`mailConnectionFailed`、`mailDeliveryFailed`、`mailTaskExpired`。邮箱相同返回 `emailUnchanged`，密码不正确返回 `wrongPassword`。部署与状态语义见[异步邮件说明](../MAIL.md)。

## 1.1.0 增加的部分

本版本增加笔记本创建幂等、图片上传回执和安全清理：

| 方法 | 路径 | 相比 1.0.0 的变化 |
| --- | --- | --- |
| POST | `/api2/client/notebook/add` | 增加可选创建键 `clientNotebookId` |
| POST | `/api2/file/pasteImage` | 增加可选上传键 `uploadId` 和可查询的上传记录 |
| GET | `/api2/file/imageUpload` | 新增上传回执查询 |
| POST | `/api2/file/cleanupImageUpload` | 新增未引用上传的显式清理 |

### 笔记本创建幂等

`POST /api2/client/notebook/add` 接受表单或 JSON：

```json
{"title":"随记","parentNotebookId":"","seq":0,"clientNotebookId":"0123456789abcdef01234567"}
```

`clientNotebookId` 是可选的 24 位十六进制 ObjectId 创建键，同时成为服务端
`NotebookId`。客户端须在第一次发送前持久保存键和请求内容。重复请求返回同一
笔记本的当前状态，不修改标题、父级、时间或 USN；遇到已删除记录返回
`creationKeyDeleted`，遇到其他所有者的同 ID 返回 `creationKeyUnavailable`，不返回其数据。
主键唯一约束处理并发重复请求。省略键沿用随机 ID 创建，不保证重试幂等。

### 图片上传回执

`POST /api2/file/pasteImage` 使用 multipart，字段为 `noteId`、`file` 和可选的
`uploadId`。目标笔记必须存在且属于当前用户。`uploadId` 使用 16–128 位 ASCII
字母、数字、下划线或短横线，建议使用持久化的随机任务 ID；同一账号内唯一。
同键、同笔记、同 SHA-256 重试返回同一文件 ID，不重复写入；不同载荷返回
`uploadKeyConflict`，清理后的键返回 `uploadAlreadyCleaned`。文件 ID 由服务端安全
随机生成，与客户端上传键独立；不能在上传记录提交前推算其 ID 并提前写入引用。
既有 Web 在新笔记保存前省略 `noteId`、`uploadId` 的上传仍可用，但不产生可清理回执。

`GET /api2/file/imageUpload?uploadId=...` 或 `?fileId=...` 仅查询当前所有者的新式上传：

```json
{"Ok":true,"Id":"0123456789abcdef01234567","Item":{"FileId":"0123456789abcdef01234567","NoteId":"abcdef0123456789abcdef01","Digest":"<64位SHA-256>","State":"pending","CanCleanup":true}}
```

摘要对应上传原始字节。`pending` 表示尚未被任何已保存引用采用；`retained` 表示
已经永久保留；`deleting` 表示清理进行中，可重试；`deleted` 表示已清理。未知、
旧式或其他账号的文件统一返回 `uploadNotFound`。数据库读取失败不会冒充“不存在”。

### 显式清理

`POST /api2/file/cleanupImageUpload` 接受 `{"fileId":"..."}` 或等价表单，成功返回
`{"Ok":true}`。仅允许清理本人新式、从未被保存引用采用的图片，不处理既有图片、
头像或附件。正文、封面、图片关系、用户资料和历史写入前先原子标记保留，清理
通过相同状态的 CAS 竞争；不使用单进程锁证明安全。已保留返回
`imageAlreadyReferenced`，即使当前正文已移除或历史已裁剪，也不自动解除保留。
清理保留数据库 tombstone，阻止后来引用及同键重传，重复清理可安全重试。
物理删除限制在应用 `files/` 下，拒绝路径穿越及符号链接。

### 数据库升级与部署

PostgreSQL 1.1.0 迁移新增上传状态、目标笔记、摘要、上传键和部分唯一索引；
MongoDB 使用稀疏唯一上传键索引，旧文档不回填、不纳入清理。部署时先停掉旧
服务端进程，再升级并恢复写入；不得让缺少引用保护的旧进程同时写入同一数据库。

## 1.0.0 原来的部分

以下接口和规则在 1.0.0 已有，后续版本继续提供。响应示例中的 `Version` 取当前
服务端版本，不表示这些字段是 1.1.0 新增的。

### 笔记修改时间

API2 客户端笔记新增/更新接口中，省略可选的 `UpdatedTime` 表示使用服务端时间。
同步版本判断使用 USN；本地缓存实现见[客户端开发指南](https://github.com/gemsnote/gemsnote-desktop/blob/main/docs/DEVELOPMENT.md)。
服务端为同次请求的笔记元数据和正文使用同一个时间，更新响应不再用另一个
“响应生成时间”覆盖实际写入时间。显式传入历史时间的既有客户端合同及旧 `/api`
保持不变；无需新增必填字段或改变路径。`CreatedTime` 仍保留原始创建时间。

### 认证

Web 和新客户端均使用 `POST /api2/auth/login`，推荐 `application/json` 请求体
（`email`、`pwd`），也接受等价表单。该接口在一次响应中返回 token、用户资料和
服务端协议版本，登录流程不需要再请求 `/api2/user/info` 或 `/api2/system/version`：

```json
{
  "Ok": true,
  "Token": "...",
  "User": {"UserId": "...", "Username": "admin", "Email": "...", "Verified": true, "Logo": "..."},
  "Server": {"Name": "gemsnote", "Version": "1.2.0", "MinVersion": ""}
}
```

受保护的 API2 请求统一使用 `Authorization: Bearer <token>`，包括 Web 的 JSON、
表单和文件上传，不依赖浏览器 Cookie 或重复登录。token 由 `crypto/rand` 生成
32 字节安全随机数并以 Base64URL 编码，不要求 JWT 格式；退出及密码变更会撤销
相应凭据。不要将 token 写入正文、日志或错误信息。
服务端暂保留 `token` 参数供已有客户端使用，新代码不要将凭据放入 URL。
`POST /api2/auth/session` 仍保留浏览器会话登录动作（`email`、`pwd`、可选 `captcha`），
但它不能代替受保护 API2 所需的 Bearer 认证。Web 实现仍发送
`X-Requested-With: XMLHttpRequest`，该请求头也不能代替身份凭据。
`/api/doLogin` 仅供旧 Leanote 兼容客户端使用。

### Web 核心接口

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

保存现有笔记时应显式发送 `isMarkdown`。省略该字段时保留现有格式；显式传 `false`
仍表示富文本，可用于恢复富文本历史版本。新建笔记省略该字段时默认富文本。

`GET /api2/bootstrap` 保持公开：缺失、失效或冲突的 API2 凭据返回 `User:null`，
携带有效 Bearer token 时返回该 token 对应的用户及数据。初始化接口和受保护接口
使用相同的身份来源，不会用旧浏览器 Cookie 把匿名 API2 请求判为“已登录”。
因此升级前留下的 Cookie 不会让页面在 `/note` 与 `/login` 之间反复跳转。

文件上传 `/api2/avatar`、`/api2/attachments/upload` 使用 `multipart/form-data`，
文件字段为 `file`，不能把文件编码为普通 JSON 请求。文件读取使用
`/api2/file/getImage`、`/api2/file/getAttach` 和 `/api2/attachments/download`。

API2 并非所有接口均接受 JSON。例如 `/api2/user/updatePwd`、
`/api2/notebook/updateNotebookTitle` 和部分共享/管理动作仍使用
`application/x-www-form-urlencoded`。客户端笔记写入 `/api2/note/addNote`、
`/api2/note/updateNote` 及 `/api2/client/*` 动作也复用表单绑定。
完整路由以 `conf/routes` 为准，请求编码可对照 `frontend/src/api.ts` 和对应控制器；
不要仅根据路径前缀推断请求格式或将浏览器动作当成 token 同步接口。

### 笔记本移动

`POST /api2/notebook/moveNotebook` 接受 JSON：
`{"notebookId":"源笔记本ID","parentNotebookId":"目标父笔记本ID"}`。
`parentNotebookId` 为空表示移到根目录。请求携带 Bearer token，返回
`{"Ok":true}` 或失败信息。只修改父级，保留子树及笔记归属；拒绝其他账号、
已删除、自身或子孙目标。旧 `/api` 不变。

### 客户端正文同步合同

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

### 共享与历史恢复

共享原文只读，只有原作者可以修改；接收者通过复制到自己的笔记本生成独立副本。
副本使用新的笔记和资源 ID，并归复制者所有，修改副本不影响原文。
历史附件保留、完整快照标记与恢复接口见 [HISTORY_RESTORE.md](HISTORY_RESTORE.md)。

### 错误

当前没有统一的 HTTP 状态码与错误包装合同。客户端必须同时检查 HTTP 状态和
端点的 JSON 业务结果：受保护 API2 缺失或失效凭据返回 HTTP 401 和
`{"Ok":false,"Msg":"NOTLOGIN"}`，不会重定向登录页；登录失败、参数或业务错误
仍可能以 HTTP 200 返回 `{"Ok":false,"Msg":"..."}`。部分防护会返回 403/405，
历史动作还可能返回布尔值。不能将 HTTP 200 一律视为成功。

`Msg` 既可能是固定错误码，也可能是本地化文本（如登录密码错误）。仅对已明确约定的
错误码进行逻辑分支，其余作为提示展示，不能假设所有 `Msg` 都是稳定、未翻译的机器码。

### 版本策略

`/api` 是 legacy 接口，当前版本不会移除。API2 `1.0.0` 已定稿，后续合同变更须升级版本。
新客户端通过 `/api2/system/version` 获取 `{"server":"gemsnote","version":"1.2.0","min_version":""}`，
或使用登录响应的 `Server` 字段（注意两者字段大小写不同）。当未来设置 `min_version`
后，低于该版本的客户端需提示升级；空值表示保留旧 Leanote 客户端兼容性。
