// Runs in a sandboxed iframe / isolated local WebView. No account or API access.
import mermaid from 'mermaid'
import DOMPurify from 'dompurify'

const options = {
  startOnLoad: false, securityLevel: 'strict' as const, theme: 'default' as const,
  htmlLabels: false, suppressErrorRendering: true, maxTextSize: 20000, maxEdges: 200,
  fontFamily: 'Arial, sans-serif',
  flowchart: {htmlLabels: false},
}
const host = document.getElementById('graph')!
let queue = Promise.resolve()
let sequence = 0
function reply(message: object) {
  const result = {channel: 'gemsnote-mermaid', ...message}
  const bridge = (window as any).MermaidResult
  if (bridge) bridge.postMessage(JSON.stringify(result))
  else window.parent.postMessage(result, '*')
}
async function draw(request: {id: number; source: string; format?: string}) {
  const {id, source} = request
  try {
    // Note-level directives cannot alter security, CSS, fonts or resource URLs.
    if (!source.trim() || source.length > 20000 || source.split('\n').length > 500 ||
        /%%\s*\{|^\s*---(?:\r?\n|$)/.test(source)) throw new Error('unsupported')
    mermaid.initialize(options)
    const {svg} = await mermaid.render(`diagram-${++sequence}`, source, host)
    const clean = DOMPurify.sanitize(svg, {
      USE_PROFILES: {svg: true, svgFilters: true},
      FORBID_TAGS: ['foreignObject', 'image', 'a', 'script'],
    })
    host.innerHTML = clean
    const element = host.querySelector('svg')!
    const box = element.viewBox.baseVal
    const width = box.width || element.getBoundingClientRect().width
    const height = box.height || element.getBoundingClientRect().height
    if (!(width > 0 && height > 0) || width > 20000 || height > 20000) throw new Error('size')
    element.setAttribute('width', String(width))
    element.setAttribute('height', String(height))
    element.style.maxWidth = 'none'
    const output = new XMLSerializer().serializeToString(element)
    if (request.format === 'png') {
      // Raster cache on Flutter: one WebView per document, not one per diagram.
      const scale = Math.min(2, 2048 / width, 2048 / height, Math.sqrt(2000000 / (width * height)))
      const canvas = document.createElement('canvas')
      canvas.width = Math.max(1, Math.ceil(width * scale))
      canvas.height = Math.max(1, Math.ceil(height * scale))
      const image = new Image()
      image.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(output)
      await image.decode()
      const context = canvas.getContext('2d')!
      context.fillStyle = '#ffffff'
      context.fillRect(0, 0, canvas.width, canvas.height)
      context.drawImage(image, 0, 0, canvas.width, canvas.height)
      reply({id, width, height, png: canvas.toDataURL('image/png')})
    } else reply({id, width, height, svg: output})
  } catch {
    // Never echo Mermaid exceptions containing user-controlled HTML.
    reply({id, error: 'renderFailed'})
  } finally { host.replaceChildren() }
}
function receive(data: any) {
  if (data?.channel !== 'gemsnote-mermaid' || !Number.isSafeInteger(data.id) ||
      typeof data.source !== 'string') return
  queue = queue.then(() => draw(data))
}
window.addEventListener('message', event => { if (event.source === window.parent) receive(event.data) })
;(window as any).renderGemsnoteMermaid = receive
reply({ready: true})
