// Browser regression for the shared SPA used by desktop. All API responses
// are mocked: these checks never authenticate against or reset a real account.
import { chromium, expect } from '@playwright/test'
import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { preview } from 'vite'

const server = await preview({
  root: fileURLToPath(new URL('..', import.meta.url)),
  configFile: false,
  preview: { host: '127.0.0.1', port: 0, open: false },
})
const baseURL = `http://127.0.0.1:${server.httpServer.address().port}`
const executablePath = [process.env.GEMSNOTE_CHROMIUM_PATH, '/usr/bin/chromium', '/usr/bin/chromium-browser', '/usr/bin/google-chrome'].find(path => path && existsSync(path))
let browser
try {
  browser = await chromium.launch({ headless: true, ...(executablePath ? { executablePath } : {}) })
  for (const language of ['zh-CN', 'en-US']) {
    for (const scenario of ['cached', 'fresh', 'cache-discovered', 'skip']) {
      const hasCache = scenario === 'cached' || scenario === 'skip'
      const page = await browser.newPage()
      const resets = []
      let loggedIn = false
      let snapshot = { Running: false }
      let polls = 0
      let mediaPolls = 0, mediaRunning = false
      await page.addInitScript(({language, hasNativeEvents}) => {
        localStorage.setItem('gemsnote:language', language)
        if (!hasNativeEvents) return
        // Use the Wails event subscription API, not DOM progress events, to
        // verify that backend progress continues updating the opened dialog.
        const listeners = new Map()
        window.runtime = { EventsOn(name, callback) {
          const callbacks = listeners.get(name) || new Set()
          callbacks.add(callback)
          listeners.set(name, callbacks)
          return () => callbacks.delete(callback)
        } }
        window.emitTestWailsEvent = (name, value) => listeners.get(name)?.forEach(callback => callback(value))
      }, {language, hasNativeEvents: scenario === 'cached'})
      await page.route('**/api2/**', async route => {
        const path = new URL(route.request().url()).pathname
        if (path === '/api2/web/resetSync') {
          // Leave the request pending so the UI is checked before any server
          // response or native progress event (including a slow preflight).
          resets.push(route)
          return
        }
        let data
        if (path === '/api2/web/downloadStatus') {
          mediaPolls++
          data = { Running: mediaRunning, Completed: 0, Total: 14 }
        } else if (path === '/api2/web/syncProgress') {
          polls++
          data = snapshot
        } else if (path === '/api2/bootstrap') {
          data = { Desktop: true, Host: 'https://gemsnote.example', User: loggedIn ? { UserId: 'test-user', Username: 'tester' } : null, Notebooks: [], Tags: [] }
        } else if (path === '/api2/auth/session') {
          loggedIn = true
          data = { Ok: true, SyncChoiceRequired: hasCache, ResetSyncRequired: !hasCache }
        } else if (path === '/api2/notes') {
          data = []
        } else {
          throw new Error(`Unexpected API request: ${path}`)
        }
        await route.fulfill({ json: data })
      })
      const chinese = language === 'zh-CN'
      await page.goto(`${baseURL}/login`)
      await page.getByLabel(chinese ? '邮箱或用户名' : 'Email or username').fill('tester')
      await page.getByLabel(chinese ? '密码' : 'Password', { exact: true }).fill('test-password')
      await page.getByRole('button', { name: chinese ? '登录' : 'Sign in', exact: true }).click()
      const choice = page.getByRole('dialog', { name: chinese ? '选择登录后的同步方式' : 'Choose sync after sign-in' })
      const progress = page.getByRole('dialog', { name: chinese ? '重新同步' : 'Reset and sync', exact: true })
      if (hasCache) {
        await expect(choice).toBeVisible()
        await expect(progress).toHaveCount(0)
        expect(resets).toHaveLength(0)
        if (scenario === 'skip') {
          await choice.getByRole('button', { name: chinese ? '暂不同步' : 'Not now', exact: true }).click()
          await expect(page).toHaveURL(`${baseURL}/note`)
          expect(resets).toHaveLength(0)
          console.log(`PASS login preserves cache without syncing: ${language}`)
          await page.close()
          continue
        }
        await choice.getByRole('button', { name: chinese ? '重新同步' : 'Reset and sync', exact: true }).click()
      }
      await expect(progress).toBeVisible()
      await expect(choice).toHaveCount(0)
      await expect(progress.locator('progress')).not.toHaveAttribute('value')
      await expect.poll(() => resets.length).toBe(1)
      const firstReset = new URLSearchParams(resets[0].request().postData())
      expect(firstReset.get('confirm')).toBe(hasCache ? 'true' : null)
      expect(firstReset.get('initial')).toBe(hasCache ? null : 'true')
      // Real completed counts keep moving without Wails events. Unknown
      // totals render an indeterminate bar instead of a frozen fake percent.
      snapshot = { Running: true, Mode: 'reset', Stage: 'notes', Current: 7, Total: 0, Percent: 40 }
      await expect(progress.locator('.sync-completed')).toContainText(chinese ? '已完成 7 项' : 'Completed 7 items')
      snapshot = { ...snapshot, Current: 21 }
      await expect(progress.locator('.sync-completed')).toContainText(chinese ? '已完成 21 项' : 'Completed 21 items')
      expect(polls).toBeGreaterThan(0)
      await expect(progress).toContainText(chinese ? '正在同步笔记' : 'Synchronizing notes')
      if (scenario === 'cached') {
        snapshot = { Running: true, Mode: 'reset', Stage: 'images', Current: 3, Total: 10, Percent: 89 }
        await page.evaluate(value => window.emitTestWailsEvent('sync-progress', value), snapshot)
        await expect(progress.locator('progress')).toHaveAttribute('value', '89')
        await expect(progress.locator('.sync-completed')).toContainText('3 / 10')
        await page.evaluate(() => window.emitTestWailsEvent('sync-finished', { Ok: true, Full: true }))
        await expect(progress).toBeVisible() // HTTP completion is authoritative.
      }

      if (hasCache || scenario === 'cache-discovered') {
        // Failure must dismiss progress even without sync-finished, restore
        // the choice and error, and leave the reset action usable for retry.
        snapshot = { Running: false }
        await resets[0].fulfill({ json: { Ok: false, Msg: hasCache ? 'TLS handshake timeout' : 'confirmationRequired' } })
        await expect(progress).toHaveCount(0)
        await expect(choice).toBeVisible()
        if (hasCache) await expect(choice.getByRole('alert')).toContainText('TLS handshake timeout')
        else await expect(choice.getByRole('alert')).toHaveCount(0)
        await choice.getByRole('button', { name: chinese ? '重新同步' : 'Reset and sync', exact: true }).click()
        await expect(progress).toBeVisible()
        await expect(choice).toHaveCount(0)
        await expect.poll(() => resets.length).toBe(2)
        expect(new URLSearchParams(resets[1].request().postData()).get('confirm')).toBe('true')
      }
      // Success must close the dialog and navigate using the HTTP result,
      // even if the native completion event was missed.
      mediaRunning = true
      await resets.at(-1).fulfill({ json: { Ok: true, Full: true, Reset: true } })
      await expect(page).toHaveURL(`${baseURL}/note`)
      await expect.poll(() => mediaPolls).toBeGreaterThan(0)
      await expect(progress).toHaveCount(0)
      await expect(choice).toHaveCount(0)
      if (scenario === 'cached') {
        await page.evaluate(() => window.emitTestWailsEvent('sync-progress', { Running: true, Mode: 'reset', Stage: 'notes', Current: 1 }))
        await expect(progress).toHaveCount(0)
      }
      console.log(`PASS login reset progress: ${language}, scenario=${scenario}`)
      await page.close()
    }
  }
} finally {
  await browser?.close()
  await new Promise((resolve, reject) => server.httpServer.close(error => error ? reject(error) : resolve()))
}
