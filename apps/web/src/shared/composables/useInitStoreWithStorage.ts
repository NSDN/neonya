import { useLogin } from '@/features/authorization/composables/useLogin'
import { useCurrentBoard } from '@/features/board/composables/useCurrentBoard'

export function useInitStoreWithStorage() {
  const { initCurrentBoardWithStorage } = useCurrentBoard()
  const { initUserInfo } = useLogin()

  const initStoreWithStorage = async () => {
    initCurrentBoardWithStorage()
    await initUserInfo()
  }

  return { initStoreWithStorage }
}
