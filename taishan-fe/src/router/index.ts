import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'
import type { App } from 'vue'
import { Layout } from '@/utils/routerHelper'
import { NO_RESET_WHITE_LIST } from '@/constants'

export const constantRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard/analysis',
    name: 'Root',
    meta: {
      hidden: true
    }
  },
  {
    path: '/redirect',
    component: Layout,
    name: 'Redirect',
    children: [
      {
        path: '/redirect/:path(.*)',
        name: 'Redirect',
        component: () => import('@/views/Redirect/Redirect.vue'),
        meta: {}
      }
    ],
    meta: {
      hidden: true,
      noTagsView: true
    }
  },
  {
    path: '/404',
    component: () => import('@/views/Error/404.vue'),
    name: 'NoFind',
    meta: {
      hidden: true,
      title: '404',
      noTagsView: true
    }
  }
]

export const asyncRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/dashboard',
    component: Layout,
    redirect: '/dashboard/analysis',
    name: '首页',
    meta: {
      title: '首页',
      icon: 'ant-design:dashboard-filled'
    },
    children: [
      {
        path: 'analysis',
        component: () => import('@/views/Dashboard/Analysis.vue'),
        name: '首页',
        meta: {
          title: '首页',
          icon: 'ant-design:dashboard-filled'
        }
      }
      // {
      //   path: 'workplace',
      //   component: () => import('@/views/Dashboard/Workplace.vue'),
      //   name: 'Workplace',
      //   meta: {
      //     title: ('router.workplace'),
      //     noCache: true
      //   }
      // }
    ]
  },
  {
    path: '/plan',
    component: Layout,
    meta: {
      title: '压力测试',
      icon: 'ep:odometer'
    },
    name: 'Plan',
    children: [
      {
        path: 'list',
        component: () => import('@/views/Plan/PlanList.vue'),
        name: 'PlanList',
        meta: {
          title: '测试计划',
          icon: 'ep:odometer'
        }
      },
      {
        path: 'detail/:id',
        component: () => import('@/views/Plan/PlanDetail.vue'),
        name: 'PlanDetail',
        meta: {
          hidden: true,
          title: '计划详情',
          icon: 'ep:odometer'
        }
      }
    ]
  },
  {
    path: '/report',
    component: Layout,
    meta: {
      title: '测试报告',
      icon: 'ep:document-copy'
    },
    name: 'Report',
    children: [
      {
        path: 'list',
        component: () => import('@/views/Report/ReportList.vue'),
        name: 'ReportList',
        meta: {
          title: '测试报告',
          icon: 'ant-design:bar-chart-outlined'
        }
      },
      {
        path: 'detail/:id',
        component: () => import('@/views/Report/ReportDetail.vue'),
        name: 'ReportDetail',
        meta: {
          hidden: true,
          title: '报告详情',
          icon: 'ep:document'
        }
      }
    ]
  },
  {
    path: '/machine',
    component: Layout,
    meta: {
      title: '机器管理',
      icon: 'ep:odometer'
    },
    name: 'Machine',
    children: [
      {
        path: 'manager',
        component: () => import('@/views/Machine/MachineManager.vue'),
        name: 'MachineList',
        meta: {
          title: '机器管理',
          icon: 'ep:monitor'
        }
      }
    ]
  }
  // {
  //   path: '/error',
  //   component: Layout,
  //   redirect: '/error/404',
  //   name: 'Error',
  //   meta: {
  //     title: t('router.errorPage'),
  //     icon: 'ci:error',
  //     alwaysShow: true
  //   },
  //   children: [
  //     {
  //       path: '404-demo',
  //       component: () => import('@/views/Error/404.vue'),
  //       name: '404Demo',
  //       meta: {
  //         title: '404'
  //       }
  //     },
  //     {
  //       path: '403-demo',
  //       component: () => import('@/views/Error/403.vue'),
  //       name: '403Demo',
  //       meta: {
  //         title: '403'
  //       }
  //     },
  //     {
  //       path: '500-demo',
  //       component: () => import('@/views/Error/500.vue'),
  //       name: '500Demo',
  //       meta: {
  //         title: '500'
  //       }
  //     }
  //   ]
  // },
  // {
  //   path: '/authorization',
  //   component: Layout,
  //   redirect: '/authorization/user',
  //   name: 'Authorization',
  //   meta: {
  //     title: t('router.authorization'),
  //     icon: 'eos-icons:role-binding',
  //     alwaysShow: true
  //   },
  //   children: [
  //     {
  //       path: 'department',
  //       component: () => import('@/views/Authorization/Department/Department.vue'),
  //       name: 'Department',
  //       meta: {
  //         title: t('router.department')
  //       }
  //     },
  //     {
  //       path: 'user',
  //       component: () => import('@/views/Authorization/User/User.vue'),
  //       name: 'User',
  //       meta: {
  //         title: t('router.user')
  //       }
  //     },
  //     {
  //       path: 'menu',
  //       component: () => import('@/views/Authorization/Menu/Menu.vue'),
  //       name: 'Menu',
  //       meta: {
  //         title: t('router.menuManagement')
  //       }
  //     },
  //     {
  //       path: 'role',
  //       component: () => import('@/views/Authorization/Role/Role.vue'),
  //       name: 'Role',
  //       meta: {
  //         title: t('router.role')
  //       }
  //     }
  //   ]
  // }
]

const router = createRouter({
  history: createWebHashHistory(),
  strict: true,
  routes: constantRouterMap as RouteRecordRaw[],
  scrollBehavior: () => ({ left: 0, top: 0 })
})

export const resetRouter = (): void => {
  router.getRoutes().forEach((route) => {
    const { name } = route
    if (name && !NO_RESET_WHITE_LIST.includes(name as string)) {
      router.hasRoute(name) && router.removeRoute(name)
    }
  })
}

export const setupRouter = (app: App<Element>) => {
  app.use(router)
}

export default router
