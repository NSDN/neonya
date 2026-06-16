import type { Result, Option } from '@/shared/utils/rust'
import type { ApiError } from '@/shared/errors'
import { request } from '@/shared/services/request'
import { API_URLS } from '@/shared/constants'
import type { Board } from '../types'

export async function getBoards(): Promise<
  Result<Option<Option<Board>[]>, ApiError>
> {
  return request({
    url: API_URLS.BOARDS,
    method: 'GET',
  })
}
