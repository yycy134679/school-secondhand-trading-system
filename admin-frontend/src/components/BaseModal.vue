<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div v-if="open" class="modal-backdrop" @click.self="handleBackdrop">
        <div class="modal" :style="{ maxWidth: width }" role="dialog" aria-modal="true">
          <header v-if="title" class="modal-header">
            <h2 class="modal-title">{{ title }}</h2>
            <button type="button" class="modal-close" aria-label="关闭" @click="emitClose">×</button>
          </header>
          <section class="modal-body">
            <slot />
          </section>
          <footer v-if="$slots.footer" class="modal-footer">
            <slot name="footer" />
          </footer>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title?: string
    width?: string
    closeOnEsc?: boolean
    closeOnBackdrop?: boolean
  }>(),
  {
    closeOnEsc: true,
    closeOnBackdrop: true,
  },
)

const emit = defineEmits<{
  (event: 'close'): void
}>()

const handleBackdrop = () => {
  if (props.closeOnBackdrop) {
    emit('close')
  }
}

const handleKeydown = (event: KeyboardEvent) => {
  if (props.closeOnEsc && event.key === 'Escape' && props.open) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.classList.remove('modal-open')
})

watch(
  () => props.open,
  (value) => {
    document.body.classList.toggle('modal-open', value)
  },
  { immediate: true },
)

const emitClose = () => emit('close')
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(2px);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  z-index: 2000;
}

.modal {
  width: 100%;
  max-width: 520px;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.26);
  border: 1px solid #e5e7eb;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 20px 0;
}

.modal-title {
  margin: 0;
  font-size: 18px;
}

.modal-close {
  border: none;
  background: transparent;
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
  color: #6b7280;
}

.modal-body {
  padding: 20px;
  color: #1f2937;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 0 20px 20px;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

:global(body.modal-open) {
  overflow: hidden;
}
</style>
