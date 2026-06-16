import { type ThreadListItem } from '../types'
// import { useRouter } from 'vue-router'
import { getThreads } from '../apis/threads'
import { useCurrentBoard } from '@/features/board/composables/useCurrentBoard'
import { Option } from '@/shared/utils/rust'
import { useThreadsStore } from '../stores/threads'
import { computed } from 'vue'

export function useThreads() {
  // const router = useRouter()
  const { currentBoardId } = useCurrentBoard()
  const threadsStore = useThreadsStore()

  const handleGetThreads = async () => {
    const result = await getThreads(currentBoardId.value)

    const optionalThreads = result.match({
      ok: inner =>
        inner.andThen(threads => {
          const computedThreads = threads.flatMap(thread =>
            thread.match({
              some: value => [value],
              none: () => []
            })
          )

          return Option.some(computedThreads)
        }),

      err: error => {
        error.notify()
        return Option.none()
      }
    })

    threadsStore.setThreads(optionalThreads)
  }

  const threads = computed<ThreadListItem[]>(() => threadsStore.threads)

  const goToThread = async (id: string) => {
    // TODO: 跳转至帖子画面。
    console.log({ id })
    // await router.push({
    // name: ROUTE_NAMES.BOARD,
    //   params: { id }
    // })
  }

  return {
    handleGetThreads,
    threads,
    goToThread
  }
}
