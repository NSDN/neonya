import { Result } from '../utils/rust'
import { ApiError } from '../errors'
import { ErrorLevel } from '../constants'
import { HTTP_DEFAULT_CONFIG } from '../constants/http'
import { optionizeDeep } from '../utils/useful'

export async function request<Response>(
  path: string,
  options?: RequestInit
): Promise<Result<Response, ApiError>> {
  try {
    const base = HTTP_DEFAULT_CONFIG.BASE_URL
    const url = new URL(path, base.endsWith('/') ? base : `${base}/`)

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': HTTP_DEFAULT_CONFIG.CONTENT_TYPE,
        ...options?.headers
      },
      signal: AbortSignal.timeout(HTTP_DEFAULT_CONFIG.TIMEOUT)
    })

    if (!response.ok) {
      const err = new ApiError({
        level: ErrorLevel.Error,
        message: `[API]: 请求失败，状态码： ${response.status}`
      })
      return Result.err(err)
    }

    const data = await response.json()
    const optionized = optionizeDeep(data) as Response

    return Result.ok(optionized)
  } catch (error) {
    let message = '未知类型的错误'

    if (error instanceof Error) {
      if (error.name === 'AbortError' || error.name === 'TimeoutError') {
        message = '[API]: 请求超时'
      } else {
        message = `[API]: 请求失败，响应为： ${error.message}`
      }
    } else {
      message = String(error)
    }

    const err = new ApiError({ level: ErrorLevel.Error, message })
    return Result.err(err)
  }
}
