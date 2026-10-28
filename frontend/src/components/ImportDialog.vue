<script setup lang="ts">
import {computed,onBeforeUnmount,ref,shallowRef} from 'vue'
import {objectId,request} from '../api'
import type {Params} from '../api'
import {t} from '../i18n'
import {importLimits,ImportError,type ImportSource} from '../markdownImport'
import {createImportJob,runImportJob,type ImportJob} from '../importJob'
import FilePicker from './FilePicker.vue'
import NotebookPicker from './NotebookPicker.vue'

const props=defineProps<{items:any[];selected:string;userId:string;desktop:boolean}>()
const emit=defineEmits<{close:[];changed:[job:ImportJob,complete:boolean]}>()
const source=shallowRef<ImportSource>(),job=shallowRef<ImportJob>(),destination=ref(props.selected),name=ref(''),busy=ref(false),reading=ref(false),error=ref(''),done=ref(0),stepTitle=ref(''),stop=ref(false),discard=ref(false)
let worker:Worker|undefined,timer:ReturnType<typeof setTimeout>|undefined,sequence=0,disposed=false
const complete=computed(()=>!!job.value&&done.value===job.value.steps.length)
const pending=computed(()=>busy.value||!!job.value&&!complete.value)
defineExpose({pending})
function cancelWorker(){worker?.terminate();worker=undefined;clearTimeout(timer)}
function message(e:unknown){return e instanceof ImportError?t(e.key):String(e)}
async function selected(file:File){
 if(busy.value||job.value)return
 const token=++sequence;cancelWorker();source.value=undefined;error.value='';reading.value=true
 try{
  if(!/\.(md|markdown|zip)$/i.test(file.name))throw new ImportError('请选择 Markdown 或 ZIP 文件')
  if(file.size>(/\.zip$/i.test(file.name)?importLimits.archive:importLimits.note))throw new ImportError(/\.zip$/i.test(file.name)?'ZIP 文件大小不能超过 32 MiB':'Markdown 内容大小超出导入限制')
  const bytes=new Uint8Array(await file.arrayBuffer())
  if(token!==sequence)return
  worker=new Worker(new URL('../markdownImport.worker.ts',import.meta.url),{type:'module'})
  timer=setTimeout(()=>{cancelWorker();error.value=t('解析文件超时，请缩小文件后重试');reading.value=false},15000)
  worker.onerror=()=>{cancelWorker();error.value=t('无法解析文件，请重新选择');reading.value=false}
  worker.onmessage=({data})=>{cancelWorker();reading.value=false;if(data.error){error.value=t(data.error);return}source.value=data.source;name.value=data.source.title;if(!data.source.zip&&!destination.value)destination.value=props.items[0]?.NotebookId||''}
  worker.postMessage({name:file.name,bytes},[bytes.buffer])
 }catch(e){if(token===sequence){error.value=message(e);reading.value=false}}
}
async function start(){
 if(!source.value||busy.value)return
 try{if(!job.value)job.value=createImportJob(source.value,props.userId,destination.value,name.value,objectId)}catch(e){error.value=message(e);return}
 busy.value=true;stop.value=false;error.value=''
 try{await runImportJob(job.value,importRequest,()=>{done.value=job.value!.next;stepTitle.value=job.value!.steps[done.value]?.title||''},()=>stop.value)}
 catch(e){error.value=message(e)}
 finally{busy.value=false;if(!disposed)emit('changed',job.value,complete.value)}
}
async function importRequest(path:string,params?:Params){
 const controller=new AbortController(),timeout=setTimeout(()=>controller.abort(),30000)
 try{return await request(path,params,{signal:controller.signal})}
 catch(e){if(controller.signal.aborted)throw new ImportError('导入请求超时，请继续导入以确认进度');throw e}
 finally{clearTimeout(timeout)}
}
function close(){if(busy.value)return;if(job.value&&!complete.value){discard.value=true;return}emit('close')}
onBeforeUnmount(()=>{disposed=true;sequence++;stop.value=true;cancelWorker()})
</script>
<template>
<section class="dialog-backdrop import-backdrop">
<div class="dialog import-dialog" role="dialog" aria-modal="true" :aria-label="t('导入 Markdown')" :aria-busy="busy||reading">
<h2>{{t('导入 Markdown')}}</h2>
<p>{{t('单个 MD 导入为笔记；ZIP 导入为新笔记本，文件夹成为子笔记本。')}}</p>
<p class="muted">{{t('仅导入 UTF-8 Markdown 文本，不导入图片和附件。MD 最大 4 MiB，ZIP 最大 32 MiB。')}}</p>
<template v-if="!job">
<FilePicker :label="t('Markdown 或 ZIP 文件')" accept=".md,.markdown,.zip" @selected="selected"/>
<p v-if="reading" role="status">{{t('正在解析文件…')}}</p>
<template v-if="source">
<label v-if="source.zip">{{t('新笔记本名称')}}<input v-model="name" maxlength="200"></label>
<p>{{t('将导入 {notes} 篇笔记、{books} 个新笔记本；跳过 {ignored} 个文件。',{notes:source.notes.length,books:source.zip?source.folders.length+1:0,ignored:source.ignored})}}</p>
<p>{{t(source.zip?'选择新笔记本的父级':'选择目标笔记本')}}</p>
<NotebookPicker v-model="destination" :items="items" :root="source.zip"/>
<p v-if="!source.zip&&!items.length">{{t('请先创建一个笔记本，再导入单个 Markdown 文件。')}}</p>
</template>
</template>
<template v-else>
<progress :value="done" :max="job.steps.length" :aria-label="t('导入进度')"></progress>
<p role="status">{{done}} / {{job.steps.length}} <span v-if="busy">{{stepTitle}}</span></p>
<p v-if="complete">{{t('导入完成')}}</p>
<p v-else>{{t('已创建的内容会保留；继续导入不会重复创建。')}}</p>
<p v-if="complete&&desktop">{{t('已保存到本地，请同步以上传到服务端。')}}</p>
</template>
<p v-if="error" role="alert" class="import-error">{{error}}</p>
<template v-if="discard">
<p>{{t('关闭后将无法继续本次导入，已导入内容会保留。确认关闭？')}}</p>
<div class="choice-dialog-actions"><button @click="emit('close')">{{t('确认关闭')}}</button><button @click="discard=false">{{t('取消')}}</button></div>
</template>
<div v-else class="choice-dialog-actions">
<button v-if="busy" :disabled="stop" @click="stop=true">{{t('暂停导入')}}</button>
<button v-else-if="!complete" class="primary" :disabled="reading||!source||(!source.zip&&!destination)||(source.zip&&!name.trim())" @click="start">{{t(job?'继续导入':'开始导入')}}</button>
<button :disabled="busy" @click="close">{{t(complete?'关闭':'取消')}}</button>
</div>
</div></section>
</template>
<style scoped>
.import-backdrop{z-index:105}.import-dialog p{overflow-wrap:anywhere}.import-dialog label{display:grid;gap:6px}.import-dialog input,.import-dialog progress{width:100%;box-sizing:border-box}.import-error{color:#a53030}.import-dialog progress{accent-color:#235c48}
</style>
