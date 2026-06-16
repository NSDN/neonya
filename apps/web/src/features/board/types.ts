import { PageType } from '@/shared/constants/enum'

export type BoardId = string

/** 版块 */
export interface Board {
  /** 版块 ID（兼路由名） */
  id: BoardId
  /** 版块名 */
  name: string
  /** 背景图片（图床地址） */
  background: string
  /** 画面类型 */
  pageType: PageType
  /** 排序锚点 */
  sortOrder: number
}
