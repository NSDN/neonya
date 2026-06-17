<script setup lang="ts">
import { CONFIG } from '@/shared/constants'
import { onMounted, ref } from 'vue'

import Sidebar from './Sidebar.vue'
import SidebarController from './SidebarController.vue'
import FloatingActionGroup from '../components/FloatingActionGroup/FloatingActionGroup.vue'
import Header from '../components/Header.vue'
import { useNaiveUIGlobalConfig } from '@/shared/composables'

const displaySidebar = ref<boolean>(true)
const controlSidebar = () => (displaySidebar.value = !displaySidebar.value)

const { initMessager } = useNaiveUIGlobalConfig()
onMounted(() => initMessager())
</script>

<template>
  <div
    id="main-layout"
    :style="`background-image: url(${CONFIG.BASE_BACKGROUND})`"
  >
    <Header />

    <div class="middle">
      <Transition name="slide">
        <Sidebar v-show="displaySidebar" class="desktop-only" />
      </Transition>

      <div class="content">
        <SidebarController class="desktop-only" @click="controlSidebar" />
        <slot />
      </div>
    </div>

    <FloatingActionGroup />
  </div>
</template>

<style scoped>
#main-layout {
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.middle {
  display: flex;
  flex: 1;
  overflow: auto;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease-in-out;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}

.middle .content {
  background: var(--color-layout-content-background);
  box-sizing: border-box;
  display: flex;
  flex: 1;
  overflow: hidden;
  padding: var(--common-content-padding);
  position: relative;
}

.desktop-only {
  display: block;
}

.mobile-only {
  display: none;
}

@media (--mobile) {
  .desktop-only {
    display: none;
  }

  .mobile-only {
    display: block;
  }
}
</style>
