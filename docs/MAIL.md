# 异步邮件与邮箱修改

邮件由 Asynq worker 发送，Go 服务端只负责验证、渲染模板和入队。Redis 是邮件任务队列，业务数据仍保存在 PostgreSQL 或 MongoDB。不启用邮件服务不会影响笔记功能，但修改邮箱、邮箱验证和找回密码无法发送邮件。

## Docker 部署

三个 Compose 文件都提供可选的 `mail` profile。使用 GHCR 镜像：

```bash
docker compose --profile mail up -d
```

这会增加 `redis` 和 `mail-worker` 两个服务。`gemsnote` 和 `mail-worker` 使用同一个 `ghcr.io/gemsnote/gemsnote:${GEMSNOTE_VERSION}` 镜像，分别运行在独立容器中。默认入口 `/opt/gemsnote/docker-entry-server.sh` 启动 Go 服务端；worker 通过 Compose 覆盖入口为 `/opt/gemsnote/docker-entry-mail-worker.sh`，以非 root 用户启动 Asynq。两个容器都通过 `redis://redis:6379/0` 连接邮件队列。新版本发布前，本地测试使用源码 Compose：

```bash
docker compose -f docker-compose.postgres.yml --profile mail up -d --build
```

源码 MongoDB Compose 使用同样的 `--profile mail` 参数。Redis 不发布宿主机端口，任务队列使用独立的 `mail-redis-data` 卷持久化。排查发送失败：

```bash
docker compose --profile mail logs --tail=100 mail-worker redis
```

源码部署时给这些命令加对应的 `-f` 参数。仅重建 Go 服务端镜像不会启动 Asynq，必须启用 `mail` profile。

### 从 Celery 版本升级

Asynq 与 Celery 的任务协议不兼容。先暂停创建账号、发送验证邮件等入队操作，保持旧 worker 运行直至旧任务处理完成，再同时升级服务端和 worker；不要混用新旧版本。无法处理的旧任务不会自动转入新队列，升级后通过页面重新发送。

切换时停止旧 Compose 项目，避免旧 `mail-redis` 与新 `redis` 同时使用同一个数据卷。以下以源码 PostgreSQL 部署为例：

```bash
docker compose -f docker-compose.postgres.yml --profile mail down --remove-orphans
docker compose -f docker-compose.postgres.yml --profile mail up -d --build
```

GHCR 部署使用 `docker compose --profile mail down --remove-orphans`，随后 `docker compose --profile mail pull` 和 `docker compose --profile mail up -d`。保持原来的项目名和部署目录，不要使用 `down -v`；服务名改为 `redis`，已有 `mail-redis-data` 卷名保留，避免丢失队列。此操作会短暂停止服务。

## Release 或源码服务端

Release 包附带原生可执行文件 `bin/gemsnote-mail-worker`（Windows 为 `.exe`），无需 Python 环境。准备 Redis，在 Release 解压目录运行：

```bash
export GEMSNOTE_MAIL_REDIS_URL='redis://127.0.0.1:6379/0'
./bin/gemsnote-mail-worker -concurrency 2
```

Windows PowerShell：

```powershell
$env:GEMSNOTE_MAIL_REDIS_URL = 'redis://127.0.0.1:6379/0'
.\bin\gemsnote-mail-worker.exe -concurrency 2
```

源码部署可在项目根目录运行 `go build -o /tmp/gemsnote-mail-worker ./mail-worker`，设置同一环境变量后执行该文件。`-concurrency` 默认为 2；worker 只需 Redis 和 SMTP，不连接业务数据库，也不读取 `app.conf`。

将 worker 配置为系统服务长期运行。Go 服务端也必须配置同一个 `GEMSNOTE_MAIL_REDIS_URL`，或在 `conf/app.conf` 中设置：

```ini
mail.redis_url=redis://127.0.0.1:6379/0
```

环境变量优先于配置文件。Redis 需专用于邮件队列，或使用独立数据库编号；可使用 `redis://:密码@主机:端口/数据库编号`。不要将 Redis 暴露到公网；其中短暂保存邮件正文、验证链接和 SMTP 凭据。敏感任务内容最多保存 10 分钟，发送完成后删除；发送状态保留 24 小时。Windows 可运行原生 Go worker，Redis 可部署在 Docker 或 Linux 主机。

## 配置 SMTP

管理员在 Web 笔记页面左下角点击头像 → **管理** → **站点与邮件**，填写站点 URL、SMTP 主机、端口、用户名、密码和 SSL。

- 站点 URL 应是用户实际可以访问的 HTTPS 地址，用于生成验证链接。
- SMTP 用户名同时作为发件地址，必须是有效邮箱。
- SSL 开启表示隐式 TLS（通常为 465）；关闭表示使用 STARTTLS（通常为 587）。不允许明文传输 SMTP 密码。
- 自建 CA 可为 worker 配置 `SSL_CERT_FILE`；不要关闭证书验证。
- 账号页发送验证邮件可以检查完整发送链路。入队成功不等于送达收件箱；`sent` 表示 SMTP 服务器已接受邮件，仍可能进入垃圾邮件或随后退信。

