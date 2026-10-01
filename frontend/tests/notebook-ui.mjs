import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:process.env.GEMSNOTE_CHROMIUM_PATH||'/usr/bin/chromium'})
try {
 for(const desktop of [false,true]) {
 const page=await browser.newPage({viewport:{width:1440,height:900}})
 await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
 const books=[{NotebookId:'a',Title:'Parent',Subs:[{NotebookId:'b',ParentNotebookId:'a',Title:'Child',Subs:[{NotebookId:'c',ParentNotebookId:'b',Title:'Grandchild'}]}]},{NotebookId:'d',Title:'Other'}]
 const mutations=[]
 await page.route('**/api2/**',async route=>{
  const path=new URL(route.request().url()).pathname
  let data={Ok:true}
  if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'u',Username:'user'},Notebooks:books,Tags:[],TotalNotes:1}
  else if(path==='/api2/notes')data=[{NoteId:'n',NotebookId:'a',Title:'Test note',IsMarkdown:true}]
  else if(path==='/api2/document')data={Note:{NoteId:'n',NotebookId:'a',UserId:'u',Title:'Test note',IsMarkdown:true,Tags:[]},Content:'# Hello',Writable:true}
  else if(path==='/api2/attachments')data={Ok:true,List:[]}
  else if(path.includes('moveNotebook')||path.includes('copyNote')||path.includes('moveNote'))mutations.push({path,body:path.includes('moveNotebook')?new URLSearchParams(route.request().postDataJSON()):new URLSearchParams(route.request().postData()||'')})
  await route.fulfill({json:data})
 })
 await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note/n`)
 await expect(page.locator('.rendered')).toContainText('Hello')
 await expect(page.locator('.content-editor')).toHaveCount(0)
 await page.getByRole('button',{name:'移动/复制',exact:true}).click()
 const dialog=page.getByRole('dialog')
 await expect(dialog.getByRole('button',{name:'Grandchild',exact:true})).toBeVisible()
 const input=dialog.getByRole('textbox')
 await input.focus()
 expect(await input.evaluate(el=>getComputedStyle(el).outlineOffset)).toBe('-2px')
 await page.evaluate(()=>{const host=document.createElement('div');host.className='mce-window';host.innerHTML='<input class="mce-textbox" aria-label="focus-check">';document.body.append(host)})
 await page.getByLabel('focus-check').focus()
 expect(await page.getByLabel('focus-check').evaluate(el=>getComputedStyle(el).boxShadow)).toContain('inset')
 await page.evaluate(()=>document.querySelector('.mce-window')?.remove())
 await input.fill('Grandchild')
 await expect(dialog.getByRole('button',{name:'Parent',exact:true})).toBeVisible()
 await dialog.getByRole('button',{name:'Grandchild',exact:true}).click()
 await dialog.getByRole('button',{name:'复制',exact:true}).click()
 expect(mutations.at(-1).body.get('notebookId')).toBe('c')
 await page.locator('.notebook-row').filter({hasText:'Parent'}).hover()
 await page.getByRole('button',{name:'Parent 菜单',exact:true}).click()
 await page.getByRole('menuitem',{name:'移动',exact:true}).click()
 const bookDialog=page.getByRole('dialog',{name:'移动笔记本'})
 await expect(bookDialog.getByRole('button',{name:'Parent',exact:true})).toHaveCount(0)
 await expect(bookDialog.getByRole('button',{name:'Child',exact:true})).toHaveCount(0)
 await bookDialog.getByRole('button',{name:'Other',exact:true}).click()
 await bookDialog.getByRole('button',{name:'移动',exact:true}).click()
 expect(mutations.at(-1).body.get('parentNotebookId')).toBe('d')
 await page.close()
 }
 console.log('Web/desktop notebook tree, move, preview and focus checks passed')
} finally {await browser.close();await server.httpServer.close()}
