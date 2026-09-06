<template>
  <warp-overlay ref="warp" @finished="navigate" />
  <button type="button" class="launch" @click="handleClick">
    <span class="launch-glow"></span>
    <span class="launch-label"><slot /></span>
    <span class="launch-arrow">→</span>
  </button>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import WarpOverlay from './WarpOverlay.vue';

const props = defineProps({
  to: { type: String, required: true },
});

const router = useRouter();
const warp = ref(null);

const handleClick = () => warp.value?.trigger();
const navigate = () => router.push(props.to);
</script>

<style scoped>
.launch {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 0.9rem;
  padding: 1rem 2.2rem;
  border: 1px solid var(--line-bright);
  border-radius: 999px;
  background: rgba(94, 242, 224, 0.06);
  color: var(--ink);
  font-family: var(--font-mono);
  font-size: 0.78rem;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  cursor: pointer;
  overflow: hidden;
  transition:
    border-color 0.3s ease,
    box-shadow 0.3s ease,
    transform 0.3s ease;
}

.launch:hover {
  transform: translateY(-2px);
  border-color: var(--aqua);
  box-shadow:
    0 0 30px rgba(94, 242, 224, 0.3),
    inset 0 0 24px rgba(94, 242, 224, 0.12);
}

.launch:active {
  transform: translateY(0);
}

/* A light sweeps across the button on hover. */
.launch-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(105deg, transparent 35%, rgba(255, 255, 255, 0.28), transparent 65%);
  transform: translateX(-120%);
  transition: transform 0.7s ease;
}

.launch:hover .launch-glow {
  transform: translateX(120%);
}

.launch-label,
.launch-arrow {
  position: relative;
}

.launch-arrow {
  color: var(--aqua);
  transition: transform 0.3s ease;
}

.launch:hover .launch-arrow {
  transform: translateX(5px);
}
</style>
