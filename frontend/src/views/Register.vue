<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-brand">Rubium</div>
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
              <i :class="hasMinLength ? 'fas fa-check' : 'fas fa-circle'"></i> 8+ символов
            </span>
            <span :class="{ valid: hasUpperCase }">
              <i :class="hasUpperCase ? 'fas fa-check' : 'fas fa-circle'"></i> Заглавная буква
            </span>
            <span :class="{ valid: hasLowerCase }">
              <i :class="hasLowerCase ? 'fas fa-check' : 'fas fa-circle'"></i> Строчная буква
            </span>
            <span :class="{ valid: isLatinOnly }">
              <i :class="isLatinOnly ? 'fas fa-check' : 'fas fa-circle'"></i> Только латиница
            </span>
          </div>
        </div>
        
        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="agreed" type="checkbox" required>
            <span>
              Я согласен с 
              <a href="#" @click.prevent="openDoc('privacy')">политикой конфиденциальности</a> 
              и 
              <a href="#" @click.prevent="openDoc('terms')">пользовательским соглашением</a>
            </span>
          </label>
        </div>
        
        <button type="submit" class="btn-login" :disabled="loading || !isPasswordValid || !agreed">
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

      <Transition name="modal">
        <div v-if="showDoc" class="modal" @click.self="showDoc = false">
          <div class="modal-card doc-modal">
            <div class="modal-header">
              <h2>{{ docTitle }}</h2>
              <button class="modal-close" @click="showDoc = false">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <div class="modal-body doc-content">
              <p v-if="docLoading" class="doc-loading">Загружаем...</p>
              <div v-else v-html="docText"></div>
            </div>
          </div>
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
const agreed = ref(false)

const showDoc = ref(false)
const docTitle = ref('')
const docText = ref('')
const docLoading = ref(false)

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

async function openDoc(type) {
  docLoading.value = true
  showDoc.value = true
  
  if (type === 'privacy') {
    docTitle.value = 'Политика конфиденциальности'
  } else if (type === 'terms') {
    docTitle.value = 'Пользовательское соглашение'
  } else if (type === 'consent') {
    docTitle.value = 'Согласие на обработку данных'
  }
  
  try {
    const response = await fetch(`/docs/${type}.txt`)
    const text = await response.text()
    docText.value = text.replace(/\n/g, '<br>')
  } catch (e) {
    docText.value = 'Не удалось загрузить документ'
  } finally {
    docLoading.value = false
  }
}

async function handleRegister() {
  if (!isPasswordValid.value) {
    error.value = 'Пароль не соответствует требованиям'
    setTimeout(() => { error.value = '' }, 5000)
    return
  }

  if (!agreed.value) {
    error.value = 'Необходимо согласие с документами'
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

.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin-top: 2px;
  accent-color: #ffffff;
}

.checkbox-label span {
  font-size: 0.8rem;
  color: #737373;
  line-height: 1.4;
}

.checkbox-label a {
  color: #a3a3a3;
  text-decoration: underline;
  text-underline-offset: 2px;
  transition: color 0.15s ease;
}

.checkbox-label a:hover {
  color: #e5e5e5;
}

.password-hints {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 8px;
}

.password-hints span {
  font-size: 0.75rem;
  color: #525252;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.15s ease;
}

.password-hints span.valid {
  color: #a3a3a3;
}

.password-hints span i {
  font-size: 0.55rem;
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

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-card {
  background: #111111;
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
  padding: 28px;
  width: 100%;
  max-width: 560px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.modal-header h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #ffffff;
}

.modal-close {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 6px;
  border-radius: 6px;
  transition: color 0.15s ease;
}

.modal-close:hover {
  color: #e5e5e5;
}

.modal-body {
  overflow-y: auto;
  flex: 1;
}

.doc-content {
  color: #a3a3a3;
  font-size: 0.85rem;
  line-height: 1.6;
}

.doc-loading {
  color: #525252;
  text-align: center;
  padding: 20px;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.25s;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.25s;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: scale(0.97);
}
</style>