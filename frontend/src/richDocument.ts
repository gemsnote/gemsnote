import {DOMParser} from '@xmldom/xmldom'

export const richWarning='复杂排版、内嵌图片和对象可能简化；图片以占位文字保留，请核对导入结果。'
export function escapeHtml(value:string){return value.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}
function safeLink(value:string){return /^(https?:\/\/|mailto:)/i.test(value)&&!/[\x00-\x20"<>]/.test(value)}
function rtfText(value:string){let result='';for(let i=0;i<value.length;i++){const n=value.charCodeAt(i);result+=n===92||n===123||n===125?'\\'+value[i]:n===10?'\\line ':n===9?'\\tab ':n>=32&&n<127?value[i]:n===13?'':`\\u${n>32767?n-65536:n}?`}return result}
function elements(node:any):any[]{return Array.from(node.childNodes||[]).filter((n:any)=>n.nodeType===1)}
function tag(node:any){return (node.localName||node.nodeName||'').toLowerCase()}
function children(node:any,name:string){return elements(node).filter(n=>tag(n)===name.toLowerCase())}
function child(node:any,name:string){return children(node,name)[0]}
function attr(node:any,name:string){if(!node)return '';for(const a of Array.from(node.attributes||[]) as any[])if((a.localName||a.name)===name)return a.value;return ''}
function css(node:any){return Object.fromEntries(String(node.getAttribute?.('style')||'').split(';').map(s=>s.split(':').map(x=>x.trim().toLowerCase())).filter(a=>a.length===2))}
function rgb(value:string):string|undefined{
 if(/^#[0-9a-f]{6}$/i.test(value))return value.toLowerCase()
 if(/^#[0-9a-f]{3}$/i.test(value))return '#'+[...value.slice(1)].map(c=>c+c).join('')
 const match=/^rgb\(\s*(\d+),\s*(\d+),\s*(\d+)\s*\)$/.exec(value)
 if(match)return '#'+match.slice(1).map(n=>Math.min(255,+n).toString(16).padStart(2,'0')).join('')
 return ({red:'#ff0000',blue:'#0000ff',green:'#008000',black:'#000000',white:'#ffffff'} as Record<string,string>)[value]
}
export function htmlToRtf(html:string,title=''):string{
 const document=new DOMParser().parseFromString('<html><body>'+html+'</body></html>','text/html'),colors:string[]=[]
 function color(value:string){const hex=rgb(value);if(!hex)return '';let i=colors.indexOf(hex);if(i<0){colors.push(hex);i=colors.length-1}return `\\cf${i+1} `}
 function render(node:any,depth=0):string{
  if(depth>100)throw Error('documentTooComplex')
  if(node.nodeType===3)return rtfText(node.data)
  if(node.nodeType!==1&&node.nodeType!==9)return ''
  const name=tag(node);if(['script','style','iframe','object','head'].includes(name))return ''
  const style=css(node),body=()=>Array.from(node.childNodes||[]).map(n=>render(n,depth+1)).join('')
  if(name==='br')return '\\line '
  if(name==='hr')return '\\par ____________________\\par '
  if(name==='img')return rtfText(`[${node.getAttribute('alt')||'Image'}] ${node.getAttribute('src')||''}`)
  if(name==='table')return elements(node).map(n=>render(n,depth+1)).join('')+'\\par '
  if(name==='tr'){
   const cells=elements(node).filter(n=>['td','th'].includes(tag(n)))
   return '\\trowd '+cells.map((_,i)=>`\\cellx${(i+1)*Math.floor(9000/Math.max(1,cells.length))} `).join('')+cells.map(n=>'\\intbl '+render(n,depth+1)+'\\cell ').join('')+'\\row '
  }
  let controls=''
  if(['b','strong','th'].includes(name)||['bold','700','800','900'].includes(style['font-weight']))controls+='\\b '
  if(['i','em'].includes(name)||style['font-style']==='italic')controls+='\\i '
  if(name==='u'||style['text-decoration']?.includes('underline'))controls+='\\ul '
  if(['s','strike','del'].includes(name)||style['text-decoration']?.includes('line-through'))controls+='\\strike '
  if(name==='sup')controls+='\\super '
  if(name==='sub')controls+='\\sub '
  const heading=/^h([1-6])$/.exec(name);if(heading)controls+=`\\b\\fs${48-(+heading[1]-1)*4} `
  const size=/^(\d+(?:\.\d+)?)(pt|px)$/.exec(style['font-size']||'');if(size)controls+=`\\fs${Math.max(12,Math.min(144,Math.round(+size[1]*(size[2]==='px'?1.5:2))))} `
  controls+=color(style.color||node.getAttribute?.('color')||'')
  const align=style['text-align']||node.getAttribute?.('align');if(align)controls+=({left:'\\ql ',center:'\\qc ',right:'\\qr ',justify:'\\qj '} as Record<string,string>)[align]||''
  let text=body()
  if(name==='a'&&safeLink(node.getAttribute('href')||''))text=`{\\field{\\*\\fldinst HYPERLINK "${rtfText(node.getAttribute('href'))}"}{\\fldrslt ${text}}}`
  if(name==='li'){
   const parent=node.parentNode;const index=elements(parent).filter(n=>tag(n)==='li').indexOf(node)+1
   text=rtfText(tag(parent)==='ol'?`${index}. `:'• ')+text
  }
  if(['p','div','li','blockquote','pre'].includes(name)||heading)text+='\\par '
  return '{'+controls+text+'}'
 }
 const body=render(document)
 const colorTable=colors.map(hex=>`\\red${parseInt(hex.slice(1,3),16)}\\green${parseInt(hex.slice(3,5),16)}\\blue${parseInt(hex.slice(5,7),16)};`).join('')
 return `{\\rtf1\\ansi\\ansicpg1252\\deff0\\uc1{\\fonttbl{\\f0\\fnil Arial;}}{\\colortbl;${colorTable}}{\\info{\\title ${rtfText(title)}}}\\f0\\fs24 ${body}}`
}

type RtfState={skip:boolean;uc:number;b:boolean;i:boolean;ul:boolean;strike:boolean;fs:number;cf:number;align:string;link:string;field?:{instruction:string};instruction:boolean}
export function rtfToHtml(bytes:Uint8Array):string{
 let source='';for(let i=0;i<bytes.length;i+=32768)source+=String.fromCharCode(...bytes.subarray(i,i+32768))
 if(!/^\{\\rtf1(?:\D|$)/.test(source)||source.length>32*1024*1024)throw Error('invalidRichDocument')
 const colorTable=/\{\\colortbl([^{}]*)\}/.exec(source)?.[1]||''
 const colors=colorTable.split(';').map(c=>{const get=(key:string)=>Math.min(255,+(new RegExp('\\\\'+key+'(\\d+)').exec(c)?.[1]||0));return `rgb(${get('red')},${get('green')},${get('blue')})`})
 let state:RtfState={skip:false,uc:1,b:false,i:false,ul:false,strike:false,fs:24,cf:0,align:'',link:'',instruction:false},stack:RtfState[]=[],fallback=0,pending:number[]=[],codepage=1252,line='',result='',surrogate='',row:string[]|undefined,table=false
 const flushParagraph=()=>{if(row){if(line)line+='<br>';return}if(table){result+='</table>';table=false}if(line){result+=`<p${state.align?' style="text-align:'+state.align+'"':''}>${line}</p>`;line=''}}
 function text(value:string){
  value=surrogate+value;surrogate='';if(value.length&&/[\uD800-\uDBFF]/.test(value.at(-1)!)){surrogate=value.at(-1)!;value=value.slice(0,-1)}if(!value)return
  if(state.instruction){if(state.field)state.field.instruction+=value;return}
  if(state.skip)return
  let html=escapeHtml(value)
  if(state.b)html='<strong>'+html+'</strong>';if(state.i)html='<em>'+html+'</em>';if(state.ul)html='<u>'+html+'</u>';if(state.strike)html='<s>'+html+'</s>'
  const style=(state.fs!==24?`font-size:${state.fs/2}pt;`:'')+(state.cf>0&&colors[state.cf]?`color:${colors[state.cf]};`:'')
  if(style)html=`<span style="${style}">${html}</span>`
  if(state.link)html=`<a href="${escapeHtml(state.link)}">${html}</a>`
  line+=html
  if(line.length+result.length>32*1024*1024)throw Error('documentTooComplex')
 }
 function flush(){if(pending.length){let decoder:TextDecoder;try{decoder=new TextDecoder(({936:'gbk',950:'big5',932:'shift_jis',65001:'utf-8'} as Record<number,string>)[codepage]||`windows-${codepage}`,{fatal:true})}catch{throw Error('unsupportedRtfEncoding')}try{text(decoder.decode(new Uint8Array(pending)))}catch{throw Error('unsupportedRtfEncoding')}pending=[]}}
 function byte(n:number){if(fallback){fallback--;return}if(state.skip&&!state.instruction)return;pending.push(n)}
 const skipped=new Set(['fonttbl','colortbl','stylesheet','info','header','headerl','headerr','footer','footerl','footerr','listtable','listoverridetable','generator','filetbl','datastore','xmlnstbl','themedata','colorschememapping'])
 let pos=0,closed=false
 while(pos<source.length){
  const ch=source[pos++]
  if(closed){if(ch.trim())throw Error('invalidRichDocument');continue}
  if(ch==='{'||ch==='}'){
   flush();fallback=0
   if(ch==='{'){if(stack.length>=128)throw Error('documentTooComplex');stack.push({...state})}
   else{if(!stack.length)throw Error('invalidRichDocument');if(state.instruction&&state.field){const match=/HYPERLINK\s+"([^"\r\n]+)"/i.exec(state.field.instruction);if(match&&safeLink(match[1]))state.field.instruction=match[1]}
    // End the paragraph before paragraph formatting reverts at a group boundary.
    if(state.align!==stack[stack.length-1].align&&line&&!row)flushParagraph()
    state=stack.pop()!;if(!stack.length)closed=true
   }
   continue
  }
  if(ch==='\r'||ch==='\n')continue
  if(ch!=='\\'){byte(ch.charCodeAt(0));continue}
  const next=source[pos++]
  if(!next)throw Error('invalidRichDocument')
  if(next==="'"){const hex=source.slice(pos,pos+2);if(!/^[a-f\d]{2}$/i.test(hex))throw Error('invalidRichDocument');byte(parseInt(hex,16));pos+=2;continue}
  if(next==='\\'||next==='{'||next==='}'){byte(next.charCodeAt(0));continue}
  flush()
  if(next==='*'){state.skip=true;continue}
  if(next==='~'){text('\u00a0');continue}
  if(next==='-'){text('\u00ad');continue}
  if(next==='_'){text('\u2011');continue}
  if(!/[a-z]/i.test(next))continue
  const start=pos-1;while(/[a-z]/i.test(source[pos]||'')&&pos<source.length)pos++
  const word=source.slice(start,pos);const number=/^-?\d+/.exec(source.slice(pos));const n=number?+number[0]:1;if(number)pos+=number[0].length;if(source[pos]===' ')pos++
  if(word==='bin'){if(!number||n<0||pos+n>source.length)throw Error('invalidRichDocument');pos+=n;continue}
  if(skipped.has(word)){state.skip=true;continue}
  if(word==='pict'||word==='object'){if(!state.skip)text('[Image / object]');state.skip=true;continue}
  if(word==='field'){state.field={instruction:''};continue}
  if(word==='fldinst'){state.instruction=true;state.skip=true;continue}
  if(word==='fldrslt'){state.skip=false;state.instruction=false;state.link=state.field&&safeLink(state.field.instruction)?state.field.instruction:'';continue}
  if(word==='uc'){state.uc=Math.max(0,Math.min(16,n));continue}
  if(word==='u'){if(n< -32768||n>65535)throw Error('invalidRichDocument');text(String.fromCharCode(n&65535));fallback=state.uc;continue}
  if(word==='ansicpg'){codepage=n;continue}
  if(state.skip)continue
  if(word==='b'||word==='i'||word==='ul'||word==='strike')state[word]=n!==0
  else if(word==='ulnone')state.ul=false
  else if(word==='fs')state.fs=Math.min(144,Math.max(12,n))
  else if(word==='cf')state.cf=n
  else if(word==='plain'){state.b=state.i=state.ul=state.strike=false;state.fs=24;state.cf=0}
  else if(['ql','qc','qr','qj'].includes(word))state.align=({ql:'left',qc:'center',qr:'right',qj:'justify'} as any)[word]
  else if(word==='pard')state.align=''
  else if(word==='par')flushParagraph()
  else if(word==='line')line+='<br>'
  else if(word==='tab')text('\u00a0\u00a0\u00a0\u00a0')
  else if(word==='trowd'){if(line)flushParagraph();row=[];if(!table){result+='<table>';table=true}}
  else if(word==='cell'&&row){row.push(line);line=''}
  else if(word==='row'&&row){if(line)row.push(line);result+='<tr>'+row.map(c=>'<td>'+c+'</td>').join('')+'</tr>';line='';row=undefined}
  else if(['emdash','endash','bullet','lquote','rquote','ldblquote','rdblquote'].includes(word))text(({emdash:'—',endash:'–',bullet:'•',lquote:'‘',rquote:'’',ldblquote:'“',rdblquote:'”'} as any)[word])
  // Leave table mode when a regular paragraph starts.
  if(table&&!row&&word==='pard'){result+='</table>';table=false}
 }
 flush();if(stack.length||!closed)throw Error('invalidRichDocument');if(table){if(row){row.push(line);result+='<tr>'+row.map(c=>'<td>'+c+'</td>').join('')+'</tr>';line=''}result+='</table>'}flushParagraph()
 return result||'<p></p>'
}

// DOCX conversion reads only package XML and never fetches relationships or runs macros.
export function docxXmlToHtml(files:Record<string,Uint8Array>):string{
 function xml(name:string){const bytes=files[name];if(!bytes)return undefined;const source=new TextDecoder('utf-8',{fatal:true}).decode(bytes);if(/<!DOCTYPE|<!ENTITY/i.test(source))throw Error('invalidRichDocument');let bad=false;const doc=new DOMParser({onError:()=>{bad=true}}).parseFromString(source,'application/xml');if(bad)throw Error('invalidRichDocument');return doc.documentElement}
 const root=xml('word/document.xml');if(!root||tag(root)!=='document')throw Error('invalidRichDocument')
 const rels=new Map<string,string>();for(const r of elements(xml('word/_rels/document.xml.rels')||{})){const url=attr(r,'Target');if(attr(r,'Type').endsWith('/hyperlink')&&safeLink(url))rels.set(attr(r,'Id'),url)}
 const styles=new Map<string,any>();for(const s of elements(xml('word/styles.xml')||{}))styles.set(attr(s,'styleId'),s)
 const numbering=new Map<string,string>(),abstract=new Map<string,string>();const nums=xml('word/numbering.xml');for(const a of children(nums||{},'abstractNum'))abstract.set(attr(a,'abstractNumId'),attr(child(child(a,'lvl'),'numFmt'),'val'));for(const n of children(nums||{},'num'))numbering.set(attr(n,'numId'),abstract.get(attr(child(n,'abstractNumId'),'val'))||'bullet')
 const enabled=(node:any,name:string)=>{const p=child(node||{},name);return p&&!['0','false','off','none'].includes(attr(p,'val'))}
 function run(node:any):string{
  const p=child(node,'rPr');let body=''
  for(const c of elements(node)){switch(tag(c)){case 't':body+=escapeHtml(c.textContent||'');break;case 'tab':body+='&nbsp;&nbsp;&nbsp;&nbsp;';break;case 'br':case 'cr':body+='<br>';break;case 'drawing':case 'pict':case 'object':body+='[Image / object]';break}}
  if(enabled(p,'b'))body='<strong>'+body+'</strong>';if(enabled(p,'i'))body='<em>'+body+'</em>';if(enabled(p,'u'))body='<u>'+body+'</u>';if(enabled(p,'strike'))body='<s>'+body+'</s>'
  const color=attr(child(p||{},'color'),'val'),size=+attr(child(p||{},'sz'),'val');const style=(/^[a-f\d]{6}$/i.test(color)?`color:#${color};`:'')+(size>=12&&size<=144?`font-size:${size/2}pt;`:'')
  return style?`<span style="${style}">${body}</span>`:body
 }
 function render(node:any,depth=0):string{
  if(depth>100)throw Error('documentTooComplex')
  const name=tag(node),all=()=>elements(node).map(n=>render(n,depth+1)).join('')
  if(name==='r')return run(node)
  if(name==='hyperlink'){const body=all(),url=rels.get(attr(node,'id'));return url?`<a href="${escapeHtml(url)}">${body}</a>`:body}
  if(name==='p'){
   const p=child(node,'pPr'),id=attr(child(p||{},'pStyle'),'val'),style=styles.get(id),heading=/heading\s*([1-6])/i.exec(id+' '+attr(child(style||{},'name'),'val'))
   const name=heading?'h'+heading[1]:'p',align=attr(child(p||{},'jc'),'val'),body=all()||'<br>',num=child(p||{},'numPr')
   const result=`<${name}${['left','center','right','both'].includes(align)?` style="text-align:${align==='both'?'justify':align}"`:''}>${body}</${name}>`
   if(num){const kind=numbering.get(attr(child(num,'numId'),'val'))==='bullet'?'ul':'ol';return `<${kind}><li>${result}</li></${kind}>`}
   return result
  }
  if(name==='tbl')return '<table>'+children(node,'tr').map(n=>render(n,depth+1)).join('')+'</table>'
  if(name==='tr')return '<tr>'+children(node,'tc').map(n=>render(n,depth+1)).join('')+'</tr>'
  if(name==='tc'){const span=+attr(child(child(node,'tcPr')||{},'gridSpan'),'val');return `<td${span>1&&span<100?` colspan="${span}"`:''}>${all()}</td>`}
  if(name.endsWith('pr')||['instrText','del','sectPr'].includes(name))return ''
  return all()
 }
 const body=child(root,'body');if(!body)throw Error('invalidRichDocument');return render(body).replace(/<\/(ol|ul)><\1>/g,'')||'<p></p>'
}
