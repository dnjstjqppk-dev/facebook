<template>
  <div class="min-h-screen bg-[#f0f2f5]">
    <Navbar />

    <main class="max-w-7xl mx-auto px-2 sm:px-4 pt-4 grid grid-cols-1 md:grid-cols-4 lg:grid-cols-4 gap-6">
      <!-- Left Sidebar (Desktop Only) -->
      <aside class="hidden md:block col-span-1 space-y-2 text-sm font-medium text-gray-700">
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <div class="w-9 h-9 rounded-full bg-blue-100 flex items-center justify-center font-bold text-blue-700">
            Me
          </div>
          <span>User Profile</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <span class="text-xl">👥</span>
          <span>Friends</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <span class="text-xl">🚩</span>
          <span>Pages</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <span class="text-xl">🌐</span>
          <span>Groups</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <span class="text-xl">🏪</span>
          <span>Marketplace</span>
        </div>
        <div class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors">
          <span class="text-xl">🔖</span>
          <span>Saved</span>
        </div>
      </aside>

      <!-- Center Main Feed Column -->
      <section class="col-span-1 md:col-span-3 lg:col-span-2 space-y-4">
        <!-- Stories Carousel Bar -->
        <div class="flex space-x-2 overflow-x-auto pb-2 scrollbar-none">
          <div class="flex-shrink-0 w-28 h-48 bg-white border border-[#e4e6eb] rounded-xl overflow-hidden shadow-sm flex flex-col justify-between p-2 relative cursor-pointer hover:opacity-95 transition">
            <div class="w-full h-32 bg-gray-200 rounded-lg flex items-center justify-center text-gray-400">
              📸
            </div>
            <div class="text-center font-semibold text-xs text-gray-800">
              Create Story
            </div>
          </div>
          <div
            v-for="i in 4"
            :key="i"
            class="flex-shrink-0 w-28 h-48 bg-gradient-to-b from-blue-400 to-indigo-600 rounded-xl overflow-hidden shadow-sm flex flex-col justify-between p-2 relative cursor-pointer hover:scale-[1.02] transition"
          >
            <div class="w-9 h-9 rounded-full border-2 border-white bg-blue-100 flex items-center justify-center font-bold text-xs text-blue-700">
              U{{ i }}
            </div>
            <span class="text-white font-bold text-xs drop-shadow">Friend {{ i }}</span>
          </div>
        </div>

        <!-- Post Creation Box -->
        <FeedComposer />

        <!-- Posts Feed Stream -->
        <PostCard
          author-name="Sarah Jenkins"
          time-ago="25m ago"
          content="Super excited to announce the launch of our new distributed architecture on Go and Vue 3! Built for global scale with zero compromise on speed or security. 🚀✨"
          :reaction-count="318"
          :comment-count="45"
          :share-count="12"
        />

        <PostCard
          author-name="Nova Engineering Team"
          time-ago="2h ago"
          content="Welcome to the functional parity social platform! Realtime messaging, dynamic news feeds, privacy-first controls, and high-concurrency microservices are now active."
          :reaction-count="892"
          :comment-count="134"
          :share-count="88"
        />
      </section>

      <!-- Right Contacts & Birthdays Column (Large Screen Only) -->
      <aside class="hidden lg:block col-span-1 space-y-4 text-sm">
        <div class="border-b border-[#e4e6eb] pb-3">
          <h4 class="font-semibold text-gray-500 mb-2">Birthdays</h4>
          <div class="flex items-center space-x-2 text-gray-800">
            <span class="text-xl">🎁</span>
            <span class="text-xs"><strong>John Doe</strong> and 2 others have birthdays today.</span>
          </div>
        </div>

        <div>
          <div class="flex items-center justify-between text-gray-500 font-semibold mb-2">
            <span>Contacts</span>
            <div class="space-x-2">
              <button class="hover:text-gray-700">🔍</button>
              <button class="hover:text-gray-700">•••</button>
            </div>
          </div>
          <div class="space-y-1">
            <div
              v-for="i in 6"
              :key="i"
              class="flex items-center space-x-3 p-2 hover:bg-gray-200 rounded-lg cursor-pointer transition-colors"
            >
              <div class="relative">
                <div class="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center font-bold text-xs text-gray-700">
                  C{{ i }}
                </div>
                <span class="absolute bottom-0 right-0 w-2.5 h-2.5 bg-green-500 border-2 border-white rounded-full"></span>
              </div>
              <span class="font-medium text-gray-800">Contact {{ i }}</span>
            </div>
          </div>
        </div>
      </aside>
    </main>
  </div>
</template>

<script setup lang="ts">
// Feed page logic
</script>
