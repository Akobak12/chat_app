<template>
  <section
    class="flex w-full h-96 bg-dragon-purple border-x-8 border-t-8 border-clay-purple"
  >
    <button
      @click="$emit('openTools')"
      class="relative self-start w-9 h-9 -top-5 left-6 pl-1 bg-dragon-purple rounded-lg rotate-90"
    >
      <img
        src="../../assets/buttons/vector.png"
        class="scale-75"
        :class="{ 'rotate-180': toolBox }"
      />
    </button>
  </section>
</template>

<script>
export default {
  props: {
    toolBox: Boolean,
  },
};
</script>

<style></style>
