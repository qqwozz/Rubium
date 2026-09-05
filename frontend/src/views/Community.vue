<template>
  <div class="community-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <span class="page-title">Каталог тетрадей</span>
        </header>

        <div class="content">
          <div class="community-header">
            <h1>Публичные тетради</h1>
            <p>Находи конспекты других учеников и делись своими</p>
          </div>

          <div class="controls-row">
            <div class="search-bar">
              <i class="fas fa-search search-icon"></i>
              <input v-model="searchQuery" type="text" placeholder="Поиск по названию или тегам..." @keydown.enter="loadNotebooks">
            </div>

            <div class="sort-row">
              <button v-for="sort in sorts" :key="sort.value" class="sort-btn" :class="{ active: currentSort === sort.value }" @click="currentSort = sort.value; loadNotebooks()">
                {{ sort.label }}
              </button>
            </div>
          </div>

          <div v-if="loading" class="loading">
            <div class="spinner"></div>
            <span>Загружаем тетради...</span>
          </div>

          <div v-else-if="notebooks.length === 0" class="empty-state">
            <div class="empty-icon"><i class="fas fa-book-open"></i></div>
            <h3>Нет публичных тетрадей</h3>
            <p>Стань первым — создай тетрадь и опубликуй её!</p>
          </div>

          <div v-else class="notebooks-grid">
            <div v-for="notebook in notebooks" :key="notebook.id" class="notebook-card" @click="openNotebook(notebook)">
              <div class="notebook-color" :style="{ background: notebook.color || '#525252' }"></div>
              <div class="notebook-info">
                <div class="notebook-title-row">
                  <div class="notebook-title">{{ notebook.title }}</div>
                  <i v-if="notebook.is_verified" class="fas fa-check-circle verified-badge-icon" title="От разработчика"></i>
                </div>
                <div v-if="notebook.description" class="notebook-description">
                  {{ truncateText(notebook.description, 60) }}
                </div>
                <div class="notebook-footer">
                  <div class="notebook-author">
                    <i class="fas fa-user"></i> {{ getAuthorFullName(notebook.author) }}
                  </div>
                  <div class="notebook-rating">
                    <i class="fas fa-star"></i> {{ formatRating(notebook.average_rating) }}
                    <span class="rating-count">({{ notebook.ratings_count || 0 }})</span>
                  </div>
                  <div v-if="notebook.tags && notebook.tags.length" class="notebook-tags">
                    <span v-for="tag in notebook.tags.slice(0, 1)" :key="tag" class="tag">{{ tag }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <Transition name="modal">
        <div v-if="selectedNotebook" class="modal" @click.self="selectedNotebook = null">
          <div class="modal-card">
            <div class="modal-header">
              <div class="modal-color" :style="{ background: selectedNotebook.color || '#525252' }"></div>
              <div class="modal-title-block">
                <h2>{{ selectedNotebook.title }}</h2>
              </div>
              <button class="modal-close" @click="selectedNotebook = null">
                <i class="fas fa-times"></i>
              </button>
            </div>

            <div class="modal-body">
              <div v-if="selectedNotebook.description" class="modal-description">
                <p>{{ selectedNotebook.description }}</p>
              </div>

              <div 
                class="modal-author" 
                @click="goToAuthorProfile(selectedNotebook.author)"
                role="button"
                tabindex="0"
              >
                <div class="author-avatar">
                  <img v-if="getAuthorAvatar(selectedNotebook.author)" :src="getAuthorAvatar(selectedNotebook.author)" alt="Аватар">
                  <span v-else>{{ getInitial(selectedNotebook.author) }}</span>
                </div>
                <div>
                  <div class="author-name">
                    {{ getAuthorFullName(selectedNotebook.author) }}
                    <i v-if="selectedNotebook.is_verified" class="fas fa-check-circle verified-badge-icon" title="От разработчика"></i>
                  </div>
                  <div class="author-email">{{ getAuthorEmail(selectedNotebook.author) }}</div>
                </div>
                <i class="fas fa-chevron-right author-arrow"></i>
              </div>

              <div class="modal-stats">
                <span class="stat-item"><i class="fas fa-star"></i> {{ formatRating(selectedNotebook.average_rating) }} ({{ selectedNotebook.ratings_count || 0 }})</span>
                <span class="stat-item"><i class="fas fa-eye"></i> {{ selectedNotebook.views_count || 0 }}</span>
              </div>

              <div v-if="selectedNotebook.tags?.length" class="modal-tags">
                <span v-for="tag in selectedNotebook.tags" :key="tag" class="tag">{{ tag }}</span>
              </div>

              <div class="modal-actions">
                <button class="btn-open" @click="incrementViews">
                  <i class="fas fa-book-open"></i> Открыть
                </button>
                <button v-if="!isOwnNotebook" class="btn-save" @click="toggleSave" :class="{ saved: isSaved }" :disabled="saving">
                  <i :class="saving ? 'fas fa-spinner fa-spin' : (isSaved ? 'fas fa-bookmark' : 'far fa-bookmark')"></i>
                  {{ isSaved ? 'Сохранено' : 'Сохранить' }}
                </button>
                <button v-if="!isOwnNotebook" class="btn-rate" @click="openRateModal">
                  <i class="fas fa-star"></i> Оценить
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <Transition name="modal">
        <div v-if="showRateModal" class="modal" @click.self="showRateModal = false">
          <div class="modal-card rate-modal">
            <h2>Оценить тетрадь</h2>
            <div class="rate-stars">
              <button v-for="star in 5" :key="star" class="star-btn" :class="{ active: star <= rateValue }" @click="rateValue = star">
                <i class="fas fa-star"></i>
              </button>
            </div>
            <div class="modal-actions">
              <button class="btn-submit" @click="submitRate">Оценить</button>
              <button class="btn-cancel" @click="showRateModal = false">Отмена</button>
            </div>
          </div>
        </div>
      </Transition>

      <Transition name="modal">
        <div v-if="notification" class="modal" @click.self="notification = null">
          <div class="modal-card notification-modal">
            <div class="notification-icon" :class="notification.type">
              <i :class="notification.type === 'success' ? 'fas fa-check-circle' : 'fas fa-exclamation-circle'"></i>
            </div>
            <p>{{ notification.message }}</p>
            <button class="btn-submit" @click="notification = null">ОК</button>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { apiFetch } from '../api/client'
import { useAuthStore } from '../stores/auth'
import { supabase } from '../api/supabase'

const router = useRouter()
const auth = useAuthStore()
const notebooks = ref([])
const loading = ref(true)
const searchQuery = ref('')
const currentSort = ref('rating')
const selectedNotebook = ref(null)
const showRateModal = ref(false)
const rateValue = ref(5)
const notification = ref(null)
const saving = ref(false)
const sidebarRef = ref(null)

const DEVELOPER_EMAILS = ['nsdmlk@yandex.ru', 'offconix@gmail.com', 'oleg.veter.08@mail.ru']

const sorts = [
  { value: 'rating', label: 'По рейтингу' },
  { value: 'newest', label: 'Новые' },
  { value: 'popular', label: 'Популярные' }
]

const isSaved = computed(() => {
  return auth.profile?.saved_notebooks?.includes(selectedNotebook.value?.id)
})

const isOwnNotebook = computed(() => {
  return selectedNotebook.value?.author?.id === auth.profile?.id
})

function getInitial(author) {
  if (!author) return 'А'
  return (author.first_name || author.email || 'А')[0].toUpperCase()
}

function getAuthorAvatar(author) {
  return author?.avatar_url || ''
}

function getAuthorFullName(author) {
  if (!author) return 'Автор'
  const parts = [author.first_name, author.last_name].filter(Boolean)
  return parts.join(' ') || author.email?.split('@')[0] || 'Автор'
}

function getAuthorEmail(author) {
  return author?.email || ''
}

function formatRating(rating) {
  return Number(rating || 0).toFixed(1)
}

function truncateText(text, maxLength) {
  if (!text) return ''
  return text.length > maxLength ? text.slice(0, maxLength) + '...' : text
}

function showNotification(message, type = 'success') {
  notification.value = { message, type }
}

async function toggleSave() {
  if (!selectedNotebook.value?.id) return
  if (saving.value) return

  if (!auth.isAuthenticated) {
    showNotification('Войди в аккаунт, чтобы сохранить тетрадь', 'error')
    return
  }

  saving.value = true
  try {
    const userId = auth.profile?.id
    if (!userId) throw new Error('Нет профиля')

    const saved = auth.profile?.saved_notebooks || []
    let updated

    if (saved.includes(selectedNotebook.value.id)) {
      updated = saved.filter(id => id !== selectedNotebook.value.id)
    } else {
      updated = [...saved, selectedNotebook.value.id]
    }

    const { error } = await supabase
      .from('rubium_users')
      .update({ saved_notebooks: updated })
      .eq('id', userId)

    if (error) throw error

    await auth.loadProfile()
  } catch (e) {
    console.error(e)
    showNotification('Ошибка при сохранении', 'error')
  } finally {
    saving.value = false
  }
}

function goToAuthorProfile(author) {
  if (!author?.id) return
  selectedNotebook.value = null
  router.push(`/user/${author.id}`)
}

async function loadNotebooks() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (searchQuery.value.trim()) params.append('search', searchQuery.value.trim())
    params.append('sort', currentSort.value)

    const data = await apiFetch(`/notebooks/community?${params}`)

    notebooks.value = (data.notebooks || []).map(n => ({
      ...n,
      is_verified: DEVELOPER_EMAILS.includes(n.author?.email)
    }))
  } catch (e) {
    console.error(e)
    showNotification('Ошибка при загрузке тетрадей', 'error')
  } finally {
    loading.value = false
  }
}

