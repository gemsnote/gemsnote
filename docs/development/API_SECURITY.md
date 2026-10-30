# API 安全审查与回归

本轮审查（2026-10-29）覆盖 `/api`、`/api2` 及它们复用的旧 Web 控制器，重点为
身份来源、资源归属、密码找回、文件上传下载、旧通配路由和图片代取。

## 已修复的问题

| 风险 | 修复 |
| --- | --- |
| 旧 token 与浏览器 Cookie 身份混用；API2 匿名文件请求可能回退到 Cookie 中残留的旧身份 | token 解析结果成为唯一身份；无效 API2 凭据不回退；旧文件入口校验 Cookie 是否仍有效 |
| 邮箱验证、邮箱变更令牌可被密码找回流程使用；过期时长共享变量存在竞态 | 查询及校验同时限定令牌用途，过期时长改为局部变量 |
| 旧通配路由提供 GET 修改入口；标签缺少登录拦截 | API2 只接受显式路由；补齐标签认证；Cookie 修改操作要求 POST 和 `X-Requested-With: XMLHttpRequest` |
| 图片复制接口信任请求中的源用户与目标用户 | 先验证当前用户对源文件的读权限，只能复制到自己的账号 |
| 附件打包将笔记标题用作磁盘路径，归档条目可含路径穿越 | 使用独立随机临时文件，清理下载名及归档条目名，响应结束删除临时文件；拒绝越过实例目录的文件路径和符号链接 |
| 远程图片代取可访问内网、无下载上限 | 仅 HTTP(S) 公网地址；连接时解析并校验每个 IP，直接连接校验后的 IP；重定向仍走相同检查；不使用环境代理；30 秒、32 MiB、最多 5 次跳转 |
| 以图片后缀上传 HTML/SVG 等主动内容 | 按实际图片头识别栅格格式，保存时使用检测出的后缀；鉴权文件响应增加 `nosniff`、沙箱 CSP 与禁止缓存头 |
| 旧 PDF 导出携带并打印应用密钥，且历史博客标记绕过权限 | 改为单笔记、5 分钟有效的 HMAC 凭据；检查当前读权限；直接传参启动工具，不拼 shell、不输出凭据；禁止本地文件读取并设置执行超时 |

## 兼容范围

- `/api` 的显式 token 调用保留原有 GET/POST 及响应格式。Cookie 不再覆盖显式 token。
- API2 已声明的合法路由、编码格式、Bearer 与 token 参数不变；未声明的通配动作和错误 HTTP 方法不再作为可调用接口。
- 内网 URL 代取、伪装图片、跨账号复制等不安全行为不保留兼容。图片仍可直接作为文件上传。
- `/note/toPdf` 的内部 `appKey` 参数不再接受长期应用密钥，只接受短期笔记凭据；PDF 导出 API 由服务端自动生成该凭据。

## 验证方式与范围

```sh
go test ./app/...
go test -race ./app/controllers/... ./app/lea/... ./app/service
```

双账号 HTTP 回归见 [scripts/test-api-security.py](../../scripts/test-api-security.py)。
它会创建用户、笔记、附件并修改测试账号密码，**只能用于一次性本地实例**。
测试 PostgreSQL 容器名须以 `gemsnote-api2-` 开头，数据库与用户均为夹具配置
（数据库 `gemsnote`，用户 `gemsnote_test`）；实例包含默认种子账号。

```sh
GEMSNOTE_DISPOSABLE_API2=1 \
GEMSNOTE_API2_URL=http://127.0.0.1:19090 \
GEMSNOTE_TEST_PG_CONTAINER=gemsnote-api2-your-disposable-test \
python3 scripts/test-api-security.py
```

本轮使用独立 PostgreSQL 实例验证身份隔离、权限拒绝、CSRF 拒绝、正常上传下载、
归档内容、恶意图片拒绝、内网请求拒绝及密码找回令牌用途。没有测试 MongoDB；
PDF 的凭据与命令构造已修复，未将外部 wkhtmltopdf 渲染器视为完成沙箱或渗透测试。
本记录是代码审查与定向回归结果，不等同于全量安全认证或生产环境渗透测试。

## Web 认证回归

Web 认证回归使用 `frontend/tests/web-auth.mjs`，连接带内置种子的独立临时
PostgreSQL/API2 实例（仅限 `127.0.0.1:19090`），验证旧 Cookie、过期/撤销 token、
Bearer 登录和页面刷新。显式确认目标是一次性环境后运行：

```sh
GEMSNOTE_DISPOSABLE_API2=1 GEMSNOTE_API2_URL=http://127.0.0.1:19090 \
  npm run test:web-auth --prefix frontend
```

不要将该用例指向实际使用的服务；测试会创建和撤销该临时实例的登录会话。

## 管理员角色与状态回归

角色权限回归：`npm --prefix frontend run test:admin-roles` 检查中英文管理界面和自己的角色按钮禁用。`scripts/test-admin-roles.py` 仅允许通过 `GEMSNOTE_DISPOSABLE_API2=1` 对临时 `127.0.0.1:19090` 实例执行，测试角色授权、撤销、自我降级拒绝和并发互相降级；它会修改测试管理员角色，禁止指向实际实例。

用户状态回归：`scripts/test-user-flags.py` 在上述临时 API2 实例中验证禁用、恢复、禁止自我禁用、Cookie／API token 撤销和并发操作；需设置 `GEMSNOTE_DISPOSABLE_API2=1` 和临时 PostgreSQL 容器名 `GEMSNOTE_TEST_PG_CONTAINER=gemsnote-api2-test-<PID>-<随机数>`。禁止对实际实例执行。中英文状态按钮由 `npm --prefix frontend run test:admin-roles` 覆盖。
