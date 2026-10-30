import {describe,it,expect} from 'vitest'
import {generateInitialPassword} from './password'

describe('initial password generation',()=>{
  it('generates distinct 12-character passwords with all four character classes',()=>{
    const passwords=Array.from({length:100},()=>generateInitialPassword())
    expect(new Set(passwords).size).toBe(100)
    for(const password of passwords){
      expect(password).toHaveLength(12)
      for(const pattern of [/[A-Z]/,/[a-z]/,/[0-9]/,/[!@#$%^&*_=+?\-]/])expect(password).toMatch(pattern)
    }
  })
})
