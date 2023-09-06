import htmx from "htmx.org";
import Alpine from "alpinejs";
import "./toTheUpBtn.js";

document.Alpine = Alpine;
Alpine.start();

function toggleMenu(flag) {
  let value = document.getElementById("menu");
  if (flag) {
    value.classList.remove("hidden");
  } else {
    value.classList.add("hidden");
  }
}
