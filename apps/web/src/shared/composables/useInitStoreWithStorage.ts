import { useCurrentBoard } from '@/features/board/composables/useCurrentBoard'

export function useInitStoreWithStorage() {
  const { initCurrentBoardWithStorage } = useCurrentBoard()

  const initStoreWithStorage = () => {
    initCurrentBoardWithStorage()
  }

  return { initStoreWithStorage }
}
