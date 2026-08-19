<template>
  <aside class="sidebar" :class="{ open: isOpen }">
    <div class="sidebar-brand"><span>Rub</span>ium</div>
    
    <nav class="sidebar-nav">
      <router-link to="/" class="nav-link">
        <i class="fas fa-house"></i> Главная
      </router-link>
      <router-link to="/trainer" class="nav-link">
        <i class="fas fa-pencil"></i> Тренажёр
      </router-link>
      <router-link to="/daily" class="nav-link">
        <i class="fas fa-calendar-day"></i> Ежедневные
      </router-link>
      <router-link to="/notebooks" class="nav-link">
        <i class="fas fa-book"></i> Мои тетради
      </router-link>
      <router-link to="/community" class="nav-link">
        <i class="fas fa-globe"></i> Каталог тетрадей
      </router-link>
      <router-link to="/profile" class="nav-link">
        <i class="fas fa-user"></i> Профиль
      </router-link>
      <router-link v-if="auth.isAdmin" to="/admin" class="nav-link admin-link">
        <i class="fas fa-shield-halved"></i> Админ
      </router-link>
    </nav>
    
    <div class="sidebar-footer">
      <router-link to="/profile" class="user-link">
        <div class="user-avatar">{{ auth.userName[0]?.toUpperCase() || 'У' }}</div>
        <div class="user-info">
          <div class="user-name">{{ auth.userName }}</div>
          <div class="user-email">{{ auth.user?.email }}</div>
        </div>
      </router-link>
      <button class="logout-btn" @click="auth.logout()">
        <i class="fas fa-sign-out-alt"></i>
      </button>
    </div>
  </aside>
  
  <div class="overlay" :class="{ show: isOpen }" @click="close"></div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const isOpen = ref(false)

const close = () => { isOpen.value = false }
const toggle = () => { isOpen.value = !isOpen.value }

defineExpose({ toggle })
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 240px;
  background: rgba(255,255,255,0.02);
  border-right: 1px solid rgba(255,255,255,0.06);
  backdrop-filter: blur(20px);
  display: flex;
  flex-direction: column;
  padding: 20px 16px;
  z-index: 100;
  transition: transform 0.3s ease;
}

.sidebar-brand {
  font-family: 'Unbounded', sans-serif;
  font-size: 1.3rem;
  font-weight: 800;
  padding: 8px 12px 20px;
  color: #F1F5F9;
}

.sidebar-brand span {
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 10px 14px;
  border-radius: 16px;
  color: #94A3B8;
  text-decoration: none;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.2s;
}

.nav-link:hover {
  background: rgba(255,255,255,0.06);
  color: #F1F5F9;
}

.nav-link.router-link-active {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
}

.nav-link i {
  width: 20px;
  text-align: center;
}

.sidebar-footer {
  border-top: 1px solid rgba(255,255,255,0.06);
  padding-top: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-link {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: inherit;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 0.8rem;
  font-weight: 600;
  color: #F1F5F9;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: 0.65rem;
  color: #64748B;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-btn {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 6px 8px;
  transition: all 0.2s;
}

.logout-btn:hover {
  background: rgba(248,113,113,0.1);
  color: #F87171;
  border-radius: 30%;
}

.overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 99;
}

.overlay.show {
  display: block;
}

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    width: 280px;
  }
  
  .sidebar.open {
    transform: translateX(0);
  }
}
</style>