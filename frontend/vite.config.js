import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/ai/v1': {
        target: 'http://localhost:5080',
        changeOrigin: true
      }
    }
  }
})