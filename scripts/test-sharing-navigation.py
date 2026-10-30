#!/usr/bin/env python3
"""Run only against an isolated seeded local PostgreSQL/API2 test instance."""
import json, os, secrets, subprocess, urllib.request, urllib.parse, urllib.error
BASE=os.environ.get('GEMSNOTE_API2_URL','http://127.0.0.1:19090').rstrip('/')
PG=os.environ.get('GEMSNOTE_TEST_PG_CONTAINER','')
assert os.environ.get('GEMSNOTE_DISPOSABLE_API2')=='1'
assert urllib.parse.urlsplit(BASE).hostname in ('127.0.0.1','localhost')
assert PG.startswith('gemsnote-api2-test-')
def sql(statement):
 return subprocess.check_output(['docker','exec',PG,'psql','-U','gemsnote_test','-d','gemsnote','-At','-v','ON_ERROR_STOP=1','-c',statement],text=True).strip()
def call(path,data=None,token=None,form=False):
 headers={}
 if token:headers['Authorization']='Bearer '+token
 if data is not None:
  headers['Content-Type']='application/x-www-form-urlencoded' if form else 'application/json'
  data=(urllib.parse.urlencode(data) if form else json.dumps(data)).encode()
 try:r=urllib.request.urlopen(urllib.request.Request(BASE+path,data,headers),timeout=15)
 except urllib.error.HTTPError as e:r=e
 with r:return r.status,json.load(r)
def api(path,data=None,token=None,form=False):
 status,result=call(path,data,token,form);assert status==200,(path,status,result);return result
admin=api('/api2/auth/login',{'email':'admin','pwd':'gemsnote'});at=admin['Token'];owner=admin['User']['UserId']
uid=secrets.token_hex(12)
# Test account only; no invitation or external mail is sent.
sql(f"INSERT INTO users(id,email,username,username_raw,pwd,verified,role) SELECT '{uid}','reader@example.test','reader','reader',pwd,true,'user' FROM users WHERE id='{owner}'")
assert api('/api2/web/adminResetPwd',{'userId':uid,'pwd':'reader-test-password'},at,True)['Ok']
rt=api('/api2/auth/login',{'email':'reader','pwd':'reader-test-password'})['Token']
book=api('/api2/bootstrap',token=at)['Notebooks'][0]['NotebookId']
target=secrets.token_hex(12)
assert api('/api2/notebook/addNotebook',{'notebookId':target,'title':'Reader shelf'},rt,True).get('NotebookId')==target
note=secrets.token_hex(12)
original=api('/api2/save',{'noteId':note,'notebookId':book,'title':'Shared original','content':'# original','isNew':True,'isMarkdown':True,'tags':'shared'},at)
share=api('/api2/share/addShareNote',{'noteId':note,'emails[0]':'reader@example.test','perm':0},at,True)
assert share['reader@example.test']['Ok'],share
out=api('/api2/bootstrap',token=at)['Sharing']['Recipients']
assert any(u['UserId']==uid and u['Username']=='reader' and note in u['NoteIds'] for u in out),out
received=api('/api2/document',{'noteId':note},rt)
assert not received['Writable'] and received['Note']['IsShared']
assert api('/api2/shared/move',{'noteId':note,'notebookId':target},rt)['Ok']
placed=api('/api2/bootstrap',token=rt)['Sharing']
assert placed['Placements'][note]==target and placed['PlacedShares'][0]['NotebookId']==target
assert api('/api2/document',{'noteId':note},rt)['Note']['NotebookId']==target
unchanged=api('/api2/document',{'noteId':note},at)
assert unchanged==original,(unchanged,original)
assert not api('/api2/shared/move',{'noteId':note,'notebookId':book},rt)['Ok']
assert not api('/api2/shared/move',{'noteId':note,'notebookId':target},at)['Ok']
assert call('/api2/shared/move',{'noteId':note,'notebookId':target})[0]==401
assert not api('/api2/save',{'noteId':note,'notebookId':target,'ownerId':owner,'usn':original['Note']['Usn'],'title':'overwrite','content':'bad','isMarkdown':True},rt).get('Writable',False)
copy=api('/api2/note/copySharedNote',{'noteIds[0]':note,'notebookId':target,'fromUserId':owner},rt,True)
assert copy['Ok'],copy
copyid=copy['Item'][0]['NoteId'];assert copyid!=note
copied=api('/api2/document',{'noteId':copyid},rt)
assert copied['Writable'] and copied['Note']['UserId']==uid and copied['Content']=='# original'
changed=api('/api2/save',{'noteId':copyid,'notebookId':target,'usn':copied['Note']['Usn'],'title':'Independent','content':'# editable copy','isMarkdown':True},rt)
assert changed['Content']=='# editable copy'
assert api('/api2/document',{'noteId':note},at)==original
snapshot=api('/api2/shared/snapshots',{},rt)
assert snapshot['Sharing']['Placements'][note]==target
assert api('/api2/share/deleteShareNote',{'noteId':note,'toUserId':uid},at,True)
assert not api('/api2/bootstrap',token=rt)['Sharing']['PlacedShares']
assert call('/api2/document',{'noteId':note},rt)[0]==403
assert not api('/api2/shared/move',{'noteId':note,'notebookId':target},rt)['Ok']
assert api('/api2/document',{'noteId':copyid},rt)['Content']=='# editable copy'
# Inherited notebook and group shares must contribute effective recipients too.
group=secrets.token_hex(12)
sql(f"INSERT INTO groups(id,user_id,title) VALUES('{group}','{owner}','Readers'); INSERT INTO group_users(id,group_id,user_id) VALUES('{secrets.token_hex(12)}','{group}','{uid}'); INSERT INTO share_notebooks(id,user_id,to_group_id,notebook_id) VALUES('{secrets.token_hex(12)}','{owner}','{group}','{book}')")
assert any(note in u['NoteIds'] for u in api('/api2/bootstrap',token=at)['Sharing']['Recipients'])
assert api('/api2/document',{'noteId':note},rt)['Writable'] is False
print('PASS shared navigation, recipient placement, owner isolation, independent copy, revocation, group/notebook grants')
