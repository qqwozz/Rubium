<template>
  <div class="courses-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <span class="page-title">Курсы</span>
        </header>

        <div class="content">
          <div class="courses-header">
            <h1>Курсы</h1>
            <p>Бесплатные и платные курсы от лучших университетов и компаний мира</p>
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
                <span class="course-badge" :class="course.free ? 'free' : 'paid'">
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
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import coursesData from '../assets/courses.json'

const courses = ref(coursesData)
const activeSubject = ref('Все')
const sidebarRef = ref(null)

const allSubjects = computed(() => {
  const subjects = new Set(['Все'])
  courses.value.forEach(c => c.subjects.forEach(s => subjects.add(s)))
  return [...subjects]
})

const filteredCourses = computed(() => {
  if (activeSubject.value === 'Все') return courses.value
  return courses.value.filter(c => c.subjects.includes(activeSubject.value))
})
</script>

<style scoped>
.courses-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #0a0a0a;
  color: #fafafa;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
}

.page-body {
  display: flex;
  flex: 1;
}

.main-content {
  margin-left: 240px;
  flex: 1;
}

.topbar {
  padding: 20px 48px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.page-title {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: #525252;
}

.content {
  max-width: 900px;
  margin: 0 auto;
  padding: 48px 48px 96px;
}

.courses-header {
  margin-bottom: 40px;
}

.courses-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
  margin-bottom: 8px;
}

.courses-header p {
  color: #737373;
  font-size: 0.95rem;
}

.filters-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 32px;
}

.filter-btn {
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid rgba(255,255,255,0.06);
  background: transparent;
  color: #737373;
  cursor: pointer;
  font-family: inherit;
  font-size: 0.8rem;
  font-weight: 500;
  transition: all 0.15s ease;
}

.filter-btn:hover {
  color: #e5e5e5;
  border-color: rgba(255,255,255,0.1);
}

.filter-btn.active {
  background: #ffffff;
  border-color: #ffffff;
  color: #0a0a0a;
}

.courses-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.course-card {
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 14px;
  padding: 20px;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
}

.course-card:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.1);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.course-university {
  font-size: 0.75rem;
  color: #a3a3a3;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.course-badge {
  font-size: 0.65rem;
  padding: 3px 8px;
  border-radius: 8px;
  font-weight: 600;
}

.course-badge.free {
  background: rgba(255,255,255,0.06);
  color: #a3a3a3;
}

.course-badge.paid {
  background: rgba(255,255,255,0.06);
  color: #737373;
}

.course-card h3 {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 8px;
  line-height: 1.3;
  color: #e5e5e5;
}

.course-card p {
  font-size: 0.85rem;
  color: #737373;
  line-height: 1.5;
  margin-bottom: 16px;
  flex: 1;
}

.course-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
}

.course-subjects {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.subject-tag {
  font-size: 0.7rem;
  padding: 3px 8px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 6px;
  color: #737373;
  font-weight: 500;
}

.course-rating {
  font-size: 0.8rem;
  color: #a3a3a3;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.course-rating i {
  font-size: 0.7rem;
  color: #525252;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #525252;
  font-size: 0.95rem;
}

@media (max-width: 768px) {
  .page-body {
    display: block;
  }

  .main-content {
    margin-left: 0;
  }

  .topbar {
    display: none;
  }

  .content {
    padding: 32px 24px 64px;
  }

  .courses-grid {
    grid-template-columns: 1fr;
  }
}
</style>