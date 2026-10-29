import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync,readFileSync} from 'node:fs'
import {unzipSync,strFromU8} from 'fflate'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [true,false]){
  const page=await browser.newPage({acceptDownloads:true,viewport:{width:1280,height:800}})
  await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
  let count=0,cancelled=true,failed=false,exported,unreadable=false,paginated=false
  const note={NoteId:'n',NotebookId:'b',UserId:'u',Title:'珠玑',IsMarkdown:true,Tags:['test'],CreatedTime:'2020-01-01T00:00:00Z',UpdatedTime:'2021-01-01T00:00:00Z'}
  const child={...note,NoteId:'c',NotebookId:'child',Title:'富文本',IsMarkdown:false}
  await page.route('**/api2/**',async route=>{
   const req=route.request(),path=new URL(req.url()).pathname
   const params=req.headers()['content-type']?.includes('application/json')?JSON.parse(req.postData()||'{}'):Object.fromEntries(new URLSearchParams(req.postData()))
   let data
   switch(path){
    case '/api2/bootstrap':data={Desktop:desktop,User:{UserId:'u'},Notebooks:[{NotebookId:'b',Title:'Book',Subs:[{NotebookId:'child',Title:'Child',ParentNotebookId:'b'}]}],Tags:[],SharedNotebooks:{}};break
    case '/api2/notes':data=params.notebookId==='child'?[child]:paginated?Array.from({length:101},(_,i)=>({...note,NoteId:'p'+i,Title:'Page '+i})).slice(((params.page||1)-1)*100,(params.page||1)*100):[note];break
    case '/api2/document':data=unreadable?{Ok:false,Msg:'readFailed'}:params.noteId==='c'?{Note:child,Content:'<h1>Child body</h1>',Writable:true}:{Note:params.noteId?.startsWith('p')?{...note,NoteId:params.noteId,Title:'Page '+params.noteId.slice(1)}:note,Content:'# 珠玑\n\n**原文**',Writable:true};break
    case '/api2/attachments':data={List:[]};break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/desktop/exportFile':count++;exported=params;data=failed?{Ok:false,Msg:'exportSaveFailed'}:{Ok:true,Cancelled:cancelled};break
    default:throw Error(`Unexpected request: ${path}`)
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${base}/note/n`)
  const button=page.getByRole('button',{name:'导出',exact:true})
  if(desktop){
   await button.click();await expect.poll(()=>count).toBe(1);await expect(page.getByRole('status')).toHaveCount(0)
   expect(exported.filename).toBe('珠玑.md');expect(Buffer.from(exported.data,'base64').toString()).toContain('categories: ["Book"]')
   cancelled=false;await button.click();await expect(page.getByRole('status')).toContainText('笔记已导出')
   failed=true;await button.click();await expect(page.locator('.editor>.error')).toContainText('导出失败');failed=false
  }else{
   const pending=page.waitForEvent('download');await button.click();const download=await pending
   expect(download.suggestedFilename()).toBe('珠玑.md');const text=readFileSync(await download.path(),'utf8');expect(text).toContain('tags: ["test"]');expect(text).toContain('lastmod: "2021-01-01');expect(text.endsWith('# 珠玑\n\n**原文**')).toBe(true)
  }
  const bookExport=async()=>{await page.locator('.notebook-row').first().hover();await page.getByRole('button',{name:'Book 菜单',exact:true}).click();await page.getByRole('menuitem',{name:'导出笔记本（ZIP）'}).click()}
  let bytes
  if(desktop){const before=count;await bookExport();await expect.poll(()=>count).toBe(before+1);expect(exported.filename).toBe('Book.zip');bytes=Buffer.from(exported.data,'base64')}
  else {const pending=page.waitForEvent('download');await bookExport();const result=await pending;expect(result.suggestedFilename()).toBe('Book.zip');bytes=readFileSync(await result.path())}
  const files=unzipSync(bytes);expect(Object.keys(files)).toEqual(['Book/','Book/珠玑.md','Book/Child/','Book/Child/富文本.rtf']);expect(strFromU8(files['Book/Child/富文本.rtf'])).toContain('Child body')
  paginated=true
  if(desktop){const before=count;await bookExport();await expect.poll(()=>count,{timeout:30000}).toBe(before+1);bytes=Buffer.from(exported.data,'base64')}
  else {const pending=page.waitForEvent('download');await bookExport();bytes=readFileSync(await (await pending).path())}
  const pagedFiles=unzipSync(bytes);expect(Object.keys(pagedFiles).filter(n=>/\.(md|rtf)$/.test(n))).toHaveLength(102);expect(pagedFiles['Book/Page 100.md']).toBeTruthy()
  unreadable=true;await bookExport();await expect(page.locator('.editor>.error')).toContainText('readFailed')
  await page.close();console.log(`${desktop?'Desktop':'Web'} Hugo note and notebook export passed`)
 }
}finally{await browser?.close();await server.close()}
