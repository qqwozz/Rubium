<template>
  <div class="rubium-tech-page">
    <MobileHeader @toggle="sidebarRef?.toggle()" />

    <div class="page-body">
      <Sidebar ref="sidebarRef" />

      <div class="main-content">
        <header class="topbar">
          <span class="page-title">Rubium Tech</span>
        </header>

        <div class="content">
          <div class="tech-header">
            <div class="tech-icon">
              <i class="fas fa-laptop-code"></i>
            </div>
            <h1>Rubium Tech</h1>
            <p>Закрытое комьюнити для тех, кто хочет получить реальный опыт разработки.</p>
          </div>

          <div class="features-grid">
            <div class="feature-item">
              <div class="feature-icon"><i class="fas fa-tasks"></i></div>
              <div class="feature-content">
                <h3>Уникальные идеи</h3>
                <p>Каждый может предложить идеи развития. От дизайна и коммерции до технических моментов.</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon"><i class="fas fa-code-branch"></i></div>
              <div class="feature-content">
                <h3>Контрибуторство</h3>
                <p>Коммиты в коммерческий продукт. Твой код видят реальные пользователи.</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon"><i class="fas fa-chart-line"></i></div>
              <div class="feature-content">
                <h3>Рост и скиллы</h3>
                <p>Помогаем прокачивать навыки через код-ревью и менторство.</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon"><i class="fas fa-folder-open"></i></div>
              <div class="feature-content">
                <h3>Портфолио</h3>
                <p>Собираем портфолио вместе. Работы в Rubium — это сильный кейс для работодателя.</p>
              </div>
            </div>
          </div>

          <div class="join-section">
            <h2>Ты не просто наблюдаешь. Ты строишь Rubium.</h2>

            <div class="join-actions">
              <a href="https://t.me/+5KS-7mS6aDw0NDJi" target="_blank" class="btn-join">
                <i class="fab fa-telegram"></i> Присоединиться
              </a>

              <div class="contacts">
                <span>По вопросам:</span>
                <a href="https://t.me/KantervilleGhost" target="_blank">@KantervilleGhost</a>
                <a href="https://t.me/onixxed" target="_blank">@onixxed</a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import Sidebar from '../components/Sidebar.vue'
import MobileHeader from '../components/MobileHeader.vue'
import { ref } from 'vue'

const sidebarRef = ref(null)
</script>

<style scoped>
.rubium-tech-page {
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

.tech-header {
  text-align: center;
  margin-bottom: 64px;
  padding-top: 16px;
}

.tech-icon {
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

.tech-header h1 {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #ffffff;
  margin-bottom: 10px;
}

.tech-header p {
  color: #737373;
  font-size: 1rem;
  max-width: 480px;
  margin: 0 auto;
  line-height: 1.6;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 32px;
  margin-bottom: 80px;
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
  color: #a3a3a3;
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

.join-section {
  text-align: center;
  padding: 48px;
  background: rgba(255,255,255,0.02);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 18px;
}

.join-section h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 28px;
  letter-spacing: -0.01em;
}

.btn-join {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: #ffffff;
  color: #0a0a0a;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9rem;
  transition: all 0.2s ease;
  margin-bottom: 20px;
}

.btn-join:hover {
  background: #e5e5e5;
}

.contacts {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
  font-size: 0.8rem;
  color: #525252;
}

.contacts a {
  color: #a3a3a3;
  text-decoration: none;
  transition: color 0.15s ease;
}

.contacts a:hover {
  color: #e5e5e5;
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

  .features-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .join-section {
    padding: 32px 24px;
  }
}
</style>