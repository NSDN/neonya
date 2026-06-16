import type { RouteRecordRaw } from 'vue-router'

import { ROUTE_NAMES, ROUTE_PATHS } from '@/shared/constants'

const boards: RouteRecordRaw[] = [
  {
    path: ROUTE_PATHS.BOARD,
    name: ROUTE_NAMES.BOARD,
    meta: { displayCreateThreadButton: true },
    component: () => import('@/features/board/views/BoardContent.vue'),
  },
]

export default boards
