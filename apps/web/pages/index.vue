<template>
  <div class="flex flex-col min-h-screen">
    <!-- Top Global Navigation Bar -->
    <header class="sticky top-0 z-50 flex items-center justify-between h-14 px-4 bg-white border-b border-gray-200 shadow-sm">
      <div class="flex items-center space-x-2">
        <div class="flex items-center justify-center w-10 h-10 bg-brand-600 rounded-full text-white font-bold text-2xl shadow">
          f
        </div>
        <div class="relative hidden sm:block">
          <input
            type="text"
            placeholder="Search Global Network"
            class="w-60 pl-9 pr-4 py-2 bg-gray-100 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
          <span class="absolute left-3 top-2.5 text-gray-400 text-sm">🔍</span>
        </div>
      </div>

      <!-- Navigation Icons -->
      <nav class="flex items-center space-x-1 sm:space-x-8 text-gray-500">
        <button class="p-2 sm:px-6 py-2 border-b-4 border-brand-600 text-brand-600 hover:bg-gray-100 rounded-lg">
          🏠
        </button>
        <button class="p-2 sm:px-6 py-2 hover:bg-gray-100 rounded-lg">
          👥
        </button>
        <button class="p-2 sm:px-6 py-2 hover:bg-gray-100 rounded-lg">
          📺
        </button>
        <button class="p-2 sm:px-6 py-2 hover:bg-gray-100 rounded-lg">
          🏪
        </button>
        <button class="p-2 sm:px-6 py-2 hover:bg-gray-100 rounded-lg">
          👨‍👩‍👦
        </button>
      </nav>

      <!-- User Menu -->
      <div class="flex items-center space-x-2">
        <button class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center hover:bg-gray-300">
          💬
        </button>
        <button class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center hover:bg-gray-300">
          🔔
        </button>
        <div class="w-10 h-10 bg-brand-600 text-white rounded-full flex items-center justify-center font-bold">
          U
        </div>
      </div>
    </header>

    <!-- Main 3-Column Layout -->
    <main class="flex-1 flex max-w-[1920px] mx-auto w-full justify-between pt-4 px-2 sm:px-6">
      <!-- Left Sidebar: Shortcuts -->
      <aside class="hidden lg:block w-72 pr-4 space-y-1 text-sm font-medium">
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <div class="w-8 h-8 rounded-full bg-brand-600 text-white flex items-center justify-center font-semibold">U</div>
          <span>User Profile</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">👥</span>
          <span>Friends (3 pending requests)</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">👨‍👩‍👧‍👦</span>
          <span>Groups</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">🏪</span>
          <span>Marketplace</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">📺</span>
          <span>Watch</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">🕒</span>
          <span>Memories</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
          <span class="text-xl">🔖</span>
          <span>Saved Posts</span>
        </div>
      </aside>

      <!-- Center: Feed & Composer -->
      <section class="flex-1 max-w-[640px] mx-auto space-y-4">
        <!-- Stories Carousel -->
        <div class="flex space-x-2 overflow-x-auto pb-2">
          <div class="w-28 h-48 bg-white rounded-xl shadow border border-gray-200 flex-shrink-0 flex flex-col justify-between p-2 cursor-pointer hover:opacity-95">
            <span class="text-sm font-semibold">+ Create Story</span>
            <div class="w-8 h-8 bg-brand-600 rounded-full text-white flex items-center justify-center mx-auto">+</div>
          </div>
        </div>

        <!-- Post Composer Card -->
        <div class="bg-white rounded-xl shadow p-4 space-y-3">
          <div class="flex items-center space-x-2">
            <div class="w-10 h-10 rounded-full bg-brand-600 text-white flex items-center justify-center font-bold">U</div>
            <input
              v-model="newPostText"
              type="text"
              placeholder="What's on your mind?"
              class="flex-1 bg-gray-100 hover:bg-gray-200 rounded-full px-4 py-2 text-sm focus:outline-none"
            />
          </div>
          <hr class="border-gray-200" />
          <div class="flex justify-around text-gray-500 font-semibold text-sm">
            <button class="flex items-center space-x-1 hover:bg-gray-100 p-2 rounded-lg">
              <span>📹</span>
              <span>Live video</span>
            </button>
            <button class="flex items-center space-x-1 hover:bg-gray-100 p-2 rounded-lg">
              <span>🖼️</span>
              <span>Photo/video</span>
            </button>
            <button class="flex items-center space-x-1 hover:bg-gray-100 p-2 rounded-lg">
              <span>😊</span>
              <span>Feeling/activity</span>
            </button>
          </div>
        </div>

        <!-- Sample Live Feed Post -->
        <article class="bg-white rounded-xl shadow p-4 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <div class="w-10 h-10 rounded-full bg-blue-500 text-white flex items-center justify-center font-bold">N</div>
              <div>
                <h4 class="font-bold text-sm">Network Announcement</h4>
                <p class="text-xs text-gray-500">10 mins ago • 🌐 Public</p>
              </div>
            </div>
            <button class="text-gray-400 hover:text-gray-600">•••</button>
          </div>
          <p class="text-sm leading-relaxed">
            Welcome to the Global Social Network. Microservices are live, PostgreSQL and Redis are bound, and functional parity is operational!
          </p>
          <div class="flex items-center justify-between text-xs text-gray-500 pt-2 border-t border-gray-100">
            <span>👍 ❤️ 42 reactions</span>
            <span>12 comments • 5 shares</span>
          </div>
          <div class="flex justify-around border-t border-b border-gray-100 py-1 text-sm font-semibold text-gray-600">
            <button class="flex-1 py-1.5 hover:bg-gray-100 rounded-lg flex items-center justify-center space-x-1">
              <span>👍</span>
              <span>Like</span>
            </button>
            <button class="flex-1 py-1.5 hover:bg-gray-100 rounded-lg flex items-center justify-center space-x-1">
              <span>💬</span>
              <span>Comment</span>
            </button>
            <button class="flex-1 py-1.5 hover:bg-gray-100 rounded-lg flex items-center justify-center space-x-1">
              <span>↗️</span>
              <span>Share</span>
            </button>
          </div>
        </article>
      </section>

      <!-- Right Sidebar: Contacts & Sponsored -->
      <aside class="hidden xl:block w-72 pl-4 space-y-4 text-sm">
        <div class="bg-white p-3 rounded-xl shadow space-y-2">
          <h4 class="font-semibold text-gray-500 text-xs uppercase tracking-wider">Sponsored</h4>
          <div class="flex items-center space-x-3">
            <div class="w-20 h-20 bg-gray-200 rounded-lg flex items-center justify-center text-xs text-gray-400">Ad Creative</div>
            <div>
              <p class="font-bold text-sm">Enterprise Cloud</p>
              <p class="text-xs text-gray-500">cloud.example.com</p>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <div class="flex items-center justify-between text-gray-500">
            <h4 class="font-semibold text-xs uppercase tracking-wider">Contacts</h4>
            <span>🔍 •••</span>
          </div>
          <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-xl cursor-pointer">
            <div class="relative">
              <div class="w-8 h-8 rounded-full bg-green-600 text-white flex items-center justify-center font-bold">K</div>
              <span class="absolute bottom-0 right-0 w-2.5 h-2.5 bg-green-500 border-2 border-white rounded-full"></span>
            </div>
            <span class="font-medium">Kevin Adams</span>
          </div>
        </div>
      </aside>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const newPostText = ref('')
</script>
