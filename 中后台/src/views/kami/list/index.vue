<template>
  <div class="art-full-height">
    <ArtSearchBar v-model="searchForm" :items="searchItems" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="genDialogVisible = true" v-ripple>生成卡密</ElButton>
          <ElButton v-if="lastCodes" plain @click="copyLastCodes" v-ripple>复制最新卡密</ElButton>
          <ElButton v-if="selectedIds.length > 0" type="danger" plain @click="handleBatchDelete" v-ripple>
            删除选中 ({{ selectedIds.length }})
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>

    <!-- 生成卡密 -->
    <ElDialog v-model="genDialogVisible" title="批量生成卡密" width="400px" align-center>
      <ElForm :model="genForm" label-width="80px">
        <ElFormItem label="数量" required>
          <ElInputNumber v-model="genForm.count" :min="1" :max="1000" />
        </ElFormItem>
        <ElFormItem label="类型" required>
          <ElRadioGroup v-model="genForm.kmyp">
            <ElRadio :value="1">30天</ElRadio>
            <ElRadio :value="2">90天</ElRadio>
            <ElRadio :value="3">365天</ElRadio>
            <ElRadio :value="4">自定义</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="genForm.kmyp === 4" label="天数" required>
          <ElInputNumber v-model="genForm.customDays" :min="1" :max="36500" />
        </ElFormItem>
        <ElFormItem label="前缀">
          <ElInput v-model="genForm.prefix" placeholder="可选，如 VIP" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="genDialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="genLoading" @click="handleGenerate">生成</ElButton>
      </template>
    </ElDialog>

    <!-- 生成结果 -->
    <ElDialog v-model="resultVisible" title="生成结果" width="500px" align-center>
      <p class="mb-3 text-sm text-g-600">共生成 {{ genResult.count }} 个卡密</p>
      <ElInput v-model="genResult.codes" type="textarea" :rows="10" readonly />
      <template #footer>
        <ElButton type="primary" @click="copyResult">复制全部</ElButton>
        <ElButton @click="resultVisible = false">关闭</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { fetchKamiList, fetchGenerateKami, fetchDeleteKami, batchDeleteKami } from '@/api/source'
  import { useTable } from '@/hooks/core/useTable'
  import { ElTag, ElMessageBox } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useClipboard } from '@vueuse/core'

  defineOptions({ name: 'KamiList' })

  const genDialogVisible = ref(false)
  const genLoading = ref(false)
  const resultVisible = ref(false)
  const genResult = ref({ count: 0, codes: '' })
  const genForm = ref({ count: 10, kmyp: 1, prefix: '', customDays: 7 })
  const lastCodes = ref('')
  const selectedIds = ref<number[]>([])
  const { copy } = useClipboard()

  const KMYP_MAP: Record<number, string> = { 1: '30天', 2: '90天', 3: '365天' }

  const searchForm = ref({ jh: '', kmyp: '' })
  const searchItems = [
    {
      label: '状态', key: 'jh', type: 'select',
      props: { placeholder: '请选择状态', options: [{ label: '未使用', value: '0' }, { label: '已使用', value: '1' }] }
    },
    {
      label: '类型', key: 'kmyp', type: 'select',
      props: { placeholder: '请选择类型', options: [{ label: '30天', value: '1' }, { label: '90天', value: '2' }, { label: '365天', value: '3' }] }
    }
  ]

  const { columns, columnChecks, data, loading, pagination, getData, searchParams, resetSearchParams, handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: fetchKamiList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'selection' },
        {
          prop: 'kami', label: '卡密', minWidth: 160,
          formatter: (row) => h(
            'span',
            {
              class: 'cursor-pointer text-theme hover:underline',
              title: '点击复制',
              onClick: () => { copy(row.kami); ElMessage.success('已复制：' + row.kami) }
            },
            row.kami
          )
        },
        {
          prop: 'kmyp', label: '类型', width: 90,
          formatter: (row) => h(ElTag, { type: 'info' }, () => KMYP_MAP[row.kmyp] || '-')
        },
        {
          prop: 'jh', label: '状态', width: 90,
          formatter: (row) => h(ElTag, { type: row.jh === 1 ? 'success' : 'warning' }, () => row.jh === 1 ? '已使用' : '未使用')
        },
        { prop: 'udid', label: '绑定UDID', minWidth: 200 },
        {
          prop: 'addtime', label: '生成时间', width: 160,
          formatter: (row) => row.addtime ? new Date(row.addtime * 1000).toLocaleString('zh-CN') : '-'
        },
        {
          prop: 'usetime', label: '使用时间', width: 160,
          formatter: (row) => row.usetime && row.usetime > 0 ? new Date(row.usetime * 1000).toLocaleString('zh-CN') : '-'
        },
        {
          prop: 'endtime', label: '到期时间', width: 160,
          formatter: (row) => row.endtime && row.endtime > 0 ? new Date(row.endtime * 1000).toLocaleString('zh-CN') : '-'
        },
        {
          prop: 'operation', label: '操作', width: 80, fixed: 'right',
          formatter: (row) => h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
        }
      ]
    }
  })

  function handleSelectionChange(rows: any[]) {
    selectedIds.value = rows.map((r) => r.id)
  }

  function handleSearch(params: Record<string, any>) {
    Object.assign(searchParams, params)
    getData()
  }

  async function handleGenerate() {
    genLoading.value = true
    try {
      const res = await fetchGenerateKami(genForm.value)
      genResult.value = res
      lastCodes.value = res.codes
      genDialogVisible.value = false
      resultVisible.value = true
      getData()
    } finally {
      genLoading.value = false
    }
  }

  function handleDelete(row: any) {
    ElMessageBox.confirm('确定删除该卡密吗？', '删除确认', { type: 'warning' }).then(async () => {
      await fetchDeleteKami(row.id)
      ElMessage.success('删除成功')
      getData()
    })
  }

  function handleBatchDelete() {
    ElMessageBox.confirm(`确定删除选中的 ${selectedIds.value.length} 条卡密吗？`, '批量删除', { type: 'warning' }).then(async () => {
      await batchDeleteKami(selectedIds.value)
      ElMessage.success('删除成功')
      selectedIds.value = []
      getData()
    })
  }

  function copyResult() {
    copy(genResult.value.codes)
    ElMessage.success('已复制到剪贴板')
  }

  function copyLastCodes() {
    copy(lastCodes.value)
    ElMessage.success('已复制到剪贴板')
  }
</script>