function openNotebook(notebook) {
  selectedNotebook.value = notebook
}

function openRateModal() {
  rateValue.value = 5
  showRateModal.value = true
}

async function incrementViews() {
  if (!selectedNotebook.value?.id) return

  if (!auth.isAuthenticated) {
    showNotification('Войди в аккаунт, чтобы открыть тетрадь', 'error')
    return
  }

  try {
    await apiFetch(`/notebooks/${selectedNotebook.value.id}/view`, {
      method: 'POST'
    })
    router.push(`/notebook/${selectedNotebook.value.id}`)
  } catch (e) {
    console.error(e)
    showNotification('Ошибка при открытии', 'error')
  }
}

async function submitRate() {
  if (!selectedNotebook.value) return

  if (!auth.isAuthenticated) {
    showNotification('Войди в аккаунт, чтобы оценить тетрадь', 'error')
    return
  }

  try {
    await apiFetch(`/notebooks/${selectedNotebook.value.id}/rate`, {
      method: 'POST',
      body: JSON.stringify({ rating: rateValue.value })
    })

    showRateModal.value = false
    selectedNotebook.value = null
    await loadNotebooks()
    showNotification('Спасибо за оценку!')
  } catch (e) {
    console.error(e)
    showNotification(e.message || 'Ошибка при оценке', 'error')
  }
}

