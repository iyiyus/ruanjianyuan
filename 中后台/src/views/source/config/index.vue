<template>
  <div class="art-full-height" style="overflow-y:auto">
    <ElCard shadow="never" v-loading="loading">
      <ElTabs v-model="activeTab">

        <!-- 软件源配置 -->
        <ElTabPane label="软件源配置" name="source">
          <ElForm :model="sourceForm" label-width="140px" style="max-width: 700px; margin-top: 16px">
            <ElFormItem label="源名称"><ElInput v-model="sourceForm.name" /></ElFormItem>
            <ElFormItem label="公告内容"><ElInput v-model="sourceForm.message" type="textarea" :rows="4" /></ElFormItem>
            <ElFormItem label="源识别标符"><ElInput v-model="sourceForm.identifier" /></ElFormItem>
            <ElFormItem label="软件来源地址"><ElInput v-model="sourceForm.sourceURL" /></ElFormItem>
            <ElFormItem label="源图标URL"><ElInput v-model="sourceForm.sourceicon" /></ElFormItem>
            <ElFormItem label="解锁发卡地址"><ElInput v-model="sourceForm.payURL" /></ElFormItem>
            <ElFormItem label="解锁接口地址"><ElInput v-model="sourceForm.unlockURL" /></ElFormItem>
            <ElFormItem label="软件源加密">
              <ElSwitch v-model="sourceForm.opencry" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="自动拉黑添加者">
              <ElSwitch v-model="sourceForm.openblack" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem label="自动拉黑破解者">
              <ElSwitch v-model="sourceForm.openblack2" active-value="1" inactive-value="0" />
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" :loading="submitLoading" @click="handleSaveSource" v-ripple>保存配置</ElButton>
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <!-- 网站配置 -->
        <ElTabPane label="网站配置" name="site">
          <ElForm :model="siteForm" label-width="140px" style="max-width: 700px; margin-top: 16px">
            <ElFormItem label="网站名称">
              <ElInput v-model="siteForm.site_name" placeholder="请输入网站名称" />
            </ElFormItem>
            <ElFormItem label="GO项目名称">
              <ElInput v-model="siteForm.go_project_name" placeholder="Go后端项目名称" />
            </ElFormItem>
            <ElFormItem label="网站LOGO">
              <div class="flex items-center gap-4">
                <ElUpload class="logo-uploader" :action="uploadAction" :headers="uploadHeaders"
                  :show-file-list="false" :on-success="handleLogoSuccess" :before-upload="beforeImgUpload" accept="image/*">
                  <img v-if="siteForm.site_logo" :src="siteForm.site_logo" class="logo-preview" />
                  <div v-else class="logo-placeholder">
                    <el-icon class="text-2xl text-g-400"><Plus /></el-icon>
                    <p class="text-xs text-g-400 mt-1">点击上传</p>
                  </div>
                </ElUpload>
                <ElButton v-if="siteForm.site_logo" link type="danger" size="small" @click="siteForm.site_logo = ''">移除</ElButton>
              </div>
            </ElFormItem>
            <ElFormItem>
              <ElButton type="primary" :loading="submitLoading" @click="handleSaveSite" v-ripple>保存配置</ElButton>
            </ElFormItem>
          </ElForm>
        </ElTabPane>

        <!-- 云存储配置 -->
        <ElTabPane label="云存储" name="storage">
          <ElForm :model="storageForm" label-width="160px" style="max-width: 750px; margin-top: 16px">
            <ElFormItem label="存储驱动">
              <ElSelect v-model="storageForm.storage_driver" style="width:200px">
                <ElOption label="本地存储" value="local" />
                <ElOption label="七牛云" value="qiniu" />
                <ElOption label="腾讯云 COS" value="tencent" />
                <ElOption label="阿里云 OSS" value="aliyun" />
                <ElOption label="WebDAV" value="webdav" />
              </ElSelect>
            </ElFormItem>

            <!-- 七牛云 -->
            <template v-if="storageForm.storage_driver === 'qiniu'">
              <ElDivider>七牛云配置</ElDivider>
              <ElFormItem label="AccessKey"><ElInput v-model="storageForm.qiniu_access_key" /></ElFormItem>
              <ElFormItem label="SecretKey"><ElInput v-model="storageForm.qiniu_secret_key" show-password /></ElFormItem>
              <ElFormItem label="Bucket"><ElInput v-model="storageForm.qiniu_bucket" /></ElFormItem>
              <ElFormItem label="访问域名"><ElInput v-model="storageForm.qiniu_domain" placeholder="https://cdn.example.com" /></ElFormItem>
              <ElFormItem label="存储区域">
                <ElSelect v-model="storageForm.qiniu_zone" style="width:200px">
                  <ElOption label="华东 (z0)" value="z0" />
                  <ElOption label="华北 (z1)" value="z1" />
                  <ElOption label="华南 (z2)" value="z2" />
                  <ElOption label="北美 (na0)" value="na0" />
                  <ElOption label="东南亚 (as0)" value="as0" />
                </ElSelect>
              </ElFormItem>
            </template>

            <!-- 腾讯云 -->
            <template v-if="storageForm.storage_driver === 'tencent'">
              <ElDivider>腾讯云 COS 配置</ElDivider>
              <ElFormItem label="SecretId"><ElInput v-model="storageForm.tencent_secret_id" /></ElFormItem>
              <ElFormItem label="SecretKey"><ElInput v-model="storageForm.tencent_secret_key" show-password /></ElFormItem>
              <ElFormItem label="Bucket"><ElInput v-model="storageForm.tencent_bucket" placeholder="mybucket-1250000000" /></ElFormItem>
              <ElFormItem label="Region"><ElInput v-model="storageForm.tencent_region" placeholder="ap-guangzhou" /></ElFormItem>
              <ElFormItem label="自定义域名"><ElInput v-model="storageForm.tencent_domain" placeholder="留空使用默认域名" /></ElFormItem>
            </template>

            <!-- 阿里云 -->
            <template v-if="storageForm.storage_driver === 'aliyun'">
              <ElDivider>阿里云 OSS 配置</ElDivider>
              <ElFormItem label="AccessKey"><ElInput v-model="storageForm.aliyun_access_key" /></ElFormItem>
              <ElFormItem label="SecretKey"><ElInput v-model="storageForm.aliyun_secret_key" show-password /></ElFormItem>
              <ElFormItem label="Bucket"><ElInput v-model="storageForm.aliyun_bucket" /></ElFormItem>
              <ElFormItem label="Endpoint"><ElInput v-model="storageForm.aliyun_endpoint" placeholder="oss-cn-hangzhou.aliyuncs.com" /></ElFormItem>
              <ElFormItem label="自定义域名"><ElInput v-model="storageForm.aliyun_domain" placeholder="留空使用默认域名" /></ElFormItem>
            </template>

            <!-- WebDAV -->
            <template v-if="storageForm.storage_driver === 'webdav'">
              <ElDivider>WebDAV 配置</ElDivider>
              <ElFormItem label="WebDAV地址"><ElInput v-model="storageForm.webdav_url" placeholder="https://dav.example.com/files/" /></ElFormItem>
              <ElFormItem label="用户名"><ElInput v-model="storageForm.webdav_username" /></ElFormItem>
              <ElFormItem label="密码"><ElInput v-model="storageForm.webdav_password" show-password /></ElFormItem>
              <ElFormItem label="访问域名"><ElInput v-model="storageForm.webdav_domain" placeholder="留空使用WebDAV地址" /></ElFormItem>
            </template>

            <ElFormItem>
              <ElButton type="primary" :loading="submitLoading" @click="handleSaveStorage" v-ripple>保存配置</ElButton>
            </ElFormItem>
          </ElForm>
        </ElTabPane>

      </ElTabs>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { fetchConfigMap, fetchUpdateConfig } from '@/api/source'
  import { useUserStore } from '@/store/modules/user'
  import { Plus } from '@element-plus/icons-vue'

  defineOptions({ name: 'SourceConfig' })

  const loading = ref(false)
  const submitLoading = ref(false)
  const activeTab = ref('source')

  const userStore = useUserStore()
  const uploadAction = '/api/upload'
  const uploadHeaders = computed(() => ({ Authorization: userStore.accessToken }))

  const sourceForm = ref({
    name: '', message: '', identifier: '', sourceURL: '',
    sourceicon: '', payURL: '', unlockURL: '',
    opencry: '0', openblack: '0', openblack2: '0'
  })

  const siteForm = ref({ site_name: '', site_logo: '', go_project_name: '' })

  const storageForm = ref({
    storage_driver: 'local',
    qiniu_access_key: '', qiniu_secret_key: '', qiniu_bucket: '', qiniu_domain: '', qiniu_zone: 'z0',
    tencent_secret_id: '', tencent_secret_key: '', tencent_bucket: '', tencent_region: '', tencent_domain: '',
    aliyun_access_key: '', aliyun_secret_key: '', aliyun_bucket: '', aliyun_endpoint: '', aliyun_domain: '',
    webdav_url: '', webdav_username: '', webdav_password: '', webdav_domain: ''
  })

  async function loadConfig() {
    loading.value = true
    try {
      const data = await fetchConfigMap()
      ;[sourceForm.value, siteForm.value, storageForm.value].forEach((form) => {
        Object.keys(form).forEach((key) => {
          if (data[key] !== undefined) (form as any)[key] = data[key]
        })
      })
    } finally {
      loading.value = false
    }
  }

  async function handleSaveSource() {
    submitLoading.value = true
    try { await fetchUpdateConfig(sourceForm.value as any); ElMessage.success('保存成功') }
    finally { submitLoading.value = false }
  }

  async function handleSaveSite() {
    submitLoading.value = true
    try { await fetchUpdateConfig(siteForm.value as any); ElMessage.success('保存成功') }
    finally { submitLoading.value = false }
  }

  async function handleSaveStorage() {
    submitLoading.value = true
    try { await fetchUpdateConfig(storageForm.value as any); ElMessage.success('保存成功') }
    finally { submitLoading.value = false }
  }

  function handleLogoSuccess(res: any) {
    if (res.code === 200) { siteForm.value.site_logo = res.data.url; ElMessage.success('上传成功') }
    else ElMessage.error(res.msg || '上传失败')
  }

  function beforeImgUpload(file: File) {
    const ok = file.type.startsWith('image/') && file.size / 1024 / 1024 < 2
    if (!file.type.startsWith('image/')) ElMessage.error('只能上传图片')
    if (file.size / 1024 / 1024 >= 2) ElMessage.error('图片不能超过2MB')
    return ok
  }

  onMounted(loadConfig)
</script>

<style scoped>
  .logo-uploader :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 8px;
    cursor: pointer;
    overflow: hidden;
    transition: border-color 0.3s;
  }
  .logo-uploader :deep(.el-upload:hover) { border-color: var(--el-color-primary); }
  .logo-preview { width: 100px; height: 100px; object-fit: contain; display: block; }
  .logo-placeholder { width: 100px; height: 100px; display: flex; flex-direction: column; align-items: center; justify-content: center; }
</style>
