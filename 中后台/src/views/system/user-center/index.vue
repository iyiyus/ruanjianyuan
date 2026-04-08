<!-- 个人中心页面 -->
<template>
  <div class="w-full h-full p-0 bg-transparent border-none shadow-none">
    <div class="relative flex-b mt-2.5 max-md:block max-md:mt-1">

      <!-- 左侧：头像卡片 -->
      <div class="w-80 mr-5 max-md:w-full max-md:mr-0">
        <div class="art-card-sm relative p-9 pb-6 overflow-hidden text-center">
          <img class="absolute top-0 left-0 w-full h-50 object-cover" src="@imgs/user/bg.webp" />
          <div class="relative z-10 mt-30 mx-auto w-20 h-20">
            <img
              :src="avatarSrc"
              class="w-20 h-20 mx-auto object-cover border-2 border-white rounded-full"
            />
            <ElUpload
              class="avatar-upload-btn"
              :action="uploadAction"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
              accept="image/*"
            >
              <div class="absolute bottom-0 right-0 w-6 h-6 bg-theme rounded-full flex-cc cursor-pointer">
                <el-icon class="text-white text-xs"><Camera /></el-icon>
              </div>
            </ElUpload>
          </div>
          <h2 class="mt-5 text-xl font-normal">{{ userInfo.userName }}</h2>
          <p class="mt-2 text-sm text-g-500">{{ userInfo.email }}</p>
        </div>
      </div>

      <!-- 右侧：表单 -->
      <div class="flex-1 overflow-hidden max-md:w-full max-md:mt-3.5">

        <!-- 基本信息 -->
        <div class="art-card-sm">
          <h1 class="p-4 text-xl font-normal border-b border-g-300">基本信息</h1>
          <ElForm :model="form" class="box-border p-5" label-position="top">
            <ElRow :gutter="16">
              <ElCol :span="12">
                <ElFormItem label="昵称">
                  <ElInput v-model="form.nickname" :disabled="!isEdit" />
                </ElFormItem>
              </ElCol>
              <ElCol :span="12">
                <ElFormItem label="邮箱">
                  <ElInput v-model="form.email" :disabled="!isEdit" />
                </ElFormItem>
              </ElCol>
            </ElRow>
            <div class="flex-c justify-end">
              <ElButton type="primary" v-ripple @click="handleProfileBtn" :loading="profileLoading">
                {{ isEdit ? '保存' : '编辑' }}
              </ElButton>
            </div>
          </ElForm>
        </div>

        <!-- 更改密码 -->
        <div class="art-card-sm my-5">
          <h1 class="p-4 text-xl font-normal border-b border-g-300">更改密码</h1>
          <ElForm :model="pwdForm" class="box-border p-5" label-position="top">
            <ElFormItem label="当前密码">
              <ElInput v-model="pwdForm.password" type="password" :disabled="!isEditPwd" show-password />
            </ElFormItem>
            <ElFormItem label="新密码">
              <ElInput v-model="pwdForm.newPassword" type="password" :disabled="!isEditPwd" show-password />
            </ElFormItem>
            <ElFormItem label="确认新密码">
              <ElInput v-model="pwdForm.confirmPassword" type="password" :disabled="!isEditPwd" show-password />
            </ElFormItem>
            <div class="flex-c justify-end">
              <ElButton type="primary" v-ripple @click="handlePwdBtn" :loading="pwdLoading">
                {{ isEditPwd ? '保存' : '编辑' }}
              </ElButton>
            </div>
          </ElForm>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import request from '@/utils/http'
  import { Camera } from '@element-plus/icons-vue'
  import defaultAvatar from '@/assets/images/user/avatar.webp'

  defineOptions({ name: 'UserCenter' })

  const userStore = useUserStore()
  const userInfo = computed(() => userStore.getUserInfo)

  const uploadAction = '/api/upload'
  const uploadHeaders = computed(() => ({ Authorization: userStore.accessToken }))

  const avatarSrc = computed(() => {
    const a = (userInfo.value as any).avatar
    return a && a.trim() ? a : defaultAvatar
  })

  // 基本信息
  const isEdit = ref(false)
  const profileLoading = ref(false)
  const form = reactive({ nickname: '', email: '' })

  // 初始化表单
  watchEffect(() => {
    form.nickname = userInfo.value.userName || ''
    form.email = userInfo.value.email || ''
  })

  async function handleProfileBtn() {
    if (!isEdit.value) { isEdit.value = true; return }
    profileLoading.value = true
    try {
      await request.post({ url: '/api/user/profile', params: { nickname: form.nickname, email: form.email } })
      ElMessage.success('保存成功')
      isEdit.value = false
    } finally {
      profileLoading.value = false
    }
  }

  // 头像上传
  function beforeAvatarUpload(file: File) {
    if (!file.type.startsWith('image/')) { ElMessage.error('只支持图片格式'); return false }
    if (file.size / 1024 / 1024 > 2) { ElMessage.error('图片不能超过2MB'); return false }
    return true
  }

  async function handleAvatarSuccess(res: any) {
    // ElUpload 的 on-success 回调拿到的是原始响应体
    const url = res?.data?.url || res?.url
    if (!url) { ElMessage.error('上传失败'); return }
    await request.post({ url: '/api/user/profile', params: { avatar: url } })
    const current = userStore.getUserInfo
    userStore.setUserInfo({ ...current, avatar: url } as any)
    ElMessage.success('头像更新成功')
  }

  // 修改密码
  const isEditPwd = ref(false)
  const pwdLoading = ref(false)
  const pwdForm = reactive({ password: '', newPassword: '', confirmPassword: '' })

  async function handlePwdBtn() {
    if (!isEditPwd.value) {
      pwdForm.password = ''; pwdForm.newPassword = ''; pwdForm.confirmPassword = ''
      isEditPwd.value = true; return
    }
    if (!pwdForm.password || !pwdForm.newPassword || !pwdForm.confirmPassword) {
      ElMessage.warning('请填写完整密码信息'); return
    }
    if (pwdForm.newPassword !== pwdForm.confirmPassword) {
      ElMessage.error('两次新密码不一致'); return
    }
    if (pwdForm.newPassword.length < 6) {
      ElMessage.error('新密码不能少于6位'); return
    }
    pwdLoading.value = true
    try {
      await request.post({ url: '/api/user/change-password', params: { oldPassword: pwdForm.password, newPassword: pwdForm.newPassword } })
      ElMessage.success('密码修改成功，请重新登录')
      isEditPwd.value = false
      setTimeout(() => useUserStore().logOut(), 1500)
    } finally {
      pwdLoading.value = false
    }
  }
</script>
