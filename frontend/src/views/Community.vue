<template>
  <div class="community-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Каталог тетрадей</span>
      </header>
      
      <div class="content">
        <div class="community-header">
          <h1><i class="fas fa-globe"></i> Публичные тетради</h1>
          <p>Находи конспекты других учеников и делись своими</p>
        </div>
        
        <div class="controls-row">
          <div class="search-bar">
            <i class="fas fa-search search-icon"></i>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Поиск по названию или тегам..."
              @keydown.enter="loadNotebooks"
            >
          </div>
          
          <div class="sort-row">
            <button 
              v-for="sort in sorts" 
              :key="sort.value"
              class="sort-btn"
              :class="{ active: currentSort === sort.value }"
              @click="currentSort = sort.value; loadNotebooks()"
            >
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
          <div 
            v-for="notebook in notebooks" 
            :key="notebook.id" 
            class="notebook-card"
            @click="openNotebook(notebook)"
          >
            <div class="notebook-color" :style="{ background: notebook.color || '#A78BFA' }"></div>
            <div class="notebook-info">
              <div class="notebook-title">{{ notebook.title }}</div>
              <div class="notebook-author">
                <div class="author-avatar-mini">{{ getInitial(notebook.author) }}</div>
                <span class="author-name">{{ getAuthorName(notebook.author) }}</span>
              </div>
              <div class="notebook-footer">
                <div class="notebook-rating">
                  <i class="fas fa-star"></i> {{ formatRating(notebook.average_rating) }}
                  <span class="rating-count">({{ countRatings(notebook.ratings) }})</span>
                </div>
                <div v-if="notebook.tags && notebook.tags.length" class="notebook-tags">
                  <span v-for="tag in notebook.tags.slice(0, 3)" :key="tag" class="tag">{{ tag }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <Transition name="modal">
      <div v-if="selectedNotebook" class="modal" @click.self="selectedNotebook = null">
        <div class="modal-card">
          <div class="modal-header">
            <div class="modal-color" :style="{ background: selectedNotebook.color || '#A78BFA' }"></div>
            <div class="modal-title-block">
              <h2>{{ selectedNotebook.title }}</h2>
            </div>
            <button class="modal-close" @click="selectedNotebook = null">
              <i class="fas fa-times"></i>
            </button>
          </div>
          
          <div class="modal-body">
            <div class="modal-author">
              <div class="author-avatar">{{ getInitial(selectedNotebook.author) }}</div>
              <div>
                <div class="author-name">{{ getAuthorName(selectedNotebook.author) }}</div>
                <div class="author-email">{{ getAuthorEmail(selectedNotebook.author) }}</div>
                <div class="author-meta">
                  <span><i class="fas fa-star"></i> {{ formatRating(selectedNotebook.average_rating) }} ({{ countRatings(selectedNotebook.ratings) }})</span>
                  <span><i class="fas fa-eye"></i> {{ selectedNotebook.views_count || 0 }}</span>
                  <span><i class="fas fa-copy"></i> {{ selectedNotebook.copies_count || 0 }}</span>
                </div>
              </div>
            </div>
            
            <div v-if="selectedNotebook.description" class="modal-description">
              <p>{{ selectedNotebook.description }}</p>
            </div>
            
            <div v-if="selectedNotebook.tags?.length" class="modal-tags">
              <span v-for="tag in selectedNotebook.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
            
            <div class="modal-actions">
              <button class="btn-open" @click="router.push(`/notebook/${selectedNotebook.id}`)">
                <i class="fas fa-book-open"></i> Открыть
              </button>
              <button class="btn-rate" @click="rateNotebook(selectedNotebook)">
                <i class="fas fa-star"></i> Оценить
              </button>
              <button class="btn-copy" @click="copyNotebook(selectedNotebook)">
                <i class="fas fa-copy"></i> Сохранить
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { supabase } from '../api/supabase'

const router = useRouter()
const notebooks = ref([])
const loading = ref(true)
const searchQuery = ref('')
const currentSort = ref('rating')
const selectedNotebook = ref(null)

const sorts = [
  { value: 'rating', label: 'По рейтингу' },
  { value: 'newest', label: 'Новые' },
  { value: 'popular', label: 'Популярные' }
]

function getInitial(author) {
  if (!author) return 'А'
  return (author.first_name || author.email || 'А')[0].toUpperCase()
}

function getAuthorName(author) {
  if (!author) return 'Автор'
  return author.first_name || author.email?.split('@')[0] || 'Автор'
}

function getAuthorEmail(author) {
  return author?.email || ''
}

function formatRating(rating) {
  return Number(rating || 0).toFixed(1)
}

function countRatings(ratings) {
  if (!ratings) return 0
  if (typeof ratings === 'string') {
    try {
      return Object.keys(JSON.parse(ratings)).length
    } catch {
      return 0
    }
  }
  return Object.keys(ratings).length
}

async function loadNotebooks() {
  loading.value = true
  try {
    let query = supabase
      .from('notebooks')
      .select(`
        *,
        author:rubium_users!user_id(id, first_name, email)
      `)
      .eq('is_public', true)
    
    if (currentSort.value === 'newest') {
      query = query.order('created_at', { ascending: false })
    } else if (currentSort.value === 'rating') {
      query = query.order('average_rating', { ascending: false })
    } else if (currentSort.value === 'popular') {
      query = query.order('views_count', { ascending: false })
    }
    
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      query = query.or(`title.ilike.%${q}%,description.ilike.%${q}%`)
    }
    
    const { data, error } = await query
    
    if (error) throw error
    
    notebooks.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openNotebook(notebook) {
  selectedNotebook.value = notebook
}

async function rateNotebook(notebook) {
  const rating = prompt('Оцени от 1 до 5:', '5')
  if (!rating) return
  
  const numRating = parseInt(rating)
  if (numRating < 1 || numRating > 5) {
    alert('Оценка от 1 до 5')
    return
  }
  
  try {
    const { data: { session } } = await supabase.auth.getSession()
    if (!session) {
      alert('Нужно авторизоваться')
      return
    }
    
    const ratings = typeof notebook.ratings === 'string' 
      ? JSON.parse(notebook.ratings || '{}') 
      : (notebook.ratings || {})
    
    ratings[session.user.id] = numRating
    
    const values = Object.values(ratings)
    const avg = values.reduce((sum, v) => sum + v, 0) / values.length
    
    const { error } = await supabase
      .from('notebooks')
      .update({ ratings, average_rating: avg })
      .eq('id', notebook.id)
    
    if (error) throw error
    
    await loadNotebooks()
    selectedNotebook.value = null
  } catch (e) {
    console.error(e)
    alert('Ошибка при оценке')
  }
}

async function copyNotebook(notebook) {
  try {
    const { data: { session } } = await supabase.auth.getSession()
    if (!session) {
      alert('Нужно авторизоваться')
      return
    }
    
    const { data: userData } = await supabase
      .from('rubium_users')
      .select('id')
      .eq('auth_id', session.user.id)
      .single()
    
    const { error } = await supabase
      .from('notebooks')
      .insert({
        user_id: userData.id,
        title: notebook.title + ' (копия)',
        description: notebook.description,
        color: notebook.color,
        tags: notebook.tags,
        is_public: false,
        content: notebook.content
      })
    
    if (error) throw error
    
    alert('Тетрадь скопирована в твой каталог!')
    selectedNotebook.value = null
  } catch (e) {
    console.error(e)
    alert('Ошибка при копировании')
  }
}

onMounted(loadNotebooks)
</script>

<style scoped>
.community-page {
  display: flex;
  min-height: 100vh;
}

.main-content {
  margin-left: 240px;
  flex: 1;
}

.topbar {
  padding: 16px 32px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.page-title {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: #64748B;
  font-family: 'JetBrains Mono', monospace;
}

.content {
  max-width: 1100px;
  margin: 0 auto;
  padding: 32px;
}

.community-header {
  margin-bottom: 28px;
}

.community-header h1 {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.community-header h1 i {
  color: #A78BFA;
  margin-right: 12px;
}

.community-header p {
  color: #94A3B8;
  font-size: 0.95rem;
}

.controls-row {
  display: flex;
  gap: 16px;
  margin-bottom: 28px;
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
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #64748B;
  font-size: 0.85rem;
}

.search-bar input {
  width: 100%;
  padding: 14px 16px 14px 44px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 14px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  transition: all 0.3s;
}

.search-bar input:focus {
  border-color: #A78BFA;
  background: rgba(167,139,250,0.05);
  box-shadow: 0 0 0 3px rgba(167,139,250,0.1);
}

.search-bar input::placeholder {
  color: #64748B;
}

.sort-row {
  display: flex;
  gap: 4px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  padding: 4px;
}

.sort-btn {
  padding: 8px 16px;
  border-radius: 9px;
  border: none;
  background: transparent;
  color: #94A3B8;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 600;
  transition: all 0.3s;
}

.sort-btn:hover {
  color: #F1F5F9;
}

.sort-btn.active {
  background: #A78BFA;
  color: #0F0F1A;
  box-shadow: 0 4px 12px rgba(167,139,250,0.3);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 60px 20px;
  color: #94A3B8;
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid rgba(167,139,250,0.2);
  border-top-color: #A78BFA;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #94A3B8;
}

.empty-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 20px;
  border-radius: 20px;
  background: rgba(167,139,250,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  color: #A78BFA;
}

.empty-state h3 {
  font-size: 1.2rem;
  font-weight: 700;
  color: #F1F5F9;
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 0.9rem;
}

.notebooks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 18px;
}

.notebook-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
  padding: 22px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
  display: flex;
  gap: 18px;
}

.notebook-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.2);
}

