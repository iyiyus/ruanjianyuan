<template>
  <div class="art-full-height">
    <ArtSearchBar v-model="searchForm" :items="searchItems" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="danger" plain @click="handleClearAll" v-ripple>清空全部</ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { fetchMonitorList, fetchDeleteMonitor, fetchClearMonitor } from '@/api/source'
  import { useTable } from '@/hooks/core/useTable'
  import { ElTag, ElMessageBox } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'

  defineOptions({ name: 'Monitor' })

  const searchForm = ref({ udid: '', identity: '' })
  const searchItems = [
    { label: 'UDID', key: 'udid', type: 'input', props: { placeholder: '请输入UDID' } },
    {
      label: '身份', key: 'identity', type: 'select',
      props: { placeholder: '请选择身份', options: [{ label: '添加者', value: '添加者' }, { label: '破解者', value: '破解者' }] }
    }
  ]

  const { columns, columnChecks, data, loading, pagination, getData, searchParams, resetSearchParams, handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: fetchMonitorList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'udid', label: 'UDID', minWidth: 220 },
        {
          prop: 'identity', label: '身份', width: 110,
          formatter: (row) => h(ElTag, { type: row.identity === '添加者' ? 'warning' : 'danger' }, () => row.identity)
        },
        { prop: 'count', label: '触发次数', width: 100 },
        { prop: 'addtime', label: '首次时间', width: 160 },
        {
          prop: 'operation', label: '操作', width: 80, fixed: 'right',
          formatter: (row) => h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
        }
      ]
    }
  })

  function handleSearch(params: Record<string, any>) {
    Object.assign(searchParams, params)
    getData()
  }

  function handleDelete(row: any) {
    ElMessageBox.confirm('确定删除该监控记录吗？', '删除确认', { type: 'warning' }).then(async () => {
      await fetchDeleteMonitor(row.id)
      ElMessage.success('删除成功')
      getData()
    })
  }

  function handleClearAll() {
    ElMessageBox.confirm('确定清空所有监控记录吗？此操作不可恢复！', '清空确认', { type: 'error' }).then(async () => {
      await fetchClearMonitor()
      ElMessage.success('清空成功')
      getData()
    })
  }
</script>
