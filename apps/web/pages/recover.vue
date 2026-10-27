<template>
  <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md text-center">
      <h1 class="text-4xl font-extrabold text-blue-600 tracking-tight">facebook</h1>
      <h2 class="mt-4 text-2xl font-bold text-gray-900">Find Your Account</h2>
      <p class="mt-2 text-sm text-gray-600">Please enter your email or mobile number to search for your account.</p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-6 shadow-md rounded-lg sm:px-10">
        <form @submit.prevent="handleRecover" class="space-y-4">
          <div>
            <label for="identifier" class="sr-only">Email or Mobile Number</label>
            <input
              id="identifier"
              v-model="identifier"
              type="text"
              required
              placeholder="Email address or mobile number"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>

          <div v-if="successMessage" class="p-3 bg-green-50 text-green-700 text-sm rounded-lg border border-green-200">
            {{ successMessage }}
          </div>
          <div v-if="errorMessage" class="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-200">
            {{ errorMessage }}
          </div>

          <div class="flex items-center justify-end space-x-3 pt-2">
            <NuxtLink
              to="/login"
              class="px-4 py-2 border border-gray-300 rounded-lg text-sm font-semibold text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </NuxtLink>
            <button
              type="submit"
              :disabled="loading"
              class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold rounded-lg transition disabled:opacity-50"
            >
              {{ loading ? 'Searching...' : 'Search' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const identifier = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleRecover = async () => {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const res = await fetch('/api/v1/auth/recovery', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ identifier: identifier.value })
    })

    if (!res.ok) {
      successMessage.value = 'If an account matches that info, we sent instructions to reset your password.'
    } else {
      successMessage.value = 'Recovery instructions have been sent to your registered contact.'
    }
  } catch (err) {
    successMessage.value = 'If an account matches that info, we sent instructions to reset your password.'
  } finally {
    loading.value = false
  }
}
</script>

