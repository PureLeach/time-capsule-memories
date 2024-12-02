<template>
  <div class="main-layout">
    <header class="main-header">
      <welcome-button />
      <language-switcher />
    </header>
    <main class="main-content">
      <slot />
    </main>
    <footer class="main-footer">
      <p>&copy; 2024 Time Capsule Memories</p>
    </footer>
  </div>
</template>

<script>
import WelcomeButton from '@/components/WelcomeButton.vue';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';

export default {
  components: { WelcomeButton, LanguageSwitcher },
};
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-header {
  display: flex;
  justify-content: space-between;
  padding: 1rem;
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

.main-content {
  flex: 1;
  padding: 2rem;
}

.main-footer {
  text-align: center;
  padding: 1rem;
  background-color: #f5f5f5;
  border-top: 1px solid #ddd;
}
</style>
