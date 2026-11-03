import {chromium,expect} from '@playwright/test'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'

const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const baseURL=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/chromium','/usr/bin/google-chrome'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [false,true])for(const language of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1280,height:800}}),errors=[]
  page.on('pageerror',error=>errors.push(error.message))
  await page.addInitScript(language=>localStorage.setItem('gemsnote:language',language),language)
  const notes=[
   {NoteId:'b',Title:'Alpha',CreatedTime:'2020-01-01T02:00:00Z',UpdatedTime:'2026-01-01'},
   {NoteId:'c',Title:'Beta',CreatedTime:'2020-01-01T09:00:00+08:00',UpdatedTime:'2026-03-01'},
   {NoteId:'a',Title:'Gamma',CreatedTime:'2020-01-01T02:00:00Z',UpdatedTime:'2026-02-01'},
  ]
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data
   switch(path){
    case '/api2/bootstrap': data={Desktop:desktop,User:{UserId:'user'},Notebooks:[],Tags:[],TotalNotes:3};break
    case '/api2/notes':data=notes;break
    case '/api2/web/downloadStatus':data={Running:false};break
    default:errors.push(`Unexpected API request: ${path}`);data={Ok:false}
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${baseURL}/note`)
  const titles=page.locator('.note-item strong')
  await expect(titles).toHaveText(['Beta','Gamma','Alpha'])
  for(const [zh,en,order] of [
   ['创建时间：升序','Created: oldest first',['Beta','Gamma','Alpha']],
   ['创建时间：降序','Created: newest first',['Gamma','Alpha','Beta']],
   ['修改时间：升序','Modified: oldest first',['Alpha','Gamma','Beta']],
   ['修改时间：降序','Modified: newest first',['Beta','Gamma','Alpha']],
   ['标题：升序','Title: A–Z',['Alpha','Beta','Gamma']],
   ['标题：降序','Title: Z–A',['Gamma','Beta','Alpha']],
  ]){
   await page.locator('.note-sort-menu summary').click()
   const option=page.getByRole('menuitemradio',{name:language==='zh-CN'?zh:en,exact:true,includeHidden:true})
   await option.click()
   await expect(option).toHaveAttribute('aria-checked','true')
   await expect(page.locator('.note-sort-menu')).not.toHaveAttribute('open','')
   await expect(titles).toHaveText(order)
  }
  expect(errors).toEqual([])
  await page.close()
  console.log(`${desktop?'Desktop':'Web'} ${language}: creation/modification/title ascending and descending passed`)
 }
}finally{await browser?.close();await server.close()}
