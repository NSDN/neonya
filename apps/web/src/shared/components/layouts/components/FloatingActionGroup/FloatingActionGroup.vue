<script setup lang="ts">
import ScrollToTopButton from './ScrollToTopButton.vue'
import NewPostButton from './NewPostButton.vue'
import BoardMenuButton from './BoardMenuButton.vue'
</script>

<template>
  <div class="floating-action-group">
    <ScrollToTopButton />
    <NewPostButton />
    <BoardMenuButton />
  </div>
</template>

<style scoped>
.floating-action-group {
  bottom: 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  position: fixed;
  right: 2rem;
  z-index: 100;
}
</style>