onMounted(loadNotebooks)
</script>

<style scoped>
.community-page {
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
  max-width: 900px;
  margin: 0 auto;
  padding: 48px 48px 96px;
}

.community-header {
  margin-bottom: 40px;
}

.community-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 8px;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.community-header p {
  color: #737373;
  font-size: 0.95rem;
}

.controls-row {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
  flex-wrap: wrap;
  align-items: center;
}

.search-bar {
  flex: 1;
  min-width: 250px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: #525252;
  font-size: 0.85rem;
}

.search-bar input {
  width: 100%;
  padding: 12px 16px 12px 40px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  transition: all 0.2s ease;
}

.search-bar input:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.search-bar input::placeholder {
  color: #525252;
}

.sort-row {
  display: flex;
  gap: 2px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  padding: 4px;
}

.sort-btn {
  padding: 8px 14px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 500;
  transition: all 0.15s ease;
}

.sort-btn:hover {
  color: #e5e5e5;
}

.sort-btn.active {
  background: #ffffff;
  color: #0a0a0a;
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

.notebooks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
  align-items: stretch;
}

.notebook-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  gap: 16px;
  height: 100%;
  min-height: 150px;
}

.notebook-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.notebook-color {
  width: 4px;
  border-radius: 2px;
  flex-shrink: 0;
  align-self: stretch;
}

.notebook-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.notebook-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 6px;
}

.notebook-title {
  font-size: 1.05rem;
  font-weight: 600;
  line-height: 1.3;
  color: #e5e5e5;
}

.verified-badge-icon {
  color: #a3a3a3;
  font-size: 0.8rem;
  flex-shrink: 0;
  margin-top: 2px;
}

