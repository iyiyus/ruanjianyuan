import { AppRouteRecord } from '@/types/router'

export const sourceRoutes: AppRouteRecord[] = [
  {
    path: '/source/apps',
    name: 'Apps',
    component: '/source/apps',
    meta: {
      title: 'App管理',
      icon: 'ri:app-store-line',
      keepAlive: true,
      roles: ['R_SUPER', 'R_ADMIN']
    }
  },
  {
    path: '/source/config',
    name: 'SourceConfig',
    component: '/source/config',
    meta: {
      title: '系统配置',
      icon: 'ri:settings-3-line',
      keepAlive: true,
      roles: ['R_SUPER']
    }
  }
]

export const kamiRoutes: AppRouteRecord[] = [
  {
    path: '/kami/list',
    name: 'KamiList',
    component: '/kami/list',
    meta: {
      title: '卡密管理',
      icon: 'ri:key-2-line',
      keepAlive: true,
      roles: ['R_SUPER']
    }
  }
]

export const securityRoutes: AppRouteRecord[] = [
  {
    path: '/security/black',
    name: 'Black',
    component: '/security/black',
    meta: {
      title: 'UDID黑名单',
      icon: 'ri:forbid-line',
      keepAlive: true,
      roles: ['R_SUPER']
    }
  },
  {
    path: '/security/monitor',
    name: 'Monitor',
    component: '/security/monitor',
    meta: {
      title: 'UDID监控',
      icon: 'ri:eye-line',
      keepAlive: true,
      roles: ['R_SUPER']
    }
  }
]

export const apiDocRoute: AppRouteRecord = {
  path: '/api-doc',
  name: 'ApiDoc',
  component: '/api-doc/index',
  meta: {
    title: '接口文档',
    icon: 'ri:file-list-3-line',
    keepAlive: true,
    roles: ['R_SUPER']
  }
}

export const importSourceRoute: AppRouteRecord = {
  path: '/import-source',
  name: 'ImportSource',
  component: '/import-source/index',
  meta: {
    title: '软件源搬运',
    icon: 'ri:download-cloud-line',
    keepAlive: true,
    roles: ['R_SUPER']
  }
}
