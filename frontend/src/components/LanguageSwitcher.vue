<template>
  <button type="button" class="language-switcher" :aria-label="switchLabel" @click="toggleLanguage">
    <span class="flag-container" :class="{ flipping: isFlipping }">
      <img :src="flagSrc" :alt="flagAlt" class="flag" />
    </span>
  </button>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/store';
import enFlag from '@/assets/flags/en.svg';
import ruFlag from '@/assets/flags/ru.svg';

const FLIP_MS = 300;

const { locale, t } = useI18n();
const appStore = useAppStore();
const isFlipping = ref(false);
let timer = null;

const isEnglish = computed(() => appStore.language === 'en');
const flagSrc = computed(() => (isEnglish.value ? enFlag : ruFlag));
const flagAlt = computed(() => (isEnglish.value ? 'English' : 'Русский'));
const switchLabel = computed(() => t('menu.switchLanguage'));

function toggleLanguage() {
  isFlipping.value = true;
  clearTimeout(timer);
  timer = setTimeout(() => {
    const next = isEnglish.value ? 'ru' : 'en';
    appStore.setLanguage(next);
    locale.value = next;
    isFlipping.value = false;
  }, FLIP_MS);
}

onBeforeUnmount(() => clearTimeout(timer));
</script>

<style scoped>
.language-switcher {
  position: fixed;
  top: 10px;
  right: 10px;
  cursor: pointer;
}
.flag-container {
  width: 30px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  transition: background 0.3s ease;
  perspective: 600px;
}
.flag-container:hover {
  background: rgba(255, 255, 255, 0.4);
}
.flag {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
  backface-visibility: hidden;
  transform-style: preserve-3d;
}
.flipping .flag {
  transform: rotateY(180deg);
}
</style>
