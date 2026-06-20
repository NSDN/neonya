import { useUserStore } from '@/stores/user'
import { computed } from 'vue'
import { ONE_SECOND, ROUTE_NAMES } from '../constants'
import { useRouter } from 'vue-router'

/** 简易个人信息 */
export function useUserSimpleInfo() {
  const router = useRouter()
  const userStore = useUserStore()

  const nickname = computed<string>(() => userStore.nickname)
  const userIcon = computed<string>(() => userStore.userIcon)

  /** 点击简易个人信息面板 */
  const handleClickSimpleInfo = async () => {
    if (userStore.loggedIn) {
      // TODO: 登入后改为跳转到个人信息页（等待个人信息页的实装）
      console.warn('TODO: 登入后改为跳转到个人信息页（等待个人信息页的实装）')
    } else {
      await router.push({ name: ROUTE_NAMES.LOGIN })
    }
  }

  let timer: number
  const clearTimer = () => timer && clearInterval(timer)

  /** 长按简易个人信息面板以登出 */
  const handleLongTimePushSimpleInfo = () => {
    // 未登入时不使长按生效
    if (!userStore.loggedIn) {
      return
    }

    // 清除可能存在的计时器
    clearTimer()

    // 计时超过 n 秒则登出，不到时间放开则会清除计时器
    timer = setTimeout((): void => {
      userStore.logout()
      clearTimer()
    }, 1 * ONE_SECOND)
  }

  /** @description 放开按键 */
  const handleReleaseKey = () => {
    clearTimer()
  }

  return {
    nickname,
    userIcon,
    handleClickSimpleInfo,
    handleLongTimePushSimpleInfo,
    handleReleaseKey
  }
}
