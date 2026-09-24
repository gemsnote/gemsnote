<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { request, upload, resolveAvatarUrl } from '../api'
import Navigation from '../components/Navigation.vue'
import FilePicker from '../components/FilePicker.vue'
import {confirmAction} from '../dialogs'
import {t} from '../i18n'

type GroupUser = { UserId: string; Username: string; Email: string }
type Group = { GroupId: string; UserId: string; Title: string; Users?: GroupUser[] }

const user = ref<any>({})
const admin = ref(false)
const message = ref('')
const username = ref('')
const email = ref('')
const serverAddress = ref('')
const emailPwd = ref('')
const oldPwd = ref('')
const pwd = ref('')
const groups = ref<Group[]>([])
const newGroupTitle = ref('')
const addUserEmails = ref<Record<string, string>>({})
const editingGroupId = ref('')
const editingGroupTitle = ref('')

function showError(error: unknown) { message.value = error instanceof Error ? error.message : String(error) }
function owned(group: Group) { return group.UserId === user.value.UserId }
async function loadGroups() {
  try {
    const result = await request<Group[]>('/web/groups')
    if (!Array.isArray(result)) throw new Error(t('分组接口返回的数据格式不正确'))
    groups.value = result
  } catch (error) { groups.value = []; showError(error) }
}

onMounted(async () => {
  try {
    const bootstrap = await request('/web/bootstrap')
    if (!bootstrap.User) { location.href = '/login'; return }
    user.value = bootstrap.User; username.value = bootstrap.User.Username; email.value = bootstrap.User.Email; admin.value = bootstrap.IsAdmin
    serverAddress.value = bootstrap.Desktop ? bootstrap.Host || '' : location.origin
    await loadGroups()
  } catch (error) { showError(error) }
})

async function update(path: string, data: any) {
  try {
    await request(path, data)
    message.value = t('已更新')
    if (path.endsWith('updatePwd')) {
      await request('/web/logout', {})
      location.href = '/login'
    }
  } catch (error) { showError(error) }
}
async function avatar(file: File) {
  try { const result = await upload('/file/uploadAvatar', file); user.value.Logo = resolveAvatarUrl(result.Id); message.value = t('头像已更新') } catch (error) { showError(error) }
}
async function addGroup() {
  const title = newGroupTitle.value.trim()
  if (!title) { message.value = t('请输入分组名称'); return }
  try { await request('/member/group/addGroup', { title }); newGroupTitle.value = ''; message.value = t('分组已创建'); await loadGroups() } catch (error) { showError(error) }
}
function startRename(group: Group) { editingGroupId.value = group.GroupId; editingGroupTitle.value = group.Title }
async function renameGroup(group: Group) {
  const title = editingGroupTitle.value.trim()
  if (!title) { message.value = t('请输入分组名称'); return }
  try { await request('/member/group/updateGroupTitle', { groupId: group.GroupId, title }); editingGroupId.value = ''; message.value = t('分组名称已更新'); await loadGroups() } catch (error) { showError(error) }
}
async function deleteGroup(group: Group) {
  if (!await confirmAction(t('确定删除分组“{name}”吗？', {name:group.Title}))) return
  try { await request('/member/group/deleteGroup', { groupId: group.GroupId }); message.value = t('分组已删除'); await loadGroups() } catch (error) { showError(error) }
}
async function addUser(group: Group) {
  const memberEmail = (addUserEmails.value[group.GroupId] || '').trim()
  if (!memberEmail) { message.value = t('请输入用户邮箱'); return }
  try { await request('/member/group/addUser', { groupId: group.GroupId, email: memberEmail }); addUserEmails.value[group.GroupId] = ''; message.value = t('用户已加入分组'); await loadGroups() } catch (error) { showError(error) }
}
async function deleteUser(group: Group, member: GroupUser) {
  if (!await confirmAction(t('确定将 {name} 移出此分组吗？', {name:member.Email || member.Username}))) return
  try { await request('/member/group/deleteUser', { groupId: group.GroupId, userId: member.UserId }); message.value = t('用户已移出分组'); await loadGroups() } catch (error) { showError(error) }
}
</script>

