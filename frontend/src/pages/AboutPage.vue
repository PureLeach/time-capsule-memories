<template>
  <main-layout>
    <section class="about">
      <header class="about-head">
        <span class="eyebrow">{{ t('about.eyebrow') }}</span>
        <h1 class="about-title">{{ t('about.title') }}</h1>
        <p class="about-lead">{{ t('about.description') }}</p>
      </header>

      <div class="about-grid">
        <article class="panel card">
          <button
            type="button"
            class="disclosure"
            :aria-expanded="isOpen"
            @click="isOpen = !isOpen"
          >
            <span class="mono">{{
              isOpen ? t('about.spoiler.close') : t('about.spoiler.open')
            }}</span>
            <span class="disclosure-icon" :class="{ 'is-open': isOpen }">+</span>
          </button>
          <transition name="reveal">
            <!-- Safe: the copy is authored in this repo, no user input reaches it. -->
            <!-- eslint-disable-next-line vue/no-v-html -->
            <div v-if="isOpen" class="disclosure-body" v-html="t('about.spoiler.text')"></div>
          </transition>

          <ol class="steps">
            <li v-for="step in steps" :key="step" class="step">
              <span class="mono step-index">{{ String(step).padStart(2, '0') }}</span>
              <span class="step-text">{{ t(`about.steps.${step}`) }}</span>
            </li>
          </ol>
        </article>

        <form class="panel card feedback" @submit.prevent="handleSubmit">
          <span class="eyebrow">{{ t('about.form.eyebrow') }}</span>
          <label for="feedback-message" class="mono feedback-label">
            {{ t('about.form.label') }}
          </label>
          <textarea
            id="feedback-message"
            v-model="message"
            class="feedback-input"
            rows="6"
            maxlength="4096"
            :placeholder="t('about.form.placeholder')"
          ></textarea>
          <div class="feedback-foot">
            <span class="mono">{{ message.length }} / 4096</span>
            <button type="submit" class="btn" :disabled="isSubmitting">
              {{ t('about.form.submit') }}
            </button>
          </div>
        </form>
      </div>
    </section>

    <transition name="reveal">
      <div v-if="showModal" class="dialog-mask" @click.self="goHome">
        <div class="panel dialog">
          <span class="eyebrow">{{ t('about.modal.tag') }}</span>
          <h2 class="dialog-title">{{ t('about.modal.title') }}</h2>
          <p class="dialog-text">{{ t('about.modal.message') }}</p>
          <button type="button" class="btn" @click="goHome">{{ t('about.modal.button') }}</button>
        </div>
      </div>
    </transition>
  </main-layout>
</template>

<script setup>
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import MainLayout from '@/layouts/MainLayout.vue';
import { submitFeedback } from '@/api/feedback';

const { t } = useI18n();
const router = useRouter();

const steps = [1, 2, 3, 4];
const isOpen = ref(false);
const message = ref('');
const isSubmitting = ref(false);
const showModal = ref(false);

async function handleSubmit() {
  if (!message.value.trim()) {
    ElMessage.warning(t('about.form.emptyMessage'));
    return;
  }

  isSubmitting.value = true;
  try {
    await submitFeedback({ message: message.value });
    message.value = '';
    showModal.value = true;
  } catch (error) {
    ElMessage.error(error.message);
  } finally {
    isSubmitting.value = false;
  }
}

function goHome() {
  showModal.value = false;
  router.push('/');
}
</script>

<style scoped>
.about {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.about-head {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
  max-width: 60ch;
}

.about-title {
  font-size: clamp(1.8rem, 3.4vw, 2.6rem);
  font-weight: 300;
}

.about-lead {
  font-size: 0.98rem;
  line-height: 1.65;
  color: var(--ink-soft);
}

.about-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
  align-items: start;
}

.card {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  padding: clamp(1.4rem, 3vw, 2rem);
}

.disclosure {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.03);
  color: var(--aqua);
  cursor: pointer;
  transition: border-color 0.25s ease;
}

.disclosure:hover {
  border-color: var(--line-bright);
}

.disclosure-icon {
  font-size: 1.1rem;
  color: var(--aqua);
  transition: transform 0.3s ease;
}

.disclosure-icon.is-open {
  transform: rotate(45deg);
}

.disclosure-body {
  font-size: 0.88rem;
  line-height: 1.7;
  color: var(--ink-soft);
  white-space: pre-line;
}

.disclosure-body :deep(a) {
  color: var(--amber);
  text-decoration: none;
  border-bottom: 1px solid rgba(255, 196, 107, 0.4);
}

.steps {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  list-style: none;
}

.step {
  display: flex;
  gap: 0.9rem;
  align-items: baseline;
  padding-left: 0.9rem;
  border-left: 1px solid var(--line);
}

.step-index {
  color: var(--aqua);
  font-size: 0.6rem;
}

.step-text {
  font-size: 0.86rem;
  line-height: 1.55;
  color: var(--ink-soft);
}

.feedback-label {
  font-size: 0.6rem;
}

.feedback-input {
  width: 100%;
  padding: 0.9rem 1rem;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.03);
  color: var(--ink);
  font-family: inherit;
  font-size: 0.9rem;
  line-height: 1.6;
  resize: none;
  transition:
    border-color 0.25s ease,
    box-shadow 0.25s ease;
}

.feedback-input::placeholder {
  color: var(--ink-faint);
}

.feedback-input:focus {
  outline: none;
  border-color: var(--aqua);
  box-shadow: 0 0 20px rgba(94, 242, 224, 0.18);
}

.feedback-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.btn {
  padding: 0.7rem 1.6rem;
  border: 1px solid var(--aqua);
  border-radius: 999px;
  background: rgba(94, 242, 224, 0.12);
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 0.66rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 0.25s ease;
}

.btn:hover:not(:disabled) {
  background: rgba(94, 242, 224, 0.2);
  box-shadow: 0 0 26px rgba(94, 242, 224, 0.3);
}

.btn:disabled {
  opacity: 0.6;
  cursor: progress;
}

.dialog-mask {
  position: fixed;
  inset: 0;
  z-index: 2200;
  display: grid;
  place-items: center;
  padding: 1.5rem;
  background: rgba(4, 5, 12, 0.8);
  backdrop-filter: blur(8px);
}

.dialog {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  width: min(400px, 100%);
  padding: 2.2rem 2rem;
  text-align: center;
}

.dialog-title {
  font-size: 1.3rem;
  font-weight: 300;
}

.dialog-text {
  font-size: 0.9rem;
  color: var(--ink-soft);
}

.reveal-enter-active,
.reveal-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.reveal-enter-from,
.reveal-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
