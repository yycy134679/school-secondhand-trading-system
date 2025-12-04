<template>
  <div class="not-found">
    <h2>页面不存在</h2>
    <p>访问的后台页面还没有实现，返回仪表盘继续操作吧。</p>
    <RouterLink class="link" to="/admin">返回仪表盘</RouterLink>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

defineOptions({ name: 'AdminNotFoundView' })
</script>

<style scoped>
.not-found {
  text-align: center;
  padding: 80px 16px;
}

.not-found h2 {
  font-size: 28px;
  margin-bottom: 12px;
}

.not-found p {
  color: #6b7280;
  margin-bottom: 24px;
}

.link {
  color: #2563eb;
  text-decoration: none;
}
</style>
