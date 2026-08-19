<template>
  <div class="profile-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Профиль</span>
      </header>
      
      <div class="content">
        <div class="profile-header">
          <div class="avatar" @click="triggerFileInput">
            <img v-if="auth.profile?.avatar_url" :src="auth.profile.avatar_url" alt="Аватар">
            <span v-else>{{ (auth.profile?.first_name || auth.userName || 'У')[0]?.toUpperCase() }}</span>
            <div class="avatar-overlay">
              <i class="fas fa-camera"></i>
            </div>
          </div>
          <input ref="fileInput" type="file" accept="image/*" hidden @change="uploadAvatar">
          
          <div class="profile-info">
            <h1>{{ getFullName() }}</h1>
            <p>{{ auth.user?.email }}</p>
            <div class="profile-badges">
              <span v-if="auth.isAdmin" class="admin-badge">
                <i class="fas fa-shield-halved"></i> Администратор
              </span>
              <span class="reg-date">
                <i class="fas fa-calendar"></i> С нами с {{ formatDate(auth.profile?.created_at) }}
              </span>
            </div>
          </div>
        </div>
        
        <div class="quick-actions">
          <router-link to="/notebooks" class="quick-card">
            <i class="fas fa-book"></i>
            <span>Мои тетради</span>
            <i class="fas fa-chevron-right"></i>
          </router-link>
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
        
        <div class="settings-section">
          <h2><i class="fas fa-gear"></i> Настройки</h2>
          
          <div class="settings-form">
            <div class="form-group">
              <label>Имя</label>
              <input v-model="editFirstName" type="text" placeholder="Имя">
            </div>
            
            <div class="form-group">
              <label>Фамилия</label>
              <input v-model="editLastName" type="text" placeholder="Фамилия">
            </div>
            
            <div class="form-group">
              <label>Email</label>
              <input v-model="editEmail" type="email" placeholder="Email" disabled>
              <small class="form-hint">Email менять нельзя — это твой логин</small>
            </div>
            
            <button class="btn-save" @click="saveProfile" :disabled="saving">
              <i v-if="saving" class="fas fa-spinner fa-spin"></i>
              <span v-else><i class="fas fa-check"></i> Сохранить</span>
            </button>
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
import { supabase } from '../api/supabase'

const router = useRouter()
const auth = useAuthStore()
const pinnedNotebook = ref(null)

const editFirstName = ref('')
const editLastName = ref('')
const editEmail = ref('')
const saving = ref(false)
const fileInput = ref(null)

function getFullName() {
  const parts = [auth.profile?.first_name, auth.profile?.last_name].filter(Boolean)
  return parts.join(' ') || auth.userName || 'Пользователь'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

function triggerFileInput() {
  fileInput.value?.click()
}

async function uploadAvatar(event) {
  const file = event.target.files?.[0]
  if (!file) return
  
  try {
    const { data: { session } } = await supabase.auth.getSession()
    if (!session) return
    
    const userId = auth.profile?.id
    if (!userId) return
    
    const fileExt = file.name.split('.').pop()
    const filePath = `avatars/${userId}.${fileExt}`
    
    const { error: uploadError } = await supabase.storage
      .from('avatars')
      .upload(filePath, file, {
        upsert: true,
        contentType: file.type
      })
    
    if (uploadError) throw uploadError
    
    const { data: { publicUrl } } = supabase.storage
      .from('avatars')
      .getPublicUrl(filePath)
    
    const { error: updateError } = await supabase
      .from('rubium_users')
      .update({ avatar_url: publicUrl })
      .eq('id', userId)
    
    if (updateError) throw updateError
    
    await auth.loadProfile()
  } catch (e) {
    console.error(e)
    alert('Ошибка загрузки аватара')
  }
}

async function loadPinnedNotebook(id) {
  try {
    const { data } = await supabase
      .from('notebooks')
      .select('*')
      .eq('id', id)
      .single()
    
    pinnedNotebook.value = data
  } catch (e) {
    console.error(e)
  }
}

function initForm() {
  editFirstName.value = auth.profile?.first_name || ''
  editLastName.value = auth.profile?.last_name || ''
  editEmail.value = auth.user?.email || ''
}

async function saveProfile() {
  saving.value = true
  try {
    const userId = auth.profile?.id
    if (!userId) throw new Error('Нет профиля')
    
    const { error } = await supabase
      .from('rubium_users')
      .update({
        first_name: editFirstName.value,
        last_name: editLastName.value
      })
      .eq('id', userId)
    
    if (error) throw error
    
    await auth.loadProfile()
  } catch (e) {
    console.error(e)
    alert('Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  initForm()
  
  if (auth.profile?.pinned_notebook_id) {
    loadPinnedNotebook(auth.profile.pinned_notebook_id)
  }
})
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
  margin-bottom: 24px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
  cursor: pointer;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.avatar:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay i {
  color: white;
  font-size: 1.2rem;
}

.profile-info h1 {
  font-size: 1.5rem;
  font-weight: 800;
  margin-bottom: 4px;
}

.profile-info p {
  color: #94A3B8;
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.profile-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.admin-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: rgba(251,191,36,0.15);
  color: #FBBF24;
  border-radius: 20px;
  font-size: 0.7rem;
  font-weight: 600;
}

.reg-date {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
  border-radius: 20px;
  font-size: 0.7rem;
  font-weight: 600;
}

.quick-actions {
  display: grid;
  gap: 12px;
  margin-bottom: 24px;
}

.quick-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  text-decoration: none;
  color: #F1F5F9;
  transition: all 0.3s;
}

.quick-card:hover {
  border-color: #A78BFA;
  transform: translateY(-2px);
}

.quick-card i:first-child {
  color: #A78BFA;
  font-size: 1.2rem;
}

.quick-card span {
  flex: 1;
  font-weight: 600;
}

.quick-card i:last-child {
  color: #64748B;
  font-size: 0.8rem;
}

.pinned-section {
  margin-bottom: 32px;
}

.pinned-section h2,
.settings-section h2 {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.pinned-section h2 i,
.settings-section h2 i {
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

.settings-section {
  margin-bottom: 32px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: 0.8rem;
  color: #94A3B8;
}

.form-group input {
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus {
  border-color: #A78BFA;
}

.form-group input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-hint {
  font-size: 0.7rem;
  color: #64748B;
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
  align-self: flex-start;
}

.btn-save:hover:not(:disabled) {
  background: #8B5CF6;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
}
</style>