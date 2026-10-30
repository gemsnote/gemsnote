"""Run against the disposable seeded PostgreSQL API2 fixture only."""
import concurrent.futures
import http.cookiejar
import json
import os
import re
import subprocess
import urllib.error
import urllib.parse
import urllib.request

if os.environ.get('GEMSNOTE_DISPOSABLE_API2') != '1':
    raise SystemExit('Requires disposable API2 fixture at 127.0.0.1:19090')
container = os.environ.get('GEMSNOTE_TEST_PG_CONTAINER', '')
if not re.fullmatch(r'gemsnote-api2-test-[0-9]+-[0-9]+', container):
    raise SystemExit('Requires the isolated PostgreSQL fixture container')
BASE = 'http://127.0.0.1:19090'

def request(path, data=None, token=None, form=False, opener=None):
    headers = {'X-Requested-With': 'XMLHttpRequest', 'Accept-Language': 'en-US'}
    if token:
        headers['Authorization'] = 'Bearer ' + token
    body = None
    if data is not None:
        body = (urllib.parse.urlencode(data) if form else json.dumps(data)).encode()
        headers['Content-Type'] = 'application/x-www-form-urlencoded' if form else 'application/json'
    try:
        response = (opener.open if opener else urllib.request.urlopen)(urllib.request.Request(BASE + path, body, headers), timeout=20)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        return response.status, json.load(response)

def login(name, password='Flag-test-123!'):
    return request('/api2/auth/login', {'email': name, 'pwd': password})[1]

def bootstrap(token):
    return request('/api2/bootstrap', token=token)[1]

def change(token, user_id, flag):
    return request('/api2/admin/flag', {'userId': user_id, 'flag': flag}, token)

def sql(statement):
    return subprocess.check_output(['docker', 'exec', '-i', container, 'psql', '-U', 'gemsnote_test', '-d', 'gemsnote', '-At', '-v', 'ON_ERROR_STOP=1'], input=statement, text=True).strip()

admin = login('admin', 'gemsnote')['Token']
admin_id = bootstrap(admin)['User']['UserId']
assert change(admin, admin_id, 'D')[1]['Msg'] == 'cannotDisableSelf'
assert change(admin, admin_id.upper(), 'D')[1]['Msg'] == 'cannotDisableSelf'
assert change(admin, admin_id, 'invalid')[1]['Msg'] == 'invalidUserFlag'
name = 'flags-test@example.test'
assert request('/api2/web/adminRegister', {'username': 'flags-test', 'email': name, 'pwd': 'Flag-test-123!', 'flag': 'D'}, admin, True)[1]['Ok']
created_id=next(u['UserId'] for u in request('/api2/admin/data',{},admin)[1]['Users'] if u['Username']=='flags-test')
assert request('/api2/web/adminResetPwd',{'userId':created_id,'pwd':'Flag-test-123!'},admin,True)[1]['Ok']
other = login(name)['Token']
user_id = bootstrap(other)['User']['UserId']
assert re.fullmatch('[a-f0-9]{24}', user_id)
users = request('/api2/admin/data', {}, admin)[1]['Users']
assert next(user for user in users if user['UserId'] == user_id)['Flag'] == 'A'
assert change(other, admin_id, 'D')[0] == 403
assert change(None, user_id, 'D')[0] == 401
legacy = request('/api/auth/login?' + urllib.parse.urlencode({'email': name, 'pwd': 'Flag-test-123!'}))[1]['Token']
cookie = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
assert request('/doLogin', {'email': name, 'pwd': 'Flag-test-123!'}, form=True, opener=cookie)[1]['Ok']
assert request('/web/bootstrap', opener=cookie)[1]['User']
notebooks = sql(f"SELECT count(*) FROM notebooks WHERE user_id='{user_id}';")
sql(f"INSERT INTO tokens(id,email,token,type) VALUES ('{user_id}','disabled-target@example.test','flag-fixture-verification',2);")
assert change(admin, user_id, 'D')[1]['Ok']
assert not bootstrap(other)['User']
assert request('/api2/user/info', token=other)[0] == 401
assert not request('/api/user/info?token=' + urllib.parse.quote(legacy))[1]['Ok']
assert not request('/web/bootstrap', opener=cookie)[1]['User']
assert login(name)['Msg'] == 'userDisabled'
assert login(name, 'wrong-password')['Msg'] != 'userDisabled'
assert not request('/api/auth/login', {'email': name, 'pwd': 'Flag-test-123!'}, form=True)[1]['Ok']
assert not request('/doLogin', {'email': name, 'pwd': 'Flag-test-123!'}, form=True, opener=cookie)[1]['Ok']
assert sql(f"SELECT count(*) FROM sessions WHERE user_id='{user_id}';") == '0'
assert sql(f"SELECT count(*) FROM tokens WHERE id='{user_id}';") == '0'
assert sql(f"SELECT count(*) FROM notebooks WHERE user_id='{user_id}';") == notebooks
# Even a lingering row (e.g. direct SQL disabling) cannot authorize a disabled user.
sql(f"INSERT INTO sessions(id,session_id,user_id) VALUES ('0123456789abcdef01234567','flag-fixture-orphan','{user_id}');")
assert request('/api2/user/info', token='flag-fixture-orphan')[0] == 401
sql(f"INSERT INTO tokens(id,email,token,type) VALUES ('{user_id}','disabled-target@example.test','flag-fixture-orphan-link',2);")
assert not request('/api2/web/verifyEmail', {'token':'flag-fixture-orphan-link','change':'true'}, form=True)[1]['Ok']
assert change(admin, user_id, 'A')[1]['Ok']
assert not request('/api2/web/verifyEmail', {'token':'flag-fixture-orphan-link','change':'true'}, form=True)[1]['Ok']
assert not request('/api2/web/verifyEmail', {'token':'flag-fixture-verification','change':'true'}, form=True)[1]['Ok']
for token in (other, legacy, 'flag-fixture-orphan'):
    assert not bootstrap(token)['User']
assert not request('/web/bootstrap', opener=cookie)[1]['User']
other = login(name)['Token']
# Racing login and disabling must never leave a newly usable token after disable.
with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
    futures = [pool.submit(login, name) for _ in range(8)]
    disabled = pool.submit(change, admin, user_id, 'D')
    results = [future.result() for future in futures]
    assert disabled.result()[1]['Ok']
for result in results:
    if result.get('Token'):
        assert not bootstrap(result['Token'])['User']
assert sql(f"SELECT count(*) FROM sessions WHERE user_id='{user_id}';") == '0'
assert change(admin, user_id, 'A')[1]['Ok']
other = login(name)['Token']
assert request('/api2/admin/role', {'userId': user_id, 'role': 'admin'}, admin)[1]['Ok']
assert change(other, user_id, 'D')[1]['Msg'] == 'cannotDisableSelf'
# Role demotion and disabling share a lock, preserving an available administrator.
with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
    first = pool.submit(change, admin, user_id, 'D')
    second = pool.submit(request, '/api2/admin/role', {'userId': admin_id, 'role': 'user'}, other)
    results = [first.result(), second.result()]
assert sum(bool(result.get('Ok')) for _, result in results) == 1, results
assert int(sql("SELECT count(*) FROM users WHERE role='admin' AND flag='A';")) >= 1
if not bootstrap(admin)['IsAdmin']:
    assert request('/api2/admin/role', {'userId': admin_id, 'role': 'admin'}, other)[1]['Ok']
    assert change(admin, user_id, 'D')[1]['Ok']
print('PASS flags: self-disable denied, legacy/API2/cookie login blocked, tokens revoked, re-enable requires login, races and data retention')
