import type { Option, Result } from '@/shared/utils/rust'
import type { ThreadListItem } from '../types'
import { request } from '@/shared/services/request'
import { API_URLS } from '@/shared/constants'
import type { ApiError } from '@/shared/errors'

/**
 * 获取帖子列表。
 *
 * @param boardId - 版块 ID
 * @returns 帖子列表
 */
export async function getThreads(
  boardId: string
): Promise<Result<Option<Option<ThreadListItem>[]>, ApiError>> {
  return request(`${API_URLS.THREADS}?boardId=${encodeURIComponent(boardId)}`, {
    method: 'GET'
  })
}
