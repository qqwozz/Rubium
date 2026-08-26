<template>
  <div class="notebooks-page">
    <Sidebar ref="sidebarRef" />
    
    <header class="mobile-header">
      <button class="mobile-menu-btn" @click="sidebarRef?.toggle()">
        <i class="fas fa-bars"></i>
      </button>
    </header>
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Тетради</span>
      </header>
      
      <div class="content">
        <div class="notebooks-header">
          <h1>Мои тетради</h1>
          <button class="btn-create" @click="openCreateModal">
            <i class="fas fa-plus"></i> Создать
          </button>
        </div>
        
        <div v-if="loading" class="loading">
          <div class="spinner"></div>
          <span>Загружаем...</span>
        </div>
        
        <div v-else-if="notebooks.length === 0" class="empty-state">
          <div class="empty-icon"><i class="fas fa-book-open"></i></div>
          <h3>У тебя пока нет тетрадей</h3>
          <p>Создай первую — начни вести конспекты</p>
        </div>
        
        <div v-else class="notebooks-grid">
          <div v-for="notebook in notebooks" :key="notebook.id" class="notebook-card" @click="openNotebook(notebook)">
            <div class="notebook-color" :style="{ background: notebook.color || '#525252' }"></div>
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
                <button v-for="color in colors" :key="color" class="color-btn"
                  :class="{ active: newColor === color }" :style="{ background: color }" @click="newColor = color">
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>Теги</label>
              <div class="tags-input-wrapper">
                <div v-for="tag in newTags" :key="tag" class="selected-tag">
                  {{ tag }}<span @click="removeNewTag(tag)">×</span>
                </div>
                <input v-model="tagInput" type="text" placeholder="Выбери тег..."
                  @input="filterTags" @keydown.enter.prevent="addTag" @keydown.backspace="removeLastTag">
              </div>
              <div v-if="filteredTags.length" class="tags-suggestions">
                <div v-for="tag in filteredTags" :key="tag" class="tag-suggestion" @click="addTag(tag)">{{ tag }}</div>
              </div>
            </div>
            <div class="form-group">
              <label class="checkbox-label">
                <input v-model="newIsPublic" type="checkbox"> Публичная тетрадь
              </label>
            </div>
            <div class="modal-actions">
              <button class="btn-cancel" @click="showCreateModal = false">Отмена</button>
              <button class="btn-save" @click="createNotebook" :disabled="!newTitle.trim()">Создать</button>
            </div>
          </div>
        </div>
      </Transition>

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
                <button v-for="color in colors" :key="color" class="color-btn"
                  :class="{ active: editColor === color }" :style="{ background: color }" @click="editColor = color">
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>Теги</label>
              <div class="tags-input-wrapper">
                <div v-for="tag in editTags" :key="tag" class="selected-tag">
                  {{ tag }}<span @click="removeEditTag(tag)">×</span>
                </div>
                <input v-model="editTagInput" type="text" placeholder="Выбери тег..."
                  @input="filterEditTags" @keydown.enter.prevent="addEditTag" @keydown.backspace="removeLastEditTag">
              </div>
              <div v-if="filteredEditTags.length" class="tags-suggestions">
                <div v-for="tag in filteredEditTags" :key="tag" class="tag-suggestion" @click="addEditTag(tag)">{{ tag }}</div>
              </div>
            </div>
            <div class="form-group">
              <label class="checkbox-label">
                <input v-model="editIsPublic" type="checkbox"> Публичная тетрадь
              </label>
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-delete" @click="deleteNotebook">Удалить</button>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import { supabase } from '../api/supabase'
import tagsRaw from '../assets/tags.txt?raw'

const router = useRouter()
const notebooks = ref([])
const loading = ref(true)
const sidebarRef = ref(null)

const showCreateModal = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const newColor = ref('#525252')
const newTags = ref([])
const newIsPublic = ref(false)
const tagInput = ref('')
const filteredTags = ref([])

const showEditModal = ref(false)
const editId = ref(null)
const editTitle = ref('')
const editDescription = ref('')
const editColor = ref('#525252')
const editTags = ref([])
const editIsPublic = ref(false)
const editTagInput = ref('')
const filteredEditTags = ref([])

const colors = ['#525252', '#737373', '#a3a3a3', '#d4d4d4', '#a78bfa', '#f472b6']
const allTags = ref([])

