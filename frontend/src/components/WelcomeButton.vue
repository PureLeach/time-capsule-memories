<template>
  <router-link to="/" class="welcome-button">
    <button class="btn">{{ $t('welcome.title') }}</button>
  </router-link>
</template>

<script>
export default {
  name: 'WelcomeButton',
};
</script>

<style scoped>
.welcome-button {
  text-decoration: none;
}

.btn {
  background: linear-gradient(45deg, rgba(255, 182, 193, 1) 0%, rgba(255, 144, 128, 1) 100%);
  color: white;
  padding: 15px 30px;
  font-size: 1.3rem;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  font-weight: bold;
  text-transform: uppercase;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  outline: none;
  transform: scale(1);
}

.btn::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 300%;
  height: 300%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, rgba(255, 255, 255, 0) 70%);
  animation: pulse 2s infinite ease-in-out;
  transform: translate(-50%, -50%);
}

.btn:hover {
  background: linear-gradient(45deg, rgba(255, 144, 128, 1) 0%, rgba(255, 182, 193, 1) 100%);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3), 0 0 20px rgba(255, 182, 193, 1), 0 0 40px rgba(255, 182, 193, 1);
  transform: scale(1.05); /* Уменьшенный наклон */
  color: #fff;
  text-shadow: 0 0 10px rgba(255, 182, 193, 1), 0 0 20px rgba(255, 182, 193, 1);
}

.btn:active {
  background: linear-gradient(45deg, rgba(255, 182, 193, 1) 0%, rgba(255, 144, 128, 1) 100%);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: scale(1); /* Без наклона */
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 0.6;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.3;
  }
  100% {
    transform: scale(1);
    opacity: 0.6;
  }
}
</style>
