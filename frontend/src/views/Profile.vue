<template>
  <div class="profile-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Профиль</span>
      </header>
      
      <div class="content">
        <div class="profile-header">
          <div class="avatar">{{ auth.userName[0]?.toUpperCase() || 'У' }}</div>
          <div class="profile-info">
            <h1>{{ auth.userName }}</h1>
            <p>{{ auth.user?.email }}</p>
            <span v-if="auth.isAdmin" class="admin-badge">
              <i class="fas fa-shield-halved"></i> Администратор
            </span>
          </div>
        </div>
        
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-value">{{ profile?.xp || 0 }}</div>
            <div class="stat-label">Опыт</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ profile?.tasks_solved || 0 }}</div>
            <div class="stat-label">Решено задач</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ profile?.streak || 0 }} 🔥</div>
            <div class="stat-label">Дней подряд</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ profile?.accuracy || 0 }}%</div>
            <div class="stat-label">Точность</div>
          </div>
        </div>
        
        <div v-if="pinnedNotebook" class="pinned-section">
          <h2><i class="fas fa-star"></i> Закреплённая тетрадь</h2>
          <div class="pinned-card" @click="router.push(`/notebook/${pinnedNotebook.id}`)">
            <div class="pinned-color" :style="{ background: pinnedNotebook.color || '#A78BFA' }"></div>
            <div>
              <div class="pinned-title">{{ pinnedNotebook.title }}</div>
              <div class="pinned-meta">{{ pinnedNotebook.pages_count || 0 }} страниц</div>
            </div>
          </div>
        </div>
        
        <div class="profile-actions">
          <button class="btn-logout" @click="auth.logout()">
            <i class="fas fa-sign-out-alt"></i> Выйти
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const profile = ref(null)
const pinnedNotebook = ref(null)

onMounted(() => {
  profile.value = auth.profile
  
  if (auth.profile?.pinned_notebook_id) {
    loadPinnedNotebook(auth.profile.pinned_notebook_id)
  }
})

async function loadPinnedNotebook(id) {
  try {
    const { apiFetch } = await import('../api/client')
    const data = await apiFetch(`/notebooks/${id}`)
    pinnedNotebook.value = data.notebook
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.profile-page {
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
  max-width: 700px;
  margin: 0 auto;
  padding: 32px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 32px;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
}

.profile-info h1 {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 4px;
}

.profile-info p {
  color: #94A3B8;
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.admin-badge {
  display: inline-block;
  padding: 4px 12px;
  background: rgba(251,191,36,0.15);
  color: #FBBF24;
  border-radius: 20px;
  font-size: 0.7rem;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 32px;
}

.stat-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px;
  text-align: center;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 0.75rem;
  color: #64748B;
}

.pinned-section {
  margin-bottom: 32px;
}

.pinned-section h2 {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.pinned-section h2 i {
  color: #FBBF24;
  margin-right: 8px;
}

.pinned-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.pinned-card:hover {
  border-color: #A78BFA;
}

.pinned-color {
  width: 8px;
  height: 40px;
  border-radius: 4px;
  flex-shrink: 0;
}

.pinned-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.pinned-meta {
  font-size: 0.8rem;
  color: #64748B;
}

.profile-actions {
  display: flex;
  justify-content: center;
}

.btn-logout {
  padding: 10px 24px;
  background: rgba(248,113,113,0.1);
  color: #F87171;
  border: 1px solid rgba(248,113,113,0.2);
  border-radius: 12px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-logout:hover {
  background: rgba(248,113,113,0.2);
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .stats-grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>