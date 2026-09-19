# Gemsnote Legacy API (`/api`)

本目录下的 `/api` 接口是 Leanote 兼容层，继续保留用于旧版 Leanote Desktop 等客户端。它使用历史上的表单参数、token/session 混合认证和部分不一致的响应结构，不再增加新功能。

新 Web 和后续客户端请使用 [API2](API2.md)。旧接口的完整动作及参数列表见 `app/controllers/api/API列表-v0.1.md`；该文件仅作为兼容参考，不代表新接口规范。

兼容层的主要限制：

- 参数以 `application/x-www-form-urlencoded` 为主，数组使用 `Field[0]` 形式。
- 个别修改操作历史上使用 GET，无法可靠表达幂等性和缓存语义。
- 成功响应有时直接返回对象，有时返回 `Ok` 包装对象。
- 错误文本可能经过本地化，客户端应按旧协议处理。

这些限制不会影响 API2；两套接口在迁移期间并行运行。
