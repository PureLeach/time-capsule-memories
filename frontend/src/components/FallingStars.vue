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
      }, 3000); // Увеличена длительность анимации
    },
    generateStars() {
      this.stars = Array.from({ length: 70 }).map(() => ({
        id: Math.random(),
        style: {
          left: `${Math.random() * 100}vw`,
          top: `-10px`,
          width: `${Math.random() * 3 + 3}px`, // Минимум 4px, максимум 8px
          height: `${Math.random() * 4 + 5}px`, // Высота звезд чуть больше
          opacity: Math.random() * 0.6 + 0.4, // Прозрачность от 0.4 до 1
          boxShadow: `0 0 5px rgba(255, 255, 255, 1)`,  // Яркое свечение
          backgroundColor: `hsl(${Math.random() * 360}, 100%, 85%)`, // Разноцветные звезды
          animationDelay: `${Math.random() * 1}s`, // Случайная задержка
          animationDuration: `${Math.random() * 2 + 1.5}s`, // Случайная длительность падения
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
  border-radius: 50%;
  animation: fall 2s linear infinite, twinkle 3s ease-in-out infinite;
  animation-fill-mode: forwards;
  transform-origin: center;
  box-shadow: 0 0 10px rgba(255, 255, 255, 1); /* Яркое свечение */
}

@keyframes fall {
  0% {
    transform: translate(0, -60%) rotate(180deg);
    /* Начинаем с верхнего края, с небольшим вращением */
  }

  50% {
    transform: translate(50vw, 50vh) rotate(90deg);
    /* Звезда будет двигаться по траектории с небольшими колебаниями */
  }

  80% {
    transform: translate(80vw, 80vh) rotate(140deg);
    opacity: 0.6; /* Прозрачность немного снижается */
  }

  100% {
    transform: translate(100vw, 100vh) rotate(160deg);
    opacity: 0; /* Звезда исчезает */
    transform: scale(0); /* Звезда уменьшается */
  }
}

@keyframes twinkle {
  0%, 100% {
    opacity: 0.8; /* Стандартная прозрачность */
  }

  50% {
    opacity: 1; /* Звезды слегка мигают */
  }
}
</style>
