<template>
  <div class="art-card p-5 flex-b mb-5 max-sm:mb-4">
    <div>
      <h2 class="text-2xl font-medium">关于系统</h2>
      <p class="text-g-700 mt-1">{{ systemName }} 是一套基于 Go + Vue3 开发的 iOS 软件源管理系统</p>
      <p class="text-g-700 mt-1">支持 App 管理、卡密授权、UDID 黑名单、云存储等完整功能</p>

      <div class="flex flex-wrap gap-3.5 max-w-150 mt-6">
        <div class="px-3.5 py-2 border border-g-300 rounded-lg text-sm bg-g-100">
          <span class="text-g-500 mr-1">后端</span><span class="text-g-800 font-medium">Go · Gin · GORM</span>
        </div>
        <div class="px-3.5 py-2 border border-g-300 rounded-lg text-sm bg-g-100">
          <span class="text-g-500 mr-1">前端</span><span class="text-g-800 font-medium">Vue3 · TypeScript · Element Plus</span>
        </div>
        <div class="px-3.5 py-2 border border-g-300 rounded-lg text-sm bg-g-100">
          <span class="text-g-500 mr-1">数据库</span><span class="text-g-800 font-medium">MySQL</span>
        </div>
        <div class="px-3.5 py-2 border border-g-300 rounded-lg text-sm bg-g-100">
          <span class="text-g-500 mr-1">存储</span><span class="text-g-800 font-medium">本地 / 七牛 / 腾讯 / 阿里 / WebDAV</span>
        </div>
        <div class="px-3.5 py-2 border border-g-300 rounded-lg text-sm bg-g-100">
          <span class="text-g-500 mr-1">认证</span><span class="text-g-800 font-medium">JWT</span>
        </div>
        <div
          class="flex-cc gap-2 h-12.5 px-5 border border-g-300 c-p rounded-lg text-sm bg-g-100 duration-300 hover:-translate-y-1"
          @click="goPage('https://qm.qq.com/q/XSuLUlM7qc')"
        >
          <ArtSvgIcon icon="ri:qq-line" class="text-lg text-g-600" />
          <span class="text-g-700">联系作者</span>
        </div>
      </div>
    </div>
    <img class="w-75 max-md:!hidden" src="@imgs/draw/draw1.png" alt="draw1" />
  </div>
</template>

<script setup lang="ts">
  import { useSiteStore } from '@/store/modules/site'

  const siteStore = useSiteStore()
  const systemName = computed(() => siteStore.siteName)

  const goPage = (url: string): void => {
    window.open(url, '_blank', 'noopener,noreferrer')
  }
</script>
