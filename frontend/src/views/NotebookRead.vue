<template>
  <div class="notebook-read-page">
    <Sidebar />
    
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
          <div v-for="(section, si) in sections" :key="section.id" class="section-item">
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
            <i class="fas fa-book-open"></i>
            <p>Выбери страницу</p>
          </div>
          
          <div v-else>
            <h1 class="page-title">{{ currentPage.title }}</h1>
            <div class="page-content" v-html="renderLatex(currentPage.content)"></div>
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

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const notebook = ref(null)
const sections = ref([])
const openSections = ref([])
const currentPage = ref(null)

const isOwner = computed(() => {
  return notebook.value?.user_id === auth.user?.id
})

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

.btn-edit {
  padding: 10px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
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

.read-panel {
  flex: 1;
  padding: 32px;
}

.read-empty {
  text-align: center;
  padding: 60px;
  color: #64748B;
}

.read-empty i {
  font-size: 3rem;
  margin-bottom: 16px;
  display: block;
}

.read-panel h1 {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 20px;
}

.page-content {
  font-size: 0.95rem;
  line-height: 1.7;
  color: #94A3B8;
  white-space: pre-wrap;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .read-layout {
    flex-direction: column;
  }
  .sections-panel {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid rgba(255,255,255,0.06);
  }
  .read-panel {
    padding: 16px;
  }
}
</style>