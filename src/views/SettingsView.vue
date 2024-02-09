<template>
  <div class="flex flex-col w-2/5 h-4/5 bg-white text-white">
    <section class="flex items-center pl-16 bg-dragon-purple w-full h-1/6 border-8 border-clay-purple">
        <h1 class="text-2xl">Settings</h1>
    </section>
    <div class="flex flex-1">
        <section class="bg-dragon-purple flex flex-col items-center w-2/5 h-full border-8 border-clay-purple">
            <button 
            class="flex items-center w-11/12 h-10 my-1 rounded-md hover:bg-[#522c58]" 
            :class="{'bg-[#522c58]': general}"
            @click="highlight('general')">
              <img src="../assets/settings/general.png" class="h-3/5 mx-2">
              <span>general</span>
            </button>
            <hr class="bg-clay-purple h-1 border-0 w-11/12 rounded-full">
            <button 
            class="flex items-center w-11/12 h-10 my-1 rounded-md hover:bg-[#522c58]" 
            :class="{'bg-[#522c58]': notifications}"
            @click="highlight('notifications')">
              <img src="../assets/settings/notifications.png" class="h-3/5 mx-2">
              <span>notifications</span>
            </button>
            <hr class="bg-clay-purple h-1 border-0 w-11/12 rounded-full">
            <button 
            class="flex items-center w-11/12 h-10 my-1 rounded-md hover:bg-[#522c58]" 
            :class="{'bg-[#522c58]': friends}"
            @click="highlight('friends')">
              <img src="../assets/settings/friends.png" class="h-3/5 mx-2">
              <span>friends</span>
            </button>
            <hr class="bg-clay-purple h-1 border-0 w-11/12 rounded-full">
            <button 
            class="flex items-center w-11/12 h-10 my-1 rounded-md hover:bg-[#522c58]" 
            :class="{'bg-[#522c58]': server}"
            @click="highlight('servers')">
              <img src="../assets/settings/servers.png" class="h-3/5 mx-2">
              <span>servers</span>
            </button>
            <hr class="bg-clay-purple h-1 border-0 w-11/12 rounded-full">
        </section>
        <section class="flex-1 bg-midnight-blue border-8 border-clay-purple">
          <!--general-->
          <section class="flex flex-col items-center pt-4" v-if="general">
            <h1 class="text-2xl mb-10">general settings</h1>
            <h2 class="text-xl text-gray-400">Visuals</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
            <div class="flex items-center mt-1 mb-5">
              <small class="mr-6">Color palette</small>
              <span @click="setPallete('purple')" class="flex justify-center items-center h-5 w-5 mr-2 bg-dragon-purple rounded-full">
                <img v-if="purple" src="../assets/check.png" class="h-1/2">
              </span>
              <span @click="setPallete('dark')" class="flex justify-center items-center h-5 w-5 mr-2 bg-green-950 rounded-full">
                <img v-if="darkGreen" src="../assets/check.png" class="h-1/2">
              </span>
              <span @click="setPallete('light')" class="flex justify-center items-center h-5 w-5 mr-2 bg-green-600 rounded-full">
                <img v-if="lightGreen" src="../assets/check.png" class="h-1/2">
            </span>
          </div>
            <h2 class="text-xl text-gray-400">Active status</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
            <div class="mt-1 mb-5">
              <small>
                <button @click="status = true" class="mr-2 h-fit w-auto px-1 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : status}">Show</button>
                <span class="mr-2">/</span>
                <button @click="status = false" class="px-1 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : !status}">Dont show</button>
              </small>
            </div>
            <h2 class="text-xl text-gray-400">Another settings</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
            <small class="mt-1">Placeholder text</small>
          </section>
          <!--notifications-->
          <section class="flex flex-col items-center pt-4" v-if="notifications">
            <h1 class="text-2xl mb-10">Notification settings</h1>
            <h2 class="text-xl text-gray-400">Sound</h2>
            <hr class="bg-clay-purple h-1 w-10/12 mb-2 border-0 rounded-full">
            <div class="mb-5">
              <button @click="sound = true" class="mr-2 px-2 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : sound}">ON</button>
              <span class="mr-2">/</span>
              <button @click="sound = false" class="px-1 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : !sound}">OFF</button>
            </div>
            <h2 class="text-xl text-gray-400">Another settings</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
            <small class="mt-1">Placeholder text</small>
          </section>
          <!--friends-->
          <section class="flex flex-col items-center pt-4" v-if="friends">
            <h1 class="text-2xl mb-10">Friends settings</h1>
            <h2 class="text-xl text-gray-400">Mute</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
            <div class="w-10/12 mt-10 border-clay-purple border-4 p-1 rounded-3xl">
              <Search :placeholder="'Search friends'" class="h-10"/>
              <User v-for="user in users" :user="user" :mute="true" :key="user.id" @muteUser="user.muted = !user.muted" class="h-12 mb-1"/>
            </div>          
          </section>
          <!--servers-->
          <section class="flex flex-col items-center pt-4" v-if="server">
            <h1 class="text-2xl mb-10">Servers settings</h1>
            <h2 class="text-xl text-gray-400">Mute</h2>
            <hr class="bg-clay-purple h-1 border-0 w-10/12 rounded-full">
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
  setup() {
    const general = ref(true)
    const notifications = ref(false)
    const friends = ref(false)
    const server = ref(false)

    const purple = ref(true)
    const darkGreen = ref(false)
    const lightGreen = ref(false)
    const status = ref(true)
    const sound = ref(true)

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
      darkGreen.value = false
      lightGreen.value = false

      if (color == "purple") {
        purple.value = true
      }
      else if (color == "light") {
        lightGreen.value = true
      }
      else if (color == "dark") {
        darkGreen.value = true
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
    return{ general, notifications, friends, server, status, sound, purple, darkGreen, lightGreen, users, servers, setPallete, highlight }
},
};
</script>

<style>

</style>