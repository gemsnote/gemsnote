import { chromium } from '@playwright/test'
import { existsSync } from 'node:fs'

const baseURL = process.env.GEMSNOTE_E2E_URL
const username = process.env.GEMSNOTE_E2E_USERNAME
const password = process.env.GEMSNOTE_E2E_PASSWORD
if (!baseURL || !username || !password) {
  throw new Error('GEMSNOTE_E2E_URL, GEMSNOTE_E2E_USERNAME and GEMSNOTE_E2E_PASSWORD are required')
}

const executablePath = [process.env.GEMSNOTE_CHROMIUM_PATH, '/usr/bin/chromium', '/usr/bin/chromium-browser', '/usr/bin/google-chrome'].find(path => path && existsSync(path))
const browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
const page = await browser.newPage()
await page.addInitScript(() => localStorage.setItem('gemsnote:language', 'zh-CN'))
page.setDefaultTimeout(15_000)
try {
  await page.goto(`${baseURL}/login`)
  await page.getByLabel('邮箱或用户名').fill(username)
  await page.getByLabel('密码').fill(password)
  await page.getByRole('button', { name: '继续' }).click()
  await page.waitForURL(/\/note/)

  await page.getByLabel('用户菜单').click()
  await page.getByRole('menuitem', { name: '语言' }).click()
  await page.getByRole('menuitemradio', { name: 'English' }).click()
  await page.getByRole('button', { name: /All notes/ }).waitFor()
  await page.getByLabel('User menu').click()
  await page.getByRole('menuitem', { name: 'Language' }).click()
  await page.getByRole('menuitemradio', { name: '简体中文' }).click()
  await page.getByRole('button', { name: /所有笔记/ }).waitFor()

  async function verifyKeyboardResize(selector, separatorName) {
    const panel = page.locator(selector)
    const originalWidth = await panel.evaluate(element => element.getBoundingClientRect().width)
    const direction = originalWidth >= 520 ? 'ArrowLeft' : 'ArrowRight'
    await page.getByRole('separator', { name: separatorName }).press(direction)
    const resizedWidth = await panel.evaluate(element => element.getBoundingClientRect().width)
    if (direction === 'ArrowRight' ? resizedWidth <= originalWidth : resizedWidth >= originalWidth) {
      throw new Error(`${separatorName} keyboard resize did not move ${direction}`)
    }
  }

  await verifyKeyboardResize('.notebooks', '调节我的空间宽度')
  const notebookMenu = page.locator('.notebook-more').first()
  if (await notebookMenu.count()) {
    await notebookMenu.locator('..').hover()
    await notebookMenu.click()
    await page.getByRole('menuitem', { name: '重命名' }).waitFor()
    await page.keyboard.press('Escape')
  }

  await page.getByLabel('排序方式').click()
  await page.getByRole('menuitemradio', { name: '标题：升序' }).click()
  await page.getByLabel('排序方式').click()
  await page.getByRole('menuitemradio', { name: '修改时间：降序' }).click()
  await verifyKeyboardResize('.note-list', '调节笔记栏宽度')
  await page.getByLabel('隐藏笔记栏').click()
  await page.getByRole('button', { name: '所有笔记' }).click()
  await page.getByLabel('隐藏笔记栏').waitFor()
  await page.getByLabel('隐藏我的空间').click()
  await page.getByLabel('展开我的空间').click()
  await page.getByLabel('搜索笔记本').waitFor()
  await page.getByLabel('隐藏我的空间').click()
  await page.setViewportSize({ width: 800, height: 700 })
  await page.getByLabel('展开我的空间').click()
  await page.getByLabel('搜索笔记本').waitFor()
  await page.setViewportSize({ width: 1280, height: 800 })
  await page.getByLabel('新建富文本笔记').waitFor()
  await page.getByLabel('新建 Markdown 笔记').waitFor()

  async function waitForAutoSave(content) {
    const response = await page.waitForResponse(result =>
      result.url().includes('/api2/save') && (result.request().postData() || '').includes(content)
    )
    const saved = await response.json()
    if (!response.ok() || !saved.Note) throw new Error(`automatic save failed: ${JSON.stringify(saved)}`)
  }

  const noteTitle = `Vue 端到端验收 ${Date.now()}`
  await page.getByTitle('新建 Markdown 笔记').click()
  const firstSave = waitForAutoSave('Vue Web UI 正常保存。')
  await page.getByLabel('笔记标题').fill(noteTitle)
  await page.getByLabel('Markdown 正文').fill('# Gemsnote\n\nVue Web UI 正常保存。')
  await firstSave

  await page.setViewportSize({ width: 640, height: 700 })
  await page.getByRole('button', { name: '所有笔记' }).click()
  await page.locator('.note-list').waitFor()
  await page.getByText(noteTitle, { exact: true }).click()
  await page.getByLabel('笔记标题').waitFor()
  await page.setViewportSize({ width: 1280, height: 800 })

  await page.locator('.note-info-line').getByText('格式：MD').waitFor()
  await page.locator('.note-tag-line').getByLabel('笔记标签').waitFor()

  const attachmentName = `gemsnote-e2e-${Date.now()}.txt`
  await page.locator('input[type="file"]').nth(1).setInputFiles({
    name: attachmentName,
    mimeType: 'text/plain',
    buffer: Buffer.from('Gemsnote attachment E2E'),
  })
  await page.getByRole('heading', { name: '附件' }).waitFor()
  await page.getByRole('link', { name: attachmentName }).waitFor()
  const deleteResponse = page.waitForResponse(response => response.url().includes('/api2/attachments/delete'))
  page.once('dialog', dialog => dialog.accept())
  await page.getByRole('listitem').filter({ hasText: attachmentName }).getByRole('button', { name: '删除' }).click()
  const deleteResult = await (await deleteResponse).json()
  if (!deleteResult.Ok) throw new Error(`attachment delete failed: ${deleteResult.Msg}`)
  await page.getByRole('link', { name: attachmentName }).waitFor({ state: 'detached' })
  await page.getByRole('button', { name: '关闭' }).click()

  // Attachment mutations increment the note USN; saving afterwards verifies
  // that the UI refreshed its optimistic-concurrency metadata.
  const secondSave = waitForAutoSave('附件后继续保存正常。')
  await page.getByLabel('Markdown 正文').fill('# Gemsnote\n\n附件后继续保存正常。')
  await secondSave

  await page.getByLabel('用户菜单').click()
  await page.getByRole('menuitem', { name: '账号' }).click()
  await page.getByRole('heading', { name: '账号管理' }).waitFor()
  await page.getByRole('link', { name: '返回' }).click()
  await page.getByLabel('用户菜单').click()
  await page.getByRole('menuitem', { name: '管理' }).click()
  await page.getByRole('heading', { name: '系统管理' }).waitFor()

  const removed = await page.request.get(`${baseURL}/blog`)
  if (removed.status() !== 410) throw new Error(`blog returned ${removed.status()}, expected 410`)
  console.log('Gemsnote Vue Web E2E passed')
} finally {
  await browser.close()
}
