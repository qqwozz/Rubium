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
            </div>

            <div v-if="profile.bio" class="bio-section">
              <p>{{ profile.bio }}</p>
            </div>

            <div class="stats-section">
              <div class="stat">
                <span class="stat-value">{{ notebooks.length }}</span>
                <span class="stat-label">Тетрадей</span>
              </div>
              <div class="stat-divider"></div>
              <div class="stat">
                <span class="stat-value">{{ averageRating }}</span>
                <span class="stat-label">Рейтинг</span>
              </div>
              <div class="stat-divider"></div>
              <div class="stat">
                <span class="stat-value">{{ totalViews }}</span>
                <span class="stat-label">Просмотров</span>
              </div>
            </div>

            <div class="notebooks-section">
              <h2>Публичные тетради</h2>
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

const route = useRoute()
const router = useRouter()
const profile = ref(null)
const notebooks = ref([])
const loading = ref(true)
const sidebarRef = ref(null)

const DEVELOPER_EMAILS = ['nsdmlk@yandex.ru', 'offconix@gmail.com', 'oleg.veter.08@mail.ru']

const isDeveloper = computed(() => {
  return DEVELOPER_EMAILS.includes(profile.value?.email)
})

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

async function loadProfile() {
  try {
    const { data: userData, error: userError } = await supabase
      .from('rubium_users')
      .select('id, first_name, last_name, avatar_url, bio, email, created_at')
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
  margin-bottom: 32px;
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

.bio-section {
  margin-bottom: 32px;
}

.bio-section p {
  color: #a3a3a3;
  font-size: 0.95rem;
  line-height: 1.7;
  display: block;
  justify-content: center;
}

.stats-section {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 40px;
  padding: 20px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  text-align: center;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: -0.01em;
}

.stat-label {
  font-size: 0.75rem;
  color: #525252;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 500;
}

.stat-divider {
  width: 1px;
  height: 36px;
  background: rgba(255,255,255,0.06);
  flex-shrink: 0;
}

.notebooks-section h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 14px;
  letter-spacing: -0.01em;
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

  .stats-section {
    gap: 12px;
    padding: 16px;
  }

  .stat-value {
    font-size: 1.3rem;
  }

  .stat-label {
    font-size: 0.68rem;
  }
}
</style>