<template>
  <div class="shell settings-shell">
    <Navigation :admin="admin" :user="user" back-only />
    <main class="settings">
      <h1>{{t('账号管理')}}</h1><p class="muted">{{ user.Email }}</p><p v-if="message" role="status" class="message">{{ message }}</p>
      <label v-if="serverAddress" class="server-address">{{t('服务器地址')}}<input :value="serverAddress" :title="serverAddress" readonly></label>
      <section class="card"><h2>{{t('个人资料')}}</h2><img v-if="user.Logo" :src="user.Logo" class="avatar" :alt="t('当前头像')"><FilePicker :label="t('更换头像')" accept="image/*" @selected="avatar"/><form @submit.prevent="update('/user/updateUsername', { username })"><label>{{t('用户名')}}<input v-model="username" required></label><button>{{t('更新用户名')}}</button></form><form @submit.prevent="update('/web/emailChange', { email, pwd: emailPwd })"><label>{{t('新邮箱')}}<input v-model="email" type="email" required></label><label>{{t('当前密码')}}<input v-model="emailPwd" type="password" required></label><button>{{t('发送邮箱验证邮件')}}</button></form><button @click="update('/user/reSendActiveEmail', {})">{{t('重新发送当前邮箱验证邮件')}}</button></section>
      <section class="card"><h2>{{t('修改密码')}}</h2><form @submit.prevent="update('/user/updatePwd', { oldPwd, pwd })"><label>{{t('原密码')}}<input v-model="oldPwd" type="password" autocomplete="current-password" required></label><label>{{t('新密码')}}<input v-model="pwd" type="password" autocomplete="new-password" required></label><button class="primary">{{t('修改并重新登录')}}</button></form></section>
      <section class="card groups">
        <h2>{{t('用户分组')}}</h2><p class="muted">{{t('将已有账号加入分组后，可在分享笔记或笔记本时选择整个分组。')}}</p>
        <form class="inline-form" @submit.prevent="addGroup"><input v-model="newGroupTitle" maxlength="100" :placeholder="t('新分组名称')" :aria-label="t('新分组名称')"><button class="primary">{{t('新建分组')}}</button></form>
        <p v-if="!groups.length" class="muted">{{t('尚无分组。')}}</p>
        <article v-for="group in groups" :key="group.GroupId" class="group">
          <header>
            <form v-if="editingGroupId === group.GroupId" class="inline-form" @submit.prevent="renameGroup(group)"><input v-model="editingGroupTitle" maxlength="100" :aria-label="t('分组名称')"><button>{{t('保存')}}</button><button type="button" @click="editingGroupId = ''">{{t('取消')}}</button></form>
            <template v-else><h3>{{ group.Title }}</h3><span class="muted">{{ group.Users?.length || 0 }} {{t('位成员')}}</span><span v-if="!owned(group)" class="muted">{{t('（他人分组）')}}</span><div v-if="owned(group)" class="actions"><button @click="startRename(group)">{{t('改名')}}</button><button class="danger" @click="deleteGroup(group)">{{t('删除')}}</button></div></template>
          </header>
          <ul v-if="group.Users?.length" class="members"><li v-for="member in group.Users" :key="member.UserId"><span>{{ member.Username || member.Email }}</span><span class="muted">{{ member.Email }}</span><button v-if="owned(group)" class="danger" @click="deleteUser(group, member)">{{t('移除')}}</button></li></ul>
          <p v-else class="muted">{{t('暂无成员。')}}</p>
          <form v-if="owned(group)" class="inline-form" @submit.prevent="addUser(group)"><input v-model="addUserEmails[group.GroupId]" type="email" :placeholder="t('已有用户的邮箱')" :aria-label="`${group.Title} · ${t('已有用户的邮箱')}`"><button>{{t('添加用户')}}</button></form>
        </article>
      </section>
    </main>
  </div>
</template>

<style scoped>
.avatar { width: 72px; height: 72px; border-radius: 50%; object-fit: cover }
.server-address { max-width: 760px }
.message { color: var(--accent, #1769aa) }
.inline-form, .actions { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }
.groups { max-width: 760px }.group { padding: 1rem 0; border-top: 1px solid #e5e7eb }.group header { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap }.group h3 { margin: 0; margin-right: .25rem }.actions { margin-left: auto }.members { list-style: none; padding: 0; margin: .75rem 0 }.members li { display: flex; align-items: center; gap: .75rem; padding: .35rem 0 }.members .muted { flex: 1 }.danger { color: #b42318 }
</style>
