<template>
  <div class="notebook-write-page">
    <ThinSidebar ref="sidebarRef" />
    
    <header class="mobile-header">
      <button class="mobile-menu-btn" @click="sidebarRef?.toggle()">
        <i class="fas fa-bars"></i>
      </button>
    </header>
    
    <div class="main-content">
      <header class="topbar">
        <button class="back-btn" @click="router.push('/notebooks')">
          <i class="fas fa-arrow-left"></i> К тетрадям
        </button>
        <span class="page-title">{{ notebook?.title || 'Тетрадь' }}</span>
        <button class="btn-save" @click="saveNotebook" :disabled="saving">
          <i v-if="saving" class="fas fa-spinner fa-spin"></i>
          <span v-else><i class="fas fa-check"></i> Сохранить</span>
        </button>
      </header>
      
      <div class="editor-layout">
        <div class="sections-panel">
          <div class="sections-header">
            <span>Разделы</span>
            <button class="btn-add-section" @click="addSection" title="Добавить раздел">
              <i class="fas fa-plus"></i>
            </button>
          </div>
          
          <div v-if="sections.length === 0" class="sections-empty">
            <div class="empty-icon"><i class="fas fa-folder-open"></i></div>
            <p>Нет разделов</p>
          </div>
          
          <div v-for="(section, si) in sections" :key="section.id" class="section-item">
            <div class="section-header" :class="{ active: openSections.includes(section.id) }">
              <i 
                class="fas fa-chevron-right section-arrow" 
                :class="{ rotated: openSections.includes(section.id) }"
                @click="toggleSection(section.id)"
              ></i>
              
              <input 
                v-if="editingSectionId === section.id"
                v-model="section.title"
                class="section-title-input"
                @blur="editingSectionId = null"
                @keydown.enter="editingSectionId = null"
                @keydown.esc="editingSectionId = null"
                ref="sectionInput"
              />
              <span 
                v-else
                class="section-title"
                @click="toggleSection(section.id)"
                @dblclick="startEditSection(section.id)"
              >{{ section.title }}</span>
              
              <div class="section-actions">
                <button class="btn-icon" @click.stop="startEditSection(section.id)" title="Переименовать">
                  <i class="fas fa-pen"></i>
                </button>
                <button class="btn-icon danger" @click.stop="deleteSection(si)" title="Удалить">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            
            <div v-if="openSections.includes(section.id)" class="pages-list">
              <div 
                v-for="page in section.pages" 
                :key="page.id"
                class="page-item"
                :class="{ active: currentPage?.id === page.id }"
                @click="selectPage(page)"
              >
                <i class="fas fa-file"></i>
                <span>{{ page.title || 'Без названия' }}</span>
                <button class="btn-icon-mini" @click.stop="deletePage(si, page.id)" title="Удалить страницу">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              
              <button class="btn-add-page" @click="addPage(si)">
                <i class="fas fa-plus"></i> Страница
              </button>
            </div>
          </div>
        </div>
        
        <div class="editor-panel">
          <div v-if="!currentPage && sections.length > 0" class="editor-empty">
            <div class="empty-icon"><i class="fas fa-file"></i></div>
            <p>Выбери или создай страницу</p>
            <button class="btn-create" @click="addPage(0)">
              <i class="fas fa-plus"></i> Создать страницу
            </button>
          </div>
          
          <div v-if="!currentPage && sections.length === 0" class="editor-empty">
            <div class="empty-icon"><i class="fas fa-book"></i></div>
            <p>Создай первый раздел</p>
            <button class="btn-create" @click="addSection">
              <i class="fas fa-plus"></i> Создать раздел
            </button>
          </div>
          
          <div v-if="currentPage" class="editor-content">
            <div class="page-header">
              <input 
                v-model="currentPage.title" 
                class="page-title-input" 
                placeholder="Заголовок страницы"
              >
            </div>
            
            <NotebookEditor v-model="currentPage.content" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onBeforeUnmount } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import NotebookEditor from '../components/NotebookEditor.vue'
