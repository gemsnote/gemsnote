<script setup lang="ts">
import {computed,onMounted,onUnmounted,ref} from 'vue'
import {request} from '../api'
import {t} from '../i18n'
import {confirmAction} from '../dialogs'
const props=defineProps<{userId:string;currentEmail:string;verified:boolean;pendingEmail?:string}>()
const emit=defineEmits<{pending:[email:string]}>()
const email=ref(''),password=ref(''),state=ref(''),error=ref(''),taskId=ref(''),busy=ref(false),confirming=ref(false)
let timer:ReturnType<typeof setTimeout>|undefined
const abort=new AbortController()
const key=`gemsnote:email-task:${props.userId}`
const messages:Record<string,string>={
 queued:'验证邮件已排队，正在等待发送。邮箱尚未修改。',
 sending:'正在发送验证邮件，请稍候。',
 sent:'验证邮件已交给邮件服务器，请检查收件箱和垃圾邮件。点击邮件中的链接后才会生效。',
 failed:'验证邮件发送失败，邮箱没有修改。请检查配置后重试。'
}
const description=computed(()=>t(messages[state.value]||''))
function remember(){try{if(taskId.value)sessionStorage.setItem(key,taskId.value);else sessionStorage.removeItem(key)}catch{}}
async function check(){
 if(!taskId.value)return
 busy.value=true;error.value=''
 try{
  const result=await request('/user/emailStatus',{taskId:taskId.value},{signal:abort.signal})
  state.value=result.State
  if(result.State==='sent'||result.State==='failed'){
   error.value=result.Error?t(result.Error):'';taskId.value='';remember();busy.value=false
  }else{timer=setTimeout(check,2000)}
 }catch(e){if(!abort.signal.aborted){error.value=e instanceof Error?e.message:String(e);busy.value=false; if(error.value===t('mailTaskNotFound')){state.value='failed';taskId.value='';remember()}}}
}
async function send(current=false){
 if(busy.value||taskId.value)return
 if(!current&&email.value.trim().toLowerCase()===props.currentEmail.toLowerCase()){error.value=t('emailUnchanged');return}
 busy.value=true;error.value='';state.value=''
 try{
  const result=await request('/web/emailChange',current?{current:true}:{email:email.value.trim(),pwd:password.value},{signal:abort.signal})
  emit('pending',result.PendingEmail||'');password.value='';taskId.value=result.TaskId;state.value='queued';remember();timer=setTimeout(check,1000)
 }catch(e){if(!abort.signal.aborted){error.value=e instanceof Error?e.message:String(e);busy.value=false}}
}
async function verifyCurrentEmail(){
 if(busy.value||taskId.value||confirming.value||props.verified)return
 confirming.value=true
 try{
  if(await confirmAction(t('是否向 {email} 重发验证邮件？',{email:props.currentEmail}))&&!abort.signal.aborted&&!props.verified)await send(true)
 }finally{confirming.value=false}
}
onMounted(()=>{try{taskId.value=sessionStorage.getItem(key)||''}catch{};if(taskId.value)void check()})
onUnmounted(()=>{clearTimeout(timer);abort.abort()})
</script>
<template>
 <div class="email-verification">
  <p class="current-email"><span>{{t('当前邮箱')}}：{{currentEmail}}</span><button v-if="!verified" type="button" class="verification-link" :disabled="busy||!!taskId||confirming" :title="t('重发验证邮件')" @click="verifyCurrentEmail">{{t('邮箱未验证')}}</button></p>
  <p v-if="pendingEmail" class="pending-email">{{t('待验证的新邮箱')}}：<strong>{{pendingEmail}}</strong></p>
  <p class="muted">{{t('修改邮箱需要验证新地址。验证成功前，继续使用当前邮箱登录。')}}</p>
  <form @submit.prevent="send()">
   <label>{{t('新邮箱')}}<input v-model="email" type="email" required :disabled="busy||!!taskId||confirming" autocomplete="email"></label>
   <label>{{t('当前密码')}}<input v-model="password" type="password" required :disabled="busy||!!taskId||confirming" autocomplete="current-password"></label>
   <button :disabled="busy||!!taskId||confirming">{{t('发送邮箱验证邮件')}}</button>
  </form>
  <p v-if="description" role="status" :class="{failed:state==='failed',success:state==='sent'}">{{description}}</p>
  <p v-if="error" role="alert" class="failed">{{error}}</p>
  <button v-if="taskId&&!busy" @click="check">{{t('重新查询发送结果')}}</button>
 </div>
</template>
<style scoped>
.current-email{display:flex;align-items:baseline;gap:.5rem;flex-wrap:wrap;overflow-wrap:anywhere}
.verification-link{padding:0;border:0;background:none;color:#246f58;text-decoration:underline;font-size:inherit}
.failed{color:#b42318}.success{color:#176b46}
</style>
