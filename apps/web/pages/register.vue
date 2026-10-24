<template>
  <div class="min-h-screen bg-[#f0f2f5] flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white p-6 sm:p-8 rounded-xl shadow-lg border border-gray-200 space-y-4">
      <div>
        <h2 class="text-3xl font-extrabold text-gray-900">Sign Up</h2>
        <p class="text-gray-500 text-sm">It's quick and easy.</p>
      </div>

      <hr class="border-gray-200" />

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <input
            v-model="firstName"
            type="text"
            required
            placeholder="First name"
            class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
          <input
            v-model="lastName"
            type="text"
            required
            placeholder="Last name"
            class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-gray-600 mb-1">Date of birth (Must be 14+)</label>
          <input
            v-model="dateOfBirth"
            type="date"
            required
            class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-gray-600 mb-1">Gender</label>
          <div class="grid grid-cols-3 gap-2 text-sm">
            <label class="flex items-center justify-between border border-gray-300 p-2 rounded-lg cursor-pointer">
              <span>Female</span>
              <input type="radio" value="FEMALE" v-model="gender" />
            </label>
            <label class="flex items-center justify-between border border-gray-300 p-2 rounded-lg cursor-pointer">
              <span>Male</span>
              <input type="radio" value="MALE" v-model="gender" />
            </label>
            <label class="flex items-center justify-between border border-gray-300 p-2 rounded-lg cursor-pointer">
              <span>Custom</span>
              <input type="radio" value="CUSTOM" v-model="gender" />
            </label>
          </div>
        </div>

        <div>
          <input
            v-model="email"
            type="email"
            required
            placeholder="Email address"
            class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div>
          <input
            v-model="password"
            type="password"
            required
            placeholder="New password (8+ chars, upper, lower, digit, symbol)"
            class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
          />
        </div>

        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 text-red-600 rounded-lg text-sm">
          {{ errorMessage }}
        </div>

        <p class="text-xs text-gray-500 leading-tight">
          By clicking Sign Up, you agree to our Terms, Privacy Policy and Cookies Policy.
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2.5 bg-green-600 hover:bg-green-700 text-white font-bold text-lg rounded-lg transition duration-150 disabled:opacity-50"
        >
          {{ loading ? 'Creating account...' : 'Sign Up' }}
        </button>
      </form>

      <div class="text-center pt-2">
        <NuxtLink to="/login" class="text-brand-600 hover:underline text-sm font-medium">
          Already have an account?
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const firstName = ref('')
const lastName = ref('')
const dateOfBirth = ref('2000-01-01')
const gender = ref('CUSTOM')
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const handleRegister = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const res = await fetch('/api/v1/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        first_name: firstName.value,
        last_name: lastName.value,
        date_of_birth: dateOfBirth.value,
        gender: gender.value,
        email: email.value,
        password: password.value
      })
    })

    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.error || 'Failed to sign up')
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

