<script setup lang="ts">
import {computed,onBeforeUnmount,onMounted,ref} from 'vue'
import {request} from '../api'
import {ApiError,errorMessage} from '../errors'
import {t} from '../i18n'
import {runNoteBatch,type BatchAction,type BatchNote} from '../noteBatch'
import NotebookPicker from './NotebookPicker.vue'
import ShareRecipientInput from './ShareRecipientInput.vue'

const props=defineProps<{action:BatchAction;notes:BatchNote[];userId:string;notebooks:any[];suggestions:string[];trash:boolean}>()
const emit=defineEmits<{close:[];changed:[ids:string[],action:BatchAction]}>()
const destination=ref(''),recipient=ref(''),busy=ref(false),done=ref(0),error=ref(''),attempted=ref(false),dialog=ref<HTMLElement>()
const completed=new Set<string>(),controller=new AbortController()
const heading=computed(()=>t(props.action==='move'?'批量移动笔记':props.action==='share'?'批量共享笔记':'批量删除笔记'))
const canStart=computed(()=>!busy.value&&(props.action==='move'?!!destination.value:props.action==='share'?!!recipient.value.trim():true))
async function api(path:string,params:any){
 const timeout=new AbortController(),abort=()=>timeout.abort(),timer=setTimeout(abort,30000)
 controller.signal.addEventListener('abort',abort,{once:true})
 if(controller.signal.aborted)abort()
 try{return await request(path,params,{signal:timeout.signal})}
 catch(e){if(timeout.signal.aborted&&!controller.signal.aborted)throw new ApiError('requestTimeout');throw e}
 finally{clearTimeout(timer);controller.signal.removeEventListener('abort',abort)}
}
async function start(){
 if(!canStart.value)return
 busy.value=true;attempted.value=true;error.value='';dialog.value?.focus()
 const previous=new Set(completed)
 try{
  const failures=await runNoteBatch(props.action,props.notes,props.userId,api,completed,{destination:destination.value,recipient:recipient.value,trash:props.trash},()=>done.value=completed.size,controller.signal)
  if(controller.signal.aborted)return
  if(failures.length)error.value=t('部分笔记操作失败，已完成 {done}/{total}；重试仅处理失败项。',{done:completed.size,total:props.notes.length})+' '+errorMessage(failures[0].error)
  const changed=[...completed].filter(id=>!previous.has(id))
  if(changed.length)emit('changed',changed,props.action)
  if(completed.size===props.notes.length)emit('close')
 }catch(e){error.value=errorMessage(e)}finally{busy.value=false}
}
function trapFocus(event:KeyboardEvent){
 if(event.key!=='Tab')return
 const controls=[...dialog.value!.querySelectorAll<HTMLElement>('input:not(:disabled),button:not(:disabled):not([hidden])')]
 if(!controls.length){event.preventDefault();dialog.value?.focus();return}
 const index=controls.indexOf(document.activeElement as HTMLElement)
 if(event.shiftKey&&index<=0){event.preventDefault();controls.at(-1)!.focus()}
 else if(!event.shiftKey&&(index<0||index===controls.length-1)){event.preventDefault();controls[0].focus()}
}
const previousFocus=document.activeElement as HTMLElement|null
onMounted(()=>dialog.value?.querySelector<HTMLElement>('input:not(:disabled),button:not(:disabled):not([hidden])')?.focus())
onBeforeUnmount(()=>{controller.abort();if(previousFocus?.isConnected)previousFocus.focus()})
</script>
<template>
<section class="dialog-backdrop batch-backdrop" @click.self="!busy&&emit('close')" @keydown.esc.stop.prevent="!busy&&emit('close')">
<div ref="dialog" class="dialog batch-dialog" role="dialog" aria-modal="true" tabindex="-1" :aria-label="heading" :aria-busy="busy" @keydown="trapFocus">
<h2>{{heading}}</h2>
<p>{{t('已选中 {count} 篇笔记',{count:notes.length})}}</p>
<ul class="batch-titles"><li v-for="note in notes" :key="note.NoteId" :title="note.Title">{{note.Title||t('未命名')}}</li></ul>
<NotebookPicker v-if="action==='move'" v-model="destination" :items="notebooks" :disabled="busy||attempted"/>
<form v-else-if="action==='share'" @submit.prevent="start">
<label>{{t('用户名或邮箱')}}<ShareRecipientInput v-model="recipient" :suggestions="suggestions" :disabled="busy||attempted"/></label>
<p class="muted">{{t('所选笔记将以只读方式共享。')}}</p>
<button type="submit" hidden :disabled="!canStart">{{t('共享')}}</button>
</form>
<p v-else>{{t(trash?'永久删除所选笔记？此操作不可撤销。':'将所选笔记放入回收站？')}}</p>
<p v-if="busy||done" role="status">{{done}} / {{notes.length}}</p>
<p v-if="error" class="error" role="alert">{{error}}</p>
<div class="choice-dialog-actions">
<button :disabled="!canStart" class="primary" @click="start">{{t(error?'重试':action==='move'?'移动':action==='share'?'共享':'删除')}}</button>
<button :disabled="busy" @click="emit('close')">{{t('取消')}}</button>
</div>
</div></section>
</template>
<style scoped>
.batch-backdrop{z-index:105}.batch-titles{max-height:130px;overflow:auto;padding-left:24px;font-size:12px}.batch-titles li{overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.batch-dialog .error{margin:12px 0}
</style>
