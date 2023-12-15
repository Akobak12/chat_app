<template>
  <div class="flex items-center space-x-4 text-gray-200">
    <img class="w-12 h-12 self-start" src="../../assets/user.png" alt="" />
    <p class="w-11/12 break-words text-left">{{ content }}</p>
  </div>
</template>

<script>
export default {
  props: {
    content: String,
  },
};
</script>
