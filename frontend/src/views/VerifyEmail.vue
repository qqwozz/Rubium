<template>
  <div class="verify-page">
    <div class="verify-card">
      <div class="verify-icon">
        <i class="fas fa-envelope-open-text"></i>
      </div>
      <h2>Подтверди email</h2>
      <p>Мы отправили письмо на <strong>{{ email }}</strong>. Перейди по ссылке, чтобы активировать аккаунт.</p>
      <button @click="resend" class="btn-resend" :disabled="cooldown > 0">
        {{ cooldown > 0 ? `Отправить повторно (${cooldown}s)` : 'Отправить повторно' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const email = ref('')
const cooldown = ref(0)
let timer = null

onMounted(() => {
  email.value = localStorage.getItem('pendingEmail') || ''
})

async function resend() {
  try {
    await auth.resendVerification(email.value)
    cooldown.value = 60
    timer = setInterval(() => {
      cooldown.value--
      if (cooldown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (e) {
    console.error(e)
  }
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.verify-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.verify-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.verify-icon {
  font-size: 3rem;
  color: #A78BFA;
  margin-bottom: 16px;
}

.verify-card h2 {
  color: #F1F5F9;
  margin-bottom: 12px;
}

.verify-card p {
  color: #94A3B8;
  font-size: 0.9rem;
  margin-bottom: 20px;
}

.verify-card strong {
  color: #F1F5F9;
}

.btn-resend {
  padding: 12px 24px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-resend:hover:not(:disabled) {
  background: #8B5CF6;
}

.btn-resend:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>