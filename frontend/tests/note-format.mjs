import {chromium,expect} from '@playwright/test'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'

const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const baseURL=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [false,true])for(const markdown of [true,false]){
  const page=await browser.newPage({viewport:{width:1280,height:800}}),errors=[],saves=[]
  page.on('pageerror',e=>errors.push(e.message))
  await page.addInitScript(()=>{
   localStorage.setItem('gemsnote:language','zh-CN')
   // Keep the rich-text fallback textarea usable without loading TinyMCE.
   window.tinymce={init(options){let content='';options.setup({initialized:true,on(){},setContent(v){content=v},getContent(){return content},remove(){}})}}
  })
  let doc={Note:{NoteId:'507f1f77bcf86cd799439013',NotebookId:'507f1f77bcf86cd799439012',UserId:'507f1f77bcf86cd799439011',Title:'格式回归',IsMarkdown:markdown,Tags:[],Usn:7},Content:markdown?'# 原始\n\n正文':'<p>原始正文</p>',Writable:true}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data
   switch(path){
    case '/api2/bootstrap':data={Desktop:desktop,User:{UserId:doc.Note.UserId},Notebooks:[{NotebookId:doc.Note.NotebookId,Title:'Notebook',NumberNotes:1}],SharedNotebooks:{},Tags:[],TotalNotes:1};break
    case '/api2/notes':data=[doc.Note];break
    case '/api2/document':data=doc;break
    case '/api2/attachments':data={List:[]};break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/save':{
     const body=route.request().postDataJSON();saves.push(body)
     // Model the original server behavior: omitting the bool changes it to false.
     doc={Note:{...doc.Note,Title:body.title,IsMarkdown:body.isMarkdown===true,Usn:doc.Note.Usn+1},Content:body.content,Writable:true};data=doc;break
    }
    default:errors.push(`Unexpected API request: ${path}`);data={Ok:false,Msg:'unexpected'}
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${baseURL}/note/${doc.Note.NoteId}`)
  await expect(page.locator('.rendered')).toBeVisible()
  await page.getByRole('button',{name:'编辑',exact:true}).click()
  const editor=page.locator(markdown?'textarea.content-editor':'.rich-editor textarea')
  const changed=markdown?'# 修改后\n\n**仍是 Markdown**':'<p><strong>仍是富文本</strong></p>'
  await editor.fill(changed)
  await page.getByRole('button',{name:'预览',exact:true}).click()
  await expect.poll(()=>saves.length).toBe(1)
  expect(saves[0].isMarkdown).toBe(markdown)
  await expect(page.locator('.note-info-line')).toContainText(markdown?'格式：MD':'格式：RT')
  await page.getByRole('button',{name:'编辑',exact:true}).click()
  await expect(editor).toHaveValue(changed)
  await expect(page.locator('.rich-editor')).toHaveCount(markdown?0:1)
  await page.getByRole('textbox',{name:'笔记标题',exact:true}).fill('修改标题，格式保持')
  await expect.poll(()=>saves.length).toBe(2)
  expect(saves[1].isMarkdown).toBe(markdown)
  expect(doc.Content).toBe(changed)
  await page.reload()
  await expect(page.locator('.note-info-line')).toContainText(markdown?'格式：MD':'格式：RT')
  await page.getByRole('button',{name:'编辑',exact:true}).click()
  await expect(editor).toHaveValue(changed)
  expect(errors).toEqual([])
  await page.close()
  console.log(`${desktop?'Desktop':'Web'} ${markdown?'Markdown':'rich text'}: body/title autosave, preview/edit and reload preserve format`)
 }
}finally{await browser?.close();await server.close()}
