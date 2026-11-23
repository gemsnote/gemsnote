import {describe, expect, it} from 'vitest'
import {fileIdFromImageSrc, stripImageSignatures} from './signedImages'

describe('short-lived image capabilities', () => {
  it('finds file ids in stored and signed image URLs', () => {
    expect(fileIdFromImageSrc('/api2/file/getImage?fileId=0123456789abcdef01234567')).toBe('0123456789abcdef01234567')
    expect(fileIdFromImageSrc('/api2/file/getImage?fileId=0123456789abcdef01234567&uid=507f1f77bcf86cd799439011&token=1700000000.abcdef')).toBe('0123456789abcdef01234567')
    expect(fileIdFromImageSrc('/file/outputImage?fileId=0123456789abcdef01234567')).toBe('')
    expect(fileIdFromImageSrc('https://example.test/logo.png')).toBe('')
    expect(fileIdFromImageSrc('')).toBe('')
  })

  it('strips capabilities before note content is saved', () => {
    const stored = '<p><img src="/api2/file/getImage?fileId=0123456789abcdef01234567"></p>'
    const rendered = stored.replace('?fileId=', '?fileId=')
      .replace('0123456789abcdef01234567"', '0123456789abcdef01234567&uid=507f1f77bcf86cd799439011&token=1700000000.abcdef"')
    expect(stripImageSignatures(rendered)).toBe(stored)
    expect(stripImageSignatures(stored)).toBe(stored)
    expect(stripImageSignatures('<p><img src="/public/upload/logo.png"></p>')).toBe('<p><img src="/public/upload/logo.png"></p>')
  })

  it('leaves unrelated images untouched', () => {
    const html = '<img src="/file/outputImage?fileId=0123456789abcdef01234567"><img src="https://example.test/a.png">'
    expect(stripImageSignatures(html)).toBe(html)
  })
})
