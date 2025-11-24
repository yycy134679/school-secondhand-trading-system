<script setup lang="ts"></script>

<template>
  <header class="app-header">
    <div class="container header-content">
      <div class="logo">
        <router-link to="/">校园二手</router-link>
      </div>

      <div class="search-bar">
        <input type="text" placeholder="搜索商品..." />
      </div>

      <div class="actions">
        <button class="btn btn-primary">发布闲置</button>
        <div class="user-info">
          <span>登录</span>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped lang="scss">
.app-header {
  height: 64px;
  background-color: var(--color-bg-white);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo a {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-primary);
  text-decoration: none;
}

.search-bar input {
  width: 400px;
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  outline: none;

  &:focus {
    border-color: var(--color-primary);
  }
}

.actions {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style>
