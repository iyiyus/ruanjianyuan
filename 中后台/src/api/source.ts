import request from '@/utils/http'

// ===== App管理 =====
export function fetchAppList(params: any) {
  return request.get<any>({ url: '/api/apps', params })
}
export function fetchCreateApp(params: any) {
  return request.post<any>({ url: '/api/apps', params })
}
export function fetchUpdateApp(id: number, params: any) {
  return request.put<any>({ url: `/api/apps/${id}`, params })
}
export function fetchDeleteApp(id: number) {
  return request.del<any>({ url: `/api/apps/${id}` })
}

// ===== 卡密管理 =====
export function fetchKamiList(params: any) {
  return request.get<any>({ url: '/api/kami', params })
}
export function fetchGenerateKami(params: { count: number; kmyp: number }) {
  return request.post<any>({ url: '/api/kami/generate', params })
}
export function fetchDeleteKami(id: number) {
  return request.del<any>({ url: `/api/kami/${id}` })
}

// ===== 黑名单 =====
export function fetchBlackList(params: any) {
  return request.get<any>({ url: '/api/black', params })
}
export function fetchCreateBlack(params: { udid: string; reason?: string }) {
  return request.post<any>({ url: '/api/black', params })
}
export function fetchDeleteBlack(id: number) {
  return request.del<any>({ url: `/api/black/${id}` })
}

// ===== 监控 =====
export function fetchMonitorList(params: any) {
  return request.get<any>({ url: '/api/monitor', params })
}
export function fetchDeleteMonitor(id: number) {
  return request.del<any>({ url: `/api/monitor/${id}` })
}
export function fetchClearMonitor() {
  return request.del<any>({ url: '/api/monitor' })
}

// ===== 配置 =====
export function fetchConfigMap() {
  return request.get<any>({ url: '/api/config/map' })
}
export function fetchUpdateConfig(configs: Record<string, string>) {
  return request.put<any>({ url: '/api/config', params: { configs } })
}

// ===== 仪表盘 =====
export function fetchDashboard() {
  return request.get<any>({ url: '/api/dashboard' })
}

// ===== 网站配置（公开，无需登录）=====
export function fetchSiteConfig() {
  return request.get<any>({ url: '/api/config/map' })
}

// ===== 上传 =====
export function uploadFile(formData: FormData) {
  return request.post<any>({ url: '/api/upload', params: formData })
}
export function uploadIPA(formData: FormData) {
  return request.post<any>({ url: '/api/upload/ipa', params: formData })
}

export function batchDeleteKami(ids: number[]) {
  return request.del<any>({ url: '/api/kami', params: { ids } })
}

// ===== 软件源搬运 =====
export function fetchRemoteSource(params: { url: string; udid?: string }) {
  return request.post<any>({ url: '/api/source/fetch', params })
}
export function importApps(apps: any[]) {
  return request.post<any>({ url: '/api/source/import', params: { apps } })
}