.notebook-color {
  width: 6px;
  border-radius: 3px;
  flex-shrink: 0;
}

.notebook-info {
  flex: 1;
  min-width: 0;
}

.notebook-title {
  font-size: 1.05rem;
  font-weight: 700;
  margin-bottom: 8px;
  line-height: 1.3;
}

.notebook-author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8rem;
  color: #94A3B8;
  margin-bottom: 14px;
}

.author-avatar-mini {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.65rem;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
}

.author-name {
  font-weight: 600;
  color: #E2E8F0;
}

.notebook-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.notebook-rating {
  font-size: 0.85rem;
  color: #FBBF24;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.notebook-rating i {
  font-size: 0.75rem;
}

.rating-count {
  color: #64748B;
  font-size: 0.75rem;
}

.notebook-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.tag {
  font-size: 0.65rem;
  padding: 4px 10px;
  background: rgba(167,139,250,0.08);
  border: 1px solid rgba(167,139,250,0.15);
  border-radius: 8px;
  color: #A78BFA;
  font-weight: 500;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 20px;
}

.modal-card {
  background: #1a1a2e;
  border: 1px solid rgba(167,139,250,0.15);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 640px;
  box-shadow: 0 24px 48px rgba(0,0,0,0.4);
}

.modal-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 28px;
}

