<template>
  <section
    class="grow w-3/4 overflow-auto bg-midnight-blue border-8 border-clay-purple"
    ref="messageContainer"
  >
    <Message
      class="m-4"
      v-for="message in messages"
      :key="message.id"
      :content="message.content"
    />
  </section>
  <section class="relative bottom-0 w-3/4 h-24 mt-5">
    <ToolBox
      class="transition-transform duration-300"
      :class="{ '-translate-y-24': toolBox }"
      @openTools="viewTools"
      :toolBox="toolBox"
    />
    <TextBar class="absolute bottom-0" @message-sent="addMessage" />
  </section>
  <div class="absolute flex -right-[23rem]">
    <section class="flex">
      <section class="w-96 h-screen bg-dragon-purple transition-transform duration-300" :class="{'-translate-x-[23rem]': members}">
        <button
          @click="viewProfile"
          class="absolute w-8 h-9 -left-6 top-20 pl-1 bg-dragon-purple rounded-lg"
        >
          <img
            src="../assets/buttons/vector.png"
            class="scale-75"
            :class="{ 'rotate-180': members }"
          />
        </button>
        <div class="flex flex-col h-full w-full p-4 text-3xl text-white">
          <h1 class="pb-3">Server 1</h1>
          <section class="flex flex-col flex-1 items-center w-full bg-midnight-blue">
            <h1 class="py-3">Users</h1>
            <hr class="bg-clay-purple h-2 border-0 w-11/12 rounded-full">
            <div class="w-full pt-3 px-6">
              <User v-for="user in users" :user="user" :id="user.id" :key="user.id" @click="selectUser(user)" class="mb-2 h-16" :class="{ 'border-clay-purple border-2': user.selected }"/>
            </div>
            <button class="absolute bottom-12 text-lg w-1/2 h-12 bg-dragon-purple rounded-xl">Invite friends</button>
          </section>
        </div>
      </section>
    </section>
    <div class="h-screen w-14 z-10 bg-midnight-blue -translate-x-[23rem]"></div>
  </div>
  

</template>

<script>
import { ref, onUnmounted, nextTick, onMounted, inject } from "vue";
import { onBeforeRouteLeave } from 'vue-router'
import Message from "../components/chat/MessageComponent.vue";
import TextBar from "../components/chat/TextBar.vue";
import ToolBox from "../components/chat/ToolBox.vue";
import User from "../components/UserComponent.vue";

export default {
  components: {
    Message,
    TextBar,
    ToolBox,
    User
  },

  setup() {
    const members = ref(false)
    const joinRoom = inject("joinRoom")

    const messages = ref([]);

    const messageContainer = ref(null);
    const toolBox = ref(false);

    const users = ref([
      {id: "3", username: "idk", lastActive: "1991", online: false, selected: false},
      {id: "4", username: "bro", lastActive: "Active now", online: true, selected: false}
    ])

    const websocket = ref(null);

    const openWebsocket = () => {
      websocket.value = new WebSocket(joinRoom.value);
      
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

    const viewProfile = () => {
      members.value = !members.value
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

    

    return { messages, messageContainer, toolBox, addMessage, viewTools, members, viewProfile, users };

  }
};
</script>
