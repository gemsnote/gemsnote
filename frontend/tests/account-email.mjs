import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync} from 'node:fs'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=['/usr/bin/google-chrome','/usr/bin/chromium'].find(existsSync)
const browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
try{
 for(const desktop of [false,true]){
  const page=await browser.newPage({locale:'zh-CN',viewport:{width:1280,height:900}})
  await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
  let failed=false,queued=0,queries=0,pendingEmail='',verified=false,renamed='',passwordUpdates=0,logoutCount=0
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={Ok:true,User:{UserId:'owner',Username:'renamed-admin',Email:'old@example.test',Verified:verified,PendingEmail:pendingEmail},IsAdmin:true,Desktop:desktop,Notebooks:[],Tags:[],TotalNotes:0,TotalStarred:0,SharedNotebooks:{}}
   if(path==='/api2/groups'||path==='/api2/notes')data=[]
   if(path==='/api2/user/info')data={UserId:'owner',Email:'old@example.test',Verified:verified,PendingEmail:pendingEmail}
   if(path==='/api2/user/updateUsername')renamed=new URLSearchParams(route.request().postData()).get('username')
   if(path==='/api2/user/updatePwd'){
    passwordUpdates++
    const body=new URLSearchParams(route.request().postData())
    expect(body.get('pwd')).toBe('NewPassword12!');expect(body.get('oldPwd')).toBe('OldPassword12!')
    expect(body.has('confirmPwd')).toBe(false)
   }
   if(path==='/api2/logout')logoutCount++
   if(path==='/api2/web/emailChange'){queued++;pendingEmail='new@example.test';data={Ok:true,TaskId:'test-job',State:'queued',PendingEmail:pendingEmail}}
   if(path==='/api2/user/emailStatus'){queries++;data=failed?{State:'failed',Error:'mailAuthenticationFailed'}:{State:'sent'}}
   await route.fulfill({json:data})
  })
  await page.goto(base+'/note')
  await page.getByRole('button',{name:'用户菜单'}).click()
  await expect(page.getByRole('menuitem',{name:'管理',exact:true})).toHaveCount(desktop?0:1)
  await page.getByRole('menuitem',{name:'账号',exact:true}).click()
  await page.getByLabel('用户名',{exact:true}).fill('new-admin')
  const saveUsername=page.getByRole('button',{name:'更新用户名',exact:true})
  expect(await saveUsername.innerText()).toBe('')
  const inputBounds=await page.getByLabel('用户名',{exact:true}).boundingBox()
  const buttonBounds=await saveUsername.boundingBox()
  expect(buttonBounds.x).toBeGreaterThan(inputBounds.x+inputBounds.width)
  await saveUsername.click()
  await expect.poll(()=>renamed).toBe('new-admin')
  await page.getByLabel('新邮箱',{exact:true}).fill('new@example.test')
  await page.getByLabel('当前密码',{exact:true}).fill('secret')
  await page.getByRole('button',{name:'发送邮箱验证邮件',exact:true}).click()
  await expect(page.getByText('验证邮件已排队，正在等待发送。邮箱尚未修改。')).toBeVisible()
  await expect(page.getByText('验证邮件已交给邮件服务器，请检查收件箱和垃圾邮件。点击邮件中的链接后才会生效。')).toBeVisible()
  await expect(page.locator('.current-email')).toContainText('old@example.test')
  await expect(page.getByRole('button',{name:'邮箱未验证',exact:true})).toBeVisible()
  await expect(page.getByText('待验证的新邮箱：new@example.test',{exact:true})).toBeVisible()
  await page.reload()
  await expect(page.getByText('待验证的新邮箱：new@example.test',{exact:true})).toBeVisible()
  expect(queued).toBe(1);expect(queries).toBe(1)
  failed=true
  const verify=page.getByRole('button',{name:'邮箱未验证',exact:true})
  await verify.click()
  const dialog=page.getByRole('dialog')
  await expect(dialog).toContainText('old@example.test')
  await dialog.getByRole('button',{name:'取消',exact:true}).click()
  expect(queued).toBe(1)
  await verify.click()
  await dialog.getByRole('button',{name:'确定',exact:true}).click()
  await expect(page.getByRole('alert')).toContainText('SMTP 认证失败')
  await expect(page.getByRole('button',{name:'发送邮箱验证邮件',exact:true})).toBeEnabled()
  expect(queued).toBe(2)
  verified=true
  await page.evaluate(()=>window.dispatchEvent(new Event('focus')))
  await expect(verify).toHaveCount(0)
  await page.getByLabel('原密码',{exact:true}).fill('OldPassword12!')
  await page.getByLabel('新密码',{exact:true}).fill('NewPassword12!')
  await page.getByLabel('确认新密码',{exact:true}).fill('DoesNotMatch!')
  await page.getByRole('button',{name:'修改并重新登录',exact:true}).click()
  await expect(page.locator('#password-error')).toHaveText('两次输入的密码不一致')
  expect(passwordUpdates).toBe(0);expect(logoutCount).toBe(0)
  await page.getByLabel('确认新密码',{exact:true}).fill('NewPassword12!')
  await page.getByRole('button',{name:'修改并重新登录',exact:true}).click()
  await page.waitForURL(base+'/login')
  expect(passwordUpdates).toBe(1);expect(logoutCount).toBe(1)
  await page.close()
  console.log(`PASS account mail status/failure and administrator menu; desktop=${desktop}`)
 }
}finally{await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
