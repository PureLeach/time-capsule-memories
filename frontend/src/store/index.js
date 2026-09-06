import { defineStore } from 'pinia';

export const SUPPORTED_LANGUAGES = ['en', 'ru'];
const DEFAULT_LANGUAGE = 'en';
const STORAGE_KEY = 'language';

// localStorage throws when cookies are blocked; the app should still work,
// just without remembering the choice.
function readStoredLanguage() {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    return SUPPORTED_LANGUAGES.includes(stored) ? stored : DEFAULT_LANGUAGE;
  } catch {
    return DEFAULT_LANGUAGE;
  }
}

export const useAppStore = defineStore('app', {
  state: () => ({
    language: readStoredLanguage(),
  }),

  actions: {
    setLanguage(language) {
      if (!SUPPORTED_LANGUAGES.includes(language)) return;

      this.language = language;
      try {
        localStorage.setItem(STORAGE_KEY, language);
      } catch {
        // Preference is kept for this session only.
      }
    },
  },
});
