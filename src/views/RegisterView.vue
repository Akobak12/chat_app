<template>
  <img src="../assets/logo.png" class="h-32" />
  <span class="mb-6 text-4xl text-white">Code Chat</span>
  <form class="flex flex-col items-center justify-between w-1/3 h-3/5 bg-midnight-blue text-white border-8 rounded-md border-clay-purple">
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
import { ref } from 'vue'

export default {
  setup() {
    const credentials = ref({
      username: "",
      email: "",
      password: "",
    });

    const register = async () => {
      try {
        const response = await axios.post("http://127.0.0.1:3030/api/register", {
          username: credentials.value.username,
          email: credentials.value.email,
          password: credentials.value.password,
        });
        console.log("register successful:", response);
        localStorage.setItem("username", credentials.value.username)
        localStorage.setItem("email", credentials.value.email)
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