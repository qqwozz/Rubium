<template>
  <div class="home-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />
    
    <div class="page-body">
      <Sidebar ref="sidebarRef" />
      
      <div class="main-content">
        <header class="topbar">
          <span class="page-title">Главная</span>
        </header>
        
        <div class="content">
          <!-- Hero -->
          <section class="hero">
            <p class="hero-greeting">Готов к новым вершинам?</p>
            <h1 class="hero-name">{{ auth.userName }}</h1>
            <div class="hero-actions">
              <router-link to="/notebooks" class="btn btn-primary">
                <i class="fa-solid fa-book"></i> К моим тетрадям
              </router-link>
              <router-link to="/community" class="btn btn-secondary">
                <i class="fa-solid fa-globe"></i> В каталог тетрадей
              </router-link>
            </div>
          </section>

          <!-- О Rubium -->
          <section class="section">
            <h2 class="section-title">О Rubium</h2>
            <p class="section-body">
              <strong>Rubium</strong> — это образовательная среда нового поколения, созданная студентами для студентов. Она объединяет подготовку к экзаменам, обмен знаниями и доступ к реальным карьерным возможностям в одном цифровом пространстве. Здесь каждый учится, растёт и строит свой путь к сильному будущему.
            </p>
          </section>

          <!-- Возможности -->
          <section class="section">
            <h2 class="section-title">Возможности</h2>
            <div class="features-grid">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-book"></i></div>
                <div class="feature-content">
                  <h3>Рабочая тетрадь</h3>
                  <p>Пиши конспекты с формулами и таблицами. Автосохранение включено.</p>
                </div>
              </div>
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-globe"></i></div>
                <div class="feature-content">
                  <h3>Сообщество</h3>
                  <p>Находи публичные тетради, сохраняй их себе и делись своими конспектами.</p>
                </div>
              </div>
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-graduation-cap"></i></div>
                <div class="feature-content">
                  <h3>Каталог курсов</h3>
                  <p>Бесплатные курсы от MIT, Stanford, Harvard и других университетов.</p>
                </div>
              </div>
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-rocket"></i></div>
                <div class="feature-content">
                  <h3>Развитие</h3>
                  <p>Подготовка к экзаменам, обмен знаниями и путь к карьере — всё в одном пространстве.</p>
                </div>
              </div>
            </div>
          </section>

          <!-- С чего начать -->
          <section class="section">
            <h2 class="section-title">С чего начать</h2>
            <div class="steps-list">
              <div class="step">
                <span class="step-num">01</span>
                <div class="step-content">
                  <h3>Создай аккаунт</h3>
                  <p>Это займёт меньше минуты</p>
                </div>
              </div>
              <div class="step">
                <span class="step-num">02</span>
                <div class="step-content">
                  <h3>Веди тетрадь</h3>
                  <p>Сохраняй конспекты и делись ими</p>
                </div>
              </div>
              <div class="step">
                <span class="step-num">03</span>
                <div class="step-content">
                  <h3>Ищи возможности</h3>
                  <p>Выбери курс в каталоге и конспектируй его тут</p>
                </div>
              </div>
            </div>
          </section>

          <!-- Telegram -->
          <section class="section">
            <h2 class="section-title">Наши Telegram-каналы</h2>
            <div class="links-list">
              <a href="https://t.me/rubium_edu" target="_blank" class="link-row">
                <div class="link-icon"><i class="fab fa-telegram"></i></div>
                <div class="link-text">
                  <h3>Rubium</h3>
                  <p>Новости, обновления и анонсы платформы. Будь в курсе всего первым.</p>
                </div>
                <i class="fas fa-arrow-up-right-from-square link-arrow"></i>
              </a>
              <a href="https://t.me/+5KS-7mS6aDw0NDJi" target="_blank" class="link-row">
                <div class="link-icon"><i class="fas fa-laptop-code"></i></div>
                <div class="link-text">
                  <h3>Rubium Tech</h3>
                  <p>Хочешь быть частью команды разработчиков? Попробуй себя в реальном проекте.</p>
                </div>
                <i class="fas fa-arrow-up-right-from-square link-arrow"></i>
              </a>
            </div>
          </section>

          <!-- Создатели -->
          <section class="section">
            <h2 class="section-title">Создатели</h2>
            <div class="devs-list">
              <div v-for="dev in developers" :key="dev.name" class="dev-row">
                <img :src="dev.avatar" :alt="dev.name" class="dev-avatar" loading="lazy">
                <div class="dev-info">
                  <h4>{{ dev.name }}</h4>
                  <div class="dev-links">
                    <a :href="`https://github.com/${dev.github}`" target="_blank">GitHub</a>
                    <a :href="`https://t.me/${dev.telegram}`" target="_blank">Telegram</a>
                    <a :href="dev.website" target="_blank">Веб-сайт</a>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const sidebarRef = ref(null)

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
  max-width: 720px;
  margin: 0 auto;
  padding: 48px 48px 96px;
}

