import {strFromU8, unzipSync} from 'fflate'

// Validate the whole input before creating anything. No archive entry is ever
// extracted to disk, and non-Markdown payloads are never decompressed.
export const importLimits = {archive:32*1024*1024, note:4*1024*1024, total:64*1024*1024, entries:5000, notes:2000, notebooks:1000, depth:20}
export interface ImportNote {path:string; title:string; content:string; folder:string}
export interface ImportSource {zip:boolean; title:string; folders:string[]; notes:ImportNote[]; ignored:number}
export class ImportError extends Error {
 constructor(public key:string){super(key);this.name='ImportError'}
}
function fail(key:string):never {throw new ImportError(key)}
const markdown=/\.(md|markdown)$/i
// Import planning is independent of the UI/system language. Preserve a stable
// path order across builds, with parent folders created before their children.
function comparePath(a:string,b:string){return a<b?-1:a>b?1:0}
function title(name:string){return name.replace(/\.(md|markdown|zip)$/i,'').trim()||'Untitled'}
function decode(bytes:Uint8Array){
 try{return new TextDecoder('utf-8',{fatal:true}).decode(bytes)}catch{return fail('Markdown 文件必须使用 UTF-8 编码')}
}
function safePath(name:string){
 if(!name||/[\x00-\x1f\x7f]/.test(name)||/^[\\/]/.test(name)||name.includes('\\')||/^[a-z]:/i.test(name))fail('ZIP 中包含不安全的文件路径')
 const parts=name.replace(/\/$/,'').split('/')
 if(parts.some(p=>!p||p==='.'||p==='..')||parts.length>importLimits.depth)fail('ZIP 中包含不安全的文件路径')
 return parts.join('/')
}
const crcTable=Uint32Array.from({length:256},(_,i)=>{for(let j=0;j<8;j++)i=(i>>>1)^((i&1)?0xedb88320:0);return i>>>0})
function crc32(bytes:Uint8Array){let crc=0xffffffff;for(const b of bytes)crc=(crc>>>8)^crcTable[(crc^b)&255];return (crc^0xffffffff)>>>0}
interface ZipEntry {name:string; path:string; directory:boolean; size:number; crc:number; keep:boolean}
function zipEntries(bytes:Uint8Array):ZipEntry[]{
 const view=new DataView(bytes.buffer,bytes.byteOffset,bytes.byteLength),len=bytes.length
 const u16=(n:number)=>view.getUint16(n,true),u32=(n:number)=>view.getUint32(n,true)
 let end=len-22
 while(end>=Math.max(0,len-65557)&&(u32(end)!==0x06054b50||end+22+u16(end+20)!==len))end--
 if(end<0||end<Math.max(0,len-65557))fail('ZIP 文件损坏或格式不受支持')
 const count=u16(end+10),offset=u32(end+16),centralSize=u32(end+12)
 if(u16(end+4)||u16(end+6)||u16(end+8)!==count||count===65535||offset+centralSize!==end)fail('ZIP 文件损坏或格式不受支持')
 if(count>importLimits.entries)fail('ZIP 文件数量超出导入限制')
 const entries:ZipEntry[]=[],seen=new Set<string>()
 let pos=offset,total=0,notes=0
 for(let i=0;i<count;i++){
  if(pos+46>end||u32(pos)!==0x02014b50)fail('ZIP 文件损坏或格式不受支持')
  const flags=u16(pos+8),method=u16(pos+10),compressed=u32(pos+20),size=u32(pos+24),nameSize=u16(pos+28),extra=u16(pos+30),comment=u16(pos+32),local=u32(pos+42)
  const next=pos+46+nameSize+extra+comment
  if(next>end||u16(pos+34)||flags&1||(u32(pos+38)>>>16&0xf000)===0xa000)fail('ZIP 文件损坏或格式不受支持')
  let name:string
  try{name=strFromU8(bytes.subarray(pos+46,pos+46+nameSize),!(flags&2048))}catch{return fail('ZIP 文件损坏或格式不受支持')}
  const path=safePath(name),directory=name.endsWith('/')
  if(seen.has(path))fail('ZIP 中存在重复文件路径')
  seen.add(path)
  if(local+30>offset||u32(local)!==0x04034b50||u16(local+8)!==method||u16(local+6)!==flags)fail('ZIP 文件损坏或格式不受支持')
  const start=local+30+u16(local+26)+u16(local+28)
  if(start+compressed>offset||start<local+30)fail('ZIP 文件损坏或格式不受支持')
  const ignored=path.split('/').some(p=>p==='__MACOSX'||p.startsWith('._')||p==='.DS_Store')
  const keep=!ignored&&!directory&&markdown.test(path)
  if(keep){
   if(method!==0&&method!==8)fail('ZIP 文件损坏或格式不受支持')
   if(size>importLimits.note||(total+=size)>importLimits.total)fail('Markdown 内容大小超出导入限制')
   if(++notes>importLimits.notes)fail('ZIP 文件数量超出导入限制')
  }
  entries.push({name,path,directory,size,crc:u32(pos+16),keep});pos=next
 }
 if(pos!==end)fail('ZIP 文件损坏或格式不受支持')
 return entries
}
export function parseMarkdownImport(name:string,bytes:Uint8Array):ImportSource {
 if(markdown.test(name)){
  if(bytes.length>importLimits.note)fail('Markdown 内容大小超出导入限制')
  const base=name.split(/[\\/]/).pop()!
  return {zip:false,title:title(base),folders:[],notes:[{path:base,title:title(base),content:decode(bytes),folder:''}],ignored:0}
 }
 if(!/\.zip$/i.test(name))fail('请选择 Markdown 或 ZIP 文件')
 if(bytes.length>importLimits.archive)fail('ZIP 文件大小不能超过 32 MiB')
 let entries:ZipEntry[],files:Record<string,Uint8Array>
 try{
  entries=zipEntries(bytes)
  const selected=new Set(entries.filter(e=>e.keep).map(e=>e.name))
  files=unzipSync(bytes,{filter:entry=>selected.has(entry.name)})
 }catch(e){if(e instanceof ImportError)throw e;return fail('ZIP 文件损坏或格式不受支持')}
 const folders=new Set<string>(),notes:ImportNote[]=[]
 function addFolder(path:string){const parts=path.split('/');for(let i=1;i<=parts.length;i++)folders.add(parts.slice(0,i).join('/'))}
 for(const entry of entries){
  if(entry.directory&&!entry.path.split('/').some(p=>p==='__MACOSX'||p.startsWith('._')))addFolder(entry.path)
  if(!entry.keep)continue
  const data=files[entry.name]
  if(!data||data.length!==entry.size||crc32(data)!==entry.crc)fail('ZIP 文件损坏或格式不受支持')
  const index=entry.path.lastIndexOf('/'),folder=index<0?'':entry.path.slice(0,index)
  if(folder)addFolder(folder)
  notes.push({path:entry.path,title:title(entry.path.slice(index+1)),folder,content:decode(data)})
 }
 if(!notes.length)fail('ZIP 中没有 Markdown 文件')
 if(folders.size>importLimits.notebooks-1)fail('ZIP 文件数量超出导入限制')
 return {zip:true,title:title(name),folders:[...folders].sort((a,b)=>a.split('/').length-b.split('/').length||comparePath(a,b)),notes:notes.sort((a,b)=>comparePath(a.path,b.path)),ignored:entries.filter(e=>!e.directory&&!e.keep).length}
}
