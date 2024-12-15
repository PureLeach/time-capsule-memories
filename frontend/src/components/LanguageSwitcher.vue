<template>
  <div class="language-switcher">
    <div
      class="flag-container"
      :class="{'flipping': isFlipping}"
      @click="toggleLanguage"
    >
      <img
        v-if="currentLocale === 'en'"
        src="@/assets/flags/en.svg"
        alt="English"
        class="flag"
      />
      <img
        v-if="currentLocale === 'ru'"
        src="@/assets/flags/ru.svg"
        alt="Russian"
        class="flag"
      />
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import { ref } from 'vue'; // Импортируем ref для создания реактивных данных

export default {
  setup() {
    const { locale } = useI18n();
    const currentLocale = locale;
    const isFlipping = ref(false); // Используем ref для реактивной переменной isFlipping

    const changeLanguage = (value) => {
      locale.value = value;
    };

    const toggleLanguage = () => {
      isFlipping.value = true;
      setTimeout(() => {
        const newLanguage = currentLocale.value === 'en' ? 'ru' : 'en';
        changeLanguage(newLanguage);
        isFlipping.value = false;
      }, 300);
    };

    return { currentLocale, toggleLanguage, isFlipping };
  },
};
</script>

<style scoped>
.language-switcher {
  position: fixed; /* Фиксированное положение */
  top: 10px; /* Отступ сверху */
  right: 10px; /* Отступ справа */
  z-index: 9999; /* Поверх других элементов */
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.flag-container {
  width: 36px; /* Уменьшаем размер */
  height: 36px; /* Уменьшаем размер */
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  transition: background 0.3s ease;
}

.flag-container:hover {
  background: rgba(255, 255, 255, 0.4);
}

.flag {
  width: 20px; /* Уменьшаем размер иконки */
  height: 20px; /* Уменьшаем размер иконки */
  transition: transform 0.3s ease;
}

.flipping .flag {
  transform: rotateY(180deg);
}
</style>
