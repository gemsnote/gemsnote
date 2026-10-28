// Real API2 authentication regression. Run only against an explicitly
// disposable loopback server initialized with the bundled admin seed.
import { chromium, expect } from '@playwright/test'
import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { preview } from 'vite'

if (process.env.GEMSNOTE_DISPOSABLE_API2 !== '1') {
  throw new Error('Set GEMSNOTE_DISPOSABLE_API2=1 for an isolated test server')
}
const backend = new URL(process.env.GEMSNOTE_API2_URL || 'http://127.0.0.1:19090')
if (backend.origin !== 'http://127.0.0.1:19090' || backend.pathname !== '/' || backend.search || backend.hash || backend.username || backend.password) {
  throw new Error('Only the disposable API2 fixture at http://127.0.0.1:19090 is allowed')
}
const server = await preview({
  root: fileURLToPath(new URL('..', import.meta.url)),
  configFile: false,
  preview: {
    host: '127.0.0.1', port: 0, open: false,
    proxy: Object.fromEntries(['/api2/', '/web/bootstrap', '/doLogin'].map(path => [path, backend.origin])),
  },
})
const baseURL = `http://127.0.0.1:${server.httpServer.address().port}`
const executablePath = [process.env.GEMSNOTE_CHROMIUM_PATH, '/usr/bin/chromium', '/usr/bin/chromium-browser', '/usr/bin/google-chrome'].find(path => path && existsSync(path))
let browser
try {
  browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
  for (const expiredToken of [false, true]) {
    const context = await browser.newContext({ locale: 'zh-CN' })
    try {
      const legacy = await context.request.post(`${baseURL}/doLogin`, { form: { email: 'admin', pwd: 'gemsnote' } })
      expect((await legacy.json()).Ok).toBe(true)
      // Establish that the cookie really is valid for the legacy Web entry.
      const cookieBootstrap = await context.request.get(`${baseURL}/web/bootstrap`)
      expect(Boolean((await cookieBootstrap.json()).User?.UserId)).toBe(true)
      const anonymous = await context.request.get(`${baseURL}/api2/bootstrap`, { headers: expiredToken ? { Authorization: 'Bearer expired-fixture-token' } : {} })
      expect((await anonymous.json()).User).toBeNull()

      const page = await context.newPage()
      if (expiredToken) {
        await page.addInitScript(() => sessionStorage.setItem('gemsnote:api2-token', 'expired-fixture-token'))
      }
      let bootstraps = 0, notes = 0
      page.on('request', request => {
        const path = new URL(request.url()).pathname
        if (path === '/api2/bootstrap') bootstraps++
        if (path === '/api2/notes') notes++
      })
      await page.goto(`${baseURL}/note`)
      await expect(page).toHaveURL(`${baseURL}/login`)
      await expect(page.getByRole('button', { name: '登录', exact: true })).toBeVisible()
      // A redirect loop would keep issuing bootstrap and note requests after
      // the first login screen appears. Allow the router/network queue to run.
      await page.waitForTimeout(500)
      expect(bootstraps).toBe(2)
      expect(notes).toBe(0)
      console.log(`PASS legacy cookie cannot loop API2 pages; expired token=${expiredToken}`)
    } finally { await context.close() }
  }

  const context = await browser.newContext({ locale: 'zh-CN' })
  try {
    const page = await context.newPage()
    await page.goto(`${baseURL}/login`)
    await page.getByLabel('邮箱或用户名').fill('admin')
    await page.getByLabel('密码', { exact: true }).fill('gemsnote')
    await page.getByRole('button', { name: '登录', exact: true }).click()
    await expect(page).toHaveURL(`${baseURL}/note`)
    await expect(page.getByLabel('搜索笔记本')).toBeVisible()
    expect(await page.evaluate(() => Boolean(sessionStorage.getItem('gemsnote:api2-token')))).toBe(true)

    // Bearer login and a full reload must work without a browser cookie.
    await context.clearCookies()
    await page.reload()
    await expect(page.getByLabel('搜索笔记本')).toBeVisible()
    await expect(page).toHaveURL(`${baseURL}/note`)
    console.log('PASS bearer-only Web login and reload')

    // A revoked token must also produce an anonymous bootstrap; a previously
    // authenticated Vue page must settle on login rather than bounce back.
    await page.evaluate(async () => {
      const token = sessionStorage.getItem('gemsnote:api2-token')
      await fetch('/api2/auth/logout', { method: 'POST', headers: { Authorization: `Bearer ${token}` } })
    })
    await page.reload()
    await expect(page).toHaveURL(`${baseURL}/login`)
    await expect(page.getByRole('button', { name: '登录', exact: true })).toBeVisible()
    await page.waitForTimeout(500)
    await expect(page).toHaveURL(`${baseURL}/login`)
    console.log('PASS revoked bearer returns to a stable login screen')
  } finally { await context.close() }
} finally {
  if (browser) await browser.close()
  await server.close()
}
