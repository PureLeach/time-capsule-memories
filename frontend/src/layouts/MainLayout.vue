<template>
  <div class="main-layout">
    <header class="main-header">
      <navigation-menu />
      <language-switcher />
    </header>
    <main class="main-content">
      <slot />
    </main>
    <footer class="main-footer">
      <p>&copy; 2024 Time Capsule of Memories</p>
    </footer>
  </div>
</template>

<script>
import NavigationMenu from '@/components/NavigationMenu.vue';
import LanguageSwitcher from '@/components/LanguageSwitcher.vue';

export default {
  components: { NavigationMenu, LanguageSwitcher },
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
  padding: 0.5rem;
  background: linear-gradient(135deg, rgba(19, 121, 104, 0.8), rgba(128, 83, 89, 0.664), rgba(25, 25, 112, 0.8));
  color: white;
  border-bottom: 1px solid rgba(221, 221, 221, 0.5);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  font-family: 'Arial', sans-serif;
  font-size: 1.1rem;
  letter-spacing: 0.5px;
  padding: 1.6rem;
}

.main-content {
  flex: 1;
  padding: 2rem;
}

.main-footer {
  text-align: center;
  padding: 0.7rem;
  background: linear-gradient(135deg, rgba(12, 14, 29, 0.8), rgba(31, 38, 51, 0.8), rgba(41, 44, 67, 0.8), rgba(74, 42, 122, 0.8));
  color: white;
  border-top: 1px solid rgba(221, 221, 221, 0.5);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  font-family: 'Arial', sans-serif;
  font-size: 1rem;
  letter-spacing: 0.5px;
}
</style>
