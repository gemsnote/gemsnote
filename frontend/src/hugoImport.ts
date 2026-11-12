import {parse as yaml} from 'yaml'
import {parse as toml} from 'smol-toml'

export interface HugoNote {title:string;content:string;tags:string[];createdTime?:string;updatedTime?:string}
export function parseHugoMarkdown(content:string,fallbackTitle:string):HugoNote {
 content=content.replace(/^\ufeff/,'')
 let metadata:unknown,end=0
 const opening=/^(---|\+\+\+)[ \t]*\r?\n/.exec(content)
 if(opening){
  const closing=new RegExp('^'+(opening[1]==='---'?'---':'\\+\\+\\+')+'[ \\t]*(?:\\r?\\n|$)','gm')
  closing.lastIndex=opening[0].length
  const match=closing.exec(content)
  if(!match)return {title:fallbackTitle,content,tags:[]}
  try{const header=content.slice(opening[0].length,match.index);metadata=opening[1]==='---'?yaml(header,{maxAliasCount:0}):toml(header)}catch{throw Error('Hugo 元数据格式无效')}
  end=closing.lastIndex
 }else if(content.startsWith('{')){
  let depth=0,quoted=false,escaped=false
  for(let i=0;i<content.length;i++){
   const ch=content[i]
   if(quoted){if(escaped)escaped=false;else if(ch==='\\')escaped=true;else if(ch==='"')quoted=false}
   else if(ch==='"')quoted=true
   else if(ch==='{')depth++
   else if(ch==='}'&&--depth===0){
    try{metadata=JSON.parse(content.slice(0,i+1))}catch{return {title:fallbackTitle,content,tags:[]}}
    end=i+1;if(content.slice(end,end+2)==='\r\n')end+=2;else if(content[end]==='\n')end++
    break
   }
  }
 }
 if(!end)return {title:fallbackTitle,content,tags:[]}
 if(!metadata||typeof metadata!=='object'||Array.isArray(metadata))throw Error('Hugo 元数据格式无效')
 const fields=Object.fromEntries(Object.entries(metadata).map(([key,value])=>[key.toLowerCase(),value]))
 const date=(value:unknown)=>{
  if(value==null||value==='')return undefined
  const text=value instanceof Date?value.toISOString():String(value)
  // Zone-less Hugo dates are interpreted as UTC on every client.
  const normalized=/^\d{4}-\d\d-\d\d$/.test(text)?text+'T00:00:00Z':text.replace(' ','T').replace(/([T]\d\d:\d\d(?::\d\d(?:\.\d+)?)?)$/,'$1Z')
  const parsed=new Date(normalized)
  if(!/^\d{4}-\d\d-\d\d/.test(text)||!Number.isFinite(parsed.getTime()))throw Error('Hugo 日期无效')
  return parsed.toISOString()
 }
 const createdTime=date(fields.date),updatedTime=date(fields.lastmod??fields.modified)
 if(!createdTime&&!updatedTime)throw Error('Hugo 元数据缺少创建或修改时间')
 if(fields.title!=null&&typeof fields.title!=='string')throw Error('Hugo 标题必须是文本')
 const rawTags=fields.tags==null?[]:typeof fields.tags==='string'?[fields.tags]:fields.tags
 if(!Array.isArray(rawTags)||rawTags.some(tag=>typeof tag!=='string'||tag.includes(',')))throw Error('Hugo 标签格式无效')
 return {title:fields.title??'',content:content.slice(end).replace(/^\r?\n/,''),tags:[...new Set(rawTags.map(tag=>tag.trim()).filter(Boolean))],createdTime:createdTime??updatedTime,updatedTime:updatedTime??createdTime}
}
