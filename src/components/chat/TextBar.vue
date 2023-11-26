<template>
  <div
    class="flex justify-end w-full h-12 py-1 bg-midnight-blue border-l-8 border-clay-purple"
  >
    <input
      @keyup.enter="sendMessage"
      v-model="newMessage"
      type="text"
      placeholder="message"
      class="relative top-0 w-4/5 h-auto pl-3 bg-gengar-purple placeholder-white text-white"
    />
    <button @click="sendMessage" class="h-full mx-5">
      <img
        src="../../assets/buttons/submit.png"
        class="h-full object-contain"
      />
    </button>
  </div>
</template>

<script>
import { ref } from "vue";

export default {
  setup(props, { emit }) {
    const newMessage = ref("");

    function sendMessage() {
      if (newMessage.value.trim() !== "") {
        emit("message-sent", newMessage.value);
        newMessage.value = "";
      }
    }

    return { newMessage, sendMessage };
  },
};
</script>

<style></style>
