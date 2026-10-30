#!/usr/bin/env python3
"""Destructive regression fixture: use ONLY a disposable local API2 server/PG."""
import base64
import http.cookiejar
import io
import json
import os
import secrets
import subprocess
import tarfile
import urllib.error
import urllib.parse
import urllib.request

BASE = os.environ.get("GEMSNOTE_API2_URL", "http://127.0.0.1:19090").rstrip("/")
PG = os.environ.get("GEMSNOTE_TEST_PG_CONTAINER", "")
assert os.environ.get("GEMSNOTE_DISPOSABLE_API2") == "1", "disposable test instance required"
assert urllib.parse.urlsplit(BASE).hostname in ("127.0.0.1", "localhost"), "local fixtures only"
assert PG.startswith("gemsnote-api2-"), "disposable PostgreSQL fixture required"
client = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))


def request(path, data=None, token=None, form=False, headers=None, opener=client):
    hdr = dict(headers or {})
    if token:
        hdr["Authorization"] = "Bearer " + token
    if isinstance(data, dict):
        hdr["Content-Type"] = "application/x-www-form-urlencoded" if form else "application/json"
        data = (urllib.parse.urlencode(data) if form else json.dumps(data)).encode()
    req = urllib.request.Request(BASE + path, data=data, headers=hdr)
    try:
        response = opener.open(req, timeout=15)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        return response.code, response.read(), response.headers


def api(path, data=None, token=None, **kwargs):
    status, raw, _ = request(path, data, token, **kwargs)
    assert status == 200, (path, status)
    return json.loads(raw)


def upload(path, content, token, **fields):
    boundary = "gemsnote-" + secrets.token_hex(12)
    parts = []
    for name, value in fields.items():
        parts.append(f'--{boundary}\r\nContent-Disposition: form-data; name="{name}"\r\n\r\n{value}\r\n'.encode())
    parts += [f'--{boundary}\r\nContent-Disposition: form-data; name="file"; filename="test.png"\r\nContent-Type: image/png\r\n\r\n'.encode(), content, f'\r\n--{boundary}--\r\n'.encode()]
    return api(path, b"".join(parts), token, headers={"Content-Type": "multipart/form-data; boundary=" + boundary})


def sql(statement):
    return subprocess.check_output(["docker", "exec", PG, "psql", "-U", "gemsnote_test", "-d", "gemsnote", "-At", "-v", "ON_ERROR_STOP=1", "-c", statement], text=True).strip()


admin = api("/api2/auth/login", {"email": "admin", "pwd": "gemsnote"})
admin_token = admin["Token"]
owner = admin["User"]["UserId"]
email = "security-" + secrets.token_hex(6) + "@example.test"
pwd = secrets.token_urlsafe(18)
assert api("/api2/web/adminRegister", {"username": "security-test", "email": email, "pwd": pwd}, admin_token, form=True)["Ok"]
created_id=next(u["UserId"] for u in api("/api2/admin/data",{},admin_token)["Users"] if u["Email"]==email)
assert api("/api2/web/adminResetPwd",{"userId":created_id,"pwd":pwd},admin_token,form=True)["Ok"]
other = api("/api2/auth/login", {"email": email, "pwd": pwd})
token = other["Token"]
user = other["User"]["UserId"]
assert request("/api2/admin/data", {}, token)[0] == 403

# A valid browser admin cookie must not override a different account's token.
assert api("/doLogin", {"email": "admin", "pwd": "gemsnote"}, form=True)["Ok"]
identity = api("/api/user/info?" + urllib.parse.urlencode({"token": token}))
assert identity["UserId"] == user, "legacy token/cookie identity confusion"
assert request("/api/user/updateUsername?username=csrf-changed")[0] == 405
assert request("/web/adminSettings?siteUrl=http://changed.invalid")[0] == 405
assert request("/adminsetting/doSiteUrl?siteUrl=http://changed.invalid")[0] == 405
invalid_token = secrets.token_hex(32)
assert not api("/api/user/info?token=" + invalid_token)["Ok"]
assert sql(f"SELECT count(*) FROM sessions WHERE session_id='{invalid_token}'") == "0"
assert request("/api2/Save", {})[0] in (401, 404), "wildcard API2 alias bypass"

