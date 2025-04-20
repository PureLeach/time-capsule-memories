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
    delay: {
      type: Number,
      default: 1500,  // Default delay time
    },
  },
  components: { FallingStars },
  setup(props) {
    const router = useRouter();
    const stars = ref(null);

    const triggerStars = () => {
      if (stars.value) stars.value.trigger();
    };

    const navigate = () => {
      setTimeout(() => router.push(props.to), props.delay);
    };

    const handleClick = () => {
      triggerStars();
      navigate();
    };

    return { handleClick, stars };
  },
};
</script>

<style scoped>
/* Main style for the button */
.start-button {
  position: relative;
  padding: 20px 40px;
  border-radius: 25px;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background: linear-gradient(135deg, #1d2a6c, #0d1b2a);
  border: 2px solid #ffffff;
  text-transform: uppercase;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
  min-width: 150px;
  max-width: 300px;
}

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

/* Smooth hover effect */
.start-button:hover {
  background: #333;
  box-shadow: 0 0 20px rgba(255, 255, 255, 1);
}
</style>
