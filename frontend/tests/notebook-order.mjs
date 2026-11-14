import {chromium,expect} from '@playwright/test'
import {preview} from 'vite'
import {fileURLToPath} from 'node:url'
const server=await preview({root:fileURLToPath(new URL('..',import.meta.url)),configFile:false,preview:{host:'127.0.0.1',port:0}})
const browser=await chromium.launch({headless:true,executablePath:'/usr/bin/google-chrome'})
try{
 for(const desktop of [false,true])for(const language of ['zh-CN','en-US']){
  const page=await browser.newPage({viewport:{width:1440,height:900}});page.setDefaultTimeout(10000)
  await page.addInitScript(l=>localStorage.setItem('gemsnote:language',l),language)
  const child={NotebookId:'child',ParentNotebookId:'b',Title:'Child',Seq:0}
  let books=[{NotebookId:'a',Title:'Alpha',Seq:0},{NotebookId:'b',Title:'Beta',Seq:0,Subs:[child]},{NotebookId:'c',Title:'Gamma',Seq:0}],calls=[]
  await page.route('**/api2/**',async route=>{
   const path=new URL(route.request().url()).pathname;let data={Running:false}
   if(path==='/api2/bootstrap')data={Desktop:desktop,User:{UserId:'u'},Notebooks:books,Tags:[],TotalNotes:0}
   if(path==='/api2/notes')data=[]
   if(path==='/api2/notebook/reorder'){
    expect(route.request().headers()['content-type']).toBe('application/json')
    const p=JSON.parse(route.request().postData());calls.push(p.notebookIds)
    books=books.map(b=>({...b,Seq:p.notebookIds.indexOf(b.NotebookId)}));data=true
   }
   await route.fulfill({json:data})
  })
  const titles=()=>page.locator('.notebook-select span[title]')
  const row=name=>page.locator('.notebook-row').filter({has:page.locator('.notebook-select span[title]',{hasText:new RegExp(`^${name}$`)})})
  await page.goto(`http://127.0.0.1:${server.httpServer.address().port}/note`)
  await expect(titles()).toHaveText(['Alpha','Beta','Gamma'])
  let box=await row('Gamma').boundingBox();await row('Alpha').dragTo(row('Gamma'),{targetPosition:{x:50,y:box.height-2}})
  await expect(titles()).toHaveText(['Beta','Gamma','Alpha']);expect(calls).toEqual([['b','c','a']])
  box=await row('Beta').boundingBox();await row('Alpha').dragTo(row('Beta'),{targetPosition:{x:50,y:2}})
  await expect(titles()).toHaveText(['Alpha','Beta','Gamma']);expect(calls[1]).toEqual(['a','b','c'])
  await row('Beta').locator('.notebook-select').click();await expect(row('Child')).toBeVisible()
  await row('Alpha').dragTo(row('Child'));expect(calls).toHaveLength(2)
  await page.reload();await expect(titles()).toHaveText(['Alpha','Beta','Gamma'])
  await page.close();console.log(`${desktop?'Desktop':'Web'} ${language}: sibling drag before/after, reload persistence, cross-parent rejection passed`)
 }
}finally{await browser.close();await server.close()}
