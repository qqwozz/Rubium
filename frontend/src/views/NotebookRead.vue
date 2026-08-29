<template>
  <div class="notebook-read-page">
    <ThinSidebar ref="sidebarRef" />
    
    <header class="mobile-header">
      <button class="mobile-menu-btn" @click="sidebarRef?.toggle()">
        <i class="fas fa-bars"></i>
      </button>
    </header>
    
    <div class="main-content">
      <header class="topbar">
        <button class="back-btn" @click="router.back()">
          <i class="fas fa-arrow-left"></i> Назад
        </button>
        <span class="page-title">{{ notebook?.title || 'Просмотр' }}</span>
        <button v-if="isOwner" class="btn-edit" @click="router.push(`/notebook/${route.params.id}/edit`)">
          <i class="fas fa-pencil"></i> Редактировать
        </button>
      </header>
      
      <div class="read-layout">
        <div class="sections-panel">
          <div v-for="section in sections" :key="section.id" class="section-item">
            <div class="section-header" @click="toggleSection(section.id)">
              <i class="fas fa-chevron-right" :class="{ rotated: openSections.includes(section.id) }"></i>
              <span>{{ section.title }}</span>
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
              </div>
            </div>
          </div>
        </div>
        
        <div class="read-panel">
          <div v-if="!currentPage" class="read-empty">
            <div class="empty-icon"><i class="fas fa-book-open"></i></div>
            <p>Выбери страницу</p>
          </div>
          
          <div v-else>
            <h1 class="page-heading">{{ currentPage.title }}</h1>
            <div class="page-content" v-html="renderContent(currentPage.content)"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import { useAuthStore } from '../stores/auth'
import katex from 'katex'
import ThinSidebar from '../components/ThinSidebar.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const sidebarRef = ref(null)

const notebook = ref(null)
const sections = ref([])
const openSections = ref([])
const currentPage = ref(null)

const isOwner = computed(() => {
  return notebook.value?.user_id === auth.user?.id
})

function renderContent(html) {
  if (!html) return ''
  
  let result = html
  
  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    try { 
      return katex.renderToString(formula, { displayMode: true, throwOnError: false }) 
    } catch { 
      return formula 
    }
  })
  
  result = result.replace(/\$([^\$\n]+?)\$/g, (_, formula) => {
    try { 
      return katex.renderToString(formula, { displayMode: false, throwOnError: false }) 
    } catch { 
      return formula 
    }
  })
  
  return `<div class="tiptap-content">${result}</div>`
}

function toggleSection(sectionId) {
  const index = openSections.value.indexOf(sectionId)
  if (index > -1) {
    openSections.value.splice(index, 1)
  } else {
    openSections.value.push(sectionId)
  }
}

function selectPage(page) {
  currentPage.value = page
}

async function loadNotebook() {
  try {
    const data = await apiFetch(`/notebooks/${route.params.id}`)
    notebook.value = data.notebook
    sections.value = data.notebook.content?.sections || []
    if (sections.value.length > 0) {
      openSections.value = [sections.value[0].id]
      if (sections.value[0].pages.length > 0) {
        currentPage.value = sections.value[0].pages[0]
      }
    }
  } catch (e) {
    console.error(e)
  }
}

onMounted(loadNotebook)
</script>

<style scoped>
.notebook-read-page {
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

.btn-edit {
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

.btn-edit:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.read-layout {
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
  color: #737373;
  transition: all 0.15s ease;
}

.section-header:hover {
  background: rgba(255,255,255,0.04);
  color: #e5e5e5;
}

.section-header i {
  font-size: 0.65rem;
  transition: transform 0.2s ease;
  color: #525252;
}

.section-header i.rotated {
  transform: rotate(90deg);
}

.section-header span {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.read-panel {
  flex: 1;
  padding: 32px 40px;
  overflow-y: auto;
  background: #0a0a0a;
}

.read-empty {
  text-align: center;
  padding: 80px 20px;
  color: #525252;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  border-radius: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #a3a3a3;
}

.page-heading {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 24px;
  color: #ffffff;
  letter-spacing: -0.02em;
}

.page-content {
  font-size: 0.95rem;
  line-height: 1.7;
  color: #a3a3a3;
}

.page-content :deep(.tiptap-content) {
  line-height: 1.7;
  color: #a3a3a3;
}

.page-content :deep(.tiptap-content h2) {
  font-size: 1.3rem;
  font-weight: 600;
  margin: 24px 0 12px;
  color: #e5e5e5;
  letter-spacing: -0.01em;
}

.page-content :deep(.tiptap-content h3) {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 20px 0 10px;
  color: #e5e5e5;
  letter-spacing: -0.01em;
}

.page-content :deep(.tiptap-content p) {
  margin-bottom: 12px;
}

.page-content :deep(.tiptap-content ul),
.page-content :deep(.tiptap-content ol) {
  margin-left: 20px;
  margin-bottom: 12px;
}

.page-content :deep(.tiptap-content li) {
  margin-bottom: 4px;
}

.page-content :deep(.tiptap-content table) {
  border-collapse: collapse;
  width: 100%;
  margin: 20px 0;
  font-size: 0.9rem;
}

.page-content :deep(.tiptap-content th),
.page-content :deep(.tiptap-content td) {
  border: 1px solid rgba(255,255,255,0.08);
  padding: 8px 12px;
}

.page-content :deep(.tiptap-content th) {
  background: rgba(255,255,255,0.04);
  font-weight: 600;
  color: #e5e5e5;
}

.page-content :deep(.tiptap-content img) {
  max-width: 100%;
  border-radius: 10px;
  margin: 12px 0;
}

.page-content :deep(.tiptap-content .katex) {
  font-size: 1.1em;
}

.page-content :deep(.katex-mathml) {
  display: none;
}

.page-content :deep(.tiptap-content blockquote) {
  border-left: 2px solid rgba(255,255,255,0.15);
  padding-left: 16px;
  margin: 16px 0;
  color: #737373;
  font-style: italic;
}

.page-content :deep(.tiptap-content code) {
  background: rgba(255,255,255,0.06);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace;
  font-size: 0.85em;
  color: #e5e5e5;
}

.page-content :deep(.tiptap-content pre) {
  background: rgba(255,255,255,0.03);
  padding: 16px;
  border-radius: 10px;
  margin: 16px 0;
  overflow-x: auto;
  border: 1px solid rgba(255,255,255,0.06);
}

.page-content :deep(.tiptap-content pre code) {
  background: none;
  padding: 0;
  color: #a3a3a3;
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
  
  .read-layout {
    flex-direction: column;
  }
  
  .sections-panel {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    max-height: 280px;
    padding: 16px;
  }
  
  .read-panel {
    padding: 20px 16px;
  }
  
  .page-heading {
    font-size: 1.25rem;
  }
}
</style>