import ThinSidebar from '../components/ThinSidebar.vue'

const route = useRoute()
const router = useRouter()
const sidebarRef = ref(null)

const notebook = ref(null)
const sections = ref([])
const openSections = ref([])
const currentPage = ref(null)
const saving = ref(false)
const editingSectionId = ref(null)
const sectionInput = ref(null)

let autoSaveTimer = null

function debouncedSave() {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => {
    saveNotebook()
  }, 1500)
}

function handleKeydown(event) {
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault()
    saveNotebook()
  }
}

onMounted(() => {
  loadNotebook()
  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  if (route.params.id) saveNotebook()
})

onBeforeRouteLeave(() => {
  if (route.params.id) saveNotebook()
})

async function loadNotebook() {
  try {
    const data = await apiFetch(`/notebooks/${route.params.id}`)
    notebook.value = data.notebook
    sections.value = data.notebook.content?.sections || []
    if (sections.value.length > 0) {
      openSections.value = [sections.value[0].id]
      if (sections.value[0].pages?.length > 0) {
        currentPage.value = sections.value[0].pages[0]
      }
    }
  } catch (e) {
    console.error(e)
  }
}

function toggleSection(sectionId) {
  const index = openSections.value.indexOf(sectionId)
  if (index > -1) {
    openSections.value.splice(index, 1)
  } else {
    openSections.value.push(sectionId)
  }
}

function addSection() {
  const newSection = {
    id: `section-${Date.now()}`,
    title: 'Новый раздел',
    pages: []
  }
  sections.value.push(newSection)
  openSections.value.push(newSection.id)
  startEditSection(newSection.id)
}

async function startEditSection(sectionId) {
  editingSectionId.value = sectionId
  await nextTick()
  const input = sectionInput.value?.[0] || document.querySelector('.section-title-input')
  if (input) {
    input.focus()
    input.select()
  }
}

function deleteSection(index) {
  const section = sections.value[index]
  if (confirm(`Удалить раздел "${section.title}"?`)) {
    sections.value.splice(index, 1)
    if (currentPage.value && section.pages?.includes(currentPage.value)) {
      currentPage.value = null
    }
  }
}

function deletePage(sectionIndex, pageId) {
  const section = sections.value[sectionIndex]
  const pageIndex = section.pages.findIndex(p => p.id === pageId)
  if (pageIndex > -1) {
    if (confirm(`Удалить страницу "${section.pages[pageIndex].title}"?`)) {
      section.pages.splice(pageIndex, 1)
      if (currentPage.value?.id === pageId) {
        currentPage.value = null
      }
    }
  }
}

function deleteCurrentPage() {
  if (!currentPage.value) return
  if (confirm(`Удалить страницу "${currentPage.value.title}"?`)) {
    for (const section of sections.value) {
      const index = section.pages.findIndex(p => p.id === currentPage.value.id)
      if (index > -1) {
        section.pages.splice(index, 1)
        break
      }
    }
    currentPage.value = null
  }
}

function addPage(sectionIndex) {
  const newPage = {
    id: `page-${Date.now()}`,
    title: 'Новая страница',
    content: '',
    source_task_id: null
  }
  sections.value[sectionIndex].pages.push(newPage)
  currentPage.value = newPage
}

function selectPage(page) {
  currentPage.value = page
}

