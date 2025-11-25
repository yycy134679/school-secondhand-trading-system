<script setup lang="ts">
withDefaults(defineProps<{
  loading?: boolean
  type?: 'fullscreen' | 'skeleton' | 'spinner'
}>(), {
  loading: false,
  type: 'spinner'
})
</script>

<template>
  <div v-if="loading" class="loading-container" :class="type">
    <!-- Fullscreen Loading -->
    <div v-if="type === 'fullscreen'" class="loading-mask">
      <div class="spinner"></div>
    </div>

    <!-- Skeleton Loading -->
    <div v-else-if="type === 'skeleton'" class="skeleton-wrapper">
      <div class="skeleton-block image"></div>
      <div class="skeleton-block title"></div>
      <div class="skeleton-block text"></div>
      <div class="skeleton-block text short"></div>
    </div>

    <!-- Inline Spinner -->
    <div v-else class="spinner-wrapper">
      <div class="spinner"></div>
    </div>
  </div>
  <slot v-else></slot>
</template>

<style scoped lang="scss">
.loading-container {
  &.fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 9999;
  }
}

.loading-mask {
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(2px);
  display: flex;
  justify-content: center;
  align-items: center;
}

.spinner-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid var(--color-primary, #0066ff);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

// Skeleton Styles
.skeleton-wrapper {
  width: 100%;
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.skeleton-block {
  background: linear-gradient(90deg, #f0f0f0 25%, #f8f8f8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 12px;

  &.image {
    width: 100%;
    padding-bottom: 100%; // 1:1 aspect ratio
    margin-bottom: 16px;
    border-radius: 8px;
  }

  &.title {
    height: 24px;
    width: 80%;
  }

  &.text {
    height: 16px;
    width: 100%;

    &.short {
      width: 60%;
    }
  }
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
