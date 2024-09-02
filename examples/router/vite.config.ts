import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
import { default as viteReact } from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    TanStackRouterVite(),
    viteReact(),
  ],
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
