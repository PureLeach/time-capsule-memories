import { defineStore } from 'pinia';

export const useAppStore = defineStore('app', {
  state: () => ({
    language: 'en',
  }),
  actions: {
    setLanguage(lang) {
      this.language = lang;
    },
  },
});
