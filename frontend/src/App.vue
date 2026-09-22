<script setup lang="ts">
import {computed,onBeforeUnmount,onMounted,ref} from 'vue'
import {t} from './i18n'

const syncVisible=ref(false)
const syncStage=ref('start')
const syncMode=ref<'full'|'reset'>('full')
const syncCurrent=ref(0)
const syncTotal=ref(100)
const unsubscribers:Array<()=>void>=[]
const syncPercent=computed(()=>syncTotal.value>0?Math.max(0,Math.min(100,Math.round(syncCurrent.value*100/syncTotal.value))):0)
const stageText=computed(()=>t(({start:'正在准备同步',notebooks:'正在同步笔记本',notes:'正在同步笔记',tags:'正在同步标签',push:'正在上传本地更改',images:'正在同步图片和附件',done:'同步完成'} as Record<string,string>)[syncStage.value]||'正在完全同步'))

function showSyncProgress(event:any){
 const value=event?.detail||event||{}
 if(value.Mode==='full'||value.Mode==='reset')syncMode.value=value.Mode
 syncStage.value=String(value.Stage||'start')
 syncCurrent.value=Number(value.Current||0)
 syncTotal.value=Number(value.Total||100)
 syncVisible.value=true
}
function finishSync(){syncVisible.value=false}
function finishWindowSync(event:any){const value=event?.detail||event||{};if(value.Full||value.Ok===false)finishSync()}

onMounted(()=>{
 window.addEventListener('sync-progress',showSyncProgress as EventListener)
 window.addEventListener('sync-result',finishWindowSync as EventListener)
 const wails=(window as any).runtime
 if(wails?.EventsOn){
  const offProgress=wails.EventsOn('sync-progress',showSyncProgress)
  const offFinished=wails.EventsOn('sync-finished',finishSync)
  if(typeof offProgress==='function')unsubscribers.push(offProgress)
  if(typeof offFinished==='function')unsubscribers.push(offFinished)
 }
})
onBeforeUnmount(()=>{
 window.removeEventListener('sync-progress',showSyncProgress as EventListener)
 window.removeEventListener('sync-result',finishWindowSync as EventListener)
 unsubscribers.splice(0).forEach(unsubscribe=>unsubscribe())
})
</script>

<template>
 <RouterView />
 <section v-if="syncVisible" class="sync-progress-backdrop" aria-live="polite">
  <div class="sync-progress-dialog" role="dialog" aria-modal="true" :aria-label="t(syncMode==='reset'?'重新同步':'完全同步')">
   <h2>{{t(syncMode==='reset'?'重新同步':'完全同步')}}</h2>
   <p>{{stageText}}</p>
   <progress :value="syncPercent" max="100"></progress>
   <small>{{syncPercent}}%</small>
  </div>
 </section>
</template>
