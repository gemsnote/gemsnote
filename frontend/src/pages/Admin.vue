<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { request } from '../api'
import Navigation from '../components/Navigation.vue'
import {t} from '../i18n'

const roleBusy=ref(false)
const data = ref<any>({ Users: [], Settings: {} })
const router = useRouter()
const user = ref<any>({})
const message = ref(''), keywords = ref(''), page = ref(1)
const email = ref(''), pwd = ref(''), resetId = ref(''), resetPwd = ref('')

async function load() {
  try { data.value = await request('/web/adminData', { keywords: keywords.value, page: page.value }) }
  catch (error) { message.value = String(error) }
}
async function action(path: string, params: any) {
  try { await request(path, params); message.value = t('操作成功'); await load() }
  catch (error) { message.value = String(error) }
}
async function setRole(target:any){
 if(roleBusy.value||target.UserId===user.value.UserId)return
 roleBusy.value=true
 try{await action('/api2/admin/role',{userId:target.UserId,role:target.Role==='admin'?'user':'admin'})}
 finally{roleBusy.value=false}
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
        <tr v-for="u in data.Users" :key="u.UserId"><td>{{u.Username}}</td><td>{{u.Email}}</td><td>{{t(u.Role==='admin'?'管理员':'普通用户')}}</td><td>{{t(u.Flag==='D'?'已禁用':'可用')}}</td><td><button :disabled="roleBusy||u.UserId===user.UserId" :title="u.UserId===user.UserId?t('不能禁用自己的账号'):''" @click="setFlag(u)">{{t(u.Flag==='D'?'启用用户':'禁用用户')}}</button><button :disabled="roleBusy||u.UserId===user.UserId" :title="u.UserId===user.UserId?t('不能修改自己的管理员权限'):''" @click="setRole(u)">{{t(u.Role==='admin'?'设为普通用户':'设为管理员')}}</button><button @click="resetId=u.UserId">{{t('重置密码')}}</button></td></tr>
      </tbody></table>
      <div class="inline"><button :disabled="page===1" @click="page--;load()">{{t('上一页')}}</button><span>{{page}}</span><button :disabled="data.Users.length<20" @click="page++;load()">{{t('下一页')}}</button></div>
      <form v-if="resetId" class="inline" @submit.prevent="action('/web/adminResetPwd',{userId:resetId,pwd:resetPwd})"><label>{{t('新密码')}}<input v-model="resetPwd" type="password" required></label><button>{{t('确认重置')}}</button><button type="button" @click="resetId=''">{{t('取消')}}</button></form>
    </section>
    <section class="card"><h2>{{t('添加用户（管理）')}}</h2><form @submit.prevent="action('/web/adminRegister',{email,pwd})"><label>{{t('邮箱')}}<input v-model="email" type="email" required></label><label>{{t('初始密码')}}<input v-model="pwd" type="password" required></label><button>{{t('创建')}}</button></form></section>
    <section class="card"><h2>{{t('站点与邮件')}}</h2><form @submit.prevent="action('/web/adminSettings',data.Settings)">
      <label>{{t('站点 URL')}}<input v-model="data.Settings.siteUrl" type="url"></label>
      <label>{{t('邮件服务器')}}<input v-model="data.Settings.emailHost"></label><label>{{t('邮件端口')}}<input v-model="data.Settings.emailPort"></label>
      <label>{{t('邮件账号')}}<input v-model="data.Settings.emailUsername"></label><label>{{t('邮件密码')}}<input v-model="data.Settings.emailPassword" type="password" :placeholder="t('留空保持原值')"></label>
      <label>{{t('邮件 SSL')}}<select v-model="data.Settings.emailSSL"><option value="1">{{t('启用')}}</option><option value="0">{{t('禁用')}}</option></select></label>
      <h3>{{t('上传限制（MB）')}}</h3>
      <label>{{t('图片')}}<input v-model="data.Settings.uploadImageSize" type="number" min="0" step="0.1"></label><label>{{t('头像')}}<input v-model="data.Settings.uploadAvatarSize" type="number" min="0" step="0.1"></label><label>{{t('附件')}}<input v-model="data.Settings.uploadAttachSize" type="number" min="0" step="0.1"></label>
      <h3>{{t('其他非博客设置')}}</h3>
      <label>{{t('PDF 导出程序路径')}}<input v-model="data.Settings.exportPdfBinPath"></label>
      <button class="primary">{{t('保存设置')}}</button>
    </form></section>
  </main></div>
</template>
