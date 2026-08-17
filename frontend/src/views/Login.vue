<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand"><span>Rub</span>ium</div>
      <p class="login-subtitle">Войди, чтобы продолжить</p>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" required placeholder="user@example.com">
        </div>
        
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" required placeholder="••••••••">
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading">
          <i v-if="loading" class="fas fa-spinner fa-spin"></i>
          <span v-else>Войти</span>
        </button>
        
        <p v-if="error" class="error">{{ error }}</p>
      </form>
      
      <p class="register-link">
        Нет аккаунта? <a href="#" @click.prevent="showRegister = true">Зарегистрируйся</a>
      </p>
    </div>
    
    <div v-if="showRegister" class="modal" @click.self="showRegister = false">
      <div class="login-card">
        <div class="login-brand"><span>Rub</span>ium</div>
        <p class="login-subtitle">Создай аккаунт</p>
        
        <form @submit.prevent="handleRegister">
          <div class="form-group">
            <label>Имя</label>
            <input v-model="firstName" type="text" required>
          </div>
          
          <div class="form-group">
            <label>Фамилия</label>
            <input v-model="lastName" type="text">
          </div>
          
          <div class="form-group">
            <label>Email</label>
            <input v-model="regEmail" type="email" required>
          </div>
          
          <div class="form-group">
            <label>Пароль</label>
            <input v-model="regPassword" type="password" required minlength="6">
          </div>
          
          <button type="submit" class="btn-login" :disabled="loading">
            <i v-if="loading" class="fas fa-spinner fa-spin"></i>
            <span v-else>Создать аккаунт</span>
          </button>
          
          <p v-if="error" class="error">{{ error }}</p>
        </form>
        
        <p class="register-link">
          <a href="#" @click.prevent="showRegister = false">Назад к входу</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

console.log('Supabase URL:', import.meta.env.VITE_SUPABASE_URL)

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const firstName = ref('')
const lastName = ref('')
const regEmail = ref('')
const regPassword = ref('')
const loading = ref(false)
const error = ref('')
const showRegister = ref(false)

async function handleLogin() {
  console.log('handleLogin called')
  loading.value = true
  error.value = ''
  
  try {
    console.log('Calling auth.login...')
    await auth.login(email.value, password.value)
    console.log('Login success')
    router.push('/')
  } catch (e) {
    console.log('Login error:', e)
    error.value = e.message || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  loading.value = true
  error.value = ''
  
  try {
    await auth.register(regEmail.value, regPassword.value, firstName.value, lastName.value)
    router.push('/')
  } catch (e) {
    error.value = e.message || 'Ошибка регистрации'
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
}

.login-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
}

.login-brand {
  font-family: 'Unbounded', sans-serif;
  font-size: 1.8rem;
  font-weight: 800;
  text-align: center;
  margin-bottom: 8px;
}

.login-brand span {
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-subtitle {
  text-align: center;
  color: #94A3B8;
  font-size: 0.9rem;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 0.8rem;
  color: #94A3B8;
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-size: 0.9rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus {
  border-color: #A78BFA;
}

.btn-login {
  width: 100%;
  padding: 12px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-family: inherit;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-login:hover:not(:disabled) {
  background: #8B5CF6;
  transform: translateY(-1px);
}

.btn-login:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error {
  color: #F87171;
  font-size: 0.8rem;
  margin-top: 12px;
  text-align: center;
}

.register-link {
  text-align: center;
  margin-top: 16px;
  font-size: 0.85rem;
  color: #94A3B8;
}

.register-link a {
  color: #A78BFA;
  text-decoration: none;
}

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  padding: 20px;
}
</style>