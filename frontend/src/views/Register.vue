<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand"><span>Rub</span>ium</div>
      <p class="login-subtitle">Создай аккаунт</p>
      
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>Имя</label>
          <input v-model="firstName" type="text" required placeholder="Илья">
        </div>
        
        <div class="form-group">
          <label>Фамилия</label>
          <input v-model="lastName" type="text" placeholder="Иванов">
        </div>
        
        <div class="form-group">
          <label>Email</label>
          <input v-model="email" type="email" required placeholder="user@example.com">
        </div>
        
        <div class="form-group">
          <label>Пароль</label>
          <input v-model="password" type="password" required minlength="6" placeholder="Минимум 6 символов">
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading">
          <i v-if="loading" class="fas fa-spinner fa-spin"></i>
          <span v-else>Создать аккаунт</span>
        </button>
        
        <p v-if="error" class="error">{{ error }}</p>
      </form>
      
      <p class="register-link">
        Уже есть аккаунт? <router-link to="/login">Войти</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleRegister() {
  loading.value = true
  error.value = ''
  
  try {
    await auth.register(email.value, password.value, firstName.value, lastName.value)
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
</style>