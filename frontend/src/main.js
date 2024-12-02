import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';
import i18n from './i18n';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import './styles/main.scss';

const app = createApp(App);

app
  .use(router)
  .use(createPinia())
  .use(i18n)
  .use(ElementPlus)
  .mount('#app');
