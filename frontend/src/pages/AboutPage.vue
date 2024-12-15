<template>
    <main-layout>
      <div class="about-page">
        <div class="about-form-container">
          <h1 class="about-title">{{ $t('about.title') }}</h1>
          <p class="about-description">{{ $t('about.description') }}</p>
          <form class="about-form">
            <label for="message" class="form-label">{{ $t('about.form.label') }}</label>
            <textarea
              id="message"
              class="form-textarea"
              rows="5"
              placeholder="{{ $t('about.form.placeholder') }}"
            ></textarea>
            <button type="submit" class="form-button">{{ $t('about.form.submit') }}</button>
          </form>
        </div>
      </div>
    </main-layout>
  </template>
  
  <script>
  import MainLayout from '@/layouts/MainLayout.vue';
  
  export default {
    name: 'AboutPage',
    components: {
      MainLayout
    },
  };
  </script>
  
  <style scoped>
  .about-page {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    /* padding: 2rem; */
    padding-top: 5rem;
    /* Убираем фон */
  }
  
  
  .about-form-container {
    max-width: 400px;
    width: 100%;
    background: radial-gradient(circle, rgba(25, 25, 112, 0.9), rgba(0, 0, 51, 0.8), rgba(47, 79, 79, 0.9));
    padding: 2rem;
    border-radius: 16px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    color: white;
    text-align: center;
  }
  
  .about-title {
    font-size: 1.5rem;
    margin-bottom: 1rem;
  }
  
  .about-description {
    font-size: 1rem;
    margin-bottom: 2rem;
  }
  
  .about-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  .form-label {
    font-size: 1rem;
    text-align: left;
  }
  
  .form-textarea {
    border: none;
    padding: 1rem;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    font-size: 1rem;
    resize: none;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.5);
  }
  
  .form-button {
    background: linear-gradient(145deg, #ffd700, #daa520);
    color: #fff8dc;
    padding: 0.8rem 2rem;
    font-size: 1rem;
    border: none;
    border-radius: 30px;
    cursor: pointer;
    font-weight: bold;
    text-transform: uppercase;
    transition: background 0.3s ease, transform 0.2s ease;
  }
  
  .form-button:hover {
    background: linear-gradient(145deg, #daa520, #ffd700);
    transform: scale(1.05);
  }
  
  .form-button:active {
    background: linear-gradient(145deg, #b8860b, #ffd700);
    transform: scale(1);
  }
  </style>
  