function loadTags() {
  allTags.value = tagsRaw.split('\n').map(t => t.trim()).filter(Boolean)
}

function filterTags() {
  const q = tagInput.value.toLowerCase().trim()
  filteredTags.value = allTags.value.filter(t => t.toLowerCase().includes(q) && !newTags.value.includes(t)).slice(0, 8)
}

function addTag(tag) {
  const cleanTag = typeof tag === 'string' ? tag : tagInput.value.trim()
  if (cleanTag && !newTags.value.includes(cleanTag)) newTags.value.push(cleanTag)
  tagInput.value = ''
  filteredTags.value = []
}

function removeNewTag(tag) { newTags.value = newTags.value.filter(t => t !== tag) }
function removeLastTag() { if (!tagInput.value && newTags.value.length) newTags.value.pop() }

function filterEditTags() {
  const q = editTagInput.value.toLowerCase().trim()
  filteredEditTags.value = allTags.value.filter(t => t.toLowerCase().includes(q) && !editTags.value.includes(t)).slice(0, 8)
}

function addEditTag(tag) {
  const cleanTag = typeof tag === 'string' ? tag : editTagInput.value.trim()
  if (cleanTag && !editTags.value.includes(cleanTag)) editTags.value.push(cleanTag)
  editTagInput.value = ''
  filteredEditTags.value = []
}

function removeEditTag(tag) { editTags.value = editTags.value.filter(t => t !== tag) }
function removeLastEditTag() { if (!editTagInput.value && editTags.value.length) editTags.value.pop() }

