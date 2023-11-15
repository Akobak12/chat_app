<!-- Chat.vue -->
<template>
  <Message class="bg-purple-950 text-gray-200 p-2 m-1 rounded-lg" v-for="message in messages" :key="message.id" :content="message.content" :isLeft="message.isLeft" :isRight="message.isRight" />
</template>

<script>
import { ref } from 'vue';
import Message from '@/components/Message.vue';


export default {
  components: {
    Message,
  },
  setup() {
    const messages = ref([
      { id: 1, content: 'Hello!', user: {name: 'Pepa', icon: ''}},
      { id: 2, content: 'Hi there!', },
      // Other message objects
    ]);

    return { messages };
  },
};
</script>

<style scoped>
/* Add any overall chat styling here */
</style>
