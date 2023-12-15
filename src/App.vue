<template>
  <main
    class="relative flex flex-col items-center bg-midnight-puprle w-screen h-screen overflow-hidden"
  >
    <div class="flex justify-center w-screen" v-if="$route.path != '/login' && $route.path != '/register'">
      <img src="./assets/logo.png" class="absolute scale-75 left-6" />
      <TopBar class="mb-12" />
    </div>

    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component"></component>
      </keep-alive>
    </router-view>
  </main>
</template>

<script>
import { provide, onUnmounted } from "vue";
import TopBar from "./components/top-bar/TopBar.vue";

export default {
  components: {
    TopBar,
  },

  setup() {
    const websocket = new WebSocket("ws://127.0.0.1:3030/ws");

    websocket.onopen = () => {
      console.log("WebSocket connection established");
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    websocket.onclose = (event) => {
      console.log("WebSocket connection closed:", event);
    };

    onUnmounted(() => {
      if (websocket.readyState === WebSocket.OPEN) {
        websocket.close();
      }
    });

    provide("websocket", websocket);

    return { websocket };
  },
};
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
