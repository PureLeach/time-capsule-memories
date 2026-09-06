<template>
  <canvas ref="canvas" class="cosmic" aria-hidden="true"></canvas>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';

const LAYERS = [
  { count: 110, size: [0.4, 1.0], speed: 0.008, depth: 6, alpha: [0.2, 0.5] },
  { count: 70, size: [0.8, 1.7], speed: 0.018, depth: 16, alpha: [0.35, 0.75] },
  { count: 28, size: [1.4, 2.6], speed: 0.032, depth: 34, alpha: [0.55, 1] },
];
const TINTS = ['#eaf2ff', '#9b7bff', '#5ef2e0', '#ffc46b'];

const canvas = ref(null);

let ctx;
let frame = 0;
let stars = [];
let comet = null;
let nextComet = 3000;
let width = 0;
let height = 0;
let dpr = 1;
// Pointer position eased toward the cursor, so the parallax glides instead of
// snapping to every mouse event.
const pointer = { x: 0, y: 0, tx: 0, ty: 0 };

const rand = (min, max) => min + Math.random() * (max - min);
const reducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches;

function seed() {
  stars = LAYERS.flatMap((layer, index) =>
    Array.from({ length: layer.count }, () => ({
      x: Math.random() * width,
      y: Math.random() * height,
      r: rand(...layer.size),
      alpha: rand(...layer.alpha),
      twinkle: rand(0.4, 1.6),
      phase: Math.random() * Math.PI * 2,
      tint: TINTS[Math.floor(Math.random() * TINTS.length)],
      layer: index,
    }))
  );
}

function resize() {
  dpr = Math.min(window.devicePixelRatio || 1, 2);
  width = window.innerWidth;
  height = window.innerHeight;
  canvas.value.width = width * dpr;
  canvas.value.height = height * dpr;
  canvas.value.style.width = `${width}px`;
  canvas.value.style.height = `${height}px`;
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  seed();
}

function spawnComet() {
  const fromLeft = Math.random() > 0.5;
  comet = {
    x: fromLeft ? -80 : width + 80,
    y: rand(0, height * 0.55),
    vx: (fromLeft ? 1 : -1) * rand(6, 10),
    vy: rand(1.5, 3.2),
    life: 1,
  };
}

function draw(time) {
  ctx.clearRect(0, 0, width, height);

  pointer.x += (pointer.tx - pointer.x) * 0.04;
  pointer.y += (pointer.ty - pointer.y) * 0.04;

  for (const star of stars) {
    const layer = LAYERS[star.layer];
    star.y += layer.speed;
    if (star.y > height + 2) {
      star.y = -2;
      star.x = Math.random() * width;
    }

    const flicker = 0.75 + 0.25 * Math.sin(time * 0.001 * star.twinkle + star.phase);
    ctx.globalAlpha = star.alpha * flicker;
    ctx.fillStyle = star.tint;
    ctx.beginPath();
    ctx.arc(
      star.x + pointer.x * layer.depth,
      star.y + pointer.y * layer.depth,
      star.r,
      0,
      Math.PI * 2
    );
    ctx.fill();
  }

  if (comet) {
    ctx.globalAlpha = comet.life;
    const gradient = ctx.createLinearGradient(
      comet.x,
      comet.y,
      comet.x - comet.vx * 18,
      comet.y - comet.vy * 18
    );
    gradient.addColorStop(0, 'rgba(255,255,255,0.95)');
    gradient.addColorStop(1, 'rgba(94,242,224,0)');
    ctx.strokeStyle = gradient;
    ctx.lineWidth = 1.6;
    ctx.beginPath();
    ctx.moveTo(comet.x, comet.y);
    ctx.lineTo(comet.x - comet.vx * 18, comet.y - comet.vy * 18);
    ctx.stroke();

    comet.x += comet.vx;
    comet.y += comet.vy;
    comet.life -= 0.006;
    if (comet.life <= 0 || comet.x < -200 || comet.x > width + 200) comet = null;
  } else if (time > nextComet) {
    spawnComet();
    nextComet = time + rand(7000, 18000);
  }

  ctx.globalAlpha = 1;
  frame = requestAnimationFrame(draw);
}

function onPointerMove(event) {
  pointer.tx = event.clientX / window.innerWidth - 0.5;
  pointer.ty = event.clientY / window.innerHeight - 0.5;
}

onMounted(() => {
  ctx = canvas.value.getContext('2d');
  resize();
  window.addEventListener('resize', resize);

  if (reducedMotion()) {
    draw(0);
    cancelAnimationFrame(frame);
    return;
  }

  window.addEventListener('pointermove', onPointerMove, { passive: true });
  frame = requestAnimationFrame(draw);
});

onBeforeUnmount(() => {
  cancelAnimationFrame(frame);
  window.removeEventListener('resize', resize);
  window.removeEventListener('pointermove', onPointerMove);
});
</script>

<style scoped>
.cosmic {
  position: fixed;
  inset: 0;
  z-index: -2;
  pointer-events: none;
}
</style>
