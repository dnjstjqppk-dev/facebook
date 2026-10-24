<template>
  <div class="bg-white rounded-lg shadow-sm border border-[#e4e6eb] p-4 mb-4">
    <div class="flex items-center space-x-3 mb-3">
      <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center font-bold text-blue-700">
        U
      </div>
      <input
        v-model="postContent"
        type="text"
        placeholder="What's on your mind?"
        class="flex-1 bg-[#f0f2f5] hover:bg-[#e4e6eb] rounded-full px-4 py-2.5 text-sm focus:outline-none cursor-pointer transition-colors"
      />
    </div>
    <hr class="border-[#e4e6eb] my-2" />
    <div class="flex justify-around text-sm font-semibold text-gray-600">
      <button class="flex items-center space-x-2 py-1.5 px-3 hover:bg-gray-100 rounded-lg transition-colors">
        <span class="text-red-500 text-lg">📹</span>
        <span>Live video</span>
      </button>
      <button class="flex items-center space-x-2 py-1.5 px-3 hover:bg-gray-100 rounded-lg transition-colors">
        <span class="text-green-500 text-lg">🖼️</span>
        <span>Photo/video</span>
      </button>
      <button class="flex items-center space-x-2 py-1.5 px-3 hover:bg-gray-100 rounded-lg transition-colors">
        <span class="text-yellow-500 text-lg">😊</span>
        <span>Feeling/activity</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const postContent = ref('')
</script>
