<script setup lang="ts">
import {errorMessage} from '../errors'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { request, upload, resolveAvatarUrl } from '../api'
import Navigation from '../components/Navigation.vue'
import EmailVerification from '../components/EmailVerification.vue'
import FeedbackMessage from '../components/FeedbackMessage.vue'
import {confirmAction} from '../dialogs'
import {t} from '../i18n'

type GroupUser = { UserId: string; Username: string; Email: string }
type Group = { GroupId: string; UserId: string; Title: string; Users?: GroupUser[] }

const user = ref<any>({})
const admin = ref(false)
type Feedback = { text: string; error: boolean }
const feedback = ref<Record<string, Feedback | undefined>>({})
const username = ref('')
const serverAddress = ref('')
const oldPwd = ref('')
const pwd = ref(''), confirmPwd = ref(''), passwordError = ref('')
const usernameSaving = ref(false), passwordSaving = ref(false), avatarUploading = ref(false)
const avatarInput = ref<HTMLInputElement>()
const avatarSrc = computed(() => resolveAvatarUrl(user.value.Logo || ''))
const groups = ref<Group[]>([])
const newGroupTitle = ref('')
const addUserIdentities = ref<Record<string, string>>({})
const editingGroupId = ref('')
const editingGroupTitle = ref('')

function setMessage(scope: string, text: string, error = false) { feedback.value[scope] = { text, error } }
function clearMessage(scope: string) { delete feedback.value[scope] }
function showError(scope: string, error: unknown) { setMessage(scope, errorMessage(error), true) }
function owned(group: Group) { return group.UserId === user.value.UserId }
async function loadGroups() {
  clearMessage('groupsLoad')
  try {
    const result = await request<Group[]>('/web/groups')
    if (!Array.isArray(result)) throw new Error(t('分组接口返回的数据格式不正确'))
    groups.value = result
  } catch (error) { showError('groupsLoad', error) }
}

async function refreshEmail() {
  if(!user.value.UserId)return
  clearMessage('email')
  try {
    const profile=await request('/user/info')
    if(profile.UserId===user.value.UserId){
      user.value.Email=profile.Email;user.value.Verified=profile.Verified
      if('PendingEmail' in profile)user.value.PendingEmail=profile.PendingEmail
    }
  }catch(error){showError('email',error)}
}
onMounted(()=>window.addEventListener('focus',refreshEmail))
onUnmounted(()=>window.removeEventListener('focus',refreshEmail))
onMounted(async () => {
  try {
    const bootstrap = await request('/web/bootstrap')
    if (!bootstrap.User) { location.href = '/login'; return }
    user.value = bootstrap.User; username.value = bootstrap.User.Username; admin.value = bootstrap.IsAdmin
    serverAddress.value = bootstrap.Desktop ? bootstrap.Host || '' : location.origin
    await loadGroups()
  } catch (error) { showError('page', error) }
})

