<template>
  <div class="product-card-skeleton">
    <div class="image-skeleton skeleton-animate"></div>
    <div class="content">
      <div class="title-skeleton skeleton-animate"></div>
      <div class="info">
        <div class="price-skeleton skeleton-animate"></div>
        <div class="time-skeleton skeleton-animate"></div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.product-card-skeleton {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--color-border, #eee);

  .image-skeleton {
    width: 100%;
    padding-top: 100%; // 1:1 Aspect Ratio
    background-color: #f0f0f0;
  }

  .content {
    padding: 12px;

    .title-skeleton {
      height: 16px;
      background-color: #f0f0f0;
      border-radius: 4px;
      margin-bottom: 12px;
    }

    .info {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .price-skeleton {
        width: 60px;
        height: 18px;
        background-color: #f0f0f0;
        border-radius: 4px;
      }

      .time-skeleton {
        width: 50px;
        height: 14px;
        background-color: #f0f0f0;
        border-radius: 4px;
      }
    }
  }
}

.skeleton-animate {
  animation: skeleton-loading 1.5s ease-in-out infinite;
  background: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #e0e0e0 50%,
    #f0f0f0 75%
  );
  background-size: 200% 100%;
}

@keyframes skeleton-loading {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