boot = api("/api2/bootstrap", token=admin_token)
book = boot["Notebooks"][0]["NotebookId"]
note = secrets.token_hex(12)
doc = api("/api2/save", {"noteId": note, "notebookId": book, "title": "../../security-overwrite", "content": "private-original", "isNew": True, "isMarkdown": True}, admin_token)
assert doc["Content"] == "private-original"
assert request("/api2/document", {"noteId": note}, token)[0] == 403
assert not api("/api/note/getNote?" + urllib.parse.urlencode({"noteId": note, "token": token})).get("NoteId")

png = base64.b64decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jRZkAAAAASUVORK5CYII=")
image = upload("/api2/file/pasteImage", png, admin_token, noteId=note)
assert image["Ok"]
image_id = image["Id"]
assert request("/api2/file/getImage?fileId=" + image_id, token=admin_token)[1] == png
assert request("/api2/file/getImage?fileId=" + image_id, token=token)[1] != png
assert not upload("/api2/avatar", b"<html><script>alert(1)</script></html>", token)["Ok"]
assert not upload("/api2/file/pasteImage", b'<svg onload="alert(1)"/>', admin_token, noteId=note)["Ok"]

attach = upload("/api2/attachments/upload", b"attachment bytes", admin_token, noteId=note)
assert attach["Ok"]
attachment_id = attach["Id"]
assert request("/api2/attach/download?attachId=" + attachment_id, token=admin_token)[1] == b"attachment bytes"
assert request("/api2/attach/download?attachId=" + attachment_id, token=token)[1] != b"attachment bytes"
status, archive, headers = request("/api2/attach/downloadAll?noteId=" + note, token=admin_token)
assert status == 200 and "../" not in headers.get("Content-Disposition", "")
with tarfile.open(fileobj=io.BytesIO(archive), mode="r:gz") as tar:
    assert tar.getnames() == ["test.png"]
    assert tar.extractfile("test.png").read() == b"attachment bytes"
status, archive, _ = request("/api/file/getAllAttachs?" + urllib.parse.urlencode({"noteId": note, "token": admin_token}))
assert status == 200
with tarfile.open(fileobj=io.BytesIO(archive), mode="r:gz") as tar:
    assert tar.extractfile("test.png").read() == b"attachment bytes"

# The public file aliases must resolve Bearer identity and never use stale cookies.
assert request("/api2/file/getImage?fileId=" + image_id)[1] != png
assert request("/api2/attach/download?attachId=" + attachment_id)[1] != b"attachment bytes"
api("/api2/auth/logout", {}, admin_token)
assert request("/api2/file/getImage?fileId=" + image_id, token=admin_token)[1] != png

assert api("/doLogin", {"email": email, "pwd": pwd}, form=True)["Ok"]
user_image = upload("/api2/file/pasteImage", png, token)
assert user_image["Ok"]
assert request("/file/outputImage?fileId="+user_image["Id"])[1] == png
copied = api("/file/copyImage", {"userId": owner, "fileId": image_id, "toUserId": user}, form=True, headers={"X-Requested-With": "XMLHttpRequest"})
assert not copied["Ok"]
assert not api("/file/copyHttpImage", {"src": BASE + "/api2/system/version"}, form=True, headers={"X-Requested-With": "XMLHttpRequest"})["Ok"]

# Activation tokens must not reset passwords; valid recovery tokens still work.
recovery = secrets.token_hex(32)
sql(f"INSERT INTO tokens(id,email,token,type) VALUES('{user}','{email}','{recovery}',1) ON CONFLICT(id) DO UPDATE SET token=EXCLUDED.token,type=1,created_time=NOW()")
assert not api("/api2/auth/password/reset", {"token": recovery, "pwd": pwd + "new"})["Ok"]
sql(f"UPDATE tokens SET type=0 WHERE id='{user}'")
assert api("/api2/auth/password/reset", {"token": recovery, "pwd": pwd + "new"})["Ok"]
assert request("/api2/user/info", token=token)[0] == 401
assert request("/file/outputImage?fileId="+user_image["Id"])[1] != png, "revoked cookie still reads files"
assert api("/api2/auth/login", {"email": email, "pwd": pwd + "new"})["Ok"]
print("API/API2 security regression passed: identities, ownership, CSRF, uploads, archives, SSRF and reset-token scope")
