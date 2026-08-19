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
          <button class="btn-create" @click="openCreateModal">
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
              <div class="notebook-title-row">
                <div class="notebook-title">{{ notebook.title }}</div>
                <button class="btn-edit" @click.stop="openEditModal(notebook)">
                  <i class="fas fa-gear"></i>
                </button>
              </div>
              <div v-if="notebook.description" class="notebook-description">{{ notebook.description }}</div>
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

    <Teleport to="body">
      <!-- Модалка создания -->
      <Transition name="modal">
        <div v-if="showCreateModal" class="modal" @click.self="showCreateModal = false">
          <div class="modal-card">
            <h2>Новая тетрадь</h2>
            
            <div class="form-group">
              <label>Название</label>
              <input v-model="newTitle" type="text" placeholder="Например: Тригонометрия ЕГЭ">
            </div>

            <div class="form-group">
              <label>Описание</label>
              <textarea v-model="newDescription" rows="3" placeholder="О чём эта тетрадь?"></textarea>
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
              <label>Теги</label>
              <div class="tags-input-wrapper">
                <div v-for="tag in newTags" :key="tag" class="selected-tag">
                  {{ tag }}
                  <span @click="removeNewTag(tag)">×</span>
                </div>
                <input 
                  v-model="tagInput" 
                  type="text" 
                  placeholder="Выбери тег..."
                  @input="filterTags"
                  @keydown.enter.prevent="addTag"
                  @keydown.backspace="removeLastTag"
                >
              </div>
              <div v-if="filteredTags.length" class="tags-suggestions">
                <div 
                  v-for="tag in filteredTags" 
                  :key="tag"
                  class="tag-suggestion"
                  @click="addTag(tag)"
                >
                  {{ tag }}
                </div>
              </div>
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
      </Transition>

      <!-- Модалка редактирования -->
      <Transition name="modal">
        <div v-if="showEditModal" class="modal" @click.self="showEditModal = false">
          <div class="modal-card">
            <h2>Настройки тетради</h2>
            
            <div class="form-group">
              <label>Название</label>
              <input v-model="editTitle" type="text">
            </div>

            <div class="form-group">
              <label>Описание</label>
              <textarea v-model="editDescription" rows="3" placeholder="О чём эта тетрадь?"></textarea>
            </div>
            
            <div class="form-group">
              <label>Цвет</label>
              <div class="colors-row">
                <button 
                  v-for="color in colors" 
                  :key="color"
                  class="color-btn"
                  :class="{ active: editColor === color }"
                  :style="{ background: color }"
                  @click="editColor = color"
                ></button>
              </div>
            </div>
            
            <div class="form-group">
              <label>Теги</label>
              <div class="tags-input-wrapper">
                <div v-for="tag in editTags" :key="tag" class="selected-tag">
                  {{ tag }}
                  <span @click="removeEditTag(tag)">×</span>
                </div>
                <input 
                  v-model="editTagInput" 
                  type="text" 
                  placeholder="Выбери тег..."
                  @input="filterEditTags"
                  @keydown.enter.prevent="addEditTag"
                  @keydown.backspace="removeLastEditTag"
                >
              </div>
              <div v-if="filteredEditTags.length" class="tags-suggestions">
                <div 
                  v-for="tag in filteredEditTags" 
                  :key="tag"
                  class="tag-suggestion"
                  @click="addEditTag(tag)"
                >
                  {{ tag }}
                </div>
              </div>
            </div>
            
            <div class="form-group">
              <label class="checkbox-label">
                <input v-model="editIsPublic" type="checkbox">
                Публичная тетрадь
              </label>
            </div>
            
            <div class="modal-actions">
              <button class="btn-delete" @click="deleteNotebook">Удалить</button>
              <button class="btn-cancel" @click="showEditModal = false">Отмена</button>
              <button class="btn-save" @click="saveNotebook" :disabled="!editTitle.trim()">Сохранить</button>
            </div>
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
import { apiFetch } from '../api/client'
import { supabase } from '../api/supabase'
import tagsRaw from '../assets/tags.txt?raw'

const router = useRouter()
const notebooks = ref([])
const loading = ref(true)

const showCreateModal = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const newColor = ref('#A78BFA')
const newTags = ref([])
const newIsPublic = ref(false)
const tagInput = ref('')
const filteredTags = ref([])

const showEditModal = ref(false)
const editId = ref(null)
const editTitle = ref('')
const editDescription = ref('')
const editColor = ref('#A78BFA')
const editTags = ref([])
const editIsPublic = ref(false)
const editTagInput = ref('')
const filteredEditTags = ref([])

const colors = ['#A78BFA', '#F472B6', '#34D399', '#FBBF24', '#F87171', '#60A5FA']

const allTags = ref([])

function loadTags() {
  allTags.value = tagsRaw.split('\n').map(t => t.trim()).filter(Boolean)
}

function filterTags() {
  const q = tagInput.value.toLowerCase().trim()
  filteredTags.value = allTags.value.filter(t => 
    t.toLowerCase().includes(q) && !newTags.value.includes(t)
  ).slice(0, 8)
}

