# 自动化构建与发布

## 持续集成

`.github/workflows/ci.yml` 会在每次 push 和 Pull Request 时执行：

1. 安装前端锁定依赖，运行 Vitest、类型检查和生产构建；
2. 下载 Go 模块依赖并运行 `go test ./...`；
3. 生成 Revel 服务端入口并编译；
4. 构建 Docker 镜像。

数据库集成测试仍需显式提供临时 MongoDB/PostgreSQL，普通 CI 中未配置数据库的集成测试会按测试代码约定自动跳过。

## 正式发布

正式版本标签必须采用 `vMAJOR.MINOR.PATCH` 格式，并与 `app/version/version.go` 中的 `Current` 完全一致。例如当前版本的发布标签为：

```bash
git tag -a v1.0.0 -m "Gemsnote 1.0.0"
git push origin v1.0.0
```

标签推送后，`.github/workflows/release.yml` 会：

1. 校验标签格式和代码版本；
2. 再次执行前端构建测试和全部 Go 测试；
3. 构建 Linux amd64、Linux arm64、Windows amd64、macOS amd64 和 macOS arm64 发布包；
4. 生成 `checksums.txt`；
5. 根据提交记录生成 Release Notes 并创建 GitHub Release；
6. 同时向 GitHub Container Registry（GHCR）推送包含 `linux/amd64` 和 `linux/arm64` 的镜像，镜像标签为 `ghcr.io/gemsnote/gemsnote:<版本号>`，例如 `ghcr.io/gemsnote/gemsnote:1.0.0`。GitHub Release 与镜像发布都在版本校验和发布包构建通过后执行。

镜像发布使用工作流自带的 `GITHUB_TOKEN`，无需配置 Docker Hub 凭据。首次发布的 GHCR 包默认可能为私有；如需公开拉取，请在 GitHub 的包设置中将其可见性改为 Public。公开后可用 `docker pull ghcr.io/gemsnote/gemsnote:1.0.0` 拉取，Docker 会自动选择匹配主机架构的镜像。本流程只发布明确的版本标签，不更新 `latest`。

如果发布工作流本身修复后需要重试已有标签，不要删除或强推标签，也不要只点旧运行记录的“重新运行”（它仍使用原来的工作流版本）。将修复后的工作流合并到 GitHub 默认分支后，在 Actions → Release → Run workflow 中填写已有标签（如 `v1.0.0`）。手动运行会从该标签检出源码，重新测试、构建并发布；请先确认该标签尚未创建 GitHub Release，避免重复创建失败。

每个发布包都包含服务端、Vue `frontend/dist`、MongoDB/PostgreSQL 迁移工具、配置、页面资源、数据库 Schema、初始数据和文档；运行发布包不需要安装 Node.js。

## 发布新版本

发布前应同时完成以下事项：

1. 修改 `app/version/version.go` 中的版本；
2. 在 `app/db/migrations.go` 添加该版本需要的幂等数据库迁移；
3. 更新文档并执行 `go test ./...`；
4. 合并到正式发布分支；
5. 创建并推送对应标签。

如果标签与应用版本不一致，Release 工作流会直接失败，不会产生错误版本的附件。

## 权限

CI 工作流只有仓库内容只读权限。Release 工作流只在正式标签触发，使用 GitHub 自动提供的 `GITHUB_TOKEN` 创建 Release 和发布 GHCR 镜像，无需配置个人访问令牌。仅镜像发布任务具有 `packages: write` 权限。

## 手工构建 Release

不使用 GitHub Actions 时，可在仓库根目录执行：

```bash
npm ci --prefix frontend
npm test --prefix frontend
npm run build --prefix frontend
go test ./...
version=$(sed -n 's/^const Current = "\([^"]*\)"/\1/p' app/version/version.go)
chmod +x scripts/build-release.sh
scripts/build-release.sh "$version" linux amd64 "$PWD/release"
```

`scripts/build-release.sh` 的参数格式为：

```bash
scripts/build-release.sh <version> <goos> <goarch> <绝对输出目录>
```

可用目标包括 `linux amd64`、`linux arm64`、`windows amd64`、`darwin amd64` 和 `darwin arm64`。POSIX 目标生成 `.tar.gz`，Windows 目标生成 `.zip`。输出目录必须使用绝对路径，因为 Windows 打包阶段会切换到临时目录。脚本会构建 Vue、生成 Revel 生产入口，并编译服务端和 `gemsnote-migrate` 工具；脚本本身不执行测试，也不会校验版本参数与 `app/version/version.go` 中的 `Current` 是否一致。

发布包包含服务端、迁移工具、`frontend/dist`、配置、数据库 Schema、初始数据和文档，运行时不需要 Node.js。包内通过 `run.sh`（Linux/macOS）或 `run.bat`（Windows）启动。发布前应检查包内 `conf`，不要带入本地密码、`.env` 或运行时数据；`files/` 不会打包，`public/upload` 为空目录。

启动前编辑 `conf/app.conf`，配置数据库连接并修改 `app.secret`。全新 PostgreSQL 数据库按部署方式执行 Schema 和 seed；MongoDB 不会自动导入 BSON，需要先按 [部署、测试与交付](../DEPLOYMENT.md) 恢复或连接数据库。已有数据库不会因重新构建 Release 而重置 seed 或密码。

## 手工构建 Desktop

Desktop 复用仓库根目录的 Vue 前端：

```bash
bash desktop-app/build-frontend.sh
cd desktop-app
wails build
```

需要 Go、Node.js/npm、Wails CLI 及当前操作系统的 WebKit/GTK 编译依赖。Desktop 输出通常位于 `desktop-app/build/bin/`，应用图标来自 `desktop-app/build/appicon.png`。Desktop 不使用服务端 Release 脚本的 CGO 交叉编译方式，跨平台时建议在对应系统原生构建。

## 发布包校验和

```bash
cd release
sha256sum gemsnote-* > checksums.txt       # Linux
shasum -a 256 gemsnote-* > checksums.txt   # macOS
```

Windows PowerShell：

```powershell
Get-FileHash .\gemsnote-*.zip -Algorithm SHA256
```
