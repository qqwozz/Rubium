<template>
  <button class="feedback-btn" @click="open" title="Сообщить о проблеме">
    <i class="fas fa-comment-dots"></i>
  </button>

  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="modal" @click.self="close">
        <div class="modal-card">
          <div class="modal-header">
            <h2>Сообщить о проблеме</h2>
            <button class="modal-close" @click="close"><i class="fas fa-times"></i></button>
          </div>

          <div class="modal-body">
            <div class="form-group">
              <label>Тип проблемы</label>
              <div class="type-row">
                <button 
                  v-for="t in types" 
                  :key="t.value" 
                  class="type-btn"
                  :class="{ active: feedbackType === t.value }"
                  @click="feedbackType = t.value"
                >
                  <i :class="t.icon"></i> {{ t.label }}
                </button>
              </div>
            </div>

            <div class="form-group">
              <label>Описание</label>
              <textarea 
                v-model="message" 
                rows="4" 
                placeholder="Опиши проблему подробнее..."
                maxlength="1000"
              ></textarea>
              <small class="form-hint">{{ message.length }}/1000</small>
            </div>
          </div>

          <div class="modal-actions">
            <button class="btn-submit" @click="submit" :disabled="!canSubmit || sending">
              <i v-if="sending" class="fas fa-spinner fa-spin"></i>
              <span v-else><i class="fas fa-paper-plane"></i> Отправить</span>
            </button>
            <button class="btn-cancel" @click="close">Отмена</button>
          </div>
        </div>
      </div>
    </Transition>

    <Transition name="modal">
      <div v-if="success" class="modal" @click.self="success = false">
        <div class="modal-card success-card">
          <div class="success-icon"><i class="fas fa-check-circle"></i></div>
          <p>Спасибо! Мы получили твоё сообщение.</p>
          <button class="btn-submit" @click="success = false">ОК</button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed } from 'vue'
import { supabase } from '../api/supabase'
import { useAuthStore } from '../stores/auth'
import { useRoute } from 'vue-router'

const auth = useAuthStore()
const route = useRoute()

const show = ref(false)
const success = ref(false)
const sending = ref(false)
const feedbackType = ref('bug')
const message = ref('')

const types = [
  { value: 'bug', label: 'Баг', icon: 'fas fa-bug' },
  { value: 'ui', label: 'Дизайн', icon: 'fas fa-palette' },
  { value: 'idea', label: 'Идея', icon: 'fas fa-lightbulb' },
  { value: 'other', label: 'Другое', icon: 'fas fa-ellipsis' }
]

const canSubmit = computed(() => message.value.trim().length >= 5)

function open() {
  show.value = true
  message.value = ''
  feedbackType.value = 'bug'
}

function close() {
  show.value = false
}

async function submit() {

  if (!auth.isAuthenticated) {
    alert('Войди в аккаунт, чтобы отправить сообщение')
    return
  }
  
  if (!canSubmit.value || sending.value) return
  sending.value = true

  try {
    const { error } = await supabase
      .from('feedback')
      .insert({
        user_id: auth.profile?.id || null,
        type: feedbackType.value,
        message: message.value.trim(),
        page_url: route.fullPath,
        user_agent: navigator.userAgent
      })

    if (error) throw error

    show.value = false
    success.value = true
  } catch (e) {
    console.error(e)
    alert('Не удалось отправить. Попробуй позже.')
  } finally {
    sending.value = false
  }
}
</script>

<style scoped>
.feedback-btn {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #111111;
  border: 1px solid rgba(255,255,255,0.08);
  color: #737373;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 500;
  box-shadow: 0 8px 24px rgba(0,0,0,0.4);
}

.feedback-btn:hover {
  background: #ffffff;
  color: #0a0a0a;
  border-color: #ffffff;
  transform: scale(1.05);
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
  max-width: 480px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h2 {
  font-size: 1.2rem;
  font-weight: 600;
  color: #ffffff;
}

.modal-close {
  background: none;
  border: none;
  color: #525252;
  cursor: pointer;
  font-size: 1rem;
  padding: 6px;
  border-radius: 6px;
  transition: all 0.15s ease;
}

.modal-close:hover {
  color: #e5e5e5;
  background: rgba(255,255,255,0.04);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 0.8rem;
  font-weight: 500;
  color: #a3a3a3;
  margin-bottom: 8px;
}

.type-row {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.type-btn {
  padding: 8px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 8px;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 500;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.type-btn:hover {
  color: #e5e5e5;
  background: rgba(255,255,255,0.05);
}

.type-btn.active {
  background: #ffffff;
  color: #0a0a0a;
  border-color: #ffffff;
}

.form-group textarea {
  width: 100%;
  padding: 10px 14px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 10px;
  color: #e5e5e5;
  font-family: inherit;
  font-size: 0.9rem;
  outline: none;
  resize: vertical;
  transition: all 0.2s ease;
  min-height: 100px;
}

.form-group textarea:focus {
  border-color: rgba(255,255,255,0.15);
  background: rgba(255,255,255,0.04);
}

.form-group textarea::placeholder {
  color: #525252;
}

.form-hint {
  display: block;
  text-align: right;
  font-size: 0.7rem;
  color: #525252;
  margin-top: 4px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.btn-cancel {
  padding: 10px 18px;
  background: rgba(255,255,255,0.04);
  color: #737373;
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.btn-cancel:hover {
  background: rgba(255,255,255,0.08);
  color: #e5e5e5;
}

.btn-submit {
  padding: 10px 18px;
  background: #ffffff;
  color: #0a0a0a;
  border: 1px solid #ffffff;
  border-radius: 10px;
  cursor: pointer;
  font-family: inherit;
  font-weight: 500;
  font-size: 0.85rem;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-submit:hover:not(:disabled) {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-submit:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.success-card {
  text-align: center;
  max-width: 320px;
}

.success-icon {
  font-size: 2.5rem;
  color: #ffffff;
  margin-bottom: 16px;
}

.success-card p {
  color: #e5e5e5;
  font-size: 0.95rem;
  margin-bottom: 20px;
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

@media (max-width: 768px) {
  .feedback-btn {
    bottom: 20px;
    right: 20px;
    width: 44px;
    height: 44px;
  }

  .modal-actions {
    flex-direction: column;
  }

  .btn-cancel,
  .btn-submit {
    width: 100%;
    justify-content: center;
  }
}
</style>