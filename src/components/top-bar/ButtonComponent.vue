<template>
  <router-link :to="route" class="flex flex-col w-24" :class="{ 'highlight-bg': isHighlighted || this.$route.path.startsWith(route) }">
    <img :src="require(`../../assets/buttons/${picture}`)" class="relative object-contain top-1/4 h-10 z-10">
    <img
    src="../../assets/buttons/bookmark.png" 
    class="relative object-contain bottom-3 h-16 z-0 transition-transform duration-500" 
    :class="{ 'highlight translate-y-4': isHighlighted || this.$route.path.startsWith(route) }">
  </router-link>
</template>

<script>
export default {
  props: {
    picture: String,
    isHighlighted: Boolean,
    route: String,
  },
};
</script>

<style>
  .highlight {
    filter: brightness(0) saturate(100%) invert(19%) sepia(8%) saturate(4180%) hue-rotate(201deg) brightness(99%) contrast(85%);
  }
  .highlight-bg {
    backdrop-filter: brightness(0) saturate(100%) invert(19%) sepia(8%) saturate(4180%) hue-rotate(201deg) brightness(99%) contrast(85%);
  }
</style>
