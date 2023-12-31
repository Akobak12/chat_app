<template>
  <div
    class="h-4/5 w-3/4 overflow-auto bg-midnight-blue border-8 border-clay-purple">
    <section class="flex flex-col items-center relative w-2/5 h-full p-1 border-clay-purple border-r-8">
      <User v-for="user in users" :user="user" :id="user.id" :key="user.id" @click="saveId(user.id)" class="mb-2"/>
      <button @click="logout" class="absolute bottom-4 text-lg text-white w-1/2 h-12 bg-dragon-purple rounded-xl">logout</button>
    </section>
    <section></section>
  </div>
</template>

<script>
import axios from "axios";
import { onMounted, ref } from "vue"
import User from "../components/UserComponent.vue";
import router from '@/router';

export default {
  components: {
    User,
  },
  setup(props, {emit}) {
    // array with rooms
    const users = ref([     
    /*{
      userName: "pablo",
      profilePicture: "",
      online: true,
      lastActive: "Active now",
      id: 1
    },
    {
      userName: "bruh",
      profilePicture: "",
      online: false,
      lastActive: "last active",
      id: 2
    },*/
    ]);

    const loadChats = async () => {
      try {
        const response = await axios.post("http://127.0.0.1:3030/get-chats", {
          username: localStorage.getItem("username"),
        });
        users.value = response.data;
        console.log("Load successful:", response);
      } catch (error) {
        console.error("Load failed:", error);
      }
    };
    
    onMounted(() => {
      loadChats()
    })

    const logout = () => {
      localStorage.setItem("isAuth", "")
      router.push({ name: "login" })
    }

    const saveId = (id) => {
      localStorage.setItem("id", id)
      emit('updateUserId', id);
    }

    return { users, logout, saveId, loadChats };
  },
};
</script>

<style>
</style>