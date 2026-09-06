import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/pages/HomePage.vue'),
    meta: { titleKey: 'pageTitles.home' },
  },
  {
    path: '/form',
    name: 'form',
    component: () => import('@/pages/FormPage.vue'),
    meta: { titleKey: 'pageTitles.form' },
  },
  {
    path: '/about',
    name: 'about',
    component: () => import('@/pages/AboutPage.vue'),
    meta: { titleKey: 'pageTitles.about' },
  },
  // A typo or a stale link, not a page.
  { path: '/:pathMatch(.*)*', redirect: { name: 'home' } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: (to, from, savedPosition) => savedPosition || { top: 0 },
});

export default router;
