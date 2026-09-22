<script setup lang="ts">
import {computed,onMounted,ref} from 'vue'
import {request,resolveAvatarUrl} from '../api'
import {language,languages,setLanguage,t} from '../i18n'
import brandMark from '../assets/gemsnote_s.png'

const props=defineProps<{admin?:boolean,user?:{Username?:string;Email?:string;Logo?:string},workspace?:boolean,collapsed?:boolean,pending?:boolean,backOnly?:boolean,desktop?:boolean}>()
const emit=defineEmits<{synced:[]}>()
const menuOpen=ref(false)
const languageOpen=ref(false)
const syncOpen=ref(false)
const logoutPrompt=ref(false)
const resetPrompt=ref(false)
const aboutOpen=ref(false)
const aboutInfo=ref({Name:'Gemsnote',Version:'',Platform:'',Arch:'',Runtime:''})
let resolveLogoutChoice:((force:boolean)=>void)|undefined
const initials=computed(()=>(props.user?.Username||props.user?.Email||t('用户')).trim().slice(0,1).toUpperCase())
const avatarSrc=computed(()=>resolveAvatarUrl(props.user?.Logo||''))
const desktop=computed(()=>!!props.desktop)
const syncing=ref(false)

function close(){menuOpen.value=false;languageOpen.value=false;syncOpen.value=false}
function updateNativeLanguage(value:typeof language.value){
  const setNative=(window as any).go?.main?.App?.SetLanguage
  if(typeof setNative==='function')void setNative(value)
}
function chooseLanguage(value:typeof language.value){setLanguage(value);updateNativeLanguage(value);close()}
function focusOut(event:FocusEvent){if(!(event.currentTarget as HTMLElement).contains(event.relatedTarget as Node|null))close()}
function askLogoutChoice(){
  close();logoutPrompt.value=true
  return new Promise<boolean>(resolve=>{resolveLogoutChoice=resolve})
}
function chooseLogout(force:boolean){
  logoutPrompt.value=false
  const resolve=resolveLogoutChoice;resolveLogoutChoice=undefined;resolve?.(force)
}
async function showAbout(){
  close()
  const getAbout=(window as any).go?.main?.App?.GetAboutInfo
  if(typeof getAbout==='function'){
    try{aboutInfo.value={...aboutInfo.value,...await getAbout()}}catch{}
  }
  aboutOpen.value=true
}
onMounted(()=>updateNativeLanguage(language.value))
async function runSync(mode:'full'|'reset'='full'){
  if(syncing.value)return
  syncing.value=true
  close()
  window.dispatchEvent(new CustomEvent('sync-progress',{detail:{Stage:'start',Current:0,Total:100,Mode:mode}}))
  try{
    if(mode==='reset'){
      const beforeReset=(window as any).__gemsnoteBeforeReset
      if(typeof beforeReset==='function'&&!await beforeReset())throw new Error(t('本地笔记保存失败，已取消重新同步'))
      await request('/web/resetSync',{confirm:true})
      const afterReset=(window as any).__gemsnoteAfterReset
      if(typeof afterReset==='function')afterReset()
      location.href='/note'
      return
    }
    await request('/web/fullSync',{})
    emit('synced')
  }catch(e){
    const message=e instanceof Error?e.message:String(e)
    if(message!=='already syncing')window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:message,Full:true,Reset:mode==='reset'}}))
  }
  finally{syncing.value=false}
}
async function runImmediateSync(){
  if(syncing.value)return
  syncing.value=true
  try{
    if(!desktop.value){
      const syncNow=(window as any).__gemsnoteSyncNow
      if(typeof syncNow==='function'&&!await syncNow())return
    }else await request('/web/sync',{})
    emit('synced')
  }catch(e){
    const message=e instanceof Error?e.message:String(e)
    if(message!=='already syncing')window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:message,Full:false}}))
  }finally{syncing.value=false}
}
async function logout(event:MouseEvent){
  event.preventDefault()
  const beforeLogout=(window as any).__gemsnoteBeforeLogout
  if(typeof beforeLogout==='function'&&!await beforeLogout())return
  if(!desktop.value){
    try{
      await request('/web/logout',{})
      location.href='/login'
    }catch(e){
      const message=e instanceof Error?e.message:String(e)
      window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:message,Full:false}}))
    }
    return
  }
  try{
    await request('/web/logout',{})
    location.href='/login'
  }catch(e){
    const message=e instanceof Error?e.message:String(e)
    if(message==='syncFailed'){
      if(!await askLogoutChoice())return
      try{await request('/web/logout',{force:true});location.href='/login'}
      catch(forceError){window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:forceError instanceof Error?forceError.message:String(forceError),Full:false}}))}
      return
    }
    window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:message,Full:false}}))
  }
}
</script>