function addTag(tag) {
  const cleanTag = typeof tag === 'string' ? tag : tagInput.value.trim()
  if (cleanTag && !newTags.value.includes(cleanTag)) {
    newTags.value.push(cleanTag)
  }
  tagInput.value = ''
  filteredTags.value = []
}

function removeNewTag(tag) {
  newTags.value = newTags.value.filter(t => t !== tag)
}

function removeLastTag() {
  if (!tagInput.value && newTags.value.length) {
    newTags.value.pop()
  }
}

function filterEditTags() {
  const q = editTagInput.value.toLowerCase().trim()
  filteredEditTags.value = allTags.value.filter(t => 
    t.toLowerCase().includes(q) && !editTags.value.includes(t)
  ).slice(0, 8)
}

function addEditTag(tag) {
  const cleanTag = typeof tag === 'string' ? tag : editTagInput.value.trim()
  if (cleanTag && !editTags.value.includes(cleanTag)) {
    editTags.value.push(cleanTag)
  }
  editTagInput.value = ''
  filteredEditTags.value = []
}

function removeEditTag(tag) {
  editTags.value = editTags.value.filter(t => t !== tag)
}

function removeLastEditTag() {
  if (!editTagInput.value && editTags.value.length) {
    editTags.value.pop()
  }
}

function openCreateModal() {
  newTitle.value = ''
  newDescription.value = ''
  newColor.value = '#A78BFA'
  newTags.value = []
  newIsPublic.value = false
  tagInput.value = ''
  filteredTags.value = []
  showCreateModal.value = true
}

function openEditModal(notebook) {
  editId.value = notebook.id
  editTitle.value = notebook.title || ''
  editDescription.value = notebook.description || ''
  editColor.value = notebook.color || '#A78BFA'
  editTags.value = notebook.tags || []
  editIsPublic.value = notebook.is_public || false
  editTagInput.value = ''
  filteredEditTags.value = []
  showEditModal.value = true
}

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

async function deleteNotebook() {
  if (!editId.value) return
  
  try {
    await apiFetch(`/notebooks/${editId.value}`, {
      method: 'DELETE'
    })
    
    showEditModal.value = false
    await loadNotebooks()
  } catch (e) {
    console.error(e)
  }
}

async function createNotebook() {
  try {
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
        description: newDescription.value,
        color: newColor.value,
        tags: newTags.value,
        is_public: newIsPublic.value,
        user_id: userData.id
      })
    })
    
    showCreateModal.value = false
    await loadNotebooks()
  } catch (e) {
    console.error(e)
  }
}

async function saveNotebook() {
  try {
    await apiFetch(`/notebooks/${editId.value}`, {
      method: 'PUT',
      body: JSON.stringify({
        title: editTitle.value,
        description: editDescription.value,
        color: editColor.value,
        tags: editTags.value,
        is_public: editIsPublic.value
      })
    })
    
    showEditModal.value = false
    await loadNotebooks()
  } catch (e) {
    console.error(e)
  }
}

function openNotebook(notebook) {
  if (!notebook.id) return
  router.push(`/notebook/${notebook.id}/edit`)
}

onMounted(() => {
  loadTags()
  loadNotebooks()
})
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

.notebook-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}

.notebook-title {
  font-size: 1rem;
  font-weight: 700;
}

.btn-edit {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  transition: all 0.2s;
  flex-shrink: 0;
}

.btn-edit:hover {
  color: #A78BFA;
  background: rgba(167,139,250,0.1);
}

.btn-delete {
  padding: 10px 20px;
  background: rgba(248,113,113,0.1);
  color: #F87171;
  border: 1px solid rgba(248,113,113,0.2);
  border-radius: 12px;
  cursor: pointer;
  font-family: inherit;
  margin-right: auto;
}

.btn-delete:hover {
  background: rgba(248,113,113,0.2);
}

.notebook-description {
  font-size: 0.8rem;
  color: #94A3B8;
  margin-bottom: 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
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
  z-index: 9999;
  padding: 20px;
}

.modal-card {
  background: #1a1a2e;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 24px 48px rgba(0,0,0,0.4);
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

.form-group input[type="text"],
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  outline: none;
  resize: vertical;
}

.form-group input[type="text"]:focus,
.form-group textarea:focus {
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

.tags-input-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
}

.tags-input-wrapper input {
  flex: 1;
  min-width: 120px;
  background: transparent;
  border: none;
  outline: none;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.85rem;
  padding: 4px;
}

.selected-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: rgba(167,139,250,0.15);
  border: 1px solid rgba(167,139,250,0.2);
  border-radius: 8px;
  color: #A78BFA;
  font-size: 0.75rem;
  font-weight: 600;
}

.selected-tag span {
  cursor: pointer;
  color: #64748B;
}

.selected-tag span:hover {
  color: #F87171;
}

.tags-suggestions {
  margin-top: 6px;
  background: #1a1a2e;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  max-height: 160px;
  overflow-y: auto;
}

.tag-suggestion {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 0.85rem;
  color: #94A3B8;
  transition: all 0.2s;
}

.tag-suggestion:hover {
  background: rgba(167,139,250,0.1);
  color: #A78BFA;
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
    padding: 16px;
  }
  .notebooks-grid {
    grid-template-columns: 1fr;
  }
}
</style>