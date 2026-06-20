import type { AuthorizationServices } from '../types'

import { API_URLS } from '@/shared/constants'
import { HTTPMethods } from '@/shared/constants/enum'
import type { ApiError } from '@/shared/errors'
import { request } from '@/shared/services/request'
import type { Option, Result } from '@/shared/utils/rust'
import { hashSecreate } from '@/shared/utils/useful'

export async function login({
  username,
  password
}: AuthorizationServices.LoginInfo): Promise<
  Result<Option<AuthorizationServices.Token>, ApiError>
> {
  return request(API_URLS.LOGIN, {
    method: HTTPMethods.POST,
    body: JSON.stringify({
      username,
      password: password && hashSecreate(password)
    })
  })
}
