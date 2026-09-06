<template>
  <button type="button" class="lang" :aria-label="t('menu.switchLanguage')" @click="toggle">
    <span
      v-for="code in SUPPORTED_LANGUAGES"
      :key="code"
      class="lang-option"
      :class="{ 'is-on': appStore.language === code }"
    >
      {{ code.toUpperCase() }}
    </span>
    <span class="lang-thumb" :style="{ transform: `translateX(${isEnglish ? 0 : 100}%)` }"></span>
  </button>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { SUPPORTED_LANGUAGES, useAppStore } from '@/store';

const { locale, t } = useI18n();
const appStore = useAppStore();

const isEnglish = computed(() => appStore.language === 'en');

function toggle() {
  const next = isEnglish.value ? 'ru' : 'en';
  appStore.setLanguage(next);
  locale.value = next;
}
</script>

<style scoped>
.lang {
  position: relative;
  display: flex;
  align-items: center;
  padding: 3px;
  border: 1px solid var(--line);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  overflow: hidden;
}

.lang-option {
  position: relative;
  z-index: 1;
  width: 38px;
  padding: 0.3rem 0;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 0.62rem;
  letter-spacing: 0.1em;
  color: var(--ink-faint);
  transition: color 0.3s ease;
}

.lang-option.is-on {
  color: var(--void-0);
}

.lang-thumb {
  position: absolute;
  top: 3px;
  bottom: 3px;
  left: 3px;
  width: 38px;
  border-radius: 999px;
  background: var(--aqua);
  box-shadow: var(--glow-aqua);
  transition: transform 0.35s cubic-bezier(0.4, 1.4, 0.5, 1);
}
</style>
