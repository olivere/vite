import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],

  build: {
    outDir: "../../dist/app",
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: "/src/main.tsx",
    },
  },

  server: {
    port: 5173,
  },
})
