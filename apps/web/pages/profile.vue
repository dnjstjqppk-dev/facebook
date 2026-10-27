<template>
  <div class="min-h-screen bg-[#f0f2f5]">
    <Navbar />

    <!-- Profile Header Banner Area -->
    <div class="bg-white border-b border-[#e4e6eb] shadow-sm">
      <div class="max-w-6xl mx-auto">
        <!-- Cover Photo -->
        <div class="h-64 sm:h-80 md:h-96 w-full bg-gradient-to-r from-blue-700 via-indigo-600 to-purple-800 rounded-b-xl relative overflow-hidden flex items-end justify-end p-4">
          <button class="bg-white/90 hover:bg-white text-gray-800 font-semibold text-xs sm:text-sm px-4 py-2 rounded-lg shadow flex items-center space-x-2 transition">
            <span>📷</span>
            <span>Edit cover photo</span>
          </button>
        </div>

        <!-- Avatar, Identity & Header Action Buttons -->
        <div class="px-6 pb-4">
          <div class="flex flex-col md:flex-row items-center md:items-end justify-between -mt-16 md:-mt-10 mb-4 gap-4">
            <!-- Avatar & Titles -->
            <div class="flex flex-col md:flex-row items-center md:items-end space-y-3 md:space-y-0 md:space-x-6 text-center md:text-left">
              <div class="w-36 h-36 sm:w-44 sm:h-44 rounded-full border-4 border-white bg-blue-100 flex items-center justify-center font-bold text-5xl text-blue-700 shadow-md relative overflow-hidden">
                <span>AM</span>
              </div>
              <div class="mb-2">
                <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 leading-tight">Alex Mercer</h1>
                <p class="text-sm text-gray-500 font-medium">1.2K friends · 84 mutual</p>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex items-center space-x-3 mb-2">
              <button class="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2 rounded-lg text-sm flex items-center space-x-2 transition shadow-sm">
                <span>➕</span>
                <span>Add to story</span>
              </button>
              <button
                @click="isEditingBio = true"
                class="bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold px-4 py-2 rounded-lg text-sm flex items-center space-x-2 transition"
              >
                <span>✏️</span>
                <span>Edit profile</span>
              </button>
              <button class="bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold px-3 py-2 rounded-lg text-sm transition">
                •••
              </button>
            </div>
          </div>

          <hr class="border-[#e4e6eb] my-1" />

          <!-- Profile Nav Tabs -->
          <nav class="flex space-x-1 sm:space-x-4 text-sm font-semibold text-gray-600 overflow-x-auto">
            <button class="py-3 px-3 border-b-4 border-blue-600 text-blue-600">Posts</button>
            <button class="py-3 px-3 border-b-4 border-transparent hover:bg-gray-100 rounded-t-lg transition">About</button>
            <button class="py-3 px-3 border-b-4 border-transparent hover:bg-gray-100 rounded-t-lg transition">Friends</button>
            <button class="py-3 px-3 border-b-4 border-transparent hover:bg-gray-100 rounded-t-lg transition">Photos</button>
            <button class="py-3 px-3 border-b-4 border-transparent hover:bg-gray-100 rounded-t-lg transition">Videos</button>
            <button class="py-3 px-3 border-b-4 border-transparent hover:bg-gray-100 rounded-t-lg transition">More ▾</button>
          </nav>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <main class="max-w-6xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-12 gap-6">
      <!-- Left Rail (Intro, Photos, Friends previews) -->
      <aside class="md:col-span-5 space-y-4">
        <!-- Intro Card -->
        <div class="bg-white p-4 rounded-lg shadow-sm border border-[#e4e6eb] space-y-3">
          <h3 class="font-bold text-lg text-gray-900">Intro</h3>
          <p class="text-sm text-gray-700 text-center py-1">
            {{ bioText }}
          </p>
          <button
            @click="isEditingBio = true"
            class="w-full bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-1.5 rounded-lg text-sm transition"
          >
            Edit bio
          </button>
          <hr class="border-[#e4e6eb]" />
          <div class="space-y-2 text-sm text-gray-700">
            <div class="flex items-center space-x-2">
              <span class="text-gray-500">💼</span>
              <span>Software Architect at <strong>Nova Labs</strong></span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-gray-500">🎓</span>
              <span>Studied Computer Science at <strong>Tech University</strong></span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-gray-500">🏠</span>
              <span>Lives in <strong>Seoul, South Korea</strong></span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-gray-500">📍</span>
              <span>From <strong>San Francisco, California</strong></span>
            </div>
          </div>
        </div>

        <!-- Photos Preview Card -->
        <div class="bg-white p-4 rounded-lg shadow-sm border border-[#e4e6eb]">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-bold text-lg text-gray-900">Photos</h3>
            <span class="text-sm text-blue-600 hover:underline cursor-pointer">See all photos</span>
          </div>
          <div class="grid grid-cols-3 gap-1 rounded-lg overflow-hidden">
            <div v-for="i in 6" :key="i" class="h-24 bg-gradient-to-br from-blue-300 to-indigo-400"></div>
          </div>
        </div>
      </aside>

      <!-- Right Timeline Stream -->
      <section class="md:col-span-7 space-y-4">
        <!-- Post Composer -->
        <FeedComposer />

        <!-- User Posts -->
        <PostCard
          author-name="Alex Mercer"
          time-ago="3h ago"
          content="Architecting the next generation distributed social network. Built with Go microservices, Nuxt 3, and a high-performance event-driven backbone! 🌐✨"
          :reaction-count="142"
          :comment-count="28"
          :share-count="6"
        />

        <PostCard
          author-name="Alex Mercer"
          time-ago="Yesterday"
          content="Milestone unlocked: All unit and integration tests passing across Identity, Social Graph, and API Gateway modules. 🚀"
          :reaction-count="254"
          :comment-count="39"
          :share-count="14"
        />
      </section>
    </main>

    <!-- Edit Bio Modal -->
    <div v-if="isEditingBio" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full p-6 space-y-4">
        <div class="flex items-center justify-between border-b pb-3">
          <h3 class="font-bold text-lg text-gray-900">Edit Bio</h3>
          <button @click="isEditingBio = false" class="text-gray-500 hover:text-gray-700 text-xl font-bold">×</button>
        </div>
        <textarea
          v-model="tempBio"
          rows="3"
          class="w-full border border-gray-300 rounded-lg p-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Describe who you are..."
        ></textarea>
        <div class="flex justify-end space-x-2 pt-2">
          <button
            @click="isEditingBio = false"
            class="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded-lg text-sm font-semibold transition"
          >
            Cancel
          </button>
          <button
            @click="saveBio"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-semibold transition"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isEditingBio = ref(false)
const bioText = ref('Distributed Systems Architect & Social Tech Enthusiast. Building the future of open communication.')
const tempBio = ref(bioText.value)

const saveBio = () => {
  bioText.value = tempBio.value
  isEditingBio.value = false
}
</script>
