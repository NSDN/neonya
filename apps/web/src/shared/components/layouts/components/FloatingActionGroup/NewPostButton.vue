<script setup lang="ts">
function handleClick() {
  window.$message.warning('TODO')
}
</script>

<template>
  <button class="new-post-button" @click="handleClick">＋</button>
</template>

<style scoped>
.new-post-button {
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
</style>
