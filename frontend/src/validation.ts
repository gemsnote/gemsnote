import { watch } from 'vue'
import { language, t } from './i18n'

// WebView validation messages otherwise follow the OS/browser locale, not
// the language selected in Gemsnote. Retain native constraint checking.
export function installLocalizedValidation() {
  function update(target: EventTarget | null) {
    if (!(target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement || target instanceof HTMLSelectElement)) return
    target.setCustomValidity('')
    const validity = target.validity
    let message = ''
    if (validity.valueMissing) message = t('请填写此字段。')
    else if (validity.typeMismatch && target instanceof HTMLInputElement) message = t(target.type === 'email' ? '请输入有效的邮箱地址。' : '请输入有效的网址。')
    else if (validity.rangeUnderflow) message = t('请输入不小于 {min} 的值。', { min: (target as HTMLInputElement).min })
    else if (validity.rangeOverflow) message = t('请输入不大于 {max} 的值。', { max: (target as HTMLInputElement).max })
    else if (validity.tooShort) message = t('请至少输入 {min} 个字符。', { min: (target as HTMLInputElement).minLength })
    else if (validity.tooLong) message = t('请最多输入 {max} 个字符。', { max: (target as HTMLInputElement).maxLength })
    else if (!validity.valid) message = t('请输入有效值。')
    target.setCustomValidity(message)
  }
  const listener = (event: Event) => update(event.target)
  for (const event of ['invalid', 'input', 'change']) document.addEventListener(event, listener, true)
  const stop = watch(language, () => document.querySelectorAll('input,textarea,select').forEach(update))
  return () => {
    stop()
    for (const event of ['invalid', 'input', 'change']) document.removeEventListener(event, listener, true)
  }
}
