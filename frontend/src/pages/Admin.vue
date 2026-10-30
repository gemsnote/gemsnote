<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { request } from '../api'
import Navigation from '../components/Navigation.vue'
import {t} from '../i18n'
import {confirmAction} from '../dialogs'

const roleBusy=ref(false)
const data = ref<any>({ Users: [], Settings: {} })
const router = useRouter()
const user = ref<any>({})
const message = ref(''), keywords = ref(''), page = ref(1)
const username = ref(''), email = ref(''), resetId = ref(''), resetPwd = ref('')

const resendBusy=ref(false)
const verification=ref<Record<string,{state:string;error:string;taskId:string}>>({})
const verificationTimers=new Map<string,ReturnType<typeof setTimeout>>()
async function checkVerification(userId:string){
 const task=verification.value[userId]
 if(!task?.taskId)return
 try{
  const status=await request('/user/emailStatus',{userId,taskId:task.taskId},{signal:mailAbort.signal})
  if(verification.value[userId]!==task)return
  task.state=status.State;task.error=status.Error?t(status.Error):''
  if(status.State==='queued'||status.State==='sending')verificationTimers.set(userId,setTimeout(()=>checkVerification(userId),2000))
  else verificationTimers.delete(userId)
 }catch(error){if(!mailAbort.signal.aborted){task.state='failed';task.error=String(error)}}
}
async function resendVerification(target:any){
 if(resendBusy.value||target.Verified||target.Flag==='D'||['queued','sending'].includes(verification.value[target.UserId]?.state))return
 resendBusy.value=true
 try{
  if(!await confirmAction(t('是否向 {email} 重发验证邮件？',{email:target.Email})))return
  if(mailAbort.signal.aborted)return
  const result=await request('/api2/admin/resendVerification',{userId:target.UserId},{signal:mailAbort.signal})
  verification.value[target.UserId]={state:'queued',error:'',taskId:result.TaskId}
  verificationTimers.set(target.UserId,setTimeout(()=>checkVerification(target.UserId),1000))
 }catch(error){if(!mailAbort.signal.aborted)verification.value[target.UserId]={state:'failed',error:String(error),taskId:''}}
 finally{resendBusy.value=false}
}

const settingsSaving=ref(false),settingsResult=ref(''),settingsError=ref('')
async function saveSettings(){
 if(settingsSaving.value)return
 settingsSaving.value=true;settingsResult.value='';settingsError.value=''
 try{
  await request('/web/adminSettings',{...data.value.Settings})
  data.value.Settings.emailPassword=''
  settingsResult.value='设置已保存'
 }catch(error){settingsError.value=error instanceof Error?error.message:String(error)}
 finally{settingsSaving.value=false}
}

const creating = ref(false)
const created = ref<{username:string;email:string;userId:string;taskId:string;state:string;error:string}|null>(null)
let mailTimer:ReturnType<typeof setTimeout>|undefined
const mailAbort=new AbortController()
onUnmounted(()=>{clearTimeout(mailTimer);verificationTimers.forEach(clearTimeout);mailAbort.abort()})
function dismissInvitation(){created.value=null;clearTimeout(mailTimer)}
async function checkInvitation() {
 const item=created.value
 if(!item?.taskId)return
 try {
  const status=await request('/user/emailStatus',{taskId:item.taskId,userId:item.userId},{signal:mailAbort.signal})
  if(created.value!==item)return
  item.state=status.State;item.error=status.Error||''
  if(status.State==='queued'||status.State==='sending')mailTimer=setTimeout(checkInvitation,2000)
 }catch(error){if(!mailAbort.signal.aborted&&created.value===item){item.state='failed';item.error=String(error)}}
}
async function createUser() {
 if(creating.value)return
 creating.value=true
 try {
  const name=username.value.trim(),address=email.value.trim()
  const result=await request('/web/adminRegister',{username:name,email:address})
  clearTimeout(mailTimer)
  created.value={username:name,email:address,userId:result.UserId,taskId:result.TaskId,state:result.MailError?'failed':'queued',error:result.MailError||''}
  username.value='';email.value=''
  message.value=t('用户已创建，请用户通过验证邮件设置密码。')
  if(result.TaskId)mailTimer=setTimeout(checkInvitation,1000)
  await load()
 }catch(error){message.value=String(error)}
 finally{creating.value=false}
}

