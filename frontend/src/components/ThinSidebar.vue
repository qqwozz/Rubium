<template>
  <aside class="thin-sidebar" :class="{ open: isOpen }">
    <div class="thin-nav">
      <router-link to="/" class="thin-link" data-tooltip="Главная">
        <i class="fas fa-house"></i>
      </router-link>
      
      <router-link to="/community" class="thin-link" data-tooltip="Каталог тетрадей">
        <i class="fas fa-globe"></i>
      </router-link>
      
      <template v-if="auth.isAuthenticated">
        <router-link to="/notebooks" class="thin-link" data-tooltip="Мои тетради">
          <i class="fas fa-book"></i>
        </router-link>
        <router-link to="/notifications" class="thin-link" data-tooltip="Уведомления" @click="markAllRead">
          <i class="fas fa-bell"></i>
          <span v-if="unreadCount > 0" class="thin-badge">{{ unreadCount }}</span>
        </router-link>
        <router-link v-if="auth.isAdmin" to="/admin" class="thin-link" data-tooltip="Админ">
          <i class="fas fa-shield-halved"></i>
        </router-link>
      </template>
      
      <router-link to="/courses" class="thin-link" data-tooltip="Курсы">
        <i class="fas fa-graduation-cap"></i>
      </router-link>
      
      <router-link to="/rubium_tech" class="thin-link" data-tooltip="Rubium Tech">
        <i class="fas fa-laptop"></i>
      </router-link>
    </div>
    
    <div class="thin-footer">
      <template v-if="auth.isAuthenticated">
        <router-link to="/profile" class="thin-link avatar-link" :data-tooltip="auth.userName">
          <div class="thin-avatar">
            <img v-if="auth.profile?.avatar_url" :src="auth.profile.avatar_url" alt="Аватар">
            <span v-else>{{ auth.userName[0]?.toUpperCase() || 'У' }}</span>
          </div>
        </router-link>
      </template>
      <template v-else>
        <router-link to="/login" class="thin-link" data-tooltip="Войти">
          <i class="fas fa-sign-in-alt"></i>
        </router-link>
      </template>
    </div>
  </aside>
  
  <div class="thin-overlay" :class="{ show: isOpen }" @click="close"></div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
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
.thin-sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 64px;
  background: #0a0a0a;
  border-right: 1px solid rgba(255,255,255,0.06);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0;
  z-index: 100;
  transition: transform 0.3s ease;
}

.thin-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  width: 100%;
}

.thin-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  width: 100%;
  padding-top: 12px;
  border-top: 1px solid rgba(255,255,255,0.06);
  margin-top: auto;
}

.thin-link {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #525252;
  font-size: 0.95rem;
  text-decoration: none;
  transition: all 0.15s ease;
  cursor: pointer;
  background: none;
  border: none;
}

.thin-link:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.thin-link.router-link-active {
  background: rgba(255,255,255,0.06);
  color: #ffffff;
}

.thin-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  background: #ffffff;
  color: #0a0a0a;
  font-size: 0.6rem;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 8px;
  min-width: 16px;
  text-align: center;
  line-height: 1.4;
}

/* Tooltip */
.thin-link[data-tooltip]::after {
  content: attr(data-tooltip);
  position: absolute;
  left: calc(100% + 10px);
  top: 50%;
  transform: translateY(-50%);
  background: #111111;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 500;
  color: #e5e5e5;
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: all 0.15s ease;
  pointer-events: none;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  z-index: 200;
}

.thin-link[data-tooltip]::before {
  content: '';
  position: absolute;
  left: calc(100% + 5px);
  top: 50%;
  transform: translateY(-50%);
  border: 5px solid transparent;
  border-right-color: #111111;
  opacity: 0;
  visibility: hidden;
  transition: all 0.15s ease;
  pointer-events: none;
  z-index: 200;
}

.thin-link:hover[data-tooltip]::after,
.thin-link:hover[data-tooltip]::before {
  opacity: 1;
  visibility: visible;
}

/* Avatar */
.avatar-link {
  padding: 0;
}

.thin-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 600;
  color: #a3a3a3;
  overflow: hidden;
}

.thin-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* Overlay for mobile */
.thin-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  z-index: 99;
}

.thin-overlay.show {
  display: block;
}

@media (max-width: 768px) {
  .thin-sidebar {
    transform: translateX(-100%);
    width: 240px;
  }
  
  .thin-sidebar.open {
    transform: translateX(0);
  }
}
</style>