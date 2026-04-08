<!-- 系统logo -->
<template>
  <div class="flex-cc">
    <img :style="logoStyle" :src="logoSrc" alt="logo" class="w-full h-full" />
  </div>
</template>

<script setup lang="ts">
  import { useSiteStore } from '@/store/modules/site'
  import defaultLogo from '@imgs/common/logo.webp'

  defineOptions({ name: 'ArtLogo' })

  interface Props {
    size?: number | string
  }

  const props = withDefaults(defineProps<Props>(), { size: 36 })
  const siteStore = useSiteStore()
  const logoSrc = computed(() => siteStore.siteLogo || defaultLogo)
  const logoStyle = computed(() => ({ width: `${props.size}px` }))
</script>
