<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <div class="page-wrapper">
        <component :is="Component" />
      </div>
    </router-view>
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
  display: flex;
  flex-direction: column;
}

.page-wrapper {
  flex: 1;
  padding-bottom: 70px;
}

@media (min-width: 769px) {
  .page-wrapper {
    padding-bottom: 60px;
  }
}
</style>
