import { createI18n } from 'vue-i18n';
import { watch } from 'vue';
import en from './en.json';
import ru from './ru.json';
import { useAppStore } from '@/store';

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en, ru },
});

// Must run after the Pinia plugin is installed.
export function initializeLanguage() {
  const appStore = useAppStore();

  i18n.global.locale.value = appStore.language;
  watch(
    () => appStore.language,
    (language) => {
      i18n.global.locale.value = language;
      document.documentElement.lang = language;
    },
    { immediate: true }
  );
}

export default i18n;
