<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { activeDialog, closeDialog } from '../dialogs'
import { t } from '../i18n'

const dialog=ref<HTMLDialogElement>(),input=ref<HTMLInputElement>(),value=ref('')
watch(activeDialog, async request=>{
  dialog.value?.close()
  if(!request)return
  value.value=request.initialValue
  await nextTick()
  dialog.value?.showModal()
  input.value?.focus()
  input.value?.select()
})
</script>

<template>
  <dialog ref="dialog" class="choice-dialog app-dialog" :aria-label="t('珠玑笔记')" @cancel.prevent="closeDialog(false)">
    <form v-if="activeDialog" @submit.prevent="closeDialog(true,value)">
      <h2>{{t('珠玑笔记')}}</h2>
      <label v-if="activeDialog.kind==='prompt'">{{activeDialog.message}}<input ref="input" v-model="value"></label>
      <p v-else>{{activeDialog.message}}</p>
      <div class="choice-dialog-actions">
        <button v-if="activeDialog.kind!=='alert'" type="button" autofocus @click="closeDialog(false)">{{t('取消')}}</button>
        <button class="primary" type="submit">{{t('确定')}}</button>
      </div>
    </form>
  </dialog>
</template>

<style scoped>
.app-dialog{border:0;margin:auto;max-height:90vh;overflow:auto}
.app-dialog::backdrop{background:#172e3266}
</style>
