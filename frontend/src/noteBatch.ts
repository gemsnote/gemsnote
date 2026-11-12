import type {Params} from './api'

export type BatchAction='move'|'share'|'delete'
export interface BatchNote {NoteId:string;UserId:string;Title?:string;IsShared?:boolean;IsTrash?:boolean}
export interface BatchFailure {ids:string[];error:unknown}
export type BatchAPI=(path:string,params:Params)=>Promise<any>
export function ownsNote(note:BatchNote,userId:string){return note.UserId===userId&&!note.IsShared}

// Completed IDs survive an explicit retry, so successful shares and moves are
// not repeated after a partial failure. Ownership is checked before any write.
export async function runNoteBatch(action:BatchAction,notes:BatchNote[],userId:string,api:BatchAPI,completed:Set<string>,options:{destination?:string;recipient?:string;trash?:boolean},progress:()=>void=()=>{},signal?:AbortSignal):Promise<BatchFailure[]>{
 if(!notes.length||!userId)throw Error('请选择笔记')
 if((action==='share'||action==='delete')&&notes.some(note=>!ownsNote(note,userId)))throw Error('只能共享或删除自己的笔记')
 if(action==='move'&&(!options.destination||options.trash)||action==='share'&&(!options.recipient?.trim()||options.trash))throw Error('操作参数无效')
 const pending=notes.filter(note=>!completed.has(note.NoteId)),jobs:{ids:string[];run:()=>Promise<void>}[]=[]
 const add=(items:BatchNote[],path:string,params:Params)=>{if(items.length)jobs.push({ids:items.map(n=>n.NoteId),run:async()=>{await api(path,params)}})}
 if(action==='move'){
  const own=pending.filter(note=>ownsNote(note,userId))
  add(own,'/note/moveNote',{noteIds:own.map(n=>n.NoteId),notebookId:options.destination})
  for(const note of pending.filter(note=>!ownsNote(note,userId)))add([note],'/shared/move',{noteId:note.NoteId,notebookId:options.destination})
 }else if(action==='delete'){
  if(options.trash)for(const note of pending)add([note],'/note/deleteTrash',{noteId:note.NoteId})
  else add(pending,'/note/deleteNote',{noteIds:pending.map(n=>n.NoteId),isShared:false})
 }else{
  const recipient=options.recipient!.trim()
  for(const note of pending)jobs.push({ids:[note.NoteId],run:async()=>{
   const response=await api('/share/addShareNote',{noteId:note.NoteId,emails:[recipient],perm:0})
   const statuses:any[]=response&&typeof response==='object'?Object.values(response):[]
   if(!statuses.length||statuses.some(status=>status?.Ok!==true))throw Error(statuses.find(status=>status?.Ok!==true)?.Msg||'sharingLoadFailed')
  }})
 }
 const failures:BatchFailure[]=[]
 for(const job of jobs){
  if(signal?.aborted)break
  try{await job.run();for(const id of job.ids)completed.add(id)}catch(error){failures.push({ids:job.ids,error})}
  progress()
 }
 return failures
}
