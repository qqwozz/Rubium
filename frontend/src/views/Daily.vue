<template>
  <div class="daily-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Ежедневные задания</span>
      </header>
      
      <div class="content">
        <div class="daily-header">
          <h1><i class="fas fa-calendar-day"></i> Ежедневные задания</h1>
          <div class="progress">
            <span class="progress-num">{{ solvedCount }}</span>
            <span class="progress-label">/ 6</span>
          </div>
        </div>
        
        <div class="subjects-grid">
          <button 
            v-for="subject in subjects" 
            :key="subject.id"
            class="subject-card"
            :class="{ active: currentSubject === subject.id, done: solvedSubjects.includes(subject.id) }"
            @click="selectSubject(subject.id)"
            :disabled="solvedSubjects.includes(subject.id) || solvedCount >= 6"
          >
            <span class="subject-icon"><i :class="'fas ' + subject.icon"></i></span>
            <span class="subject-name">{{ subject.name }}</span>
            <span class="subject-count">{{ subject.count }} задач</span>
          </button>
        </div>
        
        <div v-if="loading" class="loading">
          <i class="fas fa-spinner fa-spin"></i> Загружаем задачи...
        </div>
        
        <div v-else-if="currentTasks.length > 0" class="tasks-list">
          <div v-for="task in currentTasks" :key="task.id" class="task-card">
            <div class="task-meta">
              <span>№{{ task.task_number || '?' }}</span>
              <span>{{ task.topic }}</span>
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
              <button class="btn-hint" @click="showHint(task)">
                <i class="fas fa-lightbulb"></i>
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
              <div v-if="results[task.id].explanation" class="solution" v-html="renderLatex(results[task.id].explanation)"></div>
            </div>
          </div>
        </div>
        
        <div v-else class="empty-state">
          <i class="fas fa-hand-peace"></i>
          <p>Выбери предмет, чтобы начать</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import { apiFetch } from '../api/client'
import katex from 'katex'

const subjects = ref([
  { id: 'math', name: 'Математика', icon: 'fa-calculator', count: 0 },
  { id: 'informatics', name: 'Информатика', icon: 'fa-laptop-code', count: 0 },
  { id: 'physics', name: 'Физика', icon: 'fa-atom', count: 0 },
  { id: 'russian', name: 'Русский язык', icon: 'fa-language', count: 0 }
])

const currentSubject = ref(null)
const currentTasks = ref([])
const loading = ref(false)
const answers = ref({})
const results = ref({})
const solvedCount = ref(0)
const solvedSubjects = ref([])

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

async function loadDailyState() {
  const stored = localStorage.getItem('rubium_daily')
  if (stored) {
    try {
      const data = JSON.parse(stored)
      const today = new Date().toISOString().split('T')[0]
      if (data.date === today) {
        solvedCount.value = data.solved || 0
        solvedSubjects.value = data.subjects || []
      }
    } catch {}
  }
}

function saveDailyState() {
  const data = {
    date: new Date().toISOString().split('T')[0],
    solved: solvedCount.value,
    subjects: solvedSubjects.value
  }
  localStorage.setItem('rubium_daily', JSON.stringify(data))
}

async function selectSubject(subjectId) {
  if (solvedSubjects.value.includes(subjectId) || solvedCount.value >= 6) return
  
  currentSubject.value = subjectId
  loading.value = true
  results.value = {}
  answers.value = {}
  
  try {
    const data = await apiFetch(`/tasks?subject=${subjectId}&limit=6`)
    currentTasks.value = data.tasks || []
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
    
    if (!results.value[task.id].solved) {
      results.value[task.id].solved = true
      solvedCount.value++
      saveDailyState()
    }
  } catch (e) {
    console.error(e)
  }
}

function showHint(task) {
  results.value[task.id] = {
    correct: false,
    correct_answer: task.answer,
    explanation: task.solution || 'Подумай над условием'
  }
}

onMounted(() => {
  loadDailyState()
})
</script>

<style scoped>
.daily-page {
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

.daily-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.daily-header h1 {
  font-size: 1.5rem;
  font-weight: 800;
}

.daily-header h1 i {
  color: #FBBF24;
  margin-right: 10px;
}

.progress {
  display: flex;
  align-items: baseline;
  gap: 4px;
  background: rgba(255,255,255,0.04);
  padding: 8px 20px;
  border-radius: 50px;
}

.progress-num {
  font-size: 1.1rem;
  font-weight: 800;
}

.progress-label {
  font-size: 0.65rem;
  color: #64748B;
  font-family: 'JetBrains Mono', monospace;
}

.subjects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
  margin-bottom: 28px;
}

.subject-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  font-family: inherit;
  color: #94A3B8;
}

.subject-card:hover:not(:disabled) {
  border-color: #A78BFA;
  transform: translateY(-2px);
}

.subject-card.active {
  background: rgba(167,139,250,0.15);
  border-color: #A78BFA;
  color: #A78BFA;
}

.subject-card.done {
  opacity: 0.4;
  cursor: default;
}

.subject-icon {
  font-size: 1.3rem;
  display: block;
  margin-bottom: 6px;
}

.subject-name {
  font-size: 0.8rem;
  font-weight: 600;
  display: block;
}

.subject-count {
  font-size: 0.65rem;
  color: #64748B;
  display: block;
  margin-top: 4px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #94A3B8;
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

.btn-hint {
  padding: 12px;
  background: rgba(251,191,36,0.15);
  color: #FBBF24;
  border: none;
  border-radius: 12px;
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

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .subjects-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>