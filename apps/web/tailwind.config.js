/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ['class', '[data-theme="dark"]'],
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue'
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eef2ff',
          100: '#e0e7ff',
          500: '#3b82f6',
          600: '#1877f2', // Facebook Trust Blue Parity
          700: '#1d4ed8'
        },
        surface: {
          light: '#f0f2f5',
          card: '#ffffff',
          dark: '#18191a',
          darkCard: '#242526'
        }
      }
    }
  },
  plugins: []
}
