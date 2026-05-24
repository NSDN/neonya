<script setup lang="ts">
import { CONFIG } from '@/shared/constants'
import Header from './components/Header.vue'
</script>

<template>
  <div
    id="main-layout"
    :style="`background-image: url(${CONFIG.BASE_BACKGROUND})`"
  >
    <Header />

    <div class="middle">
      <div class="content"><RouterView /></div>
    </div>
  </div>
</template>

<style scoped>
#main-layout {
  background-size: cover;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.middle {
  display: flex;
  flex: 1 1 0;
  overflow: auto;
}

.middle .content {
  background: var(--color-layout-content-background);
  box-sizing: border-box;
  display: flex;
  flex: 1;
  padding: var(--common-content-padding);
  overflow: auto;
}
</style>
