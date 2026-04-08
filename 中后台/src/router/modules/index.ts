import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { sourceRoutes, kamiRoutes, securityRoutes, apiDocRoute, importSourceRoute } from './source'

const userCenterRoute: AppRouteRecord = {
  path: '/system/user-center',
  name: 'UserCenter',
  component: '/system/user-center',
  meta: {
    title: 'menus.system.userCenter',
    isHide: true,
    keepAlive: true,
    isHideTab: true
  }
}

export const routeModules: AppRouteRecord[] = [
  dashboardRoutes,
  ...sourceRoutes,
  ...kamiRoutes,
  ...securityRoutes,
  apiDocRoute,
  importSourceRoute,
  userCenterRoute
]
