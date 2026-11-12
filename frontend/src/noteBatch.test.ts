import {describe,it,expect} from 'vitest'
import {runNoteBatch,type BatchAPI} from './noteBatch'
const notes=[{NoteId:'a',UserId:'me',Title:'A'},{NoteId:'b',UserId:'me',Title:'B'}]
describe('batch note actions',()=>{
 it('moves owned notes in one API2-compatible array request',async()=>{
  const calls:any[]=[],done=new Set<string>()
  await runNoteBatch('move',notes,'me',async(path,params)=>{calls.push([path,params]);return true},done,{destination:'target'})
  expect(calls).toEqual([['/note/moveNote',{noteIds:['a','b'],notebookId:'target'}]])
  expect([...done]).toEqual(['a','b'])
 })
 it('moves received notes through personal placement, never editing the source',async()=>{
  const calls:any[]=[]
  await runNoteBatch('move',[...notes,{NoteId:'received',UserId:'author',IsShared:true}],'me',async(path,params)=>{calls.push([path,params]);return true},new Set(),{destination:'target'})
  expect(calls[1]).toEqual(['/shared/move',{noteId:'received',notebookId:'target'}])
 })
 it.each(['share','delete'] as const)('rejects foreign/shared %s before any writes',async action=>{
  let calls=0
  await expect(runNoteBatch(action,[...notes,{NoteId:'received',UserId:'author',IsShared:true}],'me',async()=>{calls++},new Set(),{recipient:'reader'})).rejects.toThrow('只能共享或删除自己的笔记')
  expect(calls).toBe(0)
 })
 it('shares every selected note read-only and retries only failures',async()=>{
  const calls:string[]=[],done=new Set<string>();let fail=true
  const api:BatchAPI=async(path,params)=>{
   expect(path).toBe('/share/addShareNote');expect(params.emails).toEqual(['reader']);expect(params.perm).toBe(0)
   calls.push(String(params.noteId))
   return {reader:{Ok:params.noteId==='b'&&fail?false:true,Msg:'offline'}}
  }
  const failures=await runNoteBatch('share',notes,'me',api,done,{recipient:' reader '})
  expect(failures.map(f=>f.ids)).toEqual([['b']]);expect([...done]).toEqual(['a'])
  fail=false
  expect(await runNoteBatch('share',notes,'me',api,done,{recipient:'reader'})).toEqual([])
  expect(calls).toEqual(['a','b','b'])
 })
 it('does not acknowledge an empty share response as success',async()=>{
  const done=new Set<string>()
  const failures=await runNoteBatch('share',notes,'me',async()=>({}),done,{recipient:'reader'})
  expect(done.size).toBe(0);expect(failures.length).toBe(2)
 })
 it('trashes active notes together and permanently deletes trash only on explicit request',async()=>{
  const calls:any[]=[],api:BatchAPI=async(path,params)=>{calls.push([path,params]);return true}
  await runNoteBatch('delete',notes,'me',api,new Set(),{})
  expect(calls).toEqual([['/note/deleteNote',{noteIds:['a','b'],isShared:false}]])
  calls.length=0
  await runNoteBatch('delete',notes,'me',api,new Set(),{trash:true})
  expect(calls.map(c=>c[0])).toEqual(['/note/deleteTrash','/note/deleteTrash'])
 })
 it('stops starting jobs when the view is disposed',async()=>{
  const controller=new AbortController(),calls:string[]=[]
  await runNoteBatch('share',notes,'me',async(_,params)=>{calls.push(String(params.noteId));controller.abort();return {reader:{Ok:true}}},new Set(),{recipient:'reader'},()=>{},controller.signal)
  expect(calls).toEqual(['a'])
 })
})
