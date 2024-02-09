<template>
  <div class="flex w-full mb-4 pr-2 rounded-xl bg-white text-black">
    <input type="text" :placeholder="placeholder" class="w-full bg-transparent outline-none h-full pl-4">
    <button class="search-button">
      <img src="../assets/search.png" class="h-">
    </button>
  </div>
</template>

<script>
export default {
  props: {
    placeholder: String
  }
}
</script>

<style>

</style>