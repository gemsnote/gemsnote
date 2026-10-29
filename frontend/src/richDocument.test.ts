import {it,expect} from 'vitest'
import {strToU8,zipSync} from 'fflate'
import {htmlToRtf,rtfToHtml} from './richDocument'
import {parseMarkdownImport} from './markdownImport'
import {createImportJob,runImportJob} from './importJob'
it('round trips Unicode, styles, paragraphs, links and tables through real RTF',()=>{
 const rtf=htmlToRtf('<h1>标题</h1><p><b>中文😀</b> <i>italic</i><u>under</u><a href="https://example.com">link</a></p><ul><li>item</li></ul><table><tr><td>A</td><td>B</td></tr><tr><td>C</td><td>D</td></tr></table><p>after</p>','Title {\\}')
 expect(rtf).toMatch(/^\{\\rtf1/);expect(rtf).toContain('\\trowd');
 const html=rtfToHtml(strToU8(rtf));expect(html.replace(/<[^>]+>/g,'')).toContain('中文');expect(html).toContain('😀');expect(html).toContain('<strong>');expect(html).toContain('<em>italic</em>');expect(html).toContain('href="https://example.com"');expect(html).toContain('<td>A</td><td>B</td>');expect(html).toContain('</table><p>after</p>')
})
it('decodes code pages, Unicode fallback, and skips embedded objects',()=>{
 const html=rtfToHtml(strToU8(String.raw`{\rtf1\ansi\ansicpg936 \'d6\'d0\'ce\'c4\uc1\u20013?{\*\generator hidden;}\par {\object unsafe} text}`))
 expect(html).toContain('中文');expect(html).toContain('中');expect(html).not.toContain('hidden');expect(html).not.toContain('unsafe')
 for(const rtf of ['not rtf','{\\rtf1 missing','{\\rtf1 \\bin999 x}', '{\\rtf1 text} extra'])expect(()=>rtfToHtml(strToU8(rtf))).toThrow()
})
const word='<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><w:body><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>中文😀</w:t></w:r><w:hyperlink r:id="r1"><w:r><w:t>link</w:t></w:r></w:hyperlink></w:p><w:tbl><w:tr><w:tc><w:p><w:r><w:t>cell</w:t></w:r></w:p></w:tc></w:tr></w:tbl></w:body></w:document>'
it('imports DOCX without fetching external resources or trusting hyperlink protocols',()=>{
 const bytes=zipSync({'word/document.xml':strToU8(word),'word/_rels/document.xml.rels':strToU8('<Relationships><Relationship Id="r1" Type="x/hyperlink" Target="javascript:alert(1)" TargetMode="External"/></Relationships>')})
 const source=parseMarkdownImport('Word.docx',bytes);expect(source.notes[0].isMarkdown).toBe(false);expect(source.notes[0].content).toContain('<strong>中文😀</strong>');expect(source.notes[0].content).toContain('<td><p>cell</p></td>');expect(source.notes[0].content).not.toContain('javascript');expect(source.warnings?.length).toBe(1)
 expect(()=>parseMarkdownImport('bad.docx',zipSync({'word/document.xml':strToU8('<!DOCTYPE x [<!ENTITY x SYSTEM "file:///etc/passwd">]>'+word)}))).toThrow()
})
it('imports RTF and mixed ZIP records as rich notes and resumes with the right format',async()=>{
 const rtf=htmlToRtf('<p><b>bold</b></p>')
 const source=parseMarkdownImport('docs.zip',zipSync({'Folder/Rich.rtf':strToU8(rtf),'Plain.md':strToU8('# Markdown')}))
 expect(source.notes.map(n=>n.isMarkdown)).toEqual([false,true])
 let id=0;const job=createImportJob(source,'u','parent','Docs',()=>String(++id)),writes:any[]=[]
 await runImportJob(job,async(path,params)=>{if(path==='/web/bootstrap')return {User:{UserId:'u'}};if(path==='/web/save')writes.push(params)},()=>{})
 expect(writes.map(n=>n.isMarkdown)).toEqual([false,true])
})
