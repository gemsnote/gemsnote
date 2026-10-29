import {describe,it,expect} from 'vitest'
import {zipSync,strToU8,type Zippable} from 'fflate'
import {parseMarkdownImport,importLimits} from './markdownImport'
import {createImportJob,runImportJob,type ImportAPI} from './importJob'

const zip=(entries:Zippable)=>zipSync(entries)
const md=(text:string)=>strToU8(text)
describe('Markdown import validation',()=>{
 it('uses the filename as title, strips UTF-8 BOM and preserves Markdown verbatim',()=>{
  const text='# 标题\r\n\r\n**内容** [链接](./other.md)\r\n'
  const source=parseMarkdownImport('我的笔记.MD',md('\ufeff'+text))
  expect(source.zip).toBe(false)
  expect(source.notes).toEqual([{path:'我的笔记.MD',title:'我的笔记',content:text,folder:''}])
 })
 it('imports root notes, nested notes and empty folders while skipping binary payloads and Mac metadata',()=>{
  const source=parseMarkdownImport('资料.zip',zip({'根.md':md('root'),'学习/Go/开始.markdown':md('start'),'学习/空/':new Uint8Array(),'image.png':new Uint8Array([255]),'__MACOSX/._根.md':md('ignore')}))
  expect(source.title).toBe('资料')
  expect(source.folders).toEqual(['学习','学习/Go','学习/空'])
  expect(source.notes.map(n=>n.folder)).toEqual(['学习/Go',''])
  expect(source.ignored).toBe(2)
 })
 it('plans mixed-language paths in a stable order, with parent folders first',()=>{
  const entries:Zippable={'空/子/两.md':md('two'),'Go/子/一.md':md('one'),'é.md':md('accent'),'E.md':md('latin'),'笔记.md':md('Chinese')}
  const source=parseMarkdownImport('notes.zip',zip(entries))
  const reversed=parseMarkdownImport('notes.zip',zip(Object.fromEntries(Object.entries(entries).reverse())))
  expect(source.folders).toEqual(['Go','空','Go/子','空/子'])
  expect(source.notes.map(n=>n.path)).toEqual(['E.md','Go/子/一.md','é.md','空/子/两.md','笔记.md'])
  expect(reversed.folders).toEqual(source.folders)
  expect(reversed.notes).toEqual(source.notes)
  const job=(input:typeof source)=>{let id=0;return createImportJob(input,'user','','notes',()=>String(++id))}
  expect(job(reversed).steps).toEqual(job(source).steps)
 })
 it('rejects invalid UTF-8 before any mutation',()=>{
  expect(()=>parseMarkdownImport('bad.md',new Uint8Array([255]))).toThrow('UTF-8')
  expect(()=>parseMarkdownImport('bad.zip',zip({'bad.md':new Uint8Array([255])}))).toThrow('UTF-8')
 })
 it.each(['../bad.md','/bad.md','C:/bad.md','a\\bad.md','a/../bad.md','a//bad.md'])( 'rejects unsafe path %s',path=>{
  expect(()=>parseMarkdownImport('bad.zip',zip({[path]:md('bad')}))).toThrow('不安全')
 })
 it('rejects oversized decompressed content before inflation',()=>{
  const bytes=zip({'big.md':md('small')})
  const view=new DataView(bytes.buffer)
  for(let i=0;i<bytes.length-4;i++)if(view.getUint32(i,true)===0x02014b50){view.setUint32(i+24,importLimits.note+1,true);break}
  expect(()=>parseMarkdownImport('big.zip',bytes)).toThrow('大小')
  expect(()=>parseMarkdownImport('big.md',new Uint8Array(importLimits.note+1))).toThrow('大小')
 })
 it('rejects corrupt checksums, truncated archives and encrypted archives',()=>{
  const original=zip({'note.md':md('content')})
  const corrupt=original.slice(),view=new DataView(corrupt.buffer)
  for(let i=0;i<corrupt.length-4;i++)if(view.getUint32(i,true)===0x02014b50){view.setUint32(i+16,0,true);break}
  expect(()=>parseMarkdownImport('bad.zip',corrupt)).toThrow('损坏')
  expect(()=>parseMarkdownImport('bad.zip',original.slice(0,-1))).toThrow('损坏')
  const encrypted=original.slice(),encryptedView=new DataView(encrypted.buffer)
  for(let i=0;i<encrypted.length-4;i++)if(encryptedView.getUint32(i,true)===0x02014b50){encryptedView.setUint16(i+8,1,true);break}
  expect(()=>parseMarkdownImport('bad.zip',encrypted)).toThrow('损坏')
 })
 it('rejects archives without Markdown and unsupported file extensions',()=>{
  expect(()=>parseMarkdownImport('empty.zip',zip({'a.png':md('image')}))).toThrow('没有可导入的文档')
  expect(()=>parseMarkdownImport('a.txt',md('text'))).toThrow('请选择')
 })
})

