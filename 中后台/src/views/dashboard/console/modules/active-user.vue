<template>
  <div class="art-card h-105 p-4 box-border mb-5 max-sm:mb-4">
    <ArtBarChart
      class="box-border p-2"
      barWidth="50%"
      height="13.7rem"
      :showAxisLine="false"
      :data="chartData"
      :xAxisData="xAxisLabels"
    />
    <div class="ml-1">
      <h3 class="mt-5 text-lg font-medium">卡密使用概览</h3>
      <p class="mt-1 text-sm">近9个月卡密激活数量统计</p>
    </div>
    <div class="flex-b mt-2">
      <div class="flex-1" v-for="(item, index) in list" :key="index">
        <p class="text-2xl text-g-900">{{ item.num }}</p>
        <p class="text-xs text-g-500">{{ item.name }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { fetchDashboard } from '@/api/source'

  const xAxisLabels = ref<string[]>([])
  const chartData = ref<number[]>([])

  const list = ref([
    { name: 'App总数', num: '0' },
    { name: '卡密总数', num: '0' },
    { name: '已激活', num: '0' },
    { name: '黑名单', num: '0' }
  ])

  onMounted(async () => {
    try {
      const data = await fetchDashboard()
      list.value[0].num = String(data.appCount)
      list.value[1].num = String(data.kamiTotal)
      list.value[2].num = String(data.kamiUsed)
      list.value[3].num = String(data.blackCount)
      xAxisLabels.value = data.kamiLabels || []
      chartData.value = data.kamiMonthly || []
    } catch {}
  })
</script>
