<script setup lang="ts"></script>

<template>
  <div class="product-detail-skeleton">
    <div class="detail-container">
      <!-- 左侧图片骨架 -->
      <div class="image-section">
        <div class="skeleton-box main-image"></div>
        <div class="thumbnail-list">
          <div class="skeleton-box thumbnail" v-for="i in 4" :key="i"></div>
        </div>
      </div>

      <!-- 右侧信息骨架 -->
      <div class="info-section">
        <div class="skeleton-line title"></div>
        <div class="skeleton-line price"></div>
        <div class="skeleton-tags">
          <div class="skeleton-box tag" v-for="i in 2" :key="i"></div>
        </div>
        <div class="skeleton-line meta"></div>
        <div class="skeleton-box seller-card"></div>
        <div class="skeleton-description">
          <div class="skeleton-line" v-for="i in 3" :key="i"></div>
        </div>
        <div class="skeleton-box action-button"></div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.product-detail-skeleton {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.detail-container {
  display: flex;
  gap: 40px;

  @media (max-width: 768px) {
    flex-direction: column;
  }
}

.image-section {
  flex: 0 0 60%;
  max-width: 60%;

  @media (max-width: 768px) {
    flex: 0 0 100%;
    max-width: 100%;
  }

  .main-image {
    width: 100%;
    aspect-ratio: 4/3;
    margin-bottom: 16px;
  }

  .thumbnail-list {
    display: flex;
    gap: 10px;

    .thumbnail {
      width: 80px;
      height: 80px;
      flex-shrink: 0;
    }
  }
}

.info-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.skeleton-box,
.skeleton-line {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading 1.5s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-line {
  height: 20px;

  &.title {
    height: 32px;
    width: 80%;
  }

  &.price {
    height: 36px;
    width: 40%;
  }

  &.meta {
    width: 60%;
  }
}

.skeleton-tags {
  display: flex;
  gap: 8px;

  .tag {
    width: 80px;
    height: 24px;
  }
}

.skeleton-box {
  &.seller-card {
    height: 80px;
  }

  &.action-button {
    height: 48px;
    margin-top: auto;
  }
}

.skeleton-description {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .skeleton-line {
    &:last-child {
      width: 70%;
    }
  }
}

@keyframes loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
