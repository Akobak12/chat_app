<template>
  <nav class="relative flex justify-start items-center bg-clay-purple h-20 w-full
  lg:w-3/4 lg:rounded-b-lg lg:h-16">
    <ul class="flex justify-between w-3/4 h-full pl-4
    md:pl-24">
      <Button
        v-for="button in buttons"
        :key="button.id"
        :picture="button.picture"
        :route="button.route"
        :isHighlighted="button.isHighlighted"
        @mouseover="toggleHighlight(button)"
        @mouseout="toggleHighlight(button)"
      />
    </ul>
    <button class="absolute right-3" @click="logout"><img src="../../assets/buttons/logout.png" class="h-10"></button>
  </nav>
</template>

<script>
import Button from "./ButtonComponent.vue";
import { inject, ref, watch } from "vue";
import router from '@/router';
import axios from 'axios';
export default {
  components: {
    Button,
  },
  props: {
    chatId: Number
  },
  setup(props) {
    const logoutURL = inject("logout")
    const buttons = ref([
      { id: "1", picture: "chat.png", route: `/chat/${localStorage.getItem("id")}`, isHighlighted: false },
      {
        id: "2",
        picture: "settings.png",
        route: "/settings",
        isHighlighted: false,
      },
      { id: "3", picture: "file.png", route: "/homepage", isHighlighted: false },
      {
        id: "4",
        picture: "account.png",
        route: "/servers",
        isHighlighted: false,
      },
    ]);

    watch(() => props.chatId, (newValue, oldValue) => {
      if (newValue !== oldValue) {
        buttons.value[0].route = `/chat/${newValue}`;
      }
    });

    const toggleHighlight = (button) => {
      button.isHighlighted = !button.isHighlighted;
    };

    const logout = async () => {
      try {
        const response = await axios.get(logoutURL, { withCredentials: true });
        localStorage.setItem("isAuth", false)
        router.push({ name: "login" })
        console.log("Logout successful:", response);
      } catch (error) {
        console.error("Logout failed:", error);
      }
    }

    return { buttons, logout, toggleHighlight };
  },
};
</script>

<style scoped></style>
