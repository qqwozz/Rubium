<template>
  <div class="rubium-tech-page">
    <Sidebar />
    
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
          <div class="feature-card">
            <i class="fas fa-tasks"></i>
            <h3>Реальные задачи</h3>
            <p>Выкладываем задачи из живого проекта. Не учебные примеры, а то, что реально нужно продукту.</p>
          </div>
          
          <div class="feature-card">
            <i class="fas fa-code-branch"></i>
            <h3>Контрибуторство</h3>
            <p>Коммиты в коммерческий продукт. Твой код видят реальные пользователи.</p>
          </div>
          
          <div class="feature-card">
            <i class="fas fa-chart-line"></i>
            <h3>Рост и скиллы</h3>
            <p>Помогаем прокачивать навыки через код-ревью и менторство.</p>
          </div>
          
          <div class="feature-card">
            <i class="fas fa-folder-open"></i>
            <h3>Портфолио</h3>
            <p>Собираем портфолио вместе. Работы в Rubium — это сильный кейс для работодателя.</p>
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
</template>

<script setup>
import Sidebar from '../components/Sidebar.vue'
</script>

<style scoped>
.rubium-tech-page {
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
  max-width: 800px;
  margin: 0 auto;
  padding: 32px;
}

.tech-header {
  text-align: center;
  margin-bottom: 40px;
}

.tech-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 24px;
  background: linear-gradient(135deg, rgba(167,139,250,0.15), rgba(244,114,182,0.1));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: #A78BFA;
}

.tech-header h1 {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 8px;
}

.tech-header p {
  color: #94A3B8;
  font-size: 1rem;
  max-width: 500px;
  margin: 0 auto;
  line-height: 1.6;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 40px;
}

.feature-card {
  background: rgba(255,255,255,0.04);
  border: 1px solid rgba(255,255,255,0.06);
  border-radius: 16px;
  padding: 24px;
  transition: all 0.3s;
}

.feature-card:hover {
  background: rgba(255,255,255,0.06);
  border-color: rgba(167,139,250,0.15);
  transform: translateY(-2px);
}

.feature-card i {
  font-size: 1.5rem;
  color: #A78BFA;
  margin-bottom: 12px;
  display: block;
}

.feature-card h3 {
  font-size: 1rem;
  font-weight: 700;
  margin-bottom: 8px;
}

.feature-card p {
  font-size: 0.85rem;
  color: #94A3B8;
  line-height: 1.5;
}

.join-section {
  text-align: center;
  padding: 40px;
  background: rgba(167,139,250,0.05);
  border: 1px solid rgba(167,139,250,0.1);
  border-radius: 24px;
}

.join-section h2 {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 24px;
}

.btn-join {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 14px 28px;
  background: linear-gradient(135deg, #A78BFA, #8B5CF6);
  color: #fff;
  border-radius: 16px;
  text-decoration: none;
  font-weight: 700;
  font-size: 0.9rem;
  transition: all 0.3s;
  margin-bottom: 16px;
}

.btn-join:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(167,139,250,0.3);
}

.contacts {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
  font-size: 0.8rem;
  color: #64748B;
}

.contacts a {
  color: #A78BFA;
  text-decoration: none;
}

.contacts a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  .content {
    padding: 16px;
  }
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>