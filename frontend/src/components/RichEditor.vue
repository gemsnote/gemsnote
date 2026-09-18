<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, watch } from 'vue'
import {language,t} from '../i18n'
const props=defineProps<{modelValue:string}>()
const emit=defineEmits<{ 'update:modelValue':[value:string] }>()
const target=ref<HTMLTextAreaElement>(),error=ref('')
let editor:any
declare global { interface Window { tinymce:any } }
let active=true
onMounted(async()=>{
 try {
  if(!window.tinymce) await new Promise<void>((resolve,reject)=>{const script=document.createElement('script');script.src='/tinymce/tinymce.min.js';script.onload=()=>resolve();script.onerror=()=>reject(new Error('编辑器加载失败，请刷新重试'));document.head.append(script)})
  if(!active)return
  window.tinymce.init({target:target.value,selector:undefined,menubar:false,statusbar:false,height:500,language:language.value==='zh-CN'?'zh-cn':'en-us',relative_urls:false,convert_urls:false,valid_children:'+pre[div|#text|p|br|span]',extended_valid_elements:'pre[*],code[*]',plugins:'link image lists table code paste',toolbar:'undo redo | formatselect | bold italic underline | bullist numlist | link image table | code',setup:(instance:any)=>{editor=instance;instance.on('init',()=>instance.setContent(props.modelValue));instance.on('input change undo redo',()=>emit('update:modelValue',instance.getContent()))}})
 }catch(e){error.value=e instanceof Error?e.message:String(e)}
})
watch(()=>props.modelValue,value=>{if(editor?.initialized&&editor.getContent()!==value)editor.setContent(value)})
onBeforeUnmount(()=>{active=false;editor?.remove()})
</script>
<template><div class="rich-editor"><p v-if="error" role="alert">{{error}}</p><textarea ref="target" :value="modelValue" @input="emit('update:modelValue',($event.target as HTMLTextAreaElement).value)" :aria-label="t('富文本正文')" /></div></template>
