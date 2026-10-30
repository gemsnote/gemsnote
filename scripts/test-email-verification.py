import json,urllib.request,urllib.parse,time,os,email as mail,re
from email.header import decode_header,make_header
if os.environ.get('GEMSNOTE_DISPOSABLE_API2') != '1' or not os.environ.get('MAIL_CAPTURE'):
 raise SystemExit('Requires disposable API2 fixture and local SMTP capture')
base='http://127.0.0.1:19090'
def request(path,body=None,token=None,asjson=False,language="zh-CN"):
 headers={'X-Requested-With':'XMLHttpRequest','Accept-Language':language}
 if token:headers['Authorization']='Bearer '+token
 if body is not None:
  headers['Content-Type']='application/json' if asjson else 'application/x-www-form-urlencoded'
  body=(json.dumps(body) if asjson else urllib.parse.urlencode(body)).encode()
 with urllib.request.urlopen(urllib.request.Request(base+path,body,headers),timeout=10) as response:return json.load(response)
login=request('/api2/auth/login',{'email':'admin','pwd':'gemsnote'},asjson=True)
token=login['Token'];boot=request('/api2/bootstrap',token=token);old=boot['User']['Email'];assert boot['IsAdmin'],boot
assert request('/api2/user/updateUsername',{'username':'renamed-admin'},token)['Ok']
assert request('/api2/bootstrap',token=token)['IsAdmin']
assert request('/api2/admin/data',{},token,True).get('Users') is not None
assert request('/api2/web/emailChange',{'email':'new@example.test','pwd':'wrong'},token)['Msg']=='wrongPassword'
assert request('/api2/web/emailChange',{'email':old,'pwd':'gemsnote'},token)['Msg']=='emailUnchanged'
assert request('/api2/web/adminSettings',{'emailHost':'localhost','emailPort':'10258','emailUsername':'sender@example.test','emailPassword':'test-secret','emailSSL':'1','siteUrl':base},token)['Ok']
job=request('/api2/web/emailChange',{'email':'new@example.test','pwd':'gemsnote'},token)
assert job['Ok'] and job['State']=='queued',job
assert request('/api2/web/emailChange',{'email':'new@example.test','pwd':'gemsnote'},token)['Msg']=='mailRateLimited'
profile=request('/api2/bootstrap',token=token)['User']
assert profile['Email']==old and profile['Verified']==boot['User']['Verified'] and profile['PendingEmail']=='new@example.test'
assert request('/api2/user/info',token=token)['PendingEmail']=='new@example.test'
for _ in range(50):
 result=request('/api2/user/emailStatus',{'taskId':job['TaskId']},token)
 if result['State'] in ('sent','failed'):break
 time.sleep(.2)
assert result['State']=='sent',result
profile=request('/api2/bootstrap',token=token)['User']
assert profile['Email']==old and profile['Verified']==boot['User']['Verified'] and profile['PendingEmail']=='new@example.test'
assert request('/api2/user/info',token=token)['PendingEmail']=='new@example.test'
messages=[json.loads(line) for line in open(os.environ['MAIL_CAPTURE'])]
message=next(m for m in reversed(messages) if m['to']==['new@example.test'])
parts=mail.message_from_string(message['content']).walk()
html=next(part.get_payload(decode=True).decode().replace('\r\n','\n') for part in parts if part.get_content_type()=='text/html')
verification=re.search(r'/user/updateEmail\?token=([^"&<>\s]+)',html).group(1)
assert request('/api2/auth/login',{'email':old,'pwd':'gemsnote'},asjson=True).get('Token')
assert not request('/api2/auth/login',{'email':'new@example.test','pwd':'gemsnote'},asjson=True).get('Token')
assert request('/api2/web/verifyEmail',{'token':'invalid','change':'true'})['Msg']=='verificationLinkInvalid'
assert request('/api2/web/verifyEmail',{})['Msg']=='verificationLinkMissing'
assert request('/api2/bootstrap',token=token)['User']['Email']==old
assert request('/api2/web/verifyEmail',{'token':verification,'change':'true'})['Ok']
profile=request('/api2/bootstrap',token=token)['User']
assert profile['Email']=='new@example.test' and profile['Verified'] and profile['PendingEmail']==''
assert request('/api2/auth/login',{'email':'new@example.test','pwd':'gemsnote'},asjson=True).get('Token')
assert not request('/api2/auth/login',{'email':old,'pwd':'gemsnote'},asjson=True).get('Token')
assert not request('/api2/web/verifyEmail',{'token':verification,'change':'true'})['Ok']
print('PASS live API2 admin rename, queue, password validation, cooldown, SMTP, verification and token consumption')

