import {chromium, expect} from '@playwright/test'
import {existsSync, readFileSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'
const server = await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const baseURL = `http://127.0.0.1:${server.httpServer.address().port}`
const executablePath = [process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/google-chrome','/usr/bin/chromium'].find(p=>p&&existsSync(p))
let browser
const diagrams = [
 'flowchart LR\n A[记录想法] --> B[整理笔记]',
 'sequenceDiagram\n Alice->>Bob: Hello',
 'classDiagram\n Animal <|-- Duck',
 'stateDiagram-v2\n [*] --> Ready',
 'gantt\n title 计划\n dateFormat YYYY-MM-DD\n section 工作\n 实现 :2026-11-01, 2d',
]
try {
 browser = await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 // Run the exact mobile asset as a local file too: no HTTP server or CDN.
 const mobileFile = fileURLToPath(new URL('../../mobile-app/assets/mermaid/index.html',import.meta.url))
 const webFile = fileURLToPath(new URL('../dist/assets/mermaid/index.html',import.meta.url))
 if (existsSync(mobileFile)) expect(readFileSync(mobileFile,'utf8')).toBe(readFileSync(webFile,'utf8'))
 const renderer = await browser.newPage()
 const external=[]
 renderer.on('request',r=>{if(/^https?:/.test(r.url()))external.push(r.url())})
 await renderer.addInitScript(()=>{window.results=[];window.MermaidResult={postMessage:s=>window.results.push(JSON.parse(s))}})
 await renderer.goto(`file://${existsSync(mobileFile)?mobileFile:webFile}`)
 await expect.poll(()=>renderer.evaluate(()=>window.results.some(r=>r.ready))).toBe(true)
 for (const [id,source] of diagrams.entries()) {
  await renderer.evaluate(({id,source})=>window.renderGemsnoteMermaid({channel:'gemsnote-mermaid',id,source,format:'png'}),{id,source})
  await expect.poll(()=>renderer.evaluate(id=>window.results.find(r=>r.id===id)?.png?.startsWith('data:image/png;base64,'),id)).toBe(true)
 }
 expect(external).toEqual([])
 await renderer.close()
 for (const desktop of [false,true]) {
  const page=await browser.newPage({viewport:{width:1280,height:900}})
  const errors=[],foreign=[]
  page.on('pageerror',e=>errors.push(e.message))
  page.on('request',r=>{if(/^https?:/.test(r.url())&&!r.url().startsWith(baseURL))foreign.push(r.url())})
  await page.addInitScript(()=>{if(window===window.top&&!localStorage.getItem('gemsnote:language'))localStorage.setItem('gemsnote:language','zh-CN')})
  let content='# 图表测试\n\n'+diagrams.map(s=>'```mermaid\n'+s+'\n```').join('\n\n')+'\n\n```js\nconst ordinary = true\n```'
  const note={NoteId:'507f1f77bcf86cd799439013',NotebookId:'507f1f77bcf86cd799439012',UserId:'507f1f77bcf86cd799439011',Title:'Mermaid',IsMarkdown:true,Tags:[]}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   const data=path==='/api2/bootstrap'?{Desktop:desktop,User:{UserId:note.UserId},Notebooks:[{NotebookId:note.NotebookId,Title:'Notebook',NumberNotes:1}],SharedNotebooks:{},Tags:[],TotalNotes:1}
    :path==='/api2/notes'?[note]:path==='/api2/document'?{Note:note,Content:content,Writable:true}:path==='/api2/attachments'?{List:[]}:path==='/api2/save'?{Note:note,Content:route.request().postDataJSON().content,Writable:true}:{Ok:true}
   await route.fulfill({json:data})
  })
  await page.goto(`${baseURL}/note/${note.NoteId}`)
  await expect(page.locator('.diagram-image img')).toHaveCount(5,{timeout:15000})
  expect(decodeURIComponent(await page.locator('.diagram-image img').first().getAttribute('src'))).toContain('记录想法')
  if (!desktop) await page.screenshot({path:'/tmp/gemsnote-mermaid-web.png'})
  expect(await page.locator('.diagram-image img').evaluateAll(images=>images.every(i=>i.complete&&i.naturalWidth>0))).toBe(true)
  await expect(page.locator('code.language-js')).toContainText('ordinary')
  // Shared render host has no same-origin access, no Wails bridge.
  await expect(page.locator('iframe[title="Mermaid renderer"]')).toHaveAttribute('sandbox','allow-scripts')
  await page.getByRole('button',{name:'放大图表'}).first().click()
  await expect(page.getByRole('dialog')).toBeVisible()
  await page.getByRole('button',{name:'关闭',exact:true}).click()
  await page.getByRole('button',{name:'编辑',exact:true}).click()
  await expect(page.locator('textarea.content-editor')).toHaveValue(content)
  const invalid='```mermaid\nnot valid mermaid!\n```\n\n```mermaid\n%%{init: {securityLevel: "loose"}}%%\nflowchart LR\nA-->B\n```\n\n'+
    '```mermaid\nflowchart LR\n A["<img src=https://evil.invalid/a onerror=alert(1)>"] --> B\n```'
  await page.locator('textarea.content-editor').fill(invalid)
  await page.getByRole('button',{name:'预览',exact:true}).click()
  await expect(page.getByText('图表无法渲染，请检查语法或图表大小。',{exact:true})).toHaveCount(2)
  await expect(page.locator('.mermaid-diagram')).toHaveCount(3)
  await expect(page.locator('.diagram-image img')).toHaveCount(1)
  await expect(page.locator('.mermaid-diagram pre').first()).toContainText('not valid mermaid!')
  expect(foreign).toEqual([])
  expect(errors).toEqual([])
  await page.evaluate(()=>localStorage.setItem('gemsnote:language','en-US'))
  await page.reload()
  await expect(page.getByRole('button',{name:'Enlarge diagram'})).toHaveCount(5)
  await page.close()
  console.log(`${desktop?'Desktop':'Web'}: five diagrams, Chinese, source, zoom, editing, errors and isolated rendering passed`)
 }
} finally {await browser?.close();await server.close()}
