<script setup lang="ts">
import { useBoards } from '@/features/board/composables/useBoards'
import { useCurrentBoard } from '@/features/board/composables/useCurrentBoard'
import type { BoardId } from '@/features/board/types'
import { onMounted, ref } from 'vue'

const { boards, handleGetBoards, goToBoard } = useBoards()
const { isCurrentBoard } = useCurrentBoard()

const show = ref(false)

function openSheet() {
  show.value = true
}

function closeSheet() {
  show.value = false
}

onMounted(async () => {
  await handleGetBoards()
})

function handleSelectBoard(id: BoardId) {
  goToBoard(id)
  closeSheet()
}
</script>

<template>
  <button class="mobile-only board-menu-button" @click="openSheet">☰</button>

  <Teleport to="body">
    <Transition name="sheet">
      <div v-if="show" class="bottom-sheet-overlay" @click.self="closeSheet">
        <Transition name="sheet-panel">
          <div class="bottom-sheet-panel">
            <button
              v-for="(item, index) of boards"
              :key="index.toString()"
              class="sheet-board-item"
              :class="{ actived: isCurrentBoard(item.id) }"
              @click="() => handleSelectBoard(item.id)"
            >
              {{ item.name }}
            </button>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.board-menu-button {
  background: var(--color-normal-box-background);
  border: 1px solid var(--color-pink);
  border-radius: 0.5rem;
  box-shadow: var(--shadow-normal-box-shadow);
  color: var(--color-pink);
  cursor: pointer;
  font-size: 1.5rem;
  height: 3rem;
  line-height: 1;
  padding: 0;
  width: 3rem;
}

.mobile-only {
  display: none;
}

@media (--mobile) {
  .mobile-only {
    display: block;
  }
}

.bottom-sheet-overlay {
  align-items: flex-end;
  background: rgba(0, 0, 0, 0.5);
  bottom: 0;
  display: flex;
  justify-content: center;
  left: 0;
  position: fixed;
  right: 0;
  top: 0;
  z-index: 1000;
}

.bottom-sheet-panel {
  background: var(--color-sidebar-background);
  border-radius: 1rem 1rem 0 0;
  box-sizing: border-box;
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem;
  max-height: 70vh;
  overflow: auto;
  padding: 1rem;
  width: 100%;
}

.sheet-board-item {
  background: var(--color-normal-box-background);
  border: none;
  border-radius: 0.5rem;
  box-sizing: border-box;
  color: inherit;
  flex: 1 1 calc(50% - 0.4rem);
  font-size: 1rem;
  min-width: calc(50% - 0.4rem);
  padding: 1rem 0.5rem;
}

.sheet-board-item.actived {
  border: 0.2rem solid var(--color-pink);
}

.sheet-enter-active {
  transition: opacity 0.3s ease;
}

.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
}

.sheet-panel-enter-active {
  transition: transform 0.3s ease;
}

.sheet-panel-enter-from,
.sheet-panel-leave-to {
  transform: translateY(100%);
}
</style>
