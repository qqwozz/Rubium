import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createHead } from '@unhead/vue/client'
import App from './App.vue'
import router from './router/index.js'
import { useTracking } from './composables/useTracking'
import './style.css'

const head = createHead()
const app = createApp(App)
app.use(createPinia())
app.use(head)
app.use(router)

useTracking().init()

app.mount('#app')