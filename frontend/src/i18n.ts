import { ref, watch } from 'vue'

export type Language = 'zh-CN' | 'en-US'
export const languages: { code: Language; name: string }[] = [
  { code: 'zh-CN', name: '简体中文' },
  { code: 'en-US', name: 'English' },
]

const translations: Record<string, string> = {
  '珠玑笔记': 'Gemsnote', '页面不存在': 'Page not found', '返回笔记': 'Back to notes',
  '语言': 'Language', '返回': 'Back', '用户': 'User', '用户菜单': 'User menu', '全局导航': 'Main navigation',
  '立即同步': 'Sync now', '有未同步的更改': 'Unsynced changes', '完全同步': 'Full sync',
  '同步': 'Sync', '重新同步': 'Reset and sync', '确认重新同步': 'Reset and sync',
  '选择登录后的同步方式': 'Choose sync after sign-in', '暂不同步': 'Not now',
  '检测到此账号已有本地缓存。重新同步会删除该账号的全部本地数据，并从服务端重新下载；未上传的修改会永久丢失。': 'A local cache exists for this account. Reset and sync deletes all local data for this account and downloads it again from the server. Unsynced changes will be permanently lost.',
  '本地笔记保存失败，已取消重新同步': 'Could not save the local note. Reset and sync was canceled.',
  '将删除当前账户的全部本地数据，从服务端重新下载。未上传的本地修改会永久丢失。是否继续？': 'All local data for this account will be deleted and downloaded again from the server. Unsynced local changes will be permanently lost. Continue?',
  '退出登录': 'Sign out', '退出前同步失败，是否在不同步的情况下退出？': 'Sync failed. Sign out without syncing?',
  '确定': 'Confirm',
  '管理': 'Administration', '账号': 'Account', '退出': 'Sign out',
  '关于': 'About', '关于珠玑笔记': 'About Gemsnote', '版本': 'Version', '平台': 'Platform',
  '运行环境': 'Runtime',
  '正在加载…': 'Loading…', '无法读取应用信息': 'Could not load application information', '重试': 'Retry',
  '选择文件': 'Choose file', '未选择文件': 'No file selected', '选择图片': 'Choose image', '选择附件': 'Choose attachment',
  '请填写此字段。': 'Please fill out this field.', '请输入有效的邮箱地址。': 'Please enter a valid email address.',
  '请输入有效的网址。': 'Please enter a valid URL.', '请输入有效值。': 'Please enter a valid value.',
  '请输入不小于 {min} 的值。': 'Please enter a value greater than or equal to {min}.',
  '请输入不大于 {max} 的值。': 'Please enter a value less than or equal to {max}.',
  '请至少输入 {min} 个字符。': 'Please enter at least {min} characters.',
  '请最多输入 {max} 个字符。': 'Please enter no more than {max} characters.',
  '退出前同步失败，请确认服务端可用后重试': 'Sync failed before sign-out. Check the server and try again.',
  '登录': 'Sign in', '注册': 'Register', '找回密码': 'Forgot password', '欢迎回来': 'Welcome back',
  '创建账号': 'Create an account', '服务器': 'Server', '邮箱或用户名': 'Email or username',
  '服务器地址': 'Server address',
  '密码': 'Password', '验证码': 'Verification code', '处理中…': 'Working…', '正在登录…': 'Signing in…', '正在同步…': 'Syncing…', '继续': 'Continue',
  'Gemsnote · 记录、整理、沉淀': 'Gemsnote · Capture, organize, reflect',
  '当前连接的是旧版 Leanote 服务端，请先迁移到 Gemsnote 服务端。': 'This is an old Leanote server. Migrate it to Gemsnote first.',
  '当前客户端版本过低，请先升级客户端。': 'This client is too old. Upgrade it first.',
  '当前服务端版本过低，请先升级服务端。': 'This server is too old. Upgrade it first.',
  '密码已修改，请重新登录。': 'Password changed. Please sign in again.',
  '请检查邮箱中的重置密码链接。': 'Check your email for the password reset link.',
  '账号管理': 'Account settings', '个人资料': 'Profile', '当前头像': 'Current avatar',
  '更换头像': 'Change avatar', '用户名': 'Username', '更新用户名': 'Update username',
  '新邮箱': 'New email', '当前密码': 'Current password', '发送邮箱验证邮件': 'Send verification email',
  '重新发送当前邮箱验证邮件': 'Resend verification email', '修改密码': 'Change password',
  '原密码': 'Current password', '新密码': 'New password', '修改并重新登录': 'Change password and sign in again',
  '用户分组': 'User groups', '将已有账号加入分组后，可在分享笔记或笔记本时选择整个分组。': 'Add existing users to a group to share notes or notebooks with them.',
  '新分组名称': 'New group name', '新建分组': 'Create group', '尚无分组。': 'No groups yet.',
  '分组名称': 'Group name', '保存': 'Save', '取消': 'Cancel', '位成员': 'members',
  '（他人分组）': '(another user’s group)', '改名': 'Rename', '删除': 'Delete',
  '移除': 'Remove', '暂无成员。': 'No members yet.', '已有用户的邮箱': 'Existing user email',
  '添加用户': 'Add user', '已更新': 'Updated', '头像已更新': 'Avatar updated',
  '请输入分组名称': 'Enter a group name', '分组已创建': 'Group created',
  '分组名称已更新': 'Group renamed', '分组已删除': 'Group deleted',
  '请输入用户邮箱': 'Enter a user email', '用户已加入分组': 'User added to group',
  '用户已移出分组': 'User removed from group', '分组接口返回的数据格式不正确': 'Invalid group response',
  '确定删除分组“{name}”吗？': 'Delete group “{name}”?',
  '确定将 {name} 移出此分组吗？': 'Remove {name} from this group?',
  '系统管理': 'System administration', '操作成功': 'Operation completed',
  '用户管理': 'User management', '搜索用户': 'Search users', '搜索': 'Search',
  '邮箱': 'Email', '操作': 'Actions', '重置密码': 'Reset password',
  '上一页': 'Previous', '下一页': 'Next', '确认重置': 'Confirm reset',
  '添加用户（管理）': 'Add user', '初始密码': 'Initial password', '创建': 'Create',
  '站点与邮件': 'Site and email', '站点 URL': 'Site URL', '开放注册': 'Open registration',
  '是': 'Yes', '否': 'No', '邮件服务器': 'Mail server', '邮件端口': 'Mail port',
  '邮件账号': 'Mail account', '邮件密码': 'Mail password', '留空保持原值': 'Leave blank to keep current value',
  '邮件 SSL': 'Mail SSL', '启用': 'Enabled', '禁用': 'Disabled',
  '上传限制（MB）': 'Upload limits (MB)', '图片': 'Image', '头像': 'Avatar',
  '附件': 'Attachments', '其他非博客设置': 'Other settings', 'PDF 导出程序路径': 'PDF export program path',
  '保存设置': 'Save settings', '正在验证…': 'Verifying…', '邮箱验证': 'Email verification',
  '邮箱 {email} 验证成功，请重新登录。': 'Email {email} verified. Please sign in again.',
  '返回登录': 'Back to sign in', '编辑器加载失败，请刷新重试': 'Editor failed to load. Refresh and try again.',
  '富文本正文': 'Rich-text content', 'Markdown 正文': 'Markdown content', 'HTML 正文': 'HTML content',
  '共享笔记': 'Shared notes', '共享笔记本': 'Shared notebook', '未知笔记本': 'Unknown notebook',
  '回收站': 'Trash', '已加星': 'Starred', '所有笔记': 'All notes',
  '笔记接口返回的数据格式不正确': 'Invalid notes response',
  '离线缓存': 'Offline cache', '离线缓存 {date}': 'Cached offline {date}',
  '缓存待更新': 'Cache needs updating', '未下载，请先同步': 'Not downloaded; sync first',
  '已撤销访问': 'Access revoked', '已缓存': 'Cached', '未下载': 'Not downloaded',
  '下载失败': 'Download failed', '已加入离线下载队列，稍后自动下载': 'Queued for offline download.',
  '此笔记已在其他端修改。请先导出本地内容，再重新加载笔记。': 'This note changed on another device. Export your local content before reloading.',
  '请先选择一个可编辑的共享笔记本': 'Select an editable shared notebook first',
  '未命名笔记': 'Untitled note', '子笔记本名称': 'Subnotebook name', '笔记本名称': 'Notebook name',
  '删除此笔记本？请先移动其中的笔记。': 'Delete this notebook? Move its notes first.',
  '永久删除这篇笔记？此操作不可撤销。': 'Permanently delete this note? This cannot be undone.',
  '将笔记放入回收站？': 'Move this note to Trash?',
  '从共享列表移除这篇笔记？': 'Remove this note from your shared list?',
  '图片上传失败': 'Image upload failed', '删除此附件？': 'Delete this attachment?',
  '同步失败': 'Sync failed', '完全同步成功': 'Full sync completed',
  '正在完全同步': 'Full synchronization in progress', '正在准备同步': 'Preparing synchronization',
  '正在连接服务器': 'Connecting to server', '正在清理本地缓存': 'Clearing local cache',
  '已完成 {current} / {total} 项': 'Completed {current} / {total} items', '已完成 {count} 项': 'Completed {count} items',
  '正在同步笔记本': 'Synchronizing notebooks', '正在同步笔记': 'Synchronizing notes',
  '正在同步标签': 'Synchronizing tags', '正在上传本地更改': 'Uploading local changes',
  '正在同步图片和附件': 'Synchronizing images and attachments', '同步完成': 'Synchronization completed',
  '该共享笔记已被撤销访问权限': 'Access to this shared note was revoked',
  '展开我的空间': 'Expand workspace', '新建笔记本': 'New notebook', '收起我的空间': 'Collapse workspace',
  '搜索笔记本': 'Search notebooks', '{name} 菜单': '{name} menu',
  '新增子笔记本': 'New subnotebook', '重命名': 'Rename', '共享给我': 'Shared with me',
  '服务端不支持共享离线缓存': 'Server does not support offline shared-note cache',
  '暂无共享内容': 'Nothing shared yet', '标签': 'Tags', '{count}篇笔记': '{count} notes',
  '调节我的空间宽度': 'Resize workspace', '排序方式': 'Sort by',
  '修改时间：降序': 'Modified: newest first', '修改时间：升序': 'Modified: oldest first',
  '标题：降序': 'Title: Z–A', '标题：升序': 'Title: A–Z',
  '我的空间': 'Workspace', '新建富文本笔记': 'New rich-text note',
  '新建 Markdown 笔记': 'New Markdown note', '隐藏笔记栏': 'Hide note list',
  '搜索笔记': 'Search notes', '取消加星': 'Remove star', '加星': 'Star',
  '未命名': 'Untitled', '暂无摘要': 'No preview',
  '· 只读': '· Read-only', '这里还没有笔记': 'No notes here yet',
  '调节笔记栏宽度': 'Resize note list', '关闭': 'Close', '笔记标题': 'Note title',
  '笔记本：': 'Notebook:', '创建：': 'Created:', '修改：': 'Modified:', '格式：': 'Format:',
  '预览': 'Preview', '编辑': 'Edit', '移动/复制': 'Move / copy', '共享': 'Share',
  '历史': 'History', '导出': 'Export', '恢复': 'Restore',
  '笔记基本信息': 'Note details', '笔记标签': 'Note tags', '标签，以逗号分隔': 'Tags, separated by commas',
  '移动或复制笔记': 'Move or copy note', '选择笔记本': 'Select notebook',
  '移动': 'Move', '复制': 'Copy', '对方邮箱': 'Recipient email',
  '只读': 'Read-only', '可编辑': 'Can edit', '添加共享': 'Share', '取消共享': 'Stop sharing',
  '暂无附件': 'No attachments', '下载供离线使用': 'Download for offline use',
  '历史版本': 'History', '暂无历史': 'No history yet', '恢复到编辑器': 'Restore in editor',
  '日积字句，终得珠玑。': 'Small notes become lasting insights.',
  '选择笔记或开始全新记录。': 'Select a note or start a new one.',
  '新建笔记': 'New note',
  '响应不是有效 JSON [{path}]': 'Invalid JSON response [{path}]',
  '请求失败 ({status}) [{path}]': 'Request failed ({status}) [{path}]',
  '上传失败 ({status}) [{path}]': 'Upload failed ({status}) [{path}]',
  'NOTLOGIN': 'Please sign in again', 'invalidRequest': 'Invalid request',
  'invalidJSON': 'Invalid JSON request', 'noAuth': 'Permission denied',
  'notFound': 'Not found', 'offline': 'Server unavailable',
  'syncFailed': 'Sync failed', 'invalidNotebook': 'Invalid notebook',
  'noteIdNotExists': 'Note not found', 'notebookIdNotExists': 'Notebook not found',
  'conflict': 'The note changed elsewhere',
}

function initialLanguage(): Language {
  try {
    const saved = localStorage.getItem('gemsnote:language')
    if (saved === 'zh-CN' || saved === 'en-US') return saved
  } catch { /* Storage can be disabled; keep the browser preference. */ }
  return typeof navigator !== 'undefined' && /^zh\b/i.test(navigator.language) ? 'zh-CN' : 'en-US'
}

export const language = ref<Language>(initialLanguage())

export function setLanguage(next: Language) {
  language.value = next
  try { localStorage.setItem('gemsnote:language', next) } catch { /* Private browsing. */ }
}

export function t(key: string, variables: Record<string, string | number> = {}): string {
  const text = language.value === 'en-US' ? translations[key] || key : key
  return text.replace(/\{(\w+)\}/g, (_, name: string) => String(variables[name] ?? `{${name}}`))
}

export function formatDate(value: string | number | Date): string {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '' : date.toLocaleString(language.value)
}

watch(language, value => {
  if (typeof document !== 'undefined') document.documentElement.lang = value
}, { immediate: true })
