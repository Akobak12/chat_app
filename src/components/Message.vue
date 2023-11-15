<!-- Message.vue -->
<template>
  <div class="w-max flex items-center">
    <img class="w-12 h-12" src="../assets/user.png" alt="">
    <p>{{ content }}</p>
  </div>
</template>

<script>

export default {
  props: {
    content: String,
    isLeft: Boolean,
    isRight: Boolean,
    // Other props for additional details
  },
};
</script>

<style scoped>
/* Add any specific styles for the message component here */
.message-left {
  text-align: left;
  /* Other styles for messages on the left side */
}

.message-right {
  text-align: right;
  /* Other styles for messages on the right side */
}
</style>
