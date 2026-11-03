import {chromium, expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
import {existsSync} from 'node:fs'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const base=`http://127.0.0.1:${server.httpServer.address().port}`
const executablePath=['/usr/bin/google-chrome','/usr/bin/chromium'].find(existsSync)
const browser=await chromium.launch({headless:true,...(executablePath?{executablePath}:{})})
try {
 for(const desktop of [false,true]) for(const lang of ['zh-CN','en-US']) {
  const page=await browser.newPage({locale:lang,viewport:{width:1280,height:1000}})
  await page.addInitScript(lang=>localStorage.setItem('gemsnote:language',lang),lang)
  const members=[], submitted=[]
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname
   let data={Ok:true}
   if(path==='/api2/bootstrap')data={User:{UserId:'owner',Username:'owner',Email:'owner@example.test'},Desktop:desktop,IsAdmin:false,Notebooks:[],Tags:[],SharedNotebooks:{}}
   if(path==='/api2/groups')data=[{GroupId:'mine',UserId:'owner',Title:'My group',Users:members},{GroupId:'others',UserId:'other',Title:'Other group',Users:[]}]
   if(path==='/api2/member/group/addUser'){
    expect(route.request().method()).toBe('POST')
    const fields=new URLSearchParams(route.request().postData())
    const identity=fields.get('email');submitted.push(identity)
    expect(fields.get('groupId')).toBe('mine')
    if(identity==='alice@example.test')data={Ok:false,Msg:'userExistsInGroup'}
    else members.push({UserId:`member-${members.length}`,Username:identity.includes('@')?'Second':identity,Email:identity.includes('@')?identity:'alice@example.test'})
   }
   await route.fulfill({json:data})
  })
  await page.goto(base+'/user/account')
  const group=page.locator('.group').filter({has:page.getByRole('heading',{name:'My group'})})
  const input=group.locator('input')
  const button=group.getByRole('button',{name:lang==='zh-CN'?'添加用户':'Add user',exact:true})
  await expect(input).toHaveAttribute('type','text')
  await expect(input).toHaveAttribute('placeholder',lang==='zh-CN'?'已有用户的用户名或邮箱':'Existing username or email')
  await input.fill('  Alice  ');await button.click()
  await expect.poll(()=>submitted).toEqual(['Alice'])
  await expect(group.locator('.members li')).toHaveCount(1)
  await expect(input).toHaveValue('')
  await input.fill('alice@example.test');await button.click()
  await expect(page.getByRole('status')).toHaveText(lang==='zh-CN'?'该用户已在此分组中。':'This user is already in the group.')
  await expect(group.locator('.members li')).toHaveCount(1)
  await input.fill(' second@example.test ');await button.click()
  await expect(group.locator('.members li')).toHaveCount(2)
  await input.fill('   ');await button.click()
  await expect(page.getByRole('status')).toHaveText(lang==='zh-CN'?'请输入用户名或邮箱':'Enter a username or email')
  expect(submitted).toEqual(['Alice','alice@example.test','second@example.test'])
  await expect(page.locator('.group').filter({has:page.getByRole('heading',{name:'Other group'})}).locator('input')).toHaveCount(0)
  await page.close()
  console.log(`PASS group identity form, duplicates, counts, permissions; desktop=${desktop}; lang=${lang}`)
 }
} finally {await browser.close();await new Promise(resolve=>server.httpServer.close(resolve))}
