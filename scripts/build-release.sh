#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage: $0 <goos> <goarch> <absolute-output-dir>" >&2
  echo "Version is read from app/version/version.go (Current); no version argument is accepted." >&2
}
if [[ $# -eq 1 && ( "$1" == "--help" || "$1" == "-h" ) ]]; then
  usage
  exit 0
fi
if [[ $# -ne 3 ]]; then
  usage
  exit 2
fi

target_os=$1
target_arch=$2
output_dir=$3
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
version_file="$project_root/app/version/version.go"
[[ -f "$version_file" ]] || { echo "Version file not found: $version_file" >&2; exit 1; }
version=$(sed -n 's/^const Current = "\([^"]*\)"/\1/p' "$version_file")
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "Invalid or missing Current in $version_file (expected MAJOR.MINOR.PATCH)" >&2; exit 1; }
case "$target_os/$target_arch" in
  linux/amd64|linux/arm64|darwin/amd64|darwin/arm64|windows/amd64) ;;
  *) echo "Unsupported release target: $target_os/$target_arch" >&2; exit 2 ;;
esac
[[ "$output_dir" == /* ]] || { echo "Output directory must be an absolute path: $output_dir" >&2; exit 2; }
package_name="gemsnote-${target_os}-${target_arch}-v${version}"
stage_dir=$(mktemp -d)
revel_stage=$(mktemp -d)
cleanup() {
  chmod -R u+w "$stage_dir" "$revel_stage" 2>/dev/null || true
  rm -rf "$stage_dir" "$revel_stage"
}
trap cleanup EXIT

cd "$project_root"
if [[ ! -d frontend/node_modules ]]; then
  npm ci --prefix frontend
fi
npm run build --prefix frontend
go run ./app/cmd build . "$revel_stage" prod

binary_name=gemsnote
migration_name=gemsnote-migrate
if [[ $target_os == windows ]]; then
  binary_name=gemsnote.exe
  migration_name=gemsnote-migrate.exe
fi

mkdir -p "$stage_dir/gemsnote/bin" \
  "$stage_dir/gemsnote/runtime/github.com/revel/revel" \
  "$stage_dir/gemsnote/runtime/github.com/gemsnote"

CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
  go build -trimpath -ldflags="-s -w" \
  -o "$stage_dir/gemsnote/bin/$binary_name" ./app/tmp
CGO_ENABLED=0 GOOS="$target_os" GOARCH="$target_arch" \
  go build -trimpath -ldflags="-s -w" \
  -o "$stage_dir/gemsnote/bin/$migration_name" ./tools/migration

cp -R conf messages public app/views mongodb_backup database docs \
  "$stage_dir/gemsnote/"
mkdir -p "$stage_dir/gemsnote/frontend"
cp -R frontend/dist "$stage_dir/gemsnote/frontend/"
cp README.md "$stage_dir/gemsnote/"
rm -rf "$stage_dir/gemsnote/public/upload"
mkdir -p "$stage_dir/gemsnote/public/upload"

revel_dir=$(go list -m -f '{{.Dir}}' github.com/revel/revel)
cp -R "$revel_dir/conf" "$revel_dir/templates" \
  "$stage_dir/gemsnote/runtime/github.com/revel/revel/"

if [[ $target_os == windows ]]; then
  cp scripts/release/run.bat "$stage_dir/gemsnote/run.bat"
else
  ln -s ../../.. "$stage_dir/gemsnote/runtime/github.com/gemsnote/gemsnote"
  cp scripts/release/run.sh "$stage_dir/gemsnote/run.sh"
  chmod +x "$stage_dir/gemsnote/run.sh"
fi

mkdir -p "$output_dir"
if [[ $target_os == windows ]]; then
  (
    cd "$stage_dir"
    zip -qr "$output_dir/$package_name.zip" gemsnote
  )
else
  tar -C "$stage_dir" -czf "$output_dir/$package_name.tar.gz" gemsnote
fi
