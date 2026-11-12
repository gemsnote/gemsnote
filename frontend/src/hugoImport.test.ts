import {describe,it,expect} from 'vitest'
import {zipSync,strToU8} from 'fflate'
import {parseMarkdownImport} from './markdownImport'
import {createImportJob,runImportJob} from './importJob'

const created='2020-01-02T03:04:05+08:00',updated='2021-02-03T04:05:06Z'
const headers=[
 `---\ntitle: "元数据标题"\ndate: ${created}\nlastmod: ${updated}\ntags:\n  - Go\n  - 中文\nparams:\n  custom: ignored\n---\n`,
 `+++\ntitle = "元数据标题"\ndate = ${created}\nlastmod = ${updated}\ntags = ["Go", "中文"]\n[params]\ncustom = "ignored"\n+++\n`,
 JSON.stringify({title:'元数据标题',date:created,lastmod:updated,tags:['Go','中文'],params:{custom:'{ignored}'}})+'\n'
]
describe('Hugo imports',()=>{
 it.each(headers)('uses front matter in single MD and ZIP, stripping all metadata',async header=>{
  const text=header+'\n# 正文\n\n---\nkeep body\n'
  for(const source of [parseMarkdownImport('filename.md',strToU8(text)),parseMarkdownImport('archive.zip',zipSync({'Dir/filename.md':strToU8(text),'plain.md':strToU8('plain')}))]){
   expect(source.notes[0]).toMatchObject({title:'元数据标题',createdTime:'2020-01-01T19:04:05.000Z',updatedTime:'2021-02-03T04:05:06.000Z',tags:['Go','中文'],content:'# 正文\n\n---\nkeep body\n'})
   let id=0;const writes:any[]=[]
   await runImportJob(createImportJob(source,'user','target','unused',()=>String(++id)),async(path,params)=>{if(path==='/web/bootstrap')return {User:{UserId:'user'}};writes.push({path,...params});return {}},()=>{})
   expect(writes.find(w=>w.path==='/web/save'&&w.title==='元数据标题')).toMatchObject({tags:'Go,中文',createdTime:'2020-01-01T19:04:05.000Z',updatedTime:'2021-02-03T04:05:06.000Z'})
   expect(writes.filter(w=>w.path==='/notebook/addNotebook').map(w=>w.title)).not.toContain('archive')
  }
 })
 it('imports a flat ZIP directly into the selected notebook',()=>{
  const job=createImportJob(parseMarkdownImport('archive.zip',zipSync({'a.md':strToU8('body')})),'u','target','ignored',()=> 'id')
  expect(job.steps).toHaveLength(1);expect(job.steps[0]).toMatchObject({kind:'note',parent:'target',title:'a',tags:[]})
 })
 it('keeps plain Markdown and horizontal rules verbatim',()=>{
  expect(parseMarkdownImport('filename.md',strToU8('---\n# body\n')).notes[0]).toMatchObject({title:'filename',content:'---\n# body\n',tags:[]})
 })
 it('does not fabricate a missing metadata title or replace invalid dates',()=>{
  expect(parseMarkdownImport('file.md',strToU8('---\ndate: 2020-01-02\n---\nbody')).notes[0]).toMatchObject({title:'',createdTime:'2020-01-02T00:00:00.000Z',updatedTime:'2020-01-02T00:00:00.000Z'})
  expect(()=>parseMarkdownImport('file.md',strToU8('---\ntitle: abc\ndate: invalid\n---\nbody'))).toThrow('Hugo 日期无效')
 })
})
