<template>
  <img src="../assets/logo.png" class="h-32" />
  <span class="mb-6 text-4xl text-white">Code Chat</span>
  <form class="flex flex-col items-center justify-between w-1/3 h-3/5 bg-midnight-blue text-white border-8 rounded-md border-clay-purple">
    <span class="text-4xl mt-4 mb-8">LOGIN</span>
    <input v-model="credentials.email" type="email" placeholder="Email" />
    <input v-model="credentials.password" type="password" placeholder="Password"/>
    <button @click.prevent="login" class="w-48 h-12 text-2xl rounded-lg bg-dragon-purple">Continue</button>
    <div class="flex mt-6">
      <router-link to="/register" class="text-sm mr-8"><u>Create an account</u></router-link>
      <a href="" class="text-sm ml-8"><u>Forgot your password?</u></a>
    </div>
  </form>
</template>

<script>
import axios from "axios";
import { inject, onMounted, ref } from 'vue'
import router from '@/router';

export default {
  setup() {
    const loginURL = inject("login")

    const credentials = ref({
      email: "",
      password: "",
    });

    const login = async () => {
      try {
        const response = await axios.post(loginURL, {
          email: credentials.value.email,
          password: credentials.value.password,
        });
        console.log("Login successful:", response);
        localStorage.setItem("token", response.data.token)
        localStorage.setItem("isAuth", true)
        localStorage.setItem("userID", response.data.id)
        localStorage.setItem("username", response.data.username)

        console.log("bruh", localStorage.getItem("username"))
        router.push({ name: "homepage" })
      } catch (error) {
        console.error("Login failed:", error);
      }
    };

    onMounted(() => {
      if (localStorage.getItem("isAuth") == "true") {
        router.push({ name: "homepage" });
      }
    })

    return { login, credentials };
  },
};
</script>


<style scoped>
input {
  height: 15%;
  width: 90%;
  margin-bottom: 1.5rem;
  color: black;
  font-size: 1.25rem;
  line-height: 1.75rem;
  padding-left: 0.75rem;
  border-radius: 0.375rem;
}
</style>