import type { AuthorizationServices } from '../types'

import { useRouter } from 'vue-router'
import useJWTStore from '../stores/jwt'
import { useUserStore } from '@/stores/user'
import { reactive } from 'vue'
import {
  getMessage,
  MessageKeys,
  updateObjectValue
} from '@/shared/utils/useful'
import { useLoginRedirectStore } from '../stores/loginRedirect'
import { ROUTE_NAMES, ErrorLevel } from '@/shared/constants'
import { ResponsError } from '@/shared/errors'
import { login } from '../apis/login'
import { fetchUserInfo } from '../apis/userInfo'
import { Option } from '@/shared/utils/rust'

export function useLogin() {
  const userStore = useUserStore()
  const jwtStore = useJWTStore()
  const router = useRouter()

  /** @description 登入信息 */
  const loginInfo = reactive<AuthorizationServices.LoginInfo>({
    username: '',
    password: ''
  })

  const formError = reactive<AuthorizationServices.LoginInfoError>({
    username: '',
    password: ''
  })

  /** @description 初始化用户信息。 */
  const initUserInfo = async () => {
    const token = jwtStore.jwt

    if (token.isNone()) {
      return
    }

    const userInfo = await fetchUserInfo(token.unwrap().accessToken)

    userInfo
      .andThen(info =>
        info.okOrElse(
          () =>
            new ResponsError({
              level: ErrorLevel.Error,
              message: '[Login]: 未获取到用户信息。'
            })
        )
      )
      .match({
        ok: info => userStore.setUserInfo(info),
        err: error => error.notify()
      })
  }

  /** @description 登入 */
  const executeLogin = async () => {
    // 表单验证
    const error = validateLoginForm(loginInfo)
    updateObjectValue(formError, error)
    const hasError = Object.values(formError).some(value => !!value)

    // 表单验证失败则不执行登入
    if (hasError) {
      return
    }

    const token = (await login(loginInfo)).andThen(token =>
      token.okOrElse(
        () =>
          new ResponsError({
            level: ErrorLevel.Error,
            message: '[Login]: 未获取到令牌。'
          })
      )
    )

    if (token.isErr()) {
      token.error.notify()
      return
    }

    jwtStore.setJWT(Option.some(token.value))
    await initUserInfo()

    if (!userStore.loggedIn) {
      return
    }

    // 登入成功则跳转到指定的重定向画面或首页
    const loginRedirect = useLoginRedirectStore()
    router.push(loginRedirect.redirect ?? { name: ROUTE_NAMES.HOME })
  }

  /** 画面刷新时根据令牌恢复用户信息 */

  return { loginInfo, formError, initUserInfo, executeLogin }
}

/** @description 验证登录表单 */
function validateLoginForm(
  info: AuthorizationServices.LoginInfo
): AuthorizationServices.LoginInfoError {
  const error: AuthorizationServices.LoginInfoError = {
    username: '',
    password: ''
  }

  if (info.username.length < 1) {
    error.username = getMessage(MessageKeys.NEED_USERNAME)
  }

  if (info.password.length < 1) {
    error.password = getMessage(MessageKeys.NEED_PASSWORD)
  }

  return error
}
