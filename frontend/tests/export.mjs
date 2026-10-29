import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync,readFileSync} from 'node:fs'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [true,false]){
  const page=await browser.newPage({acceptDownloads:true,viewport:{width:1280,height:800}})
  await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
  let count=0,cancelled=true,failed=false
  const note={NoteId:'n',NotebookId:'b',UserId:'u',Title:'珠玑',IsMarkdown:true,Tags:[]}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname;let data
   switch(path){
    case '/api2/bootstrap':data={Desktop:desktop,User:{UserId:'u'},Notebooks:[{NotebookId:'b',Title:'Book'}],Tags:[],SharedNotebooks:{}};break
    case '/api2/notes':data=[note];break
    case '/api2/document':data={Note:note,Content:'# 珠玑\n\n**原文**',Writable:true};break
    case '/api2/attachments':data={List:[]};break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/desktop/export':count++;expect(Object.fromEntries(new URLSearchParams(route.request().postData()))).toEqual({noteId:'n'});data=failed?{Ok:false,Msg:'exportSaveFailed'}:{Ok:true,Cancelled:cancelled};break
    default:throw Error(`Unexpected request: ${path}`)
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${base}/note/n`)
  const button=page.getByRole('button',{name:'导出',exact:true})
  if(desktop){
   await button.click();await expect.poll(()=>count).toBe(1);await expect(page.getByRole('status')).toHaveCount(0)
   cancelled=false;await button.click();await expect(page.getByRole('status')).toContainText('笔记已导出')
   failed=true;await button.click();await expect(page.locator('.editor>.error')).toContainText('导出失败')
  }else{
   const pending=page.waitForEvent('download');await button.click();const download=await pending
   expect(download.suggestedFilename()).toBe('珠玑.md');expect(readFileSync(await download.path(),'utf8')).toBe('# 珠玑\n\n**原文**');expect(count).toBe(0)
  }
  await page.close();console.log(`${desktop?'Desktop':'Web'} export passed`)
 }
}finally{await browser?.close();await server.close()}
