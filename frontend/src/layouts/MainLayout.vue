<template>
  <div class="shell">
    <cosmic-background />

    <header class="bar">
      <router-link to="/" class="brand">
        <span class="brand-mark"></span>
        <span class="brand-text">
          <span class="brand-name">TIME CAPSULE</span>
          <span class="mono brand-sub">{{ t('layout.tagline') }}</span>
        </span>
      </router-link>

      <navigation-menu />

      <div class="bar-right">
        <span class="mono clock">{{ clock }}</span>
        <language-switcher />
      </div>
    </header>

    <main class="content">
      <slot />
    </main>

    <footer class="bar bar-foot">
      <span class="mono">© {{ year }} Time Capsule of Memories</span>
      <span class="mono status"><i class="dot"></i>{{ t('layout.status') }}</span>
    </footer>
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import CosmicBackground from '@/components/CosmicBackground.vue';
import NavigationMenu from '@/components/NavigationMenu.vue';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';

const { t } = useI18n();
const year = new Date().getFullYear();

// A live UTC readout, matching the timezone the dispatcher schedules in.
const clock = ref('');
const tick = () => {
  clock.value = `${new Date().toISOString().slice(11, 19)} UTC`;
};
tick();
const timer = setInterval(tick, 1000);
onBeforeUnmount(() => clearInterval(timer));
</script>

<style scoped>
.shell {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  padding: 1.1rem clamp(1rem, 4vw, 2.6rem);
  border-bottom: 1px solid var(--line);
  background: rgba(6, 8, 18, 0.55);
  backdrop-filter: blur(16px);
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  color: inherit;
}

.brand-mark {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid var(--aqua);
  background: radial-gradient(circle at 35% 30%, #fff, var(--aqua) 45%, transparent 70%);
  box-shadow: var(--glow-aqua);
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}

.brand-name {
  font-size: 0.82rem;
  letter-spacing: 0.24em;
}

.brand-sub {
  font-size: 0.54rem;
  letter-spacing: 0.18em;
}

.bar-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.clock {
  font-size: 0.64rem;
  color: var(--aqua);
  font-variant-numeric: tabular-nums;
}

.content {
  flex: 1;
  width: var(--shell);
  margin: 0 auto;
  padding: clamp(2rem, 6vh, 4.5rem) 0;
}

.bar-foot {
  border-bottom: 0;
  border-top: 1px solid var(--line);
  font-size: 0.6rem;
}

.status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--aqua);
  box-shadow: var(--glow-aqua);
  animation: blink 2.4s ease-in-out infinite;
}

@keyframes blink {
  50% {
    opacity: 0.25;
  }
}

@media (max-width: 860px) {
  .bar {
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.8rem;
  }
  .clock {
    display: none;
  }
}
</style>
