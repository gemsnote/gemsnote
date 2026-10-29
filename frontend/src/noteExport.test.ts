import {describe,it,expect} from 'vitest'
import {unzipSync,strFromU8} from 'fflate'
import {hugoMarkdown,notebookZip,exportName} from './noteExport'
const note={NoteId:'n',Title:'标题: "引号"\n下一行',Tags:['yes','a: b','#标签'],CreatedTime:'2020-01-02T12:00:00+08:00',UpdatedTime:'2021-01-02T12:00:00Z',IsMarkdown:true}
describe('Hugo Markdown export',()=>{
 it('escapes metadata and preserves original Markdown and timestamps',()=>{
  const body='---\nbody front matter\n---\n\n# 原文\r\n'
  const result=hugoMarkdown(note,body,'笔记本')
  expect(result).toContain('title: "标题: \\"引号\\"\\n下一行"')
  expect(result).toContain('date: "2020-01-02T04:00:00.000Z"')
  expect(result).toContain('lastmod: "2021-01-02T12:00:00.000Z"')
  expect(result).toContain('tags: ["yes","a: b","#标签"]')
  expect(result).toContain('categories: ["笔记本"]')
  expect(result.endsWith(body)).toBe(true)
  expect(hugoMarkdown({...note,Title:'one\u2028two'},'', 'book')).toContain('title: "one\\u2028two"')
 })
 it('exports empty subdirectories, duplicate names and safe paths',async()=>{
  const root={NotebookId:'b',Title:'../根',Subs:[{NotebookId:'child',Title:'A'},{NotebookId:'empty',Title:'a'}]}
  const bytes=await notebookZip(root,async book=>book.NotebookId==='empty'?[]:[{Note:{...note,Title:'../same'},Content:'# test'},{Note:{...note,Title:'../same'},Content:'two'}])
  const files=unzipSync(bytes)
  expect(Object.keys(files)).toEqual(['_根/','_根/_same.md','_根/_same (2).md','_根/A/','_根/A/_same.md','_根/A/_same (2).md','_根/a (2)/'])
  expect(strFromU8(files['_根/A/_same.md'])).toContain('categories: ["A"]')
  expect(exportName('CON')).toBe('_CON')
  expect(new TextEncoder().encode(exportName('😀'.repeat(100))).length).toBeLessThanOrEqual(200)
 })
 it('fails instead of silently omitting unreadable notes or cycles',async()=>{
  await expect(notebookZip({NotebookId:'b',Title:'B'},async()=>{throw Error('offline')})).rejects.toThrow('offline')
  await expect(notebookZip({NotebookId:'b',Title:'B',Subs:[{NotebookId:'b',Title:'B'}]},async()=>[])).rejects.toThrow('invalidNotebookTree')
 })
})
