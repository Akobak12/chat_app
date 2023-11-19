<template>
  <section class="grow w-3/4 overflow-auto bg-midnight-blue border-8 border-clay-purple" ref="messageContainer">
    <Message
      class="m-4"
      v-for="message in messages"
      :key="message.id"
      :content="message.content"
    />
  </section>
  <section class="relative bottom-0 w-3/4 h-24 mt-5">
    <ToolBox class="transition-transform duration-300" :class="{'-translate-y-24': toolBox}" @openTools="viewTools" :toolBox="toolBox"/>
    <TextBar class="absolute bottom-0" @message-sent="addMessage" />
  </section>
</template>

<script>
import { nextTick, ref } from "vue";
import Message from "./Message.vue";
import TextBar from "./TextBar.vue"
import ToolBox from "./ToolBox.vue"

export default {
  components: {
    Message,
    TextBar,
    ToolBox
  },
  setup() {
    const messages = ref([
      { id: 1, content: "Hello there" },
      { id: 2, content: "General Kenobi!" },
    ]);

    const messageContainer = ref(null);
    const toolBox = ref(false)

    const addMessage = (newMessage) => {
      messages.value.push({ id: messages.value.length + 1, content: newMessage });
      scrollToBottom()
    };

    const scrollToBottom = () => {
      nextTick(() => {
        if (messageContainer.value) {
          messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
        }
      });
    }

    const viewTools = () => {
      toolBox.value = !toolBox.value
    }

    return { messages, addMessage, messageContainer, toolBox, viewTools };
  },
};
</script>

