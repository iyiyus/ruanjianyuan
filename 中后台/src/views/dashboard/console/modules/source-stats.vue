<template>
  <ElRow :gutter="16" class="mb-4">
    <ElCol :xs="12" :sm="8" :md="4" v-for="item in stats" :key="item.label">
      <ElCard shadow="never" class="text-center py-4">
        <div class="text-2xl font-bold" :style="{ color: item.color }">{{ item.value }}</div>
        <div class="text-sm text-gray-500 mt-1">{{ item.label }}</div>
      </ElCard>
    </ElCol>
  </ElRow>
</template>

<script setup lang="ts">
  import { fetchDashboard } from '@/api/source'

  const stats = ref([
    { label: 'App总数', value: 0, color: '#5D87FF' },
    { label: '卡密总数', value: 0, color: '#B48DF3' },
    { label: '已使用卡密', value: 0, color: '#60C041' },
    { label: '黑名单数', value: 0, color: '#FF5B5B' },
    { label: '监控记录', value: 0, color: '#F9901F' }
  ])

  onMounted(async () => {
    try {
      const data = await fetchDashboard()
      stats.value[0].value = data.appCount
      stats.value[1].value = data.kamiTotal
      stats.value[2].value = data.kamiUsed
      stats.value[3].value = data.blackCount
      stats.value[4].value = data.monitorCount
    } catch {}
  })
</script>
