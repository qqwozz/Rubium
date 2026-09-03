<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand">Rubium</div>
      <p class="login-subtitle">Восстановление пароля</p>
      
      <form @submit.prevent="handleReset">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" required placeholder="user@example.com">
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading">
          <i v-if="loading" class="fas fa-spinner fa-spin"></i>
          <span v-else>Отправить ссылку</span>
        </button>
      </form>
      
      <p class="register-link">
        Вспомнил пароль? <router-link to="/login">Войти</router-link>
      </p>
    </div>
    
    <Teleport to="body">
      <Transition name="banner">
        <div v-if="success" class="banner success-banner" @click="success = ''">
          <i class="fas fa-check-circle"></i>
          <span>Ссылка отправлена. Проверь почту.</span>
        </div>
      </Transition>
      <Transition name="banner">
        <div v-if="error" class="banner error-banner" @click="error = ''">
          <i class="fas fa-exclamation-circle"></i>
          <span>{{ error }}</span>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { supabase } from '../api/supabase'

const email = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

function getErrorMessage(e) {
  const msg = e.message?.toLowerCase() || ''
  
  if (msg.includes('fetch') || msg.includes('socket') || msg.includes('network')) {
    return 'Проблемы с сетью. Проверь подключение.'
  }
  if (msg.includes('rate limit') || msg.includes('too many')) {
    return 'Слишком много попыток. Подожди немного.'
  }
  return e.message || 'Ошибка'
}

async function handleReset() {
  loading.value = true
  error.value = ''
  success.value = ''
  
  try {
    const { error: resetError } = await supabase.auth.resetPasswordForEmail(email.value, {
      redirectTo: 'https://rubium.tech/update-password'
    })
    
    if (resetError) throw resetError
    
    success.value = 'Ссылка отправлена. Проверь почту.'
    setTimeout(() => { success.value = '' }, 5000)
  } catch (e) {
    error.value = getErrorMessage(e)
    setTimeout(() => { error.value = '' }, 5000)
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