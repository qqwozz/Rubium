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
          <span v-else>Сохранить</span>
        </button>
      </header>
      
      <div class="editor-layout">
        <div class="sections-panel">
          <div class="sections-header">
            <span>Разделы</span>
            <button class="btn-add-section" @click="addSection">
              <i class="fas fa-plus"></i>
            </button>
          </div>
          
          <div v-if="sections.length === 0" class="sections-empty">
            Нет разделов
          </div>
          
          <div v-for="(section, si) in sections" :key="section.id" class="section-item">
            <div class="section-header" @click="toggleSection(section.id)">
              <i class="fas fa-chevron-right" :class="{ rotated: openSections.includes(section.id) }"></i>
              <span>{{ section.title }}</span>
              <button class="btn-delete" @click.stop="deleteSection(si)">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            
            <div v-if="openSections.includes(section.id)" class="pages-list">
              <div 
                v-for="(page, pi) in section.pages" 
                :key="page.id"
                class="page-item"
                :class="{ active: currentPage?.id === page.id }"
                @click="selectPage(page)"
              >
                <i class="fas fa-file"></i>
                <span>{{ page.title || 'Без названия' }}</span>
              </div>
              
              <button class="btn-add-page" @click="addPage(si)">
                <i class="fas fa-plus"></i> Страница
              </button>
            </div>
          </div>
        </div>
        
        <div class="editor-panel">
          <div v-if="!currentPage" class="editor-empty">
            <i class="fas fa-file"></i>
            <p>Выбери или создай страницу</p>
          </div>
          
          <div v-else class="editor-content">
            <input 
              v-model="currentPage.title" 
              class="page-title-input" 
              placeholder="Заголовок страницы"
            >
            
            <textarea 
              v-model="currentPage.content" 
              class="page-content-editor"
              placeholder="Пиши конспект... Используй $...$ для формул"
              rows="20"
            ></textarea>
            
            <div v-if="currentPage.content" class="page-preview">
              <div v-html="renderLatex(currentPage.content)"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import katex from 'katex'

const route = useRoute()
const router = useRouter()

const notebook = ref(null)
const sections = ref([])
const openSections = ref([])
const currentPage = ref(null)
const saving = ref(false)

function renderLatex(text) {
  if (!text) return ''
  let result = text
  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    try { return katex.renderToString(formula, { displayMode: true, throwOnError: false }) }
    catch { return formula }
  })
  result = result.replace(/\$([^\$]*?)\$/g, (_, formula) => {
    try { return katex.renderToString(formula, { displayMode: false, throwOnError: false }) }
    catch { return formula }
  })
  return result
}

async function loadNotebook() {
  try {
    const data = await apiFetch(`/notebooks/${route.params.id}`)
    notebook.value = data.notebook
    sections.value = data.notebook.content?.sections || []
    if (sections.value.length > 0) {
      openSections.value = [sections.value[0].id]
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
}

function deleteSection(index) {
  sections.value.splice(index, 1)
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

onMounted(loadNotebook)
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
}

.btn-save:disabled {
  opacity: 0.5;
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
}

.sections-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 0.9rem;
}

.btn-add-section {
  background: rgba(167,139,250,0.15);
  border: none;
  color: #A78BFA;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  cursor: pointer;
}

.sections-empty {
  color: #64748B;
  font-size: 0.85rem;
  text-align: center;
  padding: 20px;
}

.section-item {
  margin-bottom: 8px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
}

.section-header:hover {
  background: rgba(255,255,255,0.04);
}

.section-header .fa-chevron-right {
  transition: transform 0.2s;
  font-size: 0.7rem;
}

.section-header .fa-chevron-right.rotated {
  transform: rotate(90deg);
}

.btn-delete {
  margin-left: auto;
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 4px;
}

.btn-delete:hover {
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
}

.page-item:hover {
  background: rgba(255,255,255,0.04);
}

.page-item.active {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
}

.btn-add-page {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  padding: 8px 10px;
}

.btn-add-page:hover {
  color: #A78BFA;
}

.editor-panel {
  flex: 1;
  padding: 32px;
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

.page-title-input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-size: 1.2rem;
  font-weight: 700;
  font-family: inherit;
  outline: none;
  margin-bottom: 16px;
}

.page-title-input:focus {
  border-color: #A78BFA;
}

.page-content-editor {
  width: 100%;
  padding: 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.9rem;
  outline: none;
  resize: vertical;
}

.page-content-editor:focus {
  border-color: #A78BFA;
}

.page-preview {
  margin-top: 16px;
  padding: 16px;
  background: rgba(255,255,255,0.02);
  border-radius: 12px;
  font-size: 0.95rem;
  line-height: 1.7;
  color: #94A3B8;
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
  }
  .editor-panel {
    padding: 16px;
  }
}
</style>