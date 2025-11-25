<script setup lang="ts"></script>

<template>
  <footer class="app-footer">
    <div class="container footer-content">
      <div class="copyright">&copy; 2025 校园二手交易系统</div>
      <div class="links">
        <a href="https://github.com/yycy134679/school-secondhand-trading-system" target="_blank" rel="noopener noreferrer">项目源码</a>
        <a href="#">关于我们</a>
        <a href="#">安全中心</a>
      </div>
      <div class="safety-tip">
        <span class="icon">🛡️</span>
        线下交易请注意人身财产安全，建议在校内公共场所进行交易
      </div>
    </div>
  </footer>
</template>

<style scoped lang="scss">
.app-footer {
  background-color: var(--color-bg-white, #fff);
  border-top: 1px solid var(--color-border, #e5e7eb);
  padding: 32px 0;
  margin-top: auto;
}

.footer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: var(--color-text-secondary, #666);
  font-size: 14px;
  max-width: 1200px;
  margin: 0 auto;
}

.links {
  display: flex;
  gap: 32px;

  a {
    color: var(--color-text-secondary, #666);
    text-decoration: none;
    transition: color 0.2s;

    &:hover {
      color: var(--color-primary, #0066ff);
    }
  }
}

.safety-tip {
  color: var(--color-text-tertiary, #999);
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  background-color: #f7f8fa;
  padding: 6px 16px;
  border-radius: 16px;

  .icon {
    font-size: 14px;
  }
}
</style>
