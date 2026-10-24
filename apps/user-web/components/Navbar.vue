<template>
  <header class="sticky top-0 z-50 bg-white border-b border-[#e4e6eb] h-14 px-4 flex items-center justify-between shadow-sm">
    <!-- Left: Brand Logo & Search Bar -->
    <div class="flex items-center space-x-2">
      <NuxtLink to="/feed" class="flex items-center space-x-2">
        <div class="w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold text-2xl shadow">
          N
        </div>
      </NuxtLink>
      <div class="relative hidden sm:block">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search Nova..."
          class="bg-[#f0f2f5] text-sm rounded-full pl-9 pr-4 py-2 w-60 focus:w-72 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all"
        />
        <svg class="w-4 h-4 text-gray-500 absolute left-3 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
      </div>
    </div>

    <!-- Center: Primary Navigation Tabs -->
    <nav class="hidden md:flex items-center h-full space-x-1 lg:space-x-8">
      <NuxtLink
        to="/feed"
        class="h-full flex items-center px-6 border-b-4 border-blue-600 text-blue-600 font-semibold"
      >
        <span class="text-xl">🏠</span>
      </NuxtLink>
      <NuxtLink
        to="/friends"
        class="h-full flex items-center px-6 border-b-4 border-transparent text-gray-500 hover:bg-gray-100 transition-colors"
      >
        <span class="text-xl">👥</span>
      </NuxtLink>
      <NuxtLink
        to="/watch"
        class="h-full flex items-center px-6 border-b-4 border-transparent text-gray-500 hover:bg-gray-100 transition-colors"
      >
        <span class="text-xl">📺</span>
      </NuxtLink>
      <NuxtLink
        to="/marketplace"
        class="h-full flex items-center px-6 border-b-4 border-transparent text-gray-500 hover:bg-gray-100 transition-colors"
      >
        <span class="text-xl">🏪</span>
      </NuxtLink>
      <NuxtLink
        to="/groups"
        class="h-full flex items-center px-6 border-b-4 border-transparent text-gray-500 hover:bg-gray-100 transition-colors"
      >
        <span class="text-xl">🌐</span>
      </NuxtLink>
    </nav>

    <!-- Right: Quick Action Controls & User Avatar -->
    <div class="flex items-center space-x-2">
      <button class="w-10 h-10 rounded-full bg-[#e4e6eb] hover:bg-gray-300 flex items-center justify-center font-bold text-gray-700">
        💬
      </button>
      <button class="w-10 h-10 rounded-full bg-[#e4e6eb] hover:bg-gray-300 flex items-center justify-center font-bold text-gray-700">
        🔔
      </button>
      <NuxtLink to="/profile" class="w-10 h-10 rounded-full bg-blue-100 border border-blue-300 flex items-center justify-center overflow-hidden">
        <span class="text-sm font-semibold text-blue-700">Me</span>
      </NuxtLink>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const searchQuery = ref('')
</script>
