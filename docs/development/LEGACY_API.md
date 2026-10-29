# Gemsnote Legacy API (`/api`)

本目录下的 `/api` 接口是 Leanote 兼容层，继续保留用于旧版 Leanote Desktop 等客户端。它使用历史上的表单参数、token/session 混合认证和部分不一致的响应结构，不再增加新功能。

新 Web 和客户端请使用 [API2](API2.md)。旧接口的完整动作及参数列表见[历史 API 列表](../../app/controllers/api/API列表-v0.1.md)；该文件仅作为兼容参考，不代表新接口规范。

兼容层的主要限制：

- 参数以 `application/x-www-form-urlencoded` 为主，数组使用 `Field[0]` 形式。
- 个别修改操作历史上使用 GET，无法可靠表达幂等性和缓存语义。
- 显式 token 决定当前身份，Cookie 不得覆盖它。旧 token 客户端保留原有调用方式；仅依赖 Cookie 的修改请求必须使用 POST 并携带 `X-Requested-With: XMLHttpRequest`，否则返回 405 或 403。
- 成功响应有时直接返回对象，有时返回 `Ok` 包装对象。
- 错误文本可能经过本地化，客户端应按旧协议处理。

两套接口并行运行。API2 的部分动作仍复用表单绑定和历史响应结构，不能仅凭 `/api2`
前缀就假定全部接受 JSON 或使用统一错误包装；具体合同见 API2 文档及对应控制器。

文件、图片代取及密码找回的安全边界见 [API 安全审查](API_SECURITY.md)。
