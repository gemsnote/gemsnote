import { chromium, expect } from '@playwright/test'
import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { preview } from 'vite'

// Exercise the shared frontend without window.go bindings or a real server.
const server = await preview({ root: fileURLToPath(new URL('..', import.meta.url)), configFile: false, preview: { host: '127.0.0.1', port: 0 } })
const baseURL = `http://127.0.0.1:${server.httpServer.address().port}`
const executablePath = [process.env.GEMSNOTE_CHROMIUM_PATH, '/usr/bin/chromium', '/usr/bin/chromium-browser', '/usr/bin/google-chrome'].find(path => path && existsSync(path))
let browser
try {
  browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
  for (const language of ['zh-CN', 'en-US']) {
    const page = await browser.newPage({ locale: 'en-US', viewport: { width: 1280, height: 800 } })
    const chinese = language === 'zh-CN'
    let aboutFails = true, uploadCount = 0, createdCount = 0, deletedCount = 0
    const nativeDialogs = [], pageErrors = []
    page.on('dialog', async dialog => { nativeDialogs.push(dialog.type()); await dialog.dismiss() })
    page.on('pageerror', error => pageErrors.push(error.message))
    await page.addInitScript(language => localStorage.setItem('gemsnote:language', language), language)
    await page.route('**/api2/**', async route => {
      const path = new URL(route.request().url()).pathname
      let data
      switch (path) {
        case '/api2/web/downloadStatus': data = { Running: false }; break
        case '/api2/bootstrap': data = { Desktop: true, Host: 'https://notes.example.test:2443', User: { UserId: 'user1', Username: 'tester', Email: 'tester@example.test' }, Notebooks: [], Tags: [] }; break
        case '/api2/notes': data = []; break
        case '/api2/desktop/about':
          data = aboutFails ? { Ok: false, Msg: 'info unavailable' } : { Name: 'Gemsnote', Version: '9.8.7-test', Platform: 'linux', Arch: 'amd64', Runtime: 'go1.23.0' }
          break
        case '/api2/groups': data = deletedCount ? [] : [{ GroupId: 'g1', UserId: 'user1', Title: 'test group', Users: [] }]; break
        case '/api2/avatar':
          expect(route.request().headers()['content-type']).toContain('multipart/form-data')
          expect(route.request().postDataBuffer().toString()).toContain('avatar.png')
          uploadCount++
          data = { Ok: true, Id: '507f1f77bcf86cd799439011' }
          break
        case '/api2/file/getImage': await route.fulfill({ status: 204 }); return
        case '/api2/member/group/deleteGroup': deletedCount++; data = { Ok: true }; break
        case '/api2/notebook/addNotebook':
          expect(new URLSearchParams(route.request().postData()).get('title')).toBe('New book')
          createdCount++; data = { Ok: true }; break
        default: throw new Error(`Unexpected API request: ${path}`)
      }
      await route.fulfill({ json: data })
    })
    await page.goto(`${baseURL}/note`)
    // The desktop's initial 1280x800 window must show all three columns without
    // overlap or horizontal overflow, rather than starting in compact mode.
    const columns = []
    for (const selector of ['.notebooks', '.note-list', '.editor']) {
      const column = page.locator(selector)
      await expect(column).toBeVisible()
      const bounds = await column.boundingBox()
      expect(bounds.width).toBeGreaterThan(0)
      expect(bounds.x).toBeGreaterThanOrEqual(0)
      expect(bounds.x + bounds.width).toBeLessThanOrEqual(1280)
      columns.push(bounds)
    }
    expect(columns[0].x + columns[0].width).toBeLessThanOrEqual(columns[1].x)
    expect(columns[1].x + columns[1].width).toBeLessThanOrEqual(columns[2].x)
    await page.getByRole('button', { name: chinese ? '用户菜单' : 'User menu', exact: true }).click()
    await page.getByRole('menuitem', { name: chinese ? '关于' : 'About', exact: true }).click()
    const about = page.getByRole('dialog', { name: chinese ? '关于珠玑笔记' : 'About Gemsnote', exact: true })
    await expect(about.getByRole('alert')).toContainText(chinese ? '无法读取应用信息' : 'Could not load application information')
    aboutFails = false
    await about.getByRole('button', { name: chinese ? '重试' : 'Retry' }).click()
    for (const value of ['9.8.7-test', 'linux', 'amd64', 'go1.23.0']) await expect(about).toContainText(value)
    await about.getByRole('button', { name: chinese ? '关闭' : 'Close' }).click()

    // Prompt cancel/Escape must not create anything; submit must preserve input.
    const newBook = page.getByRole('button', { name: chinese ? '新建笔记本' : 'New notebook', exact: true })
    await newBook.click()
    const dialog = page.getByRole('dialog', { name: chinese ? '珠玑笔记' : 'Gemsnote', exact: true })
    await expect(dialog).toBeVisible()
    await dialog.getByRole('button', { name: chinese ? '取消' : 'Cancel' }).click()
    expect(createdCount).toBe(0)
    await newBook.click()
    await page.keyboard.press('Escape')
    await expect(dialog).not.toBeVisible()
    expect(createdCount).toBe(0)
    await newBook.click()
    await dialog.getByRole('textbox').fill('New book')
    await dialog.getByRole('button', { name: chinese ? '确定' : 'Confirm', exact: true }).click()
    await expect.poll(() => createdCount).toBe(1)

    await page.getByRole('button', { name: chinese ? '用户菜单' : 'User menu', exact: true }).click()
    await page.getByRole('menuitem', { name: chinese ? '账号' : 'Account', exact: true }).click()
    const serverAddress = page.getByLabel(chinese ? '服务器地址' : 'Server address', { exact: true })
    await expect(serverAddress).toHaveValue('https://notes.example.test:2443')
    await expect(serverAddress).toHaveAttribute('readonly', '')
    const picker = page.getByRole('group', { name: chinese ? '更换头像' : 'Change avatar' })
    await expect(picker).toContainText(chinese ? '未选择文件' : 'No file selected')
    const fileChooserPromise = page.waitForEvent('filechooser')
    await picker.getByRole('button', { name: chinese ? '选择文件' : 'Choose file' }).click()
    const fileChooser = await fileChooserPromise
    expect(await fileChooser.element().getAttribute('accept')).toBe('image/*')
    const avatar = { name: 'avatar.png', mimeType: 'image/png', buffer: Buffer.from('test avatar bytes') }
    await fileChooser.setFiles(avatar)
    await expect.poll(() => uploadCount).toBe(1)
    await expect(picker.getByRole('status')).toHaveText('avatar.png')
    // Selecting the same file again is a valid retry.
    await picker.locator('input[type=file]').setInputFiles(avatar)
    await expect.poll(() => uploadCount).toBe(2)

    const email = page.getByLabel(chinese ? '新邮箱' : 'New email', { exact: true })
    await email.fill('not-an-email')
    await page.getByRole('button', { name: chinese ? '发送邮箱验证邮件' : 'Send verification email', exact: true }).click()
    expect(await email.evaluate(input => input.validationMessage)).toBe(chinese ? '请输入有效的邮箱地址。' : 'Please enter a valid email address.')
    await email.fill('valid@example.test')
    expect(await email.evaluate(input => input.validationMessage)).toBe('')

    // Confirmation labels and title must be app-localized, with no WebView dialog.
    const deleteGroup = page.getByRole('button', { name: chinese ? '删除' : 'Delete', exact: true })
    await deleteGroup.click()
    await dialog.getByRole('button', { name: chinese ? '取消' : 'Cancel', exact: true }).click()
    expect(deletedCount).toBe(0)
    await deleteGroup.click()
    await dialog.getByRole('button', { name: chinese ? '确定' : 'Confirm', exact: true }).click()
    await expect.poll(() => deletedCount).toBe(1)
    expect(nativeDialogs).toEqual([])
    expect(pageErrors).toEqual([])

    // Switch language from the menu in the same session, then revisit account.
    await page.getByRole('link', { name: chinese ? '返回' : 'Back', exact: true }).click()
    await page.getByRole('button', { name: chinese ? '用户菜单' : 'User menu', exact: true }).click()
    await page.getByRole('menuitem', { name: chinese ? '语言' : 'Language', exact: true }).click()
    await page.getByRole('menuitemradio', { name: chinese ? 'English' : '简体中文' }).click()
    await page.getByRole('button', { name: chinese ? 'User menu' : '用户菜单', exact: true }).click()
    await page.getByRole('menuitem', { name: chinese ? 'Account' : '账号', exact: true }).click()
    await expect(page.getByRole('button', { name: chinese ? 'Choose file' : '选择文件', exact: true })).toBeVisible()
    console.log(`PASS about, file upload, dialogs, validation and language switching: ${language}`)
    await page.close()
  }
} finally {
  await browser?.close()
  await new Promise((resolve, reject) => server.httpServer.close(error => error ? reject(error) : resolve()))
}
