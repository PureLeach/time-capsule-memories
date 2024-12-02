<template>
    <el-button
      type="primary"
      :to="to"
      @click="navigate"
      class="transition-button"
    >
      <slot />
    </el-button>
  </template>
  
  <script>
  import { useRouter } from 'vue-router';
  
  export default {
    props: {
      to: {
        type: String,
        required: true,
      },
    },
    setup(props) {
      const router = useRouter();
  
      const navigate = () => {
        router.push(props.to);
      };
  
      return { navigate };
    },
  };
  </script>
  
  <style scoped>
  .transition-button {
    transition: all 0.3s ease;
  }
  
  .transition-button:hover {
    transform: scale(1.05);
  }
  </style>
  