import { createRouter, createWebHistory } from 'vue-router';

const HomePage = () => import('@/pages/HomePage.vue');
const AboutPage = () => import('@/pages/AboutPage.vue');
const FormPage = () => import('@/pages/FormPage.vue');

const routes = [
  { path: '/', name: 'Home', component: HomePage, meta: { title: 'Home Page' } },
  { path: '/form', name: 'Form', component: FormPage, meta: { title: 'Form Page' } },
  { path: '/about', name: 'About', component: AboutPage, meta: { title: 'About Us' } },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 };
  },
});

router.beforeEach((to) => {
  document.title = to.meta.title || 'Time Capsule Memories';
});

export default router;