function openCreateModal() {
  newTitle.value = ''
  newDescription.value = ''
  newColor.value = '#525252'
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
  editColor.value = notebook.color || '#525252'
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
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function deleteNotebook() {
  if (!editId.value) return
  try {
    await apiFetch(`/notebooks/${editId.value}`, { method: 'DELETE' })
    showEditModal.value = false
    editId.value = null
    await loadNotebooks()
  } catch (e) { console.error(e) }
}

async function createNotebook() {
  try {
    const { data: { session } } = await supabase.auth.getSession()
    const { data: userData } = await supabase.from('rubium_users').select('id').eq('auth_id', session.user.id).single()
    await apiFetch('/notebooks', {
      method: 'POST',
      body: JSON.stringify({
        title: newTitle.value, description: newDescription.value, color: newColor.value,
        tags: newTags.value, is_public: newIsPublic.value, user_id: userData.id
      })
    })
    showCreateModal.value = false
    await loadNotebooks()
  } catch (e) { console.error(e) }
}

async function saveNotebook() {
  try {
    await apiFetch(`/notebooks/${editId.value}`, {
      method: 'PUT',
      body: JSON.stringify({
        title: editTitle.value, description: editDescription.value, color: editColor.value,
        tags: editTags.value, is_public: editIsPublic.value
      })
    })
    showEditModal.value = false
    await loadNotebooks()
  } catch (e) { console.error(e) }
}

function openNotebook(notebook) {
  if (!notebook.id) return
  router.push(`/notebook/${notebook.id}/edit`)
}

onMounted(() => { loadTags(); loadNotebooks() })
</script>

<style scoped>
.notebooks-page {
  display: flex;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.main-content { margin-left: 240px; flex: 1; }

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

.mobile-header { display: none; }

.notebooks-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
  flex-wrap: wrap;
  gap: 12px;
}

.notebooks-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.btn-create {
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
}

.btn-create:hover { background: #e5e5e5; border-color: #e5e5e5; }

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 20px;
  color: #737373;
}

.spinner {
  width: 32px; height: 32px;
  border: 2px solid rgba(255,255,255,0.06);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #737373;
}

.empty-icon {
  width: 64px; height: 64px;
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

.empty-state h3 { font-size: 1.1rem; font-weight: 600; color: #e5e5e5; margin-bottom: 6px; }
.empty-state p { font-size: 0.9rem; }

.notebooks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
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
}

.notebook-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.notebook-color { width: 4px; border-radius: 2px; flex-shrink: 0; }

.notebook-info { flex: 1; min-width: 0; }

.notebook-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}

.notebook-title { font-size: 1rem; font-weight: 600; color: #e5e5e5; }

.btn-edit {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  padding: 4px;
  border-radius: 6px;
  transition: all 0.15s ease;
  flex-shrink: 0;
  font-size: 0.85rem;
}

.btn-edit:hover { color: #a3a3a3; background: rgba(255,255,255,0.04); }

.notebook-description {
  font-size: 0.82rem;
  color: #737373;
  margin-bottom: 10px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notebook-meta {
  display: flex;
  gap: 12px;
  font-size: 0.78rem;
  color: #525252;
  margin-bottom: 8px;
}

.notebook-meta i { font-size: 0.7rem; margin-right: 3px; }

.notebook-tags { display: flex; gap: 4px; flex-wrap: wrap; }

.tag {
  font-size: 0.7rem;
  padding: 3px 8px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  color: #737373;
  font-weight: 500;
}

.modal {
  position: fixed; inset: 0;
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
  max-width: 480px;
}

.modal-card h2 { font-size: 1.2rem; font-weight: 600; color: #ffffff; margin-bottom: 24px; }

.form-group { margin-bottom: 18px; }

.form-group label {
  display: block;
  font-size: 0.8rem;
  font-weight: 500;
  color: #a3a3a3;
  margin-bottom: 8px;
}

.form-group input[type="text"],
.form-group textarea {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  resize: vertical;
  transition: all 0.2s ease;
}

.form-group input[type="text"]:focus,
.form-group textarea:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.form-group input::placeholder,
.form-group textarea::placeholder { color: #525252; }

.colors-row { display: flex; gap: 8px; flex-wrap: wrap; }

.color-btn {
  width: 28px; height: 28px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.15s ease;
}

.color-btn.active { border-color: #ffffff; transform: scale(1.1); }

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: #a3a3a3;
  cursor: pointer;
}

.tags-input-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
}

.tags-input-wrapper input {
  flex: 1; min-width: 120px;
  background: transparent; border: none; outline: none;
  color: #e5e5e5;
  font-family: inherit;
  font-size: 0.85rem;
  padding: 4px;
}

.selected-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  background: rgba(255,255,255,0.06);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  color: #a3a3a3;
  font-size: 0.75rem;
  font-weight: 500;
}

.selected-tag span { cursor: pointer; color: #525252; font-size: 0.85rem; line-height: 1; }
.selected-tag span:hover { color: #e5e5e5; }

.tags-suggestions {
  margin-top: 6px;
  background: #111111;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  max-height: 160px;
  overflow-y: auto;
}

.tag-suggestion {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 0.85rem;
  color: #737373;
  transition: all 0.15s ease;
}

.tag-suggestion:hover { background: rgba(255,255,255,0.04); color: #e5e5e5; }

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn-cancel {
  padding: 10px 18px;
  background: rgba(255,255,255,0.04);
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.btn-cancel:hover { background: rgba(255,255,255,0.08); color: #e5e5e5; }

.btn-save {
  padding: 10px 18px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.btn-save:hover { background: #e5e5e5; border-color: #e5e5e5; }
.btn-save:disabled { opacity: 0.4; cursor: not-allowed; }

.btn-delete {
  padding: 10px 18px;
  background: rgba(255,255,255,0.04);
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  margin-right: auto;
  transition: all 0.2s ease;
}

.btn-delete:hover {
  background: rgba(239,68,68,0.08);
  color: #ef4444;
  border-color: rgba(239,68,68,0.15);
}

.modal-enter-active, .modal-leave-active { transition: opacity 0.25s; }
.modal-enter-active .modal-card, .modal-leave-active .modal-card { transition: transform 0.25s; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.modal-enter-from .modal-card, .modal-leave-to .modal-card { transform: scale(0.97); }

@media (max-width: 768px) {
  .main-content { margin-left: 0; }
  .topbar { display: none; }
  
  .mobile-header {
    display: flex;
    align-items: flex-start;
    padding: 10px 5px;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    position: sticky;
    top: 0;
    background: #0a0a0a;
    z-index: 50;
  }
  
  .mobile-menu-btn {
    background: none;
    border: none;
    color: #737373;
    font-size: 1.1rem;
    cursor: pointer;
    padding: 4px 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.15s ease;
  }
  
  .mobile-menu-btn:hover { color: #e5e5e5; }
  
  .content { padding: 32px 24px 64px; }
  .notebooks-grid { grid-template-columns: 1fr; }
  
  .modal-actions { flex-direction: column; }
  .btn-cancel, .btn-save, .btn-delete { width: 100%; justify-content: center; }
  .btn-delete { margin-right: 0; order: 3; }
}
</style>