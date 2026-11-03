import {knownMessage,t} from './i18n'

// Control flow uses code, not the localized message (logout, reset confirmation,
// conflict handling). Keep the original error available without rendering it.
export class ApiError extends Error {
 constructor(public code:string){super(errorMessage(code));this.name='ApiError'}
}
export function errorCode(error:unknown):string {
 if(error instanceof ApiError)return error.code
 if(error&&typeof error==='object'){
  const value=error as {key?:unknown;Msg?:unknown;message?:unknown}
  for(const candidate of [value.key,value.Msg,value.message])if(typeof candidate==='string')return candidate
 }
 return typeof error==='string'?error:''
}
export function errorMessage(error:unknown):string {
 const raw=errorCode(error).replace(/^(?:Error|TypeError):\s*/,'').trim()
 const known=knownMessage(raw)
 if(known!==undefined)return known
 // Network / SQLite failures may be wrapped several times by Go's sync layer.
 if(/deadline exceeded|timed? out|timeout/i.test(raw))return t('requestTimeout')
 if(/connection refused|no such host|network is unreachable|failed to fetch|networkerror|load failed/i.test(raw))return t('connectionRefused')
 if(/certificate|x509:/i.test(raw))return t('certificateError')
 if(/database is locked|database is busy|sqlite[_ ](?:busy|locked)/i.test(raw))return t('databaseBusy')
 if(/pending attachment original missing|no such file or directory/i.test(raw))return t('localFileMissing')
 if(/no complete local content cache/i.test(raw))return t('localContentMissing')
 if(/context canceled|operation was aborted/i.test(raw))return t('operationCancelled')
 // API errors are often wrapped as "update note failed: conflict". Match a
 // complete trailing code only; never interpolate the original response body.
 const trailing=raw.match(/:\s*([a-zA-Z][a-zA-Z0-9]*)$/)?.[1]
 if(trailing){const value=knownMessage(trailing);if(value!==undefined)return value}
 const status=raw.match(/(?:HTTP\s+|status=)([45]\d\d)\b/i)?.[1]
 if(status)return t('httpFailure',{status})
 return t('operationFailed')
}