.hero {
  padding: 32px 0 64px;
}

.hero-greeting {
  font-size: 1rem;
  color: #737373;
  margin-bottom: 8px;
  line-height: 1.4;
}

.hero-name {
  font-size: 3rem;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: #ffffff;
  margin-bottom: 32px;
  line-height: 1.1;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.btn-primary {
  background: #ffffff;
  color: #0a0a0a;
  border-color: #ffffff;
}

.btn-primary:hover {
  background: #e5e5e5;
  border-color: #e5e5e5;
}

.btn-secondary {
  background: transparent;
  color: #fafafa;
  border-color: rgba(255,255,255,0.12);
}

.btn-secondary:hover {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.2);
}

.section {
  margin-bottom: 80px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 24px;
  letter-spacing: -0.01em;
}

.section-body {
  color: #a3a3a3;
  font-size: 1rem;
  line-height: 1.7;
  max-width: 640px;
}

.section-body strong {
  color: #f5f5f5;
  font-weight: 600;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 32px;
}

.feature-item {
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.feature-icon {
  width: 36px;
  height: 36px;
  border-radius: 9px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.feature-content h3 {
  font-size: 0.95rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 4px;
}

.feature-content p {
  font-size: 0.85rem;
  color: #737373;
  line-height: 1.5;
}

.steps-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.step {
  padding-top: 8px;
}

.step-num {
  display: block;
  font-size: 0.75rem;
  font-weight: 600;
  color: #404040;
  margin-bottom: 12px;
  letter-spacing: 0.05em;
}

.step-content h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 6px;
}

.step-content p {
  font-size: 0.85rem;
  color: #737373;
  line-height: 1.5;
}

.links-list {
  display: flex;
  flex-direction: column;
}

.link-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 12px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  text-decoration: none;
  color: inherit;
  transition: background 0.2s ease;
  border-radius: 8px;
}

.link-row:hover {
  background: rgba(255,255,255,0.02);
}

.link-icon {
  width: 36px;
  height: 36px;
  border-radius: 9px;
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.link-text {
  flex: 1;
}

.link-text h3 {
  font-size: 0.95rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 3px;
}

.link-text p {
  font-size: 0.82rem;
  color: #737373;
  line-height: 1.4;
}

.link-arrow {
  font-size: 0.75rem;
  color: #404040;
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.link-row:hover .link-arrow {
  color: #a3a3a3;
}

.devs-list {
  display: flex;
  flex-direction: column;
}

.dev-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.dev-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  background: #1a1a1a;
}

.dev-info h4 {
  font-size: 0.95rem;
  font-weight: 600;
  color: #e5e5e5;
  margin-bottom: 6px;
}

.dev-links {
  display: flex;
  gap: 16px;
}

.dev-links a {
  font-size: 0.8rem;
  color: #525252;
  text-decoration: none;
  transition: color 0.2s ease;
}

.dev-links a:hover {
  color: #a3a3a3;
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
  
  .hero-name {
    font-size: 2.2rem;
  }
  
  .features-grid,
  .steps-list {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  
  .hero-actions {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>