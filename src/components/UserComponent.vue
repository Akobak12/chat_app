<template>
  <section class="relative flex items-center h-16 w-full rounded-md text-white bg-midnight-puprle">
    <div class="flex items-center justify-center w-16 h-16 rounded-l-md bg-white">
      <img src="../assets/user.png" class="w-3/4" />
    </div>
    <span class="ml-3 text-lg">{{ user.userName }}</span>
    <div class="flex items-center absolute bottom-1 right-3 text-sm">
      <span class="block h-4 w-4 mr-2 bg-green-600 rounded-full" v-if="user.online"></span>
      <span class="text-gray-400">{{user.lastActive}}</span>
    </div>
  </section>
</template>

<script>
export default {
  props: {
    user: Object,
  },
  setup() {},
};
</script>

<style>
</style>