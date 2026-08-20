<template>
  <div class="notebook-write-page">
    <Sidebar />
    
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
            <i class="fas fa-folder-open"></i>
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
            <i class="fas fa-file"></i>
            <p>Выбери или создай страницу</p>
            <button class="btn-create" @click="addPage(0)">
              <i class="fas fa-plus"></i> Создать страницу
            </button>
          </div>
          
          <div v-if="!currentPage && sections.length === 0" class="editor-empty">
            <i class="fas fa-book"></i>
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
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import NotebookEditor from '../components/NotebookEditor.vue'
import { onBeforeUnmount} from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
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

const route = useRoute()
const router = useRouter()

const notebook = ref(null)
const sections = ref([])
const openSections = ref([])
const currentPage = ref(null)
const saving = ref(false)
const editingSectionId = ref(null)
const sectionInput = ref(null)

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
}

.main-content {
  margin-left: 240px;
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
  color: #94A3B8;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: color 0.2s;
}

.back-btn:hover {
  color: #F1F5F9;
}

.page-title {
  flex: 1;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: #64748B;
  font-family: 'JetBrains Mono', monospace;
}

.btn-save {
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

.btn-save:hover:not(:disabled) {
  background: #8B5CF6;
  transform: translateY(-1px);
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.editor-layout {
  display: flex;
  min-height: calc(100vh - 60px);
}

.sections-panel {
  width: 300px;
  border-right: 1px solid rgba(255,255,255,0.06);
  padding: 20px 16px;
  flex-shrink: 0;
  overflow-y: auto;
}

.sections-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 700;
  font-size: 0.9rem;
  color: #F1F5F9;
}

.btn-add-section {
  background: rgba(167,139,250,0.15);
  border: none;
  color: #A78BFA;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-add-section:hover {
  background: rgba(167,139,250,0.25);
}

.sections-empty {
  color: #64748B;
  font-size: 0.85rem;
  text-align: center;
  padding: 40px 20px;
}

.sections-empty i {
  font-size: 2rem;
  margin-bottom: 8px;
  display: block;
}

.section-item {
  margin-bottom: 6px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.2s;
  color: #94A3B8;
}

.section-header:hover {
  background: rgba(255,255,255,0.04);
  color: #F1F5F9;
}

.section-header.active {
  background: rgba(167,139,250,0.08);
  color: #A78BFA;
}

.section-arrow {
  font-size: 0.7rem;
  transition: transform 0.2s;
  flex-shrink: 0;
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
  background: rgba(255,255,255,0.05);
  border: 1px solid #A78BFA;
  border-radius: 6px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.85rem;
  font-weight: 600;
  outline: none;
}

.section-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s;
}

.section-header:hover .section-actions {
  opacity: 1;
}

.btn-icon {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 6px;
  border-radius: 6px;
  font-size: 0.75rem;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: rgba(255,255,255,0.06);
  color: #F1F5F9;
}

.btn-icon.danger:hover {
  color: #F87171;
}

.pages-list {
  margin-left: 20px;
  margin-top: 4px;
}

.page-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.8rem;
  color: #94A3B8;
  transition: all 0.2s;
}

.page-item:hover {
  background: rgba(255,255,255,0.04);
}

.page-item.active {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
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
  color: #64748B;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  font-size: 0.7rem;
  opacity: 0;
  transition: all 0.2s;
}

.page-item:hover .btn-icon-mini {
  opacity: 1;
}

.btn-icon-mini:hover {
  color: #F87171;
}

.btn-add-page {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  padding: 8px 10px;
  width: 100%;
  text-align: left;
  border-radius: 8px;
  transition: all 0.2s;
}

.btn-add-page:hover {
  background: rgba(255,255,255,0.04);
  color: #A78BFA;
}

.editor-panel {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}

.editor-empty {
  text-align: center;
  padding: 60px;
  color: #64748B;
}

.editor-empty i {
  font-size: 3rem;
  margin-bottom: 16px;
  display: block;
}

.btn-create {
  margin-top: 16px;
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

.page-header {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  margin-bottom: 16px;
}

.page-title-input {
  flex: 1;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-size: 1.2rem;
  font-weight: 700;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}

.page-title-input:focus {
  border-color: #A78BFA;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .editor-layout {
    flex-direction: column;
  }
  .sections-panel {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    max-height: 300px;
  }
  .editor-panel {
    padding: 16px;
  }
}
</style>