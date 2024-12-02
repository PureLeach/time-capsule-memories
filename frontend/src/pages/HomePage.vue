<template>
  <main-layout>
    <section class="hero">
      <h1>{{ $t('welcome.title') }}</h1>
      <p>{{ $t('welcome.description') }}</p>
      <transition name="fade" mode="out-in">
        <transition-button :to="'/form'" class="welcome-button" key="button">
          {{ $t('welcome.buttonText') }}
        </transition-button>
      </transition>
    </section>
  </main-layout>
</template>

<script>
import MainLayout from '@/layouts/MainLayout.vue';
import TransitionButton from '@/components/TransitionButton.vue';

export default {
  components: { MainLayout, TransitionButton },
};
</script>

<style scoped>
.hero {
  display: flex;
  flex-direction: column;
  background-size: cover;
  justify-content: center;
  text-align: center;
  align-items: center;
  margin: 50px auto;
  padding: 20px;
  box-sizing: border-box;
  font-family: 'Arial', sans-serif;
  color: white;
  background: url('https://example.com/space-background.jpg') no-repeat center center fixed;
}

h1 {
  font-size: 3rem;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.8), 0 0 20px rgba(255, 255, 255, 0.6);
  margin-bottom: 20px;
}

p {
  font-size: 1.2rem;
  margin-bottom: 30px;
  opacity: 0.9;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.8);
}

.welcome-button {
  padding: 20px 40px;
  background: linear-gradient(45deg, #6a0dad, #4b0082, #0000ff);
  border: none;
  border-radius: 50px;
  color: white;
  font-size: 1.5rem;
  font-weight: bold;
  cursor: pointer;
  box-shadow: 0 0 25px rgba(255, 255, 255, 0.6), 0 0 40px rgba(0, 0, 255, 0.6);
  transition: transform 0.3s ease, box-shadow 0.3s ease, background 0.3s ease;
  position: relative;
  overflow: hidden;
  z-index: 1;
}

.welcome-button::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 300%;
  height: 300%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.2) 20%, rgba(0, 0, 255, 0.8) 80%);
  transition: width 0.4s, height 0.4s;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  z-index: -1;
}

.welcome-button:hover {
  transform: scale(1.1);
  box-shadow: 0 0 35px rgba(255, 255, 255, 0.8), 0 0 50px rgba(0, 0, 255, 0.8);
}

.welcome-button:hover::before {
  width: 500%;
  height: 500%;
}

.welcome-button:active {
  transform: scale(1);
  box-shadow: 0 0 25px rgba(255, 255, 255, 0.6), 0 0 35px rgba(0, 0, 255, 0.6);
}
</style>
