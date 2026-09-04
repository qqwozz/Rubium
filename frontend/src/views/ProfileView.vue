<template>
  <div class="profile-view-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <button class="back-btn" @click="router.back()">
            <i class="fas fa-arrow-left"></i> Назад
          </button>
          <span class="page-title">Профиль</span>
        </header>

        <div class="content">
          <div v-if="loading" class="loading">
            <div class="spinner"></div>
            <span>Загружаем профиль...</span>
          </div>

          <template v-else-if="profile">
            <div class="profile-header">
              <div class="avatar">
                <img v-if="profile.avatar_url" :src="profile.avatar_url" alt="Аватар">
                <span v-else>{{ getInitials() }}</span>
              </div>

              <div class="profile-info">
                <h1>
                  {{ getFullName() }}
                  <i v-if="isDeveloper" class="fas fa-check-circle verified-badge-icon" title="Разработчик"></i>
                </h1>
                <p class="profile-email">{{ profile.email }}</p>
              </div>

              <div v-if="!isOwnProfile" class="subscribe-block">
                <button 
                  class="btn-subscribe" 
                  :class="{ subscribed: isSubscribed }"
                  @click="toggleSubscription"
                  :disabled="subscribing"
                >
                  <i :class="subscribing ? 'fas fa-spinner fa-spin' : (isSubscribed ? 'fas fa-check' : 'fas fa-plus')"></i>
                  {{ isMutual ? 'Друзья' : (isSubscribed ? 'Вы подписаны' : 'Подписаться') }}
                </button>
              </div>
            </div>

            <div v-if="profile.bio" class="bio-section">
              <p>{{ profile.bio }}</p>
            </div>

            <div class="stats-block">
              <div class="stats-row">
                <span class="stat-item">
                  <span class="stat-value">{{ subscribersCount }}</span>
                  <span class="stat-label">подписчики</span>
                </span>
                <span class="stat-divider"></span>
                <span class="stat-item">
                  <span class="stat-value">{{ subscriptionsCount }}</span>
                  <span class="stat-label">подписки</span>
                </span>
                <span class="stat-divider"></span>
                <span class="stat-item">
                  <span class="stat-value">{{ friendsCount }}</span>
                  <span class="stat-label">друзья</span>
                </span>
              </div>

              <div class="stats-row">
                <span class="stat-item">
                  <span class="stat-value">{{ notebooks.length }}</span>
                  <span class="stat-label">тетрадей</span>
                </span>
                <span class="stat-divider"></span>
                <span class="stat-item">
                  <span class="stat-value">{{ averageRating }}</span>
                  <span class="stat-label">рейтинг</span>
                </span>
                <span class="stat-divider"></span>
                <span class="stat-item">
                  <span class="stat-value">{{ totalViews }}</span>
                  <span class="stat-label">просмотров</span>
                </span>
              </div>
            </div>

            <div class="notebooks-toggle">
              <button 
                class="toggle-btn" 
                :class="{ active: activeTab === 'own' }"
                @click="activeTab = 'own'"
              >
                Тетради
              </button>
              <button 
                class="toggle-btn" 
                :class="{ active: activeTab === 'saved' }"
                @click="activeTab = 'saved'"
              >
                Сохранённые
              </button>
            </div>

            <div v-if="activeTab === 'saved'" class="notebooks-section">
              <div v-if="savedNotebooks.length > 0" class="notebooks-grid">
                <div 
                  v-for="nb in savedNotebooks" 
                  :key="nb.id" 
                  class="notebook-card"
                  @click="router.push(`/notebook/${nb.id}`)"
                >
                  <div class="notebook-color" :style="{ background: nb.color || '#525252' }"></div>
                  <div class="notebook-info">
                    <div class="notebook-title">{{ nb.title }}</div>
                    <div v-if="nb.description" class="notebook-description">
                      {{ truncateText(nb.description, 80) }}
                    </div>
                    <div class="notebook-stats">
                      <span><i class="fas fa-star"></i> {{ formatRating(nb.average_rating) }}</span>
                      <span><i class="fas fa-eye"></i> {{ nb.views_count || 0 }}</span>
                    </div>
                  </div>
                  <i class="fas fa-chevron-right card-arrow"></i>
                </div>
              </div>
              <div v-else class="empty">
                <div class="empty-icon"><i class="fas fa-bookmark"></i></div>
                <p>Нет сохранённых тетрадей</p>
              </div>
            </div>

            <div v-else class="notebooks-section">
              <div v-if="notebooks.length > 0" class="notebooks-grid">
                <div 
                  v-for="nb in notebooks" 
                  :key="nb.id" 
                  class="notebook-card"
                  @click="router.push(`/notebook/${nb.id}`)"
                >
                  <div class="notebook-color" :style="{ background: nb.color || '#525252' }"></div>
                  <div class="notebook-info">
                    <div class="notebook-title">{{ nb.title }}</div>
                    <div v-if="nb.description" class="notebook-description">
                      {{ truncateText(nb.description, 80) }}
                    </div>
                    <div class="notebook-stats">
                      <span><i class="fas fa-star"></i> {{ formatRating(nb.average_rating) }}</span>
                      <span><i class="fas fa-eye"></i> {{ nb.views_count || 0 }}</span>
                    </div>
                  </div>
                  <i class="fas fa-chevron-right card-arrow"></i>
                </div>
              </div>

              <div v-else class="empty">
                <div class="empty-icon"><i class="fas fa-book"></i></div>
                <p>Нет публичных тетрадей</p>
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
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const profile = ref(null)
const notebooks = ref([])
const savedNotebooks = ref([])
const loading = ref(true)
const sidebarRef = ref(null)
const activeTab = ref('saved')

