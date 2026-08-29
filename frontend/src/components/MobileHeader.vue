<template>
  <header class="mobile-header">
    <button class="mobile-menu-btn" @click="$emit('toggle')">
      <i class="fas fa-bars"></i>
    </button>
  </header>
</template>

<style scoped>
.mobile-header {
  display: none;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  position: sticky;
  top: 0;
  background: #0a0a0a;
  z-index: 500;
  width: 100%;
  height: 48px;
}

.mobile-menu-btn {
  background: none;
  border: none;
  color: #737373;
  font-size: 1.1rem;
  cursor: pointer;
  padding: 4px 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s ease;
  height: 100%;
}

.mobile-menu-btn:hover {
  color: #e5e5e5;
}

@media (max-width: 768px) {
  .mobile-header {
    display: flex;
  }
}
</style>