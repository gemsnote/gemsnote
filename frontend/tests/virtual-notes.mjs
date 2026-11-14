import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/google-chrome'})
try{
 for(const desktop of [false,true]){
  const page=await browser.newPage({viewport:{width:1440,height:900}});page.setDefaultTimeout(10000)
  await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
  const errors=[];page.on('pageerror',e=>errors.push(e.message))
  const notes=Array.from({length:1201},(_,i)=>({NoteId:`n${i}`,UserId:'u',NotebookId:'b',Title:`Note ${String(i).padStart(4,'0')}`,UpdatedTime:new Date(2026,0,1,0,0,1201-i).toISOString(),IsMarkdown:true}))
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname;const p=JSON.parse(route.request().postData()||'{}');let data={Running:false}
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'u'},Notebooks:[{NotebookId:'b',Title:'Book',NumberNotes:1201}],Tags:[],TotalNotes:1201}
   if(path==='/api2/notes')data=notes.slice((p.page-1)*100,p.page*100)
   if(path==='/api2/document')data={Note:notes.find(n=>n.NoteId===p.noteId),Content:'full body',Writable:true}
   if(path==='/api2/attachments')data={List:[]}
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note`)
  await expect(page.locator('.virtual-note-list')).toBeVisible()
  expect(await page.locator('.note-item').count()).toBeLessThan(40)
  await page.locator('.note-item').first().click()
  await expect(page.locator('.toolbar-title')).toHaveValue('Note 0000')
  const list=page.locator('.list-items')
  await list.evaluate(el=>el.scrollTop=el.scrollHeight)
  const last=page.locator('.note-item').filter({hasText:'Note 1200'})
  await expect(last).toBeVisible();expect(await page.locator('.note-item').count()).toBeLessThan(40)
  await last.click({modifiers:['Shift']})
  await expect(page.locator('.note-selection-bar')).toContainText('1201')
  await page.keyboard.press('Escape');await last.click()
  await expect(page.locator('.toolbar-title')).toHaveValue('Note 1200')
  await page.locator('.note-sort-menu summary').click();await page.getByRole('menuitemradio',{name:'标题：降序',exact:true}).click()
  await expect(page.locator('.note-item').first().locator('strong')).toHaveText('Note 1200')
  expect(errors).toEqual([]);await page.close();console.log(`${desktop?'Desktop':'Web'}: 1201 rows, bounded DOM, last note, cross-window Shift selection, sorting passed`)
 }
}finally{await browser.close();await server.close()}
