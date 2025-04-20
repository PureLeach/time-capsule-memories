import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';
import i18n, { initializeLanguage } from './i18n';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import './styles/main.scss';

const app = createApp(App);
const pinia = createPinia();

app
  .use(pinia)             // устанавливаем Pinia
  .use(router)
  .use(i18n)              // устанавливаем i18n
  .use(ElementPlus)
  .mount('#app');

// Инициализируем языковую логику после монтирования приложения
initializeLanguage();