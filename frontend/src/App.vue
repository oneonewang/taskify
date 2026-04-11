<template>
  <div id="app">
    <router-view />
    <AppBottom />
  </div>
</template>

<script setup lang="ts">
import AppBottom from './components/AppBottom.vue'
</script>

<style>
/* Global Reset - Additional styles in assets/styles/global.css */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  width: 100%;
  height: 100%;
}

/* Typography defaults handled by global.css */

/* Prevent FOUC */
#app {
  opacity: 1;
}

/* Add padding-bottom to main content to account for bottom nav on mobile */
@media (max-width: 768px) {
  body {
    padding-bottom: 64px;
  }
}
</style>
