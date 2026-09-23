<script setup lang="ts">
import {ref,computed,onMounted,onBeforeUnmount,watch} from 'vue'
import {useRoute,useRouter,onBeforeRouteLeave} from 'vue-router'
import {marked} from 'marked'
import DOMPurify from 'dompurify'
import {request,objectId,upload} from '../api'
import Navigation from '../components/Navigation.vue'
import brandMark from '../assets/gemsnote_s.png'
import RichEditor from '../components/RichEditor.vue'
import {language,t,formatDate} from '../i18n'
const route=useRoute(),router=useRouter(),boot=ref<any>({}),notes=ref<any[]>([]),current=ref<any>(null),authenticated=ref(false)
const notebook=ref(''),notebookSearch=ref(''),search=ref(''),tagFilter=ref(''),trash=ref(false),starred=ref(false),sharedOwner=ref('')
const title=ref(''),content=ref(''),tags=ref(''),dirty=ref(false),saving=ref(false),error=ref(''),preview=ref(false),panel=ref(''),destination=ref(''),shareEmail=ref(''),sharePerm=ref('0'),histories=ref<any[]>([]),members=ref<any[]>([]),attachments=ref<any[]>([]),sort=ref<'UpdatedTime'|'Title'>('UpdatedTime'),sortAsc=ref(false),sidebar=ref(false)
const notebooksVisible=ref(true),notesVisible=ref(true),notebooksWidth=ref(230),notesWidth=ref(290),bookMenu=ref(''),compact=ref(false)
const expandedNotebooks=ref<Set<string>>(new Set())
const mobileNotesVisible=ref(false)
const sortMenu=ref<HTMLDetailsElement>()
let layoutReady=false
function updateCompact(){compact.value=window.innerWidth<=1100}
const imageInput=ref<HTMLInputElement>(),attachInput=ref<HTMLInputElement>(),sharedQueueMsg=ref('')
const isSharedNote=computed(()=>!!current.value?.Note.IsShared)
const sharedCacheLabel=computed(()=>{if(!isSharedNote.value)return'';const state=current.value?.Note.CacheState as string|undefined;const at=current.value?.Note.CachedAt as string|undefined;switch(state){case 'ready':return at?t('离线缓存 {date}',{date:formatDate(at)}):t('离线缓存');case 'stale':return t('缓存待更新');case 'pending':return t('未下载，请先同步');case 'revoked':return t('已撤销访问');default:return''}})
function cacheStateLabel(a:any){switch(a.CacheState){case 'ready':return t('已缓存');case 'pending':return t('未下载');case 'failed':return t('下载失败');default:return''}}
async function queueSharedDownload(attachId:string){try{await request('/attach/queueSharedDownload',{attachId});sharedQueueMsg.value=t('已加入离线下载队列，稍后自动下载');attachmentsFor(current.value?.Note.NoteId||'')}catch(e){error.value=String(e)}}
let saveTimer:ReturnType<typeof setTimeout>|undefined,searchTimer:ReturnType<typeof setTimeout>|undefined,loadId=0,notesLoadId=0,savePromise:Promise<boolean>|undefined
const runtimeUnsubscribers:Array<()=>void>=[]
function compareBooks(a:any,b:any){return String(a.Title||'').localeCompare(String(b.Title||''),language.value,{sensitivity:'base',numeric:true})||String(a.NotebookId||'').localeCompare(String(b.NotebookId||''))}
function sortBookTree(items:any[]):any[]{return [...(items||[])].sort(compareBooks).map(n=>({...n,Subs:sortBookTree(n.Subs||[])}))}
function flatten(items:any[],depth=0):any[]{return (items||[]).flatMap(n=>[{...n,depth},...flatten(n.Subs,depth+1)])}
const notebookTree=computed(()=>sortBookTree(boot.value.Notebooks||[]))
const notebooks=computed(()=>flatten(notebookTree.value))
const currentNotebookTitle=computed(()=>{
 const id=current.value?.Note?.NotebookId
 if(!id)return t(isSharedNote.value?'共享笔记':'未知笔记本')
 const ownBook=notebooks.value.find(n=>n.NotebookId===id)
 if(ownBook)return ownBook.Title
 for(const books of Object.values(boot.value.SharedNotebooks||{})){
  const sharedBook=flatten(books as any[]).find(n=>n.NotebookId===id)
  if(sharedBook)return sharedBook.Title
 }
 return t(isSharedNote.value?'共享笔记本':'未知笔记本')
})
function flattenVisible(items:any[],depth=0):any[]{return (items||[]).flatMap(n=>[{...n,depth},...(expandedNotebooks.value.has(n.NotebookId)?flattenVisible(n.Subs,depth+1):[])])}
function flattenMatches(items:any[],key:string,depth=0):any[]{return (items||[]).flatMap(n=>{const children=flattenMatches(n.Subs,key,depth+1);const match=String(n.Title||'').toLowerCase().includes(key);return match||children.length?[{...n,depth},...children]:[]})}
const filteredBooks=computed(()=>{const key=notebookSearch.value.trim().toLowerCase();return key?flattenMatches(notebookTree.value,key):flattenVisible(notebookTree.value)})
const notesHeading=computed(()=>{
 if(trash.value)return t('回收站')
 if(starred.value)return t('已加星')
 if(tagFilter.value)return`#${tagFilter.value}`
 if(notebook.value){
  const books=sharedOwner.value?flatten(boot.value.SharedNotebooks?.[sharedOwner.value]||[]):notebooks.value
  return books.find((n:any)=>n.NotebookId===notebook.value)?.Title||t('共享笔记')
 }
 if(sharedOwner.value)return t('共享笔记')
 return t('所有笔记')
})
const sorted=computed(()=>[...notes.value].sort((a,b)=>{
 const result=sort.value==='Title'
  ?String(a.Title||'').localeCompare(String(b.Title||''),language.value,{sensitivity:'base',numeric:true})
  :(Date.parse(String(a.UpdatedTime||''))||0)-(Date.parse(String(b.UpdatedTime||''))||0)
 return (sortAsc.value?result:-result)||String(a.NoteId||'').localeCompare(String(b.NoteId||''))
}))
function setSort(field:'UpdatedTime'|'Title',ascending:boolean){sort.value=field;sortAsc.value=ascending;if(sortMenu.value)sortMenu.value.open=false}
const html=computed(()=>DOMPurify.sanitize(current.value?.Note.IsMarkdown?String(marked.parse(content.value,{async:false})):content.value))
const writable=computed(()=>!!current.value?.Writable)
const own=computed(()=>current.value?.Note.UserId===boot.value.User?.UserId)
const syncPending=computed(()=>dirty.value||!!boot.value.PendingChanges)
function layoutKey(){return `gemsnote:workspace:${boot.value.User?.UserId||'guest'}`}
function persistLayout(){if(!layoutReady)return;localStorage.setItem(layoutKey(),JSON.stringify({notebooksVisible:notebooksVisible.value,notesVisible:notesVisible.value,notebooksWidth:notebooksWidth.value,notesWidth:notesWidth.value}))}
function restoreLayout(){
 if(layoutReady)return
 notebooksWidth.value=Math.min(520,Math.max(180,Number(boot.value.User?.NotebookWidth)||230))
 notesWidth.value=Math.min(520,Math.max(180,Number(boot.value.User?.NoteListWidth)||290))
 // Space visibility is scoped to this workspace view. Returning from account
 // or admin should always start with the space panel visible.
 notebooksVisible.value=true
 try{const saved=JSON.parse(localStorage.getItem(layoutKey())||'null');if(saved){notesVisible.value=saved.notesVisible!==false;notebooksWidth.value=Math.min(520,Math.max(180,Number(saved.notebooksWidth)||notebooksWidth.value));notesWidth.value=Math.min(520,Math.max(180,Number(saved.notesWidth)||notesWidth.value))}}catch{}
 layoutReady=true
}
function showNotebooks(){notebooksVisible.value=true;if(window.innerWidth<=1100)sidebar.value=true;persistLayout()}
function hideNotebooks(){notebooksVisible.value=false;sidebar.value=false;persistLayout()}
function hideNotes(){notesVisible.value=false;persistLayout()}
function showNotes(){notesVisible.value=true;mobileNotesVisible.value=window.innerWidth<=700;sidebar.value=false;persistLayout()}
async function bootstrap(){boot.value=await request('/web/bootstrap');if(!boot.value.User){authenticated.value=false;await router.replace('/login');return false}authenticated.value=true;restoreLayout();return true}
async function load(){
 const token=++notesLoadId
 const path=sharedOwner.value?'/share/listShareNotes':'/web/notes'
 const all:any[]=[],seen=new Set<string>()
 try{
  for(let page=1;;page++){
   const params=sharedOwner.value
    ?{userId:sharedOwner.value,notebookId:notebook.value,page,sortField:sort.value,isAsc:sortAsc.value}
    :{notebookId:notebook.value,key:search.value,tag:tagFilter.value,trash:trash.value,starred:starred.value,page,sort:sort.value}
   const batch=await request<any[]>(path,params)
   if(token!==notesLoadId)return
   if(!Array.isArray(batch))throw new Error(t('笔记接口返回的数据格式不正确'))
   let added=0
   for(const item of batch){if(item?.NoteId&&!seen.has(item.NoteId)){seen.add(item.NoteId);all.push(item);added++}}
   // Browser and desktop return 100 notes per page; shared lists may return more.
   if(batch.length<100||added===0)break
  }
  const key=search.value.trim().toLocaleLowerCase()
  notes.value=sharedOwner.value&&key?all.filter(n=>`${n.Title||''} ${n.Desc||''}`.toLocaleLowerCase().includes(key)):all
 }catch(e){if(token===notesLoadId)error.value=String(e)}
}
function runSearch(){clearTimeout(searchTimer);starred.value=false;void load()}
function scheduleSearch(){clearTimeout(searchTimer);searchTimer=setTimeout(runSearch,300)}
async function setStar(n:any){if(sharedOwner.value)return;try{const next=!n.IsStar;await request('/web/star',{noteId:n.NoteId,starred:next});n.IsStar=next;boot.value.TotalStarred=Math.max(0,Number(boot.value.TotalStarred||0)+(next?1:-1));if(boot.value.Desktop)boot.value.PendingChanges=true;if(starred.value)await load()}catch(e){error.value=String(e)}}
async function attachmentsFor(noteId:string){try{const r:any=await request('/attach/getAttachs',{noteId});attachments.value=r.List||[]}catch(e){attachments.value=[];error.value=String(e)}}
async function refreshDocument(){if(!current.value)return;current.value=await request('/web/document',{noteId:current.value.Note.NoteId})}
async function open(id:string){if(!await flush())return;const token=++loadId;try{const doc=await request('/web/document',{noteId:id});if(token!==loadId)return;current.value=doc;mobileNotesVisible.value=false;title.value=doc.Note.Title;content.value=doc.Content||'';tags.value=(doc.Note.Tags||[]).join(',');dirty.value=false;panel.value='';await attachmentsFor(id);await router.replace(`/note/${id}`)}catch(e){error.value=String(e)}}
async function select(id='',owner='',isTrash=false,isStar=false){if(!await flush())return;clearTimeout(searchTimer);showNotes();notebook.value=id;sharedOwner.value=owner;trash.value=isTrash;starred.value=isStar;search.value='';tagFilter.value='';bookMenu.value='';await load();sidebar.value=false}
async function selectTag(tag:string){if(!await flush())return;clearTimeout(searchTimer);showNotes();notebook.value='';sharedOwner.value='';trash.value=false;starred.value=false;search.value='';tagFilter.value=tag;await load();sidebar.value=false}
function hasChildren(n:any){return Array.isArray(n.Subs)&&n.Subs.length>0}
async function selectBook(n:any){if(hasChildren(n)){const next=new Set(expandedNotebooks.value);if(next.has(n.NotebookId))next.delete(n.NotebookId);else next.add(n.NotebookId);expandedNotebooks.value=next}await select(n.NotebookId)}
let stopResize:(()=>void)|undefined
function resizeBy(panelName:'notebooks'|'notes',delta:number){if(panelName==='notebooks')notebooksWidth.value=Math.min(520,Math.max(180,notebooksWidth.value+delta));else notesWidth.value=Math.min(520,Math.max(180,notesWidth.value+delta));persistLayout()}
function resizePanel(panelName:'notebooks'|'notes',event:PointerEvent){
 const startX=event.clientX,startWidth=panelName==='notebooks'?notebooksWidth.value:notesWidth.value
 const move=(e:PointerEvent)=>{const width=Math.min(520,Math.max(180,startWidth+e.clientX-startX));if(panelName==='notebooks')notebooksWidth.value=width;else notesWidth.value=width}
 const stop=()=>{persistLayout();document.body.classList.remove('panel-resizing');window.removeEventListener('pointermove',move);window.removeEventListener('pointerup',stop);window.removeEventListener('pointercancel',stop);stopResize=undefined}
 stopResize=stop
 document.body.classList.add('panel-resizing')
 window.addEventListener('pointermove',move);window.addEventListener('pointerup',stop)
 window.addEventListener('pointercancel',stop)
}
function changed(){dirty.value=true;clearTimeout(saveTimer);saveTimer=setTimeout(()=>save(),1200)}
async function save():Promise<boolean>{
 if(savePromise)return savePromise
 if(!dirty.value||!current.value||!writable.value)return true
 const id=current.value.Note.NoteId,snapshot={title:title.value,content:content.value,tags:tags.value}
 saving.value=true;error.value=''
 savePromise=(async()=>{try{const doc=await request('/web/save',{noteId:id,usn:current.value.Note.Usn,...snapshot});if(current.value?.Note.NoteId===id){current.value=doc;dirty.value=title.value!==snapshot.title||content.value!==snapshot.content||tags.value!==snapshot.tags}if(boot.value.Desktop)boot.value.PendingChanges=true;await load();return true}catch(e){const message=e instanceof Error?e.message:String(e);error.value=(message==='conflict'||message===t('conflict'))?t('此笔记已在其他端修改。请先导出本地内容，再重新加载笔记。'):message;return false}finally{saving.value=false;savePromise=undefined}})()
 return savePromise
}
async function flush(){clearTimeout(saveTimer);if(!await save())return false;return dirty.value?save():true}
async function create(markdown:boolean){if(!await flush())return;let id=notebook.value||(!sharedOwner.value?notebooks.value[0]?.NotebookId:'');if(!id&&!sharedOwner.value){await addBook();id=notebooks.value[0]?.NotebookId}if(!id){error.value=t('请先选择一个可编辑的共享笔记本');return}try{const doc=await request('/web/save',{noteId:objectId(),notebookId:id,ownerId:sharedOwner.value,title:t('未命名笔记'),content:'',tags:'',isNew:true,isMarkdown:markdown});await bootstrap();await load();await open(doc.Note.NoteId)}catch(e){error.value=String(e)}}
async function addBook(parentNotebookId=''){const name=prompt(t(parentNotebookId?'子笔记本名称':'笔记本名称'));if(!name)return;bookMenu.value='';try{await request('/notebook/addNotebook',{notebookId:objectId(),title:name,parentNotebookId});await bootstrap()}catch(e){error.value=String(e)}}
async function bookAction(id:string,remove=false){bookMenu.value='';try{if(remove){if(!confirm(t('删除此笔记本？请先移动其中的笔记。')))return;await request('/notebook/deleteNotebook',{notebookId:id});if(notebook.value===id)notebook.value=''}else{const name=prompt(t('笔记本名称'),notebooks.value.find(n=>n.NotebookId===id)?.Title);if(!name)return;await request('/notebook/updateNotebookTitle',{notebookId:id,title:name})}await bootstrap();await load()}catch(e){error.value=String(e)}}
async function remove(){if(!current.value||!confirm(t(own.value?(trash.value?'永久删除这篇笔记？此操作不可撤销。':'将笔记放入回收站？'):'从共享列表移除这篇笔记？')))return;try{if(own.value&&trash.value)await request('/note/deleteTrash',{noteId:current.value.Note.NoteId});else if(own.value)await request('/note/deleteNote',{noteIds:[current.value.Note.NoteId],isShared:false});else await request('/share/deleteShareNoteBySharedUser',{noteId:current.value.Note.NoteId,fromUserId:current.value.Note.UserId});dirty.value=false;current.value=null;await bootstrap();await load()}catch(e){error.value=String(e)}}
async function restore(){try{await request('/web/restore',{noteId:current.value.Note.NoteId});current.value=null;await bootstrap();await load()}catch(e){error.value=String(e)}}
async function move(copy=false){if(!await flush()||!destination.value)return;try{const noteId=current.value.Note.NoteId;await request(copy?'/note/copyNote':'/note/moveNote',{noteIds:[noteId],notebookId:destination.value});panel.value='';await bootstrap();await load();await open(noteId)}catch(e){error.value=String(e)}}
async function share(){try{const result=await request('/share/addShareNote',{noteId:current.value.Note.NoteId,emails:[shareEmail.value],perm:Number(sharePerm.value)});const failures=Object.values(result).filter((r:any)=>!r.Ok);if(failures.length)throw new Error(JSON.stringify(failures));shareEmail.value='';await showShare()}catch(e){error.value=String(e)}}
async function showShare(){panel.value='share';try{const r=await request('/web/shareMembers',{noteId:current.value.Note.NoteId});members.value=r.Users||[]}catch(e){error.value=String(e)}}
async function revoke(id:string){try{await request('/share/deleteShareNote',{noteId:current.value.Note.NoteId,toUserId:id});await showShare()}catch(e){error.value=String(e)}}
async function history(){panel.value='history';try{histories.value=await request('/noteContentHistory/listHistories',{noteId:current.value.Note.NoteId})||[]}catch(e){error.value=String(e)}}
function chooseImage(){imageInput.value?.click()}
function chooseAttach(){attachInput.value?.click()}
function appendUpload(value:string){content.value+=(content.value&& !content.value.endsWith('\n')?'\n':'')+value;changed()}
async function uploadImage(event:Event){const input=event.target as HTMLInputElement,file=input.files?.[0];input.value='';if(!file||!current.value)return;try{const result:any=await upload('/file/pasteImage',file,{noteId:current.value.Note.NoteId});if(!result.Id)throw new Error(result.Msg||t('图片上传失败'));const src='/api2/file/getImage?fileId='+encodeURIComponent(result.Id);appendUpload(current.value.Note.IsMarkdown?`![${file.name}](${src})`:`<img src="${src}" alt="${file.name}">`)}catch(e){error.value=String(e)}}
async function uploadAttach(event:Event){const input=event.target as HTMLInputElement,file=input.files?.[0];input.value='';if(!file||!current.value||!await flush())return;try{await upload('/attach/uploadAttach',file,{noteId:current.value.Note.NoteId});await refreshDocument();await attachmentsFor(current.value.Note.NoteId);panel.value='attachments'}catch(e){error.value=String(e)}}
async function deleteAttach(attachId:string){if(!confirm(t('删除此附件？'))||!await flush())return;try{await request('/attach/deleteAttach',{attachId});await refreshDocument();await attachmentsFor(current.value.Note.NoteId)}catch(e){error.value=String(e)}}
function download(){const blob=new Blob([content.value],{type:current.value.Note.IsMarkdown?'text/markdown':'text/html'});const url=URL.createObjectURL(blob);const a=document.createElement('a');a.href=url;a.download=(title.value||'note')+(current.value.Note.IsMarkdown?'.md':'.html');a.click();URL.revokeObjectURL(url)}
function leave(e:BeforeUnloadEvent){if(dirty.value||saving.value){e.preventDefault();e.returnValue=''}}
function shortcut(e:KeyboardEvent){if(e.key==='Escape'){bookMenu.value='';if(sortMenu.value)sortMenu.value.open=false}if((e.ctrlKey||e.metaKey)&&e.key==='s'){e.preventDefault();save()}}
function closeBookMenu(e:PointerEvent){bookMenu.value='';if(sortMenu.value&&!sortMenu.value.contains(e.target as Node))sortMenu.value.open=false}
function handleSyncResult(result:any){
 const value=result?.detail||result||{}
 if(value.Ok===false){error.value=String(value.Msg||t('同步失败'));return}
 if(value.Full)error.value=''
}
onMounted(async()=>{(window as any).__gemsnoteBeforeLogout=flush;(window as any).__gemsnoteBeforeReset=flush;(window as any).__gemsnoteAfterReset=()=>{dirty.value=false;current.value=null;loadId++;notesLoadId++};(window as any).__gemsnoteSyncNow=async()=>{if(!await flush())return false;if(!await bootstrap())return false;await load();return true};window.addEventListener('sync-result',handleSyncResult as EventListener);updateCompact();window.addEventListener('resize',updateCompact);window.addEventListener('beforeunload',leave);window.addEventListener('keydown',shortcut);window.addEventListener('pointerdown',closeBookMenu);const wails=(window as any).runtime;if(wails?.EventsOn){const offRevoked=wails.EventsOn('shared-notes-revoked',(ids:string[])=>{if(ids?.includes(current.value?.Note?.NoteId)){error.value=t('该共享笔记已被撤销访问权限');current.value=null;attachments.value=[]}});const offSync=wails.EventsOn('sync-finished',async(result:any)=>{handleSyncResult(result);if(result?.Reset){if(result?.Ok===false){dirty.value=false;current.value=null;try{if(await bootstrap())await load()}catch(e){error.value=String(e)}}return}try{if(!await bootstrap())return;await load();if(current.value&&!dirty.value){const noteId=current.value.Note.NoteId;const doc=await request('/web/document',{noteId});if(current.value?.Note.NoteId===noteId&&!dirty.value){current.value=doc;title.value=doc.Note.Title;content.value=doc.Content||'';tags.value=(doc.Note.Tags||[]).join(',')}}}catch(e){error.value=String(e)}});if(typeof offRevoked==='function')runtimeUnsubscribers.push(offRevoked);if(typeof offSync==='function')runtimeUnsubscribers.push(offSync)}try{if(await bootstrap()){await load();if(route.params.noteId)await open(String(route.params.noteId))}}catch(e){error.value=String(e)}})
onBeforeUnmount(()=>{delete (window as any).__gemsnoteBeforeLogout;delete (window as any).__gemsnoteBeforeReset;delete (window as any).__gemsnoteAfterReset;delete (window as any).__gemsnoteSyncNow;clearTimeout(saveTimer);clearTimeout(searchTimer);stopResize?.();runtimeUnsubscribers.splice(0).forEach(unsubscribe=>unsubscribe());document.body.classList.remove('panel-resizing');window.removeEventListener('sync-result',handleSyncResult as EventListener);window.removeEventListener('resize',updateCompact);window.removeEventListener('beforeunload',leave);window.removeEventListener('keydown',shortcut);window.removeEventListener('pointerdown',closeBookMenu)})
onBeforeRouteLeave(async()=>await flush())
watch(()=>route.params.noteId,id=>{if(id&&id!==current.value?.Note.NoteId)open(String(id));else if(!id)flush().then(ok=>{if(ok)current.value=null})})
</script>
<template>
<main v-if="!authenticated" class="auth-shell" aria-busy="true"></main>
<div v-else class="shell workspace" :class="{'mobile-list-open':mobileNotesVisible}">
<aside class="notebooks" :class="{mobileOpen:sidebar,collapsed:!notebooksVisible}" :style="{width:(notebooksVisible?notebooksWidth:64)+'px'}">
<header class="panel-header">
<button v-if="!notebooksVisible" class="space-brand collapsed-brand" :title="t('展开我的空间')" :aria-label="t('展开我的空间')" @click="showNotebooks"><img :src="brandMark" alt=""></button>
<RouterLink v-else to="/note" class="space-brand" :title="t('珠玑笔记')"><img :src="brandMark" :alt="t('珠玑笔记')"><span>{{t('珠玑笔记')}}</span></RouterLink>
<div v-if="notebooksVisible" class="panel-actions"><button class="icon-button" @click="addBook()" :title="t('新建笔记本')" :aria-label="t('新建笔记本')">＋</button><button class="icon-button desktop-only" @click="hideNotebooks" :title="t('收起我的空间')" :aria-label="t('收起我的空间')">‹</button></div>
</header>
<div v-if="notebooksVisible" class="notebook-scroll">
<input class="notebook-search" v-model="notebookSearch" :placeholder="t('搜索笔记本')" :aria-label="t('搜索笔记本')">
<button :class="{selected:!notebook&&!trash&&!starred&&!sharedOwner}" @click="select()">{{t('所有笔记')}} <small>{{boot.TotalNotes||0}}</small></button>
<button :class="{selected:starred}" @click="select('','',false,true)">★ {{t('已加星')}} <small>{{boot.TotalStarred||0}}</small></button>
<div v-for="n in filteredBooks" :key="n.NotebookId" class="notebook-row" :class="{selected:notebook===n.NotebookId}" :style="{paddingLeft:4+n.depth*16+'px'}">
<button class="notebook-select" @click="selectBook(n)"><span class="tree-toggle" :aria-hidden="true">{{hasChildren(n)?(expandedNotebooks.has(n.NotebookId)?'▾':'▸'):'▱'}}</span> <span :title="n.Title">{{n.Title}}</span><small>{{n.NumberNotes}}</small></button>
<button class="notebook-more" aria-haspopup="menu" :aria-expanded="bookMenu===n.NotebookId" :aria-label="t('{name} 菜单',{name:n.Title})" @pointerdown.stop @click.stop="bookMenu=bookMenu===n.NotebookId?'':n.NotebookId">⋯</button>
<div v-if="bookMenu===n.NotebookId" class="notebook-menu" role="menu" @pointerdown.stop><button role="menuitem" @click="addBook(n.NotebookId)">{{t('新增子笔记本')}}</button><button role="menuitem" @click="bookAction(n.NotebookId)">{{t('重命名')}}</button><button role="menuitem" @click="bookAction(n.NotebookId,true)">{{t('删除')}}</button></div>
</div>
<h3>{{t('共享给我')}}</h3>
<p v-if="boot.SharedCache==='unsupported'" class="muted">{{t('服务端不支持共享离线缓存')}}</p>
<div v-for="(books,owner) in boot.SharedNotebooks" :key="String(owner)">
<button v-for="n in flatten(books)" :key="n.NotebookId" @click="select(n.IsDefault?'':n.NotebookId,String(owner))" :title="n.Title||t('共享笔记')">♧ {{n.Title||t('共享笔记')}}</button>
</div>
<p v-if="!Object.keys(boot.SharedNotebooks||{}).length" class="muted">{{t('暂无共享内容')}}</p>
<h3>{{t('标签')}}</h3>
<div class="tags">
<button v-for="tag in boot.Tags" :key="tag.Tag" :title="`${tag.Tag} (${t('{count}篇笔记',{count:tag.Count||0})})`" @click="selectTag(tag.Tag)">#{{tag.Tag}} <small>{{tag.Count||0}}</small></button>
</div>
<button @click="select('','',true)">{{t('回收站')}}</button>
</div>
<Navigation workspace :collapsed="!notebooksVisible" :admin="boot.IsAdmin" :user="boot.User" :pending="syncPending" :desktop="boot.Desktop" @synced="bootstrap"/>
<div v-if="notebooksVisible" class="resize-handle" role="separator" tabindex="0" aria-orientation="vertical" :aria-label="t('调节我的空间宽度')" :aria-valuenow="notebooksWidth" aria-valuemin="180" aria-valuemax="520" @keydown.arrow-left.prevent="resizeBy('notebooks',-10)" @keydown.arrow-right.prevent="resizeBy('notebooks',10)" @pointerdown.prevent="resizePanel('notebooks',$event)"></div>
</aside>
<section v-if="notesVisible" class="note-list" :style="{width:notesWidth+'px'}">
<header class="panel-header">
<details ref="sortMenu" class="note-sort-menu">
<summary class="icon-button" :title="t('排序方式')" :aria-label="t('排序方式')" aria-haspopup="menu">☰</summary>
<div class="note-sort-options" role="menu">
<button type="button" role="menuitemradio" :aria-checked="sort==='UpdatedTime'&&!sortAsc" @click="setSort('UpdatedTime',false)">{{t('修改时间：降序')}}</button>
<button type="button" role="menuitemradio" :aria-checked="sort==='UpdatedTime'&&sortAsc" @click="setSort('UpdatedTime',true)">{{t('修改时间：升序')}}</button>
<button type="button" role="menuitemradio" :aria-checked="sort==='Title'&&!sortAsc" @click="setSort('Title',false)">{{t('标题：降序')}}</button>
<button type="button" role="menuitemradio" :aria-checked="sort==='Title'&&sortAsc" @click="setSort('Title',true)">{{t('标题：升序')}}</button>
</div>
</details>
<button v-if="compact" class="mobile-toggle" :title="t('我的空间')" :aria-label="t('我的空间')" @click="notebooksVisible?sidebar=!sidebar:showNotebooks()">▤</button>
<h2 :title="notesHeading">{{notesHeading}}</h2>
<div class="note-header-actions">
<button class="icon-button" @click="create(false)" :disabled="trash" :title="t('新建富文本笔记')" :aria-label="t('新建富文本笔记')"><span aria-hidden="true">＋</span></button>
<button class="icon-button" @click="create(true)" :disabled="trash" :title="t('新建 Markdown 笔记')" :aria-label="t('新建 Markdown 笔记')"><span class="markdown-icon" aria-hidden="true">M＋</span></button>
<button class="icon-button desktop-only" @click="hideNotes" :title="t('隐藏笔记栏')" :aria-label="t('隐藏笔记栏')">‹</button>
</div>
</header>
<div class="note-search-row">
<input class="note-search" v-model="search" type="search" :placeholder="t('搜索笔记')" :aria-label="t('搜索笔记')" @input="scheduleSearch" @keydown.enter.prevent="runSearch">
</div>
<div class="list-items">
<div v-for="n in sorted" :key="n.NoteId" class="note-row">
<button v-if="!sharedOwner&&!trash" class="star-button" :aria-label="t(n.IsStar?'取消加星':'加星')" @click="setStar(n)">{{n.IsStar?'★':'☆'}}</button>
<button class="note-item" :class="{selected:current?.Note.NoteId===n.NoteId}" @click="open(n.NoteId)">
<strong :title="n.Title||t('未命名')"><span v-if="n.IsDirty" class="note-pending-dot" :aria-label="t('有未同步的更改')"></span>{{n.Title||t('未命名')}}</strong>
<p>{{n.Desc||t('暂无摘要')}}</p>
<small>{{formatDate(n.UpdatedTime)}} {{n.Perm===0?t('· 只读'):''}}</small>
</button>
</div>
<p v-if="!notes.length" class="empty">{{t('这里还没有笔记')}}</p>
</div>
<div class="resize-handle" role="separator" tabindex="0" aria-orientation="vertical" :aria-label="t('调节笔记栏宽度')" :aria-valuenow="notesWidth" aria-valuemin="180" aria-valuemax="520" @keydown.arrow-left.prevent="resizeBy('notes',-10)" @keydown.arrow-right.prevent="resizeBy('notes',10)" @pointerdown.prevent="resizePanel('notes',$event)"></div>
</section>
<main class="editor">
<p v-if="error" role="alert" class="error">{{error}} <button @click="error=''">{{t('关闭')}}</button>
</p>
<template v-if="current">
<header class="toolbar">
<input class="toolbar-title" v-model="title" :readonly="!writable" @input="changed" :aria-label="t('笔记标题')">
<span v-if="isSharedNote" class="shared-cache-label">{{sharedCacheLabel}}</span>
<div class="toolbar-actions">
<button @click="preview=!preview">{{t(preview?'编辑':'预览')}}</button>
<button v-if="writable&&!trash" @click="chooseImage">{{t('图片')}}</button>
<button v-if="writable&&!trash" @click="chooseAttach">{{t('附件')}}{{attachments.length?` (${attachments.length})`:''}}</button>
<button v-else-if="attachments.length" @click="panel='attachments'">{{t('附件')}}{{` (${attachments.length})`}}</button>
<button v-if="own&&!trash" @click="panel='move'">{{t('移动/复制')}}</button>
<button v-if="own&&!trash" @click="showShare">{{t('共享')}}</button>
<button v-if="own" @click="history">{{t('历史')}}</button>
<button @click="download">{{t('导出')}}</button>
<button v-if="trash&&own" @click="restore">{{t('恢复')}}</button>
<button v-if="own||writable" @click="remove">{{t('删除')}}</button>
</div>
<input ref="imageInput" class="visually-hidden" type="file" accept="image/*" @change="uploadImage">
<input ref="attachInput" class="visually-hidden" type="file" @change="uploadAttach">
</header>
<div class="note-info-line" :aria-label="t('笔记基本信息')">
<span :title="currentNotebookTitle">{{t('笔记本：')}}{{currentNotebookTitle}}</span>
<span>{{t('创建：')}}{{formatDate(current.Note.CreatedTime)}}</span>
<span>{{t('修改：')}}{{formatDate(current.Note.UpdatedTime)}}</span>
<span>{{t('格式：')}}{{current.Note.IsMarkdown?'MD':'RT'}}</span>
<span class="note-id" :title="current.Note.NoteId">{{current.Note.NoteId}}</span>
</div>
<div v-if="preview||!writable" class="rendered" v-html="html">
</div>
<RichEditor v-else-if="!current.Note.IsMarkdown" :key="`${current.Note.NoteId}:${language}`" v-model="content" @update:model-value="changed"/>
<textarea v-else v-model="content" class="content-editor" @input="changed" :aria-label="t(current.Note.IsMarkdown?'Markdown 正文':'HTML 正文')" spellcheck="false">
</textarea>
<div class="note-tag-line"><span>{{t('标签')}}</span><input class="tag-input" v-model="tags" :readonly="!writable" @input="changed" :placeholder="t('标签，以逗号分隔')" :aria-label="t('笔记标签')"></div>
<section v-if="panel" class="dialog-backdrop" @click.self="panel=''">
<div class="dialog" role="dialog" aria-modal="true">
<button class="close" @click="panel=''">{{t('关闭')}}</button>
<template v-if="panel==='move'">
<h2>{{t('移动或复制笔记')}}</h2>
<input v-model="notebookSearch" :placeholder="t('搜索笔记本')">
<select v-model="destination">
<option value="">{{t('选择笔记本')}}</option>
<option v-for="n in filteredBooks" :value="n.NotebookId">{{n.Title}}</option>
</select>
<button @click="move()">{{t('移动')}}</button>
<button @click="move(true)">{{t('复制')}}</button>
</template>
<template v-if="panel==='share'">
<h2>{{t('共享笔记')}}</h2>
<form @submit.prevent="share">
<label>{{t('对方邮箱')}}<input v-model="shareEmail" type="email" required>
</label>
<select v-model="sharePerm">
<option value="0">{{t('只读')}}</option>
<option value="1">{{t('可编辑')}}</option>
</select>
<button>{{t('添加共享')}}</button>
</form>
<p v-for="m in members" :key="m.ToUserId">{{m.Email}} · {{t(m.Perm?'可编辑':'只读')}} <button v-if="!m.NotebookHasShared" @click="revoke(m.ToUserId)">{{t('取消共享')}}</button>
</p>
</template>
<template v-if="panel==='attachments'">
<h2>{{t('附件')}}</h2>
<p v-if="sharedQueueMsg" role="status" class="muted">{{sharedQueueMsg}}</p>
<p v-if="!attachments.length" class="muted">{{t('暂无附件')}}</p>
<ul class="attachments">
<li v-for="attachment in attachments" :key="attachment.AttachId">
<a :href="`/api2/attachments/download?attachId=${attachment.AttachId}`">{{attachment.Title}}</a>
<small>{{Math.ceil((attachment.Size||0)/1024)}} KB</small>
<small v-if="isSharedNote" class="muted">{{cacheStateLabel(attachment)}}</small>
<button v-if="isSharedNote&&attachment.CacheState!=='ready'" @click="queueSharedDownload(attachment.AttachId)">{{t('下载供离线使用')}}</button>
<button v-if="writable" @click="deleteAttach(attachment.AttachId)">{{t('删除')}}</button>
</li>
</ul>
</template>
<template v-if="panel==='history'">
<h2>{{t('历史版本')}}</h2>
<p v-if="!histories.length">{{t('暂无历史')}}</p>
<details v-for="h in histories" :key="h.UpdatedTime">
<summary>{{formatDate(h.UpdatedTime)}}</summary>
<pre>{{h.Content}}</pre>
<button :disabled="!writable" @click="content=h.Content;changed();panel=''">{{t('恢复到编辑器')}}</button>
</details>
</template>
</div>
</section>
</template>
<div v-else class="welcome">
<img class="welcome-icon" :src="brandMark" :alt="t('珠玑笔记')" />
<h1>{{t('日积字句，终得珠玑。')}}</h1>
<p>{{t('选择笔记或开始全新记录。')}}</p>
<button class="primary" @click="create(true)">{{t('新建笔记')}}</button>
</div>
</main>
</div>
</template>
<style scoped>.visually-hidden{position:absolute;width:1px;height:1px;padding:0;margin:-1px;overflow:hidden;clip:rect(0,0,0,0);white-space:nowrap;border:0}.note-row{position:relative}.note-row .note-item{padding-left:40px}.star-button{position:absolute;z-index:1;left:8px;top:13px;padding:3px;border:0;background:transparent;color:#b17b18;font-size:18px;cursor:pointer}.star-button:focus-visible{outline:2px solid currentColor;outline-offset:2px;border-radius:2px}.attachments{padding:0;list-style:none}.attachments li{display:flex;align-items:center;gap:10px;padding:10px 0;border-bottom:1px solid #edf0e9}.attachments li a{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.attachments li small{margin-left:auto}.attachments li button{padding:4px 8px;font-size:12px}</style>
