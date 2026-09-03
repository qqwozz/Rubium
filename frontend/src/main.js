import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router/index.js'
import { supabase } from './api/supabase'
import './style.css'

// Проверяем hash до инициализации роутера
const params = new URLSearchParams(window.location.hash.substring(1))
const accessToken = params.get('access_token')
const refreshToken = params.get('refresh_token')

console.log('HASH:', window.location.hash)
console.log('ACCESS:', accessToken)
console.log('REFRESH:', refreshToken)

if (accessToken && refreshToken) {
  supabase.auth.setSession({ access_token: accessToken, refresh_token: refreshToken })
    .then(() => {
      console.log('SESSION SET')
      window.history.replaceState(null, '', '/update-password')
    })
    .catch(console.error)
} else {
  console.log('NO TOKENS IN HASH')
  // Пробуем из search params
  const searchParams = new URLSearchParams(window.location.search)
  const searchAccess = searchParams.get('access_token')
  const searchRefresh = searchParams.get('refresh_token')
  console.log('SEARCH ACCESS:', searchAccess)
  console.log('SEARCH REFRESH:', searchRefresh)
  
  if (searchAccess && searchRefresh) {
    supabase.auth.setSession({ access_token: searchAccess, refresh_token: searchRefresh })
      .then(() => {
        console.log('SESSION SET FROM SEARCH')
        window.history.replaceState(null, '', '/update-password')
      })
      .catch(console.error)
  }
}

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')