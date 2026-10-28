// Real downloads against an isolated PostgreSQL server and its disposable files.
import assert from 'node:assert/strict'
import {randomBytes} from 'node:crypto'
import {realpath,rename,mkdir,rmdir} from 'node:fs/promises'
import {resolve,sep} from 'node:path'

const base=process.env.GEMSNOTE_API2_URL,serverDir=process.env.GEMSNOTE_DISPOSABLE_SERVER_DIR
if(process.env.GEMSNOTE_DISPOSABLE_API2!=='1'||base!=='http://127.0.0.1:19090'||!/^\/tmp\/gemsnote-api2-live\.[a-zA-Z0-9]+\/server$/.test(serverDir||''))throw Error('Requires the disposable API2 fixture and its /tmp server directory')
const root=await realpath(serverDir)
const login=await (await fetch(`${base}/api2/auth/login`,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({email:'admin',pwd:'gemsnote'})})).json()
assert.equal(login.Ok,true)
const headers={Authorization:`Bearer ${login.Token}`}
const post=async(path,data)=>await (await fetch(base+path,{method:'POST',headers:{...headers,'Content-Type':'application/json'},body:JSON.stringify(data)})).json()
const boot=await (await fetch(`${base}/api2/bootstrap`,{headers})).json()
const noteId=randomBytes(12).toString('hex')
const note=await post('/api2/save',{noteId,notebookId:boot.Notebooks[0].NotebookId,title:'download regression',content:'# attachment',tags:'',isNew:true,isMarkdown:true})
assert.equal(note.Note.NoteId,noteId)
const bytes=Buffer.from([0,1,2,3,4,0xff,0xfe,0x80]),form=new FormData()
form.append('noteId',noteId)
form.append('file',new Blob([bytes]),'原件.gemsnote-test')
const upload=await (await fetch(`${base}/api2/attachments/upload`,{method:'POST',headers,body:form})).json()
assert.equal(upload.Ok,true)
const listing=await post('/api2/attachments',{noteId}),metadata=listing.List.find(f=>f.AttachId===upload.Id)
assert(metadata)
assert(metadata.Path.startsWith('files/'))
const filePath=await realpath(resolve(root,metadata.Path))
assert(filePath.startsWith(root+sep+'files'+sep))
const endpoints=[`/api2/file/getAttach?fileId=${upload.Id}`,`/api2/attachments/download?attachId=${upload.Id}`,`/api/file/getAttach?fileId=${upload.Id}&token=${encodeURIComponent(login.Token)}`]
for(const path of endpoints){
 const response=await fetch(base+path,{headers})
 assert.equal(response.status,200)
 assert.match(response.headers.get('content-disposition')||'',/^attachment;/)
 assert.deepEqual(Buffer.from(await response.arrayBuffer()),bytes)
 const partial=await fetch(base+path,{headers:{...headers,Range:'bytes=2-5'}})
 assert.equal(partial.status,206)
 assert.deepEqual(Buffer.from(await partial.arrayBuffer()),bytes.subarray(2,6))
}
const backup=filePath+'.download-test-backup'
await rename(filePath,backup)
try{
 for(const path of endpoints){
  const response=await fetch(base+path,{headers})
  assert.equal(response.status,404)
  assert.equal(await response.text(),'fileNotFound')
  assert.equal(response.headers.get('content-disposition'),null)
 }
 await mkdir(filePath)
 try{
  const response=await fetch(base+endpoints[0],{headers})
  assert.equal(response.status,404)
  assert.equal(await response.text(),'fileNotFound')
 }finally{await rmdir(filePath)}
}finally{await rename(backup,filePath)}
const restored=await fetch(base+endpoints[0],{headers})
assert.equal(restored.status,200)
assert.deepEqual(Buffer.from(await restored.arrayBuffer()),bytes)
console.log('Disposable PostgreSQL: API2/client, API2/Web and legacy attachment downloads, Range, missing files, directories and recovery passed')
