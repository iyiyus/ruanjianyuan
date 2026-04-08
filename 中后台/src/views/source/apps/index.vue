<template>
  <div class="art-full-height">
    <ArtSearchBar v-model="searchForm" :items="searchItems" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="showDialog('add')" v-ripple>新增App</ElButton>
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

      <ElDialog
        v-model="dialogVisible"
        :title="dialogType === 'add' ? '新增App' : '编辑App'"
        width="620px"
        align-center
      >
        <!-- IPA上传区 -->
        <ElAlert v-if="ipaUploading" title="正在解析IPA，请稍候..." type="info" :closable="false" class="mb-3" />
        <div class="mb-4 p-3 border border-dashed border-g-300 rounded-lg">
          <div class="flex items-center gap-3">
            <ElUpload :action="uploadIPAAction" :headers="uploadHeaders" :show-file-list="false"
              :before-upload="beforeIPAUpload" :on-success="handleIPASuccess" :on-error="handleIPAError"
              accept=".ipa" :disabled="ipaUploading">
              <ElButton type="primary" plain :loading="ipaUploading">
                <el-icon class="mr-1"><Upload /></el-icon>上传IPA自动解析
              </ElButton>
            </ElUpload>
            <span class="text-xs text-g-400">上传后自动填充名称、版本、图标、大小等信息</span>
          </div>
        </div>

        <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
          <ElRow :gutter="16">
            <ElCol :span="12">
              <ElFormItem label="App名称" prop="name">
                <ElInput v-model="form.name" placeholder="如 微信" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="版本号" prop="nickname">
                <ElInput v-model="form.nickname" placeholder="如 8.0.0" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="状态">
                <ElSelect v-model="form.status" style="width:100%">
                  <ElOption label="正常" value="normal" />
                  <ElOption label="隐藏" value="hidden" />
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :span="24">
              <ElFormItem label="图标URL">
                <div class="flex gap-2 w-full items-center">
                  <ElInput v-model="form.image" placeholder="图标图片地址" />
                  <img v-if="form.image" :src="fullUrl(form.image)" style="width:40px;height:40px;border-radius:8px;object-fit:cover;flex-shrink:0;border:1px solid var(--el-border-color)" @error="(e:any)=>e.target.style.display='none'" />
                  <ElUpload :action="uploadAction" :headers="uploadHeaders" :show-file-list="false"
                    :on-success="(r: any) => { if(r.code===200) form.image = r.data.url }"
                    :before-upload="beforeImgUpload" accept="image/*" style="flex-shrink:0">
                    <ElButton plain><el-icon><Upload /></el-icon></ElButton>
                  </ElUpload>
                </div>
              </ElFormItem>
            </ElCol>
            <ElCol :span="24">
              <ElFormItem label="下载链接">
                <div class="flex gap-2 w-full">
                  <ElInput v-model="form.bt1a" placeholder="ipa下载地址" />
                  <ElUpload :action="uploadIPAAction" :headers="uploadHeaders" :show-file-list="false"
                    :before-upload="beforeIPAUpload" :on-success="handleIPASuccess" :on-error="handleIPAError"
                    accept=".ipa" :disabled="ipaUploading" style="flex-shrink:0">
                    <ElButton plain :loading="ipaUploading"><el-icon><Upload /></el-icon></ElButton>
                  </ElUpload>
                </div>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="文件大小(MB)">
                <ElInput v-model="form.bt2a" placeholder="如 118.2" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="颜色">
                <div class="flex gap-2 w-full items-center">
                  <ElInput v-model="form.bt1b" placeholder="如 018084" />
                  <ElColorPicker v-model="colorPickerVal" @change="onColorChange" size="large" style="flex-shrink:0" />
                </div>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="分类类型">
                <ElSelect v-model="form.type" style="width:100%">
                  <ElOption label="默认" value="default" />
                  <ElOption label="应用" value="1" />
                  <ElOption label="游戏" value="2" />
                  <ElOption label="影音" value="3" />
                  <ElOption label="工具" value="4" />
                  <ElOption label="插件" value="5" />
                </ElSelect>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="是否加锁">
                <ElSwitch v-model="form.bt2b" active-value="1" inactive-value="0" active-text="加锁" inactive-text="免费" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="蓝奏云">
                <ElSwitch v-model="form.flag" active-value="1" inactive-value="0" active-text="是" inactive-text="否" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="权重">
                <ElInputNumber v-model="form.weigh" :min="0" style="width:100%" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="24">
              <ElFormItem label="软件说明">
                <ElInput v-model="form.keywords" type="textarea" :rows="3" placeholder="版本更新说明，换行用\n" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="24">
              <ElFormItem label="备注">
                <ElInput v-model="form.beizhu" placeholder="内部备注" />
              </ElFormItem>
            </ElCol>
          </ElRow>
        </ElForm>
        <template #footer>
          <ElButton @click="dialogVisible = false">取消</ElButton>
          <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">提交</ElButton>
        </template>
      </ElDialog>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { fetchAppList, fetchCreateApp, fetchUpdateApp, fetchDeleteApp } from '@/api/source'
  import { useTable } from '@/hooks/core/useTable'
  import { ElTag, ElImage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { DialogType } from '@/types'
  import { useUserStore } from '@/store/modules/user'
  import { Upload } from '@element-plus/icons-vue'

  // 补全相对路径
  const fullUrl = (url: string) => {
    if (!url) return ''
    if (url.startsWith('/')) return window.location.origin + url
    return url
  }

  defineOptions({ name: 'Apps' })

  const userStore = useUserStore()
  const uploadAction = '/api/upload'
  const uploadIPAAction = '/api/upload/ipa'
  const uploadHeaders = computed(() => ({ Authorization: userStore.accessToken }))
  const ipaUploading = ref(false)

  const dialogVisible = ref(false)
  const dialogType = ref<DialogType>('add')
  const submitLoading = ref(false)
  const currentId = ref(0)
  const formRef = ref<FormInstance>()

  // 颜色选择器（#RRGGBB → RRGGBB 互转）
  const colorPickerVal = computed({
    get: () => form.value.bt1b ? '#' + form.value.bt1b.replace('#', '') : '#018084',
    set: (v) => {}
  })
  function onColorChange(val: string | null) {
    if (val) form.value.bt1b = val.replace('#', '').toUpperCase()
  }

  const searchForm = ref({ name: '', status: '' })
  const searchItems = [
    { label: 'App名称', key: 'name', type: 'input', props: { placeholder: '请输入App名称' } },
    {
      label: '状态', key: 'status', type: 'select',
      props: { placeholder: '请选择状态', options: [{ label: '正常', value: 'normal' }, { label: '隐藏', value: 'hidden' }] }
    }
  ]

  const defaultForm = () => ({
    name: '', nickname: '', image: '', bt1a: '', bt1b: '018084',
    bt2a: '', bt2b: '0', keywords: '', description: '', diyname: '',
    beizhu: '', weigh: 0, status: 'normal', type: 'default', flag: '0', flag2: ''
  })
  const form = ref(defaultForm())

  const rules: FormRules = {
    name: [{ required: true, message: '请输入App名称', trigger: 'blur' }],
    nickname: [{ required: true, message: '请输入版本号', trigger: 'blur' }]
  }

  const { columns, columnChecks, data, loading, pagination, getData, searchParams, resetSearchParams, handleSizeChange, handleCurrentChange, refreshData } = useTable({
    core: {
      apiFn: fetchAppList,
      apiParams: { current: 1, size: 20, ...searchForm.value },
      columnsFactory: () => [
        { type: 'index', width: 60, label: '序号' },
        {
          prop: 'image', label: '图标', width: 70,
          formatter: (row: any) => h(ElImage, { src: fullUrl(row.image), style: 'width:40px;height:40px;border-radius:8px', previewSrcList: [fullUrl(row.image)], previewTeleported: true })
        },
        { prop: 'name', label: 'App名称', minWidth: 100 },
        { prop: 'keywords', label: '软件说明', minWidth: 150,
          formatter: (row: any) => row.keywords ? row.keywords.replace(/\\n/g, ' ') : '-'
        },
        { prop: 'nickname', label: '版本', width: 90 },
        {
          prop: 'type', label: '类型', width: 80,
          formatter: (row) => {
            const map: Record<string, string> = { default: '默认', '0': '默认', '1': '应用', '2': '游戏', '3': '影音', '4': '工具', '5': '插件' }
            return map[row.type] || row.type
          }
        },
        {
          prop: 'bt2b', label: '加锁', width: 75,
          formatter: (row) => h(ElTag, { type: row.bt2b === '1' ? 'danger' : 'success', size: 'small' }, () => row.bt2b === '1' ? '加锁' : '免费')
        },
        {
          prop: 'flag', label: '蓝奏云', width: 80,
          formatter: (row) => h(ElTag, { type: row.flag === '1' ? 'warning' : 'info', size: 'small' }, () => row.flag === '1' ? '是' : '否')
        },
        {
          prop: 'status', label: '状态', width: 75,
          formatter: (row) => h(ElTag, { type: row.status === 'normal' ? 'success' : 'info', size: 'small' }, () => row.status === 'normal' ? '正常' : '隐藏')
        },
        { prop: 'weigh', label: '权重', width: 70 },
        {
          prop: 'operation', label: '操作', width: 120, fixed: 'right',
          formatter: (row) => h('div', [
            h(ArtButtonTable, { type: 'edit', onClick: () => showDialog('edit', row) }),
            h(ArtButtonTable, { type: 'delete', onClick: () => handleDelete(row) })
          ])
        }
      ]
    }
  })

  function handleSearch(params: Record<string, any>) {
    Object.assign(searchParams, params)
    getData()
  }

  function showDialog(type: DialogType, row?: any) {
    dialogType.value = type
    if (type === 'edit' && row) {
      form.value = { ...defaultForm(), ...row }
      currentId.value = row.id
    } else {
      form.value = defaultForm()
      currentId.value = 0
    }
    nextTick(() => { dialogVisible.value = true })
  }

  async function handleSubmit() {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (!valid) return
      submitLoading.value = true
      try {
        if (dialogType.value === 'add') {
          await fetchCreateApp(form.value)
        } else {
          await fetchUpdateApp(currentId.value, form.value)
        }
        ElMessage.success(dialogType.value === 'add' ? '添加成功' : '更新成功')
        dialogVisible.value = false
        getData()
      } finally {
        submitLoading.value = false
      }
    })
  }

  function handleDelete(row: any) {
    ElMessageBox.confirm(`确定删除 "${row.name}" 吗？`, '删除确认', { type: 'warning' }).then(async () => {
      await fetchDeleteApp(row.id)
      ElMessage.success('删除成功')
      getData()
    })
  }

  function beforeIPAUpload(file: File) {
    if (!file.name.toLowerCase().endsWith('.ipa')) {
      ElMessage.error('只支持.ipa文件')
      return false
    }
    ipaUploading.value = true
    return true
  }

  function handleIPASuccess(res: any) {
    ipaUploading.value = false
    if (res.code !== 200) { ElMessage.error(res.msg || '解析失败'); return }
    const d = res.data
    if (d.name) form.value.name = d.name
    if (d.version) form.value.nickname = d.version
    if (d.identifier) form.value.diyname = d.identifier
    if (d.size) form.value.bt2a = d.size
    if (d.downloadURL) form.value.bt1a = d.downloadURL
    if (d.iconURL) form.value.image = d.iconURL
    ElMessage.success('IPA解析成功，信息已自动填充')
  }

  function handleIPAError() {
    ipaUploading.value = false
    ElMessage.error('上传失败')
  }

  function beforeImgUpload(file: File) {
    if (!file.type.startsWith('image/')) { ElMessage.error('只支持图片格式'); return false }
    if (file.size / 1024 / 1024 > 5) { ElMessage.error('图片不能超过5MB'); return false }
    return true
  }
</script>
