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
cleanup(){ kill ${smtp_pid:-} ${worker_pid:-} 2>/dev/null || true; docker stop "$container" >/dev/null 2>&1 || true; }
trap cleanup EXIT
openssl req -x509 -newkey rsa:2048 -nodes -keyout "$work/key.pem" -out "$work/cert.pem" -days 1 -subj '/CN=localhost' -addext 'subjectAltName=DNS:localhost' >"$work/openssl.log" 2>&1
docker run --rm -d --name "$container" -p 127.0.0.1:16389:6379 redis:7-alpine >/dev/null
export MAIL_CAPTURE="$work/messages.jsonl"
export GEMSNOTE_MAIL_REDIS_URL=redis://127.0.0.1:16389/0 SSL_CERT_FILE="$work/cert.pem"
"$MAIL_TEST_VENV/bin/python" mail-worker/tests/smtp_fixture.py "$work/cert.pem" "$work/key.pem" "$work/messages.jsonl" >"$work/smtp.log" 2>&1 & smtp_pid=$!
(cd mail-worker; exec "$MAIL_TEST_VENV/bin/celery" -A tasks worker --pool=solo --loglevel=INFO --queues=gemsnote.mail) >"$work/worker.log" 2>&1 & worker_pid=$!
sleep 3
GEMSNOTE_MAIL_TEST_REDIS="$GEMSNOTE_MAIL_REDIS_URL" go test ./app/service -run 'TestCelery' -v -count=1
echo "Mail integration logs: $work"
