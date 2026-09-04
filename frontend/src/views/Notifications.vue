<template>
  <div class="notifications-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <span class="page-title">Уведомления</span>
        </header>

        <div class="content">
          <div class="notifications-header">
            <h1>Уведомления</h1>
            <button v-if="notifications.length > 0" class="btn-mark-all" @click="markAllRead">
              <i class="fas fa-check-double"></i> Прочитать все
            </button>
          </div>

          <div v-if="loading" class="loading">
            <div class="spinner"></div>
            <span>Загружаем...</span>
          </div>

          <div v-else-if="notifications.length === 0" class="empty-state">
            <div class="empty-icon"><i class="fas fa-bell"></i></div>
            <h3>Нет уведомлений</h3>
            <p>Здесь будут появляться подписки, приглашения и системные сообщения</p>
          </div>

          <div v-else class="notifications-list">
            <div 
              v-for="notif in notifications" 
              :key="notif.id" 
              class="notification-item"
              :class="{ unread: !notif.is_read }"
              @click="handleNotificationClick(notif)"
            >
              <div class="notif-icon" :class="`type-${notif.type}`">
                <i :class="getNotifIcon(notif.type)"></i>
              </div>
              <div class="notif-content">
                <p class="notif-message">{{ notif.message }}</p>
                <span class="notif-time">{{ formatTime(notif.created_at) }}</span>
              </div>
              <i v-if="!notif.is_read" class="fas fa-circle unread-dot"></i>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { supabase } from '../api/supabase'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const notifications = ref([])
const loading = ref(true)
const sidebarRef = ref(null)

function getNotifIcon(type) {
  const icons = {
    subscription: 'fas fa-user-plus',
    collab_invite: 'fas fa-users',
    system: 'fas fa-gear',
    announcement: 'fas fa-bullhorn'
  }
  return icons[type] || 'fas fa-bell'
}

function formatTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return 'только что'
  if (diff < 3600000) return `${Math.floor(diff / 60000)} мин назад`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} ч назад`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} дн назад`
  
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
}

async function markAllRead() {
  if (!auth.profile?.id) return

  await supabase
    .from('notifications')
    .update({ is_read: true })
    .eq('user_id', auth.profile.id)
    .eq('is_read', false)

  notifications.value = notifications.value.map(n => ({ ...n, is_read: true }))
}

async function handleNotificationClick(notif) {
  if (!notif.is_read) {
    await supabase
      .from('notifications')
      .update({ is_read: true })
      .eq('id', notif.id)

    notif.is_read = true
  }

  if (notif.from_user_id && notif.type === 'subscription') {
    router.push(`/user/${notif.from_user_id}`)
  } else if (notif.notebook_id && notif.type === 'collab_invite') {
    router.push(`/notebook/${notif.notebook_id}`)
  }
}

async function loadNotifications() {
  if (!auth.profile?.id) return

  loading.value = true
  try {
    const { data, error } = await supabase
      .from('notifications')
      .select('*')
      .eq('user_id', auth.profile.id)
      .order('created_at', { ascending: false })
      .limit(50)

    if (error) throw error
    notifications.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (auth.profile?.id) {
    loadNotifications()
  }
})
</script>

<style scoped>
.notifications-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.page-body {
  display: flex;
  flex: 1;
}

.main-content {
  margin-left: 240px;
  flex: 1;
}

.topbar {
  padding: 20px 48px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.page-title {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: #525252;
}

.content {
  max-width: 720px;
  margin: 0 auto;
  padding: 48px 48px 96px;
}

.notifications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  flex-wrap: wrap;
  gap: 12px;
}

.notifications-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.btn-mark-all {
  padding: 8px 14px;
  background: transparent;
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 8px;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-mark-all:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 20px;
  color: #737373;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 2px solid rgba(255,255,255,0.06);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #737373;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #a3a3a3;
}

.empty-state h3 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 6px;
}

.empty-state p {
  font-size: 0.9rem;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.notification-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
}

.notification-item:hover {
  background: rgba(255,255,255,0.03);
}

.notification-item.unread {
  background: rgba(255,255,255,0.02);
}

.notif-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  flex-shrink: 0;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  color: #a3a3a3;
}

.notif-content {
  flex: 1;
  min-width: 0;
}

.notif-message {
  font-size: 0.9rem;
  color: #e5e5e5;
  margin-bottom: 3px;
  line-height: 1.4;
}

.notif-time {
  font-size: 0.75rem;
  color: #525252;
}

.unread-dot {
  font-size: 0.45rem;
  color: #ffffff;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .page-body {
    display: block;
  }

  .main-content {
    margin-left: 0;
  }

  .topbar {
    display: none;
  }

  .content {
    padding: 32px 24px 64px;
  }
}
</style>