.notebook-description {
  font-size: 0.82rem;
  color: #737373;
  line-height: 1.4;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.notebook-footer {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-top: auto;
  flex-wrap: wrap;
}

.notebook-author {
  font-size: 0.78rem;
  color: #525252;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.notebook-author i {
  font-size: 0.7rem;
}

.notebook-rating {
  font-size: 0.8rem;
  color: #a3a3a3;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.notebook-rating i {
  font-size: 0.7rem;
  color: #a3a3a3;
}

.rating-count {
  color: #525252;
  font-size: 0.75rem;
}

.notebook-tags {
  display: flex;
  margin-left: auto;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: flex-end;
  max-width: 60%;
}

.tag {
  font-size: 0.7rem;
  padding: 3px 8px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  color: #737373;
  font-weight: 500;
  white-space: nowrap;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-card {
  background: #111111;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
  padding: 28px;
  width: 100%;
  max-width: 520px;
}

.modal-header {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 20px;
}

.modal-color {
  width: 4px;
  height: stretch;
  border-radius: 2px;
  flex-shrink: 0;
}

.modal-title-block {
  flex: 1;
}

.modal-title-block h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #ffffff;
  line-height: 1.3;
}

.modal-close {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 1rem;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.modal-close:hover {
  color: #e5e5e5;
  background: rgba(255,255,255,0.04);
}

.modal-description {
  margin-bottom: 20px;
  padding: 14px;
  background: rgba(255,255,255,0.02);
  border-radius: 10px;
  border: 1px solid rgba(255,255,255,0.06);
}

.modal-description p {
  color: #a3a3a3;
  font-size: 0.9rem;
  line-height: 1.6;
}

.modal-author {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.author-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  font-weight: 600;
  color: #a3a3a3;
  flex-shrink: 0;
  overflow: hidden;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.author-name {
  font-weight: 600;
  font-size: 0.95rem;
  color: #e5e5e5;
  display: flex;
  align-items: center;
  gap: 6px;
}

.author-name .verified-badge-icon {
  font-size: 0.75rem;
  margin-top: 0;
}

.author-email {
  font-size: 0.78rem;
  color: #525252;
  margin-top: 2px;
}

.author-arrow {
  margin-left: auto;
  color: #525252;
  font-size: 0.75rem;
  flex-shrink: 0;
}

.modal-stats {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.8rem;
  color: #737373;
}

.stat-item i {
  font-size: 0.75rem;
  color: #525252;
}

.modal-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-open,
.btn-save,
.btn-rate {
  padding: 10px 18px;
  border: 1px solid transparent;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-open {
  background: #ffffff;
  color: #0a0a0a;
  border-color: #ffffff;
}

.btn-open:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-save {
  background: transparent;
  color: #fafafa;
  border-color: rgba(255,255,255,0.12);
}

.btn-save:hover:not(:disabled) {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.2);
}

.btn-save.saved {
  background: rgba(255,255,255,0.06);
  border-color: rgba(255,255,255,0.2);
  color: #ffffff;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-rate {
  background: transparent;
  color: #fafafa;
  border-color: rgba(255,255,255,0.12);
}

.btn-rate:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.2);
}

.rate-modal {
  max-width: 360px;
  text-align: center;
}

.rate-modal h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 24px;
}

.rate-stars {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 24px;
}

.star-btn {
  background: none;
  border: none;
  font-size: 1.8rem;
  color: #404040;
  cursor: pointer;
  transition: all 0.15s ease;
  padding: 0;
}

.star-btn.active {
  color: #ffffff;
}

.star-btn:hover {
  transform: scale(1.1);
}

.btn-cancel {
  padding: 10px 18px;
  background: rgba(255,255,255,0.04);
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.85rem;
}

.btn-cancel:hover {
  background: rgba(255,255,255,0.08);
  color: #e5e5e5;
}

.btn-submit {
  padding: 10px 18px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.85rem;
}

.btn-submit:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.notification-modal {
  max-width: 320px;
  text-align: center;
}

.notification-icon {
  font-size: 2.5rem;
  margin-bottom: 16px;
  color: #a3a3a3;
}

.notification-modal p {
  color: #e5e5e5;
  font-size: 0.95rem;
  margin-bottom: 20px;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.25s;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.25s;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: scale(0.97);
}

@media (max-width: 768px) {
  .page-body {
    display: block;
  }

  .sort-btn {
    width: 100%;
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
  
  .controls-row {
    flex-direction: column;
    align-items: stretch;
  }
  
  .notebooks-grid {
    grid-template-columns: 1fr;
  }
  
  .modal-actions {
    flex-direction: column;
  }
  
  .btn-open,
  .btn-save,
  .btn-rate,
  .btn-submit,
  .btn-cancel {
    width: 100%;
    justify-content: center;
  }
}
</style>