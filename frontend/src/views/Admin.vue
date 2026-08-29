<template>
  <div class="admin-page">
    <Sidebar ref="sidebarRef" />
    
    <header class="mobile-header">
      <button class="mobile-menu-btn" @click="sidebarRef?.toggle()">
        <i class="fas fa-bars"></i>
      </button>
    </header>
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Админ-панель</span>
      </header>
      
      <div class="content">
        <div class="admin-header">
          <h1>Админ-панель</h1>
          <p>Управление пользователями и тетрадями</p>
        </div>
        
        <div v-if="!auth.isAdmin" class="access-denied">
          <div class="empty-icon"><i class="fas fa-lock"></i></div>
          <p>Нет доступа</p>
        </div>
        
        <template v-else>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-value">{{ stats.users }}</div>
              <div class="stat-label">Пользователей</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ stats.notebooks }}</div>
              <div class="stat-label">Тетрадей</div>
            </div>
          </div>
          
          <div class="admin-tabs">
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
                  <div class="item-meta">{{ user.email }} · {{ formatDate(user.created_at) }}</div>
                </div>
                <div class="item-actions">
                  <button 
                    class="btn-status" 
                    :class="{ active: user.status === 'admin' }"
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
                    {{ notebook.author?.first_name || 'Автор' }} · {{ notebook.views_count || 0 }} просмотров
                  </div>
                </div>
                <div class="item-actions">
                  <button 
                    class="btn-status"
                    :class="{ active: !notebook.is_public }"
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
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  if (!auth.isAdmin) {
    router.push('/')
  }
})

const auth = useAuthStore()
const sidebarRef = ref(null)
const activeTab = ref('users')
const users = ref([])
const notebooks = ref([])
const loading = ref(false)
const stats = ref({ users: 0, notebooks: 0 })

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

async function loadStats() {
  try {
    const [usersRes, notebooksRes] = await Promise.all([
      supabase.from('rubium_users').select('id', { count: 'exact', head: true }),
      supabase.from('notebooks').select('id', { count: 'exact', head: true })
    ])
    
    stats.value.users = usersRes.count || 0
    stats.value.notebooks = notebooksRes.count || 0
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
  loadUsers()
})
</script>

<style scoped>
.admin-page {
  display: flex;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.main-content {
  margin-left: 240px;
  flex: 1;
}

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

.mobile-header {
  display: none;
}

.admin-header {
  margin-bottom: 40px;
}

.admin-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
  margin-bottom: 6px;
}

.admin-header p {
  color: #737373;
  font-size: 0.95rem;
}

.access-denied {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 40px;
}

.stat-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 20px;
  text-align: center;
}

.stat-value {
  font-size: 1.6rem;
  font-weight: 700;
  color: #ffffff;
  margin-bottom: 4px;
  letter-spacing: -0.02em;
}

.stat-label {
  font-size: 0.8rem;
  color: #737373;
  font-weight: 500;
}

.admin-tabs {
  display: flex;
  gap: 2px;
  margin-bottom: 32px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  padding: 4px;
}

.admin-tab {
  flex: 1;
  padding: 9px 12px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.15s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.admin-tab:hover {
  color: #e5e5e5;
}

.admin-tab.active {
  background: #ffffff;
  color: #0a0a0a;
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
  font-weight: 600;
  color: #e5e5e5;
}

.btn-load {
  padding: 8px 14px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 8px;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-load:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-load:disabled {
  opacity: 0.4;
  cursor: not-allowed;
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
  padding: 14px 16px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  transition: all 0.15s ease;
}

.item-row:hover {
  background: rgba(255,255,255,0.03);
  border-color: rgba(255,255,255,0.08);
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 3px;
}

.item-meta {
  font-size: 0.78rem;
  color: #525252;
}

.item-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.btn-status {
  padding: 5px 10px;
  background: rgba(255,255,255,0.04);
  color: #737373;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  font-family: inherit;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-status:hover {
  background: rgba(255,255,255,0.06);
  color: #e5e5e5;
}

.btn-status.active {
  background: rgba(255,255,255,0.08);
  color: #ffffff;
  border-color: rgba(255,255,255,0.12);
}

.btn-icon {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  padding: 6px;
  border-radius: 6px;
  font-size: 0.8rem;
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

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #525252;
  font-size: 0.9rem;
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
  
  .content {
    padding: 32px 24px 64px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .item-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .item-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>