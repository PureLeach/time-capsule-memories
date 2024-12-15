<template>
  <main-layout>
    <div class="about-page">
      <div class="about-form-container">
        <h1 class="about-title">{{ $t('about.title') }}</h1>
        <p class="about-description">{{ $t('about.description') }}</p>

        <!-- Спойлер -->
        <div class="spoiler">
          <button @click="toggleSpoiler" class="spoiler-button">
            <span>{{ isOpen ? $t('about.spoiler.close') : $t('about.spoiler.open') }}</span>
          </button>
          <transition name="fade" @before-enter="beforeEnter" @enter="enter" @leave="leave">
            <div v-if="isOpen" class="spoiler-content">
              <!-- Используем v-html для рендера HTML -->
              <p v-html="$t('about.spoiler.text')"></p>
            </div>
          </transition>
        </div>

        <form class="about-form">
          <label for="message" class="form-label">{{ $t('about.form.label') }}</label>
          <textarea id="message" class="form-textarea" rows="5" :placeholder="$t('about.form.placeholder')"></textarea>
          <button type="submit" class="form-button">{{ $t('about.form.submit') }}</button>
        </form>
      </div>
    </div>
  </main-layout>
</template>

<script>
import MainLayout from '@/layouts/MainLayout.vue';

export default {
  name: 'AboutPage',
  components: {
    MainLayout
  },
  data() {
    return {
      isOpen: false // Состояние спойлера
    };
  },
  methods: {
    toggleSpoiler() {
      this.isOpen = !this.isOpen; // Переключение состояния спойлера
    },
    beforeEnter(el) {
      el.style.opacity = 0;
      el.style.transform = 'translateY(-10px)';
    },
    enter(el, done) {
      el.offsetHeight; // Форсируем перерисовку
      el.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
      el.style.opacity = 1;
      el.style.transform = 'translateY(0)';
      done();
    },
    leave(el, done) {
      el.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
      el.style.opacity = 0;
      el.style.transform = 'translateY(-10px)';
      done();
    }
  }
};
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
  background: radial-gradient(circle, rgba(25, 25, 112, 0.9), rgba(0, 0, 51, 0.8), rgba(47, 79, 79, 0.9));
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
  left: -5px; /* Сдвигает стрелку влево */
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
  background: rgba(0, 0, 51, 0.7);
  color: #fff;
  padding: 1rem;  
  border-radius: 10px;
  margin-top: 1rem;
  font-size: 1rem;
  line-height: 1.5;
  width: 100%; /* Убедитесь, что контент спойлера будет растягиваться на всю доступную ширину */
  box-sizing: border-box; /* Чтобы паддинги не нарушали размер */
}

.spoiler-content-enter-active,
.spoiler-content-leave-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.spoiler-content-enter, .spoiler-content-leave-to {
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
  background: linear-gradient(145deg, #ffd700, #daa520);
  color: #fff8dc;
  padding: 0.8rem 2rem;
  font-size: 1rem;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  font-weight: bold;
  text-transform: uppercase;
  transition: background 0.3s ease, transform 0.2s ease;
}

.form-button:hover {
  background: linear-gradient(145deg, #daa520, #ffd700);
  transform: scale(1.05);
}

.form-button:active {
  background: linear-gradient(145deg, #b8860b, #ffd700);
  transform: scale(1);
}
</style>
