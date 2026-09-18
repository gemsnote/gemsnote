import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
export default defineConfig({
  plugins: [vue()],
  server: { proxy: Object.fromEntries(['/web/', '/api/', '/doLogin', '/doRegister', '/doFindPassword', '/findPasswordUpdate', '/captcha/', '/user/', '/notebook/', '/share/', '/attach/', '/file/', '/noteContentHistory/', '/note/', '/member/', '/tinymce/', '/public/', '/logout'].map(path => [path, 'http://127.0.0.1:9000'])) },
})
