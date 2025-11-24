<script setup lang="ts"></script>

<template>
  <footer class="app-footer">
    <div class="container footer-content">
      <div class="copyright">&copy; 2025 校园二手交易系统</div>
      <div class="links">
        <a href="#">关于我们</a>
        <a href="#">安全中心</a>
      </div>
      <div class="safety-tip">线下交易请注意人身财产安全</div>
    </div>
  </footer>
</template>

<style scoped lang="scss">
.app-footer {
  background-color: var(--color-bg-white);
  border-top: 1px solid var(--color-border);
  padding: 24px 0;
  margin-top: auto;
}

.footer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--color-text-secondary);
  font-size: 14px;
}

.links {
  display: flex;
  gap: 24px;

  a {
    color: var(--color-text-secondary);
    text-decoration: none;

    &:hover {
      color: var(--color-primary);
    }
  }
}

.safety-tip {
  color: var(--color-text-tertiary);
  font-size: 12px;
}
</style>
