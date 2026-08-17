<template>
  <div class="home-page">
    <Sidebar />
    
    <div class="main-content">
      <header class="topbar">
        <span class="page-title">Главная</span>
      </header>
      
      <div class="content">
        <div class="hero-card">
          <div class="hero-left">
            <div class="hero-greeting">С возвращением</div>
            <div class="hero-name"><span>{{ auth.userName }}</span></div>
          </div>
          <div class="hero-actions">
            <router-link to="/daily" class="hero-btn daily">
              <i class="fa-solid fa-calendar-check"></i> Ежедневные задания
            </router-link>
            <router-link to="/trainer" class="hero-btn story">
              <i class="fa-solid fa-play"></i> Начать тренировку
            </router-link>
          </div>
        </div>

        <section class="section">
          <div class="section-title"><i class="fas fa-info-circle"></i> О Rubium</div>
          <div class="about-text">
            <strong>Rubium</strong> — экосистема для подготовки к экзаменам и учёбы в университете.
            Тренажёр с заданиями из ФИПИ, <strong>рабочая тетрадь</strong> для конспектов,
            <strong>кастомные варианты</strong> ЕГЭ/ОГЭ и <strong>сообщество</strong>, где можно делиться материалами.
            Подходит для школьников 9–11 классов, студентов и преподавателей.
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fas fa-cubes"></i> Возможности</div>
          <div class="features-grid">
            <div class="feature-card">
              <span class="icon"><i class="fas fa-brain"></i></span>
              <h3>Тренажёр</h3>
              <p>Задания из открытого банка ФИПИ с мгновенной проверкой и AI-подсказками.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-book"></i></span>
              <h3>Рабочая тетрадь</h3>
              <p>Markdown-редактор с формулами, кодом и автосохранением. Веди конспекты и публикуй их в сообществе.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-file-pen"></i></span>
              <h3>Кастомные варианты</h3>
              <p>Собирай вариант из нужных номеров, решай с таймером и получай автопроверку в баллах ЕГЭ/ОГЭ.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-chart-line"></i></span>
              <h3>Статистика</h3>
              <p>Личный прогресс по темам, точность, стрик дней и рейтинг среди других учеников.</p>
            </div>
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fas fa-route"></i> Как начать</div>
          <div class="steps-grid">
            <div class="step-card">
              <div class="step-number">1</div>
              <h3>Зарегистрируйся</h3>
              <p>Создай аккаунт в Rubium</p>
            </div>
            <div class="step-card">
              <div class="step-number">2</div>
              <h3>Выбери предмет</h3>
              <p>Математика, информатика, физика, русский язык и другие</p>
            </div>
            <div class="step-card">
              <div class="step-number">3</div>
              <h3>Начни заниматься</h3>
              <p>Решай задания, веди тетрадь, создавай варианты и следи за прогрессом</p>
            </div>
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fas fa-code"></i> Разработчики</div>
          <div class="devs-grid">
            <div v-for="dev in developers" :key="dev.name" class="dev-card">
              <div class="dev-avatar">
                <img :src="dev.avatar" :alt="dev.name" loading="lazy">
              </div>
              <div class="dev-info">
                <h4>{{ dev.name }}</h4>
                <div class="dev-links">
                  <a :href="`https://github.com/${dev.github}`" target="_blank">
                    <i class="fab fa-github"></i> {{ dev.github }}
                  </a>
                  <a :href="`https://t.me/${dev.telegram}`" target="_blank">
                    <i class="fab fa-telegram"></i> {{ dev.telegram }}
                  </a>
                  <a :href="dev.website" target="_blank">
                    <i class="fa-regular fa-address-card"></i> Веб-сайт
                  </a>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import Sidebar from '../components/Sidebar.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

const developers = [
  {
    name: 'Илья Емельянов',
    github: 'nsdmlk',
    telegram: 'KantervilleGhost',
    avatar: 'https://avatars.githubusercontent.com/nsdmlk',
    website: 'https://nsdmlk.github.io'
  },
  {
    name: 'Дима Киселёв',
    github: 'qqwozz',
    telegram: 'onixxed',
    avatar: 'https://avatars.githubusercontent.com/qqwozz',
    website: 'https://qqwozz.github.io/new_qq/'
  }
]
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
  max-width: 1000px;
  margin: 0 auto;
  padding: 24px 32px 32px;
}

/* HERO */
.hero-card {
  background: linear-gradient(135deg, rgba(167,139,250,0.1), rgba(244,114,182,0.06));
  border: 1px solid rgba(167,139,250,0.1);
  border-radius: 24px;
  padding: 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
  margin-bottom: 32px;
}

.hero-greeting {
  font-size: 0.85rem;
  color: #94A3B8;
  margin-bottom: 6px;
}

.hero-name {
  font-size: 2.4rem;
  font-weight: 800;
  letter-spacing: -1px;
}

.hero-name span {
  background: linear-gradient(135deg, #A78BFA, #F472B6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  border-radius: 16px;
  font-weight: 700;
  font-size: 0.88rem;
  text-decoration: none;
  transition: all 0.3s;
  white-space: nowrap;
}

.hero-btn.daily {
  background: linear-gradient(135deg, #FBBF24, #F59E0B);
  color: #0F0F1A;
}

.hero-btn.daily:hover {
  transform: translateY(-2px);
}

.hero-btn.story {
  background: linear-gradient(135deg, #A78BFA, #8B5CF6);
  color: #fff;
}

.hero-btn.story:hover {
  transform: translateY(-2px);
}

/* SECTIONS */
.section {
  margin-bottom: 40px;
}

.section-title {
  font-size: 1.4rem;
  font-weight: 800;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-title i {
  color: #A78BFA;
  font-size: 1.2rem;
}

.about-text {
  color: #94A3B8;
  font-size: 0.95rem;
  line-height: 1.7;
}

.about-text strong {
  color: #F1F5F9;
}

/* FEATURES */
.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.feature-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px 22px;
  transition: all 0.3s;
}

.feature-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.15);
  transform: translateY(-2px);
}

.feature-card .icon {
  font-size: 1.5rem;
  color: #A78BFA;
  margin-bottom: 10px;
  display: block;
}

.feature-card h3 {
  font-size: 0.95rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.feature-card p {
  font-size: 0.8rem;
  color: #94A3B8;
  line-height: 1.5;
}

/* STEPS */
.steps-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.step-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px 22px;
  text-align: center;
}

.step-number {
  display: inline-block;
  width: 32px;
  height: 32px;
  line-height: 32px;
  border-radius: 50%;
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
  font-weight: 800;
  font-size: 0.8rem;
  margin-bottom: 10px;
}

.step-card h3 {
  font-size: 0.9rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.step-card p {
  font-size: 0.78rem;
  color: #94A3B8;
  line-height: 1.5;
}

/* DEVELOPERS */
.devs-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.dev-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px 22px;
  display: flex;
  gap: 16px;
  align-items: center;
  transition: all 0.3s;
}

.dev-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.15);
}

.dev-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  flex-shrink: 0;
  overflow: hidden;
  background: linear-gradient(135deg, #A78BFA, #F472B6);
}

.dev-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dev-info h4 {
  font-size: 0.95rem;
  font-weight: 700;
  padding-bottom: 5px;
}

.dev-links {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.dev-links a {
  color: #64748B;
  font-size: 0.75rem;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: color 0.2s;
}

.dev-links a:hover {
  color: #A78BFA;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .hero-name {
    font-size: 1.8rem;
  }
  .features-grid,
  .steps-grid,
  .devs-grid {
    grid-template-columns: 1fr;
  }
}
</style>