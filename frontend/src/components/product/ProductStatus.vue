<script setup lang="ts">
import { computed } from 'vue'
import { ProductStatus } from '@common/constants/product_status'
import type { ProductStatus as ProductStatusType } from '@common/types/product'

const props = defineProps<{
  status: ProductStatusType
}>()

const statusConfig = computed(() => {
  switch (props.status) {
    case ProductStatus.FOR_SALE:
      return { label: '在售', className: 'status-for-sale' }
    case ProductStatus.SOLD:
      return { label: '已售', className: 'status-sold' }
    case ProductStatus.DELISTED:
      return { label: '下架', className: 'status-delisted' }
    default:
      return { label: props.status, className: '' }
  }
})
</script>

<template>
  <span class="product-status" :class="statusConfig.className">
    {{ statusConfig.label }}
  </span>
</template>

<style scoped lang="scss">
.product-status {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  line-height: 1.5;
  color: #fff;

  &.status-for-sale {
    background-color: var(--color-success, #52c41a);
    color: #fff;
  }

  &.status-sold {
    background-color: var(--color-text-secondary, #999);
    color: #fff;
  }

  &.status-delisted {
    background-color: var(--color-error, #ff4d4f);
    color: #fff;
  }
}
</style>
