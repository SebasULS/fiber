import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      // Intercepta las peticiones locales que apunten a /api
      '/api': {
        target: 'https://sisacad-enrollments-backend.vercel.app',
        changeOrigin: true,
        // Elimina el prefijo /api al enviarlo al backend real
        rewrite: (path) => path.replace(/^\/api/, '') 
      }
    }
  }
})
