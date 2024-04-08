<template>
  <div class="flex flex-col w-full h-full bg-white text-white
  lg:h-4/5 lg:w-2/5 lg:mt-12">
    <section class="hidden items-center pl-16 bg-dragon-purple w-full h-1/6 border-8 border-clay-purple
    lg:flex">
        <h1 class="text-2xl">Settings</h1>
    </section>
    <div class="flex flex-1 bg-midnight-puprle">
        <section class="bg-dragon-purple flex flex-col items-center w-1/6
        sm:w-1/12
        lg:w-2/5 lg:border-8 lg:border-clay-purple">
            <button 
            class="flex justify-center items-center w-11/12 h-10 my-1 rounded-lg hover:bbtn-highligh
            lg:justify-normal" 
            :class="{'bg-btn-highligh': general}"
            @click="highlight('general')">  
              <img src="../assets/settings/general.png" class="h-3/5 mx-2">
              <span class="hidden lg:block">general</span>
            </button>
            <hr class="h-[3px] border-0 w-10/12 rounded-full bg-break-line
            lg:w-11/12">
            <button 
            class="flex justify-center items-center w-11/12 h-10 my-1 rounded-lg hover:bbtn-highligh
            lg:justify-normal" 
            :class="{'bg-btn-highligh': notifications}"
            @click="highlight('notifications')">
              <img src="../assets/settings/notifications.png" class="h-3/5 mx-2">
              <span class="hidden lg:block">notifications</span>
            </button>
            <hr class="h-[3px] border-0 w-10/12 rounded-full bg-break-line
            lg:w-11/12">
            <button 
            class="flex justify-center items-center w-11/12 h-10 my-1 rounded-lg hover:bbtn-highligh
            lg:justify-normal" 
            :class="{'bg-btn-highligh': friends}"
            @click="highlight('friends')">
              <img src="../assets/settings/friends.png" class="h-3/5 mx-2">
              <span class="hidden lg:block">friends</span>
            </button>
            <hr class="h-[3px] border-0 w-10/12 rounded-full bg-break-line
            lg:w-11/12">
            <button 
            class="flex justify-center items-center w-11/12 h-10 my-1 rounded-lg hover:bbtn-highligh
            lg:justify-normal" 
            :class="{'bg-btn-highligh': server}"
            @click="highlight('servers')">
              <img src="../assets/settings/servers.png" class="h-3/5 mx-2">
              <span class="hidden lg:block">servers</span>
            </button>
            <hr class="h-[3px] border-0 w-10/12 rounded-full bg-break-line
            lg:w-11/12">
        </section>
        <section class="flex-1 bg-midnight-blue border-8 border-clay-purple scale-95
        lg:scale-100">
          <!--general-->
          <section class="flex flex-col items-center pt-4" v-if="general">
            <h1 class="text-2xl mb-10">general settings</h1>
            <h2 class="text-xl text-gray-400">Visuals</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <div class="flex items-center mt-1 mb-5">
              <small class="mr-6">Color palette</small>
              <span @click="setPallete('purple')" class="flex justify-center items-center h-5 w-5 mr-2 bg-[#2f0136] rounded-full">
                <img v-if="purple" src="../assets/check.png" class="h-1/2">
              </span>
              <span @click="setPallete('dark')" class="flex justify-center items-center h-5 w-5 mr-2 bg-black  rounded-full">
                <img v-if="dark" src="../assets/check.png" class="h-1/2">
              </span>
          </div>
            <h2 class="text-xl text-gray-400">Active status</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <div class="mt-1 mb-5">
              <small>
                <button @click="status = true" class="mr-2 h-fit w-auto px-1 rounded-lg" :class="{'bg-[rgba(255,255,255,0.1)]' : status}">Show</button>
                <span class="mr-2">/</span>
                <button @click="status = false" class="px-1 rounded-lg" :class="{'bg-[rgba(255,255,255,0.1)]' : !status}">Dont show</button>
              </small>
            </div>
            <h2 class="text-xl text-gray-400">Another settings</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <small class="mt-1">Placeholder text</small>
          </section>
          <!--notifications-->
          <section class="flex flex-col items-center pt-4" v-if="notifications">
            <h1 class="text-2xl mb-10">Notification settings</h1>
            <h2 class="text-xl text-gray-400">Sound</h2>
            <hr class="h-1 w-10/12 mb-2 border-0 rounded bg-break-line-full">
            <div class="mb-5">
              <button @click="sound = true" class="mr-2 px-2 rounded-lg" :class="{'bg-[rgba(255,255,255,0.1)]' : sound}">ON</button>
              <span class="mr-2">/</span>
              <button @click="sound = false" class="px-1 rounded-lg" :class="{'bg-[rgba(255,255,255,0.1)]' : !sound}">OFF</button>
            </div>
            <h2 class="text-xl text-gray-400">Another settings</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <small class="mt-1">Placeholder text</small>
          </section>
          <!--friends-->
          <section class="flex flex-col items-center pt-4" v-if="friends">
            <h1 class="text-2xl mb-10">Friends settings</h1>
            <h2 class="text-xl text-gray-400">Mute</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <div class="w-10/12 mt-10 border-clay-purple border-4 p-1 rounded-3xl">
              <Search :placeholder="'Search friends'" class="h-10"/>
              <User v-for="user in users" :user="user" :mute="true" :key="user.id" @muteUser="user.muted = !user.muted" class="h-12 mb-1"/>
            </div>          
          </section>
          <!--servers-->
          <section class="flex flex-col items-center pt-4" v-if="server">
            <h1 class="text-2xl mb-10">Servers settings</h1>
            <h2 class="text-xl text-gray-400">Mute</h2>
            <hr class="h-1 border-0 w-10/12 rounded-full bg-break-line">
            <div class="w-10/12 mt-10 border-clay-purple border-4 p-1 rounded-3xl">
              <Search :placeholder="'Search servers'" class="h-10"/>
              <Server v-for="server in servers" :server="server" :mute="true" :key="server.id" @muteServer="server.muted = !server.muted" class="w-full h-2/3 mb-1"/>
            </div>          
          </section>
        </section>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import User from "../components/UserComponent.vue"
