<template>
  <div class="min-h-screen flex items-center justify-center bg-[#f0f2f5] px-4 py-12">
    <div class="max-w-4xl w-full grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
      <!-- Left Hero Pitch -->
      <div class="text-center md:text-left">
        <h1 class="text-blue-600 font-extrabold text-5xl md:text-6xl tracking-tight mb-3">
          nova
        </h1>
        <p class="text-gray-700 text-xl md:text-2xl font-normal leading-relaxed">
          Nova helps you connect and share with the people in your life.
        </p>
      </div>

      <!-- Right Auth Card -->
      <div class="bg-white p-6 rounded-lg shadow-xl border border-[#e4e6eb] max-w-md mx-auto w-full">
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <input
              v-model="email"
              type="text"
              placeholder="Email address or phone number"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-base"
              required
            />
          </div>
          <div>
            <input
              v-model="password"
              type="password"
              placeholder="Password"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-base"
              required
            />
          </div>
          <button
            type="submit"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 rounded-md text-lg transition-colors shadow-sm"
          >
            Log In
          </button>
        </form>

        <div class="text-center my-3">
          <a href="#" class="text-blue-600 hover:underline text-sm font-medium">Forgotten password?</a>
        </div>

        <hr class="border-gray-300 my-4" />

        <div class="text-center">
          <NuxtLink
            to="/feed"
            class="inline-block bg-green-500 hover:bg-green-600 text-white font-bold px-6 py-3 rounded-md text-base transition-colors"
          >
            Create New Account / Explore Feed
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')

const handleLogin = () => {
  // Direct to feed for demo navigation
  router.push('/feed')
}
</script>
