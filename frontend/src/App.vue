<script setup lang="ts">
import { ProductStatus } from '@common/constants/product_status'
import type { User } from '@common/types/user'

console.log('Product Status:', ProductStatus.FOR_SALE)
const user: User = {
  id: 1,
  account: 'test',
  nickname: 'Test User',
  isAdmin: false,
  createdAt: '',
  updatedAt: '',
}
console.log('User:', user)
</script>

<template>
  <h1>Frontend App</h1>
  <p>Check console for common module imports.</p>
</template>

<style scoped></style>
