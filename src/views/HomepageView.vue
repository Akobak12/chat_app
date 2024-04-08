<template>
  <div
    class="flex justify-center items-center h-4/5 w-11/12 mt-12 overflow-auto text-white bg-midnight-blue border-8 border-clay-purple
    lg:w-3/4"
    @click="closeWindow">
    <!--left bar-->
    <section class="flex flex-col items-center relative w-full h-full p-1 
    lg:border-clay-purple lg:border-r-8 lg:w-2/5">
      <Chat v-for="chat in chats" :chat="chat" :id="chat.id" :key="chat.id" @click="saveId(chat.id)" class="mb-2"/>
      <div class="flex absolute bottom-4 text-md w-full px-4 pt-4 border-clay-purple border-t-8
      lg:border-t-0 lg:text-lg">
        <button @click.stop="openInviteWindow" class="w-full mr-6 h-12 bg-dragon-purple rounded-xl
        lg:mx-12 xl:mx-24">Invite friends</button>
        <button @click.stop="create" class="w-full h-12 bg-dragon-purple rounded-xl
        lg:hidden">Create server</button>
      </div>
    </section>
    <!--right bar-->
    <section class="hidden justify-between flex-grow h-full p-4
    lg:flex">
      <button @click="logout" class="w-5/12 h-1/5 bg-dragon-purple rounded-lg">Logout</button>
      <button @click.stop="create" class="w-5/12 h-1/5 bg-dragon-purple rounded-lg">Create server</button>
    </section>
    <!--popup windows-->
    <div v-if="invite" class="flex justify-center items-center absolute h-screen w-screen top-0 z-50">
      <section class="absolute h-3/5 w-3/4 bg-midnight-blue border-clay-purple border-4 rounded-2xl window
      lg:w-1/4 md:w-2/4">
        <section class="p-1">
          <Search :placeholder="'Search or add friends'" class="h-12"/>
        </section>
        <section class="px-4">  
          <User v-for="user in users" :user="user" :id="user.id" :key="user.id" @click="selectUser(user)" class="mb-2 h-16" :class="{ 'border-clay-purple border-2': user.selected }"/>
        </section>
        <section class="bottom-4 absolute w-full">
          <div class="flex relative items-center mb-4 h-16 w-full bg-[#64416a] border-clay-purple border-y-8"
          :class="{ 'bg-[#282828]': darkMode }">
            <Server v-for="server in servers" :server="server" :key="server.id" class="w-fit h-2/3 m-2"/>
          </div>
          <button @click="inviteUser()" class="text-lg w-1/2 h-12 bg-dragon-purple rounded-xl">Invite</button>
        </section>
      </section>
    </div>
    
    <div v-if="server" class="flex justify-center items-center absolute h-screen w-screen top-0 z-50">
      <section class="flex flex-col items-center absolute  w-3/4 py-3 bg-midnight-blue border-clay-purple border-4 rounded-2xl window
      lg:w-1/4 md:w-2/4">
      <h1>Server specs</h1>
      <hr class="bg-clay-purple h-1 w-11/12 mb-4 border-0 rounded-full">
      <section class="px-4 mb-4">
        <input type="text" placeholder="Name" v-model="serverName" class="h-12 w-full mb-2 pl-3 text-black rounded-lg">
        <input type="text" placeholder="Bio" class="h-12 w-full mb-2 pl-3 text-black rounded-lg">
        <section class="flex justify-between items-center">
          <span class="ml-6">Logo:</span>
          <button class="flex justify-between items-center h-12 w-2/3 px-3 rounded-lg bg-dragon-purple">
          Upload image
          <img src="../assets/upload.png" class="h-1/2">
          </button>
        </section>
      </section>
      <h2 class="text-xl text-gray-400">Privacy</h2>
      <hr class="bg-clay-purple h-1 w-10/12 mb-2 border-0 rounded-full">
      <div class="mb-6">
        <button @click="privacy = true" class="mr-2 px-2 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : privacy}">Public</button>
        <span class="mr-2">/</span>
        <button @click="privacy = false" class="px-1 rounded-md" :class="{'bg-[rgba(255,255,255,0.1)]' : !privacy}">Private</button>
      </div>
      <button class="text-lg w-1/2 h-12 bg-dragon-purple rounded-xl">Create server</button>
    </section>
    </div>
    
  </div>
</template>

<script>
import axios from "axios";
import { inject, onMounted, ref } from "vue"
import Chat from "../components/ChatComponent.vue";
import User from "../components/UserComponent.vue";
import Server from "../components/ServerComponent.vue";
import Search from "../components/SearchComponent.vue"
import router from '@/router';

export default {
  components: {
    Chat,
    User,
    Server,
    Search
  },
  setup(props, {emit}) {
    //const getRooms = inject("getRooms")
    const createRoom = inject("createRoom")
    const logoutURL = inject("logout")
    const darkMode = ref(inject("darkMode"))

    const chats = ref([
      {id: "1", name: "ola"},
      {id: "2", name: "user"}
    ])
    const users = ref([
      {id: "3", username: "idk", lastActive: "1991", online: false, selected: false},
      {id: "4", username: "bro", lastActive: "Active now", online: true, selected: false}
    ])
    const servers = ref([
      {id: "1", name: "Server 1"},
      {id: "2", name: "Server 2"}

    ])

    const serverName = ""

    const invite = ref(false)
    const server = ref(false)

    const privacy = ref(true)

    const loadChats = async () => {
    /*  try {
        const response = await axios.get(getRooms, {
          username: localStorage.getItem("username"),
        });
        chats.value = response.data;
        console.log("Load successful:", response);
      } catch (error) {
        console.error("Load failed:", error);
      }*/
    };
    
    onMounted(() => {
      loadChats()
    })

    const closeWindow = (event) => {
      const windowElement = document.querySelector(".window");
      if (invite.value == true || server.value == true) {
        if (windowElement && !windowElement.contains(event.target)) {
          invite.value = false
          server.value = false
      }}
    };

    const selectUser = (user) => {
      users.value.forEach(u => {
        u.selected = (u.id === user.id) && !u.selected;
      });
      console.log(users.value)
    };

    const inviteUser = async () => {
      for (const user of users.value) {
          if (user.selected == true && !chats.value.find(room => room.id === user.id)) {
            try {
              const response = await axios.post(createRoom, {          
              id: user.id,
              name: user.username
            });
            chats.value.push(response.data)
            invite.value = false
            console.log(chats.value )
            console.log("Invite successful:", response);        
            } 
            catch (error) {
              console.error("Invite failed:", error);
              console.log(user.id, user.username)
            }     
          }           
        }   
    }

    const logout = async () => {
      try {
        const response = await axios.get(logoutURL);
        localStorage.setItem("isAuth", false)
        console.log(localStorage.getItem("isAuth"))
        router.push({ name: "login" })
        console.log("Logout successful:", response);
      } catch (error) {
        console.error("Logout failed:", error);
      }
    }

    const openInviteWindow = () => {
      invite.value = true
    }

    const create = () => {
      server.value = true
    }

    const saveId = (id) => {
      localStorage.setItem("id", id)
      emit('updateUserId', id);
    }

    return { users, servers, invite, server, privacy, chats, serverName, darkMode, inviteUser, selectUser, create, openInviteWindow, saveId, loadChats, logout, closeWindow };
  },
};
</script>

<style scoped>
.window {
  box-shadow: 0 0 0 max(100vh, 100vw) rgba(0, 0, 0, 0.4);
}
</style>