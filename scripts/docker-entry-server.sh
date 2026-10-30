#!/bin/sh
set -eu
cd /opt/gemsnote
exec /opt/gemsnote/gemsnote -importPath github.com/gemsnote/gemsnote \
    -srcPath /opt/gemsnote/runtime -runMode prod "$@"
