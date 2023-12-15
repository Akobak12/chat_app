<template>
  <Chat />
  <section class="flex absolute right-0">
    <UserProfile 
    class=" transition-transform duration-300" 
    :class="{'translate-x-60': !userInfo}" @openProfile="viewProfile"
    :userInfo="userInfo" />
    <div class=" h-screen w-14 z-10 bg-midnight-blue"></div>
  </section>

</template>

<script>
import Chat from "../components/chat/ChatWindow.vue"
import UserProfile from "../components/UserProfile.vue"
import { ref } from "vue"

export default {
  name: "ChatView",
  components: {
    Chat,
    UserProfile
  },

  setup() {
    const userInfo = ref(false)

    const viewProfile = () => {
      userInfo.value = !userInfo.value
    }

    return{ userInfo, viewProfile }
  }
};
</script>
