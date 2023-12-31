<template>
  <main
    class="relative flex flex-col items-center bg-midnight-puprle w-screen h-screen overflow-hidden"
  >
    <div class="flex justify-center w-screen" v-if="$route.path != '/login' && $route.path != '/register'">
      <img src="./assets/logo.png" class="absolute scale-75 left-6" />
      <TopBar :userId="userId" class="mb-12"/>
    </div>

    <router-view v-slot="{ Component }" @updateUserId="setUserId  " >
      <keep-alive>
        <component :is="Component"></component>
      </keep-alive>
    </router-view>
  </main>
</template>

<script>
import { provide, onUnmounted, ref } from "vue";
import TopBar from "./components/top-bar/TopBar.vue";

export default {
  components: {
    TopBar,
  },

  setup() {
    const userId = ref(null);

    const setUserId = (id) => {
      userId.value = id;
    };

    if (localStorage.getItem("id") === null){
      localStorage.setItem("id", "")
    }

    const websocket = ref();

    if (localStorage.getItem('isAuth')) {
      websocket.value = new WebSocket("ws://127.0.0.1:3030/ws");

      websocket.value.onopen = () => {
        console.log("WebSocket connection established");
      };

      websocket.value.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      websocket.value.onclose = (event) => {
        console.log("WebSocket connection closed:", event);
      };
      console.log("1", websocket.value)
    }
    console.log("2", websocket.value)
    

    onUnmounted(() => {
      if (websocket.value.readyState === WebSocket.OPEN) {
        websocket.value.close();
      }
    });

    provide("websocket", websocket);

    return { websocket, userId, setUserId };
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
