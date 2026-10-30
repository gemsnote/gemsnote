#!/bin/sh
set -eu
: "${GEMSNOTE_MAIL_REDIS_URL:?Set GEMSNOTE_MAIL_REDIS_URL to the mail Redis database}"
cd /opt/gemsnote
if [ "$(id -u)" -eq 0 ]; then
    exec su-exec worker:worker /opt/gemsnote/gemsnote-worker "$@"
fi
exec /opt/gemsnote/gemsnote-worker "$@"
