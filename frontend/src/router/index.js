import { createRouter, createWebHistory } from "vue-router";
import HomePage from "@/pages/HomePage.vue";
import AboutPage from "@/pages/AboutPage.vue";
import FormPage from "@/pages/FormPage.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: HomePage,
  },
  {
    path: "/form",
    name: "Form",
    component: FormPage,
  },
  {
    path: "/about",
    name: "About",
    component: AboutPage,
  }
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
