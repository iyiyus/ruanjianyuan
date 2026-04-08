<template>
  <div class="art-full-height">
    <ArtSearchBar v-model="searchForm" :items="searchItems" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="danger" @click="addDialogVisible = true" v-ripple>手动添加</ElButton>
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

    <ElDialog v-model="addDialogVisible" title="添加黑名单" width="400px" align-center>
      <ElForm ref="formRef" :model="addForm" :rules="rules" label-width="80px">
        <ElFormItem label="UDID" prop="udid">
          <ElInput v-model="addForm.udid" placeholder="请输入设备UDID" />
        </ElFormItem>
        <ElFormItem label="拉黑原因">
          <ElInput v-model="addForm.reason" placeholder="可选，如：盗版、破解等" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="addDialogVisible = false">取消</ElButton>
        <ElButton type="danger" :loading="addLoading" @click="handleAdd">添加</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { fetchBlackList, fetchCreateBlack, fetchDeleteBlack } from '@/api/source'
  import { useTable } from '@/hooks/core/useTable'
  import { ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'

  defineOptions({ name: 'Black' })

  const addDialogVisible = ref(false)
  const addLoading = ref(false)
  const formRef = ref<FormInstance>()
  const addForm = ref({ udid: '', reason: '' })
  const rules: FormRules = { udid: [{ required: true, message: '请输入UDID', trigger: 'blur' }] }

  const searchForm = ref({ udid: '' })
  const searchItems = [{ label: 'UDID', key: 'udid', type: 'input', props: { placeholder: '请输入UDID' } }]

  const { columns, columnChecks, data, loading, pagination, getData, searchParams, resetSearchParams, handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: fetchBlackList,
      apiParams: { current: 1, size: 20 },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        { prop: 'udid', label: 'UDID', minWidth: 220 },
        { prop: 'reason', label: '拉黑原因', minWidth: 120 },
        {
          prop: 'addtime', label: '添加时间', width: 180,
          formatter: (row) => row.addtime ? new Date(row.addtime * 1000).toLocaleString('zh-CN') : '-'
        },
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

  async function handleAdd() {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (!valid) return
      addLoading.value = true
      try {
        await fetchCreateBlack({ udid: addForm.value.udid.trim(), reason: addForm.value.reason.trim() })
        ElMessage.success('添加成功')
        addDialogVisible.value = false
        addForm.value.udid = ''
        addForm.value.reason = ''
        getData()
      } finally {
        addLoading.value = false
      }
    })
  }

  function handleDelete(row: any) {
    ElMessageBox.confirm('确定将该UDID移出黑名单吗？', '移除确认', { type: 'warning' }).then(async () => {
      await fetchDeleteBlack(row.id)
      ElMessage.success('移除成功')
      getData()
    })
  }
</script>
