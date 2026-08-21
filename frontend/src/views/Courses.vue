<template>
  <div class="courses-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Курсы</span>
      </header>
      
      <div class="content">
        <div class="courses-header">
          <h1><i class="fas fa-graduation-cap"></i> Курсы</h1>
          <p>Бесплатные и платные курсы от лучших университетов и компаний мира</p>
        </div>
        
        <div class="filters-row">
          <button 
            v-for="subject in subjects" 
            :key="subject"
            class="filter-btn"
            :class="{ active: activeSubject === subject }"
            @click="activeSubject = subject"
          >
            {{ subject }}
          </button>
        </div>
        
        <div v-if="filteredCourses.length" class="courses-grid">
          <a 
            v-for="course in filteredCourses" 
            :key="course.id"
            :href="course.url"
            target="_blank"
            class="course-card"
          >
            <div class="course-header">
              <div class="course-university">{{ course.university }}</div>
              <span class="course-free" :class="{ paid: !course.free }">
                {{ course.free ? 'Бесплатно' : 'Платно' }}
              </span>
            </div>
            
            <h3>{{ course.name }}</h3>
            <p>{{ course.description }}</p>
            
            <div class="course-footer">
              <div class="course-subjects">
                <span v-for="subj in course.subjects.slice(0, 3)" :key="subj" class="subject-tag">
                  {{ subj }}
                </span>
              </div>
              <div class="course-rating">
                <i class="fas fa-star"></i> {{ course.rating }}
              </div>
            </div>
          </a>
        </div>
        
        <div v-else class="empty-state">
          <p>Нет курсов в этой категории</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import coursesData from '../assets/courses.json'

const courses = ref(coursesData)
const activeSubject = ref('Все')
const searchQuery = ref('')

const allSubjects = computed(() => {
  const subjects = new Set(['Все'])
  courses.value.forEach(c => c.subjects.forEach(s => subjects.add(s)))
  return [...subjects]
})

const filteredCourses = computed(() => {
  let result = courses.value
  
  if (activeSubject.value !== 'Все') {
    result = result.filter(c => c.subjects.includes(activeSubject.value))
  }
  
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c => 
      c.name.toLowerCase().includes(q) ||
      c.description.toLowerCase().includes(q)
    )
  }
  
  return result
})
</script>

<style scoped>
.courses-page {
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
  max-width: 1000px;
  margin: 0 auto;
  padding: 32px;
}

.courses-header {
  margin-bottom: 24px;
}

.courses-header h1 {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 8px;
}

.courses-header h1 i {
  color: #A78BFA;
  margin-right: 12px;
}

.courses-header p {
  color: #94A3B8;
}

.filters-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.filter-btn {
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.06);
  background: transparent;
  color: #94A3B8;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.75rem;
  font-weight: 600;
  transition: all 0.2s;
}

.filter-btn:hover {
  color: #F1F5F9;
}

.filter-btn.active {
  background: rgba(167,139,250,0.15);
  border-color: #A78BFA;
  color: #A78BFA;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 16px;
}

.course-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
}

.course-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.15);
  transform: translateY(-2px);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.course-university {
  font-size: 0.7rem;
  color: #A78BFA;
  font-weight: 600;
}

.course-free {
  font-size: 0.65rem;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(52,211,153,0.1);
  color: #34D399;
  font-weight: 600;
}

.course-free.paid {
  background: rgba(251,191,36,0.1);
  color: #FBBF24;
}

.course-card h3 {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 8px;
  line-height: 1.3;
}

.course-card p {
  font-size: 0.8rem;
  color: #94A3B8;
  line-height: 1.5;
  margin-bottom: 16px;
}

.course-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.course-subjects {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.subject-tag {
  font-size: 0.65rem;
  padding: 2px 8px;
  background: rgba(167,139,250,0.08);
  border-radius: 6px;
  color: #A78BFA;
}

.course-rating {
  font-size: 0.75rem;
  color: #FBBF24;
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-state {
  text-align: center;
  padding: 60px;
  color: #94A3B8;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .courses-grid {
    grid-template-columns: 1fr;
  }
}
</style>