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
</template>

<script>
import { nextTick, ref } from "vue";
import Message from "./MessageComponent.vue";
import TextBar from "./TextBar.vue";
import ToolBox from "./ToolBox.vue";
import { onUnmounted } from "vue";

export default {
  components: {
    Message,
    TextBar,
    ToolBox,
  },

  setup() {
    const messages = ref([]);

    const websocket = new WebSocket("ws://127.0.0.1:3030/ws");

    const messageContainer = ref(null);
    const toolBox = ref(false);

    const addMessage = (newMessage) => {
      websocket.send(newMessage);
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

    websocket.onopen = () => {
      console.log("WebSocket connection established");
    };

    websocket.onmessage = (event) => {
      setMessage(event.data);
    };

    websocket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    websocket.onclose = (event) => {
      console.log("WebSocket connection closed:", event);
    };

    onUnmounted(() => {
      if (websocket.readyState === WebSocket.OPEN) {
        websocket.close();
      }
    });

    return { messages, addMessage, messageContainer, toolBox, viewTools };
  },
};
</script>
