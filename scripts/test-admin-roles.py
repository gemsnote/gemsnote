"""Destructive role regression: run only on the disposable API2 seed fixture."""
import concurrent.futures
import json
import os
import urllib.error
import urllib.parse
import urllib.request

if os.environ.get('GEMSNOTE_DISPOSABLE_API2') != '1':
    raise SystemExit('Requires GEMSNOTE_DISPOSABLE_API2=1 and a disposable server at 127.0.0.1:19090')
BASE = 'http://127.0.0.1:19090'

def request(path, data=None, token=None, json_body=True):
    headers = {'X-Requested-With': 'XMLHttpRequest'}
    if token:
        headers['Authorization'] = 'Bearer ' + token
    body = None
    if data is not None:
        headers['Content-Type'] = 'application/json' if json_body else 'application/x-www-form-urlencoded'
        body = (json.dumps(data) if json_body else urllib.parse.urlencode(data)).encode()
    try:
        response = urllib.request.urlopen(urllib.request.Request(BASE + path, body, headers), timeout=15)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        return response.status, json.load(response)

def login(username, password):
    status, result = request('/api2/auth/login', {'email': username, 'pwd': password})
    assert status == 200 and result.get('Token'), result
    return result['Token']

def bootstrap(token):
    return request('/api2/bootstrap', token=token)[1]

def change(token, user_id, role):
    return request('/api2/admin/role', {'userId': user_id, 'role': role}, token)

admin = login(os.environ.get('GEMSNOTE_TEST_ADMIN', 'admin'), 'gemsnote')
admin_id = bootstrap(admin)['User']['UserId']
assert change(admin, admin_id, 'user')[1]['Msg'] == 'cannotDemoteSelf'
assert bootstrap(admin)['IsAdmin']
assert change(admin, admin_id, 'superuser')[1]['Msg'] == 'invalidRole'
assert request('/api2/web/adminRegister', {'username': 'roles-test', 'email': 'roles-test@example.test', 'pwd': 'Role-test-123!', 'role': 'admin'}, admin, False)[1]['Ok']
# Explicit username is normalized, independently usable for login, and unique.
assert login('ROLES-TEST', 'Role-test-123!')
for invalid in ['', 'abc', 'bad name', 'bad@example.test']:
    result = request('/api2/web/adminRegister', {'username': invalid, 'email': 'invalid@example.test', 'pwd': 'Role-test-123!'}, admin, False)[1]
    assert not result['Ok'] and result['Msg'] == 'invalidUsername', result
result = request('/api2/web/adminRegister', {'username': 'Roles-Test', 'email': 'duplicate@example.test', 'pwd': 'Role-test-123!'}, admin, False)[1]
assert not result['Ok'] and result['Msg'] == 'usernameIsExisted', result
other = login('roles-test@example.test', 'Role-test-123!')
other_id = bootstrap(other)['User']['UserId']
assert not bootstrap(other)['IsAdmin'], 'new users must not inherit an injected admin role'
assert change(other, other_id, 'admin')[0] == 403
assert change(admin, other_id, 'admin')[1]['Ok']
assert bootstrap(other)['IsAdmin'], 'role elevation must apply to the existing token'
assert change(other, other_id, 'user')[1]['Msg'] == 'cannotDemoteSelf'
assert change(admin, other_id, 'user')[1]['Ok']
assert not bootstrap(other)['IsAdmin']
assert request('/api2/admin/data', {}, other)[0] == 403
assert change(admin, other_id, 'admin')[1]['Ok']
# Competing demotions must not remove both administrators.
with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
    results = list(pool.map(lambda args: change(*args), [(admin, other_id, 'user'), (other, admin_id, 'user')]))
assert sum(bool(result.get('Ok')) for _, result in results) == 1, results
assert bootstrap(admin)['IsAdmin'] != bootstrap(other)['IsAdmin']
if not bootstrap(other)['IsAdmin']:
    assert change(admin, other_id, 'admin')[1]['Ok']
    assert change(other, admin_id, 'user')[1]['Ok']
assert not bootstrap(admin)['IsAdmin'] and bootstrap(other)['IsAdmin']
assert request('/api2/admin/data', {}, admin)[0] == 403
print('PASS roles: no self-demotion, normal-user denial, grants/revocation, existing tokens and concurrent demotion')
