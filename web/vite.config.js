import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const LoginHtmlFallbackPlugin = {
  name: 'login-html-fallback',
  configureServer(server) {
    server.middlewares.use('/login', (req, res, next) => {
      req.url += '.html'
      next()
    })
  }
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), LoginHtmlFallbackPlugin],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:80/',
      },
      '/auth': {
        target: 'http://127.0.0.1:80/'
      }
    }
  }
})