async function load() {
  try { data.value = await request('/web/adminData', { keywords: keywords.value, page: page.value }) }
  catch (error) { message.value = String(error) }
}
async function action(path: string, params: any) {
  try { await request(path, params); message.value = t('操作成功'); await load() }
  catch (error) { message.value = String(error) }
}
async function setRole(target:any,event:Event){
 const select=event.target as HTMLSelectElement
 const role=select.value
 if(roleBusy.value||target.UserId===user.value.UserId||role===target.Role){select.value=target.Role;return}
 roleBusy.value=true
 try{
  await request('/api2/admin/role',{userId:target.UserId,role})
  target.Role=role
  message.value=t('操作成功')
  await load()
 }catch(error){message.value=String(error)}
 finally{select.value=target.Role;roleBusy.value=false}
}
async function setFlag(target:any){
 if(roleBusy.value||target.UserId===user.value.UserId)return
 roleBusy.value=true
 try{await action('/api2/admin/flag',{userId:target.UserId,flag:target.Flag==='D'?'A':'D'})}
 finally{roleBusy.value=false}
}
onMounted(async()=>{try{const bootstrap=await request('/web/bootstrap');if(bootstrap.Desktop){await router.replace('/note');return}user.value=bootstrap.User||{};await load()}catch(error){message.value=String(error)}})
</script>

