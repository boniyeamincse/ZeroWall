/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'zw': {
          'dark': '#0a1929',
          'darker': '#060d17',
          'light': '#0f2942',
          'accent': '#00d4ff',
          'accent-hover': '#00b8e6',
          'success': '#00c853',
          'warning': '#ffab00',
          'danger': '#ff1744',
        }
      }
    },
  },
  plugins: [],
}
