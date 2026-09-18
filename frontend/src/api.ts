import {language,t} from './i18n'

export type Params = Record<string, string | number | boolean | string[] | undefined>
// All browser data calls use the versioned API namespace.  The server keeps
// /api and the historical controller routes for legacy clients, while the
// Vue application is migrated independently.
function api2Path(path: string): string {
  if (path.startsWith('/api2/')) return path
  if (path.startsWith('/api/')) return `/api2/${path.slice('/api/'.length)}`
  const core: Record<string, string> = {
    '/web/bootstrap': '/api2/bootstrap', '/web/logout': '/api2/logout', '/doLogin': '/api2/auth/session', '/doRegister': '/api2/auth/register', '/doFindPassword': '/api2/auth/password/request', '/findPasswordUpdate': '/api2/auth/password/reset', '/web/groups': '/api2/groups', '/web/adminData': '/api2/admin/data', '/file/uploadAvatar': '/api2/avatar', '/attach/getAttachs': '/api2/attachments', '/attach/uploadAttach': '/api2/attachments/upload', '/attach/deleteAttach': '/api2/attachments/delete', '/web/notes': '/api2/notes',
    '/web/star': '/api2/star', '/web/document': '/api2/document',
    '/web/save': '/api2/save', '/web/restore': '/api2/restore'
  }
  if (core[path]) return core[path]
  return `/api2${path}`
}
export async function request<T = any>(path: string, params?: Params): Promise<T> {
  const body = new URLSearchParams()
  for (const [key, value] of Object.entries(params || {})) {
    if (value === undefined) continue
    if (Array.isArray(value)) value.forEach((v, i) => body.append(`${key}[${i}]`, v))
    else body.append(key, String(value))
  }
  const target = api2Path(path)
  const core = target === '/api2/auth/session' || target === '/api2/auth/register' || target === '/api2/auth/password/request' || target === '/api2/auth/password/reset' || target === '/api2/notes' || target === '/api2/star' || target === '/api2/document' || target === '/api2/save' || target === '/api2/restore' || target === '/api2/admin/data' || target === '/api2/attachments' || target === '/api2/attachments/delete'
  const response = await fetch(target, { method: params ? 'POST' : 'GET', credentials: 'same-origin', headers: { 'X-Requested-With': 'XMLHttpRequest', 'Accept-Language': language.value, ...(core ? {'Content-Type':'application/json'} : {}) }, body: params ? (core ? JSON.stringify(params) : body) : undefined })
  const data = await response.json().catch(() => null)
  if (data === null) throw new Error(response.ok ? t('响应不是有效 JSON [{path}]',{path:target}) : t('请求失败 ({status}) [{path}]',{status:response.status,path:target}))
  if (!response.ok || data === false || data?.Ok === false) {
    if (data?.Msg === 'NOTLOGIN') window.dispatchEvent(new Event('session-expired'))
    throw new Error(data?.Msg ? t(data.Msg) : t('请求失败 ({status}) [{path}]',{status:response.status,path:target}))
  }
  return data
}
export function objectId(): string {
  return Math.floor(Date.now() / 1000).toString(16).padStart(8, '0') + Array.from(crypto.getRandomValues(new Uint8Array(8)), b => b.toString(16).padStart(2, '0')).join('')
}

export function resolveAvatarUrl(logo: string): string {
  if (/^[a-f\d]{24}$/i.test(logo)) return `/api2/file/getImage?fileId=${logo}`
  if (/^\/?public\/upload\//.test(logo)) return `/${logo.replace(/^\/+/, '')}`
  return logo
}

export async function upload(path: string, file: File, params: Record<string, string> = {}) {
  const body = new FormData()
  body.append('file', file)
  Object.entries(params).forEach(([key, value]) => body.append(key, value))
  const target = api2Path(path)
  const response = await fetch(target, { method: 'POST', credentials: 'same-origin', headers: { 'X-Requested-With': 'XMLHttpRequest', 'Accept-Language': language.value }, body })
  const data = await response.json().catch(() => null)
  if (!response.ok || !data || data.Ok === false) throw new Error(data?.Msg ? t(data.Msg) : t('上传失败 ({status}) [{path}]',{status:response.status,path:target}))
  return data
}
