<template>
  <div class="notebooks-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Тетради</span>
      </header>
      
      <div class="content">
        <div class="notebooks-header">
          <h1><i class="fas fa-book"></i> Мои тетради</h1>
          <button class="btn-create" @click="showCreateModal = true">
            <i class="fas fa-plus"></i> Создать тетрадь
          </button>
        </div>
        
        <div v-if="loading" class="loading">
          <i class="fas fa-spinner fa-spin"></i> Загружаем...
        </div>
        
        <div v-else-if="notebooks.length === 0" class="empty-state">
          <i class="fas fa-book-open"></i>
          <p>У тебя пока нет тетрадей</p>
          <p class="empty-hint">Создай первую — начни вести конспекты</p>
        </div>
        
        <div v-else class="notebooks-grid">
          <div v-for="notebook in notebooks" :key="notebook.id" class="notebook-card" @click="openNotebook(notebook)">
            <div class="notebook-color" :style="{ background: notebook.color || '#A78BFA' }"></div>
            <div class="notebook-info">
              <div class="notebook-title">{{ notebook.title }}</div>
              <div class="notebook-meta">
                <span>{{ notebook.pages_count || 0 }} страниц</span>
                <span v-if="notebook.is_public"><i class="fas fa-globe"></i> Публичная</span>
                <span v-else><i class="fas fa-lock"></i> Приватная</span>
              </div>
              <div v-if="notebook.tags && notebook.tags.length" class="notebook-tags">
                <span v-for="tag in notebook.tags.slice(0, 3)" :key="tag" class="tag">{{ tag }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="showCreateModal" class="modal" @click.self="showCreateModal = false">
      <div class="modal-card">
        <h2>Новая тетрадь</h2>
        
        <div class="form-group">
          <label>Название</label>
          <input v-model="newTitle" type="text" placeholder="Например: Тригонометрия ЕГЭ">
        </div>
        
        <div class="form-group">
          <label>Цвет</label>
          <div class="colors-row">
            <button 
              v-for="color in colors" 
              :key="color"
              class="color-btn"
              :class="{ active: newColor === color }"
              :style="{ background: color }"
              @click="newColor = color"
            ></button>
          </div>
        </div>
        
        <div class="form-group">
          <label>Теги (через запятую)</label>
          <input v-model="newTags" type="text" placeholder="math, EGE, тригонометрия">
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="newIsPublic" type="checkbox">
            Публичная тетрадь
          </label>
        </div>
        
        <div class="modal-actions">
          <button class="btn-cancel" @click="showCreateModal = false">Отмена</button>
          <button class="btn-save" @click="createNotebook" :disabled="!newTitle.trim()">
            Создать
          </button>
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
import { supabase } from '../api/supabase'

const router = useRouter()
const notebooks = ref([])
const loading = ref(true)
const showCreateModal = ref(false)
const newTitle = ref('')
const newColor = ref('#A78BFA')
const newTags = ref('')
const newIsPublic = ref(false)

const colors = ['#A78BFA', '#F472B6', '#34D399', '#FBBF24', '#F87171', '#60A5FA']

async function loadNotebooks() {
  loading.value = true
  try {
    const data = await apiFetch('/notebooks')
    notebooks.value = data.notebooks || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function createNotebook() {
  try {
    const tags = newTags.value
      .split(',')
      .map(t => t.trim())
      .filter(Boolean)
    
    const { data: { session } } = await supabase.auth.getSession()
    const { data: userData } = await supabase
      .from('rubium_users')
      .select('id')
      .eq('auth_id', session.user.id)
      .single()
    
    await apiFetch('/notebooks', {
      method: 'POST',
      body: JSON.stringify({
        title: newTitle.value,
        color: newColor.value,
        tags,
        is_public: newIsPublic.value,
        user_id: userData.id
      })
    })
    
    showCreateModal.value = false
    newTitle.value = ''
    newTags.value = ''
    newIsPublic.value = false
    await loadNotebooks()
  } catch (e) {
    console.error(e)
  }
}

function openNotebook(notebook) {
  if (!notebook.id) return
  router.push(`/notebook/${notebook.id}/edit`)
}

onMounted(loadNotebooks)
</script>

<style scoped>
.notebooks-page {
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

.notebooks-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.notebooks-header h1 {
  font-size: 1.5rem;
  font-weight: 800;
}

.notebooks-header h1 i {
  color: #34D399;
  margin-right: 10px;
}

.btn-create {
  padding: 10px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-create:hover {
  background: #8B5CF6;
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

.empty-hint {
  font-size: 0.85rem;
  margin-top: 8px;
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
  margin-bottom: 8px;
}

.notebook-meta {
  display: flex;
  gap: 12px;
  font-size: 0.75rem;
  color: #64748B;
  margin-bottom: 8px;
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

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 20px;
}

.modal-card {
  background: #1a1a2e;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 450px;
}

.modal-card h2 {
  font-size: 1.3rem;
  font-weight: 800;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 0.8rem;
  color: #94A3B8;
  margin-bottom: 6px;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  outline: none;
}

.form-group input[type="text"]:focus {
  border-color: #A78BFA;
}

.colors-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.color-btn.active {
  border-color: white;
  transform: scale(1.1);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: #94A3B8;
  cursor: pointer;
}

.modal-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn-cancel {
  padding: 10px 20px;
  background: rgba(255,255,255,0.04);
  color: #94A3B8;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  cursor: pointer;
  font-family: inherit;
}

.btn-save {
  padding: 10px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 600;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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