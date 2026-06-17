<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    threshold?: number
  }>(),
  {
    threshold: 500
  }
)

const display = ref<boolean>(false)

function handleScroll() {
  display.value = window.scrollY > props.threshold
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <Transition name="fade">
    <button v-show="display" class="scroll-to-top" @click="scrollToTop">
      ↑
    </button>
  </Transition>
</template>

<style scoped>
.scroll-to-top {
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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
