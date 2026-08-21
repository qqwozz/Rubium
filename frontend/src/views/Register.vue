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
          <input 
            v-model="password" 
            type="password" 
            required 
            minlength="8" 
            placeholder="Минимум 8 символов"
            @input="validatePassword"
          >
          <div class="password-hints" v-if="password.length > 0">
            <span :class="{ valid: hasMinLength }">
              <i :class="hasMinLength ? 'fas fa-check-circle' : 'fas fa-circle'"></i> 8+ символов
            </span>
            <span :class="{ valid: hasUpperCase }">
              <i :class="hasUpperCase ? 'fas fa-check-circle' : 'fas fa-circle'"></i> Заглавная буква
            </span>
            <span :class="{ valid: hasLowerCase }">
              <i :class="hasLowerCase ? 'fas fa-check-circle' : 'fas fa-circle'"></i> Строчная буква
            </span>
            <span :class="{ valid: isLatinOnly }">
              <i :class="isLatinOnly ? 'fas fa-check-circle' : 'fas fa-circle'"></i> Только латиница
            </span>
          </div>
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading || !isPasswordValid">
          <i v-if="loading" class="fas fa-spinner fa-spin"></i>
          <span v-else>Создать аккаунт</span>
        </button>
      </form>
      
      <p class="register-link">
        Уже есть аккаунт? <router-link to="/login">Войти</router-link>
      </p>
    </div>
    
    <Teleport to="body">
      <Transition name="banner">
        <div v-if="error" class="banner error-banner" @click="error = ''">
          <i class="fas fa-exclamation-circle"></i>
          <span>{{ error }}</span>
          <button class="banner-close" @click.stop="error = ''">
            <i class="fas fa-times"></i>
          </button>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
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

const hasMinLength = computed(() => password.value.length >= 8)
const hasUpperCase = computed(() => /[A-Z]/.test(password.value))
const hasLowerCase = computed(() => /[a-z]/.test(password.value))
const isLatinOnly = computed(() => /^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$/.test(password.value))
const isPasswordValid = computed(() => hasMinLength.value && hasUpperCase.value && hasLowerCase.value && isLatinOnly.value)

function validatePassword() {
  if (password.value.length > 0 && !isLatinOnly.value) {
    error.value = 'Пароль может содержать только латинские буквы, цифры и спецсимволы'
    setTimeout(() => { error.value = '' }, 3000)
  }
}

async function handleRegister() {
  if (!isPasswordValid.value) {
    error.value = 'Пароль не соответствует требованиям'
    setTimeout(() => { error.value = '' }, 5000)
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    const result = await auth.register(email.value, password.value, firstName.value, lastName.value)
    
    if (result.requiresVerification) {
      localStorage.setItem('pendingEmail', email.value)
      router.push('/verify-email')
    } else {
      router.push('/')
    }
  } catch (e) {
    error.value = e.message || 'Ошибка регистрации'
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

.password-hints {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 8px;
}

.password-hints span {
  font-size: 0.75rem;
  color: #64748B;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.2s;
}

.password-hints span.valid {
  color: #34D399;
}

.password-hints span i {
  font-size: 0.6rem;
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
  margin-top: 8px;
}

.btn-login:hover:not(:disabled) {
  background: #8B5CF6;
  transform: translateY(-1px);
}

.btn-login:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.banner {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 18px;
  border-radius: 16px;
  z-index: 9999;
  max-width: 400px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
  cursor: pointer;
}

.error-banner {
  background: rgba(248,113,113,0.15);
  border: 1px solid rgba(248,113,113,0.3);
  color: #F87171;
}

.banner i {
  font-size: 1.2rem;
}

.banner span {
  flex: 1;
  font-size: 0.85rem;
  font-weight: 600;
}

.banner-close {
  background: none;
  border: none;
  color: #F87171;
  cursor: pointer;
  font-size: 0.8rem;
  padding: 4px;
  border-radius: 6px;
  transition: all 0.2s;
}

.banner-close:hover {
  background: rgba(248,113,113,0.2);
}

.banner-enter-active,
.banner-leave-active {
  transition: all 0.3s;
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