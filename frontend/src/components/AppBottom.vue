<template>
  <footer class="app-bottom">
    <div class="feature-highlights">
      <div class="feature-item">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="7" height="7"/>
          <rect x="14" y="3" width="7" height="7"/>
          <rect x="14" y="14" width="7" height="7"/>
          <rect x="3" y="14" width="7" height="7"/>
        </svg>
        <span>可视化看板</span>
      </div>
      <div class="feature-item">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
          <circle cx="9" cy="7" r="4"/>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
        </svg>
        <span>团队协作</span>
      </div>
      <div class="feature-item">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
        </svg>
        <span>实时追踪</span>
      </div>
    </div>
    <div class="bottom-content">
      <span class="copyright">© 2026 Taskify</span>
      <span class="divider">·</span>
      <span class="tagline">让协作更简单</span>
    </div>
  </footer>
</template>

<script setup lang="ts">
// Simple footer component for all pages
</script>

<style scoped>
.app-bottom {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: var(--space-3) var(--space-4);
  padding-bottom: calc(var(--space-3) + env(safe-area-inset-bottom, 0px));
  background: var(--color-bg-surface);
  border-top: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  z-index: 50;
}

.feature-highlights {
  display: flex;
  justify-content: center;
  gap: var(--space-5);
}

.feature-item {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.feature-item svg {
  color: var(--color-primary);
  opacity: 0.7;
}

.bottom-content {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.copyright {
  font-weight: var(--font-weight-medium);
}

.divider {
  opacity: 0.5;
}

.tagline {
  font-style: italic;
  opacity: 0.8;
}

/* Safe area for notched phones */
@supports (padding-bottom: env(safe-area-inset-bottom)) {
  .app-bottom {
    padding-bottom: calc(var(--space-3) + env(safe-area-inset-bottom));
  }
}
</style>
