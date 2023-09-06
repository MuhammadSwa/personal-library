/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/html/**/*.tmpl.html",
    "./web/html/base.tmpl.html",
    "./node_modules/tw-elements/dist/js/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui"), require("tw-elements/dist/plugin.cjs")],
  darkMode: "class",
};
