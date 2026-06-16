import { useBoardsStore } from '../stores/boards'
import { storage } from '@/shared/services/storage'
import { PageType, STORAGE_KEYS } from '@/shared/constants'
import type { Board, BoardId } from '../types'
import { toOption } from '@/shared/utils/useful'
import { computed } from 'vue'

export function useCurrentBoard() {
  const boardsStore = useBoardsStore()

  const currentBoardId = computed<BoardId>(() => {
    return boardsStore.currentBoard.match({
      some: board => board.id,
      none: () => boardsStore.boards[0]?.id ?? 'localization'
    })
  })

  const isCurrentBoard = (id: BoardId) => {
    return boardsStore.currentBoard.match({
      some: board => board.id === id,
      none: () => false
    })
  }

  const isComicBoard = computed<boolean>(() =>
    boardsStore.currentBoard.match({
      some: board => board.pageType === PageType.COMIC,
      none: () => false
    })
  )

  const setCurrentBoardById = (id: BoardId) => {
    const target = toOption(boardsStore.boards.find(board => board.id === id))

    target.match({
      some: board => {
        boardsStore.setCurrentBoard(target)
        storage.set<Board>(STORAGE_KEYS.CURRENT_BOARD, board)
      },

      none: () => {}
    })
  }

  const initCurrentBoardWithStorage = () => {
    storage.get<Board>(STORAGE_KEYS.CURRENT_BOARD).match({
      ok: board => boardsStore.setCurrentBoard(board),
      err: error => error.notify()
    })
  }

  return {
    currentBoardId,
    isCurrentBoard,
    isComicBoard,
    setCurrentBoardById,
    initCurrentBoardWithStorage
  }
}