.modal-color {
  width: 6px;
  height: 56px;
  border-radius: 3px;
  flex-shrink: 0;
}

.modal-title-block {
  flex: 1;
}

.modal-title-block h2 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 8px;
  line-height: 1.3;
}

.modal-close {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 6px;
  border-radius: 8px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.modal-close:hover {
  color: #F1F5F9;
  background: rgba(255,255,255,0.06);
}

.modal-author {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 24px;
}

.author-avatar {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(167,139,250,0.3);
}

.author-name {
  font-weight: 700;
  font-size: 1rem;
}

.author-email {
  font-size: 0.8rem;
  color: #64748B;
  margin-bottom: 6px;
}

.author-meta {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
  font-size: 0.78rem;
  color: #94A3B8;
}

.author-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.author-meta i {
  font-size: 0.7rem;
}

.modal-description {
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255,255,255,0.03);
  border-radius: 12px;
  border: 1px solid rgba(255,255,255,0.06);
}

.modal-description p {
  color: #94A3B8;
  font-size: 0.9rem;
  line-height: 1.6;
}

.modal-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 28px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-open,
.btn-rate,
.btn-copy {
  padding: 12px 22px;
  border: none;
  border-radius: 12px;
  font-family: inherit;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-open {
  background: linear-gradient(135deg, #A78BFA, #8B5CF6);
  color: #fff;
  box-shadow: 0 4px 16px rgba(167,139,250,0.3);
}

.btn-open:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(167,139,250,0.4);
}

.btn-rate {
  background: rgba(251,191,36,0.12);
  color: #FBBF24;
  border: 1px solid rgba(251,191,36,0.2);
}

.btn-rate:hover {
  background: rgba(251,191,36,0.2);
  transform: translateY(-2px);
}

.btn-copy {
  background: rgba(255,255,255,0.04);
  color: #94A3B8;
  border: 1px solid rgba(255,255,255,0.08);
}

.btn-copy:hover {
  background: rgba(255,255,255,0.08);
  color: #F1F5F9;
  transform: translateY(-2px);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.3s;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: scale(0.95);
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 20px;
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
  .btn-rate,
  .btn-copy {
    width: 100%;
    justify-content: center;
  }
}
</style>