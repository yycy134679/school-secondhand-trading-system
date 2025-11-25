<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-content">
        <h1 class="welcome-text">欢迎来到校园二手交易平台</h1>
        <p class="sub-text">让闲置物品流动起来，发现身边的宝藏</p>
        <div class="hot-tags">
          <span class="tag-label">热门搜索：</span>
          <router-link
            v-for="tag in hotTags"
            :key="tag"
            :to="{ path: '/search', query: { q: tag } }"
            class="hot-tag"
          >
            {{ tag }}
          </router-link>
        </div>
      </div>
    </section>

    <!-- Category Navigation -->
    <section class="category-nav">
      <div class="container">
        <div class="category-scroll">
          <router-link
            v-for="category in appStore.categories"
            :key="category.id"
            :to="`/category/${category.id}`"
            class="category-item"
          >
            <div class="category-icon">📦</div>
            <span class="category-name">{{ category.name }}</span>
          </router-link>
        </div>
      </div>
    </section>

    <div class="container main-content">
      <!-- Recommendations Section (Logged in only) -->
      <section
        v-if="userStore.isLoggedIn && productStore.homeRecommendations.length > 0"
        class="section recommendations"
      >
        <h2 class="section-title">为你推荐</h2>
        <div class="product-grid">
          <ProductCard
            v-for="product in productStore.homeRecommendations"
            :key="product.id"
            :product="product"
          />
        </div>
      </section>

      <!-- Latest Products Section -->
      <section class="section latest">
        <h2 class="section-title">最新发布</h2>

        <!-- Skeleton Loading State -->
        <div v-if="loading && productStore.homeLatest.items.length === 0" class="product-grid">
          <ProductCardSkeleton v-for="i in pageSize" :key="`skeleton-${i}`" />
        </div>

        <!-- Content -->
        <template v-else>
          <div v-if="productStore.homeLatest.items.length > 0" class="product-grid">
            <ProductCard
              v-for="product in productStore.homeLatest.items"
              :key="product.id"
              :product="product"
            />
          </div>
          <Empty
            v-else
            title="暂无商品"
            description="快来发布第一个商品吧！"
            actionLabel="去发布"
            @action="$router.push('/products/new')"
          />

          <div class="pagination-wrapper" v-if="productStore.homeLatest.total > 0">
            <Pagination
              v-model:page="currentPage"
              :pageSize="pageSize"
              :total="productStore.homeLatest.total"
            />
          </div>
        </template>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAppStore } from '@/stores/app'
import { useProductStore } from '@/stores/product'
import { useUserStore } from '@/stores/user'
import ProductCard from '@/components/product/ProductCard.vue'
import ProductCardSkeleton from '@/components/product/ProductCardSkeleton.vue'
import Pagination from '@/components/common/Pagination.vue'
import Empty from '@/components/common/Empty.vue'

const appStore = useAppStore()
const productStore = useProductStore()
const userStore = useUserStore()

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)

// 热门搜索词，优先从 store 获取，或者使用默认值
const hotTags = computed(() => {
  if (appStore.tags.length > 0) {
    return appStore.tags.slice(0, 5).map((t) => t.name)
  }
  return ['教材', '自行车', '手机', '电脑', '考研资料']
})

const fetchData = async () => {
  loading.value = true
  try {
    await productStore.fetchHomeData({
      page: currentPage.value,
      pageSize: pageSize.value,
    })
  } catch (error) {
    console.error('Failed to fetch home data', error)
  } finally {
    loading.value = false
  }
}

// 监听页码变化
watch(currentPage, () => {
  fetchData()
})

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.home-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* Hero Section */
.hero-section {
  background: linear-gradient(135deg, var(--primary-color) 0%, #4b9eff 100%);
  color: white;
  padding: 60px 0;
  text-align: center;
  margin-bottom: 30px;

  .hero-content {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .welcome-text {
    font-size: 2.5rem;
    font-weight: 700;
    margin-bottom: 16px;
  }

  .sub-text {
    font-size: 1.2rem;
    opacity: 0.9;
    margin-bottom: 30px;
  }

  .hot-tags {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;

    .tag-label {
      font-size: 0.9rem;
      opacity: 0.8;
    }

    .hot-tag {
      background: rgba(255, 255, 255, 0.2);
      padding: 6px 16px;
      border-radius: 20px;
      color: white;
      text-decoration: none;
      font-size: 0.9rem;
      transition: all 0.2s ease;
      backdrop-filter: blur(4px);

      &:hover {
        background: rgba(255, 255, 255, 0.3);
        transform: translateY(-2px);
      }
    }
  }
}

/* Category Navigation */
.category-nav {
  margin-bottom: 40px;

  .category-scroll {
    display: flex;
    gap: 20px;
    overflow-x: auto;
    padding: 10px 4px;
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE and Edge */

    &::-webkit-scrollbar {
      display: none;
    }
  }

  .category-item {
    flex: 0 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    text-decoration: none;
    color: var(--text-color);
    padding: 16px;
    background: white;
    border-radius: 12px;
    min-width: 100px;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      color: var(--primary-color);
    }

    .category-icon {
      font-size: 24px;
    }

    .category-name {
      font-size: 0.9rem;
      font-weight: 500;
    }
  }
}

/* Sections */
.section {
  margin-bottom: 40px;

  .section-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin-bottom: 24px;
    color: var(--text-color);
    display: flex;
    align-items: center;

    &::before {
      content: '';
      display: block;
      width: 4px;
      height: 24px;
      background: var(--primary-color);
      margin-right: 12px;
      border-radius: 2px;
    }
  }
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 24px;
}

.pagination-wrapper {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

/* Responsive */
@media (max-width: 768px) {
  .hero-section {
    padding: 40px 0;

    .welcome-text {
      font-size: 1.8rem;
    }
  }

  .product-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .product-grid {
    grid-template-columns: 1fr;
  }
}
</style>
