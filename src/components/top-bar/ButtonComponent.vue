<template>
  <router-link :to="route" class="flex flex-col w-24" :class="{ 'highlight-bg': (isHighlighted || this.$route.path.startsWith(route)) && darkMode == false, 
  'highlight-bg-dark': (isHighlighted || this.$route.path.startsWith(route)) && darkMode == true}">
    <img :src="require(`../../assets/buttons/${picture}`)" class="relative object-contain top-1/4 h-10 z-10">
    <img
    src="../../assets/buttons/bookmark.png"
    v-if="darkMode == false"
    class="hidden relative object-contain bottom-3 h-16 z-0 transition-transform duration-500
    lg:block" 
    :class="{ 'highlight translate-y-4': (isHighlighted || this.$route.path.startsWith(route))}">
    <img
    src="../../assets/buttons/bookmark-dark.png" 
    v-else
    class="hidden relative object-contain bottom-3 h-16 z-0 transition-transform duration-500 scale-x-75
    lg:block" 
    :class="{ 'highlight-dark translate-y-4': isHighlighted || this.$route.path.startsWith(route) }">
  </router-link>
</template>

<script>
import { inject, ref } from 'vue';
export default {
  props: {
    picture: String,
    isHighlighted: Boolean,
    route: String,
  },
  setup() {
    const darkMode = ref(inject("darkMode"))

    return {darkMode}
  }
};
</script>

<style>
  .highlight {
    filter: brightness(0) saturate(100%) invert(19%) sepia(8%) saturate(4180%) hue-rotate(201deg) brightness(99%) contrast(85%);
  }
  .highlight-bg {
    backdrop-filter: brightness(0) saturate(100%) invert(19%) sepia(8%) saturate(4180%) hue-rotate(201deg) brightness(99%) contrast(85%);
  }
  .highlight-dark {
    filter: brightness(0) saturate(100%) invert(73%) sepia(4%) saturate(0%) hue-rotate(192deg) brightness(89%) contrast(91%);
  }
  .highlight-bg-dark {
    backdrop-filter: brightness(0) saturate(100%) invert(73%) sepia(4%) saturate(0%) hue-rotate(192deg) brightness(89%) contrast(91%);
  }
</style>
