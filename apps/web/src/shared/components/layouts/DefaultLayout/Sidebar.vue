<script setup lang="ts">
import { useBoards } from '@/features/board/composables/useBoards'
import { useCurrentBoard } from '@/features/board/composables/useCurrentBoard'
import { onMounted } from 'vue'

import BoardOnSidebar from './BoardOnSidebar.vue'

const { boards, handleGetBoards, goToBoard } = useBoards()
const { isCurrentBoard } = useCurrentBoard()

onMounted(async () => {
  await handleGetBoards()
})
</script>

<template>
  <div class="sidebar">
    <div class="board-group">
      <BoardOnSidebar
        class="board-item"
        v-for="(item, index) of boards"
        :key="index.toString()"
        :item="item"
        :actived="isCurrentBoard(item.id)"
        @click="() => goToBoard(item.id)"
      />
    </div>
  </div>
</template>

<style scoped>
.sidebar {
  background: var(--color-sidebar-background);
  box-sizing: border-box;
  box-shadow: 4px 0 5px 0 #888;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  width: 16rem;
}

.board-group {
  flex: 1;
  box-sizing: border-box;
  margin: 1rem -1rem 0;
  overflow: auto;
  padding: 0.3rem 1rem;
}

.board-item {
  margin: 1rem 0 0;
}

.board-item:first-child {
  margin: 0;
}

@media (--mobile) {
  .sidebar {
    display: none;
  }
}
</style>
