import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/google-chrome'})
try {
 for(const desktop of [false,true])for(const lang of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}})
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),lang)
  let restored=false
  const note={NoteId:'trash',UserId:'me',NotebookId:'deleted',Title:'Trash note',IsTrash:true,IsMarkdown:true}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data=true
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'me'},Notebooks:[{NotebookId:'one',Title:'Target',NumberNotes:0}],Tags:[],TotalNotes:0}
   if(path==='/api2/notes')data=restored?[]:[note]
   if(path==='/api2/document')data={Note:note,Content:'body',Writable:true}
   if(path==='/api2/attachments')data={List:[]}
   if(path==='/api2/web/downloadStatus')data={Running:false}
   if(path==='/api2/note/moveNote'){
    const params=new URLSearchParams(route.request().postData());expect(params.get('noteIds[0]')).toBe('trash');expect(params.get('notebookId')).toBe('one');restored=true
   }
   if(path==='/api2/restore')throw Error('orphan restore must select a valid target first')
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note`)
  await page.getByRole('button',{name:lang==='zh-CN'?'回收站':'Trash',exact:true}).click()
  await page.locator('.note-item').click()
  await page.locator('.toolbar-actions').getByRole('button',{name:lang==='zh-CN'?'恢复':'Restore',exact:true}).click()
  const dialog=page.getByRole('dialog')
  await expect(dialog).toBeVisible()
  const confirm=dialog.getByRole('button',{name:lang==='zh-CN'?'恢复':'Restore',exact:true})
  await expect(confirm).toBeDisabled()
  await dialog.getByRole('button',{name:'Target',exact:true}).click()
  await confirm.click();await expect(dialog).toHaveCount(0);expect(restored).toBe(true)
  await page.close();console.log(`${desktop?'Desktop':'Web'} ${lang}: orphan trash restoration requires a target notebook`)
 }
}finally{await browser.close();await server.close()}
