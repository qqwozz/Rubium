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
  background: #0a0a0a;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.verify-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
  padding: 40px 32px;
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.verify-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  border-radius: 16px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: #a3a3a3;
}

.verify-card h2 {
  color: #ffffff;
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 12px;
  letter-spacing: -0.01em;
}

.verify-card p {
  color: #737373;
  font-size: 0.9rem;
  line-height: 1.6;
  margin-bottom: 24px;
}

.verify-card strong {
  color: #e5e5e5;
  font-weight: 600;
}

.btn-resend {
  padding: 10px 20px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-resend:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-resend:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>