<template>
  <div class="falling-stars" v-if="active">
    <div v-for="star in stars" :key="star.id" class="star" :style="star.style"></div>
  </div>
</template>

<script>
export default {
  name: 'FallingStars',
  data() {
    return {
      active: false,
      stars: [],
    };
  },
  methods: {
    trigger() {
      this.active = true;
      this.generateStars();
      setTimeout(() => {
        this.active = false;
      }, 1500);
    },
    generateStars() {
      this.stars = Array.from({ length: 50 }).map(() => ({
        id: Math.random(),
        style: {
          left: `${Math.random() * 100}vw`,
          top: `-10px`,  // Начинаем выше экрана, чтобы избежать замедления в начале
          width: `${Math.random() * 2 + 3}px`,  // Звезды будут от 3px до 5px в ширину
          height: `${Math.random() * 3 + 4}px`, // Высота от 10px до 15px для более вытянутого вида
          opacity: Math.random() * 0.8 + 0.5,  // Прозрачность от 0.5 до 1
          boxShadow: `0 0 5px rgba(255, 255, 255, 0.8)`,  // Лёгкая тень
          animationDelay: `${Math.random() * 1}s`, // Задержка для случайности
        },
      }));
    },
  },
};
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
  background: white;
  border-radius: 50%;
  animation: fall 2s linear infinite;
  transform-origin: center;
  animation-fill-mode: forwards;
  /* Ожидание анимации без зависания на старте */
}

@keyframes fall {
  0% {
    transform: translate(0, -60%) rotate(180deg);
    /* Начинаем с верхнего края, поворачиваем под 45 градусов */
  }

  100% {
    transform: translate(100vw, 100vh) rotate(160deg);
    /* Падаем под 45 градусов до правого нижнего угла */
    opacity: 0;
    /* Прозрачность исчезает */
  }
}
</style>