<template>
  <nav :class="workspace?['space-footer',{collapsed}]:'settings-nav'" :aria-label="t('全局导航')">
    <RouterLink v-if="!workspace" to="/note" class="settings-brand" :title="t('返回笔记')"><img :src="brandMark" :alt="t('珠玑笔记')"><span>{{t('珠玑笔记')}}</span></RouterLink>
    <button v-if="workspace&&desktop" class="footer-sync" :class="{'has-pending':pending}" :disabled="syncing" :aria-busy="syncing" :title="pending?t('有未同步的更改'):t('立即同步')" :aria-label="t('立即同步')" @click="runImmediateSync">↻<span v-if="pending" class="sync-pending-dot" aria-hidden="true"></span></button>
    <RouterLink v-if="!workspace&&backOnly" to="/note" class="settings-back">{{t('返回')}}</RouterLink>
    <div v-else class="user-menu" @focusout="focusOut" @keydown.escape="close">
      <button class="avatar-button" :aria-expanded="menuOpen" aria-haspopup="menu" :aria-label="t('用户菜单')" @click="menuOpen=!menuOpen"><img v-if="avatarSrc" :src="avatarSrc" alt=""><span v-else>{{initials}}</span></button>
      <div v-if="menuOpen" class="user-popover" role="menu">
        <button v-if="desktop" class="sync-trigger" role="menuitem" aria-haspopup="menu" :aria-expanded="syncOpen" @click="syncOpen=!syncOpen;languageOpen=false">{{t('同步')}} <span aria-hidden="true">{{syncOpen?'▾':'▸'}}</span></button>
        <div v-if="desktop&&syncOpen" class="sync-options" role="menu" :aria-label="t('同步')">
          <button role="menuitem" :disabled="syncing" @click="runSync('full')">{{t('完全同步')}}</button>
          <button role="menuitem" :disabled="syncing" @click="close();resetPrompt=true">{{t('重新同步')}}</button>
        </div>
        <RouterLink v-if="admin&&!desktop" to="/admin" role="menuitem" @click="close">{{t('管理')}}</RouterLink>
        <RouterLink to="/member" role="menuitem" @click="close">{{t('账号')}}</RouterLink>
        <button class="language-trigger" role="menuitem" aria-haspopup="menu" :aria-expanded="languageOpen" @click="languageOpen=!languageOpen;syncOpen=false">{{t('语言')}} <span aria-hidden="true">{{languageOpen?'▾':'▸'}}</span></button>
        <div v-if="languageOpen" class="language-options" role="menu" :aria-label="t('语言')">
          <button v-for="item in languages" :key="item.code" role="menuitemradio" :aria-checked="language===item.code" @click="chooseLanguage(item.code)"><span>{{item.name}}</span><span v-if="language===item.code" aria-hidden="true">✓</span></button>
        </div>
        <button v-if="desktop" class="about-trigger" role="menuitem" @click="showAbout">{{t('关于')}}</button>
        <a href="/api2/logout" role="menuitem" @click="logout">{{t('退出')}}</a>
      </div>
    </div>
  </nav>
  <Teleport to="body">
    <section v-if="logoutPrompt" class="choice-dialog-backdrop" role="presentation" @click.self="chooseLogout(false)">
      <div class="choice-dialog" role="dialog" aria-modal="true" :aria-label="t('退出登录')">
        <h2>{{t('退出登录')}}</h2>
        <p>{{t('退出前同步失败，是否在不同步的情况下退出？')}}</p>
        <div class="choice-dialog-actions">
          <button class="primary" @click="chooseLogout(true)">{{t('确定')}}</button>
          <button @click="chooseLogout(false)">{{t('取消')}}</button>
        </div>
      </div>
    </section>
    <section v-if="resetPrompt" class="choice-dialog-backdrop" role="presentation" @click.self="resetPrompt=false">
      <div class="choice-dialog" role="dialog" aria-modal="true" :aria-label="t('重新同步')">
        <h2>{{t('重新同步')}}</h2>
        <p>{{t('将删除当前账户的全部本地数据，从服务端重新下载。未上传的本地修改会永久丢失。是否继续？')}}</p>
        <div class="choice-dialog-actions">
          <button @click="resetPrompt=false">{{t('取消')}}</button>
          <button class="primary" @click="resetPrompt=false;runSync('reset')">{{t('确认重新同步')}}</button>
        </div>
      </div>
    </section>
    <section v-if="aboutOpen" class="choice-dialog-backdrop" role="presentation" @click.self="aboutOpen=false">
      <div class="choice-dialog about-dialog" role="dialog" aria-modal="true" :aria-label="t('关于珠玑笔记')">
        <img :src="brandMark" :alt="t('珠玑笔记')">
        <h2>{{t('珠玑笔记')}}</h2>
        <p>{{t('日积字句，终得珠玑。')}}</p>
        <dl>
          <div><dt>{{t('版本')}}</dt><dd>{{aboutInfo.Version||'—'}}</dd></div>
          <div><dt>{{t('平台')}}</dt><dd>{{aboutInfo.Platform}} / {{aboutInfo.Arch}}</dd></div>
          <div><dt>{{t('运行环境')}}</dt><dd>{{aboutInfo.Runtime}}</dd></div>
        </dl>
        <div class="choice-dialog-actions"><button class="primary" @click="aboutOpen=false">{{t('关闭')}}</button></div>
      </div>
    </section>
  </Teleport>
</template>
