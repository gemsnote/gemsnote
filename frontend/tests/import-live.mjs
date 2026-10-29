// Real browser/API2 persistence test. Never run against a user's database.
import {chromium,expect} from '@playwright/test'
import {zipSync,strToU8} from 'fflate'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'

if(process.env.GEMSNOTE_DISPOSABLE_API2!=='1'||process.env.GEMSNOTE_API2_URL!=='http://127.0.0.1:19090')throw Error('Requires the disposable API2/PostgreSQL fixture at 127.0.0.1:19090')
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0,proxy:{'/api2/':process.env.GEMSNOTE_API2_URL}}})
const baseURL=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 const context=await browser.newContext({viewport:{width:1280,height:800}})
 const response=await context.request.post(`${baseURL}/api2/auth/login`,{data:{email:'admin',pwd:'gemsnote'}})
 const auth=await response.json()
 expect(auth.Token).toBeTruthy()
 const headers={Authorization:`Bearer ${auth.Token}`}
 await context.addInitScript(token=>{sessionStorage.setItem('gemsnote:api2-token',token);localStorage.setItem('gemsnote:language','zh-CN')},auth.Token)
 const page=await context.newPage()
 await page.goto(`${baseURL}/note`)
 const before=await (await context.request.get(`${baseURL}/api2/bootstrap`,{headers})).json()
 const target=before.Notebooks[0]
 async function showImport(){await page.locator('.notebook-row').filter({hasText:target.Title}).first().hover();await page.getByRole('button',{name:`${target.Title} 菜单`,exact:true}).click();await page.getByRole('menuitem',{name:'导入笔记/笔记本',exact:true}).click()}
 await showImport()
 let dialog=page.getByRole('dialog',{name:'导入 Markdown',exact:true})
 await dialog.locator('input[type=file]').setInputFiles({name:'导入测试.md',mimeType:'text/markdown',buffer:Buffer.from('# 单篇\n\n完整正文')})
 await dialog.getByRole('button',{name:'开始导入',exact:true}).click()
 await expect(dialog).toContainText('导入完成')
 await dialog.getByRole('button',{name:'关闭',exact:true}).click()
 await expect(page.locator('.editor')).toContainText('完整正文')

 // Saving and switching views must retain the imported Markdown format.
 await expect(page).toHaveURL(/\/note\/[0-9a-f]{24}$/)
 const noteId=new URL(page.url()).pathname.split('/').pop()
 await page.getByRole('button',{name:'编辑',exact:true}).click()
 const edited='# 单篇修改\n\n**仍是 Markdown**'
 await page.locator('textarea.content-editor').fill(edited)
 await page.getByRole('button',{name:'预览',exact:true}).click()
 const read=async()=>await (await context.request.post(`${baseURL}/api2/document`,{headers,data:{noteId}})).json()
 await expect.poll(async()=> (await read()).Content).toBe(edited)
 expect((await read()).Note.IsMarkdown).toBe(true)
 await page.getByRole('button',{name:'编辑',exact:true}).click()
 await expect(page.locator('textarea.content-editor')).toHaveValue(edited)
 // The API also protects cached clients that omit the format. Explicit false
 // remains meaningful for restoring a historical rich-text version.
 for(const format of [undefined,false,true]){
  const current=await read()
  const saved=await (await context.request.post(`${baseURL}/api2/save`,{headers,data:{noteId,title:current.Note.Title,content:edited,tags:'',usn:current.Note.Usn,...(format===undefined?{}:{isMarkdown:format})}})).json()
  expect(saved.Note.IsMarkdown).toBe(format!==false)
  expect(saved.Content).toBe(edited)
 }
 await page.reload()
 await page.getByRole('button',{name:'编辑',exact:true}).click()
 await expect(page.locator('textarea.content-editor')).toHaveValue(edited)

 await showImport()
 dialog=page.getByRole('dialog',{name:'导入 Markdown',exact:true})
 await dialog.locator('input[type=file]').setInputFiles({name:'批量导入测试.zip',mimeType:'application/zip',buffer:Buffer.from(zipSync({'root.md':strToU8('root body'),'一级/二级/child.md':strToU8('# 子目录正文'),'ignore.png':new Uint8Array([255])}))})
 await expect(dialog).toContainText('将导入 2 篇笔记、3 个新笔记本；跳过 1 个文件。')
 await dialog.getByRole('button',{name:'开始导入',exact:true}).click()
 await expect(dialog).toContainText('导入完成')
 await dialog.getByRole('button',{name:'关闭',exact:true}).click()
 const boot=await (await context.request.get(`${baseURL}/api2/bootstrap`,{headers})).json()
 const parent=boot.Notebooks.find(b=>b.NotebookId===target.NotebookId)
 const root=parent.Subs.find(b=>b.Title==='批量导入测试')
 expect(root).toBeTruthy();expect(root.NumberNotes).toBe(1)
 const first=root.Subs.find(b=>b.Title==='一级'),second=first.Subs.find(b=>b.Title==='二级')
 expect(second.NumberNotes).toBe(1)
 const list=await (await context.request.post(`${baseURL}/api2/notes`,{headers,data:{notebookId:second.NotebookId,page:1}})).json()
 expect(list).toHaveLength(1);expect(list[0].Title).toBe('child')
 const doc=await (await context.request.post(`${baseURL}/api2/document`,{headers,data:{noteId:list[0].NoteId}})).json()
 expect(doc.Content).toBe('# 子目录正文');expect(doc.Note.IsMarkdown).toBe(true)
 await expect(page.locator('.notebook-row').filter({hasText:'二级'})).toBeVisible()
 console.log('Disposable PostgreSQL: MD/ZIP import, persistence, hierarchy, counts and format-preserving save passed')
 await context.close()
}finally{await browser?.close();await server.close()}
