/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/app/**/*.{html,js,templ}"],
  theme: {
    extend: {},
  },
  darkMode: 'class',
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["garden", "dracula"],
  },
}

