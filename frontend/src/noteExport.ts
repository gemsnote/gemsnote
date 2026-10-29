import {htmlToRtf} from './richDocument'
import {strToU8, zip} from 'fflate'

export interface ExportNote {
 NoteId: string; Title: string; Tags?: string[]; CreatedTime?: string; UpdatedTime?: string; IsMarkdown: boolean
}
export interface ExportBook {NotebookId: string; Title: string; Subs?: ExportBook[]}
const maxBytes=128*1024*1024
// JSON strings are valid YAML scalars; escape YAML-only Unicode line breaks too.
function quote(value:unknown){return JSON.stringify(value).replace(/\u0085/g,'\\u0085').replace(/\u2028/g,'\\u2028').replace(/\u2029/g,'\\u2029')}
function date(value?:string){return value&&Number.isFinite(Date.parse(value))?JSON.stringify(new Date(value).toISOString()):'null'}
export function hugoMarkdown(note:ExportNote,body:string,notebook:string):string{
 return `---\ntitle: ${quote(note.Title||'')}\ndate: ${date(note.CreatedTime)}\nlastmod: ${date(note.UpdatedTime)}\ntags: ${quote(note.Tags||[])}\ncategories: ${quote(notebook?[notebook]:[])}\n---\n\n${body}`
}
export function exportName(value:string,fallback='note'):string{
 let name=Array.from(value.replace(/[<>:"/\\|?*\x00-\x1f\x7f]/g,'_').trim().replace(/^\.+|[. ]+$/g,'')).slice(0,70).join('')||fallback
 while(new TextEncoder().encode(name).length>200)name=Array.from(name).slice(0,-1).join('')
 if(/^(con|prn|aux|nul|com[1-9]|lpt[1-9])(?:\.|$)/i.test(name))name='_'+name
 return name
}
function uniqueName(name:string,used:Set<string>,suffix=''){
 let result=name+suffix,index=2
 while(used.has(result.toLowerCase()))result=`${name} (${index++})${suffix}`
 used.add(result.toLowerCase());return result
}
export async function notebookZip(root:ExportBook,load:(book:ExportBook)=>Promise<{Note:ExportNote;Content:string}[]>|AsyncIterable<{Note:ExportNote;Content:string}>):Promise<Uint8Array>{
 const files:Record<string,Uint8Array>=Object.create(null),visited=new Set<string>();let size=0,count=0
 async function walk(book:ExportBook,path:string){
  if(visited.has(book.NotebookId))throw new Error('invalidNotebookTree')
  if(visited.size>=3000||path.split('/').length>100)throw new Error('exportTooLarge')
  visited.add(book.NotebookId);files[path+'/']=new Uint8Array()
  const used=new Set<string>(),children=(book.Subs||[]).map(child=>({child,name:uniqueName(exportName(child.Title,'notebook'),used)}))
  for await(const doc of await load(book)){
   const bytes=strToU8(doc.Note.IsMarkdown?hugoMarkdown(doc.Note,doc.Content,book.Title):htmlToRtf(doc.Content,doc.Note.Title));size+=bytes.length
   if(size>maxBytes||++count>10000)throw new Error('exportTooLarge')
   files[path+'/'+uniqueName(exportName(doc.Note.Title),used,doc.Note.IsMarkdown?'.md':'.rtf')]=bytes
  }
  for(const {child,name} of children)await walk(child,path+'/'+name)
 }
 await walk(root,exportName(root.Title,'notebook'))
 return new Promise((resolve,reject)=>zip(files,{level:6},(err,data)=>err?reject(err):resolve(data)))
}
export function base64Bytes(bytes:Uint8Array):string{
 let binary='';for(let i=0;i<bytes.length;i+=32768)binary+=String.fromCharCode(...bytes.subarray(i,i+32768));return btoa(binary)
}
