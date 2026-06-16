import { computed } from 'vue'
import { useBoardsStore } from '../stores/boards'
import { getBoards } from '../apis/boards'
import { ResponsError } from '@/shared/errors'
import { ErrorLevel, ROUTE_NAMES } from '@/shared/constants'
import { type Board, type BoardId } from '../types'
import { Option, Result } from '@/shared/utils/rust'
import { useCurrentBoard } from './useCurrentBoard'
import { useRouter } from 'vue-router'
import { useThreads } from '@/features/thread'

export function useBoards() {
  const router = useRouter()
  const boardsStore = useBoardsStore()
  const { setCurrentBoardById } = useCurrentBoard()
  const { handleGetThreads } = useThreads()

  const boards = computed<Board[]>(() => boardsStore.boards)

  function shouldSkipFetch() {
    const conditions = [boardsStore.loading, boards.value.length > 0]
    return conditions.some(Boolean)
  }

  const handleGetBoards = async () => {
    if (shouldSkipFetch()) return

    boardsStore.loading = true
    const result = await getBoards()

    result
      .andThen(data =>
        data
          .andThen(boards => {
            const computedBoards = boards.flatMap(board =>
              board.match({
                some: value => [value],
                none: () => []
              })
            )

            if (computedBoards.length === 0) {
              return Option.none()
            }

            return Option.some(computedBoards)
          })
          .match({
            some: boards => Result.ok(boards),

            none: () => {
              const error = new ResponsError({
                level: ErrorLevel.Error,
                message: '[Boards]: 没有取到版块列表。'
              })

              return Result.err(error)
            }
          })
      )
      .match({
        ok: boards => {
          boardsStore.setBoards(
            boards.sort((previous, next) => previous.sortOrder - next.sortOrder)
          )

          if (boardsStore.currentBoard.isSome()) {
            return
          }

          setCurrentBoardById(boards[0].id)
        },

        err: error => error.notify()
      })

    boardsStore.loading = false
  }

  const goToBoard = async (id: BoardId) => {
    await router.push({
      name: ROUTE_NAMES.BOARD,
      params: { id }
    })

    setCurrentBoardById(id)

    await handleGetThreads()
  }

  return {
    boards,
    handleGetBoards,
    goToBoard
  }
}
