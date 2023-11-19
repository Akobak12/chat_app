<template>
  <main class="relative flex flex-col items-center bg-midnight-puprle w-screen h-screen overflow-hidden">
    <div class="flex justify-center w-screen">
      <img src="./assets/logo.png" class="absolute scale-75 left-6">
      <TopBar class="mb-12"/>
    </div>
    
    <router-view/>
  </main>
</template>

<script>
import TopBar from "./components/top-bar/TopBar.vue";

export default {
  components: {
    TopBar
  }
}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100vh;
}

</style>
