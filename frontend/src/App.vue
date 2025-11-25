<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAppStore } from '@/stores/app'
import AppHeader from './components/common/AppHeader.vue'
import AppFooter from './components/common/AppFooter.vue'

const appStore = useAppStore()

onMounted(() => {
  appStore.initDictionaries()
})
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
