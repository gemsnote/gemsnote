import {afterEach,describe,expect,it} from 'vitest'
import {ApiError,errorCode,errorMessage} from './errors'
import {language,setLanguage,t} from './i18n'

const original=language.value
afterEach(()=>setLanguage(original))
describe('localized errors',()=>{
 it.each(['syncFailed','confirmationRequired','noPermission','noFile','invalidParentNotebook','sharedReadOnly','groupHasUsers','oldPasswordError','avatarCacheFailed','uploadKeyConflict'])('translates %s in both languages without changing the code',code=>{
  const results=[]
  for(const lang of ['zh-CN','en-US'] as const){
   setLanguage(lang)
   const error=new ApiError(code)
   expect(errorCode(error)).toBe(code)
   expect(error.message).not.toBe(code)
   expect(errorMessage(error)).toBe(error.message)
   results.push(error.message)
  }
  expect(results[0]).not.toBe(results[1])
 })
 it.each([
  ['update note failed: conflict','conflict'],
  ['get note content abc: context deadline exceeded (Client.Timeout exceeded while reading body)','requestTimeout'],
  ['server unavailable before reset: Get "https://example.test/api2": dial tcp: connection refused','connectionRefused'],
  ['SQLITE_BUSY: database is locked','databaseBusy'],
  ['pending attachment original missing for note abc','localFileMissing'],
  ['download attachment abc failed (HTTP 404)','httpFailure'],
 ])('localizes wrapped errors: %s',(raw,key)=>{
  for(const lang of ['zh-CN','en-US'] as const){setLanguage(lang);expect(errorMessage(new Error(raw))).toBe(t(key,{status:404}))}
 })
 it('preserves translated messages with dynamic values and can change language',()=>{
  setLanguage('zh-CN')
  const raw=t('请求失败 ({status}) [{path}]',{status:404,path:'/api2/attachments'})
  expect(errorMessage(new Error(raw))).toBe(raw)
  setLanguage('en-US')
  expect(errorMessage(new Error(raw))).toBe(t('请求失败 ({status}) [{path}]',{status:404,path:'/api2/attachments'}))
 })
 it('translates server Chinese errors into English',()=>{
  setLanguage('en-US')
  expect(errorMessage('不能分享给自己')).toBe(t('cannotShareSelf'))
  expect(errorMessage('笔记本下有子笔记本')).toBe(t('notebookHasChildren'))
 })
 it('does not show unknown response bodies, credentials or Error prefixes',()=>{
  for(const lang of ['zh-CN','en-US'] as const){setLanguage(lang);expect(errorMessage(new Error('unexpected failure: token=private-secret'))).toBe(t('operationFailed'))}
 })
})
