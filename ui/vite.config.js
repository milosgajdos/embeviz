import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: "/ui",
  build: {
    assetsDir: 'assets', // Specify the directory for assets
  },
})
