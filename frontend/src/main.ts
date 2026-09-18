import { createApp, h } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Auth from './pages/Auth.vue'
import Workspace from './pages/Workspace.vue'
import Account from './pages/Account.vue'
import Admin from './pages/Admin.vue'
import VerifyEmail from './pages/VerifyEmail.vue'
import './styles.css'
import {t} from './i18n'
const router = createRouter({ history: createWebHistory(), routes: [
  { path: '/', redirect: '/note' }, { path: '/index', redirect: '/note' },
  { path: '/login', component: Auth }, { path: '/register', component: Auth },
  { path: '/findPassword/:token?', component: Auth },
  { path: '/note/:noteId?', component: Workspace },
  { path: '/member/:pathMatch(.*)*', component: Account },
  { path: '/user/account', component: Account },
  { path: '/user/activeEmail', component: VerifyEmail },
  { path: '/user/updateEmail', component: VerifyEmail },
  { path: '/admin/:pathMatch(.*)*', component: Admin },
  { path: '/:pathMatch(.*)*', component: { render: () => h('main', { class: 'auth-card' }, [h('h1', t('页面不存在')), h('a', { href: '/note' }, t('返回笔记'))]) } },
] })
window.addEventListener('session-expired', () => router.push('/login'))
createApp(App).use(createPinia()).use(router).mount('#app')
