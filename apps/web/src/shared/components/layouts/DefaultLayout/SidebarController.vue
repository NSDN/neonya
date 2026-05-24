<template>
  <button class="sidebar-controller">|</button>
</template>

<style scoped>
.sidebar-controller {
  background: var(--color-sidebar-background);
  border: unset;
  box-shadow: 2px 0 6px 2px #bbb;
  color: #666;
  cursor: pointer;
  font-size: 12px;
  font-weight: bold;
  height: 100%;
  left: 0;
  opacity: 0.4;
  padding: unset;
  position: absolute;
  top: 0;
}

.sidebar-controller:hover {
  opacity: 1;
}

@media (--mobile) {
  .sidebar-controller {
    display: none;
  }
}
</style>
