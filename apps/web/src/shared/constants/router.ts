import type { RouteRecordRedirectOption } from 'vue-router'

export const ROUTE_PATHS = {
  /** 首页 */
  HOME: '/',
  /** 版块 */
  BOARD: '/boards/:id',
  /** 授权 */
  AUTHORIZATION: '/authorization',
  /** 登入 */
  LOGIN: '/authorization/login',
  /** 注册 */
  REGISTER: '/authorization/register',
  /** 404 */
  NOT_FOUND: '/:pathMatch(.*)*'
} as const

export const ROUTE_NAMES = {
  /** 首页 */
  HOME: 'Home',
  /** 版块 */
  BOARD: 'Board',
  /** 登入 */
  LOGIN: 'Login',
  /** 注册 */
  REGISTER: 'Register',
  /** 404 */
  NOT_FOUND: 'NotFound'
} as const

/** 版块根目录重定向 */
export const BOARD_ROOT_REDIRECT: RouteRecordRedirectOption = {
  name: ROUTE_NAMES.BOARD,
  params: { id: 'localization' }
} as const
