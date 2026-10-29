import ssl,sys,time,json
from aiosmtpd.controller import Controller
from aiosmtpd.smtp import AuthResult
class Handler:
 async def handle_RCPT(self,server,session,envelope,address,options):
  if address=='reject@example.test':return '550 Recipient rejected'
  envelope.rcpt_tos.append(address);return '250 OK'
 async def handle_DATA(self,server,session,envelope):
  with open(sys.argv[3],'a') as f:f.write(json.dumps({'to':envelope.rcpt_tos,'content':envelope.content.decode()})+'\n')
  return '250 Message accepted'
def authenticate(server,session,envelope,mechanism,auth_data):
 return AuthResult(success=auth_data.login==b'sender@example.test' and auth_data.password==b'test-secret',handled=False)
context=ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER);context.load_cert_chain(sys.argv[1],sys.argv[2])
controller=Controller(Handler(),hostname='127.0.0.1',port=10258,ssl_context=context,authenticator=authenticate,auth_required=True,auth_require_tls=False)
controller.start()
starttls=Controller(Handler(),hostname="127.0.0.1",port=10259,tls_context=context,require_starttls=True,authenticator=authenticate,auth_required=True)
starttls.start()
while True:time.sleep(1)
