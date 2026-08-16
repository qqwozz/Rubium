<template>
  <div class="home-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Главная</span>
      </header>
      
      <div class="content">
        <div class="welcome">
          <h1>Привет, {{ auth.userName }}</h1>
          <p>Готов учиться сегодня?</p>
        </div>
        
        <div class="cards-grid">
          <router-link to="/daily" class="feature-card">
            <div class="feature-icon" style="background: rgba(251,191,36,0.15); color: #FBBF24;">
              <i class="fas fa-calendar-day"></i>
            </div>
            <div class="feature-title">Ежедневные задания</div>
            <div class="feature-desc">6 задач в день для поддержания формы</div>
          </router-link>
          
          <router-link to="/trainer" class="feature-card">
            <div class="feature-icon" style="background: rgba(167,139,250,0.15); color: #A78BFA;">
              <i class="fas fa-pencil"></i>
            </div>
            <div class="feature-title">Тренажёр</div>
            <div class="feature-desc">Решай задачи по любым темам</div>
          </router-link>
          
          <router-link to="/notebooks" class="feature-card">
            <div class="feature-icon" style="background: rgba(52,211,153,0.15); color: #34D399;">
              <i class="fas fa-book"></i>
            </div>
            <div class="feature-title">Тетради</div>
            <div class="feature-desc">Твои конспекты и шпаргалки</div>
          </router-link>
          
          <router-link to="/community" class="feature-card">
            <div class="feature-icon" style="background: rgba(244,114,182,0.15); color: #F472B6;">
              <i class="fas fa-globe"></i>
            </div>
            <div class="feature-title">Каталог</div>
            <div class="feature-desc">Публичные тетради других учеников</div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Sidebar from '../components/Sidebar.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>

<style scoped>
.home-page {
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

.welcome {
  margin-bottom: 32px;
}

.welcome h1 {
  font-size: 1.8rem;
  font-weight: 800;
  margin-bottom: 8px;
}

.welcome p {
  color: #94A3B8;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.feature-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 24px;
  padding: 24px;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
}

.feature-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: #A78BFA;
  transform: translateY(-2px);
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  margin-bottom: 16px;
}

.feature-title {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 6px;
}

.feature-desc {
  font-size: 0.85rem;
  color: #94A3B8;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .cards-grid {
    grid-template-columns: 1fr;
  }
}
</style>