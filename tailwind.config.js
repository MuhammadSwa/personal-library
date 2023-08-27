/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/html/**/*.tmpl.html", "./web/html/base.tmpl.html"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
};
