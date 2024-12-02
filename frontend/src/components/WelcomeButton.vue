<template>
    <router-link to="/" class="welcome-button">
      <button class="btn">{{ $t('welcome.title') }}</button>
    </router-link>
  </template>
  
  <script>
  export default {
    name: 'WelcomeButton',
  };
  </script>
  
  <style scoped>
  .welcome-button {
    text-decoration: none;
  }
  
  .btn {
    background-color: #4CAF50;
    color: white;
    padding: 10px 20px;
    font-size: 1.2rem;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s;
  }
  
  .btn:hover {
    background-color: #45a049;
  }
  </style>
  