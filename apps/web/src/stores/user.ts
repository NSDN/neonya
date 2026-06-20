import { STORE_ID, CONFIG } from '@/shared/constants'
import type { AuthorizationServices } from '@/features/authorization/types'
import useJWTStore from '@/features/authorization/stores/jwt'
import { Option } from '@/shared/utils/rust'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useUserStore = defineStore(STORE_ID.USER, () => {
  const uid = ref<string>('')
  const username = ref<string>('')
  const nickname = ref<string>('游客')
  const userGroup = ref<string>('')
  const userIcon = ref<string>(CONFIG.VISITOR_ICON)

  const loggedIn = computed(() => uid.value !== '')

  const setUsername = (name: string) => (username.value = name)
  const setUserIcon = (icon: string) => (userIcon.value = icon)

  const setUserInfo = (info: AuthorizationServices.UserInfo) => {
    uid.value = info.uid
    username.value = info.username
    nickname.value = info.nickname
    userGroup.value = info.userGroup
    userIcon.value = info.icon
  }

  const logout = () => {
    const jwtStore = useJWTStore()
    jwtStore.setJWT(Option.none())
    uid.value = ''
    username.value = ''
    nickname.value = '游客'
    userGroup.value = ''
    userIcon.value = CONFIG.VISITOR_ICON
  }

  return {
    uid,
    username,
    nickname,
    userGroup,
    userIcon,
    loggedIn,
    setUsername,
    setUserIcon,
    setUserInfo,
    logout
  }
})
