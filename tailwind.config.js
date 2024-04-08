/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        "midnight-blue": "var(--midnight-blue)",
        "midnight-puprle": "var(--midnight-puprle)",
        "clay-purple": "var(--clay-purple)",
        "dragon-purple": "var(--dragon-purple)",
        "btn-highligh": "var(--btn-highligh)",
        "break-line": "var(--break-line)",
        "gengar-purple": "var(--gengar-purple)",
        
      },
      borderWidth: {
        "12": "12px"
      },
    },
  },
  plugins: [],
}

