<template>
  <section class="relative flex items-center w-full rounded-md text-white bg-midnight-puprle">
    <div class="flex items-center justify-center aspect-square h-full rounded-l-md bg-white">
      <img src="../assets/user.png" class="h-3/4" />
    </div>
    <span class="ml-3 text-lg">{{ user.username }}</span>
    <div class="flex items-center absolute bottom-1 right-3 text-sm">
      <span class="block h-4 w-4 mr-2 bg-green-600 rounded-full" v-if="user.online"></span>
      <span class="text-gray-400">{{user.lastActive}}</span>
    </div>
    <img src="../assets/mute.png" v-if="mute" @click="$emit('muteUser')" class="absolute right-2 hover:cursor-pointer" :class="{'filter: brightness-50': !user.muted}">
  </section>
</template>

<script>
export default {
  props: {
    user: Object,
    mute: Boolean,
    id: Number
  },
  setup() {

  },
};
</script>

<style>
</style>