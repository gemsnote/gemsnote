import {chromium,expect} from '@playwright/test'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/chromium','/usr/bin/google-chrome'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const language of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1280,height:800}}),errors=[]
  const label=(zh,en)=>language==='zh-CN'?zh:en
  await page.addInitScript(language=>localStorage.setItem('gemsnote:language',language),language)
  page.on('pageerror',error=>errors.push(error.message))
  let pending,versionRequests=0,minimum='',malformed=false
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data
   switch(path){
    case '/api2/bootstrap':data={User:{UserId:'user',Username:'Tester'},Desktop:false,Notebooks:[],Tags:[]};break
    case '/api2/notes':data=[];break
    case '/api2/system/version':
     expect(route.request().method()).toBe('GET');versionRequests++
     if(versionRequests===1){pending=route;return}
     data=malformed?{}:{server:'gemsnote',version:'9.8.7-test',min_version:minimum};break
    default:errors.push(`Unexpected endpoint: ${path}`);data={Ok:false}
   }
   await route.fulfill({json:data})
  })
  await page.goto(base+'/note')
  async function open(){
   await page.getByRole('button',{name:label('用户菜单','User menu'),exact:true}).click()
   await page.getByRole('menuitem',{name:label('关于','About'),exact:true}).click()
  }
  await open()
  const dialog=page.getByRole('dialog',{name:label('关于珠玑笔记','About Gemsnote')})
  await expect(dialog.getByRole('status')).toHaveText(label('正在加载…','Loading…'))
  await expect(page.locator('.user-popover')).toHaveCount(0)
  await expect.poll(()=>!!pending).toBe(true)
  await pending.fulfill({status:503,json:{Ok:false}})
  await expect(dialog.getByRole('alert')).toBeVisible()
  await dialog.getByRole('button',{name:label('重试','Retry'),exact:true}).click()
  await expect(dialog.locator('dd')).toHaveText(['gemsnote','9.8.7-test',base,label('未限制','No minimum')])
  await expect(dialog.locator('dt')).toHaveText([label('服务端名称','Server name'),label('服务端版本','Server version'),label('服务器地址','Server address'),label('最低客户端版本','Minimum client version')])
  await expect(dialog.getByRole('img')).toBeVisible()
  await dialog.getByRole('button',{name:label('关闭','Close'),exact:true}).click()
  minimum='2.0.0'
  await open()
  await expect(dialog.locator('dd').last()).toHaveText('2.0.0')
  await dialog.getByRole('button',{name:label('关闭','Close'),exact:true}).click()
  malformed=true
  await open()
  await expect(dialog.getByRole('alert')).toContainText(label('服务端版本信息不完整','The server returned incomplete version information'))
  await expect(dialog.locator('dl')).toHaveCount(0)
  expect(errors).toEqual([])
  await page.close()
  console.log(`PASS Web about: ${language}, server version/minimum, loading, retry and invalid response`)
 }
}finally{await browser?.close();await server.close()}
