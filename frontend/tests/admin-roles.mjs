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
  const users=[{UserId:'self',Username:'current-admin',Email:'self@example.test',Role:'admin'},{UserId:'other',Username:'other-user',Email:'other@example.test',Role:'user'}]
  let changes=0
  await page.route('**/api2/**',async route=>{
   const req=route.request(),path=new URL(req.url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={User:users[0],IsAdmin:true}
   if(path==='/api2/admin/data')data={Users:users,Settings:{}}
   if(path==='/api2/admin/role'){
    expect(req.headers()['content-type']).toBe('application/json')
    const params=req.postDataJSON();expect(params.userId).toBe('other');users[1].Role=params.role;changes++
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
  await page.close()
  console.log(`PASS admin roles UI and disabled self-demotion: ${language}`)
 }
}finally{await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
