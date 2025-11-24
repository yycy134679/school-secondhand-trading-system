<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeader from './components/common/AppHeader.vue'
import AppFooter from './components/common/AppFooter.vue'
</script>

<template>
  <div class="app-layout">
    <AppHeader />

    <main class="main-content container">
      <RouterView />
    </main>

    <AppFooter />
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  padding: 24px 0;
}
</style>