# Invitations ignore any supplied password; only the email owner sets one.
invited=request('/api2/web/adminRegister',{'username':'verify-user','email':'verify@example.test','pwd':'Injected12!'},token)
assert invited['Ok'] and invited['TaskId'] and not invited['MailError'],invited
for password in ['', 'Injected12!']:
 assert not request('/api2/auth/login',{'email':'verify-user','pwd':password},asjson=True).get('Token')
profile=next(u for u in request('/api2/admin/data',{},token,True)['Users'] if u['Username']=='verify-user')
assert not profile['Verified'] and 'Pwd' not in profile
for _ in range(80):
 status=request('/api2/user/emailStatus',{'taskId':invited['TaskId'],'userId':invited['UserId']},token)
 if status['State'] in ('sent','failed'):break
 time.sleep(.2)
assert status['State']=='sent',status
messages=[json.loads(line) for line in open(os.environ['MAIL_CAPTURE'])]
target=next(m for m in messages if m['to']==['verify@example.test'])
parsed_mail=mail.message_from_string(target['content'])
assert str(make_header(decode_header(parsed_mail['Subject'])))=='恭喜您已注册珠玑笔记，请验证邮箱。'
plain=next(part.get_payload(decode=True).decode().replace('\r\n','\n') for part in parsed_mail.walk() if part.get_content_type()=='text/plain')
assert 'verify-user（verify@example.test），您好：\n请点击链接验证邮箱。' in plain
assert '日积字句，终得珠玑。' in plain and '48' in plain

parts=mail.message_from_string(target['content']).walk()
html=next(part.get_payload(decode=True).decode().replace('\r\n','\n') for part in parts if part.get_content_type()=='text/html')
assert 'verify-user（verify@example.test），您好：<br>请点击链接验证邮箱。' in html
assert '日积字句，终得珠玑。' in html and 'leanote' not in html.lower()
activation=re.search(r'/user/activeEmail\?token=([^"&<>\s]+)',html).group(1)
# An unverified invitation is not itself a password-reset credential.
assert not request('/api2/auth/password/reset',{'token':activation,'pwd':'Verify-Test12!'},asjson=True)['Ok']
verified=request('/api2/web/verifyEmail',{'token':activation})
assert verified['Ok'] and verified['ResetToken'],verified
# A lost response can be retried without extending the setup deadline.
assert request('/api2/web/verifyEmail',{'token':activation})['ResetToken']==verified['ResetToken']
assert not request('/api2/auth/login',{'email':'verify-user','pwd':'Verify-Test12!'},asjson=True).get('Token')
assert request('/api2/auth/password/reset',{'token':verified['ResetToken'],'pwd':'Verify-Test12!'},asjson=True)['Ok']
created=request('/api2/auth/login',{'email':'verify-user','pwd':'Verify-Test12!'},asjson=True)
assert created.get('Token') and created['User']['Verified']
assert not request('/api2/auth/password/reset',{'token':verified['ResetToken'],'pwd':'Different12!'},asjson=True)['Ok']
assert request('/api2/web/verifyEmail',{'token':activation})['Msg']=='verificationLinkInvalid'
# Existing accounts still verify their email without starting initial setup.
assert request('/api2/web/adminRegister',{'username':'existing-user','email':'existing@example.test'},token)['Ok']
profile=next(u for u in request('/api2/admin/data',{},token,True)['Users'] if u['Username']=='existing-user')
assert request('/api2/web/adminResetPwd',{'userId':profile['UserId'],'pwd':'Existing12!'},token)['Ok']
for _ in range(80):
 messages=[json.loads(line) for line in open(os.environ['MAIL_CAPTURE'])]
 target=next((m for m in messages if m['to']==['existing@example.test']),None)
 if target:break
 time.sleep(.2)
assert target
parts=mail.message_from_string(target['content']).walk()
html=next(part.get_payload(decode=True).decode().replace('\r\n','\n') for part in parts if part.get_content_type()=='text/html')
activation=re.search(r'/user/activeEmail\?token=([^"&<>\s]+)',html).group(1)
result=request('/api2/web/verifyEmail',{'token':activation})
assert result['Ok'] and not result['ResetToken']
assert request('/api2/web/verifyEmail',{'token':activation})['Msg']=='verificationLinkInvalid'
print('PASS invitations: no assigned password, email verification, retryable setup handoff, one-time password setup, existing account activation')

