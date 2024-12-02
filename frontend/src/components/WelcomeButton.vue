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
  background: linear-gradient(90deg, rgba(255, 105, 180, 1) 0%, rgba(255, 159, 159, 1) 100%);
  color: white;
  padding: 12px 28px;
  font-size: 1.3rem;
  border: none;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.4s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.btn:hover {
  background: linear-gradient(90deg, rgba(255, 159, 159, 1) 0%, rgba(255, 105, 180, 1) 100%);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  transform: scale(1.05);
}

.btn:active {
  background: linear-gradient(90deg, rgba(255, 105, 180, 1) 0%, rgba(255, 159, 159, 1) 100%);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  transform: scale(1);
}
</style>
