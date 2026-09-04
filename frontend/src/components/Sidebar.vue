<template>
  <aside class="sidebar" :class="{ open: isOpen }">
    <div class="sidebar-brand">Rubium</div>
    
    <nav class="sidebar-nav">
      <router-link to="/" class="nav-link">
        <i class="fas fa-house"></i> Главная
      </router-link>
      <router-link to="/community" class="nav-link">
        <i class="fas fa-globe"></i> Каталог тетрадей
      </router-link>
      <template v-if="auth.isAuthenticated">
        <router-link to="/notebooks" class="nav-link">
          <i class="fas fa-book"></i> Мои тетради
        </router-link>
        <router-link to="/notifications" class="nav-link" @click="markAllRead">
          <i class="fas fa-bell"></i> Уведомления
          <span v-if="unreadCount > 0" class="badge">{{ unreadCount }}</span>
        </router-link>
        <router-link v-if="auth.isAdmin" to="/admin" class="nav-link admin-link">
          <i class="fas fa-shield-halved"></i> Админ
        </router-link>
      </template>
      <router-link to="/courses" class="nav-link">
        <i class="fas fa-graduation-cap"></i> Каталог курсов
      </router-link>
      <router-link to="/rubium_tech" class="nav-link">
        <i class="fas fa-laptop"></i> Rubium Tech
      </router-link>
    </nav>
    
    <div class="sidebar-footer">
      <template v-if="auth.isAuthenticated">
        <router-link to="/profile" class="user-link">
          <div class="user-avatar">
            <img v-if="auth.profile?.avatar_url" :src="auth.profile.avatar_url" alt="Аватар">
            <span v-else>{{ auth.userName[0]?.toUpperCase() || 'У' }}</span>
          </div>
          <div class="user-info">
            <div class="user-name">{{ auth.userName }}</div>
            <div class="user-email">{{ auth.user?.email }}</div>
          </div>
        </router-link>
        <button class="logout-btn" @click="auth.logout()">
          <i class="fas fa-sign-out-alt"></i>
        </button>
      </template>
      <template v-else>
        <router-link to="/login" class="login-btn">
          <i class="fas fa-sign-in-alt"></i> Войти
        </router-link>
      </template>
    </div>
  </aside>
  
  <div class="overlay" :class="{ show: isOpen }" @click="close"></div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { supabase } from '../api/supabase'

const auth = useAuthStore()
const isOpen = ref(false)
const unreadCount = ref(0)

const close = () => { isOpen.value = false }
const toggle = () => { isOpen.value = !isOpen.value }

defineExpose({ toggle })

async function loadUnreadCount() {
  if (!auth.profile?.id) return

  const { count } = await supabase
    .from('notifications')
    .select('*', { count: 'exact', head: true })
    .eq('user_id', auth.profile.id)
    .eq('is_read', false)

  unreadCount.value = count || 0
}

async function markAllRead() {
  if (!auth.profile?.id || unreadCount.value === 0) return

  await supabase
    .from('notifications')
    .update({ is_read: true })
    .eq('user_id', auth.profile.id)
    .eq('is_read', false)

  unreadCount.value = 0
}

watch(() => auth.profile?.id, (id) => {
  if (id) loadUnreadCount()
})

onMounted(() => {
  if (auth.profile?.id) loadUnreadCount()
})
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 240px;
  background: #0a0a0a;
  border-right: 1px solid rgba(255,255,255,0.06);
  display: flex;
  flex-direction: column;
  padding: 24px 16px 20px;
  z-index: 100;
  transition: transform 0.3s ease;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.sidebar-brand {
  font-size: 1.1rem;
  font-weight: 600;
  padding: 0 12px 28px;
  color: #ffffff;
  letter-spacing: -0.02em;
}

.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 9px 12px;
  border-radius: 8px;
  color: #737373;
  text-decoration: none;
  font-size: 0.88rem;
  font-weight: 500;
  transition: all 0.15s ease;
  position: relative;
}

.nav-link:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.nav-link.router-link-active {
  background: rgba(255,255,255,0.06);
  color: #ffffff;
}

.nav-link i {
  width: 18px;
  text-align: center;
  font-size: 0.85rem;
}

.badge {
  margin-left: auto;
  background: #ffffff;
  color: #0a0a0a;
  font-size: 0.65rem;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 10px;
  min-width: 18px;
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
  gap: 10px;
  text-decoration: none;
  color: inherit;
  min-width: 0;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.75rem;
  color: #a3a3a3;
  flex-shrink: 0;
  overflow: hidden;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 0.82rem;
  font-weight: 600;
  color: #e5e5e5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: 0.7rem;
  color: #525252;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 1px;
}

.logout-btn {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 0.85rem;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-btn:hover {
  color: #e5e5e5;
  background: rgba(255,255,255,0.04);
}

.login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 9px 16px;
  background: #ffffff;
  color: #0a0a0a;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  font-size: 0.85rem;
  width: 100%;
  transition: all 0.15s ease;
}

.login-btn:hover {
  background: #e5e5e5;
}

.overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
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