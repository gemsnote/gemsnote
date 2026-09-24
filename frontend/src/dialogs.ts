import { shallowRef } from 'vue'

type DialogRequest = {
  kind: 'alert' | 'confirm' | 'prompt'
  message: string
  initialValue: string
  resolve: (value: string | boolean | null) => void
}
export const activeDialog = shallowRef<DialogRequest | null>(null)
const queued: DialogRequest[] = []

function ask(kind: DialogRequest['kind'], message: string, initialValue = '') {
  return new Promise<string | boolean | null>(resolve => {
    const request = { kind, message, initialValue, resolve }
    if (activeDialog.value) queued.push(request)
    else activeDialog.value = request
  })
}
export async function confirmAction(message: string) { return await ask('confirm', message) === true }
export async function promptText(message: string, initialValue = '') { return await ask('prompt', message, initialValue) as string | null }
export async function showNotice(message: string) { await ask('alert', message) }
export function closeDialog(accepted: boolean, value = '') {
  const request = activeDialog.value
  if (!request) return
  activeDialog.value = queued.shift() || null
  request.resolve(request.kind === 'prompt' ? (accepted ? value : null) : accepted)
}
