<template>
  <falling-stars ref="stars" />
  <el-button type="primary" :to="to" @click="handleClick" class="transition-button">
    <slot />
  </el-button>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import FallingStars from './FallingStars.vue';

export default {
  props: {
    to: {
      type: String,
      required: true,
    },
  },
  components: { FallingStars },
  setup(props) {
    const router = useRouter();
    const stars = ref(null);

    const handleClick = () => {
      if (stars.value) stars.value.trigger();
      setTimeout(() => router.push(props.to), 1500);
    };

    return { handleClick, stars };
  },
};
</script>

<style scoped>
/* Основной стиль для кнопки */
.transition-button {
  position: relative;
  padding: 20px 40px;
  /* Увеличен размер кнопки */
  border-radius: 25px;
  font-size: 18px; /* Увеличен шрифт */
  font-weight: bold;
  color: #fff;
  background: linear-gradient(135deg, #1d2a6c, #0d1b2a); /* Тёмно-синий градиент */
  border: 2px solid #ffffff;
  text-transform: uppercase;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
  animation: gradientAnimation 10s ease-in-out infinite; /* Увеличили продолжительность анимации */
}

/* Эффект неонового свечения */
.transition-button::before {
  content: '';
  position: absolute;
  top: -5px;
  left: -5px;
  right: -5px;
  bottom: -5px;
  background: linear-gradient(45deg, #ff007f, #00aaff, #ff007f);
  z-index: -1;
  filter: blur(10px);
  opacity: 0.8;
  animation: pulse 1s infinite alternate;
}

/* Плавный эффект при наведении */
.transition-button:hover {
  background: #333;
  box-shadow: 0 0 20px rgba(255, 255, 255, 1);
}

/* Анимация пульсации */
@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 0.7;
  }
  100% {
    transform: scale(1.2);
    opacity: 1;
  }
}

/* Анимация градиента по кругу (медленный эффект переливания) */
@keyframes gradientAnimation {
  0% {
    background: linear-gradient(135deg, #1d2a6c, #0d1b2a);
  }
  25% {
    background: linear-gradient(135deg, #374c91, #1b2c53);
  }
  50% {
    background: linear-gradient(135deg, #1a355b, #2e3a77);
  }
  75% {
    background: linear-gradient(135deg, #4f6796, #1d2a6c);
  }
  100% {
    background: linear-gradient(135deg, #1d2a6c, #0d1b2a);
  }
}
</style>
