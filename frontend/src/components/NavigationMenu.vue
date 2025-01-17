<template>
    <navigation class="navigation-menu">
      <router-link to="/" class="menu-item">
        {{ $t('menu.home') }}
      </router-link>
      <router-link to="/about" class="menu-item">
        {{ $t('menu.about') }}
      </router-link>
    </navigation>
  </template>
  
  <script>
  export default {
    name: 'NavigationMenu',
  };
  </script>
  
  <style scoped>
  .navigation-menu {
    display: flex;
    gap: 1.5rem;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg, rgba(0, 0, 51, 0.9), rgba(20, 24, 82, 0.8), rgba(47, 79, 79, 0.9));
    border-radius: 16px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
    position: fixed;
    top: 3px;
    left: 1px;
  }
  
  .menu-item {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 150px;
    height: 30px;
    color: #f0f8ff;
    font-size: 1rem;
    font-weight: bold;
    text-decoration: none;
    text-transform: uppercase;
    letter-spacing: 1px;
    transition: color 0.3s ease, transform 0.3s ease;
    padding: 5px 10px;
    border-radius: 8px;
    position: relative; /* Нужно для псевдоэлементов */
  }
  
  .menu-item:hover {
    color: #87ceeb;
    transform: scale(1.1);
  }
  
  .menu-item::before,
  .menu-item::after {
    content: '';
    position: absolute;
    width: 0;
    height: 2px;
    background: #87ceeb;
    transition: width 0.4s ease, opacity 0.4s ease;
    opacity: 0;
  }
  
  .menu-item::before {
    bottom: 0;
    left: 0;
  }
  
  .menu-item::after {
    top: 0;
    right: 0;
  }
  
  .menu-item:hover::before,
  .menu-item:hover::after {
    width: 100%;
    opacity: 1;
  }
  
  .navigation-menu::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 140%;
    height: 140%;
    background: radial-gradient(circle, rgba(72, 61, 139, 0.4), transparent 70%);
    transform: translate(-50%, -50%);
    z-index: -1;
    filter: blur(15px);
  }
  </style>
  