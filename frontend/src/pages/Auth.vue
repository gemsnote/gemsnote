<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { request } from '../api'
import brandMark from '../assets/gemsnote_s.png'
import {t} from '../i18n'
import {showNotice} from '../dialogs'
const route=useRoute(), router=useRouter()
const email=ref(''), pwd=ref(''), captcha=ref(''), error=ref(''), message=ref(''), busy=ref(false), busyAction=ref<'login'|'sync'|'other'>('other'), needCaptcha=ref(false), openRegister=ref(false), captchaUrl=ref('/captcha/get'), isDesktop=ref(false), server=ref(''), syncChoice=ref(false)
const mode=computed(()=>route.path.startsWith('/findPassword')?'reset':route.path==='/register'&&!isDesktop.value?'register':'login')
const submitLabel=computed(()=>{if(busy.value){if(busyAction.value==='login')return t('正在登录…');if(busyAction.value==='sync')return t('正在同步…');return t('处理中…')}return mode.value==='login'?t('登录'):t('继续')})
onMounted(async()=>{ try{const data=await request('/web/bootstrap');openRegister.value=data.OpenRegister;needCaptcha.value=data.NeedCaptcha;isDesktop.value=!!data.Desktop;server.value=data.Host||'';if(isDesktop.value&&route.path==='/register')await router.replace('/login');else if(data.User)await router.replace('/note')}catch(e){error.value=String(e)} })
async function submit(){if(busy.value)return;busy.value=true;busyAction.value=mode.value==='login'?'login':'other';error.value='';message.value='';try{
 if(mode.value==='login'){const result:any=await request('/doLogin',{email:email.value,pwd:pwd.value,captcha:captcha.value,host:desktopHost()});if(result?.Notice){const notices:any={serverMigrationRequired:t('当前连接的是旧版 Leanote 服务端，请先迁移到 Gemsnote 服务端。'),clientUpgradeRequired:t('当前客户端版本过低，请先升级客户端。'),serverUpgradeRequired:t('当前服务端版本过低，请先升级服务端。')};message.value=notices[result.Notice]||result.Notice;await showNotice(message.value)}if(result?.InitialSyncError)throw new Error(String(result.InitialSyncError));if(isDesktop.value&&result?.SyncChoiceRequired){syncChoice.value=true;return}if(isDesktop.value&&result?.ResetSyncRequired){await resetAfterLogin(true);return}await router.replace('/note')}
 else if(mode.value==='register'){await request('/doRegister',{email:email.value,pwd:pwd.value,iu:String(route.query.iu||''),host:desktopHost()});await router.replace('/note')}
 else if(route.params.token){await request('/findPasswordUpdate',{token:String(route.params.token),pwd:pwd.value,host:desktopHost()});message.value=t('密码已修改，请重新登录。')}
 else {await request('/doFindPassword',{email:email.value,host:desktopHost()});message.value=t('请检查邮箱中的重置密码链接。')}
}catch(e){error.value=String(e);const data=await request('/web/bootstrap').catch(()=>({}));needCaptcha.value=data.NeedCaptcha}finally{busy.value=false}}
function desktopHost(){return isDesktop.value?server.value.trim():undefined}
async function resetAfterLogin(initial:unknown=false){
 if(busy.value&&busyAction.value==='sync')return
 const restoreChoice=syncChoice.value
 busy.value=true;busyAction.value='sync';error.value='';syncChoice.value=false
 // Open before the reset preflight request, which can wait on the network.
 // The choice dialog has a higher stacking order and must be closed first.
 window.dispatchEvent(new CustomEvent('sync-progress',{detail:{Stage:'start',Current:0,Total:0,Percent:0,Mode:'reset'}}))
 try{
  // Only an explicit click is confirmation to discard a cache. The initial
  // path rechecks emptiness in Go under the sync lock, before any deletion.
  await request('/web/resetSync',initial===true?{initial:true}:{confirm:true})
  window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:true,Full:true,Reset:true}}))
  await router.replace('/note')
 }catch(e){
  error.value=String(e)
  syncChoice.value=restoreChoice
  if(initial===true&&e instanceof Error&&e.message==='confirmationRequired'){
   error.value=''
   syncChoice.value=true
  }
  window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:error.value,Full:true,Reset:true}}))
 }finally{busy.value=false}
}
async function skipInitialSync(){syncChoice.value=false;await router.replace('/note')}
</script>
<template><main class="auth-shell"><section class="auth-card"><div class="brand"><img :src="brandMark" :alt="t('珠玑笔记')"> {{t('珠玑笔记')}}</div><p class="muted">{{t('Gemsnote · 记录、整理、沉淀')}}</p><h1>{{t(mode==='login'?'欢迎回来':mode==='register'?'创建账号':'找回密码')}}</h1><form @submit.prevent="submit"><label v-if="isDesktop">{{t('服务器')}}<input v-model="server" placeholder="https://your-gemsnote-server" autocomplete="url"></label><label v-if="!route.params.token">{{t('邮箱或用户名')}}<input v-model="email" autocomplete="username" required></label><label v-if="mode!=='reset'||route.params.token">{{t('密码')}}<input v-model="pwd" type="password" :autocomplete="mode==='login'?'current-password':'new-password'" required></label><label v-if="needCaptcha&&mode==='login'"><span>{{t('验证码')}}</span><img :src="captchaUrl" :alt="t('验证码')"><input v-model="captcha" required></label><p v-if="error" role="alert" class="error">{{error}}</p><p v-if="message" role="status">{{message}}</p><button class="primary" :disabled="busy">{{submitLabel}}</button></form><nav class="inline"><RouterLink v-if="!isDesktop" to="/login">{{t('登录')}}</RouterLink><RouterLink v-if="!isDesktop&&openRegister" to="/register">{{t('注册')}}</RouterLink><RouterLink to="/findPassword">{{t('找回密码')}}</RouterLink></nav></section></main><section v-if="syncChoice" class="choice-dialog-backdrop"><div class="choice-dialog" role="dialog" aria-modal="true" :aria-label="t('选择登录后的同步方式')"><h2>{{t('选择登录后的同步方式')}}</h2><p>{{t('检测到此账号已有本地缓存。重新同步会删除该账号的全部本地数据，并从服务端重新下载；未上传的修改会永久丢失。')}}</p><p v-if="error" role="alert" class="error">{{error}}</p><div class="choice-dialog-actions"><button :disabled="busy" @click="skipInitialSync">{{t('暂不同步')}}</button><button class="primary" :disabled="busy" @click="resetAfterLogin">{{t('重新同步')}}</button></div></div></section></template>
