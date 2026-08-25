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

.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin-top: 2px;
}

.checkbox-label span {
  font-size: 0.8rem;
  color: #94A3B8;
  line-height: 1.4;
}

.checkbox-label a {
  color: #A78BFA;
  text-decoration: none;
}

.checkbox-label a:hover {
  text-decoration: underline;
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

.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-card {
  background: #1a1a2e;
  border: 1px solid rgba(167,139,250,0.15);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 24px 48px rgba(0,0,0,0.4);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.modal-header h2 {
  font-size: 1.3rem;
  font-weight: 800;
}

.modal-close {
  background: none;
  border: none;
  color: #64748B;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 6px;
  border-radius: 8px;
  transition: all 0.2s;
}

.modal-close:hover {
  color: #F1F5F9;
  background: rgba(255,255,255,0.06);
}

.modal-body {
  overflow-y: auto;
  flex: 1;
}

.doc-content {
  color: #94A3B8;
  font-size: 0.85rem;
  line-height: 1.6;
}

.doc-loading {
  color: #64748B;
  text-align: center;
  padding: 20px;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
  transition: transform 0.3s;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-card,
.modal-leave-to .modal-card {
  transform: scale(0.95);
}
</style>