<template>
  <div v-show="active" class="falling-stars" aria-hidden="true">
    <div v-for="star in stars" :key="star.id" class="star" :style="star.style"></div>
  </div>
</template>

<script setup>
import { onBeforeUnmount, ref } from 'vue';

const props = defineProps({
  starCount: { type: Number, default: 70 },
  durationMs: { type: Number, default: 3000 },
});

const emit = defineEmits(['finished']);

const active = ref(false);
const stars = ref([]);
let timer = null;
let nextStarId = 0;

function randomStar() {
  return {
    id: nextStarId++,
    style: {
      left: `${Math.random() * 100}vw`,
      top: '-10px',
      width: `${Math.random() * 3 + 3}px`,
      height: `${Math.random() * 4 + 5}px`,
      opacity: Math.random() * 0.6 + 0.4,
      backgroundColor: `hsl(${Math.random() * 360}, 100%, 85%)`,
      animationDelay: `${Math.random()}s`,
      animationDuration: `${Math.random() * 2 + 1.5}s`,
    },
  };
}

function trigger() {
  clearTimeout(timer);
  stars.value = Array.from({ length: props.starCount }, randomStar);
  active.value = true;

  timer = setTimeout(() => {
    active.value = false;
    stars.value = [];
    emit('finished');
  }, props.durationMs);
}

// The timer outlives the component if the user navigates away mid-animation.
onBeforeUnmount(() => clearTimeout(timer));

defineExpose({ trigger });
</script>

<style scoped>
.falling-stars {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 9999;
}

.star {
  position: absolute;
  border-radius: 50%;
  animation:
    fall 2s linear infinite,
    twinkle 3s ease-in-out infinite;
  animation-fill-mode: forwards;
  transform-origin: center;
  box-shadow: 0 0 10px rgba(255, 255, 255, 1);
}

@keyframes fall {
  0% {
    transform: translate(0, -60%) rotate(180deg);
  }

  50% {
    transform: translate(50vw, 50vh) rotate(90deg);
  }

  80% {
    transform: translate(80vw, 80vh) rotate(140deg);
    opacity: 0.6;
  }

  100% {
    transform: translate(100vw, 100vh) rotate(160deg) scale(0);
    opacity: 0;
  }
}

@keyframes twinkle {
  0%,
  100% {
    opacity: 0.8;
  }

  50% {
    opacity: 1;
  }
}
</style>
