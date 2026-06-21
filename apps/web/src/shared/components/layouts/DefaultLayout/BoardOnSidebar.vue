<script setup lang="ts">
import { computed } from 'vue'
import type { Board } from '@/features/board/types'

const props = defineProps<{
  item: Board
  actived: boolean
}>()

const background = computed<string>(() => {
  // For Showing
  // return `url('https://i.pixiv.cat/img-master/img/2025/07/20/19/53/22/132915510_p1_master1200.jpg')`
  return `url(${props.item.background})`
})

const title = computed<string>(() => {
  return props.item.name
})
</script>

<template>
  <button class="board-on-sidebar" :class="{ actived }">
    <span class="title">{{ title }}</span>
  </button>
</template>

<style scoped>
.board-on-sidebar {
  align-items: center;
  /* 此为图片加载失败时的降级底色。 */
  background-color: var(--color-normal-box-background);
  background: v-bind(background);
  background-size: cover;
  border: none;
  border-radius: 3px;
  box-shadow: var(--common-block-box-shadow);
  box-sizing: border-box;
  display: flex;
  height: var(--common-board-item-height);
  justify-content: center;
  outline: none;
  padding: 0;
  width: 100%;
}

.board-on-sidebar.actived {
  border: 0.3rem solid var(--color-pink);
}

.title {
  align-items: center;
  background: transparent;
  color: var(--color-white);
  display: flex;
  font-size: 1.1rem;
  height: 100%;
  justify-content: center;
  text-shadow: var(--common-board-item-text-shadow);
  width: 100%;
}
</style>
