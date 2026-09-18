#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
exec "$ROOT/bin/gemsnote" \
  -importPath github.com/gemsnote/gemsnote \
  -srcPath "$ROOT/runtime" \
  -runMode prod "$@"
