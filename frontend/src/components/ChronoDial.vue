<template>
  <div class="dial" aria-hidden="true">
    <div class="orbit orbit-1"></div>
    <div class="orbit orbit-2"></div>
    <div class="orbit orbit-3"></div>
    <div class="ticks">
      <span v-for="n in 60" :key="n" :style="tickStyle(n)"></span>
    </div>
    <div class="core">
      <div class="core-pulse"></div>
    </div>
    <div class="hand" :style="{ transform: `rotate(${secondAngle}deg)` }"></div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';

const secondAngle = ref(0);
const radius = ref(138);
let frame = 0;

// Driven off the wall clock rather than a counter, so the hand stays true even
// after the tab has been backgrounded.
function tick() {
  const now = new Date();
  const seconds = now.getSeconds() + now.getMilliseconds() / 1000;
  secondAngle.value = seconds * 6;
  frame = requestAnimationFrame(tick);
}

const tickStyle = (n) => ({
  transform: `translate(-50%, -50%) rotate(${n * 6}deg) translateY(-${radius.value}px)`,
  opacity: n % 5 === 0 ? 0.85 : 0.25,
  height: n % 5 === 0 ? '10px' : '5px',
});

function measure() {
  radius.value = window.innerWidth <= 720 ? 94 : 138;
}

onMounted(() => {
  measure();
  window.addEventListener('resize', measure);
  frame = requestAnimationFrame(tick);
});
onBeforeUnmount(() => {
  cancelAnimationFrame(frame);
  window.removeEventListener('resize', measure);
});
</script>

<style scoped>
.dial {
  position: relative;
  width: 320px;
  height: 320px;
  display: grid;
  place-items: center;
}

.orbit {
  position: absolute;
  border-radius: 50%;
  border: 1px solid var(--line);
}

.orbit-1 {
  inset: 0;
  border-top-color: var(--aqua);
  animation: spin 18s linear infinite;
}

.orbit-2 {
  inset: 32px;
  border-right-color: var(--violet);
  animation: spin 26s linear infinite reverse;
}

.orbit-3 {
  inset: 64px;
  border-bottom-color: var(--amber);
  animation: spin 34s linear infinite;
}

/* Each satellite rides its own ring. */
.orbit-1::after,
.orbit-2::after {
  content: '';
  position: absolute;
  top: -4px;
  left: 50%;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  transform: translateX(-50%);
}

.orbit-1::after {
  background: var(--aqua);
  box-shadow: var(--glow-aqua);
}

.orbit-2::after {
  background: var(--violet);
  box-shadow: var(--glow-violet);
}

.ticks span {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 1px;
  background: var(--aqua);
}

.core {
  position: relative;
  width: 86px;
  height: 86px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 30%, #fff, var(--aqua) 40%, transparent 72%);
  box-shadow:
    0 0 40px rgba(94, 242, 224, 0.55),
    0 0 90px rgba(155, 123, 255, 0.35);
}

.core-pulse {
  position: absolute;
  inset: -18px;
  border-radius: 50%;
  border: 1px solid var(--aqua);
  animation: pulse 3.2s ease-out infinite;
}

.hand {
  position: absolute;
  width: 1px;
  height: 130px;
  bottom: 50%;
  transform-origin: bottom center;
  background: linear-gradient(180deg, var(--amber), transparent);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes pulse {
  0% {
    transform: scale(0.85);
    opacity: 0.8;
  }
  100% {
    transform: scale(1.9);
    opacity: 0;
  }
}

@media (max-width: 720px) {
  .dial {
    width: 220px;
    height: 220px;
  }
  .core {
    width: 62px;
    height: 62px;
  }
  .hand {
    height: 88px;
  }
}
</style>
