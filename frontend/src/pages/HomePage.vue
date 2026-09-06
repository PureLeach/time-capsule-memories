<template>
  <main-layout>
    <section class="hero">
      <div class="hero-copy">
        <span class="eyebrow">{{ t('welcome.eyebrow') }}</span>
        <h1 class="hero-title">
          <span class="line">{{ t('welcome.titleLine1') }}</span>
          <span class="line accent">{{ t('welcome.titleLine2') }}</span>
        </h1>
        <p class="hero-lead">{{ t('welcome.description') }}</p>

        <start-button to="/form">{{ t('welcome.buttonText') }}</start-button>

        <dl class="specs">
          <div v-for="spec in specs" :key="spec.key" class="spec">
            <dt class="mono">{{ t(`welcome.specs.${spec.key}.label`) }}</dt>
            <dd>{{ t(`welcome.specs.${spec.key}.value`) }}</dd>
          </div>
        </dl>
      </div>

      <div class="hero-dial">
        <chrono-dial />
      </div>
    </section>
  </main-layout>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import MainLayout from '@/layouts/MainLayout.vue';
import StartButton from '@/components/StartButton.vue';
import ChronoDial from '@/components/ChronoDial.vue';

const { t } = useI18n();
const specs = [{ key: 'delivery' }, { key: 'payload' }, { key: 'horizon' }];
</script>

<style scoped>
.hero {
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  align-items: center;
  gap: clamp(2rem, 5vw, 4rem);
  min-height: 62vh;
}

.hero-copy {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 1.5rem;
}

.hero-title {
  display: flex;
  flex-direction: column;
  font-size: clamp(2.4rem, 5.4vw, 4.1rem);
  font-weight: 300;
  line-height: 1.02;
  letter-spacing: -0.02em;
}

.line.accent {
  background: linear-gradient(100deg, var(--aqua), var(--violet) 55%, var(--amber));
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.hero-lead {
  max-width: 46ch;
  font-size: 1rem;
  line-height: 1.65;
  color: var(--ink-soft);
}

.specs {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1px;
  width: 100%;
  margin-top: 0.5rem;
  background: var(--line);
  border: 1px solid var(--line);
  border-radius: 12px;
  overflow: hidden;
}

.spec {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  padding: 0.9rem 1rem;
  background: rgba(8, 10, 22, 0.7);
}

.spec dt {
  font-size: 0.56rem;
}

.spec dd {
  font-size: 0.86rem;
  color: var(--ink);
}

.hero-dial {
  display: grid;
  place-items: center;
}

@media (max-width: 900px) {
  .hero {
    grid-template-columns: 1fr;
    text-align: center;
  }
  .hero-copy {
    align-items: center;
  }
  .hero-dial {
    order: -1;
  }
  .specs {
    grid-template-columns: 1fr;
  }
}
</style>
