<template>
  <div class="language-switcher">
    <div class="flag-container" :class="{ flipping: isFlipping }" @click="toggleLanguage">
      <img :src="flagSrc" :alt="flagAlt" class="flag" />
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n';
import { ref, computed } from 'vue';
import { useAppStore } from '@/store';

export default {
  name: 'LanguageSwitcher',
  setup() {
    const { locale } = useI18n();
    const appStore = useAppStore();
    const isFlipping = ref(false);

    const flagSrc = computed(() =>
      appStore.language === 'en'
        ? new URL('@/assets/flags/en.svg', import.meta.url).href
        : new URL('@/assets/flags/ru.svg', import.meta.url).href
    );

    const flagAlt = computed(() => (appStore.language === 'en' ? 'English' : 'Russian'));

    const toggleLanguage = () => {
      isFlipping.value = true;
      setTimeout(() => {
        const newLanguage = appStore.language === 'en' ? 'ru' : 'en';
        appStore.setLanguage(newLanguage);
        locale.value = newLanguage;
        isFlipping.value = false;
      }, 300);
    };

    return { isFlipping, flagSrc, flagAlt, toggleLanguage };
  },
};
</script>

<style scoped>
.language-switcher { position: fixed; top: 10px; right: 10px; cursor: pointer; }
.flag-container {
  width: 30px; height: 30px;
  display: flex; justify-content: center; align-items: center;
  border-radius: 50%; background: rgba(255,255,255,0.2);
  transition: background 0.3s ease; perspective: 600px;
}
.flag-container:hover { background: rgba(255,255,255,0.4); }
.flag {
  width: 20px; height: 20px;
  transition: transform 0.3s ease;
  backface-visibility: hidden;
  transform-style: preserve-3d;
}
.flipping .flag { transform: rotateY(180deg); }
</style>