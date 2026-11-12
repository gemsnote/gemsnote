import {chromium,expect} from '@playwright/test'
import {existsSync} from 'node:fs'
import {fileURLToPath} from 'node:url'
import {preview} from 'vite'

const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const baseURL=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=[process.env.GEMSNOTE_CHROMIUM_PATH,'/usr/bin/chromium','/usr/bin/google-chrome'].find(p=>p&&existsSync(p))
let browser
try{
 browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
 for(const desktop of [false,true])for(const lang of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}}),errors=[],calls=[]
  page.setDefaultTimeout(10000)
  const text=(zh,en)=>lang==='zh-CN'?zh:en
  page.on('pageerror',error=>errors.push(error.message))
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),lang)
  let notes=['Alpha','Beta','Gamma','Delta','Epsilon'].map((Title,index)=>({NoteId:String(index),Title,UserId:'me',NotebookId:'source',IsMarkdown:true,UpdatedTime:`2026-11-${String(20-index).padStart(2,'0')}T10:00:00Z`,IsTrash:false}))
  const received={NoteId:'received',Title:'Received note',UserId:'author',NotebookId:'source',IsMarkdown:true,IsShared:true,Perm:0,UpdatedTime:'2026-11-01T10:00:00Z'}
  const shareCalls=[],shared=new Set();let failShare=true,delayDocument=true,releaseDocument=()=>{}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname,body=route.request().postData()||''
   const params=body.startsWith('{')?JSON.parse(body):Object.fromEntries(new URLSearchParams(body))
   const ids=Object.entries(params).filter(([key])=>key.startsWith('noteIds[')).map(([,id])=>id)
   let data=true
   switch(path){
    case '/api2/bootstrap':{
     const count=id=>notes.filter(note=>!note.IsTrash&&note.NotebookId===id).length
     data={Desktop:desktop,User:{UserId:'me',Username:'me'},Notebooks:[{NotebookId:'source',Title:'Source',NumberNotes:count('source')},{NotebookId:'target',Title:'Target',NumberNotes:count('target'),Subs:[{NotebookId:'child',ParentNotebookId:'target',Title:'Child',NumberNotes:count('child')}]}],Tags:[],TotalNotes:notes.filter(note=>!note.IsTrash).length,Sharing:{Suggestions:[{UserId:'reader',Username:'Reader',NoteIds:[]}],Recipients:shared.size?[{UserId:'reader',Username:'Reader',NoteIds:[...shared]}]:[],Placements:{received:received.NotebookId},PlacedShares:[received]}}
     break
    }
    case '/api2/notes':data=notes.filter(note=>!!note.IsTrash===!!params.trash&&(!params.notebookId||note.NotebookId===params.notebookId));break
    case '/api2/document':{
     if(delayDocument&&params.noteId==='0'){delayDocument=false;await new Promise(resolve=>releaseDocument=resolve)}
     const note=notes.find(note=>note.NoteId===params.noteId)||received
     data={Note:note,Content:`# Body of ${note.Title}`,Writable:note.UserId==='me'};break
    }
    case '/api2/attachments':data={List:[]};break
    case '/api2/web/downloadStatus':data={Running:false};break
    case '/api2/note/moveNote':
     calls.push({path,ids,notebookId:params.notebookId})
     notes.filter(note=>ids.includes(note.NoteId)).forEach(note=>note.NotebookId=params.notebookId);break
    case '/api2/shared/move':calls.push({path,...params});received.NotebookId=params.notebookId;break
    case '/api2/note/deleteNote':calls.push({path,ids});notes.filter(note=>ids.includes(note.NoteId)).forEach(note=>note.IsTrash=true);break
    case '/api2/note/deleteTrash':calls.push({path,...params});notes=notes.filter(note=>note.NoteId!==params.noteId);break
    case '/api2/share/addShareNote':{
     const id=params.noteId;shareCalls.push(id)
     const ok=!(id==='2'&&failShare)
     if(ok)shared.add(id)
     data={Reader:{Ok:ok,Msg:ok?'':'offline'}};break
    }
    default:errors.push(`Unexpected API request: ${path}`)
   }
   await route.fulfill({json:data})
  })
  await page.goto(`${baseURL}/note`)
  const rows=page.locator('.note-item'),selected=page.locator('.note-item[aria-pressed="true"] strong')
  const row=title=>rows.filter({has:page.locator('strong',{hasText:new RegExp(`^${title}$`)})})
  const menu=page.getByRole('menu',{name:text('所选笔记操作','Selected note actions')})
  const dialog=page.getByRole('dialog')
  const action=name=>menu.getByRole('menuitem',{name,exact:true})
  await expect(rows).toHaveCount(6)

  // Preserve notebook operations while changing their visible order and grouping.
  await page.locator('.notebook-row').filter({hasText:'Source'}).hover()
  await page.getByRole('button',{name:text('Source 菜单','Source menu'),exact:true}).click()
  await expect(page.locator('.notebook-menu [role="menuitem"]')).toHaveText(lang==='zh-CN'?['新增','重命名','移动','删除','导入','导出']:['New','Rename','Move','Delete','Import','Export'])
  expect(await page.locator('.notebook-menu').evaluate(el=>[...el.children].findIndex(child=>child.getAttribute('role')==='separator'))).toBe(4)
  await page.keyboard.press('Escape')

  // A late document response must not overwrite a selection made while loading.
  const pendingDocument=page.waitForRequest(req=>req.url().endsWith('/api2/document'))
  await row('Alpha').click();await pendingDocument
  await row('Gamma').click({modifiers:['Control']})
  await expect(selected).toHaveText(['Alpha','Gamma'])
  await expect(menu).toBeVisible()
  releaseDocument()
  await expect(page.locator('.toolbar-title')).toHaveValue('Alpha')
  await expect(selected).toHaveText(['Alpha','Gamma'])
  await row('Gamma').click({modifiers:['Control']})
  await expect(selected).toHaveText(['Alpha'])
  await row('Delta').click({modifiers:['Shift']})
  await expect(selected).toHaveText(['Gamma','Delta'])
  await row('Alpha').click({modifiers:['Control','Shift']})
  await expect(selected).toHaveText(['Alpha','Beta','Gamma','Delta'])
  await page.keyboard.press('Escape')
  await expect(selected).toHaveCount(0)

  // Meta is the macOS counterpart of Ctrl. Right-click preserves a group.
  await row('Alpha').click()
  await row('Gamma').click({modifiers:['Meta']})
  await row('Alpha').click({button:'right'})
  await expect(selected).toHaveText(['Alpha','Gamma'])
  await action(text('共享','Share')).click()
  const recipient=dialog.getByRole('combobox',{name:text('用户名或邮箱','Username or email')})
  await recipient.fill('Rea')
  await dialog.getByRole('option',{name:'Reader',exact:true}).click()
  await dialog.getByRole('button',{name:text('共享','Share'),exact:true}).click()
  await expect(dialog.getByRole('alert')).toContainText('1/2')
  await expect(recipient).toBeDisabled()
  expect(shareCalls).toEqual(['0','2'])
  failShare=false
  await dialog.getByRole('button',{name:text('重试','Retry'),exact:true}).click()
  await expect(dialog).toHaveCount(0)
  expect(shareCalls).toEqual(['0','2','2'])

  // Shift follows the visual sort order rather than the original API array.
  await page.locator('.note-sort-menu summary').click()
  await page.getByRole('menuitemradio',{name:text('标题：降序','Title: Z–A'),exact:true}).click()
  await row('Gamma').click()
  await row('Delta').click({modifiers:['Shift']})
  await expect(selected).toHaveText(['Gamma','Epsilon','Delta'])
  await page.getByRole('button',{name:/^.*Target\s*0$/}).click()
  await expect(selected).toHaveCount(0)
  await page.getByRole('button',{name:/^.*Source\s*6$/}).click()
  await expect(rows).toHaveCount(6)
  await row('Alpha').click()
  await row('Beta').click({modifiers:['Control']})
  await action(text('移动','Move')).click()
  await dialog.getByRole('button',{name:'Child',exact:true}).click()
  await dialog.getByRole('button',{name:text('移动','Move'),exact:true}).click()
  await expect(dialog).toHaveCount(0)
  await expect(rows).toHaveCount(4)
  expect(calls).toContainEqual({path:'/api2/note/moveNote',ids:['1','0'],notebookId:'child'})
  await expect(page.locator('.notebook-select').filter({hasText:'Source'}).locator('small')).toHaveText('4')

  // Received shares retain only the recipient's notebook-placement operation.
  await row('Received note').click({button:'right'})
  await expect(action(text('共享','Share'))).toBeDisabled()
  await expect(action(text('删除','Delete'))).toBeDisabled()
  await action(text('移动','Move')).click()
  await dialog.getByRole('button',{name:'Child',exact:true}).click()
  await dialog.getByRole('button',{name:text('移动','Move'),exact:true}).click()
  await expect(dialog).toHaveCount(0)
  await expect(rows).toHaveCount(3)
  expect(calls).toContainEqual({path:'/api2/shared/move',noteId:'received',notebookId:'child'})
  await expect(page.locator('.notebook-select').filter({hasText:'Source'}).locator('small')).toHaveText('3')

  // Delete always confirms, and cancelling never mutates the selected notes.
  await row('Gamma').click()
  await row('Delta').click({modifiers:['Control']})
  await action(text('删除','Delete')).click()
  await dialog.getByRole('button',{name:text('取消','Cancel'),exact:true}).click()
  expect(calls.filter(call=>call.path==='/api2/note/deleteNote')).toHaveLength(0)
  await row('Gamma').click({button:'right'})
  await action(text('删除','Delete')).click()
  await dialog.getByRole('button',{name:text('删除','Delete'),exact:true}).click()
  await expect(dialog).toHaveCount(0)
  await expect(rows).toHaveCount(1)
  await expect(page.locator('.toolbar-title')).toHaveCount(0)
  expect(calls).toContainEqual({path:'/api2/note/deleteNote',ids:['2','3']})
  await expect(page.locator('.notebook-select').filter({hasText:'Source'}).locator('small')).toHaveText('1')
  expect(errors).toEqual([])
  await page.close()
  console.log(`${desktop?'Desktop':'Web'} ${lang}: Ctrl/Meta/Shift selection, menu order, move/counts, shared permissions, share retries and confirmed deletion passed`)
 }
}finally{await browser?.close();await server.close()}