<template>
  <div class="shell settings-shell"><Navigation admin :user="user" back-only/><main class="settings">
    <h1>{{t('系统管理')}}</h1><p role="status">{{ message }}</p>
    <section class="card"><h2>{{t('用户管理')}}</h2>
      <form class="inline" @submit.prevent="page=1;load()"><input v-model="keywords" :placeholder="t('搜索用户')"><button>{{t('搜索')}}</button></form>
      <table><thead><tr><th>{{t('用户名')}}</th><th>{{t('邮箱')}}</th><th>{{t('角色')}}</th><th>{{t('状态')}}</th><th>{{t('操作')}}</th></tr></thead><tbody>
        <tr v-for="u in data.Users" :key="u.UserId">
          <td>{{u.Username}}</td><td>{{u.Email}}<small v-if="u.Verified" class="email-state">{{t('邮箱已验证')}}</small><button v-else type="button" class="email-state email-link" :disabled="resendBusy||u.Flag==='D'||['queued','sending'].includes(verification[u.UserId]?.state)" :title="t('重发验证邮件')" @click="resendVerification(u)">{{t('邮箱未验证')}}</button>
            <small v-if="verification[u.UserId]" class="email-state" :class="{'email-error':verification[u.UserId].state==='failed'}" role="status">{{verification[u.UserId].state==='failed'?verification[u.UserId].error:t(verification[u.UserId].state==='sent'?'验证邮件已发送':'验证邮件正在发送，请稍候。')}}</small><small v-if="u.PendingEmail" class="email-state">{{t('待验证的新邮箱')}}：{{u.PendingEmail}}</small></td>
          <td><select :value="u.Role==='admin'?'admin':'user'" :aria-label="`${u.Username} · ${t('角色')}`" :disabled="roleBusy||u.UserId===user.UserId" :title="u.UserId===user.UserId?t('不能修改自己的管理员权限'):''" @change="setRole(u,$event)"><option value="user">{{t('普通用户')}}</option><option value="admin">{{t('管理员')}}</option></select></td>
          <td><button type="button" class="status-toggle" :class="u.Flag==='D'?'is-disabled':'is-available'" :disabled="roleBusy||u.UserId===user.UserId" :aria-label="t(u.Flag==='D'?'启用用户':'禁用用户')" :title="u.UserId===user.UserId?t('不能禁用自己的账号'):`${t(u.Flag==='D'?'已禁用':'可用')} · ${t(u.Flag==='D'?'启用用户':'禁用用户')}`" @click="setFlag(u)">
            <svg viewBox="0 0 24 24" aria-hidden="true"><path v-if="u.Flag==='D'" d="m6 6 12 12M18 6 6 18"/><path v-else d="m4 12 5 5L20 6"/></svg>
          </button></td>
          <td><button @click="resetId=u.UserId">{{t('重置密码')}}</button></td>
        </tr>
      </tbody></table>
      <div class="inline"><button :disabled="page===1" @click="page--;load()">{{t('上一页')}}</button><span>{{page}}</span><button :disabled="data.Users.length<20" @click="page++;load()">{{t('下一页')}}</button></div>
      <form v-if="resetId" class="inline" @submit.prevent="action('/web/adminResetPwd',{userId:resetId,pwd:resetPwd})"><label>{{t('新密码')}}<input v-model="resetPwd" type="password" required></label><button>{{t('确认重置')}}</button><button type="button" @click="resetId=''">{{t('取消')}}</button></form>
    </section>
    <section class="card"><h2>{{t('添加用户（管理）')}}</h2>
      <form @submit.prevent="createUser">
        <label>{{t('用户名')}}<input v-model.trim="username" required minlength="4" pattern="[A-Za-z0-9_\-]+" autocomplete="off" :disabled="creating"></label>
        <small>{{t('用户名至少4位，只能包含字母、数字、下划线和连字符')}}</small>
        <label>{{t('邮箱')}}<input v-model.trim="email" type="email" required :disabled="creating"></label>
        <small class="invitation-help">{{t('用户验证邮箱后自行设置密码，无需管理员指定密码。')}}</small>
        <button class="create-user-button" :disabled="creating">{{t('创建用户')}}</button>
      </form>
      <div v-if="created" class="created-invitation" role="status">
        <p>{{created.username}} · {{created.email}}</p>
        <p v-if="created.state==='queued'||created.state==='sending'">{{t('验证邮件正在发送，请稍候。')}}</p>
        <p v-else-if="created.state==='sent'">{{t('验证邮件已发送，请用户点击链接验证邮箱并设置密码。')}}</p>
        <template v-else-if="created.state==='failed'"><p class="error">{{t('账号已创建，但验证邮件发送失败。')}} {{t(created.error)}}</p><p>{{t('请检查邮件配置。用户可通过登录页的找回密码重新获取链接，无需重复创建账号。')}}</p></template>
        <button type="button" @click="dismissInvitation">{{t('关闭')}}</button>
      </div>
    </section>
    <section class="card"><h2>{{t('站点与邮件')}}</h2><form @submit.prevent="saveSettings" @input="settingsResult='';settingsError=''" :aria-busy="settingsSaving"><fieldset class="settings-fields" :disabled="settingsSaving">
      <label>{{t('站点 URL')}}<input v-model="data.Settings.siteUrl" type="url"></label>
      <label>{{t('邮件服务器')}}<input v-model="data.Settings.emailHost"></label><label>{{t('邮件端口')}}<input v-model="data.Settings.emailPort"></label>
      <label>{{t('邮件账号')}}<input v-model="data.Settings.emailUsername"></label><label>{{t('邮件密码')}}<input v-model="data.Settings.emailPassword" type="password" :placeholder="t('留空保持原值')"></label>
      <label>{{t('邮件 SSL')}}<select v-model="data.Settings.emailSSL"><option value="1">{{t('启用')}}</option><option value="0">{{t('禁用')}}</option></select></label>
      <h3>{{t('上传限制（MB）')}}</h3>
      <label>{{t('图片')}}<input v-model="data.Settings.uploadImageSize" type="number" min="0" step="0.1"></label><label>{{t('头像')}}<input v-model="data.Settings.uploadAvatarSize" type="number" min="0" step="0.1"></label><label>{{t('附件')}}<input v-model="data.Settings.uploadAttachSize" type="number" min="0" step="0.1"></label>
      <h3>{{t('其他非博客设置')}}</h3>
      <label>{{t('PDF 导出程序路径')}}<input v-model="data.Settings.exportPdfBinPath"></label>
      <button class="primary" :disabled="settingsSaving">{{t(settingsSaving?'正在保存…':'保存设置')}}</button>
      </fieldset>
      <p v-if="settingsResult" role="status" class="settings-saved">{{t(settingsResult)}}</p>
      <p v-if="settingsError" role="alert" class="error">{{t('设置保存失败')}}：{{settingsError}}</p>
    </form></section>
  </main></div>
</template>

<style scoped>
.settings-fields{border:0;padding:0;margin:0;min-width:0}
.settings-saved{color:#176b46;background:#e4f2e6;border-radius:6px;padding:10px;margin-bottom:0}
.email-state{display:block;margin-top:4px;overflow-wrap:anywhere}
.email-link{padding:0;border:0;background:transparent;color:#246f58;text-decoration:underline;font-size:12px;text-align:left}
.email-error{color:#b42318}
.status-toggle{display:inline-grid;place-items:center;width:36px;height:36px;padding:7px;background:transparent}
.status-toggle.is-available{color:#238044}
.status-toggle.is-disabled{color:#c53030}
.status-toggle svg{width:20px;height:20px;fill:none;stroke:currentColor;stroke-width:2.5;stroke-linecap:round;stroke-linejoin:round}
.invitation-help{display:block;margin-top:6px}
.create-user-button{display:block;margin-top:12px}
</style>
