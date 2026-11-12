import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:process.env.GEMSNOTE_CHROMIUM_PATH||'/usr/bin/google-chrome'})
try {
 for(const desktop of [false,true])for(const lang of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}})
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),lang)
  let moved=false,copied=false
  const shares=[]
  const requests=[]
  const own={NoteId:'own',NotebookId:'mine',UserId:'me',Title:'My rich text',IsMarkdown:false,UpdatedTime:'2026-10-30T10:00:00Z'}
  let catalogue=[own]
  const outgoing=['own'],notePages=[]
  const received={NoteId:'shared',NotebookId:'remote',UserId:'author',Title:'Shared Markdown',IsMarkdown:true,UpdatedTime:'2026-10-30T09:00:00Z',Perm:0,IsShared:true}
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   const body=route.request().postData()||''
   let data={Ok:true}
   switch(path){
    case '/api2/bootstrap':data={Desktop:desktop,User:{UserId:'me',Username:'me'},Notebooks:[{NotebookId:'mine',Title:'My notebook',NumberNotes:copied?2:1}],Tags:[],TotalNotes:copied?2:1,SharedNotebooks:{author:[{IsDefault:true,Title:'Author shares'}]},SharedUsers:desktop?undefined:[{UserId:'author',Username:'Original author'}],Sharing:{Owners:desktop?[{UserId:'author',Username:'Original author'}]:[],Suggestions:[{UserId:'past',Username:'past-reader',NoteIds:[]}],Recipients:[{UserId:'reader',Username:'Reader Alice',NoteIds:outgoing}],Placements:moved?{shared:'mine'}:{},PlacedShares:moved?[{...received,NotebookId:'mine'}]:[]}};break
    case '/api2/notes':{const number=JSON.parse(body).page||1;notePages.push(number);data=catalogue.slice((number-1)*100,number*100);break}
    case '/api2/share/listShareNotes':data=[received];break
    case '/api2/document':{
     const id=JSON.parse(body).noteId
     data={Note:id==='own'?own:received,Content:id==='own'?'<p>own body</p>':'# Read only content',Writable:id==='own'};break
    }
    case '/api2/attachments':data={List:[]};break
    case '/api2/web/shareMembers':data={Users:[]};break
    case '/api2/share/addShareNote':{
     const params=new URLSearchParams(body);const identity=params.get('emails[0]');shares.push(identity);data={[identity]:{Ok:true}};break
    }

    case '/api2/shared/move':requests.push({path,body:JSON.parse(body)});moved=true;break
    case '/api2/note/copySharedNote':requests.push({path,body:Object.fromEntries(new URLSearchParams(body))});copied=true;break
    case '/api2/web/downloadStatus':data={Running:false};break
   }
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note/own`)
  await expect(page.locator('.note-info-line')).toContainText('Reader Alice')
  await expect(page.locator('.note-format')).toHaveText('RT')
  await page.locator('.toolbar-actions').getByRole('button',{name:lang==='zh-CN'?'共享':'Share',exact:true}).click()
  const recipient=page.getByRole('combobox',{name:lang==='zh-CN'?'用户名或邮箱':'Username or email'})
  await recipient.fill('past')
  await expect(page.getByRole('option',{name:'past-reader',exact:true})).toBeVisible()
  await recipient.press('ArrowDown');await recipient.press('Enter')
  await expect(recipient).toHaveValue('past-reader');expect(shares).toEqual([])
  await page.getByRole('dialog').getByRole('button',{name:lang==='zh-CN'?'添加共享':'Share',exact:true}).click()
  await expect.poll(()=>shares.length).toBe(1);expect(shares[0]).toBe('past-reader')
  await recipient.fill('Reader');await page.getByRole('option',{name:'Reader Alice',exact:true}).click()
  await expect(recipient).toHaveValue('Reader Alice')
  await recipient.fill(' reader@example.test ')
  await page.getByRole('dialog').getByRole('button',{name:lang==='zh-CN'?'添加共享':'Share',exact:true}).click()
  await expect.poll(()=>shares.length).toBe(2);expect(shares[1]).toBe('reader@example.test')
  await page.getByRole('dialog').getByRole('button',{name:lang==='zh-CN'?'关闭':'Close',exact:true}).click()

  // Match the live regression: two outgoing shares at positions 5 and 232.
  // Both must appear even though their metadata is fetched on different pages.
  outgoing.push('old-own')
  catalogue=Array.from({length:251},(_,i)=>i===4?own:{...own,NoteId:i===231?'old-own':`filler-${i}`,Title:i===231?'Older shared note':`Filler ${i}`})
  await page.reload()
  await expect(page.locator('.note-item')).toHaveCount(251)
  notePages.length=0
  await page.getByRole('button',{name:/Reader Alice/}).click()
  await expect(page.locator('.note-list h2')).toHaveText('Reader Alice')
  await expect(page.locator('.note-item')).toHaveCount(2)
  await expect(page.locator('.note-list')).toContainText('Older shared note')
  expect(notePages).toEqual([1,2,3])
  catalogue=[own]
  await page.getByRole('button',{name:'♧ Author shares',exact:true}).click()
  await expect(page.locator('.note-item')).toHaveCount(1)
  await expect(page.locator('.note-format')).toHaveText('MD')
  await page.locator('.note-item').click()
  await expect(page.locator('.rendered')).toContainText('Read only content')
  await expect(page.locator('.note-info-line')).toContainText((lang==='zh-CN'?'共享自：':'Shared by: ')+'Original author')
  await expect(page.getByRole('button',{name:lang==='zh-CN'?'编辑':'Edit',exact:true})).toHaveCount(0)
  // Locate via the toolbar to keep translated button text assertions explicit.
  const action=page.locator('.toolbar-actions button').filter({hasText:lang==='zh-CN'?'移动/复制':/Move.*copy/i})
  await action.click()
  const dialog=page.getByRole('dialog')
  await dialog.getByRole('button',{name:'My notebook',exact:true}).click()
  await dialog.getByRole('button',{name:lang==='zh-CN'?'移动':'Move',exact:true}).click()
  await expect(dialog).toHaveCount(0)
  expect(requests[0]).toEqual({path:'/api2/shared/move',body:{noteId:'shared',notebookId:'mine'}})
  await expect(page.locator('.note-info-line')).toContainText('My notebook')
  await expect(page.locator('.note-info-line')).toContainText('Original author')
  await expect(page.locator('.notebook-select small')).toHaveText('2')
  await page.locator('.notebook-select').click()
  await expect(page.locator('.note-item')).toHaveCount(2)
  await action.click()
  await dialog.getByRole('button',{name:'My notebook',exact:true}).click()
  await dialog.getByRole('button',{name:lang==='zh-CN'?'复制到我的笔记本':/Copy to my notebook/i}).click()
  await expect(dialog).toHaveCount(0)
  expect(requests[1].body['noteIds[0]']).toBe('shared')
  expect(requests[1].body.fromUserId).toBe('author')
  await page.close()
 }
 console.log('PASS Web/desktop shared navigation, readonly actions, placement/counts, copy, RT/MD, Chinese/English')
}finally{await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
