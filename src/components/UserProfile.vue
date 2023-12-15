<template>
  <section class="w-64 h-screen bg-dragon-purple">
    <button
      @click="$emit('openProfile')"
      class="absolute w-8 h-9 -left-6 top-20 pl-1 bg-dragon-purple rounded-lg"
    >
      <img
        src="../assets/buttons/vector.png"
        class="scale-75"
        :class="{ 'rotate-180': userInfo }"
      />
    </button>
    <div class="px-10 text-white">
      <section class="flex flex-col items-center">
        <span class="text-2xl">{{ profile.userName }}</span>
        <div class="relative">
          <img src="../assets/user.png" class="w-full" />
          <span class="absolute right-3 bottom-3 h-6 w-6 rounded-full" :class="profile.online ? 'bg-green-600' : 'bg-red-600'"></span>
        </div>
        <span class="text-sm">{{ profile.lastActive }}</span>
      </section>
      <ul class="list-disc">
        <li class="text-left my-6 text-sm">Name: {{ profile.userName }}</li>
        <li class="text-left my-6 text-sm">
          Description: {{ profile.description }}
        </li>
        <li class="text-left my-6 text-sm">Exp: {{ profile.experience }}</li>
        <li class="text-left my-6 text-sm">joined: {{ profile.joined }}</li>
      </ul>
    </div>
  </section>
</template>

<script>
import { ref } from "vue";

export default {
  props: {
    userInfo: Boolean,
  },

  setup() {
    const profile = ref({
      userName: "pablo",
      profilePicture: "",
      description: "I am a user, with normal user hobbies",
      experience: "C#, C++, Python",
      joined: "7845",
      online: true,
      lastActive: "5656"
    });

    return { profile };
  },
};
</script>


<style>

</style>