## 用户修改邮箱

头像 → **账号**，填写新邮箱和当前密码并发送验证邮件。页面区分排队、发送中、已交给 SMTP、发送失败，并显示认证失败、收件人拒收、连接失败等原因。返回账号页仍可查询该任务；任务允许本人查询，管理员也可查询所管理账号的任务。

验证成功前，账号邮箱保持不变。用户必须打开新邮箱收到的验证链接，完成验证后才会修改邮箱，链接只允许使用一次。每个账号 60 秒内只能请求一次验证邮件，防止反复点击使已发送链接失效。离线或查询失败时可点击“重新查询发送结果”。排队超过 10 分钟会提示任务超时，请管理员检查 worker；整个 SMTP 会话限时 30 秒。SMTP 拒收、认证失败和连接失败会记录为最终失败，不自动重发邮件；Redis 基础设施错误最多重试 3 次。worker 中断后若无法确认 SMTP 是否已经发送，将标记发送失败供用户手动重发，避免自动重复投递。

“重新发送当前邮箱验证邮件”只验证现有地址，不修改邮箱。旧邮件接口的成功结果表示任务已接受；账号页使用新增任务查询接口显示最终发送结果。管理员批量邮件日志记录入队结果，最终失败可在 worker 日志按任务 ID 排查。

## 验证邮件内容与语言

创建用户、管理员重发验证邮件以及账号页验证邮箱时，按操作请求的 `Accept-Language` 选择简体中文或英文。Web 自动传递当前界面语言；未指定或无支持语言时默认中文。邮件主题、正文、按钮和页脚在入队前确定，Asynq 不会按 worker 的系统语言重新选择。

验证邮件使用程序内置的珠玑笔记／Gemsnote 模板，不再读取旧数据库中的注册、改邮箱模板及其页眉页脚，因此升级已有 Leanote 实例无需手工修改模板数据。顶部应用名链接到配置的站点 URL；问候语后换行，深绿色验证按钮下方显示小字号完整备用链接。页脚为“珠玑笔记 · 日积字句，终得珠玑。”或“Gemsnote · Small notes become lasting insights.”，同时提供对应语言的纯文本邮件部分。

新用户主题为“恭喜您已注册珠玑笔记，请验证邮箱。”或“Welcome to Gemsnote — please verify your email.”。注册／当前邮箱验证链接有效期为 48 小时；只有尚未设置密码的新用户才显示验证后设置密码的说明。修改邮箱的验证链接仍为 2 小时，不显示初始密码设置说明。

## 新用户验证邮件

用户管理列表中的“邮箱未验证”可点击，确认收件地址后重发验证邮件。页面显示发送进度、成功或失败原因；发送中及已禁用账号不可点击，重复请求受每账号 60 秒间隔限制。发送成功表示邮件服务器已接受邮件。新用户点击验证链接后设置自己的密码，已有密码的用户仅验证邮箱。

## 管理员角色

Web 的“头像 → 管理 → 用户管理”列出每位用户的角色，可将其他用户设为“管理员”或“普通用户”。新建用户默认为普通用户。用户管理还可设置账号状态：`A`（可用）或 `D`（禁用）。管理员不能禁用自己；禁用后登录、已有会话及验证链接失效，重新启用需重新登录，笔记数据保留。

管理员不能将自己改为普通用户：界面禁用自己的角色下拉框，服务端也强制拒绝此请求。其他管理员可以撤销该用户的权限，下一次管理请求立即按最新角色检查。桌面端不提供管理界面。

升级到 1.2.0 时自动增加用户角色字段。首次初始化沿用已有 `adminUserId`，或按 `adminUsername` 定位原管理员。若旧管理员已改名且未保存 ID，只有在所有旧全局配置均指向同一个有效拥有者、该账号没有明确角色且当前没有管理员时，才按该拥有者 ID 恢复管理员；配置拥有者缺失或不一致时不会猜测。只有没有角色的旧账号才会被初始化为管理员。此后以用户记录中的角色为准，改用户名或邮箱不改变角色，重启也不会重新提升已被降为普通用户的账号。

若旧配置无法唯一识别原管理员（拥有者缺失、不一致或账号不存在），启动日志会报告管理员初始化失败。此时需在 `conf/app.conf` 中设置正确的 `adminUsername` 或 `adminUserId=<24 位用户 ID>`，用于初始化原管理员。配置不再是长期权限白名单，也不能覆盖已有明确角色。

## 回归测试

使用临时 Redis、真实 Asynq worker、本地 TLS SMTP fixture 验证成功、拒收、认证失败、任务归属和超时，不连接真实邮箱：

```bash
python3 -m venv /tmp/gemsnote-mail-test-venv
/tmp/gemsnote-mail-test-venv/bin/pip install -r mail-worker/tests/requirements.txt
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
