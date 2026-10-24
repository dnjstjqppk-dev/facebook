<template>
  <article class="bg-white rounded-lg shadow-sm border border-[#e4e6eb] mb-4 overflow-hidden">
    <!-- Header -->
    <div class="p-4 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-full bg-blue-500 text-white flex items-center justify-center font-bold">
          {{ authorName.charAt(0) }}
        </div>
        <div>
          <h4 class="font-semibold text-sm leading-tight text-gray-900">{{ authorName }}</h4>
          <span class="text-xs text-gray-500">{{ timeAgo }} · 🌐</span>
        </div>
      </div>
      <button class="text-gray-500 hover:bg-gray-100 rounded-full p-1.5 transition-colors">
        •••
      </button>
    </div>

    <!-- Content -->
    <div class="px-4 pb-3 text-sm text-gray-900 whitespace-pre-line">
      {{ content }}
    </div>

    <!-- Counters -->
    <div class="px-4 py-2 border-b border-[#e4e6eb] flex items-center justify-between text-xs text-gray-500">
      <div class="flex items-center space-x-1">
        <span>👍❤️</span>
        <span>{{ reactionCount }}</span>
      </div>
      <div class="space-x-3">
        <span>{{ commentCount }} comments</span>
        <span>{{ shareCount }} shares</span>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="px-2 py-1 flex justify-around text-sm font-semibold text-gray-600">
      <button
        @click="handleReact"
        class="flex-1 flex items-center justify-center space-x-2 py-2 hover:bg-gray-100 rounded-lg transition-colors"
        :class="{ 'text-blue-600': isLiked }"
      >
        <span>👍</span>
        <span>Like</span>
      </button>
      <button class="flex-1 flex items-center justify-center space-x-2 py-2 hover:bg-gray-100 rounded-lg transition-colors">
        <span>💬</span>
        <span>Comment</span>
      </button>
      <button class="flex-1 flex items-center justify-center space-x-2 py-2 hover:bg-gray-100 rounded-lg transition-colors">
        <span>↗️</span>
        <span>Share</span>
      </button>
    </div>
  </article>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    authorName?: string
    timeAgo?: string
    content: string
    reactionCount?: number
    commentCount?: number
    shareCount?: number
  }>(),
  {
    authorName: 'Alex Mercer',
    timeAgo: '3h ago',
    reactionCount: 142,
    commentCount: 28,
    shareCount: 6
  }
)

const isLiked = ref(false)
const reactionCount = ref(props.reactionCount)

const handleReact = () => {
  isLiked.value = !isLiked.value
  if (isLiked.value) {
    reactionCount.value++
  } else {
    reactionCount.value--
  }
}
</script>
