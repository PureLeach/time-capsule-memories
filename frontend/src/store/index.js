import { defineStore } from 'pinia';

export const useAppStore = defineStore('app', {
  state: () => ({
    language: localStorage.getItem('language') || 'en',
  }),

  actions: {
    setLanguage(lang) {
      this.language = lang;
      localStorage.setItem('language', lang);
    },
  },

  getters: {
    currentLanguage: (state) => state.language,
  },
});