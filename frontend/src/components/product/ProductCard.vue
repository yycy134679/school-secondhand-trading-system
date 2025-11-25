<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { Product } from '@common/types/product'
import ProductStatus from './ProductStatus.vue'
import { formatRelativeTime, formatPrice } from '@/utils/format'

const props = defineProps<{
  product: Product
}>()

const router = useRouter()

const handleClick = () => {
  router.push(`/products/${props.product.id}`)
}
</script>

<template>
  <div class="product-card" @click="handleClick">
    <div class="image-wrapper">
      <img
        :src="product.mainImageUrl || '/placeholder.png'"
        :alt="product.title"
        class="product-image"
      />
      <div class="status-tag" v-if="product.status !== 'ForSale'">
        <ProductStatus :status="product.status" />
      </div>
    </div>
    <div class="content">
      <h3 class="title" :title="product.title">{{ product.title }}</h3>
      <div class="info">
        <span class="price">{{ formatPrice(product.price) }}</span>
        <span class="time">{{ formatRelativeTime(product.createdAt) }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.product-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  border: 1px solid var(--color-border, #eee);

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .image-wrapper {
    position: relative;
    width: 100%;
    padding-top: 100%; // 1:1 Aspect Ratio
    background-color: #f5f5f5;

    .product-image {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .status-tag {
      position: absolute;
      top: 8px;
      right: 8px;
    }
  }

  .content {
    padding: 12px;

    .title {
      margin: 0 0 8px;
      font-size: 14px;
      font-weight: 500;
      color: var(--color-text-primary, #333);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .info {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .price {
        font-size: 16px;
        font-weight: 600;
        color: var(--color-primary, #0066ff);
      }

      .time {
        font-size: 12px;
        color: var(--color-text-secondary, #999);
      }
    }
  }
}
</style>
