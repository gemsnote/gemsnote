import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:process.env.GEMSNOTE_CHROMIUM_PATH||'/usr/bin/google-chrome'})
try {
 for(const desktop of [false,true]) {
  const page=await browser.newPage({viewport:{width:1280,height:800}})
  await page.addInitScript(()=>localStorage.setItem('gemsnote:language','zh-CN'))
  let tags=['existing'],writable=true,manyTags=false
  const saves=[]
  const doc=()=>({Note:{NoteId:'n',NotebookId:'book',UserId:'u',Title:'Note',Tags:tags,IsMarkdown:true},Content:'Body',Writable:writable})
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'u',Username:'user'},Notebooks:[{NotebookId:'book',Title:'Book'}],Tags:[...[...new Set([...tags,'alpha','alpine','中文'])].map(Tag=>({Tag,Count:1})),...(manyTags?Array.from({length:12},(_,i)=>({Tag:`popular${i}`,Count:30-i})):[])]}
   else if(path==='/api2/notes')data=[doc().Note]
   else if(path==='/api2/document')data=doc()
   else if(path==='/api2/attachments')data={List:[]}
   else if(path==='/api2/save'){const body=route.request().postDataJSON();saves.push(body);tags=body.tags?body.tags.split(','):[];data=doc()}
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note/n`)
  const input=page.getByRole('combobox',{name:'笔记标签'})
  const chips=page.locator('.tag-chip')
  await expect(chips).toHaveCount(1)
  await input.fill('al')
  await expect(page.getByRole('option')).toHaveCount(2)
  await input.press('ArrowDown')
  await input.press('Enter')
  await expect(chips).toHaveCount(2)
  await expect(input).toHaveValue('')
  await input.fill('alpha')
  await input.press('Enter')
  await expect(chips).toHaveCount(2)
  await input.fill('中文')
  await input.dispatchEvent('compositionstart')
  await input.dispatchEvent('keydown',{key:'Enter',code:'Enter',isComposing:true})
  await expect(chips).toHaveCount(2)
  await input.dispatchEvent('compositionend',{data:'中文'})
  await page.getByRole('option',{name:'中文',exact:true}).click()
  await expect(chips).toHaveCount(3)
  await input.fill('one,two，three')
  await input.press('Enter')
  await expect(chips).toHaveCount(6)
  await input.fill('blurTag')
  await page.getByRole('button',{name:'编辑',exact:true}).click()
  await expect(chips).toHaveCount(7)
  await page.getByRole('button',{name:'删除标签 existing',exact:true}).click()
  await expect(chips).toHaveCount(6)
  await input.focus()
  await input.press('Control+s')
  await expect.poll(()=>saves.at(-1)?.tags).toBe('alpha,中文,one,two,three,blurTag')
  const sidebar=page.locator('.sidebar-tags')
  await expect(sidebar.getByRole('button',{name:'#blurTag 1',exact:true})).toBeVisible()
  await expect(sidebar.getByRole('button',{name:'#existing 1',exact:true})).toHaveCount(0)
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('BLUR')
  await expect(sidebar.getByRole('button')).toHaveCount(1)
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('missing')
  await expect(sidebar).toContainText('没有匹配的标签')
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('')
  await input.focus()
  // Candidate escape does not delete the draft; blur still confirms it.
  await input.fill('al')
  await input.press('Escape')
  await expect(page.getByRole('listbox')).toHaveCount(0)
  await input.fill('')
  for(let i=0;i<6;i++)await input.press('Backspace')
  await input.press('Control+s')
  await expect.poll(()=>saves.at(-1)?.tags).toBe('')
  tags=['readonly'];writable=false;manyTags=true
  await page.reload()
  await expect(page.locator('.tag-chip')).toHaveText('readonly')
  await expect(input).toHaveCount(0)
  await expect(page.getByRole('button',{name:/删除标签/})).toHaveCount(0)
  await expect(sidebar.getByRole('button')).toHaveCount(10)
  await expect(sidebar.getByRole('button').first()).toContainText('#popular0')
  await expect(sidebar.getByRole('button').last()).toContainText('#popular9')
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('popular')
  await expect(sidebar.getByRole('button')).toHaveCount(12)
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('readonly')
  await expect(sidebar.getByRole('button')).toHaveCount(1)
  await expect(sidebar.getByRole('button')).toContainText('#readonly')
  await page.getByRole('searchbox',{name:'过滤标签'}).fill('')
  await expect(sidebar.getByRole('button')).toHaveCount(10)
  await page.close()
 }
 console.log('Tag chips, IME, candidate selection, CSV saving and readonly checks passed for Web/desktop')
} finally {await browser.close();await server.httpServer.close()}
