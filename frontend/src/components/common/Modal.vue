<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

const props = withDefaults(defineProps<{
  visible: boolean
  title?: string
  width?: string
}>(), {
  width: '500px'
})

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'close'): void
}>()

const close = () => {
  emit('update:visible', false)
  emit('close')
}

const handleEsc = (e: KeyboardEvent) => {
  if (props.visible && e.key === 'Escape') {
    close()
  }
}

onMounted(() => document.addEventListener('keydown', handleEsc))
onUnmounted(() => document.removeEventListener('keydown', handleEsc))
</script>

<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="close">
      <div class="modal-container" :style="{ width }">
        <div class="modal-header">
          <h3 class="modal-title">{{ title }}</h3>
          <button class="close-btn" @click="close">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <slot></slot>
        </div>

        <div class="modal-footer" v-if="$slots.footer">
          <slot name="footer"></slot>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped lang="scss">
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
  animation: fadeIn 0.2s ease-out;
}

.modal-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.2s ease-out;
}

.modal-header {
  padding: 20px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--color-border, #e5e7eb);
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary, #1a1a1a);
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #999;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;

  &:hover {
    color: #333;
    background-color: #f7f8fa;
  }
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--color-border, #e5e7eb);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>
