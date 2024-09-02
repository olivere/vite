import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // generates .vite/manifest.json in outDir
    manifest: true,
    emptyOutDir: false,
    rollupOptions: {
      // overwrite default .html entry
      input: "/src/main.tsx",
    },
  },
})
