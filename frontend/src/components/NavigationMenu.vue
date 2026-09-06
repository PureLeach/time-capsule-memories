<template>
  <nav class="nav">
    <router-link v-for="link in links" :key="link.to" :to="link.to" class="nav-link">
      <span class="nav-index mono">{{ link.index }}</span>
      <span class="nav-label">{{ t(link.label) }}</span>
    </router-link>
  </nav>
</template>

<script setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const links = [
  { to: '/', label: 'menu.home', index: '01' },
  { to: '/form', label: 'menu.create', index: '02' },
  { to: '/about', label: 'menu.about', index: '03' },
];
</script>

<style scoped>
.nav {
  display: flex;
  gap: 0.35rem;
}

.nav-link {
  position: relative;
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  padding: 0.5rem 0.85rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--ink-soft);
  transition:
    color 0.25s ease,
    background 0.25s ease;
}

.nav-link:hover {
  color: var(--ink);
  background: rgba(255, 255, 255, 0.04);
}

.nav-index {
  font-size: 0.56rem;
  color: var(--ink-faint);
  transition: color 0.25s ease;
}

.nav-label {
  font-size: 0.82rem;
  letter-spacing: 0.06em;
}

.nav-link.router-link-exact-active {
  color: var(--ink);
}

.nav-link.router-link-exact-active .nav-index {
  color: var(--aqua);
}

/* Active tab is marked by a lit rail rather than a filled pill. */
.nav-link.router-link-exact-active::after {
  content: '';
  position: absolute;
  left: 0.85rem;
  right: 0.85rem;
  bottom: 2px;
  height: 1px;
  background: var(--aqua);
  box-shadow: var(--glow-aqua);
}

@media (max-width: 640px) {
  .nav-index {
    display: none;
  }
  .nav-link {
    padding: 0.5rem 0.55rem;
  }
}
</style>
