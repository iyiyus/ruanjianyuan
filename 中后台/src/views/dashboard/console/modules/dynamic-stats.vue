<template>
  <div class="art-card h-128 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>最新黑名单</h4>
        <p>最近被拉黑的UDID</p>
      </div>
    </div>
    <div class="h-9/10 mt-2 overflow-hidden">
      <ElScrollbar>
        <div
          class="h-17.5 leading-17.5 border-b border-g-300 text-sm overflow-hidden last:border-b-0"
          v-for="(item, index) in list"
          :key="index"
        >
          <span class="text-g-800 font-medium text-xs">{{ item.udid }}</span>
          <span class="ml-2 text-g-500 text-xs">{{ item.time }}</span>
        </div>
        <div v-if="list.length === 0" class="text-center text-g-400 text-sm mt-8">暂无黑名单记录</div>
      </ElScrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { fetchBlackList } from '@/api/source'

  const list = ref<{ udid: string; time: string }[]>([])

  onMounted(async () => {
    try {
      const res = await fetchBlackList({ current: 1, size: 8 })
      list.value = (res.records || []).map((item: any) => ({
        udid: item.udid,
        time: item.addtime ? new Date(item.addtime * 1000).toLocaleDateString('zh-CN') : '-'
      }))
    } catch {}
  })
</script>
