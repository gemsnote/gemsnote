#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/.."
: "${MAIL_TEST_VENV:?Set MAIL_TEST_VENV to a virtualenv with mail-worker/tests/requirements.txt installed}"
"$MAIL_TEST_VENV/bin/python" - <<'PORTS'
import socket
for port in (16389, 10258, 10259):
    with socket.socket() as probe:
        probe.bind(('127.0.0.1', port))
PORTS
work=$(mktemp -d /tmp/gemsnote-mail-test.XXXXXX)
container=gemsnote-mail-test-$$
smtp_pid= worker_pid=
worker_container=
cleanup() {
  if [ -n "$worker_container" ]; then
    docker logs "$worker_container" >"$work/worker.log" 2>&1 || true
    docker stop -t 40 "$worker_container" >/dev/null 2>&1 || true
  fi
  kill ${smtp_pid:-} ${worker_pid:-} 2>/dev/null || true
  if [ -n "$worker_pid" ]; then wait "$worker_pid" 2>/dev/null || true; fi
  docker stop "$container" >/dev/null 2>&1 || true
  echo "Mail integration logs: $work"
}
trap cleanup EXIT
openssl req -x509 -newkey rsa:2048 -nodes -keyout "$work/key.pem" -out "$work/cert.pem" -days 1 -subj '/CN=localhost' -addext 'subjectAltName=DNS:localhost' >"$work/openssl.log" 2>&1
docker run --rm -d --name "$container" -p 127.0.0.1:16389:6379 redis:7-alpine >/dev/null
go build -o "$work/mail-worker" ./mail-worker
export MAIL_CAPTURE="$work/messages.jsonl"
export GEMSNOTE_MAIL_REDIS_URL=redis://127.0.0.1:16389/0 SSL_CERT_FILE="$work/cert.pem"
"$MAIL_TEST_VENV/bin/python" mail-worker/tests/smtp_fixture.py "$work/cert.pem" "$work/key.pem" "$work/messages.jsonl" >"$work/smtp.log" 2>&1 & smtp_pid=$!
if [ -n "${MAIL_TEST_IMAGE:-}" ]; then
    # Linux Docker host networking connects only to the isolated fixtures above.
    worker_container="$container-worker"
    docker run --rm -d --name "$worker_container" --network host \
        -e GEMSNOTE_MAIL_REDIS_URL -e SSL_CERT_FILE=/tmp/mail-test-cert.pem \
        -v "$work/cert.pem:/tmp/mail-test-cert.pem:ro" \
        --entrypoint /opt/gemsnote/docker-entry-mail-worker.sh \
        "$MAIL_TEST_IMAGE" >/dev/null
else
    "$work/mail-worker" >"$work/worker.log" 2>&1 & worker_pid=$!
fi
sleep 3
GEMSNOTE_MAIL_TEST_REDIS="$GEMSNOTE_MAIL_REDIS_URL" go test ./app/service -run 'TestAsynq' -v -count=1
GEMSNOTE_MAIL_TEST_REDIS="$GEMSNOTE_MAIL_REDIS_URL" go test ./internal/mailqueue -run 'TestRedis' -v -count=1