async function updateUsername() {
  if (usernameSaving.value) return
  usernameSaving.value = true; clearMessage('username')
  try {
    const value = username.value.trim()
    await request('/user/updateUsername', { username: value })
    user.value.Username = value; username.value = value; setMessage('username', t('已更新'))
  } catch (error) { showError('username', error) }
  finally { usernameSaving.value = false }
}
async function updatePassword() {
  if (passwordSaving.value) return
  passwordError.value = ''
  if (pwd.value !== confirmPwd.value) { passwordError.value = t('两次输入的密码不一致'); return }
  passwordSaving.value = true
  try {
    await request('/user/updatePwd', { oldPwd: oldPwd.value, pwd: pwd.value })
    oldPwd.value = ''; pwd.value = ''; confirmPwd.value = ''
    await request('/web/logout', {})
    location.href = '/login'
  } catch (error) { passwordError.value = errorMessage(error) }
  finally { passwordSaving.value = false }
}
async function avatar(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = '' // Allow selecting the same file again after a failed upload.
  if (!file || avatarUploading.value) return
  avatarUploading.value = true; clearMessage('avatar')
  try { const result = await upload('/file/uploadAvatar', file); user.value.Logo = resolveAvatarUrl(result.Id); setMessage('avatar', t('头像已更新')) } catch (error) { showError('avatar', error) }
  finally { avatarUploading.value = false }
}
async function addGroup() {
  clearMessage('groups')
  const title = newGroupTitle.value.trim()
  if (!title) { setMessage('groups', t('请输入分组名称'), true); return }
  try { await request('/member/group/addGroup', { title }); newGroupTitle.value = ''; setMessage('groups', t('分组已创建')); await loadGroups() } catch (error) { showError('groups', error) }
}
function startRename(group: Group) { clearMessage(`group:${group.GroupId}`); editingGroupId.value = group.GroupId; editingGroupTitle.value = group.Title }
async function renameGroup(group: Group) {
  const scope = `group:${group.GroupId}`
  clearMessage(scope)
  const title = editingGroupTitle.value.trim()
  if (!title) { setMessage(scope, t('请输入分组名称'), true); return }
  try { await request('/member/group/updateGroupTitle', { groupId: group.GroupId, title }); editingGroupId.value = ''; setMessage(scope, t('分组名称已更新')); await loadGroups() } catch (error) { showError(scope, error) }
}
async function deleteGroup(group: Group) {
  if (!await confirmAction(t('确定删除分组“{name}”吗？', {name:group.Title}))) return
  const scope = `group:${group.GroupId}`
  clearMessage(scope)
  try { await request('/member/group/deleteGroup', { groupId: group.GroupId }); clearMessage(scope); setMessage('groups', t('分组已删除')); await loadGroups() } catch (error) { showError(scope, error) }
}
async function addUser(group: Group) {
  const scope = `group:${group.GroupId}`
  clearMessage(scope)
  const memberIdentity = (addUserIdentities.value[group.GroupId] || '').trim()
  if (!memberIdentity) { setMessage(scope, t('请输入用户名或邮箱'), true); return }
  try { await request('/member/group/addUser', { groupId: group.GroupId, email: memberIdentity }); addUserIdentities.value[group.GroupId] = ''; setMessage(scope, t('用户已加入分组')); await loadGroups() } catch (error) { showError(scope, error) }
}
async function deleteUser(group: Group, member: GroupUser) {
  if (!await confirmAction(t('确定将 {name} 移出此分组吗？', {name:member.Email || member.Username}))) return
  const scope = `group:${group.GroupId}`
  clearMessage(scope)
  try { await request('/member/group/deleteUser', { groupId: group.GroupId, userId: member.UserId }); setMessage(scope, t('用户已移出分组')); await loadGroups() } catch (error) { showError(scope, error) }
}
</script>

