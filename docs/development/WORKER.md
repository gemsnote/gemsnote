# gemsnote-worker 开发

`gemsnote-worker` 是独立的 Go/Asynq 后台进程，当前处理邮件任务。它与 Web 服务端共享模块和发布镜像，但独立运行，仅连接 Redis 和 SMTP，不初始化业务数据库或读取 `app.conf`。本地可使用 `go run ./gemsnote-worker -concurrency 2`，需先设置 `GEMSNOTE_MAIL_REDIS_URL`。运行和配置见[部署指南](../DEPLOYMENT.md#后台任务服务)，打包见[发布指南](RELEASE.md)。

| 代码位置 | 职责 |
| --- | --- |
| `gemsnote-worker/main.go` | worker 入口、并发数、任务处理器注册和优雅停止。 |
| `internal/mailqueue/` | Asynq 入队、私有载荷、发送状态和 SMTP 投递。 |
| `app/service/MailQueue.go` | 读取服务端配置、入队以及按账号授权查询结果。 |
| `app/service/EmailService.go` | 验证邮件业务校验、发送间隔和模板调用。 |
| `app/service/verification_mail.go` | 简体中文／英文验证邮件的 HTML 和纯文本模板。 |

Asynq 队列为 `gemsnote.mail`，任务类型为 `gemsnote.mail.send`，任务载荷仅包含任务 ID。邮件正文、验证链接和 SMTP 凭据保存在独立 Redis key 中，10 分钟过期、处理完成后删除；所属账号和发送状态保留 24 小时。不要把敏感数据写入 Asynq 任务参数、错误或日志。

前端查询状态为 `queued`、`sending`、`sent`、`failed`；`sent` 仅表示 SMTP 接受，不保证收件箱送达。SMTP 会话总超时为 30 秒，明确的 SMTP 失败不自动重发；Redis 基础设施错误最多重试 3 次。发送前的持久化标记用于避免恢复任务时重复投递；若中断后结果不确定，记录失败并要求用户手动重发，不承诺 SMTP 严格只投递一次。接口定义见 [API2](API2.md) 和 [Legacy API](LEGACY_API.md)。

## 验证邮件模板

创建用户、管理员重发验证邮件以及账号页验证邮箱时，按操作请求的 `Accept-Language` 选择简体中文或英文。Web 自动传递当前界面语言；未指定或无支持语言时默认中文。邮件主题、正文、按钮和页脚在入队前确定，Asynq 不会按 worker 的系统语言重新选择。

验证邮件使用程序内置的珠玑笔记／Gemsnote 模板，不再读取旧数据库中的注册、改邮箱模板及其页眉页脚，因此升级已有 Leanote 实例无需手工修改模板数据。顶部应用名链接到配置的站点 URL；问候语后换行，深绿色验证按钮下方显示小字号完整备用链接。页脚为“珠玑笔记 · 日积字句，终得珠玑。”或“Gemsnote · Small notes become lasting insights.”，同时提供对应语言的纯文本邮件部分。

新用户主题为“恭喜您已注册珠玑笔记，请验证邮箱。”或“Welcome to Gemsnote — please verify your email.”。注册／当前邮箱验证链接有效期为 48 小时；只有尚未设置密码的新用户才显示验证后设置密码的说明。修改邮箱的验证链接仍为 2 小时，不显示初始密码设置说明。

## worker 与账号回归测试

使用临时 Redis、真实 Asynq worker、本地 TLS SMTP fixture 验证成功、拒收、认证失败、任务归属和超时，不连接真实邮箱：

```bash
python3 -m venv /tmp/gemsnote-mail-test-venv
/tmp/gemsnote-mail-test-venv/bin/pip install -r gemsnote-worker/tests/requirements.txt
MAIL_TEST_VENV=/tmp/gemsnote-mail-test-venv bash scripts/test-mail-queue.sh
npm --prefix frontend run build
node frontend/tests/account-email.mjs
```

Linux 下可复用相同用例验证共享镜像的 worker 入口（使用临时 Redis 和 SMTP，不访问实际邮箱）：

```bash
docker build -t gemsnote:mail-test .
MAIL_TEST_IMAGE=gemsnote:mail-test MAIL_TEST_VENV=/tmp/gemsnote-mail-test-venv bash scripts/test-mail-queue.sh
```

测试需要 Docker、OpenSSL、Go、Node.js 和 Chromium/Google Chrome；SMTP fixture 及 Redis 占用回环端口 10258 / 10259（SMTP）和 16389（Redis），脚本结束后停止临时容器和进程。

实现参考：[Asynq 官方仓库](https://github.com/hibiken/asynq)。运行时无需 Python；测试中的 Python 仅用于本地 SMTP fixture。

`scripts/test-email-verification.py` 用于完整账号邀请、改邮箱、重发、验证、密码设置及中英文内容回归。它仅允许连接一次性 `127.0.0.1:19090` 实例，要求 `GEMSNOTE_DISPOSABLE_API2=1` 和本地 SMTP fixture 的 `MAIL_CAPTURE` 文件；需先准备该实例、临时 Redis、worker 与 SMTP。该脚本会修改测试账号及 SMTP 设置，不应连接实际部署。
