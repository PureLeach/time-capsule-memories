<template>
  <figure class="capsule" :class="{ 'is-sealed': sealed }">
    <div class="capsule-shell">
      <div class="capsule-glass">
        <div class="capsule-scan"></div>
        <div class="capsule-content">
          <span class="mono capsule-tag">{{ t('form.preview.tag') }}</span>
          <p class="capsule-message">{{ message || t('form.preview.placeholder') }}</p>
          <div class="capsule-meta">
            <span class="mono">{{ sender || '—' }}</span>
            <span class="capsule-arrow">→</span>
            <span class="mono">{{ recipient || '—' }}</span>
          </div>
          <div v-if="attachments" class="capsule-chips">
            <span v-for="n in attachments" :key="n" class="chip"></span>
            <span class="mono chip-count">{{
              t('form.preview.images', { count: attachments })
            }}</span>
          </div>
        </div>
      </div>
      <div class="capsule-ring capsule-ring-top"></div>
      <div class="capsule-ring capsule-ring-bottom"></div>
    </div>
    <figcaption class="mono capsule-caption">
      {{ sealed ? t('form.preview.sealed') : t('form.preview.open') }}
    </figcaption>
  </figure>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

defineProps({
  sender: { type: String, default: '' },
  recipient: { type: String, default: '' },
  message: { type: String, default: '' },
  attachments: { type: Number, default: 0 },
  sealed: { type: Boolean, default: false },
});

const { t } = useI18n();
</script>

<style scoped>
.capsule {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.9rem;
}

.capsule-shell {
  position: relative;
  width: 100%;
  max-width: 300px;
  padding: 1.1rem 0;
}

.capsule-glass {
  position: relative;
  overflow: hidden;
  border-radius: 130px / 34px;
  border: 1px solid var(--line-bright);
  background: linear-gradient(160deg, rgba(94, 242, 224, 0.1), rgba(155, 123, 255, 0.08));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.14),
    0 0 46px rgba(94, 242, 224, 0.18);
  transition:
    box-shadow 0.6s ease,
    transform 0.6s ease;
}

.is-sealed .capsule-glass {
  transform: scale(0.97);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 0 70px rgba(255, 196, 107, 0.4);
  border-color: var(--amber);
}

/* A slow sweep of light, as if the capsule were being read by a scanner. */
.capsule-scan {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, transparent, rgba(94, 242, 224, 0.16), transparent);
  transform: translateY(-100%);
  animation: scan 4.5s ease-in-out infinite;
}

.capsule-content {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
  padding: 2rem 1.6rem;
  min-height: 190px;
}

.capsule-tag {
  color: var(--aqua);
  font-size: 0.6rem;
}

.capsule-message {
  font-size: 0.86rem;
  line-height: 1.55;
  color: var(--ink-soft);
  display: -webkit-box;
  -webkit-line-clamp: 4;
  line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
  overflow-wrap: anywhere;
}

.capsule-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.62rem;
  overflow: hidden;
}

.capsule-meta .mono {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 42%;
  text-transform: none;
  letter-spacing: 0.06em;
}

.capsule-arrow {
  color: var(--amber);
}

.capsule-chips {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.chip {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 1px solid var(--line-bright);
  background: rgba(94, 242, 224, 0.14);
}

.chip-count {
  font-size: 0.56rem;
  text-transform: none;
}

.capsule-ring {
  position: absolute;
  left: 12%;
  right: 12%;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--line-bright), transparent);
}

.capsule-ring-top {
  top: 0;
}

.capsule-ring-bottom {
  bottom: 0;
}

.capsule-caption {
  font-size: 0.6rem;
  color: var(--ink-faint);
}

.is-sealed .capsule-caption {
  color: var(--amber);
}

@keyframes scan {
  0%,
  100% {
    transform: translateY(-100%);
  }
  50% {
    transform: translateY(100%);
  }
}
</style>
