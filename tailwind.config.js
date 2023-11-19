/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        "midnight-blue": "#3a3a67",
        "midnight-puprle": "#190835",
        "clay-purple": "#85648A",
        "dragon-purple": "#2f0136",
        "gengar-purple": "#7979a4"
      },
      borderWidth: {
        "12": "12px"
      }
    },
  },
  plugins: [],
}

