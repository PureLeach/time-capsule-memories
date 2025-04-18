import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist'
  },
  resolve: {
    alias: {
      '@': '/src', // Для удобных импортов
    },
  },
  server: {
    port: 8001, // Укажите желаемый порт
    strictPort: true, // Убедитесь, что Vite будет работать только на этом порту
  },
});
