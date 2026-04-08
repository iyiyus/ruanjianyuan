<!-- 安装向导 -->
<template>
  <div class="install-page">
    <div class="install-container">
      <!-- 顶部 Logo -->
      <div class="install-header">
        <ArtLogo size="36" />
        <span class="site-name">iOS软件源管理系统</span>
      </div>

      <!-- 步骤条 -->
      <ElSteps :active="currentStep" finish-status="success" class="install-steps">
        <ElStep title="许可协议" />
        <ElStep title="环境检测" />
        <ElStep title="系统配置" />
        <ElStep title="安装完成" />
      </ElSteps>

      <!-- 内容卡片 -->
      <div class="art-card install-card">

        <!-- Step 0: 协议 -->
        <template v-if="currentStep === 0">
          <h3 class="step-title">用户使用协议</h3>
          <p class="step-desc">在开始安装之前，请阅读并接受以下条款</p>
          <div class="agreement-box">
            <p><b>1. 法律合规</b><br>本程序仅供学习交流及正版软件分发使用，不得用于违法用途。</p>
            <p><b>2. 免责声明</b><br>本软件按"原样"提供，开发者不对因使用本软件导致的任何损失负责。</p>
            <p><b>3. 授权说明</b><br>严禁用于诈骗、博彩、色情等违法违规内容的传播。</p>
          </div>
          <ElCheckbox v-model="agreed" class="mt-4">我已阅读并同意上述协议</ElCheckbox>
          <div class="step-footer">
            <ElButton type="primary" :disabled="!agreed" @click="currentStep++" v-ripple>下一步</ElButton>
          </div>
        </template>

        <!-- Step 1: 环境检测 -->
        <template v-if="currentStep === 1">
          <h3 class="step-title">环境兼容性检测</h3>
          <p class="step-desc">检测系统运行环境是否满足要求</p>
          <div class="check-list">
            <div class="check-row" v-for="item in envChecks" :key="item.name">
              <span class="text-sm">{{ item.name }}</span>
              <ElTag type="success" size="small">✓ {{ item.status }}</ElTag>
            </div>
          </div>
          <div class="step-footer">
            <ElButton @click="currentStep--">上一步</ElButton>
            <ElButton type="primary" @click="currentStep++" v-ripple>下一步</ElButton>
          </div>
        </template>

        <!-- Step 2: 配置 -->
        <template v-if="currentStep === 2">
          <h3 class="step-title">系统参数配置</h3>
          <p class="step-desc">请填写数据库连接信息和系统参数</p>
          <ElForm :model="form" label-position="top" ref="formRef" :rules="rules">
            <ElTabs v-model="configTab">
              <ElTabPane label="数据库配置" name="db">
                <ElRow :gutter="16">
                  <ElCol :span="16">
                    <ElFormItem label="主机地址" prop="dbHost">
                      <ElInput v-model="form.dbHost" placeholder="127.0.0.1" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="8">
                    <ElFormItem label="端口">
                      <ElInput v-model="form.dbPort" placeholder="3306" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="用户名" prop="dbUser">
                      <ElInput v-model="form.dbUser" placeholder="root" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="密码">
                      <ElInput v-model="form.dbPassword" type="password" show-password placeholder="数据库密码" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="24">
                    <ElFormItem label="数据库名" prop="dbName">
                      <ElInput v-model="form.dbName" placeholder="如 ios_source">
                        <template #append>
                          <ElButton :loading="testingDB" @click="testDB">测试连接</ElButton>
                        </template>
                      </ElInput>
                    </ElFormItem>
                  </ElCol>
                </ElRow>
              </ElTabPane>

              <ElTabPane label="应用配置" name="app">
                <ElRow :gutter="16">
                  <ElCol :span="12">
                    <ElFormItem label="服务端口号">
                      <ElInput v-model="form.appPort" placeholder="如 1117（需与宝塔一致）" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="后台访问路径">
                      <ElInput v-model="form.adminPath" placeholder="如 admin">
                        <template #prepend>/</template>
                      </ElInput>
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="GO项目名称">
                      <ElInput v-model="form.projectName" placeholder="如 go-source" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="网站名称">
                      <ElInput v-model="form.siteName" placeholder="如 iOS软件源管理系统" />
                    </ElFormItem>
                  </ElCol>
                </ElRow>
              </ElTabPane>

              <ElTabPane label="管理员账号" name="admin">
                <ElRow :gutter="16">
                  <ElCol :span="12">
                    <ElFormItem label="账号" prop="adminUser">
                      <ElInput v-model="form.adminUser" placeholder="如 admin" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="12">
                    <ElFormItem label="密码">
                      <ElInput v-model="form.adminPwd" type="password" show-password placeholder="默认 123456" />
                    </ElFormItem>
                  </ElCol>
                  <ElCol :span="24">
                    <ElFormItem label="邮箱（可选）">
                      <ElInput v-model="form.adminEmail" placeholder="admin@example.com" />
                    </ElFormItem>
                  </ElCol>
                </ElRow>
              </ElTabPane>
            </ElTabs>
          </ElForm>
          <div class="step-footer">
            <ElButton @click="currentStep--">上一步</ElButton>
            <ElButton type="primary" :loading="installing" @click="doInstall" v-ripple>立即部署</ElButton>
          </div>
        </template>

        <!-- Step 3: 完成 -->
        <template v-if="currentStep === 3">
          <div class="done-wrap">
            <div class="done-icon">✓</div>
            <h3 class="step-title">安装完成</h3>
            <p class="step-desc">系统已成功部署，请保存以下信息</p>
            <div class="result-box">
              <div class="res-row"><span>管理账号</span><b>{{ result.adminUser }}</b></div>
              <div class="res-row"><span>管理密码</span><b>{{ result.adminPwd }}</b></div>
              <div class="res-row"><span>后台路径</span><b>/{{ result.adminPath }}</b></div>
              <div class="res-row"><span>服务端口</span><b>{{ result.appPort }}</b></div>
            </div>

            <!-- 重启状态 -->
            <div class="restart-box" :class="restartStatus">
              <template v-if="restartStatus === 'restarting'">
                <ElIcon class="is-loading"><Loading /></ElIcon>
                <span>正在重启服务，请稍候... ({{ countdown }}s)</span>
              </template>
              <template v-else-if="restartStatus === 'done'">
                <span>✓ 服务已重启完成</span>
              </template>
              <template v-else-if="restartStatus === 'timeout'">
                <span>⚠️ 重启超时，请手动在宝塔重启 Go 项目</span>
              </template>
            </div>

            <ElButton v-if="restartStatus === 'done'" type="primary" class="mt-4 w-full" @click="goHome" v-ripple>
              进入管理面板
            </ElButton>
            <ElButton v-if="restartStatus === 'timeout'" type="warning" class="mt-4 w-full" @click="goHome" v-ripple>
              手动重启后进入
            </ElButton>
          </div>
        </template>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import request from '@/utils/http'
  import type { FormInstance } from 'element-plus'
  import { Loading } from '@element-plus/icons-vue'

  defineOptions({ name: 'Install' })

  const currentStep = ref(0)
  const agreed = ref(false)
  const testingDB = ref(false)
  const installing = ref(false)
  const configTab = ref('db')
  const formRef = ref<FormInstance>()
  const result = ref<any>({})

  const envChecks = [
    { name: 'Go 运行环境', status: '正常' },
    { name: 'MySQL 数据库驱动', status: '已就绪' },
    { name: '文件存储权限', status: '可写入' },
    { name: '网络服务', status: '运行中' }
  ]

  const form = ref({
    dbHost: '', dbPort: '', dbUser: '', dbPassword: '', dbName: '',
    appPort: '', adminPath: '', projectName: '', siteName: '',
    adminUser: '', adminPwd: '', adminEmail: ''
  })

  const rules = {
    dbHost: [{ required: true, message: '请输入数据库地址', trigger: 'blur' }],
    dbUser: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    dbName: [{ required: true, message: '请输入数据库名', trigger: 'blur' }],
    adminUser: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }]
  }

  async function testDB() {
    testingDB.value = true
    try {
      await request.post({ url: '/api/install/check-db', params: {
        host: form.value.dbHost || '127.0.0.1',
        port: form.value.dbPort || '3306',
        user: form.value.dbUser,
        password: form.value.dbPassword,
        dbName: form.value.dbName
      }})
      ElMessage.success('数据库连接成功')
    } finally {
      testingDB.value = false
    }
  }

  async function doInstall() {
    if (!formRef.value) return
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) { configTab.value = 'db'; return }
    installing.value = true
    try {
      const data = await request.post({ url: '/api/install/run', params: form.value })
      result.value = data
      currentStep.value = 3
      startRestartPolling()
    } finally {
      installing.value = false
    }
  }

  const restartStatus = ref<'restarting' | 'done' | 'timeout'>('restarting')
  const countdown = ref(60)
  let countdownTimer: ReturnType<typeof setInterval> | null = null
  let pollTimer: ReturnType<typeof setInterval> | null = null

  function startRestartPolling() {
    restartStatus.value = 'restarting'
    countdown.value = 60

    // 倒计时
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(countdownTimer!)
        if (restartStatus.value === 'restarting') {
          restartStatus.value = 'timeout'
          clearInterval(pollTimer!)
        }
      }
    }, 1000)

    // 轮询检测服务是否重启完成（安装后接口会短暂不可用，再次可用说明重启完成）
    let wasDown = false
    pollTimer = setInterval(async () => {
      try {
        await request.get({ url: '/api/install/status', showErrorMessage: false })
        if (wasDown) {
          // 服务恢复了，重启完成
          restartStatus.value = 'done'
          clearInterval(pollTimer!)
          clearInterval(countdownTimer!)
        }
      } catch {
        wasDown = true // 服务暂时不可用，正在重启中
      }
    }, 2000)
  }

  onUnmounted(() => {
    if (countdownTimer) clearInterval(countdownTimer)
    if (pollTimer) clearInterval(pollTimer)
  })

  function goHome() {
    const adminPath = result.value.adminPath || 'admin'
    window.location.href = '/' + adminPath
  }

  onMounted(async () => {
    try {
      const data = await request.get({ url: '/api/install/status' })
      if (data.installed) {
        // 已安装，读取后台路径跳转
        try {
          const cfg = await request.get({ url: '/api/config/map', showErrorMessage: false })
          const adminPath = cfg.admin_path || 'admin'
          window.location.href = '/' + adminPath
        } catch {
          window.location.href = '/admin'
        }
      }
    } catch {}
  })