describe('resumable API2 import',()=>{
 function fixture(){
  let serial=0,failAfterWrite=''
  const books:any[]=[],notes=new Map<string,any>(),writes:string[]=[]
  const api:ImportAPI=async(path,params)=>{
   if(path==='/web/bootstrap')return {User:{UserId:'user'},Notebooks:books}
   if(path==='/web/document'){if(!notes.has(String(params?.noteId)))throw Error('notExists');return notes.get(String(params?.noteId))}
   const id=String(params?.noteId||params?.notebookId)
   writes.push(id)
   if(path==='/notebook/addNotebook')books.push({NotebookId:id,Title:params?.title,ParentNotebookId:params?.parentNotebookId})
   else notes.set(id,{Note:{NoteId:id,NotebookId:params?.notebookId,UserId:'user',Title:params?.title,IsMarkdown:params?.isMarkdown},Content:params?.content})
   if(failAfterWrite===path){failAfterWrite='';throw Error('connection lost')}
   return {Ok:true}
  }
  return {api,books,notes,writes,id:()=>String(++serial),fail:(path:string)=>{failAfterWrite=path}}
 }
 it('creates parents before children and notes with stable IDs and Markdown bodies',async()=>{
  const f=fixture(),source=parseMarkdownImport('Batch.zip',zip({'a/b/note.md':md('你好')}))
  const job=createImportJob(source,'user','existing','批量',f.id)
  await runImportJob(job,f.api,()=>{})
  expect(f.books.map(b=>[b.NotebookId,b.Title,b.ParentNotebookId])).toEqual([['1','批量','existing'],['2','a','1'],['3','b','2']])
  expect(f.notes.get('4')).toMatchObject({Note:{NotebookId:'3',IsMarkdown:true,Title:'note'},Content:'你好'})
  expect(job.next).toBe(job.steps.length)
 })
 it.each(['/notebook/addNotebook','/web/save'])('reconciles a lost response from %s without creating duplicates',async path=>{
  const f=fixture(),job=createImportJob(parseMarkdownImport('Batch.zip',zip({'n.md':md('body')})),'user','','Batch',f.id)
  f.fail(path)
  await expect(runImportJob(job,f.api,()=>{})).rejects.toThrow('connection lost')
  await runImportJob(job,f.api,()=>{})
  expect(f.writes).toEqual(['1','2'])
  expect(job.next).toBe(2)
 })
 it('pauses between records and resumes without re-uploading completed records',async()=>{
  const f=fixture(),job=createImportJob(parseMarkdownImport('Batch.zip',zip({'n.md':md('body')})),'user','','Batch',f.id)
  await runImportJob(job,f.api,()=>{},()=>job.next===1)
  expect(job.next).toBe(1)
  await runImportJob(job,f.api,()=>{})
  expect(f.writes).toEqual(['1','2'])
 })
 it('refuses to resume under another account',async()=>{
  const f=fixture(),job=createImportJob(parseMarkdownImport('n.md',md('body')),'other','existing','',f.id)
  await expect(runImportJob(job,f.api,()=>{})).rejects.toThrow('账号已改变')
  expect(f.writes).toEqual([])
 })
 it('does not overwrite an imported note changed after a lost response',async()=>{
  const f=fixture(),job=createImportJob(parseMarkdownImport('n.md',md('body')),'user','existing','',f.id)
  f.fail('/web/save')
  await expect(runImportJob(job,f.api,()=>{})).rejects.toThrow()
  f.notes.get('1').Content='edited elsewhere'
  await expect(runImportJob(job,f.api,()=>{})).rejects.toThrow('已发生变化')
  expect(f.writes).toEqual(['1'])
 })
})
