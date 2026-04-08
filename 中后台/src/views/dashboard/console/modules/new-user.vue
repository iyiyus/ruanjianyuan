<template>
  <div class="art-card p-5 h-128 overflow-hidden mb-5 max-sm:mb-4">
    <div class="art-card-header">
      <div class="title">
        <h4>最新App列表</h4>
        <p>最近更新的应用</p>
      </div>
    </div>
    <ArtTable
      class="w-full"
      :data="tableData"
      style="width: 100%"
      size="large"
      :border="false"
      :stripe="false"
      :header-cell-style="{ background: 'transparent' }"
    >
      <template #default>
        <ElTableColumn label="图标" width="60px">
          <template #default="scope">
            <div style="width:36px;height:36px;border-radius:8px;background:#f5f5f5;display:flex;align-items:center;justify-content:center;overflow:hidden;flex-shrink:0">
              <img style="width:36px;height:36px;object-fit:cover" :src="scope.row.image || '/assets/img/avatar.png'" alt="icon" />
            </div>
          </template>
        </ElTableColumn>
        <ElTableColumn label="名称" prop="name" />
        <ElTableColumn label="版本" prop="nickname" width="90" />
        <ElTableColumn label="状态" width="80">
          <template #default="scope">
            <ElTag :type="scope.row.bt2b === '1' ? 'danger' : 'success'" size="small">
              {{ scope.row.bt2b === '1' ? '加锁' : '免费' }}
            </ElTag>
          </template>
        </ElTableColumn>
      </template>
    </ArtTable>
  </div>
</template>

<script setup lang="ts">
  import { fetchAppList } from '@/api/source'
  import { ElTag } from 'element-plus'

  const tableData = ref<any[]>([])

  onMounted(async () => {
    try {
      const res = await fetchAppList({ current: 1, size: 6 })
      tableData.value = res.records || []
    } catch {}
  })
</script>
