import { useCurrentPlate } from '@/features/plate/composables/useCurrentPlate'

export function useInitStoreWithStorage() {
  const { initCurrentPlateWithStorage } = useCurrentPlate()

  const initStoreWithStorage = () => {
    initCurrentPlateWithStorage()
  }

  return { initStoreWithStorage }
}