import Server from "../components/ServerComponent.vue"
import Search from "../components/SearchComponent.vue"

export default {
  components: {
    User,
    Server,
    Search
  },
  setup(props, { emit }) {
    const general = ref(true)
    const notifications = ref(false)
    const friends = ref(false)
    const server = ref(false)

    const purple = ref(true)
    const dark = ref(false)
    const status = ref(true)
    const sound = ref(true)

    if (localStorage.getItem("darkMode") == "true") {
      purple.value = false
      dark.value = true
    }

    const users = ref([
      {id: "1", username: "idk", muted: false},
      {id: "2", username: "username", muted: true}
    ])

    const servers = ref([
      {id: "1", name: "server 1", muted: false},
      {id: "2", name: "server 2", muted: false}
    ])

    const setPallete = (color) => {
      purple.value = false
      dark.value = false

      if (color == "purple") {
        purple.value = true
        localStorage.setItem("darkMode", false)
        emit("new-pallete", false);
      }
      else if (color == "dark") {
        dark.value = true
        localStorage.setItem("darkMode", true)
        emit("new-pallete", true);
      }
    }

    const highlight = (button) => {
      general.value = false
      notifications.value = false
      friends.value = false
      server.value = false

      if (button == "general") {
        general.value = true
      }
      else if (button == "notifications") {
        notifications.value = true
      }
      else if (button == "friends") {
        friends.value = true
      }
      else if (button == "servers") {
        server.value = true
      }
    }
    return{ general, notifications, friends, server, status, sound, purple, dark, users, servers, setPallete, highlight }
},
};
</script>

<style scoped>

</style>
