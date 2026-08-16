<template>
  <div class="trainer-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Тренажёр</span>
      </header>
      
      <div class="content">
        <div class="trainer-header">
          <h1>Тренажёр</h1>
          <p>Решай задачи по любым темам</p>
        </div>
        
        <div class="filters">
          <select v-model="subject" class="filter-select">
            <option value="">Все предметы</option>
            <option value="math">Математика</option>
            <option value="informatics">Информатика</option>
            <option value="physics">Физика</option>
            <option value="russian">Русский язык</option>
          </select>
          
          <select v-model="examType" class="filter-select">
            <option value="">Все экзамены</option>
            <option value="EGE">ЕГЭ</option>
            <option value="OGE">ОГЭ</option>
          </select>
          
          <select v-model="difficulty" class="filter-select">
            <option value="">Любая сложность</option>
            <option value="1">1</option>
            <option value="2">2</option>
            <option value="3">3</option>
            <option value="4">4</option>
            <option value="5">5</option>
          </select>
          
          <button class="btn-load" @click="loadTasks" :disabled="loading">
            <i v-if="loading" class="fas fa-spinner fa-spin"></i>
            <span v-else>Загрузить задачи</span>
          </button>
        </div>
        
        <div v-if="tasks.length === 0 && !loading" class="empty-state">
          <i class="fas fa-pencil"></i>
          <p>Выбери фильтры и загрузи задачи</p>
        </div>
        
        <div v-else class="tasks-list">
          <div v-for="task in tasks" :key="task.id" class="task-card">
            <div class="task-meta">
              <span>№{{ task.task_number || '?' }}</span>
              <span>{{ task.topic }}</span>
              <span>{{ task.exam_type }}</span>
              <span>{{ '★'.repeat(task.difficulty || 1) }}</span>
            </div>
            
            <div class="task-content" v-html="renderLatex(task.content)"></div>
            
            <div class="task-answer">
              <input 
                v-model="answers[task.id]" 
                type="text" 
                placeholder="Твой ответ"
                class="answer-input"
                @keydown.enter="checkTask(task)"
              >
              <button class="btn-check" @click="checkTask(task)" :disabled="results[task.id]">
                Проверить
              </button>
            </div>
            
            <div v-if="results[task.id]" class="task-result" :class="results[task.id].correct ? 'correct' : 'wrong'">
              <template v-if="results[task.id].correct">
                <i class="fas fa-check-circle"></i> Правильно!
              </template>
              <template v-else>
                <i class="fas fa-times-circle"></i> Неправильно
                <div>Ответ: {{ results[task.id].correct_answer }}</div>
              </template>
              <div v-if="results[task.id].explanation" class="solution">
                {{ results[task.id].explanation }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import katex from 'katex'

const subject = ref('')
const examType = ref('')
const difficulty = ref('')
const loading = ref(false)
const tasks = ref([])
const answers = ref({})
const results = ref({})

function renderLatex(text) {
  if (!text) return ''
  let result = text
  result = result.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    try { return katex.renderToString(formula, { displayMode: true, throwOnError: false }) }
    catch { return formula }
  })
  result = result.replace(/\$([^\$]*?)\$/g, (_, formula) => {
    try { return katex.renderToString(formula, { displayMode: false, throwOnError: false }) }
    catch { return formula }
  })
  return result
}

async function loadTasks() {
  loading.value = true
  tasks.value = []
  results.value = {}
  
  try {
    const params = new URLSearchParams()
    if (subject.value) params.append('subject', subject.value)
    if (examType.value) params.append('exam_type', examType.value)
    if (difficulty.value) params.append('difficulty', difficulty.value)
    params.append('limit', '10')
    
    const data = await apiFetch(`/tasks?${params}`)
    tasks.value = data.tasks || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function checkTask(task) {
  const answer = answers.value[task.id]
  if (!answer) return
  
  try {
    const data = await apiFetch('/check', {
      method: 'POST',
      body: JSON.stringify({ task_id: task.id, answer })
    })
    
    results.value[task.id] = {
      correct: data.correct,
      correct_answer: data.correct_answer,
      explanation: data.explanation
    }
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.trainer-page {
  display: flex;
  min-height: 100vh;
}

.main-content {
  margin-left: 240px;
  flex: 1;
}

.topbar {
  padding: 16px 32px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.page-title {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: #64748B;
  font-family: 'JetBrains Mono', monospace;
}

.content {
  max-width: 900px;
  margin: 0 auto;
  padding: 32px;
}

.trainer-header {
  margin-bottom: 24px;
}

.trainer-header h1 {
  font-size: 1.6rem;
  font-weight: 800;
  margin-bottom: 6px;
}

.trainer-header p {
  color: #94A3B8;
}

.filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.filter-select {
  padding: 10px 16px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  font-size: 0.85rem;
  outline: none;
}

.btn-load {
  padding: 10px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
}

.btn-load:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
}

.empty-state i {
  font-size: 3rem;
  color: #64748B;
  margin-bottom: 16px;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 24px;
}

.task-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 12px;
  font-size: 0.7rem;
  color: #64748B;
  font-family: 'JetBrains Mono', monospace;
}

.task-content {
  font-size: 0.95rem;
  line-height: 1.7;
  color: #94A3B8;
  margin-bottom: 16px;
}

.task-answer {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.answer-input {
  flex: 1;
  min-width: 200px;
  padding: 12px 16px;
  background: rgba(255,255,255,0.03);
  border: 2px solid rgba(255,255,255,0.06);
  border-radius: 12px;
  color: #F1F5F9;
  font-family: inherit;
  outline: none;
}

.answer-input:focus {
  border-color: #A78BFA;
}

.btn-check {
  padding: 12px 20px;
  background: #A78BFA;
  color: #0F0F1A;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
}

.task-result {
  margin-top: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 0.85rem;
}

.task-result.correct {
  background: rgba(52,211,153,0.1);
  color: #34D399;
}

.task-result.wrong {
  background: rgba(248,113,113,0.1);
  color: #F87171;
}

.solution {
  margin-top: 8px;
  color: #94A3B8;
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
}
</style>