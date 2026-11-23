import {request} from './api'

// Browser <img> elements cannot send an Authorization header, so image reads
// use a short-lived, file-scoped capability issued for the signed-in user.
// Stored note content keeps plain /api2/file/getImage?fileId=... URLs; the
// signature is only added while rendering and is stripped again before saving.
const imagePath = '/api2/file/getImage'
const fallbackLifetime = 10 * 60 * 1000
const imagePattern = /\/api2\/file\/getImage\?fileId=([0-9a-fA-F]{24})(?:&(?:amp;)?(?:uid|token)=[^\s"'<>]*)*/g

const cache = new Map<string, {url: string; expires: number}>()
const pending = new Map<string, Promise<string>>()

export function fileIdFromImageSrc(src: string): string {
  try {
    const url = new URL(src, 'http://localhost')
    if (url.pathname !== imagePath) return ''
    return url.searchParams.get('fileId') || ''
  } catch {
    return ''
  }
}

export function stripImageSignatures(content: string): string {
  if (!content.includes(imagePath)) return content
  return content.replace(imagePattern, `${imagePath}?fileId=$1`)
}

async function capability(fileId: string): Promise<string> {
  const cached = cache.get(fileId)
  if (cached && cached.expires - Date.now() > 30_000) return cached.url
  const inflight = pending.get(fileId)
  if (inflight) return inflight
  const task = (async () => {
    try {
      const data: any = await request('/api2/file/imageToken', {fileId})
      if (data?.Ok && typeof data?.Item?.Url === 'string') {
        const seconds = Number(data.Item.ExpiresIn)
        const lifetime = (Number.isFinite(seconds) && seconds > 0 ? seconds : fallbackLifetime / 1000) * 1000
        cache.set(fileId, {url: data.Item.Url, expires: Date.now() + lifetime})
        return data.Item.Url as string
      }
    } catch {
      // Fall through: the plain URL may still work for token-authenticated clients.
    }
    return ''
  })().finally(() => pending.delete(fileId))
  pending.set(fileId, task)
  return task
}

export async function signImageSrc(src: string): Promise<string> {
  const fileId = fileIdFromImageSrc(src)
  if (!fileId) return src
  const signed = await capability(fileId)
  return signed || src
}

export async function signImageRefs(html: string): Promise<string> {
  if (!html.includes(imagePath)) return html
  const ids = Array.from(new Set(Array.from(html.matchAll(imagePattern), match => match[1])))
  if (!ids.length) return html
  const signed = new Map<string, string>()
  await Promise.all(ids.map(async id => {
    const url = await capability(id)
    if (url) signed.set(id, url)
  }))
  if (!signed.size) return html
  return html.replace(imagePattern, (raw, fileId: string) => signed.get(fileId) || raw)
}
