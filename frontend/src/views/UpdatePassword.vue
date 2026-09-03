<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand">Rubium</div>
      <p class="login-subtitle">Новый пароль</p>
      
      <form @submit.prevent="handleUpdate">
        <div class="form-group">
          <label>Новый пароль</label>
          <input v-model="password" type="password" required minlength="8" placeholder="Минимум 8 символов">
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading">
          <i v-if="loading" class="fas fa-spinner fa-spin"></i>
          <span v-else>Сохранить</span>
        </button>
      </form>
      
      <p v-if="error" class="error-text">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { supabase } from '../api/supabase'

const router = useRouter()
const password = ref('')
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  // Получаем токены из URL
  const hash = window.location.hash
  const params = new URLSearchParams(hash.substring(1))
  
  const accessToken = params.get('access_token')
  const refreshToken = params.get('refresh_token')
  
  if (accessToken && refreshToken) {
    try {
      const { error: sessionError } = await supabase.auth.setSession({
        access_token: accessToken,
        refresh_token: refreshToken
      })
      
      if (sessionError) throw sessionError
      
      // Очищаем URL
      window.history.replaceState(null, '', '/update-password')
    } catch (e) {
      console.error(e)
    }
  }
})

async function handleUpdate() {
  loading.value = true
  error.value = ''
  
  try {
    const { error: updateError } = await supabase.auth.updateUser({
      password: password.value
    })
    
    if (updateError) throw updateError
    
    router.push('/login')
  } catch (e) {
    error.value = e.message || 'Ошибка'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: #0a0a0a;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.login-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
}

.login-brand {
  font-size: 1.5rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 8px;
  color: #ffffff;
  letter-spacing: -0.02em;
}

.login-subtitle {
  text-align: center;
  color: #737373;
  font-size: 0.9rem;
  margin-bottom: 28px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 0.8rem;
  font-weight: 500;
  color: #a3a3a3;
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-size: 0.9rem;
  font-family: inherit;
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

.btn-login {
  width: 100%;
  padding: 12px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 8px;
}

.btn-login:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-login:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 0.85rem;
  color: #525252;
}

.register-link a {
  color: #a3a3a3;
  text-decoration: none;
  transition: color 0.15s ease;
}

.register-link a:hover {
  color: #e5e5e5;
}

.banner {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  z-index: 9999;
  max-width: 400px;
  box-shadow: 0 16px 32px rgba(0,0,0,0.4);
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
}

.error-banner {
  background: #111111;
  border: 1px solid rgba(239,68,68,0.2);
  color: #ef4444;
}

.banner i {
  font-size: 1rem;
}

.banner span {
  flex: 1;
}

.banner-close {
  background: none;
  border: none;
  color: #737373;
  cursor: pointer;
  font-size: 0.8rem;
  padding: 4px;
  border-radius: 5px;
  transition: color 0.15s ease;
}

.banner-close:hover {
  color: #e5e5e5;
}

.banner-enter-active,
.banner-leave-active {
  transition: all 0.25s ease;
}

.banner-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.banner-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>