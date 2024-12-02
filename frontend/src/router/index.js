import { createRouter, createWebHistory } from 'vue-router';
import HomePage from '@/pages/HomePage.vue';
import FormPage from '@/pages/FormPage.vue';

const routes = [
  { path: '/', component: HomePage },
  { path: '/form', component: FormPage },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