const isSubscribed = ref(false)
const isMutual = ref(false)
const subscribing = ref(false)

const subscribersCount = ref(0)
const subscriptionsCount = ref(0)
const friendsCount = ref(0)

const DEVELOPER_EMAILS = ['nsdmlk@yandex.ru', 'offconix@gmail.com', 'oleg.veter.08@mail.ru']

const isDeveloper = computed(() => {
  return DEVELOPER_EMAILS.includes(profile.value?.email)
})

const currentUserId = computed(() => auth.profile?.id)
const isOwnProfile = computed(() => currentUserId.value === profile.value?.id)

const averageRating = computed(() => {
  const rated = notebooks.value.filter(nb => nb.average_rating > 0)
  if (rated.length === 0) return '0.0'
  const sum = rated.reduce((acc, nb) => acc + nb.average_rating, 0)
  return (sum / rated.length).toFixed(1)
})

const totalViews = computed(() => {
  return notebooks.value.reduce((sum, nb) => sum + (nb.views_count || 0), 0)
})

function getFullName() {
  const parts = [profile.value?.first_name, profile.value?.last_name].filter(Boolean)
  return parts.join(' ') || 'Пользователь'
}

function getInitials() {
  const parts = [profile.value?.first_name, profile.value?.last_name].filter(Boolean)
  if (parts.length === 0) return '?'
  return parts.map(p => p[0]).join('').toUpperCase()
}

function formatRating(rating) {
  return Number(rating || 0).toFixed(1)
}

function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text
}

async function checkSubscription() {
  if (!currentUserId.value || !profile.value?.id) return
  if (currentUserId.value === profile.value.id) return

  const { data: sub } = await supabase
    .from('subscriptions')
    .select('subscriber_id')
    .eq('subscriber_id', currentUserId.value)
    .eq('subscribed_to_id', profile.value.id)
    .maybeSingle()

  isSubscribed.value = !!sub

  if (isSubscribed.value) {
    const { data: mutual } = await supabase
      .from('subscriptions')
      .select('subscriber_id')
      .eq('subscriber_id', profile.value.id)
      .eq('subscribed_to_id', currentUserId.value)
      .maybeSingle()

    isMutual.value = !!mutual
  }
}

async function loadStats() {
  if (!profile.value?.id) return

  const { count: subs } = await supabase
    .from('subscriptions')
    .select('*', { count: 'exact', head: true })
    .eq('subscribed_to_id', profile.value.id)

  subscribersCount.value = subs || 0

  const { count: subscr } = await supabase
    .from('subscriptions')
    .select('*', { count: 'exact', head: true })
    .eq('subscriber_id', profile.value.id)

  subscriptionsCount.value = subscr || 0

  const { data: mutual } = await supabase
    .from('subscriptions')
    .select('subscriber_id, subscribed_to_id')
    .or(`subscriber_id.eq.${profile.value.id},subscribed_to_id.eq.${profile.value.id}`)

  if (mutual) {
    const userSubs = mutual.filter(s => s.subscriber_id === profile.value.id).map(s => s.subscribed_to_id)
    const userSubscribers = mutual.filter(s => s.subscribed_to_id === profile.value.id).map(s => s.subscriber_id)
    
    const friends = userSubs.filter(id => userSubscribers.includes(id))
    friendsCount.value = friends.length
  }
}

async function toggleSubscription() {
  if (!auth.isAuthenticated) {
    router.push('/login')
    return
  }

  if (!currentUserId.value || !profile.value?.id || subscribing.value) return
  if (currentUserId.value === profile.value.id) return

  subscribing.value = true

  try {
    if (isSubscribed.value) {
      const { error } = await supabase
        .from('subscriptions')
        .delete()
        .eq('subscriber_id', currentUserId.value)
        .eq('subscribed_to_id', profile.value.id)

      if (error) throw error

      isSubscribed.value = false
      isMutual.value = false
      await loadStats()
    } else {
      const { error: insertError } = await supabase
        .from('subscriptions')
        .insert({
          subscriber_id: currentUserId.value,
          subscribed_to_id: profile.value.id
        })

      if (insertError) throw insertError

      isSubscribed.value = true

      const { error: notifError } = await supabase
        .from('notifications')
        .insert({
          user_id: profile.value.id,
          type: 'subscription',
          from_user_id: currentUserId.value,
          message: `${getFullName()} подписался на вас`
        })

      if (notifError) console.error(notifError)

      const { data: mutual } = await supabase
        .from('subscriptions')
        .select('subscriber_id')
        .eq('subscriber_id', profile.value.id)
        .eq('subscribed_to_id', currentUserId.value)
        .maybeSingle()

      isMutual.value = !!mutual
      await loadStats()
    }
  } catch (e) {
    console.error(e)
  } finally {
    subscribing.value = false
  }
}

