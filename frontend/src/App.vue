<script setup lang="ts">
import {computed,onBeforeUnmount,onMounted,ref} from 'vue'
import {t} from './i18n'
import AppDialog from './components/AppDialog.vue'
import {installLocalizedValidation} from './validation'
import {request} from './api'

const syncVisible=ref(false)
const syncStage=ref('start')
const syncMode=ref<'full'|'reset'>('full')
const syncCurrent=ref(0)
const syncTotal=ref(0)
const syncPercent=ref(0)
let polling:ReturnType<typeof setTimeout>|undefined
let progressGeneration=0
let awaitingRequest=false
const unsubscribers:Array<()=>void>=[]
const stageText=computed(()=>t(({start:'正在准备同步',checking:'正在连接服务器',resetting:'正在清理本地缓存',notebooks:'正在同步笔记本',notes:'正在同步笔记',tags:'正在同步标签',push:'正在上传本地更改',images:'正在同步图片和附件',done:'同步完成'} as Record<string,string>)[syncStage.value]||'正在完全同步'))
const completedText=computed(()=>syncTotal.value>0?t('已完成 {current} / {total} 项',{current:syncCurrent.value,total:syncTotal.value}):t('已完成 {count} 项',{count:syncCurrent.value}))
const indeterminate=computed(()=>syncTotal.value===0&&syncStage.value!=='done')

function applyProgress(value:any){
 if(value.Mode==='full'||value.Mode==='reset')syncMode.value=value.Mode
 syncStage.value=String(value.Stage||'start')
 syncCurrent.value=Number(value.Current||0)
 syncTotal.value=Number(value.Total||0)
 syncPercent.value=Math.max(0,Math.min(100,Number(value.Percent||0)))
}
function startPolling(){
 const generation=++progressGeneration
 clearTimeout(polling)
 const poll=async()=>{
  try{
   const value=await request('/web/syncProgress')
   // An old idle snapshot must not close a newly requested sync or reset
   // its mode. The owning HTTP request is authoritative for completion.
   if(generation===progressGeneration&&syncVisible.value&&value.Running&&value.Mode===syncMode.value)applyProgress(value)
  }catch{/* The sync request reports failures; a transient poll must not hide it. */}
  if(generation===progressGeneration&&syncVisible.value)polling=setTimeout(poll,500)
 }
 polling=setTimeout(poll,100)
}
function showSyncProgress(event:any){
 const requested=!!event?.detail
 const value=event?.detail||event||{}
 // Only an explicit UI action opens the modal. Late native events from an
 // already completed request must not reopen it (nor should auto-sync).
 if(!requested&&!syncVisible.value)return
 if(!requested&&value.Running===false)return
 if(!requested&&syncVisible.value&&value.Mode&&value.Mode!==syncMode.value)return
 const opening=!syncVisible.value
 if(requested)awaitingRequest=true
 applyProgress(value)
 syncVisible.value=true
 if(opening||requested)startPolling()
}
function finishSync(){
 syncVisible.value=false
 awaitingRequest=false
 progressGeneration++
 clearTimeout(polling)
}
function finishNativeSync(){if(!awaitingRequest)finishSync()}
function finishWindowSync(event:any){const value=event?.detail||event||{};if(value.Full||value.Ok===false)finishSync()}

onMounted(()=>{
 unsubscribers.push(installLocalizedValidation())
 window.addEventListener('sync-progress',showSyncProgress as EventListener)
 window.addEventListener('sync-result',finishWindowSync as EventListener)
 const wails=(window as any).runtime
 if(wails?.EventsOn){
  const offProgress=wails.EventsOn('sync-progress',showSyncProgress)
  const offFinished=wails.EventsOn('sync-finished',finishNativeSync)
  if(typeof offProgress==='function')unsubscribers.push(offProgress)
  if(typeof offFinished==='function')unsubscribers.push(offFinished)
 }
})
onBeforeUnmount(()=>{
 finishSync()
 window.removeEventListener('sync-progress',showSyncProgress as EventListener)
 window.removeEventListener('sync-result',finishWindowSync as EventListener)
 unsubscribers.splice(0).forEach(unsubscribe=>unsubscribe())
})
</script>

<template>
 <RouterView />
 <AppDialog />
 <section v-if="syncVisible" class="sync-progress-backdrop" aria-live="polite">
  <div class="sync-progress-dialog" role="dialog" aria-modal="true" :aria-label="t(syncMode==='reset'?'重新同步':'完全同步')">
   <h2>{{t(syncMode==='reset'?'重新同步':'完全同步')}}</h2>
   <p>{{stageText}}</p>
   <progress :value="indeterminate?undefined:syncPercent" max="100"></progress>
   <small v-if="syncCurrent>0||syncTotal>0" class="sync-completed">{{completedText}}</small>
  </div>
 </section>
</template>
