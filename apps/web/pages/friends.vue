<template>
  <div class="min-h-screen bg-[#f0f2f5]">
    <Navbar />

    <main class="max-w-7xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-4 gap-6">
      <!-- Left Friends Navigation Sidebar -->
      <aside class="col-span-1 bg-white p-4 rounded-lg shadow-sm border border-[#e4e6eb] space-y-2">
        <h2 class="text-xl font-bold text-gray-900 mb-4">Friends</h2>
        <button
          @click="activeTab = 'home'"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-colors"
          :class="activeTab === 'home' ? 'bg-blue-50 text-blue-600' : 'hover:bg-gray-100 text-gray-700'"
        >
          <span class="text-lg">👥</span>
          <span>Home</span>
        </button>
        <button
          @click="activeTab = 'requests'"
          class="w-full flex items-center justify-between px-3 py-2.5 rounded-lg text-sm font-semibold transition-colors"
          :class="activeTab === 'requests' ? 'bg-blue-50 text-blue-600' : 'hover:bg-gray-100 text-gray-700'"
        >
          <div class="flex items-center space-x-3">
            <span class="text-lg">📩</span>
            <span>Friend Requests</span>
          </div>
          <span v-if="friendRequests.length > 0" class="bg-red-500 text-white text-xs px-2 py-0.5 rounded-full font-bold">
            {{ friendRequests.length }}
          </span>
        </button>
        <button
          @click="activeTab = 'suggestions'"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-colors"
          :class="activeTab === 'suggestions' ? 'bg-blue-50 text-blue-600' : 'hover:bg-gray-100 text-gray-700'"
        >
          <span class="text-lg">✨</span>
          <span>Suggestions</span>
        </button>
        <button
          @click="activeTab = 'all'"
          class="w-full flex items-center space-x-3 px-3 py-2.5 rounded-lg text-sm font-semibold transition-colors"
          :class="activeTab === 'all' ? 'bg-blue-50 text-blue-600' : 'hover:bg-gray-100 text-gray-700'"
        >
          <span class="text-lg">📋</span>
          <span>All Friends</span>
        </button>
      </aside>

      <!-- Right Main Content Area -->
      <section class="col-span-1 md:col-span-3 space-y-6">
        <!-- Section: Friend Requests -->
        <div v-if="activeTab === 'home' || activeTab === 'requests'" class="bg-white p-6 rounded-lg shadow-sm border border-[#e4e6eb]">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold text-gray-900">Friend Requests</h3>
            <span class="text-sm text-blue-600 hover:underline cursor-pointer">See all</span>
          </div>

          <div v-if="friendRequests.length === 0" class="py-8 text-center text-gray-500 text-sm">
            No pending friend requests.
          </div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="req in friendRequests"
              :key="req.id"
              class="border border-[#e4e6eb] rounded-lg overflow-hidden shadow-sm flex flex-col justify-between"
            >
              <div class="h-40 bg-gray-200 flex items-center justify-center font-bold text-3xl text-gray-400">
                {{ req.name.charAt(0) }}
              </div>
              <div class="p-3 space-y-2">
                <h4 class="font-semibold text-sm text-gray-900">{{ req.name }}</h4>
                <p class="text-xs text-gray-500">{{ req.mutualCount }} mutual friends</p>
                <div class="space-y-1 pt-1">
                  <button
                    @click="confirmRequest(req.id)"
                    class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-1.5 rounded-md text-xs transition shadow-sm"
                  >
                    Confirm
                  </button>
                  <button
                    @click="deleteRequest(req.id)"
                    class="w-full bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-1.5 rounded-md text-xs transition"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Section: People You May Know (Suggestions) -->
        <div v-if="activeTab === 'home' || activeTab === 'suggestions'" class="bg-white p-6 rounded-lg shadow-sm border border-[#e4e6eb]">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold text-gray-900">People You May Know</h3>
            <span class="text-sm text-blue-600 hover:underline cursor-pointer">See all</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="sug in suggestions"
              :key="sug.id"
              class="border border-[#e4e6eb] rounded-lg overflow-hidden shadow-sm flex flex-col justify-between"
            >
              <div class="h-40 bg-gradient-to-tr from-indigo-100 to-blue-200 flex items-center justify-center font-bold text-3xl text-blue-600">
                {{ sug.name.charAt(0) }}
              </div>
              <div class="p-3 space-y-2">
                <h4 class="font-semibold text-sm text-gray-900">{{ sug.name }}</h4>
                <p class="text-xs text-gray-500">{{ sug.mutualCount }} mutual friends</p>
                <div class="space-y-1 pt-1">
                  <button
                    v-if="!sug.sent"
                    @click="sug.sent = true"
                    class="w-full bg-blue-100 hover:bg-blue-200 text-blue-700 font-semibold py-1.5 rounded-md text-xs transition"
                  >
                    Add friend
                  </button>
                  <button
                    v-else
                    disabled
                    class="w-full bg-gray-100 text-gray-500 font-semibold py-1.5 rounded-md text-xs cursor-default"
                  >
                    Request Sent
                  </button>
                  <button class="w-full bg-gray-200 hover:bg-gray-300 text-gray-800 font-semibold py-1.5 rounded-md text-xs transition">
                    Remove
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref('home')

const friendRequests = ref([
  { id: 'req-1', name: 'David Kim', mutualCount: 12 },
  { id: 'req-2', name: 'Elena Rostova', mutualCount: 4 }
])

const suggestions = ref([
  { id: 'sug-1', name: 'Marcus Vance', mutualCount: 19, sent: false },
  { id: 'sug-2', name: 'Chloe Zhao', mutualCount: 8, sent: false },
  { id: 'sug-3', name: 'Liam O\'Connor', mutualCount: 2, sent: false }
])

const confirmRequest = (id: string) => {
  friendRequests.value = friendRequests.value.filter(r => r.id !== id)
}

const deleteRequest = (id: string) => {
  friendRequests.value = friendRequests.value.filter(r => r.id !== id)
}
</script>
