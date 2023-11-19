<template>
  <section class="w-96 h-screen bg-midnight-blue border-l-12 border-dragon-purple">
    <button @click="$emit('openProfile')" class="absolute w-8 h-9 -left-8 top-20 pl-1 bg-dragon-purple rounded-lg">
      <img src="../assets/buttons/vector.png" class="scale-75" :class="{'rotate-180': userInfo}">
    </button>
  </section>
</template>

<script>
export default {
  props: {
    userInfo: String,
  },

}
</script>

<style>

</style>