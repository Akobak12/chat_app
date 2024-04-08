<template>
  <section 
    v-if="!call"
    class="grow w-11/12 mt-8 overflow-auto bg-midnight-blue border-8 border-clay-purple
    lg:w-3/4 lg:mt-12"
    ref="messageContainer"
  >
    <Message
      class="m-4"
      v-for="message in messages"
      :key="message.id"
      :content="message.content"
    />
  </section>

  <section 
    v-if="call"
    class="flex grow relative justify-center w-11/12 mt-8 overflow-auto bg-[#352057] border-8 border-dragon-purple
    lg:w-3/4 lg:mt-12"
    :class="{ 'bg-[#353535]': darkMode }"
  >
    <div class="flex absolute justify-between px-2 w-3/4 bottom-0 bg-[#43335D] border-2 border-dragon-purple rounded-2xl
    md:w-1/2 lg:w-1/4"
    :class="{ 'bg-[#454545]': darkMode }">
      <button><img class="h-10" src="../assets/mic.png"></button>
      <button class="flex justify-center scale-75 rounded-full h-14 aspect-square bg-red-600"><img src="../assets/endcall.png"></button>
      <button @click="deafen=!deafen">
        <img v-if="!deafen" class="w-10" src="../assets/deaf.png">
        <img v-else class="w-8 scale-75 mr-2" src="../assets/undeaf.png">
      </button>
    </div>
  </section>

  <section class="relative bottom-0 w-full h-24 mt-5
  lg:w-3/4">
    <ToolBox
      class="transition-transform duration-300"
      :class="{ '-translate-y-24': toolBox }"
      @openTools="viewTools"
      :toolBox="toolBox"
    />
    <TextBar class="absolute bottom-0" @message-sent="addMessage" />
  </section>


  <div class="absolute flex -right-64 z-30">
    <section class="flex">
      <section :class="{'-translate-x-60': userInfo}" class="w-64 h-screen bg-dragon-purple transition-transform duration-300">
        <button
          @click="viewProfile"
          class="absolute w-8 h-9 -left-6 top-20 pl-1 bg-dragon-purple rounded-lg"
        >
          <img
            src="../assets/buttons/vector.png"
            class="scale-75"
            :class="{ 'rotate-180': userInfo, 'highlight-dark': darkMode }"
          />
        </button>
        <div class="px-10 text-white">
        <section class="flex flex-col items-center">
          <span class="text-2xl">{{ profile.userName }}</span>
          <div class="relative">
            <img src="../assets/user.png" class="w-full" />
            <span class="absolute right-3 bottom-3 h-6 w-6 rounded-full" :class="profile.online ? 'bg-green-600' : 'bg-red-600'"></span>
          </div>
          <span class="text-sm">{{ profile.lastActive }}</span>
        </section>
        <ul class="list-disc">
          <li class="text-left my-6 text-sm">Name: {{ profile.userName }}</li>
          <li class="text-left my-6 text-sm">
            Description: {{ profile.description }}
          </li>
          <li class="text-left my-6 text-sm">Exp: {{ profile.experience }}</li>
          <li class="text-left my-6 text-sm">joined: {{ profile.joined }}</li>
        </ul>
      </div>
      </section>
    </section>
    <div class="hidden h-screen w-14 z-10 bg-midnight-blue -translate-x-60
    lg:block"></div>
  </div>

</template>

<script>
import { ref, onUnmounted, nextTick, onMounted, inject } from "vue";
import { onBeforeRouteLeave } from 'vue-router'
import Message from "../components/chat/MessageComponent.vue";
import TextBar from "../components/chat/TextBar.vue";
import ToolBox from "../components/chat/ToolBox.vue";

export default {
  components: {
    Message,
    TextBar,
    ToolBox,
  },

  setup() {
    const userInfo = ref(false)
    const joinRoom = inject("joinRoom")
    const call = ref(true)
    const deafen = ref(false)
    const darkMode = ref(inject("darkMode"))

    const websocket = ref(null);

    const openWebsocket = () => {
      websocket.value = new WebSocket(joinRoom.value);
      console.log(websocket.value)
      
      websocket.value.onopen = () => {
        console.log("WebSocket connection established");
      };

      websocket.value.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      websocket.value.onclose = (event) => {
        console.log("WebSocket connection closed:", event);
      };
      
      websocket.value.onmessage = (event) => {
        setMessage(event.data);
      };
    }

    onMounted(() => {
      openWebsocket()
    });

    const closeWebSocket = () => {
      if (websocket.value && websocket.value.readyState === WebSocket.OPEN) {
        websocket.value.close();
      }
    };

    onUnmounted(() => {
      closeWebSocket()
    });

    onBeforeRouteLeave(() => {
      closeWebSocket()    
    })

    const messages = ref([]);

    const messageContainer = ref(null);
    const toolBox = ref(false);

    const profile = ref({
      userName: "pablo",
      profilePicture: "",
      description: "I am a user, with normal user hobbies",
      experience: "C#, C++, Python",
      joined: "7845",
      online: true,
      lastActive: "5656"
    });

    const viewProfile = () => {
      userInfo.value = !userInfo.value
    }

    const addMessage = (newMessage) => {
      websocket.value.send(newMessage);
      scrollToBottom();
    };

    const setMessage = (newMessage) => {
      messages.value.push({
        id: messages.value.length + 1,
        content: newMessage,
      });
      scrollToBottom();
    };

    const scrollToBottom = () => {
      nextTick(() => {
        if (messageContainer.value) {
          messageContainer.value.scrollTop =
            messageContainer.value.scrollHeight;
        }
      });
    };

    const viewTools = () => {
      toolBox.value = !toolBox.value;
    };

    

    return { messages, messageContainer, toolBox, addMessage, viewTools, userInfo, viewProfile, profile, call, deafen, darkMode };

  }
};
</script>
