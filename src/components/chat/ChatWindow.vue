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
import {  nextTick, ref, inject, onMounted } from "vue";
import Message from "./MessageComponent.vue";
import TextBar from "./TextBar.vue";
import ToolBox from "./ToolBox.vue";

export default {
  components: {
    Message,
    TextBar,
    ToolBox,
  },

  setup() {
    const websocket = inject("websocket");
    const messages = ref([]);

    const messageContainer = ref(null);
    const toolBox = ref(false);

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

    websocket.value.onmessage = (event) => {
      setMessage(event.data);
    };

    onMounted(() => {
      websocket.value.onmessage = (event) => {
        setMessage(event.data);
      };
    });

    return { messages, addMessage, messageContainer, toolBox, viewTools };
  },
};
</script>
