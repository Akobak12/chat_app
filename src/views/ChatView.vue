<template>
  <Chat />
  <UserProfile
    class="absolute -right-80 z-20 transition-transform duration-300"
    :class="{ '-translate-x-80': userInfo }"
    @openProfile="viewProfile"
    :userInfo="userInfo"
  />
</template>

<script>
import Chat from "../components/chat/Chat.vue";
import UserProfile from "../components/UserProfile.vue";
import { ref } from "vue";

export default {
  name: "ChatView",
  components: {
    Chat,
    UserProfile,
  },

  setup() {
    const userInfo = ref(false);

    const viewProfile = () => {
      userInfo.value = !userInfo.value;
    };

    return { userInfo, viewProfile };
  },
};
</script>
