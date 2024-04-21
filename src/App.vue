<template>
  <main
    class="relative flex flex-col items-center bg-midnight-puprle w-screen h-screen overflow-hidden"
    :class="{dark: darkMode}"
  >
    <div class="flex justify-center w-screen" v-if="$route.path != '/login' && $route.path != '/register'">
      <img src="./assets/logo.png" class="absolute scale-75 left-0 -top-3 
      lg:scale-100 lg:top-0 lg:left-8 xl:left-12" />
      <TopBar :chatId="chatId" class=""/>
    </div>

    <router-view v-slot="{ Component }" @updateUserId="setChatId" @new-pallete="changeTheme" >
        <component :is="Component"></component>
    </router-view>
  </main>
</template>

<script>
import { provide, ref } from "vue";
import TopBar from "./components/top-bar/TopBar.vue";

export default {
  components: {
    TopBar,
  },

  setup() {
    const chatId = ref(null);

    const darkMode = ref(false)

    const register = "http://localhost:3030/api/register"
    const login = "http://localhost:3030/api/login"
    const logout = "http://localhost:3030/api/logout"
    const createRoom = "http://localhost:3030/api/ws/create-room"
    const joinRoom = ref(`ws://localhost:3030/api/ws/join-room/${chatId.value}`)
    const createFriend = "http://localhost:3030/api/ws/create-friend"
    const getFriendReq = "http://localhost:3030/api/ws/get-friend-requests"
    const getUser = "http://localhost:3030/api/ws/get-user/"
    const acceptFriend = "http://localhost:3030/api/ws/accept-friend"
    const rejectFriend = "http://localhost:3030/api/ws/reject-friend"
    const getFriends = "http://localhost:3030/api/ws/get-friends"
    const getRooms = "http://localhost:3030/api/ws/get-rooms"

    if (localStorage.getItem("darkMode") == "true") {
      darkMode.value = true
    }
    else {
      darkMode.value = false
    }

    const setChatId = (id) => {
      chatId.value = id;
      joinRoom.value = `ws://localhost:3030/api/ws/join-room/${chatId.value}`
    };

    const changeTheme = (dark) => {
      darkMode.value = dark
    }

    if (localStorage.getItem("id") === null){
      localStorage.setItem("id", "")
    }

    provide("darkMode", darkMode)

    provide("register", register)
    provide("login", login)
    provide("logout", logout)
    provide("createRoom", createRoom)
    provide("joinRoom", joinRoom)
    provide("createFriend", createFriend)
    provide("getFriendReq", getFriendReq)
    provide("getUser", getUser)
    provide("acceptFriend", acceptFriend)
    provide("rejectFriend", rejectFriend)
    provide("getFriends", getFriends)
    provide("getRooms", getRooms)
    return { chatId, setChatId, changeTheme, darkMode };
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
