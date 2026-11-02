import {Marked} from 'marked'
import DOMPurify from 'dompurify'

export function previewDocument(content: string, markdown: boolean, prefix: string) {
  const diagrams: {id: string; source: string}[] = []
  const parser = new Marked()
  parser.use({renderer: {code(token) {
    if (token.lang?.trim().split(/\s+/)[0].toLowerCase() !== 'mermaid') return false
    const id = `${prefix}-${diagrams.length}`
    diagrams.push({id, source: token.text})
    return `<div id="${id}" class="mermaid-slot"></div>`
  }}})
  return {
    html: DOMPurify.sanitize(markdown ? String(parser.parse(content, {async: false})) : content),
    diagrams,
  }
}

export type DiagramImage = {svg: string; width: number; height: number}
export class MermaidRenderer {
  private frame?: HTMLIFrameElement
  private ready?: Promise<void>
  private readyResolve?: () => void
  private readyReject?: (reason: Error) => void
  private sequence = 0
  private requested = 0
  private pending = new Map<number, {resolve: (value: DiagramImage) => void; reject: (error: Error) => void; timer: ReturnType<typeof setTimeout>}>()
  private readyTimer?: ReturnType<typeof setTimeout>
  private disposed = false
  constructor(private container: HTMLElement) { window.addEventListener('message', this.onMessage) }
  private onMessage = (event: MessageEvent) => {
    if (event.source !== this.frame?.contentWindow || event.data?.channel !== 'gemsnote-mermaid') return
    if (event.data.ready) { clearTimeout(this.readyTimer); this.readyResolve?.(); return }
    const request = this.pending.get(event.data.id)
    if (!request) return
    clearTimeout(request.timer)
    this.pending.delete(event.data.id)
    if (typeof event.data.svg === 'string' && event.data.width > 0 && event.data.height > 0) request.resolve(event.data)
    else request.reject(new Error('renderFailed'))
  }
  private start() {
    if (this.ready) return this.ready
    this.ready = new Promise<void>((resolve, reject) => { this.readyResolve = resolve; this.readyReject = reject })
    const frame = document.createElement('iframe')
    frame.title = 'Mermaid renderer'
    frame.setAttribute('sandbox', 'allow-scripts')
    frame.setAttribute('aria-hidden', 'true')
    frame.tabIndex = -1
    frame.style.cssText = 'position:absolute;width:1px;height:1px;border:0;visibility:hidden;pointer-events:none'
    frame.src = `${import.meta.env.BASE_URL}assets/mermaid/index.html`
    this.frame = frame
    this.container.append(frame)
    this.readyTimer = setTimeout(() => this.dispose(), 20000)
    return this.ready
  }
  async render(source: string): Promise<DiagramImage> {
    if (this.disposed) throw new Error('disposed')
    if (++this.requested > 32) throw new Error('limit')
    if (!source.trim() || source.length > 20000 || source.split('\n').length > 500) throw new Error('limit')
    await this.start()
    if (this.disposed) throw new Error('disposed')
    return new Promise((resolve, reject) => {
      const id = ++this.sequence
      const timer = setTimeout(() => { this.pending.delete(id); reject(new Error('timeout')) }, 30000)
      this.pending.set(id, {resolve, reject, timer})
      this.frame!.contentWindow!.postMessage({channel: 'gemsnote-mermaid', id, source}, '*')
    })
  }
  dispose() {
    if (this.disposed) return
    this.disposed = true
    clearTimeout(this.readyTimer)
    this.readyReject?.(new Error('disposed'))
    for (const request of this.pending.values()) { clearTimeout(request.timer); request.reject(new Error('disposed')) }
    this.pending.clear()
    window.removeEventListener('message', this.onMessage)
    this.frame?.remove()
  }
}
