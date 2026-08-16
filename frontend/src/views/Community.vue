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
        </div>
        
        <div class="search-bar">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Поиск по названию или тегам..."
            @keydown.enter="loadNotebooks"
          >
          <button class="btn-search" @click="loadNotebooks">
            <i class="fas fa-search"></i>
          </button>
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
        
        <div v-if="loading" class="loading">
          <i class="fas fa-spinner fa-spin"></i> Загружаем...
        </div>
        
        <div v-else-if="notebooks.length === 0" class="empty-state">
          <i class="fas fa-book"></i>
          <p>Нет публичных тетрадей</p>
        </div>
        
        <div v-else class="notebooks-grid">
          <div 
            v-for="notebook in notebooks" 
            :key="notebook.id" 
            class="notebook-card"
            @click="router.push(`/notebook/${notebook.id}`)"
          >
            <div class="notebook-color" :style="{ background: notebook.color || '#A78BFA' }"></div>
            <div class="notebook-info">
              <div class="notebook-title">
                {{ notebook.title }}
                <i v-if="notebook.is_verified" class="fas fa-check-circle verified"></i>
              </div>
              <div class="notebook-author">
                <i class="fas fa-user"></i> {{ notebook.author?.first_name || 'Автор' }}
              </div>
              <div class="notebook-rating">
                <i class="fas fa-star"></i> {{ notebook.average_rating?.toFixed(1) || 0 }}
                <span class="rating-count">({{ notebook.total_ratings || 0 }})</span>
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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'

const router = useRouter()
const notebooks = ref([])
const loading = ref(true)
const searchQuery = ref('')
const currentSort = ref('rating')

const sorts = [
  { value: 'rating', label: 'По рейтингу' },
  { value: 'newest', label: 'Новые' },
  { value: 'popular', label: 'Популярные' }
]

async function loadNotebooks() {
  loading.value = true
  try {
    const params = new URLSearchParams({ sort: currentSort.value, limit: '50' })
    if (searchQuery.value) params.append('search', searchQuery.value)
    
    const data = await apiFetch(`/notebooks/community?${params}`)
    notebooks.value = data.notebooks || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
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
  max-width: 900px;
  margin: 0 auto;
  padding: 32px;
}

.community-header {
  margin-bottom: 24px;
}

.community-header h1 {
  font-size: 1.5rem;
  font-weight: 800;
}

.community-header h1 i {
  color: #F472B6;
  margin-right: 10px;
}

.search-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.search-bar input {
  flex: 1;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
}

.search-bar input:focus {
  border-color: #A78BFA;
}

.btn-search {
  padding: 12px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  cursor: pointer;
}

.sort-row {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.sort-btn {
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.06);
  background: transparent;
  color: #94A3B8;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.75rem;
  transition: all 0.2s;
}

.sort-btn.active {
  background: rgba(167,139,250,0.15);
  border-color: #A78BFA;
  color: #A78BFA;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
}

.empty-state i {
  font-size: 3rem;
  color: #64748B;
  margin-bottom: 16px;
  display: block;
}

.notebooks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.notebook-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 20px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  gap: 16px;
}

.notebook-card:hover {
  border-color: #A78BFA;
  transform: translateY(-2px);
}

.notebook-color {
  width: 8px;
  border-radius: 4px;
  flex-shrink: 0;
}

.notebook-info {
  flex: 1;
  min-width: 0;
}

.notebook-title {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.verified {
  color: #34D399;
  font-size: 0.8rem;
}

.notebook-author {
  font-size: 0.8rem;
  color: #94A3B8;
  margin-bottom: 4px;
}

.notebook-rating {
  font-size: 0.85rem;
  color: #FBBF24;
  margin-bottom: 8px;
}

.rating-count {
  color: #64748B;
  font-size: 0.75rem;
}

.notebook-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.tag {
  font-size: 0.65rem;
  padding: 2px 8px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  color: #64748B;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .notebooks-grid {
    grid-template-columns: 1fr;
  }
}
</style>