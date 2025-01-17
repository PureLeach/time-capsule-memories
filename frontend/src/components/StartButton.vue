<template>
  <falling-stars ref="stars" />
  <el-button type="primary" :to="to" @click="handleClick" class="start-button">
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
.start-button {
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
  width: 150px; /* Фиксированная ширина */
}

/* Эффект неонового свечения */
.start-button::before {
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
}

/* Плавный эффект при наведении */
.start-button:hover {
  background: #333;
  box-shadow: 0 0 20px rgba(255, 255, 255, 1);
}

</style>
