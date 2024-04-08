<template>
  <img src="../assets/logo.png" class="h-32 my-6" />
  <form class="flex flex-col items-center justify-between w-4/5 h-3/5 bg-midnight-blue text-white border-8 rounded-md border-clay-purple
  sm:w-2/3 md:w-3/5 lg:w-1/3">
    <span class="text-4xl mt-4 mb-8">LOGIN</span>
    <input v-model="credentials.username" type="text" placeholder="Username" />
    <input v-model="credentials.email" type="email" placeholder="Email">
    <input v-model="credentials.password" type="password" placeholder="Password"/>
    <button @click.prevent="register" class="w-48 h-12 text-2xl rounded-lg bg-dragon-purple">Continue</button>
    <router-link to="/login" class="flex mt-6 text-sm"><u>Already have a account?</u></router-link>
  </form>
</template>

<script>
import axios from "axios";
import { inject, ref } from 'vue'
import router from '@/router';

export default {
  setup() {
    const registerURL = inject("register")

    const credentials = ref({
      username: "",
      email: "",
      password: "",
    });

    const register = async () => {
      try {
        const response = await axios.post(registerURL, {
          username: credentials.value.username,
          email: credentials.value.email,
          password: credentials.value.password,
        });
        console.log("register successful:", response);
        router.push({ name: "login" })
      } catch (error) {
        console.error("register failed:", error);
      }
    };

    return { register, credentials };
  },
};
</script>


<style scoped>
input {
  height: 15%;
  width: 90%;
  margin-bottom: 1rem;
  color: black;
  font-size: 1.25rem;
  line-height: 1.75rem;
  padding-left: 0.75rem;
  border-radius: 0.375rem;
}
</style>