async function saveNotebook() {
  if (!route.params.id) return
  
  saving.value = true
  try {
    await apiFetch(`/notebooks/${route.params.id}`, {
      method: 'PUT',
      body: JSON.stringify({
        content: { sections: sections.value }
      })
    })
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.notebook-write-page {
  display: flex;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.main-content {
  margin-left: 64px;
  flex: 1;
}

.topbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 32px;
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

.btn-save {
  padding: 8px 16px;
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

.btn-save:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-save:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.editor-layout {
  display: flex;
  min-height: calc(100vh - 60px);
}

.sections-panel {
  width: 280px;
  border-right: 1px solid rgba(255,255,255,0.06);
  padding: 20px 16px;
  flex-shrink: 0;
  overflow-y: auto;
  background: #0a0a0a;
}

.sections-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 0.9rem;
  color: #e5e5e5;
}

.btn-add-section {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  color: #a3a3a3;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
}

.btn-add-section:hover {
  background: rgba(255,255,255,0.08);
  color: #e5e5e5;
}

.sections-empty {
  color: #525252;
  font-size: 0.85rem;
  text-align: center;
  padding: 40px 20px;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 12px;
  border-radius: 12px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  color: #a3a3a3;
}

.section-item {
  margin-bottom: 4px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  transition: all 0.15s ease;
  color: #737373;
}

.section-header:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.section-header.active {
  background: rgba(255,255,255,0.06);
  color: #ffffff;
}

.section-arrow {
  font-size: 0.65rem;
  transition: transform 0.2s ease;
  flex-shrink: 0;
  color: #525252;
}

.section-arrow.rotated {
  transform: rotate(90deg);
}

.section-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.section-title-input {
  flex: 1;
  padding: 4px 8px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 6px;
  color: #ffffff;
  font-family: inherit;
  font-size: 0.85rem;
  font-weight: 500;
  outline: none;
}

.section-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.15s ease;
}

.section-header:hover .section-actions {
  opacity: 1;
}

.btn-icon {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  padding: 5px;
  border-radius: 5px;
  font-size: 0.75rem;
  transition: all 0.15s ease;
}

.btn-icon:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.btn-icon.danger:hover {
  color: #ef4444;
  background: rgba(239,68,68,0.06);
}

.pages-list {
  margin-left: 18px;
  margin-top: 2px;
}

.page-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.8rem;
  color: #737373;
  transition: all 0.15s ease;
}

.page-item:hover {
  background: rgba(255,255,255,0.04);
}

.page-item.active {
  background: rgba(255,255,255,0.06);
  color: #ffffff;
}

.page-item i {
  font-size: 0.7rem;
  color: #525252;
}

.page-item.active i {
  color: #a3a3a3;
}

.page-item span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-icon-mini {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 0.7rem;
  opacity: 0;
  transition: all 0.15s ease;
}

.page-item:hover .btn-icon-mini {
  opacity: 1;
}

.btn-icon-mini:hover {
  color: #ef4444;
}

.btn-add-page {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  padding: 7px 10px;
  width: 100%;
  text-align: left;
  border-radius: 6px;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-add-page:hover {
  background: rgba(255,255,255,0.04);
  color: #a3a3a3;
}

.editor-panel {
  flex: 1;
  padding: 32px 40px;
  overflow-y: auto;
  background: #0a0a0a;
}

.editor-empty {
  text-align: center;
  padding: 80px 20px;
  color: #525252;
}

.editor-empty .empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
}

.editor-empty p {
  font-size: 0.95rem;
  margin-bottom: 20px;
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

.btn-create:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.page-header {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  margin-bottom: 20px;
}

.page-title-input {
  flex: 1;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #ffffff;
  font-size: 1.15rem;
  font-weight: 600;
  font-family: inherit;
  outline: none;
  transition: all 0.2s ease;
  letter-spacing: -0.01em;
}

.page-title-input:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.page-title-input::placeholder {
  color: #525252;
}

/* Mobile Header */
.mobile-header {
  display: none;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  
  .topbar {
    display: none;
  }
  
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
  
  .mobile-menu-btn:hover {
    color: #e5e5e5;
  }
  
  .editor-layout {
    flex-direction: column;
  }
  
  .sections-panel {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    max-height: 280px;
    padding: 16px;
  }
  
  .editor-panel {
    padding: 20px 16px;
  }
  
  .page-title-input {
    font-size: 1rem;
  }
}
</style>