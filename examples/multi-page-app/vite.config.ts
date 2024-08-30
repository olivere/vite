import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    manifest: true,
    emptyOutDir: false,
    rollupOptions: {
      input: {
        main: "/src/main.tsx",
        nested: "/src/nested.tsx",
      },
    },
  },
})
