<template>
  <div class="art-card h-105 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>UDID监控趋势</h4>
        <p>近12个月异常UDID记录</p>
      </div>
    </div>
    <ArtLineChart
      height="calc(100% - 56px)"
      :data="data"
      :xAxisData="xAxisData"
      :showAreaColor="true"
      :showAxisLine="false"
    />
  </div>
</template>

<script setup lang="ts">
  import { fetchDashboard } from '@/api/source'
  const data = ref<number[]>([])
  const xAxisData = ref<string[]>([])

  onMounted(async () => {
    try {
      const res = await fetchDashboard()
      xAxisData.value = res.monitorLabels || []
      data.value = res.monitorMonthly || []
    } catch {}
  })
</script>
