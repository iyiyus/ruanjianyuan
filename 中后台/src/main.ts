import App from './App.vue'
import { createApp } from 'vue'
import { initStore } from './store'
import { initRouter } from './router'
import language from './locales'
import '@styles/core/tailwind.css'
import '@styles/index.scss'
import '@utils/sys/console.ts'
import { setupGlobDirectives } from './directives'
import { setupErrorHandle } from './utils/sys/error-handle'
import { useSiteStore } from './store/modules/site'

document.addEventListener('touchstart', function () {}, { passive: false })

const app = createApp(App)
initStore(app)
initRouter(app)
setupGlobDirectives(app)
setupErrorHandle(app)
app.use(language)
app.mount('#app')

// 加载网站配置（挂载后异步执行，不阻塞渲染）
useSiteStore().loadSiteConfig().then(() => {
  const store = useSiteStore()
  // 更新标签页标题
  if (store.siteName) {
    document.title = store.siteName
  }
  // 更新 favicon
  if (store.siteLogo && store.siteLogo.trim()) {
    const link = document.querySelector("link[rel='shortcut icon']") as HTMLLinkElement
      || document.createElement('link') as HTMLLinkElement
    link.rel = 'shortcut icon'
    link.href = store.siteLogo
    document.head.appendChild(link)
  }
})