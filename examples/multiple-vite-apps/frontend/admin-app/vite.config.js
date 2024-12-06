import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  
  build: {
    outDir: "../../dist/admin-app",
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: "/src/main.js",
    },
  },
  server: {
    port: 5174,
  },
  base: '/admin',
}))
