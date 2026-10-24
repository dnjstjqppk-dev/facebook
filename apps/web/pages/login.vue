<template>
  <div class="min-h-screen bg-[#f0f2f5] flex items-center justify-center p-4">
    <div class="max-w-4xl w-full grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
      <!-- Brand Lead -->
      <div class="space-y-4 text-center md:text-left">
        <h1 class="text-5xl font-extrabold text-brand-600 tracking-tight">facebook</h1>
        <p class="text-2xl font-normal text-gray-700 leading-snug">
          Connect with friends and the world around you on the Global Social Network.
        </p>
      </div>

      <!-- Login Card -->
      <div class="bg-white p-6 sm:p-8 rounded-xl shadow-lg border border-gray-200 space-y-4">
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <input
              v-model="email"
              type="text"
              required
              placeholder="Email or phone number"
              class="w-full px-4 py-3 bg-gray-50 border border-gray-300 rounded-lg text-base focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </div>
          <div>
            <input
              v-model="password"
              type="password"
              required
              placeholder="Password"
              class="w-full px-4 py-3 bg-gray-50 border border-gray-300 rounded-lg text-base focus:outline-none focus:ring-2 focus:ring-brand-500"
            />
          </div>

          <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 text-red-600 rounded-lg text-sm">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-brand-600 hover:bg-brand-700 text-white font-bold text-lg rounded-lg transition duration-150 disabled:opacity-50"
          >
            {{ loading ? 'Logging in...' : 'Log In' }}
          </button>
        </form>

        <div class="text-center">
          <NuxtLink to="/recover" class="text-brand-600 hover:underline text-sm font-medium">
            Forgot password?
          </NuxtLink>
        </div>

        <hr class="border-gray-200" />

        <div class="text-center pt-2">
          <NuxtLink
            to="/register"
            class="inline-block px-6 py-3 bg-green-600 hover:bg-green-700 text-white font-bold text-base rounded-lg transition duration-150"
          >
            Create new account
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
const errorMessage = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, password: password.value })
    })

    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.error || 'Failed to log in')
    }

    const data = await res.json()
    localStorage.setItem('auth_tokens', JSON.stringify(data.tokens))
    router.push('/')
  } catch (err: any) {
    errorMessage.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

