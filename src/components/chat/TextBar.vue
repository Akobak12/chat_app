<template>
  <div
    class="flex justify-end w-full h-12 py-1 bg-midnight-blue border-clay-purple
    lg:border-l-8"
    :class="{ 'bg-[#797979]': darkMode }"
  >
    <img src="../../assets/call.png" class="absolute scale-75 left-2">
    <input
      @keyup.enter="sendMessage"
      v-model="newMessage"
      type="text"
      placeholder="message"
      class="relative top-0 w-7/12 h-auto pl-3 bg-white text-black
      sm:w-8/12 md:w-9/12 lg:w-4/5"
      :class="{ 'bg-[#2c2c2c] text-white': darkMode }"
    />
    <button @click="sendMessage" class="h-full mx-5">
      <img
        src="../../assets/buttons/submit.png"
        class="h-full object-contain bg-clay-purple rounded-full"
        :class="{ 'bg-[#c9c9c9]': darkMode }"/>
    </button>
  </div>
</template>

<script>
import { inject, ref } from "vue";

export default {
  setup(props, { emit }) {
    const newMessage = ref("");
    const darkMode = ref(inject("darkMode"))

    function sendMessage() {
      if (newMessage.value.trim() !== "") {
        emit("message-sent", newMessage.value);
        newMessage.value = "";
      }
    }

    return { newMessage, darkMode, sendMessage };
  },
};
</script>

<style></style>
