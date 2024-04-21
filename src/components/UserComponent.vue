<template>
  <section class="relative flex items-center w-full rounded-md text-white bg-midnight-puprle">
    <div class="flex items-center justify-center aspect-square h-full rounded-l-md bg-white">
      <img src="../assets/user.png" class="h-3/4" />
    </div>
    <span class="ml-3 text-lg">{{ user.username }}</span>
    <div class="flex items-center absolute bottom-1 right-3 text-sm">
      <span class="block h-4 w-4 mr-2 bg-green-600 rounded-full" v-if="user.online"></span>
      <span class="text-gray-400">{{user.lastActive}}</span>
    </div>
    <img src="../assets/mute.png" v-if="mute" @click="$emit('muteUser')" class="absolute right-2 hover:cursor-pointer" :class="{'filter: brightness-50': !user.muted}">
    <div v-if="request" class="absolute flex h-3/5 right-4">
      <button class="aspect-square bg-red-600 rounded-xl mr-2" @click="rejectFriendRequest">&#10005;</button>
      <button class="aspect-square bg-green-600 rounded-xl" @click="acceptFriendRequest">&#10003;</button>
    </div>
  </section>
</template>

<script>
import { inject } from 'vue';
import axios from 'axios';
export default {
  props: {
    user: Object,
    mute: Boolean,
    request: Boolean,
    id: Number
  },
  setup(props, {emit}) {
    const acceptFriend = inject("acceptFriend")
    const rejectFriend = inject("rejectFriend")
    const acceptFriendRequest = async () => {
      try {
        const response = await axios.post(acceptFriend, {
        id: props.user.id
      },{ withCredentials:true });
      emit("reloadRequest")
      location.reload()
      console.log("Accept successful:", response);        
      } 
      catch (error) {
        console.error("Accept failed:", error);
      }     
    }

  const rejectFriendRequest = async () => {
    try {
        const response = await axios.post(rejectFriend, {
        id: props.user.id
      },{ withCredentials:true });
      emit("reloadRequest")
      location.reload()
      console.log("Reject successful:", response);        
      } 
      catch (error) {
        console.error("Reject failed:", error);
      }   
  }

    return { acceptFriendRequest, rejectFriendRequest }
  },
};
</script>

<style>
</style>