# A pending account can recover through a fresh emailed password-reset link.
assert request('/api2/web/adminRegister',{'username':'recover-user','email':'recover@example.test'},token)['Ok']
assert request('/api2/auth/password/request',{'email':' RECOVER@example.test '},asjson=True)['Ok']
recovery=None
for _ in range(80):
 messages=[json.loads(line) for line in open(os.environ['MAIL_CAPTURE'])]
 for candidate in reversed(messages):
  if candidate['to']!=['recover@example.test']:continue
  parts=mail.message_from_string(candidate['content']).walk()
  html=next(part.get_payload(decode=True).decode().replace('\r\n','\n') for part in parts if part.get_content_type()=='text/html')
  match=re.search(r'/findPassword/([^"&<>\s]+)',html)
  if match:recovery=match.group(1);break
 if recovery:break
 time.sleep(.2)
assert recovery,'recovery email missing'
assert request('/api2/auth/password/reset',{'token':recovery,'pwd':'Recovered12!'},asjson=True)['Ok']
recovered=request('/api2/auth/login',{'email':'recover@example.test','pwd':'Recovered12!'},asjson=True)
assert recovered.get('Token') and recovered['User']['Verified']
assert not request('/api2/auth/password/reset',{'token':recovery,'pwd':'Other12!'},asjson=True)['Ok']
# Queue failures retain the account and report a separate mail error.
assert request('/api2/web/adminSettings',{'emailHost':''},token)['Ok']
failed=request('/api2/web/adminRegister',{'username':'mail-failed','email':'failed@example.test'},token)
assert failed['Ok'] and failed['UserId'] and failed['MailError'] and not failed['TaskId'],failed
assert not request('/api2/auth/login',{'email':'mail-failed','pwd':''},asjson=True).get('Token')
print('PASS invitation recovery: email reset sets password and verifies mailbox; queue failure is reported separately')

# The administrator can recover the failed invitation with a confirmed resend.
assert request('/api2/web/adminSettings',{'emailHost':'localhost'},token)['Ok']
resend=request('/api2/admin/resendVerification',{'userId':failed['UserId'],'email':'ignored@example.test'},token,True,language='en-US')
assert resend['Ok'] and resend['TaskId'],resend
assert request('/api2/admin/resendVerification',{'userId':failed['UserId']},token,True)['Msg']=='mailRateLimited'
assert request('/api2/admin/resendVerification',{'userId':'invalid'},token,True)['Msg']=='userNotExists'
assert request('/api2/admin/resendVerification',{'userId':created['User']['UserId']},token,True)['Msg']=='emailAlreadyVerified'
try:
 request('/api2/admin/resendVerification',{'userId':failed['UserId']},created['Token'],True)
 raise AssertionError('ordinary user must not resend another account verification')
except urllib.error.HTTPError as error:
 assert error.code==403
for _ in range(80):
 status=request('/api2/user/emailStatus',{'taskId':resend['TaskId'],'userId':failed['UserId']},token)
 if status['State'] in ('sent','failed'):break
 time.sleep(.2)
assert status['State']=='sent',status
messages=[json.loads(line) for line in open(os.environ['MAIL_CAPTURE'])]
assert any(m['to']==['failed@example.test'] for m in messages)
assert not any(m['to']==['ignored@example.test'] for m in messages)
assert request('/api2/admin/flag',{'userId':failed['UserId'],'flag':'D'},token,True)['Ok']
assert request('/api2/admin/resendVerification',{'userId':failed['UserId']},token,True)['Msg']=='userDisabled'
print('PASS admin resend: stored recipient, permissions, disabled/verified accounts, cooldown, real delivery')

english=next(m for m in reversed(messages) if m['to']==['failed@example.test'])
parsed_mail=mail.message_from_string(english['content'])
assert str(make_header(decode_header(parsed_mail['Subject'])))=='Welcome to Gemsnote — please verify your email.'
for part in parsed_mail.walk():
 if part.get_content_type() not in ('text/plain','text/html'):continue
 body=part.get_payload(decode=True).decode().replace('\r\n','\n')
 assert 'Small notes become lasting insights.' in body and '48' in body
 assert '珠玑笔记' not in body and 'leanote' not in body.lower()
 if part.get_content_type()=='text/html':assert 'Hello mail-failed (failed@example.test),<br>Please verify' in body
 else:assert 'Hello mail-failed (failed@example.test),\nPlease verify' in body
print('PASS verification mail: Chinese creation, English resend, localized HTML/plain text, greeting line breaks, 48-hour expiry')
