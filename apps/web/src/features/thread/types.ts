import type { BoardId } from '../board'

/** 帖子共通信息 */
export interface ThreadCore {
  /** 帖子ID（由数据库自增生成） */
  id: string
  /** 作者 */
  // author: AuthorizationServices.UserInfo
  authorId: string
  /** 版块ID */
  boardId: BoardId
  /** 标题 */
  title: string
  /** 预览图链接 */
  thumbnail_link: string
  /** TAG */
  tag: string[]
  /** 创建时间 */
  createdAt: string
  /** 更新时间 */
  updatedAt: string
}

export enum ThreadType {
  ARTICLE = 'article',
  COMIC = 'comic'
}

/** 文章 */
export interface Article extends ThreadCore {
  /** 帖子类型 */
  threadType: ThreadType.ARTICLE
}

// TODO: 等做到漫画功能时具体确定剩余字段。
/** 漫画 */
export interface Comic extends ThreadCore {
  /** 帖子类型 */
  threadType: ThreadType.COMIC
}

export type Thread = Article | Comic

/** 帖子列表元素 */
export interface ThreadListItem {
  /** 帖子ID（经过编码的） */
  id: string
  /** 标题 */
  title: string
  /** 预览图链接 */
  thumbnailLink: string
  /** 更新时间 */
  updatedAt: string
}