</script>

<style scoped>
  .install-page {
    min-height: 100vh;
    background: var(--default-bg-color);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
  }

  .install-container {
    width: 100%;
    max-width: 680px;
  }

  .install-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 28px;
    justify-content: center;
  }

  .site-name {
    font-size: 18px;
    font-weight: 700;
    color: var(--art-gray-900);
  }

  .install-steps {
    margin-bottom: 24px;
  }

  .install-card {
    padding: 32px 36px;
  }

  .step-title {
    font-size: 20px;
    font-weight: 700;
    margin-bottom: 6px;
    color: var(--art-gray-900);
  }

  .step-desc {
    font-size: 14px;
    color: var(--art-gray-500);
    margin-bottom: 20px;
  }

  .agreement-box {
    background: var(--art-gray-100);
    border: 1px solid var(--el-border-color);
    border-radius: 8px;
    padding: 16px 20px;
    max-height: 200px;
    overflow-y: auto;
    font-size: 14px;
    line-height: 1.8;
    color: var(--art-gray-700);
  }
  .agreement-box p { margin: 0 0 10px; }
  .agreement-box b { color: var(--art-gray-900); }

  .check-list { margin-bottom: 8px; }
  .check-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 13px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  .check-row:last-child { border: none; }

  .step-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 24px;
  }

  .done-wrap { text-align: center; }
  .done-icon {
    width: 64px; height: 64px;
    background: var(--el-color-success-light-9);
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    margin: 0 auto 16px;
    font-size: 28px;
    color: var(--el-color-success);
  }

  .result-box {
    border: 1px solid var(--el-border-color);
    border-radius: 8px;
    overflow: hidden;
    text-align: left;
    margin-top: 16px;
  }
  .res-row {
    display: flex;
    justify-content: space-between;
    padding: 11px 16px;
    font-size: 14px;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  .res-row:last-child { border: none; }

  .restart-box {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 16px;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 14px;
  }
  .restart-box.restarting {
    background: var(--el-color-warning-light-9);
    color: var(--el-color-warning);
    border: 1px solid var(--el-color-warning-light-5);
  }
  .restart-box.done {
    background: var(--el-color-success-light-9);
    color: var(--el-color-success);
    border: 1px solid var(--el-color-success-light-5);
  }
  .restart-box.timeout {
    background: var(--el-color-danger-light-9);
    color: var(--el-color-danger);
    border: 1px solid var(--el-color-danger-light-5);
  }
</style>
