<template>
  <div class="relative flex items-center p-1 bg-midnight-puprle rounded-md">
    <img src="../assets/serverlogo.png" class="h-full mr-4">
    <span class="text-white">{{ server.name }}</span>
    <img src="../assets/mute.png" v-if="mute" @click="$emit('muteServer')" class="absolute right-2 hover:cursor-pointer" :class="{'filter: brightness-50': !server.muted}">
  </div>
</template>

<script>
export default {
  props: {
    server: String,
    mute: Boolean,
  },
  setup() {},
};
</script>

<style>
</style>