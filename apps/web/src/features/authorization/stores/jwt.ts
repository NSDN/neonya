import type { AuthorizationServices } from '../types'

import { storage } from '@/shared/services/storage'
import { STORAGE_KEYS, STORE_ID } from '@/shared/constants'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { Option } from '@/shared/utils/rust'

/** @description JWT 状态库 */
const useJWTStore = defineStore(STORE_ID.JWT, () => {
  /** @description JWT 令牌 */
  const jwt = ref<Option<AuthorizationServices.Token>>(
    storage.get<AuthorizationServices.Token>(STORAGE_KEYS.TOKEN).match({
      ok: value => value,
      err: error => {
        error.notify()
        return Option.none()
      }
    })
  )

  /**
   * @description 设置令牌
   * @param token 从后端获取到的令牌
   */
  const setJWT = (token: Option<AuthorizationServices.Token>) => {
    jwt.value = token

    // 在浏览器的 storage 中保存 token。
    token.match({
      some: value => storage.set(STORAGE_KEYS.TOKEN, value),
      none: () => storage.remove(STORAGE_KEYS.TOKEN)
    })
  }

  return { jwt, setJWT }
})

export default useJWTStore
