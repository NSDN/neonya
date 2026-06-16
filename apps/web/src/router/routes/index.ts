import type { RouteRecordRaw } from 'vue-router'

import {
  BOARD_ROOT_REDIRECT,
  ROUTE_NAMES,
  ROUTE_PATHS
} from '@/shared/constants'
import boards from './boards'
import authorization from './authorization'

export const routes: RouteRecordRaw[] = [
  {
    path: ROUTE_PATHS.HOME,
    name: ROUTE_NAMES.HOME,
    redirect: BOARD_ROOT_REDIRECT
  },

  ...boards,
  ...authorization
]
