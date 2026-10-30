import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync} from 'node:fs'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=['/usr/bin/google-chrome','/usr/bin/chromium'].find(existsSync)
const browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
try{
 for(const language of ['zh-CN','en-US']){
  const page=await browser.newPage({locale:language})
  await page.addInitScript(lang=>localStorage.setItem('gemsnote:language',lang),language)
  const users=[{UserId:'self',Username:'current-admin',Email:'self@example.test',Role:'admin',Flag:'A',Verified:true},{UserId:'other',Username:'other-user',Email:'other@example.test',Role:'user',Flag:'A'}]
  let changes=0,flagChanges=0,createdUser=null,failRole=false,resends=0,failResend=false,failSettings=false,releaseSave
  await page.route('**/api2/**',async route=>{
   const req=route.request(),path=new URL(req.url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={User:users[0],IsAdmin:true}
   if(path==='/api2/admin/data')data={Users:users,Settings:{}}
   if(path==='/api2/web/adminSettings'){
    await new Promise(resolve=>releaseSave=resolve)
    data=failSettings?{Ok:false,Msg:'noAuth'}:{Ok:true}
   }
   if(path==='/api2/admin/role'){
    expect(req.headers()['content-type']).toBe('application/json')
    const params=req.postDataJSON();expect(params.userId).toBe('other')
    if(failRole){await route.fulfill({json:{Ok:false,Msg:'noAuth'}});return}
    users[1].Role=params.role;changes++
   }
   if(path==='/api2/web/adminRegister'){
    expect(req.headers()['accept-language']).toBe(language)
    createdUser=Object.fromEntries(new URLSearchParams(req.postData()))
    expect(createdUser.username).toBe('new-user');expect(createdUser.email).toBe('new@example.test');expect(createdUser.pwd).toBeUndefined()
    data={Ok:true,UserId:'created',TaskId:'invitation',MailError:''}
   }
   if(path==='/api2/admin/resendVerification'){
    expect(req.headers()['accept-language']).toBe(language)
    expect(req.headers()['content-type']).toBe('application/json');expect(req.postDataJSON()).toEqual({userId:'other'});resends++
    data=failResend?{Ok:false,Msg:'mailQueueUnavailable'}:{Ok:true,TaskId:'resend-task',State:'queued'}
   }
   if(path==='/api2/user/emailStatus'){
    const params=Object.fromEntries(new URLSearchParams(req.postData()));expect(['created','other']).toContain(params.userId);data={State:'sent'}
   }
   if(path==='/api2/admin/flag'){
    expect(req.headers()['content-type']).toBe('application/json')
    const params=req.postDataJSON();expect(params.userId).toBe('other');users[1].Flag=params.flag;flagChanges++
   }
   await route.fulfill({json:data})
  })
  await page.goto(base+'/admin')
  const label=(zh,en)=>language==='zh-CN'?zh:en
  const self=page.getByRole('row').filter({hasText:'current-admin'})
  await expect(self.getByRole('combobox')).toBeDisabled()
  await expect(self.getByRole('combobox')).toHaveValue('admin')
  const other=page.getByRole('row').filter({hasText:'other-user'})
  const role=other.getByRole('combobox')
  await expect(role).toHaveValue('user')
  await role.selectOption('admin')
  await expect(role).toBeEnabled()
  await expect(role).toHaveValue('admin')
  await role.selectOption('user')
  await expect(role).toBeEnabled()
  await expect(role).toHaveValue('user')
  failRole=true
  await role.selectOption('admin')
  await expect(role).toBeEnabled()
  await expect(role).toHaveValue('user')
  failRole=false
  expect(changes).toBe(2)
  await expect(self.getByRole('button',{name:label('禁用用户','Disable user')})).toBeDisabled()
  await other.getByRole('button',{name:label('禁用用户','Disable user')}).click()
  await expect(other.getByRole('button',{name:label('启用用户','Enable user')})).toHaveClass(/is-disabled/)
  await other.getByRole('button',{name:label('启用用户','Enable user')}).click()
  await expect(other.getByRole('button',{name:label('禁用用户','Disable user')})).toHaveClass(/is-available/)
  expect(flagChanges).toBe(2)
  await expect(self.getByRole('button',{name:label('邮箱未验证','Email not verified'),exact:true})).toHaveCount(0)
  const resend=other.getByRole('button',{name:label('邮箱未验证','Email not verified'),exact:true})
  await resend.click()
  await expect(page.getByRole('dialog')).toContainText('other@example.test')
  await page.getByRole('dialog').getByRole('button',{name:label('取消','Cancel'),exact:true}).click()
  expect(resends).toBe(0)
  await resend.click()
  await page.getByRole('dialog').getByRole('button',{name:label('确定','Confirm'),exact:true}).click()
  await expect(resend).toBeDisabled()
  await expect(other.getByRole('status')).toHaveText(label('验证邮件已发送','Verification email sent'))
  expect(resends).toBe(1)
  failResend=true
  await resend.click()
  await page.getByRole('dialog').getByRole('button',{name:label('确定','Confirm'),exact:true}).click()
  await expect(other.getByRole('status')).toContainText(label('无法连接邮件任务队列','The mail queue is unavailable'))
  await expect(resend).toBeEnabled()

  const form=page.locator('form').filter({has:page.getByRole('button',{name:label('创建用户','Create user'),exact:true})})
  await expect(form.locator('input')).toHaveCount(2)
  await form.getByLabel(label('用户名','Username'),{exact:true}).fill('new-user')
  await form.getByLabel(label('邮箱','Email'),{exact:true}).fill('new@example.test')
  await form.getByRole('button',{name:label('创建用户','Create user'),exact:true}).click()
  await expect(page.getByText(label('验证邮件已发送，请用户点击链接验证邮箱并设置密码。','Verification email sent. Ask the user to verify their email and set a password using the link.'),{exact:true})).toBeVisible()
  await expect(form.getByLabel(label('用户名','Username'),{exact:true})).toHaveValue('')
  const settings=page.locator('form').filter({has:page.getByLabel(label('邮件服务器','Mail server'),{exact:true})})
  await settings.getByLabel(label('邮件密码','Mail password'),{exact:true}).fill('temporary-test-password')
  await settings.getByRole('button',{name:label('保存设置','Save settings'),exact:true}).click()
  await expect(settings.getByRole('button',{name:label('正在保存…','Saving…'),exact:true})).toBeDisabled()
  await expect(settings.getByLabel(label('邮件服务器','Mail server'),{exact:true})).toBeDisabled()
  await expect.poll(()=>typeof releaseSave).toBe('function');releaseSave();releaseSave=undefined
  await expect(settings.getByRole('status')).toHaveText(label('设置已保存','Settings saved'))
  await expect(settings.getByLabel(label('邮件密码','Mail password'),{exact:true})).toHaveValue('')
  failSettings=true
  await settings.getByLabel(label('邮件服务器','Mail server'),{exact:true}).fill('smtp.example.test')
  await expect(settings.getByRole('status')).toHaveCount(0)
  await settings.getByRole('button',{name:label('保存设置','Save settings'),exact:true}).click()
  await expect.poll(()=>typeof releaseSave).toBe('function');releaseSave();releaseSave=undefined
  await expect(settings.getByRole('alert')).toContainText(label('设置保存失败','Could not save settings'))
  await expect(settings.getByRole('button',{name:label('保存设置','Save settings'),exact:true})).toBeEnabled()
  // A logged-in administrator opening another account's invitation must still
  // reach setup, rather than being redirected to their own notebook page.
  let resetRequests=0
  await page.route('**/api2/web/verifyEmail',route=>route.fulfill({json:{Ok:true,Email:'new@example.test',ResetToken:'setup-proof'}}))
  await page.route('**/api2/auth/password/reset',route=>{
   expect(route.request().postDataJSON()).toMatchObject({token:'setup-proof',pwd:'New-user-pass12!'})
   resetRequests++;return route.fulfill({json:{Ok:true}})
  })
  await page.goto(base+'/user/activeEmail?token=email-proof')
  await expect(page).toHaveURL(/findPassword\/setup-proof\?setup=1/)
  await expect(page.getByRole('heading',{name:label('设置密码','Set password'),exact:true})).toBeVisible()
  await page.getByLabel(label('密码','Password'),{exact:true}).fill('New-user-pass12!')
  await page.getByLabel(label('确认密码','Confirm password'),{exact:true}).fill('different')
  await page.getByRole('button',{name:label('继续','Continue'),exact:true}).click()
  await expect(page.getByRole('alert')).toContainText(label('两次输入的密码不一致','The passwords do not match.'))
  expect(resetRequests).toBe(0)
  await page.getByLabel(label('确认密码','Confirm password'),{exact:true}).fill('New-user-pass12!')
  await page.getByRole('button',{name:label('继续','Continue'),exact:true}).click()
  await expect(page.getByRole('status')).toContainText(label('密码已修改，请重新登录。','Password changed. Please sign in again.'))
  expect(resetRequests).toBe(1)
  await page.unroute('**/api2/web/verifyEmail')
  let linkError='verificationLinkInvalid'
  await page.route('**/api2/web/verifyEmail',route=>route.fulfill({json:{Ok:false,Msg:linkError}}))
  await page.goto(base+'/user/activeEmail?token=unknown-proof')
  await expect(page.getByRole('status')).toContainText(label('可能已使用、被新邮件替换，或指向了不同的服务端','It may have been used, replaced by a newer email, or opened on a different server'))
  await expect(page.getByRole('status')).not.toContainText(label('已超过有效期','This link has expired'))
  linkError='verificationLinkExpired'
  await page.reload()
  await expect(page.getByRole('status')).toContainText(label('此链接已超过有效期','This link has expired'))
  await page.close()
  console.log(`PASS admin roles/flags UI and disabled self-demotion/self-disable: ${language}`)
 }
}finally{await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
