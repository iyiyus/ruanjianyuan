<template>
  <div class="art-card h-128 p-5 mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>最新监控</h4>
        <p>异常UDID记录</p>
      </div>
    </div>
    <div class="h-[calc(100%-40px)] overflow-auto">
      <ElScrollbar>
        <div
          class="flex-cb h-17.5 border-b border-g-300 text-sm last:border-b-0"
          v-for="(item, index) in list"
          :key="index"
        >
          <div>
            <p class="text-xs text-g-800 truncate" style="max-width:160px">{{ item.udid }}</p>
            <p class="text-g-500 mt-1 text-xs">{{ item.identity }} · 触发{{ item.count }}次</p>
          </div>
          <ElTag :type="item.identity === '添加者' ? 'warning' : 'danger'" size="small">
            {{ item.identity }}
          </ElTag>
        </div>
        <div v-if="list.length === 0" class="text-center text-g-400 text-sm mt-8">暂无监控记录</div>
      </ElScrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { fetchMonitorList } from '@/api/source'
  import { ElTag } from 'element-plus'

  const list = ref<any[]>([])

  onMounted(async () => {
    try {
      const res = await fetchMonitorList({ current: 1, size: 6 })
      list.value = res.records || []
    } catch {}
  })
</script>
