<script setup lang="ts">
import {computed,ref} from 'vue'
import {request,resolveAvatarUrl} from '../api'
import {language,languages,setLanguage,t} from '../i18n'
import brandMark from '../assets/gemsnote_s.png'

const props=defineProps<{admin?:boolean,user?:{Username?:string;Email?:string;Logo?:string},workspace?:boolean,collapsed?:boolean,pending?:boolean,backOnly?:boolean,desktop?:boolean}>()
const emit=defineEmits<{synced:[]}>()
const menuOpen=ref(false)
const languageOpen=ref(false)
const initials=computed(()=>(props.user?.Username||props.user?.Email||t('用户')).trim().slice(0,1).toUpperCase())
const avatarSrc=computed(()=>resolveAvatarUrl(props.user?.Logo||''))
const desktop=computed(()=>!!props.desktop)
const syncing=ref(false)

function close(){menuOpen.value=false;languageOpen.value=false}
function chooseLanguage(value:typeof language.value){setLanguage(value);close()}
function focusOut(event:FocusEvent){if(!(event.currentTarget as HTMLElement).contains(event.relatedTarget as Node|null))close()}
async function runSync(full=false){
  if(syncing.value)return
  syncing.value=true
  if(full){
    close()
    window.dispatchEvent(new CustomEvent('sync-progress',{detail:{Stage:'start',Current:0,Total:100}}))
  }
  try{
    if(!desktop.value&&!full){
      const syncNow=(window as any).__gemsnoteSyncNow
      if(typeof syncNow==='function'&&!await syncNow())return
    }else await request(full?'/web/fullSync':'/web/sync',{})
    emit('synced')
  }catch(e){window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:e instanceof Error?e.message:String(e),Full:full}}))}
  finally{syncing.value=false}
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
    let hasPending=!!props.pending
    if(!hasPending){
      try{hasPending=!!(await request<any>('/web/bootstrap')).PendingChanges}catch{/* logout reports the real session error below */}
    }
    const force=hasPending&&confirm(t('仍有未同步的更改。是否在不同步的情况下退出？'))
    await request('/web/logout',force?{force:true}:{})
    location.href='/login'
  }catch(e){
    const message=e instanceof Error?e.message:String(e)
    window.dispatchEvent(new CustomEvent('sync-result',{detail:{Ok:false,Msg:message==='syncFailed'?t('退出前同步失败，请确认服务端可用后重试'):message,Full:false}}))
  }
}
</script>

<template>
  <nav :class="workspace?['space-footer',{collapsed}]:'settings-nav'" :aria-label="t('全局导航')">
    <RouterLink v-if="!workspace" to="/note" class="settings-brand" :title="t('返回笔记')"><img :src="brandMark" :alt="t('珠玑笔记')"><span>{{t('珠玑笔记')}}</span></RouterLink>
    <button v-if="workspace&&desktop" class="footer-sync" :class="{'has-pending':pending}" :disabled="syncing" :aria-busy="syncing" :title="pending?t('有未同步的更改'):t('立即同步')" :aria-label="t('立即同步')" @click="runSync()">↻<span v-if="pending" class="sync-pending-dot" aria-hidden="true"></span></button>
    <RouterLink v-if="!workspace&&backOnly" to="/note" class="settings-back">{{t('返回')}}</RouterLink>
    <div v-else class="user-menu" @focusout="focusOut" @keydown.escape="close">
      <button class="avatar-button" :aria-expanded="menuOpen" aria-haspopup="menu" :aria-label="t('用户菜单')" @click="menuOpen=!menuOpen"><img v-if="avatarSrc" :src="avatarSrc" alt=""><span v-else>{{initials}}</span></button>
      <div v-if="menuOpen" class="user-popover" role="menu">
        <a v-if="desktop" href="#" role="menuitem" :aria-busy="syncing" @click.prevent="runSync(true)">{{t('完全同步')}}</a>
        <RouterLink v-if="admin&&!desktop" to="/admin" role="menuitem" @click="close">{{t('管理')}}</RouterLink>
        <RouterLink to="/member" role="menuitem" @click="close">{{t('账号')}}</RouterLink>
        <button class="language-trigger" role="menuitem" aria-haspopup="menu" :aria-expanded="languageOpen" @click="languageOpen=!languageOpen">{{t('语言')}} <span aria-hidden="true">{{languageOpen?'▾':'▸'}}</span></button>
        <div v-if="languageOpen" class="language-options" role="menu" :aria-label="t('语言')">
          <button v-for="item in languages" :key="item.code" role="menuitemradio" :aria-checked="language===item.code" @click="chooseLanguage(item.code)"><span>{{item.name}}</span><span v-if="language===item.code" aria-hidden="true">✓</span></button>
        </div>
        <a href="/api2/logout" role="menuitem" @click="logout">{{t('退出')}}</a>
      </div>
    </div>
  </nav>
</template>
