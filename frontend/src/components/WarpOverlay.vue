<template>
  <div v-show="active" class="warp" aria-hidden="true">
    <canvas ref="canvas"></canvas>
    <div class="warp-flash" :class="{ 'is-lit': flash }"></div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue';

const props = defineProps({
  durationMs: { type: Number, default: 1500 },
  streaks: { type: Number, default: 260 },
});

const emit = defineEmits(['finished']);

const canvas = ref(null);
const active = ref(false);
const flash = ref(false);

let ctx;
let frame = 0;
let timer = null;
let particles = [];
let startedAt = 0;

// Streaks radiate from the centre, so each one is polar: an angle it keeps and a
// radius that accelerates outward.
function seed(width, height) {
  const reach = Math.hypot(width, height) / 2;
  particles = Array.from({ length: props.streaks }, () => ({
    angle: Math.random() * Math.PI * 2,
    radius: Math.random() * reach * 0.25,
    speed: 2 + Math.random() * 9,
    length: 20 + Math.random() * 120,
    hue: [180, 265, 40][Math.floor(Math.random() * 3)],
    reach,
  }));
}

function draw(time) {
  const { width, height } = canvas.value;
  const dpr = Math.min(window.devicePixelRatio || 1, 2);
  const cx = width / dpr / 2;
  const cy = height / dpr / 2;
  const progress = Math.min((time - startedAt) / props.durationMs, 1);
  // Ease in, so the jump builds instead of starting at full speed.
  const thrust = progress < 0.75 ? progress / 0.75 : 1;

  ctx.clearRect(0, 0, width, height);
  ctx.lineCap = 'round';

  for (const p of particles) {
    p.radius += p.speed * (0.4 + thrust * 2.6);
    if (p.radius > p.reach) {
      p.radius = Math.random() * 40;
      p.angle = Math.random() * Math.PI * 2;
    }

    const tail = Math.max(p.radius - p.length * thrust, 0);
    const fade = Math.min(p.radius / p.reach, 1);

    ctx.strokeStyle = `hsla(${p.hue}, 100%, ${70 + fade * 25}%, ${0.15 + fade * 0.75})`;
    ctx.lineWidth = 0.6 + fade * 2;
    ctx.beginPath();
    ctx.moveTo(cx + Math.cos(p.angle) * tail, cy + Math.sin(p.angle) * tail);
    ctx.lineTo(cx + Math.cos(p.angle) * p.radius, cy + Math.sin(p.angle) * p.radius);
    ctx.stroke();
  }

  frame = requestAnimationFrame(draw);
}

function stop() {
  cancelAnimationFrame(frame);
  clearTimeout(timer);
  active.value = false;
  flash.value = false;
}

function trigger() {
  stop();

  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    emit('finished');
    return;
  }

  active.value = true;
  requestAnimationFrame(() => {
    const el = canvas.value;
    if (!el) return;

    const dpr = Math.min(window.devicePixelRatio || 1, 2);
    el.width = window.innerWidth * dpr;
    el.height = window.innerHeight * dpr;
    ctx = el.getContext('2d');
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    seed(window.innerWidth, window.innerHeight);
    startedAt = performance.now();
    frame = requestAnimationFrame(draw);

    // The flash lands at the end of the run, covering the route swap.
    timer = setTimeout(() => (flash.value = true), props.durationMs - 320);
    setTimeout(() => {
      stop();
      emit('finished');
    }, props.durationMs);
  });
}

onBeforeUnmount(stop);

defineExpose({ trigger });
</script>

<style scoped>
.warp {
  position: fixed;
  inset: 0;
  z-index: 9999;
  pointer-events: none;
}

.warp canvas {
  width: 100%;
  height: 100%;
}

.warp-flash {
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center, rgba(255, 255, 255, 0.95), rgba(94, 242, 224, 0));
  opacity: 0;
  transition: opacity 0.32s ease-in;
}

.warp-flash.is-lit {
  opacity: 1;
}
</style>
