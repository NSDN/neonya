import type { AuthorizationServices } from '../types'

import { API_URLS } from '@/shared/constants'
import { HTTPMethods } from '@/shared/constants/enum'
import type { ApiError } from '@/shared/errors'
import { request } from '@/shared/services/request'
import type { Option, Result } from '@/shared/utils/rust'

export async function fetchUserInfo(
  token: string
): Promise<Result<Option<AuthorizationServices.UserInfo>, ApiError>> {
  return request(API_URLS.USER_INFO, {
    method: HTTPMethods.GET,
    headers: {
      Authorization: `Bearer ${token}`
    }
  })
}
