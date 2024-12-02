import { createI18n } from 'vue-i18n';
import en from './en.json';
import ru from './ru.json';

const messages = { en, ru };

const i18n = createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages,
});

export default i18n;
