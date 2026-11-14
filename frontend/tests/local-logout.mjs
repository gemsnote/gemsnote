import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/google-chrome'})
try{
 for(const desktop of [false,true])for(const lang of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}}),calls=[]
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),lang)
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname;calls.push(path)
   let data={Running:false}
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'u',Username:'u'},Notebooks:[],Tags:[],PendingChanges:true}
   if(path==='/api2/notes')data=[]
   if(path==='/api2/logout')data={Ok:true}
   await route.fulfill({json:data})
  })
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note`)
  await page.getByRole('button',{name:lang==='zh-CN'?'用户菜单':'User menu',exact:true}).click()
  await page.evaluate(()=>{sessionStorage.setItem('test-local-flush-count','0');window.__gemsnoteBeforeLogout=async()=>{sessionStorage.setItem('test-local-flush-count','1');return true}})
  await page.getByRole('menuitem',{name:lang==='zh-CN'?'退出':'Sign out',exact:true}).click()
  await expect(page).toHaveURL(/\/login$/)
  expect(await page.evaluate(()=>sessionStorage.getItem('test-local-flush-count'))).toBe('1')
  expect(calls.filter(path=>path==='/api2/logout')).toHaveLength(1)
  expect(calls.some(path=>/\/sync$|\/fullSync$|getSyncState/.test(path))).toBe(false)
  await expect(page.getByRole('dialog')).toHaveCount(0)
  await page.close();console.log(`${desktop?'Desktop':'Web'} ${lang}: logout redirects without sync or confirmation`)
 }
}finally{await browser.close();await server.close()}
