import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/google-chrome'})
try {
 for(const desktop of [false,true])for(const lang of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}})
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),lang)
  const notes=[{NoteId:'a',UserId:'me',NotebookId:'one',Title:'Match one'},{NoteId:'b',UserId:'me',NotebookId:'two',Title:'Match two'}]
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data={Running:false}
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'me'},Notebooks:[{NotebookId:'one',Title:'One',NumberNotes:1}],Tags:[],TotalNotes:2}
   if(path==='/api2/notes'){const params=JSON.parse(route.request().postData());data=notes.filter(n=>!params.notebookId||n.NotebookId===params.notebookId)}
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note`)
  await page.locator('.notebook-select').click()
  const rows=page.locator('.note-item');await expect(rows).toHaveCount(1)
  await page.locator('.note-search').fill('Match');await page.locator('.note-search').press('Enter')
  const scope=page.getByRole('checkbox',{name:lang==='zh-CN'?'搜索所有笔记本':'Search all notebooks'})
  await scope.check();await expect(rows).toHaveCount(2)
  await scope.uncheck();await expect(rows).toHaveCount(1)
  await scope.check();await expect(rows).toHaveCount(2)
  await page.locator('.note-search').fill('');await page.locator('.note-search').press('Enter');await expect(rows).toHaveCount(1)
  await page.close();console.log(`${desktop?'Desktop':'Web'} ${lang}: current/all notebook search scope passed`)
 }
}finally{await browser.close();await server.close()}
