<template>
  <div class="profile-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

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
                <span v-if="auth.isAdmin" class="badge">
                  <i class="fas fa-shield-halved"></i> Администратор
                </span>
                <span class="badge secondary">
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
            <h2>Закреплённая тетрадь</h2>
            <div class="pinned-card" @click="router.push(`/notebook/${pinnedNotebook.id}`)">
              <div class="pinned-color" :style="{ background: pinnedNotebook.color || '#525252' }"></div>
              <div>
                <div class="pinned-title">{{ pinnedNotebook.title }}</div>
                <div class="pinned-meta">{{ pinnedNotebook.pages_count || 0 }} страниц</div>
              </div>
            </div>
          </div>

          <div class="settings-section">
            <h2>Настройки</h2>

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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { useAuthStore } from '../stores/auth'
import { supabase } from '../api/supabase'

const router = useRouter()
const auth = useAuthStore()
const pinnedNotebook = ref(null)
const sidebarRef = ref(null)

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
  flex-direction: column;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.page-body {
  display: flex;
  flex: 1;
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
  max-width: 720px;
  margin: 0 auto;
  padding: 48px 48px 96px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 40px;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.6rem;
  font-weight: 600;
  color: #a3a3a3;
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: border-color 0.2s ease;
}

.avatar:hover {
  border-color: rgba(255,255,255,0.15);
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.15s ease;
}

.avatar:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay i {
  color: #e5e5e5;
  font-size: 1.1rem;
}

.profile-info h1 {
  font-size: 1.4rem;
  font-weight: 700;
  color: #ffffff;
  margin-bottom: 4px;
  letter-spacing: -0.02em;
}

.profile-info p {
  color: #737373;
  font-size: 0.9rem;
  margin-bottom: 10px;
}

.profile-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 500;
  color: #a3a3a3;
}

.badge i {
  font-size: 0.7rem;
  color: #525252;
}

.quick-actions {
  display: grid;
  gap: 12px;
  margin-bottom: 40px;
}

.quick-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  text-decoration: none;
  color: #e5e5e5;
  transition: all 0.2s ease;
  font-weight: 500;
  font-size: 0.95rem;
}

.quick-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.quick-card i:first-child {
  color: #a3a3a3;
  font-size: 1.1rem;
  width: 20px;
  text-align: center;
}

.quick-card span {
  flex: 1;
}

.quick-card i:last-child {
  color: #525252;
  font-size: 0.75rem;
}

.pinned-section {
  margin-bottom: 40px;
}

.pinned-section h2,
.settings-section h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 14px;
  letter-spacing: -0.01em;
}

.pinned-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.pinned-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.pinned-color {
  width: 4px;
  height: 36px;
  border-radius: 2px;
  flex-shrink: 0;
}

.pinned-title {
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 3px;
  font-size: 0.95rem;
}

.pinned-meta {
  font-size: 0.8rem;
  color: #525252;
}

.settings-section {
  margin-bottom: 40px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 500;
  color: #a3a3a3;
}

.form-group input {
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  transition: all 0.2s ease;
}

.form-group input:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.form-group input::placeholder {
  color: #525252;
}

.form-group input:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.form-hint {
  font-size: 0.75rem;
  color: #525252;
}

.btn-save {
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
  align-self: flex-start;
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

.profile-actions {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

.btn-logout {
  padding: 10px 24px;
  background: transparent;
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
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

.btn-logout:hover {
  background: rgba(239,68,68,0.06);
  color: #ef4444;
  border-color: rgba(239,68,68,0.12);
}

@media (max-width: 768px) {
  .page-body {
    display: block;
  }

  .main-content {
    margin-left: 0;
  }

  .topbar {
    display: none;
  }

  .content {
    padding: 32px 24px 64px;
  }

  .profile-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .avatar {
    width: 64px;
    height: 64px;
    font-size: 1.4rem;
  }
}
</style>