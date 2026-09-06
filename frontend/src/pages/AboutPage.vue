<template>
  <main-layout>
    <div class="about-page">
      <div class="about-form-container">
        <h1 class="about-title">{{ t('about.title') }}</h1>
        <p class="about-description">{{ t('about.description') }}</p>

        <div class="spoiler">
          <button type="button" class="spoiler-button" @click="isOpen = !isOpen">
            <span>{{ isOpen ? t('about.spoiler.close') : t('about.spoiler.open') }}</span>
          </button>
          <transition name="spoiler-content">
            <!-- Safe: the copy is authored in this repo, no user input reaches it. -->
            <!-- eslint-disable-next-line vue/no-v-html -->
            <p v-if="isOpen" class="spoiler-content" v-html="t('about.spoiler.text')"></p>
          </transition>
        </div>

        <form class="about-form" @submit.prevent="handleSubmit">
          <label for="feedback-message" class="form-label">{{ t('about.form.label') }}</label>
          <textarea
            id="feedback-message"
            v-model="message"
            class="form-textarea"
            rows="5"
            maxlength="4096"
            :placeholder="t('about.form.placeholder')"
          ></textarea>
          <button type="submit" class="form-button" :disabled="isSubmitting">
            {{ t('about.form.submit') }}
          </button>
        </form>
      </div>
    </div>

    <div v-if="showModal" class="modal">
      <div class="modal-content">
        <h2 class="modal-title">{{ t('about.modal.title') }}</h2>
        <p class="modal-message">{{ t('about.modal.message') }}</p>
        <button type="button" class="modal-button" @click="goHome">
          {{ t('about.modal.button') }}
        </button>
      </div>
    </div>
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
.about-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding-top: 2.5rem;
}

.about-form-container {
  max-width: 640px;
  width: 100%;
  background: radial-gradient(
    circle,
    rgba(41, 123, 134, 0.9),
    rgba(2, 76, 92, 0.8),
    rgba(2, 76, 92, 0.9)
  );
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
  color: white;
  text-align: center;
}

.about-title {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.about-description {
  font-size: 1rem;
  margin-bottom: 2rem;
}

.spoiler {
  margin-top: -1.5rem;
  margin-bottom: 1rem;
  text-align: left;
  padding: 0.5rem;
  text-align: center;
}

.spoiler-button {
  background: none;
  border: none;
  color: #ffd700;
  font-size: 1rem;
  cursor: pointer;
  text-decoration: none;
  font-weight: bold;
  display: inline-flex;
  align-items: center;
  padding: 0.5rem;
  border-radius: 8px;
  transition: all 0.3s ease;
  position: relative;
}

.spoiler-button span {
  margin-left: 5px;
}

.spoiler-button::before {
  content: '';
  position: absolute;
  left: -5px;
  top: 30%;
  transform: translateY(-50%);
  width: 10px;
  height: 10px;
  border: solid 1px #ffd700;
  border-width: 2px 2px 0 0;
  transform: rotate(45deg);
  transition: transform 0.3s ease;
}

.spoiler-button:hover {
  color: #ffcc00;
  transform: scale(1.05);
}

.spoiler-button:hover::before {
  transform: rotate(135deg);
}

.spoiler-content {
  background: rgba(52, 123, 133, 0.7);
  color: #fff;
  padding: 1rem;
  border-radius: 10px;
  margin-top: 1rem;
  font-size: 1rem;
  line-height: 1.5;
  width: 100%;
  box-sizing: border-box;
}

.spoiler-content-enter-active,
.spoiler-content-leave-active {
  transition:
    opacity 0.5s ease,
    transform 0.5s ease;
}

.spoiler-content-enter,
.spoiler-content-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.about-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-label {
  font-size: 1rem;
  text-align: left;
}

.form-textarea {
  border: none;
  padding: 1rem;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  font-size: 1rem;
  resize: none;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.5);
}

.form-button {
  background: linear-gradient(145deg, #ffcc00, #ffd700, #ffcc00);
  color: #fff8dc;
  padding: 0.8rem 2rem;
  font-size: 1rem;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  font-weight: bold;
  text-transform: uppercase;
  transition:
    background 0.3s ease,
    transform 0.2s ease,
    box-shadow 0.3s ease;
  box-shadow: 0 4px 10px rgba(255, 215, 0, 0.3);
  background-size: 200% 200%;
  background-position: 100% 0;
}

.form-button:hover {
  background: linear-gradient(145deg, #ffd700, #ffcc00, #ffd700);
  transform: scale(1.05);
  box-shadow: 0 8px 20px rgba(255, 215, 0, 0.5);
  background-position: 0 0;
}

.form-button:active {
  background: linear-gradient(145deg, #b8860b, #ffd700, #b8860b);
  transform: scale(1);
  box-shadow: 0 4px 10px rgba(255, 215, 0, 0.3);
  background-position: 100% 0;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: #2c4e79;
  padding: 1rem;
  border-radius: 16px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
  color: white;
  text-align: center;
  width: 350px;
}

.modal-title {
  font-size: 1.25rem;
  margin-bottom: 1rem;
}

.modal-message {
  font-size: 1rem;
  margin-bottom: 1.5rem;
}

.modal-button {
  background: linear-gradient(145deg, #eecc0fe1, #daa520);
  color: #fff8dc;
  padding: 0.8rem 2rem;
  font-size: 1rem;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  font-weight: bold;
  text-transform: uppercase;
  transition:
    background 0.3s ease,
    transform 0.2s ease;
}

.modal-button:hover {
  background: linear-gradient(145deg, #daa520, #ffd700);
  transform: scale(1.05);
}

.modal-button:active {
  background: linear-gradient(145deg, #b8860b, #ffd700);
  transform: scale(1);
}
</style>
