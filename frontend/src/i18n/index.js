import { createI18n } from 'vue-i18n';
import { watch } from 'vue';
import en from './en.json';
import ru from './ru.json';
import { useAppStore } from '@/store';

const locales = { en, ru };

// Создаем экземпляр vue-i18n с Composition API (non-legacy)
const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: locales,
});

// Функция для инициализации и синхронизации с Pinia
export function initializeLanguage() {
  const appStore = useAppStore();
  const storedLanguage = localStorage.getItem('language') || appStore.language;

  // Устанавливаем локаль при инициализации
  i18n.global.locale.value = storedLanguage;

  // Смотрим за изменениями свойства language и обновляем i18n на лету
  watch(
    () => appStore.language,
    (newLang) => {
      i18n.global.locale.value = newLang;
    }
  );
}

export default i18n;
