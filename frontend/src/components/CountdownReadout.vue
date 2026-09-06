<template>
  <div class="readout" :class="{ 'is-idle': !parts }">
    <span class="mono readout-label">{{ t('form.countdownLabel') }}</span>
    <div v-if="parts" class="readout-grid">
      <div v-for="unit in parts" :key="unit.key" class="unit">
        <span class="unit-value">{{ unit.value }}</span>
        <span class="unit-key">{{ t(`form.unit.${unit.key}`) }}</span>
      </div>
    </div>
    <p v-else class="readout-empty mono">{{ t('form.countdownEmpty') }}</p>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps({
  target: { type: [Date, String], default: null },
});

const { t } = useI18n();
const now = ref(Date.now());
const timer = setInterval(() => (now.value = Date.now()), 1000);
onBeforeUnmount(() => clearInterval(timer));

const parts = computed(() => {
  if (!props.target) return null;

  const target = new Date(props.target).getTime();
  const delta = target - now.value;
  if (Number.isNaN(target) || delta <= 0) return null;

  const seconds = Math.floor(delta / 1000);
  const days = Math.floor(seconds / 86400);

  return [
    { key: 'years', value: Math.floor(days / 365) },
    { key: 'days', value: days % 365 },
    { key: 'hours', value: Math.floor(seconds / 3600) % 24 },
    { key: 'minutes', value: Math.floor(seconds / 60) % 60 },
  ].map((unit) => ({ ...unit, value: String(unit.value).padStart(2, '0') }));
});
</script>

<style scoped>
.readout {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}

.readout-label {
  color: var(--amber);
  letter-spacing: 0.24em;
}

.readout-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
}

.unit {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.6rem 0.2rem;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.03);
}

.unit-value {
  font-family: var(--font-mono);
  font-size: 1.35rem;
  font-variant-numeric: tabular-nums;
  color: var(--ink);
  text-shadow: 0 0 18px rgba(94, 242, 224, 0.45);
}

.unit-key {
  font-family: var(--font-mono);
  font-size: 0.56rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: var(--ink-faint);
}

.readout-empty {
  letter-spacing: 0.1em;
  text-transform: none;
  line-height: 1.5;
}
</style>