async function loadProfile() {
  try {
    if (!auth.profile && auth.isAuthenticated) {
      await auth.loadProfile()
    }

    const { data: userData, error: userError } = await supabase
      .from('rubium_users')
      .select('id, first_name, last_name, avatar_url, bio, email, saved_notebooks, created_at')
      .eq('id', route.params.id)
      .single()

    if (userError) throw userError
    profile.value = userData

    const { data: notebooksData, error: notebooksError } = await supabase
      .from('notebooks')
      .select('id, title, description, color, average_rating, ratings_count, views_count, updated_at')
      .eq('user_id', route.params.id)
      .eq('is_public', true)
      .order('updated_at', { ascending: false })

    if (notebooksError) throw notebooksError
    notebooks.value = notebooksData || []

    const savedIds = userData?.saved_notebooks || []
    if (savedIds.length > 0) {
      const { data: savedData, error: savedError } = await supabase
        .from('notebooks')
        .select('id, title, description, color, average_rating, ratings_count, views_count, updated_at')
        .in('id', savedIds)
        .eq('is_public', true)

      if (savedError) throw savedError
      savedNotebooks.value = savedData || []
    }

    await Promise.all([checkSubscription(), loadStats()])
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(loadProfile)
</script>

<style scoped>
.profile-view-page {
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

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 28px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.6rem;
  font-weight: 600;
  color: #a3a3a3;
  flex-shrink: 0;
  overflow: hidden;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-info {
  flex: 1;
  min-width: 0;
}

.profile-info h1 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: -0.02em;
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
}

.verified-badge-icon {
  color: #a3a3a3;
  font-size: 1.05rem;
  flex-shrink: 0;
}

.profile-email {
  color: #525252;
  font-size: 0.85rem;
}

.subscribe-block {
  flex-shrink: 0;
}

.btn-subscribe {
  padding: 10px 18px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.btn-subscribe:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-subscribe.subscribed {
  background: transparent;
  color: #fafafa;
  border-color: rgba(255,255,255,0.12);
}

.btn-subscribe.subscribed:hover:not(:disabled) {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.2);
}

.btn-subscribe:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.stats-block {
  margin-bottom: 28px;
  display: flex;
  align-items: center;
  gap: 5px;
  flex-direction: column;
}

.stats-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 10px;
}

.stat-item {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-shrink: 0;
}

.stat-value {
  font-size: 0.9rem;
  font-weight: 600;
  color: #737373;
}

.stat-label {
  font-size: 0.75rem;
  color: #525252;
}

.stat-divider {
  width: 1px;
  height: 14px;
  background: rgba(255,255,255,0.08);
  flex-shrink: 0;
}

.bio-section {
  margin-bottom: 28px;
}

.bio-section p {
  color: #a3a3a3;
  font-size: 0.95rem;
  line-height: 1.7;
}

.notebooks-toggle {
  display: flex;
  gap: 2px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  padding: 4px;
  margin-bottom: 24px;
  width: fit-content;
}

.toggle-btn {
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

.toggle-btn:hover {
  color: #e5e5e5;
}

.toggle-btn.active {
  background: #ffffff;
  color: #0a0a0a;
}

.notebooks-section {
  margin-bottom: 40px;
}

.notebooks-grid {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notebook-card {
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

.notebook-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.notebook-color {
  width: 4px;
  height: auto;
  border-radius: 2px;
  flex-shrink: 0;
  align-self: stretch;
}

.notebook-info {
  flex: 1;
  min-width: 0;
}

.notebook-title {
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 4px;
  font-size: 0.95rem;
}

.notebook-description {
  font-size: 0.8rem;
  color: #525252;
  margin-bottom: 6px;
  line-height: 1.4;
}

.notebook-stats {
  display: flex;
  gap: 12px;
  font-size: 0.8rem;
  color: #525252;
}

.notebook-stats i {
  margin-right: 4px;
  font-size: 0.7rem;
}

.card-arrow {
  color: #525252;
  font-size: 0.75rem;
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.notebook-card:hover .card-arrow {
  color: #a3a3a3;
}

.empty {
  text-align: center;
  padding: 60px 20px;
  color: #525252;
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

  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 16px;
  }

  .avatar {
    width: 72px;
    height: 72px;
    font-size: 1.4rem;
  }

  .profile-info h1 {
    font-size: 1.3rem;
    justify-content: center;
  }

  .bio-section {
    display: flex;
    justify-content: center;
  }

  .notebooks-toggle {
    width: 100%;
  }

  .toggle-btn {
    flex: 1;
  }

  .subscribe-block {
    width: 100%;
    display: flex;
    justify-content: center;
  }

  .stats-row {
    flex-wrap: wrap;
    gap: 10px;
  }
}
</style>