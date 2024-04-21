<template>
  <router-link :to="`/chat/${chat.id}`" class="relative flex items-center h-16 w-full rounded-md text-white bg-midnight-puprle">
    <div class="flex items-center justify-center w-16 h-16 rounded-l-md bg-white">
      <img src="../assets/user.png" class="w-3/4" />
    </div>
    <span class="ml-3 text-lg">{{ chat.username }}</span>
<!--<div class="flex items-center absolute bottom-1 right-3 text-sm">
      <span class="block h-4 w-4 mr-2 bg-green-600 rounded-full" v-if="chat.online"></span>
      <span class="text-gray-400">{{chat.lastActive}}</span>
    </div>-->
  </router-link>
</template>

<script>
export default {
  props: {
    chat: Object,
    id: Number
  },
  setup() {},
};
</script>

<style>
</style>