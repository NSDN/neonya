import { STORE_ID } from '@/shared/constants'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Board } from '../types'
import { Option, type Optionable } from '@/shared/utils/rust'
import { toOption } from '@/shared/utils/useful'

export const useBoardsStore = defineStore(STORE_ID.BOARDS, () => {
  const boards = ref<Board[]>([])

  const setBoards = (payload: Board[]) => {
    boards.value = payload
  }

  const currentBoard = ref<Option<Board>>(Option.none())

  const setCurrentBoard = (payload: Optionable<Board>) => {
    currentBoard.value = toOption(payload)
  }

  const loading = ref(false)

  return { boards, loading, setBoards, currentBoard, setCurrentBoard }
})
