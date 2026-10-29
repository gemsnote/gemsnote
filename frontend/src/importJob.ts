import type {ImportSource} from './markdownImport'
import {ImportError} from './markdownImport'
import type {Params} from './api'

type BookStep={kind:'book';id:string;title:string;parent:string}
type NoteStep={kind:'note';id:string;title:string;parent:string;content:string;isMarkdown:boolean}
export interface ImportJob {userId:string;steps:(BookStep|NoteStep)[];next:number;uncertain:boolean;rootId:string;noteId:string;bookIds:string[]}
export type ImportAPI=(path:string,params?:Params)=>Promise<any>
export function createImportJob(source:ImportSource,userId:string,parent:string,name:string,newId:()=>string):ImportJob{
 const steps:ImportJob['steps']=[],ids=new Map<string,string>(),bookIds:string[]=[]
 const addBook=(path:string,title:string,parent:string)=>{const id=newId();ids.set(path,id);bookIds.push(id);steps.push({kind:'book',id,title,parent});return id}
 const rootId=source.zip?addBook('',name.trim(),parent):parent
 if(!rootId||!userId||(source.zip&&!name.trim()))throw new ImportError('请选择目标笔记本并填写名称')
 for(const path of source.folders){const at=path.lastIndexOf('/');addBook(path,path.slice(at+1),at<0?rootId:ids.get(path.slice(0,at))!)}
 let noteId=''
 for(const note of source.notes){const id=newId();noteId=id;steps.push({kind:'note',id,title:note.title,parent:note.folder?ids.get(note.folder)!:rootId,content:note.content,isMarkdown:note.isMarkdown!==false})}
 return {userId,steps,next:0,uncertain:false,rootId,noteId:source.zip?'':noteId,bookIds}
}
function flatten(books:any[]):any[]{return (books||[]).flatMap(b=>[b,...flatten(b.Subs)])}
export async function runImportJob(job:ImportJob,api:ImportAPI,progress:()=>void,stopped:()=>boolean=()=>false){
 // The account may have changed while a partially completed dialog was open.
 const boot=await api('/web/bootstrap')
 if(boot.User?.UserId!==job.userId)throw new ImportError('登录账号已改变，请关闭导入窗口')
 for(;job.next<job.steps.length;){
  if(stopped())return
  const step=job.steps[job.next]
  if(job.uncertain){
   let found:any
   if(step.kind==='book')found=flatten(boot.Notebooks).find(b=>b.NotebookId===step.id)
   else{try{found=await api('/web/document',{noteId:step.id})}catch{/* A failed read is safe: create below uses the same ID and never overwrites. */}}
   if(found){
    const matches=step.kind==='book'
     ?found.Title===step.title&&(found.ParentNotebookId||'')===step.parent
     :found.Note?.NoteId===step.id&&found.Note?.UserId===job.userId&&found.Note?.NotebookId===step.parent&&found.Note?.Title===step.title&&found.Note?.IsMarkdown===step.isMarkdown&&found.Content===step.content
    if(!matches)throw new ImportError('导入记录已发生变化，请检查已导入内容')
    job.next++;job.uncertain=false;progress();continue
   }
  }
  job.uncertain=true
  progress()
  if(step.kind==='book')await api('/notebook/addNotebook',{notebookId:step.id,title:step.title,parentNotebookId:step.parent})
  else await api('/web/save',{noteId:step.id,notebookId:step.parent,title:step.title,content:step.content,tags:'',isNew:true,isMarkdown:step.isMarkdown})
  job.next++;job.uncertain=false;progress()
 }
}
