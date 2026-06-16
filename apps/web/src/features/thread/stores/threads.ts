import { STORE_ID } from '@/shared/constants'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { ThreadListItem } from '../types'
import type { Optionable } from '@/shared/utils/rust'
import { toOption } from '@/shared/utils/useful'

export const useThreadsStore = defineStore(STORE_ID.THREADS, () => {
  const threads = ref<ThreadListItem[]>([])

  const setThreads = (payload: Optionable<ThreadListItem[]>) => {
    toOption(payload).match({
      some: value => (threads.value = value),
      none: () => (threads.value = [])
    })
  }

  return { threads, setThreads }
})
