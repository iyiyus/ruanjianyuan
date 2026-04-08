<template>
  <div class="p-1">
    <ElCard shadow="never">
      <!-- 标题 -->
      <div class="flex items-center gap-3 pb-4 border-b border-g-300">
        <h2 class="text-xl font-medium">API接口文档</h2>
        <ElTag type="success" size="small">v1.0</ElTag>
      </div>

      <!-- 基础信息 -->
      <div class="mt-4 px-4 py-3 bg-g-100 rounded-lg text-sm text-g-600">
        接口基础地址：<span class="font-mono text-g-800">{{ baseUrl }}</span>
        &nbsp;|&nbsp; 黑名单接口无需 Token，直接访问
      </div>

      <!-- 接口列表 -->
      <div v-for="api in apiList" :key="api.path" class="mt-8">
        <!-- 接口标题 -->
        <div class="flex items-center gap-3 mb-4">
          <ElTag :type="api.method === 'GET' ? 'success' : 'warning'" size="default" class="font-mono">
            {{ api.method }}
          </ElTag>
          <span class="font-mono text-base font-medium text-g-800">{{ api.path }}</span>
          <span class="text-g-500 text-sm">{{ api.desc }}</span>
        </div>

        <!-- 请求参数 -->
        <div class="mb-4">
          <h4 class="text-sm font-medium text-g-700 mb-2">请求参数</h4>
          <ElTable v-if="api.params && api.params.length > 0" :data="api.params" border size="small">
            <ElTableColumn prop="name" label="参数名" width="140" />
            <ElTableColumn prop="type" label="类型" width="100" />
            <ElTableColumn label="必填" width="80">
              <template #default="{ row }">
                <ElTag :type="row.required ? 'danger' : 'info'" size="small">
                  {{ row.required ? '是' : '否' }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn prop="desc" label="说明" />
          </ElTable>
          <p v-else class="text-sm text-g-400">无需请求参数</p>
        </div>

        <!-- 返回示例 -->
        <h4 class="text-sm font-medium text-g-700 mb-2">返回示例</h4>
        <div class="relative">
          <pre class="bg-g-100 rounded-lg p-4 text-xs font-mono text-g-700 overflow-x-auto leading-6">{{ api.response }}</pre>
          <ElButton
            size="small" plain class="absolute top-2 right-2"
            @click="copyText(api.response)"
          >复制</ElButton>
        </div>

        <ElDivider v-if="api !== apiList[apiList.length - 1]" />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useClipboard } from '@vueuse/core'

  defineOptions({ name: 'ApiDoc' })

  const { copy } = useClipboard()
  const baseUrl = computed(() => window.location.origin)

  function copyText(text: string) {
    copy(text)
    ElMessage.success('已复制')
  }

  const apiList = [
    {
      method: 'GET',
      path: '/Blacklist',
      desc: '获取全量黑名单列表（无需认证）',
      params: [],
      response: `{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 2,
    "list": [
      {
        "udid": "00008030-001E78490133802E",
        "reason": "盗版破解",
        "addtime": 1652351345
      }
    ]
  }
}`
    },
    {
      method: 'GET',
      path: '/Blacklist/check?udid={udid}',
      desc: '查询单个UDID是否在黑名单（无需认证）',
      params: [
        { name: 'udid', type: 'string', required: true, desc: '设备UDID' }
      ],
      response: `// 在黑名单中
{
  "code": 200,
  "msg": "success",
  "data": {
    "blocked": true,
    "reason": "盗版破解"
  }
}

// 不在黑名单中
{
  "code": 200,
  "msg": "success",
  "data": {
    "blocked": false,
    "reason": ""
  }
}`
    },
    {
      method: 'GET',
      path: '/appstore',
      desc: '获取软件源数据（iOS Cydia/Sileo）',
      params: [
        { name: 'udid', type: 'string', required: false, desc: '设备UDID，用于黑名单检测和卡密验证' },
        { name: 'code', type: 'string', required: false, desc: '卡密，用于解锁付费内容' }
      ],
      response: `{
  "name": "软件源名称",
  "message": "公告内容",
  "identifier": "源标识符",
  "sourceURL": "https://example.com/appstore",
  "payURL": "购买卡密地址",
  "unlockURL": "解锁接口地址",
  "apps": [
    {
      "name": "微信",
      "version": "8.0.0",
      "type": "0",
      "versionDate": "2024-01-01T00:00:00+08:00",
      "versionDescription": "版本说明",
      "lock": "0",
      "downloadURL": "https://example.com/app.ipa",
      "isLanZouCloud": "0",
      "iconURL": "https://example.com/icon.png",
      "tintColor": "018084",
      "size": "118.2"
    }
  ]
}`
    }
  ]
</script>
