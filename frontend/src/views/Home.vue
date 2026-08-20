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
            <div class="hero-greeting">
              <i class="fas fa-mountain"></i> Готов к новым вершинам?
            </div>
            <div class="hero-name"><span>{{ auth.userName }}</span></div>
          </div>
          <div class="hero-actions">
            <router-link to="/trainer" class="hero-btn daily">
              <i class="fa-solid fa-play"></i> Начать тренировку
            </router-link>
            <router-link to="/notebooks" class="hero-btn story">
              <i class="fa-solid fa-book"></i> К моим тетрадям
            </router-link>
          </div>
        </div>

        <section class="section">
          <div class="section-title"><i class="fas fa-info-circle"></i> О Rubium</div>
          <div class="about-text">
            <strong>Rubium</strong> — это <em>образовательная среда</em> нового поколения, созданная <strong>студентами</strong> для <strong>студентов</strong>. 
            Она объединяет подготовку к экзаменам, обмен знаниями и доступ к реальным карьерным возможностям в одном <em>цифровом пространстве</em>. 
            Здесь каждый учится, растёт и строит свой путь к <strong>сильному будущему</strong>.
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fas fa-cubes"></i> Возможности</div>
          <div class="features-grid">
            <div class="feature-card">
              <span class="icon"><i class="fas fa-brain"></i></span>
              <h3>Тренажёр</h3>
              <p>Решай задания из открытого банка ФИПИ и проверяй себя мгновенно.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-book"></i></span>
              <h3>Рабочая тетрадь</h3>
              <p>Пиши конспекты с формулами и таблицами. Автосохранение включено.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-globe"></i></span>
              <h3>Сообщество</h3>
              <p>Находи публичные тетради, сохраняй их себе и делись своими конспектами.</p>
            </div>
            <div class="feature-card">
              <span class="icon"><i class="fas fa-rocket"></i></span>
              <h3>Развитие</h3>
              <p>Подготовка к экзаменам, обмен знаниями и путь к карьере — всё в одном пространстве.</p>
            </div>
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fas fa-route"></i> С чего начать</div>
          <div class="steps-grid">
            <div class="step-card">
              <div class="step-number">1</div>
              <h3>Создай аккаунт</h3>
              <p>Это займёт меньше минуты</p>
            </div>
            <div class="step-card">
              <div class="step-number">2</div>
              <h3>Открой тренажёр</h3>
              <p>Решай задания из банка ФИПИ</p>
            </div>
            <div class="step-card">
              <div class="step-number">3</div>
              <h3>Веди тетрадь</h3>
              <p>Сохраняй конспекты и делись ими</p>
            </div>
          </div>
        </section>

        <section class="section">
          <div class="section-title"><i class="fab fa-telegram"></i> Наши Telegram-каналы</div>
          <div class="tg-grid">
            <a href="https://t.me/rubium_edu" target="_blank" class="tg-card">
              <div class="tg-icon public">
                <i class="fab fa-telegram"></i>
              </div>
              <div class="tg-info">
                <h3>Rubium</h3>
                <p>Новости, обновления и анонсы платформы. Будь в курсе всего первым.</p>
              </div>
              <i class="fas fa-arrow-up-right-from-square tg-arrow"></i>
            </a>
            
            <a href="https://t.me/+5KS-7mS6aDw0NDJi" target="_blank" class="tg-card">
              <div class="tg-icon tech">
                <i class="fas fa-laptop-code"></i>
              </div>
              <div class="tg-info">
                <h3>Rubium Tech</h3>
                <p>Хочешь быть частью команды разработчиков? Попробуй себя в реальном проекте.</p>
              </div>
              <i class="fas fa-arrow-up-right-from-square tg-arrow"></i>
            </a>
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
  font-size: 1rem;
  color: #94A3B8;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
}

.hero-greeting i {
  margin-right: 8px;
  color: #A78BFA;
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

.tg-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.tg-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 20px 22px;
  display: flex;
  align-items: center;
  gap: 14px;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
}

.tg-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.15);
  transform: translateY(-2px);
}

.tg-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  flex-shrink: 0;
}

.tg-icon.public {
  background: rgba(56,189,248,0.15);
  color: #38BDF8;
}

.tg-icon.tech {
  background: rgba(167,139,250,0.15);
  color: #A78BFA;
}

.tg-info {
  flex: 1;
}

.tg-info h3 {
  font-size: 0.9rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.tg-info p {
  font-size: 0.78rem;
  color: #94A3B8;
  line-height: 1.4;
}

.tg-arrow {
  color: #64748B;
  font-size: 0.8rem;
  flex-shrink: 0;
}

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
  .tg-grid,
  .devs-grid {
    grid-template-columns: 1fr;
  }
}
</style>