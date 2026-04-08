<template>
  <div class="art-full-height" style="overflow-y:auto">
    <ElCard shadow="never" class="mb-4">
      <div class="flex items-center gap-3">
        <span class="text-sm text-g-600 flex-shrink-0">软件源地址</span>
        <ElInput v-model="form.url" placeholder="如 https://app.rrovo.cn/appstore" clearable
          @keyup.enter="handleFetch" style="flex:2" />
        <span class="text-sm text-g-600 flex-shrink-0">UDID</span>
        <ElInput v-model="form.udid" placeholder="可选，用于解锁付费内容" clearable style="flex:1" />
        <ElButton type="primary" :loading="fetching" @click="handleFetch" v-ripple>获取数据</ElButton>
      </div>
      <ElAlert class="mt-3" type="warning" :closable="false" show-icon>
        <template #title>
          仅支持未加密软件源。该软件源无法搬运。
        </template>
      </ElAlert>

      <!-- 历史记录 -->
      <div v-if="history.length > 0" class="flex items-center gap-2 mt-1">
        <span class="text-xs text-g-400">历史：</span>
        <ElTag v-for="url in history" :key="url" size="small" class="cursor-pointer" @click="loadHistory(url)">
          {{ url.length > 40 ? url.slice(0, 40) + '...' : url }}
        </ElTag>
        <ElButton link size="small" type="danger" @click="clearHistory">清空</ElButton>
      </div>
    </ElCard>

    <!-- 应用列表 -->
    <ElCard v-if="apps.length > 0" shadow="never">
      <div class="flex items-center gap-3 mb-4">
        <span class="text-sm text-g-600">共 {{ apps.length }} 个应用，已选 {{ selectedApps.length }} 个</span>
        <ElButton size="small" @click="selectAll">全选</ElButton>
        <ElButton size="small" @click="unselectAll">取消全选</ElButton>
        <ElButton type="primary" size="small" :loading="importing" :disabled="selectedApps.length === 0"
          @click="handleImport" v-ripple>
          导入选中 ({{ selectedApps.length }})
        </ElButton>
      </div>

      <div class="apps-grid">
        <div v-for="(app, idx) in pagedApps" :key="idx"
          class="app-item"
          :class="{ 'is-selected': isSelected(pageStart + idx) }"
          @click="toggleSelect(pageStart + idx)">
          <div class="flex items-start gap-3">
            <div class="relative flex-shrink-0">
              <ElCheckbox :model-value="isSelected(pageStart + idx)" @click.stop="toggleSelect(pageStart + idx)" />
            </div>
            <div style="width:44px;height:44px;border-radius:10px;background:#f5f5f5;overflow:hidden;flex-shrink:0">
              <img :src="app.iconURL" style="width:44px;height:44px;object-fit:cover"
                @error="(e: any) => e.target.src='/assets/img/avatar.png'" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span class="font-medium text-sm text-g-800 truncate">{{ app.name }}</span>
                <ElTag size="small" type="info">v{{ app.version }}</ElTag>
                <ElTag v-if="app.lock === '1'" size="small" type="danger">加锁</ElTag>
              </div>
              <p class="text-xs text-g-500 line-clamp-2">{{ app.versionDescription || '暂无说明' }}</p>
              <p class="text-xs text-g-400 mt-1">大小：{{ app.size ? app.size + ' MB' : '-' }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end mt-4">
        <ElPagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="apps.length"
          layout="total, prev, pager, next"
          background
        />
      </div>
    </ElCard>

    <ElEmpty v-else-if="fetched" description="没有获取到应用数据" />
  </div>
</template>

<script setup lang="ts">
  import { fetchRemoteSource, importApps } from '@/api/source'
  import { ElMessageBox } from 'element-plus'

  defineOptions({ name: 'ImportSource' })

  const form = ref({ url: '', udid: '' })
  const apps = ref<any[]>([])
  const selectedIndexes = ref<Set<number>>(new Set())
  const fetching = ref(false)
  const importing = ref(false)
  const fetched = ref(false)
  const history = ref<string[]>(JSON.parse(localStorage.getItem('sourceHistory') || '[]'))

  const currentPage = ref(1)
  const pageSize = 20
  const pageStart = computed(() => (currentPage.value - 1) * pageSize)
  const pagedApps = computed(() => apps.value.slice(pageStart.value, pageStart.value + pageSize))

  const selectedApps = computed(() => apps.value.filter((_, i) => selectedIndexes.value.has(i)))

  function isSelected(idx: number) { return selectedIndexes.value.has(idx) }
  function toggleSelect(idx: number) {
    const s = new Set(selectedIndexes.value)
    s.has(idx) ? s.delete(idx) : s.add(idx)
    selectedIndexes.value = s
  }
  function selectAll() { selectedIndexes.value = new Set(apps.value.map((_, i) => i)) }
  function unselectAll() { selectedIndexes.value = new Set() }

  function saveHistory(url: string) {
    let h = JSON.parse(localStorage.getItem('sourceHistory') || '[]') as string[]
    h = [url, ...h.filter(u => u !== url)].slice(0, 8)
    localStorage.setItem('sourceHistory', JSON.stringify(h))
    history.value = h
  }

  function loadHistory(url: string) {
    form.value.url = url
    handleFetch()
  }

  function clearHistory() {
    localStorage.removeItem('sourceHistory')
    history.value = []
  }

  async function handleFetch() {
    if (!form.value.url.trim()) { ElMessage.warning('请输入软件源地址'); return }
    fetching.value = true
    fetched.value = false
    apps.value = []
    selectedIndexes.value = new Set()
    try {
      const data = await fetchRemoteSource({ url: form.value.url.trim(), udid: form.value.udid.trim() || undefined })
      apps.value = data.apps || []
      fetched.value = true
      currentPage.value = 1
      saveHistory(form.value.url.trim())
      if (apps.value.length > 0) {
        selectAll()
        ElMessage.success(`获取到 ${apps.value.length} 个应用`)
      }
    } catch {
      fetched.value = true
    } finally {
      fetching.value = false
    }
  }

  async function handleImport() {
    if (selectedApps.value.length === 0) return
    await ElMessageBox.confirm(`确定导入选中的 ${selectedApps.value.length} 个应用吗？`, '导入确认', { type: 'info' })
    importing.value = true
    try {
      await importApps(selectedApps.value)
      ElMessage.success(`导入成功`)
      unselectAll()
    } finally {
      importing.value = false
    }
  }
</script>

<style scoped>
  .apps-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 12px;
  }
  .app-item {
    border: 1px solid var(--el-border-color);
    border-radius: 8px;
    padding: 12px;
    cursor: pointer;
    transition: all 0.2s;
  }
  .app-item:hover { border-color: var(--el-color-primary); }
  .app-item.is-selected { border-color: var(--el-color-primary); background: var(--el-color-primary-light-9); }
  .line-clamp-2 { display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
</style>
