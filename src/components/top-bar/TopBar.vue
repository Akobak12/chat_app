<template>
  <nav class="flex justify-start bg-clay-purple w-3/4 h-16 rounded-b-lg">
    <ul class="flex justify-between w-3/4 h-full pl-24">
      <Button
        v-for="button in buttons"
        :key="button.id"
        :picture="button.picture"
        :route="button.route"
        :isHighlighted="button.isHighlighted"
        @mouseover="toggleHighlight(button)"
        @mouseout="toggleHighlight(button)"
      />
    </ul>
  </nav>
</template>

<script>
import Button from "./Button.vue";
import { ref } from "vue";

export default {
  components: {
    Button,
  },
  setup() {
    const buttons = ref([
      { id: "1", picture: "chat.png", route: "/", isHighlighted: false },
      {
        id: "2",
        picture: "settings.png",
        route: "/settings",
        isHighlighted: false,
      },
      { id: "3", picture: "file.png", route: "/friends", isHighlighted: false },
      {
        id: "4",
        picture: "account.png",
        route: "/profile",
        isHighlighted: false,
      },
      { id: "5", picture: "add.png", route: "/servers", isHighlighted: false },
    ]);

    const toggleHighlight = (button) => {
      button.isHighlighted = !button.isHighlighted;
    };

    return { buttons, toggleHighlight };
  },
};
</script>

<<<<<<< HEAD
<style scoped></style>
=======
<style scoped></style>
>>>>>>> 34f5006 (server)
