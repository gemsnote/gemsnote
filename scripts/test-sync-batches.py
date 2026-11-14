#!/usr/bin/env python3
"""Exercise API2 metadata/body batches against a disposable PostgreSQL server."""
import gzip
from concurrent.futures import ThreadPoolExecutor
import json
import os
import subprocess
import urllib.error
import urllib.parse
import urllib.request

base = os.environ.get('GEMSNOTE_API2_URL', 'http://127.0.0.1:19440').rstrip('/')
container = os.environ.get('GEMSNOTE_TEST_PG_CONTAINER', '')
assert os.environ.get('GEMSNOTE_DISPOSABLE_API2') == '1'
assert urllib.parse.urlsplit(base).hostname in ('127.0.0.1', 'localhost')
assert container.startswith('gemsnote-api2-')
opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))


def request(path, value=None, token='', compress=False):
    headers = {'Content-Type': 'application/json'}
    if token:
        headers['Authorization'] = 'Bearer ' + token
    if compress:
        headers['Accept-Encoding'] = 'gzip'
    req = urllib.request.Request(base + path, data=None if value is None else json.dumps(value).encode(), headers=headers)
    try:
        response = opener.open(req, timeout=30)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        body = response.read()
        if response.headers.get('Content-Encoding') == 'gzip':
            body = gzip.decompress(body)
        return response.status, json.loads(body)


def sql(statement):
    return subprocess.check_output(['docker', 'exec', '-i', container, 'psql', '-U', 'gemsnote_test', '-d', 'gemsnote', '-At', '-v', 'ON_ERROR_STOP=1'], input=statement, text=True).strip()


_, login = request('/api2/auth/login', {'email': 'admin', 'pwd': 'gemsnote'})
assert login['Ok'], login
token, user = login['Token'], login['User']['UserId']
book = sql(f"SELECT id FROM notebooks WHERE user_id='{user}' AND NOT is_deleted LIMIT 1")
assert book
# All writes are confined to this disposable container, never a deployment DB.
sql(f"""
DELETE FROM note_contents WHERE id IN (SELECT lpad(to_hex(700000+i),24,'0') FROM generate_series(1,501) AS i);
DELETE FROM notes WHERE id IN (SELECT lpad(to_hex(700000+i),24,'0') FROM generate_series(1,501) AS i);
INSERT INTO notes(id,user_id,notebook_id,title,usn)
SELECT lpad(to_hex(700000+i),24,'0'),'{user}','{book}','batch-'||i,1000000+i FROM generate_series(1,501) AS i;
INSERT INTO note_contents(id,user_id,content)
SELECT lpad(to_hex(700000+i),24,'0'),'{user}',CASE WHEN i=1 THEN '' ELSE 'body-'||i END FROM generate_series(1,501) AS i;
UPDATE users SET usn=1000501 WHERE id='{user}';
""")
ids = [f'{700000+i:024x}' for i in range(1, 502)]
status, page = request('/api2/note/getSyncNotes?afterUsn=1000000&maxEntry=500', token=token)
assert status == 200 and len(page) == 500
assert page[0]['Usn'] == 1000001 and page[-1]['Usn'] == 1000500
_, tail = request('/api2/note/getSyncNotes?afterUsn=1000500&maxEntry=500', token=token)
assert [n['NoteId'] for n in tail] == [ids[-1]]
status, notes = request('/api2/note/getNoteContents', {'noteIds': ids[:20][::-1]}, token, True)
assert status == 200 and [n['NoteId'] for n in notes] == ids[:20][::-1]
assert notes[-1]['Content'] == '' and notes[-1]['Usn'] == 1000001
for requested in ([], ids[:21], [ids[0], ids[0]], ['not-an-id'], ['f' * 24]):
    _, error = request('/api2/note/getNoteContents', {'noteIds': requested}, token)
    assert not error['Ok'], (requested, error)
status, _ = request('/api2/note/getNoteContents', {'noteIds': ids[:1]})
assert status == 401
foreign = 'e' * 24
sql(f"INSERT INTO notes(id,user_id,notebook_id,title,usn) VALUES('{foreign}','{'f'*24}','{book}','foreign',1) ON CONFLICT(id) DO NOTHING")
_, error = request('/api2/note/getNoteContents', {'noteIds': [foreign]}, token)
assert not error['Ok']
sql(f"UPDATE note_contents SET content=repeat('x',5*1024*1024) WHERE id IN ('{ids[1]}','{ids[2]}')")
status, error = request('/api2/note/getNoteContents', {'noteIds': ids[1:3]}, token)
assert status == 413 and error['Msg'] == 'noteBatchTooLarge'
status, large = request('/api2/note/getNoteContents', {'noteIds': ids[1:2]}, token, True)
assert status == 200 and len(large[0]['Content']) == 5*1024*1024
sql(f"DELETE FROM note_contents WHERE id='{ids[3]}'")
status, error = request('/api2/note/getNoteContents', {'noteIds': ids[3:4]}, token)
assert status == 503 and error['Msg'] == 'noteContentUnavailable'
# Concurrent web saves must never pair a fresh USN/title with a previous body.
target=ids[4]
sql(f"UPDATE notes SET title='revision-0' WHERE id='{target}'; UPDATE note_contents SET content='revision-0' WHERE id='{target}'")
def write_revisions():
    usn=1000005
    for revision in range(1,21):
        value='revision-'+str(revision)
        status, saved=request('/api2/save', {'noteId':target,'notebookId':book,'title':value,'content':value,'tags':'','usn':usn,'isMarkdown':True}, token)
        assert status==200 and 'Note' in saved, saved
        usn=saved['Note']['Usn']
def read_revisions():
    for _ in range(40):
        status, values=request('/api2/note/getNoteContents', {'noteIds':[target]}, token)
        assert status==200 and values[0]['Title']==values[0]['Content']
with ThreadPoolExecutor(max_workers=2) as executor:
    writer=executor.submit(write_revisions)
    reader=executor.submit(read_revisions)
    writer.result()
    reader.result()
print('PASS: concurrent web writes preserve matching metadata/body snapshots')
print('PASS: 500 metadata / 20 bodies, ordered empty bodies, authentication, ownership, validation, 8 MiB budget, missing body')
