import { defineStore } from 'pinia'
import { fetchSiteConfig } from '@/api/source'

export const useSiteStore = defineStore('siteStore', () => {
  const siteName = ref('iOS软件源管理系统')
  const siteLogo = ref('')

  async function loadSiteConfig() {
    try {
      const data = await fetchSiteConfig()
      if (data.site_name) siteName.value = data.site_name
      if (data.site_logo) siteLogo.value = data.site_logo
    } catch {}
  }

  return { siteName, siteLogo, loadSiteConfig }
}, {
  persist: { key: 'site', storage: localStorage }
})
