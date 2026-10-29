import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync} from 'node:fs'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const language of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1280,height:800}}),zh=language==='zh-CN'
  await page.addInitScript(lang=>localStorage.setItem('gemsnote:language',lang),language)
  let restored=false,saves=0,retries=0
  const note={NoteId:'n',NotebookId:'b',UserId:'u',Title:'Missing files',IsMarkdown:true,Tags:[],Usn:1}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data
   switch(path){
    case '/api2/bootstrap':data={Desktop:true,User:{UserId:'u'},Notebooks:[{NotebookId:'b',Title:'Book'}],SharedNotebooks:{},Tags:[],TotalNotes:1};break
    case '/api2/notes':data=[note];break
    case '/api2/document':data={Note:note,Content:'![lost](/api2/file/getImage?fileId=missing)',Writable:true};break
    case '/api2/file/getImage':await route.fulfill({status:404,body:'fileNotFound'});return
    case '/api2/attachments':data={List:[{AttachId:'a',Title:'ELK.xml',Size:8,CacheState:restored?'ready':'missing'}]};break
    case '/api2/desktop/attachments/retry':retries++;data=restored?{Ok:true}:{Ok:false,Msg:'attachmentMissing'};break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/save':saves++;data={Ok:false,Msg:'unexpected save'};break
    default:throw Error(`Unexpected request: ${path}`)
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${base}/note/n`)
  await expect(page.locator('.rendered img')).toHaveAttribute('data-unavailable','true')
  await expect(page.locator('.rendered img')).toHaveAttribute('src',/^data:image\/svg\+xml,/)
  await page.getByRole('button',{name:zh?/^附件/:/^Attachments/}).click()
  await expect(page.locator('.missing-attachment')).toContainText('ELK.xml')
  await expect(page.locator('.attachments')).toContainText(zh?'文件已丢失':'File missing')
  await expect(page.locator('.attachments a')).toHaveCount(0)
  await page.getByRole('button',{name:zh?'重试下载':'Retry download',exact:true}).click()
  await expect(page.locator('main.editor>.error')).toContainText(zh?'服务端文件仍然缺失':'The file is still missing')
  restored=true
  await page.getByRole('button',{name:zh?'重试下载':'Retry download',exact:true}).click()
  await expect(page.locator('.attachments a')).toHaveText('ELK.xml')
  expect(retries).toBe(2);expect(saves).toBe(0)
  await page.close()
  console.log(`Missing files: placeholder, image fallback and retry passed (${language})`)
 }
}finally{await browser?.close();await server.close()}
