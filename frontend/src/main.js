import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import i18n, { initializeLanguage } from './i18n';
import './styles/main.scss';

const app = createApp(App);

app.use(createPinia()).use(router).use(i18n);

initializeLanguage();

const APP_NAME = 'Time Capsule of Memories';

router.afterEach((to) => {
  const title = to.name === 'home' ? '' : i18n.global.t(to.meta.titleKey ?? '');
  document.title = title ? `${title} · ${APP_NAME}` : APP_NAME;
});

app.mount('#app');
