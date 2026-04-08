<template>
  <ElRow :gutter="20" class="flex">
    <ElCol v-for="(item, index) in dataList" :key="index" :sm="12" :md="6" :lg="6">
      <div class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4">
        <span class="text-g-700 text-sm">{{ item.des }}</span>
        <ArtCountTo class="text-[26px] font-medium mt-2" :target="item.num" :duration="1300" />
        <div
          class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10"
        >
          <ArtSvgIcon :icon="item.icon" class="text-xl text-theme" />
        </div>
      </div>
    </ElCol>
  </ElRow>
</template>

<script setup lang="ts">
  import { fetchDashboard } from '@/api/source'

  interface CardDataItem {
    des: string
    icon: string
    num: number
  }

  const dataList = reactive<CardDataItem[]>([
    { des: 'App总数', icon: 'ri:app-store-line', num: 0 },
    { des: '卡密总数', icon: 'ri:key-2-line', num: 0 },
    { des: '已使用卡密', icon: 'ri:check-double-line', num: 0 },
    { des: 'UDID黑名单', icon: 'ri:forbid-line', num: 0 },
    { des: 'UDID监控', icon: 'ri:eye-line', num: 0 }
  ])

  onMounted(async () => {
    try {
      const data = await fetchDashboard()
      dataList[0].num = data.appCount
      dataList[1].num = data.kamiTotal
      dataList[2].num = data.kamiUsed
      dataList[3].num = data.blackCount
      dataList[4].num = data.monitorCount
    } catch {}
  })
</script>
