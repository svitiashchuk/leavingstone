/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["light", "dark", "pastel", "cupcake"],
  },
  plugins: [
    require('@tailwindcss/typography'),
    require("daisyui"),
  ],
}
