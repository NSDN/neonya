<script setup lang="ts">
import { usePlates } from '@/features/plate/composables/usePlates'
import { useCurrentPlate } from '@/features/plate/composables/useCurrentPlate'
import type { PlateId } from '@/features/plate/types'
import { onMounted, ref, watch } from 'vue'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
}>()

const { plates, handleGetPlates, goToPlate } = usePlates()
const { isCurrentPlate } = useCurrentPlate()

const visible = ref(false)

onMounted(async () => {
  await handleGetPlates()
})

watch(() => props.show, (value) => {
  if (value) {
    visible.value = true
  }
})

function handleClose() {
  emit('update:show', false)
}

function handleAfterLeave() {
  visible.value = false
}

function handleSelectPlate(id: PlateId) {
  goToPlate(id)
  emit('update:show', false)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="sheet" @after-leave="handleAfterLeave">
      <div v-if="visible && show" class="bottom-sheet-overlay" @click.self="handleClose">
        <Transition name="sheet-panel">
          <div v-if="show" class="bottom-sheet-panel">
            <div class="sheet-header">
              <span class="sheet-title">版块</span>
              <button class="sheet-close" @click="handleClose">✕</button>
            </div>
            <div class="sheet-plates">
              <button
                v-for="(item, index) of plates"
                :key="index.toString()"
                class="sheet-plate-item"
                :class="{ actived: isCurrentPlate(item.id) }"
                @click="() => handleSelectPlate(item.id)"
              >
                {{ item.name }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
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
  max-height: 70vh;
  overflow: auto;
  padding: 1rem;
  width: 100%;
}

.sheet-header {
  align-items: center;
  display: flex;
  justify-content: space-between;
  margin: 0 0 1rem;
}

.sheet-title {
  font-size: 1.2rem;
  font-weight: 600;
}

.sheet-close {
  all: unset;
  cursor: pointer;
  font-size: 1.2rem;
}

.sheet-plates {
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem;
}

.sheet-plate-item {
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

.sheet-plate-item.actived {
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
