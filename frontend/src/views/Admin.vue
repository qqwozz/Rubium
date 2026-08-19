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
          <p>Управление задачами, пользователями и тетрадями</p>
        </div>
        
        <div v-if="!auth.isAdmin" class="access-denied">
          <i class="fas fa-lock"></i>
          <p>Нет доступа</p>
        </div>
        
        <template v-else>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-value">{{ stats.tasks }}</div>
              <div class="stat-label">Задач</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">Пользователей</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ stats.notebooks }}</div>
              <div class="stat-label">Тетрадей</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ stats.today }}</div>
              <div class="stat-label">Активность сегодня</div>
            </div>
          </div>
          
          <div class="admin-tabs">
            <button 
              class="admin-tab" 
              :class="{ active: activeTab === 'tasks' }"
              @click="activeTab = 'tasks'; loadTasks()"
            >
              <i class="fas fa-tasks"></i> Задачи
            </button>
            <button 
              class="admin-tab" 
              :class="{ active: activeTab === 'users' }"
              @click="activeTab = 'users'; loadUsers()"
            >
              <i class="fas fa-users"></i> Пользователи
            </button>
            <button 
              class="admin-tab" 
              :class="{ active: activeTab === 'notebooks' }"
              @click="activeTab = 'notebooks'; loadPublicNotebooks()"
            >
              <i class="fas fa-book"></i> Тетради
            </button>
          </div>
          
          <!-- ЗАДАЧИ -->
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
            
            <div v-if="tasks.length" class="items-list">
              <div v-for="task in tasks" :key="task.id" class="item-row">
                <div class="item-info">
                  <div class="item-title">{{ task.content?.slice(0, 80) }}...</div>
                  <div class="item-meta">{{ task.subject }} | {{ task.topic }} | №{{ task.task_number }}</div>
                </div>
                <button class="btn-icon danger" @click="deleteTask(task.id)">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            
            <div v-else class="empty-state">
              <p>Нет задач</p>
            </div>
          </div>
          
          <!-- ПОЛЬЗОВАТЕЛИ -->
          <div v-if="activeTab === 'users'" class="admin-section">
            <div class="section-header">
              <h2>Пользователи</h2>
              <button class="btn-load" @click="loadUsers" :disabled="loading">
                <i v-if="loading" class="fas fa-spinner fa-spin"></i>
                <span v-else>Обновить</span>
              </button>
            </div>
            
            <div v-if="users.length" class="items-list">
              <div v-for="user in users" :key="user.id" class="item-row">
                <div class="item-info">
                  <div class="item-title">{{ user.first_name || '—' }} {{ user.last_name || '' }}</div>
                  <div class="item-meta">{{ user.email }} | Регистрация: {{ formatDate(user.created_at) }}</div>
                </div>
                <div class="item-actions">
                  <button 
                    class="btn-status" 
                    :class="{ admin: user.status === 'admin' }"
                    @click="toggleStatus(user)"
                  >
                    {{ user.status === 'admin' ? 'Админ' : 'Юзер' }}
                  </button>
                  <button class="btn-icon danger" @click="deleteUser(user)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
            
            <div v-else class="empty-state">
              <p>Нет пользователей</p>
            </div>
          </div>
          
          <!-- ТЕТРАДИ -->
          <div v-if="activeTab === 'notebooks'" class="admin-section">
            <div class="section-header">
              <h2>Публичные тетради</h2>
              <button class="btn-load" @click="loadPublicNotebooks" :disabled="loading">
                <i v-if="loading" class="fas fa-spinner fa-spin"></i>
                <span v-else>Обновить</span>
              </button>
            </div>
            
            <div v-if="notebooks.length" class="items-list">
              <div v-for="notebook in notebooks" :key="notebook.id" class="item-row">
                <div class="item-info">
                  <div class="item-title">{{ notebook.title }}</div>
                  <div class="item-meta">
                    {{ notebook.author?.first_name || 'Автор' }} | {{ notebook.views_count }} просмотров
                  </div>
                </div>
                <div class="item-actions">
                  <button 
                    class="btn-status"
                    :class="{ admin: !notebook.is_public }"
                    @click="toggleNotebookVisibility(notebook)"
                  >
                    {{ notebook.is_public ? 'Скрыть' : 'Показать' }}
                  </button>
                  <button class="btn-icon danger" @click="deleteNotebook(notebook)">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
            
            <div v-else class="empty-state">
              <p>Нет публичных тетрадей</p>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import { useAuthStore } from '../stores/auth'
