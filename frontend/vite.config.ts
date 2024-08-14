import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 8080,
  },
  resolve: {
    alias: {
      '@root': path.join(__dirname),
      '@src': path.join(__dirname, 'src'),
      '@pages': path.join(__dirname, 'src/pages'),
      '@routes': path.join(__dirname, 'src/routes'),
      '@factories': path.join(__dirname, 'src/factories'),
      '@styles': path.join(__dirname, 'src/styles')
    }
  }
})
