<template>
  <div class="connections-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <button class="back-btn" @click="router.back()">
            <i class="fas fa-arrow-left"></i> Назад
          </button>
          <span class="page-title">Подписки</span>
        </header>

        <div class="content">
          <div v-if="loading" class="loading">
            <div class="spinner"></div>
            <span>Загружаем...</span>
          </div>

          <template v-else>
            <div class="connections-header">
              <h1>{{ getHeaderTitle() }}</h1>
            </div>

            <div class="tabs-row">
              <button 
                class="tab-btn" 
                :class="{ active: activeTab === 'subscribers' }"
                @click="activeTab = 'subscribers'"
              >
                Подписчики
              </button>
              <button 
                class="tab-btn" 
                :class="{ active: activeTab === 'subscriptions' }"
                @click="activeTab = 'subscriptions'"
              >
                Подписки
              </button>
              <button 
                class="tab-btn" 
                :class="{ active: activeTab === 'friends' }"
                @click="activeTab = 'friends'"
              >
                Друзья
              </button>
            </div>

            <div v-if="currentList.length === 0" class="empty-state">
              <div class="empty-icon"><i class="fas fa-users"></i></div>
              <p>{{ getEmptyText() }}</p>
            </div>

            <div v-else class="users-list">
              <div 
                v-for="user in currentList" 
                :key="user.id" 
                class="user-card"
                @click="router.push(`/user/${user.id}`)"
              >
                <div class="user-avatar">
                  <img v-if="user.avatar_url" :src="user.avatar_url" alt="Аватар">
                  <span v-else>{{ getInitial(user) }}</span>
                </div>
                <div class="user-info">
                  <div class="user-name">{{ getFullName(user) }}</div>
                  <div class="user-bio" v-if="user.bio">{{ truncateText(user.bio, 50) }}</div>
                </div>
                <i class="fas fa-chevron-right card-arrow"></i>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { supabase } from '../api/supabase'

const route = useRoute()
const router = useRouter()
const userId = route.params.id
const activeTab = ref('subscribers')
const subscribers = ref([])
const subscriptions = ref([])
const friends = ref([])
const loading = ref(true)
const sidebarRef = ref(null)

const currentList = computed(() => {
  if (activeTab.value === 'subscribers') return subscribers.value
  if (activeTab.value === 'subscriptions') return subscriptions.value
  return friends.value
})

function getHeaderTitle() {
  if (activeTab.value === 'subscribers') return 'Подписчики'
  if (activeTab.value === 'subscriptions') return 'Подписки'
  return 'Друзья'
}

function getEmptyText() {
  if (activeTab.value === 'subscribers') return 'Нет подписчиков'
  if (activeTab.value === 'subscriptions') return 'Нет подписок'
  return 'Нет друзей'
}

function getInitial(user) {
  const parts = [user?.first_name, user?.last_name].filter(Boolean)
  if (parts.length === 0) return '?'
  return parts.map(p => p[0]).join('').toUpperCase()
}

function getFullName(user) {
  const parts = [user?.first_name, user?.last_name].filter(Boolean)
  return parts.join(' ') || 'Пользователь'
}

function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text
}

async function loadUsers(ids) {
  if (!ids || ids.length === 0) return []
  
  const { data, error } = await supabase
    .from('rubium_users')
    .select('id, first_name, last_name, avatar_url, bio')
    .in('id', ids)

  if (error) {
    console.error(error)
    return []
  }

  return data || []
}

async function loadConnections() {
  loading.value = true
  try {
    // Подписчики
    const { data: subs } = await supabase
      .from('subscriptions')
      .select('subscriber_id')
      .eq('subscribed_to_id', userId)

    const subscriberIds = (subs || []).map(s => s.subscriber_id)
    subscribers.value = await loadUsers(subscriberIds)

    // Подписки
    const { data: subscr } = await supabase
      .from('subscriptions')
      .select('subscribed_to_id')
      .eq('subscriber_id', userId)

    const subscriptionIds = (subscr || []).map(s => s.subscribed_to_id)
    subscriptions.value = await loadUsers(subscriptionIds)

    // Друзья (взаимные)
    const friendIds = subscriberIds.filter(id => subscriptionIds.includes(id))
    friends.value = await loadUsers(friendIds)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(loadConnections)
</script>

<style scoped>
.connections-page {
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
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 48px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.back-btn {
  background: none;
  border: none;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: color 0.15s ease;
  padding: 0;
}

.back-btn:hover {
  color: #e5e5e5;
}

.page-title {
  flex: 1;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: #525252;
  text-align: center;
}

.content {
  max-width: 720px;
  margin: 0 auto;
  padding: 48px 48px 96px;
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

.connections-header {
  margin-bottom: 24px;
}

.connections-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.tabs-row {
  display: flex;
  gap: 2px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  padding: 4px;
  margin-bottom: 28px;
  width: fit-content;
}

.tab-btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.15s ease;
}

.tab-btn:hover {
  color: #e5e5e5;
}

.tab-btn.active {
  background: #ffffff;
  color: #0a0a0a;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #525252;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  border-radius: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #a3a3a3;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
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
  font-weight: 600;
  font-size: 0.95rem;
  color: #e5e5e5;
  margin-bottom: 2px;
}

.user-bio {
  font-size: 0.8rem;
  color: #525252;
  line-height: 1.4;
}

.card-arrow {
  color: #525252;
  font-size: 0.75rem;
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

  .tabs-row {
    width: 100%;
  }

  .tab-btn {
    flex: 1;
  }
}
</style>