import { supabase } from '../api/supabase'

const auth = useAuthStore()
const activeTab = ref('tasks')
const tasks = ref([])
const users = ref([])
const notebooks = ref([])
const loading = ref(false)
const taskSubject = ref('')
const stats = ref({ tasks: 0, users: 0, notebooks: 0, today: 0 })

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

async function loadStats() {
  try {
    const [tasksRes, usersRes, notebooksRes] = await Promise.all([
      supabase.from('tasks').select('id', { count: 'exact', head: true }),
      supabase.from('rubium_users').select('id', { count: 'exact', head: true }),
      supabase.from('notebooks').select('id', { count: 'exact', head: true })
    ])
    
    stats.value.tasks = tasksRes.count || 0
    stats.value.users = usersRes.count || 0
    stats.value.notebooks = notebooksRes.count || 0
    
    const today = new Date().toISOString().split('T')[0]
    const { count } = await supabase
      .from('tasks')
      .select('id', { count: 'exact', head: true })
      .gte('created_at', today)
    
    stats.value.today = count || 0
  } catch (e) {
    console.error(e)
  }
}

async function loadTasks() {
  loading.value = true
  try {
    let query = supabase.from('tasks').select('*').limit(50)
    if (taskSubject.value) query = query.eq('subject', taskSubject.value)
    
    const { data, error } = await query
    if (error) throw error
    tasks.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function deleteTask(taskId) {
  if (!confirm('Удалить задачу?')) return
  
  try {
    const { error } = await supabase.from('tasks').delete().eq('id', taskId)
    if (error) throw error
    tasks.value = tasks.value.filter(t => t.id !== taskId)
    await loadStats()
  } catch (e) {
    console.error(e)
  }
}

async function loadUsers() {
  loading.value = true
  try {
    const { data, error } = await supabase
      .from('rubium_users')
      .select('*')
      .order('created_at', { ascending: false })
    
    if (error) throw error
    users.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function toggleStatus(user) {
  const newStatus = user.status === 'admin' ? 'user' : 'admin'
  
  try {
    const { error } = await supabase
      .from('rubium_users')
      .update({ status: newStatus })
      .eq('id', user.id)
    
    if (error) throw error
    user.status = newStatus
  } catch (e) {
    console.error(e)
  }
}

async function deleteUser(user) {
  if (!confirm(`Удалить пользователя ${user.email}?`)) return
  
  try {
    const { error } = await supabase.from('rubium_users').delete().eq('id', user.id)
    if (error) throw error
    users.value = users.value.filter(u => u.id !== user.id)
    await loadStats()
  } catch (e) {
    console.error(e)
  }
}

async function loadPublicNotebooks() {
  loading.value = true
  try {
    const { data, error } = await supabase
      .from('notebooks')
      .select(`
        *,
        author:rubium_users!user_id(id, first_name, email)
      `)
      .eq('is_public', true)
      .order('created_at', { ascending: false })
    
    if (error) throw error
    notebooks.value = data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function toggleNotebookVisibility(notebook) {
  try {
    const { error } = await supabase
      .from('notebooks')
      .update({ is_public: !notebook.is_public })
      .eq('id', notebook.id)
    
    if (error) throw error
    notebook.is_public = !notebook.is_public
  } catch (e) {
    console.error(e)
  }
}

async function deleteNotebook(notebook) {
  if (!confirm(`Удалить тетрадь "${notebook.title}"?`)) return
  
  try {
    const { error } = await supabase.from('notebooks').delete().eq('id', notebook.id)
    if (error) throw error
    notebooks.value = notebooks.value.filter(n => n.id !== notebook.id)
    await loadStats()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadStats()
  loadTasks()
})
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.stat-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 16px;
  text-align: center;
}

.stat-value {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.7rem;
  color: #64748B;
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

.items-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-title {
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 4px;
}

.item-meta {
  font-size: 0.7rem;
  color: #64748B;
}

.item-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.btn-status {
  padding: 6px 12px;
  background: rgba(255,255,255,0.06);
  color: #94A3B8;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.7rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-status.admin {
  background: rgba(251,191,36,0.15);
  color: #FBBF24;
  border-color: rgba(251,191,36,0.2);
}

.btn-icon {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.2s;
}

.btn-icon:hover {
  background: rgba(255,255,255,0.06);
}

.btn-icon.danger:hover {
  color: #F87171;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>