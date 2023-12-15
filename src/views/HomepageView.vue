<template>
  <div
    class="h-4/5 w-3/4 overflow-auto bg-midnight-blue border-8 border-clay-purple">
    <section class="flex flex-col items-center relative w-2/5 h-full p-1 border-clay-purple border-r-8">
      <User v-for="user in users" :user="user" :key="user.id" class="mb-2"/>
      <button @click="logout" class="absolute bottom-4 text-lg text-white w-1/2 h-12 bg-dragon-purple rounded-xl">logout</button>
    </section>
    <section></section>
  </div>
</template>

<script>
import { ref } from "vue"
import User from "../components/UserComponent.vue";
import router from '@/router';

export default {
  components: {
    User,
  },
  setup() {
    const users = ref([{
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
    },
    ]);

    const logout = () => {
      localStorage.setItem("isAuth", "")
      router.push({ name: "login" })
    }

    return { users, logout };
  },
};
</script>

<style>
</style>