<template>
  <main-layout>
    <section class="hero">
      <h1>{{ $t('welcome.title') }}</h1>
      <p>{{ $t('welcome.description') }}</p>
      <transition-button :to="'/form'">
        {{ $t('welcome.buttonText') }}
      </transition-button>
    </section>
  </main-layout>
</template>

<script>
import MainLayout from '@/layouts/MainLayout.vue';
import TransitionButton from '@/components/TransitionButton.vue';

export default {
  components: { MainLayout, TransitionButton },
};
</script>

<style scoped>
.hero {
  text-align: center;
  margin: 50px auto;
}
</style>
