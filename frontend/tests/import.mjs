import {chromium,expect} from '@playwright/test'
import {zipSync,strToU8} from 'fflate'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'

const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const baseURL=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [false,true]){
  const page=await browser.newPage({viewport:{width:1280,height:800}}),errors=[]
  page.on('pageerror',e=>errors.push(e.message))
  await page.addInitScript(language=>localStorage.setItem('gemsnote:language',language),desktop?'en-US':'zh-CN')
  const label=(zh,en)=>desktop?en:zh
  const userId='507f1f77bcf86cd799439011',initial='507f1f77bcf86cd799439012'
  const target='507f1f77bcf86cd799439014'
  const books=[{NotebookId:initial,Title:'Existing',ParentNotebookId:''},{NotebookId:target,Title:'Destination',ParentNotebookId:''}],notes=new Map(),writes=[]
  let loseResponse=false
  function tree(parent=''){return books.filter(b=>b.ParentNotebookId===parent).map(b=>({...b,Subs:tree(b.NotebookId),NumberNotes:[...notes.values()].filter(n=>n.Note.NotebookId===b.NotebookId).length}))}
  await page.route('**/api2/**',async route=>{
   const req=route.request(),path=new URL(req.url()).pathname
   const body=req.postData()||'',params=req.headers()['content-type']?.includes('application/json')?JSON.parse(body||'{}'):Object.fromEntries(new URLSearchParams(body))
   let data
   switch(path){
    case '/api2/bootstrap':data={User:{UserId:userId,Username:'tester'},Desktop:desktop,Notebooks:tree(),Tags:[],SharedNotebooks:{},TotalNotes:notes.size,PendingChanges:desktop&&notes.size>0};break
    case '/api2/notes':data=[...notes.values()].filter(n=>!params.notebookId||n.Note.NotebookId===params.notebookId).map(n=>n.Note);break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/attachments':data={List:[]};break
    case '/api2/notebook/addNotebook':
     expect(books.some(b=>b.NotebookId===params.notebookId)).toBe(false)
     books.push({NotebookId:params.notebookId,Title:params.title,ParentNotebookId:params.parentNotebookId});writes.push(path);data={Ok:true};break
    case '/api2/save':
     expect(req.headers()['content-type']).toBe('application/json')
     expect(params.isNew).toBe(true);expect(params.isMarkdown).toBe(true)
     expect(notes.has(params.noteId)).toBe(false)
     data={Note:{NoteId:params.noteId,UserId:userId,NotebookId:params.notebookId,Title:params.title,IsMarkdown:true,Tags:[],IsDirty:desktop,CreatedTime:'2026-10-28T10:00:00Z',UpdatedTime:'2026-10-28T10:00:00Z'},Content:params.content,Writable:true}
     notes.set(params.noteId,data);writes.push(path)
     if(loseResponse){loseResponse=false;await route.abort('failed');return}break
    case '/api2/document':data=notes.get(params.noteId)||{Ok:false,Msg:'notExists'};break
    default:errors.push(`Unexpected request: ${path}`);data={Ok:false,Msg:'unexpected'}
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${baseURL}/note`)
  await expect(page.locator('.notebook-more')).toHaveCount(2)
  async function showImport(){await page.locator('.notebook-row').filter({hasText:'Destination'}).hover();await page.getByRole('button',{name:label('Destination 菜单','Destination menu'),exact:true}).click();await page.getByRole('menuitem',{name:label('导入笔记/笔记本','Import notes / notebook'),exact:true}).click()}
  await expect(page.locator('.panel-header').getByRole('button',{name:label('导入 Markdown','Import Markdown'),exact:true})).toHaveCount(0)
  await showImport()
  let dialog=page.getByRole('dialog',{name:label('导入 Markdown','Import Markdown'),exact:true})
  await dialog.locator('input[type=file]').setInputFiles({name:'单篇.MD',mimeType:'text/markdown',buffer:Buffer.from('# 你好\n\n完整正文')})
  await expect(dialog).toContainText(label('将导入 1 篇笔记','Import 1 notes'))
  loseResponse=true
  await dialog.getByRole('button',{name:label('开始导入','Start import'),exact:true}).click()
  await expect(dialog.getByRole('alert')).toBeVisible()
  await dialog.getByRole('button',{name:label('继续导入','Resume import'),exact:true}).click()
  await expect(dialog).toContainText(label('导入完成','Import complete'))
  await dialog.getByRole('button',{name:label('关闭','Close'),exact:true}).click()
  await expect(page.locator('.editor')).toContainText('完整正文')
  expect(notes.size).toBe(1)
  expect([...notes.values()][0].Note.NotebookId).toBe(target)

  await showImport()
  dialog=page.getByRole('dialog',{name:label('导入 Markdown','Import Markdown'),exact:true})
  await dialog.locator('input[type=file]').setInputFiles({name:'坏包.zip',mimeType:'application/zip',buffer:Buffer.from(zipSync({'../bad.md':strToU8('unsafe')}))})
  await expect(dialog.getByRole('alert')).toContainText(label('不安全','unsafe'))
  expect(writes.length).toBe(1)
  const archive=zipSync({'root.md':strToU8('root body'),'Folder/子目录/child.md':strToU8('# child'),'ignore.png':new Uint8Array([0,255])})
  await dialog.locator('input[type=file]').setInputFiles({name:'资料.zip',mimeType:'application/zip',buffer:Buffer.from(archive)})
  await expect(dialog).toContainText(label('将导入 2 篇笔记、3 个新笔记本；跳过 1 个文件。','Import 2 notes and 3 new notebooks; skip 1 files.'))
  await dialog.getByRole('button',{name:label('开始导入','Start import'),exact:true}).click()
  await expect(dialog).toContainText(label('导入完成','Import complete'))
  if(desktop)await expect(dialog).toContainText('Saved locally. Sync to upload to the server.')
  await dialog.getByRole('button',{name:label('关闭','Close'),exact:true}).click()
  expect(books.length).toBe(5);expect(notes.size).toBe(3)
  const root=books.find(b=>b.Title==='资料'),folder=books.find(b=>b.Title==='Folder'),child=books.find(b=>b.Title==='子目录')
  expect(root.ParentNotebookId).toBe(target);expect(folder.ParentNotebookId).toBe(root.NotebookId);expect(child.ParentNotebookId).toBe(folder.NotebookId)
  expect([...notes.values()].find(n=>n.Note.Title==='child').Note.NotebookId).toBe(child.NotebookId)
  await expect(page.locator('.notebook-row').filter({hasText:'子目录'})).toBeVisible()
  await expect(page.locator('.notebook-row').filter({hasText:'资料'}).locator('small')).toHaveText('1')
  await expect(page.locator('.notebook-row').filter({hasText:'子目录'}).locator('small')).toHaveText('1')
  if(desktop)await expect(page.locator('.sync-pending-dot')).toHaveCount(1)
  expect(errors).toEqual([])
  await page.close()
  console.log(`${desktop?'Desktop (English)':'Web (Chinese)'}: MD, ZIP hierarchy, validation, retry and counts passed`)
 }
}finally{await browser?.close();await server.close()}
