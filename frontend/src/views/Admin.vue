<template>
  <div class="admin-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Админ-панель</span>
      </header>
      
      <div class="content">
        <div class="admin-header">
          <h1><i class="fas fa-shield-halved"></i> Админ-панель</h1>
          <p>Управление задачами и пользователями</p>
        </div>
        
        <div v-if="!auth.isAdmin" class="access-denied">
          <i class="fas fa-lock"></i>
          <p>Нет доступа</p>
        </div>
        
        <template v-else>
          <div class="admin-tabs">
            <button 
              class="admin-tab" 
              :class="{ active: activeTab === 'tasks' }"
              @click="activeTab = 'tasks'"
            >
              <i class="fas fa-tasks"></i> Задачи
            </button>
            <button 
              class="admin-tab" 
              :class="{ active: activeTab === 'users' }"
              @click="activeTab = 'users'"
            >
              <i class="fas fa-users"></i> Пользователи
            </button>
          </div>
          
          <div v-if="activeTab === 'tasks'" class="admin-section">
            <div class="section-header">
              <h2>Задачи</h2>
              <div class="filters">
                <select v-model="taskSubject" class="filter-select">
                  <option value="">Все предметы</option>
                  <option value="math">Математика</option>
                  <option value="informatics">Информатика</option>
                  <option value="physics">Физика</option>
                  <option value="russian">Русский язык</option>
                </select>
                <button class="btn-load" @click="loadTasks" :disabled="loading">
                  <i v-if="loading" class="fas fa-spinner fa-spin"></i>
                  <span v-else>Загрузить</span>
                </button>
              </div>
            </div>
            
            <div v-if="tasks.length" class="tasks-table">
              <div v-for="task in tasks" :key="task.id" class="task-row">
                <div class="task-info">
                  <div class="task-title">{{ task.content?.slice(0, 80) }}...</div>
                  <div class="task-meta">{{ task.subject }} | {{ task.topic }} | №{{ task.task_number }}</div>
                </div>
                <button class="btn-delete" @click="deleteTask(task.id)">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            
            <div v-else class="empty-state">
              <p>Нет задач</p>
            </div>
          </div>
          
          <div v-if="activeTab === 'users'" class="admin-section">
            <div class="section-header">
              <h2>Пользователи</h2>
            </div>
            <p class="coming-soon">Раздел в разработке</p>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import { useAuthStore } from '../stores/auth'
import { apiFetch } from '../api/client'

const auth = useAuthStore()
const activeTab = ref('tasks')
const tasks = ref([])
const loading = ref(false)
const taskSubject = ref('')

async function loadTasks() {
  loading.value = true
  tasks.value = []
  
  try {
    const params = new URLSearchParams()
    if (taskSubject.value) params.append('subject', taskSubject.value)
    params.append('limit', '50')
    
    const data = await apiFetch(`/tasks?${params}`)
    tasks.value = data.tasks || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function deleteTask(taskId) {
  if (!confirm('Удалить задачу?')) return
  
  try {
    await apiFetch(`/tasks/${taskId}`, { method: 'DELETE' })
    tasks.value = tasks.value.filter(t => t.id !== taskId)
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.admin-page {
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

.admin-header {
  margin-bottom: 24px;
}

.admin-header h1 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 6px;
}

.admin-header h1 i {
  color: #FBBF24;
  margin-right: 10px;
}

.admin-header p {
  color: #94A3B8;
}

.access-denied {
  text-align: center;
  padding: 60px;
  color: #F87171;
}

.access-denied i {
  font-size: 3rem;
  margin-bottom: 16px;
}

.admin-tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 24px;
  background: rgba(255,255,255,0.04);
  border-radius: 12px;
  padding: 4px;
}

.admin-tab {
  flex: 1;
  padding: 10px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: #94A3B8;
  cursor: pointer;
  font-family: inherit;
  font-weight: 600;
  font-size: 0.85rem;
}

.admin-tab.active {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.section-header h2 {
  font-size: 1.1rem;
  font-weight: 700;
}

.filters {
  display: flex;
  gap: 8px;
}

.filter-select {
  padding: 8px 12px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.8rem;
}

.btn-load {
  padding: 8px 16px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 8px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
}

.tasks-table {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
}

.task-info {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 4px;
}

.task-meta {
  font-size: 0.7rem;
  color: #64748B;
}

.btn-delete {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 8px;
}

.btn-delete:hover {
  color: #F87171;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
}

.coming-soon {
  color: #64748B;
  text-align: center;
  padding: 40px;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
}
</style>