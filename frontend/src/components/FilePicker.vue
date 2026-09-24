<script setup lang="ts">
import { ref } from 'vue'
import { t } from '../i18n'

defineProps<{label:string,accept?:string}>()
const emit=defineEmits<{selected:[file:File]}>()
const input=ref<HTMLInputElement>(),filename=ref('')
function selected(event:Event){
 const target=event.target as HTMLInputElement,file=target.files?.[0]
 // Allow retrying the same file after a failed upload.
 target.value=''
 if(!file)return
 filename.value=file.name
 emit('selected',file)
}
</script>

<template>
 <div class="file-picker" role="group" :aria-label="label">
  <span>{{label}}</span>
  <div class="inline">
   <button type="button" @click="input?.click()">{{t('选择文件')}}</button>
   <span class="filename" :title="filename" role="status">{{filename||t('未选择文件')}}</span>
  </div>
  <input ref="input" type="file" hidden :accept="accept" :aria-label="label" @change="selected">
 </div>
</template>

<style scoped>
.file-picker{margin:12px 0;display:grid;gap:6px}
.filename{min-width:0;max-width:100%;overflow-wrap:anywhere}
</style>
