import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { objectId, request, resolveAvatarUrl, upload } from './api'
import {language,setLanguage,t,formatDate} from './i18n'

afterEach(()=>vi.unstubAllGlobals())
describe('interface language',()=>{
 it('switches Chinese and English without changing note data',()=>{
  const original=language.value
  try{
   setLanguage('en-US')
   expect(t('所有笔记')).toBe('All notes')
   expect(t('确定删除分组“{name}”吗？',{name:'Research'})).toBe('Delete group “Research”?')
   expect(formatDate('2026-09-18T08:30:00Z')).not.toContain('Invalid')
   setLanguage('zh-CN')
   expect(t('所有笔记')).toBe('所有笔记')
  }finally{setLanguage(original)}
 })
 it('persists the selected language',()=>{
  const original=language.value
  const setItem=vi.fn()
  vi.stubGlobal('localStorage',{setItem})
  try{
   setLanguage('en-US')
   expect(setItem).toHaveBeenCalledWith('gemsnote:language','en-US')
  }finally{setLanguage(original)}
 })
})
describe('API2 adapter',()=>{
 beforeEach(()=>setLanguage('zh-CN'))
 it('restores a stable history version with an expected USN through API2',async()=>{
  const fetch=vi.fn().mockResolvedValue({ok:true,json:async()=>({Ok:true})});vi.stubGlobal('fetch',fetch)
  await request('/noteContentHistory/restore',{noteId:'note',historyId:'stable',usn:12})
  expect(fetch.mock.calls[0][0]).toBe('/api2/noteContentHistory/restore')
  const options=fetch.mock.calls[0][1]
  expect(options.method).toBe('POST')
  expect(Object.fromEntries(options.body)).toEqual({noteId:'note',historyId:'stable',usn:'12'})
  expect(options.headers['X-Requested-With']).toBe('XMLHttpRequest')
 })
 it('surfaces missing historical files instead of accepting a partial restore',async()=>{
  vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,json:async()=>({Ok:false,Msg:'historyFileMissing'})}))
  await expect(request('/noteContentHistory/restore',{noteId:'n',historyId:'h',usn:1})).rejects.toThrow('历史文件已丢失')
 })

 it('stores the API2 token and sends it as a bearer credential',async()=>{
  const storage=new Map<string,string>()
  vi.stubGlobal('sessionStorage',{getItem:(key:string)=>storage.get(key)||null,setItem:(key:string,value:string)=>storage.set(key,value),removeItem:(key:string)=>storage.delete(key)})
  const fetch=vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>({Ok:true,Token:'jwt-token'})});vi.stubGlobal('fetch',fetch)
  await request('/doLogin',{email:'admin',pwd:'secret'})
  expect(fetch.mock.calls[0][0]).toBe('/api2/auth/login')
  expect(storage.get('gemsnote:api2-token')).toBe('jwt-token')
  await request('/web/bootstrap')
  expect(fetch.mock.calls[1][1].headers.Authorization).toBe('Bearer jwt-token')
  expect(fetch.mock.calls[0][1].headers['Content-Type']).toBe('application/json')
 })
 it('encodes indexed arrays and leaves password text unchanged',async()=>{
  const fetch=vi.fn().mockResolvedValue({ok:true,json:async()=>({Ok:true})});vi.stubGlobal('fetch',fetch)
  await request('/share/addShareNote',{emails:['a@b.test','c@d.test'],pwd:'a&b+c'})
  const options=fetch.mock.calls[0][1]
  expect(options.body.get('emails[1]')).toBe('c@d.test')
  expect(options.body.get('pwd')).toBe('a&b+c')
  expect(options.credentials).toBe('same-origin')
 })
  it('does not treat an unsuccessful HTTP-200 write as saved',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>({Ok:false,Msg:'conflict'})}))
   await expect(request('/web/save',{})).rejects.toThrow('conflict')
  })
  it('uses the versioned JSON API for note writes', async()=>{
   const fetch=vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>({Ok:true})});vi.stubGlobal('fetch',fetch)
   await request('/web/save',{noteId:'0123456789abcdef01234567',title:'标题',isNew:false})
   expect(fetch.mock.calls[0][0]).toBe('/api2/save')
   expect(fetch.mock.calls[0][1].headers['Content-Type']).toBe('application/json')
   expect(JSON.parse(fetch.mock.calls[0][1].body).title).toBe('标题')
  })
  it('surfaces a friendly error when the backend answers with non-JSON',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:false,status:404,json:async()=>{throw new SyntaxError('Unexpected token <')}}))
   await expect(request('/web/document',{noteId:'x'})).rejects.toThrow('请求失败 (404)')
  })
  it('does not report success for an HTTP-200 HTML response',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>{throw new SyntaxError('Unexpected token <')}}))
   await expect(request('/user/updateUsername',{username:'new'})).rejects.toThrow('响应不是有效 JSON')
  })
  it('accepts an empty history array from API2',async()=>{
   const fetch=vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>[]});vi.stubGlobal('fetch',fetch)
   await expect(request('/noteContentHistory/listHistories',{noteId:'0123456789abcdef01234567'})).resolves.toEqual([])
   expect(fetch.mock.calls[0][0]).toBe('/api2/noteContentHistory/listHistories')
  })
  it('rejects bare-false upload responses instead of reporting success',async()=>{
   vi.stubGlobal('fetch',vi.fn().mockResolvedValue({ok:true,status:200,json:async()=>false}))
   await expect(upload('/attach/uploadAttach',new File(['x'],'a.txt'),{noteId:'x'})).rejects.toThrow('上传失败')
  })
 it('generates Mongo-compatible identifiers without duplicates',()=>{
  const ids=Array.from({length:100},objectId)
  expect(new Set(ids).size).toBe(100)
  ids.forEach(id=>expect(id).toMatch(/^[0-9a-f]{24}$/))
 })
 it('uses the correct avatar URL for server paths and desktop image IDs',()=>{
  expect(resolveAvatarUrl('public/upload/u/images/logo/new.gif')).toBe('/public/upload/u/images/logo/new.gif')
  expect(resolveAvatarUrl('/public/upload/u/images/logo/new.gif')).toBe('/public/upload/u/images/logo/new.gif')
  expect(resolveAvatarUrl('0123456789abcdef01234567')).toBe('/api2/file/getImage?fileId=0123456789abcdef01234567')
 })
})
