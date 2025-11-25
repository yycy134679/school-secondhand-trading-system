<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
  actionLabel?: string
}>()

const emit = defineEmits<{
  (e: 'action'): void
}>()
</script>

<template>
  <div class="empty-state">
    <div class="icon">
      <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
        <polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline>
        <line x1="12" y1="22.08" x2="12" y2="12"></line>
      </svg>
    </div>
    <h3 class="title">{{ title || '暂无数据' }}</h3>
    <p class="description" v-if="description">{{ description }}</p>
    <button
      v-if="actionLabel"
      class="btn btn-primary action-btn"
      @click="emit('action')"
    >
      {{ actionLabel }}
    </button>
  </div>
</template>

<style scoped lang="scss">
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
  color: var(--color-text-secondary, #666);
}

.icon {
  margin-bottom: 16px;
  color: #e5e7eb;
}

.title {
  font-size: 18px;
  font-weight: 500;
  color: var(--color-text-primary, #1a1a1a);
  margin-bottom: 8px;
}

.description {
  font-size: 14px;
  color: var(--color-text-tertiary, #999);
  margin-bottom: 24px;
  max-width: 300px;
  line-height: 1.5;
}

.action-btn {
  padding: 8px 24px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  background-color: var(--color-primary, #0066ff);
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #0052cc;
  }
}
</style>