<template>
  <div class="shell settings-shell">
    <Navigation :admin="admin" :user="user" back-only />
    <main class="settings">
      <h1>{{t('账号管理')}}</h1><FeedbackMessage :feedback="feedback.page" />
      <label v-if="serverAddress" class="server-address">{{t('服务器地址')}}<input :value="serverAddress" :title="serverAddress" readonly></label>
      <section class="card profile-card">
        <h2>{{t('个人资料')}}</h2>
        <button type="button" class="profile-avatar" :disabled="avatarUploading" :aria-busy="avatarUploading" :aria-label="t('更换头像')" :title="t('更换头像')" @click="avatarInput?.click()">
          <img v-if="avatarSrc" :src="avatarSrc" :alt="t('当前头像')">
          <span v-else>{{ (user.Username || user.Email || '?').slice(0, 1).toUpperCase() }}</span>
        </button>
        <input ref="avatarInput" type="file" hidden accept="image/*" :disabled="avatarUploading" :aria-label="t('更换头像')" @change="avatar">
        <FeedbackMessage :feedback="feedback.avatar" />
        <form @submit.prevent="updateUsername">
          <label for="account-username">{{t('用户名')}}</label>
          <div class="username-row">
            <input id="account-username" v-model="username" required autocomplete="username" :disabled="usernameSaving">
            <button class="username-save" :disabled="usernameSaving" :aria-label="t('更新用户名')" :title="t('更新用户名')">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h12l4 4v12a2 2 0 0 1-2 2Z"/><path d="M17 21v-8H7v8M7 3v5h9"/></svg>
            </button>
          </div>
        </form>
        <FeedbackMessage :feedback="feedback.username" />
        <EmailVerification v-if="user.UserId" :user-id="user.UserId" :current-email="user.Email" :verified="!!user.Verified" :pending-email="user.PendingEmail" @pending="user.PendingEmail=$event"/>
        <FeedbackMessage :feedback="feedback.email" />
      </section>
      <section class="card password-card">
        <h2>{{t('修改密码')}}</h2>
        <form @submit.prevent="updatePassword">
          <label>{{t('原密码')}}<input v-model="oldPwd" type="password" autocomplete="current-password" required :disabled="passwordSaving"></label>
          <label>{{t('新密码')}}<input v-model="pwd" type="password" autocomplete="new-password" required :disabled="passwordSaving" @input="passwordError=''"></label>
          <label>{{t('确认新密码')}}<input v-model="confirmPwd" type="password" autocomplete="new-password" required :disabled="passwordSaving" :aria-invalid="!!passwordError" :aria-describedby="passwordError?'password-error':undefined" @input="passwordError=''"></label>
          <button class="primary" :disabled="passwordSaving">{{t('修改并重新登录')}}</button>
          <FeedbackMessage id="password-error" :feedback="passwordError ? {text: passwordError, error: true} : undefined" />
        </form>
      </section>
      <section class="card groups">
        <h2>{{t('用户分组')}}</h2><p class="muted">{{t('将已有账号加入分组后，可在分享笔记或笔记本时选择整个分组。')}}</p>
        <form class="inline-form" @submit.prevent="addGroup"><input v-model="newGroupTitle" maxlength="100" :placeholder="t('新分组名称')" :aria-label="t('新分组名称')"><button class="primary">{{t('新建分组')}}</button></form>
        <FeedbackMessage :feedback="feedback.groups" />
        <FeedbackMessage :feedback="feedback.groupsLoad" />
        <p v-if="!groups.length" class="muted">{{t('尚无分组。')}}</p>
        <article v-for="group in groups" :key="group.GroupId" class="group">
          <header>
            <form v-if="editingGroupId === group.GroupId" class="inline-form" @submit.prevent="renameGroup(group)"><input v-model="editingGroupTitle" maxlength="100" :aria-label="t('分组名称')"><button>{{t('保存')}}</button><button type="button" @click="editingGroupId = ''">{{t('取消')}}</button></form>
            <template v-else><h3>{{ group.Title }}</h3><span class="muted">{{ group.Users?.length || 0 }} {{t('位成员')}}</span><span v-if="!owned(group)" class="muted">{{t('（他人分组）')}}</span><div v-if="owned(group)" class="actions"><button @click="startRename(group)">{{t('改名')}}</button><button class="danger" @click="deleteGroup(group)">{{t('删除')}}</button></div></template>
          </header>
          <ul v-if="group.Users?.length" class="members"><li v-for="member in group.Users" :key="member.UserId"><span>{{ member.Username || member.Email }}</span><span class="muted">{{ member.Email }}</span><button v-if="owned(group)" class="danger" @click="deleteUser(group, member)">{{t('移除')}}</button></li></ul>
          <p v-else class="muted">{{t('暂无成员。')}}</p>
          <form v-if="owned(group)" class="inline-form" @submit.prevent="addUser(group)"><input v-model="addUserIdentities[group.GroupId]" type="text" autocomplete="off" autocapitalize="none" :spellcheck="false" :placeholder="t('已有用户的用户名或邮箱')" :aria-label="`${group.Title} · ${t('已有用户的用户名或邮箱')}`"><button>{{t('添加用户')}}</button></form>
          <FeedbackMessage :feedback="feedback[`group:${group.GroupId}`]" />
        </article>
      </section>
    </main>
  </div>
</template>

<style scoped>
.profile-avatar { display: grid; place-items: center; width: 72px; height: 72px; padding: 0; border-radius: 50%; overflow: hidden; background: #e9eee8; color: #193b35; font-size: 28px }
.profile-avatar img { width: 100%; height: 100%; object-fit: cover }
.username-row { display: flex; align-items: stretch; gap: .5rem }
.username-row input { flex: 1; min-width: 0 }
.username-save { display: grid; place-items: center; flex: 0 0 44px; padding: 8px }
label[for="account-username"] { margin-bottom: 8px }
.inline-form, .actions { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }
.group { padding: 1rem 0; border-top: 1px solid #e5e7eb }.group header { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }.group h3 { margin: 0; margin-right: .25rem }.actions { margin-left: auto }.members { list-style: none; padding: 0; margin: .75rem 0 }.members li { display: flex; align-items: center; gap: .75rem; padding: .35rem 0 }.members .muted { flex: 1 }.danger { color: #b42318 }
</style>
