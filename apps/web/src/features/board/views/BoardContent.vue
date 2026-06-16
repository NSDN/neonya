<script setup lang="ts">
import { useThreads } from '@/features/thread'

import Announcement from '../components/Announcement.vue'
import ArticleList from '../components/ArticleList.vue'
import { useCurrentBoard } from '../composables/useCurrentBoard'
import { onMounted } from 'vue'
// import PictureList from '@/features/comic/components/PictureList.vue'

const { handleGetThreads, threads, goToThread } = useThreads()
const { isComicBoard } = useCurrentBoard()

onMounted(async () => {
  await handleGetThreads()
})
</script>

<template>
  <div id="board-content">
    <Announcement />

    <div class="search">筛选 / 检索</div>

    <div v-if="isComicBoard">漫画列表</div>
    <!-- 
    <PictureList
      v-if="board.isComicBoard"
      :list="commic.list"
      @click-item="commic.clickListItem"
    />
    -->
    <ArticleList v-else :list="threads" @click="goToThread" />
  </div>
</template>

<style scoped>
#board-content {
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 1rem;
}

.search {
  align-items: center;
  background: var(--color-normal-box-background);
  display: flex;
  height: 5rem;
  justify-content: center;
  margin: 1rem 0 0;
}
</style>
