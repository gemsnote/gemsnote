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
  const users=[{UserId:'self',Username:'current-admin',Email:'self@example.test',Role:'admin',Flag:'A'},{UserId:'other',Username:'other-user',Email:'other@example.test',Role:'user',Flag:'A'}]
  let changes=0,flagChanges=0,createdUser=null
  await page.route('**/api2/**',async route=>{
   const req=route.request(),path=new URL(req.url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={User:users[0],IsAdmin:true}
   if(path==='/api2/admin/data')data={Users:users,Settings:{}}
   if(path==='/api2/admin/role'){
    expect(req.headers()['content-type']).toBe('application/json')
    const params=req.postDataJSON();expect(params.userId).toBe('other');users[1].Role=params.role;changes++
   }
   if(path==='/api2/web/adminRegister'){
    createdUser=Object.fromEntries(new URLSearchParams(req.postData()))
    expect(createdUser.username).toBe('new-user');expect(createdUser.email).toBe('new@example.test')
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
  await expect(self.getByRole('button',{name:label('设为普通用户','Make regular user')})).toBeDisabled()
  const other=page.getByRole('row').filter({hasText:'other-user'})
  await other.getByRole('button',{name:label('设为管理员','Make administrator')}).click()
  await expect(other.getByRole('button',{name:label('设为普通用户','Make regular user')})).toBeEnabled()
  await other.getByRole('button',{name:label('设为普通用户','Make regular user')}).click()
  await expect(other.getByRole('button',{name:label('设为管理员','Make administrator')})).toBeEnabled()
  expect(changes).toBe(2)
  await expect(self.getByRole('button',{name:label('禁用用户','Disable user')})).toBeDisabled()
  await other.getByRole('button',{name:label('禁用用户','Disable user')}).click()
  await expect(other.getByRole('cell',{name:label('已禁用','Disabled'),exact:true})).toBeVisible()
  await other.getByRole('button',{name:label('启用用户','Enable user')}).click()
  await expect(other.getByRole('cell',{name:label('可用','Available'),exact:true})).toBeVisible()
  expect(flagChanges).toBe(2)
  const form=page.locator('form').filter({has:page.getByLabel(label('初始密码','Initial password'),{exact:true})})
  const password=form.getByLabel(label('初始密码','Initial password'),{exact:true})
  await expect(password).toHaveAttribute('type','text')
  await expect(password).toHaveAttribute('readonly','')
  const first=await password.inputValue()
  expect(first).toHaveLength(12)
  for(const pattern of [/[A-Z]/,/[a-z]/,/[0-9]/,/[!@#$%^&*_=+?\-]/])expect(first).toMatch(pattern)
  await form.getByRole('button',{name:label('重新生成','Regenerate'),exact:true}).click()
  const generated=await password.inputValue();expect(generated).not.toBe(first)
  await form.getByLabel(label('用户名','Username'),{exact:true}).fill('new-user')
  await form.getByLabel(label('邮箱','Email'),{exact:true}).fill('new@example.test')
  await form.getByRole('button',{name:label('创建','Create'),exact:true}).click()
  await expect(page.getByLabel(label('已创建用户的初始密码','Created user’s initial password'),{exact:true})).toHaveValue(generated)
  expect(createdUser.pwd).toBe(generated)
  expect(await password.inputValue()).not.toBe(generated)
  await expect(form.getByLabel(label('用户名','Username'),{exact:true})).toHaveValue('')
  await page.close()
  console.log(`PASS admin roles/flags UI and disabled self-demotion/self-disable: ${language}`)
 }
}finally{await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
