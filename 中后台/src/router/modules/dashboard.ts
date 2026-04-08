import { AppRouteRecord } from '@/types/router'

export const dashboardRoutes: AppRouteRecord = {
  name: 'Console',
  path: '/dashboard/console',
  component: '/dashboard/console',
  meta: {
    title: 'menus.dashboard.console',
    icon: 'ri:home-smile-2-line',
    keepAlive: false,
    fixedTab: true,
    roles: ['R_SUPER', 'R_ADMIN']
  }
}
