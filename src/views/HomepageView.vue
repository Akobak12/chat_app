<template>
  <div
    class="flex justify-center items-center h-4/5 w-11/12 mt-12 overflow-auto text-white bg-midnight-blue border-8 border-clay-purple
    lg:w-3/4"
    @click="closeWindow">
    <!--left bar-->
    <section class="flex flex-col items-center relative w-full h-full
    lg:border-clay-purple lg:border-r-8 lg:w-2/5">
      <section class="w-full h-2/5 overflow-y-scroll p-1 no-scrollbar
      lg:h-full">
        <Server v-for="server in servers" :server="server" :id="server.id" :key="server.id" class="mb-2"/>
        <Chat v-for="chat in chats" :chat="chat" :id="chat.id" :key="chat.id" @click="saveId(chat.id)" class="mb-2"/>
      </section>
      <section class="flex flex-col justify-center items-center flex-1 w-full border-clay-purple border-t-8
      lg:justify-end">
        <div class="h-3/4 w-5/6 bg-dragon-purple border-4 border-clay-purple rounded-lg
        lg:hidden">
          <h1 class="text-2xl mb-4">Friend requests</h1>
          <section class="px-4">  
            <User v-for="user in requests" :user="user" :id="user.id" :key="user.id" :request="true" @reloadRequests="loadRequestId()" class="mb-2 h-16"/>
          </section>
        </div>
        <div class="flex bottom-4 text-md w-full px-4 pt-4
        lg:text-lg lg:mb-4">
          <button @click.stop="openInviteWindow" class="w-full mr-6 h-12 bg-dragon-purple rounded-xl
          lg:mx-12 xl:mx-24">Invite friends</button>
          <button @click.stop="create" class="w-full h-12 bg-dragon-purple rounded-xl
          lg:hidden">Create server</button>
        </div>
      </section>    
    </section>
    <!--right bar-->
    <section class="hidden grid-cols-2 grid-rows-4 gap-6 justify-between flex-grow h-full p-4
    lg:grid">
      <button @click.stop="create" class="bg-dragon-purple rounded-lg">Create server</button>
      <div class="col-span-2 row-span-3 bg-dragon-purple border-4 border-clay-purple rounded-lg">
        <h1 class="text-2xl mb-4">Friend requests</h1>
        <section class="px-4">  
          <User v-for="user in requests" :user="user" :id="user.id" :key="user.id" :request="true" @reloadRequests="loadRequestId()" class="mb-2 h-16"/>
        </section>
      </div>
    </section>
    <!--popup windows-->
    <div v-if="invite" class="flex justify-center items-center absolute h-screen w-screen top-0 z-50">
      <section class="absolute w-3/4 bg-midnight-blue border-clay-purple border-4 rounded-2xl window
      lg:w-1/4 md:w-2/4">
        <section class="p-1">
          <div class="flex w-full h-12 pr-2 rounded-xl bg-white text-black">
            <input type="text" placeholder="Input user ID" v-model="searchUser" class="w-full bg-transparent outline-none h-full pl-4">
            <button class="search-button">
              <img src="../assets/search.png">
            </button>
          </div>
        </section>
        <span>{{ inviteMsg }}</span>
        <section class="bottom-4 w-full">
          <div class="flex relative items-center mb-4 h-16 w-full  border-clay-purple border-y-8"
          :class="{ 'bg-[#282828]': darkMode, 'bg-[#64416a]': !darkMode}">
            <Server v-for="server in servers" :server="server" :key="server.id" class="w-fit h-2/3 m-2"/>
          </div>
          <button @click="inviteUser" class="text-lg w-1/2 h-12 mb-4 bg-dragon-purple rounded-xl">Invite</button>
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
      <button class="text-lg w-1/2 h-12 bg-dragon-purple rounded-xl" @click="createServer">Create server</button>
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

export default {
  components: {
    Chat,
    User,
    Server,
  },
  setup(props, {emit}) {
    const searchUser = ref()
    const createFriend = inject("createFriend")
    const getFriendReq = inject("getFriendReq")
    const getUser = inject("getUser")
    const getFriends = inject("getFriends")
    const createRoom = inject("createRoom")
    const getRooms = inject("getRooms")
    const darkMode = ref(inject("darkMode"))

    const chatId = ref([])
    const chats = ref([])
    const requestId = ref([])
    const requests = ref([])
    const serverId = ref([])
    const servers = ref([])

    const serverName = ""
    const inviteMsg = ref("")

    const invite = ref(false)
    const server = ref(false)

    const privacy = ref(true)

    const loadChats = async () => {
      try {
        const response = await axios.get(getFriends, { withCredentials:true });
        chatId.value = new Set(response.data.friends)
        chats.value = await getUsers(chatId.value);
        chats.value.sort()
        console.log(chats.value)
        console.log("Chats load successful:", response);
      } catch (error) {
        console.error("Chats load failed:", error);
      }
    };

    const loadServers = async () => {
      try {
        const response = await axios.get(getRooms, { withCredentials:true });
        serverId.value = new Set(response.data)
        chats.value = await getUsers(chatId.value);
        console.log(chats.value)
        console.log("Chats load successful:", response);
      } catch (error) {
        console.error("Chats load failed:", error);
      }
    };

    const loadRequestId = async () => {
      console.log("nope")
      try {
        const response = await axios.get(getFriendReq, { withCredentials:true })
        requestId.value = new Set(response.data.requests)
        requests.value = await getUsers(requestId.value);
        console.log("req Load successful:", response)
      }
      catch (error) {
        console.log("req Load failed:", error)
      }
    }

    const getUsers = async (array) => {
      const req = [];
      const promises = [];

      for (const id of array) {
        if (id != 0) {
          const promise = axios.get(getUser + id, { withCredentials: true })
            .then(response => {
              req.push({ id: id, ...response.data });
            })
            .catch(error => {
              console.log("Load failed:", error);
            });

          promises.push(promise);
        }
      }

      await Promise.all(promises);

      return req;
    };

    const createServer = async () => {
      try {
        const response = await axios.post(createRoom, {
          name: serverName
        },{ withCredentials:true })

        console.log("Server successful:", response)
      }
      catch (error) {
        console.log("Server failed:", error)
      }
    }
    
    onMounted(() => {
      loadChats()
      loadRequestId()
    })

    const inviteUser = async () => {
      try {
        const response = await axios.post(createFriend, {
          id: parseInt(searchUser.value)
        },{ withCredentials:true });
        inviteMsg.value = "Invite sent"
        console.log("Invite successful:", response);        
      } 
      catch (error) {
        inviteMsg.value = "Error has occurred"
        console.error("Invite failed:", error); 
      }  
    }

    const closeWindow = (event) => {
      const windowElement = document.querySelector(".window");
      if (invite.value == true || server.value == true) {
        if (windowElement && !windowElement.contains(event.target)) {
          invite.value = false
          server.value = false
      }}
    };

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

    return { requests, servers, invite, server, privacy, chats, serverName, darkMode, searchUser, inviteMsg, loadServers, createServer, inviteUser, create, openInviteWindow, saveId, loadChats, closeWindow };
  },
};
</script>

<style scoped>
.window {
  box-shadow: 0 0 0 max(100vh, 100vw) rgba(0, 0, 0, 0.4);
}

.no-scrollbar::-webkit-scrollbar {
    display: none;
}
</style>