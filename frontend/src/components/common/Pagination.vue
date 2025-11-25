<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  (e: 'update:page', value: number): void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const handlePageChange = (newPage: number) => {
  if (newPage >= 1 && newPage <= totalPages.value) {
    emit('update:page', newPage)
  }
}
</script>

<template>
  <div class="pagination" v-if="total > 0">
    <button
      class="btn-nav"
      :disabled="page === 1"
      @click="handlePageChange(page - 1)"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="15 18 9 12 15 6"></polyline>
      </svg>
    </button>

    <div class="page-info">
      <span class="current">{{ page }}</span>
      <span class="separator">/</span>
      <span class="total">{{ totalPages }}</span>
    </div>

    <button
      class="btn-nav"
      :disabled="page === totalPages"
      @click="handlePageChange(page + 1)"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="9 18 15 12 9 6"></polyline>
      </svg>
    </button>

    <span class="total-count">共 {{ total }} 条</span>
  </div>
</template>

<style scoped lang="scss">
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin: 24px 0;
  font-size: 14px;
  color: var(--color-text-secondary, #666);
}

.btn-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-border, #e5e7eb);
  background-color: white;
  border-radius: 4px;
  cursor: pointer;
  color: var(--color-text-primary, #1a1a1a);
  transition: all 0.2s;

  &:hover:not(:disabled) {
    border-color: var(--color-primary, #0066ff);
    color: var(--color-primary, #0066ff);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    background-color: #f7f8fa;
  }
}

.page-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;

  .current {
    color: var(--color-primary, #0066ff);
  }
}

.total-count {
  margin-left: 8px;
  color: var(--color-text-tertiary, #999);
  font-size: 12px;
}
</style>
