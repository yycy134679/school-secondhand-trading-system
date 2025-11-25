<script setup lang="ts">
import { ProductStatus } from '@common/constants/product_status'
import type { User } from '@common/types/user'

console.log('Admin Product Status:', ProductStatus.SOLD)
const admin: User = {
  id: 1,
  account: 'admin',
  nickname: 'Admin User',
  isAdmin: true,
  createdAt: '',
  updatedAt: '',
}
console.log('Admin User:', admin)
</script>

<template>
  <h1>Admin Frontend App</h1>
  <p>Check console for common module imports.</p>
</template>

<style scoped></style>
