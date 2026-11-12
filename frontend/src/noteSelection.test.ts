import {describe,it,expect} from 'vitest'
import {selectNote} from './noteSelection'
const empty=()=>({ids:new Set<string>(),anchor:''})
describe('note selection',()=>{
 it('plain click replaces selection and sets the range anchor',()=>{
  const first=selectNote(empty(),['a','b','c'],'a')
  expect([...selectNote(first,['a','b','c'],'c').ids]).toEqual(['c'])
 })
 it('Ctrl toggles individual notes without losing other selections',()=>{
  const first=selectNote(empty(),['a','b','c'],'a')
  const second=selectNote(first,['a','b','c'],'c',{toggle:true})
  expect([...second.ids]).toEqual(['a','c'])
  expect([...selectNote(second,['a','b','c'],'a',{toggle:true}).ids]).toEqual(['c'])
 })
 it('Shift follows the current displayed order, including reverse ranges',()=>{
  const first=selectNote(empty(),['d','b','c','a'],'c')
  expect([...selectNote(first,['d','b','c','a'],'d',{range:true}).ids]).toEqual(['d','b','c'])
  expect([...selectNote(first,['a','c','b','d'],'d',{range:true}).ids]).toEqual(['c','b','d'])
 })
 it('Ctrl+Shift adds a range to an existing discontiguous selection',()=>{
  const result=selectNote({ids:new Set(['a','e']),anchor:'b'},['a','b','c','d','e'],'d',{range:true,toggle:true})
  expect([...result.ids].sort()).toEqual(['a','b','c','d','e'])
 })
 it('removed notes cannot remain selected and a missing anchor starts at the first row',()=>{
  const result=selectNote({ids:new Set(['gone']),anchor:'gone'},['a','b','c'],'b',{range:true})
  expect([...result.ids]).toEqual(['a','b'])
  expect(result.anchor).toBe('a')
 })
})
